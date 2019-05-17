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
	"net/http"

    "configcenter/src/common/metadata"
)

func (p *process) GetProcessDetailByID(ctx context.Context, ownerID string, appID string, procID string, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/%s/%s/%s", ownerID, appID, procID)

	err = p.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (p *process) GetProcessBindModule(ctx context.Context, ownerID string, businessID string, procID string, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/module/%s/%s/%s", ownerID, businessID, procID)

	err = p.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (p *process) BindModuleProcess(ctx context.Context, ownerID string, businessID string, procID string, moduleName string, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/module/%s/%s/%s/%s", ownerID, businessID, procID, moduleName)

	err = p.client.Put().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (p *process) DeleteModuleProcessBind(ctx context.Context, ownerID string, businessID string, procID string, moduleName string, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/module/%s/%s/%s/%s", ownerID, businessID, procID, moduleName)

	err = p.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (p *process) CreateProcess(ctx context.Context, ownerID string, businessID string, h http.Header, dat interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/%s/%s", ownerID, businessID)

	err = p.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (p *process) DeleteProcess(ctx context.Context, ownerID string, businessID string, procID string, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/%s/%s/%s", ownerID, businessID, procID)

	err = p.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (p *process) SearchProcess(ctx context.Context, ownerID string, businessID string, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/search/%s/%s", ownerID, businessID)

	err = p.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (p *process) UpdateProcess(ctx context.Context, ownerID string, businessID string, procID string, h http.Header, dat interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/%s/%s/%s", ownerID, businessID, procID)

	err = p.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (p *process) BatchUpdateProcess(ctx context.Context, ownerID, businessID string, h http.Header, dat interface{}) (resp *metadata.Response, err error) {
    resp = new(metadata.Response)
    subPath := fmt.Sprintf("/%s/%s", ownerID, businessID)

    err = p.client.Put().
        WithContext(ctx).
        Body(dat).
        SubResource(subPath).
        WithHeaders(h).
        Do().
        Into(resp)
    return
}

func (p *process) OperateProcessInstance(ctx context.Context, namespace string, h http.Header, dat interface{}) (resp *metadata.Response, err error) {
    resp = new(metadata.Response)
    subPath := fmt.Sprintf("/operate/%s/process", namespace)

    err = p.client.Post().
        WithContext(ctx).
        Body(dat).
        SubResource(subPath).
        WithHeaders(h).
        Do().
        Into(resp)
    return
}

func (p *process) QueryProcessOperateResult(ctx context.Context, namespace string, h http.Header, dat interface{}) (resp *metadata.Response, err error) {
    resp = new(metadata.Response)
    subPath := fmt.Sprintf("/operate/%s/process/taskresult", namespace)

    err = p.client.Post().
        WithContext(ctx).
        Body(dat).
        SubResource(subPath).
        WithHeaders(h).
        Do().
        Into(resp)
    return
}

func (p *process) CreateConfigTemp(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error) {
    resp = new(metadata.Response)
    subPath := "/conftemp"

    err = p.client.Post().
        WithContext(ctx).
        Body(dat).
        SubResource(subPath).
        WithHeaders(h).
        Do().
        Into(resp)
    return
}

func (p *process) UpdateConfigTemp(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error) {
    resp = new(metadata.Response)
    subPath := "/conftemp"

    err = p.client.Put().
        WithContext(ctx).
        Body(dat).
        SubResource(subPath).
        WithHeaders(h).
        Do().
        Into(resp)
    return
}

func (p *process) DeleteConfigTemp(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error) {
    resp = new(metadata.Response)
    subPath := "/conftemp"

    err = p.client.Delete().
        WithContext(ctx).
        Body(dat).
        SubResource(subPath).
        WithHeaders(h).
        Do().
        Into(resp)
    return
}

func (p *process) QueryConfigTemp(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error) {
    resp = new(metadata.Response)
    subPath := "/conftemp/search"

    err = p.client.Post().
        WithContext(ctx).
        Body(dat).
        SubResource(subPath).
        WithHeaders(h).
        Do().
        Into(resp)
    return
}
