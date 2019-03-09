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
	"encoding/json"
	"net/http"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"

	"github.com/emicklei/go-restful"
)

func (ps *ProcServer) GetProcessDetailByID(req *restful.Request, resp *restful.Response) {
	srvData := ps.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr
	ownerID := req.PathParameter(common.BKOwnerIDField)
	appIDStr := req.PathParameter(common.BKAppIDField)
	appID, err := strconv.Atoi(appIDStr)
	if err != nil {
		blog.Errorf("convert appid from string to int failed!, err: %s,appID:%v,rid:%s", err.Error(), appIDStr, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}
	procIDStr := req.PathParameter(common.BKProcessIDField)
	procID, err := strconv.Atoi(procIDStr)
	if err != nil {
		blog.Errorf("convert procid from string to int failed!, err: %s,procID:%s,rid:%s", err.Error(), procIDStr, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	ret, err := ps.getProcDetail(req, ownerID, appID, procID)
	if err != nil {
		blog.Errorf("GetProcessDetailByID info err: %v,appID:%v,procID:%v,rid:%s", err, appID, procID, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcSearchDetailFaile)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(ret))
}

func (ps *ProcServer) getProcDetail(req *restful.Request, ownerID string, appID, procID int) ([]map[string]interface{}, error) {
	srvData := ps.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr
	// search process
	procCondition := make(map[string]interface{})
	procCondition[common.BKOwnerIDField] = ownerID
	procCondition[common.BKAppIDField] = appID
	procCondition[common.BKProcessIDField] = procID
	searchParams := new(meta.QueryCondition)
	searchParams.Condition = procCondition
	retObj, err := ps.CoreAPI.CoreService().Instance().ReadInstance(srvData.ctx, srvData.header, common.BKInnerObjIDProc, searchParams)
	if err != nil {
		blog.Errorf("getProcDetail http do error.err:%s,input:%+v,rid:%s", err.Error(), searchParams, srvData.rid)
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !retObj.Result {
		blog.Errorf("getProcDetail http reply error.err code:%d,err msg:%s,input:%+v,rid:%s", retObj.Code, retObj.ErrMsg, searchParams, srvData.rid)
		return nil, defErr.New(retObj.Code, retObj.ErrMsg)
	}

	proc := make(map[string]interface{})
	for _, item := range retObj.Data.Info {
		for k, v := range item {
			proc[k] = v
		}
	}

	// search objectatts
	objattCondition := make(map[string]interface{})
	objattCondition[common.BKObjIDField] = common.BKInnerObjIDProc
	objattCondition[common.BKOwnerIDField] = ownerID
	attrQueryInput := new(meta.QueryCondition)
	attrQueryInput.Condition = objattCondition
	retObjAtt, err := ps.CoreAPI.CoreService().Model().ReadModelAttr(srvData.ctx, srvData.header, common.BKInnerObjIDProc, attrQueryInput)
	if err != nil {
		blog.Errorf("getProcDetail SelectObjectAttWithParams http do error.err:%s,input:%+v,rid:%s", err.Error(), searchParams, srvData.rid)
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !retObjAtt.Result {
		blog.Errorf("getProcDetail SelectObjectAttWithParams http reply error.err code:%d,err msg:%s,input:%+v,rid:%s", retObjAtt.Code, retObjAtt.ErrMsg, searchParams, srvData.rid)
		return nil, defErr.New(retObjAtt.Code, retObjAtt.ErrMsg)
	}

	reResult := make([]map[string]interface{}, 0)
	for _, item := range retObjAtt.Data.Info {
		data := make(map[string]interface{})
		propertyID := item.PropertyID
		if propertyID == common.BKChildStr {
			continue
		}

		data[common.BKPropertyIDField] = propertyID
		data[common.BKPropertyNameField] = item.PropertyName
		data[common.BKPropertyValueField] = proc[propertyID]
		reResult = append(reResult, data)
	}

	return reResult, nil
}

type instNameAsst struct {
	ID         string `json:"id"`
	ObjID      string `json:"bk_obj_id"`
	ObjIcon    string `json:"bk_obj_icon"`
	InstID     int    `json:"bk_inst_id"`
	ObjectName string `json:"bk_obj_name"`
	InstName   string `json:"bk_inst_name"`
}

func (ps *ProcServer) getInstAsst(forward http.Header, ownerID, objID string, ids []string, page map[string]interface{}) ([]instNameAsst, int, int) {
	srvData := ps.newSrvComm(forward)

	tmpIDS := make([]int, 0)
	for _, id := range ids {
		tmpID, _ := strconv.Atoi(id)
		tmpIDS = append(tmpIDS, tmpID)
	}

	input := new(meta.QueryInput)
	condition := make(map[string]interface{})
	input.Fields = ""
	if val, ok := page["fields"].(string); ok {
		input.Fields = val
	}
	input.Start = 0
	if val, ok := page["start"].(int); ok {
		input.Start = val
	}
	input.Limit = common.BKDefaultLimit
	if val, ok := page["limit"].(int); ok {
		input.Limit = val
	}

	var targetOBJ string
	var instName string
	var instID string

	switch objID {
	case common.BKInnerObjIDHost:
		targetOBJ = ""
		instName = common.BKHostInnerIPField
		instID = common.BKHostIDField
		if 0 != len(tmpIDS) {
			condition[common.BKHostIDField] = map[string]interface{}{common.BKDBIN: tmpIDS}
		}
	case common.BKInnerObjIDApp:
		targetOBJ = common.BKInnerObjIDApp
		instName = common.BKAppNameField
		instID = common.BKAppIDField
		input.Sort = common.BKAppIDField
		condition[common.BKOwnerIDField] = ownerID
		if 0 != len(tmpIDS) {
			condition[common.BKAppIDField] = map[string]interface{}{common.BKDBIN: tmpIDS}
		}
	case common.BKInnerObjIDSet:
		targetOBJ = common.BKInnerObjIDSet
		instID = common.BKSetIDField
		instName = common.BKSetNameField
		input.Sort = common.BKSetIDField
		condition[common.BKSetIDField] = map[string]interface{}{common.BKDBIN: tmpIDS}
		condition[common.BKOwnerIDField] = ownerID
	case common.BKInnerObjIDModule:
		targetOBJ = common.BKInnerObjIDModule
		instID = common.BKModuleIDField
		instName = common.BKModuleNameField
		input.Sort = common.BKModuleIDField
		condition[common.BKOwnerIDField] = ownerID
		if 0 != len(tmpIDS) {
			condition[common.BKModuleIDField] = map[string]interface{}{common.BKDBIN: tmpIDS}
		}
	case common.BKInnerObjIDPlat:
		targetOBJ = common.BKInnerObjIDPlat
		instID = common.BKCloudIDField
		instName = common.BKCloudNameField
		input.Sort = common.BKCloudIDField
		if 0 != len(tmpIDS) {
			condition[common.BKCloudIDField] = map[string]interface{}{common.BKDBIN: tmpIDS}
		}
	default:
		targetOBJ = common.BKInnerObjIDObject
		instName = common.BKInstNameField
		instID = common.BKInstIDField
		condition[common.BKOwnerIDField] = ownerID
		condition[common.BKObjIDField] = objID
		if 0 != len(tmpIDS) {
			condition[common.BKInstIDField] = map[string]interface{}{common.BKDBIN: tmpIDS}
		}
		input.Sort = common.BKInstIDField
	}

	input.Condition = condition

	var dataInfo []mapstr.MapStr
	cnt := 0
	switch objID {
	case common.BKInnerObjIDHost:
		hostRet, err := ps.CoreAPI.HostController().Host().GetHosts(srvData.ctx, srvData.header, input)
		if err != nil || (err == nil && !hostRet.Result) {
			blog.Errorf("search inst detail failed when GetHosts, err: %v,input:%+v,rid:%s", err, input, srvData.rid)
			return nil, 0, common.CCErrHostSelectInst
		}
		dataInfo = hostRet.Data.Info
		cnt = hostRet.Data.Count
	default:
		queryCondtion := &meta.QueryCondition{
			Condition: condition,
		}
		objRet, err := ps.CoreAPI.CoreService().Instance().ReadInstance(srvData.ctx, srvData.header, targetOBJ, queryCondtion)
		if err != nil || (err == nil && !objRet.Result) {
			blog.Errorf("search inst detail failed when SearchObjects, err: %v,input:%+v,rid:%s", err, input, srvData.rid)
			return nil, 0, common.CCErrObjectSelectInstFailed
		}
		cnt = objRet.Data.Count
		for _, val := range objRet.Data.Info {
			dataInfo = append(dataInfo, val)
		}
	}

	delArrayFunc := func(s []string, i int) []string {
		s[len(s)-1], s[i] = s[i], s[len(s)-1]
		return s[:len(s)-1]
	}

	rstName := []instNameAsst{}

	for _, infoItem := range dataInfo {
		dataItemVal := infoItem[instName]
		// 提取实例名
		inst := instNameAsst{}
		if dataItemValStr, ok := dataItemVal.(string); ok {
			inst.InstName = dataItemValStr
			inst.ObjID = objID
		}
		// 删除已经存在的ID
		dataItemDelVal := infoItem[instID]
		switch d := dataItemDelVal.(type) {
		case json.Number:
			if 0 != len(ids) {
				for idx, key := range ids {
					if val, err := d.Int64(); nil == err && key == strconv.Itoa(int(val)) {
						inst.ID = ids[idx]
						inst.InstID, _ = strconv.Atoi(ids[idx])
						ids = delArrayFunc(ids, idx)
						rstName = append(rstName, inst)
					}
				}
			} else if val, err := d.Int64(); nil == err {
				inst.ID = strconv.Itoa(int(val))
				inst.InstID = int(val)
				rstName = append(rstName, inst)
			}
		}
	}

	// deal the other inst name
	for _, id := range ids {
		rstName = append(rstName, instNameAsst{ID: id})
	}

	return rstName, cnt, common.CCSuccess
}
