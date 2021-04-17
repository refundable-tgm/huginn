package ldap

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"github.com/google/uuid"
	"github.com/refundable-tgm/huginn/db"
	"strings"
)

const URL = "10.2.24.151"
const Port = 389

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
	defer mongo.Close()
	if !mongo.Connect() {
		return false
	}
	longname, err := GetLongName(username, password)
	if err != nil {
		return false
	}
	if !mongo.DoesTeacherExistsByShort(username) {
		mongo.CreateTeacher(db.Teacher{
			UUID:           uuid.NewString(),
			Short:          username,
			Longname:       longname,
			SuperUser:      false,
			AV:             false,
			Administration: false,
			PEK:            false,
		})
	}
	return true
}

func GetLongName(username, password string) (string, error) {
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
		fmt.Sprintf("(&(mailNickname=%s))", username),
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
