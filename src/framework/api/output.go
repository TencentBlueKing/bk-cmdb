package api

import (
	"configcenter/src/framework/core/output"
	"configcenter/src/framework/core/output/module/model"
	"configcenter/src/framework/core/types"
	"errors"
)

// CreateCustomOutputer create a new custom outputer
func CreateCustomOutputer(name string, runFunc func(data types.MapStr) error) (output.OutputerKey, output.Puter, error) {

	if 0 == len(name) {
		return output.OutputerKey(""), nil, errors.New("the name parmeter must be set")
	}

	if nil == runFunc {
		return output.OutputerKey(""), nil, errors.New("the run function must be set")
	}

	key, sender := mgr.OutputerMgr.CreateCustomOutputer(name, runFunc)
	return key, sender, nil
}

// CreateClassification create a new classification
func CreateClassification() model.Classification {
	return mgr.OutputerMgr.CreateClassification()
}

// FindClassificationsLikeName find a array of the classification by the name
func FindClassificationsLikeName(name string) (model.ClassificationIterator, error) {
	return mgr.OutputerMgr.FindClassificationsLikeName(name)
}

// FindClassificationsByCondition find a array of the classification by the condition
func FindClassificationsByCondition(condition types.MapStr) (model.ClassificationIterator, error) {
	return mgr.OutputerMgr.FindClassificationsByCondition(condition)
}
