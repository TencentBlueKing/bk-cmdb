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
package procserver

import (
    "github.com/emicklei/go-restful"
    "configcenter/src/common/util"
    "strconv"
    "configcenter/src/common"
    "configcenter/src/common/blog"
    meta "configcenter/src/common/metadata"
    sourceAPI "configcenter/src/source_controller/api/object"
    "net/http"
    "github.com/gin-gonic/gin/json"
    "context"
)

func (ps *ProcServer) GetProcessPortByApplicationID (req *restful.Request, resp *restful.Response) {
    language := util.GetActionLanguage(req)
    defErr := ps.CCErr.CreateDefaultCCErrorIf(language)
    
    //get appID
    pathParams := req.PathParameters()
    appID, err := strconv.Atoi(pathParams[common.BKAppIDField])
    if err != nil {
        blog.Errorf("fail to get appid from pathparameter. err: %s", err.Error())
        resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
        return
    }
    
    bodyData := make([]map[string]interface{}, 0)
    if err := json.NewDecoder(req.Request.Body).Decode(&bodyData); err != nil {
        blog.Errorf("fail to decode request body. err: %s", err.Error())
        resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg:defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
        return 
    }
    
    modules := bodyData
    // 根据模块获取所有关联的进程，建立Map ModuleToProcesses
    moduleToProcessesMap := make(map[int][]interface{})
    for _, module := range modules {
        moduleName, ok := module[common.BKModuleNameField].(string)
        if !ok {
            blog.Warnf("assign error module['ModuleName'] is not string, module:%v", module)
            continue
        }
    }
}

func (ps *ProcServer) GetProcessPortByIP (req *restful.Request, resp *restful.Response) {
    language := util.GetActionLanguage(req)
    defErr := ps.CCErr.CreateDefaultCCErrorIf(language)
}

// 根据模块获取所有关联的进程，建立Map ModuleToProcesses
func (ps *ProcServer) getProcessesByModuleName(forward *sourceAPI.ForwardParam, moduleName string) ([]interface{}, error) {
    procData := make([]interface{}, 0)
    params := map[string]interface{}{
        common.BKModuleNameField: moduleName,
    }
    
    ret, err := ps.CoreAPI.ObjectController().OpenAPI().GetProcessesByModuleName(context.Background(), forward.Header, params)
    if err != nil || (err == nil && !ret.Result) {
        blog.Errorf("get process by module failed. err: %s, errcode: %d, errmsg: %s", err.Error(), ret.Code, ret.ErrMsg)
        return procData, err
    }
    
    procData = append(procData, ret.Data)
    return procData, nil
}