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
	"configcenter/src/common/mapstr"
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
	handEventDataChan      chan chanItem = make(chan chanItem, 10000)
	chnOpLock              *sync.Once
	eventRefreshModuleData = &eventRefreshModule{}
)

func _init() {
	eventRefreshModuleData.data = make(map[string]*refreshModuleData, 0)
	eventRefreshModuleData.eventChn = make(chan bool, 10)
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

	sendEventFrefreshModuleNotice()
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
	if 5 > len(eventRefreshModuleData.eventChn) {
		eventRefreshModuleData.eventChn <- true

	}
}

//var handEventDataChan chan chanItem // := make(chan chanItem, 10000)

func (lgc *Logics) HandleHostProcDataChange(ctx context.Context, eventData *metadata.EventInst) {

	switch eventData.ObjType {
	case metadata.EventObjTypeProcModule:
		handEventDataChan <- chanItem{ctx: ctx, eventData: eventData, opFunc: lgc.eventProcInstByProcModule, retry: 3}
		//lgc.handleRetry(3, ctx, eventData, lgc.refreshProcInstByProcModule)
		//lgc.refreshProcInstByProcModule(ctx, eventData)
	case metadata.EventObjTypeModuleTransfer:
		handEventDataChan <- chanItem{ctx: ctx, eventData: eventData, opFunc: lgc.eventHostModuleChangeProcHostInstNum, retry: 3}
		//lgc.handleRetry(3, ctx, eventData, lgc.eventHostModuleChangeProcHostInstNum)
		//lgc.eventHostModuleChangeProcHostInstNum(ctx, eventData)
	case common.BKInnerObjIDHost:
		handEventDataChan <- chanItem{ctx: ctx, eventData: eventData, opFunc: lgc.eventProcInstByHostInfo, retry: 3}
	case common.BKInnerObjIDProc:
		handEventDataChan <- chanItem{ctx: ctx, eventData: eventData, opFunc: lgc.eventProcInstByProcess, retry: 3}

	}
	chnOpLock.Do(lgc.bgHandle)

}

func (lgc *Logics) bgHandle() {
	go func() {
		defer lgc.bgHandle()
		for {
			select {
			case item := <-handEventDataChan:
				for idx := 0; idx < item.retry; idx++ {
					err := item.opFunc(item.ctx, item.eventData)
					if nil == err {
						break
					}
				}
			case <-eventRefreshModuleData.eventChn:
				for {
					item := getEventRefreshModuleItem()
					if nil == item {
						break
					}
					err := lgc.HandleProcInstNumByModuleID(context.Background(), item.header, item.appID, item.moduleID)
					if nil != err {
						blog.Errorf("HandleProcInstNumByModuleID  error %s", err.Error())
					}
				}
			}

		}

	}()

}

func (lgc *Logics) handleProcInstDetailStatus(ctx context.Context, header http.Header, data, conds map[string]interface{}) {

}

