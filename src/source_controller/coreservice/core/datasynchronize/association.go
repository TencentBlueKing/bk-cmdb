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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/driver/mongodb"
)

type association struct {
	base     *synchronizeAdapter
	dataType metadata.SynchronizeOperateDataType
	// instance data classify
	DataClassify string
}

func NewSynchronizeAssociationAdapter(s *metadata.SynchronizeParameter) dataTypeInterface {

	return &association{
		base:         newSynchronizeAdapter(s),
		dataType:     s.OperateDataType,
		DataClassify: s.DataClassify,
	}
}

func (a *association) SaveSynchronize(kit *rest.Kit) errors.CCError {

	// Each model is written separately for subsequent expansion,
	// each type may be processed differently.
	switch a.base.syncData.DataClassify {
	case common.SynchronizeAssociationTypeModelHost:
		return a.saveSynchronizeAssociationModuleHostConfig(kit)
	default:
		return kit.CCError.Errorf(common.CCErrCoreServiceSyncDataClassifyNotExistError, a.dataType, a.DataClassify)
	}

}

func (a *association) PreSynchronizeFilter(kit *rest.Kit) errors.CCError {
	err := a.preSynchronizeFilterBefore(kit)
	if err != nil {
		return err
	}
	return a.base.PreSynchronizeFilter(kit)
}

func (a *association) GetErrorStringArr(kit *rest.Kit) ([]metadata.ExceptionResult, errors.CCError) {

	if len(a.base.errorArray) == 0 {
		return nil, nil
	}
	err := kit.CCError.Error(common.CCErrCoreServiceSyncError)
	switch a.base.syncData.DataClassify {
	case common.SynchronizeAssociationTypeModelHost:
		var errArr []metadata.ExceptionResult
		for _, err := range a.base.errorArray {
			errMsg := fmt.Sprintf("module and host relation error. info:%#v error:%s", err.instInfo.Info, err.err.Error())
			errArr = append(errArr, metadata.ExceptionResult{
				OriginIndex: err.idx,
				Message:     errMsg,
			})
		}
		return errArr, err
	default:
		return a.base.GetErrorStringArr(kit)
	}
}

// saveSynchronizeAssociationModuleHostConfig
// Host and module relationship is special, need special implementation
func (a *association) saveSynchronizeAssociationModuleHostConfig(kit *rest.Kit) errors.CCError {
	tableName := common.BKTableNameModuleHostConfig
	for _, item := range a.base.syncData.InfoArray {

		//  branch clone not support deep copy
		// not change value
		newItem := item.Info.Clone()

		newItem.Remove(common.MetadataField)
		cnt, err := mongodb.Client().Table(tableName).Find(newItem).Count(kit.Ctx)
		if err != nil {
			blog.Errorf("saveSynchronizeAssociationModuleHostConfig query db error,err:%s.DataSign:%s,condition:%#v,rid:%s", err.Error(), a.DataClassify, newItem, kit.Rid)
			a.base.errorArray[item.ID] = synchronizeAdapterError{
				instInfo: item,
				err:      kit.CCError.Error(common.CCErrCommDBSelectFailed),
			}
			continue
		}
		if cnt == 0 {
			err := mongodb.Client().Table(tableName).Insert(kit.Ctx, item.Info)
			if err != nil {
				blog.Errorf("saveSynchronizeAssociationModuleHostConfig save data to db error,err:%s.DataSign:%s,info:%#v,rid:%s", err.Error(), a.DataClassify, item, kit.Rid)
				a.base.errorArray[item.ID] = synchronizeAdapterError{
					instInfo: item,
					err:      kit.CCError.Error(common.CCErrCommDBInsertFailed),
				}
				continue
			}
		} else {
			err := mongodb.Client().Table(tableName).Update(kit.Ctx, newItem, item.Info)
			if err != nil {
				blog.Errorf("saveSynchronizeAssociationModuleHostConfig update data to db error,err:%s.DataSign:%s,info:%#v,rid:%s", err.Error(), a.DataClassify, item, kit.Rid)
				a.base.errorArray[item.ID] = synchronizeAdapterError{
					instInfo: item,
					err:      kit.CCError.Error(common.CCErrCommDBUpdateFailed),
				}
				continue
			}
		}
	}
	return nil
}

func (a *association) preSynchronizeFilterBefore(kit *rest.Kit) errors.CCError {
	switch a.base.syncData.DataClassify {
	case common.SynchronizeAssociationTypeModelHost:
		for idx, item := range a.base.syncData.InfoArray {
			// cc_ModuleHostConfig not id field.
			item.ID = int64(idx)
		}
	default:

	}
	return nil
}
func (a *association) preSynchronizeFilterEnd(kit *rest.Kit) errors.CCError {
	return nil
}
