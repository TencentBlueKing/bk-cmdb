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

	idgen "configcenter/pkg/id-gen"
	"configcenter/pkg/tenant/logics"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal/mongo/local"

	"github.com/emicklei/go-restful/v3"
)

// SearchPlatformConfig search platform id generator setting config
func (s *Service) SearchPlatformConfig(req *restful.Request, resp *restful.Response) {
	typeId := req.PathParameter("type")
	kit := rest.NewKitFromHeader(req.Request.Header, s.CCErr)

	if !logics.ValidatePlatformTenantMode(kit.TenantID, s.Config.EnableMultiTenantMode) {
		blog.Errorf("non-system tenant cannot view this configuration, rid: %s", kit.Rid)
		result := &metadata.RespError{
			Msg:     fmt.Errorf("non-system tenant cannot view this configuration"),
			ErrCode: common.CCErrAPICheckTenantInvalid,
		}
		_ = resp.WriteError(http.StatusOK, result)
		return
	}

	cond := map[string]interface{}{
		common.BKFieldDBID: common.PlatformConfig,
	}
	platformConfig := new(metadata.PlatformConfig)
	err := s.db.Shard(kit.SysShardOpts()).Table(common.BKTableNameSystem).Find(cond).Fields(typeId).One(kit.Ctx,
		platformConfig)
	if err != nil {
		blog.Errorf("get platform %s config failed, err: %v, rid: %s", typeId, err, kit.Rid)
		result := &metadata.RespError{
			Msg: kit.CCError.Error(common.CCErrCommDBSelectFailed),
		}
		_ = resp.WriteError(http.StatusOK, result)
		return
	}

	switch typeId {
	case metadata.IDGeneratorConfig:
		conf, err := s.addIDGenInfoToConf(kit, &platformConfig.IDGenerator)
		if err != nil {
			blog.Errorf("get current id config failed, err: %v, rid: %s", err, kit.Rid)
			result := &metadata.RespError{
				Msg: kit.CCError.Error(common.CCErrCommDBSelectFailed),
			}
			_ = resp.WriteError(http.StatusOK, result)
		}
		_ = resp.WriteEntity(metadata.NewSuccessResp(conf))
	default:
		blog.Errorf("invalid type, type: %s, rid: %s", typeId, kit.Rid)
		_ = resp.WriteError(http.StatusOK, &metadata.RespError{Msg: fmt.Errorf("invalid type")})
	}

}

