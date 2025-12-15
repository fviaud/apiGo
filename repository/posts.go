package repository

import (
	"aws-Api-Go/models"

	"gorm.io/gorm"
)

type PostRepository struct {
	DB *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{
		DB: db,
	}
}

func (r *PostRepository) FindAll() ([]models.Post, error) {
	var posts []models.Post
	err := r.DB.Find(&posts).Error
	return posts, err
}

func (r *PostRepository) FindByID(id int) (*models.Post, error) {
	var post models.Post
	err := r.DB.First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepository) FindByUserID(userID int) ([]models.Post, error) {
	var posts []models.Post
	err := r.DB.Preload("User").Where("user_id = ?", userID).Find(&posts).Error
	return posts, err
}

func (r *PostRepository) Create(post *models.Post) error {
	return r.DB.Create(post).Error
}

func (r *PostRepository) Update(post *models.Post) error {
	return r.DB.Save(post).Error
}

func (r *PostRepository) PartialUpdate(id int, updates map[string]interface{}) error {
	return r.DB.Model(&models.Post{}).Where("id = ?", id).Updates(updates).Error
}

func (r *PostRepository) Delete(id int) error {
	return r.DB.Delete(&models.Post{}, id).Error
}

func (r *PostRepository) Restore(id int) error {
	return r.DB.Unscoped().Model(&models.Post{}).Where("id = ?", id).Updates(map[string]interface{}{"deleted_at": nil}).Error
}
