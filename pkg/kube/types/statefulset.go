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
	"configcenter/pkg/kube/orm"
	"time"

	"configcenter/pkg/common"
)

// StatefulSetUpdateStrategyType is a strategy according to which a StatefulSet gets updated.
type StatefulSetUpdateStrategyType string

const (
	// RollingUpdateStatefulSetStrategyType indicates that update will be
	// applied to all Pods in the StatefulSet with respect to the StatefulSet
	// ordering constraints.
	RollingUpdateStatefulSetStrategyType StatefulSetUpdateStrategyType = "RollingUpdate"

	// OnDeleteStatefulSetStrategyType pods are recreated from the StatefulSetSpec when they are manually deleted.
	OnDeleteStatefulSetStrategyType StatefulSetUpdateStrategyType = "OnDelete"
)

// RollingUpdateStatefulSetStrategy spec to control the desired behavior of rolling update.
type RollingUpdateStatefulSetStrategy struct {
	// Partition indicates the ordinal at which the StatefulSet should be partitioned for updates.
	Partition *int32 `json:"partition" bson:"partition"`

	// The maximum number of pods that can be unavailable during the update.
	// Value can be an absolute number (ex: 5) or a percentage of desired pods (ex: 10%).
	// This can not be 0 if MaxSurge is 0.
	MaxUnavailable *IntOrString `json:"max_unavailable" bson:"max_unavailable"`
}

// StatefulSet define the statefulSet struct.
type StatefulSet struct {
	Workload              `json:",inline" bson:",inline"`
	StrategyType          *StatefulSetUpdateStrategyType    `json:"strategy_type,omitempty" bson:"strategy_type"`
	RollingUpdateStrategy *RollingUpdateStatefulSetStrategy `json:"rolling_update_strategy,omitempty" bson:"rolling_update_strategy"`
}

// StatefulSetUpdateData defines the statefulSet update data common operation.
type StatefulSetUpdateData struct {
	WlCommonUpdate `json:",inline"`
	Info           StatefulSet `json:"info"`
}

// BuildUpdateData build statefulSet update data
func (d *StatefulSetUpdateData) BuildUpdateData(user string) (map[string]interface{}, error) {
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
