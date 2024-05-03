package main

import (
	customerrors "jwt_use/customErrors"
	"jwt_use/middleware"
	"jwt_use/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	routes.AuthRoutes(router)
	router.Use(middleware.Authenticate())

	port := "6001"
	err := router.Run(":" + port)

	if err != nil {
		log.Println(customerrors.ErrServerNOtConnected)
	}
}
