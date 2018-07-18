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

package util

import (
	"configcenter/src/common"
	"configcenter/src/common/metadata"
)

func NewOperation() *operation {
	return &operation{
		op: make(map[string]interface{}),
	}
}

type operation struct {
	op map[string]interface{}
}

func (o *operation) Data() map[string]interface{} {
	return o.op
}

func (o *operation) WithHostID(hostID int64) *operation {
	o.op[common.BKHostIDField] = hostID
	return o
}

func (o *operation) WithAppID(appID int64) *operation {
	o.op[common.BKAppIDField] = appID
	return o
}

func (o *operation) WithSupplierID(field int64) *operation {
	o.op[common.BKSupplierIDField] = field
	return o
}

func (o *operation) WithOwnerID(ownerID string) *operation {
	o.op[common.BKOwnerIDField] = ownerID
	return o
}

func (o *operation) WithDefaultField(d int64) *operation {
	o.op[common.BKDefaultField] = d
	return o
}

func (o *operation) WithInstID(instID int64) *operation {
	o.op[common.BKInstIDField] = instID
	return o
}

func (o *operation) WithObjID(objID string) *operation {
	o.op[common.BKObjIDField] = objID
	return o
}

func (o *operation) WithPropertyID(id string) *operation {
	o.op[common.BKObjAttIDField] = id
	return o
}

func (o *operation) WithModuleName(name string) *operation {
	o.op[common.BKModuleNameField] = name
	return o
}

func (o *operation) WithModuleIDs(id []int64) *operation {
	o.op[common.BKModuleIDField] = id
	return o
}

func (o *operation) WithModuleID(id int64) *operation {
	o.op[common.BKModuleIDField] = id
	return o
}
func (o *operation) WithPage(p metadata.BasePage) *operation {
	o.op[metadata.PageName] = p
	return o
}

func (o *operation) WithAssoObjID(id string) *operation {
	o.op[common.BKAsstObjIDField] = id
	return o
}

func (o *operation) WithAssoInstID(id map[string]interface{}) *operation {
	o.op[common.BKAsstInstIDField] = id
	return o
}

func (o *operation) WithHostInnerIP(ip string) *operation {
	o.op[common.BKHostInnerIPField] = ip
	return o
}

func (o *operation) WithCloudID(id int64) *operation {
	o.op[common.BKCloudIDField] = id
	return o
}
