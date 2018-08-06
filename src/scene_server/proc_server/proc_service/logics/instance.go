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

func (lgc *Logics) HandleHostProcDataChange(ctx context.Context, eventData *metadata.EventInst) {

	switch eventData.ObjType {
	case metadata.EventObjTypeProcModule:
		lgc.handleRetry(3, ctx, eventData, lgc.refreshProcInstByProcModule)
		//lgc.refreshProcInstByProcModule(ctx, eventData)
	case metadata.EventObjTypeModuleTransfer:
		lgc.handleRetry(3, ctx, eventData, lgc.eventHostModuleChangeProcHostInstNum)
		//lgc.eventHostModuleChangeProcHostInstNum(ctx, eventData)
	}
}

func (lgc *Logics) handleRetry(maxRetry int, ctx context.Context, eventData *metadata.EventInst, opFunc func(ctx context.Context, eventData *metadata.EventInst) error) error {
	var err error
	for idx := 0; idx < maxRetry; idx++ {
		err = opFunc(ctx, eventData)
		if nil == err {
			return nil
		}
	}
	return err
}

func (lgc *Logics) eventHostModuleChangeProcHostInstNum(ctx context.Context, eventData *metadata.EventInst) error {
	var header http.Header
	header.Set(common.BKHTTPOwnerID, eventData.OwnerID)
	header.Set(common.BKHTTPHeaderUser, common.BKProcInstanceOpUser)
	for _, hostInfos := range eventData.Data {
		var data interface{}
		if metadata.EventActionDelete == eventData.Action {
			data = hostInfos.PreData
		} else {
			data = hostInfos.CurData
		}
		mapData, err := mapstr.NewFromInterface(data)
		if nil != err {
			byteData, _ := json.Marshal(eventData)
			blog.Errorf("productHostInstanceNum event data not map[string]interface{} item %v raw josn %s", hostInfos, string(byteData))
			return err
		}
		appID, err := mapData.Int64(common.BKAppIDField)
		if nil != err {
			byteData, _ := json.Marshal(eventData)
			blog.Errorf("productHostInstanceNum event data  appID not integer  item %v raw josn %s", hostInfos, string(byteData))
			return err
		}
		moduleID, err := mapData.Int64(common.BKModuleIDField)
		if nil != err {
			byteData, _ := json.Marshal(eventData)
			blog.Errorf("productHostInstanceNum event data  moduleID not integer  item %v raw josn %s", hostInfos, string(byteData))
			return err
		}
		lgc.HandleProcInstNumByModuleID(ctx, header, appID, moduleID)

	}

	return nil
}

