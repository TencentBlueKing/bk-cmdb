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

package auditlog

import (
	"context"
	"net/http"

	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

// SaveAuditLog api of save audit log
func (inst *auditlog) SaveAuditLog(ctx context.Context, h http.Header, logs ...metadata.AuditLog) errors.CCErrorCoder {

	resp := new(metadata.Response)
	subPath := "/create/auditlog"

	err := inst.client.Post().
		WithContext(ctx).
		Body(metadata.CreateAuditLogParam{Data: logs}).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return err
	}

	return nil
}

// SearchAuditLog api of search audit log
func (inst *auditlog) SearchAuditLog(ctx context.Context, h http.Header, param metadata.QueryCondition) (
	*metadata.AuditQueryResult, errors.CCErrorCoder) {

	resp := new(metadata.AuditQueryResponse)
	subPath := "/read/auditlog"

	err := inst.client.Post().
		WithContext(ctx).
		Body(param).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return nil, err
	}

	return resp.Data, nil
}
