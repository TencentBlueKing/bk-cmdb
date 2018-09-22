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
	"regexp"
	"strconv"
	"time"

	restful "github.com/emicklei/go-restful"
	mgo "gopkg.in/mgo.v2"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	mapStr "configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (lgc *Logics) AddProperty(
	pheader http.Header, propertyInfoList []meta.NetcollectProperty) ([]meta.AddNetPropertyResult, bool) {
	ownerID := util.GetOwnerID(pheader)

	resultList := make([]meta.AddNetPropertyResult, 0)
	hasError := false

	for _, propertyInfo := range propertyInfoList {
		errMsg := ""
		result := true

		propertyID, err := lgc.addProperty(propertyInfo, pheader, ownerID)
		if nil != err {
			errMsg = err.Error()
			result = false
			hasError = true
		}

		resultList = append(resultList, meta.AddNetPropertyResult{result, errMsg, propertyID})
	}

	return resultList, hasError
}

func (lgc *Logics) SearchProperty(pheader http.Header, params *meta.NetCollSearchParams) (meta.SearchNetProperty, error) {
	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	ownerID := util.GetOwnerID(pheader)

	deviceCond, objectCond, propertyCond, netPropertyCond := lgc.classifyNetPropertyCondition(params.Condition)

	searchResult := meta.SearchNetProperty{0, []mapStr.MapStr{}}

	var (
		err                error
		objIDs             []string
		deviceIDs          []string
		propertyIDs        []string
		showFeilds         netPropertyShowFeilds
		objIDMapShowFeilds map[string]objShowFeild
	)
	// 如果有对 obj 有筛选条件
	if 0 < len(objectCond) {
		objIDs, objIDMapShowFeilds, err = lgc.getObjIDsAndShowFeilds(objectCond, pheader)
		if nil != err {
			blog.Errorf("check net device object, get net device object fail, error: %v, condition [%#v]", err, objectCond)
			return meta.SearchNetProperty{}, err
		}

		if 0 == len(objIDs) {
			return searchResult, nil
		}
		deviceCond[common.BKObjIDField] = map[string]interface{}{common.BKDBIN: objIDs}
		propertyCond[common.BKObjIDField] = map[string]interface{}{common.BKDBIN: objIDs}
	}

	// 如果有对 device 有筛选条件
	if 0 < len(deviceCond) || 0 < len(objIDs) {
		deviceCond[common.BKOwnerIDField] = ownerID

		deviceIDs, showFeilds.deviceIDMapDeviceShowFeilds, err = lgc.getDeviceIDsAndShowFeilds(deviceCond, objIDMapShowFeilds, pheader)
		objIDMapShowFeilds = nil
		if nil != err {
			blog.Errorf("check net device object, get net device object fail, error: %v, condition [%#v]", err, objectCond)
			return meta.SearchNetProperty{}, err
		}

		if 0 == len(deviceIDs) {
			return searchResult, nil
		}
		netPropertyCond[common.BKDeviceIDField] = map[string]interface{}{common.BKDBIN: deviceIDs}
	}

	// 如果有对 property 有筛选条件
	if 0 < len(propertyCond) || 0 < len(objIDs) {
		propertyIDs, showFeilds.propertyIDMapShowFeilds, err = lgc.getPropertyIDsAndShowFeilds(propertyCond, pheader)
		if nil != err {
			blog.Errorf("check net device object, get net device object fail, error: %v, condition [%#v]", err, objectCond)
			return meta.SearchNetProperty{}, err
		}

		if 0 == len(propertyIDs) {
			return searchResult, nil
		}
		netPropertyCond[common.BKPropertyIDField] = map[string]interface{}{common.BKDBIN: propertyIDs}
	}

	searchResult.Count, err = lgc.Instance.GetCntByCondition(common.BKTableNameNetcollectDevice, deviceCond)
	if nil != err {
		blog.Errorf("search net device fail, count net device by condition [%#v] error: %v", deviceCond, err)
		return meta.SearchNetProperty{}, nil
	}
	if 0 == searchResult.Count {
		return searchResult, nil
	}

	if err = lgc.findProperty(params.Fields, deviceCond, &searchResult.Info, params.Page.Sort, params.Page.Start, params.Page.Limit); nil != err {
		blog.Errorf("search net device fail, search net device by condition [%#v] error: %v", deviceCond, err)
		return meta.SearchNetProperty{}, defErr.Errorf(common.CCErrCollectNetDeviceGetFail)
	}

	lgc.addShowFieldValueIntoNetProperty(&searchResult, showFeilds)
	return searchResult, nil
}

func (lgc *Logics) DeleteProperty(req *restful.Request, resp *restful.Response) {

}

