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
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	mapStr "configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (lgc *Logics) AddDevices(pheader http.Header, deviceInfoList []meta.NetcollectDevice) ([]meta.AddDeviceResult, bool) {
	ownerID := util.GetOwnerID(pheader)

	resultList := make([]meta.AddDeviceResult, 0)
	hasError := false

	for _, deviceInfo := range deviceInfoList {
		errMsg := ""
		result := true

		deviceID, err := lgc.addDevice(deviceInfo, pheader, ownerID)
		if nil != err {
			errMsg = err.Error()
			result = false
			hasError = true
		}

		resultList = append(resultList, meta.AddDeviceResult{result, errMsg, deviceID})
	}

	return resultList, hasError
}

func (lgc *Logics) SearchDevice(pheader http.Header, params *meta.NetCollSearchParams) (meta.SearchNetDevice, error) {
	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	deviceCond := map[string]interface{}{}
	deviceCond[common.BKOwnerIDField] = util.GetOwnerID(pheader)

	objCond := map[string]interface{}{}

	// get condition, condtion of objs and condtion of device
	for _, cond := range params.Condition {
		switch cond.Operator {
		case common.BKDBEQ:
			if common.BKObjNameField == cond.Field {
				objCond[cond.Field] = cond.Value
			} else {
				deviceCond[cond.Field] = cond.Value
			}
		default:
			if common.BKObjNameField == cond.Field {
				objCond[cond.Field] = map[string]interface{}{
					cond.Operator: cond.Value,
				}
			} else {
				deviceCond[cond.Field] = map[string]interface{}{
					cond.Operator: cond.Value,
				}
			}
		}
	}

	// if condition only has bk_obj_name but not bk_obj_id
	// get net device bk_obj_id from bk_obj_name
	if _, ok := deviceCond[common.BKObjIDField]; !ok && 0 < len(objCond) {
		objIDs, err := lgc.getNetDeviceObjIDsByCond(objCond, pheader)
		if nil != err {
			blog.Errorf("search net device fail, search net device obj id by condition [%#v] error: %v", objCond, err)
			return meta.SearchNetDevice{}, defErr.Errorf(common.CCErrCollectNetDeviceGetFail)
		}
		deviceCond[common.BKObjIDField] = map[string]interface{}{
			common.BKDBIN: objIDs,
		}
	}

	searchResult := meta.SearchNetDevice{}
	var err error

	searchResult.Count, err = lgc.Instance.GetCntByCondition(common.BKTableNameNetcollectDevice, deviceCond)
	if nil != err {
		blog.Errorf("search net device fail, count net device by condition [%#v] error: %v", deviceCond, err)
		return meta.SearchNetDevice{}, nil
	}
	if 0 == searchResult.Count {
		searchResult.Info = []mapStr.MapStr{}
		return searchResult, nil
	}

	err = lgc.findDevice(params.Fields, deviceCond, &searchResult.Info, params.Page.Sort, params.Page.Start, params.Page.Limit)
	if nil != err {
		blog.Errorf("search net device fail, search net device by condition [%#v] error: %v", deviceCond, err)
		return meta.SearchNetDevice{}, defErr.Errorf(common.CCErrCollectNetDeviceGetFail)
	}

	return searchResult, nil
}

