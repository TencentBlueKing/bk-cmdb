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
	"strings"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
)

type IdentifierOperationInterface interface {
	SearchIdentifier(params types.ContextParams, objType string, param *metadata.SearchIdentifierParam) (*metadata.SearchHostIdentifierResult, error)
}

func NewIdentifier(client apimachinery.ClientSetInterface) IdentifierOperationInterface {
	return &identifier{clientSet: client}
}

type identifier struct {
	clientSet apimachinery.ClientSetInterface
}

func (g *identifier) SearchIdentifier(params types.ContextParams, objType string, param *metadata.SearchIdentifierParam) (*metadata.SearchHostIdentifierResult, error) {
	cond := condition.CreateCondition()

	or := []mapstr.MapStr{
		{
			common.BKHostInnerIPField: map[string]interface{}{
				common.BKDBIN: param.IP.Data,
			},
		}, {
			common.BKHostOuterIPField: map[string]interface{}{
				common.BKDBIN: param.IP.Data,
			},
		},
	}

	cond.NewOR().MapStrArr(or)
	if param.IP.CloudID != nil {
		cond.Field(common.BKCloudIDField).In(param.IP.CloudID)
	}

	if param.Page.Limit > common.BKMaxPageSize {
		return nil, params.Err.CCError(common.CCErrCommOverLimit)
	}
	if param.Page.Limit == 0 {
		param.Page.Limit = common.BKMaxPageSize
	}

	sortArr := make([]metadata.SearchSort, 0)
	if len(param.Page.Sort) != 0 {
		for _, field := range strings.Split(param.Page.Sort, ",") {
			field = strings.TrimSpace(field)
			if field == "" {
				continue
			}
			var isDesc bool
			switch field[0] {
			case '-':
				field = strings.TrimLeft(field, "-")
				isDesc = true
			case '+':
				field = strings.TrimLeft(field, "+")
			}
			sortArr = append(sortArr, metadata.SearchSort{
				IsDsc: isDesc,
				Field: field,
			})
		}
	}

	hostQuery := &metadata.QueryCondition{
		Condition: cond.ToMapStr(),
		Fields:    []string{common.BKHostIDField},
		Limit: metadata.SearchLimit{
			Offset: int64(param.Page.Start),
			Limit:  int64(param.Page.Limit),
		},
		SortArr: sortArr,
	}
	hostRet, err := g.clientSet.CoreService().Instance().ReadInstance(context.Background(), params.Header, common.BKInnerObjIDHost, hostQuery)
	if nil != err {
		blog.ErrorJSON("[identifier] ReadInstance query host  http do error. error:%s, input:%s,  rid:%s", err.Error(), params, params.ReqID)
		return nil, params.Err.CCErrorf(common.CCErrCommHTTPDoRequestFailed)
	}
	if !hostRet.Result {
		blog.ErrorJSON("[identifier] ReadInstance query host  http reply error. result:%s, input:%s, condition:%s, rid:%s", hostRet, params, hostQuery, params.ReqID)
		return nil, params.Err.New(hostRet.Code, hostRet.ErrMsg)
	}
	var hostIDs []int64
	for _, hostInfo := range hostRet.Data.Info {
		hostID, err := hostInfo.Int64(common.BKHostIDField)
		if err != nil {
			blog.ErrorJSON("[identifier] ReadInstance host info bk_host_id not int . error:%s, input:%s, host info:%s,  rid:%s", err.Error(), params, hostInfo, params.ReqID)
			// format: `convert %s  field %s to %s error %s`
			return nil, params.Err.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDHost, common.BKHostIDField, "int64", err.Error())
		}
		hostIDs = append(hostIDs, hostID)
	}
	queryHostIdentifier := &metadata.SearchHostIdentifierParam{HostIDs: hostIDs}
	rsp, err := g.clientSet.CoreService().Host().FindIdentifier(context.Background(), params.Header, queryHostIdentifier)
	if nil != err {
		blog.ErrorJSON("[identifier]  SearchIdentifier http do error. error:%s, input:%s,  rid:%s", err.Error(), params, params.ReqID)
		return nil, params.Err.CCErrorf(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.ErrorJSON("[identifier]  SearchIdentifier http reply error , reply:%s, input:%s, condition:%s, rid:%s", rsp, param, queryHostIdentifier, params.ReqID)
		return nil, params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return rsp, nil
}
