package api

import (
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"github.com/masterkusok/websocketCollab/internal/businnesLogic"
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
}

type Client struct {
	Ip string
}

type Session struct {
	Id                  int
	Connect, Disconnect chan *websocket.Conn
	Broadcast           chan Message
	storage             *SessionStorage
	numberOfClients     int
	repository          *businnesLogic.DocumentRepository
	activeConnections   map[*websocket.Conn]*Client
	documentId          uint
	document            *businnesLogic.Document
}

func CreateSession(document *businnesLogic.Document, repository *businnesLogic.DocumentRepository, storage *SessionStorage) *Session {
	return &Session{
		Id:                1,
		documentId:        document.ID,
		numberOfClients:   0,
		Connect:           make(chan *websocket.Conn),
		Disconnect:        make(chan *websocket.Conn),
		Broadcast:         make(chan Message),
		activeConnections: make(map[*websocket.Conn]*Client),
		document:          document,
		repository:        repository,
		storage:           storage,
	}
}

func (s *Session) RunSession() {
	log.Printf("RUNNING SESSION ID: %d\nDocument initial text:\n%s\n", s.Id, s.document.Text)
	for {
		select {
		case connection := <-s.Connect:
			if client, ok := s.activeConnections[connection]; ok {
				log.Printf("client: %s trying to connect, while connection already exists", client.Ip)
				continue
			}
			s.numberOfClients++
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
			log.Printf("CURRENT TEXT:\n%s\n", s.document.GetText())
			s.syncDoc()

		case connection := <-s.Disconnect:
			if _, ok := s.activeConnections[connection]; !ok {
				log.Printf("Ip: %s trying to disconnect, while connection does not exists", connection.IP())
				continue
			}
			delete(s.activeConnections, connection)
			s.numberOfClients--
			log.Printf("Disconnected from session id: %d\n", s.Id)

			if s.numberOfClients == 0 {
				log.Printf("Stopping session with id %d\n", s.Id)
				err := (*s.repository).UpdateDocument(s.document)
				if err != nil {
					fmt.Printf("CANT UPDATE DOCUMENT, session id: %d\n", s.Id)
				}
				s.storage.Remove(s.Id)
				return
			}

		}
	}
}

func (s *Session) syncDoc() {
	for conn := range s.activeConnections {
		err := conn.WriteMessage(websocket.TextMessage, []byte(s.document.GetText()))
		if err != nil {
			log.Fatal(err)
		}
	}
}
