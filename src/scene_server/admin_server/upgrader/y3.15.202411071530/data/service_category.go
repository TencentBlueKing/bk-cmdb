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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/admin_server/upgrader/tools"
	"configcenter/src/storage/dal"
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

func addServiceCategoryData(kit *rest.Kit, db dal.Dal) error {
	parentServiceCategory := make([]interface{}, 0)
	for _, value := range parentCategory {
		category := ServiceCategory{
			Name:      value,
			IsBuiltIn: true,
		}
		parentServiceCategory = append(parentServiceCategory, category)
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

	parentIDs, err := tools.InsertData(kit, db.Shard(kit.ShardOpts()), common.BKTableNameServiceCategory,
		parentServiceCategory, needField)
	if err != nil {
		blog.Errorf("insert service category data for table %s failed, err: %v", common.BKTableNameServiceCategory, err)
		return err
	}
	for key, value := range parentIDs {
		name := strings.Split(key, "*")[0]
		parentIDs[name] = value
	}

	// add sub category data
	subCategoryData := make([]interface{}, 0)
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
			subCategoryData = append(subCategoryData, category)
		}
	}

	needField = &tools.InsertOptions{
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

	subIds, err := tools.InsertData(kit, db.Shard(kit.ShardOpts()), common.BKTableNameServiceCategory, subCategoryData,
		needField)
	if err != nil {
		blog.Errorf("insert service category data for table %s failed, err: %v", common.BKTableNameServiceCategory, err)
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