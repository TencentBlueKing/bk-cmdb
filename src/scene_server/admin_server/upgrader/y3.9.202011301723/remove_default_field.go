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

package y3_9_202011301723

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

// removeDefaultField remove custom layer instance default field
func removeDefaultField(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	toBeRemovedFields := []string{common.BKDefaultField}

	// get all mainline objects
	mainlineFilter := map[string]interface{}{common.AssociationKindIDField: common.AssociationKindMainline}
	associations := make([]metadata.Association, 0)
	err := db.Table(common.BKTableNameObjAsst).Find(mainlineFilter).All(ctx, &associations)
	if nil != err {
		blog.Errorf("search mainline associations failed, err: %v", err)
		return err
	}

	objectIDs := make([]string, 0)
	for _, asst := range associations {
		if !common.IsInnerModel(asst.ObjectID) {
			objectIDs = append(objectIDs, asst.ObjectID)
		}

		if !common.IsInnerModel(asst.AsstObjID) {
			objectIDs = append(objectIDs, asst.AsstObjID)
		}
	}
	objectIDs = util.StrArrayUnique(objectIDs)

	// remove all custom layer instance default fields
	tableName := common.BKTableNameBaseInst
	instIDField := common.BKInstIDField
	for _, objectID := range objectIDs {
		instFilter := map[string]interface{}{
			common.BKDefaultField: map[string]interface{}{common.BKDBExists: true},
			common.BKObjIDField:   objectID,
		}

		for {
			insts := make([]map[string]interface{}, 0)
			err := db.Table(tableName).Find(instFilter).Start(0).Limit(common.BKMaxPageSize).
				Fields(instIDField).All(ctx, &insts)
			if err != nil {
				blog.Errorf("get instance ids to remove fields failed, err: %v", err)
				return err
			}

			if len(insts) == 0 {
				break
			}

			instIDs := make([]int64, len(insts))
			for index, inst := range insts {
				instID, err := util.GetInt64ByInterface(inst[instIDField])
				if err != nil {
					blog.Errorf("get instance id failed, inst: %+v, err: %v", inst, err)
					return err
				}
				instIDs[index] = instID
			}

			instFilter := map[string]interface{}{
				instIDField: map[string]interface{}{
					common.BKDBIN: instIDs,
				},
			}

			if err := db.Table(tableName).DropColumns(ctx, instFilter, toBeRemovedFields); err != nil {
				blog.Errorf("remove custom layer instance fields(%+v) failed, err: %v", toBeRemovedFields, err)
				return err
			}

			if len(insts) < common.BKMaxPageSize {
				break
			}
		}
	}
	return nil
}
