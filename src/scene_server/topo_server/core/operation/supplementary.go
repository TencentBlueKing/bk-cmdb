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
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

// Supplementary supplementary methods
type Supplementary interface {
	Audit(params types.ContextParams, client apimachinery.ClientSetInterface, obj model.Object, inst InstOperationInterface) AuditInterface
	Validator(inst InstOperationInterface) ValidatorInterface
}

// NewSupplementary create a supplementary instance
func NewSupplementary() Supplementary {
	return &supplementary{}
}

type supplementary struct {
}

func (s *supplementary) Audit(params types.ContextParams, client apimachinery.ClientSetInterface, obj model.Object, inst InstOperationInterface) AuditInterface {
	return &auditLog{
		params: params,
		client: client,
		inst:   inst,
		obj:    obj,
	}
}

func (s *supplementary) Validator(inst InstOperationInterface) ValidatorInterface {
	return &valid{
		inst: inst,
	}
}
