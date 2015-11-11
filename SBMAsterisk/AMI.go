package SBMAsterisk

import (
	//"database/sql"
	"fmt"
	"log"
	"net"
	//"time"
	//"regexp"
	"strings"

	// PostgreSQL
	//_ "github.com/lib/pq"

	// Asterisk ARI
	"code.google.com/p/gami"

	//"github.com/BestianRU/SABModules/SBMConnect"
	"github.com/BestianRU/SABModules/SBMSystem"
	"github.com/BestianRU/SABModules/SBMText"
)

const (
	astLRNumber       = 0
	astLUChannel      = 0
	astLUUserProf     = 1
	astLUBridgeProf   = 2
	astLUNumber       = 3
	astReadErrorRetry = 2
)

type AAMI struct {
	c         net.Conn
	G         *gami.Asterisk
	LRNumber  int
	LUChannel int
	LUUser    int
	LUBridge  int
	LUNumber  int
	Retry     int
	//jsonConf  SBMSystem.ReadJSONConfig
}

func (_s *AAMI) Init(conf SBMSystem.ReadJSONConfig) int {
	var err error

	//_s.jsonConf = conf

	_s.c, err = net.Dial("tcp", fmt.Sprintf("%s:%d", conf.Conf.AST_ARI_Host, conf.Conf.AST_ARI_Port))
	if err != nil {
		log.Printf("Asterisk AMI::Dial() error: %v", err)
		return -1
	}

	_s.G = gami.NewAsterisk(&_s.c, nil)

	err = _s.G.Login(conf.Conf.AST_ARI_User, conf.Conf.AST_ARI_Pass)
	if err != nil {
		log.Printf("Asterisk AMI::Login() error: %v\n", err)
		return -1
	}

	_s.LRNumber = astLRNumber

	_s.LUChannel = astLUChannel
	_s.LUUser = astLUUserProf
	_s.LUBridge = astLUBridgeProf
	_s.LUNumber = astLUNumber

	_s.Retry = astReadErrorRetry

	return 0
}

func (_s *AAMI) Close() {
	_s.G.Logoff()
}

func (_s *AAMI) Query(query ...interface{}) [][]string {
	var (
		err     error
		y       []string
		z       [][]string
		marker  = int(-1)
		queryGo = string("")
	)

	for _, x := range query {
		queryGo = fmt.Sprintf("%s %v", queryGo, x)
	}
	queryGo = SBMText.RemoveDoubleSpace(queryGo)
	//fmt.Printf("---%s---\n\n", queryGo)

	ast_get := make(chan gami.Message, 10000)

	ast_cb := func(m gami.Message) {
		ast_get <- m
	}

	err = _s.G.Command(queryGo, &ast_cb)
	if err != nil {
		log.Printf("Asterisk ARI::Command() error: %v\n", err)
		return nil
	}

	for x1, x2 := range <-ast_get {
		if x1 == "CmdData" {
			x := strings.Split(x2, "\n")
			marker = -1
			for i := 0; i < len(x); i++ {
				if strings.Contains(x[i], "--END COMMAND--") {
					marker = -1
				}
				if marker >= 0 {
					y = append(y, SBMText.RemoveDoubleSpace(x[i]))
				}
				if strings.Contains(x[i], "==========") {
					marker = 0
				}
			}
		}
	}

	for i := 0; i < len(y); i++ {
		z = append(z, strings.Split(y[i], " "))
	}

	return z

}

func (_s *AAMI) Kick(conf, number string) int {
	ch, i := _s._getCh(conf, number)
	if i < 0 {
		return -1
	}
	_s.Query("confbridge kick", conf, ch)
	return 0
}

func (_s *AAMI) Mute(conf, number, mode string) (string, int) {
	var (
		m       gami.Message
		res     = int(-1)
		message string
	)

	ch, i := _s._getCh(conf, number)
	if i < 0 {
		return "", -1
	}
	cch := make(chan gami.Message, 1000)
	cb := func(m gami.Message) {
		if m["EventList"] == "Complete" || m["Response"] == "Error" {
			_s.G.DelCallback(m)
			close(cch)
		} else {
			cch <- m
		}
	}

	if SBMText.Low(mode) == "mute" {
		m = gami.Message{"Action": "ConfbridgeMute", "Conference": conf, "Channel": ch}
	} else {
		m = gami.Message{"Action": "ConfbridgeUnmute", "Conference": conf, "Channel": ch}
	}
	_s.G.HoldCallbackAction(m, &cb)
	for x1, x2 := range <-cch {
		switch {
		case x1 == "Response":
			if SBMText.Low(x2) == "success" {
				res = 1
			}
		case x1 == "Message":
			message = x2
		}
	}

	return message, res
}

