package user

import (
	"log"

	"github.com/google/uuid"
)

type Service struct {
	logger   *log.Logger
	userRepo Repository
}

func NewService(logger *log.Logger) *Service {
	return &Service{
		logger:   logger,
		userRepo: NewMemoryRepository(logger),
	}
}

func (s *Service) RegisterUser(userUuid uuid.UUID, email string) error {
	return s.userRepo.AddUser(userUuid, email)
}

func (s *Service) GetUser(userUuid uuid.UUID) (*User, error) {
	return s.userRepo.GetUser(userUuid)
}
