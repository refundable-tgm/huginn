package rest

import (
	"encoding/base64"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	mongo "github.com/refundable-tgm/huginn/db"
	"github.com/refundable-tgm/huginn/files"
	"github.com/refundable-tgm/huginn/ldap"
	"github.com/refundable-tgm/huginn/untis"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// AuthWall drops every token which doesnt provide a valid token
func AuthWall() gin.HandlerFunc {
	return func(con *gin.Context) {
		ok, err := TokenValid(con.Request)
		if !ok && err != nil {
			con.JSON(http.StatusUnauthorized, Error{err.Error()})
			con.Abort()
			return
		}
		con.Next()
	}
}

// Login represents the login endpoint
// @Summary Login a user
// @Description Login a user using username and password
// @ID login
// @Accept json
// @Produce json
// @Param user body User true "Account Information"
// @Success 200 {object} TokenPair
// @Failure 401 {object} Error
// @Failure 422 {object} Error
// @Router /login [post]
func Login(con *gin.Context) {
	u := User{}
	if err := con.ShouldBindJSON(&u); err != nil {
		con.JSON(http.StatusUnprocessableEntity, Error{"invalid request structure provided"})
		return
	}
	if !ldap.AuthenticateUserCredentials(u.Username, u.Password) {
		con.JSON(http.StatusUnauthorized, Error{"this credentials do not resolve into an authorized login"})
		return
	}
	token, err := CreateToken(u.Username)
	if err != nil {
		con.JSON(http.StatusUnprocessableEntity, Error{err.Error()})
		return
	}
	SaveToken(u.Username, token)
	untis.CreateClient(u.Username, u.Password)
	out := TokenPair{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
	con.JSON(http.StatusOK, out)
}

// Logout represents the logout endpoint
// @Summary Logs out a user
// @Description Destroys the session of a user
// @ID logout
// @Accept json
// @Produce json
// @Param Authorization header string true "Access Token" default(Bearer <Add access token here>)
// @Success 200 {object} Information
// @Failure 401 {object} Error
// @Router /logout [post]
func Logout(con *gin.Context) {
	auth, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, Error{"you are not logged in"})
		return
	}
	DeleteToken(auth.AccessUUID)
	untis.GetClient(auth.Username).DeleteClient()
	con.JSON(http.StatusOK, Information{"logged out"})
}

// Refresh represents the refresh endpoint
// @Summary Refreshes the token pair of a session
// @Description Creates a new token pair when a valid refresh token is provided
// @ID refresh
// @Accept json
// @Produce json
// @Param token body RefreshToken true "Refresh Token"
// @Success 201 {object} TokenPair
// @Failure 401 {object} Error
// @Failure 403 {object} Error
// @Failure 422 {object} Error
// @Router /refresh [post]
func Refresh(con *gin.Context) {
	body := RefreshToken{}
	if err := con.ShouldBindJSON(&body); err != nil {
		con.JSON(http.StatusUnprocessableEntity, Error{err.Error()})
		return
	}
	refresh := body.Token
	token, err := jwt.Parse(refresh, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(refreshSecret), nil
	})

	if err != nil {
		con.JSON(http.StatusUnauthorized, Error{"token expired"})
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		con.JSON(http.StatusUnauthorized, Error{"token unvalid"})
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uuid, ok := claims["refresh_uuid"].(string)
		if !ok {
			con.JSON(http.StatusUnprocessableEntity, Error{"couldn't extract uuid"})
			return
		}
		username, ok := claims["username"].(string)
		if !ok {
			con.JSON(http.StatusUnprocessableEntity, Error{"couldn't extract username"})
			return
		}
		DeleteToken(uuid)
		tok, err := CreateToken(username)
		if err != nil {
			con.JSON(http.StatusForbidden, Error{err.Error()})
			return
		}
		SaveToken(username, tok)
		tokens := TokenPair{
			tok.AccessToken,
			tok.RefreshToken,
		}
		con.JSON(http.StatusCreated, tokens)
	} else {
		con.JSON(http.StatusUnauthorized, Error{"refresh token expired"})
	}
}

// GetTeacherByShort represents the get teacher by short name endpoint
// @Summary Returns a teacher with the specified short name
// @Description Searches for the Teacher with the specified name and returns the data
// @ID get-teacher-by-short
// @Accept json
// @Produce json
// @Param Authorization header string true "Access Token" default(Bearer <Add access token here>)
// @Param name query string true "Short Name of Teacher"
// @Success 200 {object} db.Teacher
// @Failure 401 {object} Error
// @Failure 422 {object} Error
// @Failure 500 {object} Error
// @Router /getTeacherByShort [get]
func GetTeacherByShort(con *gin.Context) {
	_, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, Error{"you are not logged in"})
		return
	}
	query := con.Request.URL.Query()
	if query.Get("name") == "" {
		con.JSON(http.StatusUnprocessableEntity, Error{"invalid request structure provided"})
		return
	}
	name := query.Get("name")
	db := mongo.MongoDatabaseConnector{}
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, Error{"database didn't respond"})
		return
	}
	defer db.Close()
	teacher := db.GetTeacherByShort(name)
	con.JSON(http.StatusOK, teacher)
}

// GetTeacher represents the get teacher endpoint
// @Summary Returns a teacher with the specified UUID
// @Description Searches for the Teacher with the specified uuid and returns the data
// @ID get-teacher
// @Accept json
// @Produce json
// @Param Authorization header string true "Access Token" default(Bearer <Add access token here>)
// @Param uuid query string true "UUID of Teacher"
// @Success 200 {object} db.Teacher
// @Failure 401 {object} Error
// @Failure 422 {object} Error
// @Failure 500 {object} Error
// @Router /getTeacher [get]
func GetTeacher(con *gin.Context) {
	_, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, Error{"you are not logged in"})
		return
	}
	query := con.Request.URL.Query()
	if query.Get("uuid") == "" {
		con.JSON(http.StatusUnprocessableEntity, Error{"invalid request structure provided"})
		return
	}
	uuid := query.Get("uuid")
	db := mongo.MongoDatabaseConnector{}
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, Error{"database didn't respond"})
		return
	}
	defer db.Close()
	teacher := db.GetTeacherByUUID(uuid)
	con.JSON(http.StatusOK, teacher)
}

