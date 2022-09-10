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

// GameDeploymentUpdateStrategyType defines strategies for pods in-place update.
type GameDeploymentUpdateStrategyType string

const (
	// RollingGameDeploymentUpdateStrategyType indicates that we always delete Pod and create new Pod
	// during Pod update, which is the default behavior.
	RollingGameDeploymentUpdateStrategyType GameDeploymentUpdateStrategyType = "RollingUpdate"

	// InPlaceGameDeploymentUpdateStrategyType indicates that we will in-place update Pod instead of
	// recreating pod. Currently we only allow image update for pod spec. Any other changes to the pod spec will be
	// rejected by kube-apiserver
	InPlaceGameDeploymentUpdateStrategyType GameDeploymentUpdateStrategyType = "InplaceUpdate"

	// HotPatchGameDeploymentUpdateStrategyType indicates that we will hot patch container image with pod being active.
	// Currently we only allow image update for pod spec. Any other changes to the pod spec will be
	// rejected by kube-apiserver
	HotPatchGameDeploymentUpdateStrategyType GameDeploymentUpdateStrategyType = "HotPatchUpdate"
)

// RollingUpdateGameDeployment gameDeployment update strategy
type RollingUpdateGameDeployment struct {
	// Partition is the desired number of pods in old revisions. It means when partition
	// is set during pods updating, (replicas - partition) number of pods will be updated.
	Partition *int32 `json:"partition" bson:"partition"`

	// The maximum number of pods that can be unavailable during the update.
	// Value can be an absolute number (ex: 5) or a percentage of desired pods (ex: 10%).
	// This can not be 0 if MaxSurge is 0.
	MaxUnavailable *IntOrString `json:"max_unavailable" bson:"max_unavailable"`

	// The maximum number of pods that can be scheduled above the desired number of pods.
	// Value can be an absolute number (ex: 5) or a percentage of desired pods (ex: 10%).
	// This can not be 0 if MaxUnavailable is 0.
	MaxSurge *IntOrString `json:"max_surge" bson:"max_surge"`
}

// GameDeployment define the gameDeployment struct.
type GameDeployment struct {
	Workload              `json:",inline" bson:",inline"`
	StrategyType          *GameDeploymentUpdateStrategyType `json:"strategy_type,omitempty" bson:"strategy_type"`
	RollingUpdateStrategy *RollingUpdateGameDeployment      `json:"rolling_update_strategy,omitempty" bson:"rolling_update_strategy"`
}

// GameDeployUpdateData defines the gameDeployment update data common operation.
type GameDeployUpdateData struct {
	WlCommonUpdate `json:",inline"`
	Info           GameDeployment `json:"info"`
}

// BuildUpdateData build gameDeployment update data
func (d *GameDeployUpdateData) BuildUpdateData(user string) (map[string]interface{}, error) {
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
