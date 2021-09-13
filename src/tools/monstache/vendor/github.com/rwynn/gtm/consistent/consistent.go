package consistent

import (
	"errors"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/rwynn/gtm"
	"github.com/serialx/hashring"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConfigOptions struct {
	Workers []string
}

var EmptyWorkers = errors.New("config not found or workers empty")
var InvalidWorkers = errors.New("workers must be an array of string")
var WorkerMissing = errors.New("the specified worker was not found in the config")

// returns an operation filter which uses a consistent hash to determine
// if the operation will be accepted for processing. can be used to distribute work.
// name:	the name of the worker creating this filter. e.g. "Harry"
// configFile:	a file path to a TOML document.  the document should contain
// a property named 'Workers' which is a list of all the workers participating. e.g.
// workers = [ "Tom", "Dick", "Harry" ]
func ConsistentHashFilterFromFile(name string, configFile string) (gtm.OpFilter, error) {
	var config ConfigOptions
	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		return nil, EmptyWorkers
	} else {
		return ConsistentHashFilter(name, config.Workers)
	}
}

// returns an operation filter which uses a consistent hash to determine
// if the operation will be accepted for processing. can be used to distribute work.
// name:	the name of the worker creating this filter. e.g. "Harry"
// document:	a map with a string key 'workers' which has a corresponding
//				slice of string representing the available workers
func ConsistentHashFilterFromDocument(name string, document map[string]interface{}) (gtm.OpFilter, error) {
	workers := document["workers"]
	return ConsistentHashFilter(name, workers.([]string))
}

// returns an operation filter which uses a consistent hash to determine
// if the operation will be accepted for processing. can be used to distribute work.
// name:	the name of the worker creating this filter. e.g. "Harry"
// workers:	a slice of strings representing the available worker names
func ConsistentHashFilter(name string, workers []string) (gtm.OpFilter, error) {
	if len(workers) == 0 {
		return nil, EmptyWorkers
	}
	found := false
	ring := hashring.New(workers)
	for _, worker := range workers {
		if worker == name {
			found = true
		}
	}
	if !found {
		return nil, WorkerMissing
	}
	return func(op *gtm.Op) bool {
		if op.Id != nil {
			var idStr string
			switch op.Id.(type) {
			case primitive.ObjectID:
				idStr = op.Id.(primitive.ObjectID).Hex()
			default:
				idStr = fmt.Sprintf("%v", op.Id)
			}
			who, ok := ring.GetNode(idStr)
			if ok {
				return name == who
			} else {
				return false
			}
		} else {
			return true
		}
	}, nil
}
