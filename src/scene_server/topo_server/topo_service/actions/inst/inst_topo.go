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
	"configcenter/src/common/blog"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	scenecommon "configcenter/src/scene_server/common"
	api "configcenter/src/source_controller/api/object"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
	restful "github.com/emicklei/go-restful"
)

// getAsstParentObject get the association parent object, the return key is engilish property name, value is the objectid
func (cli *instAction) getAsstParentObject(forward *api.ForwardParam, objID, ownerID string) (map[string]string, int) {

	//  the key is the field id , the value is the object id
	rstMap := map[string]string{}

	// get module
	cli.objcli.SetAddress(cli.CC.ObjCtrl())
	attAsstCond := map[string]interface{}{}
	attAsstCond[common.BKAsstObjIDField] = objID
	attAsstCond[common.BKOwnerIDField] = ownerID

	searchAttAsstCond, err := json.Marshal(attAsstCond)
	if nil != err {
		return rstMap, common.CCErrCommJSONMarshalFailed
	}

	asstRst, asstRstErr := cli.objcli.SearchMetaObjectAsst(forward, searchAttAsstCond)
	if nil != asstRstErr {
		blog.Error("failed to read the object asst, error is %s ", asstRstErr.Error())
		return nil, common.CCErrTopoInstSelectFailed
	}

	for _, asstItem := range asstRst {
		rstMap[asstItem.ObjectAttID] = asstItem.ObjectID
	}

	// rstmap: key is the bk_property_id  value is the association object id
	return rstMap, common.CCSuccess
}

