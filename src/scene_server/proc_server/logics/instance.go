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

func (lgc *Logics) bgHandle(ctx context.Context) {
	lgc.reshReshInitChan(ctx, lgc.procHostInst)
	go func() {
		defer lgc.bgHandle(ctx)
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
				// use new header, so, new logics struct
				newLgc := lgc.NewFromHeader(item.Header)
				err := newLgc.HandleProcInstNumByModuleID(ctx, item.AppID, item.ModuleID)
				if nil != err {
					blog.Errorf("HandleProcInstNumByModuleID  error %s,rid:%s", err.Error(), lgc.rid)
				}

			}

		}

	}()

}

func (lgc *Logics) HandleProcInstNumByModuleID(ctx context.Context, appID, moduleID int64) error {
	//defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	maxInstID, procInst, err := lgc.getProcInstInfoByModuleID(ctx, appID, moduleID)
	if nil != err {
		blog.Errorf("handleInstanceNum getProcInstInfoByModuleID error %s", err.Error())
		return err
	}
	var hostInfos map[int64]*metadata.GseHost
	hostInfos, err = lgc.getHostByModuleID(ctx, moduleID)
	if nil != err {
		blog.Errorf("handleInstanceNum getHostByModuleID error %s,rid:%s", err.Error(), lgc.rid)
		return err
	}
	setID, procIDs, err := lgc.getModuleBindProc(ctx, appID, moduleID, lgc.header)
	if nil != err {
		return err
	}
	instProc := make([]*metadata.ProcInstanceModel, 0)
	procInfos, err := lgc.getProcInfoByID(ctx, procIDs)
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

	err = lgc.setProcInstDetallStatusUnregister(ctx, appID, moduleID)
	if nil != err {
		return err
	}
	err = lgc.handleProcInstNumDataHandle(ctx, appID, moduleID, instProc)
	if nil != err {
		return err
	}
	gseHost := make([]metadata.GseHost, 0)
	for _, host := range hostInfos {
		gseHost = append(gseHost, *host)
	}
	if 0 != len(gseHost) {
		for _, info := range procInfos {
			err = lgc.RegisterProcInstanceToGse(ctx, info.AppID, moduleID, info.ProcID, gseHost, info.ProcInfo)
			if nil != err {
				blog.Errorf("RegisterProcInstanceToGse error %s,rid:%s", err.Error(), lgc.rid)
				continue
			}
		}
	}

	err = lgc.unregisterProcInstDetall(ctx, appID, moduleID)
	if nil != err {
		return err
	}

	return nil
}

