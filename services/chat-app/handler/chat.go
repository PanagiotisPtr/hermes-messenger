package handler

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/panagiotisptr/hermes-messenger/chat-app/middleware"
	"github.com/panagiotisptr/hermes-messenger/protos"
)

func GetChatHandler(
	authenticationClient protos.AuthenticationClient,
	messagingClient protos.MessagingClient,
) httpHandler {
	return middleware.Authentication(authenticationClient, func(w http.ResponseWriter, r *http.Request) {
		userUuid := strings.TrimPrefix(r.URL.Path, "/chat/")
		myUuid := r.Context().Value("UserUuid").(string)

		fmt.Println("UserUuid:", userUuid)
		fmt.Println("MyUuid:", myUuid)

		t, err := template.ParseFiles("templates/chat.html")
		if err != nil {
			fmt.Printf("Failed to load template file for results. Error: %v", err)
			return
		}

		messagesResponse, err := messagingClient.GetMessagesBetweenUsers(
			context.Background(),
			&protos.GetMessagesBetweenUsersRequest{
				From: myUuid,
				To:   userUuid,
			},
		)
		if err != nil {
			fmt.Printf("Failed to get messages for users. Error: %v", err)
			return
		}

		t.Execute(w, struct {
			UserUuid     string
			WebSocketUrl string
			Messages     []*protos.Message
		}{
			UserUuid:     userUuid,
			WebSocketUrl: "ws://localhost/ws",
			Messages:     messagesResponse.Messages,
		})
	})
}
