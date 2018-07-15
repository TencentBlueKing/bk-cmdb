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
package service

import (
    "net/http"
    "strconv"
    "fmt"
    "context"
    "strings"
    
    "configcenter/src/common/util"
    "configcenter/src/common"
    "configcenter/src/common/blog"
    meta "configcenter/src/common/metadata"
    "configcenter/src/scene_server/validator"
    "configcenter/src/source_controller/api/metadata"
    "configcenter/src/common/auditoplog"
    "configcenter/src/common/paraparse"
    
    "github.com/gin-gonic/gin/json"
    "github.com/emicklei/go-restful"
)

func (ps *ProcServer) CreateProcess(req *restful.Request, resp *restful.Response) {
    user := util.GetActionUser(req)
    language := util.GetActionLanguage(req)
    defErr := ps.CCErr.CreateDefaultCCErrorIf(language)
    
    ownerID := req.PathParameter(common.BKOwnerIDField)
    appIDStr := req.PathParameter(common.BKAppIDField)
    appID, err := strconv.Atoi(appIDStr)
    if err != nil {
        blog.Errorf("convert appid from string to int failed!, err: %s", err.Error())
        resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
        return 
    }
    procIDStr := req.PathParameter(common.BKProcIDField)
    procID, err := strconv.Atoi(procIDStr)
    if err != nil {
        blog.Errorf("convert procid from string to int failed!, err: %s", err.Error())
        resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
        return
    }
    
    input := make(map[string]interface{})
    if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
        blog.Errorf("create process failed! decode request body err: %v", err)
        resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg:defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
        return 
    }

    input[common.BKAppIDField] = appID
    valid := validator.NewValidMap(common.BKDefaultOwnerID, common.BKInnerObjIDProc, req.Request.Header, ps.Engine)
    if err := valid.ValidMap(input, common.ValidCreate, 0); err != nil {
        blog.Errorf("fail to valid input parameters. err:%s", err.Error())
        resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg:defErr.Error(common.CCErrCommFieldNotValid)})
        return
    }
    
    // take snapshot before operation
    preProcDetail, err := ps.getProcDetail(req, ownerID, appID, procID)
    if err != nil {
        blog.Errorf("get process instance detail failed. err:%s", err.Error())
        //resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg:defErr.Error(common.CCErrAuditSaveLogFaile)})
        //return
    }
    
    input[common.BKOwnerIDField] = ownerID
    ret, err := ps.CoreAPI.ObjectController().Instance().CreateObject(context.Background(), common.BKInnerObjIDProc, req.Request.Header, input)
    if err != nil || (err == nil && !ret.Result) {
        blog.Errorf("create process failed . err: %s", err.Error())
        resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg:defErr.Error(common.CCErrProcCreateProcessFaile)})
        return
    }
    
    //save change log
    instID, ok := ret.Data["data."+common.BKProcIDField].(int)
    if !ok {
        blog.Errorf("instance id not found")
    }
    
    curDetail, err := ps.getProcDetail(req, ownerID, appID, instID)
    if err != nil {
        blog.Errorf("get process instance detail failed. err:%s", err.Error())
        //resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg:defErr.Error(common.CCErrAuditSaveLogFaile)})
        //return
    }
    
    ps.addProcLog(ownerID, appIDStr, user, preProcDetail, curDetail, instID, req.Request.Header)
    
    // return success
    result := make(map[string]interface{})
    result[common.BKProcIDField] = ret.Data[common.BKProcIDField]
    resp.WriteEntity(meta.NewSuccessResp(result))
}