// getObjectAsst read association objectid the return key is engilish property name, value is the objectid
func (cli *instAction) getObjectAsst(forward *api.ForwardParam, objID, ownerID string) (map[string]string, int) {

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
	rests, restErr := cli.objcli.SearchMetaObjectAtt(forward, searchData)
	if nil != restErr {
		blog.Error("failed to read the object att, error is %s ", restErr.Error())
		return nil, common.CCErrTopoInstSelectFailed
	}

	// 组织模型名和对应的字段
	for _, item := range rests {

		switch item.PropertyType {

		case common.FieldTypeSingleAsst:

			asst := map[string]interface{}{}
			asst[common.BKObjAttIDField] = item.PropertyID
			asst[common.BKOwnerIDField] = item.OwnerID
			asst[common.BKObjIDField] = item.ObjectID

			searchData, jsErr := json.Marshal(asst)

			if nil != jsErr {
				blog.Error("failed to marshal the data[%+v], error info is %s", searchData, jsErr.Error())
			}

			asstRst, asstRstErr := cli.objcli.SearchMetaObjectAsst(forward, searchData)
			if nil != asstRstErr {
				blog.Error("failed to read the object asst, error is %s ", asstRstErr.Error())
				return nil, common.CCErrTopoInstSelectFailed
			}

			if len(asstRst) > 0 { // only one association map
				rstmap[item.PropertyID] = asstRst[0].AsstObjID
			}

		case common.FieldTypeMultiAsst:

			asst := map[string]interface{}{}
			asst[common.BKObjAttIDField] = item.PropertyID
			asst[common.BKOwnerIDField] = item.OwnerID
			asst[common.BKObjIDField] = item.ObjectID
			searchData, jsErr := json.Marshal(asst)

			if nil != jsErr {
				blog.Error("failed to marshal the data[%+v], error info is %s", searchData, jsErr.Error())
			}

			asstRst, asstRstErr := cli.objcli.SearchMetaObjectAsst(forward, searchData)
			if nil != asstRstErr {
				blog.Error("failed to read the object asst, error is %s ", asstRstErr.Error())
				return nil, common.CCErrTopoInstSelectFailed
			}

			if len(asstRst) > 0 { // only association with one object by one field
				rstmap[item.PropertyID] = asstRst[0].AsstObjID
			}
		}

	}

	// rstmap: key is the bk_property_id  value is the association object id
	return rstmap, common.CCSuccess
}
func (cli *instAction) getCommonChildInstTopo(req *restful.Request, objID, ownerID, instRes string, page map[string]interface{}) ([]commonInstTopo, int) {

	// set address
	cli.objcli.SetAddress(cli.CC.ObjCtrl())
	forward := &api.ForwardParam{Header: req.Request.Header}

	// read the association filed about objID
	rstmap, errorno := cli.getObjectAsst(forward, objID, ownerID)
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
						objItems, objErr := cli.objcli.SearchMetaObject(forward, objConditionStr)
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

func (cli *instAction) getCommonParentInstTopo(req *restful.Request, objID, ownerID, instRes string, page map[string]interface{}) ([]commonInstTopo, int) {

	// set address
	cli.objcli.SetAddress(cli.CC.ObjCtrl())
	forward := &api.ForwardParam{Header: req.Request.Header}

	// read the parent association field about objID
	rstMap, errNo := cli.getAsstParentObject(forward, objID, ownerID)
	if common.CCSuccess != errNo {
		blog.Error("failed to search the parent association with ownerid(%s) objectid(%s)", ownerID, objID)
		return nil, errNo
	}

	blog.Debug("the parent association object:%+v", rstMap)
	// inst result
	rstInst := make([]commonInstTopo, 0)

	for _, prevObjID := range rstMap {

		// search prev association objid
		objCondition := map[string]interface{}{}
		objCondition[common.BKOwnerIDField] = ownerID
		objCondition[common.BKObjIDField] = prevObjID

		objConditionStr, _ := json.Marshal(objCondition)

		// get objid information
		objItems, objErr := cli.objcli.SearchMetaObject(forward, objConditionStr)
		if nil != objErr {
			blog.Error("failed to search objects, error info is %s", objErr.Error())
			return nil, common.CCErrTopoObjectSelectFailed
		}

		if 0 == len(objItems) {
			blog.Error("failed to search the objsect by the condition ownerid(%s) objid(%s)", ownerID, prevObjID)
			return nil, common.CCSuccess
		}

		// set common object name
		commonInst := commonInstTopo{}
		commonInst.ObjectName = objItems[0].ObjectName
		commonInst.ObjID = objItems[0].ObjectID
		commonInst.ObjIcon = objItems[0].ObjIcon
		commonInst.ID = strconv.Itoa(objItems[0].ID)

		rstInst = append(rstInst, commonInst)

		// search the insts

		js := gjson.Parse(instRes)
		rstItems := js.Get("data.info").Array()

		for _, valItem := range rstItems {
			blog.Infof("the value:%v", valItem)

			// construct the object id
			currInstID := valItem.Get(common.BKInstIDField).String()
			currObjectID := valItem.Get(common.BKObjIDField).String()

			// search parent association inst id
			objCondition := map[string]interface{}{}
			objCondition[common.BKAsstObjIDField] = currObjectID
			objCondition[common.BKAsstInstIDField] = currInstID
			objCondition[common.BKObjIDField] = prevObjID

			asstInstRes, asstInstResErr := scenecommon.SearchInstAssociation(cli.CC.ObjCtrl(), objCondition, req)
			if nil != asstInstResErr {
				blog.Errorf("failed to request, error info is %s", asstInstResErr.Error())
				return rstInst, common.CCErrTopoInstSelectFailed
			}

			// extract the inst id for the prev object
			targetInstIDS := make([]string, 0)
			gjson.Get(asstInstRes, "data.info.#."+common.BKInstIDField).ForEach(func(key, value gjson.Result) bool {

				targetInstIDS = append(targetInstIDS, fmt.Sprintf("%d", value.Int()))
				return true
			})

			// search the prev object insts
			retData, cnt, retErr := cli.getInstAsst(req, ownerID, prevObjID, targetInstIDS, map[string]interface{}{
				"start": 0,
				"limit": common.BKNoLimit,
				"sort":  "",
			})
			if common.CCSuccess != retErr {
				blog.Error("failed to get inst details")
				return nil, retErr
			}
			commonInst.Count = cnt
			commonInst.Children = append(commonInst.Children, retData...)

		}

	}

	sort.Sort(instTopoSort(rstInst))
	return rstInst, common.CCSuccess
}

func (cli *instAction) SelectAssociationTopo(req *restful.Request, resp *restful.Response) {
	blog.Info("select inst association topo (parent|child) topo")

	// read language
	language := util.GetActionLanguage(req)

	// generate error object by language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	cli.CallResponseEx(func() (int, interface{}, error) {

		//{owner_id}/object/{object_id}/inst/{inst_id}
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

		// get common topo child inst
		rstTopoChild, rstErr := cli.getCommonChildInstTopo(req, objID, ownerID, instRes, js.Page)
		blog.Debug("result topo child : %+v", rstTopoChild)
		if common.CCSuccess != rstErr {
			return http.StatusInternalServerError, "", defErr.Error(rstErr)
		}

		// get common topo parent inst
		rstTopoParent, rstErr := cli.getCommonParentInstTopo(req, objID, ownerID, instRes, js.Page)
		blog.Debug("result topo parent : %+v", rstTopoParent)
		if common.CCSuccess != rstErr {
			return http.StatusInternalServerError, "", defErr.Error(rstErr)
		}

		return http.StatusOK, rstTopoChild, nil

	}, resp)
}
