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

func AuthWall() gin.HandlerFunc {
	return func(con *gin.Context) {
		ok, err := TokenValid(con.Request)
		if !ok {
			con.JSON(http.StatusUnauthorized, err.Error())
			con.Abort()
			return
		}
		con.Next()
	}
}

func Login(con *gin.Context) {
	u := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	if err := con.ShouldBindJSON(&u); err != nil {
		con.JSON(http.StatusUnprocessableEntity, "invalid request structure provided")
		return
	}
	if !ldap.AuthenticateUserCredentials(u.Username, u.Password) {
		con.JSON(http.StatusUnauthorized, "this credentials do not resolve into an authorized login")
		return
	}
	token, err := CreateToken(u.Username)
	if err != nil {
		con.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	SaveToken(u.Username, token)
	untis.CreateClient(u.Username, u.Password)
	out := map[string]string{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	}
	con.JSON(http.StatusOK, out)
}

func Logout(con *gin.Context) {
	auth, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, "you are not logged in")
		return
	}
	DeleteToken(auth.AccessUUID)
	untis.GetClient(auth.Username).DeleteClient()
	con.JSON(http.StatusOK, "logged out")
}

func Refresh(con *gin.Context) {
	body := map[string]string{}
	if err := con.ShouldBindJSON(&body); err != nil {
		con.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	refresh := body["refresh_token"]
	token, err := jwt.Parse(refresh, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(refreshSecret), nil
	})

	if err != nil {
		con.JSON(http.StatusUnauthorized, "token expired")
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		con.JSON(http.StatusUnauthorized, err)
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uuid, ok := claims["refresh_uuid"].(string)
		if !ok {
			con.JSON(http.StatusUnprocessableEntity, err)
			return
		}
		username, ok := claims["username"].(string)
		if !ok {
			con.JSON(http.StatusUnprocessableEntity, err)
			return
		}
		DeleteToken(uuid)
		tok, err := CreateToken(username)
		if err != nil {
			con.JSON(http.StatusForbidden, err.Error())
			return
		}
		SaveToken(username, tok)
		tokens := map[string]string{
			"access_token":  tok.AccessToken,
			"refresh_token": tok.RefreshToken,
		}
		con.JSON(http.StatusCreated, tokens)
	} else {
		con.JSON(http.StatusUnauthorized, "refresh token expired")
	}
}

func GetTeacherByShort(con *gin.Context) {
	_, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, "you are not logged in")
		return
	}
	query := con.Request.URL.Query()
	if query.Get("name") == "" {
		con.JSON(http.StatusUnprocessableEntity, "invalid request structure provided")
		return
	}
	name := query.Get("name")
	db := mongo.MongoDatabaseConnector{}
	defer db.Close()
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, "database didn't respond")
		return
	}
	teacher := db.GetTeacherByShort(name)
	con.JSON(http.StatusOK, teacher)
}

func GetTeacher(con *gin.Context) {
	_, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, "you are not logged in")
		return
	}
	query := con.Request.URL.Query()
	if query.Get("uuid") == "" {
		con.JSON(http.StatusUnprocessableEntity, "invalid request structure provided")
		return
	}
	uuid := query.Get("uuid")
	db := mongo.MongoDatabaseConnector{}
	defer db.Close()
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, "database didn't respond")
		return
	}
	teacher := db.GetTeacherByUUID(uuid)
	con.JSON(http.StatusOK, teacher)
}

func SetTeacherPermissions(con *gin.Context) {
	auth, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, "you are not logged in")
		return
	}
	perm := struct {
		SuperUser bool `json:"super_user"`
		Administration bool `json:"administration"`
		AV bool `json:"av"`
		PEK bool `json:"pek"`
	}{}
	if err := con.ShouldBindJSON(&perm); err != nil {
		con.JSON(http.StatusUnprocessableEntity, "invalid request structure provided")
		return
	}
	query := con.Request.URL.Query()
	if query.Get("uuid") == "" {
		con.JSON(http.StatusUnprocessableEntity, "invalid request structure provided")
		return
	}
	uuid := query.Get("uuid")
	db := mongo.MongoDatabaseConnector{}
	defer db.Close()
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, "database didn't respond")
		return
	}
	requester := db.GetTeacherByShort(auth.Username)
	if !(requester.PEK || requester.Administration || requester.AV || requester.SuperUser) {
		con.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	teacher := db.GetTeacherByUUID(uuid)
	teacher.SuperUser = perm.SuperUser
	teacher.Administration = perm.Administration
	teacher.PEK = perm.PEK
	teacher.Administration = perm.Administration
	if db.UpdateTeacher(uuid, teacher) {
		con.JSON(http.StatusOK, "permissions updated")
	} else {
		con.JSON(http.StatusInternalServerError, "permissions couldn't be updated")
	}
}

