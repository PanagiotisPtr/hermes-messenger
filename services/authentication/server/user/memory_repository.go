package user

import (
	"fmt"

	"go.uber.org/zap"
)

type MemoryRepository struct {
	logger *zap.Logger
	users  map[string]User
}

func ProvideMemoryUserRepository(logger *zap.Logger) Repository {
	return &MemoryRepository{
		logger: logger,
		users:  make(map[string]User),
	}
}

func (mr *MemoryRepository) CreateUser(user User) error {
	if _, ok := mr.users[user.Email]; ok {
		return fmt.Errorf("User with email %s already exists", user.Email)
	}
	mr.users[user.Email] = user

	return nil
}

func (mr *MemoryRepository) GetUserByEmail(email string) (User, error) {
	user, ok := mr.users[email]
	if !ok {
		return User{}, fmt.Errorf("Could not find user with email %s", email)
	}

	return user, nil
}

func (mr *MemoryRepository) DeleteUser(user User) error {
	delete(mr.users, user.Email)

	return nil
}
