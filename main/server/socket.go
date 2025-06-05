package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// handleWebSocket handles WebSocket connections
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	// Send welcome message
	welcomeMsg := fmt.Sprintf("Welcome to WebSocket Server! Connected at: %s",
		time.Now().Format(time.RFC3339))
	err = conn.WriteMessage(websocket.TextMessage, []byte(welcomeMsg))
	if err != nil {
		log.Printf("Error writing welcome message: %v", err)
		return
	}

	log.Printf("Client connected: %s", r.RemoteAddr)

	// Simple message handling loop
	for {
		// Read message
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		// Log received message
		log.Printf("Received: %s", message)

		// Echo the message back with timestamp
		response := fmt.Sprintf("Echo at %s: %s",
			time.Now().Format("15:04:05"), string(message))

		if err := conn.WriteMessage(messageType, []byte(response)); err != nil {
			log.Printf("Error writing message: %v", err)
			break
		}
	}
}

// SetupWebSocketServer initializes the WebSocket server
func SetupWebSocketServer() {
	http.HandleFunc("/ws", handleWebSocket)
	log.Println("WebSocket server enabled at ws://localhost:8080/ws")
}
