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
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (lgc *Logics) getAssociationData(ctx context.Context, header http.Header, objID string, instAsstArr []*metadata.InstAsst) (map[string]map[int64][]PropertyPrimaryVal, error) {

	// map[objID][]instID
	asstObjIDIDArr := make(map[string][]int64)
	for _, instAsst := range instAsstArr {
		_, ok := asstObjIDIDArr[instAsst.AsstObjectID]
		if !ok {
			asstObjIDIDArr[instAsst.AsstObjectID] = make([]int64, 0)
		}
		asstObjIDIDArr[instAsst.AsstObjectID] = append(asstObjIDIDArr[instAsst.AsstObjectID], instAsst.AsstInstID)
	}

	// map[objID]map[inst_id][]Property
	retAsstObjIDInstInfoMap := make(map[string]map[int64][]PropertyPrimaryVal)
	for itemObjID, asstInstIDArr := range asstObjIDIDArr {
		objPrimaryInfo, err := lgc.fetchInstAssocationData(ctx, header, itemObjID, asstInstIDArr)
		if err != nil {
			return nil, err
		}
		retAsstObjIDInstInfoMap[itemObjID] = objPrimaryInfo
	}

	return retAsstObjIDInstInfoMap, nil
}

func (lgc *Logics) fetchAssocationData(ctx context.Context, header http.Header, objID string, instIDArr []int64) ([]*metadata.InstAsst, error) {

	ccErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	input := &metadata.SearchAssociationInstRequest{}
	cond := condition.CreateCondition()
	cond.Field(common.BKObjIDField).Eq(objID)
	cond.Field(common.BKInstIDField).In(instIDArr)
	input.Condition = cond.ToMapStr()

	result, err := lgc.CoreAPI.ApiServer().SearchAssociationInst(ctx, header, input)
	if err != nil {
		blog.Errorf("GetAssocationData fetch %s association  error:%s, input;%+v, rid:%s", objID, err.Error(), input, util.GetHTTPCCRequestID(header))
		return nil, ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !result.Result {
		blog.Errorf("GetAssocationData fetch %s association  error code:%s, error msg:%s, input;%+v, rid:%s", objID, result.Code, result.ErrMsg, input, util.GetHTTPCCRequestID(header))
		return nil, ccErr.New(result.Code, result.ErrMsg)
	}

	return result.Data, nil
}

func (lgc *Logics) fetchInstAssocationData(ctx context.Context, header http.Header, objID string, instIDArr []int64) (map[int64][]PropertyPrimaryVal, error) {

	ccErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	propertyArr, err := lgc.getObjectPrimaryFieldByObjID(objID, header)
	if err != nil {
		return nil, err
	}
	var dbFields []string
	for _, property := range propertyArr {
		dbFields = append(dbFields, property.ID)
	}
	instIDKey := metadata.GetInstIDFieldByObjID(objID)

	dbFields = append(dbFields, instIDKey)

	instAsstCond := condition.CreateCondition()
	instAsstCond.Field(instIDKey).In(instIDArr)
	instAsstCond.SetFields(dbFields)

	instResult, err := lgc.CoreAPI.ApiServer().SearchInsts(ctx, header, objID, instAsstCond)
	if err != nil {
		blog.Errorf("GetAssocationData fetch %s association instance error:%s, input;%+v, rid:%s", objID, err.Error(), instAsstCond, util.GetHTTPCCRequestID(header))
		return nil, ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !instResult.Result {
		blog.Errorf("FetchInstAssocationData fetch %s association instance error code:%s, error msg:%s, input;%+v, rid:%s", objID, instResult.Code, instResult.ErrMsg, instAsstCond, util.GetHTTPCCRequestID(header))
		return nil, ccErr.New(instResult.Code, instResult.ErrMsg)
	}

	retAsstInstInfo := make(map[int64][]PropertyPrimaryVal, 0)
	for _, inst := range instResult.Data.Info {
		instID, err := inst.Int64(instIDKey)
		if err != nil {
			blog.Warnf("FetchInstAssocationData get %s instance %s field error, err:%s, inst:%+v, rid:%s", objID, instIDKey, err.Error(), inst, util.GetHTTPCCRequestID(header))
			continue
		}
		isSkip := false
		var primaryKeysVal []PropertyPrimaryVal
		for _, key := range propertyArr {
			// use display , use string
			val, err := inst.String(key.ID)
			if err != nil {
				blog.Warnf("FetchInstAssocationData get %s instance %s field error, err:%s, inst:%+v, rid:%s", objID, key, err.Error(), inst, util.GetHTTPCCRequestID(header))
				isSkip = true
				break
			}
			primaryKeysVal = append(primaryKeysVal, PropertyPrimaryVal{
				ID:     key.ID,
				Name:   key.Name,
				StrVal: val,
			})
		}
		if isSkip {
			continue
		}
		retAsstInstInfo[instID] = primaryKeysVal

	}
	return retAsstInstInfo, nil

}
