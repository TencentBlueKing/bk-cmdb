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

package identifier

import (
	"configcenter/src/scene_server/event_server/types"
)

func GetDistInst(e *types.EventInst) []types.DistInst {
	distinst := types.DistInst{
		EventInst: *e,
	}
	distinst.ID = 0
	var ds []types.DistInst
	var m map[string]interface{}
	if e.EventType == types.EventTypeInstData && e.ObjType == common.BKINnerObjIDObject {
		var ok bool

		if e.Action == "delete" && len(e.Data) > 0 {
			m, ok = e.Data[0].PreData.(map[string]interface{})
		} else {
			m, ok = e.Data[0].CurData.(map[string]interface{})
		}
		if !ok {
			return nil
		}

		if m[common.BKObjIDField] != nil {
			distinst.ObjType = m[common.BKObjIDField].(string)
		}

	}
	ds = append(ds, distinst)

	// add new dist if event belong to hostidentifier
	if diffFields, ok := hostIndentDiffFiels[e.ObjType]; ok && e.Action == types.EventActionUpdate && e.EventType == types.EventTypeInstData {
		for dataIndex := range e.Data {
			curdata := e.Data[dataIndex].CurData.(map[string]interface{})
			predata := e.Data[dataIndex].PreData.(map[string]interface{})
			if checkDifferent(curdata, predata, diffFields...) {
				hostIdentify := types.DistInst{
					EventInst: *e,
				}
				hostIdentify.Data = nil
				hostIdentify.EventType = types.EventTypeRelation
				hostIdentify.ObjType = "hostidentifier"

				instID, _ := curdata[common.GetInstIDField(e.ObjType)].(int)
				if instID == 0 {
					// this should wound happen -_-
					blog.Errorf("conver instID faile the raw is %v", curdata[common.GetInstIDField(e.ObjType)])
					continue
				}
				count := 0
				identifiers := GetIdentifierCache().getCache(e.ObjType, instID)
				total := len(identifiers)
				// pack identifiers into 1 distribution to prevent send too many messages
				for ident := range identifiers {
					count++
					d := types.EventData{PreData: *ident}
					for _, field := range diffFields {
						ident.Set(field, curdata[field])
					}
					d.CurData = *ident
					hostIdentify.Data = append(hostIdentify.Data, d)
					// each group is divided into 1000 units in order to limit the message size
					if count%1000 == 0 || count == total {
						ds = append(ds, hostIdentify)
						hostIdentify.Data = nil
					}
				}
			}
		}
	} else if e.EventType == types.EventTypeRelation && distinst.ObjType == "moduletransfer" {

	}

	return ds
}
