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
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// INVALIDID invalid id used as return value
const INVALIDID uint64 = 0

// by checking if bk_obj_id and bk_obj_name function parameter are valid net device object or not
// one of bk_obj_id and bk_obj_name can be empty and will return both bk_obj_id if no error
func (lgc *Logics) checkNetObject(pheader http.Header, objID string, objName string) (string, string, error) {
	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	if "" == objName && "" == objID {
		blog.Errorf("[NetCollect] check net device object, empty bk_obj_id and bk_obj_name")
		return "", "", defErr.Errorf(common.CCErrCommParamsNeedSet, common.BKObjIDField)
	}

	objCond := map[string]interface{}{
		common.BKClassificationIDField: common.BKNetwork,
	}

	if "" != objName {
		objCond[common.BKObjNameField] = objName
	}
	if "" != objID {
		objCond[common.BKObjIDField] = objID
	}

	objResult, err := lgc.CoreAPI.CoreService().Model().ReadModel(context.Background(), pheader, &meta.QueryCondition{Condition: objCond})
	if nil != err {
		blog.Errorf("[NetCollect] check net device object, get net device object fail, error: %v, condition [%#v]", err, objCond)
		return "", "", defErr.Errorf(common.CCErrObjectSelectInstFailed)
	}

	if !objResult.Result {
		blog.Errorf("[NetCollect] check net device object, errors: %s, condition [%#v]", objResult.ErrMsg, objCond)
		return "", "", defErr.New(objResult.Code, objResult.ErrMsg)
	}

	if 0 == len(objResult.Data.Info) {
		blog.Errorf("[NetCollect] check net device object, device object is not exist, condition [%#v]", objCond)
		return "", "", defErr.Errorf(common.CCErrCollectObjIDNotNetDevice)
	}

	return objResult.Data.Info[0].Spec.ObjectID, objResult.Data.Info[0].Spec.ObjectName, nil
}

// by checking if bk_property_id and bk_property_name function parameter are valid net device object property or not
// one of bk_property_id and bk_property_name can be empty and will return bk_property_id value if no error
func (lgc *Logics) checkNetObjectProperty(pheader http.Header, netDeviceObjID, propertyID, propertyName string) (string, error) {
	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	if "" == netDeviceObjID {
		blog.Errorf("[NetCollect] check net device object, empty bk_obj_id")
		return "", defErr.Errorf(common.CCErrCommParamsNeedSet, common.BKObjIDField)
	}

	if "" == propertyName && "" == propertyID {
		blog.Errorf("[NetCollect] check net device object, empty bk_property_id and bk_property_name")
		return "", defErr.Errorf(common.CCErrCommParamsNeedSet, common.BKPropertyIDField)
	}

	propertyCond := map[string]interface{}{
		common.BKObjIDField: netDeviceObjID}

	if "" != propertyName {
		propertyCond[common.BKPropertyNameField] = propertyName
	}
	if "" != propertyID {
		propertyCond[common.BKPropertyIDField] = propertyID
	}

	attrResult, err := lgc.CoreAPI.CoreService().Model().ReadModelAttrByCondition(context.Background(), pheader, &meta.QueryCondition{Condition: propertyCond})
	if nil != err {
		blog.Errorf("[NetCollect] get object attribute fail, error: %v, condition [%#v]", err, propertyCond)
		return "", defErr.Errorf(common.CCErrTopoObjectAttributeSelectFailed)
	}
	if !attrResult.Result {
		blog.Errorf("[NetCollect] check net device object property, errors: %s", attrResult.ErrMsg)
		return "", defErr.New(attrResult.Code, attrResult.ErrMsg)
	}

	if 0 == len(attrResult.Data.Info) {
		blog.Errorf("[NetCollect] check net device object property, property is not exist, condition [%#v]", propertyCond)
		return "", defErr.Errorf(common.CCErrCollectNetDeviceObjPropertyNotExist)
	}

	return attrResult.Data.Info[0].PropertyID, nil
}

// by checking if bk_device_id and bk_device_name function parameter are valid net device or not
// one of bk_device_id and bk_device_name can be empty and will return bk_device_id and bk_obj_id value if no error
// bk_obj_id is used to check property
func (lgc *Logics) checkNetDeviceExist(pheader http.Header, deviceID uint64, deviceName string) (uint64, string, error) {
	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	if "" == deviceName && 0 == deviceID {
		blog.Errorf("[NetCollect] check net device exist fail, empty device_id and device_name")
		return 0, "", defErr.Errorf(common.CCErrCommParamsNeedSet, common.BKDeviceIDField)
	}

	deviceCond := map[string]interface{}{common.BKOwnerIDField: util.GetOwnerID(pheader)}

	if "" != deviceName {
		deviceCond[common.BKDeviceNameField] = deviceName
	}
	if 0 != deviceID {
		deviceCond[common.BKDeviceIDField] = deviceID
	}

	deviceData := meta.NetcollectDevice{}
	if err := lgc.db.Table(common.BKTableNameNetcollectDevice).Find(deviceCond).Fields(common.BKDeviceIDField, common.BKObjIDField).
		One(lgc.ctx, &deviceData); nil != err {

		blog.Errorf("[NetCollect] check net device exist fail, error: %v, condition [%#v]", err, deviceCond)

		if lgc.db.IsNotFoundError(err) {
			return 0, "", defErr.Error(common.CCErrCollectNetDeviceGetFail)
		}
		return 0, "", err
	}

	return deviceData.DeviceID, deviceData.ObjectID, nil
}

// get net property id by device ID and property ID
func (lgc *Logics) getNetPropertyID(propertyID string, deviceID uint64, ownerID string) (uint64, error) {
	queryParams := map[string]interface{}{
		common.BKDeviceIDField:   deviceID,
		common.BKPropertyIDField: propertyID,
		common.BKOwnerIDField:    ownerID,
	}

	result := meta.NetcollectProperty{}
	if err := lgc.db.Table(common.BKTableNameNetcollectProperty).Find(queryParams).Fields(common.BKNetcollectPropertyIDField).
		One(lgc.ctx, &result); nil != err {

		blog.Errorf(
			"[NetCollect] get net property ID by propertyID and deviceID, error: %v, params: [%#+v]", err, queryParams)
		return 0, err
	}
	blog.Errorf(
		"[NetCollect] get net property ID by propertyID and deviceID, params: [%#+v]", queryParams)

	return result.NetcollectPropertyID, nil
}

// ConvertStringToID check param ID is a num string and convert to num
func (lgc *Logics) ConvertStringToID(stringID string) (int64, error) {
	if "" == stringID || "0" == stringID {
		return 0, fmt.Errorf("invalid stringID")
	}

	ID, err := strconv.ParseInt(stringID, 10, 64)
	if nil != err {
		return 0, err
	}

	return ID, nil
}