// SetTeacherPermissions represents the set teacher permissions endpoint
// @Summary Sets the permissions of a Teacher
// @Description Sets the permissions of a Teacher to update their access rights
// @ID set-teacher-permissions
// @Accept json
// @Produce json
// @Param Authorization header string true "Access Token" default(Bearer <Add access token here>)
// @Param perm body Permissions true "Permission data of the teacher"
// @Param uuid query string true "UUID of the teacher whos permissions will be changed"
// @Success 200 {object} db.Teacher
// @Failure 401 {object} Error
// @Failure 422 {object} Error
// @Failure 500 {object} Error
// @Router /setTeacherPermissions [post]
func SetTeacherPermissions(con *gin.Context) {
	auth, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, Error{"you are not logged in"})
		return
	}
	perm := Permissions{}
	if err := con.ShouldBindJSON(&perm); err != nil {
		con.JSON(http.StatusUnprocessableEntity, Error{"invalid request structure provided"})
		return
	}
	query := con.Request.URL.Query()
	if query.Get("uuid") == "" {
		con.JSON(http.StatusUnprocessableEntity, Error{"invalid request structure provided"})
		return
	}
	uuid := query.Get("uuid")
	db := mongo.MongoDatabaseConnector{}
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, Error{"database didn't respond"})
		return
	}
	defer db.Close()
	requester := db.GetTeacherByShort(auth.Username)
	if !(requester.PEK || requester.Administration || requester.AV || requester.SuperUser) {
		con.JSON(http.StatusUnauthorized, Error{"unauthorized"})
		return
	}
	teacher := db.GetTeacherByUUID(uuid)
	teacher.SuperUser = perm.SuperUser
	teacher.Administration = perm.Administration
	teacher.PEK = perm.PEK
	teacher.Administration = perm.Administration
	if db.UpdateTeacher(uuid, teacher) {
		con.JSON(http.StatusOK, Information{"permissions updated"})
	} else {
		con.JSON(http.StatusInternalServerError, Error{"permissions couldn't be updated"})
	}
}

// GetActiveApplications represents the get active applications endpoint
// @Summary Returns all active applications
// @Description Returns all active applications as a list of applications
// @ID get-all-active-applications
// @Accept json
// @Produce json
// @Param Authorization header string true "Access Token" default(Bearer <Add access token here>)
// @Param username query string false "Filter to only show applications of this teacher"
// @Success 200 {array} db.Application
// @Failure 401 {object} Error
// @Failure 500 {object} Error
// @Router /getActiveApplications [get]
func GetActiveApplications(con *gin.Context) {
	auth, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, Error{"you are not logged in"})
		return
	}
	db := mongo.MongoDatabaseConnector{}
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, Error{"database didn't respond"})
		return
	}
	defer db.Close()
	_ = con.Request.ParseForm()
	query := con.Request.URL.Query()
	_, applyFilter := con.Request.Form["username"]
	filter := query.Get("username")
	requestTeacher := db.GetTeacherByShort(auth.Username)
	if !(requestTeacher.Administration || requestTeacher.AV || requestTeacher.SuperUser || requestTeacher.PEK || (applyFilter && requestTeacher.Short == filter)) {
		con.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	applications := db.GetActiveApplications()
	if applyFilter {
		res := make([]mongo.Application, 0)
		for _, app := range applications {
			if app.Kind == mongo.SchoolEvent {
				teachers := app.SchoolEventDetails.Teachers
				for _, t := range teachers {
					if t.Shortname == filter {
						res = append(res, app)
						break
					}
				}
			} else if app.Kind == mongo.Training {
				if app.TrainingDetails.Organizer == filter {
					res = append(res, app)
				}
			} else if app.Kind == mongo.OtherReason {
				if app.OtherReasonDetails.Filer == filter {
					res = append(res, app)
				}
			}
		}
		con.JSON(http.StatusOK, res)
		return
	}
	con.JSON(http.StatusOK, applications)
}

// GetAllApplications represents the get all applications endpoint
// @Summary Returns all applications
// @Description Returns all applications as a list of applications
// @ID get-all-applications
// @Accept json
// @Produce json
// @Param Authorization header string true "Access Token" default(Bearer <Add access token here>)
// @Param username query string false "Filter to only show applications of this teacher"
// @Success 200 {array} db.Application
// @Failure 401 {object} Error
// @Failure 500 {object} Error
// @Router /getAllApplications [get]
func GetAllApplications(con *gin.Context) {
	auth, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, Error{"you are not logged in"})
		return
	}
	db := mongo.MongoDatabaseConnector{}
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, Error{"database didn't respond"})
		return
	}
	defer db.Close()
	_ = con.Request.ParseForm()
	query := con.Request.URL.Query()
	_, applyFilter := con.Request.Form["username"]
	filter := query.Get("username")
	requestTeacher := db.GetTeacherByShort(auth.Username)
	if !(requestTeacher.Administration || requestTeacher.AV || requestTeacher.SuperUser || requestTeacher.PEK || (applyFilter && requestTeacher.Short == filter)) {
		con.JSON(http.StatusUnauthorized, Error{"unauthorized"})
		return
	}
	applications := db.GetAllApplications()
	if applyFilter {
		res := make([]mongo.Application, 0)
		for _, app := range applications {
			if app.Kind == mongo.SchoolEvent {
				teachers := app.SchoolEventDetails.Teachers
				for _, t := range teachers {
					if t.Shortname == filter {
						res = append(res, app)
						break
					}
				}
			} else if app.Kind == mongo.Training {
				if app.TrainingDetails.Organizer == filter {
					res = append(res, app)
				}
			} else if app.Kind == mongo.OtherReason {
				if app.OtherReasonDetails.Filer == filter {
					res = append(res, app)
				}
			}
		}
		con.JSON(http.StatusOK, res)
		return
	}
	con.JSON(http.StatusOK, applications)
}

// GetNews represents the get news endpoint
// @Summary Returns the news
// @Description Returns the 10 last changed applications
// @ID get-news
// @Accept json
// @Produce json
// @Param Authorization header string true "Access Token" default(Bearer <Add access token here>)
// @Success 200 {array} News
// @Failure 401 {object} Error
// @Failure 500 {object} Error
// @Router /getNews [get]
func GetNews(con *gin.Context) {
	auth, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, Error{"unauthorized"})
		return
	}
	db := mongo.MongoDatabaseConnector{}
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, Error{"database didn't respond"})
		return
	}
	defer db.Close()
	user := auth.Username
	applications := db.GetActiveApplications()
	res := make([]mongo.Application, 0)
	for _, app := range applications {
		if app.Kind == mongo.SchoolEvent {
			teachers := app.SchoolEventDetails.Teachers
			for _, t := range teachers {
				if t.Shortname == user {
					res = append(res, app)
					break
				}
			}
		} else if app.Kind == mongo.Training {
			if app.TrainingDetails.Organizer == user {
				res = append(res, app)
			}
		} else if app.Kind == mongo.OtherReason {
			if app.OtherReasonDetails.Filer == user {
				res = append(res, app)
			}
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].LastChanged.After(res[j].LastChanged)
	})
	if len(res) > 10 {
		res = res[0:10]
	}
	news := make([]News, 0)
	for _, app := range res {
		news = append(news, News{app.UUID, app.Name, app.Progress, app.LastChanged.String()})
	}
	con.JSON(http.StatusOK, news)
}