func GetActiveApplications(con *gin.Context) {
	auth, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, "you are not logged in")
		return
	}
	db := mongo.MongoDatabaseConnector{}
	defer db.Close()
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, "database didn't respond")
		return
	}
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

func GetAllApplication(con *gin.Context) {
	auth, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, "you are not logged in")
		return
	}
	db := mongo.MongoDatabaseConnector{}
	defer db.Close()
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, "database didn't respond")
		return
	}
	_ = con.Request.ParseForm()
	query := con.Request.URL.Query()
	_, applyFilter := con.Request.Form["username"]
	filter := query.Get("username")
	requestTeacher := db.GetTeacherByShort(auth.Username)
	if !(requestTeacher.Administration || requestTeacher.AV || requestTeacher.SuperUser || requestTeacher.PEK || (applyFilter && requestTeacher.Short == filter)) {
		con.JSON(http.StatusUnauthorized, "unauthorized")
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

func GetNews(con *gin.Context) {
	auth, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	db := mongo.MongoDatabaseConnector{}
	defer db.Close()
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, "database didn't respond")
		return
	}
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
	news := make([]struct {
		UUID  string `json:"uuid"`
		Title string `json:"title"`
		State int    `json:"state"`
	}, 0)
	for _, app := range res {
		news = append(news, struct {
			UUID  string `json:"uuid"`
			Title string `json:"title"`
			State int    `json:"state"`
		}{app.UUID, app.Name, app.Progress})
	}
	con.JSON(http.StatusOK, news)
}

func GetApplication(con *gin.Context) {
	auth, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, "you are not logged in")
		return
	}
	db := mongo.MongoDatabaseConnector{}
	defer db.Close()
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, "database didn't respond")
		return
	}
	query := con.Request.URL.Query()
	uuid := query.Get("uuid")
	if query.Get("uuid") == "" {
		con.JSON(http.StatusUnprocessableEntity, "invalid request structure provided")
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
		con.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	con.JSON(http.StatusOK, application)
}

func GetAdminApplication(con *gin.Context) {
	auth, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, "you are not logged in")
		return
	}
	db := mongo.MongoDatabaseConnector{}
	defer db.Close()
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, "database didn't respond")
		return
	}
	teacher := db.GetTeacherByShort(auth.Username)
	if !(teacher.PEK || teacher.Administration || teacher.AV || teacher.SuperUser) {
		con.JSON(http.StatusUnauthorized, "unauthorized")
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

func CreateApplication(con *gin.Context) {
	app := mongo.Application{}
	if err := con.ShouldBindJSON(&app); err != nil {
		con.JSON(http.StatusUnprocessableEntity, "invalid request structure provided")
		return
	}
	_, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, "you are not logged in")
		return
	}
	db := mongo.MongoDatabaseConnector{}
	defer db.Close()
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, "database didn't respond")
		return
	}
	if db.CreateApplication(app) {
		con.JSON(http.StatusOK, "success; application created")
	} else {
		con.JSON(http.StatusInternalServerError, "error; application not created")
	}
}

func UpdateApplication(con *gin.Context) {
	app := mongo.Application{}
	if err := con.ShouldBindJSON(&app); err != nil {
		con.JSON(http.StatusUnprocessableEntity, "invalid request structure provided")
		return
	}
	auth, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, "you are not logged in")
		return
	}
	db := mongo.MongoDatabaseConnector{}
	defer db.Close()
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, "database didn't respond")
		return
	}
	query := con.Request.URL.Query()
	uuid := query.Get("uuid")
	if query.Get("uuid") == "" {
		con.JSON(http.StatusUnprocessableEntity, "invalid request structure provided")
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
		con.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	if db.UpdateApplication(uuid, app) {
		con.JSON(http.StatusOK, "success; application updated")
	} else {
		con.JSON(http.StatusInternalServerError, "error; application not updated")
	}
}

