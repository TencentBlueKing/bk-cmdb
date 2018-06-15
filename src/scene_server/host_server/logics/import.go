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
    "net/http"
    "context"
    "fmt"

    "configcenter/src/common/metadata"
    "configcenter/src/common"
    "configcenter/src/scene_server/host_server/service"
    "configcenter/src/common/util"
    "configcenter/src/common/blog"
)

type importInstance struct {
    defaultFields map[string]*metadata.ObjAttDes
    asstDes       []metadata.ObjectAssociations
}

func(lgc *Logics) AddHost(appID, moduleID int64, ownerID string, pheader http.Header, hostInfos map[int64]map[string]interface{}, importType metadata.HostInputType) ([]string, []string, []string, error) {
    
    instance := new(importInstance)
    
    
    instance.defaultFields = 
}

func(lgc *Logics) getHostFields(ownerID string, pheader http.Header) (map[string]*metadata.ObjAttDes, error) {
    page := metadata.BasePage{ Start: 0, Limit: common.BKNoLimit}
    searchOp := service.NewOperation().WithObjID(common.BKInnerObjIDHost).WithOwnerID(ownerID).WithPage(page).Data()
    result, err := lgc.CoreAPI.ObjectController().Meta().SelectObjectAttWithParams(context.Background(), pheader, searchOp)
    if err != nil || (err == nil && !result.Result) {
        return nil, fmt.Errorf("search host obj log failed, err: %v, result err: %s", err, result.ErrMsg)
    }
    
    atts := result.Data
    dat := make(map[string]*metadata.ObjAttDes)
    for idx, a := range atts {
        if !util.IsAssocateProperty(a.PropertyType) {
            continue
        }
        // read property group
        condition := map[string]interface{}{
            "bk_object_att_id":    a.PropertyID, // tmp.PropertyGroup,
            common.BKOwnerIDField: a.OwnerID,
            "bk_obj_id":           a.ObjectID,
        }
        objasstval, jserr := json.Marshal(condition)
        if nil != jserr {
            blog.Error("mashar json failed, error information is %v", jserr)
            return nil, jserr
        }
        asstMsg, err := client.SearchMetaObjectAsst(forward, objasstval)
        if nil != err {
            return nil, err
        }
        if 0 < len(asstMsg) {
            atts[idx].AssociationID = asstMsg[0].AsstObjID // by the rules, only one id
            atts[idx].AsstForward = asstMsg[0].AsstForward // by the rules, only one id
        }
    }
    return atts, nil
    
}