// GetApplication represents the get application endpoint
// @Summary Returns an Application
// @Description Returns the Application matching the given UUID
// @ID get-application
// @Accept json
// @Produce json
// @Param Authorization header string true "Access Token" default(Bearer <Add access token here>)
// @Param uuid query string true "The UUID of the specifying Application"
// @Success 200 {object} db.Application
// @Failure 401 {object} Error
// @Failure 422 {object} Error
// @Failure 500 {object} Error
// @Router /getApplication [get]
func GetApplication(con *gin.Context) {
	auth, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, Error{"you are not logged in"})
		return
	}
	db := mongo.MongoDatabaseConnector{}
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, Error{"database didn't respond"})
		return
	}
	defer db.Close()
	query := con.Request.URL.Query()
	uuid := query.Get("uuid")
	if query.Get("uuid") == "" {
		con.JSON(http.StatusUnprocessableEntity, Error{"invalid request structure provided"})
		return
	}
	requestTeacher := db.GetTeacherByShort(auth.Username)
	application := db.GetApplication(uuid)
	var in bool
	if application.Kind == mongo.SchoolEvent {
		teachers := application.SchoolEventDetails.Teachers
		for _, t := range teachers {
			if t.Shortname == requestTeacher.Short {
				in = true
				break
			}
		}
	} else if application.Kind == mongo.Training {
		if application.TrainingDetails.Organizer == requestTeacher.Short {
			in = true
		}
	} else if application.Kind == mongo.OtherReason {
		if application.OtherReasonDetails.Filer == requestTeacher.Short {
			in = true
		}
	}
	if !(in || requestTeacher.Administration || requestTeacher.AV || requestTeacher.PEK || requestTeacher.SuperUser) {
		con.JSON(http.StatusUnauthorized, Error{"unauthorized"})
		return
	}
	con.JSON(http.StatusOK, application)
}

// GetAdminApplications represents the get admin applications endpoint
// @Summary Returns all admin applications
// @Description Returns all applications currently needing a review by an admin
// @ID get-admin-applications
// @Accept json
// @Produce json
// @Param Authorization header string true "Access Token" default(Bearer <Add access token here>)
// @Success 200 {array} db.Application
// @Failure 401 {object} Error
// @Failure 500 {object} Error
// @Router /getAdminApplication [get]
func GetAdminApplications(con *gin.Context) {
	auth, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, Error{"you are not logged in"})
		return
	}
	db := mongo.MongoDatabaseConnector{}
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, Error{"database didn't respond"})
		return
	}
	defer db.Close()
	teacher := db.GetTeacherByShort(auth.Username)
	if !(teacher.PEK || teacher.Administration || teacher.AV || teacher.SuperUser) {
		con.JSON(http.StatusUnauthorized, Error{"unauthorized"})
		return
	}
	applications := db.GetAllApplications()
	res := make([]mongo.Application, 0)
	for _, app := range applications {
		if app.Kind == mongo.SchoolEvent {
			if app.Progress == mongo.SEInProcess && (teacher.Administration || teacher.AV || teacher.SuperUser) {
				res = append(res, app)
			}
			if app.Progress == mongo.SECostsInProcess {
				res = append(res, app)
			}
		} else if app.Kind == mongo.Training || app.Kind == mongo.OtherReason {
			if app.Progress == mongo.TInProcess && (teacher.Administration || teacher.AV || teacher.SuperUser) {
				res = append(res, app)
			}
			if app.Progress == mongo.TCostsInProcess {
				res = append(res, app)
			}
		}
	}
	con.JSON(http.StatusOK, res)
}

// CreateApplication represents the create applications endpoint
// @Summary Creates a new application
// @Description Creates the provided application in the system
// @ID create-application
// @Accept json
// @Produce json
// @Param Authorization header string true "Access Token" default(Bearer <Add access token here>)
// @Param application body db.Application true "The Application Data"
// @Success 200 {object} Information
// @Failure 401 {object} Error
// @Failure 422 {object} Error
// @Failure 500 {object} Error
// @Router /createApplication [post]
func CreateApplication(con *gin.Context) {
	app := mongo.Application{}
	if err := con.ShouldBindJSON(&app); err != nil {
		con.JSON(http.StatusUnprocessableEntity, Error{"invalid request structure provided"})
		return
	}
	_, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, Error{"you are not logged in"})
		return
	}
	db := mongo.MongoDatabaseConnector{}
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, Error{"database didn't respond"})
		return
	}
	defer db.Close()
	if db.CreateApplication(app) {
		con.JSON(http.StatusOK, Information{"success; application created"})
	} else {
		con.JSON(http.StatusInternalServerError, Error{"error; application not created"})
	}
}

// UpdateApplication represents the update applications endpoint
// @Summary Updates an existing application
// @Description Updates an application identified by a uuid with the data in the body in the system
// @ID update-application
// @Accept json
// @Produce json
// @Param Authorization header string true "Access Token" default(Bearer <Add access token here>)
// @Param application body db.Application true "The application data to update"
// @Param uuid query string true "Identifier of the application to update"
// @Success 200 {object} Information
// @Failure 401 {object} Error
// @Failure 422 {object} Error
// @Failure 500 {object} Error
// @Router /updateApplication [put]
func UpdateApplication(con *gin.Context) {
	app := mongo.Application{}
	if err := con.ShouldBindJSON(&app); err != nil {
		con.JSON(http.StatusUnprocessableEntity, Error{"invalid request structure provided"})
		return
	}
	auth, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, Error{"you are not logged in"})
		return
	}
	db := mongo.MongoDatabaseConnector{}
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, Error{"database didn't respond"})
		return
	}
	defer db.Close()
	query := con.Request.URL.Query()
	uuid := query.Get("uuid")
	if query.Get("uuid") == "" {
		con.JSON(http.StatusUnprocessableEntity, Error{"invalid request structure provided"})
		return
	}
	requestTeacher := db.GetTeacherByShort(auth.Username)
	application := db.GetApplication(uuid)
	var in bool
	if application.Kind == mongo.SchoolEvent {
		teachers := application.SchoolEventDetails.Teachers
		for _, t := range teachers {
			if t.Shortname == requestTeacher.Short {
				in = true
				break
			}
		}
	} else if application.Kind == mongo.Training {
		if application.TrainingDetails.Organizer == requestTeacher.Short {
			in = true
		}
	} else if application.Kind == mongo.OtherReason {
		if application.OtherReasonDetails.Filer == requestTeacher.Short {
			in = true
		}
	}
	if !(in || requestTeacher.Administration || requestTeacher.AV || requestTeacher.PEK || requestTeacher.SuperUser) {
		con.JSON(http.StatusUnauthorized, Error{"unauthorized"})
		return
	}
	if db.UpdateApplication(uuid, app) {
		con.JSON(http.StatusOK, Information{"success; application updated"})
	} else {
		con.JSON(http.StatusInternalServerError, Error{"error; application not updated"})
	}
}

