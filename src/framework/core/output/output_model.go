package output

import (
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/output/module/model"
)

// CreateClassification create a new classification
func (cli *manager) CreateClassification() model.Classification {
	return model.CreateClassification()
}

// FindClassificationsLikeName find a array of the classification by the name
func (cli *manager) FindClassificationsLikeName(name string) (model.ClassificationIterator, error) {
	return model.FindClassificationsLikeName(name)
}

// FindClassificationsByCondition find a array of the classification by the condition
func (cli *manager) FindClassificationsByCondition(condition *common.Condition) (model.ClassificationIterator, error) {
	return model.FindClassificationsByCondition(condition)
}
