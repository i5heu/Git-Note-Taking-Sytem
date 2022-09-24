package files

import (
	connectionManager "Tyche/ConnectionManager"
	"time"

	"github.com/google/uuid"
)

const (
	writeAuthorityExpireTime = 10 * time.Second
)

type WriteAuthorityQueue struct {
	time    time.Time
	request WriteAuthorityRequest
}

var writeAuthorityQueue = []WriteAuthorityQueue{}
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
		writeAuthorityQueue = append(writeAuthorityQueue, WriteAuthorityQueue{
			time:    time.Now(),
			request: request,
		})

		request.Connection.WebSocket.WriteJSON(WriteAuthorityResponse{
			ThreadID:    request.ThreadID,
			MessageType: "WRITE_AUTHORITY",
			Status:      "WAITING",
		})
	}
}

func WriteAuthorityQueueBackgroundWorker() {
	for range time.Tick(time.Millisecond * 33) {
		// assert new write authority if no one has it
		if len(writeAuthorityQueue) > 0 && writeAuthority == uuid.Nil {
			newWriteAuthority := writeAuthorityQueue[0]
			writeAuthority = newWriteAuthority.request.Connection.UUID

			newWriteAuthority.request.Connection.WebSocket.WriteJSON(WriteAuthorityResponse{
				ThreadID:    newWriteAuthority.request.ThreadID,
				MessageType: "WRITE_AUTHORITY",
				Status:      "GRANTED",
			})

			newWriteAuthority.request.BackRef <- WriteAuthorityResponse{
				ThreadID:    newWriteAuthority.request.ThreadID,
				MessageType: "WRITE_AUTHORITY",
				Status:      "GRANTED",
			}
		}

		// remove write authority if it has expired
		if writeAuthority != uuid.Nil {
			for i, request := range writeAuthorityQueue {
				if request.request.Connection.UUID == writeAuthority {
					if time.Since(request.time) > writeAuthorityExpireTime {
						writeAuthority = uuid.Nil
						writeAuthorityQueue = append(writeAuthorityQueue[:i], writeAuthorityQueue[i+1:]...)
					}
					break
				}
			}
		}
	}
}

func CheckIfUUIDIsInWriteAuthorityQueue(uuid uuid.UUID) bool {
	for _, request := range writeAuthorityQueue {
		if request.request.Connection.UUID == uuid {
			return true
		}
	}

	return false
}

func CheckIfUUIDIsWriteAuthority(uuid uuid.UUID) bool {
	return writeAuthority == uuid
}
