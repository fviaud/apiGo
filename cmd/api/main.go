package main

import (
	"aws-Api-Go/database"
	"aws-Api-Go/models"
	"aws-Api-Go/utils"
	"log"
	"net/http"
	"strconv"

	"aws-Api-Go/middleware"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
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

	r.GET("/users", func(c *gin.Context) {

		var users []models.User

		if err := client.Find(&users).Error; err != nil {
			logger.Error("Error fetching users from database",
				zap.Error(err),
			)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		if len(users) == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "No users found"})
			return
		} else {
			logger.Debug("Fetched users",
				zap.Int("user_count", len(users)),
			)
		}
		c.JSON(http.StatusOK, users)
	})

	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		userID, err := strconv.Atoi(id)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var user models.User

		if err := client.First(&user, userID).Error; err != nil {
			logger.Error("Error fetching user from database",
				zap.Error(err),
			)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		c.JSON(http.StatusOK, user)
	})

	// r.POST("/users", func(c *gin.Context) {
	// 	var newUser models.User
	// 	if err := c.ShouldBindJSON(&newUser); err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 		return
	// 	}
	// 	newUser.ID = len(users) + 1
	// 	users = append(users, newUser)
	// 	c.JSON(http.StatusCreated, newUser)
	// })

	if err := r.Run(":" + os.Getenv("PORT")); err != nil {
		log.Fatal(err)
	}
}
