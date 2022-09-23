/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package types

import (
	"time"

	"configcenter/src/common"
	"configcenter/src/kube/orm"
)

// CronJob define the cronJob struct.
type CronJob struct {
	WorkloadBase    `json:",inline" bson:",inline"`
	Labels          *map[string]string `json:"labels,omitempty" bson:"labels"`
	Selector        *LabelSelector     `json:"selector,omitempty" bson:"selector"`
	Replicas        *int64             `json:"replicas,omitempty" bson:"replicas"`
	MinReadySeconds *int64             `json:"min_ready_seconds,omitempty" bson:"min_ready_seconds"`
}

// GetWorkloadBase get workload base
func (c *CronJob) GetWorkloadBase() WorkloadBase {
	return c.WorkloadBase
}

// SetWorkloadBase set workload base
func (c *CronJob) SetWorkloadBase(wl WorkloadBase) {
	c.WorkloadBase = wl
}

// CronJobUpdateData defines the cronJob update data common operation.
type CronJobUpdateData struct {
	WlCommonUpdate `json:",inline"`
	Info           CronJob `json:"info"`
}

// BuildUpdateData build cronJob update data
func (d *CronJobUpdateData) BuildUpdateData(user string) (map[string]interface{}, error) {
	now := time.Now().Unix()
	opts := orm.NewFieldOptions().AddIgnoredFields(wlIgnoreField...)
	updateData, err := orm.GetUpdateFieldsWithOption(d.Info, opts)
	if err != nil {
		return nil, err
	}
	updateData[common.LastTimeField] = now
	updateData[common.ModifierField] = user
	return updateData, err
}
