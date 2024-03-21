package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

type Server struct {
	upgrader websocket.Upgrader
}

func NewUpgrader() *Server {
    s := &Server{}
    s.upgrader = websocket.Upgrader{
        ReadBufferSize:  1024,
        WriteBufferSize: 1024,
        CheckOrigin: func(r *http.Request) bool {
            // Allow all origins for the time being
            return true
        },
    }

    return s
}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New incoming connection from client!")
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	
	defer conn.Close()

    _, message, err := conn.ReadMessage()
    if err != nil {
        log.Println(err)
        return
    }
    fmt.Println("Message received:", string(message))

    err = conn.WriteMessage(websocket.TextMessage, message)
    if err != nil {
        log.Println(err)
        return
    }
    fmt.Println("Message sent:", string(message))
}

func main() {
	server := NewUpgrader()

	r := chi.NewRouter()
	r.Get("/ws", server.handleWS)

	fmt.Println("Server is running on port 4869")
	log.Fatal(http.ListenAndServe(":4869", r))
}