// add a device
func (lgc *Logics) addDevice(deviceInfo meta.NetcollectDevice, pheader http.Header, ownerID string) (int64, error) {
	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	if "" == deviceInfo.DeviceModel {
		blog.Errorf("add net device fail, device_model is empty")
		return -1, defErr.Errorf(common.CCErrCommParamsLostField, common.BKDeviceModelField)
	}

	if "" == deviceInfo.BkVendor {
		blog.Errorf("add net device fail, bk_vendor is empty")
		return -1, defErr.Errorf(common.CCErrCommParamsLostField, common.BKVendorField)
	}

	if "" == deviceInfo.DeviceName {
		blog.Errorf("add net device fail, device_name is empty")
		return -1, defErr.Errorf(common.CCErrCommParamsLostField, common.BKDeviceModelField)
	}

	// check if device_name exist
	isExist, err := lgc.checkIfNetDeviceNameExist(deviceInfo.DeviceName, ownerID)
	if nil != err {
		blog.Errorf("add net device fail, error: %v", err)
		return -1, defErr.Errorf(common.CCErrCollectNetDeviceCreateFail)
	}
	if isExist {
		blog.Errorf("add net device fail, error: duplicate device_name")
		return -1, defErr.Errorf(common.CCErrCommDuplicateItem)
	}

	// check if bk_object_id and bk_object_name are net device object
	isNetDevice, err := lgc.checkIfNetDeviceObject(&deviceInfo, pheader)
	if nil != err {
		blog.Errorf("add net device fail, error: %v", err)
		return -1, err
	}
	if !isNetDevice {
		blog.Errorf("add net device fail, object name [%s] and object ID [%s] is not netcollect device",
			deviceInfo.ObjectName, deviceInfo.ObjectID)
		return -1, defErr.Errorf(common.CCErrCollectObjIDNotNetDevice)
	}

	// add to the storage
	now := time.Now()
	deviceInfo.CreateTime = &now
	now = time.Now()
	deviceInfo.LastTime = &now
	deviceInfo.OwnerID = ownerID

	deviceInfo.DeviceID, err = lgc.Instance.GetIncID(common.BKTableNameNetcollectDevice)
	if nil != err {
		blog.Errorf("add net device, failed to get id, error: %v", err)
		return -1, defErr.Errorf(common.CCErrCollectNetDeviceCreateFail)
	}

	_, err = lgc.Instance.Insert(common.BKTableNameNetcollectDevice, deviceInfo)
	if nil != err {
		blog.Error("failed to insert net device, error: %v", err)
		return -1, defErr.Errorf(common.CCErrCollectNetDeviceCreateFail)
	}

	return deviceInfo.DeviceID, nil
}

func (lgc *Logics) findDevice(fields []string, condition, result interface{}, sort string, skip, limit int) error {
	if err := lgc.Instance.GetMutilByCondition(common.BKTableNameNetcollectDevice, fields, condition, result, sort, skip, limit); err != nil {
		blog.Errorf("failed to query the inst, error info %s", err.Error())
		return err
	}

	return nil
}

// check the deviceInfo if is a net object
// by checking if bk_obj_id and bk_obj_name function parameter are valid net device object or not
func (lgc *Logics) checkIfNetDeviceObject(deviceInfo *meta.NetcollectDevice, pheader http.Header) (bool, error) {
	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	if "" == deviceInfo.ObjectName && "" == deviceInfo.ObjectID {
		blog.Errorf("check net device object ID, empty bk_obj_id and bk_obj_name")
		return false, defErr.Errorf(common.CCErrCommParamsLostField, common.BKObjIDField)
	}

	// netDeviceObjID is a net device objID from deviceInfo.ObjectName
	var (
		netDeviceObjID string
		err            error
	)

	// one of objectName and objectID must not be empty
	// if objectName is not empty, get net device objectID from deviceInfo.ObjectName
	if "" != deviceInfo.ObjectName {
		netDeviceObjID, err = lgc.getNetDeviceObjIDFromObjName(deviceInfo.ObjectName, pheader)
		if nil != err {
			blog.Errorf("check net device object ID, error: %v", err)
			return false, err
		}

		// if deviceInfo.ObjectName is not empty, netDeviceObjID must not be empty
		if "" == netDeviceObjID {
			blog.Errorf("check net device object ID, get empty objID from object name [%s]", deviceInfo.ObjectName)
			return false, defErr.Errorf(common.CCErrCollectObjNameNotNetDevice)
		}

		// return true if objectID is empty or netDeviceObjID and objectID are the same
		switch {
		case "" == deviceInfo.ObjectID:
			deviceInfo.ObjectID = netDeviceObjID
			return true, nil
		case netDeviceObjID == deviceInfo.ObjectID:
			return true, nil
		}

		// objectID are not empty, netDeviceObjID and objectID are the same
		blog.Errorf("check net device object ID, get objID [%s] from object name [%s], different between object ID [%s]",
			netDeviceObjID, deviceInfo.ObjectID, deviceInfo.ObjectName)
		return false, defErr.Errorf(common.CCErrCollectDiffObjIDAndName)
	}

	// objectName from function parameter is empty and objectID must not be empty
	// check if objectID is one of net collect device, such as bk_router, bk_switch...
	isNetDeviceObjID, err := lgc.checkIfObjIDIsNetDevice(deviceInfo.ObjectID, pheader)
	if nil != err {
		return false, err
	}

	return isNetDeviceObjID, nil
}

