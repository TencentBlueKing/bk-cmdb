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
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

type ReturnResult struct {
	Result   bool   `json:"result"`
	ErrMsg   string `json:"error_msg"`
	DeviceID int64  `json:"device_id"`
}

func (lgc *Logics) AddDevices(pheader http.Header, deviceInfoList []meta.NetcollectDevice) ([]ReturnResult, bool) {
	ownerID := util.GetOwnerID(pheader)

	resultList := make([]ReturnResult, 0)
	hasError := false

	for _, deviceInfo := range deviceInfoList {
		errMsg := ""
		result := true

		deviceID, err := lgc.addDevice(deviceInfo, pheader, ownerID)
		if nil != err {
			errMsg = err.Error()
			result = false
		}

		resultList = append(resultList, ReturnResult{result, errMsg, deviceID})
		if nil != err && !hasError {
			hasError = true
		}
	}

	return resultList, hasError
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
	isExist, err := lgc.isNetDeviceNameExist(deviceInfo.DeviceName, ownerID)
	if nil != err {
		blog.Errorf("add net device fail, error: %v", err)
		return -1, defErr.Errorf(common.CCErrCollectNetDeviceCreateFail)
	}
	if isExist {
		blog.Errorf("add net device fail, error: duplicate device_name")
		return -1, defErr.Errorf(common.CCErrCommDuplicateItem)
	}

	// check if bk_object_id and bk_object_name are net device
	isNetDevice, err := lgc.isNetDeviceObjectID(&deviceInfo, pheader)
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
	deviceInfo.CreateTime = new(time.Time)
	*deviceInfo.CreateTime = time.Now()
	deviceInfo.LastTime = new(time.Time)
	*deviceInfo.LastTime = time.Now()
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

// check if bk_obj_id and bk_obj_name function parameter are valid or not
func (lgc *Logics) isNetDeviceObjectID(deviceInfo *meta.NetcollectDevice, pheader http.Header) (bool, error) {
	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	if "" == deviceInfo.ObjectName && "" == deviceInfo.ObjectID {
		blog.Errorf("check net device object ID, empty bk_obj_id and bk_obj_name")
		return false, defErr.Errorf(common.CCErrCommParamsLostField, common.BKObjIDField)
	}

	// one of objectName and objectID must not be empty
	objectIDFromObjectName, err := lgc.getObjectIDFromObjectName(deviceInfo.ObjectName, pheader)
	if nil != err {
		blog.Errorf("check net device object ID, error: %v", err)
		return false, err
	}
	if "" != objectIDFromObjectName && "" != deviceInfo.ObjectID && objectIDFromObjectName != deviceInfo.ObjectID {
		blog.Errorf("check net device object ID, get objID [%s] from object name [%s], different between object ID [%s]",
			objectIDFromObjectName, deviceInfo.ObjectID, deviceInfo.ObjectName)
		return false, defErr.Errorf(common.CCErrCollectDiffObjIDAndName)
	}

	// one of objectName and objectID must not be empty
	// if objectID from function parameter is empty, get objectID from objectName
	if "" == deviceInfo.ObjectID {
		deviceInfo.ObjectID = objectIDFromObjectName
	}

	// check if objectID is one of net collect device, such as bk_router, bk_switch...
	return meta.IsNetDeviceObject(deviceInfo.ObjectID), nil
}

// get value of bk_obj_id by bk_obj_name
func (lgc *Logics) getObjectIDFromObjectName(objectName string, pheader http.Header) (string, error) {
	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	if "" == objectName {
		return "", nil
	}
	cond := make(map[string]interface{}, 0)
	// cond := condition.CreateCondition().Field(common.BKObjNameField).Eq(objectName)
	cond[common.BKObjNameField] = objectName

	objResult, err := lgc.CoreAPI.ObjectController().Meta().SelectObjects(context.Background(), pheader, cond)
	if nil != err {
		blog.Errorf("check net device object ID, search objectName fail, %v", err)
		return "", defErr.Errorf(common.CCErrObjectSelectInstFailed)
	}

	if !objResult.Result {
		blog.Errorf("check net device object ID, errors: %s", objResult.ErrMsg)
		return "", defErr.Errorf(common.CCErrObjectSelectInstFailed)
	}

	if nil == objResult.Data || 0 == len(objResult.Data) {
		blog.Errorf("check net device object ID, [%s] is not exist", objectName)
		return "", nil
	}

	return objResult.Data[0].ObjectID, nil
}

// check if net device name exist
func (lgc *Logics) isNetDeviceNameExist(deviceName string, ownerID string) (bool, error) {
	queryParams := common.KvMap{common.BKDeviceNameField: deviceName, common.BKOwnerIDField: ownerID}

	rowCount, err := lgc.Instance.GetCntByCondition(common.BKTableNameNetcollectDevice, queryParams)
	if nil != err {
		blog.Errorf("check if net device name exist, query device fail, error information is %v, params:%v",
			err, queryParams)
		return false, err
	}

	if 0 != rowCount {
		blog.Infof("check if net device name exist, bk_device_name is [%s] device is exist", deviceName)
		return true, nil
	}

	return false, nil
}
