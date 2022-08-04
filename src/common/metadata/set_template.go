// Package metadata TODO
/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017,-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */
package metadata

import (
	"time"

	"configcenter/src/common"
	"configcenter/src/common/errors"
	"configcenter/src/common/util"
)

// SetTemplate 集群模板
type SetTemplate struct {
	ID    int64  `field:"id" json:"id" bson:"id"`
	Name  string `field:"name" json:"name" bson:"name"`
	BizID int64  `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id"`

	// 通用字段
	Creator         string    `field:"creator" json:"creator" bson:"creator"`
	Modifier        string    `field:"modifier" json:"modifier" bson:"modifier"`
	CreateTime      time.Time `field:"create_time" json:"create_time" bson:"create_time"`
	LastTime        time.Time `field:"last_time" json:"last_time" bson:"last_time"`
	SupplierAccount string    `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`
}

// Validate TODO
func (st SetTemplate) Validate(errProxy errors.DefaultCCErrorIf) (key string, err error) {
	st.Name, err = util.ValidTopoNameField(st.Name, "name", errProxy)
	if err != nil {
		return "name", err
	}
	return "", nil
}

// SetServiceTemplateRelation 拓扑模板与服务模板多对多关系, 记录拓扑模板的构成
type SetServiceTemplateRelation struct {
	BizID             int64  `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id"`
	SetTemplateID     int64  `field:"set_template_id" json:"set_template_id" bson:"set_template_id"`
	ServiceTemplateID int64  `field:"service_template_id" json:"service_template_id" bson:"service_template_id"`
	SupplierAccount   string `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`
}

// SetTemplateAttr set template attributes, used to generate set, should not include non-editable fields
type SetTemplateAttr struct {
	ID int64 `json:"id" bson:"id"`

	BizID         int64       `json:"bk_biz_id" bson:"bk_biz_id"`
	SetTemplateID int64       `json:"set_template_id" bson:"set_template_id"`
	AttributeID   int64       `json:"bk_attribute_id" bson:"bk_attribute_id"`
	PropertyValue interface{} `json:"bk_property_value" bson:"bk_property_value"`

	Creator         string    `json:"creator" bson:"creator"`
	Modifier        string    `json:"modifier" bson:"modifier"`
	CreateTime      time.Time `json:"create_time" bson:"create_time"`
	LastTime        time.Time `json:"last_time" bson:"last_time"`
	SupplierAccount string    `json:"bk_supplier_account" bson:"bk_supplier_account"`
}

// Validate SetTemplateAttr
func (s *SetTemplateAttr) Validate() errors.RawErrorInfo {
	if s.BizID == 0 {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{common.BKAppIDField}}
	}

	if s.SetTemplateID == 0 {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{
			common.BKServiceTemplateIDField}}
	}

	if s.AttributeID == 0 {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{
			common.BKAttributeIDField}}
	}

	if s.PropertyValue == nil {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{
			common.BKPropertyTypeField}}
	}

	return errors.RawErrorInfo{}
}

// SetTempAttrData set template attributes data
type SetTempAttrData struct {
	Attributes []SetTemplateAttr `json:"attributes"`
}