// UpdatePlatformConfig update platform setting config
func (s *Service) UpdatePlatformConfig(req *restful.Request, resp *restful.Response) {
	typeId := req.PathParameter("type")
	kit := rest.NewKitFromHeader(req.Request.Header, s.CCErr)

	if !logics.ValidatePlatformTenantMode(kit.TenantID, s.Config.EnableMultiTenantMode) {
		blog.Errorf("non-system tenant cannot view this configuration, rid: %s", kit.Rid)
		result := &metadata.RespError{
			Msg:     fmt.Errorf("non-system tenant cannot view this configuration"),
			ErrCode: common.CCErrAPICheckTenantInvalid,
		}
		_ = resp.WriteError(http.StatusOK, result)
		return
	}

	config := new(metadata.PlatformConfig)
	if err := json.NewDecoder(req.Request.Body).Decode(config); err != nil {
		blog.Errorf("decode param failed, err: %v, body: %v, rid: %s", err, req.Request.Body, kit.Rid)
		rErr := resp.WriteError(http.StatusOK, &metadata.RespError{
			Msg: kit.CCError.Error(common.CCErrCommJSONUnmarshalFailed),
		})
		if rErr != nil {
			blog.Errorf("response request url: %s failed, err: %v, rid: %s", req.Request.RequestURI, rErr, kit.Rid)
			return
		}
		return
	}

	updateData := make(map[string]interface{})
	switch typeId {
	case metadata.IDGeneratorConfig:
		if err := config.IDGenerator.Validate(); err != nil {
			blog.Errorf("validate param failed, err: %v, input: %v, rid: %s", err, config, kit.Rid)
			rErr := resp.WriteError(http.StatusOK, &metadata.RespError{Msg: err})
			if rErr != nil {
				blog.Errorf("response request url: %s failed, err: %v, rid: %s", req.Request.RequestURI, rErr, kit.Rid)
				return
			}
			return
		}

		if err := s.validateIDGenConf(kit, &config.IDGenerator); err != nil {
			_ = resp.WriteError(http.StatusOK, &metadata.RespError{Msg: err})
			return
		}
		updateData[typeId] = config.IDGenerator
	default:
		blog.Errorf("invalid type, type: %s, rid: %s", typeId, kit.Rid)
		_ = resp.WriteError(http.StatusOK, fmt.Errorf("invalid type"))
		return
	}

	cond := map[string]interface{}{common.BKFieldDBID: common.PlatformConfig}
	preConf := make(mapstr.MapStr)
	err := s.db.Shard(kit.SysShardOpts()).Table(common.BKTableNameSystem).Find(cond).Fields(typeId).One(kit.Ctx,
		&preConf)
	if err != nil {
		_ = resp.WriteError(http.StatusOK, &metadata.RespError{Msg: err})
		return
	}

	err = s.db.Shard(kit.SysShardOpts()).Table(common.BKTableNameSystem).Update(kit.Ctx, cond, updateData)
	if err != nil {
		blog.Errorf("update platform config %s failed, err: %v, rid: %s", typeId, err, kit.Rid)
		result := &metadata.RespError{
			Msg: kit.CCError.Error(common.CCErrCommDBUpdateFailed),
		}
		rErr := resp.WriteError(http.StatusOK, result)
		if rErr != nil {
			blog.Errorf("response request url: %s failed, err: %v, rid: %s", req.Request.RequestURI, rErr, kit.Rid)
			return
		}
		return
	}

	if err = s.savePlatformSettingUpdateAudit(kit, preConf, updateData); err != nil {
		_ = resp.WriteError(http.StatusOK, &metadata.RespError{Msg: err})
		return
	}

	err = resp.WriteEntity(metadata.NewSuccessResp(fmt.Sprintf("update platform config %s success", typeId)))
	if err != nil {
		blog.Errorf("response request url: %s failed, err: %v, rid: %s", req.Request.RequestURI, err, kit.Rid)
		return
	}
}

