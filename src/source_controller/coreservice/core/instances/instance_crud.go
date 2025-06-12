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

	"configcenter/pkg/inst/logics"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/common/valid"
	"configcenter/src/storage/driver/mongodb"
)

func (m *instanceManager) batchSave(kit *rest.Kit, objID string, params []mapstr.MapStr) ([]uint64, error) {
	instTableName, err := logics.GetObjInstTableFromCache(kit, m.clientSet, objID)
	if err != nil {
		blog.Errorf("get object(%s) instance table name failed, err: %v, rid: %s", objID, err, kit.Rid)
		return nil, err
	}
	instIDFieldName := common.GetInstIDField(objID)
	ids, err := getSequences(kit, instTableName, len(params))
	if err != nil {
		return nil, err
	}
	mappings := make([]mapstr.MapStr, 0)
	ts := time.Now()
	insertDefaultAreaHosts := make([]mapstr.MapStr, 0)
	for idx := range params {
		if objID == common.BKInnerObjIDHost {
			params[idx], err = metadata.ConvertHostSpecialStringToArray(params[idx])
			if err != nil {
				blog.Errorf("convert host special string to array failed, err: %v, rid: %s", err, kit.Rid)
				return nil, err
			}

			defaultAreaHost, isDefaultAreaStaticHost, err := m.validDefaultAreaHost(kit,
				params[idx], int64(ids[idx]))
			if err != nil {
				blog.Errorf("valid default area host failed, err: %v, rid: %s", err, kit.Rid)
				return nil, err
			}
			if isDefaultAreaStaticHost {
				insertDefaultAreaHosts = append(insertDefaultAreaHosts, defaultAreaHost)
			}
		}
		// build new object instance data.
		if !valid.IsInnerObject(objID) {
			params[idx][common.BKObjIDField] = objID
		}
		params[idx].Set(instIDFieldName, ids[idx])
		params[idx].Set(common.CreateTimeField, ts)
		params[idx].Set(common.LastTimeField, ts)
		params[idx].Set(common.BKCreatedBy, kit.User)
		params[idx].Set(common.BKCreatedAt, ts)
		params[idx].Set(common.BKUpdatedAt, ts)
		if !metadata.IsCommon(objID) {
			continue
		}
		// build new object mapping data for inner object instance.
		mapping := make(mapstr.MapStr, 0)
		mapping[instIDFieldName] = ids[idx]
		mapping[common.BKObjIDField] = objID
		mappings = append(mappings, mapping)
	}

	if len(mappings) != 0 {
		// save new object mappings data for inner object instance.
		if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameObjectBaseMapping).Insert(kit.Ctx,
			mappings); err != nil {
			return nil, err
		}
	}
	if len(insertDefaultAreaHosts) > 0 {
		err = mongodb.Shard(kit.SysShardOpts()).Table(common.BKTableNameDefaultAreaHost).Insert(kit.Ctx,
			insertDefaultAreaHosts)
		if err != nil {
			blog.Errorf("save default area host failed, err: %v, rid: %s", err, kit.Rid)
			if mongodb.IsDuplicatedError(err) {
				return nil, kit.CCError.CCErrorf(common.CCErrCommDuplicateItem)
			}
			return nil, err
		}
	}
	// save object instances.
	err = mongodb.Shard(kit.ShardOpts()).Table(instTableName).Insert(kit.Ctx, params)
	if err != nil {
		blog.Errorf("save instances failed, err: %v, objID: %s, instances: %v, rid: %s", err, objID, params,
			kit.Rid)
		if mongodb.IsDuplicatedError(err) {
			return nil, kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, mongodb.GetDuplicateKey(err))
		}
		return nil, err
	}
	return ids, nil
}

