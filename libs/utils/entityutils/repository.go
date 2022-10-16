package entityutils

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/entityutils/filter"
)

type Repository[T any] interface {
	Find(context.Context, filter.Filter) (<-chan *T, error)

	FindOne(context.Context, filter.Filter) (*T, error)

	Create(context.Context, *T) (*T, error)

	Update(context.Context, filter.Filter, *T) (int64, error)

	Delete(context.Context, filter.Filter) (int64, error)
}

type RepoHelper struct{}

type OperationType int64

const (
	CreateOp OperationType = iota
	UpdateOp
	DeleteOp
)

func (rh *RepoHelper) UpdateMeta(
	ctx context.Context,
	m *Meta,
	op OperationType,
) error {
	userId, ok := ctx.Value("userId").(uuid.UUID)
	if !ok {
		return fmt.Errorf("missing userId from context")
	}

	now := time.Now()
	switch op {
	case CreateOp:
		m.CreatedAt = &now
		m.CreatedBy = &userId
		m.UpdatedAt = &now
		m.UpdatedBy = &userId
	case UpdateOp:
		m.UpdatedAt = &now
		m.UpdatedBy = &userId
	case DeleteOp:
		m.DeletedAt = &now
		m.DeletedBy = &userId
	default:
		return fmt.Errorf("unknown operation type")
	}

	return nil
}
