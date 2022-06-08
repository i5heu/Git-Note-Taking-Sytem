package channels

import (
	files "base/Files"
	queue "base/Queue"
	registry "base/Registry"
)

var RegistryChan = make(chan registry.PoolJob, 50)
var QueueChan = make(chan queue.QueueJob, 50)
var FileChan = make(chan files.FileJob, 250)