func (lgc *Logics) HandleProcInstNumByModuleID(ctx context.Context, header http.Header, appID, moduleID int64) error {
	maxInstID, procInst, err := lgc.getProcInstInfoByModuleID(ctx, appID, moduleID, header)
	if nil != err {
		return err
	}
	var hostInfos map[int64]*metadata.GseHost
	/*if isHost {
		//hostIDs = append(hostIDs, hostID)
		hostInfo, err := lgc.getHostByHostID(ctx, header, hostID)
		if nil != err {
			blog.Errorf("handleInstanceNum getHostByHostID error %s", err.Error())
			return err
		}
		hostInfos[hostID] = hostInfo
	} else {*/
	// get hostid by module
	hostInfos, err = lgc.getHostByModuleID(ctx, header, moduleID)
	if nil != err {
		blog.Errorf("handleInstanceNum getHostByModuleID error %s", err.Error())
		return err
	}
	setID, procIDs, err := lgc.getModuleBindProc(ctx, appID, moduleID, header)
	if nil != err {
		return err
	}
	instProc := make([]*metadata.ProcInstanceModel, 0)
	procInfos, err := lgc.getProcInfoByID(ctx, procIDs, header)
	for procID, info := range procInfos {
		for hostID, _ := range hostInfos {
			_, ok := procInst[getInlineProcInstKey(hostID, procID)]
			if !ok {
				maxInstID += 1
			}
			for numIdx := int64(1); numIdx < info.ProcNum+1; numIdx++ {
				procIdx := (maxInstID-1)*uint64(info.ProcNum) + uint64(numIdx)
				item := new(metadata.ProcInstanceModel)
				item.ApplicationID = appID
				item.SetID = setID
				item.ModuleID = moduleID
				item.FuncID = info.FunID
				item.HostID = hostID
				item.HostInstanID = maxInstID
				item.ProcInstanceID = procIdx
				item.ProcID = procID
				item.HostProcID = uint64(numIdx)
				instProc = append(instProc, item)
			}
		}
	}
	delConds := common.KvMap{common.BKAppIDField: appID, common.BKModuleIDField: moduleID, common.BKProcIDField: procIDs}
	/*if isHost {
		delConds[common.BKHostIDField] = hostID
	}*/
	ret, err := lgc.CoreAPI.ProcController().DeleteProcInstanceModel(ctx, header, delConds)
	if nil != err {
		blog.Errorf("handleInstanceNum create proc instance error:%s", err.Error())
		return fmt.Errorf("delete proc instance error:%s", err.Error())
	}
	if !ret.Result {
		blog.Errorf("handleInstanceNum create proc instance return err msg %s", ret.ErrMsg)
		return fmt.Errorf("delete proc instance return err msg %s", ret.ErrMsg)
	}
	if 0 < len(instProc) {
		ret, err := lgc.CoreAPI.ProcController().CreateProcInstanceModel(ctx, header, instProc)
		if nil != err {
			blog.Errorf("handleInstanceNum create proc instance error:%s", err.Error())
			return fmt.Errorf("create proc instance error:%s", err.Error())
		}
		if !ret.Result {
			blog.Errorf("handleInstanceNum create proc instance return err msg %s", ret.ErrMsg)
			return fmt.Errorf("create proc instance return err msg %s", ret.ErrMsg)
		}

	}

	return nil
}

func (lgc *Logics) HandleProcInstNumByModuleName(ctx context.Context, header http.Header, appID int64, moduleName string) error {
	moduleConds := new(metadata.SearchParams)
	moduleConds.Condition = map[string]interface{}{
		common.BKAppIDField: appID,
		moduleName:          moduleName,
	}
	moduleConds.Page = map[string]interface{}{"start": 0, "limit": 1}
	moduleInfos, err := lgc.CoreAPI.TopoServer().Instance().InstSearch(ctx, util.GetOwnerID(header), common.BKInnerObjIDModule, header, moduleConds)
	if nil != err {
		blog.Errorf("HandleProcInstNumByModuleName get module by name %s error %s", moduleName, err.Error())
		return err
	}
	if !moduleInfos.Result {
		blog.Errorf("HandleProcInstNumByModuleName get module by name %s error %s ", moduleName, moduleInfos.ErrMsg)
		return fmt.Errorf("%s", moduleInfos.ErrMsg)
	}
	for _, module := range moduleInfos.Data.Info {
		moduleID, err := module.Int64(common.BKModuleIDField)
		if nil != err {
			byteData, _ := json.Marshal(module)
			blog.Errorf("refreshProcInstByProcModule get module by name %s not module id error  item %v ", moduleName, string(byteData))
			return err
		}
		err = lgc.HandleProcInstNumByModuleID(ctx, header, appID, moduleID)
		if nil != err {
			byteData, _ := json.Marshal(module)
			blog.Errorf("refreshProcInstByProcModule get module by name %s not module id error  item %v ", moduleName, string(byteData))
			return err
		}
	}
	return nil
}

func getInlineProcInstKey(hostID, procID int64) string {
	return fmt.Sprintf("%d-%d", hostID, procID)
}

