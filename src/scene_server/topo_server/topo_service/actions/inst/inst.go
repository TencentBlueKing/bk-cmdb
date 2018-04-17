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

package inst

import (
	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/bkbase"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/errors"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/topo_service/actions/object"
	"configcenter/src/scene_server/topo_server/topo_service/manager"
	"configcenter/src/scene_server/validator"
	"configcenter/src/source_controller/api/auditlog"
	"configcenter/src/source_controller/api/metadata"
	api "configcenter/src/source_controller/api/object"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
	restful "github.com/emicklei/go-restful"
)

var inst = &instAction{}

// instAction
type instAction struct {
	base.BaseAction
	objcli         *api.Client
	metaHelperFunc *object.HelperFunction
}

func init() {

	// register action
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/inst/{owner_id}/{obj_id}", Params: nil, Handler: inst.CreateInst})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/inst/{owner_id}/{obj_id}/{inst_id}", Params: nil, Handler: inst.DeleteInst})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/inst/{owner_id}/{obj_id}/{inst_id}", Params: nil, Handler: inst.UpdateInst})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/inst/search/{owner_id}/{obj_id}", Params: nil, Handler: inst.SelectInsts})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/inst/association/search/owner/{owner_id}/object/{obj_id}", Params: nil, Handler: inst.SelectInstsByAssociation})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/inst/search/{owner_id}/{obj_id}/{inst_id}", Params: nil, Handler: inst.SelectInst})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/inst/search/topo/owner/{owner_id}/object/{object_id}/inst/{inst_id}", Params: nil, Handler: inst.SelectTopo})

	// create cc object
	inst.CreateAction()
	inst.metaHelperFunc = object.CreateHelperFunction()
	inst.objcli = api.NewClient("")
}

