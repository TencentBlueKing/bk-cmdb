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

package topographics

import (
	"context"
	"net/http"

	"configcenter/src/common/errors"

	"configcenter/src/common/metadata"
)

// SearchTopoGraphics search topo's graphics
func (t *meta) SearchTopoGraphics(ctx context.Context, h http.Header, dat *metadata.TopoGraphics) (
	[]metadata.TopoGraphics, error) {

	subPath := "/topographics/search"
	resp := new(metadata.SearchTopoGraphicsResult)

	err := t.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err = resp.CCError(); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// UpdateTopoGraphics update topo's graphics
func (t *meta) UpdateTopoGraphics(ctx context.Context, h http.Header, dat []metadata.TopoGraphics) error {
	data := map[string]interface{}{
		"data": dat,
	}
	subPath := "/topographics/update"
	resp := new(metadata.UpdateResult)
	err := t.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return errors.CCHttpError
	}

	if err = resp.CCError(); err != nil {
		return err
	}

	return nil
}
