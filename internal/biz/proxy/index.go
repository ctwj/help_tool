package proxy

import "sync"

var index uint64
var indexLock sync.Mutex

func GetRequestIndex() uint64 {
	indexLock.Lock()
	defer indexLock.Unlock()
	index += 1
	return index
}
