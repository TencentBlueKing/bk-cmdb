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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

func (lgc *Logics) getModuleNameByID(ctx context.Context, ID int64) (name string, appID int64, setID int64, err error) {
	supplierID := lgc.ownerID
	defErr := lgc.ccErr
	dat := new(metadata.QueryCondition)
	dat.Condition = mapstr.MapStr{common.BKModuleIDField: ID}
	dat.Fields = []string{common.BKModuleNameField, common.BKAppIDField, common.BKSetIDField}
	dat.Limit.Limit = common.BKNoLimit
	ret, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(ctx, lgc.header, common.BKInnerObjIDModule, dat)
	if nil != err {
		blog.Errorf("getModuleNameByID moduleID %v supplierID %s  http do error:%s,query:%+v,rid:%s", ID, supplierID, err.Error(), dat, lgc.rid)
		return "", 0, 0, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !ret.Result {
		blog.Errorf("getModuleNameByID moduleID %v supplierID %s  http reply error:%s,query:%+v,rid:%s", ID, supplierID, ret.ErrMsg, dat, lgc.rid)
		return "", 0, 0, defErr.New(ret.Code, ret.ErrMsg)

	}
	if 0 == ret.Data.Count {
		blog.Errorf("getModuleNameByID moduleID %v supplierID %s  not found module info,query:%+v,rid:%s", ID, supplierID, dat, lgc.rid)
		return "", 0, 0, nil
	}
	byteModule, _ := json.Marshal(ret.Data.Info[0])
	name, err = ret.Data.Info[0].String(common.BKModuleNameField)
	if nil != err {
		blog.Warnf("getModuleNameByID moduleID %v supplierID %s  get module name error:%s, raw:%s,query:%+v,rid:%s", ID, supplierID, err.Error(), string(byteModule), dat, lgc.rid)
		return
	}
	appID, err = ret.Data.Info[0].Int64(common.BKAppIDField)
	if nil != err {
		blog.Errorf("getModuleNameByID moduleID %v supplierID %s  get appID name error:%s raw:%s,query:%+v,rid:%s", ID, supplierID, err.Error(), string(byteModule), dat, lgc.rid)
		return "", 0, 0, err
	}
	setID, err = ret.Data.Info[0].Int64(common.BKSetIDField)
	if nil != err {
		blog.Errorf("getModuleNameByID moduleID %v supplierID %s  get set name error:%s raw:%s,query:%+v,rid:%s", ID, supplierID, err.Error(), string(byteModule), dat, lgc.rid)
		return "", 0, 0, err
	}

	return

}

func (lgc *Logics) getModuleIDByProcID(ctx context.Context, appID, procID int64) ([]int64, error) {
	condition := make(map[string]interface{}, 0)
	condition[common.BKProcessIDField] = procID
	defErr := lgc.ccErr
	// get process by module
	ret, err := lgc.CoreAPI.ProcController().GetProc2Module(ctx, lgc.header, condition)
	if nil != err {
		blog.Errorf("getModuleIDByProcID  GetProc2Module appID %d moduleID %s  http do error:%s,query:%+v,rid:%s", appID, procID, err.Error(), condition, lgc.rid)
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !ret.Result {
		blog.Errorf("getModuleIDByProcID  GetProc2Module appID %d moduleID %s  http reply error:%s,query:%+v,rid:%s", appID, procID, ret.ErrMsg, condition, lgc.rid)
		return nil, defErr.New(ret.Code, ret.ErrMsg)
	}
	var moduleIDs []int64
	for _, item := range ret.Data {
		ids, err := lgc.HandleProcInstNumByModuleName(ctx, appID, item.ModuleName)
		if nil != err {
			blog.Errorf("getModuleIDByProcID get module id by module name %s  in application id  %d error %s,rid:%s", item.ModuleName, item.ApplicationID, err.Error(), lgc.rid)
			return nil, err
		}
		moduleIDs = append(moduleIDs, ids...)
	}

	return moduleIDs, nil
}

func (lgc *Logics) GetModuleIDByHostID(ctx context.Context, hostID int64) ([]metadata.ModuleHost, error) {
	defErr := lgc.ccErr
	dat := map[string][]int64{
		common.BKHostIDField: []int64{hostID},
	}
	ret, err := lgc.CoreAPI.HostController().Module().GetModulesHostConfig(ctx, lgc.header, dat)
	if nil != err {
		blog.Errorf("GetModuleIDByHostID appID %d module id %d GetModulesHostConfig http do error:%s,query:%+v,rid:%s", hostID, err.Error(), dat, lgc.rid)
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !ret.Result {
		blog.Errorf("GetModuleIDByHostID appID %d module id %d GetModulesHostConfig reply error:%s,query:%+v,rid:%s", hostID, ret.ErrMsg, dat, lgc.rid)
		return nil, defErr.New(ret.Code, ret.ErrMsg)
	}

	return ret.Data, nil
}

func (lgc *Logics) GetModueleIDByAppID(ctx context.Context, appID int64) ([]int64, error) {
	supplierID := lgc.ownerID
	defErr := lgc.ccErr
	dat := new(metadata.QueryCondition)
	dat.Condition = mapstr.MapStr{common.BKAppIDField: appID}
	dat.Fields = []string{common.BKModuleNameField, common.BKAppIDField, common.BKSetIDField, common.BKModuleIDField}
	dat.Limit.Limit = common.BKNoLimit
	ret, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(ctx, lgc.header, common.BKInnerObjIDModule, dat)
	if nil != err {
		blog.Errorf("GetModueleIDByAppID appID %v supplierID %s  http do error:%s,query:%+v,rid:%s", appID, supplierID, err.Error(), dat, lgc.rid)
		return make([]int64, 0), defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !ret.Result {
		blog.Errorf("GetModueleIDByAppID appID %v supplierID %s  http reply error:%s,query:%+v,rid:%s", appID, supplierID, ret.ErrMsg, dat, lgc.rid)
		return make([]int64, 0), defErr.New(ret.Code, ret.ErrMsg)

	}
	if 0 == ret.Data.Count {
		blog.V(5).Infof("GetModueleIDByAppID appID %v supplierID %s  not found module info,query:%+v,rid:%s", appID, supplierID, dat, lgc.rid)
		return make([]int64, 0), nil
	}
	moduleIDs := make([]int64, 0)
	for _, module := range ret.Data.Info {
		moduleID, err := module.Int64(common.BKModuleIDField)
		if nil != err {
			byteModule, _ := json.Marshal(module)
			blog.Errorf("GetModueleIDByAppID moduleID %v supplierID %s  get set name error:%s raw:%s,query:%+v,rid:%s", appID, supplierID, err.Error(), string(byteModule), dat, lgc.rid)
			return nil, err
		}
		moduleIDs = append(moduleIDs, moduleID)
	}

	return moduleIDs, nil
}

func (lgc *Logics) GetAppList(ctx context.Context, fields []string) ([]mapstr.MapStr, error) {
	defErr := lgc.ccErr
	dat := new(metadata.QueryCondition)
	dat.Fields = fields
	dat.Limit.Limit = common.BKNoLimit
	ret, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(ctx, lgc.header, common.BKInnerObjIDApp, dat)
	if nil != err {
		blog.Errorf("GetAppList  http do error:%s,query:%+v,rid:%s", err.Error(), dat, lgc.rid)
		return make([]mapstr.MapStr, 0), defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !ret.Result {
		blog.Errorf("GetAppList  http reply error,err code:%d,err msg:%s,query:%+v,rid:%s", ret.Code, ret.ErrMsg, dat, lgc.rid)
		return make([]mapstr.MapStr, 0), defErr.New(ret.Code, ret.ErrMsg)

	}
	if 0 == ret.Data.Count {
		blog.V(5).Infof("GetAppList  not found app info,query:%+v,rid:%s", dat, lgc.rid)
		return make([]mapstr.MapStr, 0), nil
	}

	return ret.Data.Info, nil
}
