package friends

import (
	"fmt"
	"log"
	"panagiotisptr/friends/server/connection"
	"panagiotisptr/friends/server/connection/status"

	"github.com/google/uuid"
)

type Friend struct {
	UserUuid uuid.UUID
	Status   string
}

type Service struct {
	logger     *log.Logger
	repository connection.Repository
}

func NewService(logger *log.Logger) *Service {
	return &Service{
		logger:     logger,
		repository: connection.NewMemoryRepository(logger),
	}
}

func (s *Service) AddFriend(userUuid uuid.UUID, friendUuid uuid.UUID) error {
	if userUuid == friendUuid {
		return fmt.Errorf("A user can't be a friend with themselves")
	}

	connection, err := s.repository.GetConnection(userUuid, friendUuid)
	if err == nil {
		if connection.Status != status.Accepted {
			return s.repository.UpdateConnectionStatus(*connection, status.Pending)
		} else {
			return nil
		}
	}

	return s.repository.AddConnection(userUuid, friendUuid)
}

func (s *Service) RemoveFriend(userUuid uuid.UUID, friendUuid uuid.UUID) error {
	connection, err := s.repository.GetConnection(userUuid, friendUuid)
	if err != nil {
		return err
	}

	err = s.repository.UpdateConnectionStatus(*connection, status.Rejected)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetFriends(userUuid uuid.UUID) ([]Friend, error) {
	friends := make([]Friend, 0)
	connections, err := s.repository.GetConnectionsForUser(userUuid)
	if err != nil {
		return friends, err
	}

	for _, connection := range connections {
		friendUuid := connection.From
		if friendUuid == userUuid {
			friendUuid = connection.To
		}
		friends = append(friends, Friend{
			UserUuid: friendUuid,
			Status:   connection.Status,
		})
	}

	return friends, nil
}