func (lgc *Logics) HandleProcInstNumByModuleID(ctx context.Context, header http.Header, appID, moduleID int64) error {
	maxInstID, procInst, err := lgc.getProcInstInfoByModuleID(ctx, appID, moduleID, header)
	if nil != err {
		return err
	}
	var hostInfos map[int64]*metadata.GseHost
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
	isExistHostInst := make(map[string]metadata.ProcInstanceModel)
	gseRegister := make([]metadata.GseProcRequest, 0)
	for procID, info := range procInfos {
		gseHost :=make([]metadata.GseHost, 0)
		for hostID, host := range hostInfos {
			procInstInfo, ok := procInst[getInlineProcInstKey(hostID, procID)]
			hostInstID := uint64(0)
			if !ok {
				maxInstID++
				hostInstID = maxInstID
			} else {
				hostInstID = procInstInfo.HostInstanID
				isExistHostInst[getInlineProcInstKey(hostID, procID)] = procInstInfo
			}
			gseHost = append(gseHost, *host)
			instProc = append(instProc, GetProcInstModel(appID, setID, moduleID, hostID, procID, info.FunID, info.FunID, hostInstID)...)
		}
		if 0 < len(gseHost) {
			gseProcReq := new(metadata.GseProcRequest)
			gseProcReq.AppID = appID
			gseProcReq.Hosts = gseHost 
			gseProcReq.Meta = info.Meta
			gseProcReq.ModuleID = moduleID
			gseProcReq.OpType = 1 
			gseProcReq.ProcID = procID 
			gseProcReq.Spec = info.Spec
		}
	}

	err = lgc.setProcInstDetallStatusUnregister(ctx, header, appID, moduleID, unregisterProcDetail)
	if nil != err {
		return err
	}
	err = lgc.handleProcInstNumDataHandle(ctx, header, appID, moduleID, procIDs, instProc)
	if nil != err {
		return err
	}

	gseRegister, gseUnRegister []*metadata.GseProcRequest
	for procID, info := range procInfos {
		for hostID, _ := range hostInfos {
			procInstInfo, ok := procInst[getInlineProcInstKey(hostID, procID)]
			hostInstID := uint64(0)
			if !ok {
				maxInstID++
				hostInstID = maxInstID
			} else {
				hostInstID = procInstInfo.HostInstanID
				isExistHostInst[getInlineProcInstKey(hostID, procID)] = procInstInfo
			}
			instProc = append(instProc, GetProcInstModel(appID, setID, moduleID, hostID, procID, info.FunID, info.FunID, hostInstID)...)

		}
	}

	return nil
}

// setProcInstDetallStatusUnregister modify process instance status to unregister in cmdb table
func (lgc *Logics) setProcInstDetallStatusUnregister(ctx context.Context, header http.Header, appID, moduleID int64, unregister []metadata.ProcInstanceModel) error {

	if 0 != len(unregister) {
		unregisterProcDetail := make([]interface{}, 0)
		for _, item := range unregister {
			unregisterProcDetail = append(unregisterProcDetail, common.KvMap{common.BKAppIDField: item.ApplicationID, common.BKModuleIDField: item.ModuleID, common.BKHostIDField: item.HostID, common.BKProcIDField: item.ProcID})
		}
		dat := new(metadata.ModifyProcInstanceStatus)
		dat.Conds = map[string]interface{}{common.BKDBOR: unregisterProcDetail}
		dat.Data = map[string]interface{}{"status": metadata.ProcInstanceDetailStatusUnRegisterFailed}
		ret, err := lgc.CoreAPI.ProcController().ModifyProcInstanceDetail(ctx, header, dat)
		if nil != err {
			blog.Errorf("setProcInstDetallStatusUnregister  proc instance error:%s", err.Error())
			return fmt.Errorf("unregister  proc instance error:%s", err.Error())
		}
		if !ret.Result {
			blog.Errorf("setProcInstDetallStatusUnregister  proc instance return err msg %s", ret.ErrMsg)
			return fmt.Errorf("unregister proc instance return err msg %s", ret.ErrMsg)
		}

	}
	return nil
}