func (lgc *Logics) getProcInstInfoByModuleID(ctx context.Context, appID, moduleID int64, header http.Header) (maxHostInstID uint64, procInst map[string]metadata.ProcInstanceModel, err error) {
	dat := new(metadata.QueryInput)
	dat.Condition = common.KvMap{common.BKModuleIDField: moduleID, common.BKAppIDField: appID}
	dat.Limit = common.BKNoLimit

	ret, err := lgc.CoreAPI.ProcController().GetProcInstanceModel(ctx, header, dat)
	if nil != err {
		blog.Errorf("getProcInstInfoByModuleID http error, error:%s", err.Error())
		return 0, nil, err
	}
	if !ret.Result {
		blog.Errorf("getProcInstInfoByModuleID  reply error error:%s", ret.ErrMsg)
		return 0, nil, fmt.Errorf(ret.ErrMsg)
	}
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
	name, appID, setID, err = lgc.getModuleNameByID(ctx, moduleID, header)
	if nil != err {
		blog.Errorf("getModuleBindProc error:%s", err.Error())
		return 0, nil, err
	}
	dat := common.KvMap{common.BKModuleNameField: name, common.BKAppIDField: appID}
	ret, err := lgc.CoreAPI.ProcController().GetProc2Module(ctx, header, dat)
	if nil != err {
		blog.Errorf("getModuleBindProc moduleID %d supplierID %s  http do error:%s", moduleID, supplierID, err.Error())
		return 0, nil, err
	}
	if !ret.Result {
		blog.Errorf("getModuleBindProc moduleID %d supplierID %s  http reply error:%s", moduleID, supplierID, ret.ErrMsg)
		return 0, nil, fmt.Errorf(ret.ErrMsg)
	}
	for _, proc := range ret.Data {
		procID = append(procID, proc.ProcessID)
	}

	return setID, procID, nil
}

///

func (lgc *Logics) refreshProcInstByProcModule(ctx context.Context, eventData *metadata.EventInst) error {
	if metadata.EventTypeRelation != eventData.EventType {
		return nil
	}
	var header http.Header
	header.Set(common.BKHTTPOwnerID, eventData.OwnerID)
	header.Set(common.BKHTTPHeaderUser, common.BKProcInstanceOpUser)

	for _, data := range eventData.Data {
		var iData interface{}
		if metadata.EventActionDelete == eventData.Action {
			// delete process bind module relation, unregister process info
			iData = data.PreData
		} else {
			// compare  pre-change data with the current data and find the newly added data to register process info
			iData = data.CurData
		}
		mapData, err := mapstr.NewFromInterface(iData)
		if nil != err {
			byteData, _ := json.Marshal(eventData)
			blog.Errorf("refreshProcInstByProcModule event data not map[string]interface{} item %v raw josn %s", data, string(byteData))
			return err
		}
		appID, err := mapData.Int64(common.BKAppIDField)
		if nil != err {
			byteData, _ := json.Marshal(eventData)
			blog.Errorf("refreshProcInstByProcModule event data  appID not integer  item %v raw josn %s", data, string(byteData))
			return err
		}
		moduleName, err := mapData.String(common.BKModuleNameField)
		if nil != err {
			byteData, _ := json.Marshal(eventData)
			blog.Errorf("refreshProcInstByProcModule event data  appID not integer  item %v raw josn %s", data, string(byteData))
			return err
		}
		err = lgc.HandleProcInstNumByModuleName(ctx, header, appID, moduleName)
		if nil != err {
			byteData, _ := json.Marshal(eventData)
			blog.Errorf("refreshProcInstByProcModule HandleProcInstNumByModuleName error %s item %v raw josn %s", err.Error(), data, string(byteData))
			return err
		}
	}
	return nil
}

func (lgc *Logics) refreshProcInstByProcess(ctx context.Context, eventData *metadata.EventInst) {
	switch eventData.Action {
	case metadata.EventActionCreate:
		// create proccess not refresh process instance , because not bind module
	case metadata.EventActionUpdate:
		// refresh process instance register again
	case metadata.EventActionDelete:
		// delete all register process
	}
}

