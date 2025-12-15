package handlers

import (
	"aws-Api-Go/models"
	"aws-Api-Go/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PostHandler struct {
	Repo   *repository.PostRepository
	Logger *zap.Logger
}

func NewPostHandler(db *gorm.DB, logger *zap.Logger) *PostHandler {
	return &PostHandler{
		Repo:   repository.NewPostRepository(db),
		Logger: logger,
	}
}

func (h *PostHandler) GetPosts(c *gin.Context) {
	posts, err := h.Repo.FindAll()
	if err != nil {
		h.Logger.Error("Error fetching posts from database", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if len(posts) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No posts found"})
		return
	}

	h.Logger.Debug("Fetched posts", zap.Int("post_count", len(posts)))
	c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) GetPostByID(c *gin.Context) {
	id := c.Param("id")
	postID, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	post, err := h.Repo.FindByID(postID)
	if err != nil {
		h.Logger.Error("Error fetching post from database", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) GetPostsByUserID(c *gin.Context) {
	userID := c.Param("userID")
	uid, err := strconv.Atoi(userID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	posts, err := h.Repo.FindByUserID(uid)
	if err != nil {
		h.Logger.Error("Error fetching posts by user ID from database", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if len(posts) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No posts found for this user"})
		return
	}

	h.Logger.Debug("Fetched posts by user ID", zap.Int("user_id", uid), zap.Int("post_count", len(posts)))
	c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var newPost models.Post
	if err := c.ShouldBindJSON(&newPost); err != nil {
		h.Logger.Error("Error binding JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.Create(&newPost); err != nil {
		h.Logger.Error("Error creating post in database", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.Logger.Debug("Created new post", zap.Int("post_id", int(newPost.ID)))
	c.JSON(http.StatusCreated, newPost)
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	id := c.Param("id")
	postID, err := strconv.Atoi(id)
	if err != nil {
		h.Logger.Error("Invalid post ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var updatedPost models.Post
	if err := c.ShouldBindJSON(&updatedPost); err != nil {
		h.Logger.Error("Error binding JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedPost.ID = uint(postID)

	if err := h.Repo.Update(&updatedPost); err != nil {
		h.Logger.Error("Error updating post in database", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	h.Logger.Debug("Updated post", zap.Int("post_id", int(updatedPost.ID)))
	c.JSON(http.StatusOK, updatedPost)
}

func (h *PostHandler) PartialUpdatePost(c *gin.Context) {
	id := c.Param("id")
	postID, err := strconv.Atoi(id)
	if err != nil {
		h.Logger.Error("Invalid post ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var postUpdates map[string]interface{}
	if err := c.ShouldBindJSON(&postUpdates); err != nil {
		h.Logger.Error("Error binding JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.PartialUpdate(postID, postUpdates); err != nil {
		h.Logger.Error("Error partially updating post in database", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	h.Logger.Debug("Partially updated post", zap.Int("post_id", postID))

	updatedPost, err := h.Repo.FindByID(postID)
	if err != nil {
		h.Logger.Error("Error fetching updated post from database", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, updatedPost)
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	id := c.Param("id")
	postID, err := strconv.Atoi(id)
	if err != nil {
		h.Logger.Error("Invalid post ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	if err := h.Repo.Delete(postID); err != nil {
		h.Logger.Error("Error deleting post from database", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	h.Logger.Debug("Deleted post", zap.Int("post_id", postID))
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
}

func (h *PostHandler) RestorePost(c *gin.Context) {
	id := c.Param("id")
	postID, err := strconv.Atoi(id)
	if err != nil {
		h.Logger.Error("Invalid post ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	if err := h.Repo.Restore(postID); err != nil {
		h.Logger.Error("Error restoring post in database", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	h.Logger.Debug("Restored post", zap.Int("post_id", postID))
	c.JSON(http.StatusOK, gin.H{"message": "Post restored"})
}
