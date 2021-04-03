package rest

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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
