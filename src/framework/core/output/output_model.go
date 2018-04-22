package output

import (
	"configcenter/src/framework/core/output/module/model"
	"configcenter/src/framework/core/types"
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
func (cli *manager) FindClassificationsByCondition(condition types.MapStr) (model.ClassificationIterator, error) {
	return model.FindClassificationsByCondition(condition)
}
