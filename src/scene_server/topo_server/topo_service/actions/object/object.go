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

package object

import (
	"configcenter/src/common"
	"configcenter/src/common/bkbase"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/topo_service/manager"
	"strings"

	api "configcenter/src/source_controller/api/object"
	"encoding/json"

	"github.com/bitly/go-simplejson"

	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/emicklei/go-restful"
)

var obj = &objectAction{}

// HostAction
type objectAction struct {
	base.BaseAction
	mgr manager.Manager
}

func init() {

	// register action
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/object/batch", Params: nil, Handler: obj.CreateObjectBatch})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/object/search/batch", Params: nil, Handler: obj.SearchObjectBatch})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/object", Params: nil, Handler: obj.CreateObject})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/objects", Params: nil, Handler: obj.SelectObjectWithParams})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/objects/topo", Params: nil, Handler: obj.SelectObjectTopo})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/object/{id}", Params: nil, Handler: obj.UpdateObject})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/object/{id}", Params: nil, Handler: obj.DeleteObject})

	// create cc object
	obj.CreateAction()

	// set httpclient
	manager.SetManager(obj)

}

// SetManager implement the manager's Hooker interface
func (cli *objectAction) SetManager(mgr manager.Manager) error {
	cli.mgr = mgr
	return nil
}

var errEmpty = fmt.Errorf("empty string")

