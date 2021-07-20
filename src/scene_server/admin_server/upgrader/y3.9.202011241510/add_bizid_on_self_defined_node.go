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

package y3_9_202011241510

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

// 补全自定义层级的bk_biz_id字段
func addBizIDOnSelfDefinedNode(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	// 获取主线模型里的bk_asst_id->bk_obj_id的映射关系
	filter := map[string]interface{}{
		"bk_asst_id": "bk_mainline",
	}
	asstObjRes := make([]map[string]string, 0)
	err := db.Table(common.BKTableNameObjAsst).Find(filter).Fields("bk_asst_obj_id", "bk_obj_id").All(ctx, &asstObjRes)
	if err != nil {
		blog.Errorf("find asstObjRes failed, filter:%#v, err:%v", filter, err)
		return err
	}
	asstMap := make(map[string]string)
	for _, r := range asstObjRes {
		asstMap[r["bk_asst_obj_id"]] = r["bk_obj_id"]
	}

	// 获取biz到set之间从上到下的所有自定义层级模型列表objArr
	objArr := make([]string, 0)
	nodeID := "biz"
	for asstMap[nodeID] != "set" {
		objArr = append(objArr, asstMap[nodeID])
		nodeID = asstMap[nodeID]
	}

	// 查询需要补充的bk_biz_id的自定义层级实例数量
	filter = map[string]interface{}{
		"bk_obj_id": map[string]interface{}{
			"$in": objArr,
		},
		"bk_biz_id": map[string]interface{}{
			"$exists": false,
		},
	}
	cnt, err := db.Table(common.BKTableNameBaseInst).Find(filter).Count(ctx)
	if err != nil {
		blog.Errorf("find count of obj without bizID failed, filter:%#v, err:%v", filter, err)
		return err
	}
	// 如果没有需要补充bk_biz_id字段的实例则退出
	if cnt == 0 {
		return nil
	}

	// 补全cc_ObjectBase表里的自定义层级的bk_biz_id字段
	// 补全业务下的第一层级的bk_biz_id
	filter = map[string]interface{}{
		"bk_obj_id": objArr[0],
		"bk_biz_id": map[string]interface{}{
			"$exists": false,
		},
	}
	instParentRes := make([]map[string]int64, 0)
	err = db.Table(common.BKTableNameBaseInst).Find(filter).Fields("bk_inst_id", "bk_parent_id").All(ctx, &instParentRes)
	if err != nil {
		blog.Errorf("find instParentRes failed, filter:%#v, err:%v", filter, err)
		return err
	}

	doc := make(map[string]interface{})
	for _, record := range instParentRes {
		filter = map[string]interface{}{
			"bk_inst_id": record["bk_inst_id"],
		}
		doc = map[string]interface{}{
			"bk_biz_id": record["bk_parent_id"],
		}
		err = db.Table(common.BKTableNameBaseInst).Update(ctx, filter, doc)
		if err != nil {
			blog.Errorf("update first level node failed, filter:%#v, doc:%#v, err:%v", filter, doc, err)
			return err
		}
	}

	// 补全其他自定义层级模型实例的bk_biz_id
	for i := 1; i < len(objArr); i++ {
		// 获取bk_inst_id->bk_parent_id的映射关系
		filter = map[string]interface{}{
			"bk_obj_id": objArr[i],
			"bk_biz_id": map[string]interface{}{
				"$exists": false,
			},
		}
		instParentRes := make([]map[string]int64, 0)
		err = db.Table(common.BKTableNameBaseInst).Find(filter).Fields("bk_inst_id", "bk_parent_id").All(ctx, &instParentRes)
		if err != nil {
			blog.Errorf("find instParentRes failed, filter:%#v, err:%v", filter, err)
			return err
		}
		instParentMap := make(map[int64]int64)
		parentIDs := make([]int64, 0)
		for _, r := range instParentRes {
			instParentMap[r["bk_inst_id"]] = r["bk_parent_id"]
			parentIDs = append(parentIDs, r["bk_parent_id"])
		}
		parentIDs = util.IntArrayUnique(parentIDs)

		// 获取bk_parent_id->bk_biz_id的映射关系
		filter = map[string]interface{}{
			"bk_obj_id": objArr[i-1],
			"bk_inst_id": map[string]interface{}{
				"$in": parentIDs,
			},
		}
		parentBizRes := make([]map[string]int64, 0)
		err = db.Table(common.BKTableNameBaseInst).Find(filter).Fields("bk_inst_id", "bk_biz_id").All(ctx, &parentBizRes)
		if err != nil {
			blog.Errorf("find parentBizRes failed, filter:%#v, err:%v", filter, err)
			return err
		}
		parentBizMap := make(map[int64]int64)
		for _, r := range parentBizRes {
			parentBizMap[r["bk_inst_id"]] = r["bk_biz_id"]
		}

		// 补全自定义层级模型实例的bk_biz_id
		for instID, parentID := range instParentMap {
			bizID := parentBizMap[parentID]
			filter = map[string]interface{}{
				"bk_inst_id": instID,
			}
			doc = map[string]interface{}{
				"bk_biz_id": bizID,
			}
			err = db.Table(common.BKTableNameBaseInst).Update(ctx, filter, doc)
			if err != nil {
				blog.Errorf("update other level node failed, filter:%#v, doc:%#v, err:%v", filter, doc, err)
				return err
			}
		}
	}

	return nil
}
