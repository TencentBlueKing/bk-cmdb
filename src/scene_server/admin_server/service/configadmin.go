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
	"time"

	idgen "configcenter/pkg/id-gen"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal/mongo/local"

	"github.com/emicklei/go-restful/v3"
)

// SearchConfigAdmin search the config
func (s *Service) SearchConfigAdmin(req *restful.Request, resp *restful.Response) {
	rHeader := req.Request.Header
	rid := httpheader.GetRid(rHeader)
	defErr := s.CCErr.CreateDefaultCCErrorIf(httpheader.GetLanguage(rHeader))
	kit := rest.NewKitFromHeader(rHeader, s.CCErr)

	cond := map[string]interface{}{
		"_id": common.ConfigAdminID,
	}

	ret := struct {
		Config string `json:"config"`
	}{}
	err := s.db.Shard(kit.SysShardOpts()).Table(common.BKTableNameSystem).Find(cond).
		Fields(common.ConfigAdminValueField).One(s.ctx, &ret)
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
	rid := httpheader.GetRid(rHeader)
	defErr := s.CCErr.CreateDefaultCCErrorIf(httpheader.GetLanguage(rHeader))
	kit := rest.NewKitFromHeader(rHeader, s.CCErr)

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

	err = s.db.Shard(kit.SysShardOpts()).Table(common.BKTableNameSystem).Update(s.ctx, cond, data)
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
	rid := httpheader.GetRid(rHeader)
	defErr := s.CCErr.CreateDefaultCCErrorIf(httpheader.GetLanguage(rHeader))
	kit := rest.NewKitFromHeader(rHeader, s.CCErr)

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

	if err := s.validateIDGenConf(kit, &config.IDGenerator); err != nil {
		_ = resp.WriteError(http.StatusOK, &metadata.RespError{Msg: err})
		return
	}

	preConf, err := s.searchCurrentConfig(kit)
	if err != nil {
		_ = resp.WriteError(http.StatusOK, &metadata.RespError{Msg: err})
		return
	}

	err = s.updatePlatformSetting(kit, config)
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

	if err = s.savePlatformSettingUpdateAudit(kit, preConf, config); err != nil {
		_ = resp.WriteError(http.StatusOK, &metadata.RespError{Msg: err})
		return
	}

	err = resp.WriteEntity(metadata.NewSuccessResp("udpate config admin success"))
	if err != nil {
		blog.Errorf("response request url: %s failed, err: %v, rid: %s", req.Request.RequestURI, err, rid)
		return
	}
}

// validateIDGenConf validate id generator config
func (s *Service) validateIDGenConf(kit *rest.Kit, conf *metadata.IDGeneratorConf) error {
	if len(conf.InitID) == 0 {
		return nil
	}

	// check if init id types are valid, and get current sequence ids by sequence names
	seqNames := make([]string, 0)
	for typ := range conf.InitID {
		seqName, exists := idgen.GetIDGenSequenceName(typ)
		if !exists {
			blog.Errorf("id generator config type %s is invalid, rid: %s", kit.Rid)
			return fmt.Errorf("id generator type %s is invalid", typ)
		}
		seqNames = append(seqNames, seqName)
	}

	idGenCond := map[string]interface{}{
		"_id": map[string]interface{}{common.BKDBIN: seqNames},
	}

	idGens := make([]local.Idgen, 0)
	err := s.db.Shard(kit.SysShardOpts()).Table(common.BKTableNameIDgenerator).Find(idGenCond).Fields("_id",
		"SequenceID").All(s.ctx, &idGens)
	if err != nil {
		blog.Errorf("get id generator data failed, err: %v, cond: %+v, rid: %s", err, idGenCond, kit.Rid)
		return err
	}

	seqNameIDMap := make(map[string]uint64)
	for _, data := range idGens {
		seqNameIDMap[data.ID] = data.SequenceID
	}

	// check if init id config is greater than current sequence id
	for typ, id := range conf.InitID {
		seqName, _ := idgen.GetIDGenSequenceName(typ)

		if id <= seqNameIDMap[seqName] {
			blog.Errorf("id gen type %s id %d <= current id: %d, rid: %s", typ, id, seqNameIDMap[seqName], kit.Rid)
			return fmt.Errorf("id generator type %s id %d is invalid", typ, id)
		}
	}

	return nil
}

