package output

import (
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/log"
	"configcenter/src/framework/core/output/module/model"
	"configcenter/src/framework/core/types"
	"sync"
)

// MapOutputer the outputer storage
type MapOutputer map[OutputerKey]Outputer

// output implements the Outputer interface
type manager struct {
	outputerLock sync.RWMutex
	outputers    MapOutputer
}

// CreateClassification create a new classification
func (cli *manager) CreateClassification() model.Classification {
	return model.CreateClassification()
}

// FindClassificationsLikeName find a array of the classification by the name
func (cli *manager) FindClassificationsLikeName(name string) (model.ClassificationIterator, error) {
	return model.FindClassificationsLikeName(name)
}

// FindClassificationsByCondition find a array of the classification by the condition
func (cli *manager) FindClassificationsByCondition(condition types.MapStr) (model.ClassificationIterator, error) {
	return model.FindClassificationsByCondition(condition)
}

/** the following  methods are used to maintence the custom outputer */

func (cli *manager) AddOutputer(target Outputer) OutputerKey {

	cli.outputerLock.Lock()

	key := OutputerKey(common.UUID())
	cli.outputers[key] = target

	cli.outputerLock.Unlock()

	return key
}

func (cli *manager) RemoveOutputer(key OutputerKey) {

	cli.outputerLock.Lock()

	if item, ok := cli.outputers[key]; ok {
		if err := item.Stop(); nil != err {
			log.Errorf("failed to stop the outputer (%s), stop to reove it, error info is %s", item.Name(), err.Error())
		} else {
			log.Infof("remove the outputer(%s)", item.Name())
			delete(cli.outputers, key)
		}
	}

	cli.outputerLock.Unlock()
}
func (cli *manager) FetchOutputer(key OutputerKey) Puter {

	cli.outputerLock.RLock()
	defer func() {
		cli.outputerLock.RUnlock()
	}()

	if item, ok := cli.outputers[key]; ok {
		return item
	}

	return nil
}
func (cli *manager) CreateCustomOutputer(name string, run func(data types.MapStr) error) (OutputerKey, Puter) {

	log.Infof("creater custom outputer:%s", name)
	wrapper := &customWrapper{
		name:    name,
		runFunc: run,
	}

	key := cli.AddOutputer(wrapper)

	return key, wrapper
}
