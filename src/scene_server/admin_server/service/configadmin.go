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
	"net/http"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful/v3"
)

// SearchConfigAdmin search the config
func (s *Service) SearchConfigAdmin(req *restful.Request, resp *restful.Response) {
	rHeader := req.Request.Header
	rid := util.GetHTTPCCRequestID(rHeader)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(rHeader))

	cond := map[string]interface{}{
		"_id": common.ConfigAdminID,
	}

	ret := struct {
		Config string `json:"config"`
	}{}
	err := s.db.Table(common.BKTableNameSystem).Find(cond).Fields(common.ConfigAdminValueField).One(s.ctx, &ret)
	if err != nil {
		blog.Errorf("SearchConfigAdmin failed, err: %+v, rid: %s", err, rid)
		result := &metadata.RespError{
			Msg: defErr.Error(common.CCErrCommDBSelectFailed),
		}
		_ = resp.WriteError(http.StatusOK, result)
		return
	}
	conf := metadata.ConfigAdmin{}
	if err := json.Unmarshal([]byte(ret.Config), &conf); err != nil {
		blog.Errorf("SearchConfigAdmin failed, Unmarshal err: %v, config:%+v,rid:%s", err, ret.Config, rid)
		_ = resp.WriteError(http.StatusOK, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	_ = resp.WriteEntity(metadata.NewSuccessResp(conf))
}

// UpdateConfigAdmin udpate the config
func (s *Service) UpdateConfigAdmin(req *restful.Request, resp *restful.Response) {
	rHeader := req.Request.Header
	rid := util.GetHTTPCCRequestID(rHeader)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(rHeader))

	config := new(metadata.ConfigAdmin)
	if err := json.NewDecoder(req.Request.Body).Decode(config); err != nil {
		blog.Errorf("UpdateConfigAdmin failed, decode body err: %v, body:%+v,rid:%s", err, req.Request.Body, rid)
		_ = resp.WriteError(http.StatusOK, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if err := config.Validate(); err != nil {
		blog.Errorf("UpdateConfigAdmin failed, Validate err: %v, input:%+v,rid:%s", err, config, rid)
		_ = resp.WriteError(http.StatusOK, &metadata.RespError{Msg: err})
		return
	}

	bytes, err := json.Marshal(config)
	if err != nil {
		blog.Errorf("UpdateConfigAdmin failed, Marshal err: %v, input:%+v,rid:%s", err, config, rid)
		_ = resp.WriteError(http.StatusOK, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONMarshalFailed)})
		return
	}

	cond := map[string]interface{}{
		"_id": common.ConfigAdminID,
	}
	data := map[string]interface{}{
		common.ConfigAdminValueField: string(bytes),
		common.LastTimeField:         time.Now(),
	}

	err = s.db.Table(common.BKTableNameSystem).Update(s.ctx, cond, data)
	if err != nil {
		blog.Errorf("UpdateConfigAdmin failed, update err: %+v, rid: %s", err, rid)
		result := &metadata.RespError{
			Msg: defErr.Error(common.CCErrCommDBUpdateFailed),
		}
		_ = resp.WriteError(http.StatusOK, result)
		return
	}
	_ = resp.WriteEntity(metadata.NewSuccessResp("update config admin success"))
}

// UpdatePlatformSettingConfig update platform_setting.
func (s *Service) UpdatePlatformSettingConfig(req *restful.Request, resp *restful.Response) {
	rHeader := req.Request.Header
	rid := util.GetHTTPCCRequestID(rHeader)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(rHeader))

	config := new(metadata.PlatformSettingConfig)
	if err := json.NewDecoder(req.Request.Body).Decode(config); err != nil {
		blog.Errorf("decode param failed, err: %v, body: %v, rid: %s", err, req.Request.Body, rid)
		rErr := resp.WriteError(http.StatusOK, &metadata.RespError{
			Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed),
		})
		if rErr != nil {
			blog.Errorf("response request url: %s failed, err: %v, rid: %s", req.Request.RequestURI, rErr, rid)
			return
		}
		return
	}

	if err := config.Validate(); err != nil {
		blog.Errorf("validate param failed, err: %v, input: %v, rid: %s", err, config, rid)
		rErr := resp.WriteError(http.StatusOK, &metadata.RespError{Msg: err})
		if rErr != nil {
			blog.Errorf("response request url: %s failed, err: %v, rid: %s", req.Request.RequestURI, rErr, rid)
			return
		}
		return
	}

	err := s.updatePlatformSetting(config)
	if err != nil {
		blog.Errorf("update config admin failed, err: %v, rid: %s", err, rid)
		result := &metadata.RespError{
			Msg: defErr.Error(common.CCErrCommDBUpdateFailed),
		}
		rErr := resp.WriteError(http.StatusOK, result)
		if rErr != nil {
			blog.Errorf("response request url: %s failed, err: %v, rid: %s", req.Request.RequestURI, rErr, rid)
			return
		}
		return
	}

	err = resp.WriteEntity(metadata.NewSuccessResp("udpate config admin success"))
	if err != nil {
		blog.Errorf("response request url: %s failed, err: %v, rid: %s", req.Request.RequestURI, err, rid)
		return
	}
}

