package sessions

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/masterkusok/websocketCollab/internal/documents"
	"log"
)

const (
	INSERT string = "INSERT"
	DELETE string = "DELETE"
)

type Message struct {
	CMD      string  `json:"cmd"`
	Position float32 `json:"position"`
	Value    byte    `json:"value"`
	// Connection *websocket.Conn `json:"connection"`
}

type Client struct {
	Ip string
}

type Session struct {
	Id                  int
	documentId          int
	document            *documents.Document
	numberOfClients     int
	Connect, Disconnect chan *websocket.Conn
	Broadcast           chan Message
	activeConnections   map[*websocket.Conn]*Client
}

func CreateSession(documentId int) *Session {
	return &Session{
		Id:                1,
		documentId:        documentId,
		numberOfClients:   0,
		Connect:           make(chan *websocket.Conn),
		Disconnect:        make(chan *websocket.Conn),
		Broadcast:         make(chan Message),
		activeConnections: make(map[*websocket.Conn]*Client),
		document:          documents.CreateDocument("temp"),
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
			switch message.CMD {
			case INSERT:
				s.document.Insert(message.Position, message.Value)
			case DELETE:
				s.document.Delete(message.Position)
			}
			log.Printf("CLIENT, changed file\n")
			s.syncDoc()
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

func (s *Session) syncDoc() {
	for conn, _ := range s.activeConnections {
		err := conn.WriteMessage(websocket.TextMessage, []byte(s.document.Text()))
		if err != nil {
			log.Fatal(err)
		}
	}
}
