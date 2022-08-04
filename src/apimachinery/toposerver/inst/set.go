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

package inst

import (
	"context"
	"net/http"

	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/paraparse"
)

// CreateSet TODO
func (t *instanceClient) CreateSet(ctx context.Context, appID int64, h http.Header, dat mapstr.MapStr) (mapstr.MapStr,
	errors.CCErrorCoder) {

	resp := new(metadata.CreateInstResult)
	subPath := "/set/%d"

	err := t.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath, appID).
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

// DeleteSet TODO
func (t *instanceClient) DeleteSet(ctx context.Context, appID, setID int64, h http.Header) errors.CCErrorCoder {
	resp := new(metadata.Response)
	subPath := "/set/%d/%d"

	err := t.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, appID, setID).
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

// UpdateSet TODO
func (t *instanceClient) UpdateSet(ctx context.Context, appID, setID int64, h http.Header,
	dat map[string]interface{}) errors.CCErrorCoder {

	resp := new(metadata.Response)
	subPath := "/set/%d/%d"

	err := t.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath, appID, setID).
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

// SearchSet TODO
func (t *instanceClient) SearchSet(ctx context.Context, ownerID string, appID string, h http.Header, s *params.SearchParams) (resp *metadata.SearchInstResult, err error) {
	resp = new(metadata.SearchInstResult)
	subPath := "/set/search/%s/%s"

	err = t.client.Post().
		WithContext(ctx).
		Body(s).
		SubResourcef(subPath, ownerID, appID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// SearchSetBatch TODO
func (t *instanceClient) SearchSetBatch(ctx context.Context, appID string, h http.Header, s *metadata.SearchInstBatchOption) (resp *metadata.MapArrayResponse, err error) {
	resp = new(metadata.MapArrayResponse)
	subPath := "/findmany/set/bk_biz_id/%s"

	err = t.client.Post().
		WithContext(ctx).
		Body(s).
		SubResourcef(subPath, appID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
