package logics

import (
	"net/http"
	"regexp"
	"strconv"
	"time"

	restful "github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
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

	if err = lgc.checkIfNetDeviceExist(&propertyInfo, pheader); nil != err {
		blog.Errorf("add net collect property fail, error: %v", err)
		return -1, err
	}

	if err = lgc.checkIfNetProperty(&propertyInfo, pheader); nil != err {
		blog.Errorf("add net collect property fail, error: %v", err)
		return -1, err
	}

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

func (lgc *Logics) checkIfNetProperty(propertyInfo *meta.NetcollectProperty, pheader http.Header) error {
	var err error
	propertyInfo.PropertyID, err = lgc.checkNetObjectProperty(propertyInfo.ObjectID, propertyInfo.PropertyID, propertyInfo.PropertyName, pheader)
	return err
}

func (lgc *Logics) checkIfNetDeviceExist(propertyInfo *meta.NetcollectProperty, pheader http.Header) error {
	var err error
	propertyInfo.DeviceID, propertyInfo.ObjectID, err = lgc.checkNetDeviceExist(propertyInfo.DeviceID, propertyInfo.DeviceName, pheader)
	return err
}

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

const periodRegexp = "^\\d*[DHMS]$"

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
