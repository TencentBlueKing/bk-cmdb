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
	"configcenter/src/framework/core/output/module/model"
)

// CreateClassification create a new classification
func CreateClassification(name string) model.Classification {
	return mgr.OutputerMgr.CreateClassification(name)
}

// FindClassificationsLikeName find a array of the classification by the name
func FindClassificationsLikeName(name string) (model.ClassificationIterator, error) {
	return mgr.OutputerMgr.FindClassificationsLikeName("", name)
}

// FindClassificationsLikeNameWithOwner find a array of the classification by the name
func FindClassificationsLikeNameWithOwner(supplierAccount, name string) (model.ClassificationIterator, error) {
	return mgr.OutputerMgr.FindClassificationsLikeName(supplierAccount, name)
}

// FindClassificationsByCondition find a array of the classification by the condition
func FindClassificationsByCondition(cond common.Condition) (model.ClassificationIterator, error) {
	return mgr.OutputerMgr.FindClassificationsByCondition("", cond)
}

// FindClassificationsByConditionWithOwner find a array of the classification by the condition
func FindClassificationsByConditionWithOwner(supplierAccount string, cond common.Condition) (model.ClassificationIterator, error) {
	return mgr.OutputerMgr.FindClassificationsByCondition(supplierAccount, cond)
}

// GetModel get the model
func GetModel(supplierAccount, classificationID, objID string) (model.Model, error) {
	return mgr.OutputerMgr.GetModel(supplierAccount, classificationID, objID)
}

// GetBusinessModel return a Business object
func GetBusinessModel(supplierAccount string) (model.Model, error) {
	return GetModel(supplierAccount, "bk_organization", "biz")
}

// GetSetModel get a set object
func GetSetModel(supplierAccount string) (model.Model, error) {
	return GetModel(supplierAccount, "bk_biz_topo", "set")
}

// GetModuleModel get a module object
func GetModuleModel(supplierAccount string) (model.Model, error) {
	return GetModel(supplierAccount, "bk_biz_topo", "module")
}

// GetHostModel get a host object
func GetHostModel(supplierAccount string) (model.Model, error) {
	return GetModel(supplierAccount, "bk_host_manage", "host")
}

// GetPlatModel return a plat object
func GetPlatModel(supplierAccount string) (model.Model, error) {
	return GetModel(supplierAccount, "bk_host_manage", "plat")
}
