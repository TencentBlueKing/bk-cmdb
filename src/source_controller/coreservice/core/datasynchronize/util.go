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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
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
	isStr       bool
}

type synchronizeAdapter struct {
	dbProxy    dal.RDB
	syncData   *metadata.SynchronizeParameter
	errorArray map[int64]synchronizeAdapterError
}

type dataTypeInterface interface {
	PreSynchronizeFilter(kit *rest.Kit) errors.CCError
	GetErrorStringArr(kit *rest.Kit) ([]metadata.ExceptionResult, errors.CCError)
	SaveSynchronize(kit *rest.Kit) errors.CCError
}

func newSynchronizeAdapter(syncData *metadata.SynchronizeParameter, dbProxy dal.RDB) *synchronizeAdapter {
	return &synchronizeAdapter{
		syncData:   syncData,
		dbProxy:    dbProxy,
		errorArray: make(map[int64]synchronizeAdapterError, 0),
	}
}

func (s *synchronizeAdapter) PreSynchronizeFilter(kit *rest.Kit) errors.CCError {
	if s.syncData.SynchronizeFlag == "" {
		// TODO  return error not synchronize sign
		return kit.CCError.Errorf(common.CCErrCommParamsNeedSet, "synchronize_flag")
	}
	if s.syncData.InfoArray == nil {
		// TODO return error not found synchroize data
		return kit.CCError.Errorf(common.CCErrCommParamsNeedSet, "instance_info_array")
	}
	var syncDataArr []*metadata.SynchronizeItem
	for _, item := range s.syncData.InfoArray {
		if !item.Info.IsEmpty() {
			syncDataArr = append(syncDataArr, item)
		}
	}
	s.syncData.InfoArray = syncDataArr
	// synchronize data need to write data,append synchronize sign to metadata
	if s.syncData.OperateType != metadata.SynchronizeOperateTypeUpdate {
		// set synchronize sign to instance metadata
		for _, item := range s.syncData.InfoArray {
			if item.Info.Exists(common.MetadataField) {
				mData, err := item.Info.MapStr(common.MetadataField)
				if err != nil {
					blog.Errorf("preSynchronizeFilter get %s field error, inst info:%#v,rid:%s", common.MetadataField, item, kit.Rid)
					s.errorArray[item.ID] = synchronizeAdapterError{
						instInfo: item,
						err:      kit.CCError.CCErrorf(common.CCErrCommInstFieldConvertFail, s.syncData.DataClassify, common.MetadataField, "mapstr", err.Error()),
					}
					continue
				}
				mData.Set(common.MetaDataSynchronizeField, mapstr.MapStr{
					common.MetaDataSynchronizeFlagField:    s.syncData.SynchronizeFlag,
					common.MetaDataSynchronizeVersionField: s.syncData.Version,
				})
			} else {
				item.Info.Set(common.MetadataField,
					mapstr.MapStr{common.MetaDataSynchronizeField: mapstr.MapStr{
						common.MetaDataSynchronizeFlagField:    s.syncData.SynchronizeFlag,
						common.MetaDataSynchronizeVersionField: s.syncData.Version,
					}})
			}
		}
	}

	return nil
}

func (s *synchronizeAdapter) GetErrorStringArr(kit *rest.Kit) ([]metadata.ExceptionResult, errors.CCError) {
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
	return errArr, kit.CCError.Error(common.CCErrCoreServiceSyncError)
}

func (s *synchronizeAdapter) saveSynchronize(kit *rest.Kit, dbParam synchronizeAdapterDBParameter) {
	switch s.syncData.OperateType {
	case metadata.SynchronizeOperateTypeDelete:
		s.deleteSynchronize(kit, dbParam)
	case metadata.SynchronizeOperateTypeUpdate, metadata.SynchronizeOperateTypeAdd, metadata.SynchronizeOperateTypeRepalce:
		s.replaceSynchronize(kit, dbParam)

	}
}

