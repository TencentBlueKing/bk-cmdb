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
	"io"
	"strings"

	"configcenter/src/framework/common"
)

// CreatePlat create a new plat object
func CreatePlat(supplierAccount string) (*PlatWrapper, error) {
	targetModel, err := GetModel(supplierAccount, "bk_host_manage", "plat")
	if nil != err {
		return nil, err
	}

	platInst := mgr.OutputerMgr.InstOperation().CreatePlatInst(targetModel)
	platInst.SetValue(fieldSupplierAccount, supplierAccount)
	platInst.SetValue(fieldObjectID, plat)
	return &PlatWrapper{plat: platInst}, err
}

// GetPlatID get the plat id
func GetPlatID(supplierAccount, platName string) (int64, error) {

	cond := CreateCondition().Field(fieldSupplierAccount).Eq(supplierAccount).Field(fieldPlatName).In([]string{platName})
	iter, iterErr := FindPlatByCondition(supplierAccount, cond)
	if nil != iterErr {
		if strings.Contains(iterErr.Error(), "empty") {
			return -1, io.EOF
		}
		return -1, iterErr
	}

	platID := -1
	err := iter.ForEach(func(plat *PlatWrapper) error {
		id, err := plat.GetID()
		if nil != err {
			return err
		}
		platID = id
		return nil
	})

	if nil != err {
		return -1, err
	}

	if -1 == platID {
		return -1, io.EOF
	}
	return int64(platID), nil
}

// FindPlatLikeName find all insts by the name
func FindPlatLikeName(supplierAccount, platName string) (*PlatIteratorWrapper, error) {
	targetModel, err := GetModel(supplierAccount, "bk_host_manage", "plat")
	if nil != err {
		return nil, err
	}
	iter, err := mgr.OutputerMgr.InstOperation().FindPlatsLikeName(targetModel, platName)
	return &PlatIteratorWrapper{plat: iter}, err
}

// FindPlatByCondition find all insts by the condition
func FindPlatByCondition(supplierAccount string, cond common.Condition) (*PlatIteratorWrapper, error) {
	targetModel, err := GetModel(supplierAccount, "bk_host_manage", "plat")
	if nil != err {
		return nil, err
	}
	iter, err := mgr.OutputerMgr.InstOperation().FindPlatsByCondition(targetModel, cond)
	return &PlatIteratorWrapper{plat: iter}, err
}
