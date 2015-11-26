package SBMConnect

import (
	"log"
	//"strings"

	// LDAP
	"github.com/go-ldap/ldap"

	"github.com/BestianRU/SABModules/SBMSystem"
)

var LDAPCounter = int(0)

type LDAP struct {
	CS int
	D  *ldap.Conn
}

func (_s *LDAP) Init(conf SBMSystem.ReadJSONConfig) int {
	var (
		attemptCounter = int(0)
		err            error
	)

	_s.CS = -1

	for {
		if attemptCounter > len(conf.Conf.LDAP_URL)*2 {
			log.Printf("LDAP Init SRV ***** Error connect to all LDAP servers !!!")
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

	_s.CS = 0
	return 0
}

func (_s *LDAP) InitS(rLog SBMSystem.LogFile, user, password, server string) int {
	var err error

	_s.CS = -1

	log.Printf("LDAP Init SRV ***** Trying connect to server %s with login %s.", server, user)

	_s.D, err = ldap.Dial("tcp", server)
	if err != nil {
		log.Printf("LDAP::Dial() error: %v\n", err)
		return -1
	}

	//L.Debug()

	err = _s.D.Bind(user, password)
	if err != nil {
		log.Printf("LDAP::Bind() error: %v\n", err)
		return -1
	}

	log.Printf("LDAP Init SRV ***** Success! Connected to server %s with login %s.", server, user)

	_s.CS = 0
	return 0
}

func (_s *LDAP) CheckGroupMember(rLog SBMSystem.LogFile, user, group, baseDN string) int {
	const (
		recurs_count = 10
	)

	log.Printf("LDAP CheckGroupMember...")

	userDN := _s._getBaseDN(rLog, user, baseDN)
	groupDN := _s._getBaseDN(rLog, group, baseDN)

	if userDN == "" || groupDN == "" {
		return -1
	}

	if _s._checkGroupMember(rLog, userDN, groupDN, baseDN, 1) == 0 {
		return 0
	} else {
		return _s._checkGroupMember(rLog, userDN, groupDN, baseDN, recurs_count)
	}

	return -1
}

func (_s *LDAP) _checkGroupMember(rLog SBMSystem.LogFile, userDN, groupDN, baseDN string, recurse_count int) int {
	var (
		uattr  = []string{"memberOf"}
		result = int(-1)
	)

	if userDN == "" || groupDN == "" {
		return -1
	}

	if recurse_count <= 0 {
		return -1
	}

	lsearch := ldap.NewSearchRequest(userDN, 0, ldap.NeverDerefAliases, 0, 0, false, "(objectclass=*)", uattr, nil)
	sr, err := _s.D.Search(lsearch)
	if err != nil {
		log.Printf("LDAP::Search() error: %v\n", err)
	}

	if len(sr.Entries) > 0 {
		for _, entry := range sr.Entries {
			for _, attr := range entry.Attributes {
				if attr.Name == "memberOf" {
					for _, x := range attr.Values {
						if groupDN == x {
							return 0
						} else {
							if x != userDN {
								result = _s._checkGroupMember(rLog, x, groupDN, baseDN, recurse_count-1)
								if result == 0 {
									return 0
								}
							}
						}
					}
				}
			}
		}
	}
	return -1
}

func (_s *LDAP) _getBaseDN(rLog SBMSystem.LogFile, search, basedn string) string {
	var uattr = []string{"dn"}

	lsearch := ldap.NewSearchRequest(basedn, 2, ldap.NeverDerefAliases, 0, 0, false, search, uattr, nil)
	sr, err := _s.D.Search(lsearch)
	if err != nil {
		log.Printf("LDAP::Search() error: %v\n", err)
	}

	if len(sr.Entries) > 0 {
		for _, entry := range sr.Entries {
			return entry.DN
		}
	}
	return ""
}

func (_s *LDAP) Close() {
	if _s.CS != -1 {
		_s.D.Close()
	}
}
