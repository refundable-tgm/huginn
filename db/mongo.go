package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const TeacherCollection = "Teacher"
const ApplicationCollection = "Collection"

type MongoDatabaseConnector struct {
	database string
	client   *mongo.Client
	context  context.Context
	closer   context.CancelFunc
}

func (m *MongoDatabaseConnector) Connect() bool {
	uri, db, ok := resolveURI()
	if ok {
		client, err := mongo.NewClient(options.Client().ApplyURI(uri))
		if err != nil {
			log.Println(err)
			return false
		}
		ctx, cf := context.WithTimeout(context.Background(), 10*time.Minute)
		err = client.Connect(ctx)
		if err != nil {
			log.Println(err)
			cf()
			return false
		}
		m.client = client
		m.database = db
		m.context = ctx
		m.closer = cf
		return true
	}
	return false
}

func (m MongoDatabaseConnector) Close() (ok bool) {
	err := m.client.Disconnect(m.context)
	m.closer()
	if err != nil {
		log.Println(err)
		return false
	}
	return true
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
	database = os.Getenv("MONGO_DATABASE")
	usernameFilePath := os.Getenv("MONGO_USERNAME_FILE")
	passwordFilePath := os.Getenv("MONGO_PASSWORD_FILE")
	username, err := ioutil.ReadFile(usernameFilePath)
	if err != nil {
		log.Println(err)
		return "", "", false
	}
	password, err := ioutil.ReadFile(passwordFilePath)
	if err != nil {
		log.Println(err)
		return "", "", false
	}
	return "mongodb://" + string(username) + ":" + string(password) + "@" + "mongo:27017", database, true
}
