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
	"gopkg.in/redis.v5"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
)

func (s *coreService) TransferHostToInnerModule(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	inputData := &metadata.TransferHostToInnerModule{}
	if err := data.MarshalJSONInto(inputData); nil != err {
		blog.Errorf("TransferHostToDefaultModule MarshalJSONInto error, err:%s, input:%v, rid: %s", err.Error(), data, params.ReqID)
		return nil, err
	}
	exceptionArr, err := s.core.HostOperation().TransferToInnerModule(params, inputData)
	if err != nil {
		blog.ErrorJSON("TransferHostToDefaultModule  error. err:%s, exception:%s, rid:%s", err.Error(), exceptionArr, params.ReqID)
		return exceptionArr, err
	}
	return nil, nil
}

func (s *coreService) TransferHostToNormalModule(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	inputData := &metadata.HostsModuleRelation{}
	if err := data.MarshalJSONInto(inputData); nil != err {
		blog.Errorf("TransferHostModule MarshalJSONInto error, err:%s, input:%v, rid:%s", err.Error(), data, params.ReqID)
		return nil, err
	}
	exceptionArr, err := s.core.HostOperation().TransferToNormalModule(params, inputData)
	if err != nil {
		blog.ErrorJSON("TransferHostModule  error. err:%s, exception:%s, rid:%s", err.Error(), exceptionArr, params.ReqID)
		return exceptionArr, err
	}
	return nil, nil
}

func (s *coreService) TransferHostToAnotherBusiness(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	inputData := &metadata.TransferHostsCrossBusinessRequest{}
	if err := data.MarshalJSONInto(inputData); nil != err {
		blog.Errorf("TransferHostCrossBusiness MarshalJSONInto error, err:%s, input:%v, rid:%s", err.Error(), data, params.ReqID)
		return nil, err
	}
	exceptionArr, err := s.core.HostOperation().TransferToAnotherBusiness(params, inputData)
	if err != nil {
		blog.ErrorJSON("TransferHostCrossBusiness  error. err:%s, input:%s, exception:%s, rid:%s", err.Error(), data, exceptionArr, params.ReqID)
		return exceptionArr, err
	}
	return nil, nil
}

func (s *coreService) RemoveFromModule(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	inputData := &metadata.RemoveHostsFromModuleOption{}
	if err := data.MarshalJSONInto(inputData); nil != err {
		blog.Errorf("RemoveFromModule MarshalJSONInto error, err:%s, input:%v, rid:%s", err.Error(), data, params.ReqID)
		return nil, err
	}
	exceptionArr, err := s.core.HostOperation().RemoveFromModule(params, inputData)
	if err != nil {
		blog.ErrorJSON("RemoveFromModule error. err:%s, input:%s, exception:%s, rid:%s", err.Error(), data, exceptionArr, params.ReqID)
		return exceptionArr, err
	}
	return nil, nil
}

func (s *coreService) GetHostModuleRelation(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	inputData := &metadata.HostModuleRelationRequest{}
	if err := data.MarshalJSONInto(inputData); nil != err {
		blog.Errorf("GetHostModuleRelation MarshalJSONInto error, err:%s, input:%v, rid:%s", err.Error(), data, params.ReqID)
		return nil, err
	}
	relationArr, err := s.core.HostOperation().GetHostModuleRelation(params, inputData)
	if err != nil {
		blog.ErrorJSON("GetHostModuleRelation  error. err:%s, rid:%s", err.Error(), params.ReqID)
		return nil, err
	}
	return relationArr, nil
}

func (s *coreService) DeleteHostFromSystem(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	inputData := &metadata.DeleteHostRequest{}
	if err := data.MarshalJSONInto(inputData); nil != err {
		blog.Errorf("DeleteHost MarshalJSONInto error, err:%s, input:%v, rid:%s", err.Error(), data, params.ReqID)
		return nil, err
	}
	exceptionArr, err := s.core.HostOperation().DeleteFromSystem(params, inputData)
	if err != nil {
		blog.ErrorJSON("DeleteHost  error. err:%s, exception:%s, rid:%s", err.Error(), exceptionArr, params.ReqID)
		return exceptionArr, err
	}
	return nil, nil
}

func (s *coreService) HostIdentifier(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	inputData := &metadata.SearchHostIdentifierParam{}
	if err := data.MarshalJSONInto(inputData); nil != err {
		blog.Errorf("Identifier MarshalJSONInto error, err:%s, input:%#v, rid:%s", err.Error(), data, params.ReqID)
		return nil, err
	}
	hostIdentifierArr, err := s.core.HostOperation().Identifier(params, inputData)

	if err != nil {
		blog.InfoJSON("Identifier host identifier handle error. err:%s, input:%s, rid:%s", err.Error(), inputData, params.ReqID)
		return nil, err
	}

	return metadata.SearchHostIdentifierData{Info: hostIdentifierArr, Count: len(hostIdentifierArr)}, nil
}

