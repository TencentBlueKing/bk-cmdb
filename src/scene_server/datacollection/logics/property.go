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
	"configcenter/src/common/util"
)

// AddProperty create new net property
func (lgc *Logics) AddProperty(
	pHeader http.Header, propertyInfo meta.NetcollectProperty) (meta.AddNetPropertyResult, error) {

	netPropertyID, err := lgc.addProperty(pHeader, propertyInfo, util.GetOwnerID(pHeader))
	if nil != err {
		return meta.AddNetPropertyResult{NetcollectPropertyID: INVALIDID}, err
	}

	return meta.AddNetPropertyResult{NetcollectPropertyID: netPropertyID}, nil
}

func (lgc *Logics) UpdateProperty(pHeader http.Header, netPropertyID uint64, netPropertyInfo meta.NetcollectProperty) error {

	return lgc.updateProperty(pHeader, netPropertyInfo, netPropertyID, util.GetOwnerID(pHeader))
}

// BatchCreateProperty bacth create or update net propertys
func (lgc *Logics) BatchCreateProperty(
	pHeader http.Header, propertyInfoList []meta.NetcollectProperty) ([]meta.BatchAddNetPropertyResult, bool) {
	ownerID := util.GetOwnerID(pHeader)

	resultList := make([]meta.BatchAddNetPropertyResult, 0)
	hasError := false

	for _, propertyInfo := range propertyInfoList {
		errMsg := ""
		result := true

		propertyID, err := lgc.addOrUpdateProperty(pHeader, propertyInfo, ownerID)
		if nil != err {
			errMsg = err.Error()
			result = false
			hasError = true
		}

		resultList = append(resultList, meta.BatchAddNetPropertyResult{
			Result:               result,
			ErrMsg:               errMsg,
			NetcollectPropertyID: propertyID,
		})
	}

	return resultList, hasError
}

