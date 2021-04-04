package rest

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	mongo "github.com/refundable-tgm/huginn/db"
	"github.com/refundable-tgm/huginn/ldap"
	"net/http"
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
	out := map[string]string{
		"access_token": token.AccessToken,
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
		if  err != nil {
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
	resp := map[string]string {
		"short": teacher.Short,
		"long": teacher.Longname,
		"uuid": teacher.UUID,
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
