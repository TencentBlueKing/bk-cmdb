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

package instapi

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpcli "configcenter/src/common/http/httpclient"
	api "configcenter/src/source_controller/api/object"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
	restful "github.com/emicklei/go-restful"
)

// instNameAsst  association inst name
type instNameAsst struct {
	ID         string `json:"id"`
	ObjID      string `json:"bk_obj_id"`
	ObjIcon    string `json:"bk_obj_icon"`
	ObjectID   int    `json:"bk_inst_id"`
	ObjectName string `json:"bk_obj_name"`
	Name       string `json:"bk_inst_name"`
}

// Inst public inst
var Inst = &instHelper{}

// instHelper
type instHelper struct {
	hostCtrl, objCtrl string
	objcli            *api.Client
}

func (cli *instHelper) InitInstHelper(hostCtrl, objCtrl string) error {
	cli.hostCtrl = hostCtrl
	cli.objCtrl = objCtrl
	cli.objcli = api.NewClient("")
	return nil
}

// getObjectAsst read association objectid the return key is engilish property name, value is the objectid
func (cli *instHelper) GetObjectAsst(objID, ownerID string) (map[string]string, int) {

	rstmap := map[string]string{}

	//read module

	att := map[string]interface{}{}
	att[common.BKObjIDField] = objID
	att[common.BKOwnerIDField] = ownerID
	searchData, jsErr := json.Marshal(att)
	if nil != jsErr {
		blog.Error("failed to marshal the data[%+v], error info is %s", searchData, jsErr.Error())
		return nil, common.CCErrCommJSONMarshalFailed
	}
	cli.objcli.SetAddress(cli.objCtrl)
	rests, restErr := cli.objcli.SearchMetaObjectAtt(searchData)
	if nil != restErr {
		blog.Error("failed to read the object att, error is %s ", restErr.Error())
		return nil, common.CCErrTopoInstSelectFailed
	}

	// com object name and field
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

			if len(asstRst) > 0 { // one filed only asst to one module
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
				blog.Error("failed to read the object asst, error is %s ", restErr.Error())
				return nil, common.CCErrTopoInstSelectFailed
			}

			if len(asstRst) > 0 { //one filed only asst to one module
				rstmap[item.PropertyID] = asstRst[0].AsstObjID
			}
		}

	}

	return rstmap, common.CCSuccess
}

