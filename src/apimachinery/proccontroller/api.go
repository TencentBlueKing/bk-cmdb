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

package proccontroller

import (
	"context"
	"net/http"

    "configcenter/src/common/metadata"
)

func (p *procctrl) CreateProc2Module(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/module"

	err = p.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (p *procctrl) GetProc2Module(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.ProcModuleResult, err error) {
	resp = new(metadata.ProcModuleResult)
	subPath := "/module/search"

	err = p.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (p *procctrl) DeleteProc2Module(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/module"

	err = p.client.Delete().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (p *procctrl) CreateConfTemp(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error) {
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

func (p *procctrl) UpdateConfTemp(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error) {
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

func (p *procctrl) DeleteConfTemp(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error) {
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

func (p *procctrl) QueryConfTemp(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error) {
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

func (p *procctrl) CreateProcInstanceModel(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error) {
    resp = new(metadata.Response)
    subPath := "/instance/model"
    
    err = p.client.Post().
        WithContext(ctx).
        Body(dat).
        SubResource(subPath).
        WithHeaders(h).
        Do().
        Into(resp)
    
    return
}

func (p *procctrl) DeleteProcInstanceModel(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error) {
    resp = new(metadata.Response)
    subPath := "/instance/model"

    err = p.client.Delete().
        WithContext(ctx).
        Body(dat).
        SubResource(subPath).
        WithHeaders(h).
        Do().
        Into(resp)

    return
}

func (p *procctrl) GetProcInstanceModel(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.ProcInstModelResult, err error) {
    resp = new(metadata.ProcInstModelResult)
    subPath := "/instance/model/search"
    
    err = p.client.Post().
        WithContext(ctx).
        Body(dat).
        SubResource(subPath).
        WithHeaders(h).
        Do().
        Into(resp)
    
    return 
}

