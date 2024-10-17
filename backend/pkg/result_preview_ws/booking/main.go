package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var globalMessage *[]byte
var bookingMessage *[]byte

// Upgrader to handle HTTP to WebSocket upgrade
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var clients = make(map[*websocket.Conn]bool)
var mu sync.Mutex

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

// Handle WebSocket connections.
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	// Add the new connection to the clients map.
	mu.Lock()
	clients[conn] = true
	mu.Unlock()

	// Clean up when the connection is closed.
	defer func() {
		mu.Lock()
		delete(clients, conn)
		mu.Unlock()
	}()

	// Listen for messages from the WebSocket connection.
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Connection closed or read error: %v", err)
			break
		}

		// Process the incoming message.
		if messageType == websocket.TextMessage || messageType == websocket.BinaryMessage {
			var bookingResponse BookingResponse
			err := json.Unmarshal(message, &bookingResponse)
			if err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}

			// Handle the BookingResponse data (e.g., log it, store it, or process it).
			log.Printf("Received BookingResponse: %+v", bookingResponse)

			// Optionally, send a response back to the client.
			responseMessage := fmt.Sprintf("Booking %s for customer %s received.", bookingResponse.BookingID, bookingResponse.CustomerID)
			err = conn.WriteMessage(websocket.TextMessage, []byte(responseMessage))
			if err != nil {
				log.Printf("Failed to send response message: %v", err)
			}
		} else {
			log.Printf("Unsupported message type: %d", messageType)
		}
	}
}

func main() {
	// WebSocket endpoint for mobile app connections.
	http.HandleFunc("/booking", handleWebSocket)

	// Endpoint for backend service to send data.
	http.HandleFunc("/trigger_booking", handleBackendData)

	// Start the server.
	log.Println("WebSocket server is running on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

//func main() {
//	http.HandleFunc("/hairfast", handleConnections)
//
//	http.HandleFunc("/trigger_hairfast", handleTrigger)
//
//	http.HandleFunc("/booking", handleBookingConnections)
//
//	http.HandleFunc("/trigger_booking", handleBookingTrigger)
//
//	// Start the WebSocket server on localhost:8080
//	log.Println("WebSocket server started on :8080")
//	err := http.ListenAndServe(":8080", nil)
//	if err != nil {
//		log.Fatalf("Error starting server: %v", err)
//	}
//}
