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
	"strings"

	"configcenter/pkg/tenant"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/index"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/admin_server/logics"
	"configcenter/src/scene_server/admin_server/service/utils"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/driver/mongodb"

	"github.com/emicklei/go-restful/v3"
)

func (s *Service) addTenant(req *restful.Request, resp *restful.Response) {

	rHeader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(httpheader.GetLanguage(rHeader))
	kit := rest.NewKitFromHeader(rHeader, s.CCErr)

	cli := logics.GetNewTenantCli(kit)
	if cli == nil {
		blog.Errorf("get new tenant client failed, rid: %s", kit.Rid)
		result := &metadata.RespError{
			Msg: defErr.Errorf(common.CCErrorUnknownOrUnrecognizedError, fmt.Errorf("get new tenant client failed")),
		}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}

	if err := addTableIndexes(kit, cli); err != nil {
		blog.Errorf("create table and indexes for tenant %s failed, err: %v, rid: %s", kit.TenantID, err, kit.Rid)
		result := &metadata.RespError{
			Msg: defErr.Errorf(common.CCErrorUnknownOrUnrecognizedError, err.Error()),
		}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}

	if err := addDataFromTemplate(kit, cli); err != nil {
		blog.Errorf("create init data for tenant %s failed, err: %v", kit.TenantID, err)
		result := &metadata.RespError{
			Msg: defErr.Errorf(common.CCErrorUnknownOrUnrecognizedError, err.Error()),
		}
		resp.WriteError(http.StatusInternalServerError, result)
		return
	}

	// add tenant db relation
	dbName := logics.GetNewTenantDBName()
	data := &tenant.Tenant{
		TenantID: kit.TenantID,
		Database: dbName,
		Status:   tenant.EnabledStatus,
	}
	err := mongodb.Shard(kit.SysShardOpts()).Table(common.BKTableNameTenant).Insert(kit.Ctx, data)
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

// addTableIndexes add table indexes
func addTableIndexes(kit *rest.Kit, db local.DB) error {
	tableIndexes := index.TableIndexes()
	for _, object := range common.Objects {
		instAsstTable := common.GetObjectInstAsstTableName(object, kit.TenantID)
		tableIndexes[instAsstTable] = index.InstanceAssociationIndexes()
	}

	for table, index := range tableIndexes {
		if err := logics.CreateTable(kit, db, table); err != nil {
			blog.Errorf("create table %s failed, err: %v, rid: %s", table, err, kit.Rid)
			return err
		}

		if err := logics.CreateIndexes(kit, db, table, index); err != nil {
			blog.Errorf("create table %s failed, err: %v, rid: %s", table, err, kit.Rid)
			return err
		}
	}

	return nil
}

func addDataFromTemplate(kit *rest.Kit, db local.DB) error {

	templateMap := map[string][]metadata.TemplateData{}
	dataMap := map[string][]mapstr.MapStr{}
	for ty := range typeHandlerMap {
		lastId := 0
		hasMore := true
		for hasMore {
			filter := mapstr.MapStr{
				"type": ty,
				"id":   map[string]interface{}{common.BKDBGT: lastId},
			}
			result := make([]metadata.TemplateData, 0)
			err := mongodb.Shard(kit.SysShardOpts()).Table(common.BKTableNameTenantTemplate).Find(filter).
				Sort("id").Limit(uint64(common.BKMaxInstanceLimit)).All(kit.Ctx, &result)
			if err != nil {
				blog.Errorf("get template data for type %s failed, err: %v, rid: %s", ty, err, kit.Rid)
				return err
			}

			if len(result) > 0 {
				templateMap[ty] = append(templateMap[ty], result...)
				lastId = int(result[len(result)-1].ID)
				for _, item := range result {
					dataMap[ty] = append(dataMap[ty], item.Data)
				}
			}
			hasMore = len(result) == common.BKMaxInstanceLimit
		}
	}

	for ty, initor := range typeHandlerMap {
		if err := initor(kit, db, dataMap[ty]); err != nil {
			blog.Errorf("add template data failed for type %s, err: %v, rid: %s", ty, err, kit.Rid)
			return err
		}
	}

	if err := addSvrCategory(kit, db); err != nil {
		blog.Errorf("add template data failed for type service_category, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	if err := insertUniqueKey(kit, db); err != nil {
		blog.Errorf("add template data failed for type unique keys, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

func addSvrCategory(kit *rest.Kit, db local.DB) error {

	filter := mapstr.MapStr{
		common.BKTenantTemplateTypeField: metadata.TemplateTypeServiceCategory,
	}
	result := make([]metadata.SvrCategoryTmp, 0)
	err := mongodb.Shard(kit.SysShardOpts()).Table(common.BKTableNameTenantTemplate).Find(filter).All(kit.Ctx, &result)
	if err != nil {
		blog.Errorf("get template data for types service_category failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	subCategory := make(map[string][]mapstr.MapStr, 0)
	parentCategory := make([]mapstr.MapStr, 0)
	for _, item := range result {
		if item.IsParent {
			item.Data[common.BKParentIDField] = 0
			parentCategory = append(parentCategory, item.Data)
		} else {
			subCategory[item.ParentName] = append(subCategory[item.ParentName], item.Data)
		}
	}

	insertOps := &utils.InsertOptions{
		UniqueFields: []string{common.BKFieldName, common.BKParentIDField, common.BKAppIDField},
		IgnoreKeys:   []string{common.BKFieldID, common.BKRootIDField},
		IDField:      []string{common.BKFieldID, common.BKRootIDField},
		AuditTypeField: &utils.AuditResType{
			AuditType:    metadata.PlatformSetting,
			ResourceType: metadata.ServiceCategoryRes,
		},
		AuditDataField: &utils.AuditDataField{
			BizIDField:   "bk_biz_id",
			ResIDField:   "id",
			ResNameField: "name",
		},
	}

	parentIDs, err := utils.InsertData(kit, db, common.BKTableNameServiceCategory, parentCategory, insertOps)
	if err != nil {
		blog.Errorf("insert service category data for table %s failed, err: %v", common.BKTableNameServiceCategory, err)
		return err
	}
	for key, value := range parentIDs {
		name := strings.Split(key, "*")[0]
		parentIDs[name] = value
	}

	var subInsertData []mapstr.MapStr
	for key := range subCategory {
		parentID, err := util.GetInt64ByInterface(parentIDs[key])
		if err != nil {
			blog.Errorf("get parent id int64 failed, err: %v", err)
			return err
		}
		for _, subValue := range subCategory[key] {
			subValue[common.BKParentIDField] = parentID
			subValue[common.BKRootIDField] = parentID
			subInsertData = append(subInsertData, subValue)
		}
	}

	insertOps.UniqueFields = []string{common.BKFieldID}
	insertOps.IDField = []string{common.BKFieldID}
	_, err = utils.InsertData(kit, db, common.BKTableNameServiceCategory, subInsertData, insertOps)
	if err != nil {
		blog.Errorf("insert service category data for table %s failed, err: %v", common.BKTableNameServiceCategory, err)
		return err
	}
	return nil
}
