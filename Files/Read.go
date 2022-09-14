package files

import (
	connectionManager "Tyche/ConnectionManager"
	"encoding/json"
	"io/ioutil"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
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
		log.Error().Err(err).Msg("Error during message unmarshalling")
		con.WebSocket.WriteMessage(websocket.TextMessage, []byte("Error during message unmarshalling"))
		return
	}

	log.Debug().Str("path", readFileMessage.Path).Msg("Reading file")

	cleanPath := getCleanRepoPath(readFileMessage.Path)

	// check if the file is a directory
	if isDir(cleanPath) {
		sendDirectoryListing(cleanPath, con, threadID)
	} else {
		sendFile(con, cleanPath, threadID)
	}
}

func sendFile(con connectionManager.Connection, cleanPath string, threadID uuid.UUID) {
	file, err := ioutil.ReadFile(cleanPath)
	if err != nil {
		log.Error().Err(err).Msg("Error during file reading")
		con.WebSocket.WriteMessage(websocket.TextMessage, []byte("Error during file reading: "+err.Error()))
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
		log.Error().Err(err).Msg("Error reading directory")
		con.WebSocket.WriteMessage(websocket.TextMessage, []byte("Error reading directory"))
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
