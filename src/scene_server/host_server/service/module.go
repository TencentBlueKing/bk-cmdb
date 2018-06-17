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
	"encoding/json"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"github.com/emicklei/go-restful"
	"strconv"
)

func (s *Service) AddHostMutiltAppModuleRelation(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	user := util.GetUser(pheader)

	result, err := s.CoreAPI.ObjectController().Privilege().GetSystemFlag(context.Background(), common.BKDefaultOwnerID, common.HostCrossBizField, pheader)
	if err != nil || (err == nil && !result.Result) {
		blog.Errorf("add host multiple app module relation failed, err: %v, result err: %v", err, result.ErrMsg)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrHostNotAllowedToMutiBiz)})
		return
	}

	params := new(metadata.CloudHostModuleParams)
	if err := json.NewDecoder(req.Request.Body).Decode(&params); err != nil {
		blog.Errorf("add host multiple app module relation failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	module, err := s.Logics.GetModuleByModuleID(pheader, params.ApplicationID, params.ModuleID)
	if err != nil {
		blog.Errorf("add host multiple app module relation, but get module[%v] failed, err: %v", params.ModuleID, err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrTopoModuleSelectFailed)})
		return
	}

	if len(module) == 0 {
		blog.Errorf("add host multiple app module relation, but get invalid module")
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrTopoMulueIDNotfoundFailed)})
		return
	}

	defaultAppID, err := s.Logics.GetDefaultAppID(common.BKDefaultOwnerID, pheader)
	if err != nil {
		blog.Errorf("add host multiple app module relation, but get default appid failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrTopoAppSearchFailed)})
		return
	}

	var errMsg, succ []string
	var hostIDArr []int

	for index, hostInfo := range params.HostInfoArr {
		cond := NewOperation().WithHostInnerIP(hostInfo.IP).WithCloudID(strconv.Itoa(hostInfo.CloudID)).Data()
		query := &metadata.QueryInput{
			Condition: cond,
			Start:     0,
			Limit:     common.BKNoLimit,
			Sort:      common.BKHostIDField,
		}
		hResult, err := s.CoreAPI.HostController().Host().GetHosts(context.Background(), pheader, query)
		if err != nil || (err == nil && !hResult.Result) {
			blog.Errorf("add host multiple app module relation, but get hosts failed, err: %v", err)
			errMsg = append(errMsg, s.Language.Languagef("host_ip_not_exist", hostInfo.IP))
			continue
		}

		hostList := hResult.Data.Info
		if len(hostList) == 0 {
			blog.Errorf("add host multiple app module relation, but get 0 hosts ")
			errMsg = append(errMsg, s.Language.Languagef("host_ip_not_exist", hostInfo.IP))
			continue
		}

		//check if host in this module
		hostID, err := util.GetInt64ByInterface(hostList[0][common.BKHostIDField])
		if nil != err {
			blog.Error("add host multiple app module relation, but get invalid host id[%v], err:%v", hostList[0][common.BKHostIDField], err.Error())
			errMsg = append(errMsg, s.Language.Languagef("host_ip_not_exist", hostInfo.IP))
			continue
		}
		moduleHostCond := map[string][]int64{common.BKHostIDField: []int64{hostID}}
		confs, err := s.Logics.GetConfigByCond(pheader, moduleHostCond)
		if err != nil {
			blog.Error("add host multiple app module relation, but get host config failed, err:%v", err)
			errMsg = append(errMsg, s.Language.Languagef("host_ip_not_exist", hostInfo.IP))
			continue
		}

		for _, conf := range confs {
			if conf[common.BKAppIDField] == defaultAppID {
				p := metadata.ModuleHostConfigParams{
					ApplicationID: defaultAppID,
					HostID:        hostID,
				}
				hResult, err := s.CoreAPI.HostController().Module().DelDefaultModuleHostConfig(context.Background(), pheader, &p)
				if err != nil || (err == nil && !hResult.Result) {
					blog.Errorf("add host multiple app module relation, but delete default module host conf failed, err: %v", err)
					errMsg = append(errMsg, s.Language.Languagef("host_ip_not_exist", hostInfo.IP))
					continue
				}
			}

			if conf[common.BKModuleIDField] == params.ModuleID {
				blog.Errorf("add host multiple app module relation, but host already exist in module")
				errMsg = append(errMsg, s.Language.Languagef("host_str_belong_module", hostInfo.IP))
				continue
			}
		}

		//add host to this module
		// TODO: to be continued

	}

}

func (s *Service) HostModuleRelation(req *restful.Request, resp *restful.Response) {

}

func (s *Service) MoveHost2EmptyModule(req *restful.Request, resp *restful.Response) {

}

func (s *Service) MoveHost2FaultModule(req *restful.Request, resp *restful.Response) {

}

func (s *Service) MoveHostToResourcePool(req *restful.Request, resp *restful.Response) {

}

func (s *Service) AssignHostToApp(req *restful.Request, resp *restful.Response) {

}

func (s *Service) AssignHostToAppModule(req *restful.Request, resp *restful.Response) {

}
