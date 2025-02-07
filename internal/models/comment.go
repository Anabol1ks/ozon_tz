package models

import (
	"time"
)

type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PostID    uint      `gorm:"not null" json:"post_id"`
	AuthorID  uint      `gorm:"not null" json:"author_id"`
	ParentID  *uint     `json:"parent_id"`
	Content   string    `gorm:"not null;size:2000" json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
