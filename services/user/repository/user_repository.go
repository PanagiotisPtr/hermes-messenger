package repository

import (
	"github.com/panagiotisptr/hermes-messenger/libs/utils/entityutils"
	"github.com/panagiotisptr/hermes-messenger/user/model"
)

type UserRepository interface {
	entityutils.Repository[model.User]
}
