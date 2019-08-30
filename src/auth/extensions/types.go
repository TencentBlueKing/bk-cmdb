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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

type AuthManager struct {
	clientSet apimachinery.ClientSetInterface
	Authorize auth.Authorize

	RegisterModelAttributeEnabled bool
	RegisterModelUniqueEnabled    bool
	RegisterModuleEnabled         bool
	RegisterSetEnabled            bool
	RegisterAuditCategoryEnabled  bool
	SkipReadAuthorization         bool
}

func NewAuthManager(clientSet apimachinery.ClientSetInterface, Authorize auth.Authorize) *AuthManager {
	return &AuthManager{
		clientSet:                     clientSet,
		Authorize:                     Authorize,
		RegisterModelAttributeEnabled: false,
		RegisterModelUniqueEnabled:    false,
		RegisterModuleEnabled:         false,
		RegisterSetEnabled:            false,
		SkipReadAuthorization:         true,
		RegisterAuditCategoryEnabled:  false,
	}
}

type InstanceSimplify struct {
	InstanceID int64  `field:"bk_inst_id"`
	Name       string `field:"bk_inst_name"`
	BizID      int64  `field:"bk_biz_id"`
	ObjectID   string `field:"bk_obj_id"`
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
		return 0, nil
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
	BKAppIDField      int64  `field:"bk_biz_id"`
	BKAppNameField    string `field:"bk_biz_name"`
	BKSupplierIDField int64  `field:"bk_supplier_id"`
	BKOwnerIDField    string `field:"bk_supplier_account"`
	IsDefault         int64  `field:"default"`

	Maintainer string `field:"bk_biz_maintainer"`
	Producer   string `field:"bk_biz_productor"`
	Tester     string `field:"bk_biz_tester"`
	Developer  string `field:"bk_biz_developer"`
	Operator   string `field:"operator"`
}

// Parse load the data from mapstr attribute into ObjectUnique instance
func (business *BusinessSimplify) Parse(data mapstr.MapStr) (*BusinessSimplify, error) {

	err := mapstr.SetValueToStructByTags(business, data)
	if nil != err {
		return nil, err
	}

	return business, err
}

type SetSimplify struct {
	BKAppIDField   int64  `field:"bk_biz_id"`
	BKSetIDField   int64  `field:"bk_set_id"`
	BKSetNameField string `field:"bk_set_name"`
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
	BKAppIDField      int64  `field:"bk_biz_id"`
	BKModuleIDField   int64  `field:"bk_module_id"`
	BKModuleNameField string `field:"bk_module_name"`
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
	BKAppIDField       int64  `field:"bk_biz_id"`
	BKModuleIDField    int64  `field:"bk_module_id"`
	BKSetIDField       int64  `field:"bk_set_id"`
	BKHostIDField      int64  `field:"bk_host_id"`
	BKHostNameField    string `field:"bk_host_name"`
	BKHostInnerIPField string `field:"bk_host_innerip"`
}

func (is *HostSimplify) Parse(data mapstr.MapStr) (*HostSimplify, error) {

	err := mapstr.SetValueToStructByTags(is, data)
	if nil != err {
		return nil, err
	}

	return is, err
}

type PlatSimplify struct {
	BKCloudIDField   int64  `field:"bk_cloud_id"`
	BKCloudNameField string `field:"bk_cloud_name"`
}

func (is *PlatSimplify) Parse(data mapstr.MapStr) (*PlatSimplify, error) {

	err := mapstr.SetValueToStructByTags(is, data)
	if nil != err {
		return nil, err
	}

	return is, err
}

type AuditCategorySimplify struct {
	BKAppIDField    int64  `field:"bk_biz_id"`
	BKOpTargetField string `field:"op_target"`
}

func (is *AuditCategorySimplify) Parse(data mapstr.MapStr) (*AuditCategorySimplify, error) {

	err := mapstr.SetValueToStructByTags(is, data)
	if nil != err {
		return nil, err
	}

	return is, err
}

type ModelUniqueSimplify struct {
	ID         uint64 `field:"id" json:"id" bson:"id"`
	ObjID      string `field:"bk_obj_id" json:"bk_obj_id" bson:"bk_obj_id"`
	Ispre      bool   `field:"ispre" json:"ispre" bson:"ispre"`
	BusinessID int64
}

func (cls *ModelUniqueSimplify) Parse(data mapstr.MapStr) (*ModelUniqueSimplify, error) {

	err := mapstr.SetValueToStructByTags(cls, data)
	if nil != err {
		return nil, err
	}

	// parse business id
	cls.BusinessID, err = metadata.ParseBizIDFromData(data)
	if nil != err {
		return nil, err
	}

	return cls, err
}

type ProcessSimplify struct {
	ProcessID    int64  `field:"bk_process_id"`
	ProcessName  string `field:"bk_process_name"`
	BKAppIDField int64  `field:"bk_biz_id"`
}

func (is *ProcessSimplify) Parse(data mapstr.MapStr) (*ProcessSimplify, error) {

	err := mapstr.SetValueToStructByTags(is, data)
	if nil != err {
		return nil, err
	}

	return is, err
}

type DynamicGroupSimplify struct {
	BKAppIDField int64  `field:"bk_biz_id"`
	ID           string `field:"id"`
	Name         string `field:"name"`
}

func (is *DynamicGroupSimplify) Parse(data mapstr.MapStr) (*DynamicGroupSimplify, error) {

	err := mapstr.SetValueToStructByTags(is, data)
	if nil != err {
		return nil, err
	}

	return is, err
}
