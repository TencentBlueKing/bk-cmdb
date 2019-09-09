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

import "configcenter/src/framework/core/errors"

type CreateSetTemplateOption struct {
	Name               string  `field:"name" json:"name" bson:"name"`
	ServiceTemplateIDs []int64 `field:"service_template_ids" json:"service_template_ids" bson:"service_template_ids"`
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
	SetTemplateIDs []int64 `field:"set_template_ids" json:"set_template_ids" bson:"set_template_ids"`
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
	Page BasePage `json:"page"`
}

type DiffSetTplWithInstOption struct {
	SetIDs []int64 `field:"bk_set_ids" json:"bk_set_ids" bson:"bk_set_ids" mapstructure:"bk_set_ids"`
}

type SyncSetTplToInstOption struct {
	SetIDs []int64 `field:"bk_set_ids" json:"bk_set_ids" bson:"bk_set_ids" mapstructure:"bk_set_ids"`
}

type SetModuleDiff struct {
	ModuleID            int64  `json:"bk_module_id"`
	ModuleName          string `json:"bk_module_name"`
	ServiceTemplateID   int64  `json:"service_template_id"`
	ServiceTemplateName string `json:"service_template_name"`
	DiffType            string `json:"diff_type"`
}

type SetDiff struct {
	ModuleDiffs []SetModuleDiff `json:"module_diffs"`
	SetID       int64           `json:"bk_set_id"`
	SetDetail   SetInst         `json:"set_detail"`
}

var (
	ModuleDiffAdd     = "add"
	ModuleDiffRemove  = "remove"
	ModuleDiffChanged = "changed"
)
