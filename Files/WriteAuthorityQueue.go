package files

import (
	connectionManager "Tyche/ConnectionManager"
	"time"

	"github.com/google/uuid"
)

var WriteAuthorityQueue = []WriteAuthorityRequest{}
var writeAuthority = uuid.Nil

type WriteAuthorityRequest struct {
	ThreadID   uuid.UUID
	Connection connectionManager.Connection
	BackRef    chan WriteAuthorityResponse
}

type WriteAuthorityResponse struct {
	ThreadID    uuid.UUID `json:"threadID"`
	MessageType string    `json:"type"`
	Status      string    `json:"status"`
}

func WriteAuthorityQueueWorker(channel chan WriteAuthorityRequest) {
	for request := range channel {
		WriteAuthorityQueue = append(WriteAuthorityQueue, request)

		request.Connection.WebSocket.WriteJSON(WriteAuthorityResponse{
			ThreadID:    request.ThreadID,
			MessageType: "WRITE_AUTHORITY",
			Status:      "WAITING",
		})
	}
}

func WriteAuthorityQueueWorkerAuthority() {
	for range time.Tick(time.Millisecond * 5) {
		if len(WriteAuthorityQueue) > 0 && writeAuthority == uuid.Nil {
			newWriteAuthority := WriteAuthorityQueue[0]
			writeAuthority = newWriteAuthority.Connection.UUID

			newWriteAuthority.Connection.WebSocket.WriteJSON(WriteAuthorityResponse{
				ThreadID:    newWriteAuthority.ThreadID,
				MessageType: "WRITE_AUTHORITY",
				Status:      "GRANTED",
			})

			newWriteAuthority.BackRef <- WriteAuthorityResponse{
				ThreadID:    newWriteAuthority.ThreadID,
				MessageType: "WRITE_AUTHORITY",
				Status:      "GRANTED",
			}
		}
	}
}

func CheckIfUUIDIsInWriteAuthorityQueue(uuid uuid.UUID) bool {
	for _, request := range WriteAuthorityQueue {
		if request.Connection.UUID == uuid {
			return true
		}
	}

	return false
}

func CheckIfUUIDIsWriteAuthority(uuid uuid.UUID) bool {
	return writeAuthority == uuid
}
