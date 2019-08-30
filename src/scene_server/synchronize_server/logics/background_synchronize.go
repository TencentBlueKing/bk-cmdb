/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package logics

import (
	"context"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/synchronize_server/app/options"
	"configcenter/src/scene_server/synchronize_server/logics/exception/file"
)

var (
	nextDayTrigger int64 = 24 * 60
)

type synchronizeItemInterface interface {
	synchronizeItemException(ctx context.Context, exceptionMap map[string][]metadata.ExceptionResult)
	synchronizeInstanceTask(ctx context.Context) (errorInfoArr []metadata.ExceptionResult, err errors.CCError)
	synchronizeModelTask(ctx context.Context) ([]metadata.ExceptionResult, errors.CCError)
	synchronizeAssociationTask(ctx context.Context) ([]metadata.ExceptionResult, errors.CCError)
	synchronizeItemClearData(ctx context.Context) (map[string][]metadata.ExceptionResult, errors.CCError)
}

type synchronizeItem struct {
	lgc           *Logics
	config        *options.ConfigItem
	baseCondition mapstr.MapStr
	// all objID
	objIDMap map[string]bool
	// modelClassifitionMap sync model classifition
	modelClassifitionMap map[string]bool
	appIDArr             []int64
	version              int64
}

func (lgc *Logics) NewSynchronizeItem(version int64, syncConfig *options.ConfigItem) synchronizeItemInterface {

	ret := &synchronizeItem{
		lgc:                  lgc,
		config:               syncConfig,
		baseCondition:        mapstr.New(),
		objIDMap:             make(map[string]bool, 0),
		modelClassifitionMap: make(map[string]bool, 0),
		appIDArr:             make([]int64, 0),
		version:              version,
	}
	ret.configPretreatment()
	return ret
}

func (s *synchronizeItem) configPretreatment() {
	if len(s.config.ObjectIDArr) > 0 && s.config.WhiteList {
		objectIDArr := []string{common.BKInnerObjIDApp, common.BKInnerObjIDSet, common.BKInnerObjIDModule, common.BKInnerObjIDHost, common.BKInnerObjIDProc, common.BKInnerObjIDPlat}
		s.config.ObjectIDArr = append(s.config.ObjectIDArr, objectIDArr...)
	}
	return
}

func (s *synchronizeItem) synchronizeItemClearData(ctx context.Context) (map[string][]metadata.ExceptionResult, errors.CCError) {
	errorInfoArr := make(map[string][]metadata.ExceptionResult)

	clearDataInput := &metadata.SynchronizeClearDataParameter{
		Version:         s.version,
		SynchronizeFlag: s.config.SynchronizeFlag,
		Tamestamp:       time.Now().Unix(),
	}
	clearDataInput.GenerateSign(common.SynchronizeSignPrefix)
	result, err := s.lgc.CoreAPI.CoreService().Synchronize().SynchronizeClearData(ctx, s.lgc.header, clearDataInput)
	if err != nil {
		blog.Errorf("SynchronizeItem SynchronizeClearData error, config:%#v,err:%s,version:%d,rid:%s", s.config, err.Error(), s.version, s.lgc.rid)
		return errorInfoArr, nil
	}
	if !result.Result {
		var errorInfoPart []metadata.ExceptionResult
		err = s.lgc.ccErr.New(result.Code, result.ErrMsg)
		errorInfoPart = append(errorInfoPart, metadata.ExceptionResult{
			Code:        int64(result.Code),
			Message:     err.Error(),
			Data:        nil,
			OriginIndex: 0,
		})
		errorInfoArr["clear_data"] = errorInfoPart

	}
	return errorInfoArr, nil
}