// updatePlatformSetting update current configuration to database.
func (s *Service) updatePlatformSetting(config *metadata.PlatformSettingConfig) error {

	bytes, err := json.Marshal(config)
	if err != nil {
		return err
	}

	updateCond := map[string]interface{}{
		"_id": common.ConfigAdminID,
	}

	data := map[string]interface{}{
		common.ConfigAdminValueField: string(bytes),
		common.LastTimeField:         time.Now(),
	}

	err = s.db.Table(common.BKTableNameSystem).Update(s.ctx, updateCond, data)
	if err != nil {
		return err
	}

	return nil
}

// searchCurrentConfig get the current configuration in the database.
func (s *Service) searchCurrentConfig(rid string) (*metadata.PlatformSettingConfig, error) {

	cond := map[string]interface{}{"_id": common.ConfigAdminID}
	ret := make(map[string]interface{})

	err := s.db.Table(common.BKTableNameSystem).Find(cond).Fields(common.ConfigAdminValueField).One(s.ctx, &ret)
	if err != nil {
		blog.Errorf("search platform db config failed, err: %v, rid: %s", err, rid)
		return nil, err
	}
	if ret[common.ConfigAdminValueField] == nil {
		blog.Errorf("get config failed, rid: %s", rid)
		return nil, err
	}

	if _, ok := ret[common.ConfigAdminValueField].(string); !ok {
		blog.Errorf("db config type is error,rid: %s", rid)
		return nil, err
	}

	conf := new(metadata.PlatformSettingConfig)
	if err := json.Unmarshal([]byte(ret[common.ConfigAdminValueField].(string)), conf); err != nil {
		blog.Errorf("search platform config fail, unmarshal err: %v, config: %+v,rid: %s", err, ret, rid)
		return nil, err
	}
	return conf, nil
}

// searchInitConfig get init config.
func (s *Service) searchInitConfig(rid string) (*metadata.PlatformSettingConfig, error) {
	conf := new(metadata.PlatformSettingConfig)

	if err := json.Unmarshal([]byte(metadata.InitAdminConfig), conf); err != nil {
		blog.Errorf("search initial config unmarshal fail, err: %v, rid: %s", err, rid)
		return nil, err
	}

	if err := conf.EncodeWithBase64(); err != nil {
		blog.Errorf("initial config  encode bases64 fail,err: %v, rid: %s", err, rid)
		return nil, err
	}
	return conf, nil
}

// SearchPlatformSettingConfig search the platform config.typeId:current db's config ,typeId:initial initial config.
func (s *Service) SearchPlatformSettingConfig(req *restful.Request, resp *restful.Response) {
	rHeader := req.Request.Header
	rid := util.GetHTTPCCRequestID(rHeader)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(rHeader))
	typeId := req.PathParameter("type")

	conf := new(metadata.PlatformSettingConfig)
	var err error
	switch typeId {

	case "current":
		conf, err = s.searchCurrentConfig(rid)
		if err != nil {
			rErr := resp.WriteError(http.StatusOK, &metadata.RespError{
				Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
			if rErr != nil {
				blog.Errorf("response url: %s failed, err: %v, rid: %s", req.Request.RequestURI, rErr, rid)
				return
			}
			return
		}
	case "initial":
		conf, err = s.searchInitConfig(rid)
		if err != nil {
			rErr := resp.WriteError(http.StatusOK, &metadata.RespError{
				Msg: defErr.Error(common.CCErrCommParamsInvalid),
			})
			if rErr != nil {
				blog.Errorf("response url: %s failed, err: %v, rid: %s", req.Request.RequestURI, rErr, rid)
				return
			}
			return
		}

	default:
		rErr := resp.WriteError(http.StatusOK, &metadata.RespError{
			Msg: defErr.CCErrorf(common.CCErrCommParamsInvalid, "type"),
		})

		if rErr != nil {
			blog.Errorf("response url: %s failed, err: %v, rid: %s", req.Request.RequestURI, rErr, rid)
			return
		}
		return
	}

	err = resp.WriteEntity(metadata.NewSuccessResp(conf))
	if err != nil {
		blog.Errorf("response url: %s failed, err: %v, rid: %s", req.Request.RequestURI, err, rid)
		return
	}
}
