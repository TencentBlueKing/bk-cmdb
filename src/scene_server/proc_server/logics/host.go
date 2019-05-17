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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (lgc *Logics) getHostByModuleID(ctx context.Context, moduleID int64) (map[int64]*metadata.GseHost, error) {
	dat := map[string][]int64{
		common.BKModuleIDField: []int64{moduleID},
	}
	supplierID := lgc.ownerID
	intSupplierID, err := util.GetInt64ByInterface(supplierID)
	defErr := lgc.ccErr
	if nil != err {
		blog.Errorf("getHostByModuleID supplierID %s  not interger", supplierID)
		return nil, err
	}

	ret, err := lgc.CoreAPI.HostController().Module().GetModulesHostConfig(ctx, lgc.header, dat)
	if nil != err {
		blog.Errorf("getHostByModuleID GetModulesHostConfig http do error. moduleID %d supplierID %s  error:%s,input:%+v,rid:%s", moduleID, supplierID, err.Error(), dat, lgc.rid)
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !ret.Result {
		blog.Errorf("getHostByModuleID  GetModulesHostConfig http reply error.moduleID %d supplierID %s err code:%d,err msg:%s,input:%+v,rid:%s", moduleID, supplierID, ret.Code, ret.ErrMsg, dat, lgc.rid)
		return nil, defErr.New(ret.Code, ret.ErrMsg)
	}
	if 0 == len(ret.Data) {
		blog.V(5).Infof("getHostByModuleID moduleID %d supplierID %s GetModulesHostConfig len equal 0,input:%+v,rid:%s", moduleID, supplierID, dat, lgc.rid)
		return nil, nil
	}
	var hostIDs []int64
	for _, item := range ret.Data {
		hostIDs = append(hostIDs, item.HostID)
	}
	opt := new(metadata.QueryInput)
	opt.Condition = mapstr.MapStr{common.BKHostIDField: common.KvMap{common.BKDBIN: hostIDs}}
	opt.Fields = fmt.Sprintf("%s,%s,%s", common.BKHostIDField, common.BKHostInnerIPField, common.BKCloudIDField)
	opt.Limit = common.BKNoLimit
	hosts, err := lgc.CoreAPI.HostController().Host().GetHosts(ctx, lgc.header, opt)
	if nil != err {
		blog.Errorf("getHostByModuleID GetHosts http do error.moduleID %d hostID:%v supplierID %s GetHosts http do error:%s,input:%+v,rid:%s", moduleID, hostIDs, supplierID, err.Error(), opt, lgc.rid)
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !hosts.Result {
		blog.Errorf("getHostByModuleID GetHosts http reply error. moduleID %d hostID:%v supplierID %s GetHosts http reply error:%s,input:%+v,rid:%s", moduleID, hostIDs, supplierID, hosts.ErrMsg, opt, lgc.rid)
		return nil, defErr.New(hosts.Code, hosts.ErrMsg)
	}

	hostInfos := make(map[int64]*metadata.GseHost, len(hosts.Data.Info))
	for _, host := range hosts.Data.Info {
		item := new(metadata.GseHost)

		hostID, err := util.GetInt64ByInterface(host[common.BKHostIDField])
		if nil != err {
			blog.Errorf("getHostByModuleID hostInfo %+v  hostID   not interger,rid:%s", host, lgc.rid)
			return nil, err
		}
		cloudID, err := util.GetInt64ByInterface(host[common.BKCloudIDField])
		if nil != err {
			byteHost, _ := json.Marshal(host)
			blog.Errorf("getHostByModuleID  hostInfo %v  cloudID  not interger, host:%s,rid:%s", host, string(byteHost), lgc.rid)
			return nil, err
		}
		innerIP, ok := host[common.BKHostInnerIPField].(string)
		if !ok {
			byteHost, _ := json.Marshal(host)
			blog.Errorf("getHostByModuleID  hostInfo %v  innerip  not found, host:%s,rid:%s", host, string(byteHost), lgc.rid)
			return nil, err
		}
		item.HostID = hostID
		item.BkCloudId = cloudID
		item.Ip = innerIP
		item.BkSupplierId = intSupplierID

		hostInfos[hostID] = item
	}

	return hostInfos, nil
}