// UpdateGlobalConfig update global general config, like topo level
func (s *Service) UpdateGlobalConfig(req *restful.Request, resp *restful.Response) {
	typeId := req.PathParameter("type")
	kit := rest.NewKitFromHeader(req.Request.Header, s.CCErr)

	config := new(metadata.GlobalSettingConfig)
	if err := json.NewDecoder(req.Request.Body).Decode(config); err != nil {
		blog.Errorf("decode param failed, err: %v, body: %v, rid: %s", err, req.Request.Body, kit.Rid)
		rErr := resp.WriteError(http.StatusOK, &metadata.RespError{
			Msg: kit.CCError.Error(common.CCErrCommJSONUnmarshalFailed),
		})
		if rErr != nil {
			blog.Errorf("response request url: %s failed, err: %v, rid: %s", req.Request.RequestURI, rErr, kit.Rid)
			return
		}
		return
	}

	preConf := make(mapstr.MapStr)
	err := s.db.Shard(kit.ShardOpts()).Table(common.BKTableNameGlobalConfig).Find(mapstr.MapStr{}).Fields(typeId).One(kit.Ctx,
		&preConf)
	if err != nil {
		_ = resp.WriteError(http.StatusOK, &metadata.RespError{Msg: err})
		return
	}

	updateData := make(map[string]interface{})
	switch typeId {
	case metadata.BackendConfig:
		if err = config.Backend.Validate(); err != nil {
			blog.Errorf("validate param failed, err: %v, input: %v, rid: %s", err, config, kit.Rid)
			rErr := resp.WriteError(http.StatusOK, &metadata.RespError{Msg: err})
			if rErr != nil {
				blog.Errorf("response request url: %s failed, err: %v, rid: %s", req.Request.RequestURI, rErr, kit.Rid)
				return
			}
			return
		}
		updateData[typeId] = config.Backend
	default:
		blog.Errorf("invalid type, type: %s, rid: %s", typeId, kit.Rid)
		_ = resp.WriteError(http.StatusOK, fmt.Errorf("invalid type"))
		return
	}

	err = s.db.Shard(kit.ShardOpts()).Table(common.BKTableNameGlobalConfig).Update(kit.Ctx, mapstr.MapStr{}, updateData)
	if err != nil {
		blog.Errorf("update general config %s failed, err: %v, rid: %s", typeId, err, kit.Rid)
		result := &metadata.RespError{
			Msg: kit.CCError.Error(common.CCErrCommDBUpdateFailed),
		}
		rErr := resp.WriteError(http.StatusOK, result)
		if rErr != nil {
			blog.Errorf("response request url: %s failed, err: %v, rid: %s", req.Request.RequestURI, rErr, kit.Rid)
			return
		}
		return
	}

	if err = s.savePlatformSettingUpdateAudit(kit, preConf, updateData); err != nil {
		_ = resp.WriteError(http.StatusOK, &metadata.RespError{Msg: err})
		return
	}

	err = resp.WriteEntity(metadata.NewSuccessResp(fmt.Sprintf("update general %s config success", typeId)))
	if err != nil {
		blog.Errorf("response request url: %s failed, err: %v, rid: %s", req.Request.RequestURI, err, kit.Rid)
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

func (s *Service) savePlatformSettingUpdateAudit(kit *rest.Kit, preConf, curConf interface{}) error {

	id, err := s.db.Shard(kit.SysShardOpts()).NextSequence(s.ctx, common.BKTableNameAuditLog)
	if err != nil {
		blog.Errorf("generate next audit log id failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	audit := metadata.AuditLog{
		ID:              int64(id),
		AuditType:       metadata.PlatformSetting,
		ResourceType:    metadata.PlatformSettingRes,
		User:            kit.User,
		Action:          metadata.AuditUpdate,
		OperateFrom:     metadata.FromUser,
		OperationDetail: &metadata.GenericOpDetail{Data: preConf, UpdateFields: curConf},
		OperationTime:   metadata.Now(),
		AppCode:         httpheader.GetAppCode(kit.Header),
		RequestID:       kit.Rid,
	}

	if err = s.db.Shard(kit.ShardOpts()).Table(common.BKTableNameAuditLog).Insert(s.ctx, audit); err != nil {
		blog.Errorf("save audit log failed, err: %v, audit: %+v, rid: %s", err, audit, kit.Rid)
		return err
	}

	return nil
}

// addIDGenInfoToConf add current id generator info to current config
func (s *Service) addIDGenInfoToConf(kit *rest.Kit, conf *metadata.IDGeneratorConf) (
	*metadata.IDGeneratorConf, error) {

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

	conf.CurrentID = make(map[idgen.IDGenType]uint64)
	for _, typ := range idgen.GetAllIDGenTypes() {
		seqName, _ := idgen.GetIDGenSequenceName(typ)
		conf.CurrentID[typ] = seqNameIDMap[seqName]
	}

	return conf, nil
}

// SearchGlobalConfig search current global config, include maxTopoLevel, set, idle_pool, etc.
func (s *Service) SearchGlobalConfig(req *restful.Request, resp *restful.Response) {
	kit := rest.NewKitFromHeader(req.Request.Header, s.CCErr)

	conf := make(mapstr.MapStr)
	err := s.db.Shard(kit.ShardOpts()).Table(common.BKTableNameGlobalConfig).Find(mapstr.MapStr{}).One(kit.Ctx, &conf)
	if err != nil {
		rErr := resp.WriteError(http.StatusOK, &metadata.RespError{
			Msg: kit.CCError.Error(common.CCErrCommJSONUnmarshalFailed)})
		if rErr != nil {
			blog.Errorf("response url: %s failed, err: %v, rid: %s", req.Request.RequestURI, rErr, kit.Rid)
			return
		}
		return
	}
	delete(conf, common.TenantID)

	err = resp.WriteEntity(metadata.NewSuccessResp(conf))
	if err != nil {
		blog.Errorf("response url: %s failed, err: %v, rid: %s", req.Request.RequestURI, err, kit.Rid)
		return
	}
}
