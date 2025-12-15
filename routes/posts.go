package routes

import (
	"aws-Api-Go/handlers"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func SetupPostRoutes(r *gin.Engine, db *gorm.DB, logger *zap.Logger) {
	postHandler := handlers.NewPostHandler(db, logger)

	// Post routes
	r.GET("/posts", postHandler.GetPosts)
	r.GET("/posts/:id", postHandler.GetPostByID)
	r.GET("/posts/user/:userID", postHandler.GetPostsByUserID)
	r.POST("/posts", postHandler.CreatePost)
	r.PUT("/posts/:id", postHandler.UpdatePost)
	r.PATCH("/posts/:id", postHandler.PartialUpdatePost)
	r.DELETE("/posts/:id", postHandler.DeletePost)
	r.POST("/posts/:id/restore", postHandler.RestorePost)
}
