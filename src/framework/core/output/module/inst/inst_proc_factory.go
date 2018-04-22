package inst

import (
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/output/module/model"
)

func createProc(target model.Model) (Inst, error) {
	return nil, nil
}

// findProcsLikeName find all insts by inst name
func findProcsLikeName(target model.Model, businessName string) (Iterator, error) {
	// TODO:按照名字读取特定模型的实例集合，实例名字要模糊匹配
	return nil, nil
}

// findProcsByCondition find all insts by condition
func findProcsByCondition(target model.Model, condition *common.Condition) (Iterator, error) {
	// TODO:按照条件读取所有实例
	return nil, nil
}
