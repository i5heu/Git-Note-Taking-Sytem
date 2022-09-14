package webSocketGeneralTypes

import (
	connectionManager "Tyche/ConnectionManager"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type ErrorMsg struct {
	ThreadID    uuid.UUID `json:"thread_id"`
	MessageType string    `json:"type"`
	Data        string    `json:"data"`
}

func SendError(con connectionManager.Connection, threadID uuid.UUID, err error) {
	log.Error().Err(err)

	con.WebSocket.WriteJSON(ErrorMsg{
		ThreadID:    threadID,
		MessageType: "error",
		Data:        err.Error(),
	})
}
