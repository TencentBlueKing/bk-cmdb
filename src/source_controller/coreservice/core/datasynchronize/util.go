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
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

type synchronizeAdapterError struct {
	err      errors.CCError
	instInfo *metadata.SynchronizeItem
	idx      int64
}

type synchronizeAdapterDBParameter struct {
	tableName   string
	InstIDField string
}

type synchronizeAdapter struct {
	dbProxy    dal.RDB
	syncData   *metadata.SynchronizeParameter
	errorArray map[int64]synchronizeAdapterError
}

type dataTypeInterface interface {
	PreSynchronizeFilter(ctx core.ContextParams) errors.CCError
	GetErrorStringArr(ctx core.ContextParams) ([]metadata.ExceptionResult, errors.CCError)
	SaveSynchronize(ctx core.ContextParams) errors.CCError
}

func newSynchronizeAdapter(syncData *metadata.SynchronizeParameter, dbProxy dal.RDB) *synchronizeAdapter {
	return &synchronizeAdapter{
		syncData:   syncData,
		dbProxy:    dbProxy,
		errorArray: make(map[int64]synchronizeAdapterError, 0),
	}
}

func (s *synchronizeAdapter) PreSynchronizeFilter(ctx core.ContextParams) errors.CCError {
	if s.syncData.SynchronizeFlag == "" {
		// TODO  return error not synchronize sign
		return ctx.Error.Errorf(common.CCErrCommParamsNeedSet, "synchronize_flag")
	}
	if s.syncData.InfoArray == nil {
		// TODO return error not found synchroize data
		return ctx.Error.Errorf(common.CCErrCommParamsNeedSet, "instance_info_array")
	}
	var syncDataArr []*metadata.SynchronizeItem
	for _, item := range s.syncData.InfoArray {
		if !item.Info.IsEmpty() {
			syncDataArr = append(syncDataArr, item)
		}
	}
	s.syncData.InfoArray = syncDataArr
	// synchronize data need to write data,append synchronize sign to metada
	if s.syncData.OperateType != metadata.SynchronizeOperateTypeUpdate {
		// set synchroize sign to instance metadata
		for _, item := range s.syncData.InfoArray {
			if item.Info.Exists(common.MetadataField) {
				metadata, err := item.Info.MapStr(common.MetadataField)
				if err != nil {
					blog.Errorf("preSynchronizeFilter get %s field error, inst info:%#v,rid:%s", common.MetadataField, item, ctx.ReqID)
					s.errorArray[item.ID] = synchronizeAdapterError{
						instInfo: item,
						err:      ctx.Error.Errorf(common.CCErrCommInstFieldConvFail, s.syncData.DataClassify, common.MetadataField, "mapstr", err.Error()),
					}
					continue
				}
				metadata.Set(common.MetaDataSynchronizeFlagField, s.syncData.SynchronizeFlag)
				metadata.Set(common.MetaDataSynchronizeVersionField, s.syncData.Version)
			} else {
				item.Info.Set(common.MetadataField,
					mapstr.MapStr{
						common.MetaDataSynchronizeFlagField:    s.syncData.SynchronizeFlag,
						common.MetaDataSynchronizeVersionField: s.syncData.Version,
					})
			}
		}
	}

	return nil
}

func (s *synchronizeAdapter) GetErrorStringArr(ctx core.ContextParams) ([]metadata.ExceptionResult, errors.CCError) {
	if len(s.errorArray) == 0 {
		return nil, nil
	}
	var errArr []metadata.ExceptionResult
	for _, err := range s.errorArray {
		errMsg := fmt.Sprintf("[%s] instID:[%d] error:%s", s.syncData.DataClassify, err.instInfo.ID, err.err.Error())
		errArr = append(errArr, metadata.ExceptionResult{
			OriginIndex: err.instInfo.ID,
			Message:     errMsg,
		})
	}
	return errArr, ctx.Error.Error(common.CCErrCoreServiceSyncError)
}

