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

//
// Associative instance query
//

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/api/metadata"
	api "configcenter/src/source_controller/api/object"
	"encoding/json"
	"fmt"
	simplejson "github.com/bitly/go-simplejson"
	restful "github.com/emicklei/go-restful"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strings"
)

// ConditionItem subcondition
type ConditionItem struct {
	Field    string      `json:"field,omitempty"`
	Operator string      `json:"operator,omitempty"`
	Value    interface{} `json:"value,omitempty"`
}

// AssociationParams  association params
type AssociationParams struct {
	Page      metadata.BasePage          `json:"page,omitempty"`
	Fields    map[string][]string        `json:"fields,omitempty"`
	Condition map[string][]ConditionItem `json:"condition,omitempty"`
}

func (cli *instAction) createInstAssociation(req *restful.Request, instAsst []*metadata.InstAsst) error {
	if 0 == len(instAsst) {
		return nil
	}
	for _, item := range instAsst {
		uURL := cli.CC.ObjCtrl() + "/object/v1/insts/" + common.BKTableNameInstAsst
		inputJSON, err := json.Marshal(item)
		if nil != err {
			return err
		}
		objRes, err := httpcli.ReqHttp(req, uURL, common.HTTPCreate, []byte(inputJSON))
		if nil != err {
			return err
		}

		if _, ok := cli.IsSuccess([]byte(objRes)); !ok {
			blog.Error("failed to create the inst asst , error info is %s", objRes)
			continue
		}
	}
	return nil
}

func (cli *instAction) extractDataFromAssociationField(instID int64, input map[string]interface{}, asstDes []api.ObjAsstDes) []*metadata.InstAsst {
	// extract the data for the associated field
	asstFieldVal := make([]*metadata.InstAsst, 0)
	for idxItem, item := range asstDes {
		if inputVal, ok := input[item.ObjectAttID]; ok {
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

					asstInst := &metadata.InstAsst{}
					asstInst.InstID = int64(instID)
					asstInst.AsstInstID = iID
					asstInst.AsstObjectID = asstDes[idxItem].AsstObjID
					asstInst.ObjectID = asstDes[idxItem].ObjectID
					asstFieldVal = append(asstFieldVal, asstInst)
				}

			case int64, int:
				asstInst := &metadata.InstAsst{}
				asstInst.InstID = int64(instID)
				asstInst.AsstInstID, _ = util.GetInt64ByInterface(t)
				asstInst.AsstObjectID = asstDes[idxItem].AsstObjID
				asstInst.ObjectID = asstDes[idxItem].ObjectID
				asstFieldVal = append(asstFieldVal, asstInst)
			case json.Number:
				asstInst := &metadata.InstAsst{}
				asstInst.InstID = int64(instID)
				asstInst.AsstInstID, _ = t.Int64()
				asstInst.AsstObjectID = asstDes[idxItem].AsstObjID
				asstInst.ObjectID = asstDes[idxItem].ObjectID
				asstFieldVal = append(asstFieldVal, asstInst)

			default:
				blog.Warnf("the target data (%v) type is not a string ", t)
			}
		}
	}

	return asstFieldVal
}

func (cli *instAction) updateInstAssociation(req *restful.Request, instID int, ownerID, objID string, input map[string]interface{}) error {

	// get association fields
	asst := map[string]interface{}{}
	asst[common.BKOwnerIDField] = ownerID
	asst[common.BKObjIDField] = objID
	searchData, _ := json.Marshal(asst)
	cli.objcli.SetAddress(cli.CC.ObjCtrl())
	asstDes, asstErr := cli.objcli.SearchMetaObjectAsst(&api.ForwardParam{Header: req.Request.Header}, searchData)
	if nil != asstErr {
		blog.Error("failed to search the obj asst, search condition(%+v) error info is %s", asst, asstErr.Error())
		return asstErr
	}

	for _, asst := range asstDes {
		if _, ok := input[asst.ObjectAttID]; ok {
			err := cli.deleteInstAssociation(req, instID, ownerID, objID, asst.AsstObjID)
			if nil != err {
				blog.Errorf("faild to delete the old inst association, error info is %s", err.Error())
				return err
			}
		}
	}

	asstFieldVal := cli.extractDataFromAssociationField(int64(instID), input, asstDes)

	return cli.createInstAssociation(req, asstFieldVal)

}