func (lgc *Logics) eventHostModuleChangeProcHostInstNum(ctx context.Context, eventData *metadata.EventInst) error {
	var header http.Header = make(http.Header, 0)
	header.Set(common.BKHTTPOwnerID, eventData.OwnerID)
	header.Set(common.BKHTTPHeaderUser, common.BKProcInstanceOpUser)
	// use new header, so, new logics struct
	newLgc := lgc.NewFromHeader(header)
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
			blog.Errorf("productHostInstanceNum event data not map[string]interface{} item %v raw josn %s,rid:%s", hostInfos, string(byteData), newLgc.rid)
			return err
		}
		appID, err := mapData.Int64(common.BKAppIDField)
		if nil != err {
			byteData, _ := json.Marshal(eventData)
			blog.Errorf("productHostInstanceNum event data  appID not integer  item %v raw josn %s,rid:%s", hostInfos, string(byteData), newLgc.rid)
			return err
		}
		moduleID, err := mapData.Int64(common.BKModuleIDField)
		if nil != err {
			byteData, _ := json.Marshal(eventData)
			blog.Errorf("productHostInstanceNum event data  moduleID not integer  item %v raw josn %s,rid:%s", hostInfos, string(byteData), newLgc.rid)
			return err
		}
		newLgc.addEventRefreshModuleItem(appID, moduleID)
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
	// use new header, so, new logics struct
	newLgc := lgc.NewFromHeader(header)

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
			blog.Errorf("eventProcInstByProcModule event data not map[string]interface{} item %v raw josn %s,rid:%s", data, string(byteData), newLgc.rid)
			return err
		}
		appID, err := mapData.Int64(common.BKAppIDField)
		if nil != err {
			byteData, _ := json.Marshal(eventData)
			blog.Errorf("eventProcInstByProcModule event data  appID not integer  item %v raw josn %s,rid:%s", data, string(byteData), newLgc.rid)
			return err
		}
		moduleName, err := mapData.String(common.BKModuleNameField)
		if nil != err {
			byteData, _ := json.Marshal(eventData)
			blog.Errorf("eventProcInstByProcModule event data  appID not integer  item %v raw josn %s,rid:%s", data, string(byteData), newLgc.rid)
			return err
		}

		moduleID, err := newLgc.HandleProcInstNumByModuleName(ctx, appID, moduleName)
		if nil != err {
			byteData, _ := json.Marshal(eventData)
			blog.Errorf("eventProcInstByProcModule HandleProcInstNumByModuleName error %s item %v raw josn %s,rid:%s", err.Error(), data, string(byteData), newLgc.rid)
			return err
		}
		newLgc.addEventRefreshModuleItems(ctx, appID, moduleID)
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
	// use new header, so, new logics struct
	newLgc := lgc.NewFromHeader(header)

	for _, data := range eventData.Data {
		mapData, err := mapstr.NewFromInterface(data.CurData)
		if nil != err {
			byteData, _ := json.Marshal(eventData)
			blog.Errorf("eventProcInstByProcess event data not map[string]interface{} item %v raw josn %s,rid:%s", data, string(byteData), newLgc.rid)
			return err
		}
		procID, err := mapData.Int64(common.BKProcessIDField)
		if nil != err {
			byteData, _ := json.Marshal(eventData)
			blog.Errorf("eventProcInstByProcess process id not integer error %s item %v raw josn %s,rid:%s", err.Error(), data, string(byteData), newLgc.rid)
			return err
		}
		appID, err := mapData.Int64(common.BKAppIDField)
		if nil != err {
			byteData, _ := json.Marshal(eventData)
			blog.Errorf("eventProcInstByProcess application id not integer error %s item %v raw josn %s,rid:%s", err.Error(), data, string(byteData), newLgc.rid)
			return err
		}
		mdouleID, err := newLgc.getModuleIDByProcID(ctx, appID, procID)
		if nil != err {
			byteData, _ := json.Marshal(eventData)
			blog.Errorf("eventProcInstByProcess get process bind module info  by appID %d, procID %d error %s item %v raw josn %s,rid:%s", appID, procID, err.Error(), data, string(byteData), lgc.rid)
			return err
		}
		newLgc.addEventRefreshModuleItems(ctx, appID, mdouleID)
	}
	return nil
}

func (lgc *Logics) eventProcInstByHostInfo(ctx context.Context, eventData *metadata.EventInst) error {

	if metadata.EventActionUpdate == eventData.Action {
		var header http.Header = make(http.Header, 0)
		header.Set(common.BKHTTPOwnerID, eventData.OwnerID)
		header.Set(common.BKHTTPHeaderUser, common.BKProcInstanceOpUser)
		// use new header, so, new logics struct
		newLgc := lgc.NewFromHeader(header)

		// host clouid id change
		for _, data := range eventData.Data {
			mapCurData, err := mapstr.NewFromInterface(data.CurData)
			if nil != err {
				byteData, _ := json.Marshal(eventData)
				blog.Errorf("eventProcInstByHostInfo event current data not map[string]interface{} item %v raw josn %s,rid:%s", data, string(byteData), newLgc.rid)
				return err
			}

			mapPreData, err := mapstr.NewFromInterface(data.PreData)
			if nil != err {
				byteData, _ := json.Marshal(eventData)
				blog.Errorf("eventProcInstByHostInfo event pre-data not map[string]interface{} item %v raw josn %s,rid:%s", data, string(byteData), newLgc.rid)
				return err
			}
			curData, err := mapCurData.Int64(common.BKCloudIDField)
			if nil != err {
				byteData, _ := json.Marshal(eventData)
				blog.Errorf("eventProcInstByHostInfo event current data cloud id not int item %v raw josn %s,rid:%s", data, string(byteData), newLgc.rid)
				return err
			}
			preData, err := mapPreData.Int64(common.BKCloudIDField)
			if nil != err {
				byteData, _ := json.Marshal(eventData)
				blog.Errorf("eventProcInstByHostInfo event pre-data  cloud id not int item %v raw josn %s,rid:%s", data, string(byteData), newLgc.rid)
				return err
			}
			if curData != preData {
				hostID, err := mapCurData.Int64(common.BKHostIDField)
				if nil != err {
					byteData, _ := json.Marshal(eventData)
					blog.Errorf("eventProcInstByHostInfo event hostID not int item %v raw josn %s,rid:%s", data, string(byteData), newLgc.rid)
					return err
				}
				hostModule, err := newLgc.GetModuleIDByHostID(ctx, hostID)
				if nil != err {
					byteData, _ := json.Marshal(eventData)
					blog.Errorf("eventProcInstByHostInfo event hostID %s get module err :%s  item %v raw josn %s,rid:%s", hostID, err.Error(), data, string(byteData), newLgc.rid)
					return err
				}
				for _, item := range hostModule {
					newLgc.addEventRefreshModuleItem(item.AppID, item.ModuleID)
				}
			}
		}
	}
	return nil
}