func (lgc *Logics) handleProcInstNumDataHandle(ctx context.Context, header http.Header, appID, moduleID int64, procIDs []int64, instProc []*metadata.ProcInstanceModel) error {
	delConds := common.KvMap{common.BKAppIDField: appID, common.BKModuleIDField: moduleID, common.BKProcIDField: procIDs}

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

func (lgc *Logics) HandleProcInstNumByModuleName(ctx context.Context, header http.Header, appID int64, moduleName string) ([]int64, error) {
	moduleConds := new(metadata.SearchParams)
	moduleConds.Condition = map[string]interface{}{
		common.BKAppIDField: appID,
		moduleName:          moduleName,
	}
	moduleConds.Page = map[string]interface{}{"start": 0, "limit": 1}
	moduleInfos, err := lgc.CoreAPI.TopoServer().Instance().InstSearch(ctx, util.GetOwnerID(header), common.BKInnerObjIDModule, header, moduleConds)
	if nil != err {
		blog.Errorf("HandleProcInstNumByModuleName get module by name %s error %s", moduleName, err.Error())
		return nil, err
	}
	if !moduleInfos.Result {
		blog.Errorf("HandleProcInstNumByModuleName get module by name %s error %s ", moduleName, moduleInfos.ErrMsg)
		return nil, fmt.Errorf("%s", moduleInfos.ErrMsg)
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

func (lgc *Logics) getModuleIDByProcID(ctx context.Context, appID, procID int64, header http.Header) ([]int64, error) {
	condition := make(map[string]interface{}, 0)
	condition[common.BKProcIDField] = procID
	// get process by module
	ret, err := lgc.CoreAPI.ProcController().GetProc2Module(context.Background(), header, condition)
	if nil != err {
		blog.Errorf("getModuleIDByProcID  GetProc2Module appID %d moduleID %s  http do error:%s", appID, procID, err.Error())
		return nil, err
	}
	if !ret.Result {
		blog.Errorf("getModuleIDByProcID  GetProc2Module appID %d moduleID %s  http reply error:%s", appID, procID, ret.ErrMsg)
		return nil, fmt.Errorf(ret.ErrMsg)
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
		addEventRefreshModuleItem(appID, moduleID, header)

	}
	sendEventFrefreshModuleNotice()
	return nil
}

func (lgc *Logics) eventProcInstByProcModule(ctx context.Context, eventData *metadata.EventInst) error {
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
			blog.Errorf("eventProcInstByProcModule event data not map[string]interface{} item %v raw josn %s", data, string(byteData))
			return err
		}
		appID, err := mapData.Int64(common.BKAppIDField)
		if nil != err {
			byteData, _ := json.Marshal(eventData)
			blog.Errorf("eventProcInstByProcModule event data  appID not integer  item %v raw josn %s", data, string(byteData))
			return err
		}
		moduleName, err := mapData.String(common.BKModuleNameField)
		if nil != err {
			byteData, _ := json.Marshal(eventData)
			blog.Errorf("eventProcInstByProcModule event data  appID not integer  item %v raw josn %s", data, string(byteData))
			return err
		}

		moduleID, err := lgc.HandleProcInstNumByModuleName(ctx, header, appID, moduleName)
		if nil != err {
			byteData, _ := json.Marshal(eventData)
			blog.Errorf("eventProcInstByProcModule HandleProcInstNumByModuleName error %s item %v raw josn %s", err.Error(), data, string(byteData))
			return err
		}
		addEventRefreshModuleItems(appID, moduleID, header)
	}
	return nil
}

func (lgc *Logics) eventProcInstByProcess(ctx context.Context, eventData *metadata.EventInst) error {
	if metadata.EventActionCreate != eventData.Action {
		// create proccess not refresh process instance , because not bind module
		return nil
	}
	var header http.Header
	header.Set(common.BKHTTPOwnerID, eventData.OwnerID)
	header.Set(common.BKHTTPHeaderUser, common.BKProcInstanceOpUser)

	for _, data := range eventData.Data {
		mapData, err := mapstr.NewFromInterface(data.CurData)
		if nil != err {
			byteData, _ := json.Marshal(eventData)
			blog.Errorf("eventProcInstByProcess event data not map[string]interface{} item %v raw josn %s", data, string(byteData))
			return err
		}
		procID, err := mapData.Int64(common.BKProcIDField)
		if nil != err {
			byteData, _ := json.Marshal(eventData)
			blog.Errorf("eventProcInstByProcess process id not integer error %s item %v raw josn %s", err.Error(), data, string(byteData))
			return err
		}
		appID, err := mapData.Int64(common.BKAppIDField)
		if nil != err {
			byteData, _ := json.Marshal(eventData)
			blog.Errorf("eventProcInstByProcess application id not integer error %s item %v raw josn %s", err.Error(), data, string(byteData))
			return err
		}
		mdouleID, err := lgc.getModuleIDByProcID(ctx, appID, procID, header)
		if nil != err {
			byteData, _ := json.Marshal(eventData)
			blog.Errorf("eventProcInstByProcess get process bind module info  by appID %d, procID %d error %s item %v raw josn %s", appID, procID, err.Error(), data, string(byteData))
			return err
		}
		addEventRefreshModuleItems(appID, mdouleID, header)
	}
	sendEventFrefreshModuleNotice()
	return nil
}

func (lgc *Logics) eventProcInstByHostInfo(ctx context.Context, eventData *metadata.EventInst) error {

	if metadata.EventActionUpdate == eventData.Action {
		var header http.Header
		header.Set(common.BKHTTPOwnerID, eventData.OwnerID)
		header.Set(common.BKHTTPHeaderUser, common.BKProcInstanceOpUser)

		// host clouid id change
		for _, data := range eventData.Data {
			mapCurData, err := mapstr.NewFromInterface(data.CurData)
			if nil != err {
				byteData, _ := json.Marshal(eventData)
				blog.Errorf("eventProcInstByHostInfo event current data not map[string]interface{} item %v raw josn %s", data, string(byteData))
				return err
			}

			mapPreData, err := mapstr.NewFromInterface(data.PreData)
			if nil != err {
				byteData, _ := json.Marshal(eventData)
				blog.Errorf("eventProcInstByHostInfo event pre-data not map[string]interface{} item %v raw josn %s", data, string(byteData))
				return err
			}
			curData, err := mapCurData.Int64(common.BKCloudIDField)
			if nil != err {
				byteData, _ := json.Marshal(eventData)
				blog.Errorf("eventProcInstByHostInfo event current data cloud id not int item %v raw josn %s", data, string(byteData))
				return err
			}
			preData, err := mapPreData.Int64(common.BKCloudIDField)
			if nil != err {
				byteData, _ := json.Marshal(eventData)
				blog.Errorf("eventProcInstByHostInfo event pre-data  cloud id not int item %v raw josn %s", data, string(byteData))
				return err
			}
			if curData != preData {
				hostID, err := mapCurData.Int64(common.BKHostIDField)
				if nil != err {
					byteData, _ := json.Marshal(eventData)
					blog.Errorf("eventProcInstByHostInfo event hostID not int item %v raw josn %s", data, string(byteData))
					return err
				}
				hostModule, err := lgc.GetModuleIDByHostID(ctx, header, hostID)
				if nil != err {
					byteData, _ := json.Marshal(eventData)
					blog.Errorf("eventProcInstByHostInfo event hostID %s get module err :%s  item %v raw josn %s", hostID, err.Error(), data, string(byteData))
					return err
				}
				for _, item := range hostModule {
					addEventRefreshModuleItem(item.AppID, item.ModuleID, header)
				}
			}
		}
	}
	return nil
}

func (lgc *Logics) GetModuleIDByHostID(ctx context.Context, header http.Header, hostID int64) ([]metadata.ModuleHost, error) {
	dat := map[string][]int64{
		common.BKHostIDField: []int64{hostID},
	}
	ret, err := lgc.CoreAPI.HostController().Module().GetModulesHostConfig(ctx, header, dat)
	if nil != err {
		blog.Errorf("GetModuleIDByHostID appID %d module id %d GetModulesHostConfig http do error:%s", hostID, err.Error())
		return nil, err
	}
	if !ret.Result {
		blog.Errorf("GetModuleIDByHostID appID %d module id %d GetModulesHostConfig reply error:%s", hostID, ret.ErrMsg)
		return nil, fmt.Errorf(ret.ErrMsg)
	}

	return ret.Data, nil
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

func (lgc *Logics) getSetNameBySetID_deletefunc(ctx context.Context, setID int64, header http.Header) (string, error) {
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

func (lgc *Logics) getHostByHostID_deletefunc(ctx context.Context, header http.Header, hostID int64) (*metadata.GseHost, error) {

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