// DeleteApplication represents the delete applications endpoint
// @Summary Deletes an existing application
// @Description Deletes an application identified by a uuid
// @ID delete-application
// @Accept json
// @Produce json
// @Param Authorization header string true "Access Token" default(Bearer <Add access token here>)
// @Param uuid query string true "Identifier of the application to delete"
// @Success 200 {object} Information
// @Failure 401 {object} Error
// @Failure 422 {object} Error
// @Failure 500 {object} Error
// @Router /deleteApplication [delete]
func DeleteApplication(con *gin.Context) {
	auth, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, Error{"you are not logged in"})
		return
	}
	db := mongo.MongoDatabaseConnector{}
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, Error{"database didn't respond"})
		return
	}
	defer db.Close()
	query := con.Request.URL.Query()
	uuid := query.Get("uuid")
	if query.Get("uuid") == "" {
		con.JSON(http.StatusUnprocessableEntity, Error{"invalid request structure provided"})
		return
	}
	requestTeacher := db.GetTeacherByShort(auth.Username)
	application := db.GetApplication(uuid)
	var in bool
	if application.Kind == mongo.SchoolEvent {
		teachers := application.SchoolEventDetails.Teachers
		for _, t := range teachers {
			if t.Shortname == requestTeacher.Short {
				in = true
				break
			}
		}
	} else if application.Kind == mongo.Training {
		if application.TrainingDetails.Organizer == requestTeacher.Short {
			in = true
		}
	} else if application.Kind == mongo.OtherReason {
		if application.OtherReasonDetails.Filer == requestTeacher.Short {
			in = true
		}
	}
	if !(in || requestTeacher.Administration || requestTeacher.AV || requestTeacher.PEK || requestTeacher.SuperUser) {
		con.JSON(http.StatusUnauthorized, Error{"unauthorized"})
		return
	}
	if db.DeleteApplication(uuid) {
		con.JSON(http.StatusOK, Information{"success; application deleted"})
	} else {
		con.JSON(http.StatusInternalServerError, Error{"error; application not deleted"})
	}
}

// GetAbsenceFormForClasses represents get absence form for classes endpoint
// @Summary Generates an absence form for classes
// @Description Generates an absence form for classes and returns it
// @ID get-absence-form-for-classes
// @Accept json
// @Produce json
// @Param Authorization header string true "Access Token" default(Bearer <Add access token here>)
// @Param uuid query string true "Identifier of the application to generate the pdf from"
// @Param classes query []string false "Filter for classes"
// @Success 200 {object} PDF
// @Failure 401 {object} Error
// @Failure 422 {object} Error
// @Failure 500 {object} Error
// @Router /getAbsenceFormForClasses [get]
func GetAbsenceFormForClasses(con *gin.Context) {
	auth, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, Error{"you are not logged in"})
	}
	db := mongo.MongoDatabaseConnector{}
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, Error{"database didn't respond"})
		return
	}
	defer db.Close()
	_ = con.Request.ParseForm()
	query := con.Request.URL.Query()
	if _, hasUUID := con.Request.Form["uuid"]; !hasUUID {
		con.JSON(http.StatusUnprocessableEntity, Error{"invalid request structure provided"})
		return
	}
	uuid := query.Get("uuid")
	_, applyClassFilter := con.Request.Form["classes"]
	classes := make([]string, 0)
	if applyClassFilter {
		classes = query["classes"]
	}
	application := db.GetApplication(uuid)
	requestTeacher := db.GetTeacherByShort(auth.Username)
	var in bool
	if application.Kind == mongo.SchoolEvent {
		teachers := application.SchoolEventDetails.Teachers
		for _, t := range teachers {
			if t.Shortname == requestTeacher.Short {
				in = true
				break
			}
		}
	} else if application.Kind == mongo.Training {
		if application.TrainingDetails.Organizer == requestTeacher.Short {
			in = true
		}
	} else if application.Kind == mongo.OtherReason {
		if application.OtherReasonDetails.Filer == requestTeacher.Short {
			in = true
		}
	}
	if !(in || requestTeacher.Administration || requestTeacher.AV || requestTeacher.PEK || requestTeacher.SuperUser) {
		con.JSON(http.StatusUnauthorized, Error{"you have no permission to do this"})
		return
	}
	path, err := files.GenerateFileEnvironment(application)
	if err != nil {
		con.JSON(http.StatusInternalServerError, Error{"couldn't create directories"})
		return
	}
	paths, err := files.GenerateAbsenceFormForClass(path, auth.Username, application)
	if err != nil {
		con.JSON(http.StatusInternalServerError, Error{"couldn't create pdfs"})
		return
	}

	pdfs := make(map[string]string)
	for _, p := range paths {
		splits := strings.Split(p, string(filepath.Separator))
		filename := strings.Split(splits[len(splits)-1], ".")[0]
		names := strings.Split(filename, "_")
		class := names[len(names)-1]
		pdfs[class] = p
	}
	if applyClassFilter {
		for class := range pdfs {
			in := false
			for _, c := range classes {
				if c == class {
					in = true
				}
			}
			if !in {
				delete(pdfs, class)
			}
		}
	}
	pp := make([]string, 0)
	for _, p := range pdfs {
		pp = append(pp, p)
	}
	created := filepath.Join(filepath.Dir(pp[0]), fmt.Sprintf(files.ClassAbsenceFormFileName, "merge"))
	err = api.MergeCreateFile(pp, created, pdfcpu.NewDefaultConfiguration())
	if err != nil {
		con.JSON(http.StatusInternalServerError, Error{"couldn't save merged pdf"})
		return
	}
	file, err := ioutil.ReadFile(created)
	if err != nil {
		con.JSON(http.StatusInternalServerError, Error{"couldn't read merged pdf"})
		return
	}
	enc := base64.StdEncoding.EncodeToString(file)
	res := PDF{enc}
	err = os.Remove(created)
	if err != nil {
		con.JSON(http.StatusInternalServerError, Error{"couldn't delete merged pdf"})
		return
	}
	con.JSON(http.StatusOK, res)
}

