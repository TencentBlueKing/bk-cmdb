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

package operation

import (
	"configcenter/src/scene_server/topo_server/core/compatiblev2"
	"configcenter/src/scene_server/topo_server/core/types"
)

// CompatibleV2OperationInterface compatibleV2 methods
type CompatibleV2OperationInterface interface {
	Business(params types.LogicParams) compatiblev2.BusinessInterface
	Module(params types.LogicParams) compatiblev2.ModuleInterface
	Set(params types.LogicParams) compatiblev2.SetInterface
}

// NewCompatibleV2Operation create a new compatiblev2 operation instance
func NewCompatibleV2Operation() CompatibleV2OperationInterface {
	return &compatiblev2Operation{}
}

type compatiblev2Operation struct {
}

func (c *compatiblev2Operation) Business(params types.LogicParams) compatiblev2.BusinessInterface {
	return nil
}
func (c *compatiblev2Operation) Module(params types.LogicParams) compatiblev2.ModuleInterface {
	return nil
}

func (c *compatiblev2Operation) Set(params types.LogicParams) compatiblev2.SetInterface {
	return nil
}
