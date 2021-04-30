package db

import (
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

// TeacherCollection is the name of the collection in which the Teacher data is stored in
const TeacherCollection = "Teacher"

// ApplicationCollection is the name of the collection in which the Application data is stored in
const ApplicationCollection = "Application"

// SuperUserPath is the path to a file containing the name of the first Teacher to become a super user
const SuperUserPath = "/vol/files/.superuser"

// The MongoDatabaseConnector saves data used for the mongo db connection
type MongoDatabaseConnector struct {
	// the name of the database in the mongo db server
	database string
	// the client used in this connection
	client *mongo.Client
	// the created context of the client
	context context.Context
	// CancelFunc of the context
	closer context.CancelFunc
}

// Connect the MongoDatabaseConnector with the given MongoDB server
// returns whether this operation was successful
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

// Close the Connection to the MongoDB
// returns whether this operation was successful
func (m MongoDatabaseConnector) Close() (ok bool) {
	err := m.client.Disconnect(m.context)
	m.closer()
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// CreateApplication creates a new application in the collection in the database
func (m MongoDatabaseConnector) CreateApplication(application Application) bool {
	application.UUID = uuid.New().String()
	collection := m.client.Database(m.database).Collection(ApplicationCollection)
	insert, err := collection.InsertOne(m.context, application)
	if err != nil {
		log.Println(err)
		return false
	}
	log.Println("Inserted a new application with the UUID: ", application.UUID,
		"; the Title: ", application.Name, "; under the ID: ", insert.InsertedID)
	return true
}

// GetApplication returns a specific application described and identified by its uuid
func (m MongoDatabaseConnector) GetApplication(uuid string) (application Application) {
	collection := m.client.Database(m.database).Collection(ApplicationCollection)
	if err := collection.FindOne(m.context, bson.M{"uuid": uuid}).Decode(&application); err != nil {
		log.Println(err)
		return
	}
	return application
}

// GetAllApplications analyzes all applications contained in the collection and returns them as an array
func (m MongoDatabaseConnector) GetAllApplications() (applications []Application) {
	collection := m.client.Database(m.database).Collection(ApplicationCollection)
	cursor, err := collection.Find(m.context, bson.M{})
	if err != nil {
		log.Println(err)
		return
	}
	if err = cursor.All(m.context, &applications); err != nil {
		log.Println(err)
		return
	}
	return
}

// GetActiveApplications returns all currently active applications stored in the database
func (m MongoDatabaseConnector) GetActiveApplications() (applications []Application) {
	filter := bson.M{
		"$or": []bson.M{
			{"$and": []bson.M{
				{"kind": bson.M{"$in": []int{Training, OtherReason}}},
				{"progress": bson.M{"$in": []int{TRejected, TInProcess, TConfirmed, TRunning, TCostsPending, TCostsInProcess}}},
			}},
			{"$and": []bson.M{
				{"kind": SchoolEvent},
				{"progress": bson.M{"$in": []int{SERejected, SEInSubmission, SEInProcess, SEConfirmed, SERunning, SECostsPending, SECostsInProcess}}},
			}},
		},
	}
	collection := m.client.Database(m.database).Collection(ApplicationCollection)
	cursor, err := collection.Find(m.context, filter)
	if err != nil {
		log.Println(err)
		return
	}
	err = cursor.All(m.context, &applications)
	if err != nil {
		log.Println(err)
		return
	}
	return applications
}

// UpdateApplication updates an application with the matching uuid and updates it with the data in the update struct
// returns true whether one Application was modified, false if an error occurred or no Application was modified
func (m MongoDatabaseConnector) UpdateApplication(uuid string, update Application) bool {
	update.UUID = uuid
	collection := m.client.Database(m.database).Collection(ApplicationCollection)
	result, err := collection.ReplaceOne(m.context, bson.M{"uuid": uuid}, update)
	if err != nil {
		log.Println(err)
		return false
	}
	return result.ModifiedCount == 1
}

// DeleteApplication deletes an application described by the given uuid
// returns true if a document was deleted, false if not or if an error occurred
func (m MongoDatabaseConnector) DeleteApplication(uuid string) bool {
	collection := m.client.Database(m.database).Collection(ApplicationCollection)
	result, err := collection.DeleteOne(m.context, bson.M{"uuid": uuid})
	if err != nil {
		log.Println(err)
		return false
	}
	return result.DeletedCount == 1
}

// DoesApplicationExist searches the database for a Application identified by a given UUID
// and checks whether an Application can be found whilst performing this search.
// It will return true if the Application was found, false if an error occurred or none was found.
func (m MongoDatabaseConnector) DoesApplicationExist(uuid string) bool {
	application := Application{}
	collection := m.client.Database(m.database).Collection(ApplicationCollection)
	if err := collection.FindOne(m.context, bson.M{"uuid": uuid}).Decode(&application); err != nil {
		return false
	}
	return true
}

// CreateTeacher creates a new application in the system
// it will return true if this operation was successful and false if not
func (m MongoDatabaseConnector) CreateTeacher(teacher Teacher) bool {
	collection := m.client.Database(m.database).Collection(TeacherCollection)
	if teacher.Short == getInitUserName() {
		teacher.SuperUser = true
	}
	insert, err := collection.InsertOne(m.context, teacher)
	if err != nil {
		log.Println(err)
		return false
	}
	log.Println("Inserted a new teacher with the UUID: ", teacher.UUID,
		"; the shortname: ", teacher.Short, "; under the ID: ", insert.InsertedID)
	return true
}

// GetTeacherByShort returns a teacher identified by a given short name
func (m MongoDatabaseConnector) GetTeacherByShort(short string) (teacher Teacher) {
	collection := m.client.Database(m.database).Collection(TeacherCollection)
	if err := collection.FindOne(m.context, bson.M{"short": short}).Decode(&teacher); err != nil {
		log.Println(err)
		return
	}
	return teacher
}

// DoesTeacherExistByShort searches the database for a Teacher identified by a shortname
// and checks whether a teacher can be found whilst performing this search.
// It will return true if the teacher was found, false if an error occurred or none was found.
func (m MongoDatabaseConnector) DoesTeacherExistByShort(short string) bool {
	teacher := Teacher{}
	collection := m.client.Database(m.database).Collection(TeacherCollection)
	if err := collection.FindOne(m.context, bson.M{"short": short}).Decode(&teacher); err != nil {
		return false
	}
	return true
}

// GetTeacherByUUID returns a teacher identified by a given UUID
func (m MongoDatabaseConnector) GetTeacherByUUID(uuid string) (teacher Teacher) {
	collection := m.client.Database(m.database).Collection(TeacherCollection)
	if err := collection.FindOne(m.context, bson.M{"uuid": uuid}).Decode(&teacher); err != nil {
		log.Println(err)
		return
	}
	return teacher
}

// DoesTeacherExistByUUID searches the database for a Teacher identified by a uuid
// and checks whether a teacher can be found whilst performing this search.
// It will return true if the teacher was found, false if an error occurred or none was found.
func (m MongoDatabaseConnector) DoesTeacherExistByUUID(uuid string) bool {
	teacher := Teacher{}
	collection := m.client.Database(m.database).Collection(TeacherCollection)
	if err := collection.FindOne(m.context, bson.M{"uuid": uuid}).Decode(&teacher); err != nil {
		return false
	}
	return true
}

// UpdateTeacher updates a teacher with the matching uuid and updates it with the data in the update struct
// returns true whether one Teacher was modified, false if an error occurred or no Teacher was modified
func (m MongoDatabaseConnector) UpdateTeacher(uuid string, update Teacher) bool {
	update.UUID = uuid
	collection := m.client.Database(m.database).Collection(TeacherCollection)
	result, err := collection.ReplaceOne(m.context, bson.M{"uuid": uuid}, update)
	if err != nil {
		log.Println(err)
		return false
	}
	return result.ModifiedCount == 1
}

// DeleteTeacher deletes one teacher described by a given short name
// returns true if a document was deleted, false if none or an error occurred
func (m MongoDatabaseConnector) DeleteTeacher(uuid string) (ok bool) {
	collection := m.client.Database(m.database).Collection(TeacherCollection)
	result, err := collection.DeleteOne(m.context, bson.M{"uuid": uuid})
	if err != nil {
		log.Println(err)
		return false
	}
	return result.DeletedCount == 1
}

// Constructs the URI out of the given information of the docker secrets
// returns the constructed URI, the database name, and whether the operation was successful
// if it was not successful the URI and the database name are empty strings
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
	usernameString := strings.TrimSuffix(string(username), "\n")
	passwordString := strings.TrimSuffix(string(password), "\n")
	return "mongodb://" + usernameString + ":" + passwordString + "@" + "mongo:27017/?authSource=" + database, database, true
}

// getInitUserName returns the in the config file set username to set the first super user
func getInitUserName() string {
	file, err := ioutil.ReadFile(SuperUserPath)
	if err != nil {
		return ""
	}
	name := strings.TrimSuffix(string(file), "\n")
	return name
}
