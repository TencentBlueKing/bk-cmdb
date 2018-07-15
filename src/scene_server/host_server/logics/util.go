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

package logics

import (
	"encoding/json"
	"strings"

	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func ExtractDataFromAssociationField(instID int64, input map[string]interface{}, asstDes []metadata.Association) []*metadata.InstAsst {
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

			case int64, int, float64:
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
