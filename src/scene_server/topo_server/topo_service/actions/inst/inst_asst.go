package inst

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

//
// Associative instance query
//

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/api/metadata"
	"encoding/json"
	"strings"
)

func (cli *instAction) createInstAssociation(instAsst []interface{}) error {
	if 0 == len(instAsst) {
		return nil
	}
	return cli.CC.InstCli.InsertMuti(metadata.InstAsst{}.TableName(), instAsst...)
}

func (cli *instAction) updateInstAssociation(instID int, ownerID, objID string, input map[string]interface{}) error {

	// get association fields
	asst := map[string]interface{}{}
	asst[common.BKOwnerIDField] = ownerID
	asst[common.BKObjIDField] = objID
	searchData, _ := json.Marshal(asst)
	cli.objcli.SetAddress(cli.CC.ObjCtrl())
	asstDes, asstErr := cli.objcli.SearchMetaObjectAsst(searchData)
	if nil != asstErr {
		blog.Error("failed to search the obj asst, search condition(%+v) error info is %s", asst, asstErr.Error())
		return asstErr
	}
	// extract the data for the associated field
	asstFieldVal := make([]interface{}, 0)
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

					asstInst := metadata.InstAsst{}
					asstInst.InstID = int64(instID)
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

	err := cli.deleteInstAssociation(instID, ownerID, objID)
	if nil != err {
		blog.Errorf("faild to delete the old inst association, error info is %s", err.Error())
		return err
	}

	return cli.createInstAssociation(asstFieldVal)

}

func (cli *instAction) deleteInstAssociation(instID int, ownerID, objID string) error {

	return cli.CC.InstCli.DelByCondition(metadata.InstAsst{}.TableName(), map[string]interface{}{
		common.BKInstIDField: instID,
		common.BKObjIDField:  objID,
	})
}
