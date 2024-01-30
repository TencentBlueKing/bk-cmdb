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
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/types"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful/v3"
)

var allConfigNames = map[string]bool{
	"redis":    true,
	"mongodb":  true,
	"common":   true,
	"extra":    true,
	"error":    true,
	"language": true,
	"all":      true,
}

var configHelpInfo = fmt.Sprintf("config_name must be one of the [redis, mongodb, common, extra, error, language, all]")

func (s *Service) refreshConfig(req *restful.Request, resp *restful.Response) {
	rHeader := req.Request.Header
	rid := util.GetHTTPCCRequestID(rHeader)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(rHeader))

	input := new(struct {
		ConfigName string `json:"config_name"`
	})
	if err := json.NewDecoder(req.Request.Body).Decode(input); err != nil {
		blog.Errorf("refreshConfig failed, decode body err: %v ,body:%+v,rid:%s", err, req.Request.Body, rid)
		resp.WriteError(http.StatusOK, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	configName := "all"
	if input.ConfigName != "" {
		if ok := allConfigNames[input.ConfigName]; !ok {
			blog.Errorf("refreshConfig failed, config_name is wrong, %s, input:%#v, rid:%s", configHelpInfo, input, rid)
			resp.WriteError(http.StatusOK,
				&metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsInvalid, configHelpInfo)})
			return
		}
		configName = input.ConfigName
	}

	var err error
	switch configName {
	case "redis", "mongodb", "common", "extra":
		filePath := filepath.Join(s.Config.Configures.Dir, configName+".yaml")
		key := types.CC_SERVCONF_BASEPATH + "/" + configName
		err = s.ConfigCenter.WriteConfigure(filePath, key)
	case "error":
		err = s.ConfigCenter.WriteErrorRes2Center(s.Config.Errors.Res)
	case "language":
		err = s.ConfigCenter.WriteLanguageRes2Center(s.Config.Language.Res)
	case "all":
		err = s.ConfigCenter.WriteAllConfs2Center(s.Config.Configures.Dir, s.Config.Errors.Res, s.Config.Language.Res)
	default:
		blog.Errorf("refreshConfig failed, config_name is wrong, %s, input:%#v, rid:%s", configHelpInfo, input, rid)
		resp.WriteError(http.StatusOK,
			&metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsInvalid, configHelpInfo)})
		return
	}

	if err != nil {
		blog.Warnf("refreshConfig failed, input:%#v, error:%v, rid:%s", input, err, rid)
		resp.WriteError(http.StatusOK, &metadata.RespError{Msg: err})
	}

	blog.Infof("refresh config success, input:%#v", input)
	resp.WriteEntity(metadata.NewSuccessResp("refresh config success"))
}
