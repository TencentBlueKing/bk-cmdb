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
	"errors"
	"fmt"

	"configcenter/src/common"
	ccErr "configcenter/src/common/errors"
	"configcenter/src/common/time"
	"configcenter/src/common/util"
)

type CreateSetTemplateOption struct {
	Name               string  `field:"name" json:"name" bson:"name" mapstructure:"name"`
	ServiceTemplateIDs []int64 `field:"service_template_ids" json:"service_template_ids" bson:"service_template_ids" mapstructure:"service_template_ids"`
}

type UpdateSetTemplateOption struct {
	Name               string  `field:"name" json:"name" bson:"name"`
	ServiceTemplateIDs []int64 `field:"service_template_ids" json:"service_template_ids" bson:"service_template_ids"`
}

func (option UpdateSetTemplateOption) Validate(errProxy ccErr.DefaultCCErrorIf) (string, error) {
	if len(option.Name) == 0 && option.ServiceTemplateIDs == nil {
		return "", errors.New("at least one update field not empty")
	}

	if len(option.Name) > 0 {
		var err error
		option.Name, err = util.ValidTopoNameField(option.Name, "name", errProxy)
		if err != nil {
			return "name", err
		}
	}
	return "", nil
}

type SetTemplateResult struct {
	BaseResp
	Data SetTemplate `field:"data" json:"data" bson:"data"`
}

type DeleteSetTemplateOption struct {
	SetTemplateIDs []int64 `field:"set_template_ids" json:"set_template_ids" bson:"set_template_ids" mapstructure:"set_template_ids"`
}

type ListSetTemplateOption struct {
	SetTemplateIDs []int64  `field:"set_template_ids" json:"set_template_ids" bson:"set_template_ids"`
	Page           BasePage `json:"page"`
}

type ListSetTemplateResult struct {
	BaseResp
	Data MultipleSetTemplateResult `json:"data"`
}

type MultipleSetTemplateResult struct {
	Count int64         `json:"count"`
	Info  []SetTemplate `json:"info"`
}

type SetTemplateWithStatistics struct {
	SetInstanceCount int64       `json:"set_instance_count"`
	SetTemplate      SetTemplate `json:"set_template"`
}

type MultipleSetTemplateWithStatisticsResult struct {
	Count int64                       `json:"count"`
	Info  []SetTemplateWithStatistics `json:"info"`
}

type CountSetTplInstOption struct {
	SetTemplateIDs []int64 `field:"set_template_ids" json:"set_template_ids" bson:"set_template_ids"`
}

type CountSetTplInstItem struct {
	SetTemplateID    int64 `json:"set_template_id" bson:"_id"`
	SetInstanceCount int64 `json:"set_instances_count" bson:"set_instances_count"`
}

type ListSetByTemplateOption struct {
	Page   BasePage `json:"page"`
	SetIDs []int64  `field:"bk_set_ids" json:"bk_set_ids" bson:"bk_set_ids" mapstructure:"bk_set_ids"`
}

type DiffSetTplWithInstOption struct {
	SetIDs []int64 `field:"bk_set_ids" json:"bk_set_ids" bson:"bk_set_ids" mapstructure:"bk_set_ids"`
}

type SyncSetTplToInstOption struct {
	SetIDs []int64 `field:"bk_set_ids" json:"bk_set_ids" bson:"bk_set_ids" mapstructure:"bk_set_ids"`
}

type SetSyncStatusOption struct {
	SetIDs []int64 `field:"bk_set_ids" json:"bk_set_ids" bson:"bk_set_ids" mapstructure:"bk_set_ids"`
}

type SetModuleDiff struct {
	ModuleID            int64  `json:"bk_module_id" mapstructure:"bk_module_id"`
	ModuleName          string `json:"bk_module_name" mapstructure:"bk_module_name"`
	ServiceTemplateID   int64  `json:"service_template_id" mapstructure:"service_template_id"`
	ServiceTemplateName string `json:"service_template_name" mapstructure:"service_template_name"`
	DiffType            string `json:"diff_type" mapstructure:"diff_type"`
}

type SetDiff struct {
	ModuleDiffs []SetModuleDiff            `json:"module_diffs"`
	SetID       int64                      `json:"bk_set_id"`
	SetDetail   SetInst                    `json:"set_detail"`
	TopoPath    []TopoInstanceNodeSimplify `json:"topo_path"`
	NeedSync    bool                       `json:"need_sync"`
}

