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
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader/tools"
	"configcenter/src/storage/dal/mongo/local"
)

func addSelfIncrIDData(kit *rest.Kit, db local.DB) error {

	objIDs := []string{"host", "set", "module", "bk_project", "biz", "process", "plat", "bk_biz_set_obj"}
	ids := make([]string, 0)
	for _, object := range objIDs {
		ids = append(ids, metadata.GetIDRule(object))
	}
	ids = append(ids, metadata.GetIDRule(common.GlobalIDRule))

	needAddIDs := make([]mapstr.MapStr, 0)
	curTime := time.Now()
	for _, id := range ids {
		addID := mapstr.MapStr{
			common.BKFieldDBID:     id,
			common.BKFieldSeqID:    0,
			common.CreateTimeField: curTime,
			common.LastTimeField:   curTime,
		}
		needAddIDs = append(needAddIDs, addID)
	}

	needField := &tools.InsertOptions{
		UniqueFields: []string{"_id"},
		IgnoreKeys:   []string{"_id"},
	}

	_, err := tools.InsertData(kit, db, common.BKTableNameIDgenerator, needAddIDs, needField)
	if err != nil {
		blog.Errorf("insert data for table %s failed, err: %v", common.BKTableNameIDgenerator, err)
		return err
	}

	return nil
}
