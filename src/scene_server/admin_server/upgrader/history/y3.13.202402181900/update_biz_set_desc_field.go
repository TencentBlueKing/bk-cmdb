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

package y3_13_202402181900

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/storage/dal"
)

const (
	// bizSetID 相关业务集：全业务-蓝盾测试部署专用
	bizSetID = 9992001
	oldName  = "description"
	newName  = "bk_biz_set_desc"
)

func updateBizSetDescField(ctx context.Context, db dal.RDB) error {
	cond := map[string]interface{}{
		common.BKBizSetIDField: bizSetID,
		newName:                mapstr.MapStr{"$exists": true},
	}
	count, err := db.Table(common.BKTableNameBaseBizSet).Find(cond).Count(ctx)
	if err != nil {
		blog.Errorf("find business set failed, filter: %v, err: %v", cond, err)
		return err
	}
	if count != 0 {
		return nil
	}

	filter := map[string]interface{}{
		common.BKBizSetIDField: bizSetID,
	}
	if err := db.Table(common.BKTableNameBaseBizSet).RenameColumn(ctx, filter, oldName, newName); err != nil {
		blog.Errorf("rename biz set description field failed, filter: %v, err: %v", filter, err)
		return err
	}
	return nil
}
