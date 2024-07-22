package api

import (
	"github.com/masterkusok/websocketCollab/internal/businnesLogic"
)

type Handler struct {
	documentRepository businnesLogic.DocumentRepository
	activeSessions     *SessionStorage
}

func NewHandler(dr businnesLogic.DocumentRepository) *Handler {
	return &Handler{documentRepository: dr, activeSessions: NewSessionStorage()}
}