func DeleteApplication(con *gin.Context) {
	auth, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, "you are not logged in")
		return
	}
	db := mongo.MongoDatabaseConnector{}
	defer db.Close()
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, "database didn't respond")
		return
	}
	query := con.Request.URL.Query()
	uuid := query.Get("uuid")
	if query.Get("uuid") == "" {
		con.JSON(http.StatusUnprocessableEntity, "invalid request structure provided")
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
		con.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	if db.DeleteApplication(uuid) {
		con.JSON(http.StatusOK, "success; application deleted")
	} else {
		con.JSON(http.StatusInternalServerError, "error; application not deleted")
	}
}

func GetAbsenceFormForClasses(con *gin.Context) {
	auth, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, "you are not logged in")
	}
	db := mongo.MongoDatabaseConnector{}
	defer db.Close()
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, "database didn't respond")
		return
	}
	_ = con.Request.ParseForm()
	query := con.Request.URL.Query()
	if _, hasUUID := con.Request.Form["uuid"]; !hasUUID {
		con.JSON(http.StatusUnprocessableEntity, "invalid request structure provided")
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
		con.JSON(http.StatusUnauthorized, "you have no permission to do this")
		return
	}
	path, err := files.GenerateFileEnvironment(application)
	if err != nil {
		con.JSON(http.StatusInternalServerError, "couldn't create directories")
		return
	}
	paths, err := files.GenerateAbsenceFormForClass(path, auth.Username, application)
	if err != nil {
		con.JSON(http.StatusInternalServerError, "couldn't create pdfs")
		return
	}

	pdfs := make(map[string]string)
	for _, p := range paths {
		splits := strings.Split(p, string(filepath.Separator))
		filename := strings.Split(splits[len(splits) - 1], ".")[0]
		names := strings.Split(filename, "_")
		class := names[len(names) - 1]
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
		con.JSON(http.StatusInternalServerError, "couldn't save merged pdf")
		return
	}
	file, err := ioutil.ReadFile(created)
	if err != nil {
		con.JSON(http.StatusInternalServerError, "couldn't read merged pdf")
		return
	}
	enc := base64.StdEncoding.EncodeToString(file)
	res := map[string]string{
		"pdf": enc,
	}
	err = os.Remove(created)
	if err != nil {
		con.JSON(http.StatusInternalServerError, "couldn't delete merged pdf")
		return
	}
	con.JSON(http.StatusOK, res)
}

func GetAbsenceFormForTeacher(con *gin.Context) {
	auth, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, "you are not logged in")
	}
	db := mongo.MongoDatabaseConnector{}
	defer db.Close()
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, "database didn't respond")
		return
	}
	_ = con.Request.ParseForm()
	query := con.Request.URL.Query()
	if _, hasUUID := con.Request.Form["uuid"]; !hasUUID {
		con.JSON(http.StatusUnprocessableEntity, "invalid request structure provided")
		return
	}
	uuid := query.Get("uuid")
	if _, hasTeacher := con.Request.Form["teacher"]; !hasTeacher {
		con.JSON(http.StatusUnprocessableEntity, "invalid request structure provided")
		return
	}
	teacher := query.Get("teacher")
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
		con.JSON(http.StatusUnauthorized, "you have no permission to do this")
		return
	}
	path, err := files.GenerateFileEnvironment(application)
	if err != nil {
		con.JSON(http.StatusInternalServerError, "couldn't create directories")
		return
	}
	path, err = files.GenerateAbsenceFormForTeacher(path, auth.Username, teacher, application)
	if err != nil {
		con.JSON(http.StatusInternalServerError, "couldn't create pdfs")
		return
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		con.JSON(http.StatusInternalServerError, "couldn't read generated pdf")
		return
	}
	enc := base64.StdEncoding.EncodeToString(file)
	res := map[string]string{
		"pdf": enc,
	}
	con.JSON(http.StatusOK, res)
}

func GetCompensationForEducationalSupportForm(con *gin.Context) {
	auth, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, "you are not logged in")
	}
	db := mongo.MongoDatabaseConnector{}
	defer db.Close()
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, "database didn't respond")
		return
	}
	_ = con.Request.ParseForm()
	query := con.Request.URL.Query()
	if _, hasUUID := con.Request.Form["uuid"]; !hasUUID {
		con.JSON(http.StatusUnprocessableEntity, "invalid request structure provided")
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
		con.JSON(http.StatusUnauthorized, "you have no permission to do this")
		return
	}
	path, err := files.GenerateFileEnvironment(application)
	if err != nil {
		con.JSON(http.StatusInternalServerError, "couldn't create directories")
		return
	}
	path, err = files.GenerateCompensationForEducationalSupport(path, application)
	if err != nil {
		con.JSON(http.StatusInternalServerError, "couldn't create pdfs")
		return
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		con.JSON(http.StatusInternalServerError, "couldn't read generated pdf")
		return
	}
	enc := base64.StdEncoding.EncodeToString(file)
	res := map[string]string{
		"pdf": enc,
	}
	con.JSON(http.StatusOK, res)
}