func (cli *instAction) deleteInstAssociation(req *restful.Request, instID int, ownerID, objID, asstObjID string) error {

	uURL := cli.CC.ObjCtrl() + "/object/v1/insts/" + common.BKTableNameInstAsst
	input := map[string]interface{}{
		common.BKInstIDField: instID,
		common.BKObjIDField:  objID,
	}
	if "" != asstObjID {
		input[common.BKAsstObjIDField] = asstObjID
	}
	inputJSON, err := json.Marshal(input)
	if nil != err {
		return err
	}
	objRes, err := httpcli.ReqHttp(req, uURL, common.HTTPDelete, []byte(inputJSON))
	if nil != err {
		return err
	}

	if _, ok := cli.IsSuccess([]byte(objRes)); !ok {
		blog.Error("failed to delete the inst asst , error info is %s", objRes)
		return fmt.Errorf("failed to delete the inst asst, %s:%d", objID, instID)
	}

	return nil
}

func (cli *instAction) searchAssociationInst(req *restful.Request, objID string, searchParams map[string]interface{}) ([]int64, error) {
	// search the association object insts
	sURL := cli.CC.ObjCtrl() + "/object/v1/insts/" + util.GetObjByType(objID) + "/search"
	inputJSON, jsErr := json.Marshal(searchParams)

	if nil != jsErr {
		blog.Error("failed to marshal the data[%+v], error info is %s", searchParams, jsErr.Error())
		return nil, jsErr
	}

	objRes, err := httpcli.ReqHttp(req, sURL, common.HTTPSelectPost, inputJSON)
	blog.Debug("search url(%s) inst params: %s", sURL, string(inputJSON))

	if nil != err {
		blog.Errorf("failed to select the insts, error info is %s", err.Error())
		return nil, err
	}

	js, err := simplejson.NewJson([]byte(objRes))
	if nil != err {
		blog.Errorf("the input json is invalid, error info is %s", err.Error())
		return nil, err
	}

	input, jsErr := js.Map()
	if nil != jsErr {
		blog.Errorf("the input json is invalid, error info is %s", jsErr.Error())
		return nil, jsErr
	}

	instIDS := make([]int64, 0)
	if data, ok := input["data"].(map[string]interface{}); ok {
		if info, infoOk := data["info"].([]interface{}); infoOk {
			for _, infoItem := range info {

				if dataItem, dataItemOk := infoItem.(map[string]interface{}); dataItemOk {

					instID, keyItemOk := dataItem[util.GetObjIDByType(objID)]
					if keyItemOk {
						id, idErr := util.GetInt64ByInterface(instID)
						if nil != idErr {
							blog.Errorf("can not convert the instid to int64, error info is %s", idErr.Error())
							continue
						}

						instIDS = append(instIDS, id)
					}

				}

			}
		}
	}

	return instIDS, nil
}

