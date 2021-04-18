package rest

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"   // gin swagger middleware
	"github.com/swaggo/gin-swagger/swaggerFiles" // swagger files
	"log"
	"net/http"
	"strconv"
)

// Port is the port this api will listen to
const Port = 8080

// StartService starts the rest service
// @title Refundable
// @version 1.1
// @description This REST-API provides the backend of Refundable
// @contact.name Michael Beier - Entwickler
// @contact.url https://mbeier.at
// @contact.email admin@mbeier.at
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @license.name MIT
// @license.url https://github.com/refundable-tgm/huginn/blob/master/LICENSE
// @host localhost:8080
// @BasePath /
// @query.collection.format multi
func StartService() {
	InitTokenManager()
	router := gin.Default()
	//gin.SetMode(gin.ReleaseMode)
	router.POST("/login", Login)
	router.POST("/logout", AuthWall(), Logout)
	router.POST("/login/refresh", Refresh)
	router.GET("/getTeacherByShort", AuthWall(), GetTeacherByShort)
	router.GET("/getTeacher", AuthWall(), GetTeacher)
	router.GET("/setTeacherPermissions", AuthWall(), SetTeacherPermissions)
	router.GET("/getActiveApplications", AuthWall(), GetActiveApplications)
	router.GET("/getAllApplications", AuthWall(), GetAllApplications)
	router.GET("/getNews", AuthWall(), GetNews)
	router.GET("/getAdminApplications", AuthWall(), GetAdminApplications)
	router.GET("/getApplication", AuthWall(), GetApplication)
	router.POST("/createApplication", AuthWall(), CreateApplication)
	router.PUT("/updateApplication", AuthWall(), UpdateApplication)
	router.DELETE("/deleteApplication", AuthWall(), DeleteApplication)
	router.GET("/getAbsenceFormForClasses", AuthWall(), GetAbsenceFormForClasses)
	router.GET("/getAbsenceFormForTeacher", AuthWall(), GetAbsenceFormForTeacher)
	router.GET("/getCompensationForEducationalSupportForm", AuthWall(), GetCompensationForEducationalSupportForm)
	router.GET("/getTravelInvoiceForm", AuthWall(), GetTravelInvoiceForm)
	router.GET("/getBusinessTripApplicationForm", AuthWall(), GetBusinessTripApplicationForm)
	router.GET("/getTravelInvoiceExcel", AuthWall(), GetTravelInvoiceExcel)
	router.GET("/getBusinessTripApplicationExcel", AuthWall(), GetBusinessTripApplicationExcel)
	router.POST("/saveBillingReceipt", AuthWall(), SaveBillingReceipt)

	router.NoRoute(func(context *gin.Context) {
		context.JSON(http.StatusNotFound, Error{"this endpoint doesn't exist"})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/", func(context *gin.Context) {
		context.Header("Access-Control-Allow-Origin", "*")
		context.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	log.Fatal(router.Run(":" + strconv.Itoa(Port)))
}
