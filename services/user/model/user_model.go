package model

import (
	"github.com/google/uuid"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/entityutils"
)

// User struct representing a User object
type User struct {
	ID        *uuid.UUID `bson:"_id" json:"id"`
	Email     string     `bson:"email" json:"email"`
	FirstName string     `bson:"firstName" json:"firstName"`
	LastName  string     `bson:"lastName" json:"lastName"`

	entityutils.Meta
}
