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
	meta "configcenter/src/common/metadata"
	parse "configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	hutil "configcenter/src/scene_server/host_server/util"
)

// SearchObjectAttributes returns attributes of target object.
func (lgc *Logics) SearchObjectAttributes(kit *rest.Kit, bizID int64, objectID string) ([]meta.Attribute, error) {
	query := &meta.QueryCondition{
		Condition: map[string]interface{}{
			common.BKDBOR: []map[string]interface{}{
				{
					common.BKObjIDField: objectID,
					common.BKAppIDField: 0,
				},
				{
					common.BKObjIDField: objectID,
					common.BKAppIDField: bizID,
				},
			},
		},
	}

	result, err := lgc.CoreAPI.CoreService().Model().ReadModelAttr(kit.Ctx, kit.Header, objectID, query)
	if err != nil {
		blog.Errorf("search object attributes failed, err: %+v, objID: %s, input: %+v, rid: %s", err, objectID, query, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("search object attributes failed, errcode: %d, errmsg: %s, objID: %s, input: %+v, rid: %s",
			result.Code, result.ErrMsg, objectID, query, kit.Rid)
		return nil, kit.CCError.New(result.Code, result.ErrMsg)
	}
	return result.Data.Info, nil
}

func (lgc *Logics) GetTopoIDByName(kit *rest.Kit, c *meta.HostToAppModule) (int64, int64, int64, errors.CCError) {
	if "" == c.AppName || "" == c.SetName || "" == c.ModuleName {
		return 0, 0, 0, nil
	}

	appInfo, appErr := lgc.GetSingleApp(kit, mapstr.MapStr{common.BKAppNameField: c.AppName})
	if nil != appErr {
		return 0, 0, 0, appErr
	}

	appID, err := appInfo.Int64(common.BKAppIDField)
	if err != nil {
		blog.Errorf("GetTopoIDByName convert %s %s to integer error, app info:%+v, input:%+v,rid:%s", common.BKInnerObjIDApp, common.BKAppIDField, appInfo, c, kit.Rid)
		return 0, 0, 0, kit.CCError.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDApp, common.BKAppIDField, "int", err.Error())
	}

	appIdItem := meta.ConditionItem{
		Field:    common.BKAppIDField,
		Operator: common.BKDBEQ,
		Value:    appID,
	}
	setNameItem := meta.ConditionItem{
		Field:    common.BKSetNameField,
		Operator: common.BKDBEQ,
		Value:    c.SetName,
	}

	setCond := []meta.ConditionItem{appIdItem, setNameItem}
	setIDs, setErr := lgc.GetSetIDByCond(kit, setCond)
	if nil != setErr {
		return 0, 0, 0, setErr
	}
	if 0 == len(setIDs) || 0 >= setIDs[0] {
		blog.V(5).Infof("getTopoIDByName get set info not found; applicationName: %s, setName: %s, rid:%s", c.AppName, c.SetName, kit.Rid)
		return 0, 0, 0, nil
	}
	setID := setIDs[0]

	setIDConds := meta.ConditionItem{
		Field:    common.BKSetIDField,
		Operator: common.BKDBEQ,
		Value:    setID,
	}

	moduleNameCond := meta.ConditionItem{
		Field:    common.BKModuleNameField,
		Operator: common.BKDBEQ,
		Value:    c.ModuleName,
	}

	moduleCond := []meta.ConditionItem{appIdItem, setIDConds, moduleNameCond}
	moduleIDs, moduleErr := lgc.GetModuleIDByCond(kit, moduleCond)
	if nil != moduleErr {
		return 0, 0, 0, err
	}
	if 0 == len(moduleIDs) || 0 >= moduleIDs[0] {
		blog.V(5).Infof("getTopoIDByName get module info not found; applicationName: %s, setName: %s, moduleName: %s,rid:%s", c.AppName, c.SetName, c.ModuleName, kit.Rid)
		return 0, 0, 0, nil
	}
	moduleID := moduleIDs[0]

	return appID, setID, moduleID, nil
}

func (lgc *Logics) GetSetIDByObjectCond(kit *rest.Kit, appID int64, objectCond []meta.ConditionItem) ([]int64, errors.CCError) {
	objectIDArr := make([]int64, 0)
	condition := make([]meta.ConditionItem, 0)

	instItem := meta.ConditionItem{}
	var hasInstID bool
	for _, i := range objectCond {
		if i.Field != common.BKInstIDField {
			continue
		}
		value, err := util.GetInt64ByInterface(i.Value)
		if nil != err {
			return nil, err
		}
		hasInstID = true
		instItem.Field = common.BKInstParentStr
		instItem.Operator = i.Operator
		instItem.Value = i.Value

		objectIDArr = append(objectIDArr, value)
	}
	condition = append(condition, instItem)
	if !hasInstID {
		blog.Errorf("mainline miss bk_inst_id parameters. input:%#v, rid:%s", objectCond, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrHostSearchNeedObjectInstIDErr)
	}

	nodeFaultItem := meta.ConditionItem{}
	nodeFaultItem.Field = common.BKDefaultField
	nodeFaultItem.Operator = common.BKDBNE
	nodeFaultItem.Value = common.DefaultResSetFlag

	appIDItem := meta.ConditionItem{
		Field:    common.BKAppIDField,
		Operator: common.BKDBEQ,
		Value:    appID,
	}
	condition = append(condition, appIDItem)
	condition = append(condition, nodeFaultItem)

	topoRoot, err := lgc.CoreAPI.CoreService().Mainline().SearchMainlineInstanceTopo(kit.Ctx, kit.Header, appID, false)
	if err != nil {
		return nil, kit.CCError.Error(common.CCErrTopoMainlineSelectFailed)
	}

	for {
		sSetIDArr, err := lgc.GetSetIDByCond(kit, condition)
		if err != nil {
			return nil, err
		}

		if 0 != len(sSetIDArr) {
			return sSetIDArr, nil
		}

		sObjectIDArr := make([]int64, 0)
		for _, id := range objectIDArr {
			path := topoRoot.TraversalFindNode(common.BKInnerObjIDObject, id)
			if len(path) == 0 {
				continue
			}
			node := path[0]
			for _, childNode := range node.Children {
				sObjectIDArr = append(sObjectIDArr, childNode.InstanceID)
			}
		}
		objectIDArr = sObjectIDArr
		if 0 == len(sObjectIDArr) {
			return []int64{}, nil
		}

		conc := meta.ConditionItem{
			Field:    common.BKInstParentStr,
			Operator: common.BKDBIN,
			Value:    sObjectIDArr,
		}
		condition = make([]meta.ConditionItem, 0)
		condition = append(condition, conc)
		condition = append(condition, appIDItem)
		condition = append(condition, nodeFaultItem)
	}

}

