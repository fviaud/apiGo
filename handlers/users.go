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

type UserHandler struct {
	Repo   *repository.UserRepository
	Logger *zap.Logger
}

func NewUserHandler(db *gorm.DB, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		Repo:   repository.NewUserRepository(db),
		Logger: logger,
	}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.Repo.FindAll()
	if err != nil {
		h.Logger.Error("Error fetching users from database", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No users found"})
		return
	}

	h.Logger.Debug("Fetched users", zap.Int("user_count", len(users)))
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.Repo.FindByID(userID)
	if err != nil {
		h.Logger.Error("Error fetching user from database", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		h.Logger.Error("Error binding JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.Create(&newUser); err != nil {
		h.Logger.Error("Error creating user in database", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.Logger.Debug("Created new user", zap.Int("user_id", int(newUser.ID)))
	c.JSON(http.StatusCreated, newUser)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		h.Logger.Error("Invalid user ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var updatedUser models.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		h.Logger.Error("Error binding JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser.ID = uint(userID)

	if err := h.Repo.Update(&updatedUser); err != nil {
		h.Logger.Error("Error updating user in database", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	h.Logger.Debug("Updated user", zap.Int("user_id", int(updatedUser.ID)))
	c.JSON(http.StatusOK, updatedUser)
}

func (h *UserHandler) PartialUpdateUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		h.Logger.Error("Invalid user ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var userUpdates map[string]interface{}
	if err := c.ShouldBindJSON(&userUpdates); err != nil {
		h.Logger.Error("Error binding JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.PartialUpdate(userID, userUpdates); err != nil {
		h.Logger.Error("Error partially updating user in database", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	h.Logger.Debug("Partially updated user", zap.Int("user_id", userID))

	updatedUser, err := h.Repo.FindByID(userID)
	if err != nil {
		h.Logger.Error("Error fetching updated user from database", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		h.Logger.Error("Invalid user ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.Repo.Delete(userID); err != nil {
		h.Logger.Error("Error deleting user from database", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	h.Logger.Debug("Deleted user", zap.Int("user_id", userID))
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

func (h *UserHandler) RestoreUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		h.Logger.Error("Invalid user ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.Repo.Restore(userID); err != nil {
		h.Logger.Error("Error restoring user in database", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	h.Logger.Debug("Restored user", zap.Int("user_id", userID))
	c.JSON(http.StatusOK, gin.H{"message": "User restored"})
}