// get value of net device bk_obj_id from bk_obj_name
// return error if bk_obj_name is empty or bk_obj_name is not net device
func (lgc *Logics) getNetDeviceObjIDFromObjName(objectName string, pheader http.Header) (string, error) {
	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	if "" == objectName {
		return "", defErr.Errorf(common.CCErrCollectObjNameNotNetDevice)
	}

	cond := make(map[string]interface{}, 0)
	cond[common.BKObjNameField] = objectName
	cond[common.BKClassificationIDField] = common.BKNetwork

	objResult, err := lgc.CoreAPI.ObjectController().Meta().SelectObjects(context.Background(), pheader, cond)
	if nil != err {
		blog.Errorf("check net device object ID, search objectName fail, %v", err)
		return "", defErr.Errorf(common.CCErrObjectSelectInstFailed)
	}

	if !objResult.Result {
		blog.Errorf("check net device object ID, errors: %s", objResult.ErrMsg)
		return "", defErr.Errorf(objResult.Code)
	}

	if nil == objResult.Data || 0 == len(objResult.Data) {
		blog.Errorf("check net device object ID, object Name[%s] is not exist", objectName)
		return "", defErr.Errorf(common.CCErrCollectObjNameNotNetDevice)
	}

	return objResult.Data[0].ObjectID, nil
}

// check if net device name exist
func (lgc *Logics) checkIfNetDeviceNameExist(deviceName string, ownerID string) (bool, error) {
	queryParams := common.KvMap{common.BKDeviceNameField: deviceName, common.BKOwnerIDField: ownerID}

	rowCount, err := lgc.Instance.GetCntByCondition(common.BKTableNameNetcollectDevice, queryParams)
	if nil != err {
		blog.Errorf("check if net device name exist, query device fail, error information is %v, params:%v",
			err, queryParams)
		return false, err
	}

	if 0 != rowCount {
		blog.V(4).Infof("check if net device name exist, bk_device_name is [%s] device is exist", deviceName)
		return true, nil
	}

	return false, nil
}

// check if objID is net device object
func (lgc *Logics) checkIfObjIDIsNetDevice(objID string, pheader http.Header) (bool, error) {
	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	if "" == objID {
		return false, nil
	}

	cond := make(map[string]interface{}, 0)
	cond[common.BKObjIDField] = objID

	objResult, err := lgc.CoreAPI.ObjectController().Meta().SelectObjects(context.Background(), pheader, cond)
	if nil != err {
		blog.Errorf("check if objID is net device object, search objID [%s] fail, %v", objID, err)
		return false, defErr.Errorf(common.CCErrObjectSelectInstFailed)
	}

	if !objResult.Result {
		blog.Errorf("check if objID is net device object, errors: %s", objResult.ErrMsg)
		return false, defErr.Errorf(objResult.Code)
	}

	return nil != objResult.Data && 0 < len(objResult.Data), nil
}

// get net device obj ID
func (lgc *Logics) getNetDeviceObjIDsByCond(objCond map[string]interface{}, pheader http.Header) ([]string, error) {
	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	objIDs := []string{}

	if _, ok := objCond[common.BKObjNameField]; ok {
		objCond[common.BKClassificationIDField] = common.BKNetwork
		objResult, err := lgc.CoreAPI.ObjectController().Meta().SelectObjects(context.Background(), pheader, objCond)
		if nil != err {
			blog.Errorf("check net device object ID, search objectName fail, %v", err)
			return nil, err
		}

		if !objResult.Result {
			blog.Errorf("check net device object ID, errors: %s", objResult.ErrMsg)
			return nil, defErr.Errorf(objResult.Code)
		}

		if nil != objResult.Data {
			for _, data := range objResult.Data {
				objIDs = append(objIDs, data.ObjectID)
			}
		}
	}

	return objIDs, nil
}
