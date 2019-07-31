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

	"github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
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
