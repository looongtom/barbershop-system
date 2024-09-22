package main

import (
	"DoAn/pkg/previewimage/api"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

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
		message := fmt.Sprintf("Current time: %v", time.Now().Format(time.RFC3339))
		err := ws.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Printf("Error writing message to websocket: %v", err)
			break
		}

		// Wait for 1 second before sending the next message
		time.Sleep(1 * time.Second)
	}
}

func main() {
	//http.HandleFunc("/ws", handleConnections)
	//
	//// Start the WebSocket server on localhost:8080
	//log.Println("WebSocket server started on :8081")
	//err := http.ListenAndServe(":8081", nil)
	//if err != nil {
	//	log.Fatalf("Error starting server: %v", err)
	//}

	var result api.HairFastResult
	result = api.HairFastResult{
		"1",
		"2",
		"3",
		"4",
	}
	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Printf("Failed to marshal HairFastResult to JSON: %v", err)
		return
	}
	// Create WebSocket client connection
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/trigger_hairfast", http.Header{})
	if err != nil {
		log.Printf("Failed to connect to WebSocket server: %v", err)
		return
	}
	defer conn.Close()

	// Send the data
	err = conn.WriteMessage(websocket.TextMessage, resultJSON)
	if err != nil {
		log.Printf("Failed to send message to WebSocket server: %v", err)
	} else {
		log.Printf("Successfully sent message to WebSocket server: %s", resultJSON)
	}
}
