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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/admin_server/upgrader"

	"github.com/emicklei/go-restful"
)

func (s *Service) migrate(req *restful.Request, resp *restful.Response) {
	rHeader := req.Request.Header
	rid := util.GetHTTPCCRequestID(rHeader)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(rHeader))
	ownerID := common.BKDefaultOwnerID
	updateCfg := &upgrader.Config{
		OwnerID:      ownerID,
		SupplierID:   common.BKDefaultSupplierID,
		User:         common.CCSystemOperatorUserName,
		CCApiSrvAddr: s.ccApiSrvAddr,
	}

	preVersion, finishedVersions, err := upgrader.Upgrade(s.ctx, s.db, updateCfg)
	if err != nil {
		blog.Errorf("db upgrade failed, err: %+v, rid: %s", err, rid)
		result := &metadata.RespError{
			Msg: defErr.Errorf(common.CCErrCommMigrateFailed, err.Error()),
		}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}

	currentVersion := preVersion
	if len(finishedVersions) > 0 {
		currentVersion = finishedVersions[len(finishedVersions)-1]
	}

	type MigrationResponse struct {
		metadata.BaseResp `json:",inline"`
		Data              interface{} `json:"data"`
		PreVersion        string      `json:"pre_version"`
		CurrentVersion    string      `json:"current_version"`
		FinishedVersions  []string    `json:"finished_migrations"`
	}
	result := MigrationResponse{
		BaseResp: metadata.BaseResp{
			Result:      true,
			Code:        0,
			ErrMsg:      "",
			Permissions: nil,
		},
		Data:             "migrate success",
		PreVersion:       preVersion,
		CurrentVersion:   currentVersion,
		FinishedVersions: finishedVersions,
	}
	resp.WriteEntity(result)
}
