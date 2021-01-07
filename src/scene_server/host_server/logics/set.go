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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	parse "configcenter/src/common/paraparse"
)

func (lgc *Logics) GetSetIDByCond(kit *rest.Kit, cond []metadata.ConditionItem) ([]int64, errors.CCError) {
	condc := make(map[string]interface{})
	if err := parse.ParseCommonParams(cond, condc); err != nil {
		blog.Errorf("ParseCommonParams failed, err: %+v, rid: %s", err, kit.Rid)
		return nil, err
	}

	query := &metadata.QueryCondition{
		Condition: mapstr.NewFromMap(condc),
		Fields:    []string{common.BKSetIDField},
		Page:      metadata.BasePage{Start: 0, Limit: common.BKNoLimit, Sort: common.BKSetIDField},
	}

	result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDSet, query)
	if err != nil {
		blog.Errorf("GetSetIDByCond http do error, err:%s,objID:%s,input:%+v,rid:%s", err.Error(), common.BKInnerObjIDSet, query, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetSetIDByCond http reponse error, err code:%d, err msg:%s,objID:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, common.BKInnerObjIDSet, query, kit.Rid)
		return nil, kit.CCError.New(result.Code, result.ErrMsg)
	}

	setIDArr := make([]int64, 0)
	for _, i := range result.Data.Info {
		setID, err := i.Int64(common.BKSetIDField)
		if err != nil {
			blog.Errorf("GetSetIDByCond convert %s %s to integer error, set info:%+v, input:%+v,rid:%s", common.BKInnerObjIDSet, common.BKSetIDField, i, query, kit.Rid)
			return nil, kit.CCError.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDSet, common.BKSetIDField, "int", err.Error())
		}
		setIDArr = append(setIDArr, setID)
	}
	return setIDArr, nil
}

// ExecuteSetDynamicGroup searches sets base on conditions without filling topology informations.
func (lgc *Logics) ExecuteSetDynamicGroup(kit *rest.Kit, setCommonSearch *metadata.SetCommonSearch,
	fields []string, disableCounter bool) (*metadata.InstDataInfo, errors.CCError) {

	// search parameters with condition.
	queryParams := &metadata.QueryCondition{Fields: fields, Page: setCommonSearch.Page, Condition: mapstr.New(), DisableCounter: disableCounter}

	// parse set search conditions.
	for _, searchCondition := range setCommonSearch.Condition {
		condc := make(map[string]interface{})
		if err := parse.ParseCommonParams(searchCondition.Condition, condc); err != nil {
			blog.Errorf("search set failed, can't parse condition, err: %+v, cond: %+v, rid: %s", err, searchCondition.Condition, kit.Rid)
			return nil, kit.CCError.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		// add field conditions to query params.
		for field, value := range condc {
			queryParams.Condition.Set(field, value)
		}
	}
	queryParams.Condition.Set(common.BKAppIDField, setCommonSearch.AppID)

	// search set with conditions.
	result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDSet, queryParams)
	if err != nil {
		blog.Errorf("search set failed, err: %+v, input: %+v, rid: %s", err, queryParams, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("search set failed, errcode: %d, errmsg: %s, input: %+v, rid: %s", result.Code, result.ErrMsg, queryParams, kit.Rid)
		return nil, kit.CCError.New(result.Code, result.ErrMsg)
	}
	return &result.Data, nil
}

func (lgc *Logics) GetSetMapByCond(kit *rest.Kit, fields []string, cond mapstr.MapStr) (map[int64]mapstr.MapStr, errors.CCError) {
	query := &metadata.QueryCondition{
		Condition: cond,
		Fields:    fields,
		Page:      metadata.BasePage{Sort: common.BKSetIDField},
	}

	result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDSet, query)
	if err != nil {
		blog.Errorf("GetSetMapByCond http do error, err:%s,objID:%s,input:%+v,rid:%s", err.Error(), common.BKInnerObjIDSet, query, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetSetMapByCond http reponse error, err code:%d, err msg:%s,objID:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, common.BKInnerObjIDSet, query, kit.Rid)
		return nil, kit.CCError.New(result.Code, result.ErrMsg)
	}

	setMap := make(map[int64]mapstr.MapStr)
	for _, i := range result.Data.Info {
		setID, err := i.Int64(common.BKSetIDField)
		if err != nil {
			blog.Errorf("GetSetMapByCond convert %s %s to integer error, set info:%+v, input:%+v,rid:%s", common.BKInnerObjIDSet, common.BKSetIDField, i, query, kit.Rid)
			return nil, kit.CCError.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDSet, common.BKSetIDField, "int", err.Error())
		}

		setMap[setID] = i
	}
	return setMap, nil
}

// GetSetIDsByTopo get set IDs by custom layer node
func (lgc *Logics) GetSetIDsByTopo(kit *rest.Kit, objID string, instID int64) ([]int64, error) {
	if objID == common.BKInnerObjIDApp || objID == common.BKInnerObjIDSet || objID == common.BKInnerObjIDModule {
		blog.Errorf("get set IDs by topo failed, obj(%s) is a inner object, rid: %s", objID, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKObjIDField)
	}

	// get mainline association, generate map of object and its child
	asstRes, err := lgc.CoreAPI.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header, &metadata.QueryCondition{
		Condition: map[string]interface{}{common.AssociationKindIDField: common.AssociationKindMainline}})
	if err != nil {
		blog.Errorf("get set IDs by topo failed, get mainline association err: %s, rid: %s", err.Error(), kit.Rid)
		return nil, err
	}

	if !asstRes.Result {
		blog.Errorf("get set IDs by topo failed, get mainline association err: %s, rid: %s", asstRes.ErrMsg, kit.Rid)
		return nil, asstRes.CCError()
	}

	childObjMap := make(map[string]string)
	for _, asst := range asstRes.Data.Info {
		childObjMap[asst.AsstObjID] = asst.ObjectID
	}

	childObj := childObjMap[objID]
	if childObj == "" {
		blog.Errorf("get set IDs by topo failed, obj(%s) is not a mainline object, rid: %s", objID, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKObjIDField)
	}

	// traverse down topo till set, get set ids
	instIDs := []int64{instID}
	for {
		idField := common.GetInstIDField(childObj)
		query := &metadata.QueryCondition{
			Condition: map[string]interface{}{common.BKParentIDField: map[string]interface{}{common.BKDBIN: instIDs}},
			Fields:    []string{idField},
			Page:      metadata.BasePage{Limit: common.BKNoLimit},
		}

		instRes, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, childObj, query)
		if err != nil {
			blog.Errorf("get set IDs by topo failed, read instance err: %s, objID: %s, instIDs: %+v, rid: %s", err.Error(), childObj, instIDs, kit.Rid)
			return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if !instRes.Result {
			blog.Errorf("get set IDs by topo failed, read instance err: %s, objID: %s, instIDs: %+v, rid: %s", instRes.ErrMsg, childObj, instIDs, kit.Rid)
			return nil, kit.CCError.New(instRes.Code, instRes.ErrMsg)
		}

		if len(instRes.Data.Info) == 0 {
			return []int64{}, nil
		}

		instIDs = make([]int64, len(instRes.Data.Info))
		for index, inst := range instRes.Data.Info {
			id, err := inst.Int64(idField)
			if err != nil {
				blog.Errorf("get set IDs by topo failed, parse inst id err: %s, inst: %#v, rid: %s", err.Error(), inst, kit.Rid)
				return nil, err
			}
			instIDs[index] = id
		}

		if childObj == common.BKInnerObjIDSet {
			break
		}
		childObj = childObjMap[childObj]
	}

	return instIDs, nil
}
