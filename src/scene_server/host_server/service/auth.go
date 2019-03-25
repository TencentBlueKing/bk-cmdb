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
	"context"
	"errors"
	"fmt"
	"net/http"

	authmeta "configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/json"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
)

// verifyBusinessPermission will write response directly if authorized forbidden
func (s *Service) verifyBusinessPermission(requestHeader *http.Header, resp *restful.Response, businessID int64, action authmeta.Action) (shouldContinue bool) {
	srvData := s.newSrvComm(*requestHeader)

	// check authorization by call interface
	decision, err := s.Authorizer.CanDoBusinessAction(requestHeader, businessID, action)
	if decision.Authorized == false {
		blog.Errorf("check business authorization failed, reason: %v, err: %v", decision.Reason, err)
		resp.WriteError(http.StatusForbidden, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
		return
	}
	return true
}

// verifyBusinessPermission will write response directly if authorized forbidden
func (s *Service) verifyModulePermission(requestHeader *http.Header, resp *restful.Response, moduleID int64,
	action authmeta.Action) (shouldContinue bool) {

	srvData := s.newSrvComm(*requestHeader)

	businessID, err := s.getBusinessIDByModuleID(requestHeader, moduleID)
	if err != nil {
		blog.Errorf("check module:%d action:%s authorization failed, getBusinessIDByModuleID failed, reason: %v, err: %v",
			moduleID, action, err)
		resp.WriteError(http.StatusForbidden, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
		return
	}

	// check authorization by call interface
	decision, err := s.Authorizer.CanDoModuleAction(requestHeader, businessID, moduleID, action)
	if decision.Authorized == false {
		blog.Errorf("check module:%d action:%s authorization failed, reason: %v, err: %v",
			moduleID, action, decision.Reason, err)
		resp.WriteError(http.StatusForbidden, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
		return
	}
	return true
}

// verifyHostPermission will write response directly if authorized forbidden
func (s *Service) verifyHostPermission(requestHeader *http.Header, resp *restful.Response, hostIDArr *[]int64, action authmeta.Action) (shouldContinue bool) {
	srvData := s.newSrvComm(*requestHeader)
	shouldContinue = false

	var businessID int64
	if len(*hostIDArr) == 0 {
		decision, err := s.Authorizer.CanDoResourceActionWithLayers(requestHeader, authmeta.HostInstance, businessID, [][]authmeta.Item{}, action)
		if decision.Authorized == false {
			blog.Errorf("check host authorization failed, reason: %v, err: %v", decision.Reason, err)
			resp.WriteError(http.StatusForbidden, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommParamsInvalid)})
			return
		}
		return true
	}

	// step1. get app id by host id
	businessID, layers, err := s.getHostLayers(requestHeader, hostIDArr)
	if err != nil {
		blog.Errorf("get host layers by hostID failed, hostIDArr: %+v, err: %v", hostIDArr, err)
		resp.WriteError(http.StatusForbidden, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommParamsInvalid)})
		return false
	}
	// step2. check authorization by call interface
	decision, err := s.Authorizer.CanDoResourceActionWithLayers(requestHeader, authmeta.HostInstance, businessID, layers, action)
	if decision.Authorized == false {
		blog.Errorf("check host authorization failed, reason: %v, err: %v", decision.Reason, err)
		resp.WriteError(http.StatusForbidden, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommParamsInvalid)})
		return false
	}
	return true
}

func (s *Service) getInnerIPByHostIDs(rHeader http.Header, hostIDArr *[]int64) (hostIDInnerIPMap map[int64]string, err error) {
	hostIDInnerIPMap = map[int64]string{}

	cond := condition.CreateCondition()
	cond.Field(common.BKHostIDField).In(*hostIDArr)
	query := &metadata.QueryCondition{
		Fields:    []string{common.BKHostInnerIPField, common.BKHostIDField},
		Condition: cond.ToMapStr(),
		Limit:     metadata.SearchLimit{Limit: common.BKNoLimit},
	}
	hosts, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(
		context.Background(), rHeader, common.BKInnerObjIDHost, query)
	if err != nil {
		return nil, fmt.Errorf("get host:%+v layer failed, err: %+v", hostIDArr, err)
	}
	for _, host := range hosts.Data.Info {
		hostID, e := util.GetInt64ByInterface(host[common.BKHostIDField])
		if e != nil {
			return nil, fmt.Errorf("get host:%+v layer failed, err: %+v", hostIDArr, e)
		}
		innerIP := util.GetStrByInterface(host[common.BKHostInnerIPField])
		hostIDInnerIPMap[hostID] = innerIP
	}
	return hostIDInnerIPMap, nil
}

