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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	apigwcli "configcenter/src/common/resource/apigw"
	"configcenter/src/scene_server/admin_server/logics"
	"configcenter/src/storage/dal/mongo/local"
	daltypes "configcenter/src/storage/dal/types"
	"configcenter/src/storage/driver/mongodb"
	"github.com/emicklei/go-restful/v3"
)

func (s *Service) addTenant(req *restful.Request, resp *restful.Response) {
	rHeader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(httpheader.GetLanguage(rHeader))
	kit := rest.NewKitFromHeader(rHeader, s.CCErr)

	if !s.Config.EnableMultiTenantMode {
		blog.Errorf("multi-tenant mode is not enabled, cannot add tenant, rid: %s", kit.Rid)
		result := &metadata.RespError{
			Msg: defErr.Errorf(common.CCErrCommAddTenantErr,
				fmt.Errorf("multi-tenant mode is not enabled, cannot add tenant")),
		}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}

	_, exist := tenant.GetTenant(kit.TenantID)
	if exist {
		resp.WriteEntity(metadata.NewSuccessResp("tenant exist"))
		return
	}

	if !s.Config.DisableVerifyTenant {
		// get all tenants from bk-user
		tenants, err := apigwcli.Client().User().GetTenants(kit.Ctx, kit.Header)
		if err != nil {
			blog.Errorf("get tenants from bk-user failed, err: %v, rid: %s", err, kit.Rid)
			result := &metadata.RespError{
				Msg: defErr.Errorf(common.CCErrCommAddTenantErr, fmt.Errorf("get tenants from bk-user failed")),
			}
			resp.WriteError(http.StatusInternalServerError, result)
		}

		tenantMap := make(map[string]types.Status)
		for _, tenant := range tenants {
			tenantMap[tenant.ID] = tenant.Status
		}

		if status, ok := tenantMap[kit.TenantID]; !ok || status != types.EnabledStatus {
			blog.Errorf("tenant %s invalid, rid: %s", kit.TenantID, kit.Rid)
			result := &metadata.RespError{
				Msg: defErr.Errorf(common.CCErrCommAddTenantErr,
					fmt.Errorf("tenant %s invalid", kit.TenantID)),
			}
			resp.WriteError(http.StatusInternalServerError, result)
			return
		}
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

	if err = addTableIndexes(kit, cli); err != nil {
		blog.Errorf("create table and indexes for tenant %s failed, err: %v, rid: %s", kit.TenantID, err, kit.Rid)
		result := &metadata.RespError{
			Msg: defErr.Errorf(common.CCErrCommAddTenantErr, err.Error()),
		}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}

	// add default area
	if err = addDefaultArea(kit, cli); err != nil {
		blog.Errorf("add default area failed, err: %v, rid: %s", err, kit.Rid)
		result := &metadata.RespError{
			Msg: defErr.Errorf(common.CCErrCommAddTenantErr, err.Error()),
		}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}

	if err = addDataFromTemplate(kit, cli); err != nil {
		blog.Errorf("create init data for tenant %s failed, err: %v", kit.TenantID, err)
		result := &metadata.RespError{
			Msg: defErr.Errorf(common.CCErrCommAddTenantErr, err.Error()),
		}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}

	if err = addResPool(kit, cli); err != nil {
		blog.Errorf("add default resouce pool failed, err: %v, rid: %s", err, kit.Rid)
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

	// refresh tenants, ignore refresh tenants error
	if err = logics.RefreshTenants(s.CoreAPI); err != nil {
		blog.Errorf("refresh tenants failed, err: %v, rid: %s", err, kit.Rid)
	}

	resp.WriteEntity(metadata.NewSuccessResp("add tenant success"))
}

var defaultCloudAreas = []metadata.CloudArea{
	{
		CloudID:   common.BKDefaultDirSubArea,
		CloudName: common.DefaultCloudName,
		Status:    "1",
		Default:   int64(common.BuiltIn),
	},
	{
		CloudID:   common.UnassignedCloudAreaID,
		CloudName: common.UnassignedCloudAreaName,
		Default:   int64(common.BuiltIn),
	},
}

// addDefaultArea add default cloud areas
func addDefaultArea(kit *rest.Kit, db local.DB) error {
	cond := map[string]interface{}{common.BKDefaultField: common.BuiltIn}
	existCloudAreas := make([]metadata.CloudArea, 0)
	err := db.Table(common.BKTableNameBasePlat).Find(cond).Fields(common.BKCloudIDField).All(kit.Ctx, &existCloudAreas)
	if err != nil {
		blog.Errorf("get default area count failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	if len(existCloudAreas) == len(defaultCloudAreas) {
		return nil
	}

	existCloudAreaMap := make(map[int64]struct{})
	for _, area := range existCloudAreas {
		existCloudAreaMap[area.CloudID] = struct{}{}
	}

	now := time.Now()
	createCloudAreas := make([]metadata.CloudArea, 0)
	for _, cloudArea := range defaultCloudAreas {
		_, exists := existCloudAreaMap[cloudArea.CloudID]
		if exists {
			continue
		}
		cloudArea.Default = int64(common.BuiltIn)
		cloudArea.Creator = common.CCSystemOperatorUserName
		cloudArea.LastEditor = common.CCSystemOperatorUserName
		cloudArea.CreateTime = now
		cloudArea.LastTime = now
		createCloudAreas = append(createCloudAreas, cloudArea)
	}

	if err = db.Table(common.BKTableNameBasePlat).Insert(kit.Ctx, createCloudAreas); err != nil {
		blog.Errorf("add default area failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	return nil
}

func addResPoolData(kit *rest.Kit, db local.DB, objID string, data mapstr.MapStr) (int64, error) {
	table := common.GetInstTableName(objID, kit.TenantID)
	idField := common.GetInstIDField(objID)

	cond := mapstr.MapStr{common.BKDefaultField: data[common.BKDefaultField]}
	existData := make([]map[string]int64, 0)
	err := db.Table(table).Find(cond).Fields(idField).All(kit.Ctx, &existData)
	if err != nil {
		blog.Errorf("get exist resource pool %s failed, err: %v, rid: %s", objID, err, kit.Rid)
		return 0, err
	}

	if len(existData) > 0 {
		return existData[0][idField], nil
	}

	id, err := mongodb.Shard(kit.SysShardOpts()).NextSequence(kit.Ctx, table)
	if err != nil {
		blog.Errorf("get next sequence for table %s failed, err: %v, rid: %s", table, err, kit.Rid)
		return 0, err
	}

	data[idField] = id
	err = db.Table(table).Insert(kit.Ctx, data)
	if err != nil {
		blog.Errorf("create resource pool %s data(%+v) failed, err: %v, rid: %s", objID, data, err, kit.Rid)
		return 0, err
	}
	return int64(id), nil
}

func addResPool(kit *rest.Kit, db local.DB) error {
	now := time.Now()

	bizID, err := addResPoolData(kit, db, common.BKInnerObjIDApp, mapstr.MapStr{
		common.BKAppNameField:     common.DefaultAppName,
		common.BKMaintainersField: "admin",
		common.BKProductPMField:   "admin",
		common.BKTimeZoneField:    "Asia/Shanghai",
		common.BKLanguageField:    "1",
		common.BKLifeCycleField:   common.DefaultAppLifeCycleNormal,
		common.BKDefaultField:     common.DefaultAppFlag,
		common.BKDeveloperField:   "",
		common.BKTesterField:      "",
		common.BKOperatorField:    "",
		common.CreateTimeField:    now,
		common.LastTimeField:      now,
	})
	if err != nil {
		return err
	}

	setID, err := addResPoolData(kit, db, common.BKInnerObjIDSet, mapstr.MapStr{
		common.BKAppIDField:         bizID,
		common.BKInstParentStr:      bizID,
		common.BKSetNameField:       common.DefaultResSetName,
		common.BKDefaultField:       common.DefaultResSetFlag,
		common.BKSetEnvField:        "3",
		common.BKSetStatusField:     "1",
		common.BKSetDescField:       "",
		common.BKSetTemplateIDField: 0,
		common.BKSetCapacityField:   nil,
		common.BKDescriptionField:   "",
		common.CreateTimeField:      now,
		common.LastTimeField:        now,
	})
	if err != nil {
		return err
	}

	// get default service category
	cond := map[string]interface{}{
		common.BKFieldName:     common.DefaultServiceCategoryName,
		common.BKParentIDField: mapstr.MapStr{common.BKDBNE: 0},
	}
	defCategory := new(metadata.ServiceCategory)
	err = db.Table(common.BKTableNameServiceCategory).Find(cond).Fields(common.BKFieldID).One(kit.Ctx, &defCategory)
	if err != nil {
		blog.Errorf("get default service category by cond(%+v) failed, err: %v, rid: %s", cond, err, kit.Rid)
		return err
	}

	_, err = addResPoolData(kit, db, common.BKInnerObjIDModule, mapstr.MapStr{
		common.BKAppIDField:             bizID,
		common.BKSetIDField:             setID,
		common.BKInstParentStr:          setID,
		common.BKModuleNameField:        common.DefaultResModuleName,
		common.BKDefaultField:           common.DefaultResModuleFlag,
		common.BKServiceTemplateIDField: common.ServiceTemplateIDNotSet,
		common.BKSetTemplateIDField:     common.SetTemplateIDNotSet,
		common.BKServiceCategoryIDField: defCategory.ID,
		common.BKModuleTypeField:        "1",
		common.BKOperatorField:          "",
		common.BKBakOperatorField:       "",
		common.HostApplyEnabledField:    false,
		common.CreateTimeField:          now,
		common.LastTimeField:            now,
	})
	if err != nil {
		return err
	}

	return nil
}

// addTableIndexes add table indexes
func addTableIndexes(kit *rest.Kit, db local.DB) error {
	for table, index := range index.TableIndexes() {
		if err := addOneTableIndexes(kit, db, table, index); err != nil {
			return err
		}
	}

	for _, object := range common.BKInnerObjects {
		instAsstTable := common.GetObjectInstAsstTableName(object, kit.TenantID)
		if err := addOneTableIndexes(kit, db, instAsstTable, index.InstanceAssociationIndexes()); err != nil {
			return err
		}
	}

	return nil
}

// addOneTableIndexes add table indexes for one table
func addOneTableIndexes(kit *rest.Kit, db local.DB, table string, indexes []daltypes.Index) error {
	dbCli := db
	if common.IsPlatformTable(table) {
		dbCli = mongodb.Shard(kit.SysShardOpts())
	}

	if err := logics.CreateTable(kit, dbCli, table); err != nil {
		blog.Errorf("create table %s failed, err: %v, rid: %s", table, err, kit.Rid)
		return err
	}

	if err := logics.CreateIndexes(kit, dbCli, table, indexes); err != nil {
		blog.Errorf("create index %s failed, err: %v, rid: %s", table, err, kit.Rid)
		return err
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