// SelectInstsByAssociation associated with the search
func (cli *instAction) SelectInstsByAssociation(req *restful.Request, resp *restful.Response) {

	blog.Info("select association insts")
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

		var js AssociationParams
		err := json.Unmarshal([]byte(value), &js)
		if nil != err {
			blog.Error("failed to unmarshal the data[%s], error is %s", value, err.Error())
			return http.StatusBadRequest, "", defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		// used to search insts
		instCondition := map[string]interface{}{
			common.BKOwnerIDField: ownerID,
		}
		if util.GetObjByType(objID) == common.BKINnerObjIDObject {
			instCondition[common.BKObjIDField] = objID
		}

		// targetInstIDS used to search the instance of the objID
		targetInstIDS := make([]int64, 0)
		for keyObjID, objs := range js.Condition {

			// Extract the ID of the instance according to the associated object.
			condition := map[string]interface{}{}
			if util.GetObjByType(keyObjID) == common.BKINnerObjIDObject {
				condition[common.BKObjIDField] = keyObjID
				condition[common.BKOwnerIDField] = ownerID
			}

			for _, objCondition := range objs {

				if objCondition.Operator != common.BKDBEQ {

					if objID == keyObjID {
						// deal self condition
						instCondition[objCondition.Field] = map[string]interface{}{
							objCondition.Operator: objCondition.Value,
						}
					} else {
						// deal association condition
						condition[objCondition.Field] = map[string]interface{}{
							objCondition.Operator: objCondition.Value,
						}
					}
				} else {
					if objID == keyObjID {
						// deal self condition
						switch t := objCondition.Value.(type) {
						case string:
							instCondition[objCondition.Field] = map[string]interface{}{
								common.BKDBEQ: params.SpeceialCharChange(t),
							}
						default:
							instCondition[objCondition.Field] = objCondition.Value
						}

					} else {
						// deal association condition
						condition[objCondition.Field] = objCondition.Value
					}
				}

			}

			if objID == keyObjID {
				// no need to search the association objects
				continue
			}

			searchParams["fields"] = ""
			if fields, ok := js.Fields[keyObjID]; ok {
				searchParams["fields"] = strings.Join(fields, ",")
			}

			searchParams["condition"] = condition
			searchParams["start"] = 0
			searchParams["limit"] = common.BKNoLimit
			searchParams["sort"] = ""

			asstInstIDS, instIDSErr := cli.searchAssociationInst(req, keyObjID, searchParams)
			if nil != instIDSErr {
				blog.Errorf("can no search the association insts, error info is %s", instIDSErr.Error())
				continue
			}

			// search the association insts
			uURL := cli.CC.ObjCtrl() + "/object/v1/insts/" + common.BKTableNameInstAsst + "/search"
			input := map[string]interface{}{
				"page": map[string]interface{}{
					"start": 0,
					"sort":  "",
					"limit": common.BKNoLimit,
				},
				"condition": map[string]interface{}{
					"bk_asst_inst_id": map[string]interface{}{
						common.BKDBIN: asstInstIDS,
					},
					"bk_asst_obj_id": keyObjID,
					"bk_obj_id":      objID,
				},
				"fields": "",
			}

			inputJSON, _ := json.Marshal(input)
			objRes, err := httpcli.ReqHttp(req, uURL, common.HTTPSelectPost, []byte(inputJSON))
			if nil != err {
				blog.Errorf("failed to search the inst association, condition is %s ,error is %s", string(inputJSON), err.Error())
				continue
			}

			gjson.Get(objRes, "data.info.#."+common.BKInstIDField).ForEach(func(key, value gjson.Result) bool {

				targetInstIDS = append(targetInstIDS, value.Int())
				return true
			})

		}

		if instIDParam, ok := instCondition[common.BKInstIDField]; ok {
			delete(instCondition, common.BKInstIDField)
			instCondition[util.GetObjIDByType(objID)] = instIDParam
		}

		// search all the inst by the condition
		if 0 != len(targetInstIDS) {
			instCondition[util.GetObjIDByType(objID)] = map[string]interface{}{
				common.BKDBIN: targetInstIDS,
			}
		} else if 0 != len(js.Condition) {
			if _, ok := js.Condition[objID]; !ok {
				instCondition[util.GetObjIDByType(objID)] = map[string]interface{}{
					common.BKDBIN: targetInstIDS,
				}
			}
		}

		searchParams["condition"] = instCondition
		searchParams["fields"] = ""
		if fields, ok := js.Fields[objID]; ok {
			searchParams["fields"] = strings.Join(fields, ",")
		}

		searchParams["start"] = js.Page.Start
		searchParams["limit"] = js.Page.Limit
		searchParams["sort"] = js.Page.Sort

		// search insts
		sURL := cli.CC.ObjCtrl() + "/object/v1/insts/" + util.GetObjByType(objID) + "/search"
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
