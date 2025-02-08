package graph

import (
	"context"
	"errors"
	"strconv"

	"github.com/Anabol1ks/ozon_tz/graph/model"
	"github.com/Anabol1ks/ozon_tz/internal/models"
)

type CommentResolver interface {
	Children(ctx context.Context, obj *model.Comment) ([]*model.Comment, error)
}

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, title string, content string, authorID string) (*model.Post, error) {
	authorIDUint, _ := strconv.ParseUint(authorID, 10, 64)
	post := &models.Post{Title: title, Content: content, AuthorID: uint(authorIDUint)}
	if err := r.Store.CreatePost(post); err != nil {
		return nil, err
	}
	return dbPostToGraphQL(post), nil
}

func (r *commentResolver) Children(ctx context.Context, obj *model.Comment) ([]*model.Comment, error) {
	commentID, _ := strconv.ParseUint(obj.ID, 10, 64)
	comments, err := r.Store.GetCommentChildren(uint(commentID))
	if err != nil {
		return nil, err
	}

	result := make([]*model.Comment, len(comments))
	for i, comment := range comments {
		result[i] = dbCommentToGraphQL(comment)
	}
	return result, nil
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, postID string, parentID *string, authorID string, content string) (*model.Comment, error) {
	postIDUint, _ := strconv.ParseUint(postID, 10, 64)
	authorIDUint, _ := strconv.ParseUint(authorID, 10, 64)

	post, err := r.Store.GetPost(uint(postIDUint))
	if err != nil {
		return nil, err
	}

	if post.DisableComments {
		return nil, errors.New("comments are disabled for this post")
	}

	comment := &models.Comment{
		PostID:   uint(postIDUint),
		AuthorID: uint(authorIDUint),
		Content:  content,
	}

	if parentID != nil {
		parentIDUint, _ := strconv.ParseUint(*parentID, 10, 64)
		// Verify parent comment exists
		_, err := r.Store.GetComment(uint(parentIDUint))
		if err != nil {
			return nil, errors.New("parent comment not found")
		}
		parentIDUintVal := uint(parentIDUint)
		comment.ParentID = &parentIDUintVal
	}

	if err := r.Store.CreateComment(comment); err != nil {
		return nil, err
	}

	r.CommentObserversM.Lock()
	if channels, ok := r.CommentObservers[postID]; ok {
		for _, ch := range channels {
			// Отправляем комментарий в канал (если канал открыт)
			select {
			case ch <- dbCommentToGraphQL(comment):
			default: // Если канал заблокирован (подписчик не слушает), пропускаем
			}
		}
	}
	r.CommentObserversM.Unlock()

	return dbCommentToGraphQL(comment), nil
}

// ToggleComments is the resolver for the toggleComments field.
func (r *mutationResolver) ToggleComments(ctx context.Context, postID string, disable bool) (*model.Post, error) {
	postIDUint, _ := strconv.ParseUint(postID, 10, 64)
	post, err := r.Store.GetPost(uint(postIDUint))
	if err != nil {
		return nil, err
	}

	post.DisableComments = disable
	if err := r.Store.UpdatePost(post); err != nil {
		return nil, err
	}
	return dbPostToGraphQL(post), nil
}

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, username string) (*model.User, error) {
	if username == "" {
		return nil, errors.New("username не может быть пустым")
	}

	dbUser := &models.User{Username: username}
	if err := r.Store.CreateUser(dbUser); err != nil {
		return nil, err
	}

	return &model.User{
		ID:        strconv.FormatUint(uint64(dbUser.ID), 10),
		Username:  dbUser.Username,
		CreatedAt: dbUser.CreatedAt.String(),
	}, nil
}

// GetPosts is the resolver for the getPosts field.
func (r *queryResolver) GetPosts(ctx context.Context) ([]*model.Post, error) {
	posts, err := r.Store.GetPosts()
	if err != nil {
		return nil, err
	}

	result := make([]*model.Post, len(posts))
	for i, post := range posts {
		result[i] = dbPostToGraphQL(post)
	}
	return result, nil
}

// GetPost is the resolver for the getPost field.
func (r *queryResolver) GetPost(ctx context.Context, id string) (*model.Post, error) {
	postID, _ := strconv.ParseUint(id, 10, 64)
	post, err := r.Store.GetPost(uint(postID))
	if err != nil {
		return nil, err
	}
	return dbPostToGraphQL(post), nil
}

// GetComments is the resolver for the getComments field.
func (r *queryResolver) GetComments(ctx context.Context, postID string, limit *int32, offset *int32) ([]*model.Comment, error) {
	postIDUint, err := strconv.ParseUint(postID, 10, 64)
	if err != nil {
		return nil, err
	}

	comments, err := r.Store.GetComments(uint(postIDUint), limit, offset)
	if err != nil {
		return nil, err
	}

	result := make([]*model.Comment, len(comments))
	for i, comment := range comments {
		result[i] = dbCommentToGraphQL(comment)
	}

	return result, nil
}

// OnNewComment is the resolver for the onNewComment field.
func (r *subscriptionResolver) OnNewComment(ctx context.Context, postID string) (<-chan *model.Comment, error) {
	commentChan := make(chan *model.Comment, 1)

	postKey := postID

	r.CommentObserversM.Lock()
	r.CommentObservers[postKey] = append(r.CommentObservers[postKey], commentChan)
	r.CommentObserversM.Unlock()

	go func() {
		<-ctx.Done()

		r.CommentObserversM.Lock()
		defer r.CommentObserversM.Unlock()

		// Фильтруем каналы, удаляя текущий
		channels := r.CommentObservers[postKey]
		newChannels := make([]chan *model.Comment, 0, len(channels))
		for _, ch := range channels {
			if ch != commentChan {
				newChannels = append(newChannels, ch)
			}
		}

		// Если больше нет подписчиков, удаляем ключ
		if len(newChannels) > 0 {
			r.CommentObservers[postKey] = newChannels
		} else {
			delete(r.CommentObservers, postKey)
		}

		close(commentChan) // Закрываем канал
	}()

	return commentChan, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

// Comment returns CommentResolver implementation.
func (r *Resolver) Comment() CommentResolver { return &commentResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
type commentResolver struct{ *Resolver }

func dbPostToGraphQL(dbPost *models.Post) *model.Post {
	return &model.Post{
		ID:              strconv.FormatUint(uint64(dbPost.ID), 10),
		Title:           dbPost.Title,
		Content:         dbPost.Content,
		DisableComments: dbPost.DisableComments,
		CreatedAt:       dbPost.CreatedAt.String(),
	}
}

func dbCommentToGraphQL(dbComment *models.Comment) *model.Comment {
	return &model.Comment{
		ID:        strconv.FormatUint(uint64(dbComment.ID), 10),
		Content:   dbComment.Content,
		CreatedAt: dbComment.CreatedAt.String(),
	}
}
