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
	"net/http"

	"configcenter/src/common"
)

type CreateSetTemplateOption struct {
	Name               string  `field:"name" json:"name" bson:"name" mapstructure:"name"`
	ServiceTemplateIDs []int64 `field:"service_template_ids" json:"service_template_ids" bson:"service_template_ids" mapstructure:"service_template_ids"`
}

type UpdateSetTemplateOption struct {
	Name               string  `field:"name" json:"name" bson:"name"`
	ServiceTemplateIDs []int64 `field:"service_template_ids" json:"service_template_ids" bson:"service_template_ids"`
}

func (option UpdateSetTemplateOption) Validate() (string, error) {
	if len(option.Name) == 0 && option.ServiceTemplateIDs == nil {
		return "", errors.New("at least one update field not empty")
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
	Header      http.Header                `json:"header"`
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
	BizID         int64      `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id" mapstructure:"bk_biz_id"`
	SetIDs        []int64    `field:"bk_set_ids" json:"bk_set_ids" bson:"bk_set_ids" mapstructure:"bk_set_ids"`
	TaskIDs       []string   `field:"task_ids" json:"task_ids" bson:"task_ids" mapstructure:"task_ids"`
	SearchKey     string     `field:"search" json:"search" bson:"search" mapstructure:"search"`
	SetTemplateID int64      `field:"set_template_id" json:"set_template_id" bson:"set_template_id" mapstructure:"set_template_id"`
	Creator       string     `field:"creator" json:"creator,omitempty" bson:"creator" mapstructure:"creator"`
	StartTime     *Time      `field:"start_time" json:"start_time,omitempty" bson:"create_time" mapstructure:"start_time"`
	EndTime       *Time      `field:"end_time" json:"end_time,omitempty" bson:"end_time" mapstructure:"end_time"`
	Status        SyncStatus `field:"status" json:"status" bson:"status" mapstructure:"status"`
	Page          BasePage   `field:"page" json:"page" bson:"page" mapstructure:"page"`
}

func (option ListSetTemplateSyncStatusOption) ToFilter() map[string]interface{} {
	filter := map[string]interface{}{
		common.BKAppIDField:         option.BizID,
		common.BKSetTemplateIDField: option.SetTemplateID,
	}
	if option.SetIDs != nil {
		filter[common.BKSetIDField] = map[string]interface{}{
			common.BKDBIN: option.SetIDs,
		}
	}
	if option.TaskIDs != nil {
		filter[common.BKTaskIDField] = map[string]interface{}{
			common.BKDBIN: option.TaskIDs,
		}
	}
	if len(option.Status) != 0 {
		filter[common.BKStatusField] = option.Status
	}
	if len(option.Creator) != 0 {
		filter[common.CreatorField] = option.Creator
	}
	if option.StartTime != nil {
		filter[common.CreateTimeField] = map[string]interface{}{
			common.BKDBGTE: option.StartTime,
		}
	}
	if option.EndTime != nil {
		filter[common.LastTimeField] = map[string]interface{}{
			common.BKDBLTE: option.EndTime,
		}
	}
	if len(option.SearchKey) != 0 {
		filter[common.BKSetNameField] = map[string]interface{}{
			common.BKDBLIKE: fmt.Sprintf(".*%s.*", option.SearchKey),
		}
	}
	return filter
}

type MultipleSetTemplateSyncStatus struct {
	Count int64                   `json:"count"`
	Info  []SetTemplateSyncStatus `field:"info" json:"info" bson:"info" mapstructure:"info"`
}
