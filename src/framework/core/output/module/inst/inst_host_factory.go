package inst

import (
	"configcenter/src/framework/core/output/module/model"
	"configcenter/src/framework/core/types"
)

func createHost(target model.Model) (Inst, error) {
	return nil, nil
}

// findHostsLikeName find all insts by inst name
func findHostsLikeName(target model.Model, businessName string) (Iterator, error) {
	// TODO:按照名字读取特定模型的实例集合，实例名字要模糊匹配
	return nil, nil
}

// findHostsByCondition find all insts by condition
func findHostsByCondition(target model.Model, condition types.MapStr) (Iterator, error) {
	// TODO:按照条件读取所有实例
	return nil, nil
}
