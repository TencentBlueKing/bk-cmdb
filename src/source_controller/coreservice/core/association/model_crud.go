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

package association

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql"
	"configcenter/src/storage/driver/mongodb"
)

func (m *associationModel) count(kit *rest.Kit, cond universalsql.Condition) (cnt uint64, err error) {

	cnt, err = mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameObjAsst).Find(cond.ToMapStr()).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("count object association failed, err: %v, cond: %v, rid: %s", err, cond, kit.Rid)
		return 0, err
	}
	return cnt, err
}

func (m *associationModel) isExists(kit *rest.Kit, cond universalsql.Condition) (oneResult *metadata.Association,
	exists bool, err error) {

	oneResult = &metadata.Association{}
	err = mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameObjAsst).Find(cond.ToMapStr()).One(kit.Ctx, oneResult)
	if err != nil && !mongodb.IsNotFoundError(err) {
		blog.Errorf("find object association failed, err: %v, cond: %v, rid: %s", err, cond, kit.Rid)
		return oneResult, false, kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	return oneResult, !mongodb.IsNotFoundError(err), nil
}

func (m *associationModel) save(kit *rest.Kit, assoParam *metadata.Association) (id uint64, err error) {

	id, err = mongodb.Shard(kit.SysShardOpts()).NextSequence(kit.Ctx, common.BKTableNameObjAsst)
	if err != nil {
		blog.Errorf("get sequence ID failed, err: %v, cond: %v, rid: %s", err, kit.Rid)
		return id, err
	}

	assoParam.ID = int64(id)
	err = mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameObjAsst).Insert(kit.Ctx, assoParam)
	if err != nil {
		blog.Errorf("insert object association failed, err: %v, param: %v, rid: %s", err, assoParam, kit.Rid)
		return 0, err
	}
	return id, err
}

func (m *associationModel) update(kit *rest.Kit, data mapstr.MapStr, cond universalsql.Condition) (cnt uint64,
	err error) {

	cnt, err = mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameObjAsst).UpdateMany(kit.Ctx,
		cond.ToMapStr(), data)
	if err != nil {
		blog.Errorf("update object association failed, err: %v, cond: %v, rid: %s", err, cond, kit.Rid)
		return 0, err
	}
	return cnt, err
}

func (m *associationModel) delete(kit *rest.Kit, cond universalsql.Condition) (cnt uint64, err error) {

	cnt, err = mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameObjAsst).DeleteMany(kit.Ctx, cond.ToMapStr())
	if err != nil {
		blog.Errorf("delete object association failed, err: %v, cond: %v, rid: %s", err, cond, kit.Rid)
		return 0, err
	}
	return cnt, err
}

func (m *associationModel) search(kit *rest.Kit, cond universalsql.Condition) ([]metadata.Association, error) {

	dataResult := []metadata.Association{}
	err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameObjAsst).Find(cond.ToMapStr()).All(kit.Ctx,
		&dataResult)
	if err != nil {
		blog.Errorf("find object association failed, err: %v, cond: %v, rid: %s", err, cond, kit.Rid)
		return dataResult, err
	}
	return dataResult, err
}

func (m *associationModel) searchReturnMapStr(kit *rest.Kit, cond universalsql.Condition) ([]mapstr.MapStr, error) {
	dataResult := []mapstr.MapStr{}
	err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameObjAsst).Find(cond.ToMapStr()).All(kit.Ctx,
		&dataResult)
	if err != nil {
		blog.Errorf("find object association failed, err: %v, cond: %v, rid: %s", err, cond, kit.Rid)
		return dataResult, err
	}
	return dataResult, err
}
