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

package synchronize

import (
	"encoding/json"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/mongodb"
)

// parseDataArr parse data array to actual type
func parseDataArr[T any](kit *rest.Kit, data any) ([]T, error) {
	var results []T

	switch val := data.(type) {
	case []T:
		results = val
	case []json.RawMessage:
		arr, err := util.DecodeRawJsonArr[T](val)
		if err != nil {
			blog.Errorf("decode raw json array(%+v) failed, err: %v, rid: %s", val, err, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "data")
		}
		results = arr
	}

	return results, nil
}

// validateDependency validate if dependency exists
func validateDependency[T comparable](kit *rest.Kit, table, idField string, ids []T) error {
	if len(ids) == 0 {
		return nil
	}
	ids = util.ArrUnique(ids)

	cond := mapstr.MapStr{
		idField: mapstr.MapStr{
			common.BKDBIN: ids,
		},
	}

	count, err := mongodb.Client().Table(table).Find(cond).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("count %s dependency failed, err: %v, cond: %+v, rid: %s", table, err, cond, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	if int(count) != len(ids) {
		blog.Errorf("not all %s dependency exists, count: %d, ids: %+v, rid: %s", table, count, ids, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommNotFound)
	}

	return nil
}

// getMainlineParentObj get mainline parent object, returns mainline parent object id and if it is a mainline object
func getMainlineParentObj(kit *rest.Kit, objID string) (string, bool, error) {
	cond := map[string]interface{}{
		common.AssociationKindIDField: common.AssociationKindMainline,
		common.BKObjIDField:           objID,
	}

	asst := new(metadata.Association)
	err := mongodb.Client().Table(common.BKTableNameObjAsst).Find(cond).Fields(common.BKAsstObjIDField).
		One(kit.Ctx, &asst)
	if err != nil {
		if mongodb.Client().IsNotFoundError(err) {
			return "", false, nil
		}

		blog.Errorf("get %s mainline parent object failed, err: %v, cond: %+v, rid: %s", objID, err, cond, kit.Rid)
		return "", false, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	return asst.AsstObjID, true, nil
}
