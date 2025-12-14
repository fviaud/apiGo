package routes

import (
	"aws-Api-Go/handlers"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func SetupUserRoutes(r *gin.Engine, db *gorm.DB, logger *zap.Logger) {
	userHandler := handlers.NewUserHandler(db, logger)

	// User routes
	r.GET("/users", userHandler.GetUsers)
	r.GET("/users/:id", userHandler.GetUserByID)
	r.POST("/users", userHandler.CreateUser)
	r.PUT("/users/:id", userHandler.UpdateUser)
	r.PATCH("/users/:id", userHandler.PartialUpdateUser)
	r.DELETE("/users/:id", userHandler.DeleteUser)
}