func (lgc *Logics) addProperty(propertyInfo meta.NetcollectProperty, pheader http.Header, ownerID string) (int64, error) {
	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	if "" == propertyInfo.OID { // check oid
		blog.Errorf("add net collect property fail, oid is empty")
		return -1, defErr.Errorf(common.CCErrCommParamsLostField, common.BKOIDField)
	}

	// check period
	var err error
	if "" != propertyInfo.Period && common.Infinite != propertyInfo.Period {
		propertyInfo.Period, err = lgc.formatPeriod(propertyInfo.Period, defErr)
		if nil != err {
			return -1, err
		}
	}

	// check action
	if "" != propertyInfo.Action && !lgc.isValidAction(propertyInfo.Action) {
		blog.Errorf("add net collect property fail, action [%s] must be 'get' or 'walk' ")
		return -1, defErr.Errorf(common.CCErrCommParamsInvalid, common.BKActionField)
	}

	// check device
	if err = lgc.checkIfNetDeviceExist(&propertyInfo, pheader); nil != err {
		blog.Errorf("add net collect property fail, error: %v", err)
		return -1, err
	}

	// check property
	if err = lgc.checkIfNetProperty(&propertyInfo, pheader); nil != err {
		blog.Errorf("add net collect property fail, error: %v", err)
		return -1, err
	}

	// check if data duplication
	isExist, err := lgc.checkNetPropertyExist(propertyInfo.DeviceID, propertyInfo.PropertyID, ownerID)
	if nil != err {
		blog.Errorf("add net collect property fail, error: %v", err)
		return -1, defErr.Errorf(common.CCErrCollectNetPropertyCreateFail)
	}
	if isExist {
		blog.Errorf("add net collect property fail, error: duplicate [deviceID propertyID]")
		return -1, defErr.Errorf(common.CCErrCommDuplicateItem)
	}

	now := time.Now()
	propertyInfo.CreateTime = &now
	now = time.Now()
	propertyInfo.LastTime = &now
	propertyInfo.OwnerID = ownerID
	// set default value
	if "" == propertyInfo.Action {
		propertyInfo.Action = common.ActionGet
	}
	if "" == propertyInfo.Period {
		propertyInfo.Period = common.Infinite
	}

	propertyInfo.NetcollectPropertyID, err = lgc.Instance.GetIncID(common.BKTableNameNetcollectProperty)
	if nil != err {
		blog.Errorf("add net collect property, failed to get id, error: %v", err)
		return -1, defErr.Errorf(common.CCErrCollectNetDeviceCreateFail)
	}

	if _, err = lgc.Instance.Insert(common.BKTableNameNetcollectProperty, propertyInfo); nil != err {
		blog.Errorf("failed to insert net collect property, error: %v", err)
		return -1, defErr.Errorf(common.CCErrCollectNetDeviceCreateFail)
	}

	return propertyInfo.NetcollectPropertyID, nil
}

// check if bk_property_id is valid and from object of net device
// if bk_property_id is valid, propertyInfo will get bk_property_id of property
func (lgc *Logics) checkIfNetProperty(propertyInfo *meta.NetcollectProperty, pheader http.Header) error {
	var err error
	propertyInfo.PropertyID, err = lgc.checkNetObjectProperty(propertyInfo.ObjectID, propertyInfo.PropertyID, propertyInfo.PropertyName, pheader)
	return err
}

// check if device exist or not
// if device exist, propertyInfo will get bk_device_id of device
func (lgc *Logics) checkIfNetDeviceExist(propertyInfo *meta.NetcollectProperty, pheader http.Header) error {
	var err error
	propertyInfo.DeviceID, propertyInfo.ObjectID, err = lgc.checkNetDeviceExist(propertyInfo.DeviceID, propertyInfo.DeviceName, pheader)
	return err
}

// check if there is the same propertyInfo
func (lgc *Logics) checkNetPropertyExist(deviceID int64, propertyID, ownerID string) (bool, error) {
	queryParams := common.KvMap{
		common.BKDeviceIDField: deviceID, common.BKPropertyIDField: propertyID, common.BKOwnerIDField: ownerID}

	rowCount, err := lgc.Instance.GetCntByCondition(common.BKTableNameNetcollectProperty, queryParams)
	if nil != err {
		blog.Errorf("check if net deviceID and propertyID exist, query device fail, error information is %v, params:%v",
			err, queryParams)
		return false, err
	}

	if 0 != rowCount {
		blog.V(4).Infof(
			"check if net deviceID and propertyID exist, bk_device_id is [%s] bk_property_id [%s] device is exist",
			deviceID, propertyID)
		return true, nil
	}

	return false, nil
}

const periodRegexp = "^\\d*[DHMS]$" // period regexp to check period