func (s *synchronizeAdapter) replaceSynchronize(kit *rest.Kit, dbParam synchronizeAdapterDBParameter) {
	for _, item := range s.syncData.InfoArray {
		_, ok := s.errorArray[item.ID]
		if ok {
			continue
		}

		var conds mapstr.MapStr
		// can be combined
		mergeInstID, exist, err := s.getSameInfo(kit, dbParam.InstIDField, dbParam.tableName, item)
		if err != nil {
			blog.Errorf("replaceSynchronize getSameInfo error. err:%s, DataClassify:%s, info:%#v, rid:%s", err.Error(), s.syncData.DataClassify, item, kit.Rid)
			s.errorArray[item.ID] = synchronizeAdapterError{
				instInfo: item,
				err:      err,
			}
			continue
		}
		if exist {
			// The same data already exists, merging the existing data.
			conds = mapstr.MapStr{dbParam.InstIDField: mergeInstID}
		} else {
			exist, err = s.existSynchronizeID(kit, dbParam.tableName, mapstr.MapStr{dbParam.InstIDField: item.ID})
			if err != nil {
				blog.Errorf("replaceSynchronize existSynchronizeID error. err:%s, DataClassify:%s, info:%#v, exist:%v, rid:%s", err.Error(), s.syncData.DataClassify, item, exist, kit.Rid)
				s.errorArray[item.ID] = synchronizeAdapterError{
					instInfo: item,
					err:      err,
				}
				continue
			}
			if exist {
				conds = mapstr.MapStr{dbParam.InstIDField: item.ID}
			}
		}

		blog.V(6).Infof("replaceSynchronize DataClassify:%s, info:%#v, table:%s, version:%v, exist:%v, rid:%s", s.syncData.DataClassify, item, dbParam.tableName, s.syncData.Version, exist, kit.Rid)
		if exist {
			// Existing data, does not update the ID field
			delete(item.Info, dbParam.InstIDField)
			err := s.dbProxy.Table(dbParam.tableName).Update(kit.Ctx, conds, item.Info)
			if err != nil {
				blog.Errorf("replaceSynchronize update info error,err:%s.DataClassify:%s,condition:%#v,info:%#v,rid:%s", err.Error(), s.syncData.DataClassify, conds, item, kit.Rid)
				s.errorArray[item.ID] = synchronizeAdapterError{
					instInfo: item,
					err:      kit.CCError.Error(common.CCErrCommDBUpdateFailed),
				}
				continue
			}
		} else {
			err := s.dbProxy.Table(dbParam.tableName).Insert(kit.Ctx, item.Info)
			if err != nil {
				blog.Errorf("replaceSynchronize insert info error,err:%s.DataClassify:%s,info:%#v,rid:%s", err.Error(), s.syncData.DataClassify, item, kit.Rid)
				s.errorArray[item.ID] = synchronizeAdapterError{
					instInfo: item,
					err:      kit.CCError.Error(common.CCErrCommDBInsertFailed),
				}
				continue
			}
		}
	}
}

func (s *synchronizeAdapter) deleteSynchronize(kit *rest.Kit, dbParam synchronizeAdapterDBParameter) {
	var instIDArr []int64
	for _, item := range s.syncData.InfoArray {
		instIDArr = append(instIDArr, item.ID)
	}
	err := s.dbProxy.Table(dbParam.tableName).Delete(kit.Ctx, mapstr.MapStr{dbParam.InstIDField: mapstr.MapStr{common.BKDBIN: instIDArr}})
	if err != nil {
		blog.Errorf("deleteSynchronize delete info error,err:%s.DataClassify:%s,instIDArr:%#v,rid:%s", err.Error(), s.syncData.DataClassify, instIDArr, kit.Rid)
		for _, item := range s.syncData.InfoArray {
			s.errorArray[item.ID] = synchronizeAdapterError{
				instInfo: item,
				err:      kit.CCError.Error(common.CCErrCommDBDeleteFailed),
			}
		}
	}
}

func (s *synchronizeAdapter) existSynchronizeID(kit *rest.Kit, tableName string, conds mapstr.MapStr) (bool, errors.CCError) {
	cnt, err := s.dbProxy.Table(tableName).Find(conds).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("existSynchronizeID error. DataClassify:%s,conds:%#v,rid:%s", s.syncData.DataClassify, conds, kit.Rid)
		return false, kit.CCError.Error(common.CCErrCommDBSelectFailed)
	}
	if cnt > 0 {
		return true, nil
	}
	return false, nil

}

