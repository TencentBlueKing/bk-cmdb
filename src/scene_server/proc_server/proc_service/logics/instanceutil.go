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
	"sync"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

type chanItem struct {
	ctx       context.Context
	eventData *metadata.EventInst
	opFunc    func(ctx context.Context, eventData *metadata.EventInst) error
	retry     int
}

type refreshModuleData struct {
	header   http.Header
	appID    int64
	moduleID int64
}
type eventRefreshModule struct {
	data     map[string]*refreshModuleData
	eventChn chan bool
	sync.RWMutex
}

var (
	handEventDataChan      chan chanItem
	chnOpLock              *sync.Once = new(sync.Once)
	initDataLock           *sync.Once = new(sync.Once)
	eventRefreshModuleData            = &eventRefreshModule{}
	maxRefreshModuleData   int        = 10
	maxEventDataChan       int        = 10000
	retry                  int        = 3
)

func reshReshInitChan(maxEvent, maxRefresh int) {
	if 0 != maxEvent {
		maxEventDataChan = maxEvent
	}
	if 0 != maxRefresh {
		maxRefreshModuleData = maxRefresh
	}
	handEventDataChan = make(chan chanItem, maxEventDataChan)
	eventRefreshModuleData.data = make(map[string]*refreshModuleData, 0)
	eventRefreshModuleData.eventChn = make(chan bool, maxRefreshModuleData)
}

func getEventRefrshModuleKey(appID, moduleID int64) string {
	return fmt.Sprintf("%d-%d", appID, moduleID)
}

func addEventRefreshModuleItem(appID, moduleID int64, header http.Header) {
	defer eventRefreshModuleData.Unlock()
	eventRefreshModuleData.Lock()
	eventRefreshModuleData.data[getEventRefrshModuleKey(appID, moduleID)] = &refreshModuleData{appID: appID, moduleID: moduleID, header: header}
}

func addEventRefreshModuleItems(appID int64, moduleIDs []int64, header http.Header) {
	defer eventRefreshModuleData.Unlock()
	eventRefreshModuleData.Lock()
	for _, moduleID := range moduleIDs {
		eventRefreshModuleData.data[getEventRefrshModuleKey(appID, moduleID)] = &refreshModuleData{appID: appID, moduleID: moduleID, header: header}
	}
}

func getEventRefreshModuleItem() *refreshModuleData {
	defer eventRefreshModuleData.Unlock()
	eventRefreshModuleData.Lock()
	for key, item := range eventRefreshModuleData.data {
		delete(eventRefreshModuleData.data, key)
		return item
	}
	return nil
}

func sendEventFrefreshModuleNotice() {
	if maxRefreshModuleData > len(eventRefreshModuleData.eventChn) {
		go func() { eventRefreshModuleData.eventChn <- true }()
	}
}

