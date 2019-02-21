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

	"bk-cmdb/src/common/blog"
	"configcenter/src/common"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/synchronize_server/app/options"
)

func getVersion() string {
	// TODO
	return "----"
}

// Synchronize  synchronize data
func (lgc *Logics) SynchronizeItem(ctx context.Context, syncConfig *options.ConfigItem) (exceptions []metadata.ExceptionResult, err error) {
	version := getVersion()
	lgc.synchronizeInstanceTask(ctx, syncConfig, version, nil)
	lgc.synchronizeModelTask(ctx, syncConfig, version, nil)
	lgc.synchronizeAssociationTask(ctx, syncConfig, version, nil)
	return nil, nil
}

func (lgc *Logics) synchronizeInstanceTask(ctx context.Context, syncConfig *options.ConfigItem, version string, baseCondition mapstr.MapStr) ([]metadata.ExceptionResult, errors.CCError) {
	inst := lgc.NewFetchInst(syncConfig, baseCondition)
	// app
	lgc.synchronizeInstance(ctx, common.BKInnerObjIDApp, inst, version, syncConfig)
	// set
	lgc.synchronizeInstance(ctx, common.BKInnerObjIDSet, inst, version, syncConfig)
	// model
	lgc.synchronizeInstance(ctx, common.BKInnerObjIDModule, inst, version, syncConfig)
	// proc
	lgc.synchronizeInstance(ctx, common.BKInnerObjIDProc, inst, version, syncConfig)
	// plat
	lgc.synchronizeInstance(ctx, common.BKInnerObjIDPlat, inst, version, syncConfig)
	// host
	lgc.synchronizeInstance(ctx, common.BKInnerObjIDHost, inst, version, syncConfig)
	// object
	lgc.synchronizeInstance(ctx, common.BKInnerObjIDObject, inst, version, syncConfig)

	return nil, nil
}

func (lgc *Logics) synchronizeInstance(ctx context.Context, objID string, inst *FetchInst, version string, syncConfig *options.ConfigItem) ([]metadata.ExceptionResult, error) {
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
		input.Version = version
		input.SynchronizeFlag = syncConfig.SynchronizeFlag
		// synchronize api
		pageErrInfoArr, err := lgc.sycnhronizePartInstance(ctx, input)
		if err != nil {
			return nil, err
		}
		if len(pageErrInfoArr) > 0 {
			errorInfoArr = append(errorInfoArr, pageErrInfoArr...)
		}

		start += limit
		if start < int64(info.Count) {
			break
		}
	}

	return errorInfoArr, nil
}

