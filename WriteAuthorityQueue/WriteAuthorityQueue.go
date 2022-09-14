package writeauthorityqueue

import connectionManager "Tyche/ConnectionManager"

var WriteAuthorityRequestIncoming = make(chan WriteAuthorityRequest, 50)

type WriteAuthorityRequest struct {
	connection connectionManager.ConnectionJob
}