func (ps *ProcServer) UpdateProcess(req *restful.Request, resp *restful.Response) {
    user := util.GetActionUser(req)
    language := util.GetActionLanguage(req)
    defErr := ps.CCErr.CreateDefaultCCErrorIf(language)

    ownerID := req.PathParameter(common.BKOwnerIDField)
    appIDStr := req.PathParameter(common.BKAppIDField)
    appID, err := strconv.Atoi(appIDStr)
    if err != nil {
        blog.Errorf("convert appid from string to int failed!, err: %s", err.Error())
        resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
        return
    }
    procIDStr := req.PathParameter(common.BKProcIDField)
    procID, err := strconv.Atoi(procIDStr)
    if err != nil {
        blog.Errorf("convert procid from string to int failed!, err: %s", err.Error())
        resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
        return
    }
    
    procData := make(map[string]interface{})
    if err := json.NewDecoder(req.Request.Body).Decode(&procData); err != nil {
        blog.Errorf("create process failed! decode request body err: %v", err)
        resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg:defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
        return
    }
    
    procData[common.BKAppIDField] = appID
    valid := validator.NewValidMap(common.BKDefaultOwnerID, common.BKInnerObjIDProc, req.Request.Header, ps.Engine)
    if err := valid.ValidMap(procData, common.ValidCreate, 0); err != nil {
        blog.Errorf("fail to valid input parameters. err:%s", err.Error())
        resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg:defErr.Error(common.CCErrCommFieldNotValid)})
        return
    }

    // take snapshot before operation
    preProcDetail, err := ps.getProcDetail(req, ownerID, appID, procID)
    if err != nil {
        blog.Errorf("get process instance detail failed. err:%s", err.Error())
    }
    
    input := make(map[string]interface{})
    condition := make(map[string]interface{})
    condition[common.BKOwnerIDField] = ownerID
    condition[common.BKAppIDField] = appID
    condition[common.BKProcIDField] = procID
    input["condition"] = condition
    input["data"] = procData
    ret, err := ps.CoreAPI.ObjectController().Instance().UpdateObject(context.Background(), common.BKInnerObjIDProc, req.Request.Header, input)
    if err != nil || (err == nil && !ret.Result) {
        blog.Errorf("update process failed . err: %s", err.Error())
        resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg:defErr.Error(common.CCErrProcUpdateProcessFaile)})
        return
    }
    
    // save operation log
    /// take snapshot before operation
    curDetail, err := ps.getProcDetail(req, ownerID, appID, procID)
    if err != nil {
        blog.Errorf("get process instance detail failed. err:%s", err.Error())
    }
    ps.addProcLog(ownerID, appIDStr, user, preProcDetail, curDetail, procID, req.Request.Header)

    resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (ps *ProcServer) BatchUpdateProcess(req *restful.Request, resp *restful.Response) {
    //user := util.GetActionUser(req)
    language := util.GetActionLanguage(req)
    defErr := ps.CCErr.CreateDefaultCCErrorIf(language)

    ownerID := req.PathParameter(common.BKOwnerIDField)
    appIDStr := req.PathParameter(common.BKAppIDField)
    appID, err := strconv.Atoi(appIDStr)
    if err != nil {
        blog.Errorf("convert appid from string to int failed!, err: %s", err.Error())
        resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
        return
    }

    procData := make(map[string]interface{})
    if err := json.NewDecoder(req.Request.Body).Decode(&procData); err != nil {
        blog.Errorf("create process failed! decode request body err: %v", err)
        resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg:defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
        return
    }
    
    // split process str by comma
    var procIDArr []string
    if procIDStr, ok := procData[common.BKProcessIDField].(string); !ok {
        resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg:defErr.Error(common.CCErrProcFieldValidFaile)})
        return
    } else {
        procIDArr = strings.Split(procIDStr, ",")
    }
    
    // update condition
    procData[common.BKAppIDField] = appID
    delete(procData, common.BKProcessIDField)

    // forbidden edit process name and func id
    delete(procData, common.BKProcessNameField)
    delete(procData, common.BKFuncIDField)

    // parse process id and valid
    var iProcIDArr []int
    //auditContentArr := make([]metadata.Content, len(procIDArr))
    
    valid := validator.NewValidMap(common.BKDefaultOwnerID, common.BKInnerObjIDProc, req.Request.Header, ps.Engine)
    
    for _, procIDStr := range procIDArr {
        procID, err := strconv.Atoi(procIDStr)
        if err != nil {
            resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
            return
        }
        
        if err := valid.ValidMap(procData, common.ValidUpdate, procID); err != nil {
            blog.Errorf("fail to valid proc parameters. err:%s", err.Error())
            resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg:defErr.Error(common.CCErrCommFieldNotValid)})
            return
        }
        
        iProcIDArr = append(iProcIDArr, procID)
    }

    // update processes
    input := make(map[string]interface{})
    condition := make(map[string]interface{})
    condition[common.BKOwnerIDField] = ownerID
    condition[common.BKAppIDField] = appID
    condition[common.BKProcIDField] = map[string]interface{}{
        common.BKDBIN: iProcIDArr,
    }

    input["condition"] = condition
    input["data"] = procData
    ret, err := ps.CoreAPI.ObjectController().Instance().UpdateObject(context.Background(), common.BKInnerObjIDProc, req.Request.Header, input)
    if err != nil || (err == nil && !ret.Result) {
        blog.Errorf("update process failed . err: %s", err.Error())
        resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg:defErr.Error(common.CCErrProcUpdateProcessFaile)})
        return
    }

    resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (ps *ProcServer) DeleteProcess(req *restful.Request, resp *restful.Response) {
    user := util.GetActionUser(req)
    language := util.GetActionLanguage(req)
    defErr := ps.CCErr.CreateDefaultCCErrorIf(language)

    ownerID := req.PathParameter(common.BKOwnerIDField)
    appIDStr := req.PathParameter(common.BKAppIDField)
    appID, err := strconv.Atoi(appIDStr)
    if err != nil {
        blog.Errorf("convert appid from string to int failed!, err: %s", err.Error())
        resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
        return
    }
    procIDStr := req.PathParameter(common.BKProcIDField)
    procID, err := strconv.Atoi(procIDStr)
    if err != nil {
        blog.Errorf("convert procid from string to int failed!, err: %s", err.Error())
        resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
        return
    }

    // take snapshot before operation
    preProcDetail, err := ps.getProcDetail(req, ownerID, appID, procID)
    if err != nil {
        blog.Errorf("get process instance detail failed. err:%s", err.Error())
    }

    conditon := make(map[string]interface{})
    conditon[common.BKAppIDField] = appID
    conditon[common.BKProcIDField] = procID
    conditon[common.BKOwnerIDField] = ownerID
    ret , err := ps.CoreAPI.ObjectController().Instance().DelObject(context.Background(), common.BKInnerObjIDProc, req.Request.Header, conditon)
    if err != nil || (err == nil && !ret.Result) {
        blog.Errorf("update process failed . err: %s", err.Error())
        resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg:defErr.Error(common.CCErrProcDeleteProcessFaile)})
        return
    }

    // save operation log
    ps.addProcLog(ownerID, appIDStr, user, preProcDetail, nil, procID, req.Request.Header)

    resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (ps *ProcServer) SearchProcess(req *restful.Request, resp *restful.Response) {
    language := util.GetActionLanguage(req)
    defErr := ps.CCErr.CreateDefaultCCErrorIf(language)

    ownerID := req.PathParameter(common.BKOwnerIDField)
    appIDStr := req.PathParameter(common.BKAppIDField)
    appID, err := strconv.Atoi(appIDStr)
    if err != nil {
        blog.Errorf("convert appid from string to int failed!, err: %s", err.Error())
        resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
        return
    }
    
    var srchparam params.SearchParams
    if err := json.NewDecoder(req.Request.Body).Decode(&srchparam); err != nil {
        blog.Errorf("decode request body err: %v", err)
        resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg:defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
        return
    }

    condition := srchparam.Condition
    condition[common.BKOwnerIDField] = ownerID
    condition[common.BKAppIDField] = appID
    if processName, ok := condition["bk_process_name"]; ok {
        processNameStr, ok := processName.(string)
        if ok {
            condition["bk_process_name"] = map[string]interface{}{common.BKDBLIKE: params.SpeceialCharChange(processNameStr)}

        } else {
            condition["bk_process_name"] = map[string]interface{}{common.BKDBLIKE: processName}

        }
    }
    page := srchparam.Page
    searchParams := new(meta.QueryInput)
    searchParams.Condition = condition
    searchParams.Fields = strings.Join(js.Fields, ",")
    searchParams.Start = page["start"].(int)
    searchParams.Limit = page["limit"].(int)
    searchParams.Sort = page["sort"].(string)
    
    // query process by module name
    if moduleName, ok := condition[common.BKModuleNameField]; ok {
        reqParam := make(map[string]interface{})
        reqParam[common.BKAppIDField] = appID
        reqParam[common.BKModuleNameField] = moduleName
        
        blog.Infof("get process arr by module(%s), reqParam: %+v", moduleName, reqParam)
        ret, err := ps.CoreAPI.ProcController().GetProc2Module(context.Background(), req.Request.Header, reqParam)
        if err != nil || (err == nil && !ret.Result) {
            blog.Errorf("query process by module failed. err: %s", err.Error())
            resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg:defErr.Error(common.CCErrProcSearchProcessFaile)})
            return
        }
        
        // parse procid array
        procIdArr := make([]int, 0)
        for _, item := range ret.Data {
            procIdArr = append(procIdArr, item.ProcessID)
        }
        
        // update condition
        condition[common.BKProcIDField] = map[string]interface{}{
            "$in": procIdArr,
        }
        delete(condition, common.BKModuleNameField)
        searchParams.Condition = condition
    }
    
    // search process
    ret, err := ps.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDProc, req.Request.Header, searchParams)
    if err != nil || (err == nil && !ret.Result) {
        blog.Errorf("search process failed . err: %s", err.Error())
        resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg:defErr.Error(common.CCErrProcSearchProcessFaile)})
        return
    }

    resp.WriteEntity(meta.NewSuccessResp(ret.Data))
}

