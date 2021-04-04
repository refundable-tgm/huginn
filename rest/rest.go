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
	router.POST("/logout", AuthWall(), Logout)
	router.POST("/login/refresh", Refresh)
	router.GET("/getLongName", AuthWall(), GetLongName)
	router.GET("/getTeacher", AuthWall(), GetTeacher)
	router.GET("/getActiveApplications", AuthWall(), GetActiveApplications)
	router.GET("/getAllApplications", AuthWall(), GetAllApplication)
	router.GET("/getNews", AuthWall(), GetNews)
	router.GET("/getApplication", AuthWall(), GetApplication)
	log.Fatal(router.Run(":" + strconv.Itoa(Port)))
}
