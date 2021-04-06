package rest

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	mongo "github.com/refundable-tgm/huginn/db"
	"github.com/refundable-tgm/huginn/ldap"
	"github.com/refundable-tgm/huginn/untis"
	"net/http"
	"sort"
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
	err = untis.GetClient(auth.Username).Close()
	if err != nil {
		con.JSON(http.StatusInternalServerError, "error logging out of untis")
		return
	}
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

func GetLongName(con *gin.Context) {
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
	resp := map[string]string{
		"short": teacher.Short,
		"long":  teacher.Longname,
		"uuid":  teacher.UUID,
	}
	con.JSON(http.StatusOK, resp)
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
	query := con.Request.URL.Query()
	applyFilter := query.Get("username") == ""
	filter := query.Get("username")
	requestTeacher := db.GetTeacherByShort(auth.Username)
	if !(requestTeacher.Administration || requestTeacher.AV || requestTeacher.PEK || (applyFilter && requestTeacher.Short == filter)) {
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
	query := con.Request.URL.Query()
	applyFilter := query.Get("username") == ""
	filter := query.Get("username")
	requestTeacher := db.GetTeacherByShort(auth.Username)
	if !(requestTeacher.Administration || requestTeacher.AV || requestTeacher.PEK || (applyFilter && requestTeacher.Short == filter)) {
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
	if !(in || requestTeacher.Administration || requestTeacher.AV || requestTeacher.PEK) {
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
	if !(teacher.PEK || teacher.Administration || teacher.AV) {
		con.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	applications := db.GetAllApplications()
	res := make([]mongo.Application, 0)
	for _, app := range applications {
		if app.Kind == mongo.SchoolEvent {
			if app.Progress == mongo.SEInProcess && (teacher.Administration || teacher.AV) {
				res = append(res, app)
			}
			if app.Progress == mongo.SECostsInProcess {
				res = append(res, app)
			}
		} else if app.Kind == mongo.Training || app.Kind == mongo.OtherReason {
			if app.Progress == mongo.TInProcess && (teacher.Administration || teacher.AV) {
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
		con.JSON(http.StatusOK, "error; application not created")
	}
}
