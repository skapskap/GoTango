package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
}

type Message struct {
	author  *Client
	content []byte
}

type Server struct {
	router     *chi.Mux
	upgrader   websocket.Upgrader
	clients    []*Client
	message    chan Message
	register   chan *Client
	unregister chan *Client
}

func (c *Client) write(message []byte) {
	c.conn.WriteMessage(websocket.TextMessage, message)
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{conn: conn}
}

func (s *Server) listen() {
	for {
		select {
		case message := <-s.message:
			for _, client := range s.clients {
				if client != message.author {
					client.write(message.content)
				}
			}
		case client := <-s.register:
			s.clients = append(s.clients, client)
		case client := <-s.unregister:
			for i, c := range s.clients {
				if c == client {
					s.clients = append(s.clients[:i], s.clients[i+1:]...)
					break
				}
			}
		}
	}
}

func readMessages(server *Server, client *Client) {
	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				log.Println("client disconnected")
				server.unregister <- client
				return
			}
			log.Println("Unexpected error")
			log.Println(err)
			continue
		}

		log.Println("message received: ", string(message))
		server.message <- Message{author: client, content: message}
	}
}


func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New incoming connection from client!")
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
        log.Println("Failed to upgrade connection")
		log.Println(err)
		return
	}
	
	log.Println("client connected")

	client := NewClient(conn)
	s.register <- client

	go readMessages(s, client)
}

func NewUpgrader() *Server {
    s := &Server{}

    s.router = chi.NewRouter()
    s.router.Get("/ws", s.handleWS)

    s.upgrader = websocket.Upgrader{
        ReadBufferSize:  1024,
        WriteBufferSize: 1024,
        CheckOrigin: func(r *http.Request) bool {
            // Allow all origins for the time being
            return true
        },
    }

    s.message = make(chan Message)
	s.register = make(chan *Client)
	s.unregister = make(chan *Client)

	return s
}

func main() {
	server := NewUpgrader()
    go server.listen()

	fmt.Println("Server is running on port 4869")
	log.Fatal(http.ListenAndServe(":4869", server.router))
}