// SearchProperty get net devices by conditions
func (lgc *Logics) SearchProperty(pHeader http.Header, params *meta.NetCollSearchParams) (*meta.SearchNetProperty, error) {
	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pHeader))
	rid := util.GetHTTPCCRequestID(pHeader)

	// classify condition
	deviceCond, objectCond, propertyCond, netPropertyCond := lgc.classifyNetPropertyCondition(params.Condition)

	searchResult := meta.SearchNetProperty{Count: 0, Info: []meta.NetcollectProperty{}}

	var (
		err                error
		objIDs             []string
		deviceIDs          []uint64
		propertyIDs        []string
		showFields         netPropertyShowFields // to be displayed field of netProperty that be got from other tables
		objIDMapShowFields map[string]objShowField
	)
	// if property has filter condition
	if 0 < len(propertyCond) {
		// get propertyID and value of fields to be shown by property condition
		objIDs, propertyIDs, showFields.propertyIDMapShowFields, err = lgc.getPropertyIDsAndShowFields(pHeader, propertyCond)
		if nil != err {
			blog.Errorf("[NetProperty] search net property, get property fail, error: %v, condition [%#v], rid: %s", err, propertyCond, rid)
			return nil, err
		}

		// if find any propertyIDs matched condition, will must not find any property propetry
		if 0 == len(propertyIDs) || 0 == len(objIDs) {
			return &searchResult, nil
		}

		// propertyIDs as filter conditoin of net property
		netPropertyCond[common.BKPropertyIDField] = map[string]interface{}{common.BKDBIN: propertyIDs}
		objectCond[common.BKObjIDField] = map[string]interface{}{common.BKDBIN: objIDs}
	}

	// if obj has filter condition
	if 0 < len(objectCond) {
		// get objID and value of fields to be shown by obj condition
		objIDs, objIDMapShowFields, err = lgc.getObjIDsAndShowFields(pHeader, objectCond)
		if nil != err {
			blog.Errorf("[NetProperty] search net property, get net object fail, error: %v, condition [%#v], rid: %s", err, objectCond, rid)
			return nil, err
		}

		// if not find any objID matched condition, will not find any device propetry
		if 0 == len(objIDs) {
			return &searchResult, nil
		}

		// if could get object from object condition, condition of device and property will not empty
		// objIDs as filter condition of device and property
		deviceCond[common.BKObjIDField] = map[string]interface{}{common.BKDBIN: objIDs}
		propertyCond[common.BKObjIDField] = map[string]interface{}{common.BKDBIN: objIDs}
	}

	// if device has filter condition
	if 0 < len(deviceCond) {
		if 0 == len(objIDMapShowFields) {
			_, objIDMapShowFields, err = lgc.getObjIDsAndShowFields(pHeader, map[string]interface{}{})
			if nil != err {
				return nil, err
			}
			if 0 == len(objIDMapShowFields) {
				blog.Errorf("[NetProperty] search net object failed, could not get any net object, rid: %s", rid)
				return nil, defErr.Error(common.CCErrCollectNetPropertyGetFail)
			}
		}

		// get deviceID and value of fields to be shown by device condition
		deviceIDs, showFields.deviceIDMapDeviceShowFields, err = lgc.getDeviceIDsAndShowFields(
			pHeader, deviceCond, objIDMapShowFields)
		if nil != err {
			blog.Errorf("[NetProperty] search net property, get net device fail, error: %v, condition [%#v], rid: %s", err, deviceCond, rid)
			return nil, err
		}

		// if find any deviceIDs matched condition, will must not find any device propetry
		if 0 == len(deviceIDs) {
			return &searchResult, nil
		}

		// deviceIDs as filter conditoin of net property
		netPropertyCond[common.BKDeviceIDField] = map[string]interface{}{common.BKDBIN: deviceIDs}
	}

	netPropertyCond[common.BKOwnerIDField] = util.GetOwnerID(pHeader)
	searchResult.Count, err = lgc.db.Table(common.BKTableNameNetcollectProperty).Find(netPropertyCond).Count(lgc.ctx)
	if nil != err {
		blog.Errorf("[NetProperty] search net property fail, count net property by condition [%#v] error: %v, rid: %s", propertyCond, err, rid)
		return nil, err
	}
	if 0 == searchResult.Count {
		return &searchResult, nil
	}

	// field device_id and bk_property_id must be in params.Fields
	// to help add value of fields from other tables into search result
	if 0 != len(params.Fields) {
		params.Fields = append(params.Fields, []string{common.BKDeviceIDField, common.BKPropertyIDField}...)
	}

	if err = lgc.findProperty(params.Fields, netPropertyCond, &searchResult.Info, params.Page.Sort, params.Page.Start, params.Page.Limit); nil != err {
		blog.Errorf("[NetProperty] search net property fail, search net property by condition [%#v] error: %v, rid: %s", propertyCond, err, rid)
		return nil, defErr.Error(common.CCErrCollectNetPropertyGetFail)
	}

	// if net property are not empty, should add property and device shown info to the net property result
	deviceShowFieldLen := len(showFields.deviceIDMapDeviceShowFields)
	propertyShowFieldLen := len(showFields.propertyIDMapShowFields)

	// if object condition cond and device condition is empty, device shown fields will be empty
	// if property condition is empty, property shown fields will be empty
	if 0 == deviceShowFieldLen || 0 == propertyShowFieldLen {
		deviceIDs, propertyIDs = lgc.getDeviceIDsAndPropertyIDsFromNetPropertys(searchResult.Info)
	}

	if 0 == deviceShowFieldLen {
		showFields.deviceIDMapDeviceShowFields, err = lgc.getDeviceShowField(pHeader, deviceIDs)
		if nil != err {
			blog.Errorf("[NetProperty] search net property, get device show info fail, error: %v, rid: %s", err, rid)
			return nil, defErr.Error(common.CCErrCollectNetPropertyGetFail)
		}
	}
	if 0 == propertyShowFieldLen {
		showFields.propertyIDMapShowFields, err = lgc.getPropertyShowField(pHeader, propertyIDs)
		if nil != err {
			blog.Errorf("[NetProperty] search net property, get device show info fail, error: %v, rid: %s", err, rid)
			return nil, defErr.Error(common.CCErrCollectNetPropertyGetFail)
		}
	}

	// add value of fields from other tables into search result
	lgc.addShowFieldValueIntoNetProperty(searchResult.Info, showFields)

	return &searchResult, nil
}