func (s *synchronizeAdapter) saveSynchronize(ctx core.ContextParams, dbParam synchronizeAdapterDBParameter) {
	switch s.syncData.OperateType {
	case metadata.SynchronizeOperateTypeDelete:
		s.deleteSynchronize(ctx, dbParam)
	case metadata.SynchronizeOperateTypeUpdate, metadata.SynchronizeOperateTypeAdd, metadata.SynchronizeOperateTypeRepalce:
		s.replaceSynchronize(ctx, dbParam)

	}
}

func (s *synchronizeAdapter) replaceSynchronize(ctx core.ContextParams, dbParam synchronizeAdapterDBParameter) {
	for _, item := range s.syncData.InfoArray {
		_, ok := s.errorArray[item.ID]
		if ok {
			continue
		}
		conds := mapstr.MapStr{dbParam.InstIDField: item.ID}
		exist, err := s.existSynchronizeID(ctx, dbParam.tableName, conds)
		if err != nil {
			blog.Errorf("replaceSynchronize existSynchronizeID error.DataClassify:%s,info:%#v,rid:%s", s.syncData.DataClassify, item, ctx.ReqID)
			s.errorArray[item.ID] = synchronizeAdapterError{
				instInfo: item,
				err:      err,
			}
			continue
		}
		if exist {
			err := s.dbProxy.Table(dbParam.tableName).Update(ctx, conds, item.Info)
			if err != nil {
				blog.Errorf("replaceSynchronize update info error,err:%s.DataClassify:%s,condition:%#v,info:%#v,rid:%s", err.Error(), s.syncData.DataClassify, conds, item, ctx.ReqID)
				s.errorArray[item.ID] = synchronizeAdapterError{
					instInfo: item,
					err:      ctx.Error.Error(common.CCErrCommDBUpdateFailed),
				}
				continue
			}
		} else {
			err := s.dbProxy.Table(dbParam.tableName).Insert(ctx, item.Info)
			if err != nil {
				blog.Errorf("replaceSynchronize insert info error,err:%s.DataClassify:%s,info:%#v,rid:%s", err.Error(), s.syncData.DataClassify, item, ctx.ReqID)
				s.errorArray[item.ID] = synchronizeAdapterError{
					instInfo: item,
					err:      ctx.Error.Error(common.CCErrCommDBInsertFailed),
				}
				continue
			}
		}
	}
}

func (s *synchronizeAdapter) deleteSynchronize(ctx core.ContextParams, dbParam synchronizeAdapterDBParameter) {
	var instIDArr []int64
	for _, item := range s.syncData.InfoArray {
		instIDArr = append(instIDArr, item.ID)
	}
	err := s.dbProxy.Table(dbParam.tableName).Delete(ctx, mapstr.MapStr{dbParam.InstIDField: mapstr.MapStr{common.BKDBIN: instIDArr}})
	if err != nil {
		blog.Errorf("deleteSynchronize delete info error,err:%s.DataClassify:%s,instIDArr:%#v,rid:%s", err.Error(), s.syncData.DataClassify, instIDArr, ctx.ReqID)
		for _, item := range s.syncData.InfoArray {
			s.errorArray[item.ID] = synchronizeAdapterError{
				instInfo: item,
				err:      ctx.Error.Error(common.CCErrCommDBDeleteFailed),
			}
		}
	}
}

func (s *synchronizeAdapter) existSynchronizeID(ctx core.ContextParams, tableName string, conds mapstr.MapStr) (bool, errors.CCError) {
	cnt, err := s.dbProxy.Table(tableName).Find(conds).Count(ctx)
	if err != nil {
		blog.Errorf("existSynchronizeID error. DataClassify:%s,conds:%#v,rid:%s", s.syncData.DataClassify, conds, ctx.ReqID)
		return false, ctx.Error.Error(common.CCErrCommDBSelectFailed)
	}
	if cnt > 0 {
		return true, nil
	}
	return false, nil

}
