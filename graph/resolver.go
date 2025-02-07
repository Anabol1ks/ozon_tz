package graph

import (
	"sync"

	"github.com/Anabol1ks/ozon_tz/graph/model"
	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB                *gorm.DB
	CommentObservers  map[string][]chan *model.Comment
	CommentObserversM sync.Mutex
}