func (lgc *Logics) DeleteProperty(pHeader http.Header, netPropertyID uint64) error {
	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pHeader))
	rid := util.GetHTTPCCRequestID(pHeader)

	netPropertyCond := map[string]interface{}{
		common.BKOwnerIDField:              util.GetOwnerID(pHeader),
		common.BKNetcollectPropertyIDField: netPropertyID}

	if err := lgc.db.Table(common.BKTableNameNetcollectProperty).Delete(lgc.ctx, netPropertyCond); nil != err {
		blog.Errorf("[NetProperty] delete net property with id [%d] failed, err: %v, params: %#v, rid: %s", netPropertyID, err, netPropertyCond, rid)
		return defErr.Error(common.CCErrCollectNetPropertyDeleteFail)
	}

	blog.V(5).Infof("[NetProperty] delete net property with id [%d] success", netPropertyID)

	return nil
}

func (lgc *Logics) addProperty(pHeader http.Header, netPropertyInfo meta.NetcollectProperty, ownerID string) (uint64, error) {
	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pHeader))
	rid := util.GetHTTPCCRequestID(pHeader)

	isExist, err := lgc.checkNetProperty(pHeader, &netPropertyInfo, ownerID)
	if nil != err {
		blog.Errorf("[NetProperty] add net property fail, %v, rid: %s", err, rid)
		return INVALIDID, err
	}
	if isExist { // exist the same [deviceID + propertyID], duplicate data
		blog.Errorf("[NetProperty] add net property fail, error: duplicate propertyID and deviceID, rid: %s", rid)
		return INVALIDID, defErr.Errorf(common.CCErrCommDuplicateItem, "property_id+device_id")
	}

	// add to the storage
	netPropertyID, err := lgc.addNewNetProperty(netPropertyInfo, ownerID)
	if nil != err {
		blog.Errorf("[NetProperty] add net property fail, error: %v, rid: %s", err, rid)
		return INVALIDID, defErr.Error(common.CCErrCollectNetPropertyCreateFail)
	}

	blog.V(5).Infof("[NetProperty] add net property, netPropertyInfo [%#+v], rid: %s", netPropertyInfo, rid)

	return netPropertyID, nil
}

func (lgc *Logics) updateProperty(
	pHeader http.Header, netPropertyInfo meta.NetcollectProperty, netPropertyID uint64, ownerID string) error {
	rid := util.GetHTTPCCRequestID(pHeader)

	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pHeader))

	isExist, err := lgc.checkNetProperty(pHeader, &netPropertyInfo, ownerID)
	if nil != err {
		blog.Errorf("[NetProperty] upate net property fail, %v, rid: %s", err, rid)
		return err
	}

	if isExist { // check if duplicate data
		propertyID, err := lgc.getNetPropertyID(netPropertyInfo.PropertyID, netPropertyInfo.DeviceID, ownerID)

		// error is not 'not found'
		if nil != err && lgc.db.IsNotFoundError(err) {
			blog.Errorf("[NetProperty] update net property fail, error: %v, deviceID [%d], propertyID [%s], rid: %s",
				err, netPropertyInfo.DeviceID, netPropertyInfo.PropertyID, rid)

			return defErr.Error(common.CCErrCollectNetPropertyUpdateFail)
		}

		// exist the same [deviceID + propertyID], duplicate data
		if nil == err && propertyID != netPropertyID {

			blog.Errorf("[NetProperty] update net property fail, duplicate deviceID [%d] and propertyID [%s], rid: %s",
				netPropertyInfo.DeviceID, netPropertyInfo.PropertyID, rid)

			return defErr.Errorf(common.CCErrCommDuplicateItem, "property_id+device_id")
		}
	}

	// update to the storage
	netPropertyInfo.OwnerID = ownerID
	if err := lgc.updateExistingPropertyByNetPropertyID(netPropertyInfo, netPropertyID); nil != err {
		blog.Errorf("[NetProperty] upadte net property fail, error: %v, rid: %s", err, rid)
		return defErr.Error(common.CCErrCollectNetPropertyUpdateFail)
	}

	blog.V(5).Infof("[NetProperty] update net property by net property id [%d], netPropertyInfo [%#+v], rid: %s",
		netPropertyID, netPropertyInfo, rid)

	return nil
}

