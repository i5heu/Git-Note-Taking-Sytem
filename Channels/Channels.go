package channels

import connectionManager "Tyche/ConnectionManager"

var ConnectionManagerChan = make(chan connectionManager.ConnectionJob, 50)