func (s *synchronizeItem) synchronizeItemException(ctx context.Context, exceptionMap map[string][]metadata.ExceptionResult) {
	exceptionFile, err := file.NewException(s.config.SynchronizeFlag, s.version, s.config.ExceptionFileCount)
	if err != nil {
		blog.Errorf("synchronizeItemException  error, config:%#v,err:%s,version:%d,rid:%s", s.config, err.Error(), s.version, s.lgc.rid)
		return
	}
	for exceptionType, exceptions := range exceptionMap {
		err := exceptionFile.WriteStringln("synchronize " + exceptionType + " exception start")
		if err != nil {
			blog.Errorf("synchronizeItemException %s error, config:%#v,err:%s,version:%d,rid:%s", exceptionType, s.config, err.Error(), s.version, s.lgc.rid)
		}
		err = exceptionFile.Write(ctx, exceptions)
		if err != nil {
			blog.Errorf("synchronizeItemException %s error, config:%#v,err:%s,version:%d,rid:%s", exceptionType, s.config, err.Error(), s.version, s.lgc.rid)
		}
		err = exceptionFile.WriteStringln("synchronize " + exceptionType + " exception end")
		if err != nil {
			blog.Errorf("synchronizeItemException %s error, config:%#v,err:%s,version:%d,rid:%s", exceptionType, s.config, err.Error(), s.version, s.lgc.rid)
		}
	}

}

func (s *synchronizeItem) synchronizeInstanceTask(ctx context.Context) (errorInfoArr []metadata.ExceptionResult, err errors.CCError) {

	inst := s.lgc.NewFetchInst(s.config, s.baseCondition)
	err = inst.Pretreatment()
	if err != nil {
		blog.Errorf("instance Pretreatment error. err:%s, rid:%s", err.Error(), s.lgc.rid)
		return nil, err
	}
	var partErrorInfoArr []metadata.ExceptionResult

	// must first business.  get synchronize app information,
	// set,modelneeds to be synchronized, and the required model can be filtered by the service,
	//  the model needs the business information in the subsequent service.
	if _, ok := s.objIDMap[common.BKInnerObjIDApp]; ok {
		partErrorInfoArrItem, err := s.synchronizeInstance(ctx, common.BKInnerObjIDApp, inst)
		if err != nil {
			blog.Errorf("synchronizeInstanceTask synchronize %s error,err:%s,rid:%s", common.BKInnerObjIDApp, err.Error(), s.lgc.rid)
			return nil, err
		}
		if len(partErrorInfoArr) > 0 {
			partErrorInfoArr = append(partErrorInfoArr, partErrorInfoArrItem...)
		}
		inst.SetAppIDArr(s.appIDArr)
	}
	for objID := range s.objIDMap {
		if objID == common.BKInnerObjIDApp {
			// last already synchronize
			continue
		}
		partErrorInfoArrItem, err := s.synchronizeInstance(ctx, objID, inst)
		if err != nil {
			blog.Errorf("synchronizeInstanceTask synchronize %s error,err:%s,rid:%s", objID, err.Error(), s.lgc.rid)
			return nil, err
		}
		if len(partErrorInfoArrItem) > 0 {
			partErrorInfoArr = append(partErrorInfoArr, partErrorInfoArrItem...)
		}

	}

	return
}

func (s *synchronizeItem) synchronizeInstance(ctx context.Context, objID string, inst *FetchInst) ([]metadata.ExceptionResult, error) {
	var start int64 = 0
	limit := int64(defaultLimit)
	var errorInfoArr []metadata.ExceptionResult

	for {
		info, err := inst.Fetch(ctx, objID, start, limit)
		if err != nil {
			return nil, err
		}

		input := &metadata.SynchronizeDataInfo{}
		input.OperateDataType = metadata.SynchronizeOperateDataTypeInstance
		input.DataClassify = objID
		input.InfoArray = info.Info
		input.Version = s.version
		input.SynchronizeFlag = s.config.SynchronizeFlag
		// synchronize api
		pageErrInfoArr, err := s.sycnhronizePartInstance(ctx, input)
		if err != nil {
			return nil, err
		}
		if len(pageErrInfoArr) > 0 {
			errorInfoArr = append(errorInfoArr, pageErrInfoArr...)
		}

		start += limit
		if start >= int64(info.Count) {
			break
		}
	}
	if common.BKInnerObjIDApp == objID {
		inst.SetAppIDArr(s.appIDArr)
	}

	return errorInfoArr, nil
}

