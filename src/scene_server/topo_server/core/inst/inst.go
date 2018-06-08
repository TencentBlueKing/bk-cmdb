/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package inst

import (
	"configcenter/src/apimachinery"
	frtypes "configcenter/src/common/mapstr"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

type inst struct {
	clientSet apimachinery.ClientSetInterface
	params    types.LogicParams
}

func (cli *inst) Create() error {
	return nil
}
func (cli *inst) Update() error {
	return nil
}
func (cli *inst) Delete() error {
	return nil
}
func (cli *inst) IsExists() (bool, error) {
	return false, nil
}
func (cli *inst) Save() error {
	return nil
}

func (cli *inst) GetObject() model.Object {
	return nil
}

func (cli *inst) GetInstID() int {
	return 0
}

func (cli *inst) GetInstName() string {

	return ""
}

func (cli *inst) ToMapStr() (frtypes.MapStr, error) {
	return nil, nil
}
func (cli *inst) SetValue(key string, value interface{}) error {
	return nil
}

func (cli *inst) SetValues(values frtypes.MapStr) error {
	return nil
}

func (cli *inst) GetValues() (frtypes.MapStr, error) {
	return nil, nil
}