func (cli *objectAction) updateObjectAttribute(tmpItem *api.ObjAttDes, jsObjAttr *simplejson.Json) (*api.ObjAttDes, error) {

	if !tmpItem.IsPre {
		// is not the inner field
		if jsTmp, ok := jsObjAttr.CheckGet("bk_asst_obj_id"); ok {
			if tmp, tmpErr := jsTmp.String(); nil == tmpErr {
				tmpItem.AssociationID = tmp
			} else {
				blog.Error("can not parse the bk_asst_obj_id, error info is %s", tmpErr.Error())
				return tmpItem, tmpErr
			}
		}

		if jsTmp, ok := jsObjAttr.CheckGet("editable"); ok {
			if tmp, tmpErr := jsTmp.Bool(); nil == tmpErr {
				tmpItem.Editable = tmp
			} else {
				blog.Error("can not parse the editable, error info is %s", tmpErr.Error())
				return tmpItem, tmpErr
			}
		}

		if jsTmp, ok := jsObjAttr.CheckGet(common.BKIsOnlyField); ok {
			if tmp, tmpErr := jsTmp.Bool(); nil == tmpErr {
				tmpItem.IsOnly = tmp
			} else {
				blog.Error("can not parse the isonly, error info is %s", tmpErr.Error())
				return tmpItem, tmpErr
			}
		}

		if jsTmp, ok := jsObjAttr.CheckGet("isrequired"); ok {
			if tmp, tmpErr := jsTmp.Bool(); nil == tmpErr {
				tmpItem.IsRequired = tmp
			} else {
				blog.Error("can not parse the isrequired, error info is %s", tmpErr.Error())
				return tmpItem, tmpErr
			}
		}

		if jsTmp, ok := jsObjAttr.CheckGet("isreadonly"); ok {
			if tmp, tmpErr := jsTmp.Bool(); nil == tmpErr {
				tmpItem.IsReadOnly = tmp
			} else {
				blog.Error("can not parse the isreadonly, error info is %s", tmpErr.Error())
				return tmpItem, tmpErr
			}
		}

		if jsTmp, ok := jsObjAttr.CheckGet("bk_issystem"); ok {
			if tmp, tmpErr := jsTmp.Bool(); nil == tmpErr {
				tmpItem.IsSystem = tmp
			} else {
				blog.Error("can not parse the bk_issystem, error info is %s", tmpErr.Error())
				return tmpItem, tmpErr
			}
		}
		if jsTmp, ok := jsObjAttr.CheckGet("bk_isapi"); ok {
			if tmp, tmpErr := jsTmp.Bool(); nil == tmpErr {
				tmpItem.IsAPI = tmp
			} else {
				blog.Error("can not parse the bk_issystem, error info is %s", tmpErr.Error())
				return tmpItem, tmpErr
			}
		}

	}

	// base property field reset
	if jsTmp, ok := jsObjAttr.CheckGet("bk_property_id"); ok {
		if tmp, tmpErr := jsTmp.String(); nil == tmpErr {
			tmp = strings.TrimSpace(tmp)
			if "" == tmp {
				blog.Error("bk_property_id could not be empty")
				return tmpItem, errEmpty
			}
			tmpItem.PropertyID = tmp
		} else {
			blog.Error("can not parse the bk_property_ids, error info is %s", tmpErr.Error())
			return tmpItem, tmpErr
		}
	}
	if jsTmp, ok := jsObjAttr.CheckGet("bk_property_name"); ok {
		if tmp, tmpErr := jsTmp.String(); nil == tmpErr {
			if "" == tmp {
				blog.Error("bk_property_name could not be empty")
				return tmpItem, errEmpty
			}
			tmpItem.PropertyName = tmp
		} else {
			blog.Error("can not parse the bk_property_name, error info is %s", tmpErr.Error())
			return tmpItem, tmpErr
		}
	}

	if jsTmp, ok := jsObjAttr.CheckGet("unit"); ok {
		if tmp, tmpErr := jsTmp.String(); nil == tmpErr {
			tmpItem.Unit = tmp
		} else {
			blog.Error("can not parse the unit, error info is %s", tmpErr.Error())
			return tmpItem, tmpErr
		}
	}

	if jsTmp, ok := jsObjAttr.CheckGet("placeholder"); ok {
		if tmp, tmpErr := jsTmp.String(); nil == tmpErr {
			tmpItem.Placeholder = tmp
		} else {
			blog.Error("can not parse the placeholder, error info is %s", tmpErr.Error())
			return tmpItem, tmpErr
		}
	}

	if jsTmp, ok := jsObjAttr.CheckGet("bk_property_type"); ok {
		if tmp, tmpErr := jsTmp.String(); nil == tmpErr {
			if "" == tmp {
				blog.Error("bk_property_type could not be empty")
				return tmpItem, errEmpty
			}
			tmpItem.PropertyType = tmp
		} else {
			blog.Error("can not parse the bk_property_type, error info is %s", tmpErr.Error())
			return tmpItem, tmpErr
		}
	}

	if jsTmp, ok := jsObjAttr.CheckGet("option"); ok {
		switch tmpItem.PropertyType {
		case common.FieldTypeEnum:
			if tmp, tmpErr := jsTmp.Array(); nil == tmpErr {
				tmpItem.Option = tmp
			} else {
				blog.Error("can not parse the option, error info is %s", tmpErr.Error())
				return tmpItem, tmpErr
			}
		default:
			tmpItem.Option = jsTmp
		}

		/*if tmp, tmpErr := jsTmp.String(); nil == tmpErr {
			tmpItem.Option = tmp
		} else {
			blog.Error("can not parse the option, error info is %s", tmpErr.Error())
			return tmpItem, tmpErr
		}*/
	}

	return tmpItem, nil
}

