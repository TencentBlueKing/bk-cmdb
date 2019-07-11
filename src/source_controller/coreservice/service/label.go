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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/selector"
	"configcenter/src/source_controller/coreservice/core"
)

func (s *coreService) AddLabels(ctx core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	inputData := selector.LabelAddRequest{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		blog.InfoJSON("RemoveLabels failed, MarshalJSONInto failed, data: %s, err: %s, rid: %s", data, err, ctx.ReqID)
		return nil, ctx.Error.CCError(common.CCErrCommJSONUnmarshalFailed)
	}
	if err := s.core.LabelOperation().AddLabel(ctx, inputData.TableName, inputData.Option); err != nil {
		blog.Errorf("AddLabels failed, table: %s, option: %+v, err: %s, rid: %s", inputData.TableName, inputData.Option, err.Error(), ctx.ReqID)
		return nil, err
	}
	return nil, nil
}

func (s *coreService) RemoveLabels(ctx core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	inputData := selector.LabelRemoveRequest{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		blog.InfoJSON("RemoveLabels failed, MarshalJSONInto failed, data: %s, err: %s, rid: %s", data, err, ctx.ReqID)
		return nil, ctx.Error.CCError(common.CCErrCommJSONUnmarshalFailed)
	}
	if err := s.core.LabelOperation().RemoveLabel(ctx, inputData.TableName, inputData.Option); err != nil {
		blog.Errorf("RemoveLabels failed, table: %s, option: %+v, err: %s, rid: %s", inputData.TableName, inputData.Option, err.Error(), ctx.ReqID)
		return nil, err
	}
	return nil, nil
}
