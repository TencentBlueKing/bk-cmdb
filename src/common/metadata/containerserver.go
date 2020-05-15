/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under,
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package metadata

import (
	"fmt"

	"configcenter/src/common/mapstr"
	"configcenter/src/common/querybuilder"
)

// CreatePod request for CreatePod
type CreatePod struct {
	BizID int64         `json:"bk_biz_id" mapstructure:"bk_biz_id"`
	Pod   mapstr.MapStr `json:"pod" mapstructure:"pod"`
}

// CreateManyPod request for CreateManyPod
type CreateManyPod struct {
	BizID   int64           `json:"bk_biz_id" mapstructure:"bk_biz_id"`
	PodList []mapstr.MapStr `json:"pod_list" mapstructure:"pod_list"`
}

// UpdatePod parameter for UpdatePod
type UpdatePod struct {
	BizID int64 `json:"bk_biz_id" mapstructure:"bk_biz_id"`
	UpdateOption
}

// // UpdateManyPod request for UpdateManyPod
// type UpdateManyPod struct {
// 	BizID   int64                 `json:"bk_biz_id"`
// 	Updates []UpdatePodProperties `json:"updates"`
// }

// DeletePod parameter for DeletePod
type DeletePod struct {
	BizID int64 `json:"bk_biz_id" mapstructure:"bk_biz_id"`
	DeleteOption
}

// DeletePodOption option for delete pod
type DeletePodOption struct {
	PodID string `json:"bk_pod_id"`
}

// // DeleteManyPod request for DeleteManyPod
// type DeleteManyPod struct {
// 	BizID   int64
// 	Deletes []DeletePodOption `json:"deletes"`
// }

// // DeleteManyPodResult response for DeleteManyPod
// type DeleteManyPodResult struct {
// 	BaseResp
// 	DeletedCount
// }

// FindPod request for FindPod
type FindPod struct {
	BizID int64 `json:"bk_biz_id"`
}

// ListPods request for ListPod
type ListPods struct {
	BizID             int64                     `json:"bk_biz_id"`
	SetIDs            []int64                   `json:"bk_set_ids"`
	ModuleIDs         []int64                   `json:"bk_module_ids"`
	PodPropertyFilter *querybuilder.QueryFilter `json:"pod_property_filter"`
	Fields            []string                  `json:"fields"`
	Page              BasePage                  `json:"page"`
}

// Validate whether ListPods is valid
func (option ListPods) Validate() (errKey string, err error) {
	if key, err := option.Page.Validate(false); err != nil {
		return fmt.Sprintf("page.%s", key), err
	}

	if option.PodPropertyFilter != nil {
		if key, err := option.PodPropertyFilter.Validate(); err != nil {
			return fmt.Sprintf("pod_property_filter.%s", key), err
		}
		if option.PodPropertyFilter.GetDeep() > querybuilder.MaxDeep {
			return "pod_property_filter.rules", fmt.Errorf("exceed max query condition deepth: %d", querybuilder.MaxDeep)
		}
	}
	return "", nil
}

// ListPodsResult response for ListPod
type ListPodsResult struct {
	Count int             `json:"count"`
	Info  []mapstr.MapStr `json:"info"`
}