func (_s *AAMI) List() []map[string]string {
	var (
		z []map[string]string
		i int
	)
	cch := make(chan []map[string]string, 1000)
	cb := func(m gami.Message) {
		if m["EventList"] == "Complete" || m["Response"] == "Error" {
			//fmt.Printf("DEBUG: %v\n\n", z)
			cch <- z
			_s.G.DelCallback(m)
		} else {
			x := make(map[string]string, len(m))
			for x1, x2 := range m {
				x[x1] = x2
			}
			if x["Event"] == "ConfbridgeListRooms" && len(x["Conference"]) > 0 {
				z = append(z, x)
			}
		}
	}

	zz := _s.Query("confbridge list")
	//fmt.Println("\n\nXSimpleList: ", zz)
	for i = 0; i < _s.Retry; i++ {
		m := gami.Message{"Action": "ConfbridgeListRooms"}
		_s.G.HoldCallbackAction(m, &cb)

		for x := range <-cch {
			if 1 != 1 {
				fmt.Println(x)
			}
		}

		if len(z) == len(zz) {
			break
		}
		log.Printf("Asterisk AMI::ConfList read error! Attempt %d of %d\n", i+1, _s.Retry)
		z = nil
		z = make([]map[string]string, 0)
	}

	if i >= _s.Retry {
		z = nil
		z = make([]map[string]string, 0)
		for _, x := range zz {
			//fmt.Println(x, x[_s.LRNumber])
			xx := make(map[string]string, 1)
			xx["Conference"] = x[_s.LRNumber]
			z = append(z, xx)
		}
		//fmt.Println("Z:", z)
	}

	return z
}

func (_s *AAMI) ListUsers(conf string) []map[string]string {
	var (
		z []map[string]string
		i int
	)
	cch := make(chan []map[string]string, 1000)
	cb := func(m gami.Message) {
		if m["EventList"] == "Complete" || m["Response"] == "Error" {
			//fmt.Printf("DEBUG1: %v\n\n", z)
			cch <- z
			_s.G.DelCallback(m)
		} else {
			x := make(map[string]string, len(m))
			for x1, x2 := range m {
				x[x1] = x2
			}
			//fmt.Printf("DEBUG0: %v\n\n", x)
			if x["Event"] == "ConfbridgeList" && x["Conference"] == conf {
				//fmt.Printf("DEBUG2: %v\n\n", z)
				z = append(z, x)
			}
		}
	}

	zz := _s.Query("confbridge list " + conf)
	//fmt.Println("\n\nSimpleList: ", zz)
	for i = 0; i < _s.Retry; i++ {

		m := gami.Message{"Action": "ConfbridgeList", "Conference": conf}
		_s.G.HoldCallbackAction(m, &cb)

		for _, x := range <-cch {
			//fmt.Printf("DEBUG3: %v\n", x)
			if 1 != 1 {
				fmt.Println(x)
			}
		}
		if len(z) == len(zz) && zz != nil {
			break
		}
		log.Printf("Asterisk AMI::ConfUserList read error! Attempt %d of %d\n", i+1, _s.Retry)
		z = nil
		z = make([]map[string]string, 0)
	}

	if i >= _s.Retry {
		/*var pg SBMConnect.PgSQL
		pg.Init(_s.jsonConf, "")
		defer pg.Close()*/
		z = nil
		z = make([]map[string]string, 0)
		for _, x := range zz {
			//fmt.Println(x)
			xx := make(map[string]string, 4)
			xx["Conference"] = conf
			xx["CallerIDNum"] = x[3]
			xx["Channel"] = x[0]
			/*rows, err := pg.D.Query(fmt.Sprintf("select x.cid_name from ldapx_persons x, ldapx_phones y, (select a.phone, count(a.phone) as phone_count from ldapx_phones as a, ldapx_persons as b where a.pers_id=b.uid and b.contract=0 and a.pass=2 and b.lang=1 group by a.phone order by a.phone) as subq where x.uid=y.pers_id and y.pass=2 and x.lang=1 and subq.phone=y.phone and subq.phone_count<2 and y.phone like '%s%%' and x.contract=0 group by x.cid_name, y.phone order by y.phone;", xx["CallerIDNum"]))
			if err != nil {
				log.Printf("PgSQL::Query() error: %v\n", err)
				xx["CallerIDName"] = "-"
			} else {
				rows.Next()
				rows.Scan(*xx["CallerIDName"])
			}*/
			z = append(z, xx)
		}
		//fmt.Println("Z:", z)
	}

	return z
}

func (_s *AAMI) _getCh(conf, number string) (string, int) {
	x := _s.Query("confbridge list", conf)
	if x == nil {
		return "", -1
	}
	for i := 0; i < len(x); i++ {
		if x[i][_s.LUNumber] == _s._pnModify(number) || x[i][_s.LUNumber] == number {
			return x[i][_s.LUChannel], 0
		}
	}
	return "", -1
}

func (_s *AAMI) _pnModify(x string) string {
	return x
}
