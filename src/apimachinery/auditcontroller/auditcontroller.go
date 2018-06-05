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
    "fmt"
    "context"
    
    "configcenter/src/apimachinery/rest"
    "configcenter/src/apimachinery/util"
    "configcenter/src/common/core/cc/api"
    "configcenter/src/source_controller/common/commondata"
)

type AuditCtrlInterface interface {
    AddBusinessLog(ctx context.Context, businessID string, user string, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error)
    GetAuditLog(ctx context.Context, h util.Headers, opt *commondata.ObjQueryInput) (resp *api.BKAPIRsp, err error)
    
    AddHostLog(ctx context.Context, businessID string, user string, h util.Headers, log interface{}) (resp *api.BKAPIRsp, err error)
    AddHostLogs(ctx context.Context, businessID string, user string, h util.Headers, logs interface{}) (resp *api.BKAPIRsp, err error)
   
    AddModuleLog(ctx context.Context, businessID string, user string, h util.Headers, log interface{}) (resp *api.BKAPIRsp, err error)
    AddModuleLogs(ctx context.Context, businessID string, user string, h util.Headers, logs interface{}) (resp *api.BKAPIRsp, err error)
    
    AddObjectLog(ctx context.Context, businessID string, user string, h util.Headers, log interface{}) (resp *api.BKAPIRsp, err error)
    AddObjectLogs(ctx context.Context, businessID string, user string, h util.Headers, logs interface{}) (resp *api.BKAPIRsp, err error)
    
    AddProcLog(ctx context.Context, businessID string, user string, h util.Headers, log interface{}) (resp *api.BKAPIRsp, err error)
    AddProcLogs(ctx context.Context, businessID string, user string, h util.Headers, logs interface{}) (resp *api.BKAPIRsp, err error)
   
    AddSetLog(ctx context.Context, businessID string, user string, h util.Headers, log interface{}) (resp *api.BKAPIRsp, err error)
    AddSetLogs(ctx context.Context, businessID string, user string, h util.Headers, logs interface{}) (resp *api.BKAPIRsp, err error)
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