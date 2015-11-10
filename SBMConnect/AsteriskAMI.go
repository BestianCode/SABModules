package SBMConnect

import (
	//"database/sql"
	"fmt"
	"log"
	"net"
	"time"
	//"regexp"
	"strings"

	// MySQL
	//_ "github.com/ziutek/mymysql/godrv"

	// Asterisk ARI
	"code.google.com/p/gami"

	"github.com/BestianRU/SABModules/SBMSystem"
	"github.com/BestianRU/SABModules/SBMText"
)

const (
	astNChannel    = 0
	astNUserProf   = 1
	astNBridgeProf = 2
	astNNumber     = 3
)

type AAMI struct {
	c        net.Conn
	G        *gami.Asterisk
	NChannel int
	NUser    int
	NBridge  int
	NNumber  int
}

func (_s *AAMI) Init(conf SBMSystem.ReadJSONConfig) int {
	var err error

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

	_s.NChannel = astNChannel
	_s.NUser = astNUserProf
	_s.NBridge = astNBridgeProf
	_s.NNumber = astNNumber

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

func (_s *AAMI) Mute(conf, number, mode string) int {
	ch, i := _s._getCh(conf, number)
	if i < 0 {
		return -1
	}
	if SBMText.Low(mode) == "mute" {
		_s.Query("confbridge mute", conf, ch)
	} else {
		_s.Query("confbridge unmute", conf, ch)
	}
	return 0
}

func (_s *AAMI) List() []map[string]string {
	var (
		z []map[string]string
		i int
	)

	cch := make(chan []map[string]string)
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

	for i = 0; i < 10; i++ {
		m := gami.Message{"Action": "ConfbridgeListRooms"}
		_s.G.HoldCallbackAction(m, &cb)

		for x := range <-cch {
			if 1 != 1 {
				fmt.Println("1", x)
			}
		}

		if len(z) == len(_s.Query("confbridge list")) {
			break
		}
		fmt.Println("DEBUG: ConfList read error! Try again!")
		time.Sleep(time.Second)
		z = nil
	}

	if i >= 10 {
		fmt.Println("DEBUG: ConfList read error! Stop... :(", i)
		return nil
	}

	return z
}

func (_s *AAMI) _getCh(conf, number string) (string, int) {
	x := _s.Query("confbridge list", conf)
	if x == nil {
		return "", -1
	}
	for i := 0; i < len(x); i++ {
		if x[i][astNNumber] == _s._pnModify(number) || x[i][astNNumber] == number {
			return x[i][astNChannel], 0
		}
	}
	return "", -1
}

func (_s *AAMI) _pnModify(x string) string {
	return x
}