func GetTravelInvoiceForm(con *gin.Context) {
	auth, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, "you are not logged in")
	}
	db := mongo.MongoDatabaseConnector{}
	defer db.Close()
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, "database didn't respond")
		return
	}
	_ = con.Request.ParseForm()
	query := con.Request.URL.Query()
	if _, hasUUID := con.Request.Form["uuid"]; !hasUUID {
		con.JSON(http.StatusUnprocessableEntity, "invalid request structure provided")
		return
	}
	uuid := query.Get("uuid")
	if _, hasShort := con.Request.Form["short"]; !hasShort {
		con.JSON(http.StatusUnprocessableEntity, "invalid request structure provided")
		return
	}
	short := query.Get("short")
	if _, hasTIID := con.Request.Form["ti_id"]; !hasTIID {
		con.JSON(http.StatusUnprocessableEntity, "invalid request structure provided")
		return
	}
	tiID, err := strconv.Atoi(query.Get("ti_id"))
	if err != nil {
		con.JSON(http.StatusUnprocessableEntity, "invalid ti_id provided")
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
		con.JSON(http.StatusUnauthorized, "you have no permission to do this")
		return
	}
	path, err := files.GenerateFileEnvironment(application)
	if err != nil {
		con.JSON(http.StatusInternalServerError, "couldn't create directories")
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
		con.JSON(http.StatusInternalServerError, "couldn't create pdfs")
		return
	}

	if applyMergeReceipts {
		pp := append(make([]string, 0), path)
		uploadFolder := filepath.Join(filepath.Dir(path), files.UploadFolderName)
		ff, err := ioutil.ReadDir(uploadFolder)
		if err != nil {
			con.JSON(http.StatusInternalServerError, "couldn't read upload directory")
			return
		}
		for _, file := range ff {
			data := strings.Split(file.Name(), "_")
			if data[1] == short {
				pp = append(pp, filepath.Join(uploadFolder, file.Name()))
			}
		}
		created := filepath.Join(filepath.Dir(path), fmt.Sprintf(files.TravelInvoicePDFFileName, short + "_merge"))
		err = api.MergeCreateFile(pp, created, pdfcpu.NewDefaultConfiguration())
		if err != nil {
			con.JSON(http.StatusInternalServerError, "couldn't save merged pdf")
			return
		}
		file, err := ioutil.ReadFile(created)
		if err != nil {
			con.JSON(http.StatusInternalServerError, "couldn't read generated pdf")
			return
		}
		enc := base64.StdEncoding.EncodeToString(file)
		res := map[string]string{
			"pdf": enc,
		}
		err = os.Remove(created)
		if err != nil {
			con.JSON(http.StatusInternalServerError, "couldn't delete merged pdf")
			return
		}
		con.JSON(http.StatusOK, res)
	} else {
		file, err := ioutil.ReadFile(path)
		if err != nil {
			con.JSON(http.StatusInternalServerError, "couldn't read generated pdf")
			return
		}
		enc := base64.StdEncoding.EncodeToString(file)
		res := map[string]string{
			"pdf": enc,
		}
		con.JSON(http.StatusOK, res)
	}
}

