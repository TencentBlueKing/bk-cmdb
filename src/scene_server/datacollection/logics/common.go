package logics

import (
	"context"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"
)

// by checking if bk_obj_id and bk_obj_name function parameter are valid net device object or not
// one of bk_obj_id and bk_obj_name can be empty and will return both value if no error
func (lgc *Logics) checkNetObject(objID string, objName string, pheader http.Header) (string, string, error) {
	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	if "" == objName && "" == objID {
		blog.Errorf("check net device object, empty bk_obj_id and bk_obj_name")
		return "", "", defErr.Errorf(common.CCErrCommParamsLostField, common.BKObjIDField)
	}

	objCond := map[string]interface{}{
		common.BKOwnerIDField:          util.GetOwnerID(pheader),
		common.BKClassificationIDField: common.BKNetwork}

	if "" != objName {
		objCond[common.BKObjNameField] = objName
	}
	if "" != objID {
		objCond[common.BKObjIDField] = objID
	}

	objResult, err := lgc.CoreAPI.ObjectController().Meta().SelectObjects(context.Background(), pheader, objCond)
	if nil != err {
		blog.Errorf("check net device object, get net device object fail, error: %v, condition [%#v]", err, objCond)
		return "", "", defErr.Errorf(common.CCErrObjectSelectInstFailed)
	}

	if !objResult.Result {
		blog.Errorf("check net device object, errors: %s, condition [%#v]", objResult.ErrMsg, objCond)
		return "", "", defErr.Errorf(objResult.Code)
	}

	if nil == objResult.Data || 0 == len(objResult.Data) {
		blog.Errorf("check net device object, device object is not exist, condition [%#v]", objCond)
		return "", "", defErr.Errorf(common.CCErrCollectObjNameNotNetDevice)
	}

	return objResult.Data[0].ObjectID, objResult.Data[0].ObjectName, nil
}

// by checking if bk_obj_id and bk_obj_name function parameter are valid net device object or not
// one of bk_obj_id and bk_obj_name can be empty and will return both value if no error
func (lgc *Logics) checkNetObjectProperty(netDeviceObjID, propertyName, propertyID string, pheader http.Header) (string, string, error) {
	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	if "" == netDeviceObjID {
		blog.Errorf("check net device object, empty bk_obj_id")
		return "", "", defErr.Errorf(common.CCErrCommParamsLostField, common.BKObjIDField)
	}

	if "" == propertyName && "" == propertyID {
		blog.Errorf("check net device object, empty bk_property_id and bk_property_name")
		return "", "", defErr.Errorf(common.CCErrCommParamsLostField, common.BKPropertyIDField)
	}

	propertyCond := map[string]interface{}{
		common.BKOwnerIDField: util.GetOwnerID(pheader),
		common.BKObjIDField:   netDeviceObjID}

	if "" == propertyName {
		propertyCond[common.BKPropertyNameField] = propertyName
	}
	if "" == propertyID {
		propertyCond[common.BKPropertyIDField] = propertyID
	}

	attrResult, err := lgc.CoreAPI.ObjectController().Meta().SelectObjectAttWithParams(context.Background(), pheader, propertyCond)
	if nil != err {
		blog.Errorf("get object attribute fail, error: %v, condition [%#v]", err, propertyCond)
		return "", "", err
	}
	if !attrResult.Result {
		blog.Errorf("check net device object property, errors: %s", attrResult.ErrMsg)
		return "", "", defErr.Errorf(attrResult.Code)
	}

	if nil == attrResult.Data || 0 == len(attrResult.Data) {
		blog.Errorf("check net device object property, property is not exist, condition [%#v]", propertyCond)
		return "", "", defErr.Errorf(common.CCErrCollectNetDeviceObjPropertyNotExist)
	}

	return attrResult.Data[0].PropertyID, attrResult.Data[0].PropertyName, nil
}

// by checking if bk_obj_id and bk_obj_name function parameter are valid net device object or not
// one of bk_obj_id and bk_obj_name can be empty and will return both value if no error
func (lgc *Logics) checkNetDeviceExist(deviceName, deviceID string, pheader http.Header) (string, string, error) {
	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	if "" == deviceName && "" == deviceID {
		blog.Errorf("check net device exist fail, empty device_id and device_name")
		return "", "", defErr.Errorf(common.CCErrCommParamsLostField, common.BKDeviceIDField)
	}

	deviceCond := map[string]interface{}{common.BKOwnerIDField: util.GetOwnerID(pheader)}

	if "" == deviceName {
		deviceCond[common.BKDeviceNameField] = deviceName
	}
	if "" == deviceID {
		deviceCond[common.BKDeviceIDField] = deviceID
	}

	attrResult := map[string]interface{}{}
	err := lgc.Instance.GetOneByCondition(
		common.BKTableNameNetcollectDevice, []string{common.BKDeviceIDField, common.BKDeviceNameField}, deviceCond, &attrResult)
	if nil != err {
		blog.Errorf("get object attribute fail, error: %v, condition [%#v]", err, deviceCond)
		return "", "", err
	}

	return attrResult[common.BKDeviceIDField].(string), attrResult[common.BKDeviceNameField].(string), nil
}
