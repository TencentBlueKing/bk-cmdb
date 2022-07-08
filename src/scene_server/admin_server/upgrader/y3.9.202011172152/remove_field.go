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

package y3_9_202011172152

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

// removeField remove redundant bk_childid and bk_parentid fields
func removeField(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	toBeRemovedFields := []string{"bk_childid", "bk_parentid"}

	// get all mainline objects
	mainlineFilter := map[string]interface{}{common.AssociationKindIDField: common.AssociationKindMainline}
	associations := make([]metadata.Association, 0)
	err := db.Table(bkTableNameObjAsst).Find(mainlineFilter).All(ctx, &associations)
	if nil != err {
		blog.Errorf("search mainline associations failed, err: %v", err)
		return err
	}

	objectIDs := make([]string, 0)
	for _, asst := range associations {
		objectIDs = append(objectIDs, asst.ObjectID, asst.AsstObjID)
	}
	objectIDs = util.StrArrayUnique(objectIDs)

	// remove all mainline objects' bk_childid and bk_parentid attributes
	attrFilter := map[string]interface{}{
		common.BKObjIDField:      map[string]interface{}{common.BKDBIN: objectIDs},
		common.BKPropertyIDField: map[string]interface{}{common.BKDBIN: toBeRemovedFields},
	}
	if err := db.Table(bkTableNameObjAttDes).Delete(ctx, attrFilter); err != nil {
		blog.Errorf("remove mainline objects' attributes(%+v) failed, err: %v", toBeRemovedFields, err)
		return err
	}

	// remove all mainline objects' instances' bk_childid and bk_parentid fields
	for _, objectID := range objectIDs {
		instFilter := map[string]interface{}{
			common.BKDBOR: []map[string]interface{}{{
				"bk_childid": map[string]interface{}{common.BKDBExists: true},
			}, {
				"bk_parentid": map[string]interface{}{common.BKDBExists: true},
			}},
		}
		if !common.IsInnerModel(objectID) {
			instFilter[common.BKObjIDField] = objectID
		}
		tableName := GetInstTableName(objectID)
		instIDField := common.GetInstIDField(objectID)

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
				blog.Errorf("remove mainline objects' instances fields(%+v) failed, err: %v", toBeRemovedFields, err)
				return err
			}

			if len(insts) < common.BKMaxPageSize {
				break
			}
		}
	}
	return nil
}

// GetInstTableName returns inst data table name
func GetInstTableName(objID string) string {
	switch objID {
	case bkInnerObjIDApp:
		return bkTableNameBaseApp
	case bkInnerObjIDSet:
		return bkTableNameBaseSet
	case bkInnerObjIDModule:
		return bkTableNameBaseModule
	case bkInnerObjIDHost:
		return bkTableNameBaseHost
	case bkInnerObjIDProc:
		return bkTableNameBaseProcess
	case bkInnerObjIDPlat:
		return bkTableNameBasePlat
	default:
		return bkTableNameBaseInst
	}
}

const (
	// bkTableNameInstAsst the table name of the inst association
	bkTableNameInstAsst = "cc_InstAsst"

	bkTableNameBaseApp     = "cc_ApplicationBase"
	bkTableNameBaseHost    = "cc_HostBase"
	bkTableNameBaseModule  = "cc_ModuleBase"
	bkTableNameBaseInst    = "cc_ObjectBase"
	bkTableNameBasePlat    = "cc_PlatBase"
	bkTableNameBaseSet     = "cc_SetBase"
	bkTableNameBaseProcess = "cc_Process"
	bkTableNameObjAsst     = "cc_ObjAsst"
	bkTableNameObjAttDes   = "cc_ObjAttDes"
)

const (
	// bkInnerObjIDApp the inner object
	bkInnerObjIDApp = "biz"

	// bkInnerObjIDSet the inner object
	bkInnerObjIDSet = "set"

	// bkInnerObjIDModule the inner object
	bkInnerObjIDModule = "module"

	// bkInnerObjIDHost the inner object
	bkInnerObjIDHost = "host"

	// bkInnerObjIDObject the inner object
	bkInnerObjIDObject = "object"

	// bkInnerObjIDProc the inner object
	bkInnerObjIDProc = "process"

	// bkInnerObjIDPlat the inner object
	bkInnerObjIDPlat = "plat"
)