// GetAbsenceFormForTeacher represents get absence form for teacher endpoint
// @Summary Generates an absence form for a teacher
// @Description Generates an absence form for a teacher and returns it
// @ID get-absence-form-for-teacher
// @Accept json
// @Produce json
// @Param Authorization header string true "Access Token" default(Bearer <Add access token here>)
// @Param uuid query string true "Identifier of the application to generate the pdf from"
// @Param teacher query string false "untis name of the teacher, if not provided logged in teacher will be used"
// @Success 200 {object} PDF
// @Failure 401 {object} Error
// @Failure 422 {object} Error
// @Failure 500 {object} Error
// @Router /getAbsenceFormForTeacher [get]
func GetAbsenceFormForTeacher(con *gin.Context) {
	auth, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, Error{"you are not logged in"})
	}
	db := mongo.MongoDatabaseConnector{}
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, Error{"database didn't respond"})
		return
	}
	defer db.Close()
	_ = con.Request.ParseForm()
	query := con.Request.URL.Query()
	if _, hasUUID := con.Request.Form["uuid"]; !hasUUID {
		con.JSON(http.StatusUnprocessableEntity, Error{"invalid request structure provided"})
		return
	}
	uuid := query.Get("uuid")
	_, applyTeacher := con.Request.Form["teacher"]
	teacher := ""
	if applyTeacher {
		teacher = query.Get("teacher")
	}
	application := db.GetApplication(uuid)
	requestTeacher := db.GetTeacherByShort(auth.Username)
	var in bool
	if application.Kind == mongo.SchoolEvent {
		teachers := application.SchoolEventDetails.Teachers
		for _, t := range teachers {
			if t.Shortname == requestTeacher.Short {
				in = true
				break
			}
		}
	} else if application.Kind == mongo.Training {
		if application.TrainingDetails.Organizer == requestTeacher.Short {
			in = true
		}
	} else if application.Kind == mongo.OtherReason {
		if application.OtherReasonDetails.Filer == requestTeacher.Short {
			in = true
		}
	}
	if !(!applyTeacher && in) {
		con.JSON(http.StatusUnauthorized, Error{"you have no permission to do this"})
		return
	}
	if !(applyTeacher && (requestTeacher.Administration || requestTeacher.AV || requestTeacher.PEK || requestTeacher.SuperUser)) {
		con.JSON(http.StatusUnauthorized, Error{"you have no permission to do this"})
		return
	}
	path, err := files.GenerateFileEnvironment(application)
	if err != nil {
		con.JSON(http.StatusInternalServerError, Error{"couldn't create directories"})
		return
	}
	if applyTeacher {
		reqTeacher := db.GetTeacherByShort(teacher)
		path, err = files.GenerateAbsenceFormForTeacher(path, auth.Username, reqTeacher.Longname, application)
	} else {
		path, err = files.GenerateAbsenceFormForTeacher(path, auth.Username, "self", application)
	}
	if err != nil {
		con.JSON(http.StatusInternalServerError, Error{"couldn't create pdfs"})
		return
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		con.JSON(http.StatusInternalServerError, Error{"couldn't read generated pdf"})
		return
	}
	enc := base64.StdEncoding.EncodeToString(file)
	res := PDF{enc}
	con.JSON(http.StatusOK, res)
}

// GetCompensationForEducationalSupportForm represents get compensation for educational support form endpoint
// @Summary Generates a compensation for educational support form for all teachers
// @Description Generates a compensation for educational support form for all teachers and returns it
// @ID get-compensation-for-educational-support-form
// @Accept json
// @Produce json
// @Param Authorization header string true "Access Token" default(Bearer <Add access token here>)
// @Param uuid query string true "Identifier of the application to generate the pdf from"
// @Success 200 {object} PDF
// @Failure 401 {object} Error
// @Failure 422 {object} Error
// @Failure 500 {object} Error
// @Router /getCompensationForEducationalSupportForm [get]
func GetCompensationForEducationalSupportForm(con *gin.Context) {
	auth, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, Error{"you are not logged in"})
	}
	db := mongo.MongoDatabaseConnector{}
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, Error{"database didn't respond"})
		return
	}
	defer db.Close()
	_ = con.Request.ParseForm()
	query := con.Request.URL.Query()
	if _, hasUUID := con.Request.Form["uuid"]; !hasUUID {
		con.JSON(http.StatusUnprocessableEntity, Error{"invalid request structure provided"})
		return
	}
	uuid := query.Get("uuid")
	application := db.GetApplication(uuid)
	requestTeacher := db.GetTeacherByShort(auth.Username)
	var in bool
	if application.Kind == mongo.SchoolEvent {
		teachers := application.SchoolEventDetails.Teachers
		for _, t := range teachers {
			if t.Shortname == requestTeacher.Short {
				in = true
				break
			}
		}
	} else if application.Kind == mongo.Training {
		if application.TrainingDetails.Organizer == requestTeacher.Short {
			in = true
		}
	} else if application.Kind == mongo.OtherReason {
		if application.OtherReasonDetails.Filer == requestTeacher.Short {
			in = true
		}
	}
	if !(in || requestTeacher.Administration || requestTeacher.AV || requestTeacher.PEK || requestTeacher.SuperUser) {
		con.JSON(http.StatusUnauthorized, Error{"you have no permission to do this"})
		return
	}
	path, err := files.GenerateFileEnvironment(application)
	if err != nil {
		con.JSON(http.StatusInternalServerError, Error{"couldn't create directories"})
		return
	}
	path, err = files.GenerateCompensationForEducationalSupport(path, application)
	if err != nil {
		con.JSON(http.StatusInternalServerError, Error{"couldn't create pdfs"})
		return
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		con.JSON(http.StatusInternalServerError, Error{"couldn't read generated pdf"})
		return
	}
	enc := base64.StdEncoding.EncodeToString(file)
	res := PDF{enc}
	con.JSON(http.StatusOK, res)
}

