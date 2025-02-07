package models

import "time"

type Post struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Title           string    `gorm:"not null" json:"title"`
	Content         string    `gorm:"not null" json:"content"`
	AuthorID        uint      `gorm:"not null" json:"author_id"`
	DisableComments bool      `gorm:"default:false" json:"disable_comments"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
