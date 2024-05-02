package main

import (
	"fmt"
	customerrors "jwt_use/customErrors"
	"jwt_use/routes"
	"log"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	router := gin.Default()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	port := "6001"
	err := router.Run(":" + port)

	if err != nil {
		log.Println(customerrors.ErrServerNOtConnected)
	}
	err = godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		// Handle error
	}

	fmt.Print(">>>>>>>", os.Getenv("SECRETKEY"))
	// log.Printf("Successfully running on port :%v", port)
}
