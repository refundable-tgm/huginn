package rest

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

const Port = 8080

func StartService() {
	InitTokenManager()
	router := gin.Default()
	//gin.SetMode(gin.ReleaseMode)
	router.POST("/login", Login)
	router.POST("/logout", Logout)
	log.Fatal(router.Run(":" + strconv.Itoa(Port)))
}
