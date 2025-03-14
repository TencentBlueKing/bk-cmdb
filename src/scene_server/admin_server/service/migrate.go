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
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"configcenter/pkg/tenant"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	commontype "configcenter/src/common/types"
	"configcenter/src/common/util"
	"configcenter/src/common/version"
	"configcenter/src/common/watch"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/source_controller/cacheservice/event"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/dal/mongo/sharding"
	daltypes "configcenter/src/storage/dal/types"
	"configcenter/src/storage/driver/mongodb"
	streamtypes "configcenter/src/storage/stream/types"

	"github.com/emicklei/go-restful/v3"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *Service) migrateDatabase(req *restful.Request, resp *restful.Response) {
	kit := rest.NewKitFromHeader(req.Request.Header, s.CCErr)

	result, err := upgrader.Upgrade(kit, s.db, nil)
	if err != nil {
		blog.Errorf("db upgrade failed, err: %v", err)
		result := &metadata.RespError{
			Msg: kit.CCError.Errorf(common.CCErrCommMigrateFailed, err.Error()),
		}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}

	if err = s.createWatchDBChainCollections(kit); err != nil {
		blog.Errorf("create watch db chain collections failed, err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{
			Msg: kit.CCError.Errorf(common.CCErrCommMigrateFailed, err.Error()),
		})
		return
	}

	resp.WriteEntity(metadata.NewSuccessResp(result))
}

// dbChainTTLTime the ttl time seconds of the db event chain, used to set the ttl index of mongodb
const dbChainTTLTime = 5 * 24 * 60 * 60

