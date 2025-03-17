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

package data

import (
	"strconv"
	"strings"

	tenanttmp "configcenter/pkg/types/tenant-template"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/admin_server/upgrader/tools"
	"configcenter/src/storage/dal/mongo/local"
)

var (
	defaultServiceCategoryID int64
	parentCategory           = []string{
		"数据库",
		"消息队列",
		"HTTP 服务",
		"存储",
		"Default",
	}
	subCategoryMap = map[string][]string{
		"数据库":     {"Mysql", "Redis", "Oracle", "SQLServer", "MongoDB", "Etcd", "Zookeeper"},
		"消息队列":    {"Kafka", "RabbitMQ"},
		"HTTP 服务": {"Nginx", "Apache", "Tomcat"},
		"存储":      {"Ceph", "NFS"},
		"Default": {"Default"},
	}
)

func addServiceCategoryData(kit *rest.Kit, db local.DB) error {
	parentServiceCategory := make([]mapstr.MapStr, 0)
	tmpData := make([]tenanttmp.SvrCategoryTmp, 0)
	for _, value := range parentCategory {
		category := ServiceCategory{
			Name:      value,
			IsBuiltIn: true,
		}
		tmpData = append(tmpData, tenanttmp.SvrCategoryTmp{
			Name:       value,
			ParentName: "",
		})
		item, err := tools.ConvStructToMap(category)
		if err != nil {
			blog.Errorf("convert struct to map failed, err: %v", err)
			return err
		}
		parentServiceCategory = append(parentServiceCategory, item)
	}

	// add parent category data
	needField := &tools.InsertOptions{
		UniqueFields: []string{common.BKFieldName, common.BKParentIDField, common.BKAppIDField},
		IgnoreKeys:   []string{common.BKFieldID, common.BKRootIDField},
		IDField:      []string{common.BKFieldID, common.BKRootIDField},
		AuditTypeField: &tools.AuditResType{
			AuditType:    metadata.PlatformSetting,
			ResourceType: metadata.ServiceCategoryRes,
		},
		AuditDataField: &tools.AuditDataField{
			BizIDField:   "bk_biz_id",
			ResIDField:   "id",
			ResNameField: "name",
		},
	}

	parentIDs, err := tools.InsertData(kit, db, common.BKTableNameServiceCategory, parentServiceCategory, needField)
	if err != nil {
		blog.Errorf("insert service category data for table %s failed, err: %v", common.BKTableNameServiceCategory, err)
		return err
	}

	fieldMap := map[string]int64{
		common.BKFieldName:     0,
		common.BKParentIDField: 1,
		common.BKAppIDField:    2,
	}
	existParentIDs := make(map[string]interface{}, 0)
	for key, value := range parentIDs {
		name := strings.Split(key, "*")[fieldMap[common.BKFieldName]]
		parentID, err := strconv.ParseInt(strings.Split(key, "*")[fieldMap[common.BKParentIDField]], 10, 64)
		if err != nil {
			blog.Errorf("convert interface to int64 failed, err: %v", err)
			return err
		}

		if parentID == 0 {
			existParentIDs[name] = value
		}
	}

	svrTmpData := make([]tenanttmp.TenantTmpData[tenanttmp.SvrCategoryTmp], 0)
	for _, item := range tmpData {
		svrTmpData = append(svrTmpData, tenanttmp.TenantTmpData[tenanttmp.SvrCategoryTmp]{
			Type:  tenanttmp.TemplateTypeServiceCategory,
			IsPre: true,
			Data:  item,
		})
	}
	err = tools.InsertSvrCategoryTmp(kit, db, svrTmpData)
	if err != nil {
		blog.Errorf("insert template data failed, err: %v", err)
		return err
	}
	if err = addSubSrvCategoryData(kit, db, existParentIDs); err != nil {
		blog.Errorf("add sub service category data failed, err: %v", err)
		return err
	}

	return nil
}

func addSubSrvCategoryData(kit *rest.Kit, db local.DB, parentIDs map[string]interface{}) error {
	// add sub category data
	subCategoryData := make([]mapstr.MapStr, 0)
	tmpData := make([]tenanttmp.SvrCategoryTmp, 0)
	for key, value := range subCategoryMap {
		parentID, err := util.GetInt64ByInterface(parentIDs[key])
		if err != nil {
			blog.Errorf("get parent id int64 failed, err: %v", err)
			return err
		}
		for _, subValue := range value {
			category := ServiceCategory{
				Name:      subValue,
				RootID:    parentID,
				ParentID:  parentID,
				IsBuiltIn: true,
				BizID:     0,
			}
			item, err := tools.ConvStructToMap(category)
			if err != nil {
				blog.Errorf("convert struct to map failed, err: %v", err)
				return err
			}
			subCategoryData = append(subCategoryData, item)
			tmpData = append(tmpData, tenanttmp.SvrCategoryTmp{
				Name:       subValue,
				ParentName: key,
			})
		}
	}

	needField := &tools.InsertOptions{
		UniqueFields: []string{common.BKFieldName, common.BKParentIDField, common.BKAppIDField},
		IgnoreKeys:   []string{common.BKFieldID},
		IDField:      []string{common.BKFieldID},
		AuditTypeField: &tools.AuditResType{
			AuditType:    metadata.PlatformSetting,
			ResourceType: metadata.ServiceCategoryRes,
		},
		AuditDataField: &tools.AuditDataField{
			BizIDField:   "bk_biz_id",
			ResIDField:   "id",
			ResNameField: "name",
		},
	}

	subIds, err := tools.InsertData(kit, db, common.BKTableNameServiceCategory, subCategoryData, needField)
	if err != nil {
		blog.Errorf("insert service category data for table %s failed, err: %v", common.BKTableNameServiceCategory, err)
		return err
	}

	svrTmpData := make([]tenanttmp.TenantTmpData[tenanttmp.SvrCategoryTmp], 0)
	for _, item := range tmpData {
		svrTmpData = append(svrTmpData, tenanttmp.TenantTmpData[tenanttmp.SvrCategoryTmp]{
			Type:  tenanttmp.TemplateTypeServiceCategory,
			IsPre: true,
			Data:  item,
		})
	}
	err = tools.InsertSvrCategoryTmp(kit, db, svrTmpData)
	if err != nil {
		blog.Errorf("insert template data failed, err: %v", err)
		return err
	}

	parentStr := "Default"
	defaultParentID, err := util.GetInt64ByInterface(parentIDs[parentStr])
	if err != nil {
		blog.Errorf("get default parent id int64 failed, err: %v", err)
		return err
	}
	uniqueFields := []string{"Default", strconv.FormatInt(defaultParentID, 10), "0"}
	subUniqueStr := strings.Join(uniqueFields, "*")
	// get default service category id
	defaultServiceCategoryID, err = util.GetInt64ByInterface(subIds[subUniqueStr])
	if err != nil {
		blog.Errorf("get default service category id failed, err: %v", err)
		return err
	}
	blog.Infof("add service category data success, default id: %d", defaultServiceCategoryID)

	return nil
}

type ServiceCategory struct {
	ID        int64  `bson:"id"`
	Name      string `bson:"name"`
	RootID    int64  `bson:"bk_root_id"`
	ParentID  int64  `bson:"bk_parent_id"`
	IsBuiltIn bool   `bson:"is_built_in"`
	BizID     int64  `bson:"bk_biz_id"`
}
