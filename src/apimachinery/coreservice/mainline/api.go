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

package mainline

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/common/metadata"
	"configcenter/src/framework/core/errors"
)

func (m *mainline) SearchMainlineModelTopo(ctx context.Context, h http.Header, withDetail bool) (resp *metadata.TopoModelNode, err error) {
	ret := new(metadata.SearchTopoModelNodeResult)
	// resp = new(metadata.TopoModelNode)
	subPath := "/read/mainline/model"

	input := map[string]bool{}
	input["with_detail"] = withDetail

	err = m.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return nil, err
	}
	if ret.Result == false || ret.Code != 0 {
		return nil, errors.New(ret.ErrMsg)
	}

	return &ret.Data, nil
}

func (m *mainline) SearchMainlineInstanceTopo(ctx context.Context, h http.Header, bkBizID int64, withDetail bool) (resp *metadata.TopoInstanceNode, err error) {
	input := map[string]bool{}
	input["with_detail"] = withDetail

	ret := new(metadata.SearchTopoInstanceNodeResult)
	subPath := fmt.Sprintf("/read/mainline/instance/%d", bkBizID)
	err = m.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return nil, err
	}
	if ret.Result == false || ret.Code != 0 {
		return nil, errors.New(ret.ErrMsg)
	}

	return &ret.Data, nil
}