// GetTravelInvoiceForm represents get travel invoice form endpoint
// @Summary Generates a travel invoice for a teacher
// @Description Generates a travel invoice form for a teacher and returns it
// @ID get-travel-invoice-form
// @Accept json
// @Produce json
// @Param Authorization header string true "Access Token" default(Bearer <Add access token here>)
// @Param uuid query string true "Identifier of the application to generate the pdf from"
// @Param short query string true "Short name of the teacher this should be generated for"
// @Param ti_id query int true "ID of the Travel Invoice data"
// @Param receipts query bool false "If provided the pdf will include all receipt"
// @Success 200 {object} PDF
// @Failure 401 {object} Error
// @Failure 422 {object} Error
// @Failure 500 {object} Error
// @Router /getTravelInvoiceForm [get]
func GetTravelInvoiceForm(con *gin.Context) {
	auth, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, Error{"you are not logged in"})
	}
	db := mongo.MongoDatabaseConnector{}
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, Error{"database didn't respond"})
		return
	}
	defer db.Close()
	_ = con.Request.ParseForm()
	query := con.Request.URL.Query()
	if _, hasUUID := con.Request.Form["uuid"]; !hasUUID {
		con.JSON(http.StatusUnprocessableEntity, Error{"invalid request structure provided"})
		return
	}
	uuid := query.Get("uuid")
	if _, hasShort := con.Request.Form["short"]; !hasShort {
		con.JSON(http.StatusUnprocessableEntity, Error{"invalid request structure provided"})
		return
	}
	short := query.Get("short")
	if _, hasTIID := con.Request.Form["ti_id"]; !hasTIID {
		con.JSON(http.StatusUnprocessableEntity, Error{"invalid request structure provided"})
		return
	}
	tiID, err := strconv.Atoi(query.Get("ti_id"))
	if err != nil {
		con.JSON(http.StatusUnprocessableEntity, Error{"invalid ti_id provided"})
		return
	}
	_, applyMergeReceipts := con.Request.Form["receipts"]
	application := db.GetApplication(uuid)
	requestTeacher := db.GetTeacherByShort(auth.Username)
	var in bool
	if application.Kind == mongo.SchoolEvent {
		teachers := application.SchoolEventDetails.Teachers
		for _, t := range teachers {
			if t.Shortname == requestTeacher.Short {
				in = true
				break
			}
		}
	} else if application.Kind == mongo.Training {
		if application.TrainingDetails.Organizer == requestTeacher.Short {
			in = true
		}
	} else if application.Kind == mongo.OtherReason {
		if application.OtherReasonDetails.Filer == requestTeacher.Short {
			in = true
		}
	}
	if !(in || requestTeacher.Administration || requestTeacher.AV || requestTeacher.PEK || requestTeacher.SuperUser) {
		con.JSON(http.StatusUnauthorized, Error{"you have no permission to do this"})
		return
	}
	path, err := files.GenerateFileEnvironment(application)
	if err != nil {
		con.JSON(http.StatusInternalServerError, Error{"couldn't create directories"})
		return
	}
	var ti mongo.TravelInvoice
	for _, tis := range application.TravelInvoices {
		if tis.ID == tiID {
			ti = tis
			break
		}
	}
	path, err = files.GenerateTravelInvoice(path, short, ti, application.UUID)
	if err != nil {
		con.JSON(http.StatusInternalServerError, Error{"couldn't create pdfs"})
		return
	}

	if applyMergeReceipts {
		pp := append(make([]string, 0), path)
		uploadFolder := filepath.Join(filepath.Dir(path), files.UploadFolderName)
		ff, err := ioutil.ReadDir(uploadFolder)
		if err != nil {
			con.JSON(http.StatusInternalServerError, Error{"couldn't read upload directory"})
			return
		}
		for _, file := range ff {
			data := strings.Split(file.Name(), "_")
			if data[1] == short {
				pp = append(pp, filepath.Join(uploadFolder, file.Name()))
			}
		}
		created := filepath.Join(filepath.Dir(path), fmt.Sprintf(files.TravelInvoicePDFFileName, short+"_merge"))
		err = api.MergeCreateFile(pp, created, pdfcpu.NewDefaultConfiguration())
		if err != nil {
			con.JSON(http.StatusInternalServerError, Error{"couldn't save merged pdf"})
			return
		}
		file, err := ioutil.ReadFile(created)
		if err != nil {
			con.JSON(http.StatusInternalServerError, Error{"couldn't read generated pdf"})
			return
		}
		enc := base64.StdEncoding.EncodeToString(file)
		res := PDF{enc}
		err = os.Remove(created)
		if err != nil {
			con.JSON(http.StatusInternalServerError, Error{"couldn't delete merged pdf"})
			return
		}
		con.JSON(http.StatusOK, res)
	} else {
		file, err := ioutil.ReadFile(path)
		if err != nil {
			con.JSON(http.StatusInternalServerError, Error{"couldn't read generated pdf"})
			return
		}
		enc := base64.StdEncoding.EncodeToString(file)
		res := PDF{enc}
		con.JSON(http.StatusOK, res)
	}
}