func (m *instanceManager) save(kit *rest.Kit, objID string, inputParam mapstr.MapStr) (uint64, error) {
	instTableName, err := logics.GetObjInstTableFromCache(kit, m.clientSet, objID)
	if err != nil {
		blog.Errorf("get object(%s) instance table name failed, err: %v, rid: %s", objID, err, kit.Rid)
		return 0, err
	}
	ids, err := getSequences(kit, instTableName, 1)
	if err != nil {
		return 0, err
	}

	if objID == common.BKInnerObjIDHost {
		var err error
		inputParam, err = metadata.ConvertHostSpecialStringToArray(inputParam)
		if err != nil {
			return 0, err
		}
		insertDefaultAreaHost, isDefaultAreaStaticHost, err := m.validDefaultAreaHost(kit, inputParam, int64(ids[0]))
		if err != nil {
			blog.Errorf("valid default area host failed, err: %v, rid: %s", err, kit.Rid)
			return 0, err
		}

		if isDefaultAreaStaticHost {
			err = mongodb.Shard(kit.SysShardOpts()).Table(common.BKTableNameDefaultAreaHost).Insert(kit.Ctx,
				insertDefaultAreaHost)
			if err != nil {
				blog.Errorf("save default area host failed, err: %v, rid: %s", err, kit.Rid)
				if mongodb.IsDuplicatedError(err) {
					return 0, kit.CCError.CCErrorf(common.CCErrCommDuplicateItem)
				}
				return 0, err
			}
		}
	}

	// build new object instance data.
	instIDFieldName := common.GetInstIDField(objID)
	inputParam[instIDFieldName] = ids[0]
	if !valid.IsInnerObject(objID) {
		inputParam[common.BKObjIDField] = objID
	}
	ts := time.Now()
	inputParam.Set(common.CreateTimeField, ts)
	inputParam.Set(common.LastTimeField, ts)

	inputParam.Set(common.BKCreatedBy, kit.User)
	inputParam.Set(common.BKCreatedAt, ts)
	inputParam.Set(common.BKUpdatedAt, ts)

	// build and save new object mapping data for inner object instance.
	if metadata.IsCommon(objID) {
		mapping := make(mapstr.MapStr, 0)
		mapping[instIDFieldName] = ids[0]
		mapping[common.BKObjIDField] = objID

		// save instance object type mapping.
		if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameObjectBaseMapping).Insert(kit.Ctx,
			mapping); err != nil {
			return 0, err
		}
	}

	// save object instance.
	err = mongodb.Shard(kit.ShardOpts()).Table(instTableName).Insert(kit.Ctx, inputParam)
	if err != nil {
		blog.Errorf("save instance error. err: %v, objID: %s, instance: %+v, rid: %s", err, objID, inputParam,
			kit.Rid)

		if mongodb.IsDuplicatedError(err) {
			return ids[0], kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, mongodb.GetDuplicateKey(err))
		}
		return 0, err
	}

	return ids[0], nil
}

func getSequences(kit *rest.Kit, table string, count int) ([]uint64, error) {
	if count <= 0 {
		return nil, kit.CCError.CCError(common.CCErrCommHTTPInputInvalid)
	}

	ids, err := mongodb.Shard(kit.SysShardOpts()).NextSequences(kit.Ctx, table, count)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (m *instanceManager) update(kit *rest.Kit, objID string, data mapstr.MapStr, cond mapstr.MapStr) errors.CCError {
	if objID == common.BKInnerObjIDHost {
		var err error
		data, err = metadata.ConvertHostSpecialStringToArray(data)
		if err != nil {
			return err
		}
	}

	tableName, err := logics.GetObjInstTableFromCache(kit, m.clientSet, objID)
	if err != nil {
		blog.Errorf("get object(%s) instance table name failed, err: %v, rid: %s", objID, err, kit.Rid)
		return err
	}
	if !valid.IsInnerObject(objID) {
		cond.Set(common.BKObjIDField, objID)
	}
	ts := time.Now()
	data.Set(common.LastTimeField, ts)
	data.Set(common.BKUpdatedBy, kit.User)
	data.Set(common.BKUpdatedAt, ts)

	data.Remove(common.BKObjIDField)
	err = mongodb.Shard(kit.ShardOpts()).Table(tableName).Update(kit.Ctx, cond, data)
	if err != nil {
		blog.Errorf("update instance error. err: %v, objID: %s, instance: %+v, cond: %+v, rid: %s", err, objID,
			data, cond, kit.Rid)
		if mongodb.IsDuplicatedError(err) {
			return kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, mongodb.GetDuplicateKey(err))
		}
		return kit.CCError.Error(common.CCErrCommDBUpdateFailed)
	}
	return nil
}

