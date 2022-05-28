package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/panagiotisptr/hermes-messenger/chat-app/protos"
)

func GetLoginHandler(authenticationClient protos.AuthenticationClient) httpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		email, err := getFormValue(r.Form, "email")
		if err != nil {
			http.Redirect(w, r, "/register.html", 301)
			fmt.Println(err)
			return
		}

		password, err := getFormValue(r.Form, "password")
		if err != nil {
			http.Redirect(w, r, "/register.html", 301)
			fmt.Println(err)
			return
		}

		response, err := authenticationClient.Authenticate(
			context.Background(),
			&protos.AuthenticateRequest{
				Email:    email,
				Password: password,
			},
		)

		if err != nil {
			fmt.Println("Login error: ", err)
		}

		accessTokenCookie := &http.Cookie{
			Name:     "AccessToken",
			Value:    response.AccessToken,
			Expires:  time.Now().Add(time.Hour),
			Secure:   true,
			HttpOnly: true,
			Path:     "/",
		}

		refreshTokenCookie := &http.Cookie{
			Name:     "RefreshToken",
			Value:    response.RefreshToken,
			Expires:  time.Now().Add(time.Hour * 24),
			Secure:   true,
			HttpOnly: true,
			Path:     "/",
		}

		http.SetCookie(w, accessTokenCookie)
		http.SetCookie(w, refreshTokenCookie)

		http.Redirect(w, r, "/index.html", 301)
	}
}
