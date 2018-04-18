package api

import (
	"configcenter/src/framework/core/output"
	"errors"
)

// CreateCustomOutputer create a new custom outputer
func CreateCustomOutputer(name string, runFunc func(data output.MapStr) error) (output.OutputerKey, output.Puter, error) {

	if 0 == len(name) {
		return output.OutputerKey(""), nil, errors.New("the name parmeter must be set")
	}

	if nil == runFunc {
		return output.OutputerKey(""), nil, errors.New("the run function must be set")
	}

	key, sender := mgr.OutputerMgr.CreateCustomOutputer(name, runFunc)
	return key, sender, nil
}
