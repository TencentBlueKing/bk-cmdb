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
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"configcenter/pkg/tenant"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal/mongo"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/dal/mongo/sharding"

	"github.com/emicklei/go-restful/v3"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

func (s *Service) initShardingApi(api *restful.WebService) {
	api.Route(api.GET("/find/system/sharding_db_config").To(s.GetShardingDBConfig))
	api.Route(api.PUT("/update/system/sharding_db_config").To(s.UpdateShardingDBConfig))
	api.Route(api.GET("/find/system/tenant_db_relation").To(s.GetTenantDBRelation))
}

// ShardingDBConfig is the sharding db config for api
type ShardingDBConfig struct {
	MasterDB     string                   `json:"master_db"`
	ForNewTenant string                   `json:"for_new_tenant"`
	SlaveDB      map[string]SlaveDBConfig `json:"slave_db"`
}

// SlaveDBConfig is the slave db config for api
type SlaveDBConfig struct {
	Name     string        `json:"name"`
	Disabled bool          `json:"disabled"`
	Config   *mongo.Config `json:"config"`
}

// GetShardingDBConfig get sharding db config
func (s *Service) GetShardingDBConfig(req *restful.Request, resp *restful.Response) {
	kit := rest.NewKitFromHeader(req.Request.Header, s.CCErr)

	conf, err := s.getShardingDBConf(kit)
	if err != nil {
		_ = resp.WriteError(http.StatusOK, &metadata.RespError{Msg: err})
		return
	}

	result := &ShardingDBConfig{
		MasterDB:     conf.MasterDB,
		ForNewTenant: conf.ForNewTenant,
		SlaveDB:      make(map[string]SlaveDBConfig),
	}

	for uuid, mongoConf := range conf.SlaveDB {
		uri, err := s.crypto.Decrypt(mongoConf.URI)
		if err != nil {
			blog.Errorf("decrypt %s slave mongo uri failed, err: %v, rid: %s", uuid, err, kit.Rid)
			_ = resp.WriteError(http.StatusOK,
				&metadata.RespError{Msg: kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "uri")})
			return
		}

		connStr, err := connstring.Parse(uri)
		if err != nil {
			blog.Errorf("parse %s mongo config uri failed, err: %v, rid: %s", uuid, err, kit.Rid)
			_ = resp.WriteError(http.StatusOK,
				&metadata.RespError{Msg: kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "uri")})
			return
		}

		addr := strings.Join(connStr.Hosts, ",")
		port := ""
		if len(connStr.Hosts) == 1 {
			addrPort := strings.Split(connStr.Hosts[0], ":")
			if len(addrPort) == 2 {
				addr = addrPort[0]
				port = addrPort[1]
			}
		}

		result.SlaveDB[uuid] = SlaveDBConfig{
			Name:     mongoConf.Name,
			Disabled: mongoConf.Disabled,
			Config: &mongo.Config{
				Address:       addr,
				Port:          port,
				User:          connStr.Username,
				Database:      connStr.Database,
				Mechanism:     connStr.AuthMechanism,
				MaxOpenConns:  mongoConf.MaxOpenConns,
				MaxIdleConns:  mongoConf.MaxIdleConns,
				RsName:        mongoConf.RsName,
				SocketTimeout: mongoConf.SocketTimeout,
			},
		}
	}

	_ = resp.WriteEntity(metadata.NewSuccessResp(result))
}

