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

	"github.com/emicklei/go-restful"

	authmeta "configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
    "configcenter/src/common/condition"
    "configcenter/src/common/metadata"
    "configcenter/src/common/util"
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
func (s *Service) verifyModulePermission(requestHeader *http.Header, resp *restful.Response, moduleID int64, action authmeta.Action) (shouldContinue bool) {
    srvData := s.newSrvComm(*requestHeader)

    // check authorization by call interface
    decision, err := s.Authorizer.CanDoModuleAction(requestHeader, moduleID, action)
    if decision.Authorized == false {
        blog.Errorf("check module:%d action:%s authorization failed, reason: %v, err: %v", moduleID, action, decision.Reason, err)
        resp.WriteError(http.StatusForbidden, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
        return
    }
    return true
}

// verifyHostPermission will write response directly if authorized forbidden
func (s *Service) verifyHostPermission(requestHeader *http.Header, resp *restful.Response, hostIDArr *[]int64, action authmeta.Action) (shouldContinue bool) {
	srvData := s.newSrvComm(*requestHeader)
	shouldContinue = false

	// check authorization
	var businessID int64
	var err error
	if len(*hostIDArr) > 0 {
        // step1. get app id by host id
        businessID, err = s.getHostOwenedApplicationID(*requestHeader, hostIDArr)
        if err != nil {
            blog.Errorf("check host authorization failed, get businessID by hostID failed, hosts:%+v, err: %v", hostIDArr, err)
            resp.WriteError(http.StatusForbidden, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
            return
        }
	}

	// step2. check authorization by call interface
	decision, err := s.Authorizer.CanDoHostAction(requestHeader, businessID, hostIDArr, action)
	if decision.Authorized == false {
		blog.Errorf("check host authorization failed, reason: %v, err: %v", decision.Reason, err)
		resp.WriteError(http.StatusForbidden, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommParamsInvalid)})
		return
	}
	return true
}

// get business id by hostID
func (s *Service) getHostOwenedApplicationID(rHeader http.Header, hostIDArr *[]int64) (int64, error) {
	srvData := s.newSrvComm(rHeader)
	cond := map[string][]int64{common.BKHostIDField: *hostIDArr}
	details, err := srvData.lgc.GetConfigByCond(srvData.ctx, cond)
	if err != nil {
		blog.Errorf("get app id by host id failed, err: %v,hosts:[%+v],rid:%s", err, hostIDArr, srvData.rid)
		return -1, err
	}
	if len(details) == 0 {
		blog.Errorf("get app id by host id failed, get empty result, hosts:[%+v],rid:%s", hostIDArr, srvData.rid)
		err := fmt.Errorf("get app id by host id failed, get empty result, hosts: %+v", *hostIDArr)
		return -1, err
	}
	businessID := details[0][common.BKAppIDField]
	for _, detail := range details {
		bizID := detail[common.BKAppIDField]
		if bizID != businessID {
			return -1, errors.New("hosts don't belongs to same business")
		}
	}
	return businessID, nil
}

// get resource layers id by hostID
func (s *Service) getHostLayers(rHeader http.Header, hostIDArr *[]int64) (bkBizID int64, batchLayers [][]authmeta.Item, err error) {
    batchLayers = make([][]authmeta.Item, 0)
    
    cond := condition.CreateCondition()
    cond.Field(common.BKHostIDField).In(hostIDArr)
    query := &metadata.QueryCondition{
        Fields: []string{common.BKAppIDField, common.BKModuleIDField, common.BKSetIDField, common.BKHostIDField},
        Condition: cond.ToMapStr(),
        Limit:     metadata.SearchLimit{Limit: common.BKNoLimit},
    }
    result, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(context.Background(), rHeader, common.BKTableNameModuleHostConfig, query)
    if err != nil {
        err = fmt.Errorf("get host:%+v layer failed, err: %+v", hostIDArr, err)
        return
    }
    if len(result.Data.Info) == 0 {
        err = fmt.Errorf("get host:%+v layer failed, get host module config by host id not found, maybe hostID invalid", hostIDArr)
        return
    }
    bkBizID, err = util.GetInt64ByInterface(result.Data.Info[0][common.BKAppIDField])
    if err != nil {
        err = fmt.Errorf("get host:%+v layer failed, err: %+v", hostIDArr, err)
        return
    }
    
    bizTopoTreeRoot, err := s.Engine.CoreAPI.CoreService().Mainline().SearchMainlineInstanceTopo(context.Background(), rHeader, bkBizID, true)
    if err != nil {
        err = fmt.Errorf("get host:%+v layer failed, err: %+v", hostIDArr, err)
        return
    }
    
    for _, item := range result.Data.Info {
        moduleID, err := util.GetInt64ByInterface(item[common.BKModuleIDField])
        if err != nil {
            err = fmt.Errorf("get host:%+v layer failed, err: %+v", hostIDArr, err)
        }
        path := bizTopoTreeRoot.TraversalFindModule(moduleID)
        layers := make([]authmeta.Item, 0)
        for i := len(path) -1; i >= 0; i-- {
            node := path[i]
            item := authmeta.Item{
                Name: node.Name(),
                InstanceID: node.InstanceID,
            }
            layers = append(layers, item)
        }
        batchLayers = append(batchLayers, layers)
    }
    
    return
}

func (s *Service) registerHostToCurrentBusiness(requestHeader *http.Header, hostIDArr *[]int64) error {
	// get app id by host id
	businessID, err := s.getHostOwenedApplicationID(*requestHeader, hostIDArr)
	if err != nil {
		blog.Errorf("get businessID by hostID failed, hosts:%+v, err: %v", hostIDArr, err)
		return err
	}
	err = s.Authorizer.RegisterHosts(requestHeader, businessID, hostIDArr)
	return err
}

func (s *Service) deregisterHostFromCurrentBusiness(requestHeader *http.Header, hostIDArr *[]int64) error {
	// get app id by host id
	businessID, err := s.getHostOwenedApplicationID(*requestHeader, hostIDArr)
	if err != nil {
		blog.Errorf("get businessID by hostID failed, hosts:%+v, err: %v", hostIDArr, err)
		return err
	}
	err = s.Authorizer.DeregisterHosts(requestHeader, businessID, hostIDArr)
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
	return nil
}

func (s *Service) deregisterPlat(requestHeader *http.Header, platID int64, businessID int64) error {
	// TODO implement me

	// get app id by host id
	// err = s.Authorizer.DeregisterHosts(req, businessID, hostIDArr)
	return nil
}