func (m *instanceManager) getInsts(kit *rest.Kit, objID string, cond mapstr.MapStr) (origins []mapstr.MapStr,
	exists bool, err error) {

	origins = make([]mapstr.MapStr, 0)
	tableName, err := logics.GetObjInstTableFromCache(kit, m.clientSet, objID)
	if err != nil {
		blog.Errorf("get object(%s) instance table name failed, err: %v, rid: %s", objID, err, kit.Rid)
		return origins, !mongodb.IsNotFoundError(err), err
	}
	if !valid.IsInnerObject(objID) {
		cond.Set(common.BKObjIDField, objID)
	}
	if objID == common.BKInnerObjIDHost {
		hosts := make([]metadata.HostMapStr, 0)
		err = mongodb.Shard(kit.ShardOpts()).Table(tableName).Find(cond).All(kit.Ctx, &hosts)
		for _, host := range hosts {
			origins = append(origins, mapstr.MapStr(host))
		}
	} else {
		err = mongodb.Shard(kit.ShardOpts()).Table(tableName).Find(cond).All(kit.Ctx, &origins)
	}
	return origins, !mongodb.IsNotFoundError(err), err
}

func (m *instanceManager) getInstDataByID(kit *rest.Kit, objID string, instID int64) (origin mapstr.MapStr,
	err error) {

	tableName, err := logics.GetObjInstTableFromCache(kit, m.clientSet, objID)
	if err != nil {
		blog.Errorf("get object(%s) instance table name failed, err: %v, rid: %s", objID, err, kit.Rid)
		return nil, err
	}
	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: common.GetInstIDField(objID), Val: instID})

	if common.IsObjectInstShardingTable(tableName) {
		cond.Element(&mongo.Eq{Key: common.BKObjIDField, Val: objID})
	}

	if objID == common.BKInnerObjIDHost {
		host := make(metadata.HostMapStr)
		err = mongodb.Shard(kit.ShardOpts()).Table(tableName).Find(cond.ToMapStr()).One(kit.Ctx, &host)
		origin = mapstr.MapStr(host)
	} else {
		err = mongodb.Shard(kit.ShardOpts()).Table(tableName).Find(cond.ToMapStr()).One(kit.Ctx, &origin)
	}
	if err != nil {
		return nil, err
	}
	return origin, nil
}

func (m *instanceManager) countInstance(kit *rest.Kit, objID string, cond mapstr.MapStr) (count uint64, err error) {
	tableName, err := logics.GetObjInstTableFromCache(kit, m.clientSet, objID)
	if err != nil {
		blog.Errorf("get object(%s) instance table name failed, err: %v, rid: %s", objID, err, kit.Rid)
		return 0, err
	}
	if cond == nil {
		cond = make(map[string]interface{})
	}

	if common.IsObjectInstShardingTable(tableName) {
		objIDCond, ok := cond[common.BKObjIDField]
		if ok && objIDCond != objID {
			blog.V(9).Infof("countInstance condition's bk_obj_id: %s not match objID: %s, rid: %s", objIDCond, objID,
				kit.Rid)
			return 0, nil
		}
		cond[common.BKObjIDField] = objID
	}

	count, err = mongodb.Shard(kit.ShardOpts()).Table(tableName).Find(cond).Count(kit.Ctx)

	return count, err
}

// validDefaultAreaHost valid the default area host, ip is not allowed to be duplicated
func (m *instanceManager) validDefaultAreaHost(kit *rest.Kit, instanceData mapstr.MapStr, hostID int64) (mapstr.MapStr,
	bool, error) {

	insertData := make(mapstr.MapStr)
	needDeal, err := logics.IsDefaultAreaStaticHost(instanceData)
	if err != nil {
		return insertData, false, err
	}
	if !needDeal {
		return insertData, false, nil
	}

	ip, isIPExist := instanceData[common.BKHostInnerIPField]

	ipv6, isIPV6Exist := instanceData[common.BKHostInnerIPv6Field]

	if !isIPExist && !isIPV6Exist {
		blog.Errorf("invalid default area host, ip and ipv6 is not exist, rid: %s", kit.Rid)
		return nil, false, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "bk_host_innerip")
	}

	return mapstr.MapStr{
		common.BKHostIDField:        hostID,
		common.BKCloudIDField:       common.BKDefaultDirSubArea,
		common.BKHostInnerIPField:   ip,
		common.BKHostInnerIPv6Field: ipv6,
		common.TenantID:             kit.TenantID,
	}, true, nil
}
