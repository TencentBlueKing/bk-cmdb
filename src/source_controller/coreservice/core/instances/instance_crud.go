/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package instances

import (
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/driver/mongodb/instancemapping"
)

func (m *instanceManager) save(kit *rest.Kit, objID string, inputParam mapstr.MapStr) (uint64, error) {
	if objID == common.BKInnerObjIDHost {
		inputParam = metadata.ConvertHostSpecialStringToArray(inputParam)
	}

	instTableName := common.GetInstTableName(objID, kit.SupplierAccount)
	id, err := mongodb.Client().NextSequence(kit.Ctx, instTableName)
	if err != nil {
		return 0, err
	}

	// build new object instance data.
	instIDFieldName := common.GetInstIDField(objID)
	inputParam[instIDFieldName] = id
	if !util.IsInnerObject(objID) {
		inputParam[common.BKObjIDField] = objID
	}
	ts := time.Now()
	inputParam.Set(common.BKOwnerIDField, kit.SupplierAccount)
	inputParam.Set(common.CreateTimeField, ts)
	inputParam.Set(common.LastTimeField, ts)

	// build and save new object mapping data for inner object instance.
	if metadata.IsCommon(objID) {
		mapping := make(mapstr.MapStr, 0)
		mapping[instIDFieldName] = id
		mapping[common.BKObjIDField] = objID
		mapping[common.BkSupplierAccount] = kit.SupplierAccount

		// save instance object type mapping.
		if err := instancemapping.Create(kit.Ctx, mapping); err != nil {
			return 0, err
		}
	}

	// save object instance.
	err = mongodb.Client().Table(instTableName).Insert(kit.Ctx, inputParam)
	if err != nil {
		blog.ErrorJSON("save instance error. err: %s, objID: %s, instance: %s, rid: %s",
			err.Error(), objID, inputParam, kit.Rid)
		if mongodb.Client().IsDuplicatedError(err) {
			return id, kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, mongodb.GetDuplicateKey(err))
		}
		return 0, err
	}

	return id, nil
}

func (m *instanceManager) update(kit *rest.Kit, objID string, data mapstr.MapStr, cond mapstr.MapStr) errors.CCError {
	if objID == common.BKInnerObjIDHost {
		data = metadata.ConvertHostSpecialStringToArray(data)
	}
	tableName := common.GetInstTableName(objID, kit.SupplierAccount)
	if !util.IsInnerObject(objID) {
		cond.Set(common.BKObjIDField, objID)
	}
	ts := time.Now()
	data.Set(common.LastTimeField, ts)
	data.Remove(common.BKObjIDField)
	err := mongodb.Client().Table(tableName).Update(kit.Ctx, cond, data)
	if err != nil {
		blog.ErrorJSON("update instance error. err: %s, objID: %s, instance: %s, cond: %s, rid: %s",
			err.Error(), objID, data, cond, kit.Rid)
		if mongodb.Client().IsDuplicatedError(err) {
			return kit.CCError.CCError(common.CCErrCommDuplicateItem)
		}
		return kit.CCError.Error(common.CCErrCommDBUpdateFailed)
	}
	return nil
}

func (m *instanceManager) getInsts(kit *rest.Kit, objID string, cond mapstr.MapStr) (origins []mapstr.MapStr, exists bool, err error) {
	origins = make([]mapstr.MapStr, 0)
	tableName := common.GetInstTableName(objID, kit.SupplierAccount)
	if !util.IsInnerObject(objID) {
		cond.Set(common.BKObjIDField, objID)
	}
	if objID == common.BKInnerObjIDHost {
		hosts := make([]metadata.HostMapStr, 0)
		err = mongodb.Client().Table(tableName).Find(cond).All(kit.Ctx, &hosts)
		for _, host := range hosts {
			origins = append(origins, mapstr.MapStr(host))
		}
	} else {
		err = mongodb.Client().Table(tableName).Find(cond).All(kit.Ctx, &origins)
	}
	return origins, !mongodb.Client().IsNotFoundError(err), err
}

func (m *instanceManager) getInstDataByID(kit *rest.Kit, objID string, instID int64) (origin mapstr.MapStr, err error) {
	tableName := common.GetInstTableName(objID, kit.SupplierAccount)

	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: common.GetInstIDField(objID), Val: instID})

	if common.IsObjectInstShardingTable(common.GetInstTableName(objID, kit.SupplierAccount)) {
		cond.Element(&mongo.Eq{Key: common.BKObjIDField, Val: objID})
	}

	if objID == common.BKInnerObjIDHost {
		host := make(metadata.HostMapStr)
		err = mongodb.Client().Table(tableName).Find(cond.ToMapStr()).One(kit.Ctx, &host)
		origin = mapstr.MapStr(host)
	} else {
		err = mongodb.Client().Table(tableName).Find(cond.ToMapStr()).One(kit.Ctx, &origin)
	}
	if nil != err {
		return nil, err
	}
	return origin, nil
}

func (m *instanceManager) countInstance(kit *rest.Kit, objID string, cond mapstr.MapStr) (count uint64, err error) {
	tableName := common.GetInstTableName(objID, kit.SupplierAccount)

	if common.IsObjectInstShardingTable(tableName) {
		objIDCond, ok := cond[common.BKObjIDField]
		if ok && objIDCond != objID {
			blog.V(9).Infof("countInstance condition's bk_obj_id: %s not match objID: %s, rid: %s", objIDCond, objID, kit.Rid)
			return 0, nil
		}
		cond[common.BKObjIDField] = objID
	}

	cond = util.SetQueryOwner(cond, kit.SupplierAccount)
	count, err = mongodb.Client().Table(tableName).Find(cond).Count(kit.Ctx)

	return count, err
}
