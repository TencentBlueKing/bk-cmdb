/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package y3_10_202205182148

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

const (
	bizSetID = 9991001
	oldName  = "biz_set_maintainer"
	newName  = "bk_biz_maintainer"
)

func updateDefaultBusinessSetFieldName(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	cond := map[string]interface{}{
		common.BKBizSetIDField: bizSetID,
		newName:                mapstr.MapStr{"$exists": true},
	}

	count, err := db.Table(common.BKTableNameBaseBizSet).Find(cond).Count(ctx)
	if err != nil {
		blog.Errorf("find business set failed, filter: %v, err: %v", cond, err)
		return err
	}

	// if default business set have bk_biz_maintainer field, skip update biz_set_maintainer filed
	if count != 0 {
		return nil
	}

	filter := map[string]interface{}{
		common.BKBizSetIDField: bizSetID,
	}

	if err := db.Table(common.BKTableNameBaseBizSet).RenameColumn(ctx, filter, oldName, newName); err != nil {
		blog.Errorf("rename default business set column failed, filter: %v, err: %v", filter, err)
		return err
	}
	return nil
}