func (lgc *Logics) unregisterProcInstDetall(ctx context.Context, header http.Header, appID, moduleID int64, unregister []metadata.ProcInstanceModel) error {
	if 0 == len(unregister) {
		return nil
	}
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	dat := new(metadata.QueryInput)
	dat.Condition = map[string]interface{}{common.BKDBOR: unregister}
	dat.Limit = 200
	dat.Start = 0
	for {
		ret, err := lgc.CoreAPI.ProcController().GetProcInstanceDetail(ctx, header, dat)
		if nil != err {
			blog.Errorf("unregisterProcInstDetall  proc instance error:%s", err.Error())
			return defErr.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		if !ret.Result {
			blog.Errorf("unregisterProcInstDetall  proc instance return err msg %s", ret.ErrMsg)
			return defErr.New(ret.Code, ret.ErrMsg)
		}
		if 0 == ret.Data.Count {
			return nil
		}
		for _, item := range ret.Data.Info {
			gseProc := new(metadata.GseProcRequest)
			gseProc.OpType = common.GSEProcOPUnregister
			gseProc.AppID = item.AppID
			gseProc.ProcID = item.ProcID
			gseProc.Hosts = item.Hosts
			gseProc.Meta = item.Meta
			gseProc.ModuleID = item.ModuleID
			gseProc.Spec = item.Spec
			err := lgc.unregisterProcInstanceToGse(gseProc, header)
			if nil != err {
				blog.Errorf("unregisterProcInstanceToGse  proc instance error:%s, gse proc %+v", err.Error(), gseProc)
				return err
			}
			if !ret.Result {
				blog.Errorf("unregisterProcInstanceToGse  proc instance return err msg %s, gse proc %+v", ret.ErrMsg, gseProc)
				return err
			}
		}
	}

	return nil

}

// setProcInstDetallStatusUnregister modify process instance status to unregister in cmdb table
func (lgc *Logics) setProcInstDetallStatusUnregister(ctx context.Context, header http.Header, appID, moduleID int64, unregister []metadata.ProcInstanceModel) error {

	if 0 != len(unregister) {
		return nil
	}
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	unregisterProcDetail := make([]interface{}, 0)
	for _, item := range unregister {
		unregisterProcDetail = append(unregisterProcDetail, common.KvMap{common.BKAppIDField: item.ApplicationID, common.BKModuleIDField: item.ModuleID, common.BKHostIDField: item.HostID, common.BKProcIDField: item.ProcID})
	}
	dat := new(metadata.ModifyProcInstanceDetail)
	dat.Conds = map[string]interface{}{common.BKDBOR: unregisterProcDetail}
	dat.Data = map[string]interface{}{"status": metadata.ProcInstanceDetailStatusUnRegisterFailed}
	ret, err := lgc.CoreAPI.ProcController().ModifyProcInstanceDetail(ctx, header, dat)
	if nil != err {
		blog.Errorf("setProcInstDetallStatusUnregister  proc instance error:%s", err.Error())
		return defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !ret.Result {
		blog.Errorf("setProcInstDetallStatusUnregister  proc instance return err msg %s", ret.ErrMsg)
		return defErr.New(ret.Code, ret.ErrMsg)
	}

	return nil
}

func (lgc *Logics) handleProcInstNumDataHandle(ctx context.Context, header http.Header, appID, moduleID int64, instProc []*metadata.ProcInstanceModel) error {
	delConds := common.KvMap{common.BKAppIDField: appID, common.BKModuleIDField: moduleID}
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	ret, err := lgc.CoreAPI.ProcController().DeleteProcInstanceModel(ctx, header, delConds)
	if nil != err {
		blog.Errorf("handleInstanceNum create proc instance error:%s", err.Error())
		return defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !ret.Result {
		blog.Errorf("handleInstanceNum create proc instance return err msg %s", ret.ErrMsg)
		return defErr.New(ret.Code, ret.ErrMsg)
	}
	if 0 < len(instProc) {
		ret, err := lgc.CoreAPI.ProcController().CreateProcInstanceModel(ctx, header, instProc)
		if nil != err {
			blog.Errorf("handleInstanceNum create proc instance error:%s", err.Error())
			return defErr.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		if !ret.Result {
			blog.Errorf("handleInstanceNum create proc instance return err msg %s", ret.ErrMsg)
			return defErr.New(ret.Code, ret.ErrMsg)
		}
	}

	return nil
}

func (lgc *Logics) HandleProcInstNumByModuleName(ctx context.Context, header http.Header, appID int64, moduleName string) ([]int64, error) {
	moduleConds := new(metadata.SearchParams)
	moduleConds.Condition = map[string]interface{}{
		common.BKAppIDField:      appID,
		common.BKModuleNameField: moduleName,
	}
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	moduleConds.Page = map[string]interface{}{"start": 0, "limit": common.BKNoLimit}
	moduleInfos, err := lgc.CoreAPI.TopoServer().Instance().InstSearch(ctx, util.GetOwnerID(header), common.BKInnerObjIDModule, header, moduleConds)
	if nil != err {
		blog.Errorf("HandleProcInstNumByModuleName get module by name %s error %s", moduleName, err.Error())
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !moduleInfos.Result {
		blog.Errorf("HandleProcInstNumByModuleName get module by name %s error %s ", moduleName, moduleInfos.ErrMsg)
		return nil, defErr.New(moduleInfos.Code, moduleInfos.ErrMsg)
	}
	moduleIDs := make([]int64, 0)
	for _, module := range moduleInfos.Data.Info {
		moduleID, err := module.Int64(common.BKModuleIDField)
		if nil != err {
			byteData, _ := json.Marshal(module)
			blog.Errorf("refreshProcInstByProcModule get module by name %s not module id error  item %v ", moduleName, string(byteData))
			return nil, err
		}
		moduleIDs = append(moduleIDs, moduleID)
	}

	return moduleIDs, nil
}

func (lgc *Logics) getProcInstInfoByModuleID(ctx context.Context, appID, moduleID int64, header http.Header) (maxHostInstID uint64, procInst map[string]metadata.ProcInstanceModel, err error) {
	dat := new(metadata.QueryInput)
	dat.Condition = common.KvMap{common.BKModuleIDField: moduleID, common.BKAppIDField: appID}
	dat.Limit = common.BKNoLimit
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	ret, err := lgc.CoreAPI.ProcController().GetProcInstanceModel(ctx, header, dat)
	if nil != err {
		blog.Errorf("getProcInstInfoByModuleID http error, error:%s", err.Error())
		return 0, nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !ret.Result {
		blog.Errorf("getProcInstInfoByModuleID  reply error error:%s", ret.ErrMsg)
		return 0, nil, defErr.New(ret.Code, ret.ErrMsg)
	}
	procInst = make(map[string]metadata.ProcInstanceModel, 0)
	for _, item := range ret.Data.Info {
		if maxHostInstID < item.HostInstanID {
			maxHostInstID = item.HostInstanID
		}
		procInst[getInlineProcInstKey(int64(item.HostID), int64(item.ProcID))] = item
	}

	return maxHostInstID, procInst, nil

}

func (lgc *Logics) getModuleBindProc(ctx context.Context, appID, moduleID int64, header http.Header) (setID int64, procID []int64, err error) {
	supplierID := util.GetOwnerID(header)
	var name string
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	name, appID, setID, err = lgc.getModuleNameByID(ctx, moduleID, header)
	if nil != err {
		blog.Errorf("getModuleBindProc error:%s", err.Error())
		return 0, nil, err
	}
	dat := common.KvMap{common.BKModuleNameField: name, common.BKAppIDField: appID}
	ret, err := lgc.CoreAPI.ProcController().GetProc2Module(ctx, header, dat)
	if nil != err {
		blog.Errorf("getModuleBindProc moduleID %d supplierID %s  http do error:%s", moduleID, supplierID, err.Error())
		return 0, nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !ret.Result {
		blog.Errorf("getModuleBindProc moduleID %d supplierID %s  http reply error:%s", moduleID, supplierID, ret.ErrMsg)
		return 0, nil, defErr.New(ret.Code, ret.ErrMsg)
	}
	for _, proc := range ret.Data {
		procID = append(procID, proc.ProcessID)
	}

	return setID, procID, nil
}

func (lgc *Logics) getModuleIDByProcID(ctx context.Context, appID, procID int64, header http.Header) ([]int64, error) {
	condition := make(map[string]interface{}, 0)
	condition[common.BKProcIDField] = procID
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	// get process by module
	ret, err := lgc.CoreAPI.ProcController().GetProc2Module(context.Background(), header, condition)
	if nil != err {
		blog.Errorf("getModuleIDByProcID  GetProc2Module appID %d moduleID %s  http do error:%s", appID, procID, err.Error())
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !ret.Result {
		blog.Errorf("getModuleIDByProcID  GetProc2Module appID %d moduleID %s  http reply error:%s", appID, procID, ret.ErrMsg)
		return nil, defErr.New(ret.Code, ret.ErrMsg)
	}
	var moduleIDs []int64
	for _, item := range ret.Data {
		ids, err := lgc.HandleProcInstNumByModuleName(ctx, header, appID, item.ModuleName)
		if nil != err {
			blog.Errorf("getModuleIDByProcID get module id by module name %s  in application id  %d error %s", item.ModuleName, item.ApplicationID, err.Error())
			return nil, err
		}
		moduleIDs = append(moduleIDs, ids...)
	}

	return moduleIDs, nil
}

func (lgc *Logics) getHostByModuleID(ctx context.Context, header http.Header, moduleID int64) (map[int64]*metadata.GseHost, error) {
	dat := map[string][]int64{
		common.BKModuleIDField: []int64{moduleID},
	}
	supplierID := util.GetOwnerID(header)
	intSupplierID, err := util.GetInt64ByInterface(supplierID)
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	if nil != err {
		blog.Errorf("getHostByModuleID supplierID %s  not interger", supplierID)
		return nil, err
	}

	ret, err := lgc.CoreAPI.HostController().Module().GetModulesHostConfig(ctx, header, dat)
	if nil != err {
		blog.Errorf("getHostByModuleID moduleID %d supplierID %s GetModulesHostConfig http do error:%s", moduleID, supplierID, err.Error())
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !ret.Result {
		blog.Errorf("getHostByModuleID moduleID %d supplierID %s GetModulesHostConfig http reply error:%s", moduleID, supplierID, ret.ErrMsg)
		return nil, defErr.New(ret.Code, ret.ErrMsg)
	}
	if 0 == len(ret.Data) {
		blog.V(5).Infof("getHostByModuleID moduleID %d supplierID %s GetModulesHostConfig len equal 0", moduleID, supplierID)
		return nil, nil
	}
	var hostIDs []int64
	for _, item := range ret.Data {
		hostIDs = append(hostIDs, item.HostID)
	}
	opt := new(metadata.QueryInput)
	opt.Condition = common.KvMap{common.BKHostIDField: common.KvMap{common.BKDBIN: hostIDs}}
	opt.Fields = fmt.Sprintf("%s,%s,%s", common.BKHostIDField, common.BKHostInnerIPField, common.BKCloudIDField)
	opt.Limit = common.BKNoLimit
	hosts, err := lgc.CoreAPI.HostController().Host().GetHosts(ctx, header, opt)
	if nil != err {
		blog.Errorf("getHostByModuleID moduleID %d hostID:%v supplierID %s GetHosts http do error:%s", moduleID, hostIDs, supplierID, err.Error())
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !hosts.Result {
		blog.Errorf("getHostByModuleID moduleID %d hostID:%v supplierID %s GetHosts http reply error:%s", moduleID, hostIDs, supplierID, hosts.ErrMsg)
		return nil, defErr.New(hosts.Code, hosts.ErrMsg)
	}

	hostInfos := make(map[int64]*metadata.GseHost, len(hosts.Data.Info))
	for _, host := range hosts.Data.Info {
		item := new(metadata.GseHost)

		hostID, err := util.GetInt64ByInterface(host[common.BKHostIDField])
		if nil != err {
			blog.Errorf("getHostByModuleID hostInfo %v  hostID   not interger", host)
			return nil, err
		}
		cloudID, err := util.GetInt64ByInterface(host[common.BKCloudIDField])
		if nil != err {
			byteHost, _ := json.Marshal(host)
			blog.Errorf("getHostByModuleID  hostInfo %v  cloudID  not interger, json:%s", host, string(byteHost))
			return nil, err
		}
		innerIP, ok := host[common.BKHostInnerIPField].(string)
		if !ok {
			byteHost, _ := json.Marshal(host)
			blog.Errorf("getHostByModuleID  hostInfo %v  innerip  not found, json:%s", host, string(byteHost))
			return nil, err
		}
		item.BkCloudId = cloudID
		item.Ip = innerIP
		item.BkSupplierId = intSupplierID

		hostInfos[hostID] = item
	}

	return hostInfos, nil
}

func (lgc *Logics) getProcInfoByID(ctx context.Context, procID []int64, header http.Header) (map[int64]*metadata.InlineProcInfo, error) {
	supplierID := util.GetOwnerID(header)
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	if 0 == len(procID) {
		return nil, nil
	}
	gseProc := make(map[int64]*metadata.InlineProcInfo, 0)
	dat := new(metadata.QueryInput)
	dat.Condition = common.KvMap{common.BKProcIDField: common.KvMap{common.BKDBIN: procID}}
	dat.Limit = common.BKNoLimit
	ret, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(ctx, common.BKInnerObjIDProc, header, dat)
	if nil != err {
		blog.Errorf("getProcInfoByID procID %v supplierID %s  http do error:%s", procID, supplierID, err.Error())
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !ret.Result {
		blog.Errorf("getProcInfoByID procID %v supplierID %s  http reply error:%s", procID, supplierID, ret.ErrMsg)
		return nil, defErr.New(ret.Code, ret.ErrMsg)

	}
	if 0 == ret.Data.Count {
		blog.Errorf("getProcInfoByID procID %v supplierID %s  not found process info", procID, supplierID)
		return nil, nil
	}
	for _, proc := range ret.Data.Info {
		procID, err := proc.Int64(common.BKProcIDField)
		if nil != err {
			byteHost, _ := json.Marshal(proc)
			blog.Errorf("getHostByModuleID  proc %v  procID  not interger, json:%s", proc, string(byteHost))
			return nil, err
		}
		item := new(metadata.InlineProcInfo)

		item.ProcNum = 1
		procNumI, ok := proc.Get(common.BKProcInstNum)
		if ok && nil != procNumI {
			item.ProcNum, err = proc.Int64(common.BKProcInstNum)
			if nil != err {
				byteHost, _ := json.Marshal(proc)
				blog.Errorf("getHostByModuleID  proc %v  procNum  not interger, json:%s", proc, string(byteHost))
				return nil, err
			}
		}
		item.AppID, err = proc.Int64(common.BKAppIDField)
		if nil != err {
			byteHost, _ := json.Marshal(proc)
			blog.Errorf("getHostByModuleID  proc %v  AppID  not interger, json:%s", proc, string(byteHost))
			return nil, err
		}
		item.FunID, err = proc.Int64(common.BKFuncIDField)
		if nil != err {
			byteHost, _ := json.Marshal(proc)
			blog.Errorf("getHostByModuleID  proc %v  AppID  not interger, json:%s", proc, string(byteHost))
			return nil, err
		}

		gseProc[procID] = item
	}

	return gseProc, nil
}

func (lgc *Logics) getModuleNameByID(ctx context.Context, ID int64, header http.Header) (name string, appID int64, setID int64, err error) {
	supplierID := util.GetOwnerID(header)
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	dat := new(metadata.QueryInput)
	dat.Condition = common.KvMap{common.BKModuleIDField: ID}
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
