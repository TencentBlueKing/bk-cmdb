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
	"configcenter/src/common/version"
	"configcenter/src/storage/dal"
)

// clear drop tables common.AllTables from db
func (s *Service) clear(req *restful.Request, resp *restful.Response) {
	rHeader := req.Request.Header
	rid := util.GetHTTPCCRequestID(rHeader)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(rHeader))

	if version.CCRunMode == version.CCRunModeProduct {
		blog.Errorf("clear production data forbidden, rid: %s", rid)
		result := &metadata.RespError{
			Msg: defErr.Error(common.CCErrCommMigrateFailed),
		}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}

	err := clearDatabase(s.db)
	if nil != err {
		blog.Errorf("clear database failed, err: %+v, rid: %s", err, rid)
		result := &metadata.RespError{
			Msg: defErr.Error(common.CCErrCommMigrateFailed),
		}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}

	resp.WriteEntity(metadata.NewSuccessResp(nil))
}

func clearDatabase(db dal.RDB) error {
	// clear mongodb
	for _, tableName := range common.AllTables {
		db.DropTable(tableName)
	}

	// TODO clear redis

	return nil
}
