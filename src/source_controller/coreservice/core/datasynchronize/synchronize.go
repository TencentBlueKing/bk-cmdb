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

package datasynchronize

import (
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

type SynchronizeManager struct {
	dbProxy   dal.RDB
	dependent OperationDependences
}

// New create a new model manager instance
func New(dbProxy dal.RDB, dependent OperationDependences) core.DataSynchronizeOperation {
	return &SynchronizeManager{
		dbProxy:   dbProxy,
		dependent: dependent,
	}
}

func (s *SynchronizeManager) SynchronizeInstanceAdapter(ctx core.ContextParams, syncData *metadata.SynchronizeParameter) ([]metadata.ExceptionResult, error) {
	syncDataAdpater := NewSynchronizeInstanceAdapter(syncData, s.dbProxy)
	err := syncDataAdpater.PreSynchronizeFilter(ctx)
	if err != nil {
		return nil, err
	}
	syncDataAdpater.SaveSynchronize(ctx)
	return syncDataAdpater.GetErrorStringArr(ctx)

}

func (s *SynchronizeManager) SynchronizeModelAdapter(ctx core.ContextParams, syncData *metadata.SynchronizeParameter) ([]metadata.ExceptionResult, error) {
	syncDataAdpater := NewSynchronizeModelAdapter(syncData, s.dbProxy)
	err := syncDataAdpater.PreSynchronizeFilter(ctx)
	if err != nil {
		return nil, err
	}
	syncDataAdpater.SaveSynchronize(ctx)
	return syncDataAdpater.GetErrorStringArr(ctx)

}

func (s *SynchronizeManager) SynchronizeAssociationAdapter(ctx core.ContextParams, syncData *metadata.SynchronizeParameter) ([]metadata.ExceptionResult, error) {
	syncDataAdpater := NewSynchronizeAssociationAdapter(syncData, s.dbProxy)
	err := syncDataAdpater.PreSynchronizeFilter(ctx)
	if err != nil {
		return nil, err
	}
	syncDataAdpater.SaveSynchronize(ctx)
	return syncDataAdpater.GetErrorStringArr(ctx)

}

func (s *SynchronizeManager) GetAssociationInfo(ctx core.ContextParams, fetch *metadata.SynchronizeFetchInfoParameter) ([]mapstr.MapStr, uint64, error) {
	fetchAdapter := NewSynchronizeFetchAdapter(fetch, s.dbProxy)
	return fetchAdapter.Fetch(ctx)
}
