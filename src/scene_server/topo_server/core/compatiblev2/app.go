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

package compatiblev2

import (
	"strings"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
)

// BusinessInterface business methods
type BusinessInterface interface {
	SearchAllApp(fields string, cond mapstr.MapStr) (*metadata.InstResult, error)
}

// NewBusiness create a new business instance
func NewBusiness(params types.ContextParams, client apimachinery.ClientSetInterface) BusinessInterface {
	return &business{
		params: params,
		client: client,
	}
}

type business struct {
	params types.ContextParams
	client apimachinery.ClientSetInterface
}

func (b *business) SearchAllApp(fields string, cond mapstr.MapStr) (*metadata.InstResult, error) {
	fieldArr := strings.Split(fields, ",")
	inputParam := metadata.QueryCondition{
		Fields: fieldArr,
		Limit: metadata.SearchLimit{
			Offset: 0,
			Limit:  common.BKNoLimit,
		},
		SortArr:   nil,
		Condition: cond,
	}
	rsp, err := b.client.CoreService().Instance().ReadInstance(b.params.Context, b.params.Header, common.BKInnerObjIDApp, &inputParam)
	if nil != err {
		blog.Errorf("[compatiblev2-biz] failed to request object controller, error info is %s, rid: %s", err.Error(), b.params.ReqID)
		return nil, err
	}

	if !rsp.Result {
		blog.Errorf("[compatiblev2-biz] failed to search the business, error info is %s, rid: %s", rsp.ErrMsg, b.params.ReqID)
		return nil, b.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	result := &metadata.InstResult{
		Count: rsp.Data.Count,
		Info:  rsp.Data.Info,
	}
	return result, nil
}
