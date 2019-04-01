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
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (lgc *Logics) getModuleNameByID(ctx context.Context, ID int64, header http.Header) (name string, appID int64, setID int64, err error) {
	supplierID := util.GetOwnerID(header)
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	dat := new(metadata.QueryInput)
	dat.Condition = mapstr.MapStr{common.BKModuleIDField: ID}
	dat.Fields = fmt.Sprintf("%s,%s,%s", common.BKModuleNameField, common.BKAppIDField, common.BKSetIDField)
	dat.Limit = common.BKNoLimit
	ret, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(ctx, common.BKInnerObjIDModule, header, dat)
	if nil != err {
		blog.Errorf("getModuleNameByID moduleID %v supplierID %s  http do error:%s", ID, supplierID, err.Error())
		return "", 0, 0, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !ret.Result {
		blog.Errorf("getModuleNameByID moduleID %v supplierID %s  http reply error:%s", ID, supplierID, ret.ErrMsg)
		return "", 0, 0, defErr.New(ret.Code, ret.ErrMsg)

	}
	if 0 == ret.Data.Count {
		blog.Errorf("getModuleNameByID moduleID %v supplierID %s  not found module info", ID, supplierID)
		return "", 0, 0, nil
	}
	byteModule, _ := json.Marshal(ret.Data.Info[0])
	name, err = ret.Data.Info[0].String(common.BKModuleNameField)
	if nil != err {
		blog.Errorf("getModuleNameByID moduleID %v supplierID %s  get module name error:%s, raw:%s", ID, supplierID, err.Error(), string(byteModule))
		return
	}
	appID, err = ret.Data.Info[0].Int64(common.BKAppIDField)
	if nil != err {
		blog.Errorf("getModuleNameByID moduleID %v supplierID %s  get appID name error:%s raw:%s", ID, supplierID, err.Error(), string(byteModule))
		return "", 0, 0, err
	}
	setID, err = ret.Data.Info[0].Int64(common.BKSetIDField)
	if nil != err {
		blog.Errorf("getModuleNameByID moduleID %v supplierID %s  get set name error:%s raw:%s", ID, supplierID, err.Error(), string(byteModule))
		return "", 0, 0, err
	}

	return

}

func (lgc *Logics) getModuleIDByProcID(ctx context.Context, appID, procID int64, header http.Header) ([]int64, error) {
	condition := make(map[string]interface{}, 0)
	condition[common.BKProcessIDField] = procID
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	rid := util.GetHTTPCCRequestID(header)
	// get process by module
	ret, err := lgc.CoreAPI.ProcController().GetProc2Module(context.Background(), header, condition)
	if nil != err {
		blog.Errorf("getModuleIDByProcID  GetProc2Module appID %d moduleID %d  http do error:%s,rid:%s", appID, procID, err.Error(), rid)
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !ret.Result {
		blog.Errorf("getModuleIDByProcID  GetProc2Module appID %d moduleID %d  http reply error:%s,rid:%s", appID, procID, ret.ErrMsg, rid)
		return nil, defErr.New(ret.Code, ret.ErrMsg)
	}
	var moduleIDs []int64
	for _, item := range ret.Data {
		ids, err := lgc.HandleProcInstNumByModuleName(ctx, header, appID, item.ModuleName)
		if nil != err {
			blog.Errorf("getModuleIDByProcID get module id by module name %s  in application id  %d error %s,rid:%s", item.ModuleName, item.ApplicationID, err.Error(), rid)
			return nil, err
		}
		moduleIDs = append(moduleIDs, ids...)
	}

	return moduleIDs, nil
}

func (lgc *Logics) GetModuleIDByHostID(ctx context.Context, header http.Header, hostID int64) ([]metadata.ModuleHost, error) {
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	dat := map[string][]int64{
		common.BKHostIDField: []int64{hostID},
	}
	ret, err := lgc.CoreAPI.HostController().Module().GetModulesHostConfig(ctx, header, dat)
	if nil != err {
		blog.Errorf("GetModuleIDByHostID hostID id %d GetModulesHostConfig http do error:%s", hostID, err.Error())
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !ret.Result {
		blog.Errorf("GetModuleIDByHostID  hostID %d GetModulesHostConfig reply error:%s", hostID, ret.ErrMsg)
		return nil, defErr.New(ret.Code, ret.ErrMsg)
	}

	return ret.Data, nil
}

func (lgc *Logics) GetModueleIDByAppID(ctx context.Context, header http.Header, appID int64) ([]int64, error) {
	supplierID := util.GetOwnerID(header)
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	dat := new(metadata.QueryInput)
	dat.Condition = mapstr.MapStr{common.BKAppIDField: appID}
	dat.Fields = fmt.Sprintf("%s,%s,%s,%s", common.BKModuleNameField, common.BKAppIDField, common.BKSetIDField, common.BKModuleIDField)
	dat.Limit = common.BKNoLimit
	ret, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(ctx, common.BKInnerObjIDModule, header, dat)
	if nil != err {
		blog.Errorf("GetModueleIDByAppID appID %v supplierID %s  http do error:%s", appID, supplierID, err.Error())
		return make([]int64, 0), defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !ret.Result {
		blog.Errorf("GetModueleIDByAppID appID %v supplierID %s  http reply error:%s", appID, supplierID, ret.ErrMsg)
		return make([]int64, 0), defErr.New(ret.Code, ret.ErrMsg)

	}
	if 0 == ret.Data.Count {
		blog.V(5).Infof("GetModueleIDByAppID appID %v supplierID %s  not found module info", appID, supplierID)
		return make([]int64, 0), nil
	}
	moduleIDs := make([]int64, 0)
	for _, module := range ret.Data.Info {
		moduleID, err := module.Int64(common.BKModuleIDField)
		if nil != err {
			byteModule, _ := json.Marshal(module)
			blog.Errorf("GetModueleIDByAppID moduleID %v supplierID %s  get set name error:%s raw:%s", appID, supplierID, err.Error(), string(byteModule))
			return nil, err
		}
		moduleIDs = append(moduleIDs, moduleID)
	}

	return moduleIDs, nil
}

func (lgc *Logics) GetAppList(ctx context.Context, header http.Header, fields string) ([]mapstr.MapStr, error) {
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	dat := new(metadata.QueryInput)
	if "" != strings.TrimSpace(fields) {
		dat.Fields = fields
	}
	dat.Limit = common.BKNoLimit
	ret, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(ctx, common.BKInnerObjIDApp, header, dat)
	if nil != err {
		blog.Errorf("GetAppList  http do error:%s", err.Error())
		return make([]mapstr.MapStr, 0), defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !ret.Result {
		blog.Errorf("GetAppList  http reply error:%s", ret.ErrMsg)
		return make([]mapstr.MapStr, 0), defErr.New(ret.Code, ret.ErrMsg)

	}
	if 0 == ret.Data.Count {
		blog.V(5).Infof("GetAppList  not found app info")
		return make([]mapstr.MapStr, 0), nil
	}

	return ret.Data.Info, nil
}
