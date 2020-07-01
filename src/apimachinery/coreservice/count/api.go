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

package count

import (
	"context"
	"net/http"

	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (s *count) GetCountByFilter(ctx context.Context, h http.Header, table string, filters []map[string]interface{}) ([]int64, errors.CCErrorCoder) {
	rid := util.ExtractRequestIDFromContext(ctx)

	resp := &struct {
		metadata.BaseResp `json:",inline"`
		Data              []int64 `json:"data"`
	}{}
	body := struct {
		Table   string                   `json:"table"`
		Filters []map[string]interface{} `json:"filters"`
	}{
		Table:   table,
		Filters: filters,
	}
	subPath := "/find/resource/count"

	httpDoErr := s.client.Post().
		WithContext(ctx).
		Body(body).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	if httpDoErr != nil {
		blog.Errorf("GetCountByFilter failed, http request failed, err: %+v, rid: %s", httpDoErr, rid)
		return nil, errors.CCHttpError
	}
	if resp.Result == false || resp.Code != 0 {
		return nil, errors.New(resp.Code, resp.ErrMsg)
	}

	return resp.Data, nil
}
