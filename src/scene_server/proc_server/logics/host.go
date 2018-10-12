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
		blog.V(3).Infof("getHostByModuleID moduleID %d supplierID %s GetModulesHostConfig len equal 0", moduleID, supplierID)
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
		item.HostID = hostID
		item.BkCloudId = cloudID
		item.Ip = innerIP
		item.BkSupplierId = intSupplierID

		hostInfos[hostID] = item
	}

	return hostInfos, nil
}
