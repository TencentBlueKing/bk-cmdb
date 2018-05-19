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
	//httpcli "configcenter/src/common/http/httpclient"
	api "configcenter/src/source_controller/api/object"
	"encoding/json"
	"fmt"
	//"github.com/tidwall/gjson"
	"sort"
	"strconv"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
	restful "github.com/emicklei/go-restful"
)

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
