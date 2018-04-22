package inst

import (
	"configcenter/src/framework/core/output/module/model"
	"configcenter/src/framework/core/types"
)

func createPlat(target model.Model) (Inst, error) {
	return nil, nil
}

// findPlatsLikeName find all insts by inst name
func findPlatsLikeName(target model.Model, businessName string) (Iterator, error) {
	// TODO:按照名字读取特定模型的实例集合，实例名字要模糊匹配
	return nil, nil
}

// findPlatsByCondition find all insts by condition
func findPlatsByCondition(target model.Model, condition types.MapStr) (Iterator, error) {
	// TODO:按照条件读取所有实例
	return nil, nil
}
