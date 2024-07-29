/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package core

import (
	"errors"
	"fmt"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/querybuilder"
	"configcenter/src/common/util"
)

// GetCloudArea search total cloud area id and name return an array of cloud name and a name-id map
func (d *Client) GetCloudArea(kit *rest.Kit, names ...string) ([]string, map[string]int64, error) {
	cloudArea := make([]mapstr.MapStr, 0)
	input := metadata.CloudAreaSearchParam{
		SearchCloudOption: metadata.SearchCloudOption{
			Fields: []string{common.BKCloudIDField, common.BKCloudNameField},
			Page:   metadata.BasePage{Start: 0, Limit: common.BKMaxPageSize},
		},
	}
	if len(names) != 0 {
		input.SearchCloudOption.Condition = mapstr.MapStr{
			common.BKCloudNameField: mapstr.MapStr{common.BKDBIN: names},
		}
	}
	for {
		rsp, err := d.ApiClient.SearchCloudArea(kit.Ctx, kit.Header, input)
		if err != nil {
			blog.Errorf("search cloud area failed, err: %v, rid: %s", err, kit.Rid)
			return nil, nil, err
		}

		cloudArea = append(cloudArea, rsp.Info...)
		if len(rsp.Info) < common.BKMaxPageSize {
			break
		}

		input.SearchCloudOption.Page.Start += common.BKMaxPageSize
	}

	if len(cloudArea) == 0 {
		blog.Errorf("search cloud area failed, return empty, rid: %s", kit.Rid)
		return nil, nil, kit.CCError.CCError(common.CCErrTopoCloudNotFound)
	}

	cloudAreaArr := make([]string, 0)
	cloudAreaMap := make(map[string]int64)
	for _, item := range cloudArea {
		areaName, err := item.String(common.BKCloudNameField)
		if err != nil {
			blog.Errorf("get type of string cloud name failed, err: %v, rid: %s", err, kit.Rid)
			return nil, nil, err
		}
		cloudAreaArr = append(cloudAreaArr, areaName)

		areaID, err := item.Int64(common.BKCloudIDField)
		if err != nil {
			blog.Errorf("get type of int64 cloud id failed, err: %v, rid: %s", err, kit.Rid)
			return nil, nil, err
		}
		// cloud area name is unique
		cloudAreaMap[areaName] = areaID
	}

	return cloudAreaArr, cloudAreaMap, nil
}