func (s *Service) createWatchDBChainCollections(kit *rest.Kit) error {
	watchDBToDBRelation, err := s.getWatchDBToDBRelation(kit)
	if err != nil {
		return err
	}

	// create watch token table and init the watch token for dbs
	if err := s.createWatchToken(kit, watchDBToDBRelation); err != nil {
		return err
	}

	// create watch chain node table and init the last token info as empty for all collections
	cursorTypes := watch.ListCursorTypes()
	for _, cursorType := range cursorTypes {
		key, err := event.GetResourceKeyWithCursorType(cursorType)
		if err != nil {
			blog.Errorf("get resource key with cursor type %s failed, err: %v, rid: %s", cursorType, err, kit.Rid)
			return err
		}

		err = tenant.ExecForAllTenants(func(tenantID string) error {
			// TODO 在新增租户初始化时同时增加watch相关表，并刷新cache的tenant
			kit = kit.NewKit().WithTenant(tenantID)
			exists, err := s.watchDB.Shard(kit.ShardOpts()).HasTable(s.ctx, key.ChainCollection())
			if err != nil {
				blog.Errorf("check if table %s exists failed, err: %v, rid: %s", key.ChainCollection(), err, kit.Rid)
				return err
			}

			if !exists {
				err = s.watchDB.Shard(kit.ShardOpts()).CreateTable(s.ctx, key.ChainCollection())
				if err != nil && !mongodb.IsDuplicatedError(err) {
					blog.Errorf("create table %s failed, err: %v, rid: %s", key.ChainCollection(), err, kit.Rid)
					return err
				}
			}

			if err = s.createWatchIndexes(kit, cursorType, key); err != nil {
				return err
			}

			if err = s.createLastWatchEvent(kit, key); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}

		// TODO 在新增DB时同时增加db relation和token数据
		err = s.createWatchTokenForEventKey(kit, key, watchDBToDBRelation)
		if err != nil {
			return err
		}
	}
	return nil
}

// getWatchDBToDBRelation get watch db uuid to db uuids relation
func (s *Service) getWatchDBToDBRelation(kit *rest.Kit) (map[string][]string, error) {
	// get all db uuids
	uuidMap := make(map[string]struct{})
	err := s.db.ExecForAllDB(func(db local.DB) error {
		dbClient, ok := db.(*local.Mongo)
		if !ok {
			return fmt.Errorf("db to be watched is not an instance of local mongo")
		}
		uuidMap[dbClient.GetMongoClient().UUID()] = struct{}{}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// get watch db relations
	relations := make([]sharding.WatchDBRelation, 0)
	if err := s.watchDB.Shard(kit.SysShardOpts()).Table(common.BKTableNameWatchDBRelation).Find(nil).
		All(kit.Ctx, &relations); err != nil {
		blog.Errorf("get watch db relation failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	watchDBToDBRelation := make(map[string][]string)
	for _, relation := range relations {
		watchDBToDBRelation[relation.WatchDB] = append(watchDBToDBRelation[relation.WatchDB], relation.DB)
		delete(uuidMap, relation.DB)
	}

	// get default watch db uuid for new db to be watched
	cond := map[string]any{common.MongoMetaID: common.ShardingDBConfID}
	conf := new(sharding.ShardingDBConf)
	err = s.watchDB.Shard(kit.SysShardOpts()).Table(common.BKTableNameSystem).Find(cond).One(kit.Ctx, &conf)
	if err != nil {
		blog.Errorf("get sharding db conf failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}
	defaultWatchDBUUID := conf.ForNewData

	// create watch db relation for dbs without watch db
	newRelations := make([]sharding.WatchDBRelation, 0)
	for uuid := range uuidMap {
		watchDBToDBRelation[defaultWatchDBUUID] = append(watchDBToDBRelation[defaultWatchDBUUID], uuid)
		newRelations = append(newRelations, sharding.WatchDBRelation{
			WatchDB: defaultWatchDBUUID,
			DB:      uuid,
		})
	}

	if len(newRelations) > 0 {
		err = s.watchDB.Shard(kit.SysShardOpts()).Table(common.BKTableNameWatchDBRelation).Insert(kit.Ctx, newRelations)
		if err != nil {
			blog.Errorf("create watch db relations(%+v) failed, err: %v, rid: %s", newRelations, err, kit.Rid)
			return nil, err
		}
	}

	return watchDBToDBRelation, nil
}

func (s *Service) createWatchToken(kit *rest.Kit, watchDBToDBRelation map[string][]string) error {
	return s.watchDB.ExecForAllDB(func(watchDB local.DB) error {
		// create watch token table to store the last watch token info for db and every collection
		exists, err := watchDB.HasTable(s.ctx, common.BKTableNameWatchToken)
		if err != nil {
			blog.Errorf("check if table %s exists failed, err: %v, rid: %s", common.BKTableNameWatchToken, err, kit.Rid)
			return err
		}

		if !exists {
			err = watchDB.CreateTable(s.ctx, common.BKTableNameWatchToken)
			if err != nil && !mongodb.IsDuplicatedError(err) {
				blog.Errorf("create table %s failed, err: %v, rid: %s", common.BKTableNameWatchToken, err, kit.Rid)
				return err
			}
		}

		// get all exist db watch tokens
		mongo, ok := watchDB.(*local.Mongo)
		if !ok {
			return fmt.Errorf("db is not *local.Mongo type")
		}
		uuids := watchDBToDBRelation[mongo.GetMongoClient().UUID()]
		if len(uuids) == 0 {
			return nil
		}

		filter := map[string]interface{}{
			common.MongoMetaID: map[string]interface{}{common.BKDBIN: uuids},
		}

		existUUIDs, err := watchDB.Table(common.BKTableNameWatchToken).Distinct(kit.Ctx, common.MongoMetaID, filter)
		if err != nil {
			blog.Errorf("check if dbs(%+v) watch token exists failed, err: %v, rid: %s", uuids, err, kit.Rid)
			return err
		}

		existUUIDMap := make(map[string]struct{})
		for _, uuid := range existUUIDs {
			existUUIDMap[util.GetStrByInterface(uuid)] = struct{}{}
		}

		// create watch token for dbs to be watched
		for _, uuid := range uuids {
			if _, exists := existUUIDMap[uuid]; exists {
				continue
			}

			data := mapstr.MapStr{
				common.MongoMetaID:  uuid,
				common.BKTokenField: "",
				common.BKStartAtTimeField: streamtypes.TimeStamp{
					Sec:  uint32(time.Now().Unix()),
					Nano: 0,
				},
			}
			if err = watchDB.Table(common.BKTableNameWatchToken).Insert(s.ctx, data); err != nil {
				blog.Errorf("create db watch token failed, err: %v, data: %+v, rid: %s", err, data, kit.Rid)
				return err
			}
		}
		return nil
	})
}

func (s *Service) createWatchIndexes(kit *rest.Kit, cursorType watch.CursorType, key event.Key) error {
	indexes := []daltypes.Index{
		{Name: "index_id", Keys: bson.D{{common.BKFieldID, -1}}, Background: true, Unique: true},
		{Name: "index_cursor", Keys: bson.D{{common.BKCursorField, -1}}, Background: true, Unique: true},
		{Name: "index_cluster_time", Keys: bson.D{{common.BKClusterTimeField, -1}}, Background: true,
			ExpireAfterSeconds: dbChainTTLTime},
	}

	if cursorType == watch.ObjectBase || cursorType == watch.MainlineInstance || cursorType == watch.InstAsst {
		subResourceIndex := daltypes.Index{
			Name: "index_sub_resource", Keys: bson.D{{common.BKSubResourceField, 1}}, Background: true,
		}
		indexes = append(indexes, subResourceIndex)
	}

	existIndexArr, err := s.watchDB.Shard(kit.ShardOpts()).Table(key.ChainCollection()).Indexes(s.ctx)
	if err != nil {
		blog.Errorf("get exist indexes for table %s failed, err: %v, rid: %s", key.ChainCollection(), err, kit.Rid)
		return err
	}

	existIdxMap := make(map[string]bool)
	for _, index := range existIndexArr {
		existIdxMap[index.Name] = true
	}

	for _, index := range indexes {
		if _, exist := existIdxMap[index.Name]; exist {
			continue
		}

		err = s.watchDB.Shard(kit.ShardOpts()).Table(key.ChainCollection()).CreateIndex(s.ctx, index)
		if err != nil && !mongodb.IsDuplicatedError(err) {
			blog.Errorf("create indexes for table %s failed, err: %v, rid: %s", key.ChainCollection(), err, kit.Rid)
			return err
		}
	}
	return nil
}

func (s *Service) createWatchTokenForEventKey(kit *rest.Kit, key event.Key,
	watchDBToDBRelation map[string][]string) error {

	// create watch token of this key for every db
	err := s.watchDB.ExecForAllDB(func(db local.DB) error {
		mongo, ok := db.(*local.Mongo)
		if !ok {
			return fmt.Errorf("db is not *local.Mongo type")
		}

		for _, uuid := range watchDBToDBRelation[mongo.GetMongoClient().UUID()] {
			if err := s.createWatchTokenForDB(kit, db, uuid, key); err != nil {
				blog.Errorf("init %s key %s watch token failed, err: %v, rid: %s", uuid, key.Namespace(), err, kit.Rid)
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) createWatchTokenForDB(kit *rest.Kit, watchDB local.DB, uuid string, key event.Key) error {
	filter := map[string]interface{}{
		"_id": watch.GenDBWatchTokenID(uuid, key.Collection()),
	}

	count, err := watchDB.Table(common.BKTableNameWatchToken).Find(filter).Count(s.ctx)
	if err != nil {
		blog.Errorf("check if last watch token exists failed, err: %v, filter: %+v", err, filter)
		return err
	}

	if count > 0 {
		return nil
	}

	if key.Collection() == event.HostIdentityKey.Collection() {
		// host identity's watch token is different with other identity.
		// only set coll is ok, the other fields is useless
		data := mapstr.MapStr{
			"_id":                                     watch.GenDBWatchTokenID(uuid, key.Collection()),
			common.BKTableNameBaseHost:                new(streamtypes.TokenInfo),
			common.BKTableNameModuleHostConfig:        new(streamtypes.TokenInfo),
			common.BKTableNameBaseProcess:             new(streamtypes.TokenInfo),
			common.BKTableNameProcessInstanceRelation: new(streamtypes.TokenInfo),
		}
		err = watchDB.Table(common.BKTableNameWatchToken).Insert(s.ctx, data)
		if err != nil {
			blog.Errorf("init last watch token failed, err: %v, data: %+v, rid: %s", err, data, kit.Rid)
			return err
		}
		return nil
	}

	if key.Collection() == event.BizSetRelationKey.Collection() {
		// biz set relation's watch token is generated in the same way with the host identity's watch token
		data := mapstr.MapStr{
			"_id":                        watch.GenDBWatchTokenID(uuid, key.Collection()),
			common.BKTableNameBaseApp:    new(streamtypes.TokenInfo),
			common.BKTableNameBaseBizSet: new(streamtypes.TokenInfo),
		}
		err = watchDB.Table(common.BKTableNameWatchToken).Insert(s.ctx, data)
		if err != nil {
			blog.Errorf("init last biz set rel watch token failed, err: %v, data: %+v, rid: %s", err, data, kit.Rid)
			return err
		}
		return nil
	}

	data := mapstr.MapStr{
		common.MongoMetaID:  watch.GenDBWatchTokenID(uuid, key.Collection()),
		common.BKTokenField: "",
		common.BKStartAtTimeField: streamtypes.TimeStamp{
			Sec:  uint32(time.Now().Unix()),
			Nano: 0,
		},
	}
	if err = watchDB.Table(common.BKTableNameWatchToken).Insert(s.ctx, data); err != nil {
		blog.Errorf("init last watch token failed, err: %v, data: %+v, rid: %s", err, data, kit.Rid)
		return err
	}
	return nil
}

func (s *Service) createLastWatchEvent(kit *rest.Kit, key event.Key) error {
	filter := map[string]interface{}{
		"_id": key.Collection(),
	}

	count, err := s.watchDB.Shard(kit.ShardOpts()).Table(common.BKTableNameLastWatchEvent).Find(filter).Count(s.ctx)
	if err != nil {
		blog.Errorf("check if last watch event exists failed, err: %v, filter: %+v, rid: %s", err, filter, kit.Rid)
		return err
	}

	if count > 0 {
		return nil
	}

	data := watch.LastChainNodeData{
		Coll:   key.Collection(),
		ID:     0,
		Cursor: "",
	}
	if err = s.watchDB.Shard(kit.ShardOpts()).Table(common.BKTableNameLastWatchEvent).Insert(s.ctx, data); err != nil {
		blog.Errorf("create last watch event failed, err: %v, data: %+v, rid: %s", err, data, kit.Rid)
		return err
	}
	return nil
}

func (s *Service) migrateSpecifyVersion(req *restful.Request, resp *restful.Response) {
	rHeader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(httpheader.GetLanguage(rHeader))
	kit := rest.NewKitFromHeader(rHeader, s.CCErr)
	input := new(MigrateSpecifyVersionRequest)
	if err := json.NewDecoder(req.Request.Body).Decode(input); err != nil {
		blog.Errorf("migrateSpecifyVersion failed, decode body err: %v, body: %+v, rid: %s", err, req.Request.Body,
			kit.Rid)
		_ = resp.WriteError(http.StatusOK, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if input.CommitID != version.CCGitHash {
		_ = resp.WriteError(http.StatusOK,
			&metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsInvalid, "commit_id")})
		return
	}

	err := upgrader.UpgradeSpecifyVersion(kit, s.db, input.Version)
	if err != nil {
		blog.Errorf("db upgrade specify failed, err: %+v, rid: %s", err, kit.Rid)
		result := &metadata.RespError{
			Msg: defErr.Errorf(common.CCErrCommMigrateFailed, err.Error()),
		}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}

	result := MigrationResponse{
		BaseResp: metadata.BaseResp{
			Result:      true,
			Code:        0,
			ErrMsg:      "",
			Permissions: nil,
		},
		Data: "migrate success. version: " + input.Version,
	}
	resp.WriteEntity(result)

}

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
	rid := httpheader.GetRid(rHeader)
	defErr := s.CCErr.CreateDefaultCCErrorIf(httpheader.GetLanguage(rHeader))

	input := new(struct {
		ConfigName string `json:"config_name"`
	})
	if err := json.NewDecoder(req.Request.Body).Decode(input); err != nil {
		blog.Errorf("refreshConfig failed, decode body err: %v, body: %+v, rid:%s", err, req.Request.Body, rid)
		resp.WriteError(http.StatusOK, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	configName := "all"
	if input.ConfigName != "" {
		if ok := allConfigNames[input.ConfigName]; !ok {
			blog.Errorf("refreshConfig failed, configHelpInfo: %s, input: %#v, rid: %s", configHelpInfo, input, rid)
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
		key := commontype.CC_SERVCONF_BASEPATH + "/" + configName
		err = s.ConfigCenter.WriteConfigure(filePath, key)
	case "error":
		err = s.ConfigCenter.WriteErrorRes2Center(s.Config.Errors.Res)
	case "language":
		err = s.ConfigCenter.WriteLanguageRes2Center(s.Config.Language.Res)
	case "all":
		err = s.ConfigCenter.WriteAllConfs2Center(s.Config.Configures.Dir, s.Config.Errors.Res, s.Config.Language.Res)
	default:
		blog.Errorf("refreshConfig failed, config_name is wrong, configHelpInfo: %s, input: %#v, rid: %s",
			configHelpInfo, input, rid)
		resp.WriteError(http.StatusOK,
			&metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsInvalid, configHelpInfo)})
		return
	}

	if err != nil {
		blog.Warnf("refreshConfig failed, input: %#v, error: %v, rid: %s", input, err, rid)
		resp.WriteError(http.StatusOK, &metadata.RespError{Msg: err})
	}

	blog.Infof("refresh config success, input: %#v", input)
	resp.WriteEntity(metadata.NewSuccessResp("refresh config success"))
}
