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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

type AuditOperationInterface interface {
	Query(kit *rest.Kit, query metadata.QueryInput) (interface{}, error)
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

func (a *audit) Query(kit *rest.Kit, query metadata.QueryInput) (interface{}, error) {
	rsp, err := a.clientSet.CoreService().Audit().SearchAuditLog(kit.Ctx, kit.Header, query)
	if nil != err {
		blog.Errorf("[audit] failed to request core service, error info is %s, rid: %s", err.Error(), kit.Rid)
		return nil, kit.CCError.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rsp.Result {
		blog.Errorf("[audit] search audit log failed, error info is %s, rid: %s", rsp.ErrMsg, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrAuditSelectFailed)
	}

	return rsp.Data, nil
}
