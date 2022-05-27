package channels

import (
	queue "base/Queue"
	registry "base/Registry"
)

var RegistryChan = make(chan registry.PoolJob, 50)
var QueueChan = make(chan queue.QueueJob, 50)
