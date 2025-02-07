package storage

import (
	"github.com/Anabol1ks/ozon_tz/internal/models"
	"gorm.io/gorm"
)

type PostgresStorage struct {
	db *gorm.DB
}

func NewPostgresStorage(db *gorm.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func (s *PostgresStorage) CreateUser(user *models.User) error {
	return s.db.Create(user).Error
}

func (s *PostgresStorage) GetUser(id uint) (*models.User, error) {
	var user models.User
	err := s.db.First(&user, id).Error
	return &user, err
}

func (s *PostgresStorage) CreatePost(post *models.Post) error {
	return s.db.Create(post).Error
}

func (s *PostgresStorage) GetPost(id uint) (*models.Post, error) {
	var post models.Post
	err := s.db.First(&post, id).Error
	return &post, err
}

func (s *PostgresStorage) GetPosts() ([]*models.Post, error) {
	var posts []*models.Post
	err := s.db.Find(&posts).Error
	return posts, err
}

func (s *PostgresStorage) CreateComment(comment *models.Comment) error {
	return s.db.Create(comment).Error
}

func (s *PostgresStorage) GetComment(id uint) (*models.Comment, error) {
	var comment models.Comment
	err := s.db.First(&comment, id).Error
	return &comment, err
}

func (s *PostgresStorage) GetComments(postID uint, limit, offset *int32) ([]*models.Comment, error) {
	var comments []*models.Comment
	query := s.db.Where("post_id = ? AND parent_id IS NULL", postID)
	if limit != nil {
		query = query.Limit(int(*limit))
	}
	if offset != nil {
		query = query.Offset(int(*offset))
	}
	err := query.Find(&comments).Error
	return comments, err
}

func (s *PostgresStorage) GetCommentChildren(parentID uint) ([]*models.Comment, error) {
	var comments []*models.Comment
	err := s.db.Where("parent_id = ?", parentID).Find(&comments).Error
	return comments, err
}

func (s *PostgresStorage) UpdatePost(post *models.Post) error {
	return s.db.Save(post).Error
}
