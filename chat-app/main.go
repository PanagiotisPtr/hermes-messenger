package main

import (
	"fmt"
	"log"
	"net/http"
	"panagiotisptr/chat-app/handler"
	"panagiotisptr/chat-app/protos"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}

func main() {
	authenticationConn, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer authenticationConn.Close()

	authenticationClient := protos.NewAuthenticationClient(authenticationConn)

	messagingConn, err := grpc.Dial(":9090", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer messagingConn.Close()

	messagingClient := protos.NewMessagingClient(messagingConn)

	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/api/login", handler.GetLoginHandler(authenticationClient))
	http.HandleFunc("/api/register", handler.GetRegisterHandler(authenticationClient))
	http.HandleFunc("/chat/", handler.GetChatHandler(authenticationClient, messagingClient))
	http.HandleFunc("/message/", handler.GetMessageHandler(authenticationClient, messagingClient))

	http.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		for _, c := range r.Cookies() {
			fmt.Println(c.Name)
			fmt.Println(c.Value)
		}
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
		}
		log.Println("Client Connected")

		reader(ws)
	})

	http.ListenAndServe(":80", nil)
}