func (lgc *Logics) refreshProcInstByHostInfo(ctx context.Context, eventData *metadata.EventInst) {
	if metadata.EventTypeRelation == eventData.EventType {
		if metadata.EventActionDelete == eventData.Action {
			// delete host from module , unregister module bind all process info
		} else {
			// compare pre-change data with the current data and find the newly added data to register process info
		}
	} else {
		// host fields supperid, cloud id, innerip   change , register process info  agin
	}
}

func (lgc *Logics) refreshProcInstModuleByHostInfo(ctx context.Context, eventData *metadata.EventInst) {
	if metadata.EventActionUpdate == eventData.EventType {
		// module change name, unregister pre-change module name bind process , then register current module name bind process
	} else {

	}
}

func (lgc *Logics) operateProcInstByHost(ctx context.Context, hostID int64, hostInfo metadata.GseHost) error {
	return nil
}

func (lgc *Logics) refreshProcInstModuleWithModule(ctx context.Context, moduleID int64, supplierID string, isregister bool) error {

	var header http.Header
	header.Set(common.BKHTTPOwnerID, supplierID)
	header.Set(common.BKHTTPHeaderUser, common.BKProcInstanceOpUser)

	return nil
}

func (lgc *Logics) getHostByModuleID(ctx context.Context, header http.Header, moduleID int64) (map[int64]*metadata.GseHost, error) {
	dat := map[string][]int64{
		common.BKModuleIDField: []int64{moduleID},
	}
	supplierID := util.GetOwnerID(header)
	intSupplierID, err := util.GetInt64ByInterface(supplierID)
	if nil != err {
		blog.Errorf("getHostByModuleID supplierID %s  not interger", supplierID)
		return nil, err
	}

	ret, err := lgc.CoreAPI.HostController().Module().GetModulesHostConfig(ctx, header, dat)
	if nil != err {
		blog.Errorf("getHostByModuleID moduleID %d supplierID %s GetModulesHostConfig http do error:%s", moduleID, supplierID, err.Error())
		return nil, err
	}
	if !ret.Result {
		blog.Errorf("getHostByModuleID moduleID %d supplierID %s GetModulesHostConfig http reply error:%s", moduleID, supplierID, ret.ErrMsg)
		return nil, fmt.Errorf(ret.ErrMsg)
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
		return nil, err
	}
	if !hosts.Result {
		blog.Errorf("getHostByModuleID moduleID %d hostID:%v supplierID %s GetHosts http reply error:%s", moduleID, hostIDs, supplierID, hosts.ErrMsg)
		return nil, fmt.Errorf(hosts.ErrMsg)
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

	gseProc := make(map[int64]*metadata.InlineProcInfo, 0)

	dat := new(metadata.QueryInput)
	dat.Condition = common.KvMap{common.BKProcIDField: common.KvMap{common.BKDBIN: procID}}
	dat.Limit = common.BKNoLimit
	ret, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(ctx, common.BKInnerObjIDProc, header, dat)
	if nil != err {
		blog.Errorf("getProcInfoByID procID %v supplierID %s  http do error:%s", procID, supplierID, err.Error())
		return nil, err
	}
	if !ret.Result {
		blog.Errorf("getProcInfoByID procID %v supplierID %s  http reply error:%s", procID, supplierID, ret.ErrMsg)
		return nil, fmt.Errorf(ret.ErrMsg)

	}
	if 0 == ret.Data.Count {
		blog.Errorf("getProcInfoByID procID %v supplierID %s  not found process info", procID, supplierID)
		return nil, nil
	}
	for _, proc := range ret.Data.Info {
		procID, err := proc.Int64(common.BKProcField)
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

func (lgc *Logics) getSetNameBySetID(ctx context.Context, setID int64, header http.Header) (string, error) {
	supplierID := util.GetOwnerID(header)

	dat := new(metadata.QueryInput)
	dat.Condition = common.KvMap{common.BKSetIDField: common.KvMap{common.BKDBIN: setID}}
	dat.Limit = common.BKNoLimit
	ret, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(ctx, common.BKInnerObjIDSet, header, dat)
	if nil != err {
		blog.Errorf("getSetNameBySetID setID %v supplierID %s  http do error:%s", setID, supplierID, err.Error())
		return "", err
	}
	if !ret.Result {
		blog.Errorf("getSetNameBySetID setID %v supplierID %s  http reply error:%s", setID, supplierID, ret.ErrMsg)
		return "", fmt.Errorf(ret.ErrMsg)

	}
	if 0 == ret.Data.Count {
		blog.Errorf("getSetNameBySetID setID %v supplierID %s  not found set info", setID, supplierID)
		return "", nil
	}
	name, err := ret.Data.Info[0].String(common.BKSetNameField)
	if nil != err {
		blog.Errorf("getSetNameBySetID moduleID %v supplierID %s  get set name error:%s", setID, supplierID, err.Error())
		return "", err
	}
	return name, err

}

func (lgc *Logics) getModuleNameByID(ctx context.Context, ID int64, header http.Header) (name string, appID int64, setID int64, err error) {
	supplierID := util.GetOwnerID(header)

	dat := new(metadata.QueryInput)
	dat.Condition = common.KvMap{common.BKModuleIDField: ID}
	dat.Fields = fmt.Sprintf("%s,%s,%s", common.BKModuleNameField, common.BKAppIDField, common.BKSetIDField)
	dat.Limit = common.BKNoLimit
	ret, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(ctx, common.BKInnerObjIDModule, header, dat)
	if nil != err {
		blog.Errorf("getModuleNameByID moduleID %v supplierID %s  http do error:%s", ID, supplierID, err.Error())
		return "", 0, 0, err
	}
	if !ret.Result {
		blog.Errorf("getModuleNameByID moduleID %v supplierID %s  http reply error:%s", ID, supplierID, ret.ErrMsg)
		return "", 0, 0, fmt.Errorf(ret.ErrMsg)

	}
	if 0 == ret.Data.Count {
		blog.Errorf("getModuleNameByID moduleID %v supplierID %s  not found module info", ID, supplierID)
		return "", 0, 0, nil
	}
	byteModule, _ := json.Marshal(ret.Data.Info[0])
	name, err = ret.Data.Info[0].String(common.BKSetNameField)
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

func (lgc *Logics) getHostByHostID(ctx context.Context, header http.Header, hostID int64) (*metadata.GseHost, error) {

	supplierID := util.GetOwnerID(header)
	intSupplierID, err := util.GetInt64ByInterface(supplierID)
	if nil != err {
		blog.Errorf("getHostByHostID supplierID %s  not interger", supplierID)
		return nil, err
	}

	opt := new(metadata.QueryInput)
	opt.Condition = common.KvMap{common.BKHostIDField: hostID}
	opt.Fields = fmt.Sprintf("%s,%s,%s", common.BKHostIDField, common.BKHostInnerIPField, common.BKCloudIDField)
	opt.Limit = common.BKNoLimit
	hosts, err := lgc.CoreAPI.HostController().Host().GetHosts(ctx, header, opt)
	if nil != err {
		blog.Errorf("getHostByHostID moduleID %d hostID:%v supplierID %s GetHosts http do error:%s", hostID, supplierID, err.Error())
		return nil, err
	}
	if !hosts.Result {
		blog.Errorf("getHostByHostID moduleID %d hostID:%v supplierID %s GetHosts http reply error:%s", hostID, supplierID, hosts.ErrMsg)
		return nil, fmt.Errorf(hosts.ErrMsg)
	}

	host := hosts.Data.Info[0]

	item := new(metadata.GseHost)

	hostID, err = util.GetInt64ByInterface(host[common.BKHostIDField])
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

	return item, nil
}
