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
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

type association struct {
	base     *synchronizeAdapter
	dataType metadata.SynchronizeDataType
	dbProxy  dal.RDB
	dataSign string
}

func NewSynchronizeAssociationAdapter(s *metadata.SynchronizeParameter, dbProxy dal.RDB) dataTypeInterface {

	return &association{
		base:     newSynchronizeAdapter(s, dbProxy),
		dataType: s.DataType,
		dataSign: s.DataSign,
		dbProxy:  dbProxy,
	}
}

func (a *association) SaveSynchronize(ctx core.ContextParams) errors.CCError {
	// Each model is written separately for subsequent expansion,
	// each type may be processed differently.
	switch a.base.syncData.DataSign {
	case common.SynchronizeAssociationTypeModelHost:
		return a.saveSynchronizeAssociationModuleHostConfig(ctx)
	default:
		return ctx.Error.Errorf(common.CCErrCoreServiceSyncDataSignNotExistError, a.dataType, a.dataSign)
	}
	return nil
}

func (a *association) PreSynchronizeFilter(ctx core.ContextParams) errors.CCError {
	err := a.preSynchronizeFilterBefore(ctx)
	if err != nil {
		return err
	}
	return a.base.PreSynchronizeFilter(ctx)
}

func (a *association) GetErrorStringArr(ctx core.ContextParams) ([]string, errors.CCError) {

	if len(a.base.errorArray) == 0 {
		return make([]string, 0), nil
	}
	err := ctx.Error.Error(common.CCErrCoreServiceSyncError)
	switch a.base.syncData.DataSign {
	case common.SynchronizeAssociationTypeModelHost:
		var errStrArr []string
		for _, err := range a.base.errorArray {
			errMsg := fmt.Sprintf("module and host relation error. info:%#v error:%s", err.instInfo.Info, err.err.Error())
			errStrArr = append(errStrArr, errMsg)
		}
		return errStrArr, err
	default:
		return a.base.getErrorStringArr(ctx), err
	}
}

// saveSynchronizeAssociationModuleHostConfig
// Host and module relationship is special, need special implementation
func (a *association) saveSynchronizeAssociationModuleHostConfig(ctx core.ContextParams) errors.CCError {
	for _, item := range a.base.syncData.InfoArray {
		// todo v3.3.x branch clone not support deep copy
		// not change value
		newItem, err := item.Info.Clone()
		if err != nil {
			blog.Errorf("saveSynchronizeAssociationModuleHostConfig copy data error. err:%s, inst info:%#v,rid:%s", err.Error(), item, ctx.ReqID)
			a.base.errorArray[item.ID] = synchronizeAdapterError{
				instInfo: item,
				err:      ctx.Error.Errorf(common.CCErrCommInstFieldConvFail, a.dataSign, "info", "copy mapstr", err.Error()),
			}
			continue
		}
		newItem.Remove(common.MetadataField)
		cnt, err := a.dbProxy.Table(common.BKTableNameProcInstanceModel).Find(newItem).Count(ctx)
		if err != nil {
			blog.Errorf("saveSynchronizeAssociationModuleHostConfig query db error,err:%s.DataSign:%s,condition:%#v,rid:%s", err.Error(), a.dataSign, newItem, ctx.ReqID)
			a.base.errorArray[item.ID] = synchronizeAdapterError{
				instInfo: item,
				err:      ctx.Error.Error(common.CCErrCommDBSelectFailed),
			}
			continue
		}
		if cnt == 0 {
			err := a.dbProxy.Table(common.BKTableNameProcInstanceModel).Insert(ctx, item)
			if err != nil {
				blog.Errorf("saveSynchronizeAssociationModuleHostConfig save data to db error,err:%s.DataSign:%s,info:%#v,rid:%s", err.Error(), a.dataSign, item, ctx.ReqID)
				a.base.errorArray[item.ID] = synchronizeAdapterError{
					instInfo: item,
					err:      ctx.Error.Error(common.CCErrCommDBInsertFailed),
				}
				continue
			}
		}
	}
	return nil
}

func (a *association) preSynchronizeFilterBefore(ctx core.ContextParams) errors.CCError {
	switch a.base.syncData.DataSign {
	case common.SynchronizeAssociationTypeModelHost:
		for idx, item := range a.base.syncData.InfoArray {
			// cc_ModuleHostConfig not id field.
			item.ID = int64(idx)
		}
	default:

	}
	return nil
}
func (a *association) preSynchronizeFilterEnd(ctx core.ContextParams) errors.CCError {
	return nil
}
