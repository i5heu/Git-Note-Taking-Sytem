package main

import (
	"encoding/json"
	"net/http"
	"os"

	channels "Tyche/Channels"
	config "Tyche/Config"
	connectionManager "Tyche/ConnectionManager"
	files "Tyche/Files"
	gitManager "Tyche/GitManager"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Print("Starting the application...")

	config := config.GetConfig()
	log.Info().Msg("Config loaded")

	gitManager.GitManager(config)
	log.Info().Msg("Repo prepared")

	// Queues
	go connectionManager.ConnectionManagerWorker(channels.ConnectionManagerChan)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", socketHandler)
	http.ListenAndServe("localhost:8080", nil)
}

var upgrader = websocket.Upgrader{} // use default options

type Message struct {
	ThreadID    uuid.UUID       `json:"thread_id"`
	MessageType string          `json:"type"`
	Data        json.RawMessage `json:"data"`
}

type ResponseMessage struct {
	ThreadID    uuid.UUID `json:"thread_id"`
	MessageType string    `json:"type"`
	Data        string    `json:"data"`
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}
	defer conn.Close()

	messageCount := -1
	authenticated := false
	connection := connectionManager.Connection{}

	// The event loop
	for {
		messageCount++
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Error().AnErr("Error during message reading", err)
			break
		}

		data := Message{}
		err = json.Unmarshal(message, &data)
		if err != nil {
			log.Error().Err(err).Msg("Error during message unmarshalling")
			conn.WriteMessage(messageType, []byte("Error during message unmarshalling"))
			break
		}

		// if data.thread_id is not set reject the message
		if data.ThreadID == uuid.Nil {
			log.Error().Msg("ThreadID is not set")
			conn.WriteMessage(messageType, []byte("ThreadID is not set"))
			continue
		}

		if authenticated == true {
			switch data.MessageType {
			case "PING":
				conn.WriteJSON(ResponseMessage{
					ThreadID:    data.ThreadID,
					MessageType: "PONG",
					Data:        "PONG",
				})
			case "READ":
				files.ReadFile(data.ThreadID, data.Data, connection)
			default:
				conn.WriteMessage(messageType, []byte("Unknown message type"))
			}
		} else if data.MessageType == "REGISTER" && authenticated == false {
			log.Debug().Int("messageCount", messageCount).Msg("Registering a new connection")

			if AuthenticateConnection(&conn, data.Data) {
				log.Debug().Msg("Connection authenticated")
				ResultChan := make(chan []connectionManager.Connection)
				channels.ConnectionManagerChan <- connectionManager.ConnectionJob{Action: connectionManager.Register, WebSocket: conn, Result: ResultChan}
				resultCon := <-ResultChan
				connection = resultCon[0]
				authenticated = true
				result, err := json.Marshal(ResponseMessage{
					ThreadID:    data.ThreadID,
					MessageType: "REGISTER.OK",
					Data:        resultCon[0].UUID.String(),
				})
				if err != nil {
					log.Error().Err(err).Msg("Error during marshalling")
					conn.WriteMessage(messageType, []byte("Error during marshalling"))
					break
				}
				conn.WriteMessage(messageType, result)
			} else {
				log.Debug().Msg("Authentication failed")
				conn.WriteMessage(messageType, []byte("Authentication failed"))
				break
			}
		} else {
			log.Debug().Int("messageCount", messageCount).Msg("You are not authenticated")
			conn.WriteMessage(messageType, []byte("You are not authenticated"))
			break
		}
	}
}
