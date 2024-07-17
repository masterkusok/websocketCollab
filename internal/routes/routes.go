package routes

import (
	"encoding/json"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/masterkusok/websocketCollab/internal/sessions"
	"strconv"
	"sync"
)

var (
	activeSessions map[int]*sessions.Session
	mutex          sync.Mutex
)

func CreateRouting(app *fiber.App) {
	activeSessions = map[int]*sessions.Session{}
	app.Get("api/v1/document/:id", websocket.New(connectToDocumentHandler))
}

func connectToDocumentHandler(conn *websocket.Conn) {
	id, err := strconv.Atoi(conn.Params("id", "-1"))
	if err != nil {
		log.Info(conn.Close())
		return
	}
	if id == -1 {
		log.Info(conn.Close())
		return
	}
	session, ok := activeSessions[id]
	if !ok {
		session = sessions.CreateSession(id)
		activeSessions[id] = session
		go session.RunSession()
	}

	session.Connect <- conn
	defer func() { session.Disconnect <- conn }()

	for {
		messageType, jsonMessage, err := conn.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}
		message := sessions.Message{}
		err = json.Unmarshal(jsonMessage, &message)
		if err != nil {
			log.Fatal(err)
		}
		if messageType == websocket.TextMessage {
			session.Broadcast <- message
		} else {
			log.Fatal("unexpected type of message")
		}
	}

}