func (lgc *Logics) addOrUpdateProperty(
	pHeader http.Header, netPropertyInfo meta.NetcollectProperty, ownerID string) (uint64, error) {
	rid := util.GetHTTPCCRequestID(pHeader)
	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pHeader))

	// check if data is valid and duplicate
	isExist, err := lgc.checkNetProperty(pHeader, &netPropertyInfo, ownerID)
	if nil != err {
		blog.Errorf("[NetProperty] batch add net property fail, error: %v, rid: %s", err, rid)
		return INVALIDID, err
	}
	if isExist { // update
		// get updated net property ID
		netPropertyID, err := lgc.getNetPropertyID(netPropertyInfo.PropertyID, netPropertyInfo.DeviceID, ownerID)
		if nil != err {
			blog.Errorf("[NetProperty] batch add net proeprty, failed to get id, error: %v, rid: %s", err, rid)
			return INVALIDID, defErr.Error(common.CCErrCollectNetPropertyUpdateFail)
		}

		netPropertyInfo.OwnerID = ownerID
		netPropertyInfo.NetcollectPropertyID = netPropertyID
		// update to the storage
		if err = lgc.updateNetPropertyByPropertyIDAndDeviceID(netPropertyInfo); nil != err {
			blog.Errorf("[NetProperty] batch add net proeprty, update net property failed, error: %v, rid: %s", err, rid)
			return INVALIDID, defErr.Error(common.CCErrCollectNetPropertyUpdateFail)
		}

		blog.V(5).Infof("[NetProperty] batch add net proeprty, update net property by [propertyID+deviceID], netPropertyInfo [%#+v], rid: %s",
			netPropertyInfo, rid)

		return netPropertyID, nil
	}

	// add to the storage
	netPropertyID, err := lgc.addNewNetProperty(netPropertyInfo, ownerID)
	if nil != err {
		blog.Errorf("[NetProperty] batch add net proeprty, add net collect property failed, error: %v, rid: %s", err, rid)
		return INVALIDID, defErr.Error(common.CCErrCollectNetPropertyCreateFail)
	}

	blog.V(5).Infof("[NetProperty] batch add net proeprty, add net property, netPropertyInfo [%#+v], rid: %s", netPropertyInfo, rid)

	return netPropertyID, nil
}

func (lgc *Logics) checkNetProperty(
	pHeader http.Header, netPropertyInfo *meta.NetcollectProperty, ownerID string) (isExist bool, err error) {
	rid := util.GetHTTPCCRequestID(pHeader)

	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pHeader))

	// check oid
	if "" == netPropertyInfo.OID {
		blog.Errorf("[NetProperty] check net collect property fail, oid is empty, rid: %s", rid)
		return false, defErr.Errorf(common.CCErrCommParamsNeedSet, common.BKOIDField)
	}

	// check period
	if "" != netPropertyInfo.Period && common.Infinite != netPropertyInfo.Period {
		netPropertyInfo.Period, err = util.FormatPeriod(netPropertyInfo.Period)
		if nil != err {
			blog.Errorf("[NetProperty] check net collect property, format period [%s] fail, error: %v, rid: %s", netPropertyInfo.Period, err, rid)
			return false, defErr.Error(common.CCErrCollectPeriodFormatFail)
		}
	}

	// check action
	if "" != netPropertyInfo.Action && !lgc.isValidAction(netPropertyInfo.Action) {
		blog.Errorf("[NetProperty] check net collect property, check action fail, action [%s] must be 'get' or 'walk', rid: %s", rid)
		return false, defErr.Errorf(common.CCErrCommParamsInvalid, common.BKActionField)
	}

	// check device existence
	// if device exist, propertyInfo will get device_id of device
	netPropertyInfo.DeviceID, netPropertyInfo.ObjectID, err = lgc.checkNetDeviceExist(
		pHeader, netPropertyInfo.DeviceID, netPropertyInfo.DeviceName)
	if nil != err {
		blog.Errorf("[NetProperty] check net collect property, check device fail, error: %v, rid: %s", err, rid)
		return false, err
	}

	// check if bk_property_id is valid and from object of net device
	// if bk_property_id is valid, propertyInfo will get bk_property_id of property
	netPropertyInfo.PropertyID, err = lgc.checkNetObjectProperty(
		pHeader, netPropertyInfo.ObjectID, netPropertyInfo.PropertyID, netPropertyInfo.PropertyName)
	if nil != err {
		blog.Errorf("[NetProperty] check net collect property, check property fail, error: %v, rid: %s", err, rid)
		return false, err
	}

	// check if data duplication
	isExist, err = lgc.checkNetPropertyExist(netPropertyInfo.DeviceID, netPropertyInfo.PropertyID, ownerID)
	if nil != err {
		blog.Errorf("[NetProperty] check net collect property, check data duplication fail, error: %v, rid: %s", err, rid)
		return false, defErr.Error(common.CCErrCollectNetPropertyCreateFail)
	}

	return isExist, nil
}

