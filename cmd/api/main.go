package main

import (
	"aws-Api-Go/database"
	"aws-Api-Go/middleware"
	"aws-Api-Go/routes"
	"aws-Api-Go/utils"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// var users = mocks.Users

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	logger := utils.InitLogger()
	defer logger.Sync()

	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	client := database.GetClient()
	defer database.CloseConnection()

	r.Use(middleware.LogMiddleware(logger))

	// Setup routes
	routes.SetupUserRoutes(r, client, logger)
	routes.SetupPostRoutes(r, client, logger)

	// Start server
	if err := r.Run(":" + os.Getenv("PORT")); err != nil {
		log.Fatal(err)
	}

}
