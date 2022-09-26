package files

import (
	connectionManager "Tyche/ConnectionManager"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const (
	writeAuthorityExpireTime = 10 * time.Second
)

type WriteAuthorityAction string

const (
	Request WriteAuthorityAction = "REQUEST"
	Revoke  WriteAuthorityAction = "REVOKE"
)

type WriteAuthorityQueue struct {
	time    time.Time
	request WriteAuthorityRequest
}

var writeAuthorityQueue = []WriteAuthorityQueue{}
var writeAuthority = uuid.Nil

type WriteAuthorityRequest struct {
	ThreadID   uuid.UUID
	Data       json.RawMessage
	Connection connectionManager.Connection
	BackRef    chan WriteAuthorityResponse
}

type Data struct {
	Action string `json:"action"`
}

type WriteAuthorityResponse struct {
	ThreadID    uuid.UUID `json:"thread_id"`
	MessageType string    `json:"type"`
	Status      string    `json:"status"`
}

func WriteAuthorityQueueWorker(channel chan WriteAuthorityRequest) {
	for request := range channel {
		data := Data{}
		err := json.Unmarshal(request.Data, &data)
		if err != nil {
			request.BackRef <- WriteAuthorityResponse{
				ThreadID:    request.ThreadID,
				MessageType: "WRITE_AUTHORITY",
				Status:      "ERROR",
			}
			log.Error().Err(err).Msg("Failed to unmarshal data")
			close(request.BackRef)
			continue
		}

		if data.Action == string(Request) {
			writeAuthorityQueue = append(writeAuthorityQueue, WriteAuthorityQueue{
				time:    time.Now(),
				request: request,
			})

			request.BackRef <- WriteAuthorityResponse{
				ThreadID:    request.ThreadID,
				MessageType: "WRITE_AUTHORITY",
				Status:      "WAITING",
			}

			// remove from queue
			for i, request := range writeAuthorityQueue {
				if request.request.Connection.UUID == writeAuthority {
					writeAuthorityQueue = append(writeAuthorityQueue[:i], writeAuthorityQueue[i+1:]...)
					break
				}
			}
			continue
		} else if data.Action == string(Revoke) {
			if writeAuthority == request.Connection.UUID {
				writeAuthority = uuid.Nil
			}

			continue
		}

		request.BackRef <- WriteAuthorityResponse{
			ThreadID:    request.ThreadID,
			MessageType: "WRITE_AUTHORITY",
			Status:      "NO_ACTION_SELECTED",
		}
	}
}

func WriteAuthorityQueueBackgroundWorker() {
	for range time.Tick(time.Millisecond * 33) {
		// assert new write authority if no one has it
		if len(writeAuthorityQueue) > 0 && writeAuthority == uuid.Nil {
			newWriteAuthority := writeAuthorityQueue[0]
			writeAuthority = newWriteAuthority.request.Connection.UUID

			newWriteAuthority.request.BackRef <- WriteAuthorityResponse{
				ThreadID:    newWriteAuthority.request.ThreadID,
				MessageType: "WRITE_AUTHORITY",
				Status:      "GRANTED",
			}
			close(newWriteAuthority.request.BackRef)
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