func (sd *SetDiff) UpdateNeedSyncField() {
	sd.NeedSync = false
	for _, module := range sd.ModuleDiffs {
		if module.DiffType != ModuleDiffUnchanged {
			sd.NeedSync = true
			break
		}
	}
}

type SetTplDiffResult struct {
	Difference      []SetDiff       `json:"difference"`
	ModuleHostCount map[int64]int64 `json:"module_host_count"`
}

type SyncModuleTask struct {
	Set         SetInst                    `json:"set"`
	SetTopoPath []TopoInstanceNodeSimplify `json:"set_topo_path"`
	ModuleDiff  SetModuleDiff              `json:"module_diff"`
}

var (
	ModuleDiffAdd       = "add"
	ModuleDiffRemove    = "remove"
	ModuleDiffChanged   = "changed"
	ModuleDiffUnchanged = "unchanged"
)

type DeleteSetTemplateSyncStatusOption struct {
	SetIDs []int64 `field:"bk_set_ids" json:"bk_set_ids" bson:"bk_set_ids" mapstructure:"bk_set_ids"`
	BizID  int64   `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id" mapstructure:"bk_biz_id"`
}

type ListSetTemplateSyncStatusOption struct {
	BizID         int64           `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id" mapstructure:"bk_biz_id"`
	SetIDs        []int64         `field:"bk_set_ids" json:"bk_set_ids" bson:"bk_set_ids" mapstructure:"bk_set_ids"`
	TaskIDs       []string        `field:"task_ids" json:"task_ids" bson:"task_ids" mapstructure:"task_ids"`
	SearchKey     string          `field:"search" json:"search" bson:"search" mapstructure:"search"`
	SetTemplateID int64           `field:"set_template_id" json:"set_template_id" bson:"set_template_id" mapstructure:"set_template_id"`
	Creator       string          `field:"creator" json:"creator,omitempty" bson:"creator" mapstructure:"creator"`
	StartTime     *time.Time      `field:"start_time" json:"start_time,omitempty" bson:"create_time" mapstructure:"start_time"`
	EndTime       *time.Time      `field:"end_time" json:"end_time,omitempty" bson:"end_time" mapstructure:"end_time"`
	Status        []APITaskStatus `field:"status" json:"status" bson:"status" mapstructure:"status"`
	Page          BasePage        `field:"page" json:"page" bson:"page" mapstructure:"page"`
}

// ToSetCond parse option to query set condition
func (option ListSetTemplateSyncStatusOption) ToSetCond(errProxy ccErr.DefaultCCErrorIf) (map[string]interface{},
	ccErr.CCErrorCoder) {

	if option.BizID == 0 {
		return nil, errProxy.CCErrorf(common.CCErrCommParamsNeedSet, common.BKAppIDField)
	}

	if option.SetTemplateID == 0 {
		return nil, errProxy.CCErrorf(common.CCErrCommParamsNeedSet, common.BKServiceTemplateIDField)
	}

	filter := map[string]interface{}{
		common.BKSetTemplateIDField: option.SetTemplateID,
		common.BKAppIDField:         option.BizID,
	}

	if option.SetIDs != nil {
		filter[common.BKSetIDField] = map[string]interface{}{
			common.BKDBIN: option.SetIDs,
		}
	}

	if len(option.SearchKey) != 0 {
		filter[common.BKSetNameField] = map[string]interface{}{
			common.BKDBLIKE:    fmt.Sprintf(".*%s.*", option.SearchKey),
			common.BKDBOPTIONS: "i",
		}
	}

	return filter, nil
}

