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
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/emicklei/go-restful"
	redis "gopkg.in/redis.v5"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/eventclient"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

const (
	FavouriteCollection  = "cc_HostFavourite"
	HostBaseCollection   = "cc_HostBase"
	ModuleHostCollection = "cc_ModuleHostConfig"
	UserQueryCollection  = "cc_UserAPI"
)

func (s *Service) GetHostByID(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	ownerID := util.GetOwnerID(pheader)

	pathParams := req.PathParameters()
	hostID, err := strconv.Atoi(pathParams["bk_host_id"])
	if err != nil {
		blog.Errorf("get host by id, but got invalid host id, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommParamsIsInvalid)})
		return
	}

	var result map[string]interface{}
	condition := common.KvMap{common.BKHostIDField: hostID}
	condition = util.SetModOwner(condition, ownerID)
	fields := make([]string, 0)
	err = s.Instance.GetOneByCondition(HostBaseCollection, fields, condition, &result)
	if err != nil {
		blog.Error("get host by id[%s] failed, err: %v", hostID, err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
		return
	}

	resp.WriteEntity(meta.HostInstanceResult{
		BaseResp: meta.SuccessBaseResp,
		Data:     result,
	})
}

func (s *Service) GetHosts(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	lang := s.Core.Language.CreateDefaultCCLanguageIf(util.GetLanguage(pheader))
	ownerID := util.GetOwnerID(pheader)

	objType := common.BKInnerObjIDHost
	var dat meta.ObjQueryInput
	if err := json.NewDecoder(req.Request.Body).Decode(&dat); err != nil {
		blog.Errorf("get hosts failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	condition := util.ConvParamsTime(dat.Condition)
	condition = util.SetModOwner(condition, ownerID)
	fieldArr := strings.Split(dat.Fields, ",")
	result := make([]map[string]interface{}, 0)

	err := s.Logics.GetObjectByCondition(lang, objType, fieldArr, condition, &result, dat.Sort, dat.Start, dat.Limit)
	if err != nil {
		blog.Error("get object failed type:%s,input:%v error:%v", objType, dat, err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostSelectInst)})
		return
	}

	count, err := s.Instance.GetCntByCondition(common.GetInstTableName(objType), condition)
	if err != nil {
		blog.Error("get object failed type:%s ,input: %v error: %v", objType, dat, err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostSelectInst)})
		return
	}

	resp.WriteEntity(meta.GetHostsResult{
		BaseResp: meta.SuccessBaseResp,
		Data: meta.HostInfo{
			Count: count,
			Info:  result,
		},
	})
}

func (s *Service) AddHost(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	ownerID := util.GetOwnerID(pheader)

	objType := common.BKInnerObjIDHost
	input := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("add host failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	input[common.CreateTimeField] = time.Now()
	input = util.SetModOwner(input, ownerID)
	var idName string
	id, err := s.Logics.CreateObject(objType, input, &idName)
	if err != nil {
		blog.Errorf("create object type:%s ,data: %v error: %v", objType, input, err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostCreateInst)})
		return
	}

	// record event
	originData := map[string]interface{}{}
	if err := s.Logics.GetObjectByID(objType, nil, id, originData, ""); err != nil {
		blog.Error("create event error:%v", err)
	} else {
		ec := eventclient.NewEventContextByReq(pheader, s.Cache)
		err := ec.InsertEvent(meta.EventTypeInstData, "host", meta.EventActionCreate, originData, nil)
		if err != nil {
			blog.Error("add host, but create event error:%v", err)
		}
	}

	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     map[string]int64{idName: id},
	})
}

func (s *Service) GetHostSnap(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(language)

	hostID := req.PathParameter("bk_host_id")
	key := common.RedisSnapKeyPrefix + hostID
	result, err := s.Cache.Get(key).Result()
	if nil != err && err != redis.Nil {
		blog.Error("get host snapshot failed, hostid: %v, err: %v ", hostID, err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostGetSnapshot)})
		return
	}

	resp.WriteAsJson(meta.GetHostSnapResult{
		BaseResp: meta.SuccessBaseResp,
		Data: meta.HostSnap{
			Data: result,
		},
	})
}