// get resource layers id by hostID(layers is a data structure for call iam)
func (s *Service) getHostLayers(requestHeader *http.Header, hostIDArr *[]int64) (
	bkBizID int64, batchLayers [][]authmeta.Item, err error) {
	batchLayers = make([][]authmeta.Item, 0)

	cond := condition.CreateCondition()
	cond.Field(common.BKHostIDField).In(*hostIDArr)
	query := &metadata.QueryCondition{
		Fields:    []string{common.BKAppIDField, common.BKModuleIDField, common.BKSetIDField, common.BKHostIDField},
		Condition: cond.ToMapStr(),
		Limit:     metadata.SearchLimit{Limit: common.BKNoLimit},
	}
	result, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(
		context.Background(), *requestHeader, common.BKTableNameModuleHostConfig, query)
	if err != nil {
		err = fmt.Errorf("get host:%+v layer failed, err: %+v", hostIDArr, err)
		return
	}
	blog.V(5).Infof("get host module config: %+v", result.Data.Info)
	if len(result.Data.Info) == 0 {
		err = fmt.Errorf("get host:%+v layer failed, get host module config by host id not found, maybe hostID invalid",
			hostIDArr)
		return
	}
	bkBizID, err = util.GetInt64ByInterface(result.Data.Info[0][common.BKAppIDField])
	if err != nil {
		err = fmt.Errorf("get host:%+v layer failed, err: %+v", hostIDArr, err)
		return
	}

	bizTopoTreeRoot, err := s.Engine.CoreAPI.CoreService().Mainline().SearchMainlineInstanceTopo(
		context.Background(), *requestHeader, bkBizID, true)
	if err != nil {
		err = fmt.Errorf("get host:%+v layer failed, err: %+v", hostIDArr, err)
		return
	}

	bizTopoTreeRootJSON, err := json.Marshal(bizTopoTreeRoot)
	if err != nil {
		err = fmt.Errorf("json encode bizTopoTreeRootJSON failed: %+v", err)
		return
	}
	blog.V(5).Infof("bizTopoTreeRoot: %s", bizTopoTreeRootJSON)

	dataInfo, err := json.Marshal(result.Data.Info)
	if err != nil {
		err = fmt.Errorf("json encode dataInfo failed: %+v", err)
		return
	}
	blog.V(5).Infof("dataInfo: %s", dataInfo)

	hostIDs := make([]int64, 0)
	for _, item := range result.Data.Info {
		hostID, e := util.GetInt64ByInterface(item[common.BKHostIDField])
		if e != nil {
			err = fmt.Errorf("extract hostID from host info failed, host: %+v, err: %+v", item, e)
			return
		}
		hostIDs = append(hostIDs, hostID)
	}
	hostIDInnerIPMap, err := s.getInnerIPByHostIDs(*requestHeader, &hostIDs)
	if err != nil {
		err = fmt.Errorf("get host:%+v InnerIP failed, err: %+v", hostIDs, err)
		return
	}

	for _, item := range result.Data.Info {
		bizID, err := util.GetInt64ByInterface(item[common.BKAppIDField])
		if err != nil {
			err = fmt.Errorf("get host:%+v layer failed, get bk_app_id field failed, err: %+v", item, err)
		}
		if bizID != bkBizID {
			continue
		}
		moduleID, err := util.GetInt64ByInterface(item[common.BKModuleIDField])
		if err != nil {
			err = fmt.Errorf("get host:%+v layer failed, err: %+v", hostIDArr, err)
		}
		path := bizTopoTreeRoot.TraversalFindModule(moduleID)
		blog.V(9).Infof("traversal find module: %d result: %+v", moduleID, path)

		hostID, err := util.GetInt64ByInterface(item[common.BKHostIDField])
		if err != nil {
			err = fmt.Errorf("get host:%+v layer failed, err: %+v", item, err)
		}

		// prepare host layer
		var innerIP string
		var exist bool
		innerIP, exist = hostIDInnerIPMap[hostID]
		if exist == false {
			innerIP = fmt.Sprintf("host:%d", hostID)
		}
		hostLayer := authmeta.Item{
			Type:       authmeta.HostInstance,
			Name:       innerIP,
			InstanceID: hostID,
		}

		// layers from topo instance tree
		layers := make([]authmeta.Item, 0)
		for i := len(path) - 1; i >= 0; i-- {
			node := path[i]
			item := authmeta.Item{
				Name:       node.Name(),
				InstanceID: node.InstanceID,
				Type:       authmeta.GetResourceTypeByObjectType(node.ObjectID),
			}
			layers = append(layers, item)
		}
		layers = append(layers, hostLayer)
		blog.V(9).Infof("layers from traversal find module:%d result: %+v", moduleID, layers)
		batchLayers = append(batchLayers, layers)
	}

	return
}

