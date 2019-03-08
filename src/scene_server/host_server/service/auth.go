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
	"errors"
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful"

	auth_meta "configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
)

// verifyBusinessPermission will write response directly if authorized forbbiden
func (s *Service) verifyBusinessPermission(req *restful.Request, resp *restful.Response, businessID int64, action auth_meta.Action) (shouldContinue bool) {
	rHeader := req.Request.Header
	srvData := s.newSrvComm(rHeader)

	// check authorization by call interface
	decision, err := s.Authorizer.CanDoBusinessAction(req, businessID, action)
	if decision.Authorized == false {
		blog.Errorf("check business authorization failed, reason: %v, err: %v", decision.Reason, err)
		resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommParamsInvalid)})
		return
	}
	return true
}

// verifyHostPermission will write response directly if authorized forbbiden
func (s *Service) verifyHostPermission(req *restful.Request, resp *restful.Response, hostIDArr *[]int64, action auth_meta.Action) (shouldContinue bool) {
	rHeader := req.Request.Header
	srvData := s.newSrvComm(rHeader)
	shouldContinue = false

	// check authorization
	// step1. get app id by host id
	businessID, err := s.getHostOwenedApplicationID(rHeader, hostIDArr)
	if err != nil {
		blog.Errorf("check host authorization failed, get businessID by hostID failed, hosts:%+v, err: %v", hostIDArr, err)
		resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommParamsInvalid)})
		return
	}

	// step2. check authorization by call interface
	decision, err := s.Authorizer.CanDoHostAction(req, businessID, hostIDArr, action)
	if decision.Authorized == false {
		blog.Errorf("check host authorization failed, reason: %v, err: %v", decision.Reason, err)
		resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommParamsInvalid)})
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

func (s *Service) registerHostToCurrentBusiness(req *restful.Request, hostIDArr *[]int64) error {
	rHeader := req.Request.Header

	// get app id by host id
	businessID, err := s.getHostOwenedApplicationID(rHeader, hostIDArr)
	if err != nil {
		blog.Errorf("get businessID by hostID failed, hosts:%+v, err: %v", hostIDArr, err)
		return err
	}
	err = s.Authorizer.RegisterHosts(req, businessID, hostIDArr)
	return err
}

func (s *Service) deregisterHostFromCurrentBusiness(req *restful.Request, hostIDArr *[]int64) error {
	rHeader := req.Request.Header

	// get app id by host id
	businessID, err := s.getHostOwenedApplicationID(rHeader, hostIDArr)
	if err != nil {
		blog.Errorf("get businessID by hostID failed, hosts:%+v, err: %v", hostIDArr, err)
		return err
	}
	err = s.Authorizer.DeregisterHosts(req, businessID, hostIDArr)
	return err
}

// verifyCreatePlatPermission will write response directly if authorized forbbiden
func (s *Service) verifyCreatePlatPermission(req *restful.Request, resp *restful.Response) (shouldContinue bool) {
	shouldContinue = true
	return shouldContinue
}

// verifyPlatPermission will write response directly if authorized forbbiden
func (s *Service) verifyPlatPermission(req *restful.Request, resp *restful.Response, platIDArr *[]int64, action auth_meta.Action) (shouldContinue bool) {
	shouldContinue = true
	// TODO finish this method
	return shouldContinue
}

func (s *Service) registerPlat(req *restful.Request, platID int64, businessID int64) error {
	// TODO implement me

	// get app id by host id
	// err = s.Authorizer.RegisterHosts(req, businessID, hostIDArr)
	return nil
}

func (s *Service) deregisterPlat(req *restful.Request, platID int64, businessID int64) error {
	// TODO implement me

	// get app id by host id
	// err = s.Authorizer.DeregisterHosts(req, businessID, hostIDArr)
	return nil
}
