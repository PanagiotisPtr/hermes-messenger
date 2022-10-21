package model

import (
	"github.com/google/uuid"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/entityutils"
)

type Friend struct {
	ID       *uuid.UUID `bson:"_id" json:"id"`
	UserID   uuid.UUID  `bson:"userId" json:"userId"`
	FriendID uuid.UUID  `bson:"friendId" json:"friendId"`
	Status   string     `bson:"status" json:"status"`

	entityutils.Meta
}
