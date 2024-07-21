package api

import (
	"encoding/json"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
)

func CreateRouting(app *fiber.App, handler *Handler) {
	app.Get("api/v1/documents/:id", websocket.New(handler.connectToDocumentHandler))
	app.Post("api/v1/documents", handler.createDocumentHandler)
}

func (h *Handler) createDocumentHandler(c *fiber.Ctx) error {
	doc, err := h.documentRepository.CreateDocument()
	if err != nil {
		return err
	}
	result, _ := json.Marshal(doc)
	return c.SendString(string(result))
}

func (h *Handler) connectToDocumentHandler(conn *websocket.Conn) {
	var id string
	var err error
	var message Message
	var activeSession *Session

	defer func() {
		conn.Close()
	}()

	id = conn.Params("id", "-1")
	if id == "-1" {
		log.Printf("DIDNT PASSED ID, ip: %s\n", conn.IP())
		return
	}

	intId, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("ID IS NOT INT, ip: %s\n", conn.IP())
		return
	}

	doc, err := h.documentRepository.GetById(intId)
	if err != nil {
		log.Printf("ID NOT FOUND, ip: %s\n", conn.IP())
	}

	if !h.activeSessions.HasKey(intId) {
		activeSession = CreateSession(doc, &h.documentRepository, h.activeSessions)
		go activeSession.RunSession()
		h.activeSessions.Add(intId, activeSession)
		log.Printf("RUN NEW SESSION, DOCUMENT ID: %d, ip: %s\n", intId, conn.IP())
	} else {
		activeSession = h.activeSessions.Get(intId)
	}

	defer func() {
		activeSession.Disconnect <- conn
	}()

	activeSession.Connect <- conn
	for {
		_, messageBytes, err := conn.ReadMessage()
		if err != nil {
			return
		}

		err = json.Unmarshal(messageBytes, &message)
		if err != nil {
			log.Printf("INVALID MESSAGE STRUCTURE, ip: %s\nError text%s\n", conn.IP(), err)
		}

		activeSession.Broadcast <- message
	}
}
