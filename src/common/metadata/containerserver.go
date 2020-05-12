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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/querybuilder"
)

// CreatePod request for CreatePod
type CreatePod struct {
	BizID int64         `json:"bk_biz_id"`
	Pod   mapstr.MapStr `json:"pod"`
}

// CreatePodResult Response for CreatePod
type CreatePodResult struct {
	CreatedOneOptionResult
}

// CreateManyPod request for CreateManyPod
type CreateManyPod struct {
	BizID   int64           `json:"bk_biz_id"`
	PodList []mapstr.MapStr `json:"pod_map"`
}

// CreateManyPodResult response for CreateManyPod
type CreateManyPodResult struct {
	CreatedManyOptionResult
}

// UpdatePodProperties parameter for UpdatePod
type UpdatePodProperties struct {
	PodID      string                 `json:"bk_pod_id"`
	Properties map[string]interface{} `json:"properties"`
}

// UpdatePod request for UpdatePod
type UpdatePod struct {
	BizID      int64                  `json:"bk_biz_id"`
	Properties map[string]interface{} `json:"properties"`
}

// UpdatePodResult response for UpdatePod
type UpdatePodResult struct {
	UpdatedOptionResult
}

// UpdateManyPod request for UpdateManyPod
type UpdateManyPod struct {
	BizID   int64                 `json:"bk_biz_id"`
	Updates []UpdatePodProperties `json:"updates"`
}

// UpdateManyPodResult response for UpdateManyPod
type UpdateManyPodResult struct {
	UpdatedOptionResult
}

// DeletePodResult response for delete pod
type DeletePodResult struct {
	BaseResp
}

// DeletePodOption option for delete pod
type DeletePodOption struct {
	PodID string `json:"bk_pod_id"`
}

// DeleteManyPod request for DeleteManyPod
type DeleteManyPod struct {
	BizID   int64
	Deletes []DeletePodOption `json:"deletes"`
}

// DeleteManyPodResult response for DeleteManyPod
type DeleteManyPodResult struct {
	BaseResp
	DeletedCount
}

// ListPod request for ListPod
type ListPod struct {
	BizID              int64                     `json:"bk_biz_id"`
	SetCond            []ConditionItem           `json:"set_cond"`
	SetIDs             []int64                   `json:"bk_set_ids"`
	ModuleIDs          []int64                   `json:"bk_module_ids"`
	HostPropertyFilter *querybuilder.QueryFilter `json:"pod_property_filter"`
	Fields             []string                  `json:"fields"`
	Page               BasePage                  `json:"page"`
}

// ListPodResult response for ListPod
type ListPodResult struct {
	Count int                      `json:"count"`
	Info  []map[string]interface{} `json:"info"`
}
