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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader/tools"
	"configcenter/src/storage/dal"
)

var parentCategory = []string{
	"数据库",
	"消息队列",
	"HTTP 服务",
	"存储",
	"Default",
}

var subCategoryMap = map[string][]string{
	"数据库":     {"Mysql", "Redis", "Oracle", "SQLServer", "MongoDB", "Etcd", "Zookeeper"},
	"消息队列":    {"Kafka", "RabbitMQ"},
	"HTTP 服务": {"Nginx", "Apache", "Tomcat"},
	"存储":      {"Ceph", "NFS"},
	"Default": {"Default"},
}

var defaultID int

func addServiceCategoryData(kit *rest.Kit, db dal.Dal) error {
	var parentServiceCategory []interface{}
	rootID := 0
	parentNameMap := map[string]int{}
	var serviceCategory []tools.AuditField
	for _, value := range parentCategory {
		rootID++
		category := ServiceCategory{
			Name:      value,
			IsBuiltIn: true,
			RootID:    int64(rootID),
		}
		parentNameMap[category.Name] = int(category.RootID)
		parentServiceCategory = append(parentServiceCategory, category)
		serviceCategory = append(serviceCategory, tools.AuditField{
			AuditType:    metadata.PlatformSetting,
			ResourceType: metadata.ServiceCategoryRes,
			ResourceName: category.Name,
		})
	}
	// add parent category data

	cmpField := &tools.CmpFiled{
		Unique:     []string{common.BKFieldName, common.BKParentIDField, common.BKAppIDField},
		IgnoreKeys: make([]string, 0),
		ID:         common.BKFieldID,
	}
	_, err := tools.InsertData(kit, db.Shard(kit.ShardOpts()), common.BKTableNameServiceCategory,
		parentServiceCategory, cmpField, serviceCategory)
	if err != nil {
		blog.Errorf("insert service category data for table %s failed, err: %v", common.BKTableNameServiceCategory, err)
		return err
	}

	// add sub category data
	var subCategoryData []interface{}
	var subServiceCategory []tools.AuditField
	countPos := 0
	for key, value := range subCategoryMap {
		for _, subValue := range value {
			category := ServiceCategory{
				Name:      subValue,
				RootID:    int64(parentNameMap[key]),
				ParentID:  int64(parentNameMap[key]),
				IsBuiltIn: true,
				BizID:     0,
			}
			subCategoryData = append(subCategoryData, category)
			subServiceCategory = append(subServiceCategory, tools.AuditField{
				AuditType:    metadata.PlatformSetting,
				ResourceType: metadata.ServiceCategoryRes,
				ResourceName: category.Name,
			})
			countPos++
		}
		if key == "Default" {
			defaultID = countPos
		}
	}

	cmpField = &tools.CmpFiled{
		Unique:     []string{common.BKFieldName, common.BKParentIDField, common.BKAppIDField},
		IgnoreKeys: make([]string, 0),
		ID:         metadata.AttributeFieldID,
	}
	_, err = tools.InsertData(kit, db.Shard(kit.ShardOpts()), common.BKTableNameServiceCategory, subCategoryData,
		cmpField, subServiceCategory)
	if err != nil {
		blog.Errorf("insert service category data for table %s failed, err: %v", common.BKTableNameServiceCategory, err)
		return err
	}

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
