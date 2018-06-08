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

package process


import (
    "context"
    "fmt"

    "configcenter/src/apimachinery/util"
    "configcenter/src/common/core/cc/api"
)

func (p *process) GetProcessDetailByID(ctx context.Context, appID string, procID string, h util.Headers) (resp *api.BKAPIRsp, err error) {
    resp = new(api.BKAPIRsp)
    subPath := fmt.Sprintf("/%s/%s/%s", h.OwnerID, appID, procID)

    err = p.client.Get().
        WithContext(ctx).
        Body(nil).
        SubResource(subPath).
        WithHeaders(h.ToHeader()).
        Do().
        Into(resp)
    return
}

func (p *process) GetProcessBindModule(ctx context.Context, businessID string, procID string, h util.Headers) (resp *api.BKAPIRsp, err error) {
    resp = new(api.BKAPIRsp)
    subPath := fmt.Sprintf("/module/%s/%s/%s", h.OwnerID, businessID, procID)

    err = p.client.Get().
        WithContext(ctx).
        Body(nil).
        SubResource(subPath).
        WithHeaders(h.ToHeader()).
        Do().
        Into(resp)
    return
}

func (p *process) BindModuleProcess(ctx context.Context, businessID string, procID string, moduleName string, h util.Headers) (resp *api.BKAPIRsp, err error) {
    resp = new(api.BKAPIRsp)
    subPath := fmt.Sprintf("/module/%s/%s/%s/%s", h.OwnerID, businessID, procID, moduleName)

    err = p.client.Put().
        WithContext(ctx).
        Body(nil).
        SubResource(subPath).
        WithHeaders(h.ToHeader()).
        Do().
        Into(resp)
    return
}

func (p *process) DeleteModuleProcessBind(ctx context.Context, businessID string, procID string, moduleName string, h util.Headers) (resp *api.BKAPIRsp, err error) {
    resp = new(api.BKAPIRsp)
    subPath := fmt.Sprintf("/module/%s/%s/%s/%s", h.OwnerID, businessID, procID, moduleName)

    err = p.client.Delete().
        WithContext(ctx).
        Body(nil).
        SubResource(subPath).
        WithHeaders(h.ToHeader()).
        Do().
        Into(resp)
    return
}

func (p *process) CreateProcess(ctx context.Context, businessID string, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error) {
    resp = new(api.BKAPIRsp)
    subPath := fmt.Sprintf("/%s/%s", h.OwnerID, businessID)

    err = p.client.Post().
        WithContext(ctx).
        Body(dat).
        SubResource(subPath).
        WithHeaders(h.ToHeader()).
        Do().
        Into(resp)
    return
}

func (p *process) DeleteProcess(ctx context.Context, businessID string, procID string, h util.Headers) (resp *api.BKAPIRsp, err error) {
    resp = new(api.BKAPIRsp)
    subPath := fmt.Sprintf("/%s/%s/%s", h.OwnerID, businessID, procID)

    err = p.client.Delete().
        WithContext(ctx).
        Body(nil).
        SubResource(subPath).
        WithHeaders(h.ToHeader()).
        Do().
        Into(resp)
    return
}

func (p *process) SearchProcess(ctx context.Context, businessID string, h util.Headers) (resp *api.BKAPIRsp, err error) {
    resp = new(api.BKAPIRsp)
    subPath := fmt.Sprintf("/search/%s/%s", h.OwnerID, businessID)

    err = p.client.Post().
        WithContext(ctx).
        Body(nil).
        SubResource(subPath).
        WithHeaders(h.ToHeader()).
        Do().
        Into(resp)
    return
}

func (p *process) UpdateProcess(ctx context.Context, businessID string, procID string, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error) {
    resp = new(api.BKAPIRsp)
    subPath := fmt.Sprintf("/%s/%s/%s", h.OwnerID, businessID, procID)

    err = p.client.Put().
        WithContext(ctx).
        Body(dat).
        SubResource(subPath).
        WithHeaders(h.ToHeader()).
        Do().
        Into(resp)
    return
}



