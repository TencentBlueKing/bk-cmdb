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
	"context"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
)

type IdentifierOperationInterface interface {
	SearchIdentifier(params types.ContextParams, objType string, param *metadata.SearchIdentifierParam) ([]metadata.HostIdentifier, error)
}

func NewIdentifier(client apimachinery.ClientSetInterface) IdentifierOperationInterface {
	return &identifier{clientSet: client}
}

type identifier struct {
	clientSet apimachinery.ClientSetInterface
}

func (g *identifier) SearchIdentifier(params types.ContextParams, objType string, param *metadata.SearchIdentifierParam) ([]metadata.HostIdentifier, error) {
	rsp, err := g.clientSet.ObjectController().Identifier().SearchIdentifier(context.Background(), params.Header, objType, param)
	if nil != err {
		return nil, err
	}

	if !rsp.Result {
		blog.Errorf("[identifier] failed to search the identifier , error info is %s", rsp.ErrMsg)
		return nil, params.Err.New(common.CCErrObjectSelectIdentifierFailed, rsp.ErrMsg)
	}

	return rsp.Data.Info, nil
}
