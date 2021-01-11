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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/querybuilder"
	"configcenter/src/common/util"
)

func (lgc *Logics) getAssociationData(ctx context.Context, header http.Header, objID string, instAsstArr []*metadata.InstAsst, modelBizID int64) (map[string]map[int64][]PropertyPrimaryVal, error) {

	// map[objID][]instID
	asstObjIDIDArr := make(map[string][]int64)
	for _, instAsst := range instAsstArr {
		_, ok := asstObjIDIDArr[instAsst.AsstObjectID]
		if !ok {
			asstObjIDIDArr[instAsst.AsstObjectID] = make([]int64, 0)
		}
		asstObjIDIDArr[instAsst.AsstObjectID] = append(asstObjIDIDArr[instAsst.AsstObjectID], instAsst.AsstInstID)
		_, ok = asstObjIDIDArr[instAsst.ObjectID]
		if !ok {
			asstObjIDIDArr[instAsst.ObjectID] = make([]int64, 0)
		}
		asstObjIDIDArr[instAsst.ObjectID] = append(asstObjIDIDArr[instAsst.ObjectID], instAsst.InstID)
	}

	// map[objID]map[inst_id][]Property
	retAsstObjIDInstInfoMap := make(map[string]map[int64][]PropertyPrimaryVal)
	for itemObjID, asstInstIDArr := range asstObjIDIDArr {
		objPrimaryInfo, err := lgc.fetchInstAssocationData(ctx, header, itemObjID, asstInstIDArr, modelBizID)
		if err != nil {
			return nil, err
		}
		retAsstObjIDInstInfoMap[itemObjID] = objPrimaryInfo
	}
	return retAsstObjIDInstInfoMap, nil
}

func (lgc *Logics) fetchAssocationData(ctx context.Context, header http.Header, objID string, instIDArr []int64, modelBizID int64) ([]*metadata.InstAsst, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	ccErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	input := &metadata.SearchAssociationInstRequest{}

	//实例作为关联关系源模型
	cond := condition.CreateCondition()
	cond.Field(common.BKObjIDField).Eq(objID)
	cond.Field(common.BKInstIDField).In(instIDArr)
	input.Condition = cond.ToMapStr()
	if modelBizID > 0 {
		input.Condition.Set(common.BKAppIDField, modelBizID)
	}
	bkObjRst, err := lgc.CoreAPI.ApiServer().SearchAssociationInst(ctx, header, input)
	if err != nil {
		blog.Errorf("GetAssocationData fetch %s association  error:%s, input;%+v, rid: %s", objID, err.Error(), input, rid)
		return nil, ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !bkObjRst.Result {
		blog.Errorf("GetAssocationData fetch %s association  error code:%s, error msg:%s, input;%+v, rid:%s", objID, bkObjRst.Code, bkObjRst.ErrMsg, input, rid)
		return nil, ccErr.New(bkObjRst.Code, bkObjRst.ErrMsg)
	}

	//实例作为关联关系目标模型
	cond = condition.CreateCondition()
	cond.Field(common.BKAsstObjIDField).Eq(objID)
	cond.Field(common.BKAsstInstIDField).In(instIDArr)
	input.Condition = cond.ToMapStr()
	bkAsstObjRst, err := lgc.CoreAPI.ApiServer().SearchAssociationInst(ctx, header, input)
	if err != nil {
		blog.Errorf("GetAssocationData fetch %s association  error:%s, input;%+v, rid: %s", objID, err.Error(), input, rid)
		return nil, ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !bkAsstObjRst.Result {
		blog.Errorf("GetAssocationData fetch %s association  error code:%s, error msg:%s, input;%+v, rid:%s", objID, bkAsstObjRst.Code, bkAsstObjRst.ErrMsg, input, rid)
		return nil, ccErr.New(bkAsstObjRst.Code, bkAsstObjRst.ErrMsg)
	}
	result := append(bkObjRst.Data[:], bkAsstObjRst.Data...)

	return result, nil
}

func (lgc *Logics) fetchInstAssocationData(ctx context.Context, header http.Header, objID string, instIDArr []int64, modelBizID int64) (map[int64][]PropertyPrimaryVal, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	ccErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	propertyArr, err := lgc.getObjectPrimaryFieldByObjID(objID, header, modelBizID)
	if err != nil {
		return nil, err
	}
	var dbFields []string
	for _, property := range propertyArr {
		dbFields = append(dbFields, property.ID)
	}
	instIDKey := metadata.GetInstIDFieldByObjID(objID)

	dbFields = append(dbFields, instIDKey)

	insts := make([]mapstr.MapStr, 0)
	switch objID {
	case common.BKInnerObjIDHost:
		option := metadata.ListHostsWithNoBizParameter{
			HostPropertyFilter: &querybuilder.QueryFilter{
				Rule: querybuilder.CombinedRule{
					Condition: querybuilder.ConditionOr,
					Rules: []querybuilder.Rule{
						querybuilder.AtomRule{
							Field:    common.BKHostIDField,
							Operator: querybuilder.OperatorIn,
							Value:    instIDArr,
						},
					},
				},
			},
			Fields: dbFields,
		}

		resp, err := lgc.CoreAPI.ApiServer().ListHostWithoutApp(ctx, header, option)
		if err != nil {
			blog.ErrorJSON(" fetchInstAssocationData failed, ListHostWithoutApp err:%s, option: %s, rid: %s", err, option, rid)
			return nil, ccErr.CCError(common.CCErrCommHTTPDoRequestFailed)
		}
		if !resp.Result {
			blog.ErrorJSON(" fetchInstAssocationData failed, ListHostWithoutApp resp:%s, option: %s, rid: %s", resp, option, rid)
			return nil, resp.CCError()
		}

		for _, inst := range resp.Data.Info {
			insts = append(insts, mapstr.NewFromMap(inst))
		}
	default:
		option := mapstr.MapStr{
			"condition": mapstr.MapStr{
				instIDKey: mapstr.MapStr{
					common.BKDBIN: instIDArr,
				},
			},
			"fields": dbFields,
		}

		resp, err := lgc.CoreAPI.ApiServer().GetInstDetail(ctx, header, objID, option)
		if err != nil {
			blog.ErrorJSON(" fetchInstAssocationData failed, GetInstDetail err:%v, option: %s, rid: %s", err, option, rid)
			return nil, ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		if !resp.Result {
			blog.ErrorJSON(" fetchInstAssocationData failed, GetInstDetail resp:%s, option: %s, rid: %s", resp, option, rid)
			return nil, resp.CCError()
		}
		insts = resp.Data.Info
	}

	retAsstInstInfo := make(map[int64][]PropertyPrimaryVal, 0)
	for _, inst := range insts {
		instID, err := inst.Int64(instIDKey)
		if err != nil {
			blog.Warnf("FetchInstAssocationData get %s instance %s field error, err:%s, inst:%+v, rid:%s", objID, instIDKey, err.Error(), inst, rid)
			continue
		}
		isSkip := false
		var primaryKeysVal []PropertyPrimaryVal
		for _, key := range propertyArr {
			// use display , use string
			val, err := inst.String(key.ID)
			if err != nil {
				blog.Warnf("FetchInstAssocationData get %s instance %s field error, err:%s, inst:%+v, rid:%s", objID, key, err.Error(), inst, rid)
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
