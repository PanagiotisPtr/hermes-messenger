package model

import (
	"github.com/google/uuid"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/entityutils"
)

type FriendStatus string

const (
	FriendStatusAccepted FriendStatus = "accepted"
	FriendStatusPending  FriendStatus = "pending"
	FriendStatusDeclined FriendStatus = "declined"
)

// Friend struct representing a Friend object
type Friend struct {
	ID       *uuid.UUID   `bson:"_id" json:"id"`
	UserId   *uuid.UUID   `bson:"userId" json:"userId"`
	FriendId *uuid.UUID   `bson:"friendId" json:"friendId"`
	Status   FriendStatus `bson:"status" json:"status"`

	entityutils.Meta
}
