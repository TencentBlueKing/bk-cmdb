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
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
)

func (s *coreService) TransferHostToDefaultModule(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	inputData := &metadata.TransferHostToInnerModule{}
	if err := data.MarshalJSONInto(inputData); nil != err {
		blog.Errorf("TransferHostToDefaultModule MarshalJSONInto error, err:%s,input:%v,rid:%s", err.Error(), data, params.ReqID)
		return nil, err
	}
	exceptionArr, err := s.core.HostOperation().TransferHostToInnerModule(params, inputData)
	if err != nil {
		blog.ErrorJSON("TransferHostToDefaultModule  error. err:%s, exception:%s, rid:%s", err.Error(), exceptionArr, params.ReqID)
		return exceptionArr, err
	}
	return nil, nil
}

func (s *coreService) TransferHostModule(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	inputData := &metadata.HostsModuleRelation{}
	if err := data.MarshalJSONInto(inputData); nil != err {
		blog.Errorf("TransferHostModule MarshalJSONInto error, err:%s,input:%v,rid:%s", err.Error(), data, params.ReqID)
		return nil, err
	}
	exceptionArr, err := s.core.HostOperation().TransferHostModule(params, inputData)
	if err != nil {
		blog.ErrorJSON("TransferHostModule  error. err:%s, exception:%s, rid:%s", err.Error(), exceptionArr, params.ReqID)
		return exceptionArr, err
	}
	return nil, nil
}

func (s *coreService) TransferHostCrossBusiness(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	inputData := &metadata.TransferHostsCrossBusinessRequest{}
	if err := data.MarshalJSONInto(inputData); nil != err {
		blog.Errorf("TransferHostCrossBusiness MarshalJSONInto error, err:%s,input:%v,rid:%s", err.Error(), data, params.ReqID)
		return nil, err
	}
	exceptionArr, err := s.core.HostOperation().TransferHostCrossBusiness(params, inputData)
	if err != nil {
		blog.ErrorJSON("TransferHostCrossBusiness  error. err:%s, input:%s, exception:%s, rid:%s", err.Error(), data, exceptionArr, params.ReqID)
		return exceptionArr, err
	}
	return nil, nil
}

func (s *coreService) GetHostModuleRelation(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	inputData := &metadata.HostModuleRelationRequest{}
	if err := data.MarshalJSONInto(inputData); nil != err {
		blog.Errorf("GetHostModuleRelation MarshalJSONInto error, err:%s,input:%v,rid:%s", err.Error(), data, params.ReqID)
		return nil, err
	}
	relationArr, err := s.core.HostOperation().GetHostModuleRelation(params, inputData)
	if err != nil {
		blog.ErrorJSON("GetHostModuleRelation  error. err:%s, rid:%s", err.Error(), params.ReqID)
		return nil, err
	}
	return relationArr, nil
}

func (s *coreService) DeleteHost(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	inputData := &metadata.DeleteHostRequest{}
	if err := data.MarshalJSONInto(inputData); nil != err {
		blog.Errorf("DeleteHost MarshalJSONInto error, err:%s,input:%v,rid:%s", err.Error(), data, params.ReqID)
		return nil, err
	}
	exceptionArr, err := s.core.HostOperation().DeleteHost(params, inputData)
	if err != nil {
		blog.ErrorJSON("DeleteHost  error. err:%s, exception:%s, rid:%s", err.Error(), exceptionArr, params.ReqID)
		return exceptionArr, err
	}
	return nil, nil
}

func (s *coreService) Identifier(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	inputData := &metadata.SearchHostIdentifierParam{}
	if err := data.MarshalJSONInto(inputData); nil != err {
		blog.Errorf("Identifier MarshalJSONInto error, err:%s,input:%#v,rid:%s", err.Error(), data, params.ReqID)
		return nil, err
	}
	hostIdentifierArr, err := s.core.HostOperation().Identifier(params, inputData)

	if err != nil {
		blog.InfoJSON("Identifier host identifier handle error. err:%s, input:%s, rid:%s", err.Error(), inputData, params.ReqID)
		return nil, err
	}

	return metadata.SearchHostIdentifierData{Info: hostIdentifierArr, Count: len(hostIdentifierArr)}, nil
}
