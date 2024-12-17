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

var cloudAreaData = []CloudArea{
	{
		CloudID:   common.BKDefaultDirSubArea,
		CloudName: "Default Area",
		Status:    "1",
		Default:   int64(common.BuiltIn),
	},
	{
		CloudID:   common.UnassignedCloudAreaID,
		CloudName: common.UnassignedCloudAreaName,
		Default:   int64(common.BuiltIn),
	},
}

func addCloudAreaData(kit *rest.Kit, db dal.Dal) error {
	cloudData := make([]interface{}, 0)
	for _, data := range cloudAreaData {
		data.Time = tools.NewTime()
		cloudData = append(cloudData, data)
	}

	needField := &tools.InsertOptions{
		UniqueFields: []string{common.BKCloudNameField},
		IgnoreKeys:   []string{common.BKCloudIDField},
		IDField:      []string{common.BKCloudIDField},
		AuditTypeField: &tools.AuditResType{
			AuditType:    metadata.ModelType,
			ResourceType: metadata.ModuleRes,
		},
		AuditDataField: &tools.AuditDataField{
			ResIDField:   common.BKCloudIDField,
			ResNameField: "bk_cloud_name",
		},
	}

	_, err := tools.InsertData(kit, db.Shard(kit.ShardOpts()), common.BKTableNameBasePlat, cloudData, needField)
	if err != nil {
		blog.Errorf("insert data for table %s failed, err: %v", common.BKTableNameBasePlat, err)
		return err
	}
	return nil
}

type CloudArea struct {
	CloudID     int64  ` bson:"bk_cloud_id"`
	CloudName   string ` bson:"bk_cloud_name"`
	Status      string ` bson:"bk_status"`
	CloudVendor string ` bson:"bk_cloud_vendor"`
	VpcID       string ` bson:"bk_vpc_id"`
	VpcName     string ` bson:"bk_vpc_name"`
	Region      string ` bson:"bk_region"`
	*tools.Time `bson:",inline"`
	Default     int64 ` bson:"default"`
}
