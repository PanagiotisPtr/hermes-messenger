package server

import (
	"context"
	"log"
	"panagiotisptr/messaging/protos"
	"panagiotisptr/messaging/server/messaging"
)

type MessagingServer struct {
	logger  *log.Logger
	service *messaging.MessagingService
	protos.UnimplementedMessagingServer
}

func NewMessagingServer(l *log.Logger) *MessagingServer {
	return &MessagingServer{
		logger:  l,
		service: messaging.NewMessagingService(l),
	}
}

func (m *MessagingServer) SendMessage(
	ctx context.Context,
	request *protos.SendMessageRequest,
) (*protos.SendMessageResponse, error) {
	err := m.service.SendMessage(
		request.From,
		request.To,
		request.Content,
	)

	return &protos.SendMessageResponse{
		Sent: err == nil,
	}, err
}

func (m *MessagingServer) GetMessagesBetweenUsers(
	context context.Context,
	request *protos.GetMessagesBetweenUsersRequest,
) (*protos.GetMessagesBetweenUsersResponse, error) {
	messages, err := m.service.GetMessagesBetweenUsers(
		request.From,
		request.To,
	)
	if err != nil {
		return &protos.GetMessagesBetweenUsersResponse{
			Messages: make([]*protos.Message, 0),
		}, err
	}

	protoMessages := make([]*protos.Message, len(messages))
	for index, message := range messages {
		protoMessages[index] = &protos.Message{
			From:      message.From.String(),
			To:        message.To.String(),
			Timestamp: message.Timestamp,
			Content:   message.Content,
		}
	}

	return &protos.GetMessagesBetweenUsersResponse{
		Messages: protoMessages,
	}, nil
}
