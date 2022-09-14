package files

import (
	connectionManager "Tyche/ConnectionManager"
	webSocketGeneralTypes "Tyche/WebSocketGeneralTypes"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

func WriteFile(threadID uuid.UUID, message json.RawMessage, con connectionManager.Connection) {
	if !CheckIfUUIDIsInWriteAuthorityQueue(con.UUID) {
		webSocketGeneralTypes.SendError(con, threadID, errors.New("no write authority requested"))
		return
	}

	if !CheckIfUUIDIsWriteAuthority(con.UUID) {
		webSocketGeneralTypes.SendError(con, threadID, errors.New("no write authority pls wait"))
		return
	}
}
