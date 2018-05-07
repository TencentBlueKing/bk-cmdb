/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package api

import (
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/output/module/inst"
	"configcenter/src/framework/core/output/module/model"
)

func getModel(supplierAccount, classificationID, objID string) (model.Model, error) {
	condInner := CreateCondition().Field(model.ClassificationID).Eq(classificationID)
	clsIter, err := mgr.OutputerMgr.FindClassificationsByCondition(condInner)
	if nil != err {
		return nil, err
	}
	var targetModel model.Model
	err = clsIter.ForEach(func(item model.Classification) error {

		condInner = CreateCondition().Field(model.ObjectID).Eq(objID).
			Field(model.SupplierAccount).Eq(supplierAccount).
			Field(model.ClassificationID).Eq(item.GetID())

		modelIter, err := item.FindModelsByCondition(condInner)
		if nil != err {
			return err
		}

		err = modelIter.ForEach(func(modelItem model.Model) error {
			targetModel = modelItem
			return nil
		})

		return nil
	})

	if nil != err {
		return nil, err
	}

	return targetModel, err
}

// CreateClassification create a new classification
func CreateClassification(name string) model.Classification {
	return mgr.OutputerMgr.CreateClassification(name)
}

// FindClassificationsLikeName find a array of the classification by the name
func FindClassificationsLikeName(name string) (model.ClassificationIterator, error) {
	return mgr.OutputerMgr.FindClassificationsLikeName(name)
}

// FindClassificationsByCondition find a array of the classification by the condition
func FindClassificationsByCondition(condition common.Condition) (model.ClassificationIterator, error) {
	return mgr.OutputerMgr.FindClassificationsByCondition(condition)
}

// CreateBusiness create a new Business object
func CreateBusiness(supplierAccount string) (*BusinessWrapper, error) {

	targetModel, err := getModel(supplierAccount, "bk_organization", "biz")
	if nil != err {
		return nil, err
	}
	businessInst, err := mgr.OutputerMgr.CreateInst(targetModel)
	wrapper := &BusinessWrapper{business: businessInst}
	wrapper.SetSupplierAccount(supplierAccount)
	return wrapper, err
}

// FindBusinessLikeName find all insts by the name
func FindBusinessLikeName(supplierAccount, businessName string) (inst.Iterator, error) {
	targetModel, err := getModel(supplierAccount, "bk_organization", "biz")
	if nil != err {
		return nil, err
	}

	return mgr.OutputerMgr.FindInstsLikeName(targetModel, businessName)
}

// FindBusinessByCondition find all insts by the condition
func FindBusinessByCondition(supplierAccount string, cond common.Condition) (inst.Iterator, error) {
	targetModel, err := getModel(supplierAccount, "bk_organization", "biz")
	if nil != err {
		return nil, err
	}
	return mgr.OutputerMgr.FindInstsByCondition(targetModel, cond)
}

// CreateSet create a new set object
func CreateSet(supplierAccount string) (*SetWrapper, error) {

	targetModel, err := getModel(supplierAccount, "bk_biz_topo", "set")
	if nil != err {
		return nil, err
	}

	setInst, err := mgr.OutputerMgr.CreateInst(targetModel)

	return &SetWrapper{set: setInst}, err
}

// FindSetLikeName find all insts by the name
func FindSetLikeName(supplierAccount, setName string) (inst.Iterator, error) {
	targetModel, err := getModel(supplierAccount, "bk_biz_topo", "set")
	if nil != err {
		return nil, err
	}

	return mgr.OutputerMgr.FindInstsLikeName(targetModel, setName)
}

// FindSetByCondition find all insts by the condition
func FindSetByCondition(supplierAccount string, cond common.Condition) (inst.Iterator, error) {
	targetModel, err := getModel(supplierAccount, "bk_biz_topo", "set")
	if nil != err {
		return nil, err
	}
	return mgr.OutputerMgr.FindInstsByCondition(targetModel, cond)
}

// CreateModule create a new module object
func CreateModule(supplierAccount string) (*ModuleWrapper, error) {

	targetModel, err := getModel(supplierAccount, "bk_biz_topo", "module")
	if nil != err {
		return nil, err
	}

	moduleInst, err := mgr.OutputerMgr.CreateInst(targetModel)
	return &ModuleWrapper{module: moduleInst}, err
}

// FindModuleLikeName find all insts by the name
func FindModuleLikeName(supplierAccount, moduleName string) (inst.Iterator, error) {
	targetModel, err := getModel(supplierAccount, "bk_biz_topo", "module")
	if nil != err {
		return nil, err
	}

	return mgr.OutputerMgr.FindInstsLikeName(targetModel, moduleName)
}

// FindModuleByCondition find all insts by the condition
func FindModuleByCondition(supplierAccount string, cond common.Condition) (inst.Iterator, error) {
	targetModel, err := getModel(supplierAccount, "bk_biz_topo", "module")
	if nil != err {
		return nil, err
	}
	return mgr.OutputerMgr.FindInstsByCondition(targetModel, cond)
}

// CreateHost create a new host object
func CreateHost(supplierAccount string) (*HostWrapper, error) {
	targetModel, err := getModel(supplierAccount, "bk_host_manage", "host")
	if nil != err {
		return nil, err
	}

	hostInst, err := mgr.OutputerMgr.CreateInst(targetModel)
	return &HostWrapper{host: hostInst}, err
}

// FindHostLikeName find all insts by the name
func FindHostLikeName(supplierAccount, hostName string) (inst.Iterator, error) {
	targetModel, err := getModel(supplierAccount, "bk_host_manage", "host")
	if nil != err {
		return nil, err
	}
	return mgr.OutputerMgr.FindInstsLikeName(targetModel, hostName)
}

// FindHostByCondition find all insts by the condition
func FindHostByCondition(supplierAccount string, cond common.Condition) (inst.Iterator, error) {
	targetModel, err := getModel(supplierAccount, "bk_host_manage", "host")
	if nil != err {
		return nil, err
	}
	return mgr.OutputerMgr.FindInstsByCondition(targetModel, cond)
}

// CreateCommonInst create a common inst object
func CreateCommonInst(target model.Model) (inst.Inst, error) {
	return mgr.OutputerMgr.CreateInst(target)
}

// FindInstsLikeName find all insts by the name
func FindInstsLikeName(target model.Model, instName string) (inst.Iterator, error) {
	return mgr.OutputerMgr.FindInstsLikeName(target, instName)
}

// FindInstsByCondition find all insts by the condition
func FindInstsByCondition(target model.Model, cond common.Condition) (inst.Iterator, error) {
	return mgr.OutputerMgr.FindInstsByCondition(target, cond)
}
