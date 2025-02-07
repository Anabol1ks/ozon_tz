package graph

import (
	"sync"

	"github.com/Anabol1ks/ozon_tz/graph/model"
	"github.com/Anabol1ks/ozon_tz/pkg/storage"
	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB                *gorm.DB
	Store             storage.Storage
	CommentObservers  map[string][]chan *model.Comment
	CommentObserversM sync.Mutex
}
