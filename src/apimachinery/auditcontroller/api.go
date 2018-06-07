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
    
    "configcenter/src/apimachinery/util"
    "configcenter/src/common/core/cc/api"
    "configcenter/src/source_controller/common/commondata"
)

func(t *auditctl) AddBusinessLog(ctx context.Context, businessID string, user string, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error) {
    resp = new(api.BKAPIRsp)
    subPath := fmt.Sprintf("/app/%s/%s/%s", h.OwnerID, businessID, user)

    err = t.client.Post().
        WithContext(ctx).
        Body(dat).
        SubResource(subPath).
        WithHeaders(h.ToHeader()).
        Do().
        Into(resp)
    return
}

func(t *auditctl) GetAuditLog(ctx context.Context, h util.Headers, opt *commondata.ObjQueryInput) (resp *api.BKAPIRsp, err error) {
    resp = new(api.BKAPIRsp)
    subPath := "/search"

    err = t.client.Post().
        WithContext(ctx).
        Body(opt).
        SubResource(subPath).
        WithHeaders(h.ToHeader()).
        Do().
        Into(resp)
    return
}

func(t *auditctl) AddHostLog(ctx context.Context, businessID string, user string, h util.Headers, log interface{}) (resp *api.BKAPIRsp, err error) {
    resp = new(api.BKAPIRsp)
    subPath := fmt.Sprintf("/host/%s/%s/%s", h.OwnerID, businessID, user)

    err = t.client.Post().
        WithContext(ctx).
        Body(log).
        SubResource(subPath).
        WithHeaders(h.ToHeader()).
        Do().
        Into(resp)
    return
}

func(t *auditctl) AddHostLogs(ctx context.Context, businessID string, user string, h util.Headers, logs interface{}) (resp *api.BKAPIRsp, err error) {
    resp = new(api.BKAPIRsp)
    subPath := fmt.Sprintf("/hosts/%s/%s/%s", h.OwnerID, businessID, user)

    err = t.client.Post().
        WithContext(ctx).
        Body(logs).
        SubResource(subPath).
        WithHeaders(h.ToHeader()).
        Do().
        Into(resp)
    return
}

func(t *auditctl) AddModuleLog(ctx context.Context, businessID string, user string, h util.Headers, log interface{}) (resp *api.BKAPIRsp, err error) {
    resp = new(api.BKAPIRsp)
    subPath := fmt.Sprintf("/module/%s/%s/%s", h.OwnerID, businessID, user)

    err = t.client.Post().
        WithContext(ctx).
        Body(log).
        SubResource(subPath).
        WithHeaders(h.ToHeader()).
        Do().
        Into(resp)
    return
}

func(t *auditctl) AddModuleLogs(ctx context.Context, businessID string, user string, h util.Headers, logs interface{}) (resp *api.BKAPIRsp, err error) {
    resp = new(api.BKAPIRsp)
    subPath := fmt.Sprintf("/modules/%s/%s/%s", h.OwnerID, businessID, user)

    err = t.client.Post().
        WithContext(ctx).
        Body(logs).
        SubResource(subPath).
        WithHeaders(h.ToHeader()).
        Do().
        Into(resp)
    return
}

func(t *auditctl) AddObjectLog(ctx context.Context, businessID string, user string, h util.Headers, log interface{}) (resp *api.BKAPIRsp, err error) {
    resp = new(api.BKAPIRsp)
    subPath := fmt.Sprintf("/obj/%s/%s/%s", h.OwnerID, businessID, user)

    err = t.client.Post().
        WithContext(ctx).
        Body(log).
        SubResource(subPath).
        WithHeaders(h.ToHeader()).
        Do().
        Into(resp)
    return
}

func(t *auditctl) AddObjectLogs(ctx context.Context, businessID string, user string, h util.Headers, logs interface{}) (resp *api.BKAPIRsp, err error) {
    resp = new(api.BKAPIRsp)
    subPath := fmt.Sprintf("/objs/%s/%s/%s", h.OwnerID, businessID, user)

    err = t.client.Post().
        WithContext(ctx).
        Body(logs).
        SubResource(subPath).
        WithHeaders(h.ToHeader()).
        Do().
        Into(resp)
    return
}

func(t *auditctl) AddProcLog(ctx context.Context, businessID string, user string, h util.Headers, log interface{}) (resp *api.BKAPIRsp, err error) {
    resp = new(api.BKAPIRsp)
    subPath := fmt.Sprintf("/proc/%s/%s/%s", h.OwnerID, businessID, user)

    err = t.client.Post().
        WithContext(ctx).
        Body(log).
        SubResource(subPath).
        WithHeaders(h.ToHeader()).
        Do().
        Into(resp)
    return
}

func(t *auditctl) AddProcLogs(ctx context.Context, businessID string, user string, h util.Headers, logs interface{}) (resp *api.BKAPIRsp, err error) {
    resp = new(api.BKAPIRsp)
    subPath := fmt.Sprintf("/procs/%s/%s/%s", h.OwnerID, businessID, user)

    err = t.client.Post().
        WithContext(ctx).
        Body(logs).
        SubResource(subPath).
        WithHeaders(h.ToHeader()).
        Do().
        Into(resp)
    return
}

func(t *auditctl) AddSetLog(ctx context.Context, businessID string, user string, h util.Headers, log interface{}) (resp *api.BKAPIRsp, err error) {
    resp = new(api.BKAPIRsp)
    subPath := fmt.Sprintf("/set/%s/%s/%s", h.OwnerID, businessID, user)

    err = t.client.Post().
        WithContext(ctx).
        Body(log).
        SubResource(subPath).
        WithHeaders(h.ToHeader()).
        Do().
        Into(resp)
    return
}

func(t *auditctl) AddSetLogs(ctx context.Context, businessID string, user string, h util.Headers, logs interface{}) (resp *api.BKAPIRsp, err error) {
    resp = new(api.BKAPIRsp)
    subPath := fmt.Sprintf("/sets/%s/%s/%s", h.OwnerID, businessID, user)

    err = t.client.Post().
        WithContext(ctx).
        Body(logs).
        SubResource(subPath).
        WithHeaders(h.ToHeader()).
        Do().
        Into(resp)
    return
}
