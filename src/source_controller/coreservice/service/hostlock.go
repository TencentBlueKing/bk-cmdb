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

package service

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
)

func (s *coreService) LockHost(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	input := new(metadata.HostLockRequest)
	if err := data.MarshalJSONInto(input); err != nil {
		blog.Errorf("LockHost failed, decode body failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommHTTPReadBodyFailed)
	}

	err := s.core.HostOperation().LockHost(params, input)
	if nil != err {
		blog.Errorf("LockHost failed, lock host handle failed, err: %+v, input:%+v, rid:%s", err, input, params.ReqID)
		return nil, err
	}

	return nil, nil
}

func (s *coreService) UnlockHost(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	input := new(metadata.HostLockRequest)
	if err := data.MarshalJSONInto(input); err != nil {
		blog.Errorf("UnlockHost failed, decode body failed, err: %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommHTTPReadBodyFailed)
	}
	err := s.core.HostOperation().UnlockHost(params, input)
	if nil != err {
		blog.Errorf("UnlockHost failed, unlock host handle failed, err: %s, input:%+v, rid:%s", err, input, params.ReqID)
		return nil, err
	}

	return nil, nil
}

func (s *coreService) QueryLockHost(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	input := new(metadata.QueryHostLockRequest)
	if err := data.MarshalJSONInto(input); err != nil {
		blog.Errorf("QueryLockHost failed, decode body failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommHTTPReadBodyFailed)
	}
	hostLockArr, err := s.core.HostOperation().QueryHostLock(params, input)
	if nil != err {
		blog.Errorf("QueryLockHost failed, query host handle failed, err: %s, input:%+v, rid: %s", err.Error(), input, params.ReqID)
		return nil, err
	}
	result := metadata.HostLockQueryResponse{}
	result.Data.Info = hostLockArr
	result.Data.Count = int64(len(hostLockArr))
	return result.Data, nil
}
