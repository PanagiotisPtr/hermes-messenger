package entityutils

import (
	"context"

	"github.com/panagiotisptr/hermes-messenger/libs/utils/entityutils/filter"
)

type Repository[T any] interface {
	// Find finds all entities (as a stream/channel) that match a filter
	Find(context.Context, filter.Filter) (<-chan *T, error)

	// FindOne finds a single entity that matches a filter
	FindOne(context.Context, filter.Filter) (*T, error)

	// Create creates an entity based on the attributes passed
	// it then returns the newly created entity
	Create(context.Context, *T) (*T, error)

	// Update updates all entities that match a filter
	// and returns the number of entities that got updated
	Update(context.Context, filter.Filter, *T) (int64, error)

	// Delete deletes all entities that match a filter
	// and returns the number of entities that got deleted
	Delete(context.Context, filter.Filter) (int64, error)
}

type OperationType int64

const (
	CreateOp OperationType = iota
	UpdateOp
	DeleteOp
)
