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

    "configcenter/src/common/backbone"
    "configcenter/src/common"
    "configcenter/src/source_controller/api/metadata"
    "fmt"
    "configcenter/src/framework/core/errors"
    "configcenter/src/scene_server/host_server/service"
)

type Logics struct {
    *backbone.Engine
}

func(lgc *Logics) GetHostAttributes(ownerID string, header http.Header) ([]metadata.Header, error){
    searchOp := service.NewOperation().WithObjID(common.BKInnerObjIDHost).WithOwnerID(ownerID).Data()
    result, err := lgc.CoreAPI.ObjectController().Meta().SelectObjectAttWithParams(context.Background(), header, searchOp)
    if err != nil || (err == nil && !result.Result){
        return nil, fmt.Errorf("search host obj log failed, err: %v, result err: %s", err, result.Message)
    }

    hostAttrArr, ok := result.Data.([]interface{})
    if !ok {
        return nil, errors.New("invalid response data")
    }
    
    headers := make([]metadata.Header, 0)
    for _, i := range hostAttrArr {
        attr := i.(map[string]interface{})
        data := metadata.Header{}
        propertyID := attr[common.BKPropertyIDField].(string)
        if propertyID == common.BKChildStr {
            continue
        }
        data.PropertyID = propertyID
        data.PropertyName = attr[common.BKPropertyNameField].(string)

        headers = append(headers, data)
    }
    
    return headers, nil
}

func(lgc *Logics) GenerateHostLogs(ownerID string, hostID string, logHeaders []metadata.Header, pheader http.Header) (*metadata.Content, error) {
    ctnt := new(metadata.Content)
    ctnt.Headers = logHeaders
    
    // get host details
    result, err := lgc.CoreAPI.HostController().Host().GetHostByID(context.Background(), hostID, pheader)
    if err != nil || (err == nil && !result.Result) {
        return nil, fmt.Errorf("get host pre data failed, err, %v, %v", err, result.ErrMsg)
    }
    
    hostInfo, ok := result.Data.(map[string]interface{})
    if !ok {
        return nil, errors.New("invalid host info data")
    }
    
    // get host association
    opt := service.NewOperation().WithOwnerID(ownerID).WithObjID(common.BKInnerObjIDHost)
    assResult, err := lgc.CoreAPI.ObjectController().Meta().SelectObjectAttWithParams(context.Background(), pheader, opt)
    if err != nil || (err == nil && !assResult.Result) {
        return nil, fmt.Errorf("get host association failed, err, %v, %v", err, result.ErrMsg)
    }
    
    
    
}


