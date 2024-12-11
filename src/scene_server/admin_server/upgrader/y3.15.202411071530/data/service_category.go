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

var defaultID int64

func addServiceCategoryData(kit *rest.Kit, db dal.Dal) error {
	var parentServiceCategory []interface{}
	rootID := 0
	parentNameMap := map[string]int{}
	var serviceCategoryAudit []tools.AuditType
	for _, value := range parentCategory {
		rootID++
		category := ServiceCategory{
			Name:      value,
			IsBuiltIn: true,
			RootID:    int64(rootID),
		}
		parentNameMap[category.Name] = int(category.RootID)
		parentServiceCategory = append(parentServiceCategory, category)
		serviceCategoryAudit = append(serviceCategoryAudit, tools.AuditType{
			AuditType:    metadata.PlatformSetting,
			ResourceType: metadata.ServiceCategoryRes,
		})
	}
	// add parent category data

	cmpField := &tools.CmpFiled{
		UniqueFields: []string{common.BKFieldName, common.BKParentIDField, common.BKAppIDField},
		IgnoreKeys:   make([]string, 0),
		IDField:      common.BKFieldID,
	}
	auditDataField := &tools.AuditDataField{
		BusinessID:   "bk_biz_id",
		ResourceID:   "id",
		ResourceName: "name",
	}
	parentIDs, err := tools.InsertData(kit, db.Shard(kit.ShardOpts()), common.BKTableNameServiceCategory,
		parentServiceCategory, cmpField, serviceCategoryAudit, auditDataField)
	if err != nil {
		blog.Errorf("insert service category data for table %s failed, err: %v", common.BKTableNameServiceCategory, err)
		return err
	}

	// add sub category data
	var subCategoryData []interface{}
	var subServiceCategory []tools.AuditType
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
			subServiceCategory = append(subServiceCategory, tools.AuditType{
				AuditType:    metadata.PlatformSetting,
				ResourceType: metadata.ServiceCategoryRes,
			})
		}
	}

	cmpField = &tools.CmpFiled{
		UniqueFields: []string{common.BKFieldName, common.BKParentIDField, common.BKAppIDField},
		IgnoreKeys:   []string{metadata.AttributeFieldID},
		IDField:      metadata.AttributeFieldID,
	}
	subAuditDataField := &tools.AuditDataField{
		BusinessID:   "bk_biz_id",
		ResourceID:   "id",
		ResourceName: "name",
	}
	ids, err := tools.InsertData(kit, db.Shard(kit.ShardOpts()), common.BKTableNameServiceCategory, subCategoryData,
		cmpField, subServiceCategory, subAuditDataField)
	if err != nil {
		blog.Errorf("insert service category data for table %s failed, err: %v", common.BKTableNameServiceCategory, err)
		return err
	}

	uniqueFields := []string{"Default", "0", "0"}
	parentStr := strings.Join(uniqueFields, "*")
	defaultParentID := parentIDs[parentStr].(uint64)
	uniqueFields = []string{"Default", strconv.FormatUint(defaultParentID, 10), "0"}
	subUniqueStr := strings.Join(uniqueFields, "*")
	defaultID, err = util.GetInt64ByInterface(ids[subUniqueStr])
	if err != nil {
		blog.Errorf("get default service category id failed, err: %v", err)
		return err
	}
	blog.Infof("add service category data success, default id: %d", defaultID)

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
