package storage

import (
	"errors"
	"sync"
	"time"

	"github.com/Anabol1ks/ozon_tz/internal/models"
)

type MemoryStorage struct {
	users    map[uint]*models.User
	posts    map[uint]*models.Post
	comments map[uint]*models.Comment
	lastID   uint
	mu       sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		users:    make(map[uint]*models.User),
		posts:    make(map[uint]*models.Post),
		comments: make(map[uint]*models.Comment),
		lastID:   0,
	}
}

func (s *MemoryStorage) nextID() uint {
	s.lastID++
	return s.lastID
}

func (s *MemoryStorage) CreateUser(user *models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	user.ID = s.nextID()
	user.CreatedAt = time.Now()
	s.users[user.ID] = user
	return nil
}

func (s *MemoryStorage) GetUser(id uint) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if user, ok := s.users[id]; ok {
		return user, nil
	}
	return nil, errors.New("user not found")
}

func (s *MemoryStorage) CreatePost(post *models.Post) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	post.ID = s.nextID()
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	s.posts[post.ID] = post
	return nil
}

func (s *MemoryStorage) GetPost(id uint) (*models.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if post, ok := s.posts[id]; ok {
		return post, nil
	}
	return nil, errors.New("post not found")
}

func (s *MemoryStorage) GetPosts() ([]*models.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	posts := make([]*models.Post, 0, len(s.posts))
	for _, post := range s.posts {
		posts = append(posts, post)
	}
	return posts, nil
}

func (s *MemoryStorage) CreateComment(comment *models.Comment) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	comment.ID = s.nextID()
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	s.comments[comment.ID] = comment
	return nil
}

func (s *MemoryStorage) GetComment(id uint) (*models.Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if comment, ok := s.comments[id]; ok {
		return comment, nil
	}
	return nil, errors.New("comment not found")
}

func (s *MemoryStorage) GetComments(postID uint, limit, offset *int32) ([]*models.Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var comments []*models.Comment
	for _, comment := range s.comments {
		if comment.PostID == postID && comment.ParentID == nil {
			comments = append(comments, comment)
		}
	}

	if offset != nil {
		start := int(*offset)
		if start >= len(comments) {
			return []*models.Comment{}, nil
		}
		comments = comments[start:]
	}

	if limit != nil {
		end := int(*limit)
		if end > len(comments) {
			end = len(comments)
		}
		comments = comments[:end]
	}

	return comments, nil
}

func (s *MemoryStorage) GetCommentChildren(parentID uint) ([]*models.Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var children []*models.Comment
	for _, comment := range s.comments {
		if comment.ParentID != nil && *comment.ParentID == parentID {
			children = append(children, comment)
		}
	}
	return children, nil
}

func (s *MemoryStorage) UpdatePost(post *models.Post) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.posts[post.ID]; !ok {
		return errors.New("post not found")
	}
	post.UpdatedAt = time.Now()
	s.posts[post.ID] = post
	return nil
}
