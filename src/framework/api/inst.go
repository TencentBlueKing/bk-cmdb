package api

import (
	"configcenter/src/framework/core/output/module/inst"
	"configcenter/src/framework/core/output/module/model"
)

// CreateBusiness create a new Business object
func CreateBusiness(supplierAccount string) (inst.Inst, error) {

	// TODO:需要根据supplierAccount 获取业务模型
	return mgr.OutputerMgr.CreateInst(nil)
}

// CreateSet create a new set object
func CreateSet() (inst.Inst, error) {
	// TODO:需要根据supplierAccount 获取集群模型定义
	return mgr.OutputerMgr.CreateInst(nil)
}

// CreateModule create a new module object
func CreateModule() (inst.Inst, error) {
	// TODO:需要根据supplierAccount 获取模块的定义
	return mgr.OutputerMgr.CreateInst(nil)
}

// CreateCommonInst create a common inst object
func CreateCommonInst(target model.Model) (inst.Inst, error) {
	// TODO:根据model 创建普通实例
	return mgr.OutputerMgr.CreateInst(nil)
}
