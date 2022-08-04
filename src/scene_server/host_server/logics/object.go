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
		blog.Errorf("search object attributes failed, err: %+v, objID: %s, input: %+v, rid: %s", err, objectID,
			query, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	return result.Info, nil
}

// GetTopoIDByName TODO
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
	setIDs, setErr := lgc.GetSetIDByCond(kit, meta.ConditionWithTime{Condition: setCond})
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
	moduleIDs, moduleErr := lgc.GetModuleIDByCond(kit, meta.ConditionWithTime{Condition: moduleCond})
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

// GetSetIDByObjectCond get set ids by mainline node conditions
func (lgc *Logics) GetSetIDByObjectCond(kit *rest.Kit, appID int64, objectCond []meta.ConditionItem) ([]int64,
	errors.CCError) {

	objectIDArr := make([]int64, 0)

	// parse mainline object condition to get inst ids filter, only allows condition of 'bk_inst_id $eq value' form
	var hasInstID bool
	for _, i := range objectCond {
		if i.Field != common.BKInstIDField {
			continue
		}
		if i.Operator != common.BKDBEQ {
			continue
		}

		value, err := util.GetInt64ByInterface(i.Value)
		if err != nil {
			return nil, err
		}
		objectIDArr = append(objectIDArr, value)

		hasInstID = true
	}

	if !hasInstID {
		blog.Errorf("mainline miss bk_inst_id parameters. input:%#v, rid:%s", objectCond, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrHostSearchNeedObjectInstIDErr)
	}

	// get inst ids corresponding object to ids map and mainline child to parent map
	instObjMap, err := lgc.CoreAPI.CoreService().Instance().GetInstanceObjectMapping(kit.Ctx, kit.Header, objectIDArr)
	if err != nil {
		blog.Errorf("get instance mappings %v failed, err: %v, rid: %s", objectIDArr, err, kit.Rid)
		return nil, err
	}

	objInstIDMap := make(map[string]int64)
	for _, mapping := range instObjMap {
		// returns no set ids if more than one inst id equal condition is set for one object
		if _, exists := objInstIDMap[mapping.ObjectID]; exists {
			return make([]int64, 0), nil
		}
		objInstIDMap[mapping.ObjectID] = mapping.ID
	}

	mainlineMap, ccErr := lgc.searchMainlineRelationMap(kit)
	if err != nil {
		blog.Errorf("get mainline association failed, err: %v, rid: %s", err, kit.Rid)
		return nil, ccErr
	}

	// loop from the first mainline object under biz to set, filters out the set ids under the mainline instances
	filter := make(map[string]interface{})
	for object := mainlineMap[common.BKInnerObjIDApp]; ; object = mainlineMap[object] {
		instID, exists := objInstIDMap[object]
		if len(filter) == 0 && !exists {
			continue
		}

		filter[common.BKAppIDField] = appID

		if exists {
			filter[common.BKInstIDField] = instID
		}

		filteredIDs, err := lgc.getInstIDsByCond(kit, object, filter)
		if err != nil {
			blog.Errorf("get object[%s] inst ids failed, cond: %v, err: %v, rid: %s", object, filter, err, kit.Rid)
			return nil, err
		}

		if len(filteredIDs) == 0 {
			return make([]int64, 0), nil
		}

		if object == common.BKInnerObjIDSet {
			return filteredIDs, nil
		}

		filter = map[string]interface{}{
			common.BKParentIDField: mapstr.MapStr{common.BKDBIN: filteredIDs},
		}
	}
}

// getObjectByParentID TODO
// deprecated, please use CoreAPI.CoreService().Mainline().SearchMainlineInstanceTopo instead
func (lgc *Logics) getObjectByParentID(kit *rest.Kit, valArr []int64) ([]int64, errors.CCError) {
	instIDArr := make([]int64, 0)
	condCell, sCond := mapstr.New(), mapstr.New()
	condCell.Set(common.BKDBIN, valArr)
	sCond.Set(common.BKInstParentStr, condCell)

	query := &meta.QueryCondition{
		Condition: sCond,
	}
	// TODO common.BKInnerObjIDObject is not a valid value to search mainline topo instance,
	// it will act as bk_obj_id=object condition
	result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDObject,
		query)
	if err != nil {
		blog.Errorf("getObjectByParentID http do error, err:%s,objID:%s,input:%+v,rid:%s", err.Error(),
			common.BKInnerObjIDObject, query, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if result.Count == 0 {
		return instIDArr, nil
	}

	for _, info := range result.Info {
		id, err := info.Int64(common.BKInstIDField)
		if err != nil {
			blog.Errorf("getObjectByParentID failed, get int64 `bk_inst_id` field failed, instance: %+v, input: %+v, "+
				"err: %+v, rid:%s", info, query, err, kit.Rid)
			return nil, kit.CCError.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDObject,
				common.BKInstIDField, "int", err.Error())
		}
		instIDArr = append(instIDArr, id)
	}

	return instIDArr, nil
}

// GetObjectInstByCond search object instance by condition
func (lgc *Logics) GetObjectInstByCond(kit *rest.Kit, objID string, cond []meta.ConditionItem) ([]int64,
	errors.CCError) {
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
		blog.Errorf("GetObjectInstByCond http do error, err:%s,objID:%s,input:%+v,rid:%s", err.Error(), objID, query,
			kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if result.Count == 0 {
		return instIDArr, nil
	}

	for _, info := range result.Info {
		id, err := info.Int64(outField)
		if err != nil {
			blog.Errorf("getObjectByParentID convert %s %s to integer error, inst info:%+v, input:%+v,rid:%s", objID,
				outField, info, query, kit.Rid)
			return nil, kit.CCError.Errorf(common.CCErrCommInstFieldConvertFail, objID, outField, "int", err.Error())
		}
		instIDArr = append(instIDArr, id)
	}

	return instIDArr, nil
}

// GetHostIDByInstID TODO
func (lgc *Logics) GetHostIDByInstID(kit *rest.Kit, asstObjId string, instIDArr []int64) ([]int64, errors.CCError) {
	cond := hutil.NewOperation().WithObjID(common.BKInnerObjIDHost).
		WithAssoObjID(asstObjId).WithAssoInstID(map[string]interface{}{common.BKDBIN: instIDArr}).Data()

	query := &meta.InstAsstQueryCondition{
		Cond:  meta.QueryCondition{Condition: cond},
		ObjID: common.BKInnerObjIDHost,
	}
	result, err := lgc.CoreAPI.CoreService().Association().ReadInstAssociation(kit.Ctx, kit.Header, query)
	if err != nil {
		blog.Errorf("GetHostIDByInstID http do error, err:%s,objID:%s,input:%+v,rid:%s", err.Error(),
			common.BKTableNameInstAsst, query, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	hostIDs := make([]int64, 0)
	for _, val := range result.Info {
		hostIDs = append(hostIDs, val.InstID)
	}

	return hostIDs, nil
}
