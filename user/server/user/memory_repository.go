package user

import (
	"fmt"
	"log"

	"github.com/google/uuid"
)

type MemoryRepository struct {
	logger *log.Logger
	users  map[uuid.UUID]*User
}

func NewMemoryRepository(logger *log.Logger) *MemoryRepository {
	return &MemoryRepository{
		logger: logger,
		users:  make(map[uuid.UUID]*User),
	}
}

func (mr *MemoryRepository) AddUser(userUuid uuid.UUID, email string) error {
	if _, ok := mr.users[userUuid]; ok {
		return fmt.Errorf("There's already a user with uuid: %s", userUuid.String())
	}

	mr.users[userUuid] = &User{
		Uuid:  userUuid,
		Email: email,
	}

	return nil
}

func (mr *MemoryRepository) GetUser(userUuid uuid.UUID) (*User, error) {
	user, ok := mr.users[userUuid]
	if !ok {
		return nil, fmt.Errorf("Could not find user with uuid: %s", userUuid.String())
	}

	return user, nil
}
