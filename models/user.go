package models

type User struct {
	ID        int    `json:"id"`
	Firstname string `json:"firstname" binding:"required"`
	Lastname  string `json:"lastname" binding:"required"`
	Email     string `json:"email" binding:"required,email" gorm:"unique"`
}
