package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Password string `gorm:"not null" json:"-"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
}

type Template struct {
	gorm.Model
	Title       string `gorm:"not null" json:"title"`
	Description string `json:"description"`
	Content     string `gorm:"type:text;not null" json:"content"` // YAML content
	AuthorID    uint   `json:"author_id"`
	Author      User   `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Downloads   int    `gorm:"default:0" json:"downloads"`
	Status      string `gorm:"default:published" json:"status"` // published, draft, archived
}
