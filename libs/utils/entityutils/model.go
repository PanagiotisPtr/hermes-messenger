package entityutils

import (
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