// ToStatusCond parse option to query sync status condition
func (option ListSetTemplateSyncStatusOption) ToStatusCond(errProxy ccErr.DefaultCCErrorIf) (*QueryCondition,
	ccErr.CCErrorCoder) {

	if len(option.SetIDs) == 0 {
		return nil, errProxy.CCErrorf(common.CCErrCommParamsNeedSet, common.BKSetIDField)
	}

	condition := &QueryCondition{
		Page: option.Page,
		Fields: []string{common.BKInstIDField, common.CreateTimeField, common.LastTimeField, common.CreatorField,
			common.BKStatusField},
		Condition: map[string]interface{}{
			common.BKInstIDField: map[string]interface{}{
				common.BKDBIN: option.SetIDs,
			},
			common.BKTaskTypeField: common.SyncSetTaskFlag,
		},
	}

	if len(option.TaskIDs) != 0 {
		condition.Condition[common.BKTaskIDField] = map[string]interface{}{
			common.BKDBIN: option.TaskIDs,
		}
	}

	if len(option.Status) != 0 {
		condition.Condition[common.BKStatusField] = map[string]interface{}{
			common.BKDBIN: option.Status,
		}
	}

	if len(option.Creator) != 0 {
		condition.Condition[common.CreatorField] = option.Creator
	}

	timeConditionItem := make([]TimeConditionItem, 0)
	if option.StartTime != nil {
		timeConditionItem = append(timeConditionItem, TimeConditionItem{
			Field: common.CreateTimeField,
			Start: option.StartTime,
		})
	}

	if option.EndTime != nil {
		timeConditionItem = append(timeConditionItem, TimeConditionItem{
			Field: common.LastTimeField,
			End:   option.EndTime,
		})
	}

	if len(timeConditionItem) > 0 {
		condition.TimeCondition = &TimeCondition{
			Operator: "and",
			Rules:    timeConditionItem,
		}
	}

	return condition, nil
}

type SetUpdateToDateStatus struct {
	SetID    int64 `json:"bk_set_id"`
	NeedSync bool  `json:"need_sync"`
}

type SetTemplateUpdateToDateStatus struct {
	Sets          []SetUpdateToDateStatus `json:"sets"`
	SetTemplateID int64                   `json:"set_template_id"`
	NeedSync      bool                    `json:"need_sync"`
}

type BatchCheckSetInstUpdateToDateStatusOption struct {
	SetTemplateIDs []int64 `field:"set_template_ids" json:"set_template_ids" bson:"set_template_ids" mapstructure:"set_template_ids"`
}

// SetTempAttr simplified set template attribute
type SetTempAttr struct {
	AttributeID   int64       `json:"bk_attribute_id"`
	PropertyValue interface{} `json:"bk_property_value"`
}

// Validate SetTempAttr
func (s *SetTempAttr) Validate() ccErr.RawErrorInfo {
	if s.AttributeID == 0 {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{
			common.BKAttributeIDField}}
	}

	if s.PropertyValue == nil {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{
			common.BKPropertyValueField}}
	}

	return ccErr.RawErrorInfo{}
}

// UpdateSetTempAttrOption update set template attributes option
type UpdateSetTempAttrOption struct {
	BizID      int64         `json:"bk_biz_id"`
	ID         int64         `json:"id"`
	Attributes []SetTempAttr `json:"attributes"`
}

// Validate UpdateSetTempAttrOption
func (s *UpdateSetTempAttrOption) Validate() ccErr.RawErrorInfo {
	if s.BizID == 0 {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{common.BKAppIDField}}
	}

	if s.ID == 0 {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{common.BKFieldID}}
	}

	if len(s.Attributes) == 0 {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{"attributes"}}
	}

	for _, attribute := range s.Attributes {
		if rawErr := attribute.Validate(); rawErr.ErrCode != 0 {
			return rawErr
		}
	}
	return ccErr.RawErrorInfo{}
}

// DeleteSetTempAttrOption delete set template attributes option
type DeleteSetTempAttrOption struct {
	BizID        int64   `json:"bk_biz_id"`
	ID           int64   `json:"id"`
	AttributeIDs []int64 `json:"bk_attribute_ids"`
}

// Validate DeleteSetTempAttrOption
func (s *DeleteSetTempAttrOption) Validate() ccErr.RawErrorInfo {
	if s.BizID == 0 {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{common.BKAppIDField}}
	}

	if s.ID == 0 {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{common.BKFieldID}}
	}

	if len(s.AttributeIDs) == 0 {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{"bk_attribute_ids"}}
	}

	return ccErr.RawErrorInfo{}
}

// ListSetTempAttrOption list set template attributes option
type ListSetTempAttrOption struct {
	BizID int64 `json:"bk_biz_id"`
	ID    int64 `json:"id"`
}

// Validate ListSetTempAttrOption
func (s *ListSetTempAttrOption) Validate() ccErr.RawErrorInfo {
	if s.BizID == 0 {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{common.BKAppIDField}}
	}

	if s.ID == 0 {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{common.BKFieldID}}
	}

	return ccErr.RawErrorInfo{}
}
