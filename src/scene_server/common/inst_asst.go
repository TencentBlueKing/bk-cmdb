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

package common

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/api"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/api/metadata"
	objapi "configcenter/src/source_controller/api/object"
	"encoding/json"
	"fmt"
	"strings"

	restful "github.com/emicklei/go-restful"
)

//ExtractDataFromAssociationField  get asst data
func ExtractDataFromAssociationField(instID int64, input map[string]interface{}, asstDes []objapi.ObjAsstDes) []*metadata.InstAsst {
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

func CreateInstAssociation(objAddr string, req *restful.Request, instAsst []*metadata.InstAsst) error {
	if 0 == len(instAsst) {
		return nil
	}
	for _, item := range instAsst {
		uURL := objAddr + "/object/v1/insts/" + common.BKTableNameInstAsst
		inputJSON, err := json.Marshal(item)
		if nil != err {
			return err
		}
		objRes, err := httpcli.ReqHttp(req, uURL, common.HTTPCreate, []byte(inputJSON))
		if nil != err {
			return err
		}

		if _, ok := IsSuccess([]byte(objRes)); !ok {
			blog.Error("failed to create the inst asst , error info is %s", objRes)
			continue
		}
	}
	return nil
}

// IsSuccess check the response
func IsSuccess(rst []byte) (*api.BKAPIRsp, bool) {

	var rstRes api.BKAPIRsp
	if jserr := json.Unmarshal(rst, &rstRes); nil != jserr {
		blog.Error("can not unmarshal the result , error: %s", jserr.Error())
		return &rstRes, false
	}

	if rstRes.Code != common.CCSuccess {
		return &rstRes, false
	}

	return &rstRes, true

}

func UpdateInstAssociation(objAddr string, req *restful.Request, instID int, ownerID, objID string, input map[string]interface{}) error {

	// get association fields
	asst := map[string]interface{}{}
	asst[common.BKOwnerIDField] = ownerID
	asst[common.BKObjIDField] = objID
	searchData, _ := json.Marshal(asst)
	objCli := objapi.NewClient("")
	objCli.SetAddress(objAddr)
	asstDes, asstErr := objCli.SearchMetaObjectAsst(&objapi.ForwardParam{Header: req.Request.Header}, searchData)
	if nil != asstErr {
		blog.Error("failed to search the obj asst, search condition(%+v) error info is %s", asst, asstErr.Error())
		return asstErr
	}

	for _, asst := range asstDes {
		if _, ok := input[asst.ObjectAttID]; ok {
			err := DeleteInstAssociation(objAddr, req, instID, ownerID, objID, asst.AsstObjID)
			if nil != err {
				blog.Errorf("faild to delete the old inst association, error info is %s", err.Error())
				return err
			}
		}
	}

	asstFieldVal := ExtractDataFromAssociationField(int64(instID), input, asstDes)

	return CreateInstAssociation(objAddr, req, asstFieldVal)

}

func DeleteInstAssociation(objAddr string, req *restful.Request, instID int, ownerID, objID, asstObjID string) error {

	uURL := objAddr + "/object/v1/insts/" + common.BKTableNameInstAsst
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

	if _, ok := IsSuccess([]byte(objRes)); !ok {
		blog.Error("failed to delete the inst asst , error info is %s", objRes)
		return fmt.Errorf("failed to delete the inst asst, %s:%d", objID, instID)
	}

	return nil
}

// SearchInstAssociation search the association
func SearchInstAssociation(objAdd string, cond map[string]interface{}, req *restful.Request) (string, error) {

	// search the association insts
	uURL := objAdd + "/object/v1/insts/" + common.BKTableNameInstAsst + "/search"
	input := map[string]interface{}{
		"page": map[string]interface{}{
			"start": 0,
			"sort":  "",
			"limit": common.BKNoLimit,
		},
		"condition": cond,
		"fields":    "",
	}

	inputJSON, _ := json.Marshal(input)
	return httpcli.ReqHttp(req, uURL, common.HTTPSelectPost, []byte(inputJSON))
}
