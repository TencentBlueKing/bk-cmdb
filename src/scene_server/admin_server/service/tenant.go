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
	"fmt"
	"net/http"

	"configcenter/pkg/tenant"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/index"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/logics"
	"configcenter/src/scene_server/admin_server/upgrader/tools"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/dal/mongo/sharding"
	"configcenter/src/storage/driver/mongodb"
	"github.com/emicklei/go-restful/v3"
)

var (
	objects = []string{
		common.BKInnerObjIDApp,
		common.BKInnerObjIDModule,
		common.BKProcessObjectName,
		common.BKInnerObjIDHost,
		common.BKInnerObjIDProject,
		common.BKInnerObjIDBizSet,
		common.BKInnerObjIDPlat,
		common.BKInnerObjIDSet,
	}
)

func (s *Service) addTenant(req *restful.Request, resp *restful.Response) {

	rHeader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(httpheader.GetLanguage(rHeader))
	kit := rest.NewKitFromHeader(rHeader, s.CCErr)

	cli, err := logics.GetNewTenantCli(kit)
	if cli == nil || err != nil {
		blog.Errorf("get new tenant client failed, err: %v", err)
		result := &metadata.RespError{
			Msg: defErr.Errorf(common.CCErrorUnknownOrUnrecognizedError, err.Error()),
		}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}

	dbName, err := checkTenantInfo(kit)
	if err != nil {
		blog.Errorf("tenant %s already exist in the db, err: %v", kit.TenantID, err)
		result := &metadata.RespError{
			Msg: defErr.Errorf(common.CCErrorUnknownOrUnrecognizedError, err.Error()),
		}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}

	if err = addTableIndexes(kit, cli); err != nil {
		blog.Errorf("create table and indexes for tenant %s failed, err: %v", kit.TenantID, err)
		result := &metadata.RespError{
			Msg: defErr.Errorf(common.CCErrorUnknownOrUnrecognizedError, err.Error()),
		}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}

	if err = addDataFromTemplate(kit, cli); err != nil {
		blog.Errorf("create init data for tenant %s failed, err: %v", kit.TenantID, err)
		result := &metadata.RespError{
			Msg: defErr.Errorf(common.CCErrorUnknownOrUnrecognizedError, err.Error()),
		}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}

	// add tenant db relation
	data := &tenant.Tenant{
		TenantID: kit.TenantID,
		Database: dbName,
	}
	err = mongodb.Shard(kit.SysShardOpts()).Table(common.BKTableNameTenant).Insert(kit.Ctx, data)
	if err != nil {
		blog.Errorf("add tenant db relations failed, err: %v, rid: %s", err, kit.Rid)
		result := &metadata.RespError{
			Msg: defErr.Errorf(common.CCErrorUnknownOrUnrecognizedError, err.Error()),
		}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}

	resp.WriteEntity(metadata.NewSuccessResp("add tenant success"))
}

// checkTenantInfo add tenant db relation info
func checkTenantInfo(kit *rest.Kit) (string, error) {
	conf, err := GetDBConfig(kit)
	if err != nil {
		return "", err
	}
	if len(conf.ForNewTenant) == 0 {
		blog.Errorf("invalid new tenant conf, tenant db uuid is empty, rid: %s", kit.Rid)
		return "", fmt.Errorf("invalid new tenant conf, tenant db uuid is empty")
	}

	cond := &mapstr.MapStr{
		common.TenantID: kit.TenantID,
		"database":      conf.ForNewTenant,
	}
	count, err := mongodb.Shard(kit.SysShardOpts()).Table(common.BKTableNameTenant).Find(cond).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("get tenant db relations failed, err: %v, rid: %s", err, kit.Rid)
		return "", err
	}
	if err != nil {
		return "", err
	}
	if count > 0 {
		blog.Errorf("invalid new tenant conf, tenant db relation exist, rid: %s", kit.Rid)
		return "", fmt.Errorf("invalid new tenant conf, tenant db relation exist")
	}

	return conf.ForNewTenant, nil
}

// GetDBConfig get db config
func GetDBConfig(kit *rest.Kit) (*sharding.ShardingDBConf, error) {
	cond := map[string]any{common.MongoMetaID: common.ShardingDBConfID}
	config := new(sharding.ShardingDBConf)
	err := mongodb.Shard(kit.SysShardOpts()).Table(common.BKTableNameSystem).Find(cond).One(kit.Ctx, &config)
	if err != nil {
		blog.Errorf("get tenant db config failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	return config, nil
}

// addTableIndexes add table indexes
func addTableIndexes(kit *rest.Kit, db *local.Mongo) error {
	tableIndexes := index.TableIndexes()
	for _, object := range objects {
		instAsstTable := common.GetObjectInstAsstTableName(object, kit.TenantID)
		tableIndexes[instAsstTable] = index.InstanceAssociationIndexes()
	}

	for table, index := range tableIndexes {
		if err := tools.CreateTable(kit, db, table); err != nil {
			blog.Errorf("create table %s failed, err: %v", table, err)
			return err
		}

		if err := tools.CreateIndexes(kit, db, table, index); err != nil {
			blog.Errorf("create table %s failed, err: %v", table, err)
			return err
		}
	}

	return nil
}

func addDataFromTemplate(kit *rest.Kit, db *local.Mongo) error {

	count, err := mongodb.Shard(kit.SysShardOpts()).Table(common.BKTableNameTenantTemplate).Find(mapstr.MapStr{}).
		Count(kit.Ctx)
	if err != nil {
		blog.Errorf("get template data count failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	templateData := make(map[string][]mapstr.MapStr)
	for offset := 0; offset < int(count); offset += common.BKMaxInstanceLimit {
		result := make([]tools.TemplateData, 0)
		err = mongodb.Shard(kit.SysShardOpts()).Table(common.BKTableNameTenantTemplate).Find(mapstr.MapStr{}).
			Start(uint64(offset)).Limit(uint64(common.BKMaxInstanceLimit)).All(kit.Ctx, &result)
		if err != nil {
			blog.Errorf("get template data failed, err: %v, rid: %s", err, kit.Rid)
			return err
		}

		for _, data := range result {
			table := typeTableMap[data.Type]
			if dataList, exists := templateData[table]; exists {
				templateData[table] = append(dataList, data.Data)
			} else {
				templateData[table] = []mapstr.MapStr{data.Data}
			}
		}
	}

	for _, table := range tableSeq {
		if err = dataInitiator[table](kit, db, table, templateData[table], tableFieldsMap[table]); err != nil {
			blog.Errorf("add template data failed for table %s, err: %v, rid: %s", table, err, kit.Rid)
			return err
		}
	}
	return nil
}
