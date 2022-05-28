package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/panagiotisptr/hermes-messenger/chat-app/protos"
)

func GetRegisterHandler(authenticationClient protos.AuthenticationClient) httpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		email, err := getFormValue(r.Form, "email")
		if err != nil {
			http.Redirect(w, r, "/register.html", 301)
			return
		}

		password, err := getFormValue(r.Form, "password")
		if err != nil {
			http.Redirect(w, r, "/register.html", 301)
			return
		}

		confirmPassword, err := getFormValue(r.Form, "confirmPassword")
		if err != nil {
			http.Redirect(w, r, "/register.html", 301)
			return
		}

		if password != confirmPassword {
			http.Redirect(w, r, "/register.html", 301)
			return
		}

		response, err := authenticationClient.Register(
			context.Background(),
			&protos.RegisterRequest{
				Email:    email,
				Password: password,
			},
		)

		if err != nil {
			fmt.Println(err)
			return
		}

		if !response.Success {
			http.Redirect(w, r, "/register.html", 301)
			return
		}

		http.Redirect(w, r, "/index.html", 301)
	}
}
