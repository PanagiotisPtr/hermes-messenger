package repository

import (
	"github.com/panagiotisptr/hermes-messenger/friends/model"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/entityutils"
)

// FriendRepository repository for the Friend model
type FriendRepository interface {
	entityutils.Repository[model.Friend]
}