// deprecated, please use CoreAPI.CoreService().Mainline().SearchMainlineInstanceTopo instead
func (lgc *Logics) getObjectByParentID(kit *rest.Kit, valArr []int64) ([]int64, errors.CCError) {
	instIDArr := make([]int64, 0)
	condCell, sCond := mapstr.New(), mapstr.New()
	condCell.Set(common.BKDBIN, valArr)
	sCond.Set(common.BKInstParentStr, condCell)

	query := &meta.QueryCondition{
		Condition: sCond,
	}
	// TODO common.BKInnerObjIDObject is not a valid value to search mainline topo instance, it will act as bk_obj_id=object condition
	result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDObject, query)
	if err != nil {
		blog.Errorf("getObjectByParentID http do error, err:%s,objID:%s,input:%+v,rid:%s", err.Error(), common.BKInnerObjIDObject, query, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("getObjectByParentID http response error, err code:%d, err msg:%s,objID:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, common.BKInnerObjIDObject, query, kit.Rid)
		return nil, kit.CCError.New(result.Code, result.ErrMsg)
	}

	if result.Data.Count == 0 {
		return instIDArr, nil
	}

	for _, info := range result.Data.Info {
		id, err := info.Int64(common.BKInstIDField)
		if err != nil {
			blog.Errorf("getObjectByParentID failed, get int64 `bk_inst_id` field failed, instance: %+v, input: %+v, err: %+v, rid:%s", info, query, err, kit.Rid)
			return nil, kit.CCError.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDObject, common.BKInstIDField, "int", err.Error())
		}
		instIDArr = append(instIDArr, id)
	}

	return instIDArr, nil
}

func (lgc *Logics) GetObjectInstByCond(kit *rest.Kit, objID string, cond []meta.ConditionItem) ([]int64, errors.CCError) {
	instIDArr := make([]int64, 0)
	condc := make(map[string]interface{})
	if err := parse.ParseCommonParams(cond, condc); err != nil {
		blog.Errorf("GetObjectInstByCond failed, ParseCommonParams failed, err: %+v, rid: %s", err, kit.Rid)
		return nil, err
	}

	var outField, objType string
	if objID == common.BKInnerObjIDPlat {
		outField = common.BKCloudIDField
		objType = common.BKInnerObjIDPlat
	} else {
		condc[common.BKObjIDField] = objID
		outField = common.BKInstIDField
		objType = common.BKInnerObjIDObject
	}

	query := &meta.QueryCondition{
		Condition: mapstr.NewFromMap(condc),
		Page:      meta.BasePage{Sort: common.BKAppIDField},
	}
	result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, objType, query)
	if err != nil {
		blog.Errorf("GetObjectInstByCond http do error, err:%s,objID:%s,input:%+v,rid:%s", err.Error(), objID, query, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetObjectInstByCond http response error, err code:%d, err msg:%s,objID:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, objID, query, kit.Rid)
		return nil, kit.CCError.New(result.Code, result.ErrMsg)
	}

	if result.Data.Count == 0 {
		return instIDArr, nil
	}

	for _, info := range result.Data.Info {
		id, err := info.Int64(outField)
		if err != nil {
			blog.Errorf("getObjectByParentID convert %s %s to integer error, inst info:%+v, input:%+v,rid:%s", objID, outField, info, query, kit.Rid)
			return nil, kit.CCError.Errorf(common.CCErrCommInstFieldConvertFail, objID, outField, "int", err.Error())
		}
		instIDArr = append(instIDArr, id)
	}

	return instIDArr, nil
}

func (lgc *Logics) GetHostIDByInstID(kit *rest.Kit, asstObjId string, instIDArr []int64) ([]int64, errors.CCError) {
	cond := hutil.NewOperation().WithObjID(common.BKInnerObjIDHost).
		WithAssoObjID(asstObjId).WithAssoInstID(map[string]interface{}{common.BKDBIN: instIDArr}).Data()

	query := &meta.QueryCondition{
		Condition: cond,
	}
	result, err := lgc.CoreAPI.CoreService().Association().ReadInstAssociation(kit.Ctx, kit.Header, query)
	if err != nil {
		blog.Errorf("GetHostIDByInstID http do error, err:%s,objID:%s,input:%+v,rid:%s", err.Error(), common.BKTableNameInstAsst, query, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetHostIDByInstID http response error, err code:%d, err msg:%s,objID:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, common.BKTableNameInstAsst, query, kit.Rid)
		return nil, kit.CCError.New(result.Code, result.ErrMsg)
	}

	hostIDs := make([]int64, 0)
	for _, val := range result.Data.Info {
		hostIDs = append(hostIDs, val.InstID)
	}

	return hostIDs, nil
}
