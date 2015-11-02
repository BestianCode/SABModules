package SBMSystem

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

type ReadJSONConfig struct {
	Config_file string
	Daemon_mode string
	Conf        struct {
		PG_DSN              string
		MY_DSN              string
		AST_SQLite_DB       string
		AST_CID_Group       string
		AST_Num_Start       string
		AST_ARI_Host        string
		AST_ARI_Port        int
		AST_ARI_User        string
		AST_ARI_Pass        string
		Oracle_SRV          [][]string
		MSSQL_DSN           [][]string
		LDAP_URL            [][]string
		ROOT_OU             string
		ROOT_DN             [][]string
		Sleep_Time          int
		Sleep_cycles        int
		LOG_File            string
		PID_File            string
		TRANS_NAMES         [][]string
		BlackList_OU        []string
		WLB_SessTimeOut     int
		WLB_Listen_IP       string
		WLB_Listen_PORT     int
		WLB_LDAP_ATTR       [][]string
		WLB_SQL_PreFetch    string
		WLB_MailBT          string
		WLB_SQLite_DB       string
		WLB_HTML_Title      string
		AD_LDAP             [][]string
		AD_ScriptDir        string
		AD_LDAP_PARENT      [][]string
		TRANS_POS           [][]string
		SABRealm            string
		WLB_DavDNTreeDepLev int
	}
}

func (_s *ReadJSONConfig) Init() {
	_s.Conf.LOG_File = "./AmnesiacDefault.log" // Default log file
	_s.Config_file = "./AmnesiacDefault.json"  // Default configuration file
	_s.Daemon_mode = "NO"                      // Default start in foreground

	_s._parseCommandLine()
	_s._readConfigFile()

	log.Printf(".\n")
	log.Printf("Configuration file: %s\n", _s.Config_file)
	log.Printf("          Log file: %s\n", _s.Conf.LOG_File)
	log.Printf("          PID file: %s\n", _s.Conf.PID_File)
	log.Printf("       Daemon mode: %s\n", _s.Daemon_mode)
	log.Printf("Go!\n")
	log.Printf(".\n")
}

func (_s *ReadJSONConfig) _parseCommandLine() {
	cp := flag.String("config", _s.Config_file, "Path to Configuration file")
	dp := flag.String("daemon", _s.Daemon_mode, "Fork as system daemon (YES or NO)")
	flag.Parse()
	_s.Config_file = *cp
	_s.Daemon_mode = *dp
}

func (_s *ReadJSONConfig) _readConfigFile() {
	f, err := os.Open(_s.Config_file)
	if err != nil {
		log.Fatalf("Error open Configuration file: %s (%v)\n", _s.Config_file, err)
	}

	c := json.NewDecoder(f)
	err = c.Decode(&_s.Conf)
	if err != nil {
		log.Fatalf("Error read Configuration file: %s (%v)\n", _s.Config_file, err)
	}
	f.Close()
}
