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
	"configcenter/src/storage"
)

func (s *Service) clear(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	if version.CCRunMode == version.CCRunModeProduct {
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrCommMigrateFailed)})
		return
	}

	err := clearDatabase(s.db)
	if nil != err {
		blog.Errorf("clear error: %v", err)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrCommMigrateFailed)})
		return
	}

	resp.WriteEntity(metadata.NewSuccessResp(nil))
}

var (
	tablenames = []string{
		"cc_ApplicationBase",
		"cc_History",
		"cc_HostBase",
		"cc_HostFavourite",
		"cc_ModuleBase",
		"cc_ModuleHostConfig",
		"cc_ObjectBase",
		"cc_OperationLog",
		"cc_PlatBase",
		"cc_Privilege",
		"cc_Proc2Module",
		"cc_Process",
		"cc_SetBase",
		"cc_Subscription",
		"cc_UserAPI",
		"cc_UserCustom",
		"cc_UserGroup",
		"cc_UserGroupPrivilege",
		"cc_idgenerator",
		"cc_ObjAsst",
		"cc_ObjAttDes",
		"cc_ObjClassification",
		"cc_ObjDes",
		"cc_PropertyGroup",
		"cc_InstAsst",
		"cc_System",
	}
)

func clearDatabase(instData storage.DI) error {

	// clear mongodb
	for _, tablename := range tablenames {
		instData.DropTable(tablename)
	}

	// clear redis

	return nil
}
