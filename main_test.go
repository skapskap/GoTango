package main

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestWebSocketConnection(t *testing.T) {
	server := NewUpgrader()
	go server.listen()

	testServer := httptest.NewServer(server.router)
	defer testServer.Close()

	conn := makeConnection(testServer)
	defer conn.Close()

	time.Sleep(50 * time.Millisecond)

	if len(server.clients) != 1 {
		t.Fatalf("expected 1 client, got %d", len(server.clients))
	}
}

func TestWebSocketMessaging(t *testing.T) {
	server := NewUpgrader()
	go server.listen()

	testServer := httptest.NewServer(server.router)
	defer testServer.Close()

	conn1 := makeConnection(testServer)
	conn2 := makeConnection(testServer)
	defer conn1.Close()
	defer conn2.Close()

	message := []byte("Hello, server!")
	if err := conn1.WriteMessage(websocket.TextMessage, message); err != nil {
		t.Fatal(err)
	}

	time.Sleep(50 * time.Millisecond)

	_, received, err := conn2.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}

	if string(received) != string(message) {
		t.Fatalf("expected message %s, got %s", message, received)
	}
}

func makeConnection(testServer *httptest.Server) *websocket.Conn {
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial("ws"+testServer.URL[4:]+"/ws", nil)

	if err != nil {
		panic(err)
	}

	return conn
}