// 00002H --> 2H
// 0000D/0M ---> ∞
// empty string / ∞ ---> ∞
// regexp matched: positive integer (include positive integer begin with more the one '0') + [D/H/M/S]
// eg. 0H, 000H, 0002H, 32M，34S...
// examples of no matched:  1.4H, -2H, +2H ...
func (lgc *Logics) formatPeriod(period string, defErr errors.DefaultCCErrorIf) (string, error) {
	if common.Infinite == period || "" == period {
		return common.Infinite, nil
	}

	ok, _ := regexp.Match(periodRegexp, []byte(period))
	if !ok {
		return "", defErr.Errorf(common.CCErrCommParamsInvalid, common.BKPeriodField)
	}

	num, err := strconv.Atoi(period[:len(period)-1])
	if nil != err {
		return "", defErr.Error(common.CCErrCollectPeridFormatFail)
	}
	if 0 == num {
		return common.Infinite, nil
	}

	return strconv.Itoa(num) + period[len(period)-1:], nil
}

func (lgc *Logics) isValidAction(action string) bool {
	return common.ActionGet == action || common.ActionWalk == action
}

func (lgc *Logics) findProperty(fields []string, condition, result interface{}, sort string, skip, limit int) error {
	if err := lgc.Instance.GetMutilByCondition(common.BKTableNameNetcollectProperty, fields, condition, result, sort, skip, limit); err != nil {
		blog.Errorf("failed to query the inst, error info %s", err.Error())
		return err
	}

	return nil
}

func (lgc *Logics) classifyNetPropertyCondition(conditionList []meta.ConditionItem) (map[string]interface{}, map[string]interface{}, map[string]interface{}, map[string]interface{}) {
	var deviceCond, objectCond, propertyCond, netPropertyCond map[string]interface{}

	for _, cond := range conditionList {
		if cond.Operator == common.BKDBEQ {
			switch cond.Field {
			case meta.AttributeFieldUnit:
				fallthrough
			case common.BKPropertyNameField:
				fallthrough
			case common.BKPropertyIDField:
				propertyCond[cond.Field] = cond.Value
			case common.BKObjIDField:
				fallthrough
			case common.BKObjNameField:
				objectCond[cond.Field] = cond.Value
			case common.BKDeviceIDField:
				fallthrough
			case common.BKDeviceNameField:
				fallthrough
			case common.BKDeviceModelField:
				deviceCond[cond.Field] = cond.Value
			default:
				netPropertyCond[cond.Field] = cond.Value
			}

		} else {
			switch cond.Field {
			case meta.AttributeFieldUnit:
				fallthrough
			case common.BKPropertyNameField:
				fallthrough
			case common.BKPropertyIDField:
				propertyCond[cond.Field] = map[string]interface{}{cond.Operator: cond.Value}
			case common.BKObjIDField:
				fallthrough
			case common.BKObjNameField:
				objectCond[cond.Field] = map[string]interface{}{cond.Operator: cond.Value}
			case common.BKDeviceIDField:
				fallthrough
			case common.BKDeviceNameField:
				fallthrough
			case common.BKDeviceModelField:
				deviceCond[cond.Field] = map[string]interface{}{cond.Operator: cond.Value}
			default:
				netPropertyCond[cond.Field] = map[string]interface{}{cond.Operator: cond.Value}
			}
		}

	}

	return deviceCond, objectCond, propertyCond, netPropertyCond
}

// id map feilds
type netPropertyShowFeilds struct {
	deviceIDMapDeviceShowFeilds map[string]deviceShowFeild
	propertyIDMapShowFeilds     map[string]propertyShowFeild
}

type objShowFeild struct {
	objName string
}
type deviceShowFeild struct {
	deviceName  string
	deviceModel string
	objID       string
	objName     string
}

type propertyShowFeild struct {
	unit         string
	propertyName string
}

// get obj ID list and get feild to show by map (bk_obj_id --> bk_obj_name)
func (lgc *Logics) getObjIDsAndShowFeilds(objectCond map[string]interface{}, pheader http.Header) ([]string, map[string]objShowFeild, error) {
	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	objectCond[common.BKClassificationIDField] = common.BKNetwork

	objResult, err := lgc.CoreAPI.ObjectController().Meta().SelectObjects(context.Background(), pheader, objectCond)
	if nil != err {
		blog.Errorf("check net device object, get net device object fail, error: %v, condition [%#v]", err, objectCond)
		return nil, nil, defErr.Errorf(common.CCErrObjectSelectInstFailed)
	}

	if !objResult.Result {
		blog.Errorf("check net device object, errors: %s, condition [%#v]", objResult.ErrMsg, objectCond)
		return nil, nil, defErr.Errorf(objResult.Code)
	}

	if nil == objResult.Data || 0 == len(objResult.Data) {
		return nil, nil, nil
	}

	objIDs := []string{}
	objIDMapobjName := map[string]objShowFeild{}
	for _, obj := range objResult.Data {
		objIDs = append(objIDs, obj.ObjectID)
		objIDMapobjName[obj.ObjectID] = objShowFeild{obj.ObjectName}
	}
	objResult = nil

	return objIDs, objIDMapobjName, nil
}

