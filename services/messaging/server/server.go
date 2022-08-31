package server

import (
	"context"

	"github.com/google/uuid"
	"github.com/panagiotisptr/hermes-messenger/messaging/server/messaging"
	"github.com/panagiotisptr/hermes-messenger/protos"
	"go.uber.org/zap"
)

type MessagingServer struct {
	logger  *zap.Logger
	service *messaging.Service
	protos.UnimplementedMessagingServer
}

func ProvideUserServer(
	logger *zap.Logger,
	service *messaging.Service,
) (*MessagingServer, error) {
	return &MessagingServer{
		logger:  logger,
		service: service,
	}, nil
}

func (m *MessagingServer) SendMessage(
	ctx context.Context,
	request *protos.SendMessageRequest,
) (*protos.SendMessageResponse, error) {
	fromUuid, err := uuid.Parse(request.From)
	if err != nil {
		return &protos.SendMessageResponse{}, err
	}

	toUuid, err := uuid.Parse(request.To)
	if err != nil {
		return &protos.SendMessageResponse{}, err
	}
	err = m.service.SendMessage(
		ctx,
		fromUuid,
		toUuid,
		request.Content,
	)

	return &protos.SendMessageResponse{}, err
}

func (m *MessagingServer) GetMessages(
	ctx context.Context,
	request *protos.GetMessagesRequest,
) (*protos.GetMessagesResponse, error) {
	messages := make([]*protos.Message, 0)
	fromUuid, err := uuid.Parse(request.From)
	if err != nil {
		return &protos.GetMessagesResponse{
			Messages: messages,
		}, err
	}

	toUuid, err := uuid.Parse(request.To)
	if err != nil {
		return &protos.GetMessagesResponse{
			Messages: messages,
		}, err
	}

	ms, err := m.service.GetMessagesBetweenUsers(
		ctx,
		fromUuid,
		toUuid,
		request.Size,
		request.Offset,
	)

	for _, m := range ms {
		messages = append(messages, &protos.Message{
			From:      m.From.String(),
			To:        m.To.String(),
			Timestamp: m.Timestamp,
			Content:   m.Content,
		})
	}

	return &protos.GetMessagesResponse{
		Messages: messages,
	}, err
}
