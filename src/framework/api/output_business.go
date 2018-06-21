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
)

// CreateBusiness create a new Business object
func CreateBusiness(supplierAccount string) (*BusinessWrapper, error) {

	targetModel, err := GetBusinessModel(supplierAccount)
	if nil != err {
		return nil, err
	}
	businessInst := mgr.OutputerMgr.InstOperation().CreateBusinessInst(targetModel)
	wrapper := &BusinessWrapper{business: businessInst}
	//wrapper.SetSupplierAccount(supplierAccount)
	return wrapper, err
}

// FindBusinessLikeName find all insts by the name
func FindBusinessLikeName(supplierAccount, businessName string) (*BusinessIteratorWrapper, error) {
	targetModel, err := GetBusinessModel(supplierAccount)
	if nil != err {
		return nil, err
	}

	iter, err := mgr.OutputerMgr.InstOperation().FindBusinessLikeName(targetModel, businessName)
	return &BusinessIteratorWrapper{business: iter}, err
}

// FindBusinessByCondition find all insts by the condition
func FindBusinessByCondition(supplierAccount string, cond common.Condition) (*BusinessIteratorWrapper, error) {
	targetModel, err := GetBusinessModel(supplierAccount)
	if nil != err {
		return nil, err
	}
	iter, err := mgr.OutputerMgr.InstOperation().FindBusinessByCondition(targetModel, cond)
	return &BusinessIteratorWrapper{business: iter}, err
}