func (s *Service) savePlatformSettingUpdateAudit(kit *rest.Kit,
	preConf, curConf *metadata.PlatformSettingConfig) error {

	id, err := s.db.Shard(kit.SysShardOpts()).NextSequence(s.ctx, common.BKTableNameAuditLog)
	if err != nil {
		blog.Errorf("generate next audit log id failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	audit := metadata.AuditLog{
		ID:              int64(id),
		AuditType:       metadata.PlatformSetting,
		SupplierAccount: kit.SupplierAccount,
		User:            kit.User,
		ResourceType:    metadata.PlatformSettingRes,
		Action:          metadata.AuditUpdate,
		OperateFrom:     metadata.FromUser,
		OperationDetail: &metadata.GenericOpDetail{Data: preConf, UpdateFields: curConf},
		OperationTime:   metadata.Now(),
		AppCode:         httpheader.GetAppCode(kit.Header),
		RequestID:       kit.Rid,
	}

	if err = s.db.Shard(kit.SysShardOpts()).Table(common.BKTableNameAuditLog).Insert(s.ctx, audit); err != nil {
		blog.Errorf("save audit log failed, err: %v, audit: %+v, rid: %s", err, audit, kit.Rid)
		return err
	}

	return nil
}

// updatePlatformSetting update current configuration to database.
func (s *Service) updatePlatformSetting(kit *rest.Kit, config *metadata.PlatformSettingConfig) error {
	config.IDGenerator.CurrentID = nil

	// 校验业务是否存在
	bizCountCond := map[string]interface{}{
		common.BKAppIDField: config.Backend.SnapshotBizID,
	}
	count, err := s.db.Shard(kit.SysShardOpts()).Table(common.BKTableNameBaseApp).Find(bizCountCond).Count(s.ctx)
	if err != nil {
		blog.Errorf("update config to db failed, count biz error, err: %v, condition: %v, rid: %s", err,
			bizCountCond, kit.Rid)
		return err
	}
	if count == 0 {
		blog.Errorf("update config to db failed, can not find biz, condition: %v, rid: %s", bizCountCond, kit.Rid)
		return errors.New(common.CCErrCommParamsIsInvalid, "snapshot_biz_id")
	}

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

	err = s.db.Shard(kit.SysShardOpts()).Table(common.BKTableNameSystem).Update(s.ctx, updateCond, data)
	if err != nil {
		return err
	}

	return nil
}

// searchCurrentConfig get the current configuration in the database.
func (s *Service) searchCurrentConfig(kit *rest.Kit) (*metadata.PlatformSettingConfig, error) {

	cond := map[string]interface{}{"_id": common.ConfigAdminID}
	ret := make(map[string]interface{})

	err := s.db.Shard(kit.SysShardOpts()).Table(common.BKTableNameSystem).Find(cond).
		Fields(common.ConfigAdminValueField).One(s.ctx, &ret)
	if err != nil {
		blog.Errorf("search platform db config failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}
	if ret[common.ConfigAdminValueField] == nil {
		blog.Errorf("get config failed, rid: %s", kit.Rid)
		return nil, err
	}

	if _, ok := ret[common.ConfigAdminValueField].(string); !ok {
		blog.Errorf("db config type is error,rid: %s", kit.Rid)
		return nil, err
	}

	conf := new(metadata.PlatformSettingConfig)
	if err := json.Unmarshal([]byte(ret[common.ConfigAdminValueField].(string)), conf); err != nil {
		blog.Errorf("search platform config fail, unmarshal err: %v, config: %+v,rid: %s", err, ret, kit.Rid)
		return nil, err
	}

	conf, err = s.addIDGenInfoToConf(kit, conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

// addIDGenInfoToConf add current id generator info to current config
func (s *Service) addIDGenInfoToConf(kit *rest.Kit, conf *metadata.PlatformSettingConfig) (
	*metadata.PlatformSettingConfig, error) {

	idGenCond := map[string]interface{}{
		"_id": map[string]interface{}{common.BKDBIN: idgen.GetAllIDGenSeqNames()},
	}

	idGens := make([]local.Idgen, 0)
	err := s.db.Shard(kit.SysShardOpts()).Table(common.BKTableNameIDgenerator).Find(idGenCond).Fields("_id",
		"SequenceID").All(s.ctx, &idGens)
	if err != nil {
		blog.Errorf("list id generators failed, err: %v, cond: %+v, rid: %s", err, idGenCond, kit.Rid)
		return nil, err
	}

	seqNameIDMap := make(map[string]uint64)
	for _, idGen := range idGens {
		seqNameIDMap[idGen.ID] = idGen.SequenceID
	}

	conf.IDGenerator.CurrentID = make(map[idgen.IDGenType]uint64)
	for _, typ := range idgen.GetAllIDGenTypes() {
		seqName, _ := idgen.GetIDGenSequenceName(typ)
		conf.IDGenerator.CurrentID[typ] = seqNameIDMap[seqName]
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
	rid := httpheader.GetRid(rHeader)
	defErr := s.CCErr.CreateDefaultCCErrorIf(httpheader.GetLanguage(rHeader))
	typeId := req.PathParameter("type")
	kit := rest.NewKitFromHeader(rHeader, s.CCErr)

	conf := new(metadata.PlatformSettingConfig)
	var err error
	switch typeId {

	case "current":
		conf, err = s.searchCurrentConfig(kit)
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
