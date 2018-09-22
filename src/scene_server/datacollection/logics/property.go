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
	"net/http"
	"time"

	restful "github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/blog"
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

func (lgc *Logics) SearchProperty(req *restful.Request, resp *restful.Response) {

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
		propertyInfo.Period, err = util.FormatPeriod(propertyInfo.Period)
		if nil != err {
			blog.Errorf("add net collect property, format perid [%s] fail, error: %v", err)
			return -1, defErr.Errorf(common.CCErrCollectPeridFormatFail)
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

func (lgc *Logics) isValidAction(action string) bool {
	return common.ActionGet == action || common.ActionWalk == action
}