func (lgc *Logics) addNewNetProperty(netPropertyInfo meta.NetcollectProperty, ownerID string) (netPropertyID uint64, err error) {
	now := util.GetCurrentTimePtr()
	netPropertyInfo.CreateTime = now
	netPropertyInfo.LastTime = now
	netPropertyInfo.OwnerID = ownerID

	// set default value
	if "" == netPropertyInfo.Action {
		netPropertyInfo.Action = common.SNMPActionGet
	}
	if "" == netPropertyInfo.Period {
		netPropertyInfo.Period = common.Infinite
	}

	netPropertyInfo.NetcollectPropertyID, err = lgc.db.NextSequence(lgc.ctx, common.BKTableNameNetcollectProperty)
	if nil != err {
		return INVALIDID, fmt.Errorf("failed to get id, %v", err)
	}

	if err = lgc.db.Table(common.BKTableNameNetcollectProperty).Insert(lgc.ctx, netPropertyInfo); nil != err {
		return INVALIDID, err
	}

	return netPropertyInfo.NetcollectPropertyID, nil
}

func (lgc *Logics) updateNetPropertyByPropertyIDAndDeviceID(netPropertyInfo meta.NetcollectProperty) error {
	rid := util.ExtractRequestIDFromContext(lgc.ctx)
	queryParams := map[string]interface{}{
		common.BKDeviceIDField:   netPropertyInfo.DeviceID,
		common.BKPropertyIDField: netPropertyInfo.PropertyID,
		common.BKOwnerIDField:    netPropertyInfo.OwnerID,
	}

	netPropertyInfo.LastTime = util.GetCurrentTimePtr()

	if err := lgc.db.Table(common.BKTableNameNetcollectProperty).Update(lgc.ctx, queryParams, netPropertyInfo); nil != err {
		blog.Errorf("[NetProperty] update net property fail, error: %v, params: [%#+v], netPropertyInfo: [%#+v], rid: %s",
			err, queryParams, netPropertyInfo, rid)
		return err
	}

	return nil
}

func (lgc *Logics) updateExistingPropertyByNetPropertyID(netPropertyInfo meta.NetcollectProperty, netPropertyID uint64) error {
	rid := util.ExtractRequestIDFromContext(lgc.ctx)
	queryParams := map[string]interface{}{
		common.BKNetcollectPropertyIDField: netPropertyID,
		common.BKOwnerIDField:              netPropertyInfo.OwnerID,
	}

	netPropertyInfo.LastTime = util.GetCurrentTimePtr()
	netPropertyInfo.NetcollectPropertyID = netPropertyID

	if err := lgc.db.Table(common.BKTableNameNetcollectProperty).Update(lgc.ctx, queryParams, netPropertyInfo); nil != err {
		blog.Errorf("[NetProperty] update net property fail, error: %v, params: [%#+v], netPropertyInfo: [%#+v], rid: %s",
			err, queryParams, netPropertyInfo, rid)
		return err
	}

	return nil
}

// check if there is the same propertyInfo
func (lgc *Logics) checkNetPropertyExist(deviceID uint64, propertyID, ownerID string) (bool, error) {
	rid := util.ExtractRequestIDFromContext(lgc.ctx)
	queryParams := map[string]interface{}{
		common.BKDeviceIDField:   deviceID,
		common.BKPropertyIDField: propertyID,
		common.BKOwnerIDField:    ownerID,
	}
	rowCount, err := lgc.db.Table(common.BKTableNameNetcollectProperty).Find(queryParams).Count(lgc.ctx)
	if nil != err {
		blog.Errorf("[NetProperty] check if net deviceID and propertyID exist, query device fail, error information is %v, params:%v, rid: %s",
			err, queryParams, rid)
		return false, err
	}

	if 0 != rowCount {
		blog.V(5).Infof(
			"[NetProperty] check if net deviceID and propertyID exist, device_id[%s] and bk_property_id[%s] device is exist, rid: %s",
			deviceID, propertyID, rid)
		return true, nil
	}

	return false, nil
}

func (lgc *Logics) isValidAction(action string) bool {
	return common.SNMPActionGet == action || common.SNMPActionGetNext == action
}