func (s *synchronizeAdapter) getSameInfo(kit *rest.Kit, instIDField, tableName string, info *metadata.SynchronizeItem) (int64, bool, errors.CCError) {

	bsi := NewBuildSameInfo(info, s.syncData)
	err := bsi.BuildSameInfoBaseCond(kit)
	if err != nil {
		return 0, false, err
	}

	switch tableName {
	case common.BKTableNameObjDes:
		err = bsi.BuildSameInfoObjDescCond(kit)
		if err != nil {
			return 0, false, err
		}
	case common.BKTableNameObjClassification:
		err = bsi.BuildSameInfoObjClassificationCond(kit)
		if err != nil {
			return 0, false, err
		}
	case common.BKTableNameObjAttDes:
		err = bsi.BuildSameInfoObjAttrDescCond(kit)
		if err != nil {
			return 0, false, err
		}
	case common.BKTableNamePropertyGroup:
		err = bsi.BuildSameInfoObjAttrGroupCond(kit)
		if err != nil {
			return 0, false, err
		}
	default:
		// merged data is not supported
		return 0, false, err
	}

	inst := mapstr.New()
	err = s.dbProxy.Table(tableName).Find(bsi.Condition()).One(kit.Ctx, &inst)
	if err != nil && !s.dbProxy.IsNotFoundError(err) {
		blog.Errorf("existSameInfo query db error. err:%s, DataClassify:%s,info:%#v,condition:%#v, rid:%s", err.Error(), bsi.syncData.DataClassify, info.Info, bsi.Condition(), kit.Rid)
		return 0, false, kit.CCError.Error(common.CCErrCommDBSelectFailed)
	}

	blog.V(6).Infof("getSameInfo DataClassify:%s, info:%#v, condition:%#v, inst:%#v, rid:%s", bsi.syncData.DataClassify, info.Info, bsi.Condition(), inst, kit.Rid)
	// not found data
	if len(inst) == 0 {
		return 0, false, nil
	}

	instID, err := inst.Int64(instIDField)
	if err != nil {
		blog.Errorf("buildSameInfoBaseCond get inst error. DataClassify:%s,info:%#v,rid:%s", bsi.syncData.DataClassify, info.Info, kit.Rid)
		return 0, false, kit.CCError.Errorf(common.CCErrCommInstFieldConvertFail, "propery data", instIDField, "int", err.Error())
	}
	return instID, true, nil

}

type buildSameInfo struct {
	info     *metadata.SynchronizeItem
	cond     mapstr.MapStr
	syncData *metadata.SynchronizeParameter
}

func NewBuildSameInfo(info *metadata.SynchronizeItem, syncData *metadata.SynchronizeParameter) *buildSameInfo {
	return &buildSameInfo{
		info:     info,
		cond:     mapstr.New(),
		syncData: syncData,
	}
}

func (bsi *buildSameInfo) BuildSameInfoBaseCond(kit *rest.Kit) errors.CCError {
	info := bsi.info
	if info.Info.Exists(common.MetadataField) {
		metadataVal, err := info.Info.MapStr(common.MetadataField)
		if err != nil {
			blog.Errorf("buildSameInfoBaseCond get metadata error. DataClassify:%s,info:%#v,rid:%s", bsi.syncData.DataClassify, info.Info, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommInstFieldConvertFail, "propery", common.MetadataField, "map", err.Error())
		}
		if metadataVal.Exists(metadata.LabelBusinessID) {
			str, err := metadataVal.String(metadata.LabelBusinessID)
			if err != nil {
				blog.Errorf("buildSameInfoBaseCond get metadata.bk_biz_id error. DataClassify:%s,info:%#v,rid:%s", bsi.syncData.DataClassify, info.Info, kit.Rid)
				return kit.CCError.CCErrorf(common.CCErrCommInstFieldConvertFail, "propery", common.MetadataField, "map", err.Error())
			}
			bsi.cond.Set("metadata.label.bk_biz_id", str)
		} else {
			bsi.cond.Merge(metadata.BizLabelNotExist)
		}

	} else {
		bsi.cond.Merge(metadata.BizLabelNotExist)
	}
	ownerID, err := info.Info.String(common.BKOwnerIDField)
	if err != nil {
		blog.Errorf("buildSameInfoBaseCond get ownerID error. DataClassify:%s,info:%#v,rid:%s", bsi.syncData.DataClassify, info.Info, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommInstFieldConvertFail, "propery", common.BKOwnerIDField, "string", err.Error())
	}
	bsi.cond = util.SetQueryOwner(bsi.cond, ownerID)
	return nil
}

