package entityutils

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Meta struct {
	CreatedBy *uuid.UUID `bson:"createdBy" json:"createdBy"`
	UpdatedBy *uuid.UUID `bson:"updatedBy" json:"updatedBy"`
	DeletedBy *uuid.UUID `bson:"deletedBy" json:"deletedBy"`

	CreatedAt *time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt *time.Time `bson:"updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `bson:"deletedAt" json:"deletedAt"`
}

func (m *Meta) UpdateMeta(
	ctx context.Context,
	op OperationType,
) error {
	userId, ok := ctx.Value("user-id").(uuid.UUID)
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