func (lgc *Logics) findProperty(fields []string, condition, result interface{}, sort string, skip, limit int) error {
	rid := util.ExtractRequestIDFromContext(lgc.ctx)
	if err := lgc.db.Table(common.BKTableNameNetcollectProperty).Find(condition).Fields(fields...).Sort(sort).Start(uint64(skip)).Limit(uint64(limit)).All(lgc.ctx, result); err != nil {
		blog.Errorf("[NetProperty] failed to query the inst, error info %s, rid: %s", err.Error(), rid)
		return err
	}

	return nil
}

func (lgc *Logics) classifyNetPropertyCondition(
	conditionList []meta.ConditionItem) (map[string]interface{}, map[string]interface{}, map[string]interface{}, map[string]interface{}) {

	deviceCond := map[string]interface{}{}
	objectCond := map[string]interface{}{}
	propertyCond := map[string]interface{}{}
	netPropertyCond := map[string]interface{}{}

	for _, cond := range conditionList {
		if cond.Operator == common.BKDBEQ {
			switch cond.Field {
			case meta.AttributeFieldUnit, common.BKPropertyNameField, common.BKPropertyIDField:
				propertyCond[cond.Field] = cond.Value
			case common.BKObjIDField, common.BKObjNameField:
				objectCond[cond.Field] = cond.Value
			case common.BKDeviceIDField, common.BKDeviceNameField, common.BKDeviceModelField:
				deviceCond[cond.Field] = cond.Value
			default:
				netPropertyCond[cond.Field] = cond.Value
			}
		} else {
			switch cond.Field {
			case meta.AttributeFieldUnit, common.BKPropertyNameField, common.BKPropertyIDField:
				propertyCond[cond.Field] = map[string]interface{}{cond.Operator: cond.Value}
			case common.BKObjIDField, common.BKObjNameField:
				objectCond[cond.Field] = map[string]interface{}{cond.Operator: cond.Value}
			case common.BKDeviceIDField, common.BKDeviceNameField, common.BKDeviceModelField:
				deviceCond[cond.Field] = map[string]interface{}{cond.Operator: cond.Value}
			default:
				netPropertyCond[cond.Field] = map[string]interface{}{cond.Operator: cond.Value}
			}
		}
	}

	return deviceCond, objectCond, propertyCond, netPropertyCond
}

type netPropertyShowFields struct {
	deviceIDMapDeviceShowFields map[uint64]deviceShowField   // id map value group of device fields
	propertyIDMapShowFields     map[string]propertyShowField // propertyID:objID map value group of property fields
}

type objShowField struct {
	objName string
}

type deviceShowField struct {
	deviceName  string
	deviceModel string
	objID       string
	objName     string
}

type propertyShowField struct {
	unit         string
	propertyName string
}

// get obj ID list and get field to show by map (bk_obj_id --> bk_obj_name)
func (lgc *Logics) getObjIDsAndShowFields(pHeader http.Header, objectCond map[string]interface{}) ([]string, map[string]objShowField, error) {
	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pHeader))
	rid := util.GetHTTPCCRequestID(pHeader)
	objectCond[common.BKClassificationIDField] = common.BKNetwork

	objResult, err := lgc.CoreAPI.CoreService().Model().ReadModel(context.Background(), pHeader, &meta.QueryCondition{Condition: objectCond})
	if nil != err {
		blog.Errorf("[NetProperty] get net device object fail, error: %v, condition [%#v], rid: %s", err, objectCond, rid)
		return nil, nil, defErr.Error(common.CCErrObjectSelectInstFailed)
	}
	if !objResult.Result {
		blog.Errorf("[NetProperty] get net device object fail, errors: %s, condition [%#v], rid: %s", objResult.ErrMsg, objectCond, rid)
		return nil, nil, defErr.New(objResult.Code, objResult.ErrMsg)
	}

	if 0 == len(objResult.Data.Info) {
		return nil, nil, nil
	}

	objIDs := make([]string, 0)
	objIDMapobjName := map[string]objShowField{}
	for _, obj := range objResult.Data.Info {
		objIDs = append(objIDs, obj.Spec.ObjectID)
		objIDMapobjName[obj.Spec.ObjectID] = objShowField{obj.Spec.ObjectName}
	}

	return objIDs, objIDMapobjName, nil
}

