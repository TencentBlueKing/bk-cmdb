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

type AuditOperationInterface interface {
	Query(params types.ContextParams, query metadata.QueryInput) (interface{}, error)
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

func (a *audit) Query(params types.ContextParams, query metadata.QueryInput) (interface{}, error) {
	rsp, err := a.clientSet.CoreService().Audit().SearchAuditLog(context.Background(), params.Header, query)
	if nil != err {
		blog.Errorf("[audit] failed request audit controller, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rsp.Result {
		blog.Errorf("[audit] failed request audit controller, error info is %s, rid: %s", rsp.ErrMsg, params.ReqID)
		return nil, params.Err.New(common.CCErrAuditTakeSnapshotFailed, rsp.ErrMsg)
	}

	for index := range rsp.Data.Info {
		if desc := params.Lang.Language("auditlog_" + rsp.Data.Info[index].OpDesc); len(desc) > 0 {
			rsp.Data.Info[index].OpDesc = desc
		}
	}

	return rsp.Data, nil
}