// CreateObjectBatch batch to create the object attributes
func (cli *objectAction) CreateObjectBatch(req *restful.Request, resp *restful.Response) {
	blog.Info("batch to create objectatt information")
	/**
	input:
	the inputed data format
	{
		"objid":{
			"meta":{
				"key":"val"
			},
			"attr":{
				"0":{
					"key":"val"
				}
			}
		}
	}
	*/

	/**
	result:
	{
		"objid": {

		"success":[],
		"update_failed":[],
		"insert_failed":[],
		}

	}
	*/

	// get the owner id
	ownerID := util.GetActionOnwerID(req)

	// get the language
	language := util.GetActionLanguage(req)

	// get the error info by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	// execute
	cli.CallResponseEx(func() (int, interface{}, error) {
		forward := &api.ForwardParam{Header: req.Request.Header}
		js, jsErr := simplejson.NewFromReader(req.Request.Body)
		if nil != jsErr {
			blog.Error("unmarshal the json, error info is %s", jsErr.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		jsMap, jsMapErr := js.Map()
		if nil != jsMapErr {
			blog.Error("failed to map the json, error info is %s", jsMapErr.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		result := map[string]interface{}{}

		// parse the json get the object id
		for objID := range jsMap {

			subResult := map[string]interface{}{}

			// check the object
			condition := map[string]interface{}{}
			condition[common.BKOwnerIDField] = ownerID
			condition[common.BKObjIDField] = objID

			conditionVal, _ := json.Marshal(condition)
			if items, err := cli.mgr.SelectObject(forward, conditionVal, defErr); nil != err {
				blog.Error("failed to search, the error info is :%s", err.Error())
				subResult["errors"] = fmt.Sprintf("the object(%s) is invalid", objID)
				result[objID] = subResult
				continue
			} else if 0 == len(items) {
				// TODO: may be need to create the object in the future version
				blog.Error("not found the  objid: %s", objID)
				subResult["errors"] = fmt.Sprintf("the object(%s) is invalid", objID)
				result[objID] = subResult
				continue
			}

			// update the object attribute
			conditionAtt := map[string]interface{}{}
			jsObjAttr := js.Get(objID).Get("attr")
			attrMap, mapErr := jsObjAttr.Map()
			if nil != mapErr {
				blog.Error("can not convert to map, error info is %s", mapErr.Error())
				subResult["errors"] = fmt.Sprintf("lost the attr field in the input stream")
				result[objID] = subResult
				continue
			}

			for keyIdx := range attrMap {

				propertyID, err := jsObjAttr.Get(keyIdx).Get("bk_property_id").String()
				if nil != err {
					blog.Error("failed to parse the bk_property_id, error info is %s", err.Error())
					if failed, ok := subResult["insert_failed"]; ok {
						failedArr := failed.([]string)
						failedArr = append(failedArr, fmt.Sprintf("line:%s msg: lost the field(bk_property_id", keyIdx))
						subResult["insert_failed"] = failedArr
					} else {
						subResult["insert_failed"] = []string{
							fmt.Sprintf("line:%s msg: lost the field(bk_property_id", keyIdx),
						}
					}
					result[objID] = subResult
					continue
				}

				// check the property id
				conditionAtt[common.BKOwnerIDField] = ownerID
				conditionAtt[common.BKObjIDField] = objID
				conditionAtt["bk_property_id"] = propertyID

				conditionAttVal, _ := json.Marshal(conditionAtt)
				if items, err := cli.mgr.SelectObjectAtt(forward, conditionAttVal, defErr); nil != err {
					blog.Error("failed to search the object attribute, the condition is %+v, error info is %s", conditionAtt, err.Error())
					if failed, ok := subResult["insert_failed"]; ok {
						failedArr := failed.([]string)
						failedArr = append(failedArr, fmt.Sprintf("line:%s msg: %s", keyIdx, err.Error()))
						subResult["insert_failed"] = failedArr
					} else {
						subResult["insert_failed"] = []string{
							fmt.Sprintf("line:%s msg: %s", keyIdx, err.Error()),
						}
					}
					result[objID] = subResult
					continue

				} else if 0 != len(items) {

					// need to update
					for _, tmpItem := range items {

						item, itemErr := cli.updateObjectAttribute(&tmpItem, jsObjAttr.Get(keyIdx))
						if nil != itemErr {
							blog.Error("failed to reset the object attribute, error info is %s ", itemErr.Error())
							if failed, ok := subResult["update_failed"]; ok {
								failedArr := failed.([]string)
								failedArr = append(failedArr, fmt.Sprintf("line:%s msg: %s", keyIdx, itemErr.Error()))
								subResult["update_failed"] = failedArr
							} else {
								subResult["update_failed"] = []string{
									fmt.Sprintf("line:%s msg: %s", keyIdx, itemErr.Error()),
								}
							}
							result[objID] = subResult
							continue
						}

						itemVal, _ := json.Marshal(item)
						blog.Debug("the new attribute:%s", string(itemVal))
						if updateErr := cli.mgr.UpdateObjectAtt(forward, item.ID, itemVal, defErr); nil != updateErr {
							blog.Error("failed to update the object attribute, error info is %s", updateErr.Error())
							if failed, ok := subResult["update_failed"]; ok {
								failedArr := failed.([]string)
								failedArr = append(failedArr, fmt.Sprintf("line:%s msg: %s", keyIdx, updateErr.Error()))
								subResult["update_failed"] = failedArr
							} else {
								subResult["update_failed"] = []string{
									fmt.Sprintf("line:%s msg: %s", keyIdx, updateErr.Error()),
								}
							}
							result[objID] = subResult
							continue
						}
					}

				} else {
					// need to create
					tmpItem := &api.ObjAttDes{}
					tmpItem.ObjectID = objID
					tmpItem.OwnerID = ownerID
					item, itemErr := cli.updateObjectAttribute(tmpItem, jsObjAttr.Get(keyIdx))
					if nil != itemErr {
						blog.Error("failed to reset the object attribute, error info is %s ", itemErr.Error())
						if failed, ok := subResult["insert_failed"]; ok {
							failedArr := failed.([]string)
							failedArr = append(failedArr, fmt.Sprintf("line:%s msg: %s", keyIdx, itemErr.Error()))
							subResult["insert_failed"] = failedArr
						} else {
							subResult["insert_failed"] = []string{
								fmt.Sprintf("line:%s msg: %s", keyIdx, itemErr.Error()),
							}
						}
						result[objID] = subResult
						continue
					}

					if _, insertErr := cli.mgr.CreateObjectAtt(forward, *item, defErr); nil != insertErr {
						blog.Error("failed to create the object attribute, error info is %s", insertErr.Error())
						if failed, ok := subResult["insert_failed"]; ok {
							failedArr := failed.([]string)
							failedArr = append(failedArr, fmt.Sprintf("line:%s msg: %s", keyIdx, insertErr.Error()))
							subResult["insert_failed"] = failedArr
						} else {
							subResult["insert_failed"] = []string{
								fmt.Sprintf("line:%s msg: %s", keyIdx, insertErr.Error()),
							}
						}
						result[objID] = subResult
						continue
					}

				} // end else  create attribute

				if failed, ok := subResult["success"]; ok {
					failedArr := failed.([]string)
					failedArr = append(failedArr, keyIdx)
					subResult["success"] = failedArr
				} else {
					subResult["success"] = []string{
						keyIdx,
					}
				}

				result[objID] = subResult

			} // end foreach objid
		}

		return http.StatusOK, result, nil
	}, resp)
}

// SearchObjectBatch batch to search the object attributes
func (cli *objectAction) SearchObjectBatch(req *restful.Request, resp *restful.Response) {
	blog.Info("batch to search objectatt information")

	/**
	input:
		{
	    "objid": {
	        "meta": {
	            "key": "val"
	        },
	        "attr": {
	            "key": "val"
	        }
	    },
	    "failed":[
	        {
	           "objids":"error info"
	        }
			]

		}

	*/

	/*
			result:
			{
		    "objid": {
		        "meta": {
		            "key": "val"
		        },
		        "attr": [{
		            "key": "val"
		        }]
		    },
		    "failed":[
		        {
		           "objids":"error info"
		        }
		    ]

			}

	*/
	// get the owner id
	ownerID := util.GetActionOnwerID(req)

	// get the language
	language := util.GetActionLanguage(req)

	// get the error info by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	// execute
	cli.CallResponseEx(func() (int, interface{}, error) {

		js, jsErr := simplejson.NewFromReader(req.Request.Body)
		if nil != jsErr {
			blog.Error("unmarshal the json, error info is %s", jsErr.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		conditions, conErr := js.Get("condition").Array()
		if nil != conErr {
			blog.Error("failed to get the conditions, error info is %s", conErr.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		result := map[string]interface{}{}
		// parse the conditons
		for _, item := range conditions {
			switch objID := item.(type) {
			case string:
				{
					conditionAttr := map[string]interface{}{}
					conditionAttr[common.BKOwnerIDField] = ownerID
					conditionAttr[common.BKObjIDField] = objID
					conditionAttrVal, _ := json.Marshal(conditionAttr)
					forward := &api.ForwardParam{Header: req.Request.Header}
					attrItems, attrErr := cli.mgr.SelectObjectAtt(forward, conditionAttrVal, defErr)
					if nil != attrErr {
						blog.Error("failed to search the attribute of the object(%s) ownerID(%s), error info is %s", objID, ownerID, attrErr.Error())
						if failed, ok := result["failed"]; ok {
							failedArr := failed.([]map[string]string)
							failedArr = append(failedArr, map[string]string{
								objID: attrErr.Error(),
							})
						} else {
							result["failed"] = []map[string]string{
								map[string]string{
									objID: attrErr.Error(),
								},
							}
						}

						continue
					}

					if 0 != len(attrItems) {
						result[objID] = map[string]interface{}{
							"attr": attrItems,
						}
					} else {
						// not found the attribute, set the default(empty one)
						result[objID] = api.ObjAttDes{}
					}
				}
			default:
				blog.Error("unsupport the type for the data(%v)", objID)
				if failed, ok := result["failed"]; ok {
					failedArr := failed.([]map[string]string)
					failedArr = append(failedArr, map[string]string{
						fmt.Sprint(objID): "unsupport the type",
					})
				} else {
					result["failed"] = []map[string]string{
						map[string]string{
							fmt.Sprint(objID): "unsupport the type",
						},
					}
				}
				return http.StatusBadRequest, nil, defErr.Errorf(common.CCErrCommParamsNeedString, "condition")
			}
		}

		return http.StatusOK, result, nil
	}, resp)
}

func (cli *objectAction) CreateObject(req *restful.Request, resp *restful.Response) {

	blog.Info("create object information")

	// get the language
	language := util.GetActionLanguage(req)

	// get the error info by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	// logics
	cli.CallResponseEx(func() (int, interface{}, error) {
		forward := &api.ForwardParam{Header: req.Request.Header}
		// read body
		val, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("failed to read request body, error info is %s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		id, idErr := cli.mgr.CreateObject(forward, val, defErr)
		if nil == idErr {
			return http.StatusOK, map[string]int{"id": id}, nil
		}

		blog.Error("create object failed, error information is %v", idErr.Error())
		return http.StatusOK, nil, defErr.Error(common.CCErrTopoObjectCreateFailed)

	}, resp)
}
func (cli *objectAction) SelectObjectTopo(req *restful.Request, resp *restful.Response) {

	blog.Info("select object topo")

	// get the language
	language := util.GetActionLanguage(req)

	// get the error info by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	forward := &api.ForwardParam{Header: req.Request.Header}

	isFrom := func(fromObjID, toObjID, ownerID string) (bool, error) {
		search := map[string]interface{}{}
		search[common.BKOwnerIDField] = ownerID
		search[common.BKObjIDField] = fromObjID

		asstRst, asstErr := cli.mgr.SelectObjectAsst(forward, search, defErr)
		if nil != asstErr {
			blog.Error("failed to search object topo, error info is %s", asstErr.Error())
			return false, asstErr
		}

		for _, asstItem := range asstRst {
			if asstItem.AsstObjID == toObjID {
				return true, nil
			}
		}

		return false, nil
	}
	// execute
	cli.CallResponseEx(func() (int, interface{}, error) {

		// read body
		val, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("failed to read request body, error info is %s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		// search
		rst, rstErr := cli.mgr.SelectObject(forward, val, defErr)
		if nil != rstErr {
			blog.Error("failed to select the object, error info is %s", rstErr.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoObjectSelectFailed)
		}

		resultTopo := make([]ObjectTopo, 0)
		for _, rstItem := range rst {
			search := map[string]interface{}{}
			search[common.BKOwnerIDField] = rstItem.OwnerID
			search[common.BKObjIDField] = rstItem.ObjectID

			asstRst, asstErr := cli.mgr.SelectObjectAsst(forward, search, defErr)
			if nil != asstErr {
				blog.Error("failed to search object topo, error info is %s", asstErr.Error())
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoObjectSelectFailed)
			}
			blog.Debug("asst rst :%+v", asstRst)
			for _, asstItem := range asstRst {
				if asstItem.ObjectAttID == common.BKChildStr {
					continue
				}
				//asstItem.ObjectID
				condition := make(map[string]interface{})
				condition[common.BKObjIDField] = asstItem.AsstObjID
				condition[common.BKOwnerIDField] = asstItem.OwnerID
				conditionStr, _ := json.Marshal(condition)
				objItems, objErr := cli.mgr.SelectObject(forward, conditionStr, defErr)
				if nil != objErr {
					blog.Error("failed to search object topo, error info is %s", objErr.Error())
					return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoObjectSelectFailed)
				}
				if 0 < len(objItems) {
					tmp := ObjectTopo{}
					tmp.Label = asstItem.ObjectAttID
					tmp.LabelName = asstItem.AsstName
					tmp.From.ObjID = rstItem.ObjectID
					tmp.From.ClassificationID = rstItem.ObjCls
					tmp.From.Position = rstItem.Position
					tmp.From.OwnerID = rstItem.OwnerID
					tmp.From.ObjName = rstItem.ObjectName
					tmp.To.OwnerID = asstItem.OwnerID
					tmp.To.ObjID = asstItem.AsstObjID
					tmp.To.ClassificationID = objItems[0].ObjCls
					tmp.To.Position = objItems[0].Position
					tmp.To.ObjName = objItems[0].ObjectName
					ok, okErr := isFrom(asstItem.AsstObjID, rstItem.ObjectID, rstItem.OwnerID)
					if nil != okErr {
						return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoObjectSelectFailed)
					}
					if ok {
						tmp.Arrows = "to,from"
					} else {
						tmp.Arrows = "to"
					}

					blog.Debug("tmp association:%+v", tmp)
					resultTopo = append(resultTopo, tmp)

				}
			}
		}
		blog.Debug("result: %+v", resultTopo)
		return http.StatusOK, resultTopo, nil

	}, resp)

}
func (cli *objectAction) SelectObjectWithParams(req *restful.Request, resp *restful.Response) {

	blog.Info("select object information with params")

	// get the language
	language := util.GetActionLanguage(req)

	// get the error info by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	forward := &api.ForwardParam{Header: req.Request.Header}

	// logics
	cli.CallResponseEx(func() (int, interface{}, error) {

		// read body
		val, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("failed to read request body, error info is %s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		// search
		rst, rstErr := cli.mgr.SelectObject(forward, val, defErr)
		if nil == rstErr {
			return http.StatusOK, rst, nil
		}

		//cli.mgr.SelectObjectAsst()
		blog.Error("failed to search object, error info is %s", rstErr.Error())
		return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoObjectSelectFailed)

	}, resp)
}

func (cli *objectAction) UpdateObject(req *restful.Request, resp *restful.Response) {

	blog.Info("update object information")

	// get the language
	language := util.GetActionLanguage(req)

	// get the error info by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	forward := &api.ForwardParam{Header: req.Request.Header}

	// execute
	cli.CallResponseEx(func() (int, interface{}, error) {

		// read body
		val, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("failed to read request body, error info is %s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		// get params
		id, idErr := strconv.Atoi(req.PathParameter("id"))
		if nil != idErr {
			blog.Error("the path param(id:%s), error info is %s", req.PathParameter("id"), idErr.Error())
			return http.StatusBadRequest, nil, defErr.Errorf(common.CCErrCommParamsNeedInt, "id")
		}

		// execute
		updateErr := cli.mgr.UpdateObject(forward, id, val, defErr)
		if nil == updateErr {
			return http.StatusOK, nil, nil
		}

		blog.Error("failed to update object, error info is %s", updateErr.Error())
		return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoObjectUpdateFailed)

	}, resp)
}

func (cli *objectAction) DeleteObject(req *restful.Request, resp *restful.Response) {

	blog.Info("delete object information")

	// get the language
	language := util.GetActionLanguage(req)

	// get the error info by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	forward := &api.ForwardParam{Header: req.Request.Header}

	// execute
	cli.CallResponseEx(func() (int, interface{}, error) {

		// read body
		val, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("failed to read request body, error info is %s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		// get params
		id, idErr := strconv.Atoi(req.PathParameter("id"))
		if nil != idErr {
			blog.Error("the path param(id:%s), error info is %s", req.PathParameter("id"), idErr.Error())
			return http.StatusBadRequest, nil, defErr.Errorf(common.CCErrCommParamsNeedInt, "id")
		}

		// logics
		delErr := cli.mgr.DeleteObject(forward, id, val, defErr)
		if nil == delErr {
			return http.StatusOK, nil, nil
		}
		blog.Errorf("failed to delete object(%d), error info is %s", id, delErr.Error())
		return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoObjectDeleteFailed)
	}, resp)
}