// get device ID list and get field to show by map (device_id --> bk_device_name, ...)
// add obj show field into device show fields
func (lgc *Logics) getDeviceIDsAndShowFields(
	pHeader http.Header, deviceCond map[string]interface{}, objIDMapShowFields map[string]objShowField) ([]uint64, map[uint64]deviceShowField, error) {

	rid := util.ExtractRequestIDFromContext(lgc.ctx)
	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pHeader))

	deviceCond[common.BKOwnerIDField] = util.GetOwnerID(pHeader)
	deviceField := []string{common.BKDeviceIDField, common.BKDeviceNameField, common.BKDeviceModelField, common.BKObjIDField}
	deviceResult := make([]meta.NetcollectDevice, 0)

	if err := lgc.findDevice(deviceField, deviceCond, &deviceResult, "", 0, 0); nil != err {
		blog.Errorf("[NetProperty] search net device fail by condition [%#v], error: %v, rid: %s", deviceCond, err, rid)
		if !lgc.db.IsNotFoundError(err) {
			return nil, nil, nil
		}
		return nil, nil, defErr.Error(common.CCErrCollectNetDeviceGetFail)
	}

	deviceIDs, deviceIDMapDeviceShowFields := lgc.assembleDeviceShowFieldValue(deviceResult, objIDMapShowFields)

	if 0 == len(deviceIDs) {
		return nil, nil, nil
	}

	return deviceIDs, deviceIDMapDeviceShowFields, nil
}

// get device IDs from device list
// assemble value of device list: [deviceID] map [deviceName, deviceModel, objID, objName]
// objName is taken from objIDMapShowFields
func (lgc *Logics) assembleDeviceShowFieldValue(deviceData []meta.NetcollectDevice, objIDMapShowFields map[string]objShowField) (
	deviceIDs []uint64, deviceIDMapDeviceShowFields map[uint64]deviceShowField) {

	if nil == deviceData || 0 == len(deviceData) {
		return deviceIDs, deviceIDMapDeviceShowFields
	}

	deviceIDMapDeviceShowFields = map[uint64]deviceShowField{}

	for _, device := range deviceData {
		// get device IDs from device list
		deviceIDs = append(deviceIDs, device.DeviceID)
		// assemble value of device list: [deviceID] map [deviceName, deviceModel, objID, objName]
		deviceIDMapDeviceShowFields[device.DeviceID] = deviceShowField{
			deviceName:  device.DeviceName,
			deviceModel: device.DeviceModel,
			objID:       device.ObjectID,
			objName:     objIDMapShowFields[device.ObjectID].objName,
		}
	}

	return deviceIDs, deviceIDMapDeviceShowFields
}

// get objectID, property ID list and get field to show by map (bk_property_id --> bk_property_name, ...)
func (lgc *Logics) getPropertyIDsAndShowFields(pHeader http.Header, propertyCond map[string]interface{}) ([]string, []string, map[string]propertyShowField, error) {
	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pHeader))
	rid := util.GetHTTPCCRequestID(pHeader)

	attrResult, err := lgc.CoreAPI.CoreService().Model().ReadModelAttrByCondition(context.Background(), pHeader, &meta.QueryCondition{Condition: propertyCond})
	if nil != err {
		blog.Errorf("[NetProperty] get property fail, error: %v, condition [%#v], rid: %s", err, propertyCond, rid)
		return nil, nil, nil, defErr.Error(common.CCErrTopoObjectAttributeSelectFailed)
	}
	if !attrResult.Result {
		blog.Errorf("[NetProperty] get property fail, error: %s, rid: %s", attrResult.ErrMsg, rid)
		return nil, nil, nil, defErr.New(attrResult.Code, attrResult.ErrMsg)
	}

	objIDs, propertyIDs, propertyIDMapPropertyShowFields := lgc.assembleAttrShowFieldValue(attrResult.Data.Info)

	if 0 == len(objIDs) || 0 == len(propertyIDs) || 0 == len(propertyIDMapPropertyShowFields) {
		blog.Errorf("[NetProperty] get property fail, property is not exist, condition [%#v], rid: %s", propertyCond, rid)
		return nil, nil, nil, nil
	}

	return objIDs, propertyIDs, propertyIDMapPropertyShowFields, nil
}

