package output

import (
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/output/module/inst"
	"configcenter/src/framework/core/output/module/model"
)

// CreateInst create a instance for the model
func (cli *manager) CreateInst(target model.Model) (inst.Inst, error) {
	return inst.CreateInst(target)
}

// FindInstsLikeName find all insts by the name
func (cli *manager) FindInstsLikeName(target model.Model, instName string) (inst.Iterator, error) {
	return inst.FindInstsLikeName(target, instName)
}

// FindInstsByCondition find all insts by the condition
func (cli *manager) FindInstsByCondition(target model.Model, condition *common.Condition) (inst.Iterator, error) {
	return inst.FindInstsByCondition(target, condition)
}
