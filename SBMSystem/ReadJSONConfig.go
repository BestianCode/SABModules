package SBMSystem

import (
	"encoding/json"
	"flag"
	"fmt"
	//"log"
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
		LogLevel            int
	}
}

func (_s *ReadJSONConfig) Init() {
	_s.Conf.LOG_File = "./AmnesiacDefault.log" // Default log file
	_s.Config_file = "./AmnesiacDefault.json"  // Default configuration file
	_s.Daemon_mode = "NO"                      // Default start in foreground

	_s._parseCommandLine()
	_s._readConfigFile()

	fmt.Printf("Configuration file: %s\n", _s.Config_file)
	fmt.Printf("          Log file: %s\n", _s.Conf.LOG_File)
	fmt.Printf("          PID file: %s\n", _s.Conf.PID_File)
	fmt.Printf("       Daemon mode: %s\n", _s.Daemon_mode)
	fmt.Printf("\n")
	fmt.Printf("\n")
}

func (_s *ReadJSONConfig) _parseCommandLine() {
	cp := flag.String("config", _s.Config_file, "Path to Configuration file")
	dp := flag.String("daemon", _s.Daemon_mode, "Fork as system daemon (YES or NO)")
	flag.Parse()
	_s.Config_file = *cp
	_s.Daemon_mode = *dp

	//fmt.Println(*cp, "\n", *dp, "\n", os.Args, "\n")
}

func (_s *ReadJSONConfig) _readConfigFile() {
	f, err := os.Open(_s.Config_file)
	if err != nil {
		fmt.Printf("Error open Configuration file: %s (%v)\n", _s.Config_file, err)
		os.Exit(1)
	}

	c := json.NewDecoder(f)
	err = c.Decode(&_s.Conf)
	if err != nil {
		fmt.Printf("Error read Configuration file: %s (%v)\n", _s.Config_file, err)
		os.Exit(2)
	}
	f.Close()
}
