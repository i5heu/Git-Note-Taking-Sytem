package channels

import (
	connectionManager "Tyche/ConnectionManager"
	files "Tyche/Files"
)

var ConnectionManagerChan = make(chan connectionManager.ConnectionJob, 50)
var WriteAuthorityQueueChan = make(chan files.WriteAuthorityRequest, 50)