func (s *synchronizeItem) sycnhronizePartInstance(ctx context.Context, input *metadata.SynchronizeDataInfo) ([]metadata.ExceptionResult, error) {
	lgc := s.lgc
	var errorInfoArr []metadata.ExceptionResult
	synchronizeParameter := &metadata.SynchronizeParameter{
		OperateType:     metadata.SynchronizeOperateTypeRepalce,
		OperateDataType: input.OperateDataType,
		DataClassify:    input.DataClassify,
		Version:         input.Version,
		SynchronizeFlag: input.SynchronizeFlag,
	}
	if len(input.InfoArray) == 0 {
		return errorInfoArr, nil
	}

	for _, item := range input.InfoArray {
		IDField := common.GetInstIDField(input.DataClassify)
		id, err := item.Int64(IDField)
		if err != nil {
			//  CCErrCommInstFieldConvFail  convert %s  field %s to %s error %s
			ccErrInst := lgc.ccErr.Errorf(common.CCErrCommInstFieldConvertFail, input.DataClassify, IDField, "int64", err.Error())
			errorInfoArr = append(errorInfoArr, metadata.ExceptionResult{
				Code:        common.CCErrCommInstFieldConvertFail,
				Message:     ccErrInst.Error(),
				Data:        item,
				OriginIndex: 0,
			})
		}
		if common.BKInnerObjIDApp == input.DataClassify {
			s.appIDArr = append(s.appIDArr, id)
		}
		synchronizeParameter.InfoArray = append(synchronizeParameter.InfoArray, &metadata.SynchronizeItem{ID: id, Info: item})
	}

	result, err := lgc.CoreAPI.CoreService().Synchronize().SynchronizeInstance(ctx, lgc.header, synchronizeParameter)
	if err != nil {
		blog.Errorf("sycnhronizePartInstance http do error, error: %s,DataSign: %s,DataTeyp: %s,rid:%s", err.Error(), input.DataClassify, input.OperateDataType, lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		if len(result.Data.Exceptions) == 0 {
			errorInfoArr = append(errorInfoArr, metadata.ExceptionResult{
				Code:        int64(result.Code),
				Message:     result.ErrMsg,
				Data:        input.InfoArray,
				OriginIndex: 0,
			})
		} else {
			errorInfoArr = append(errorInfoArr, result.Data.Exceptions...)
		}
	}

	return errorInfoArr, nil
}

func (s *synchronizeItem) synchronizeModelTask(ctx context.Context) ([]metadata.ExceptionResult, errors.CCError) {
	var errorInfoArr []metadata.ExceptionResult
	objIDCond := condition.CreateCondition()
	if len(s.config.ObjectIDArr) > 0 {
		if s.config.WhiteList {
			objIDCond.Field(common.BKObjIDField).In(s.config.ObjectIDArr)
		} else {
			objIDCond.Field(common.BKObjIDField).NotIn(s.config.ObjectIDArr)
		}
	}
	conditionMapStr := objIDCond.ToMapStr()

	model := s.lgc.NewFetchModel(s.config, s.baseCondition)
	err := model.Pretreatment()
	if err != nil {
		blog.Errorf("model Pretreatment error. err:%s, rid:%s", err.Error(), s.lgc.rid)
		return nil, err
	}
	classifyArr := []string{
		common.SynchronizeModelTypeBase,
		common.SynchronizeModelTypeAttribute,
		common.SynchronizeModelTypeAttributeGroup,
	}
	for _, dataClassify := range classifyArr {
		partErrorInfoArr, err := s.synchronizeModel(ctx, model, conditionMapStr, dataClassify)
		if err != nil {
			return nil, err
		}
		if len(partErrorInfoArr) > 0 {
			errorInfoArr = append(errorInfoArr, partErrorInfoArr...)
		}

	}

	if len(s.modelClassifitionMap) > 0 {
		classifitionStrIDArr := getMapStrBoolKey(s.modelClassifitionMap)

		classificationCond := condition.CreateCondition()
		classificationCond.Field(common.BKClassificationIDField).In(classifitionStrIDArr)

		// classifition special. model sync end. can sync classification
		dataClassify := common.SynchronizeModelTypeClassification
		partErrorInfoArr, err := s.synchronizeModel(ctx, model, conditionMapStr, dataClassify)
		if err != nil {
			return nil, err
		}
		if len(partErrorInfoArr) > 0 {
			errorInfoArr = append(errorInfoArr, partErrorInfoArr...)
		}
	}

	return errorInfoArr, nil
}

func (s *synchronizeItem) synchronizeModel(ctx context.Context, model *FetchModel, cond mapstr.MapStr, dataClassify string) ([]metadata.ExceptionResult, error) {
	var start int64 = 0
	limit := int64(defaultLimit)
	var errorInfoArr []metadata.ExceptionResult

	for {

		info, err := model.Fetch(ctx, dataClassify, cond, start, limit)
		if err != nil {
			return nil, err
		}

		input := &metadata.SynchronizeDataInfo{}
		input.OperateDataType = metadata.SynchronizeOperateDataTypeModel
		input.DataClassify = dataClassify
		input.InfoArray = info.Info
		input.Version = s.version
		input.SynchronizeFlag = s.config.SynchronizeFlag
		// synchronize api
		pageErrInfoArr, err := s.sycnhronizePartModel(ctx, input)
		if err != nil {
			return nil, err
		}
		if len(pageErrInfoArr) > 0 {
			errorInfoArr = append(errorInfoArr, pageErrInfoArr...)
		}
		start += limit
		if start >= int64(info.Count) {
			break
		}
	}

	return errorInfoArr, nil
}

func (s *synchronizeItem) sycnhronizePartModel(ctx context.Context, input *metadata.SynchronizeDataInfo) ([]metadata.ExceptionResult, errors.CCError) {
	var errorInfoArr []metadata.ExceptionResult
	synchronizeParameter := &metadata.SynchronizeParameter{
		OperateType:     metadata.SynchronizeOperateTypeRepalce,
		OperateDataType: input.OperateDataType,
		DataClassify:    input.DataClassify,
		Version:         input.Version,
		SynchronizeFlag: input.SynchronizeFlag,
	}

	for _, item := range input.InfoArray {
		id, err := item.Int64(common.BKFieldID)
		if err != nil {
			//  CCErrCommInstFieldConvFail  convert %s  field %s to %s error %s
			ccErrInst := s.lgc.ccErr.Errorf(common.CCErrCommInstFieldConvertFail, "model attribute", common.BKFieldID, "int64", err.Error())
			errorInfoArr = append(errorInfoArr, metadata.ExceptionResult{
				Code:        common.CCErrCommInstFieldConvertFail,
				Message:     ccErrInst.Error(),
				Data:        item,
				OriginIndex: 0,
			})
		}
		objID, err := item.String(common.BKObjIDField)
		if err == nil && objID != "" {
			s.objIDMap[objID] = true
		}
		classifitionStrID, err := item.String(common.BKClassificationIDField)
		if err == nil && classifitionStrID != "" {
			s.modelClassifitionMap[classifitionStrID] = true
		}
		synchronizeParameter.InfoArray = append(synchronizeParameter.InfoArray, &metadata.SynchronizeItem{ID: id, Info: item})
	}
	// 当配置不需要同步模型的时候。不去保存数据，
	// 但是需要使用上面的功能获取需要同步模型id。 s.ojjIDMap
	if s.config.IgnoreModelAttr {
		blog.V(4).Infof("sycnhronizePartModel skip. rid:%s, version:%v", s.lgc.rid, s.version)
		return errorInfoArr, nil
	}
	result, err := s.lgc.CoreAPI.CoreService().Synchronize().SynchronizeModel(ctx, s.lgc.header, synchronizeParameter)
	if err != nil {
		blog.Errorf("sycnhronizePartModel http do error, error: %s,DataSign: %s,DataTeyp: %s,rid:%s", err.Error(), input.DataClassify, input.OperateDataType, s.lgc.rid)
		return nil, s.lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)

	}
	if !result.Result {
		if len(result.Data.Exceptions) == 0 {
			errorInfoArr = append(errorInfoArr, metadata.ExceptionResult{
				Code:        int64(result.Code),
				Message:     result.ErrMsg,
				Data:        input.InfoArray,
				OriginIndex: 0,
			})
		} else {
			errorInfoArr = append(errorInfoArr, result.Data.Exceptions...)
		}
	}
	return errorInfoArr, nil
}

