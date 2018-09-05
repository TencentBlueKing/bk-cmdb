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
			if 0 == len(item.Hosts) {
				continue
			}
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

	if 0 == len(unregister) {
		return nil
	}
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	unregisterProcDetail := make([]interface{}, 0)
	for _, item := range unregister {
		unregisterProcDetail = append(unregisterProcDetail, common.KvMap{common.BKAppIDField: item.ApplicationID, common.BKModuleIDField: item.ModuleID, common.BKHostIDField: item.HostID, common.BKProcessIDField: item.ProcID})
	}
	dat := new(metadata.ModifyProcInstanceDetail)
	dat.Conds = map[string]interface{}{common.BKDBOR: unregisterProcDetail}
	dat.Data = map[string]interface{}{common.BKStatusField: metadata.ProcInstanceDetailStatusUnRegisterFailed}
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
		procInst[getInlineProcInstKey(int64(item.HostID), int64(item.ModuleID))] = item
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
