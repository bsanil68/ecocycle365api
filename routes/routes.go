package routes

import (
	"ecocycleapis/controller"
	"ecocycleapis/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up the application routes.
func SetupRoutes(router *gin.Engine) {

	log.Println("Configuration loaded")

	// Define the route for your REST service

	router.Use(middleware.CORSMiddleware())
	//router.POST("/storeiotdata", ctrl.StoreDeviceMeasuerments)
	router.POST("/storeData", controller.StoreHandler)
	log.Println("Routes registered successfully")
}