//getInstDetail get inst detail
func (cli *instHelper) getInstAsst(req *restful.Request, ownerID, objID string, IDs []string, page map[string]interface{}) ([]instNameAsst, int, int) {

	tmpIDs := []int{}
	for _, ID := range IDs {
		tmpID, _ := strconv.Atoi(ID)
		tmpIDs = append(tmpIDs, tmpID)
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

	var targetobj string
	var instName string
	var instID string

	targetpre := cli.objCtrl + "/object/v1/insts/"
	switch objID {
	case common.BKInnerObjIDHost:
		targetobj = ""
		instName = common.BKHostInnerIPField
		instID = common.BKHostIDField
		targetpre = cli.hostCtrl + "/host/v1/hosts"
		if 0 != len(tmpIDs) {
			condition[common.BKHostIDField] = map[string]interface{}{"$in": tmpIDs}
		}
	case common.BKInnerObjIDApp:
		targetobj = common.BKInnerObjIDApp
		instName = common.BKAppNameField
		instID = common.BKAppIDField
		searchParams["sort"] = common.BKAppIDField
		condition[common.BKOwnerIDField] = ownerID
		if 0 != len(tmpIDs) {
			condition[common.BKAppIDField] = map[string]interface{}{"$in": tmpIDs}
		}
	case common.BKInnerObjIDSet:
		targetobj = common.BKInnerObjIDSet
		instID = common.BKSetIDField
		instName = common.BKSetNameField
		searchParams["sort"] = common.BKSetIDField
		condition[common.BKSetIDField] = map[string]interface{}{"$in": tmpIDs}
		condition[common.BKOwnerIDField] = ownerID
	case common.BKInnerObjIDModule:
		targetobj = common.BKObjIDField
		instID = common.BKModuleIDField
		instName = common.BKModuleNameField
		searchParams["sort"] = common.BKModuleIDField
		condition[common.BKOwnerIDField] = ownerID
		if 0 != len(tmpIDs) {
			condition[common.BKModuleIDField] = map[string]interface{}{"$in": tmpIDs}
		}
	case common.BKInnerObjIDPlat:
		targetobj = common.BKInnerObjIDPlat
		instID = common.BKCloudIDField
		instName = common.BKCloudNameField
		searchParams["sort"] = common.BKCloudIDField
		if 0 != len(tmpIDs) {
			condition[common.BKCloudIDField] = map[string]interface{}{"$in": tmpIDs}
		}
	default:
		targetobj = common.BKINnerObjIDObject
		instName = common.BKInstNameField
		instID = common.BKInstIDField
		condition[common.BKOwnerIDField] = ownerID
		condition[common.BKObjIDField] = objID
		if 0 != len(tmpIDs) {
			condition[common.BKInstIDField] = map[string]interface{}{"$in": tmpIDs}
		}
		searchParams["sort"] = common.BKInstIDField

	}

	searchParams["condition"] = condition

	//search
	sURL := targetpre + targetobj + "/search"
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

	delarry := func(s []string, i int) []string {
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

						// get inst name
						inst := instNameAsst{}
						if dataItemValStr, convOk := dataItemVal.(string); convOk {
							inst.Name = dataItemValStr
							inst.ObjID = objID
						}

						// delete exsit ID
						if dataItemVal, dataItemValOk := dataItem[instID]; dataItemValOk {

							//fmt.Println(key, dataItemVal)
							switch d := dataItemVal.(type) {
							case json.Number:
								if 0 != len(IDs) {
									for idx, key := range IDs {
										if val, err := d.Int64(); nil == err && key == strconv.Itoa(int(val)) {
											//fmt.Println("key:", key, "val:", val)
											inst.ID = IDs[idx]
											inst.ObjectID, _ = strconv.Atoi(IDs[idx])
											IDs = delarry(IDs, idx)
											rstName = append(rstName, inst)
											goto next
										}
									}
								} else if val, err := d.Int64(); nil == err {
									inst.ID = strconv.Itoa(int(val))
									inst.ObjectID = int(val)
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

	// get the InstName name
	for _, ID := range IDs {
		rstName = append(rstName, instNameAsst{ID: ID})
	}

	return rstName, cnt, common.CCSuccess
}

//GetInstDetailsSub get inst detail sub
func (cli *instHelper) GetInstDetailsSub(req *restful.Request, objID, ownerID string, input map[string]interface{}, page map[string]interface{}) (map[string]interface{}, int) {

	// get objID asst model and field
	rstmap, errorno := cli.GetObjectAsst(objID, ownerID)
	if common.CCSuccess != errorno {
		return nil, errorno
	}

	if info, infoOk := input["info"].([]interface{}); infoOk {
		for _, infoItem := range info {

			if dataItem, dataItemOk := infoItem.(map[string]interface{}); dataItemOk {

				// key 是关联字段，val 是字段关联的模型ID
				for key, objID := range rstmap {

					if keyItem, keyItemOk := dataItem[key]; keyItemOk {

						if nil == keyItem {
							continue
						}

						keyItemStr := fmt.Sprintf("%v", keyItem)
						blog.Debug("keyitemstr:%s", keyItemStr)
						retData, _, retErr := cli.getInstAsst(req, ownerID, objID, strings.Split(keyItemStr, ","), page)
						if common.CCSuccess != retErr {
							blog.Error("failed to get inst details")
						}
						dataItem[key] = retData
						fmt.Println("retData in: ", retData)

					}

				}
			}
		}
	}

	return input, common.CCSuccess
}

//GetInstDetails get inst detail
func (cli *instHelper) GetInstDetails(req *restful.Request, objID, ownerID, instStr string, page map[string]interface{}) (string, int) {

	rstmap, errorno := cli.GetObjectAsst(objID, ownerID)
	if common.CCSuccess != errorno {
		return "", errorno
	}
	js, err := simplejson.NewJson([]byte(instStr))
	if nil != err {
		blog.Error("the input json is invalid, error info is %s", err.Error())
		return "", common.CCErrCommJSONUnmarshalFailed
	}

	input, jsErr := js.Map()
	if nil != jsErr {
		blog.Error("the input json is invalid, error info is %s", jsErr.Error())
		return "", common.CCErrCommJSONUnmarshalFailed
	}

	blog.Debug("debug, input:%+v", input)

	if data, ok := input["data"].(map[string]interface{}); ok {

		// key 是关联字段，val 是字段关联的模型ID
		for key, val := range rstmap {

			if keyItem, keyItemOk := data[key]; keyItemOk {

				if nil == keyItem {
					continue
				}

				keyItemStr := fmt.Sprintf("%v", keyItem)
				blog.Debug("keyitemstr:%s", keyItemStr)
				retData, _, retErr := cli.getInstAsst(req, ownerID, val, strings.Split(keyItemStr, ","), page)
				if common.CCSuccess != retErr {
					blog.Error("failed to get inst details")
				}
				data[key] = retData
				blog.Debug("retData in: ", retData)

			}

		}

	}

	retData, _ := json.Marshal(input)
	return string(retData), common.CCSuccess
}
