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
	"configcenter/src/common"
	"configcenter/src/common/errors"
)

// CreateServiceCategoryOption TODO
type CreateServiceCategoryOption struct {
	BizID    int64  `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id"`
	Name     string `field:"name" json:"name,omitempty" bson:"name"`
	ParentID int64  `field:"bk_parent_id" json:"bk_parent_id,omitempty" bson:"bk_parent_id"`
}

// CreateServiceTemplateOption TODO
type CreateServiceTemplateOption struct {
	BizID             int64  `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id"`
	Name              string `field:"name" json:"name,omitempty" bson:"name"`
	ServiceCategoryID int64  `field:"service_category_id" json:"service_category_id,omitempty" bson:"service_category_id"`
	HostApplyEnabled  bool   `field:"host_apply_enabled" json:"host_apply_enabled" bson:"host_apply_enabled"`
}

// CreateSvcTempAllInfoOption create service template all info option
type CreateSvcTempAllInfoOption struct {
	BizID             int64             `json:"bk_biz_id"`
	Name              string            `json:"name"`
	ServiceCategoryID int64             `json:"service_category_id"`
	Attributes        []SvcTempAttr     `json:"attributes"`
	Processes         []ProcessTemplate `json:"processes"`
}

// SvcTempAttr simplified service template attribute
type SvcTempAttr struct {
	AttributeID   int64       `json:"bk_attribute_id"`
	PropertyValue interface{} `json:"bk_property_value"`
}

// Validate SvcTempAttr
func (s *SvcTempAttr) Validate() errors.RawErrorInfo {
	if s.AttributeID == 0 {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{
			common.BKAttributeIDField}}
	}

	if s.PropertyValue == nil {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{
			common.BKPropertyValueField}}
	}

	return errors.RawErrorInfo{}
}

// CreateSvcTempAttrsOption create service template attributes option
type CreateSvcTempAttrsOption struct {
	BizID             int64         `json:"bk_biz_id"`
	ServiceTemplateID int64         `json:"service_template_id"`
	Attributes        []SvcTempAttr `json:"attributes"`
}

// UpdateSvcTempAllInfoOption update service template all info option
type UpdateSvcTempAllInfoOption struct {
	ID                int64             `json:"id"`
	BizID             int64             `json:"bk_biz_id"`
	Name              string            `json:"name"`
	ServiceCategoryID int64             `json:"service_category_id"`
	Attributes        []SvcTempAttr     `json:"attributes"`
	Processes         []ProcessTemplate `json:"processes"`
}

// GetSvcTempAllInfoOption get service template all info option
type GetSvcTempAllInfoOption struct {
	ID    int64 `json:"id"`
	BizID int64 `json:"bk_biz_id"`
}

// GetSvcTempAllInfoResult get service template all info result
type GetSvcTempAllInfoResult struct {
	BaseResp `json:",inline"`
	Data     *SvcTempAllInfo `json:"data"`
}

// SvcTempAllInfo service template all info
type SvcTempAllInfo struct {
	ID                int64                 `json:"id"`
	BizID             int64                 `json:"bk_biz_id"`
	Name              string                `json:"name"`
	ServiceCategoryID int64                 `json:"service_category_id"`
	Attributes        []ServiceTemplateAttr `json:"attributes"`
	Processes         []ProcessTemplate     `json:"processes"`
}

// UpdateServiceTemplateOption TODO
type UpdateServiceTemplateOption struct {
	BizID             int64  `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id"`
	ID                int64  `field:"id" json:"id,omitempty" bson:"id"`
	Name              string `field:"name" json:"name,omitempty" bson:"name"`
	ServiceCategoryID int64  `field:"service_category_id" json:"service_category_id,omitempty" bson:"service_category_id"`
}

// RemoveFromModuleHost TODO
type RemoveFromModuleHost struct {
	MoveToIdle        bool    `field:"move_to_idle" json:"move_to_idle"`
	HostID            int64   `field:"bk_host_id" json:"bk_host_id"`
	RemoveFromModules []int64 `field:"bk_module_ids" json:"bk_module_ids"`
	FinalModules      []int64 `field:"final_module_ids" json:"final_module_ids"`
}

// ServiceInstanceDeletePreview TODO
type ServiceInstanceDeletePreview struct {
	ToMoveModuleHosts []RemoveFromModuleHost `field:"to_move_module_hosts" json:"to_move_module_hosts"`
	HostApplyPlan     HostApplyPlanResult    `field:"host_apply_plan" json:"host_apply_plan"`
}

// UpdateServTempAttrOption update service template attributes option
type UpdateServTempAttrOption struct {
	BizID      int64         `json:"bk_biz_id"`
	ID         int64         `json:"id"`
	Attributes []SvcTempAttr `json:"attributes"`
}

// Validate UpdateServTempAttrOption
func (s *UpdateServTempAttrOption) Validate() errors.RawErrorInfo {
	if s.BizID == 0 {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{common.BKAppIDField}}
	}

	if s.ID == 0 {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{"id"}}
	}

	if len(s.Attributes) == 0 {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{"attributes"}}
	}

	for _, attribute := range s.Attributes {
		if rawErr := attribute.Validate(); rawErr.ErrCode != 0 {
			return rawErr
		}
	}
	return errors.RawErrorInfo{}
}

// DeleteServTempAttrOption delete service template attributes option
type DeleteServTempAttrOption struct {
	BizID        int64   `json:"bk_biz_id"`
	ID           int64   `json:"id"`
	AttributeIDs []int64 `json:"bk_attribute_ids"`
}

// Validate DeleteServiceTemplateAttrOption
func (s *DeleteServTempAttrOption) Validate() errors.RawErrorInfo {
	if s.BizID == 0 {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{common.BKAppIDField}}
	}

	if s.ID == 0 {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{common.BKFieldID}}
	}

	if len(s.AttributeIDs) == 0 {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{"bk_attribute_ids"}}
	}

	return errors.RawErrorInfo{}
}

// ListServTempAttrOption list service template attributes option
type ListServTempAttrOption struct {
	BizID  int64    `json:"bk_biz_id"`
	ID     int64    `json:"id"`
	Fields []string `json:"fields"`
}

// Validate ListServTempAttrOption
func (s *ListServTempAttrOption) Validate() errors.RawErrorInfo {
	if s.BizID == 0 {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{common.BKAppIDField}}
	}

	if s.ID == 0 {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{"id"}}
	}

	return errors.RawErrorInfo{}
}
