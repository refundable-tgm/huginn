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
// @BasePath /api
// @query.collection.format multi
func StartService() {
	InitTokenManager()
	router := gin.Default()
	//gin.SetMode(gin.ReleaseMode)
	api := router.Group("/api")
	{
		api.POST("/login", Login)
		api.POST("/logout", AuthWall(), Logout)
		api.POST("/login/refresh", Refresh)
		api.GET("/getTeacherByShort", AuthWall(), GetTeacherByShort)
		api.GET("/getTeacher", AuthWall(), GetTeacher)
		api.GET("/setTeacherPermissions", AuthWall(), SetTeacherPermissions)
		api.GET("/getActiveApplications", AuthWall(), GetActiveApplications)
		api.GET("/getAllApplications", AuthWall(), GetAllApplications)
		api.GET("/getNews", AuthWall(), GetNews)
		api.GET("/getAdminApplications", AuthWall(), GetAdminApplications)
		api.GET("/getApplication", AuthWall(), GetApplication)
		api.POST("/createApplication", AuthWall(), CreateApplication)
		api.PUT("/updateApplication", AuthWall(), UpdateApplication)
		api.DELETE("/deleteApplication", AuthWall(), DeleteApplication)
		api.GET("/getAbsenceFormForClasses", AuthWall(), GetAbsenceFormForClasses)
		api.GET("/getAbsenceFormForTeacher", AuthWall(), GetAbsenceFormForTeacher)
		api.GET("/getCompensationForEducationalSupportForm", AuthWall(), GetCompensationForEducationalSupportForm)
		api.GET("/getTravelInvoiceForm", AuthWall(), GetTravelInvoiceForm)
		api.GET("/getBusinessTripApplicationForm", AuthWall(), GetBusinessTripApplicationForm)
		api.GET("/getTravelInvoiceExcel", AuthWall(), GetTravelInvoiceExcel)
		api.GET("/getBusinessTripApplicationExcel", AuthWall(), GetBusinessTripApplicationExcel)
		api.POST("/saveBillingReceipt", AuthWall(), SaveBillingReceipt)
	}

	router.NoRoute(func(context *gin.Context) {
		context.JSON(http.StatusNotFound, Error{"this endpoint doesn't exist"})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Fatal(router.Run(":" + strconv.Itoa(Port)))
}
