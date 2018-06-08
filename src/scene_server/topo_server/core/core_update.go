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

package core

import (
	"configcenter/src/common/condition"
	frtypes "configcenter/src/common/mapstr"
	"configcenter/src/scene_server/topo_server/core/types"
)

func (cli *core) UpdateClassification(params types.LogicParams, data frtypes.MapStr, cond condition.Condition) error {

	return nil
}

func (cli *core) UpdateObject(params types.LogicParams, data frtypes.MapStr, cond condition.Condition) error {
	return nil
}

func (cli *core) UpdateObjectAttribute(params types.LogicParams, data frtypes.MapStr, cond condition.Condition) error {
	return nil
}

func (cli *core) UpdateObjectGroup(params types.LogicParams, data frtypes.MapStr, cond condition.Condition) error {
	return nil
}

func (cli *core) UpdateInst(params types.LogicParams, data frtypes.MapStr, cond condition.Condition) error {
	return nil
}

func (cli *core) UpdateAssociation(params types.LogicParams, data frtypes.MapStr, cond condition.Condition) error {
	return nil
}
