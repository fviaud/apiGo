package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Firstname string `json:"firstname" binding:"required"`
	Lastname  string `json:"lastname" binding:"required"`
	Email     string `json:"email" binding:"required,email" gorm:"unique"`
	Posts     []Post `json:"posts,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
