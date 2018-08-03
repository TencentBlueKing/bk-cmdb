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
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (lgc *Logics) RefreshProcInstance(ctx context.Context, eventData *metadata.EventInst) {
	switch eventData.ObjType {
	case metadata.EventObjTypeProcModule:
		lgc.refreshProcInstByProcModule(ctx, eventData)
	}
}

func (lgc *Logics) refreshProcInstByProcModule(ctx context.Context, eventData *metadata.EventInst) {
	if metadata.EventTypeRelation != eventData.EventType {
		return
	}
	if metadata.EventActionDelete == eventData.Action {
		// delete process bind module relation, unregister process info
	} else {
		// compare  pre-change data with the current data and find the newly added data to register process info

	}

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

func (lgc *Logics) getHostByModuleID(ctx context.Context, moduleID int64, header http.Header) error {
	dat := map[string][]int64{
		common.BKModuleIDField: []int64{moduleID},
	}
	supplierID := util.GetOwnerID(header)
	intSupplierID, err := util.GetInt64ByInterface(supplierID)
	if nil != err {
		blog.Errorf("getHostByModuleID supplierID %s  not interger", supplierID)
		return err
	}

	ret, err := lgc.CoreAPI.HostController().Module().GetModulesHostConfig(ctx, header, dat)
	if nil != err {
		blog.Errorf("getHostByModuleID moduleID %d supplierID %s GetModulesHostConfig http do error:%s", moduleID, supplierID, err.Error())
		return err
	}
	if !ret.Result {
		blog.Errorf("getHostByModuleID moduleID %d supplierID %s GetModulesHostConfig http reply error:%s", moduleID, supplierID, ret.ErrMsg)
		return fmt.Errorf(ret.ErrMsg)

	}
	if 0 == len(ret.Data) {
		blog.V(5).Infof("getHostByModuleID moduleID %d supplierID %s GetModulesHostConfig len equal 0", moduleID, supplierID)
		return nil
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
		return err
	}
	if !hosts.Result {
		blog.Errorf("getHostByModuleID moduleID %d hostID:%v supplierID %s GetHosts http reply error:%s", moduleID, hostIDs, supplierID, hosts.ErrMsg)
		return fmt.Errorf(hosts.ErrMsg)
	}

	hostInfos := make(map[int64]*metadata.GseHost, len(hosts.Data.Info))
	for _, host := range hosts.Data.Info {
		item := new(metadata.GseHost)

		hostID, err := util.GetInt64ByInterface(host[common.BKHostIDField])
		if nil != err {
			blog.Errorf("getHostByModuleID hostInfo %v  hostID   not interger", host)
			return err
		}
		cloudID, err := util.GetInt64ByInterface(host[common.BKCloudIDField])
		if nil != err {
			strHost, _ := json.Marshal(host)
			blog.Errorf("getHostByModuleID  hostInfo %v  cloudID  not interger, json:%s", host, string(strHost))
			return err
		}
		innerIP, ok := host[common.BKHostInnerIPField].(string)
		if !ok {
			strHost, _ := json.Marshal(host)
			blog.Errorf("getHostByModuleID  hostInfo %v  innerip  not found, json:%s", host, string(strHost))
			return err
		}
		item.BkCloudId = cloudID
		item.Ip = innerIP
		item.BkSupplierId = intSupplierID

		hostInfos[hostID] = item
	}

	return nil
}

func (lgc *Logics) getProcInfoByID(ctx context.Context, procID []int64, header http.Header) error {
	supplierID := util.GetOwnerID(header)

	gseProc := make(map[int64]*metadata.InlineProcInfo, 0)

	dat := new(metadata.QueryInput)
	dat.Condition = common.KvMap{common.BKProcIDField: common.KvMap{common.BKDBIN: procID}}
	dat.Limit = common.BKNoLimit
	ret, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(ctx, common.BKInnerObjIDProc, header, dat)
	if nil != err {
		blog.Errorf("getProcInfoByID procID %v supplierID %s  http do error:%s", procID, supplierID, err.Error())
		return err
	}
	if !ret.Result {
		blog.Errorf("getProcInfoByID procID %v supplierID %s  http reply error:%s", procID, supplierID, ret.ErrMsg)
		return fmt.Errorf(ret.ErrMsg)

	}
	if 0 == ret.Data.Count {
		blog.Errorf("getProcInfoByID procID %v supplierID %s  not found process info", procID, supplierID)
		return nil
	}
	for _, proc := range ret.Data.Info {
		procID, err := proc.Int64(common.BKProcField)
		if nil != err {
			strProc, _ := json.Marshal(proc)
			blog.Errorf("getHostByModuleID  proc %v  procID  not interger, json:%s", proc, string(strProc))
			return err
		}
		item := new(metadata.InlineProcInfo)

		item.ProcNum = 1
		procNumI, ok := proc.Get(common.BKProcInstNum)
		if ok && nil != procNumI {
			item.ProcNum, err = proc.Int64(common.BKProcInstNum)
			if nil != err {
				strProc, _ := json.Marshal(proc)
				blog.Errorf("getHostByModuleID  proc %v  procNum  not interger, json:%s", proc, string(strProc))
				return err
			}
		}
		item.AppID, err = proc.Int64(common.BKAppIDField)
		if nil != err {
			strProc, _ := json.Marshal(proc)
			blog.Errorf("getHostByModuleID  proc %v  AppID  not interger, json:%s", proc, string(strProc))
			return err
		}
		item.FunID, err = proc.Int64(common.BKFuncIDField)
		if nil != err {
			strProc, _ := json.Marshal(proc)
			blog.Errorf("getHostByModuleID  proc %v  AppID  not interger, json:%s", proc, string(strProc))
			return err
		}

		gseProc[procID] = item
	}

	return nil
}

func (lgc *Logics) getSetNameBySetID(ctx context.Context, setID int64, header http.Header) (string, error) {
	supplierID := util.GetOwnerID(header)

	dat := new(metadata.QueryInput)
	dat.Condition = common.KvMap{common.BKProcIDField: common.KvMap{common.BKDBIN: setID}}
	dat.Limit = common.BKNoLimit
	ret, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(ctx, common.BKInnerObjIDProc, header, dat)
	if nil != err {
		blog.Errorf("getMapSetBySetID procID %v supplierID %s  http do error:%s", setID, supplierID, err.Error())
		return "", err
	}
	if !ret.Result {
		blog.Errorf("getMapSetBySetID procID %v supplierID %s  http reply error:%s", setID, supplierID, ret.ErrMsg)
		return "", fmt.Errorf(ret.ErrMsg)

	}
	if 0 == ret.Data.Count {
		blog.Errorf("getMapSetBySetID procID %v supplierID %s  not found process info", setID, supplierID)
		return "", nil
	}
	return ret.Data.Info[0].String(common.BKSetNameField)

}

func (lgc *Logics) getHostByHostID(ctx context.Context, hostID int64, header http.Header) (*metadata.GseHost, error) {
	dat := map[string][]int64{
		common.BKModuleIDField: []int64{hostID},
	}
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
		strHost, _ := json.Marshal(host)
		blog.Errorf("getHostByModuleID  hostInfo %v  cloudID  not interger, json:%s", host, string(strHost))
		return nil, err
	}
	innerIP, ok := host[common.BKHostInnerIPField].(string)
	if !ok {
		strHost, _ := json.Marshal(host)
		blog.Errorf("getHostByModuleID  hostInfo %v  innerip  not found, json:%s", host, string(strHost))
		return nil, err
	}
	item.BkCloudId = cloudID
	item.Ip = innerIP
	item.BkSupplierId = intSupplierID

	return item, nil
}

func (lgc *Logics) getProcIDByHostID(ctx context.Context, hostID int64, header http.Header) {

}