// TransferHostModuleDep is a TransferHostModule dependence
func (s *coreService) TransferHostModuleDep(ctx core.ContextParams, input *metadata.HostsModuleRelation) ([]metadata.ExceptionResult, error) {
	exceptionArr, err := s.core.HostOperation().TransferToNormalModule(ctx, input)
	if err != nil {
		blog.ErrorJSON("TransferHostModule  error. err:%s, exception:%s, rid:%s", err.Error(), exceptionArr, ctx.ReqID)
		return exceptionArr, err
	}
	return nil, nil
}

func (s *coreService) GetHostByID(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	hostID, err := strconv.Atoi(pathParams("bk_host_id"))
	if err != nil {
		blog.Errorf("GetHostByID failed, get host by id, but got invalid host id, hostID: %s, err: %+v, rid: %s", hostID, err, params.ReqID)
		return nil, params.Error.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKHostIDField)
	}

	result := make(map[string]interface{}, 0)
	condition := common.KvMap{common.BKHostIDField: hostID}
	condition = util.SetModOwner(condition, params.SupplierAccount)
	err = s.db.Table(common.BKTableNameBaseHost).Find(condition).One(params.Context, &result)
	// TODO: return error for not found and deal error with all callers
	if err != nil && !s.db.IsNotFoundError(err) {
		blog.Errorf("GetHostByID failed, get host by id[%d] failed, err: %+v, rid: %s", hostID, err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommDBSelectFailed)
	}

	return result, nil
}

func (s *coreService) GetHosts(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	var dat metadata.ObjQueryInput
	if err := data.MarshalJSONInto(&dat); err != nil {
		blog.Errorf("GetHosts failed, get hosts failed with decode body err: %+v, rid: %s", err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommJSONUnmarshalFailed)
	}

	condition := util.ConvParamsTime(dat.Condition)
	condition = util.SetModOwner(condition, params.SupplierAccount)
	fieldArr := util.SplitStrField(dat.Fields, ",")
	result, err := s.getObjectByCondition(params, common.BKInnerObjIDHost, fieldArr, condition, dat.Sort, dat.Start, dat.Limit)
	if err != nil {
		blog.Errorf("get object failed type:%s,input:%v error:%v, rid: %s", common.BKInnerObjIDHost, dat, err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrHostSelectInst)
	}

	count, err := s.db.Table(common.BKTableNameBaseHost).Find(condition).Count(params.Context)
	if err != nil {
		blog.Errorf("get object failed type:%s ,input: %v error: %v, rid: %s", common.BKInnerObjIDHost, dat, err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrHostSelectInst)
	}

	return metadata.HostInfo{
		Count: int(count),
		Info:  result,
	}, nil
}

func (s *coreService) getObjectByCondition(params core.ContextParams, objType string, fields []string, condition interface{}, sort string, skip, limit int) ([]mapstr.MapStr, error) {
	results := make([]mapstr.MapStr, 0)
	tName := common.GetInstTableName(objType)

	dbInst := s.db.Table(tName).Find(condition).Sort(sort).Start(uint64(skip)).Limit(uint64(limit))
	if 0 < len(fields) {
		dbInst.Fields(fields...)
	}
	if err := dbInst.All(params.Context, &results); err != nil {
		blog.Errorf("failed to query the inst , error info %s, rid: %s", err.Error(), params.ReqID)
		return nil, err
	}

	// translate language for default name
	if m, ok := defaultNameLanguagePkg[objType]; nil != params.Lang && ok {
		for index, info := range results {
			l := m[fmt.Sprint(info["default"])]
			if len(l) >= 3 {
				results[index][l[1]] = util.FirstNotEmptyString(params.Lang.Language(l[0]), fmt.Sprint(info[l[1]]), fmt.Sprint(info[l[2]]))
			}
		}
	}

	return results, nil
}

func (s *coreService) GetHostSnap(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	hostID := pathParams(common.BKHostIDField)
	key := common.RedisSnapKeyPrefix + hostID
	result, err := s.cache.Get(key).Result()
	if nil != err && err != redis.Nil {
		blog.Errorf("get host snapshot failed, hostID: %v, err: %v, rid: %s", hostID, err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrHostGetSnapshot)
	}

	return metadata.HostSnap{
		Data: result,
	}, nil
}
