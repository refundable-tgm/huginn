package ldap

import (
	"github.com/google/uuid"
	"github.com/refundable-tgm/huginn/db"
)

const URL = ""

func AuthenticateUserCredentials(username, password string) bool {
	//ldap.Dial("tcp", fmt.Sprintf())
	mongo := db.MongoDatabaseConnector{}
	defer mongo.Close()
	if !mongo.Connect() {
		return false
	}
	longname := "" // get from ldap
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