// GetHost get host instance
func (d *Client) GetHost(kit *rest.Kit, cond interface{}) ([]mapstr.MapStr, error) {
	hostCond, ok := cond.(mapstr.MapStr)
	if !ok {
		blog.Errorf("get host but condition parse failed, condition: %v, rid: %s", cond, kit.Rid)
		return nil, errors.New("get host but condition parse failed")
	}
	result, err := d.ApiClient.GetHostData(kit.Ctx, kit.Header, hostCond)
	if err != nil {
		blog.Errorf("get host failed, condition: %+v, err: %+v, rid: %s", cond, err, kit.Rid)
		return nil, err
	}

	if err := result.CCError(); err != nil {
		blog.Errorf("get host failed, condition: %+v, err: %+v, rid: %s", cond, err, kit.Rid)
		return nil, err
	}

	hosts := result.Data.Info

	if len(hosts) == 0 {
		blog.Errorf("not find host, cond: %#v, rid: %s", cond, kit.Rid)
		return nil, nil
	}

	// 此函数需要放在handleMainlineTopo之前，因为在handleMainlineTopo函数中会调整集群数据
	setInfo, err := d.getHostSetInfo(kit, hosts)
	if err != nil {
		blog.Errorf("get host set info failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	hosts, err = handleMainlineTopo(kit, hosts)
	if err != nil {
		blog.Errorf("handle host biz and module name failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	objIDs, err := d.getCustomTopoObjIDs(kit)
	if err != nil {
		blog.Errorf("get custom topo objID failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	if len(objIDs) != 0 {
		hosts, err = d.handleCustomTopo(kit, hosts, setInfo, objIDs)
		if err != nil {
			blog.Errorf("handle custom topo data failed, err: %v, rid: %s", err, kit.Rid)
			return nil, err
		}
	}

	topoObjIDs := []string{TopoObjID, common.BKInnerObjIDApp}
	topoObjIDs = append(topoObjIDs, objIDs...)
	topoObjIDs = append(topoObjIDs, []string{common.BKInnerObjIDSet, common.BKInnerObjIDModule}...)

	hosts, err = handleHostResult(kit, hosts, topoObjIDs)
	if err != nil {
		blog.Errorf("handle host result failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	return hosts, nil
}

func handleMainlineTopo(kit *rest.Kit, hosts []mapstr.MapStr) ([]mapstr.MapStr, error) {
	var err error
	for idx, host := range hosts {
		hosts[idx], err = handleData(kit, host, common.BKInnerObjIDApp)
		if err != nil {
			return nil, err
		}

		hosts[idx], err = handleData(kit, host, common.BKInnerObjIDSet)
		if err != nil {
			return nil, err
		}

		// add topo path data
		moduleMap, ok := host[common.BKInnerObjIDModule].([]interface{})
		if ok {
			topos := util.GetStrValsFromArrMapInterfaceByKey(moduleMap, common.TopoModuleName)
			if len(topos) > 0 {
				hosts[idx][TopoObjID] = strings.Join(topos, ", ")
			}
		}

		hosts[idx], err = handleData(kit, host, common.BKInnerObjIDModule)
		if err != nil {
			return nil, err
		}
	}

	return hosts, nil
}

func handleData(kit *rest.Kit, host mapstr.MapStr, objID string) (mapstr.MapStr, error) {
	objMap, exist := host[objID].([]interface{})
	if !exist {
		blog.Errorf("get %s map data from host data failed, not exist, data: %#v, rid: %s", objID, host, kit.Rid)
		return nil, fmt.Errorf("get %s map data from host data failed, not exist, data: %#v", objID, host)
	}

	var nameStr string
	for _, obj := range objMap {
		rowMap, err := mapstr.NewFromInterface(obj)
		if err != nil {
			blog.Errorf("get %s data from host data failed, err: %v, rid: %s", objID, err, kit.Rid)
			return nil, err
		}

		name, err := rowMap.String(common.GetInstNameField(objID))
		if err != nil {
			blog.Errorf("get %s name from host data failed, err: %v, rid: %s", objID, err, kit.Rid)
			return nil, fmt.Errorf("get %s name from host data failed, err: %v", obj, err)
		}

		if nameStr == "" {
			nameStr = name
			continue
		}

		nameStr += "," + name
	}

	host.Set(objID, nameStr)

	return host, nil
}

func (d *Client) getHostSetInfo(kit *rest.Kit, hosts []mapstr.MapStr) (*hostSetInfo, error) {
	res, err := d.ApiClient.SearchPlatformSetting(kit.Ctx, kit.Header, "current")
	if err != nil {
		return nil, err
	}

	hostSetMap := make(map[int64][]int64, 0)
	allSetIDs := make([]int64, 0)
	for _, data := range hosts {
		rowMap, err := mapstr.NewFromInterface(data[common.BKInnerObjIDHost])
		if err != nil {
			blog.Errorf("get host map data failed, hostData: %#v, err: %v, rid: %s", data, err, kit.Rid)
			return nil, err
		}

		hostID, err := rowMap.Int64(common.BKHostIDField)
		if err != nil {
			blog.Errorf("get host id failed, host id: %s, err: %v, rid: %s", hostID, err, kit.Rid)
			return nil, err
		}

		setMap, exist := data[common.BKInnerObjIDSet].([]interface{})
		if !exist {
			blog.Errorf("get set map data from host data, not exist, data: %#v, rid: %s", data, kit.Rid)
			return nil, errors.New("from host data get set map, not exist")
		}

		hostSetIds := make([]int64, 0)
		for _, set := range setMap {
			rowMap, err := mapstr.NewFromInterface(set)
			if err != nil {
				blog.Errorf("get set data from host data failed, err: %v, rid: %s", err, kit.Rid)
				return nil, err
			}

			setName, err := rowMap.String(common.BKSetNameField)
			if err != nil {
				blog.Errorf("get set name from host data failed, err: %v, rid: %s", err, kit.Rid)
				return nil, fmt.Errorf("from host data get set name, not exist, rid: %s", kit.Rid)
			}

			setID, err := rowMap.Int64(common.BKSetIDField)
			if err != nil {
				blog.Errorf("get set id from host data failed, err: %v, rid: %s", err, kit.Rid)
				return nil, err
			}

			if setName != string(res.Data.BuiltInSetName) {
				allSetIDs = append(allSetIDs, setID)
				hostSetIds = append(hostSetIds, setID)
			}

		}
		hostSetMap[hostID] = hostSetIds
	}

	return &hostSetInfo{setIDs: allSetIDs, hostSetMap: hostSetMap}, nil
}

func (d *Client) handleCustomTopo(kit *rest.Kit, hosts []mapstr.MapStr, hostSetInfo *hostSetInfo, objIDs []string) (
	[]mapstr.MapStr, error) {

	setInfo, err := d.getTopoInstData(kit, hostSetInfo.setIDs, common.BKInnerObjIDSet)
	if err != nil {
		return nil, err
	}

	parentIDs := setInfo.parentIDs
	instIDParentIDMap := setInfo.instIdParentIDMap
	hostTopoIDMap := hostSetInfo.hostSetMap

	for _, objID := range objIDs {
		parentData, err := d.getTopoInstData(kit, parentIDs, objID)

		if err != nil {
			blog.Errorf("get topo instance data failed, cond: %#v, err: %v, rid: %s", parentIDs, err, kit.Rid)
			return nil, err
		}

		hostIDMap := make(map[int64][]int64, 0)
		hostTopoNameMap := make(map[int64]string, 0)

		for hostID, topoIDs := range hostTopoIDMap {
			nameStr := ""
			for _, topoID := range topoIDs {
				hostIDMap[hostID] = append(hostIDMap[hostID], instIDParentIDMap[topoID])
				parentID := instIDParentIDMap[topoID]
				parentName := parentData.instIdNameMap[parentID]
				if nameStr == "" {
					nameStr = parentName
					continue
				}

				nameStr += "," + parentName
			}
			hostTopoNameMap[hostID] = nameStr
		}

		for _, host := range hosts {
			rowMap, err := mapstr.NewFromInterface(host[common.BKInnerObjIDHost])
			if err != nil {
				blog.Errorf("get host map data failed, hostData: %#v, err: %v, rid: %s", host, err, kit.Rid)
				return nil, err
			}

			hostID, err := rowMap.Int64(common.BKHostIDField)
			if err != nil {
				blog.Errorf("get host id failed, host id: %s, err: %v, rid: %s", hostID, err, kit.Rid)
				return nil, err
			}

			host[objID] = hostTopoNameMap[hostID]
		}

		instIDParentIDMap = parentData.instIdParentIDMap
		hostTopoIDMap = hostIDMap
		parentIDs = parentData.parentIDs
	}

	return hosts, nil
}

func (d *Client) getTopoInstData(kit *rest.Kit, instIDs []int64, objID string) (*topoInstData, error) {
	idField := common.GetInstIDField(objID)
	nameField := common.GetInstNameField(objID)

	query := &metadata.QueryCondition{
		Condition: mapstr.MapStr{idField: mapstr.MapStr{common.BKDBIN: instIDs}},
		Fields:    []string{idField, nameField, common.BKInstParentStr},
	}

	insts, err := d.ApiClient.ReadInstance(kit.Ctx, kit.Header, objID, query)
	if err != nil {
		blog.Errorf("get custom level inst data failed, query cond: %#v, err: %v, rid: %s", query, err, kit.Rid)
		return nil, err
	}

	parentIDs := make([]int64, 0)
	instIdParentIdMap := make(map[int64]int64, 0)
	instIdNameMap := make(map[int64]string, 0)

	for _, inst := range insts.Data.Info {
		parentID, err := inst.Int64(common.BKParentIDField)
		if err != nil {
			blog.Errorf("get inst parent id failed, err: %v, rid: %s", err, kit.Rid)
			return nil, err
		}
		parentIDs = append(parentIDs, parentID)

		instID, err := inst.Int64(idField)
		if err != nil {
			blog.Errorf("get inst id failed, err: %v, rid: %s", err, kit.Rid)
			return nil, err
		}
		instIdParentIdMap[instID] = parentID

		instName, err := inst.String(nameField)
		if err != nil {
			blog.Errorf("get inst name failed, err: %v, rid: %s", err, kit.Rid)
			return nil, err
		}
		instIdNameMap[instID] = instName
	}

	return &topoInstData{parentIDs: parentIDs, instIdParentIDMap: instIdParentIdMap, instIdNameMap: instIdNameMap}, nil
}

// handleHostResult 原数据是像{"host":{}, "set":{}}这样的结构，现在将主机这一层的数据拿出来，把拓扑相关的数据也作为主机的字段
func handleHostResult(kit *rest.Kit, hostsWithTopo []mapstr.MapStr, topoObjIDs []string) ([]mapstr.MapStr, error) {
	result := make([]mapstr.MapStr, len(hostsWithTopo))
	for idx, data := range hostsWithTopo {
		host, err := mapstr.NewFromInterface(data[common.BKInnerObjIDHost])
		if err != nil {
			blog.Errorf("get host data failed, data: %#v, err: %v, rid: %s", data, err, kit.Rid)
			return nil, err
		}

		for _, objID := range topoObjIDs {
			host[IDPrefix+objID] = data[objID]
		}

		result[idx] = host
	}

	return result, nil
}

// GetSameIPRes get hosts that same ip in db
func (d *Client) GetSameIPRes(kit *rest.Kit, hostInfos map[int]map[string]interface{}) (*SameIPRes, error) {
	result := &SameIPRes{V4Map: map[string]struct{}{}, V6Map: map[string]struct{}{}}

	// step1. extract all innerIP from hostInfos
	innerIPs, innerIPv6s := make([]string, 0), make([]string, 0)
	for _, host := range hostInfos {
		innerIP, ok := host[common.BKHostInnerIPField].(string)
		if ok && innerIP != "" {
			innerIPs = append(innerIPs, strings.Split(innerIP, ",")...)
		}
		innerIPv6, ok := host[common.BKHostInnerIPv6Field].(string)
		if ok && innerIPv6 != "" {
			innerIPv6s = append(innerIPv6s, strings.Split(innerIPv6, ",")...)
		}
	}
	if len(innerIPs) == 0 && len(innerIPv6s) == 0 {
		return result, nil
	}

	// step2. query host info by innerIPs
	rules := make([]querybuilder.Rule, 0)
	if len(innerIPs) > 0 {
		rules = append(rules, querybuilder.AtomRule{
			Field: common.BKHostInnerIPField, Operator: querybuilder.OperatorIn, Value: innerIPs,
		})
	}
	if len(innerIPv6s) > 0 {
		rules = append(rules, querybuilder.AtomRule{
			Field: common.BKHostInnerIPv6Field, Operator: querybuilder.OperatorIn, Value: innerIPv6s,
		})
	}

	option := metadata.ListHostsWithNoBizParameter{
		HostPropertyFilter: &querybuilder.QueryFilter{
			Rule: querybuilder.CombinedRule{Condition: querybuilder.ConditionOr, Rules: rules},
		},
		Fields: []string{common.BKHostIDField, common.BKHostInnerIPField, common.BKHostInnerIPv6Field,
			common.BKCloudIDField},
	}
	resp, err := d.ApiClient.ListHostWithoutApp(kit.Ctx, kit.Header, option)
	if err != nil {
		blog.Errorf("list host without app failed, option: %+v, err: %v, rid: %s", option, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if !resp.Result {
		blog.Errorf("ListHostWithoutApp resp:%#v, option: %d, rid: %s", resp, option, kit.Rid)
		return nil, resp.CCError()
	}

	// step3. arrange data as a map, cloudKey: hostID
	for _, host := range resp.Data.Info {
		ipv4 := util.GetStrByInterface(host[common.BKHostInnerIPField])
		if ipv4 != "" {
			keyV4 := HostCloudKey(ipv4, host[common.BKCloudIDField])
			result.V4Map[keyV4] = struct{}{}
		}

		ipv6 := util.GetStrByInterface(host[common.BKHostInnerIPv6Field])
		if ipv6 != "" {
			keyV6 := HostCloudKey(ipv6, host[common.BKCloudIDField])
			result.V6Map[keyV6] = struct{}{}
		}
	}

	return result, nil
}

// HostCloudKey generate a cloudKey for host that is unique among clouds by appending the cloudID.
func HostCloudKey(ip, cloudID interface{}) string {
	return fmt.Sprintf("%v-%v", ip, cloudID)
}

// GetExistingHost get existing host
func (d *Client) GetExistingHost(kit *rest.Kit, hosts map[int]map[string]interface{}) (map[int64]SimpleHost, error) {
	// step1. extract all innerIP from hostInfos
	var hostIDs []int64
	for _, host := range hosts {
		hostID, ok := host[common.BKHostIDField]
		if !ok {
			blog.Errorf("host can not find %s field, data: %v, rid: %s", common.BKHostIDField, host, kit.Rid)
			return nil, fmt.Errorf("host can not find %s field, data: %v", common.BKHostIDField, host)
		}

		hostIDVal, err := util.GetInt64ByInterface(hostID)
		if err != nil {
			blog.Errorf("host %s field is invalid, value: %v, rid: %s", common.BKHostIDField, hostID, kit.Rid)
			return nil, err
		}
		hostIDs = append(hostIDs, hostIDVal)
	}

	if len(hostIDs) == 0 {
		return make(map[int64]SimpleHost), nil
	}

	// step2. query host info by hostIDs
	rules := []querybuilder.Rule{
		querybuilder.AtomRule{Field: common.BKHostIDField, Operator: querybuilder.OperatorIn, Value: hostIDs},
	}
	option := metadata.ListHostsWithNoBizParameter{
		HostPropertyFilter: &querybuilder.QueryFilter{
			Rule: querybuilder.CombinedRule{Condition: querybuilder.ConditionOr, Rules: rules},
		},
		Fields: []string{
			common.BKHostIDField,
			common.BKHostInnerIPField,
			common.BKHostInnerIPv6Field,
			common.BKAgentIDField,
		},
	}
	resp, err := d.ApiClient.ListHostWithoutApp(kit.Ctx, kit.Header, option)
	if err != nil {
		blog.Errorf("list host without app failed, err: %v, option: %v, rid: %s", err, option, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if !resp.Result {
		blog.Errorf("list host without app failed, err: %v, option: %v, rid: %s", err, option, kit.Rid)
		return nil, resp.CCError()
	}

	// step3. arrange data as a map, cloudKey: hostID
	hostMap := make(map[int64]SimpleHost, 0)
	for _, host := range resp.Data.Info {
		hostID, ok := host[common.BKHostIDField]
		if !ok {
			blog.Errorf("host can not find %s field, data: %v, rid: %s", common.BKHostIDField, host, kit.Rid)
			return nil, fmt.Errorf("host can not find %s field, data: %v", common.BKHostIDField, host)
		}

		hostIDVal, err := util.GetInt64ByInterface(hostID)
		if err != nil {
			blog.Errorf("host %s field is invalid, value: %v, rid: %s", common.BKHostIDField, hostID, kit.Rid)
			return nil, err
		}

		hostMap[hostIDVal] = SimpleHost{
			Ip:      util.GetStrByInterface(host[common.BKHostInnerIPField]),
			Ipv6:    util.GetStrByInterface(host[common.BKHostInnerIPv6Field]),
			AgentID: util.GetStrByInterface(host[common.BKAgentIDField]),
		}
	}

	return hostMap, nil
}