func (lgc *Logics) sycnhronizePartInstance(ctx context.Context, input *metadata.SynchronizeDataInfo) ([]metadata.ExceptionResult, error) {
	var errorInfoArr []metadata.ExceptionResult
	synchronizeParameter := &metadata.SynchronizeParameter{
		OperateType:     metadata.SynchronizeOperateTypeRepalce,
		OperateDataType: input.OperateDataType,
		DataClassify:    input.DataClassify,
		Version:         input.Version,
		SynchronizeFlag: input.SynchronizeFlag,
	}

	for _, item := range input.InfoArray {
		IDField := common.GetInstIDField(input.DataClassify)
		id, err := item.Int64(IDField)
		if err != nil {
			//  CCErrCommInstFieldConvFail  convert %s  field %s to %s error %s
			ccErrInst := lgc.ccErr.Errorf(common.CCErrCommInstFieldConvFail, input.DataClassify, IDField, "int64", err.Error())
			errorInfoArr = append(errorInfoArr, metadata.ExceptionResult{
				Code:        common.CCErrCommInstFieldConvFail,
				Message:     ccErrInst.Error(),
				Data:        item,
				OriginIndex: 0,
			})
		}
		synchronizeParameter.InfoArray = append(synchronizeParameter.InfoArray, &metadata.SynchronizeItem{ID: id, Info: item})
	}

	result, err := lgc.CoreAPI.CoreService().Synchronize().SynchronizeInstance(ctx, lgc.header, synchronizeParameter)
	if err != nil {
		blog.Errorf("sycnhronizePartInstance http do error, error: %s,DataSign: %s,DataTeyp: %s,rid:%s", err.Error(), input.DataClassify, input.OperateDataType, lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)

		//  CCErrCommInstFieldConvFail  convert %s  field %s to %s error %s
		/*errorInfoArr = append(errorInfoArr, metadata.ExceptionResult{
			Code:        common.CCErrCommHTTPDoRequestFailed,
			Message:     err.Error(),
			Data:        input.InfoArray,
			OriginIndex: 0,
		})*/
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

func (lgc *Logics) synchronizeModelTask(ctx context.Context, syncConfig *options.ConfigItem, version string, baseCondition mapstr.MapStr) ([]metadata.ExceptionResult, errors.CCError) {
	var errorInfoArr []metadata.ExceptionResult

	model := lgc.NewFetchModel(syncConfig, nil)
	classifyArr := []string{
		common.SynchronizeModelTypeBase,
		common.SynchronizeModelTypeClassification,
		common.SynchronizeModelTypeAttribute,
		common.SynchronizeModelTypeAttributeGroup,
	}
	for _, dataClassify := range classifyArr {
		partErrorInfoArr, err := lgc.synchronizeModel(ctx, model, version, dataClassify, syncConfig)
		if err != nil {
			return nil, err
		}
		if len(partErrorInfoArr) > 0 {
			errorInfoArr = append(errorInfoArr, partErrorInfoArr...)
		}

	}

	return errorInfoArr, nil
}

func (lgc *Logics) synchronizeModel(ctx context.Context, model *FetchModel, dataClassify string, version string, syncConfig *options.ConfigItem) ([]metadata.ExceptionResult, error) {
	var start int64 = 0
	limit := int64(defaultLimit)
	var errorInfoArr []metadata.ExceptionResult

	for {

		info, err := model.Fetch(ctx, dataClassify, start, limit)
		if err != nil {
			return nil, err
		}

		input := &metadata.SynchronizeDataInfo{}
		input.OperateDataType = metadata.SynchronizeOperateDataTypeModel
		input.DataClassify = dataClassify
		input.InfoArray = info.Info
		input.Version = version
		input.SynchronizeFlag = syncConfig.SynchronizeFlag
		// synchronize api
		pageErrInfoArr, err := lgc.sycnhronizePartModel(ctx, input)
		if err != nil {
			return nil, err
		}
		if len(pageErrInfoArr) > 0 {
			errorInfoArr = append(errorInfoArr, pageErrInfoArr...)
		}
		start += limit
		if start <= int64(info.Count) {
			break
		}
	}

	return errorInfoArr, nil
}

func (lgc *Logics) sycnhronizePartModel(ctx context.Context, input *metadata.SynchronizeDataInfo) ([]metadata.ExceptionResult, errors.CCError) {
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
			ccErrInst := lgc.ccErr.Errorf(common.CCErrCommInstFieldConvFail, "model attribute", common.BKFieldID, "int64", err.Error())
			errorInfoArr = append(errorInfoArr, metadata.ExceptionResult{
				Code:        common.CCErrCommInstFieldConvFail,
				Message:     ccErrInst.Error(),
				Data:        item,
				OriginIndex: 0,
			})
		}
		synchronizeParameter.InfoArray = append(synchronizeParameter.InfoArray, &metadata.SynchronizeItem{ID: id, Info: item})
	}

	result, err := lgc.CoreAPI.CoreService().Synchronize().SynchronizeModel(ctx, lgc.header, synchronizeParameter)
	if err != nil {
		blog.Errorf("sycnhronizePartModel http do error, error: %s,DataSign: %s,DataTeyp: %s,rid:%s", err.Error(), input.DataClassify, input.OperateDataType, lgc.rid)
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

func (lgc *Logics) synchronizeAssociationTask(ctx context.Context, syncConfig *options.ConfigItem, version string, baseCondition mapstr.MapStr) ([]metadata.ExceptionResult, errors.CCError) {
	var errorInfoArr []metadata.ExceptionResult

	model := lgc.NewFetchAssociation(syncConfig, nil)
	classifyArr := []string{
		common.SynchronizeAssociationTypeModelHost,
	}
	for _, dataClassify := range classifyArr {
		partErrorInfoArr, err := lgc.sycnhronizeAssociation(ctx, model, version, dataClassify, syncConfig)
		if err != nil {
			return nil, err
		}
		if len(partErrorInfoArr) > 0 {
			errorInfoArr = append(errorInfoArr, partErrorInfoArr...)
		}

	}

	return errorInfoArr, nil
}

func (lgc *Logics) sycnhronizeAssociation(ctx context.Context, association *FetchAssociation, dataClassify string, version string, syncConfig *options.ConfigItem) ([]metadata.ExceptionResult, error) {
	var start int64 = 0
	limit := int64(defaultLimit)
	var errorInfoArr []metadata.ExceptionResult

	for {

		info, err := association.Fetch(ctx, dataClassify, start, limit)
		if err != nil {
			return nil, err
		}

		input := &metadata.SynchronizeDataInfo{}
		input.OperateDataType = metadata.SynchronizeOperateDataTypeModel
		input.DataClassify = dataClassify
		input.InfoArray = info.Info
		input.Version = version
		input.SynchronizeFlag = syncConfig.SynchronizeFlag
		// synchronize api
		pageErrInfoArr, err := lgc.sycnhronizePartAssociation(ctx, input)
		if err != nil {
			return nil, err
		}
		if len(pageErrInfoArr) > 0 {
			errorInfoArr = append(errorInfoArr, pageErrInfoArr...)
		}
		start += limit
		if start <= int64(info.Count) {
			break
		}
	}

	return errorInfoArr, nil
}

func (lgc *Logics) sycnhronizePartAssociation(ctx context.Context, input *metadata.SynchronizeDataInfo) ([]metadata.ExceptionResult, errors.CCError) {
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
			ccErrInst := lgc.ccErr.Errorf(common.CCErrCommInstFieldConvFail, "model associate", common.BKFieldID, "int64", err.Error())
			errorInfoArr = append(errorInfoArr, metadata.ExceptionResult{
				Code:        common.CCErrCommInstFieldConvFail,
				Message:     ccErrInst.Error(),
				Data:        item,
				OriginIndex: 0,
			})
		}
		synchronizeParameter.InfoArray = append(synchronizeParameter.InfoArray, &metadata.SynchronizeItem{ID: id, Info: item})
	}

	result, err := lgc.CoreAPI.CoreService().Synchronize().SynchronizeAssociation(ctx, lgc.header, synchronizeParameter)
	if err != nil {
		blog.Errorf("sycnhronizePartAssociation http do error, error: %s,DataSign: %s,DataTeyp: %s,rid:%s", err.Error(), input.DataClassify, input.OperateDataType, lgc.rid)
		//  CCErrCommInstFieldConvFail  convert %s  field %s to %s error %s
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
