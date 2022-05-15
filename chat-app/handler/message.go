package handler

import (
	"context"
	"fmt"
	"net/http"
	"panagiotisptr/chat-app/middleware"
	"panagiotisptr/chat-app/protos"
	"strings"
)

func GetMessageHandler(
	authenticationClient protos.AuthenticationClient,
	messagingClient protos.MessagingClient,
) httpHandler {
	return middleware.Authentication(authenticationClient, func(w http.ResponseWriter, r *http.Request) {
		userUuid := strings.TrimPrefix(r.URL.Path, "/message/")
		myUuid := r.Context().Value("UserUuid").(string)

		r.ParseForm()

		fmt.Println("userUuid:", userUuid)

		content, err := getFormValue(r.Form, "content")
		if err != nil {
			http.Redirect(w, r, "/chat/"+userUuid, 301)
			return
		}

		fmt.Println("Content:", content)

		messagingClient.SendMessage(
			context.Background(),
			&protos.SendMessageRequest{
				From:    myUuid,
				To:      userUuid,
				Content: content,
			},
		)

		http.Redirect(w, r, "/chat/"+userUuid, 301)
	})
}
