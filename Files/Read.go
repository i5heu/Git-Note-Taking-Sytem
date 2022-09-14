package files

import (
	connectionManager "Tyche/ConnectionManager"
	webSocketGeneralTypes "Tyche/WebSocketGeneralTypes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/google/uuid"
)

type ReadFileMessage struct {
	Path string `json:"path"`
}

type ReadFileResponse struct {
	ThreadID uuid.UUID `json:"thread_id"`
	IsDir    bool      `json:"isDir"`
	Data     string    `json:"data"`
	Children []string  `json:"children"`
}

func ReadFile(threadID uuid.UUID, message json.RawMessage, con connectionManager.Connection) {
	readFileMessage := ReadFileMessage{}
	err := json.Unmarshal(message, &readFileMessage)
	if err != nil {
		webSocketGeneralTypes.SendError(con, threadID, err)
		return
	}

	cleanPath := getCleanRepoPath(readFileMessage.Path)
	if _, err := os.Stat(cleanPath); errors.Is(err, os.ErrNotExist) {
		webSocketGeneralTypes.SendError(con, threadID, err)
		return
	}

	if isDir(cleanPath) {
		sendDirectoryListing(cleanPath, con, threadID)
	} else {
		sendFile(con, cleanPath, threadID)
	}
}

func sendFile(con connectionManager.Connection, cleanPath string, threadID uuid.UUID) {
	file, err := ioutil.ReadFile(cleanPath)
	if err != nil {
		webSocketGeneralTypes.SendError(con, threadID, err)
		return
	}

	con.WebSocket.WriteJSON(ReadFileResponse{
		ThreadID: threadID,
		IsDir:    false,
		Data:     string(file),
	})
}

func sendDirectoryListing(cleanPath string, con connectionManager.Connection, threadID uuid.UUID) {
	files, err := ioutil.ReadDir(cleanPath)
	if err != nil {
		webSocketGeneralTypes.SendError(con, threadID, err)
		return
	}

	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	con.WebSocket.WriteJSON(ReadFileResponse{
		ThreadID: threadID,
		IsDir:    true,
		Children: fileNames,
	})
}
