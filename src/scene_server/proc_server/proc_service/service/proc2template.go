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
	"context"
	"net/http"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	types "configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
	//	"github.com/gin-gonic/gin/json"
)

func (ps *ProcServer) GetProcBindTemplate(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := ps.CCErr.CreateDefaultCCErrorIf(language)

	pathParams := req.PathParameters()
	appIDStr := pathParams[common.BKAppIDField]
	appID, _ := strconv.Atoi(appIDStr)
	procIDStr := pathParams[common.BKProcessIDField]
	procID, _ := strconv.Atoi(procIDStr)

	// search object instance
	condition := make(map[string]interface{})
	condition[common.BKAppIDField] = appID
	input := new(meta.QueryInput)
	input.Condition = condition

	tempRet, err := ps.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDConfigTemp, req.Request.Header, input)
	if err != nil || !tempRet.Result {
		blog.Errorf("fail to GetProcBindTemplate when do searchobject. err:%v, errcode:%d, errmsg:%s", err, tempRet.Code, tempRet.ErrMsg)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrObjectSelectInstFailed)})
		return
	}

	condition[common.BKProcessIDField] = procID

	// get process to templation by condition
	proc2TempRet, err := ps.CoreAPI.ProcController().SearchProc2Template(context.Background(), req.Request.Header, condition)
	if err != nil || !proc2TempRet.Result {
		blog.Errorf("fail to GetProcessTemplate when do GetProc2Template. err:%v, errcode:%d, errmsg:%s", err, proc2TempRet.Code, proc2TempRet.ErrMsg)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcSelectBindToMoudleFaile)})
		return
	}

	result := make([]interface{}, 0)

	for _, temp := range tempRet.Data.Info {
		iTempID, err := util.GetInt64ByInterface(temp[common.BKTemlateIDField])
		if nil != err {
			continue
		}
		iTempName, ok := temp[common.BKTemplateNameField].(string)
		if false == ok {
			continue
		}
		iFileName, false := temp[common.BKFileNameField].(string)
		if false == ok {
			continue
		}
		isBind := 0
		for _, proc2Temp := range proc2TempRet.Data {
			jTempID, err := util.GetInt64ByInterface(proc2Temp[common.BKTemlateIDField])
			if nil != err {
				continue
			}
			if iTempID == jTempID {
				isBind = 1
			}
		}
		result = append(result, types.MapStr{common.BKTemplateNameField: iTempName, common.BKFileNameField: iFileName, "is_bind": isBind})
	}
	resp.WriteEntity(meta.NewSuccessResp(result))
}
