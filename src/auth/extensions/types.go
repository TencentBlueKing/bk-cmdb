/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package extensions

import (
	"configcenter/src/apimachinery"
	"configcenter/src/auth"
	"configcenter/src/common"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

type AuthManager struct {
	clientSet apimachinery.ClientSetInterface
	Authorize auth.Authorize
	// Err is used for return error messages of specific language on running
	Err errors.DefaultCCErrorIf
}

func NewAuthManager(clientSet apimachinery.ClientSetInterface, Authorize auth.Authorize) *AuthManager {
	return &AuthManager{
		clientSet: clientSet,
		Authorize: Authorize,
	}
}

type ClassificationSimplify struct {
	Name             string
	ID               int64
	ClassificationID string
}

type InstanceSimplify struct {
	ID         int64  `json:"id"`
	InstanceID string `json:"bk_inst_id"`
	Name       string `json:"bk_inst_name"`
	BizID     int64
}

// Parse load the data from mapstr attribute into ObjectUnique instance
func (is *InstanceSimplify) Parse(data mapstr.MapStr) (*InstanceSimplify, error) {

	err := mapstr.SetValueToStructByTags(is, data)
	if nil != err {
		return nil, err
	}

	bizID, err := is.ParseBizID(data)
	is.BizID = bizID
	return is, err
}

func (is *InstanceSimplify) ParseBizID(data mapstr.MapStr) (int64, error) {
	/*
	{
	  "metadata": {
		"label": {
		  "bk_biz_id": "2"
		}
	  }
	}
	 */
	 
    metaInterface, exist := data[common.MetadataField]
	if !exist {
		return 0, metadata.LabelKeyNotExistError
	}
    
    metaValue, ok := metaInterface.(map[string]interface{})
	if !ok {
		return 0, metadata.LabelKeyNotExistError
	}

	labelInterface, exist := metaValue["label"]
	if !exist {
		return 0, metadata.LabelKeyNotExistError
	}

	labelValue, ok := labelInterface.(map[string]interface{})
	if !ok {
		return 0, metadata.LabelKeyNotExistError
	}

	bizID, exist := labelValue[common.BKAppIDField]
	if !exist {
		return 0, metadata.LabelKeyNotExistError
	}

	return util.GetInt64ByInterface(bizID)
}


type BusinessSimplify struct {
	BKAppIDField      int64  `json:"bk_biz_id"`
	BKAppNameField    string `json:"bk_biz_name"`
	BKSupplierIDField int64  `json:"bk_supplier_id"`
	BKOwnerIDField    string `json:"bk_supplier_account"`
}

// Parse load the data from mapstr attribute into ObjectUnique instance
func (is *BusinessSimplify) Parse(data mapstr.MapStr) (*BusinessSimplify, error) {

	err := mapstr.SetValueToStructByTags(is, data)
	if nil != err {
		return nil, err
	}

	return is, err
}

type SetSimplify struct {
	BKAppIDField   int64  `json:"bk_biz_id"`
	BKSetIDField   int64  `json:"bk_set_id"`
	BKSetNameField string `json:"bk_set_name"`
}

// Parse load the data from mapstr attribute into ObjectUnique instance
func (is *SetSimplify) Parse(data mapstr.MapStr) (*SetSimplify, error) {

	err := mapstr.SetValueToStructByTags(is, data)
	if nil != err {
		return nil, err
	}

	return is, err
}

type ModuleSimplify struct {
	BKAppIDField      int64  `json:"bk_biz_id"`
	BKModuleIDField   int64  `json:"bk_module_id"`
	BKModuleNameField string `json:"bk_module_name"`
}

// Parse load the data from mapstr attribute into ObjectUnique instance
func (is *ModuleSimplify) Parse(data mapstr.MapStr) (*ModuleSimplify, error) {

	err := mapstr.SetValueToStructByTags(is, data)
	if nil != err {
		return nil, err
	}

	return is, err
}

type HostSimplify struct {
	BKAppIDField      int64  `json:"bk_biz_id"`
	BKModuleIDField   int64  `json:"bk_module_id"`
	BKSetIDField   int64  `json:"bk_set_id"`
	BKHostIDField   int64  `json:"bk_host_id"`
	BKHostNameField string `json:"bk_host_name"`
	BKHostInnerIPField string `json:"bk_host_innerip"`
}

func (is *HostSimplify) Parse(data mapstr.MapStr) (*HostSimplify, error) {

	err := mapstr.SetValueToStructByTags(is, data)
	if nil != err {
		return nil, err
	}

	return is, err
}
