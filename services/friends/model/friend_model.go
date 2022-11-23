package model

import (
	"github.com/google/uuid"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/entityutils"
)

type FriendStatus int64

const (
	FriendStatusPendingFriend FriendStatus = iota
	FriendStatusPendingUser
	FriendStatusAccepted
	FriendStatusDeclined
)

// Friend struct representing a Friend object
type Friend struct {
	ID       *uuid.UUID   `bson:"_id" json:"id"`
	UserId   uuid.UUID    `bson:"userId" json:"userId"`
	FriendId uuid.UUID    `bson:"friendId" json:"friendId"`
	Status   FriendStatus `bson:"status" json:"status"`

	entityutils.Meta
}
