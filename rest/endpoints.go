package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/refundable-tgm/huginn/ldap"
	"net/http"
)

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