// GetBusinessTripApplicationForm represents get business application form endpoint
// @Summary Generates a business trip application form for a teacher
// @Description Generates a business trip application form for a teacher and returns it
// @ID get-business-trip-application-form
// @Accept json
// @Produce json
// @Param Authorization header string true "Access Token" default(Bearer <Add access token here>)
// @Param uuid query string true "Identifier of the application to generate the form from"
// @Param short query string true "Short name of the teacher this should be generated for"
// @Param bta_id query int true "ID of the Business Trip Application data"
// @Success 200 {object} PDF
// @Failure 401 {object} Error
// @Failure 422 {object} Error
// @Failure 500 {object} Error
// @Router /getBusinessTripApplicationForm [get]
func GetBusinessTripApplicationForm(con *gin.Context) {
	auth, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, Error{"you are not logged in"})
	}
	db := mongo.MongoDatabaseConnector{}
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, Error{"database didn't respond"})
		return
	}
	defer db.Close()
	_ = con.Request.ParseForm()
	query := con.Request.URL.Query()
	if _, hasUUID := con.Request.Form["uuid"]; !hasUUID {
		con.JSON(http.StatusUnprocessableEntity, Error{"invalid request structure provided"})
		return
	}
	uuid := query.Get("uuid")
	if _, hasShort := con.Request.Form["short"]; !hasShort {
		con.JSON(http.StatusUnprocessableEntity, Error{"invalid request structure provided"})
		return
	}
	short := query.Get("short")
	if _, hasBTAID := con.Request.Form["bta_id"]; !hasBTAID {
		con.JSON(http.StatusUnprocessableEntity, Error{"invalid request structure provided"})
		return
	}
	btaID, err := strconv.Atoi(query.Get("bta_id"))
	if err != nil {
		con.JSON(http.StatusUnprocessableEntity, Error{"invalid bta_id provided"})
		return
	}
	application := db.GetApplication(uuid)
	requestTeacher := db.GetTeacherByShort(auth.Username)
	var in bool
	if application.Kind == mongo.SchoolEvent {
		teachers := application.SchoolEventDetails.Teachers
		for _, t := range teachers {
			if t.Shortname == requestTeacher.Short {
				in = true
				break
			}
		}
	} else if application.Kind == mongo.Training {
		if application.TrainingDetails.Organizer == requestTeacher.Short {
			in = true
		}
	} else if application.Kind == mongo.OtherReason {
		if application.OtherReasonDetails.Filer == requestTeacher.Short {
			in = true
		}
	}
	if !(in || requestTeacher.Administration || requestTeacher.AV || requestTeacher.PEK || requestTeacher.SuperUser) {
		con.JSON(http.StatusUnauthorized, Error{"you have no permission to do this"})
		return
	}
	path, err := files.GenerateFileEnvironment(application)
	if err != nil {
		con.JSON(http.StatusInternalServerError, Error{"couldn't create directories"})
		return
	}
	var bta mongo.BusinessTripApplication
	for _, btas := range application.BusinessTripApplications {
		if btas.ID == btaID {
			bta = btas
			break
		}
	}
	path, err = files.GenerateBusinessTripApplication(path, short, bta, application.UUID)
	if err != nil {
		con.JSON(http.StatusInternalServerError, Error{"couldn't create pdf"})
		return
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		con.JSON(http.StatusInternalServerError, Error{"couldn't read generated pdf"})
		return
	}
	enc := base64.StdEncoding.EncodeToString(file)
	res := PDF{enc}
	con.JSON(http.StatusOK, res)
}

// GetTravelInvoiceExcel represents get travel invoice excel endpoint
// @Summary Generates a travel invoice excel for a teacher
// @Description Generates a travel invoice excel for a teacher and returns it
// @ID get-travel-invoice-excel
// @Accept json
// @Produce json
// @Param Authorization header string true "Access Token" default(Bearer <Add access token here>)
// @Param uuid query string true "Identifier of the application to generate the excel from"
// @Param short query string true "Short name of the teacher this should be generated for"
// @Param ti_id query int true "ID of the Travel Invoice data"
// @Success 200 {object} Excel
// @Failure 401 {object} Error
// @Failure 422 {object} Error
// @Failure 500 {object} Error
// @Router /getTravelInvoiceExcel [get]
func GetTravelInvoiceExcel(con *gin.Context) {
	auth, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, Error{"you are not logged in"})
	}
	db := mongo.MongoDatabaseConnector{}
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, Error{"database didn't respond"})
		return
	}
	defer db.Close()
	_ = con.Request.ParseForm()
	query := con.Request.URL.Query()
	if _, hasUUID := con.Request.Form["uuid"]; !hasUUID {
		con.JSON(http.StatusUnprocessableEntity, Error{"invalid request structure provided"})
		return
	}
	uuid := query.Get("uuid")
	if _, hasShort := con.Request.Form["short"]; !hasShort {
		con.JSON(http.StatusUnprocessableEntity, Error{"invalid request structure provided"})
		return
	}
	short := query.Get("short")
	if _, hasTIID := con.Request.Form["ti_id"]; !hasTIID {
		con.JSON(http.StatusUnprocessableEntity, Error{"invalid request structure provided"})
		return
	}
	tiID, err := strconv.Atoi(query.Get("ti_id"))
	if err != nil {
		con.JSON(http.StatusUnprocessableEntity, Error{"invalid ti_id provided"})
		return
	}
	application := db.GetApplication(uuid)
	requestTeacher := db.GetTeacherByShort(auth.Username)
	var in bool
	if application.Kind == mongo.SchoolEvent {
		teachers := application.SchoolEventDetails.Teachers
		for _, t := range teachers {
			if t.Shortname == requestTeacher.Short {
				in = true
				break
			}
		}
	} else if application.Kind == mongo.Training {
		if application.TrainingDetails.Organizer == requestTeacher.Short {
			in = true
		}
	} else if application.Kind == mongo.OtherReason {
		if application.OtherReasonDetails.Filer == requestTeacher.Short {
			in = true
		}
	}
	if !(in || requestTeacher.Administration || requestTeacher.AV || requestTeacher.PEK || requestTeacher.SuperUser) {
		con.JSON(http.StatusUnauthorized, Error{"you have no permission to do this"})
		return
	}
	path, err := files.GenerateFileEnvironment(application)
	if err != nil {
		con.JSON(http.StatusInternalServerError, Error{"couldn't create directories"})
		return
	}
	var ti mongo.TravelInvoice
	for _, tis := range application.TravelInvoices {
		if tis.ID == tiID {
			ti = tis
			break
		}
	}
	path, err = files.GenerateTravelInvoiceExcel(path, short, ti)
	if err != nil {
		con.JSON(http.StatusInternalServerError, Error{"couldn't create excel"})
		return
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		con.JSON(http.StatusInternalServerError, Error{"couldn't read generated excel"})
		return
	}
	enc := base64.StdEncoding.EncodeToString(file)
	res := Excel{enc}
	con.JSON(http.StatusOK, res)
}

