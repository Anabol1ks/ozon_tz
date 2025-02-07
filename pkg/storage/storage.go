package storage

import (
	"github.com/Anabol1ks/ozon_tz/internal/models"
)

type Storage interface {
	CreateUser(*models.User) error
	GetUser(id uint) (*models.User, error)
	CreatePost(*models.Post) error
	GetPost(id uint) (*models.Post, error)
	GetPosts() ([]*models.Post, error)
	CreateComment(*models.Comment) error
	GetComment(id uint) (*models.Comment, error)
	GetComments(postID uint, limit, offset *int32) ([]*models.Comment, error)
	GetCommentChildren(parentID uint) ([]*models.Comment, error)
	UpdatePost(*models.Post) error
}
