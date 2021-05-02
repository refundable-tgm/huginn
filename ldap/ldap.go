package ldap

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"github.com/google/uuid"
	"github.com/refundable-tgm/huginn/db"
	"github.com/refundable-tgm/huginn/untis"
	"strings"
)

// URL is the ip address of the tgm ldap server to connect to
const URL = "10.2.24.151"

// Port of the tgm ldap server. In this case it is the default port
const Port = 389

// AuthenticateUserCredentials authenicates a user given by username and password through the tgm ldap server.
// Furthermore if it is the first login of a user it will create a new Teacher instance and save it to the local database.
// It will return true if the credentials are valid and able to produce a successful login operation on the ldap server
// Otherwise if any connection error occurs or the credentials aren't valid this method will return false
func AuthenticateUserCredentials(username, password string) bool {
	cred := username + "@tgm.ac.at"
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", URL, Port))
	if err != nil {
		return false
	}
	err = l.Bind(cred, password)
	if err != nil {
		return false
	}
	mongo := db.MongoDatabaseConnector{}
	if !mongo.Connect() {
		return false
	}
	defer mongo.Close()
	longname, err := GetLongName(username, password, username)
	if err != nil {
		return false
	}
	if !mongo.DoesTeacherExistByShort(username) {
		client := untis.CreateClient(username, password)
		err = client.Authenticate()
		if err != nil {
			_ = client.Close()
			return false
		}
		id, err := client.ResolveTeacherID(longname)
		if err != nil {
			return false
		}
		untisAb, err := client.ResolveTeachers([]int{id})
		if err != nil {
			return false
		}
		if !mongo.CreateTeacher(db.Teacher{
			UUID:           uuid.NewString(),
			Short:          username,
			Longname:       longname,
			SuperUser:      false,
			AV:             false,
			Administration: false,
			PEK:            false,
			Untis:          untisAb[0],
		}) {
			_ = client.Close()
			return false
		}
		err = client.Close()
		if err != nil {
			return false
		}
	}
	return true
}

// GetLongName will find out the full name (name + surname) of a teacher identified by key through their saved file on the active directory
// ldap server. If the search operation was successful the full name is returned. Otherwise any error occurred will be
// returned.
func GetLongName(username, password, key string) (string, error) {
	cred := username + "@tgm.ac.at"
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", URL, Port))
	if err != nil {
		return "", err
	}
	err = l.Bind(cred, password)
	if err != nil {
		return "", err
	}
	search := ldap.NewSearchRequest("DC=tgm,DC=ac,DC=at",
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(&(mailNickname=%s))", key),
		[]string{"dn"},
		nil,
	)
	res, err := l.Search(search)
	if err != nil {
		return "", err
	}
	if len(res.Entries) != 1 {
		return "", fmt.Errorf("user does not exist or too many entries returned")
	}
	userdn := res.Entries[0]
	cn := strings.Split(userdn.DN, ",")[0]
	name := strings.Split(cn, "=")[1]
	return name, nil
}
