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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"

	"github.com/emicklei/go-restful"
)

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
	objAttCondition := make(map[string]interface{})
	objAttCondition[common.BKObjIDField] = common.BKInnerObjIDProc
	objAttCondition[common.BKOwnerIDField] = ownerID
	attrQueryInput := new(meta.QueryCondition)
	attrQueryInput.Condition = objAttCondition
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