// GetBusinessTripApplicationExcel represents get business application excel endpoint
// @Summary Generates a business trip application excel for a teacher
// @Description Generates a business trip application excel for a teacher and returns it
// @ID get-business-trip-application-excel
// @Accept json
// @Produce json
// @Param Authorization header string true "Access Token" default(Bearer <Add access token here>)
// @Param uuid query string true "Identifier of the application to generate the excel from"
// @Param short query string true "Short name of the teacher this should be generated for"
// @Param bta_id query int true "ID of the Business Trip Application data"
// @Success 200 {object} Excel
// @Failure 401 {object} Error
// @Failure 422 {object} Error
// @Failure 500 {object} Error
// @Router /getBusinessTripApplicationExcel [get]
func GetBusinessTripApplicationExcel(con *gin.Context) {
	auth, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, Error{"you are not logged in"})
	}
	db := mongo.MongoDatabaseConnector{}
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, Error{"database didn't respond"})
		return
	}
	defer db.Close()
	_ = con.Request.ParseForm()
	query := con.Request.URL.Query()
	if _, hasUUID := con.Request.Form["uuid"]; !hasUUID {
		con.JSON(http.StatusUnprocessableEntity, Error{"invalid request structure provided"})
		return
	}
	uuid := query.Get("uuid")
	if _, hasShort := con.Request.Form["short"]; !hasShort {
		con.JSON(http.StatusUnprocessableEntity, Error{"invalid request structure provided"})
		return
	}
	short := query.Get("short")
	if _, hasBTAID := con.Request.Form["bta_id"]; !hasBTAID {
		con.JSON(http.StatusUnprocessableEntity, Error{"invalid request structure provided"})
		return
	}
	btaID, err := strconv.Atoi(query.Get("bta_id"))
	if err != nil {
		con.JSON(http.StatusUnprocessableEntity, Error{"invalid bta_id provided"})
		return
	}
	application := db.GetApplication(uuid)
	requestTeacher := db.GetTeacherByShort(auth.Username)
	var in bool
	if application.Kind == mongo.SchoolEvent {
		teachers := application.SchoolEventDetails.Teachers
		for _, t := range teachers {
			if t.Shortname == requestTeacher.Short {
				in = true
				break
			}
		}
	} else if application.Kind == mongo.Training {
		if application.TrainingDetails.Organizer == requestTeacher.Short {
			in = true
		}
	} else if application.Kind == mongo.OtherReason {
		if application.OtherReasonDetails.Filer == requestTeacher.Short {
			in = true
		}
	}
	if !(in || requestTeacher.Administration || requestTeacher.AV || requestTeacher.PEK || requestTeacher.SuperUser) {
		con.JSON(http.StatusUnauthorized, Error{"you have no permission to do this"})
		return
	}
	path, err := files.GenerateFileEnvironment(application)
	if err != nil {
		con.JSON(http.StatusInternalServerError, Error{"couldn't create directories"})
		return
	}
	var bta mongo.BusinessTripApplication
	for _, btas := range application.BusinessTripApplications {
		if btas.ID == btaID {
			bta = btas
			break
		}
	}
	path, err = files.GenerateBusinessTripApplicationExcel(path, short, bta)
	if err != nil {
		con.JSON(http.StatusInternalServerError, Error{"couldn't create excel"})
		return
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		con.JSON(http.StatusInternalServerError, Error{"couldn't read generated excel"})
		return
	}
	enc := base64.StdEncoding.EncodeToString(file)
	res := Excel{enc}
	con.JSON(http.StatusOK, res)
}

// SaveBillingReceipt represents get save billing receipt endpoint
// @Summary Saves a billing receipt
// @Description Saves a billing receipt in the context of an application
// @ID save-billing-receipt
// @Accept json
// @Produce json
// @Param Authorization header string true "Access Token" default(Bearer <Add access token here>)
// @Param uuid query string true "Identifier of the application to generate the excel from"
// @Param short query string true "Short name of the teacher this should be generated for"
// @Param files body PDFs true "The files to save as an array of the base64 decoded file contents"
// @Success 200 {object} Information
// @Failure 401 {object} Error
// @Failure 422 {object} Error
// @Failure 500 {object} Error
// @Router /saveBillingReceipt [post]
func SaveBillingReceipt(con *gin.Context) {
	auth, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, Error{"you are not logged in"})
	}
	r := PDFs{}
	if err := con.ShouldBindJSON(&r); err != nil {
		con.JSON(http.StatusUnprocessableEntity, Error{"invalid request structure provided"})
		return
	}
	db := mongo.MongoDatabaseConnector{}
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, Error{"database didn't respond"})
		return
	}
	defer db.Close()
	_ = con.Request.ParseForm()
	query := con.Request.URL.Query()
	if _, hasUUID := con.Request.Form["uuid"]; !hasUUID {
		con.JSON(http.StatusUnprocessableEntity, Error{"invalid request structure provided"})
		return
	}
	uuid := query.Get("uuid")
	if _, hasShort := con.Request.Form["short"]; !hasShort {
		con.JSON(http.StatusUnprocessableEntity, Error{"invalid request structure provided"})
		return
	}
	short := query.Get("short")
	application := db.GetApplication(uuid)
	requestTeacher := db.GetTeacherByShort(auth.Username)
	var in bool
	if application.Kind == mongo.SchoolEvent {
		teachers := application.SchoolEventDetails.Teachers
		for _, t := range teachers {
			if t.Shortname == requestTeacher.Short {
				in = true
				break
			}
		}
	} else if application.Kind == mongo.Training {
		if application.TrainingDetails.Organizer == requestTeacher.Short {
			in = true
		}
	} else if application.Kind == mongo.OtherReason {
		if application.OtherReasonDetails.Filer == requestTeacher.Short {
			in = true
		}
	}
	if !in {
		con.JSON(http.StatusUnauthorized, Error{"you have no permission to do this"})
		return
	}
	path, err := files.GenerateFileEnvironment(application)
	if err != nil {
		con.JSON(http.StatusInternalServerError, Error{"couldn't create directories"})
		return
	}
	ff, err := ioutil.ReadDir(filepath.Join(path, files.UploadFolderName))
	if err != nil {
		con.JSON(http.StatusInternalServerError, Error{"couldn't read upload directory"})
		return
	}
	counter := 1
	for _, file := range ff {
		data := strings.Split(file.Name(), "_")
		if data[1] == short {
			counter++
		}
	}
	for i, pdf := range r.Files {
		name := fmt.Sprintf(files.ReceiptFileName, i+counter, short)
		dec, err := base64.StdEncoding.DecodeString(pdf.Content)
		if err != nil {
			con.JSON(http.StatusInternalServerError, Error{fmt.Sprintf("couldn't decode the pdf file: %v", name)})
			return
		}
		file, err := os.Create(filepath.Join(files.BasePath, files.UploadFolderName, name))
		if err != nil {
			con.JSON(http.StatusInternalServerError, Error{fmt.Sprintf("couldn't create the pdf file: %v", name)})
			return
		}
		if _, err := file.Write(dec); err != nil {
			con.JSON(http.StatusInternalServerError, Error{fmt.Sprintf("couldn't write the pdf file: %v", name)})
			return
		}
		if err := file.Sync(); err != nil {
			con.JSON(http.StatusInternalServerError, Error{fmt.Sprintf("couldn't sync the pdf file: %v", name)})
			return
		}
		_ = file.Close()
	}
	con.JSON(http.StatusOK, Information{"saving successful"})
}
