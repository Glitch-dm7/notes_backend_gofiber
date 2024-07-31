package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// user model
type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key"`
	Username string    `json:"username"`
	Password string    `json:"password"` // Hashed password
	Notes    []Note    `gorm:"foreignkey:UserID"`
}

// notes moddel
type Note struct {
	ID 				uuid.UUID `gorm:"type:uuid"`
	Title			string		`json:"title"`
	Subtitle 	string		`json:"subtitle"`
	Text 			string		`json:"text"`
	UserID    uuid.UUID	`gorm:"type:uuid"` // Foreign key
	gorm.Model
}