/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package inst

import (
	"context"
	"net/http"

	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

// CreateProject create project
func (t *instanceClient) CreateProject(ctx context.Context, h http.Header, opt *metadata.CreateProjectOption) (
	*metadata.ProjectDataResp, errors.CCErrorCoder) {

	resp := new(metadata.ProjectInstResp)
	subPath := "/createmany/project"

	err := t.client.Post().
		WithContext(ctx).
		Body(opt).
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

	return &resp.Data, nil
}

// UpdateProject update project
func (t *instanceClient) UpdateProject(ctx context.Context, h http.Header,
	opt *metadata.UpdateProjectOption) errors.CCErrorCoder {

	resp := new(metadata.Response)
	subPath := "/updatemany/project"

	err := t.client.Put().
		WithContext(ctx).
		Body(opt).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return errors.CCHttpError
	}

	return resp.CCError()
}

// SearchProject search project
func (t *instanceClient) SearchProject(ctx context.Context, h http.Header,
	opt *metadata.SearchProjectOption) (*metadata.InstResult, errors.CCErrorCoder) {

	resp := new(metadata.QueryInstResult)
	subPath := "/findmany/project"

	err := t.client.Post().
		WithContext(ctx).
		Body(opt).
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

	return &resp.Data, nil
}

// DeleteProject delete project
func (t *instanceClient) DeleteProject(ctx context.Context, h http.Header,
	opt *metadata.DeleteProjectOption) errors.CCErrorCoder {

	resp := new(metadata.Response)
	subPath := "/deletemany/project"

	err := t.client.Delete().
		WithContext(ctx).
		Body(opt).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return errors.CCHttpError
	}

	return resp.CCError()
}

// UpdateProjectID update project bk_project_id
func (t *instanceClient) UpdateProjectID(ctx context.Context, h http.Header,
	opt *metadata.UpdateProjectIDOption) errors.CCErrorCoder {

	resp := new(metadata.Response)
	subPath := "/update/project/bk_project_id"

	err := t.client.Put().
		WithContext(ctx).
		Body(opt).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return errors.CCHttpError
	}

	return resp.CCError()
}
