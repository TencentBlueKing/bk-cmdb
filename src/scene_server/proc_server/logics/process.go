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
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (lgc *Logics) GetProcbyProcIDArr(ctx context.Context, procID []int64, header http.Header) ([]mapstr.MapStr, error) {
	condition := map[string]interface{}{
		common.BKProcessIDField: mapstr.MapStr{common.BKDBIN: procID},
	}

	reqParam := new(metadata.QueryInput)
	reqParam.Condition = condition
	ret, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(ctx, common.BKInnerObjIDProc, header, reqParam)
	if err != nil || (err == nil && !ret.Result) {
		return nil, fmt.Errorf("get process by procID(%+v) failed. err: %v", procID, err)
	}
	if !ret.Result {
		return nil, lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header)).New(ret.Code, ret.ErrMsg)
	}

	if len(ret.Data.Info) < 1 {
		return nil, fmt.Errorf("there is no process with procID(%d)", procID)
	}

	return ret.Data.Info, nil
}

func (lgc *Logics) getProcInfoByID(ctx context.Context, procID []int64, header http.Header) (map[int64]*metadata.InlineProcInfo, error) {
	ownerID := util.GetOwnerID(header)
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	if 0 == len(procID) {
		return nil, nil
	}
	gseProc := make(map[int64]*metadata.InlineProcInfo, 0)
	dat := new(metadata.QueryInput)
	dat.Condition = mapstr.MapStr{common.BKProcessIDField: mapstr.MapStr{common.BKDBIN: procID}}
	dat.Limit = common.BKNoLimit
	ret, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(ctx, common.BKInnerObjIDProc, header, dat)
	if nil != err {
		blog.Errorf("getProcInfoByID procID %v supplierID %s  http do error:%s, logID::%s", procID, ownerID, err.Error(), util.GetHTTPCCRequestID(header))
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !ret.Result {
		blog.Errorf("getProcInfoByID procID %v supplierID %s  http reply error:%s, logID:%s", procID, ownerID, ret.ErrMsg, util.GetHTTPCCRequestID(header))
		return nil, defErr.New(ret.Code, ret.ErrMsg)

	}
	if 0 == ret.Data.Count {
		blog.Errorf("getProcInfoByID procID %v supplierID %s  not found process info, logID:%s", procID, ownerID, util.GetHTTPCCRequestID(header))
		return nil, nil
	}
	for _, proc := range ret.Data.Info {
		procID, err := proc.Int64(common.BKProcessIDField)
		if nil != err {
			byteHost, _ := json.Marshal(proc)
			blog.Errorf("getHostByModuleID  proc %v  procID  not interger, json:%s, logID:%s", proc, string(byteHost), util.GetHTTPCCRequestID(header))
			return nil, err
		}
		item := new(metadata.InlineProcInfo)

		item.ProcNum = 1 //ProcNum not set, use default value 1
		procNumI, ok := proc.Get(common.BKProcInstNum)
		if ok && nil != procNumI {
			item.ProcNum, err = proc.Int64(common.BKProcInstNum)
			if nil != err {
				byteHost, _ := json.Marshal(proc)
				blog.Errorf("getHostByModuleID  proc %v  procNum  not interger, json:%s, logID:%s", proc, string(byteHost), util.GetHTTPCCRequestID(header))
				return nil, err
			}
		}
		item.AppID, err = proc.Int64(common.BKAppIDField)
		if nil != err {
			byteHost, _ := json.Marshal(proc)
			blog.Errorf("getHostByModuleID  proc info  AppID  not interger, error:%s, json:%s, logID:%s", err.Error(), string(byteHost), util.GetHTTPCCRequestID(header))
			return nil, err
		}
		item.FunID, err = proc.Int64(common.BKFuncIDField)
		if nil != err {
			continue
		}
		item.ProcID = procID
		item.ProcInfo = proc

		gseProc[procID] = item
	}

	return gseProc, nil
}