func (bsi *buildSameInfo) BuildSameInfoObjDescCond(kit *rest.Kit) errors.CCError {
	info := bsi.info
	objID, err := info.Info.String(common.BKObjIDField)
	if err != nil {
		blog.Errorf("buildSameInfoObjDescCond get bk_obj_id error. DataClassify:%s,info:%#v,rid:%s", bsi.syncData.DataClassify, info.Info, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommInstFieldConvertFail, "propery", common.BKObjIDField, "string", err.Error())
	}

	bsi.cond.Set(common.BKObjIDField, objID)
	return nil
}

func (bsi *buildSameInfo) BuildSameInfoObjAttrDescCond(kit *rest.Kit) errors.CCError {
	info := bsi.info
	objID, err := info.Info.String(common.BKObjIDField)
	if err != nil {
		blog.Errorf("buildSameInfoObjAttrDescCond get bk_obj_id error. DataClassify:%s,info:%#v,rid:%s", bsi.syncData.DataClassify, info.Info, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommInstFieldConvertFail, "propery", common.BKObjIDField, "string", err.Error())
	}
	propertyID, err := info.Info.String(common.BKPropertyIDField)
	if err != nil {
		blog.Errorf("buildSameInfoObjAttrDescCond get bk_obj_name error. DataClassify:%s,info:%#v,rid:%s", bsi.syncData.DataClassify, info.Info, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommInstFieldConvertFail, "propery", common.BKPropertyIDField, "string", err.Error())
	}

	bsi.cond.Set(common.BKObjIDField, objID)
	bsi.cond.Set(common.BKPropertyIDField, propertyID)
	return nil
}

func (bsi *buildSameInfo) BuildSameInfoObjAttrGroupCond(kit *rest.Kit) errors.CCError {
	info := bsi.info
	objID, err := info.Info.String(common.BKObjIDField)
	if err != nil {
		blog.Errorf("existSameInfo get bk_obj_id error. DataClassify:%s,info:%#v,rid:%s", bsi.syncData.DataClassify, info.Info, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommInstFieldConvertFail, "propery", common.BKObjIDField, "string", err.Error())
	}
	groupID, err := info.Info.String(common.BKPropertyGroupIDField)
	if err != nil {
		blog.Errorf("existSameInfo get bk_group_id error. DataClassify:%s,info:%#v,rid:%s", bsi.syncData.DataClassify, info.Info, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommInstFieldConvertFail, "propery", common.BKPropertyGroupIDField, "string", err.Error())
	}

	bsi.cond.Set(common.BKObjIDField, objID)
	bsi.cond.Set(common.BKPropertyGroupIDField, groupID)
	return nil
}

func (bsi *buildSameInfo) BuildSameInfoObjClassificationCond(kit *rest.Kit) errors.CCError {
	info := bsi.info
	classificationID, err := info.Info.String(common.BKClassificationIDField)
	if err != nil {
		blog.Errorf("existSameInfo get bk_classification_id error. DataClassify:%s,info:%#v,rid:%s", bsi.syncData.DataClassify, info.Info, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommInstFieldConvertFail, "propery", common.BKClassificationIDField, "string", err.Error())
	}

	bsi.cond.Set(common.BKClassificationIDField, classificationID)
	return nil
}

func (bsi *buildSameInfo) Condition() mapstr.MapStr {
	return bsi.cond
}
