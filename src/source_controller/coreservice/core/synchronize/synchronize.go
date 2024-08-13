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

// Package synchronize defines multiple cmdb synchronize logics
package synchronize

import (
	"configcenter/pkg/synchronize/types"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	errutil "configcenter/src/common/util/errors"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/driver/mongodb"
)

var _ core.SynchronizeOperation = (*synchronizer)(nil)

type synchronizer struct {
	synchronizerMap map[types.ResType]SyncI
}

// New create a new synchronizer instance
func New() core.SynchronizeOperation {
	s := &synchronizer{
		synchronizerMap: map[types.ResType]SyncI{
			types.Biz:            new(bizSyncer),
			types.Set:            new(setSyncer),
			types.Module:         new(moduleSyncer),
			types.Host:           new(hostSyncer),
			types.HostRelation:   new(hostRelSyncer),
			types.ObjectInstance: new(objInstSyncer),
			types.InstAsst:       new(instAsstSyncer),
		},
	}
	return s
}

// SyncI is the cmdb synchronize logics interface
type SyncI interface {
	ParseDataArr(kit *rest.Kit, data any) (any, error)
	Validate(kit *rest.Kit, subRes string, data any) error
	TableName(subRes, supplierAccount string) string
}

func (s *synchronizer) getSyncer(res types.ResType) (SyncI, errors.RawErrorInfo) {
	syncer, exists := s.synchronizerMap[res]
	if !exists {
		return nil, errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{types.ResTypeField},
		}
	}

	return syncer, errors.RawErrorInfo{}
}

// CreateData create cmdb synchronize data
func (s *synchronizer) CreateData(kit *rest.Kit, opt *types.CreateSyncDataOption) error {
	syncer, rawErr := s.getSyncer(opt.ResourceType)
	if rawErr.ErrCode != 0 {
		return rawErr.ToCCError(kit.CCError)
	}

	data, err := syncer.ParseDataArr(kit, opt.Data)
	if err != nil {
		return err
	}

	err = syncer.Validate(kit, opt.SubResource, data)
	if err != nil {
		return err
	}

	table := syncer.TableName(opt.SubResource, kit.SupplierAccount)
	err = mongodb.Client().Table(table).Insert(kit.Ctx, data)
	if err != nil {
		blog.Errorf("create sync data failed, err: %v, data: %+v, rid: %s", err, data, kit.Rid)
		return errutil.ConvDBInsertError(kit, mongodb.Client(), err)
	}

	return nil
}