func GetBusinessTripApplicationForm(con *gin.Context) {
	auth, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, "you are not logged in")
	}
	db := mongo.MongoDatabaseConnector{}
	defer db.Close()
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, "database didn't respond")
		return
	}
	_ = con.Request.ParseForm()
	query := con.Request.URL.Query()
	if _, hasUUID := con.Request.Form["uuid"]; !hasUUID {
		con.JSON(http.StatusUnprocessableEntity, "invalid request structure provided")
		return
	}
	uuid := query.Get("uuid")
	if _, hasShort := con.Request.Form["short"]; !hasShort {
		con.JSON(http.StatusUnprocessableEntity, "invalid request structure provided")
		return
	}
	short := query.Get("short")
	if _, hasBTAID := con.Request.Form["bta_id"]; !hasBTAID {
		con.JSON(http.StatusUnprocessableEntity, "invalid request structure provided")
		return
	}
	btaID, err := strconv.Atoi(query.Get("bta_id"))
	if err != nil {
		con.JSON(http.StatusUnprocessableEntity, "invalid bta_id provided")
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
		con.JSON(http.StatusUnauthorized, "you have no permission to do this")
		return
	}
	path, err := files.GenerateFileEnvironment(application)
	if err != nil {
		con.JSON(http.StatusInternalServerError, "couldn't create directories")
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
		con.JSON(http.StatusInternalServerError, "couldn't create pdf")
		return
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		con.JSON(http.StatusInternalServerError, "couldn't read generated pdf")
		return
	}
	enc := base64.StdEncoding.EncodeToString(file)
	res := map[string]string{
		"pdf": enc,
	}
	con.JSON(http.StatusOK, res)
}

func GetTravelInvoiceExcel(con *gin.Context) {
	auth, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, "you are not logged in")
	}
	db := mongo.MongoDatabaseConnector{}
	defer db.Close()
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, "database didn't respond")
		return
	}
	_ = con.Request.ParseForm()
	query := con.Request.URL.Query()
	if _, hasUUID := con.Request.Form["uuid"]; !hasUUID {
		con.JSON(http.StatusUnprocessableEntity, "invalid request structure provided")
		return
	}
	uuid := query.Get("uuid")
	if _, hasShort := con.Request.Form["short"]; !hasShort {
		con.JSON(http.StatusUnprocessableEntity, "invalid request structure provided")
		return
	}
	short := query.Get("short")
	if _, hasTIID := con.Request.Form["ti_id"]; !hasTIID {
		con.JSON(http.StatusUnprocessableEntity, "invalid request structure provided")
		return
	}
	tiID, err := strconv.Atoi(query.Get("ti_id"))
	if err != nil {
		con.JSON(http.StatusUnprocessableEntity, "invalid ti_id provided")
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
		con.JSON(http.StatusUnauthorized, "you have no permission to do this")
		return
	}
	path, err := files.GenerateFileEnvironment(application)
	if err != nil {
		con.JSON(http.StatusInternalServerError, "couldn't create directories")
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
		con.JSON(http.StatusInternalServerError, "couldn't create excel")
		return
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		con.JSON(http.StatusInternalServerError, "couldn't read generated excel")
		return
	}
	enc := base64.StdEncoding.EncodeToString(file)
	res := map[string]string{
		"excel": enc,
	}
	con.JSON(http.StatusOK, res)
}

func GetBusinessTripApplicationExcel(con *gin.Context) {
	auth, err := ExtractTokenMeta(con.Request)
	if err != nil {
		con.JSON(http.StatusUnauthorized, "you are not logged in")
	}
	db := mongo.MongoDatabaseConnector{}
	defer db.Close()
	if !db.Connect() {
		con.JSON(http.StatusInternalServerError, "database didn't respond")
		return
	}
	_ = con.Request.ParseForm()
	query := con.Request.URL.Query()
	if _, hasUUID := con.Request.Form["uuid"]; !hasUUID {
		con.JSON(http.StatusUnprocessableEntity, "invalid request structure provided")
		return
	}
	uuid := query.Get("uuid")
	if _, hasShort := con.Request.Form["short"]; !hasShort {
		con.JSON(http.StatusUnprocessableEntity, "invalid request structure provided")
		return
	}
	short := query.Get("short")
	if _, hasBTAID := con.Request.Form["bta_id"]; !hasBTAID {
		con.JSON(http.StatusUnprocessableEntity, "invalid request structure provided")
		return
	}
	btaID, err := strconv.Atoi(query.Get("bta_id"))
	if err != nil {
		con.JSON(http.StatusUnprocessableEntity, "invalid bta_id provided")
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
		con.JSON(http.StatusUnauthorized, "you have no permission to do this")
		return
	}
	path, err := files.GenerateFileEnvironment(application)
	if err != nil {
		con.JSON(http.StatusInternalServerError, "couldn't create directories")
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
		con.JSON(http.StatusInternalServerError, "couldn't create excel")
		return
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		con.JSON(http.StatusInternalServerError, "couldn't read generated excel")
		return
	}
	enc := base64.StdEncoding.EncodeToString(file)
	res := map[string]string{
		"excel": enc,
	}
	con.JSON(http.StatusOK, res)
}
