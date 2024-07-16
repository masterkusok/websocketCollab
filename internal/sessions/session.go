package sessions

import (
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"log"
)

type Client struct {
	Ip string
}

type Session struct {
	Id                  int
	documentId          int
	numberOfClients     int
	Connect, Disconnect chan *websocket.Conn
	Broadcast           chan string
	activeConnections   map[*websocket.Conn]*Client
}

func CreateSession(documentId int) *Session {
	return &Session{
		Id:                1,
		documentId:        documentId,
		numberOfClients:   0,
		Connect:           make(chan *websocket.Conn),
		Disconnect:        make(chan *websocket.Conn),
		Broadcast:         make(chan string),
		activeConnections: make(map[*websocket.Conn]*Client),
	}
}

func (s *Session) RunSession() {
	for {
		select {
		case connection := <-s.Connect:
			if client, ok := s.activeConnections[connection]; ok {
				log.Printf("client: %s trying to connect, while connection already exists", client.Ip)
				continue
			}
			client := Client{Ip: connection.IP()}
			s.activeConnections[connection] = &client

		case message := <-s.Broadcast:
			fmt.Printf("Session with id %d recieved message:\n%s\n", s.Id, message)
		case connection := <-s.Disconnect:
			if _, ok := s.activeConnections[connection]; !ok {
				log.Printf("Ip: %s trying to disconnect, while connection does not exists", connection.IP())
				continue
			}
			delete(s.activeConnections, connection)
			err := connection.Close()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
