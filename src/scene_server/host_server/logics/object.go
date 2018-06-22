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
	"fmt"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	parse "configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	hutil "configcenter/src/scene_server/host_server/util"
)

// get the object attributes
func (lgc *Logics) GetObjectAttributes(ownerID, objID string, pheader http.Header, page meta.BasePage) ([]meta.Attribute, error) {
	opt := hutil.NewOperation().WithOwnerID(ownerID).WithPage(page).WithObjID(objID).Data()
	result, err := lgc.CoreAPI.ObjectController().Meta().SelectObjectAttWithParams(context.Background(), pheader, opt)
	if err != nil || (err == nil && !result.Result) {
		return nil, fmt.Errorf("get object attribute failed, err: %v, %v", err, result.ErrMsg)
	}

	return result.Data, nil
}

func (lgc *Logics) GetTopoIDByName(pheader http.Header, c *meta.HostToAppModule) (int64, int64, int64, error) {
	if "" == c.AppName || "" == c.SetName || "" == c.ModuleName {
		return 0, 0, 0, nil
	}

	appInfo, appErr := lgc.GetSingleApp(pheader, common.KvMap{common.BKAppNameField: c.AppName, common.BKOwnerIDField: c.OwnerID})
	if nil != appErr {
		blog.Errorf("getTopoIDByName get app info error; %s", appErr.Error())
		return 0, 0, 0, lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader)).Error(common.CCErrCommHTTPDoRequestFailed)
	}

	appID, err := appInfo.Int64(common.BKAppIDField)
	if err != nil {
		return 0, 0, 0, lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader)).Error(common.CCErrCommParamsInvalid)
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
	setIDs, setErr := lgc.GetSetIDByCond(pheader, setCond)
	if nil != setErr {
		blog.Errorf("getTopoIDByName get app info error; %s", setErr.Error())
		return 0, 0, 0, lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader)).Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if 0 == len(setIDs) || 0 >= setIDs[0] {
		blog.Info("getTopoIDByName get set info not found; applicationName: %s, setName: %s", c.AppName, c.SetName)
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
	moduleIDs, moduleErr := lgc.GetModuleIDByCond(pheader, moduleCond)
	if nil != moduleErr {
		blog.Errorf("getTopoIDByName get app info error; %s", setErr.Error())
		return 0, 0, 0, lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader)).Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if 0 == len(moduleIDs) || 0 >= moduleIDs[0] {
		blog.Info("getTopoIDByName get module info not found; applicationName: %s, setName: %s, moduleName: %s", c.AppName, c.SetName, c.ModuleName)
		return 0, 0, 0, nil
	}
	moduleID := moduleIDs[0]

	return appID, setID, moduleID, nil
}

func (lgc *Logics) GetSetIDByObjectCond(pheader http.Header, appID int64, objectCond []meta.ConditionItem) ([]int64, error) {
	objectIDArr := make([]int64, 0)
	condition := make([]meta.ConditionItem, 0)

	instItem := meta.ConditionItem{}
	for _, i := range objectCond {
		if i.Field != common.BKInstIDField {
			continue
		}
		value, err := util.GetInt64ByInterface(i.Value)
		if nil != err {
			return nil, err
		}
		instItem.Field = common.BKInstParentStr
		instItem.Operator = i.Operator
		instItem.Value = i.Value
		objectIDArr = append(objectIDArr, value)
	}
	condition = append(condition, instItem)

	appIDItem := meta.ConditionItem{
		Field:    common.BKAppIDField,
		Operator: common.BKDBEQ,
		Value:    appID,
	}
	condition = append(condition, appIDItem)

	for {
		sSetIDArr, err := lgc.GetSetIDByCond(pheader, condition)
		if err != nil {
			return nil, err
		}

		if 0 != len(sSetIDArr) {
			return sSetIDArr, nil
		}

		sObjectIDArr, err := lgc.getObjectByParentID(pheader, objectIDArr)
		if err != nil {
			return nil, err
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
	}

}

func (lgc *Logics) getObjectByParentID(pheader http.Header, valArr []int64) ([]int64, error) {
	instIDArr := make([]int64, 0)
	condCell, sCond := make(map[string]interface{}), make(map[string]interface{})
	condCell[common.BKDBIN] = valArr
	sCond[common.BKInstParentStr] = condCell

	query := &meta.QueryInput{
		Condition: sCond,
		Start:     0,
		Limit:     common.BKNoLimit,
	}
	result, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKINnerObjIDObject, pheader, query)
	if err != nil || (err == nil && !result.Result) {
		return nil, fmt.Errorf("get object failed, err: %v, %v", err, result.ErrMsg)
	}

	if result.Data.Count == 0 {
		return instIDArr, nil
	}

	for _, info := range result.Data.Info {
		id, err := info.Int64(common.BKInstIDField)
		if err != nil {
			return nil, fmt.Errorf("invalid obj id: %v", err)
		}
		instIDArr = append(instIDArr, id)
	}

	return instIDArr, nil
}

func (lgc *Logics) GetObjectInstByCond(pheader http.Header, objID string, cond []meta.ConditionItem) ([]int64, error) {
	instIDArr := make([]int64, 0)
	condc := make(map[string]interface{})
	parse.ParseCommonParams(cond, condc)

	var outField, objType string
	if objID == common.BKInnerObjIDPlat {
		outField = common.BKCloudIDField
		objType = common.BKInnerObjIDPlat
	} else {
		condc[common.BKObjIDField] = objID
		outField = common.BKInstIDField
		objType = common.BKINnerObjIDObject
	}

	query := &meta.QueryInput{
		Condition: condc,
		Start:     0,
		Limit:     1,
		Sort:      common.BKAppIDField,
	}
	result, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), objType, pheader, query)
	if err != nil || (err == nil && !result.Result) {
		return nil, fmt.Errorf("%v, %v", err, result.ErrMsg)
	}

	if result.Data.Count == 0 {
		return instIDArr, nil
	}

	for _, info := range result.Data.Info {
		id, err := info.Int64(outField)
		if err != nil {
			return nil, err
		}
		instIDArr = append(instIDArr, id)
	}

	return instIDArr, nil
}
func (lgc *Logics) GetHostIDByInstID(pheader http.Header, asstObjId string, instIDArr []int64) ([]int64, error) {
	cond := hutil.NewOperation().WithObjID(common.BKInnerObjIDHost).
		WithAssoObjID(asstObjId).WithAssoInstID(map[string]interface{}{common.BKDBIN: instIDArr}).Data()

	query := &meta.QueryInput{
		Condition: cond,
		Start:     0,
		Limit:     common.BKNoLimit,
	}
	result, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKTableNameInstAsst, pheader, query)
	if err != nil || (err == nil && !result.Result) {
		return nil, fmt.Errorf("%v, %v", err, result.ErrMsg)
	}

	hostIDs := make([]int64, 0)
	for _, val := range result.Data.Info {
		id, err := val.Int64(common.BKInstIDField)
		if err != nil {
			return nil, err
		}
		hostIDs = append(hostIDs, id)
	}

	return hostIDs, nil
}