func (ps *ProcServer) addProcLog(ownerID, appID, user string, preProcDetails, curProcDetails []map[string]interface{}, instanceID int, header http.Header) {
    headers := []metadata.Header{}
    curData := map[string]interface{}{}
    preData := map[string]interface{}{}
    for _, detail := range preProcDetails {
        preData[detail[common.BKPropertyIDField].(string)] = detail[common.BKPropertyValueField]
        if nil == curProcDetails {
            headers = append(headers,
                metadata.Header{
                    PropertyID:   fmt.Sprint(detail[common.BKPropertyIDField]),
                    PropertyName: fmt.Sprint(detail[common.BKPropertyNameField]),
                })
        }
        }
    
    for _, detail := range curProcDetails {
        curData[detail[common.BKPropertyIDField].(string)] = detail[common.BKPropertyValueField]
        headers = append(headers,
            metadata.Header{
                PropertyID: fmt.Sprint(detail[common.BKPropertyIDField]),
                PropertyName: fmt.Sprint(detail[common.BKPropertyNameField]),
            })
    }

    auditContent := metadata.Content{
        CurData: curData,
        PreData: preData,
        Headers: headers,
    }

    log := common.KvMap{common.BKContentField: auditContent, common.BKOpDescField:"create process", common.BKOpTypeField: auditoplog.AuditOpTypeAdd, "inst_id": instanceID}
    ps.CoreAPI.AuditController().AddProcLog(context.Background(), ownerID, appID, user, header, log)
}