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

import "configcenter/src/framework/common"

// CreateHost create a new host object
func CreateHost(supplierAccount string) (*HostWrapper, error) {
	targetModel, err := GetModel(supplierAccount, "bk_host_manage", "host")
	if nil != err {
		return nil, err
	}

	hostInst := mgr.OutputerMgr.InstOperation().CreateHostInst(targetModel)
	return &HostWrapper{
		supplierAccount: supplierAccount,
		host:            hostInst}, err
}

// FindHostLikeName find all insts by the name
func FindHostLikeName(supplierAccount, hostName string) (*HostIteratorWrapper, error) {
	targetModel, err := GetModel(supplierAccount, "bk_host_manage", "host")
	if nil != err {
		return nil, err
	}
	iter, err := mgr.OutputerMgr.InstOperation().FindHostsLikeName(targetModel, hostName)
	return &HostIteratorWrapper{host: iter}, err
}

// FindHostByCondition find all insts by the condition
func FindHostByCondition(supplierAccount string, cond common.Condition) (*HostIteratorWrapper, error) {
	targetModel, err := GetModel(supplierAccount, "bk_host_manage", "host")
	if nil != err {
		return nil, err
	}
	//fmt.Println("the host model:", targetModel)
	iter, err := mgr.OutputerMgr.InstOperation().FindHostsByCondition(targetModel, cond)
	return &HostIteratorWrapper{host: iter}, err
}
