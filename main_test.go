package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestUseWS(t *testing.T) {
	server := NewUpgrader()

	// Create a test HTTP server
	testServer := httptest.NewServer(http.HandlerFunc(server.handleWS))
	defer testServer.Close()

	// Convert the test server URL to a WebSocket URL
	wsURL := "ws" + strings.TrimPrefix(testServer.URL, "http")

	// Create a WebSocket client
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket server: %v", err)
	}
	defer ws.Close()

	// Send a message to the WebSocket server
	message := "Hello, server!"
	err = ws.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		t.Fatalf("Failed to send message to WebSocket server: %v", err)
	}

	// Read the response from the WebSocket server
	_, response, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read message from WebSocket server: %v", err)
	}

	// Assert the response from the server
	expectedResponse := "Hello, client!"
	if string(response) != expectedResponse {
		t.Errorf("Unexpected response from server. Got: %s, Expected: %s", string(response), expectedResponse)
	}
}
