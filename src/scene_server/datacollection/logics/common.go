package logics

import (
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"
)

// by checking if bk_obj_id and bk_obj_name function parameter are valid net device object or not
// one of bk_obj_id and bk_obj_name can be empty and will return both value if no error
func (lgc *Logics) checkIfNetObject(objID string, objName string, pheader http.Header) (string, string, bool, error) {
	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	if "" == objName && "" == objID {
		blog.Errorf("check net device object, empty bk_obj_id and bk_obj_name")
		return "", "", false, defErr.Errorf(common.CCErrCommParamsLostField, common.BKObjIDField)
	}

	// netDeviceObjID is a net device objID from objName
	var (
		netDeviceObjID string
		err            error
	)

	// one of objectName and objectID must not be empty
	// if objectName is not empty, get net device objectID from objName
	if "" != objName {
		netDeviceObjID, err = lgc.getNetDeviceObjIDFromObjName(objName, pheader)
		if nil != err {
			blog.Errorf("check net device object, error: %v", err)
			return "", "", false, err
		}

		// if objName is not empty, netDeviceObjID must not be empty
		if "" == netDeviceObjID {
			blog.Errorf("check net device object, get empty objID from object name [%s]", objName)
			return "", "", false, defErr.Errorf(common.CCErrCollectObjNameNotNetDevice)
		}

		// return true if objectID is empty or netDeviceObjID and objectID are the same
		switch {
		case "" == objID:
			objID = netDeviceObjID
			return objID, objName, true, nil
		case netDeviceObjID == objID:
			return objID, objName, true, nil
		}

		// objectID are not empty, netDeviceObjID and objectID are the same
		blog.Errorf("check net device object, get objID [%s] from object name [%s], different between object ID [%s]",
			netDeviceObjID, objID, objName)
		return "", "", false, defErr.Errorf(common.CCErrCollectDiffObjIDAndName)
	}

	// objectName from function parameter is empty and objectID must not be empty
	// check if objectID is one of net collect device, such as bk_router, bk_switch...
	isNetDeviceObjID, err := lgc.checkIfObjIDIsNetDevice(objID, pheader)
	if nil != err {
		return "", "", false, err
	}

	return objID, objName, isNetDeviceObjID, nil
}
