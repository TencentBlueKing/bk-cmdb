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

package y3_9_202008121631

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func addModelFieldIshide(ctx context.Context, db dal.RDB, _ *upgrader.Config) (err error) {
	modelObjIDs := []string{common.BKInnerObjIDProc, common.BKInnerObjIDPlat}

	// 在cc_ObjDes表中将指定模型实例增加ishide=true字段
	for _, objID := range modelObjIDs {
		cond := mapstr.MapStr{
			common.BKObjIDField: objID,
		}
		doc := mapstr.MapStr{
			"ishide": true,
		}

		if err := db.Table(common.BKTableNameObjDes).Update(ctx, cond, doc); err != nil {
			blog.ErrorJSON("failed to model add field, objID: %s, err: %s", err, objID)
			return err
		}
	}

	return nil
}
