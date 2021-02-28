package db

const TeacherCollection = "Teacher"
const ApplicationCollection = "Collection"

type MongoDatabaseConnector struct {
	name     string
	password string
}

func (m MongoDatabaseConnector) Connect() bool {
	return false
}

func (m MongoDatabaseConnector) Close() (ok bool) {
	return false
}

func (m MongoDatabaseConnector) CreateApplication() (ok bool) {
	return false
}

func (m MongoDatabaseConnector) GetApplication() (application Application) {
	return Application{}
}

func (m MongoDatabaseConnector) UpdateApplication() (ok bool) {
	return false
}

func (m MongoDatabaseConnector) DeleteApplication() (ok bool) {
	return false
}

func (m MongoDatabaseConnector) CreateTeacher() (ok bool) {
	return false
}

func (m MongoDatabaseConnector) GetTeacher() (teacher Teacher) {
	return Teacher{}
}

func (m MongoDatabaseConnector) UpdateTeacher() (ok bool) {
	return false
}

func (m MongoDatabaseConnector) DeleteTeacher() (ok bool) {
	return false
}

func resolveURI() (URI string, database string, ok bool) {
	return
}
