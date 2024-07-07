package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Note struct {
	Title			string			`json:"title"`
	ID 				uuid.UUID 	`gorm:"type:uuid"`
	Subtitle 	string			`json:"subtitle"`
	Text 			string			`json:"text"`
	gorm.Model
}