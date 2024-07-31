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

package fullsynccond

import (
	types "configcenter/pkg/cache/full-sync-cond"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/mongodb"
)

// CreateFullSyncCond create full sync cache condition
func (f *FullSyncCond) CreateFullSyncCond(kit *rest.Kit, opt *types.CreateFullSyncCondOpt) (int64, error) {
	if opt.IsAll {
		// check if one resource with sub resource has only one is_all=true condition
		cond := mapstr.MapStr{
			types.ResourceField: opt.Resource,
			types.IsAllField:    true,
		}
		cond = util.SetModOwner(cond, kit.SupplierAccount)
		if opt.SubResource != "" {
			cond[types.SubResField] = opt.SubResource
		}

		cnt, err := mongodb.Client().Table(types.BKTableNameFullSyncCond).Find(cond).Count(kit.Ctx)
		if err != nil {
			blog.Errorf("count is_all=true full sync cond failed, err: %v, cond: %+v, rid: %s", err, cond, kit.Rid)
			return 0, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
		}

		if cnt > 0 {
			blog.Errorf("is_all=true full sync cond exists, cannot create another one, opt: %+v, rid: %s", opt, kit.Rid)
			return 0, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, types.IsAllField)
		}
	} else {
		// check if all is_all=false conditions count <= types.NotAllCondLimit
		cond := mapstr.MapStr{
			types.IsAllField: false,
		}
		cnt, err := mongodb.Client().Table(types.BKTableNameFullSyncCond).Find(cond).Count(kit.Ctx)
		if err != nil {
			blog.Errorf("count is_all=true full sync cond failed, err: %v, cond: %+v, rid: %s", err, cond, kit.Rid)
			return 0, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
		}

		if cnt >= types.NotAllCondLimit {
			blog.Errorf("is_all=false full sync cond exceeds limit %d, cannot create another one, opt: %+v, rid: %s",
				types.NotAllCondLimit, opt, kit.Rid)
			return 0, kit.CCError.CCErrorf(common.CCErrCommXXExceedLimit, "not all conditions", types.NotAllCondLimit)
		}
	}

	id, err := mongodb.Client().NextSequence(kit.Ctx, types.BKTableNameFullSyncCond)
	if err != nil {
		blog.Errorf("generate full sync cond id failed, err: %v, rid: %s", err, kit.Rid)
		return 0, kit.CCError.CCErrorf(common.CCErrCommGenerateRecordIDFailed)
	}

	data := &types.FullSyncCond{
		ID:              int64(id),
		Resource:        opt.Resource,
		SubResource:     opt.SubResource,
		IsAll:           opt.IsAll,
		Interval:        opt.Interval,
		Condition:       opt.Condition,
		SupplierAccount: kit.SupplierAccount,
	}

	err = mongodb.Client().Table(types.BKTableNameFullSyncCond).Insert(kit.Ctx, data)
	if err != nil {
		blog.Errorf("insert full sync cond failed, err: %v, data: %+v, rid: %s", err, data, kit.Rid)
		return 0, kit.CCError.CCErrorf(common.CCErrCommDBInsertFailed)
	}

	return int64(id), nil
}

// UpdateFullSyncCond update full sync cache condition
func (f *FullSyncCond) UpdateFullSyncCond(kit *rest.Kit, opt *types.UpdateFullSyncCondOpt) error {
	cond := mapstr.MapStr{
		types.IDField: opt.ID,
	}
	cond = util.SetModOwner(cond, kit.SupplierAccount)

	data := mapstr.MapStr{
		types.IntervalField: opt.Data.Interval,
	}

	err := mongodb.Client().Table(types.BKTableNameFullSyncCond).Update(kit.Ctx, cond, data)
	if err != nil {
		blog.Errorf("update full sync cond failed, err: %v, cond: %+v, data: %+v, rid: %s", err, cond, data, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBUpdateFailed)
	}

	return nil
}

// DeleteFullSyncCond delete full sync cache condition
func (f *FullSyncCond) DeleteFullSyncCond(kit *rest.Kit, opt *types.DeleteFullSyncCondOpt) error {
	delCond := mapstr.MapStr{
		types.IDField: opt.ID,
	}
	delCond = util.SetModOwner(delCond, kit.SupplierAccount)

	err := mongodb.Client().Table(types.BKTableNameFullSyncCond).Delete(kit.Ctx, delCond)
	if err != nil {
		blog.Errorf("delete full sync cond %d failed, err: %v, rid: %s", opt.ID, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBDeleteFailed)
	}

	return nil
}

// ListFullSyncCond list full sync cache condition
func (f *FullSyncCond) ListFullSyncCond(kit *rest.Kit, opt *types.ListFullSyncCondOpt) (*types.ListFullSyncCondRes,
	error) {

	listCond := make(mapstr.MapStr)

	if len(opt.Resource) > 0 {
		listCond[types.ResourceField] = opt.Resource
	}

	if (opt.SubResource) != "" {
		listCond[types.SubResField] = opt.SubResource
	}

	if len(opt.IDs) > 0 {
		listCond[types.IDField] = mapstr.MapStr{
			common.BKDBIN: opt.IDs,
		}
	}
	listCond = util.SetQueryOwner(listCond, kit.SupplierAccount)

	result := make([]types.FullSyncCond, 0)
	err := mongodb.Client().Table(types.BKTableNameFullSyncCond).Find(listCond).All(kit.Ctx, &result)
	if err != nil {
		blog.Errorf("list full sync cond failed, err: %v, cond: %+v, rid: %s", err, listCond, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	return &types.ListFullSyncCondRes{Info: result}, nil
}

// GetFullSyncCond get full sync cache condition
func (f *FullSyncCond) GetFullSyncCond(kit *rest.Kit, id int64) (*types.FullSyncCond, error) {
	if id <= 0 {
		blog.Errorf("get full sync cond with invalid id, rid: %s", kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, types.IDField)
	}

	cond := mapstr.MapStr{
		types.IDField: id,
	}
	cond = util.SetQueryOwner(cond, kit.SupplierAccount)

	fullSyncCond := new(types.FullSyncCond)
	err := mongodb.Client().Table(types.BKTableNameFullSyncCond).Find(cond).One(kit.Ctx, &fullSyncCond)
	if err != nil {
		blog.Errorf("get full sync cond failed, err: %v, cond: %+v, rid: %s", err, cond, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	return fullSyncCond, nil
}
