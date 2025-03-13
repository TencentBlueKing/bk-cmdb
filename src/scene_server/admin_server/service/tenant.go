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
	"time"

	"configcenter/pkg/tenant"
	"configcenter/pkg/tenant/types"
	tenanttmp "configcenter/pkg/types/tenant-template"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/index"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/logics"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/driver/mongodb"

	"github.com/emicklei/go-restful/v3"
)

func (s *Service) addTenant(req *restful.Request, resp *restful.Response) {

	rHeader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(httpheader.GetLanguage(rHeader))
	kit := rest.NewKitFromHeader(rHeader, s.CCErr)

	_, exist := tenant.GetTenant(kit.TenantID)
	if exist {
		resp.WriteEntity(metadata.NewSuccessResp("tenant exist"))
		return
	}

	cli, dbUUID, err := logics.GetNewTenantCli(kit, mongodb.Dal())
	if err != nil {
		blog.Errorf("get new tenant db failed, err: %v, rid: %s", err, kit.Rid)
		result := &metadata.RespError{
			Msg: defErr.Errorf(common.CCErrCommAddTenantErr, fmt.Errorf("get new tenant db failed")),
		}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}

	if err := addTableIndexes(kit, cli); err != nil {
		blog.Errorf("create table and indexes for tenant %s failed, err: %v, rid: %s", kit.TenantID, err, kit.Rid)
		result := &metadata.RespError{
			Msg: defErr.Errorf(common.CCErrCommAddTenantErr, err.Error()),
		}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}

	// add default area
	if err := addDefaultArea(kit, cli); err != nil {
		blog.Errorf("add default area failed, err: %v, rid: %s", err, kit.Rid)
		result := &metadata.RespError{
			Msg: defErr.Errorf(common.CCErrCommAddTenantErr, err.Error()),
		}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}

	if err := addDataFromTemplate(kit, cli); err != nil {
		blog.Errorf("create init data for tenant %s failed, err: %v", kit.TenantID, err)
		result := &metadata.RespError{
			Msg: defErr.Errorf(common.CCErrCommAddTenantErr, err.Error()),
		}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}

	// add tenant db relation
	data := &types.Tenant{
		TenantID: kit.TenantID,
		Database: dbUUID,
		Status:   types.EnabledStatus,
	}
	err = mongodb.Shard(kit.SysShardOpts()).Table(common.BKTableNameTenant).Insert(kit.Ctx, data)
	if err != nil {
		blog.Errorf("add tenant db relations failed, err: %v, rid: %s", err, kit.Rid)
		result := &metadata.RespError{
			Msg: defErr.Errorf(common.CCErrCommAddTenantErr, err.Error()),
		}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}

	resp.WriteEntity(metadata.NewSuccessResp("add tenant success"))
}

func addDefaultArea(kit *rest.Kit, db local.DB) error {
	// add default area
	cond := map[string]interface{}{"bk_cloud_name": "Default Area"}
	cnt, err := db.Table(common.BKTableNameBasePlat).Find(cond).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("get default area count failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	if cnt == 1 {
		return nil
	}

	err = db.Table(common.BKTableNameBasePlat).Insert(kit.Ctx, metadata.CloudArea{
		Creator:    common.CCSystemOperatorUserName,
		LastEditor: common.CCSystemOperatorUserName,
		CloudID:    common.BKDefaultDirSubArea,
		CloudName:  "Default Area",
		Default:    int64(common.BuiltIn),
		CreateTime: time.Now(),
		LastTime:   time.Now(),
	})
	if err != nil {
		blog.Errorf("add default area failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	return nil
}

// addTableIndexes add table indexes
func addTableIndexes(kit *rest.Kit, db local.DB) error {
	tableIndexes := index.TableIndexes()
	for _, object := range tenanttmp.BKInnerObjects {
		instAsstTable := common.GetObjectInstAsstTableName(object, kit.TenantID)
		tableIndexes[instAsstTable] = index.InstanceAssociationIndexes()
	}

	for table, index := range tableIndexes {
		dbCli := db
		if common.IsPlatformTable(table) {
			dbCli = mongodb.Shard(kit.SysShardOpts())
		}
		if err := logics.CreateTable(kit, dbCli, table); err != nil {
			blog.Errorf("create table %s failed, err: %v, rid: %s", table, err, kit.Rid)
			return err
		}

		if err := logics.CreateIndexes(kit, dbCli, table, index); err != nil {
			blog.Errorf("create table %s failed, err: %v, rid: %s", table, err, kit.Rid)
			return err
		}
	}

	return nil
}

func addDataFromTemplate(kit *rest.Kit, db local.DB) error {

	for _, ty := range tenanttmp.AllTemplateTypes {
		if err := typeHandlerMap[ty](kit, db); err != nil {
			blog.Errorf("add template data failed for type %s, err: %v, rid: %s", ty, err, kit.Rid)
			return err
		}
	}

	return nil
}
