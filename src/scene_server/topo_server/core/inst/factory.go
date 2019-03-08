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
	"configcenter/src/common/mapstr"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

// New create a new inst factory
func New(clientSet apimachinery.ClientSetInterface) Factory {
	return &factory{
		clientSet: clientSet,
	}
}

// CreateInst convert the inst into the Inst interface
func CreateInst(params types.ContextParams, clientSet apimachinery.ClientSetInterface, obj model.Object, instItems []mapstr.MapStr) []Inst {
	results := make([]Inst, 0)
	for _, item := range instItems {
		tmpInst := &inst{
			clientSet: clientSet,
			params:    params,
			target:    obj,
			datas:     mapstr.New(),
		}
		tmpInst.SetValues(item)
		results = append(results, tmpInst)
	}
	return results
}

type factory struct {
	clientSet apimachinery.ClientSetInterface
}

func (cli *factory) CreateInst(params types.ContextParams, obj model.Object) Inst {
	return &inst{
		datas:     mapstr.New(),
		params:    params,
		clientSet: cli.clientSet,
		target:    obj,
	}
}
