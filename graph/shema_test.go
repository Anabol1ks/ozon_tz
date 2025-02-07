package graph

import (
	"context"
	"fmt"
	"testing"

	"github.com/Anabol1ks/ozon_tz/internal/models"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestCreateUser(t *testing.T) {
	db := setupTestDB(t)
	resolver := &Resolver{DB: db}
	mutation := &mutationResolver{resolver}
	ctx := context.Background()

	user, err := mutation.CreateUser(ctx, "testuser")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Username)

	// Test empty username
	_, err = mutation.CreateUser(ctx, "")
	assert.Error(t, err)
}

func TestCreateAndGetPost(t *testing.T) {
	db := setupTestDB(t)
	resolver := &Resolver{DB: db}
	mutation := &mutationResolver{resolver}
	query := &queryResolver{resolver}
	ctx := context.Background()

	// Create user first
	user, _ := mutation.CreateUser(ctx, "testuser")

	// Create post
	post, err := mutation.CreatePost(ctx, "Test Title", "Test Content", user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, "Test Title", post.Title)

	// Get post
	fetchedPost, err := query.GetPost(ctx, post.ID)
	assert.NoError(t, err)
	assert.Equal(t, post.ID, fetchedPost.ID)
}

func TestCreateComment(t *testing.T) {
	db := setupTestDB(t)
	resolver := &Resolver{DB: db}
	mutation := &mutationResolver{resolver}
	ctx := context.Background()

	// Setup
	user, _ := mutation.CreateUser(ctx, "testuser")
	post, _ := mutation.CreatePost(ctx, "Test Post", "Content", user.ID)

	// Test comment creation
	comment, err := mutation.CreateComment(ctx, post.ID, nil, user.ID, "Test Comment")
	assert.NoError(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, "Test Comment", comment.Content)
}

func TestToggleComments(t *testing.T) {
	db := setupTestDB(t)
	resolver := &Resolver{DB: db}
	mutation := &mutationResolver{resolver}
	ctx := context.Background()

	// Setup
	user, _ := mutation.CreateUser(ctx, "testuser")
	post, _ := mutation.CreatePost(ctx, "Test Post", "Content", user.ID)

	// Test disable
	updatedPost, err := mutation.ToggleComments(ctx, post.ID, true)
	assert.NoError(t, err)
	assert.True(t, updatedPost.DisableComments)

	// Test enable
	updatedPost, err = mutation.ToggleComments(ctx, post.ID, false)
	assert.NoError(t, err)
	assert.False(t, updatedPost.DisableComments)
}

func TestGetPosts(t *testing.T) {
	db := setupTestDB(t)
	resolver := &Resolver{DB: db}
	mutation := &mutationResolver{resolver}
	query := &queryResolver{resolver}
	ctx := context.Background()

	// Setup
	user, _ := mutation.CreateUser(ctx, "testuser")
	_, _ = mutation.CreatePost(ctx, "Post 1", "Content 1", user.ID)
	_, _ = mutation.CreatePost(ctx, "Post 2", "Content 2", user.ID)

	// Test get all posts
	posts, err := query.GetPosts(ctx)
	assert.NoError(t, err)
	assert.Len(t, posts, 2)
}

func TestGetPostWithComments(t *testing.T) {
	db := setupTestDB(t)
	resolver := &Resolver{DB: db}
	mutation := &mutationResolver{resolver}
	query := &queryResolver{resolver}
	ctx := context.Background()

	// Setup test data
	user, _ := mutation.CreateUser(ctx, "testuser")
	post, _ := mutation.CreatePost(ctx, "Test Post", "Content", user.ID)

	// Create comments hierarchy
	comment1, _ := mutation.CreateComment(ctx, post.ID, nil, user.ID, "Parent comment")
	comment2, _ := mutation.CreateComment(ctx, post.ID, &comment1.ID, user.ID, "Child comment")

	// Test post retrieval
	fetchedPost, err := query.GetPost(ctx, post.ID)
	assert.NoError(t, err)
	assert.Equal(t, post.ID, fetchedPost.ID)

	// Test comments retrieval
	comments, err := query.GetComments(ctx, post.ID, nil, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, comments)

	// Verify parent comment
	assert.Equal(t, comment1.Content, comments[0].Content)

	// Verify child comment relationship
	children, err := (&commentResolver{resolver}).Children(ctx, comments[0])
	assert.NoError(t, err)
	assert.NotEmpty(t, children)
	assert.Equal(t, comment2.Content, children[0].Content)
}

func TestPaginatedComments(t *testing.T) {
	db := setupTestDB(t)
	resolver := &Resolver{DB: db}
	mutation := &mutationResolver{resolver}
	query := &queryResolver{resolver}
	ctx := context.Background()

	// Setup
	user, _ := mutation.CreateUser(ctx, "testuser")
	post, _ := mutation.CreatePost(ctx, "Test Post", "Content", user.ID)

	// Create multiple comments
	for i := 0; i < 15; i++ {
		_, _ = mutation.CreateComment(ctx, post.ID, nil, user.ID, fmt.Sprintf("Comment %d", i))
	}

	// Test pagination
	limit := int32(5)
	offset := int32(0)
	comments, err := query.GetComments(ctx, post.ID, &limit, &offset)
	assert.NoError(t, err)
	assert.Len(t, comments, 5)
}
