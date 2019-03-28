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
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

//var handEventDataChan chan chanItem // := make(chan chanItem, 10000)

func (lgc *Logics) HandleHostProcDataChange(ctx context.Context, eventData *metadata.EventInst) {

	switch eventData.ObjType {
	case metadata.EventObjTypeProcModule:
		handEventDataChan <- chanItem{ctx: ctx, eventData: eventData, opFunc: lgc.eventProcInstByProcModule}
	case metadata.EventObjTypeModuleTransfer:
		handEventDataChan <- chanItem{ctx: ctx, eventData: eventData, opFunc: lgc.eventHostModuleChangeProcHostInstNum}
	case common.BKInnerObjIDHost:
		handEventDataChan <- chanItem{ctx: ctx, eventData: eventData, opFunc: lgc.eventProcInstByHostInfo}
	case common.BKInnerObjIDProc:
		handEventDataChan <- chanItem{ctx: ctx, eventData: eventData, opFunc: lgc.eventProcInstByProcess}
	}

}

func (lgc *Logics) bgHandle() {
	lgc.reshReshInitChan(lgc.ProcHostInst)
	go func() {
		defer lgc.bgHandle()
		for {
			select {
			case item := <-handEventDataChan:
				err := item.opFunc(item.ctx, item.eventData)
				if nil != err {
					blog.Error(err.Error())
				}
			case item := <-refreshHostInstModuleIDChan:
				if nil == item {
					break
				}
				err := lgc.HandleProcInstNumByModuleID(context.Background(), item.Header, item.AppID, item.ModuleID)
				if nil != err {
					blog.Errorf("HandleProcInstNumByModuleID  error %s", err.Error())
				}

			}

		}

	}()

}

func (lgc *Logics) HandleProcInstNumByModuleID(ctx context.Context, header http.Header, appID, moduleID int64) error {
	//defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	maxInstID, procInst, err := lgc.getProcInstInfoByModuleID(ctx, appID, moduleID, header)
	if nil != err {
		blog.Errorf("handleInstanceNum getProcInstInfoByModuleID error %s", err.Error())
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
	for procID, info := range procInfos {
		for hostID := range hostInfos {
			procInstInfo, ok := procInst[getInlineProcInstKey(hostID, moduleID)]
			var hostInstID uint64
			if !ok {
				maxInstID++
				hostInstID = maxInstID
				isExistHostInst[getInlineProcInstKey(hostID, moduleID)] = procInstInfo
			} else {
				hostInstID = procInstInfo.HostInstanID
			}
			instProc = append(instProc, GetProcInstModel(appID, setID, moduleID, hostID, procID, info.FunID, info.ProcNum, hostInstID)...)
		}

	}

	err = lgc.setProcInstDetallStatusUnregister(ctx, header, appID, moduleID)
	if nil != err {
		return err
	}
	err = lgc.handleProcInstNumDataHandle(ctx, header, appID, moduleID, instProc)
	if nil != err {
		return err
	}
	gseHost := make([]metadata.GseHost, 0)
	for _, host := range hostInfos {
		gseHost = append(gseHost, *host)
	}
	if 0 != len(gseHost) {
		for _, info := range procInfos {
			err = lgc.RegisterProcInstanceToGse(info.AppID, moduleID, info.ProcID, gseHost, info.ProcInfo, header)
			if nil != err {
				blog.Errorf("RegisterProcInstanceToGse error %s", err.Error())
				continue
			}
		}
	}

	err = lgc.unregisterProcInstDetall(ctx, header, appID, moduleID)
	if nil != err {
		return err
	}

	return nil
}

func (lgc *Logics) eventHostModuleChangeProcHostInstNum(ctx context.Context, eventData *metadata.EventInst) error {
	var header http.Header = make(http.Header, 0)
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
		lgc.addEventRefreshModuleItem(appID, moduleID, header)
	}
	return nil
}

func (lgc *Logics) eventProcInstByProcModule(ctx context.Context, eventData *metadata.EventInst) error {
	if metadata.EventTypeRelation != eventData.EventType {
		return nil
	}
	var header http.Header = make(http.Header, 0)
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
		lgc.addEventRefreshModuleItems(appID, moduleID, header)
	}
	return nil
}

func (lgc *Logics) eventProcInstByProcess(ctx context.Context, eventData *metadata.EventInst) error {
	if metadata.EventActionUpdate != eventData.Action {
		// create proccess not refresh process instance , because not bind module
		return nil
	}
	var header http.Header = make(http.Header, 0)
	header.Set(common.BKHTTPOwnerID, eventData.OwnerID)
	header.Set(common.BKHTTPHeaderUser, common.BKProcInstanceOpUser)

	for _, data := range eventData.Data {
		mapData, err := mapstr.NewFromInterface(data.CurData)
		if nil != err {
			byteData, _ := json.Marshal(eventData)
			blog.Errorf("eventProcInstByProcess event data not map[string]interface{} item %v raw josn %s", data, string(byteData))
			return err
		}
		procID, err := mapData.Int64(common.BKProcessIDField)
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
		lgc.addEventRefreshModuleItems(appID, mdouleID, header)
	}
	return nil
}

func (lgc *Logics) eventProcInstByHostInfo(ctx context.Context, eventData *metadata.EventInst) error {

	if metadata.EventActionUpdate == eventData.Action {
		var header http.Header = make(http.Header, 0)
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
					lgc.addEventRefreshModuleItem(item.AppID, item.ModuleID, header)
				}
			}
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
