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

package resourcedir

import (
	"context"
	"net/http"

	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

func (st *ResourceDirectory) CreateResourceDirectory(ctx context.Context, header http.Header, data map[string]interface{}) (*metadata.CreateOneDataResult, errors.CCErrorCoder) {
	ret := new(metadata.CreatedOneOptionResult)
	subPath := "/create/resource/directory"

	err := st.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("CreateResourceDirectory failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.Result == false || ret.Code != 0 {
		return nil, errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return &ret.Data, nil
}

func (st *ResourceDirectory) UpdateResourceDirectory(ctx context.Context, header http.Header, moduleID int64, option metadata.UpdateSetTemplateOption) (*metadata.UpdatedCount, errors.CCErrorCoder) {
	ret := new(metadata.UpdatedOptionResult)
	subPath := "/update/resource/directory/%d"

	err := st.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath, moduleID).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("UpdateResourceDirectory failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.Result == false || ret.Code != 0 {
		return nil, errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return &ret.Data, nil
}

func (st *ResourceDirectory) SearchResourceDirectory(ctx context.Context, header http.Header, data map[string]interface{}) (*metadata.InstDataInfo, errors.CCErrorCoder) {
	ret := new(metadata.QueryConditionResult)
	subPath := "/findmany/resource/directory"

	err := st.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("SearchResourceDirectory failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.Result == false || ret.Code != 0 {
		return nil, errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return &ret.Data, nil
}

func (st *ResourceDirectory) DeleteResourceDirectory(ctx context.Context, header http.Header, moduleID int64) (*metadata.DeletedCount, errors.CCErrorCoder) {
	ret := new(metadata.DeletedOptionResult)
	subPath := "/delete/resource/directory/%d"

	err := st.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, moduleID).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("DeleteResourceDirectory failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.Result == false || ret.Code != 0 {
		return nil, errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return &ret.Data, nil
}
