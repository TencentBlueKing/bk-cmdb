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
	"configcenter/src/apimachinery"

	"configcenter/src/scene_server/topo_server/core/compatiblev2"
	"configcenter/src/scene_server/topo_server/core/types"
)

// CompatibleV2OperationInterface compatibleV2 methods
type CompatibleV2OperationInterface interface {
	Business(params types.ContextParams) compatiblev2.BusinessInterface
	Module(params types.ContextParams) compatiblev2.ModuleInterface
	Set(params types.ContextParams) compatiblev2.SetInterface
}

// NewCompatibleV2Operation create a new compatiblev2 operation instance
func NewCompatibleV2Operation(client apimachinery.ClientSetInterface) CompatibleV2OperationInterface {
	return &compatiblev2Operation{
		client: client,
	}
}

type compatiblev2Operation struct {
	client apimachinery.ClientSetInterface
}

func (c *compatiblev2Operation) Business(params types.ContextParams) compatiblev2.BusinessInterface {
	return compatiblev2.NewBusiness(params, c.client)
}
func (c *compatiblev2Operation) Module(params types.ContextParams) compatiblev2.ModuleInterface {
	return compatiblev2.NewModule(params, c.client)
}

func (c *compatiblev2Operation) Set(params types.ContextParams) compatiblev2.SetInterface {
	return compatiblev2.NewSet(params, c.client)
}
