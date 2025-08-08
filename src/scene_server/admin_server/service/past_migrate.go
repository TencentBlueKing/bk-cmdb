/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package service

import (
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader/history"

	"github.com/emicklei/go-restful/v3"
)

// migrate old upgrader to v3.14
func (s *Service) migrate(req *restful.Request, resp *restful.Response) {
	rHeader := req.Request.Header
	rid := httpheader.GetRid(rHeader)
	defErr := s.CCErr.CreateDefaultCCErrorIf(httpheader.GetLanguage(rHeader))
	updateCfg := &history.Config{
		TenantID: "0",
		User:     common.CCSystemOperatorUserName,
	}

	preVersion, finishedVersions, err := history.Upgrade(s.ctx, s.oldMigrateDB, s.cache, s.iam, updateCfg)
	if err != nil {
		blog.Errorf("db upgrade failed, err: %v, rid: %s", err, rid)
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

// MigrationResponse TODO
type MigrationResponse struct {
	metadata.BaseResp `json:",inline"`
	Data              interface{} `json:"data"`
	PreVersion        string      `json:"pre_version"`
	CurrentVersion    string      `json:"current_version"`
	FinishedVersions  []string    `json:"finished_migrations"`
}

// MigrateSpecifyVersionRequest TODO
type MigrateSpecifyVersionRequest struct {
	CommitID  string `json:"commit_id"`
	TimeStamp int64  `json:"time_stamp"`
	Version   string `json:"version"`
}
