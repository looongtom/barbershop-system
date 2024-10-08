package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var globalMessage *[]byte
var bookingMessage *[]byte

// Upgrader to handle HTTP to WebSocket upgrade
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections by default
		return true
	},
}

// WebSocket handler that will be called on each connection
func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP request to a WebSocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Error upgrading to websocket: %v", err)
		return
	}
	defer ws.Close()

	// Infinite loop for sending messages
	for {
		if globalMessage != nil {
			err := ws.WriteMessage(websocket.TextMessage, *globalMessage)
			if err != nil {
				log.Printf("Error writing message to websocket: %v", err)
				break
			}
			globalMessage = nil
		}
	}
}

func handleBookingConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP request to a WebSocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Error upgrading to websocket: %v", err)
		return
	}
	defer ws.Close()

	// Infinite loop for sending messages
	for {
		if bookingMessage != nil {
			err := ws.WriteMessage(websocket.TextMessage, *bookingMessage)
			if err != nil {
				log.Printf("Error writing message to websocket: %v", err)
				break
			}
			bookingMessage = nil
		}
	}
}

// WebSocket handler that will be called on each connection
func handleTrigger(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Error upgrading to websocket: %v", err)
		return
	}
	defer ws.Close()
	_, message, err := ws.ReadMessage()
	if err != nil {
		log.Printf("Error reading message: %v", err)
	}
	log.Printf("message: %s", message)
	globalMessage = &message
}

func handleBookingTrigger(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Error upgrading to websocket: %v", err)
		return
	}
	defer ws.Close()
	_, message, err := ws.ReadMessage()
	if err != nil {
		log.Printf("Error reading message: %v", err)
	}
	log.Printf("message: %s", message)
	bookingMessage = &message
}

func main() {
	http.HandleFunc("/hairfast", handleConnections)

	http.HandleFunc("/trigger_hairfast", handleTrigger)

	http.HandleFunc("/booking", handleBookingConnections)

	http.HandleFunc("/trigger_booking", handleBookingTrigger)

	// Start the WebSocket server on localhost:8080
	log.Println("WebSocket server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