func (s *Service) registerHostToCurrentBusiness(requestHeader *http.Header, hostIDArr *[]int64) error {
	// get app id by host id
	businessID, layers, err := s.getHostLayers(requestHeader, hostIDArr)
	if err != nil {
		blog.Errorf("get businessID by hostID failed, hosts:%+v, err: %v", hostIDArr, err)
		return fmt.Errorf("get layers by host failed, err: %+v", err)
	}

	layersJSON, err := json.Marshal(layers)
	if err != nil {
		blog.Errorf("json encode layers failed, layers: %+v err: %+v", layers, err)
		return fmt.Errorf("json encode layers failed, layers: %+v err: %+v", layers, err)
	}
	blog.V(7).Infof("host layers: %s", layersJSON)

	err = s.Authorizer.RegisterResourceWithLayers(requestHeader, authmeta.HostInstance, businessID, &layers)
	return err
}

func (s *Service) deregisterHostFromCurrentBusiness(requestHeader *http.Header, hostIDArr *[]int64) error {
	// get app id by host id
	businessID, layers, err := s.getHostLayers(requestHeader, hostIDArr)
	if err != nil {
		blog.Errorf("get businessID by hostID failed, hosts:%+v, err: %v", hostIDArr, err)
		return err
	}
	err = s.Authorizer.DeregisterResourceWithLayers(requestHeader, authmeta.HostInstance, businessID, &layers)
	return err
}

// verifyCreatePlatPermission will write response directly if authorized forbidden
func (s *Service) verifyCreatePlatPermission(requestHeader *http.Header, resp *restful.Response) (shouldContinue bool) {
	shouldContinue = true
	return shouldContinue
}

// verifyPlatPermission will write response directly if authorized forbidden
func (s *Service) verifyPlatPermission(requestHeader *http.Header, resp *restful.Response, platIDArr *[]int64, action authmeta.Action) (shouldContinue bool) {
	shouldContinue = true
	// TODO finish this method
	return shouldContinue
}

func (s *Service) registerPlat(requestHeader *http.Header, platID int64, businessID int64) error {
	// TODO implement me

	// get app id by host id
	// err = s.Authorizer.RegisterHosts(req, businessID, hostIDArr)
	return errors.New("not impleted yet.")
}

func (s *Service) deregisterPlat(requestHeader *http.Header, platID int64, businessID int64) error {
	// TODO implement me

	// get app id by host id
	// err = s.Authorizer.DeregisterHosts(req, businessID, hostIDArr)
	return errors.New("not impleted yet.")
}

