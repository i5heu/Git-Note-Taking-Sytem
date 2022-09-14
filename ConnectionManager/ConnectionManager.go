package connectionManager

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ConnectionJob struct {
	Action    Action
	WebSocket *websocket.Conn
	Result    chan []Connection
}

type Connection struct {
	UUID      uuid.UUID
	WebSocket *websocket.Conn
}

type Action int64

const (
	Register Action = iota
	Get
	Delete
)

func ConnectionManagerWorker(jobs chan ConnectionJob) {
	var connections []Connection

	for job := range jobs {
		switch job.Action {
		case Register:
			connToInsert := Connection{
				UUID:      uuid.New(),
				WebSocket: job.WebSocket,
			}
			connections = append(connections, connToInsert)
			job.Result <- []Connection{connToInsert}
		case Get:
			job.Result <- connections
		case Delete:
			for i, connection := range connections {
				if connection.WebSocket == job.WebSocket {
					connections = append(connections[:i], connections[i+1:]...)
				}
			}
		}
	}
}