func (s *Service) GetHostModulesIDs(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	ownerID := util.GetOwnerID(pheader)

	params := meta.ModuleHostConfigParams{}
	if err := json.NewDecoder(req.Request.Body).Decode(&params); err != nil {
		blog.Error("get host module id failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	condition := map[string]interface{}{common.BKAppIDField: params.ApplicationID, common.BKHostIDField: params.HostID}
	condition = util.SetModOwner(condition, ownerID)
	moduleIDs, err := s.Logics.GetModuleIDsByHostID(condition) //params.HostID, params.ApplicationID)
	if nil != err {
		blog.Errorf("get host module id failed, err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrGetModule)})
		return
	}

	resp.WriteEntity(meta.GetHostModuleIDsResult{
		BaseResp: meta.SuccessBaseResp,
		Data:     moduleIDs,
	})
}

func (s *Service) AddModuleHostConfig(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	ownerID := util.GetOwnerID(pheader)

	params := meta.ModuleHostConfigParams{}
	if err := json.NewDecoder(req.Request.Body).Decode(&params); err != nil {
		blog.Errorf("add module host config failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	ec := eventclient.NewEventContextByReq(req.Request.Header, s.Cache)
	for _, moduleID := range params.ModuleID {
		_, err := s.Logics.AddSingleHostModuleRelation(ec, params.HostID, moduleID, params.ApplicationID, ownerID)
		if nil != err {
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostTransferModule)})
			return
		}
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (s *Service) DelModuleHostConfig(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	ownerID := util.GetOwnerID(pheader)

	params := meta.ModuleHostConfigParams{}
	var moduleIDs []int64
	if err := json.NewDecoder(req.Request.Body).Decode(&params); err != nil {
		blog.Errorf("del module host config failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	getModuleParams := make(map[string]interface{}, 2)
	getModuleParams[common.BKHostIDField] = params.HostID
	getModuleParams[common.BKAppIDField] = params.ApplicationID

	if 0 == len(params.ModuleID) {
		var err error
		moduleIDs, err = s.Logics.GetModuleIDsByHostID(getModuleParams) //params.HostID, params.ApplicationID)
		if nil != err {
			blog.Errorf("delete module host config failed, %v", err)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrGetOriginHostModuelRelationship)})
			return
		}
	} else {
		moduleIDs = params.ModuleID
	}

	ec := eventclient.NewEventContextByReq(req.Request.Header, s.Cache)
	for _, moduleID := range moduleIDs {
		_, err := s.Logics.DelSingleHostModuleRelation(ec, params.HostID, moduleID, params.ApplicationID, ownerID)
		if nil != err {
			blog.Errorf("delete module host config, but delete module relation failed, err: %v", err)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrDelOriginHostModuelRelationship)})
			return
		}
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (s *Service) DelDefaultModuleHostConfig(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	ownerID := util.GetOwnerID(pheader)

	params := meta.ModuleHostConfigParams{}
	if err := json.NewDecoder(req.Request.Body).Decode(&params); err != nil {
		blog.Errorf("del default module host config failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	defaultModuleIDs, err := s.Logics.GetDefaultModuleIDs(params.ApplicationID)
	if nil != err {
		blog.Errorf("defaultModuleIds appID:%d, error:%v", params.ApplicationID, err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrGetModule)})
		return
	}

	//delete default host module relation
	ec := eventclient.NewEventContextByReq(req.Request.Header, s.Cache)
	for _, defaultModuleID := range defaultModuleIDs {
		_, err := s.Logics.DelSingleHostModuleRelation(ec, params.HostID, defaultModuleID, params.ApplicationID, ownerID)
		if nil != err {
			blog.Errorf("del default module host config failed, with relation, err:%v", err)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrDelDefaultModuleHostConfig)})
			return
		}
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (s *Service) MoveHost2ResourcePool(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	ownerID := util.GetOwnerID(pheader)

	ec := eventclient.NewEventContextByReq(req.Request.Header, s.Cache)
	params := new(meta.ParamData)
	if err := json.NewDecoder(req.Request.Body).Decode(&params); err != nil {
		blog.Errorf("move host to resourece pool failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	idleModuleID, err := s.Logics.GetIDleModuleID(params.ApplicationID)
	if nil != err {
		blog.Error("get default module failed, error:%s", err.Error())
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrGetModule)})
		return
	}

	errHostIDs, faultHostIDs, err := s.Logics.CheckHostInIDle(params.ApplicationID, idleModuleID, params.HostID)
	if nil != err {
		blog.Error("get host relationship failed, err: %s", err.Error())
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrGetModule)})
		return
	}

	if 0 != len(errHostIDs) {
		blog.Errorf("move host to resource pool, but it does not belongs to free module, hostid: %v", errHostIDs)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrNotBelongToIdleModule)})
		return
	}

	var succ, addErr, delErr []int64
	for _, hostID := range params.HostID {
		//host not belong to other biz, add new host
		if !util.ContainsInt(faultHostIDs, hostID) {
			_, err = s.Logics.AddSingleHostModuleRelation(ec, hostID, params.OwnerModuleID, params.OwnerAppplicationID, ownerID)
			if nil != err {
				addErr = append(addErr, hostID)
				continue
			}
		}

		//delete origin relation
		_, err := s.Logics.DelSingleHostModuleRelation(ec, hostID, idleModuleID, params.ApplicationID, ownerID)
		if nil != err {
			delErr = append(delErr, hostID)
			continue
		}
		succ = append(succ, hostID)
	}

	if 0 != len(addErr) || 0 != len(delErr) {
		addErr = append(addErr, delErr...)
		blog.Errorf("move host to resource pool, success: %v, failed: %v", succ, addErr)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrTransfer2ResourcePool)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (s *Service) AssignHostToApp(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	ownerID := util.GetOwnerID(pheader)

	ec := eventclient.NewEventContextByReq(req.Request.Header, s.Cache)
	params := new(meta.AssignHostToAppParams)
	if err := json.NewDecoder(req.Request.Body).Decode(params); err != nil {
		blog.Errorf("assign host to app failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	getModuleParams := make(map[string]interface{})
	for _, hostID := range params.HostID {
		// delete relation in default app module
		_, err := s.Logics.DelSingleHostModuleRelation(ec, hostID, params.OwnerModuleID, params.OwnerApplicationID, ownerID)
		if nil != err {
			blog.Errorf("assign host to app, but delete host module relationship failed, err: %v")
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrTransferHostFromPool)})
			return
		}

		getModuleParams[common.BKHostIDField] = hostID
		moduleIDs, err := s.Logics.GetModuleIDsByHostID(getModuleParams)
		if nil != err {
			blog.Errorf("assign host to app, but get module failed, err: %v", err)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrGetModule)})
			return
		}

		// delete from empty module, no relation
		if 0 < len(moduleIDs) {
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrAlreadyAssign)})
			return
		}

		// add new host
		_, err = s.Logics.AddSingleHostModuleRelation(ec, hostID, params.ModuleID, params.ApplicationID, ownerID)
		if nil != err {
			blog.Errorf("assign host to app, but add single host module relation failed, err: %v", err)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrTransferHostFromPool)})
		}
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (s *Service) GetModulesHostConfig(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	ownerID := util.GetOwnerID(pheader)

	var params = make(map[string][]int)
	if err := json.NewDecoder(req.Request.Body).Decode(&params); err != nil {
		blog.Errorf("del module host config failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	query := make(map[string]interface{})
	for key, val := range params {
		conditon := make(map[string]interface{})
		conditon[common.BKDBIN] = val
		query[key] = conditon
	}

	query = util.SetModOwner(query, ownerID)
	fields := []string{common.BKAppIDField, common.BKHostIDField, common.BKSetIDField, common.BKModuleIDField}
	var result []meta.ModuleHost
	err := s.Instance.GetMutilByCondition(ModuleHostCollection, fields, query, &result, common.BKHostIDField, 0, common.BKNoLimit)
	if err != nil {
		blog.Error("get module host config failed, err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
		return
	}

	resp.WriteEntity(meta.HostConfig{
		BaseResp: meta.SuccessBaseResp,
		Data:     result,
	})
}