func (s *Service) getShardingDBConf(kit *rest.Kit) (*sharding.ShardingDBConf, error) {
	cond := map[string]any{common.MongoMetaID: common.ShardingDBConfID}
	conf := new(sharding.ShardingDBConf)
	err := s.db.Shard(kit.SysShardOpts()).Table(common.BKTableNameSystem).Find(cond).One(kit.Ctx, &conf)
	if err != nil {
		blog.Errorf("get sharding db config failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	return conf, nil
}

// UpdateShardingDBReq is the update sharding db config request
type UpdateShardingDBReq struct {
	ForNewTenant  string                       `json:"for_new_tenant,omitempty"`
	CreateSlaveDB []SlaveDBConfig              `json:"create_slave_db,omitempty"`
	UpdateSlaveDB map[string]UpdateSlaveDBInfo `json:"update_slave_db,omitempty"`
}

// UpdateSlaveDBInfo is the update slave db info
type UpdateSlaveDBInfo struct {
	Name     string        `json:"name,omitempty"`
	Disabled *bool         `json:"disabled,omitempty"`
	Config   *mongo.Config `json:"config,omitempty"`
}

// UpdateShardingDBConfig update sharding db config
func (s *Service) UpdateShardingDBConfig(req *restful.Request, resp *restful.Response) {
	kit := rest.NewKitFromHeader(req.Request.Header, s.CCErr)

	conf := new(UpdateShardingDBReq)
	if err := json.NewDecoder(req.Request.Body).Decode(conf); err != nil {
		blog.Errorf("decode param failed, err: %v, body: %v, rid: %s", err, req.Request.Body, kit.Rid)
		_ = resp.WriteError(http.StatusOK,
			&metadata.RespError{Msg: kit.CCError.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	preConf, err := s.getShardingDBConf(kit)
	if err != nil {
		_ = resp.WriteError(http.StatusOK, &metadata.RespError{Msg: err})
		return
	}

	updateConf, err := s.genUpdatedShardingDBConf(kit, preConf, conf)
	if err != nil {
		_ = resp.WriteError(http.StatusOK, &metadata.RespError{Msg: err})
		return
	}

	if err = s.saveUpdateShardingDBAudit(kit, preConf, updateConf); err != nil {
		_ = resp.WriteError(http.StatusOK, &metadata.RespError{Msg: err})
		return
	}

	cond := map[string]any{common.MongoMetaID: common.ShardingDBConfID}
	updateData := map[string]any{
		"for_new_tenant": updateConf.ForNewTenant,
		"slave_db":       updateConf.SlaveDB,
	}
	err = s.db.Shard(kit.SysShardOpts()).Table(common.BKTableNameSystem).Update(s.ctx, cond, updateData)
	if err != nil {
		blog.Errorf("update sharding db config failed, err: %v, rid: %s", err, kit.Rid)
		_ = resp.WriteError(http.StatusOK, &metadata.RespError{Msg: kit.CCError.Error(common.CCErrCommDBUpdateFailed)})
		return
	}

	_ = resp.WriteEntity(metadata.NewSuccessResp(nil))
}

func (s *Service) genUpdatedShardingDBConf(kit *rest.Kit, dbConf *sharding.ShardingDBConf, conf *UpdateShardingDBReq) (
	*sharding.ShardingDBConf, error) {

	nameUUIDMap := make(map[string]string)
	for uuid, dbSlaveConf := range dbConf.SlaveDB {
		nameUUIDMap[dbSlaveConf.Name] = uuid
	}

	// update slave db config
	for uuid, config := range conf.UpdateSlaveDB {
		dbSlaveConf, exists := dbConf.SlaveDB[uuid]
		if !exists {
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, uuid)
		}

		if config.Name != "" {
			delete(nameUUIDMap, dbSlaveConf.Name)
			_, exists := nameUUIDMap[config.Name]
			if exists {
				return nil, kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, config.Name)
			}
			nameUUIDMap[config.Name] = uuid
			dbSlaveConf.Name = config.Name
		}

		if config.Disabled != nil {
			dbSlaveConf.Disabled = *config.Disabled
		}

		if config.Config != nil {
			var err error
			dbSlaveConf, err = s.genDBSlaveConf(kit, dbSlaveConf.Name, dbSlaveConf.Disabled, config.Config)
			if err != nil {
				return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, uuid)
			}
		}
		dbConf.SlaveDB[uuid] = dbSlaveConf
	}

	// create new slave db, generate new uuid
	for _, config := range conf.CreateSlaveDB {
		_, exists := nameUUIDMap[config.Name]
		if exists {
			return nil, kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, config.Name)
		}

		dbSlaveConf, err := s.genDBSlaveConf(kit, config.Name, config.Disabled, config.Config)
		if err != nil {
			return nil, err
		}

		newUUID := uuid.NewString()
		dbConf.SlaveDB[newUUID] = dbSlaveConf
		nameUUIDMap[config.Name] = newUUID
	}

	// update new tenant db config, check if the new tenant db config exists
	if conf.ForNewTenant != "" {
		// use uuid to specify the new tenant db config for db that already exists
		_, uuidExists := dbConf.SlaveDB[conf.ForNewTenant]
		if conf.ForNewTenant == dbConf.MasterDB || uuidExists {
			dbConf.ForNewTenant = conf.ForNewTenant
			return dbConf, nil
		}

		// use name to specify the new tenant db config for new db that doesn't have uuid before creation
		uuid, nameExists := nameUUIDMap[conf.ForNewTenant]
		if nameExists {
			dbConf.ForNewTenant = uuid
			return dbConf, nil
		}

		blog.Errorf("add new tenant db %s is invalid, rid: %s", conf.ForNewTenant, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "for_new_tenant")
	}
	return dbConf, nil
}

// genDBSlaveConf generate slave db config
func (s *Service) genDBSlaveConf(kit *rest.Kit, name string, disabled bool, conf *mongo.Config) (local.MongoConf,
	error) {

	mongoConf := conf.GetMongoConf()
	if err := mongoConf.Validate(5 * time.Second); err != nil {
		blog.Errorf("validate %s mongo config failed, err: %v, rid: %s", name, err, kit.Rid)
		return mongoConf, err
	}

	mongoConf.Name = name
	mongoConf.Disabled = disabled

	// encrypt mongodb uri
	uri, err := s.crypto.Encrypt(mongoConf.URI)
	if err != nil {
		blog.Errorf("encrypt %s mongo config failed, err: %v, rid: %s", name, err, kit.Rid)
		return mongoConf, err
	}
	mongoConf.URI = uri
	return mongoConf, nil
}

func (s *Service) saveUpdateShardingDBAudit(kit *rest.Kit, preConf, curConf *sharding.ShardingDBConf) error {
	id, err := s.db.Shard(kit.SysShardOpts()).NextSequence(kit.Ctx, common.BKTableNameAuditLog)
	if err != nil {
		blog.Errorf("generate next audit log id failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	audit := metadata.AuditLog{
		ID:              int64(id),
		AuditType:       metadata.Sharding,
		TenantID:        kit.TenantID,
		User:            kit.User,
		ResourceType:    metadata.ShardingRes,
		Action:          metadata.AuditUpdate,
		OperateFrom:     metadata.FromUser,
		OperationDetail: &metadata.GenericOpDetail{Data: preConf, UpdateFields: curConf},
		OperationTime:   metadata.Now(),
		AppCode:         httpheader.GetAppCode(kit.Header),
		RequestID:       kit.Rid,
	}

	if err = s.db.Shard(kit.SysShardOpts()).Table(common.BKTableNameAuditLog).Insert(kit.Ctx, audit); err != nil {
		blog.Errorf("save sharding db config audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

// GetTenantDBRelation get tenant db relation
func (s *Service) GetTenantDBRelation(req *restful.Request, resp *restful.Response) {
	kit := rest.NewKitFromHeader(req.Request.Header, s.CCErr)

	relations := make([]tenant.Tenant, 0)
	err := s.db.Shard(kit.SysShardOpts()).Table(common.BKTableNameTenant).Find(nil).Fields("tenant_id", "database").
		All(kit.Ctx, &relations)
	if err != nil {
		blog.Errorf("get tenant db relations failed, err: %v, rid: %s", err, kit.Rid)
		_ = resp.WriteError(http.StatusOK,
			&metadata.RespError{Msg: kit.CCError.CCError(common.CCErrCommDBSelectFailed)})
		return
	}

	_ = resp.WriteEntity(metadata.NewSuccessResp(relations))
}
