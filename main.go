package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Message struct {
	Text     string `json:"text"`
	Username string `json:"username"`
}

type ChatRoom struct {
	Messages []Message `json:"messages"`
}

var chatRoom = ChatRoom{Messages: make([]Message, 0)}
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleChatRoom(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Add the new client to the chat room
	client := &Client{Conn: conn, Username: ""}
	clients[client] = true

	// Listen for incoming messages from the client
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			delete(clients, client)
			break
		}

		if messageType == websocket.TextMessage {
			// Parse the incoming message
			var message Message
			err := json.Unmarshal(p, &message)
			if err != nil {
				log.Println(err)
				continue
			}

			// Set the username of the client if it hasn't been set yet
			if client.Username == "" {
				client.Username = message.Username
			}

			// Add the username to the message
			message.Username = client.Username

			// Add the message to the chat room
			chatRoom.Messages = append(chatRoom.Messages, message)

			// Broadcast the updated chat room to all connected clients
			for c := range clients {
				err := c.Conn.WriteJSON(chatRoom)
				if err != nil {
					log.Println(err)
					c.Conn.Close()
					delete(clients, c)
				}
			}
		}
	}
}


type Client struct {
	Conn     *websocket.Conn
	Username string
}

var clients = make(map[*Client]bool)

func main() {
	fmt.Println("Starting messaging API server...")
	http.HandleFunc("/chat", handleChatRoom)
	http.ListenAndServe(":8080", nil)
}