func (s *synchronizeItem) synchronizeAssociationTask(ctx context.Context) ([]metadata.ExceptionResult, errors.CCError) {
	var errorInfoArr []metadata.ExceptionResult

	association := s.lgc.NewFetchAssociation(s.config, s.baseCondition)

	classifyArr := []string{
		common.SynchronizeAssociationTypeModelHost,
	}
	association.SetAppIDArr(s.appIDArr)
	for _, dataClassify := range classifyArr {
		partErrorInfoArr, err := s.sycnhronizeAssociation(ctx, association, dataClassify)
		if err != nil {
			return nil, err
		}
		if len(partErrorInfoArr) > 0 {
			errorInfoArr = append(errorInfoArr, partErrorInfoArr...)
		}

	}

	return errorInfoArr, nil
}

func (s *synchronizeItem) sycnhronizeAssociation(ctx context.Context, association *FetchAssociation, dataClassify string) ([]metadata.ExceptionResult, error) {
	var start int64 = 0
	limit := int64(defaultLimit)
	var errorInfoArr []metadata.ExceptionResult

	for {

		info, err := association.Fetch(ctx, dataClassify, start, limit)
		if err != nil {
			return nil, err
		}

		input := &metadata.SynchronizeDataInfo{}
		input.OperateDataType = metadata.SynchronizeOperateDataTypeAssociation
		input.DataClassify = dataClassify
		input.InfoArray = info.Info
		input.Version = s.version
		input.SynchronizeFlag = s.config.SynchronizeFlag
		// synchronize api
		pageErrInfoArr, err := s.sycnhronizePartAssociation(ctx, input)
		if err != nil {
			return nil, err
		}
		if len(pageErrInfoArr) > 0 {
			errorInfoArr = append(errorInfoArr, pageErrInfoArr...)
		}
		start += limit
		if start >= int64(info.Count) {
			break
		}
	}

	return errorInfoArr, nil
}