// get device ID list and get feild to show by map (bk_device_id --> bk_device_name, ...)
func (lgc *Logics) getDeviceIDsAndShowFeilds(deviceCond map[string]interface{}, objIDMapShowFeilds map[string]objShowFeild, pheader http.Header) ([]string, map[string]deviceShowFeild, error) {
	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	deviceFeild := []string{common.BKDeviceIDField, common.BKDeviceNameField, common.BKDeviceModelField}
	deviceResult := []mapStr.MapStr{}

	if err := lgc.findDevice(deviceFeild, deviceCond, &deviceResult, "", 0, 0); nil != err {
		blog.Errorf("search net device fail, search net device by condition [%#v] error: %v", deviceCond, err)
		if mgo.ErrNotFound == err {
			return nil, nil, nil
		}
		return nil, nil, defErr.Errorf(common.CCErrCollectNetDeviceGetFail)
	}

	if 0 == len(deviceResult) {
		return nil, nil, nil
	}

	deviceIDs := []string{}
	deviceIDMapDeviceShowFeilds := map[string]deviceShowFeild{}
	for _, device := range deviceResult {
		deviceID := device[common.BKDeviceIDField].(string)
		deviceIDs = append(deviceIDs, deviceID)
		deviceIDMapDeviceShowFeilds[deviceID] = deviceShowFeild{
			device[common.BKDeviceNameField].(string),
			device[common.BKDeviceModelField].(string),
			device[common.BKObjIDField].(string),
			objIDMapShowFeilds[device[common.BKObjIDField].(string)].objName,
		}
	}
	deviceResult = nil

	return deviceIDs, deviceIDMapDeviceShowFeilds, nil
}

// get property ID list and get feild to show by map (bk_property_id --> bk_property_name, ...)
func (lgc *Logics) getPropertyIDsAndShowFeilds(propertyCond map[string]interface{}, pheader http.Header) ([]string, map[string]propertyShowFeild, error) {
	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	attrResult, err := lgc.CoreAPI.ObjectController().Meta().SelectObjectAttWithParams(context.Background(), pheader, propertyCond)
	if nil != err {
		blog.Errorf("get object attribute fail, error: %v, condition [%#v]", err, propertyCond)
		if mgo.ErrNotFound == err {
			return nil, nil, nil
		}
		return nil, nil, defErr.Errorf(common.CCErrTopoObjectAttributeSelectFailed)
	}
	if !attrResult.Result {
		blog.Errorf("check net device object property, errors: %s", attrResult.ErrMsg)
		return nil, nil, defErr.Errorf(attrResult.Code)
	}

	if nil == attrResult.Data || 0 == len(attrResult.Data) {
		blog.Errorf("check net device object property, property is not exist, condition [%#v]", propertyCond)
		return nil, nil, nil
	}

	propertyIDs := []string{}
	propertyIDMapDeviceShowFeilds := map[string]propertyShowFeild{}
	for _, property := range attrResult.Data {
		propertyIDs = append(propertyIDs, property.PropertyID)
		propertyIDMapDeviceShowFeilds[property.PropertyID] = propertyShowFeild{
			property.Unit,
			property.PropertyName,
		}
	}
	attrResult = nil

	return propertyIDs, propertyIDMapDeviceShowFeilds, nil
}

func (lgc *Logics) addShowFieldValueIntoNetProperty(
	searchNetProperty *meta.SearchNetProperty, netPropShowFeilds netPropertyShowFeilds) {

	for _, netProperty := range searchNetProperty.Info {
		deviceID := netProperty[common.BKDeviceIDField].(string)
		propertyID := netProperty[common.BKPropertyIDField].(string)

		deviceValue := netPropShowFeilds.deviceIDMapDeviceShowFeilds[deviceID]
		propertyValue := netPropShowFeilds.propertyIDMapShowFeilds[propertyID]

		netProperty[common.BKDeviceModelField] = deviceValue.deviceModel
		netProperty[common.BKDeviceNameField] = deviceValue.deviceName
		netProperty[common.BKObjIDField] = deviceValue.objID
		netProperty[common.BKObjNameField] = deviceValue.objName

		netProperty[meta.AttributeFieldUnit] = propertyValue.unit
		netProperty[common.BKPropertyNameField] = propertyValue.propertyName
	}
}
