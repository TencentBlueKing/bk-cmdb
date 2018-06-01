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

package actions

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/util"
	sourceAPI "configcenter/src/source_controller/api/object"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"
)

func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/{" + common.BKOwnerIDField + "}/{" + common.BKAppIDField + "}/{" + common.BKProcIDField + "}", Params: nil, Handler: process.GetProcessDetailByID})
	process.CreateAction()
}

//GetProcessDetailByID  get process detail info
func (cli *procAction) GetProcessDetailByID(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		pathParams := req.PathParameters()
		ownerID := pathParams[common.BKOwnerIDField]
		appIDStr := pathParams[common.BKAppIDField]
		appID, _ := strconv.Atoi(appIDStr)
		procIDStr := pathParams[common.BKProcIDField]
		procID, _ := strconv.Atoi(procIDStr)

		reResult, err := cli.getProcDetail(req, ownerID, appID, procID)
		if err != nil {
			blog.Error("GetProcessDetailByID info error :%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrProcSearchDetailFaile)
		}
		return http.StatusOK, reResult, nil
	}, resp)
}

func (cli *procAction) getProcDetail(req *restful.Request, ownerID string, appID, procID int) ([]map[string]interface{}, error) {
	condition := make(map[string]interface{})
	condition[common.BKOwnerIDField] = ownerID
	condition[common.BKAppIDField] = appID
	condition[common.BKProcIDField] = procID
	searchParams := make(map[string]interface{})
	searchParams["condition"] = condition
	procInfoJson, _ := json.Marshal(searchParams)
	gProcURL := cli.CC.ObjCtrl() + "/object/v1/insts/process/search"
	blog.Info("get obj detail query url: %v", gProcURL)
	blog.Info("get obj detail query params: %v", string(procInfoJson))
	gProcRe, err := httpcli.ReqHttp(req, gProcURL, common.HTTPSelectPost, []byte(procInfoJson))
	blog.Info("get obj detail query return: %v", gProcRe)
	if nil != err {
		blog.Error("GetProcessDetailByID info error :%v", err)
		return nil, err
	}
	var getResult ProcessResult
	err = json.Unmarshal([]byte(gProcRe), &getResult)
	if nil != err {
		blog.Error("GetProcessDetailByID info error :%v", err)
		return nil, err
	}
	if false == getResult.Result {
		blog.Error("GetProcDetailByID  info error :%v", err)
		return nil, err
	}
	procData, ok := getResult.Data.(map[string]interface{})
	if false == ok {
		blog.Error("GetProcDetailByID  info error :%v", err)
		return nil, err
	}
	procInfoArr := procData["info"].([]interface{})
	proc := make(map[string]interface{}, 0)
	for _, j := range procInfoArr {
		procCell := j.(map[string]interface{})
		for k, v := range procCell {
			proc[k] = v
		}
	}
	gProcAttrURL := cli.CC.ObjCtrl() + "/object/v1/meta/objectatts"
	searchBody := make(map[string]interface{})
	searchBody[common.BKObjIDField] = common.BKInnerObjIDProc
	searchBody[common.BKOwnerIDField] = ownerID
	searchJson, _ := json.Marshal(searchBody)
	gProcAttrRe, err := httpcli.ReqHttp(req, gProcAttrURL, common.HTTPSelectPost, []byte(searchJson))
	if nil != err {
		blog.Error("GetProcDetailByID  attr error :%v", err)
		return nil, err
	}

	js, err := simplejson.NewJson([]byte(gProcAttrRe))
	gProcAttr, _ := js.Map()
	gAttrResult := gProcAttr["result"].(bool)
	if false == gAttrResult {
		blog.Error("GetProcDetailByID  attr error :%v", err)
		return nil, err
	}
	procAttrArr := gProcAttr["data"].([]interface{})
	forward := &sourceAPI.ForwardParam{Header: req.Request.Header}
	rstmap, errorno := cli.getObjectAsst(forward, common.BKInnerObjIDProc, ownerID)
	if common.CCSuccess != errorno {
		return nil, fmt.Errorf("get object asst faile")
	}

	fmt.Printf("rst map : %v\n", rstmap)

	reResult := make([]map[string]interface{}, 0)
	for _, i := range procAttrArr {
		attr := i.(map[string]interface{})
		data := make(map[string]interface{})
		propertyID := attr[common.BKPropertyIDField].(string)
		if propertyID == common.BKChildStr {
			continue
		}
		data[common.BKPropertyIDField] = propertyID
		data[common.BKPropertyNameField] = attr[common.BKPropertyNameField]

		data[common.BKPropertyValueField] = proc[propertyID]
		// key is the association object filed，val is association object id
		if val, ok := rstmap[propertyID]; ok {
			keyItemStr := fmt.Sprintf("%v", proc[propertyID])
			blog.Debug("keyitemstr:%s", keyItemStr)
			retData, _, retErr := cli.getInstAsst(req, ownerID, val, strings.Split(keyItemStr, ","), nil)
			if common.CCSuccess != retErr {
				blog.Error("failed to get inst details")
			}
			data[common.BKPropertyValueField] = retData
		}
		reResult = append(reResult, data)
	}
	return reResult, nil
}

// instNameAsst  association inst name
type instNameAsst struct {
	ID         string `json:"id"`
	ObjID      string `json:"bk_obj_id"`
	ObjIcon    string `json:"bk_obj_icon"`
	InstID     int    `json:"bk_inst_id"`
	ObjectName string `json:"bk_obj_name"`
	InstName   string `json:"bk_inst_name"`
}

// getObjectAsst read association objectid the return key is engilish property name, value is the objectid
func (cli *procAction) getObjectAsst(forward *sourceAPI.ForwardParam, objID, ownerID string) (map[string]string, int) {

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
			asst["bk_object_att_id"] = item.PropertyID
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
			asst["bk_object_att_id"] = item.PropertyID
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

			if len(asstRst) > 0 { // only association with one object by one filed
				rstmap[item.PropertyID] = asstRst[0].AsstObjID
			}
		}

	}

	// rstmap: key is the bk_property_id  value is the association object id
	return rstmap, common.CCSuccess
}

func (cli *procAction) getInstAsst(req *restful.Request, ownerID, objID string, ids []string, page map[string]interface{}) ([]instNameAsst, int, int) {

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