// get obj IDs and property IDs , assemble value of attribute list:[propertyID : objID] map [property show fields]
func (lgc *Logics) assembleAttrShowFieldValue(attrData []meta.Attribute) (
	objIDs []string, propertyIDs []string, propertyIDMapPropertyShowFields map[string]propertyShowField) {

	if nil == attrData || 0 == len(attrData) {
		return []string{}, []string{}, map[string]propertyShowField{}
	}

	// get obj IDs and property IDs from attribute list
	propertyIDs, objIDs = []string{}, []string{}
	// assemble value of attribute list: [propertyID : objID] map [property unit, property name]
	propertyIDMapPropertyShowFields = map[string]propertyShowField{}

	for _, property := range attrData {
		propertyIDs = append(propertyIDs, property.PropertyID)
		objIDs = append(objIDs, property.ObjectID)

		propertyIDMapPropertyShowFields[propertyMapKey(property.PropertyID, property.ObjectID)] = propertyShowField{
			unit:         property.Unit,
			propertyName: property.PropertyName,
		}
	}

	return objIDs, propertyIDs, propertyIDMapPropertyShowFields
}

// add group value of device and property to net property
func (lgc *Logics) addShowFieldValueIntoNetProperty(
	netPropertys []meta.NetcollectProperty, netPropShowFields netPropertyShowFields) {

	for index := range netPropertys {
		netPropertys := &netPropertys[index]

		deviceValue := netPropShowFields.deviceIDMapDeviceShowFields[netPropertys.DeviceID]

		// add group value of device
		netPropertys.DeviceModel = deviceValue.deviceModel
		netPropertys.DeviceName = deviceValue.deviceName
		netPropertys.ObjectID = deviceValue.objID
		netPropertys.ObjectName = deviceValue.objName

		propertyID := netPropertys.PropertyID
		propertyValue := netPropShowFields.propertyIDMapShowFields[propertyMapKey(propertyID, deviceValue.objID)]

		// add group value of property
		netPropertys.Unit = propertyValue.unit
		netPropertys.PropertyName = propertyValue.propertyName
	}
}

func (lgc *Logics) getDeviceIDsAndPropertyIDsFromNetPropertys(
	netProperty []meta.NetcollectProperty) (deviceIDs []uint64, propertyIDs []string) {

	for index := range netProperty {
		deviceIDs = append(deviceIDs, netProperty[index].DeviceID)
		propertyIDs = append(propertyIDs, netProperty[index].PropertyID)
	}

	return deviceIDs, propertyIDs
}

// get device shown info by deviceIDs
func (lgc *Logics) getDeviceShowField(pHeader http.Header, deviceIDs []uint64) (map[uint64]deviceShowField, error) {
	_, objIDMapShowFields, err := lgc.getObjIDsAndShowFields(pHeader, map[string]interface{}{})
	if nil != err {
		return nil, err
	}
	if 0 == len(objIDMapShowFields) {
		return nil, fmt.Errorf("search net object failed, could not get any net object")
	}

	deviceCond := map[string]interface{}{
		common.BKDeviceIDField: map[string]interface{}{common.BKDBIN: deviceIDs},
		common.BKOwnerIDField:  util.GetOwnerID(pHeader),
	}
	_, deviceIDMapDeviceShowFields, err := lgc.getDeviceIDsAndShowFields(pHeader, deviceCond, objIDMapShowFields)
	if nil != err {
		return nil, err
	}

	if 0 == len(deviceIDMapDeviceShowFields) {
		return nil, fmt.Errorf("search net device failed, could not get any net device by condition [%#+v]", deviceCond)
	}

	return deviceIDMapDeviceShowFields, nil
}

// get property shown info by propertyIDs
func (lgc *Logics) getPropertyShowField(pHeader http.Header, propertyIDs []string) (map[string]propertyShowField, error) {
	propertyCond := map[string]interface{}{
		common.BKPropertyIDField: map[string]interface{}{common.BKDBIN: propertyIDs},
	}

	_, _, propertyIDMapPropertyShowFields, err := lgc.getPropertyIDsAndShowFields(pHeader, propertyCond)
	if nil != err {
		return nil, err
	}

	if 0 == len(propertyIDMapPropertyShowFields) {
		return nil, fmt.Errorf("search property failed, could not get any property by condition [%#+v]", propertyCond)
	}

	return propertyIDMapPropertyShowFields, nil
}

func propertyMapKey(propertyID, objID string) string {
	return fmt.Sprintf("%s:%s", propertyID, objID)
}
