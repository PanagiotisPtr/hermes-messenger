package entityutils

import (
	"context"

	"github.com/panagiotisptr/hermes-messenger/libs/utils/entityutils/filter"
)

type Repository[T any] interface {
	Find(context.Context, filter.Filter) (<-chan *T, error)

	FindOne(context.Context, filter.Filter) (*T, error)

	Create(context.Context, *T) (*T, error)

	Update(context.Context, filter.Filter, *T) (int64, error)

	Delete(context.Context, filter.Filter) (int64, error)
}

type OperationType int64

const (
	CreateOp OperationType = iota
	UpdateOp
	DeleteOp
)
