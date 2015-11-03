package SBMConnect

import (
	"log"

	// LDAP
	"github.com/go-ldap/ldap"

	"github.com/BestianRU/SABModules/SBMSystem"
)

var LDAPCounter = int(0)

type LDAP struct {
	D *ldap.Conn
}

func (_s *LDAP) Init(conf SBMSystem.ReadJSONConfig) int {
	var (
		attemptCounter = int(0)
		err            error
		//l   *ldap.Conn
	)

	for {
		if attemptCounter > len(conf.Conf.LDAP_URL)*2 {
			log.Printf("LDAP Init SRV ***** Error connect to all LDAP servers...")
			return -1
		}

		if LDAPCounter > len(conf.Conf.LDAP_URL)-1 {
			LDAPCounter = 0
		}

		log.Printf("LDAP Init SRV ***** Trying connect to server %d of %d: %s", LDAPCounter+1, len(conf.Conf.LDAP_URL), conf.Conf.LDAP_URL[LDAPCounter][0])
		_s.D, err = ldap.Dial("tcp", conf.Conf.LDAP_URL[LDAPCounter][0])
		if err != nil {
			LDAPCounter++
			attemptCounter++
			continue
		}

		log.Printf("LDAP Init SRV ***** Success! Connected to server %d of %d: %s", LDAPCounter+1, len(conf.Conf.LDAP_URL), conf.Conf.LDAP_URL[LDAPCounter][0])
		LDAPCounter++
		break
	}

	//_s.D.Debug()

	err = _s.D.Bind(conf.Conf.LDAP_URL[0][1], conf.Conf.LDAP_URL[0][2])
	if err != nil {
		log.Printf("LDAP::Bind() error: %v\n", err)
		return -1
	}

	return 0
}

func (_s *LDAP) Close() {
	_s.D.Close()
}
