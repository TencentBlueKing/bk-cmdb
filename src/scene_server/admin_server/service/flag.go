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
	"net/http"
	"strings"
	"time"

	"github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

var (
	// 允许用户设置的key
	userConfigKeyMap = map[string]bool{
		"blueking_modify": true,
	}
	// 过期时间
	userConfigDefaultExpireHour = 6
)

// SetSystemConfiguration used for set variable in cc_System table
func (s *Service) SetSystemConfiguration(req *restful.Request, resp *restful.Response) {
	rHeader := req.Request.Header
	rid := util.GetHTTPCCRequestID(rHeader)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(rHeader))
	ownerID := common.BKDefaultOwnerID

	blog.Infof("set system configuration on table %s start, rid: %s", common.BKTableNameSystem, rid)
	cond := map[string]interface{}{
		common.HostCrossBizField: common.HostCrossBizValue,
	}
	data := map[string]interface{}{
		common.HostCrossBizField: common.HostCrossBizValue + ownerID,
	}

	err := s.db.Table(common.BKTableNameSystem).Update(s.ctx, cond, data)
	if nil != err {
		blog.Errorf("set system configuration on table %s failed, err: %+v, rid: %s", common.BKTableNameSystem, err, rid)
		result := &metadata.RespError{
			Msg: defErr.Error(common.CCErrCommMigrateFailed),
		}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}
	resp.WriteEntity(metadata.NewSuccessResp("modify system config success"))
}

func (s *Service) UserConfigSwitch(req *restful.Request, resp *restful.Response) {
	rid, _, defErr := s.getCommObject(req.Request.Header)

	canModify := strings.ToLower(req.PathParameter("can"))
	key := req.PathParameter("key")
	blCanModify := false

	if _, ok := userConfigKeyMap[key]; !ok {
		result := &metadata.RespError{
			Msg: defErr.Errorf(common.CCErrCommParamsIsInvalid, key),
		}
		resp.WriteError(http.StatusBadRequest, result)
		return
	}
	switch canModify {
	case "true":
		blCanModify = true
	case "false":
		blCanModify = false
	default:
		result := &metadata.RespError{
			Msg: defErr.Errorf(common.CCErrCommParamsNeedBool, "can"),
		}
		resp.WriteError(http.StatusBadRequest, result)
		return
	}
	cond := map[string]interface{}{
		"type": metadata.CCSystemUserConfigSwitch,
	}
	data := map[string]metadata.SysUserConfigItem{
		key: metadata.SysUserConfigItem{
			Flag:     blCanModify,
			ExpireAt: time.Now().Unix() + int64(userConfigDefaultExpireHour*3600),
		},
	}

	err := s.db.Table(common.BKTableNameSystem).Upsert(s.ctx, cond, data)
	if err != nil {
		blog.ErrorJSON("UserConfigSwitch set key %s value %s error. err:%s, rid:%s", key, canModify, err, rid)
		resp.WriteError(http.StatusBadGateway, defErr.Error(common.CCErrCommDBUpdateFailed))
		return
	}
	resp.WriteEntity(metadata.NewSuccessResp("modify system user config success"))

}

func (s *Service) getCommObject(header http.Header) (ownerID, rid string, defErr errors.DefaultCCErrorIf) {
	rid = util.GetHTTPCCRequestID(header)
	defErr = s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	ownerID = common.BKDefaultOwnerID
	return
}
