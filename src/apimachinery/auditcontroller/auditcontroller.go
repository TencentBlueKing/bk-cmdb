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

package auditcontroller

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common/metadata"
)

type AuditCtrlInterface interface {
	AddBusinessLog(ctx context.Context, ownerID string, businessID string, user string, h http.Header, dat interface{}) (resp *metadata.Response, err error)
	GetAuditLog(ctx context.Context, h http.Header, opt *metadata.QueryInput) (resp *metadata.Response, err error)

	AddHostLog(ctx context.Context, ownerID string, businessID string, user string, h http.Header, log interface{}) (resp *metadata.Response, err error)
	AddHostLogs(ctx context.Context, ownerID string, businessID string, user string, h http.Header, logs interface{}) (resp *metadata.Response, err error)

	AddModuleLog(ctx context.Context, ownerID string, businessID string, user string, h http.Header, log interface{}) (resp *metadata.Response, err error)
	AddModuleLogs(ctx context.Context, ownerID string, businessID string, user string, h http.Header, logs interface{}) (resp *metadata.Response, err error)

	AddObjectLog(ctx context.Context, ownerID string, businessID string, user string, h http.Header, log interface{}) (resp *metadata.Response, err error)
	AddObjectLogs(ctx context.Context, ownerID string, businessID string, user string, h http.Header, logs interface{}) (resp *metadata.Response, err error)

	AddProcLog(ctx context.Context, ownerID string, businessID string, user string, h http.Header, log interface{}) (resp *metadata.Response, err error)
	AddProcLogs(ctx context.Context, ownerID string, businessID string, user string, h http.Header, logs interface{}) (resp *metadata.Response, err error)

	AddSetLog(ctx context.Context, ownerID string, businessID string, user string, h http.Header, log interface{}) (resp *metadata.Response, err error)
	AddSetLogs(ctx context.Context, ownerID string, businessID string, user string, h http.Header, logs interface{}) (resp *metadata.Response, err error)
}

func NewAuditCtrlInterface(c *util.Capability, version string) AuditCtrlInterface {
	base := fmt.Sprintf("/audit/%s", version)
	return &auditctl{
		client: rest.NewRESTClient(c, base),
	}
}

type auditctl struct {
	client rest.ClientInterface
}
