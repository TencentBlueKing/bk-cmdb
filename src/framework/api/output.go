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

	cond := CreateCondition().Field(model.ClassificationID).Eq("bk_organization")
	clsIter, err := mgr.OutputerMgr.FindClassificationsByCondition(cond)
	if nil != err {
		return nil, err
	}
	var targetModel model.Model
	err = clsIter.ForEach(func(item model.Classification) error {

		cond = CreateCondition().Field(model.ObjectID).Eq("biz").
			Field(model.SupplierAccount).Eq(supplierAccount).
			Field(model.ClassificationID).Eq(item.GetID())

		modelIter, err := item.FindModelsByCondition(cond)
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
	businessInst, err := mgr.OutputerMgr.CreateInst(targetModel)
	wrapper := &BusinessWrapper{business: businessInst}
	wrapper.SetSupplierAccount(supplierAccount)
	return wrapper, err
}

// CreateSet create a new set object
func CreateSet(supplierAccount string) (*SetWrapper, error) {

	cond := CreateCondition().Field(model.ClassificationID).Eq("bk_biz_topo")
	clsIter, err := mgr.OutputerMgr.FindClassificationsByCondition(cond)
	if nil != err {
		return nil, err
	}
	var targetModel model.Model
	err = clsIter.ForEach(func(item model.Classification) error {

		cond = CreateCondition().Field(model.ObjectID).Eq("set").
			Field(model.SupplierAccount).Eq(supplierAccount).
			Field(model.ClassificationID).Eq(item.GetID())

		modelIter, err := item.FindModelsByCondition(cond)
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

	setInst, err := mgr.OutputerMgr.CreateInst(targetModel)

	return &SetWrapper{set: setInst}, err
}

// CreateModule create a new module object
func CreateModule(supplierAccount string) (*ModuleWrapper, error) {

	cond := CreateCondition().Field(model.ClassificationID).Eq("bk_biz_topo")
	clsIter, err := mgr.OutputerMgr.FindClassificationsByCondition(cond)
	if nil != err {
		return nil, err
	}
	var targetModel model.Model
	err = clsIter.ForEach(func(item model.Classification) error {

		cond = CreateCondition().Field(model.ObjectID).Eq("module").
			Field(model.SupplierAccount).Eq(supplierAccount).
			Field(model.ClassificationID).Eq(item.GetID())

		modelIter, err := item.FindModelsByCondition(cond)
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
	moduleInst, err := mgr.OutputerMgr.CreateInst(targetModel)
	return &ModuleWrapper{module: moduleInst}, err
}

// CreateHost create a new host object
func CreateHost(supplierAccount string) (*HostWrapper, error) {

	cond := CreateCondition().Field(model.ClassificationID).Eq("bk_host_manage")
	clsIter, err := mgr.OutputerMgr.FindClassificationsByCondition(cond)
	if nil != err {
		return nil, err
	}
	var targetModel model.Model
	err = clsIter.ForEach(func(item model.Classification) error {

		cond = CreateCondition().Field(model.ObjectID).Eq("host").
			Field(model.SupplierAccount).Eq(supplierAccount).
			Field(model.ClassificationID).Eq(item.GetID())

		modelIter, err := item.FindModelsByCondition(cond)
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

	hostInst, err := mgr.OutputerMgr.CreateInst(targetModel)
	return &HostWrapper{host: hostInst}, err
}

// CreateCommonInst create a common inst object
func CreateCommonInst(target model.Model) (inst.Inst, error) {
	return mgr.OutputerMgr.CreateInst(target)
}
