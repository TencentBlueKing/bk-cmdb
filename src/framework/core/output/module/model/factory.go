package model

import (
	"configcenter/src/framework/core/types"
)

// CreateClassification create a new Classification instance
func CreateClassification() Classification {
	return &classification{}
}

// FindClassificationsLikeName find a array of the classification by the name
func FindClassificationsLikeName(name string) (ClassificationIterator, error) {
	// TODO: 按照名字模糊查找
	return nil, nil
}

// FindClassificationsByCondition find a array of the classification by the condition
func FindClassificationsByCondition(condition types.MapStr) (ClassificationIterator, error) {
	// TODO: 按照条件搜索
	return nil, nil
}