func (lgc *Logics) HandleProcInstNumByModuleName(ctx context.Context, appID int64, moduleName string) ([]int64, error) {
	moduleConds := new(metadata.SearchParams)
	moduleConds.Condition = map[string]interface{}{
		common.BKAppIDField:      appID,
		common.BKModuleNameField: moduleName,
	}
	defErr := lgc.ccErr
	moduleConds.Page = map[string]interface{}{"start": 0, "limit": common.BKNoLimit}
	moduleInfos, err := lgc.CoreAPI.TopoServer().Instance().InstSearch(ctx, lgc.ownerID, common.BKInnerObjIDModule, lgc.header, moduleConds)
	if nil != err {
		blog.Errorf("HandleProcInstNumByModuleName http do error. get module by name %s error %s,input:%+v,rid:%s", moduleName, err.Error(), moduleConds, lgc.rid)
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !moduleInfos.Result {
		blog.Errorf("HandleProcInstNumByModuleName http reply error. get module by name %s error. err code:%d,err msg:%s,input:%+v,rid:%s", moduleName, moduleInfos.Code, moduleInfos.ErrMsg, moduleConds, lgc.rid)
		return nil, defErr.New(moduleInfos.Code, moduleInfos.ErrMsg)
	}
	moduleIDs := make([]int64, 0)
	for _, module := range moduleInfos.Data.Info {
		moduleID, err := module.Int64(common.BKModuleIDField)
		if nil != err {
			byteData, _ := json.Marshal(module)
			blog.Errorf("refreshProcInstByProcModule get module by name %s not module id error  item %+v,input:%+v,rid:%s", moduleName, string(byteData), moduleConds, lgc.rid)
			return nil, err
		}
		moduleIDs = append(moduleIDs, moduleID)
	}

	return moduleIDs, nil
}

func (lgc *Logics) DeleteProcInstanceModel(ctx context.Context, appId, procId, moduleName string) error {
	condition := make(map[string]interface{})
	condition[common.BKAppIDField] = appId
	condition[common.BKProcessIDField] = procId
	condition[common.BKModuleNameField] = moduleName

	ret, err := lgc.CoreAPI.ProcController().DeleteProcInstanceModel(ctx, lgc.header, condition)
	if nil != err {
		blog.Errorf("DeleteProcInstanceModel DeleteProcInstanceModel http do error. error %s,input:%+v,rid:%s", err.Error(), condition, lgc.rid)
		return lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !ret.Result {
		blog.Errorf("DeleteProcInstanceModel DeleteProcInstanceModel http reply error.err code:%d,err msg:%s,input:%+v,rid:%s", ret.Code, ret.ErrMsg, condition, lgc.rid)
		return lgc.ccErr.New(ret.Code, ret.ErrMsg)
	}

	return nil
}
