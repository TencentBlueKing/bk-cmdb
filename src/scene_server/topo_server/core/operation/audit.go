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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
)

const CCTimeTypeParseFlag = "cc_time_type"

type AuditOperationInterface interface {
	Query(params types.ContextParams, data mapstr.MapStr) (interface{}, error)
}

// NewAuditOperation create a new inst operation instance
func NewAuditOperation(client apimachinery.ClientSetInterface) AuditOperationInterface {
	return &audit{
		clientSet: client,
	}
}

type audit struct {
	clientSet apimachinery.ClientSetInterface
}

func (a *audit) TranslateOpLanguage(params types.ContextParams, input interface{}) mapstr.MapStr {

	data, err := mapstr.NewFromInterface(input)
	if nil != err {
		blog.Errorf("failed to transate, error info is %s", err.Error())
		return data
	}

	info, err := data.MapStrArray("info")
	if nil != err {
		return data
	}

	for _, row := range info {

		opDesc, err := row.String(common.BKOpDescField)
		if nil != err {
			continue
		}
		newDesc := params.Lang.Language("auditlog_" + opDesc)
		if "" == newDesc {
			continue
		}
		row.Set(common.BKOpDescField, newDesc)
	}
	return data
}

func (a *audit) Query(params types.ContextParams, data mapstr.MapStr) (interface{}, error) {

	query := &metadata.QueryInput{}
	if err := data.MarshalJSONInto(query); nil != err {
		blog.Errorf("[audit] failed to parse the input (%#v), error info is %s", data, err.Error())
		return nil, params.Err.New(common.CCErrCommJSONUnmarshalFailed, err.Error())
	}

	iConds := query.Condition
	if nil == iConds {
		query.Condition = common.KvMap{common.BKOwnerIDField: params.SupplierAccount}
	} else {
		conds := iConds.(map[string]interface{})
		times, ok := conds[common.BKOpTimeField].([]interface{})
		if ok {
			if 2 != len(times) {
				blog.Error("search operation log input params times error, info: %v", times)
				return nil, params.Err.Error(common.CCErrCommParamsInvalid)
			}

			conds[common.BKOpTimeField] = common.KvMap{"$gte": times[0], "$lte": times[1], CCTimeTypeParseFlag: "1"}
			//delete(conds, "Time")
		}
		conds[common.BKOwnerIDField] = params.SupplierAccount
		query.Condition = conds
	}
	if 0 == query.Limit {
		query.Limit = common.BKDefaultLimit
	}

	rsp, err := a.clientSet.AuditController().GetAuditLog(context.Background(), params.Header, query)
	if nil != err {
		blog.Errorf("[audit] failed request audit conroller, error info is %s", err.Error())
		return nil, params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rsp.Result {
		blog.Errorf("[audit] failed request audit controller, error info is %s", rsp.ErrMsg)
		return nil, params.Err.New(common.CCErrAuditTakeSnapshotFaile, rsp.ErrMsg)
	}

	return a.TranslateOpLanguage(params, rsp.Data), nil
}
