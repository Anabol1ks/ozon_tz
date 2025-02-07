package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"not null;unique" json:"username"`
	Password  string    `gorm:"not null" json:"password"`
	CreatedAt time.Time `json:"created_at"`
}