func (s *Service) getBusinessIDByModuleID(requestHeader *http.Header, moduleID int64) (bkBizID int64, err error) {
	// get business ID by module ID
	cond := condition.CreateCondition()
	cond.Field(common.BKModuleIDField).Eq(moduleID)
	query := &metadata.QueryCondition{
		Fields:    []string{common.BKAppIDField, common.BKModuleIDField},
		Condition: cond.ToMapStr(),
		Limit:     metadata.SearchLimit{Limit: common.BKNoLimit},
	}
	result, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(
		context.Background(), *requestHeader, common.BKTableNameBaseModule, query)
	if err != nil {
		err = fmt.Errorf("get business ID by module ID, moduleID: %d, err: %+v", moduleID, err)
		return
	}
	blog.V(5).Infof("get business ID by module ID, moduleID: %d, result: %+v", moduleID, result.Data.Info)
	if len(result.Data.Info) == 0 {
		err = fmt.Errorf("get business ID by module ID failed, moduleID: %d not found", moduleID)
		return
	}
	if len(result.Data.Info) > 1 {
		err = fmt.Errorf("get business ID by module ID failed, get multiple records found, moduleID: %d, ret: %+v",
			moduleID, result.Data.Info)
		return
	}
	businessID, exist := result.Data.Info[0][common.BKAppIDField]
	if exist == false {
		err = fmt.Errorf("get business ID by module ID failed, %s field not found, moduleID: %d, ret: %+v",
			common.BKAppIDField, moduleID, result.Data.Info)
		return
	}
	bkBizID, err = util.GetInt64ByInterface(businessID)
	if err != nil {
		err = fmt.Errorf("get business ID by module ID failed, parse %s field to int64 failed, moduleID: %d, ret: %+v",
			common.BKAppIDField, moduleID, result.Data.Info)
		return
	}
	return bkBizID, err
}

// get business id by hostID
func (s *Service) getHostOwenedApplicationID(requestHeader *http.Header, hostIDArr *[]int64) (bkBizID int64, err error) {
	// get business ID by module ID
	cond := condition.CreateCondition()
	cond.Field(common.BKHostIDField).In(hostIDArr)
	query := &metadata.QueryCondition{
		Fields:    []string{common.BKAppIDField, common.BKHostIDField},
		Condition: cond.ToMapStr(),
		Limit:     metadata.SearchLimit{Limit: common.BKNoLimit},
	}
	result, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(
		context.Background(), *requestHeader, common.BKTableNameModuleHostConfig, query)
	if err != nil {
		err = fmt.Errorf("get business ID by host ID, hostIDArr: %+v, err: %+v", hostIDArr, err)
		return
	}
	blog.V(5).Infof("get business ID by host ID, hostIDArr: %+v, result: %+v", hostIDArr, result.Data.Info)
	if len(result.Data.Info) == 0 {
		err = fmt.Errorf("get business ID by host ID failed, hostIDArr: %+v not found", hostIDArr)
		return
	}
	if len(result.Data.Info) > 1 {
		err = fmt.Errorf("get business ID by host ID failed, get multiple records found, hostIDArr: %+v, ret: %+v",
			hostIDArr, result.Data.Info)
		return
	}

	for _, hostConfig := range result.Data.Info {
		businessID, exist := hostConfig[common.BKAppIDField]
		if exist == false {
			err = fmt.Errorf("get business ID by host ID failed, %s field not found, hostIDArr: %+v, ret: %+v",
				common.BKAppIDField, hostIDArr, result.Data.Info)
			return
		}
		bkBizID, err = util.GetInt64ByInterface(businessID)
		if err != nil {
			err = fmt.Errorf("get business ID by host ID failed, parse %s field to int64 failed, hostIDArr: %+v, ret: %+v",
				common.BKAppIDField, hostIDArr, result.Data.Info)
			return
		}
		if businessID == 0 {
			businessID = bkBizID
		} else {
			if businessID != bkBizID {
				err = fmt.Errorf(
					"get business ID by host ID failed, get multiple business id by hostID, hostIDArr: %+v, ret: %+v",
					hostIDArr, result.Data.Info)
				return
			}
		}
	}

	return bkBizID, err
}
