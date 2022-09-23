package files

import (
	connectionManager "Tyche/ConnectionManager"
	webSocketGeneralTypes "Tyche/WebSocketGeneralTypes"
	"encoding/json"
	"errors"
	"os"

	"github.com/google/uuid"
)

type WriteFileResponse struct {
	ThreadID uuid.UUID `json:"thread_id"`
	Error    string    `json:"error"`
	Written  bool      `json:"written"`
}

type WriteFileMessage struct {
	Path    string `json:"path"`
	Content string `json:"Content"`
}

func WriteFile(threadID uuid.UUID, message json.RawMessage, con connectionManager.Connection) {
	if !CheckIfUUIDIsInWriteAuthorityQueue(con.UUID) {
		webSocketGeneralTypes.SendError(con, threadID, errors.New("no write authority requested"))
		return
	}

	if !CheckIfUUIDIsWriteAuthority(con.UUID) {
		webSocketGeneralTypes.SendError(con, threadID, errors.New("no write authority pls wait"))
		return
	}

	writeFileMessage := WriteFileMessage{}
	err := json.Unmarshal(message, &writeFileMessage)
	if err != nil {
		webSocketGeneralTypes.SendError(con, threadID, err)
		return
	}

	cleanPath := getCleanRepoPath(writeFileMessage.Path)
	if _, err := os.Stat(cleanPath); errors.Is(err, os.ErrNotExist) {
		webSocketGeneralTypes.SendError(con, threadID, err)
		return
	}

	if isDir(cleanPath) {
		webSocketGeneralTypes.SendError(con, threadID, errors.New("path is a directory"))
		return
	}

	file, err := os.Create(cleanPath)
	if err != nil {
		webSocketGeneralTypes.SendError(con, threadID, err)
		return
	}

	_, err = file.WriteString(writeFileMessage.Content)
	if err != nil {
		webSocketGeneralTypes.SendError(con, threadID, err)
		return
	}

	err = file.Close()
	if err != nil {
		webSocketGeneralTypes.SendError(con, threadID, err)
		return
	}

	con.WebSocket.WriteJSON(WriteFileResponse{
		ThreadID: threadID,
		Error:    "",
		Written:  true,
	})
}
