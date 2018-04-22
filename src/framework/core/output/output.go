package output

import (
	"sync"
)

// MapOutputer the outputer storage
type MapOutputer map[OutputerKey]Outputer

// output implements the Outputer interface
type manager struct {
	outputerLock sync.RWMutex
	outputers    MapOutputer
}