func (cli *instAction) subCreateInst(req *restful.Request, defErr errors.DefaultCCErrorIf, targetInput map[string]interface{}, ownerID, objID string, isBatch bool, asstDes []api.ObjAsstDes, attDes []api.ObjAttDes) (int, interface{}, bool, error) {

	nonExistsFiled := make([]api.ObjAttDes, 0)
	ignorItems := make([]string, 0)
	for _, item := range attDes {
		if _, ok := targetInput[item.PropertyID]; !ok {
			nonExistsFiled = append(nonExistsFiled, item)
			ignorItems = append(ignorItems, item.PropertyID)
		}
	}
	ignorItems = append(ignorItems, common.BKInstParentStr)
	ignorItems = append(ignorItems, common.BKAppIDField)
	blog.Debug("the ignore items:%+v", ignorItems)
	valid := validator.NewValidMapWithKeyFileds(ownerID, objID, cli.CC.ObjCtrl(), ignorItems, defErr)
	user := util.GetActionUser(req)
	isUpdate := false
	blog.Debug("the non exists filed items:%+v", nonExistsFiled)

	// extract the data for the associated field
	asstFieldVal := make([]interface{}, 0)
	for idxItem, item := range asstDes {
		if inputVal, ok := targetInput[item.ObjectAttID]; ok {
			switch t := inputVal.(type) {
			case string:
				asstIDS := strings.Split(t, ",")
				for _, id := range asstIDS {
					if 0 == len(id) {
						continue
					}
					iID, iIDErr := util.GetInt64ByInterface(id)
					if nil != iIDErr {
						blog.Error("can not convert the data (%s) into int64, error info is %s", id, iIDErr.Error())
						continue
					}

					asstInst := metadata.InstAsst{}
					id, err := cli.CC.InstCli.GetIncID(asstInst.TableName())
					asstInst.ID = id
					if nil != err {
						blog.Error("faild to create id, error info is %s", err.Error())
					}
					asstInst.AsstInstID = iID
					asstInst.AsstObjectID = asstDes[idxItem].AsstObjID
					asstInst.ObjectID = objID
					asstFieldVal = append(asstFieldVal, asstInst)
				}

			default:
				blog.Warnf("the target data (%v) type is not a string ", t)
			}
		}
	}

	// check
	_, err := valid.ValidMap(targetInput, common.ValidCreate, 0)
	targetMethod := common.HTTPSelectPost
	input := make(map[string]interface{})
	switch e := err.(type) {
	case nil:
		// clear the association field
		if isBatch {
			for _, item := range asstDes {
				if _, ok := targetInput[item.ObjectAttID]; ok {
					delete(targetInput, item.ObjectAttID)
				}
			}

			// set the nonexist
			for _, j := range nonExistsFiled {
				propertyID := j.PropertyID
				fieldType := j.PropertyType
				switch fieldType {
				case common.FiledTypeSingleChar:
					targetInput[propertyID] = ""
				case common.FiledTypeLongChar:
					targetInput[propertyID] = ""
				default:
					targetInput[propertyID] = nil
				}
			}
		}

		input = targetInput

	case errors.CCErrorCoder:
		if e.GetCode() == common.CCErrCommDuplicateItem && isBatch {

			isUpdate = true
			condition := make(map[string]interface{})
			condition[common.BKOwnerIDField] = ownerID
			condition[common.BKObjIDField] = objID
			condition[common.BKInstNameField] = targetInput[common.BKInstNameField]

			if _, ok := targetInput[common.BKInstNameField]; !ok {
				blog.Error("lost the 'InstName' field, the error data is %+v", targetInput)
				return http.StatusBadRequest, nil, isUpdate, defErr.Errorf(common.CCErrCommParamsLostField, common.BKInstNameField)
			}

			if _, err = valid.ValidMap(targetInput, common.ValidUpdate, 0); nil != err {
				switch e := err.(type) {
				case nil:
					break
				case errors.CCErrorCoder:
					if e.GetCode() == common.CCErrCommDuplicateItem {
						break
					}
					blog.Error("failed valid the input data, error info is %s", err.Error())
					return http.StatusBadRequest, nil, isUpdate, err
				default:
					blog.Error("failed valid the input data, error info is %s", err.Error())
					return http.StatusBadRequest, nil, isUpdate, err
				}
			}

			targetMethod = common.HTTPUpdate
			// clear the association field
			if isBatch {
				for _, item := range asstDes {
					if _, ok := targetInput[item.ObjectAttID]; ok {
						delete(targetInput, item.ObjectAttID)
					}
				}
			}

			input["data"] = targetInput
			input["condition"] = condition

		} else {
			blog.Error("failed valid the input data, error info is %s", err.Error())
			return http.StatusBadRequest, nil, isUpdate, err
		}
	default:
		blog.Error("failed valid the input data, error info is %s", err.Error())
		return http.StatusBadRequest, nil, isUpdate, err
	}

	// take snapshot before operation if is update
	preData := map[string]interface{}{}
	var instID int
	var retStrErr int
	if targetMethod == common.HTTPUpdate {
		preData, retStrErr = cli.getInstDeteilByCondition(req, objID, ownerID, input["condition"].(map[string]interface{}))
		if common.CCSuccess != retStrErr {
			blog.Errorf("get inst detail error: %v", retStrErr)
			return http.StatusInternalServerError, nil, isUpdate, defErr.Error(retStrErr)
		}

		instID, _ = strconv.Atoi(fmt.Sprint(preData[common.BKInstIDField]))
	}

	// set default InstaName value if not set
	if _, ok := targetInput[common.BKInstNameField]; !ok {
		searchObjIDCond := make(map[string]interface{})
		searchObjIDCond[common.BKObjIDField] = objID
		searchObjIDCond[common.BKOwnerIDField] = ownerID
		searchObjIDCondVal, _ := json.Marshal(searchObjIDCond)
		cli.objcli.SetAddress(cli.CC.ObjCtrl())
		objName := objID
		rstItems, rstErr := cli.objcli.SearchMetaObject(searchObjIDCondVal)
		if nil != rstErr {
			blog.Error("failed to fetch the object, error info is %s", rstErr.Error())
		} else if len(rstItems) > 0 {
			objName = rstItems[0].ObjectName
		}
		input[common.BKInstNameField] = fmt.Sprintf("默认 %s)", objName)
	}

	input[common.BKOwnerIDField] = ownerID
	input[common.BKObjIDField] = objID
	input[common.BKDefaultField] = 0
	input[common.CreateTimeField] = util.GetCurrentTimeStr()

	inputJSON, jsErr := json.Marshal(input)
	if nil != jsErr {
		blog.Error("the input json is invalid, error info is %s", jsErr.Error())
		return http.StatusBadRequest, nil, isUpdate, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	cURL := cli.CC.ObjCtrl() + "/object/v1/insts/object"
	blog.Debug("inst:%v", string(inputJSON))

	instRes, err := httpcli.ReqHttp(req, cURL, targetMethod, inputJSON)
	if nil != err {
		blog.Error("create inst failed, errors:%s", err.Error())
		return http.StatusInternalServerError, nil, isUpdate, defErr.Error(common.CCErrTopoInstCreateFailed)
	}

	rsp, ok := cli.IsSuccess([]byte(instRes))
	if !ok {
		return http.StatusInternalServerError, nil, isUpdate, fmt.Errorf("%+v", rsp.Message)
	}

	{
		// save change log
		if targetMethod == common.HTTPSelectPost {
			instID = int(gjson.Get(instRes, "data."+common.BKInstIDField).Int())
		}
		headers := []metadata.Header{}
		for _, item := range attDes {
			headers = append(headers, metadata.Header{
				PropertyID:   item.PropertyID,
				PropertyName: item.PropertyName,
			})
		}

		blog.Infof("new instID = %v", instID)
		if instID == 0 {
			blog.Infof("response data: %s", instRes)
		}

		curData, retStrErr := cli.getInstDetail(req, int(instID), objID, ownerID)
		if common.CCSuccess != retStrErr {
			blog.Errorf("get inst detail error: %v, instid(%v) objid(%s) ownerid(%s)", retStrErr, instID, objID, ownerID)
			return http.StatusInternalServerError, nil, isUpdate, defErr.Error(retStrErr)
		}
		auditContent := metadata.Content{
			PreData: preData,
			CurData: curData,
			Headers: headers,
		}
		if targetMethod == common.HTTPSelectPost {
			auditlog.NewClient(cli.CC.AuditCtrl()).AuditObjLog(instID, auditContent, "create inst", objID, ownerID, "0", user, auditoplog.AuditOpTypeAdd)
		} else {
			auditlog.NewClient(cli.CC.AuditCtrl()).AuditObjLog(instID, auditContent, "update inst", objID, ownerID, "0", user, auditoplog.AuditOpTypeModify)
		}

	}

	// set the inst association table
	for idxItem, item := range asstFieldVal {
		switch t := item.(type) {
		case metadata.InstAsst:
			t.InstID = int64(instID)
			asstFieldVal[idxItem] = t
		}
	}
	if err := cli.createInstAssociation(asstFieldVal); nil != err {
		blog.Errorf("failed to create the inst association, error info is %s ", err.Error())
	}

	return http.StatusOK, rsp.Data, isUpdate, nil
}

// CreateInst create a inst
func (cli *instAction) CreateInst(req *restful.Request, resp *restful.Response) {

	blog.Info("create inst")

	// get the language
	language := util.GetActionLanguage(req)
	// get error object by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	// logics
	cli.CallResponseEx(func() (int, interface{}, error) {

		//create default module
		value, err := ioutil.ReadAll(req.Request.Body)
		if nil != err {
			blog.Error("read request body failed, error:%s", err.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		if 0 == len(value) {
			blog.Error("the request body is empty")
			return http.StatusBadRequest, "", defErr.Error(common.CCErrCommHTTPBodyEmpty)
		}

		js, err := simplejson.NewJson(value)
		if nil != err {
			blog.Error("failed to unmarshal the json, error info is %s", err.Error())
			return http.StatusBadRequest, "", defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		input, jsErr := js.Map()
		if nil != jsErr {
			blog.Error("the input json is invalid, error info is %s", jsErr.Error())
			return http.StatusBadRequest, "", defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		ownerID := req.PathParameter("owner_id")
		objID := req.PathParameter("obj_id")

		// the batch data structure define map[int]map[string]interface{}
		innerBatchInfo := &struct {
			// BatchInfo batch info
			BatchInfo *map[int]map[string]interface{} `json:"BatchInfo"`
		}{}

		if _, batchInfoOK := input["BatchInfo"]; batchInfoOK {

			delete(input, "BatchInfo")
			if jsErr := json.Unmarshal(value, innerBatchInfo); nil != jsErr {
				blog.Error("the data structure has changed, error info is %s ", jsErr.Error())
				return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
			}
			blog.Debug("the batch info %+v", *innerBatchInfo.BatchInfo)
		}

		cli.objcli.SetAddress(cli.CC.ObjCtrl())

		// define create inst function
		createFunc := cli.subCreateInst

		// get object attributes fields
		att := map[string]interface{}{}
		att[common.BKObjIDField] = objID
		att[common.BKOwnerIDField] = ownerID
		searchData, jsErr := json.Marshal(att)
		if nil != jsErr {
			blog.Error("failed to marshal the data[%+v], error info is %s", searchData, jsErr.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommJSONMarshalFailed)
		}
		attdes, restErr := cli.objcli.SearchMetaObjectAtt(searchData)
		if nil != restErr {
			blog.Error("failed to read the object att, error is %s ", restErr.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoInstSelectFailed)
		}

		// get association fields
		asst := map[string]interface{}{}
		asst[common.BKOwnerIDField] = ownerID
		asst[common.BKObjIDField] = objID
		searchData, _ = json.Marshal(asst)
		cli.objcli.SetAddress(cli.CC.ObjCtrl())
		asstDes, asstErr := cli.objcli.SearchMetaObjectAsst(searchData)
		if nil != asstErr {
			blog.Error("failed to search the obj asst, search condition(%+v) error info is %s", asst, asstErr.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoInstCreateFailed)
		}

		// create batch
		if nil != innerBatchInfo.BatchInfo {

			type batchResult struct {
				Errors       []string `json:"error"`
				Success      []string `json:"success"`
				UpdateErrors []string `json:"update_error"`
			}

			rsts := &batchResult{}
			for colIDx, colInput := range *innerBatchInfo.BatchInfo {

				delete(colInput, "import_from")

				if _, _, isUpdate, rstErr := createFunc(req, defErr, colInput, ownerID, objID, true, asstDes, attdes); nil != rstErr {
					if !isUpdate {
						blog.Debug("failed to create inst, error info is %s", rstErr.Error())
						rsts.Errors = append(rsts.Errors, fmt.Sprintf("Line:%d Error:%s", colIDx, rstErr.Error()))
					} else {
						blog.Debug("failed to update inst, error info is %s", rstErr.Error())
						rsts.UpdateErrors = append(rsts.UpdateErrors, fmt.Sprintf("Line:%d Error:%s", colIDx, rstErr.Error()))
					}

				} else {
					rsts.Success = append(rsts.Success, fmt.Sprintf("%d", colIDx))
				}
			}
			blog.Debug("the results:%+v", rsts)
			return http.StatusOK, rsts, nil
		}

		// create single inst
		status, rst, _, err := createFunc(req, defErr, input, ownerID, objID, false, asstDes, attdes)
		return status, rst, err

	}, resp)

}

// DeleteInst delete insts
func (cli *instAction) DeleteInst(req *restful.Request, resp *restful.Response) {

	blog.Info("delete inst")

	// get language
	language := util.GetActionLanguage(req)
	// get the error by language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	// logics
	cli.CallResponseEx(func() (int, interface{}, error) {

		ownerID := req.PathParameter("owner_id")
		objID := req.PathParameter("obj_id")
		user := util.GetActionUser(req)
		instID, convErr := strconv.Atoi(req.PathParameter("inst_id"))
		if nil != convErr {
			blog.Error("the instid[%s], must be int value, error info is %s", req.PathParameter("inst_id"), convErr.Error())
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsNeedInt, "inst_id")
		}

		// chceck the inst wether it has 'host' inst child
		type nextInst struct {
			instID  int
			ownerID string
			objID   string
		}

		nextOne := nextInst{instID: instID, ownerID: ownerID, objID: objID}
		willDelete := []nextInst{}
		willDelete = append(willDelete, nextInst{ownerID: ownerID, instID: instID, objID: objID})

		// define inner function to parse the object deleted
		var parseChildFunc func(child []manager.TopoInstRst) (int, interface{}, error)
		parseChildFunc = func(child []manager.TopoInstRst) (int, interface{}, error) {

			for _, instItem := range child {
				// store all child inst
				if 0 == strings.Compare(instItem.ObjID, common.BKInnerObjIDHost) {
					blog.Warn("can not delete the inst which has the 'host' child")
					return http.StatusBadRequest, "", defErr.Error(common.CCErrTopoInstHasHostChild)
				}
				willDelete = append(willDelete, nextInst{ownerID: ownerID, instID: instItem.InstID, objID: instItem.ObjID})
				//  if it is the last one, then will try next group
				if len(instItem.Child) != 0 {
					retStatus, retRst, retErr := parseChildFunc(instItem.Child)
					if nil != retErr {
						return retStatus, retRst, retErr
					}
				}
			}

			return http.StatusOK, nil, nil
		}

		// read all child inst, level -1 is no limit
		topoInstItems, topoErr := cli.metaHelperFunc.SelectInstTopo(nextOne.ownerID, nextOne.objID, 0, nextOne.instID, -1, req)
		if nil != topoErr {
			blog.Error("failed to get the inst topo , error info is %s", topoErr.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrTopoInstDeleteFailed)
		}
		for _, instItem := range topoInstItems {

			if 0 == strings.Compare(instItem.ObjID, common.BKInnerObjIDHost) {
				blog.Warn("can not delete the inst which has the 'host' child")
				return http.StatusBadRequest, "", defErr.Error(common.CCErrTopoInstHasHostChild)
			}

			willDelete = append(willDelete, nextInst{ownerID: ownerID, instID: instItem.InstID, objID: instItem.ObjID})
			//  if it is the last one, then will try next group
			if len(instItem.Child) != 0 {
				retStatus, retRst, retErr := parseChildFunc(instItem.Child)
				if nil != retErr {
					blog.Error("failed to delete inst, error info is %s", retErr.Error())
					return retStatus, retRst, retErr
				}
			}
		}

		attDesCache := map[string][]metadata.Header{}

		blog.Debug("will delete:%+v", willDelete)

		// already id
		ids := map[int]struct{}{}
		//delete module
		for _, delItem := range willDelete {

			if _, ok := ids[delItem.instID]; ok {
				continue
			}
			ids[delItem.instID] = struct{}{}
			input := make(map[string]interface{})
			input[common.BKOwnerIDField] = delItem.ownerID
			input[common.BKObjIDField] = delItem.objID
			input[common.BKInstIDField] = delItem.instID

			// delete the association
			if err := cli.deleteInstAssociation(delItem.instID, delItem.ownerID, delItem.objID); nil != err {
				blog.Errorf("failed to delete the association (%d %s %s), error info is %s", delItem.instID, delItem.ownerID, delItem.objID, err.Error())
			}

			// take snapshot before operation
			preData, retStrErr := cli.getInstDetail(req, delItem.instID, delItem.objID, delItem.ownerID)
			if common.CCSuccess != retStrErr {
				blog.Errorf("get inst detail error: %v", retStrErr)
				return http.StatusInternalServerError, "", defErr.Error(retStrErr)
			}

			uURL := cli.CC.ObjCtrl() + "/object/v1/insts/object"

			inputJSON, jsErr := json.Marshal(input)
			if nil != jsErr {
				blog.Error("unmarshal the data[%+v], error is %s", input, jsErr.Error())
				return http.StatusInternalServerError, "", defErr.Error(common.CCErrTopoInstDeleteFailed)
			}

			objRes, err := httpcli.ReqHttp(req, uURL, common.HTTPDelete, []byte(inputJSON))
			if nil != err {
				blog.Error("failed to delete the inst, error info is %s ", err.Error())
				return http.StatusInternalServerError, "", defErr.Error(common.CCErrTopoInstDeleteFailed)
			}
			if 0 == len(objRes) {
				return http.StatusOK, nil, nil
			}
			if _, ok := cli.IsSuccess([]byte(objRes)); !ok {
				blog.Error("failed to delete the inst , error inst is %s", objRes)
				return http.StatusInternalServerError, objRes, defErr.Error(common.CCErrTopoInstDeleteFailed)
			}

			{
				// save change log
				if attDesCache[delItem.objID] == nil {
					attDes, attErr := cli.getObjAttDes(delItem.ownerID, delItem.objID)
					if common.CCSuccess != attErr {
						return http.StatusInternalServerError, objRes, defErr.Error(attErr)
					}
					headers := []metadata.Header{}
					for _, item := range attDes {
						headers = append(headers, metadata.Header{
							PropertyID:   item.PropertyID,
							PropertyName: item.PropertyName,
						})
					}
					attDesCache[objID] = headers
				}

				auditContent := metadata.Content{
					PreData: preData,
					Headers: attDesCache[delItem.objID],
				}
				auditlog.NewClient(cli.CC.AuditCtrl()).AuditObjLog(delItem.instID, auditContent, "delete inst", delItem.objID, ownerID, "0", user, auditoplog.AuditOpTypeDel)
			}

		}

		return http.StatusOK, nil, nil
	}, resp)

}

// UpdateInst update inst
func (cli *instAction) UpdateInst(req *restful.Request, resp *restful.Response) {

	blog.Info("update inst")

	// get the language
	language := util.GetActionLanguage(req)

	// get the error by language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	// logics
	cli.CallResponseEx(func() (int, interface{}, error) {

		ownerID := req.PathParameter("owner_id")
		objID := req.PathParameter("obj_id")
		user := util.GetActionUser(req)
		instID, convErr := strconv.Atoi(req.PathParameter("inst_id"))
		if nil != convErr {
			blog.Error("the instid[%s], must be int value, error info is %s", req.PathParameter("inst_id"), convErr.Error())
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsNeedInt, "inst_id")
		}

		//update object
		input := make(map[string]interface{})

		condition := make(map[string]interface{})
		condition[common.BKOwnerIDField] = ownerID
		condition[common.BKObjIDField] = objID
		condition[common.BKInstIDField] = instID

		value, readErr := ioutil.ReadAll(req.Request.Body)
		if nil != readErr {
			blog.Error("failed to read the body , error info is %s", readErr.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		js, err := simplejson.NewJson([]byte(value))
		if nil != err {
			blog.Error("failed to create json object, error info is %s", err.Error())
			return http.StatusBadRequest, "", defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		data, jsErr := js.Map()
		if nil != jsErr {
			blog.Error("failed to create json object, error info is %s", jsErr.Error())
			return http.StatusBadRequest, "", defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		valid := validator.NewValidMap(ownerID, objID, cli.CC.ObjCtrl(), defErr)
		_, err = valid.ValidMap(data, common.ValidUpdate, instID)
		if nil != err {
			blog.Error("failed to valid the params, error info is %s ", err.Error())
			return http.StatusBadRequest, "", err
		}

		input["condition"] = condition
		input["data"] = data

		// take snapshot before update
		preData, retStrErr := cli.getInstDetail(req, instID, objID, ownerID)
		if common.CCSuccess != retStrErr {
			blog.Errorf("get inst detail error: %v", retStrErr)
			return http.StatusInternalServerError, "", defErr.Error(retStrErr)
		}

		// set the inst association table
		if err := cli.updateInstAssociation(instID, ownerID, objID, data); nil != err {
			blog.Errorf("failed to update the inst association, error info is %s ", err.Error())
		}

		// update the inst value
		uURL := cli.CC.ObjCtrl() + "/object/v1/insts/object"

		inputJSON, jsErr := json.Marshal(input)
		if nil != jsErr {
			blog.Error("failed to create json object, error info is %s", jsErr.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommJSONMarshalFailed)
		}

		objRes, err := httpcli.ReqHttp(req, uURL, common.HTTPUpdate, []byte(inputJSON))
		if nil != err {
			blog.Error("failed to update the inst, error info is %s", err.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrTopoInstUpdateFailed)
		}

		{
			// save change log
			headers, attErr := cli.getHeader(ownerID, objID)
			if common.CCSuccess != attErr {
				return http.StatusInternalServerError, objRes, defErr.Error(attErr)
			}
			curData, retStrErr := cli.getInstDetail(req, instID, objID, ownerID)
			if common.CCSuccess != retStrErr {
				blog.Errorf("get inst detail error: %v", retStrErr)
				return http.StatusInternalServerError, objRes, defErr.Error(retStrErr)
			}
			auditContent := metadata.Content{
				PreData: preData,
				CurData: curData,
				Headers: headers,
			}
			auditlog.NewClient(cli.CC.AuditCtrl()).AuditObjLog(instID, auditContent, "update inst", objID, ownerID, "0", user, auditoplog.AuditOpTypeModify)
		}

		return http.StatusOK, objRes, nil

	}, resp)

}

func (cli *instAction) getHeader(ownerID, objID string) ([]metadata.Header, int) {
	attDes, attErr := cli.getObjAttDes(ownerID, objID)
	if common.CCSuccess != attErr {
		return nil, attErr
	}
	headers := []metadata.Header{}
	for _, item := range attDes {
		headers = append(headers, metadata.Header{
			PropertyID:   item.PropertyID,
			PropertyName: item.PropertyName,
		})
	}
	return headers, common.CCSuccess
}

func (cli *instAction) getObjAttDes(ownerID, objID string) ([]api.ObjAttDes, int) {
	condition := map[string]interface{}{}
	condition[common.BKOwnerIDField] = ownerID
	condition[common.BKObjIDField] = objID
	searchData, _ := json.Marshal(condition)
	cli.objcli.SetAddress(cli.CC.ObjCtrl())
	attDes, restErr := cli.objcli.SearchMetaObjectAtt(searchData)
	if nil != restErr {
		blog.Error("failed to read the object att, error is %s ", restErr.Error())
		return nil, common.CCErrTopoInstSelectFailed
	}
	return attDes, 0
}

// getObjectAsst read association objectid the return key is engilish property name, value is the objectid
func (cli *instAction) getObjectAsst(objID, ownerID string) (map[string]string, int) {

	rstmap := map[string]string{}

	// get module
	cli.objcli.SetAddress(cli.CC.ObjCtrl())
	att := map[string]interface{}{}
	att[common.BKObjIDField] = objID
	att[common.BKOwnerIDField] = ownerID
	searchData, jsErr := json.Marshal(att)
	if nil != jsErr {
		blog.Error("failed to marshal the data[%+v], error info is %s", searchData, jsErr.Error())
		return nil, common.CCErrCommJSONMarshalFailed
	}
	rests, restErr := cli.objcli.SearchMetaObjectAtt(searchData)
	if nil != restErr {
		blog.Error("failed to read the object att, error is %s ", restErr.Error())
		return nil, common.CCErrTopoInstSelectFailed
	}

	// 组织模型名和对应的字段
	for _, item := range rests {

		switch item.PropertyType {

		case common.FiledTypeSingleAsst:

			asst := map[string]interface{}{}
			asst["bk_object_att_id"] = item.PropertyID
			asst[common.BKOwnerIDField] = item.OwnerID
			asst[common.BKObjIDField] = item.ObjectID

			searchData, jsErr := json.Marshal(asst)

			if nil != jsErr {
				blog.Error("failed to marshal the data[%+v], error info is %s", searchData, jsErr.Error())
			}

			asstRst, asstRstErr := cli.objcli.SearchMetaObjectAsst(searchData)
			if nil != asstRstErr {
				blog.Error("failed to read the object asst, error is %s ", asstRstErr.Error())
				return nil, common.CCErrTopoInstSelectFailed
			}

			if len(asstRst) > 0 { // only one association map
				rstmap[item.PropertyID] = asstRst[0].AsstObjID
			}

		case common.FieldTypeMultiAsst:

			asst := map[string]interface{}{}
			asst["bk_object_att_id"] = item.PropertyID
			asst[common.BKOwnerIDField] = item.OwnerID
			asst[common.BKObjIDField] = item.ObjectID
			searchData, jsErr := json.Marshal(asst)

			if nil != jsErr {
				blog.Error("failed to marshal the data[%+v], error info is %s", searchData, jsErr.Error())
			}

			asstRst, asstRstErr := cli.objcli.SearchMetaObjectAsst(searchData)
			if nil != asstRstErr {
				blog.Error("failed to read the object asst, error is %s ", asstRstErr.Error())
				return nil, common.CCErrTopoInstSelectFailed
			}

			if len(asstRst) > 0 { // only association with one object by one filed
				rstmap[item.PropertyID] = asstRst[0].AsstObjID
			}
		}

	}

	// rstmap: key is the bk_property_id  value is the association object id
	return rstmap, common.CCSuccess
}

func (cli *instAction) getInstAsst(req *restful.Request, ownerID, objID string, ids []string, page map[string]interface{}) ([]instNameAsst, int, int) {

	tmpIDS := []int{}
	for _, id := range ids {
		tmpID, _ := strconv.Atoi(id)
		tmpIDS = append(tmpIDS, tmpID)
	}

	condition := make(map[string]interface{})
	searchParams := make(map[string]interface{})

	searchParams["fields"] = ""
	if val, ok := page["fields"]; ok {
		searchParams["fields"] = val
	}
	searchParams["start"] = 0
	if val, ok := page["start"]; ok {
		searchParams["start"] = val
	}
	searchParams["limit"] = common.BKDefaultLimit
	if val, ok := page["limit"]; ok {
		searchParams["limit"] = val
	}

	var targetOBJ string
	var instName string
	var instID string

	targetPre := cli.CC.ObjCtrl() + "/object/v1/insts/"
	switch objID {
	case common.BKInnerObjIDHost:
		targetOBJ = ""
		instName = common.BKHostInnerIPField
		instID = common.BKHostIDField
		targetPre = cli.CC.HostCtrl() + "/host/v1/hosts"
		if 0 != len(tmpIDS) {
			condition[common.BKHostIDField] = map[string]interface{}{common.BKDBIN: tmpIDS}
		}
	case common.BKInnerObjIDApp:
		targetOBJ = common.BKInnerObjIDApp
		instName = common.BKAppNameField
		instID = common.BKAppIDField
		searchParams["sort"] = common.BKAppIDField
		condition[common.BKOwnerIDField] = ownerID
		if 0 != len(tmpIDS) {
			condition[common.BKAppIDField] = map[string]interface{}{common.BKDBIN: tmpIDS}
		}
	case common.BKInnerObjIDSet:
		targetOBJ = common.BKInnerObjIDSet
		instID = common.BKSetIDField
		instName = common.BKSetNameField
		searchParams["sort"] = common.BKSetIDField
		condition[common.BKSetIDField] = map[string]interface{}{common.BKDBIN: tmpIDS}
		condition[common.BKOwnerIDField] = ownerID
	case common.BKInnerObjIDModule:
		targetOBJ = common.BKInnerObjIDModule
		instID = common.BKModuleIDField
		instName = common.BKModuleNameField
		searchParams["sort"] = common.BKModuleIDField
		condition[common.BKOwnerIDField] = ownerID
		if 0 != len(tmpIDS) {
			condition[common.BKModuleIDField] = map[string]interface{}{common.BKDBIN: tmpIDS}
		}
	case common.BKInnerObjIDPlat:
		targetOBJ = common.BKInnerObjIDPlat
		instID = common.BKCloudIDField
		instName = common.BKCloudNameField
		searchParams["sort"] = common.BKCloudIDField
		if 0 != len(tmpIDS) {
			condition[common.BKCloudIDField] = map[string]interface{}{common.BKDBIN: tmpIDS}
		}
	default:
		targetOBJ = common.BKINnerObjIDObject
		instName = common.BKInstNameField
		instID = common.BKInstIDField
		condition[common.BKOwnerIDField] = ownerID
		condition[common.BKObjIDField] = objID
		if 0 != len(tmpIDS) {
			condition[common.BKInstIDField] = map[string]interface{}{common.BKDBIN: tmpIDS}
		}
		searchParams["sort"] = common.BKInstIDField

	}

	searchParams["condition"] = condition

	//search
	sURL := targetPre + targetOBJ + "/search"
	inputJSON, _ := json.Marshal(searchParams)
	objRes, err := httpcli.ReqHttp(req, sURL, common.HTTPSelectPost, []byte(inputJSON))
	blog.Debug("search inst api(%s) detail params: %s ", sURL, string(inputJSON))
	if nil != err {
		blog.Error("search inst defail failed, error: %v", err)
		return nil, 0, common.CCErrCommJSONMarshalFailed
	}

	js, err := simplejson.NewJson([]byte(objRes))
	if nil != err {
		blog.Error("the input json is invalid, error info is %s", err.Error())
		return nil, 0, common.CCErrCommJSONUnmarshalFailed
	}

	input, jsErr := js.Map()
	if nil != jsErr {
		blog.Error("the input json is invalid, error info is %s", jsErr.Error())
		return nil, 0, common.CCErrCommJSONUnmarshalFailed
	}

	delArryFunc := func(s []string, i int) []string {
		s[len(s)-1], s[i] = s[i], s[len(s)-1]
		return s[:len(s)-1]
	}

	rstName := []instNameAsst{}

	cnt := 0
	if data, ok := input["data"].(map[string]interface{}); ok {

		switch reflect.TypeOf(data["count"]).Kind() {
		case reflect.String:
			cnt, _ = strconv.Atoi(reflect.ValueOf(data["count"]).String())
		case reflect.Int:
			cnt = int(reflect.ValueOf(data["count"]).Int())
		default:
			blog.Warn("unknown the count type, value is %v type is %s", data["count"], reflect.TypeOf(data["count"]).Kind())
		}

		if info, infoOk := data["info"].([]interface{}); infoOk {
			for _, infoItem := range info {

				if dataItem, dataItemOk := infoItem.(map[string]interface{}); dataItemOk {

					if dataItemVal, dataItemValOk := dataItem[instName]; dataItemValOk {

						// 提取实例名
						inst := instNameAsst{}
						if dataItemValStr, convOk := dataItemVal.(string); convOk {
							inst.InstName = dataItemValStr
							inst.ObjID = objID
						}

						// 删除已经存在的ID
						if dataItemVal, dataItemValOk := dataItem[instID]; dataItemValOk {

							//fmt.Println(key, dataItemVal)
							switch d := dataItemVal.(type) {
							case json.Number:
								if 0 != len(ids) {
									for idx, key := range ids {
										if val, err := d.Int64(); nil == err && key == strconv.Itoa(int(val)) {
											//fmt.Println("key:", key, "val:", val)
											inst.ID = ids[idx]
											inst.InstID, _ = strconv.Atoi(ids[idx])
											ids = delArryFunc(ids, idx)
											rstName = append(rstName, inst)
											goto next
										}
									}
								} else if val, err := d.Int64(); nil == err {
									inst.ID = strconv.Itoa(int(val))
									inst.InstID = int(val)
									rstName = append(rstName, inst)
								}

							}
						next:
						}

					}
				}
			}
		}
	}

	// deal the other inst name
	for _, id := range ids {
		rstName = append(rstName, instNameAsst{ID: id})
	}

	return rstName, cnt, common.CCSuccess
}

func (cli *instAction) getCommonInstTopo(req *restful.Request, objID, ownerID, instRes string, page map[string]interface{}) ([]commonInstTopo, int) {

	// set address
	cli.objcli.SetAddress(cli.CC.ObjCtrl())

	// read the association filed about objID
	rstmap, errorno := cli.getObjectAsst(objID, ownerID)
	if common.CCSuccess != errorno {
		blog.Error("failed to search the association with ownerid(%s) objectid(%s)", ownerID, objID)
		return nil, errorno
	}

	js, err := simplejson.NewJson([]byte(instRes))
	if nil != err {
		blog.Error("the input json is invalid, error info is %s", err.Error())
		return nil, common.CCErrCommJSONUnmarshalFailed
	}

	input, jsErr := js.Map()
	if nil != jsErr {
		blog.Error("the input json is invalid, error info is %s", jsErr.Error())
		return nil, common.CCErrCommJSONUnmarshalFailed
	}

	blog.Debug("input: %+v", input)
	blog.Debug("rstmap:%+v", rstmap)
	// inst result
	rstInst := make([]commonInstTopo, 0)

	// parse the data
	if data, ok := input["data"].(map[string]interface{}); ok {
		if info, infoOk := data["info"].([]interface{}); infoOk {

			for _, infoItem := range info {

				if dataItem, dataItemOk := infoItem.(map[string]interface{}); dataItemOk {

					// key 是关联字段，val 是字段关联的模型ID
					for key, val := range rstmap {

						// search association objid
						objCondition := map[string]interface{}{}
						objCondition[common.BKOwnerIDField] = ownerID
						objCondition[common.BKObjIDField] = val
						objConditionStr, _ := json.Marshal(objCondition)

						// get objid information
						objItems, objErr := cli.objcli.SearchMetaObject(objConditionStr)
						if nil != objErr {
							blog.Error("failed to search objects, error info is %s", objErr.Error())
							return nil, common.CCErrCommHTTPDoRequestFailed
						}

						if 0 == len(objItems) {
							blog.Error("failed to search the objsect by the condition ownerid(%s) objid(%s)", ownerID, val)
							return nil, common.CCErrTopoObjectSelectFailed
						}

						// set common object name
						commonInst := commonInstTopo{}
						commonInst.InstName = objItems[0].ObjectName
						commonInst.ObjID = val
						commonInst.ObjIcon = objItems[0].ObjIcon
						commonInst.ID = strconv.Itoa(objItems[0].ID)

						if keyItem, keyItemOk := dataItem[key]; keyItemOk {

							keyItemStr := fmt.Sprintf("%v", keyItem)

							blog.Debug("keyitemstr:%s", keyItemStr)

							// search association insts
							retData, cnt, retErr := cli.getInstAsst(req, ownerID, val, strings.Split(keyItemStr, ","), map[string]interface{}{
								"start": 0,
								"limit": common.BKNoLimit,
								"sort":  "",
							})
							if common.CCSuccess != retErr {
								blog.Error("failed to get inst details")
								continue
							}
							commonInst.Count = cnt
							commonInst.Children = append(commonInst.Children, retData...)
							//dataItem[key] = retData
						}

						// append the result
						sort.Sort(instAsstSort(commonInst.Children))
						rstInst = append(rstInst, commonInst)
					}

				}
			}
		}
	}

	sort.Sort(instTopoSort(rstInst))
	return rstInst, common.CCSuccess
}

func (cli *instAction) getInstDetails(req *restful.Request, objID, ownerID, instStr string, page map[string]interface{}) (map[string]interface{}, int) {

	blog.Debug("ownerid(%s) objid(%s)", ownerID, objID)
	// read object association map
	rstmap, errorno := cli.getObjectAsst(objID, ownerID)
	if common.CCSuccess != errorno {
		return nil, errorno
	}

	js, err := simplejson.NewJson([]byte(instStr))
	if nil != err {
		blog.Error("the input json is invalid, error info is %s", err.Error())
		return nil, common.CCErrCommJSONUnmarshalFailed
	}

	input, jsErr := js.Map()
	if nil != jsErr {
		blog.Error("the input json is invalid, error info is %s", jsErr.Error())
		return nil, common.CCErrCommJSONUnmarshalFailed
	}

	if data, ok := input["data"].(map[string]interface{}); ok {
		if info, infoOk := data["info"].([]interface{}); infoOk {
			for _, infoItem := range info {

				if dataItem, dataItemOk := infoItem.(map[string]interface{}); dataItemOk {

					// key is the association object filed，val is association object id
					for key, val := range rstmap {

						if keyItem, keyItemOk := dataItem[key]; keyItemOk {

							keyItemStr := fmt.Sprintf("%v", keyItem)
							blog.Debug("keyitemstr:%s", keyItemStr)
							retData, _, retErr := cli.getInstAsst(req, ownerID, val, strings.Split(keyItemStr, ","), page)
							if common.CCSuccess != retErr {
								blog.Error("failed to get inst details")
							}
							dataItem[key] = retData

						}

					}
				}
			}
		}
	}

	return input, common.CCSuccess
}

func (cli *instAction) getInstDeteilByCondition(req *restful.Request, objID string, ownerID string, condition map[string]interface{}) (map[string]interface{}, int) {
	objType := ""
	targetpre := cli.CC.ObjCtrl() + "/object/v1/insts/"
	switch objID {
	case common.BKInnerObjIDHost:
		objType = ""
		targetpre = cli.CC.HostCtrl() + "/host/v1/hosts"
	case common.BKInnerObjIDModule:
		objType = common.BKInnerObjIDModule
		condition[common.BKOwnerIDField] = ownerID
	case common.BKInnerObjIDApp:
		objType = common.BKInnerObjIDApp
		condition[common.BKOwnerIDField] = ownerID
	case common.BKInnerObjIDSet:
		objType = common.BKInnerObjIDSet
		condition[common.BKOwnerIDField] = ownerID
	default:
		objType = common.BKINnerObjIDObject
		condition[common.BKOwnerIDField] = ownerID
	}
	searchParam, jsErr := json.Marshal(map[string]interface{}{
		"condition": condition,
	})
	if jsErr != nil {
		return nil, common.CCErrCommJSONMarshalFailed
	}
	sURL := targetpre + objType + "/search"
	objRes, err := httpcli.ReqHttp(req, sURL, "POST", searchParam)
	if nil != err {
		blog.Error("failed to select the insts, error info is %s", err.Error())
		return nil, common.CCErrTopoInstSelectFailed
	}
	ret, instErr := cli.getInstDetails(req, objID, ownerID, objRes, nil)
	if instErr != common.CCSuccess {
		blog.Error("failed to select the insts, objID(%s) ownerid(%s)", objID, ownerID)
		return nil, instErr
	}

	blog.Debug("the return data:%+v", ret)
	if data, ok := ret["data"].(map[string]interface{}); ok {
		if info, infoOk := data["info"].([]interface{}); infoOk {
			if len(info) > 0 && info[0] != nil {
				if ret, ok := info[0].(map[string]interface{}); ok {
					return ret, common.CCSuccess
				}
			}
		}
	}
	return nil, common.CCErrTopoInstSelectFailed
}

func (cli *instAction) getInstDetail(req *restful.Request, instID int, objID, ownerID string) (interface{}, int) {
	blog.Debug("instid(%d), ownerid(%s) objid(%s)", instID, ownerID, objID)
	condition := map[string]interface{}{}
	switch objID {
	case common.BKInnerObjIDHost:
		condition[common.BKHostIDField] = instID
	case common.BKInnerObjIDModule:
		condition[common.BKModuleIDField] = instID
		condition[common.BKOwnerIDField] = ownerID
	case common.BKInnerObjIDApp:
		condition[common.BKAppIDField] = instID
		condition[common.BKOwnerIDField] = ownerID
	case common.BKInnerObjIDSet:
		condition[common.BKSetIDField] = instID
		condition[common.BKOwnerIDField] = ownerID
	default:
		condition[common.BKObjIDField] = objID
		condition[common.BKInstIDField] = instID
		condition[common.BKOwnerIDField] = ownerID
	}
	return cli.getInstDeteilByCondition(req, objID, ownerID, condition)
}

// SelectTopo search common inst topo
func (cli *instAction) SelectTopo(req *restful.Request, resp *restful.Response) {

	blog.Info("select inst topo")

	// read language
	language := util.GetActionLanguage(req)

	// generate error object by language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	cli.CallResponseEx(func() (int, interface{}, error) {

		//{ownerid}/object/{objectid}/inst/{instid}
		ownerID := req.PathParameter("owner_id")
		objID := req.PathParameter("object_id")

		instID, convErr := strconv.Atoi(req.PathParameter("inst_id"))
		if nil != convErr {
			blog.Error("failed to convert, the error info is %s", convErr.Error())
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsNeedInt, "inst_id")
		}

		value, readErr := ioutil.ReadAll(req.Request.Body)
		if nil != readErr {
			blog.Error("failed to read the body, error is %s", readErr.Error())
			return http.StatusBadRequest, "", defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		var js params.SearchParams
		if 0 != len(value) {
			err := json.Unmarshal([]byte(value), &js)
			if nil != err {
				blog.Error("failed to unmarshal the data[%s], error is %s", value, err.Error())
				return http.StatusBadRequest, "", defErr.Error(common.CCErrCommJSONUnmarshalFailed)
			}
		}

		condition := map[string]interface{}{}
		objType := ""
		targetpre := cli.CC.ObjCtrl() + "/object/v1/insts/"
		switch objID {
		case common.BKInnerObjIDHost:
			objType = ""
			targetpre = cli.CC.HostCtrl() + "/host/v1/hosts"
			condition[common.BKHostIDField] = instID
		case common.BKInnerObjIDModule:
			objType = common.BKInnerObjIDModule
			condition[common.BKModuleIDField] = instID
			condition[common.BKOwnerIDField] = ownerID
		case common.BKInnerObjIDApp:
			objType = common.BKInnerObjIDApp
			condition[common.BKAppIDField] = instID
			condition[common.BKOwnerIDField] = ownerID
		case common.BKInnerObjIDSet:
			objType = common.BKInnerObjIDSet
			condition[common.BKSetIDField] = instID
			condition[common.BKOwnerIDField] = ownerID
		default:
			objType = common.BKINnerObjIDObject
			condition[common.BKObjIDField] = objID
			condition[common.BKInstIDField] = instID
			condition[common.BKOwnerIDField] = ownerID
		}

		// construct the search params

		searchParams := make(map[string]interface{})
		searchParams["condition"] = condition
		searchParams["fields"] = ""
		searchParams["start"] = 0
		searchParams["limit"] = common.BKNoLimit
		searchParams["sort"] = ""

		//search insts
		sURL := targetpre + objType + "/search"

		inputJSON, jsErr := json.Marshal(searchParams)
		if nil != jsErr {
			blog.Error("failed to marshal the data[%+v], error info is %s", searchParams, jsErr.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommJSONMarshalFailed)
		}

		instRes, err := httpcli.ReqHttp(req, sURL, "POST", []byte(inputJSON))
		if nil != err {
			blog.Error("failed to select the insts, error info is %s", err.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrTopoInstSelectFailed)
		}

		// get common topo inst
		retstr, retStrErr := cli.getCommonInstTopo(req, objID, ownerID, instRes, js.Page)
		blog.Debug("retstr: %+v", retstr)
		if common.CCSuccess != retStrErr {
			return http.StatusInternalServerError, "", defErr.Error(retStrErr)
		}

		return http.StatusOK, retstr, nil

	}, resp)

}

// SelectInst 查询实例
func (cli *instAction) SelectInst(req *restful.Request, resp *restful.Response) {
	blog.Info("select inst")

	// get language
	language := util.GetActionLanguage(req)

	// get the error by language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	// logics
	cli.CallResponseEx(func() (int, interface{}, error) {

		ownerID := req.PathParameter("owner_id")
		objID := req.PathParameter("obj_id")
		instID, convErr := strconv.Atoi(req.PathParameter("inst_id"))
		if nil != convErr {
			blog.Error("failed to convert, the error info is %s", convErr.Error())
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsNeedInt, "instid")
		}

		value, readErr := ioutil.ReadAll(req.Request.Body)
		if nil != readErr {
			blog.Error("failed to read the body, error is %s", readErr.Error())
			return http.StatusBadRequest, "", defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		var js params.SearchParams
		err := json.Unmarshal([]byte(value), &js)
		if nil != err {
			blog.Error("failed to unmarshal the data[%s], error is %s", value, err.Error())
			return http.StatusBadRequest, "", defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		condition := params.ParseAppSearchParams(js.Condition)

		condition[common.BKOwnerIDField] = ownerID
		condition[common.BKObjIDField] = objID
		condition[common.BKInstIDField] = instID

		page := js.Page

		searchParams := make(map[string]interface{})

		searchParams["condition"] = condition
		searchParams["fields"] = strings.Join(js.Fields, ",")
		searchParams["start"] = page["start"]
		searchParams["limit"] = page["limit"]
		searchParams["sort"] = page["sort"]

		//search
		sURL := cli.CC.ObjCtrl() + "/object/v1/insts/object/search"

		inputJSON, jsErr := json.Marshal(searchParams)
		if nil != jsErr {
			blog.Error("failed to marshal the data[%+v], error info is %s", searchParams, jsErr.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommJSONMarshalFailed)
		}

		objRes, err := httpcli.ReqHttp(req, sURL, "POST", []byte(inputJSON))

		if nil != err {
			blog.Error("failed to select the insts, error info is %s", err.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrTopoInstSelectFailed)
		}

		retstr, retStrErr := cli.getInstDetails(req, objID, ownerID, objRes, js.Page)
		if common.CCSuccess != retStrErr {
			return http.StatusInternalServerError, "", defErr.Error(retStrErr)
		}

		return http.StatusOK, retstr["data"], nil

	}, resp)

}

// SelectInsts search insts by condition
func (cli *instAction) SelectInsts(req *restful.Request, resp *restful.Response) {
	blog.Info("select insts")

	// get language
	language := util.GetActionLanguage(req)

	// get error info by language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	// logics
	cli.CallResponseEx(func() (int, interface{}, error) {
		ownerID := req.PathParameter("owner_id")
		objID := req.PathParameter("obj_id")

		searchParams := make(map[string]interface{})

		value, readErr := ioutil.ReadAll(req.Request.Body)

		if nil != readErr {
			blog.Error("failed to read the body , error info is %s", readErr.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		if 0 != len(value) {

			var js params.SearchParams
			err := json.Unmarshal([]byte(value), &js)
			if nil != err {
				blog.Error("failed to unmarshal the data[%s], error is %s", value, err.Error())
				return http.StatusBadRequest, "", defErr.Error(common.CCErrCommJSONUnmarshalFailed)
			}

			condition := params.ParseAppSearchParams(js.Condition)

			condition[common.BKOwnerIDField] = ownerID
			condition[common.BKObjIDField] = objID

			page := js.Page

			searchParams["condition"] = condition
			searchParams["fields"] = strings.Join(js.Fields, ",")
			searchParams["start"] = page["start"]
			searchParams["limit"] = page["limit"]
			searchParams["sort"] = page["sort"]

		} else {
			condition := make(map[string]interface{}, 0)
			condition[common.BKOwnerIDField] = ownerID
			condition[common.BKObjIDField] = objID
			searchParams["condition"] = condition
			searchParams["fields"] = ""
			searchParams["start"] = 0
			searchParams["limit"] = common.BKDefaultLimit
			searchParams["sort"] = ""

		}

		//search
		sURL := cli.CC.ObjCtrl() + "/object/v1/insts/object/search"
		inputJSON, jsErr := json.Marshal(searchParams)

		if nil != jsErr {
			blog.Error("failed to marshal the data[%+v], error info is %s", searchParams, jsErr.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommJSONMarshalFailed)
		}

		objRes, err := httpcli.ReqHttp(req, sURL, common.HTTPSelectPost, inputJSON)
		blog.Debug("search url(%s) inst params: %s", sURL, string(inputJSON))

		if nil != err {
			blog.Error("failed to select the insts, error info is %s", err.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrTopoInstSelectFailed)
		}

		retStr, retStrErr := cli.getInstDetails(req, objID, ownerID, objRes, map[string]interface{}{
			"start": 0,
			"limit": common.BKNoLimit,
			"sort":  "",
		})

		if common.CCSuccess != retStrErr {
			return http.StatusInternalServerError, "", defErr.Error(retStrErr)
		}

		return http.StatusOK, retStr["data"], nil
	}, resp)

}