func (s *synchronizeItem) sycnhronizePartAssociation(ctx context.Context, input *metadata.SynchronizeDataInfo) ([]metadata.ExceptionResult, errors.CCError) {
	var errorInfoArr []metadata.ExceptionResult

	synchronizeParameter := &metadata.SynchronizeParameter{
		OperateType:     metadata.SynchronizeOperateTypeRepalce,
		OperateDataType: input.OperateDataType,
		DataClassify:    input.DataClassify,
		Version:         input.Version,
		SynchronizeFlag: input.SynchronizeFlag,
	}
	for _, item := range input.InfoArray {
		id, err := item.Int64(common.BKFieldID)
		if err != nil {
			//  CCErrCommInstFieldConvFail  convert %s  field %s to %s error %s
			ccErrInst := s.lgc.ccErr.Errorf(common.CCErrCommInstFieldConvertFail, "model associate", common.BKFieldID, "int64", err.Error())
			errorInfoArr = append(errorInfoArr, metadata.ExceptionResult{
				Code:        common.CCErrCommInstFieldConvertFail,
				Message:     ccErrInst.Error(),
				Data:        item,
				OriginIndex: 0,
			})
		}
		synchronizeParameter.InfoArray = append(synchronizeParameter.InfoArray, &metadata.SynchronizeItem{ID: id, Info: item})
	}

	result, err := s.lgc.CoreAPI.CoreService().Synchronize().SynchronizeAssociation(ctx, s.lgc.header, synchronizeParameter)
	if err != nil {
		blog.Errorf("sycnhronizePartAssociation http do error, error: %s,DataSign: %s,DataTeyp: %s,rid:%s", err.Error(), input.DataClassify, input.OperateDataType, s.lgc.rid)
		//  CCErrCommInstFieldConvFail  convert %s  field %s to %s error %s
		return nil, s.lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)

	}
	if !result.Result {
		if len(result.Data.Exceptions) == 0 {
			errorInfoArr = append(errorInfoArr, metadata.ExceptionResult{
				Code:        int64(result.Code),
				Message:     result.ErrMsg,
				Data:        input.InfoArray,
				OriginIndex: 0,
			})
		} else {
			errorInfoArr = append(errorInfoArr, result.Data.Exceptions...)
		}
	}
	return errorInfoArr, nil
}

func getMapStrBoolKey(data map[string]bool) []string {
	var keyArr []string
	for key := range data {
		keyArr = append(keyArr, key)
	}

	return keyArr
}
