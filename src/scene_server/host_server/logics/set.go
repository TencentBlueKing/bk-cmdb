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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	parse "configcenter/src/common/paraparse"
)

func (lgc *Logics) GetSetIDByCond(ctx context.Context, cond []metadata.ConditionItem) ([]int64, errors.CCError) {
	condc := make(map[string]interface{})
	if err := parse.ParseCommonParams(cond, condc); err != nil {
		blog.Warnf("ParseCommonParams failed, err: %+v, rid: %s", err, lgc.rid)
	}

	query := &metadata.QueryCondition{
		Condition: mapstr.NewFromMap(condc),
		Fields:    []string{common.BKSetIDField},
		Page:      metadata.BasePage{Start: 0, Limit: common.BKNoLimit, Sort: common.BKSetIDField},
	}

	result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(ctx, lgc.header, common.BKInnerObjIDSet, query)
	if err != nil {
		blog.Errorf("GetSetIDByCond http do error, err:%s,objID:%s,input:%+v,rid:%s", err.Error(), common.BKInnerObjIDSet, query, lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetSetIDByCond http reponse error, err code:%d, err msg:%s,objID:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, common.BKInnerObjIDSet, query, lgc.rid)
		return nil, lgc.ccErr.New(result.Code, result.ErrMsg)
	}

	setIDArr := make([]int64, 0)
	for _, i := range result.Data.Info {
		setID, err := i.Int64(common.BKSetIDField)
		if err != nil {
			blog.Errorf("GetSetIDByCond convert %s %s to integer error, set info:%+v, input:%+v,rid:%s", common.BKInnerObjIDSet, common.BKSetIDField, i, query, lgc.rid)
			return nil, lgc.ccErr.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDSet, common.BKSetIDField, "int", err.Error())
		}
		setIDArr = append(setIDArr, setID)
	}
	return setIDArr, nil
}

func (lgc *Logics) GetSetMapByCond(ctx context.Context, fields []string, cond mapstr.MapStr) (map[int64]mapstr.MapStr, errors.CCError) {
	query := &metadata.QueryCondition{
		Condition: cond,
		Fields:    fields,
		Page:      metadata.BasePage{Sort: common.BKSetIDField},
	}

	result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(ctx, lgc.header, common.BKInnerObjIDSet, query)
	if err != nil {
		blog.Errorf("GetSetMapByCond http do error, err:%s,objID:%s,input:%+v,rid:%s", err.Error(), common.BKInnerObjIDSet, query, lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetSetMapByCond http reponse error, err code:%d, err msg:%s,objID:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, common.BKInnerObjIDSet, query, lgc.rid)
		return nil, lgc.ccErr.New(result.Code, result.ErrMsg)
	}

	setMap := make(map[int64]mapstr.MapStr)
	for _, i := range result.Data.Info {
		setID, err := i.Int64(common.BKSetIDField)
		if err != nil {
			blog.Errorf("GetSetMapByCond convert %s %s to integer error, set info:%+v, input:%+v,rid:%s", common.BKInnerObjIDSet, common.BKSetIDField, i, query, lgc.rid)
			return nil, lgc.ccErr.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDSet, common.BKSetIDField, "int", err.Error())
		}

		setMap[setID] = i
	}
	return setMap, nil
}

// GetSetIDsByTopo get set IDs by custom layer node
func (lgc *Logics) GetSetIDsByTopo(ctx context.Context, objID string, instID int64) ([]int64, error) {
	if objID == common.BKInnerObjIDApp || objID == common.BKInnerObjIDSet || objID == common.BKInnerObjIDModule {
		blog.Errorf("get set IDs by topo failed, obj(%s) is a inner object, rid: %s", objID, lgc.rid)
		return nil, lgc.ccErr.CCErrorf(common.CCErrCommParamsInvalid, common.BKObjIDField)
	}

	// get mainline association, generate map of object and its child
	asstRes, err := lgc.CoreAPI.CoreService().Association().ReadModelAssociation(ctx, lgc.header, &metadata.QueryCondition{
		Condition: map[string]interface{}{common.AssociationKindIDField: common.AssociationKindMainline}})
	if err != nil {
		blog.Errorf("get set IDs by topo failed, get mainline association err: %s, rid: %s", err.Error(), lgc.rid)
		return nil, err
	}

	if !asstRes.Result {
		blog.Errorf("get set IDs by topo failed, get mainline association err: %s, rid: %s", asstRes.ErrMsg, lgc.rid)
		return nil, asstRes.CCError()
	}

	childObjMap := make(map[string]string)
	for _, asst := range asstRes.Data.Info {
		childObjMap[asst.AsstObjID] = asst.ObjectID
	}

	childObj := childObjMap[objID]
	if childObj == "" {
		blog.Errorf("get set IDs by topo failed, obj(%s) is not a mainline object, rid: %s", objID, lgc.rid)
		return nil, lgc.ccErr.CCErrorf(common.CCErrCommParamsInvalid, common.BKObjIDField)
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

		instRes, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(ctx, lgc.header, childObj, query)
		if err != nil {
			blog.Errorf("get set IDs by topo failed, read instance err: %s, objID: %s, instIDs: %+v, rid: %s", err.Error(), childObj, instIDs, lgc.rid)
			return nil, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if !instRes.Result {
			blog.Errorf("get set IDs by topo failed, read instance err: %s, objID: %s, instIDs: %+v, rid: %s", instRes.ErrMsg, childObj, instIDs, lgc.rid)
			return nil, lgc.ccErr.New(instRes.Code, instRes.ErrMsg)
		}

		if len(instRes.Data.Info) == 0 {
			return []int64{}, nil
		}

		instIDs = make([]int64, len(instRes.Data.Info))
		for index, inst := range instRes.Data.Info {
			id, err := inst.Int64(idField)
			if err != nil {
				blog.Errorf("get set IDs by topo failed, parse inst id err: %s, inst: %#v, rid: %s", err.Error(), inst, lgc.rid)
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
