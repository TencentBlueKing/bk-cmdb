/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"fmt"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/source_controller/coreservice/core"
)

func (s *coreService) SearchMainlineModelTopo(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	detail := struct {
		WithDetail bool `field:"with_detail"`
	}{}
	if err := mapstr.SetValueToStructByTags(&detail, data); err != nil {
		blog.Errorf("decode body %+v failed, err: %v, rid: %s", data, err, params.ReqID)
		return nil, fmt.Errorf("decode body %+v failed, err: %v", data, err)
	}

	result, err := s.core.TopoOperation().SearchMainlineModelTopo(params.Context, params.Header, detail.WithDetail)
	if err != nil {
		blog.Errorf("search mainline model topo failed, %+v, rid: %s", err, params.ReqID)
		return nil, fmt.Errorf("search mainline model topo failed, %+v", err)
	}
	return result, nil
}

func (s *coreService) SearchMainlineInstanceTopo(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	bkBizID := pathParams(common.BKAppIDField)
	if len(bkBizID) == 0 {
		return nil, fmt.Errorf("field %s empty", common.BKAppIDField)
	}
	bizID, err := strconv.ParseInt(bkBizID, 10, 64)
	if err != nil {
		blog.Errorf("field %s with value:%s invalid, %v, rid: %s", common.BKAppIDField, bkBizID, err, params.ReqID)
		return nil, fmt.Errorf("field %s with valued:%s invalid, %v", common.BKAppIDField, bkBizID, err)
	}

	blog.V(2).Infof("decode body %+v, rid: %s", data, params.ReqID)
	detail := struct {
		WithDetail bool `field:"with_detail"`
	}{}
	if err := mapstr.SetValueToStructByTags(&detail, data); err != nil {
		blog.Errorf("decode body %+v failed, err: %v, rid: %s", data, err, params.ReqID)
		return nil, fmt.Errorf("decode body %+v failed, err: %v", data, err)
	}
	blog.V(2).Infof("decode result: %+v, rid: %s", detail, params.ReqID)

	result, err := s.core.TopoOperation().SearchMainlineInstanceTopo(params.Context, params.Header, bizID, detail.WithDetail)
	if err != nil {
		blog.Errorf("search mainline instance topo by business:%d failed, %+v, rid: %s", bizID, err, params.ReqID)
		return nil, fmt.Errorf("search mainline instance topo by business:%d failed, %+v", bizID, err)
	}
	return result, nil
}
