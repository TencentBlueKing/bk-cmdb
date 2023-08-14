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

package db

import (
	"bytes"
	errs "errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"configcenter/pkg/filter"
	filtertools "configcenter/pkg/tools/filter"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/resource/esb"
	"configcenter/src/common/util"
	"configcenter/src/web_server/middleware/user/plugins"
)

// Dao used to process excel-related data
type Dao struct {
	Engine *backbone.Engine
}

// GetSortedColProp get sort column property
func (d *Dao) GetSortedColProp(kit *rest.Kit, cond mapstr.MapStr) ([]ColProp, error) {
	colProps, err := d.getObjColProp(kit, cond)
	if err != nil {
		blog.Errorf("get object column property failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	objID, err := cond.String(common.BKObjIDField)
	if err != nil {
		blog.Errorf("get objID from condition failed, instCond: %v, err: %v, rid: %s", cond, err, kit.Rid)
		return nil, err
	}

	filterProp := getFilterProp(objID)
	for idx := range colProps {
		if util.InStrArr(filterProp, colProps[idx].ID) {
			colProps[idx].NotExport = true
		}
	}

	bizID, err := cond.Int64(common.BKAppIDField)
	if err != nil {
		bizID = 0
	}
	groups, err := d.getObjGroup(kit, objID, bizID)
	if err != nil {
		blog.Errorf("get object attribute group failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	colProps, err = sortColProp(colProps, groups)
	if err != nil {
		blog.Errorf("sort column property failed, column property: %v, attribute group: %v, err: %v, rid: %s",
			colProps, groups, err, kit.Rid)
		return nil, err
	}

	return colProps, nil
}

// getFilterProp 不需要展示字段
func getFilterProp(objID string) []string {
	switch objID {
	case common.BKInnerObjIDHost:
		return []string{common.BKSetNameField, common.BKModuleNameField, common.BKAppNameField}
	default:
		return []string{common.CreateTimeField}
	}
}

// getObjColProp get object column property
func (d *Dao) getObjColProp(kit *rest.Kit, cond mapstr.MapStr) ([]ColProp, error) {
	attrs, err := d.Engine.CoreAPI.ApiServer().ModelQuote().GetObjectAttrWithTable(kit.Ctx, kit.Header, cond)
	if err != nil {
		blog.Errorf("get object fields failed, condition: %v, err: %v ,rid: %s", cond, err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	result := make([]ColProp, len(attrs))
	for idx, attr := range attrs {
		colProp := ColProp{ID: attr.PropertyID, Name: attr.PropertyName, PropertyType: attr.PropertyType,
			IsRequire: attr.IsRequired, Option: attr.Option, Group: attr.PropertyGroup, RefSheet: attr.PropertyName,
		}
		result[idx] = colProp
	}

	return result, nil
}

// getObjGroup get object property group
func (d *Dao) getObjGroup(kit *rest.Kit, objID string, bizID int64) ([]metadata.AttributeGroup, error) {
	cond := mapstr.MapStr{
		common.BKObjIDField: objID,
		common.BKAppIDField: bizID,
		metadata.PageName: mapstr.MapStr{metadata.PageStart: 0, metadata.PageLimit: common.BKNoLimit,
			metadata.PageSort: common.BKPropertyGroupIndexField},
	}

	result, err := d.Engine.CoreAPI.ApiServer().GetObjectGroup(kit.Ctx, kit.Header, kit.SupplierAccount, objID, cond)
	if err != nil {
		blog.Errorf("get %s fields group failed, err:%+v, rid: %s", objID, err, kit.Rid)
		return nil, fmt.Errorf("get attribute group failed, err: %+v", err)
	}

	if !result.Result {
		blog.Errorf("get %s fields group result failed. code: %d, message: %s, rid: %s", objID, result.Code,
			result.ErrMsg, kit.Rid)

		return nil, fmt.Errorf("get attribute group result false, result: %+v", result)
	}

	return result.Data, nil
}

// GetCloudArea search total cloud area id and name return an array of cloud name and a name-id map
func (d *Dao) GetCloudArea(kit *rest.Kit) ([]string, map[string]int64, error) {
	cloudArea := make([]mapstr.MapStr, 0)
	start := 0
	for {
		input := metadata.CloudAreaSearchParam{
			SearchCloudOption: metadata.SearchCloudOption{
				Fields: []string{common.BKCloudIDField, common.BKCloudNameField},
				Page:   metadata.BasePage{Start: start, Limit: common.BKMaxPageSize},
			},
		}
		rsp, err := d.Engine.CoreAPI.ApiServer().SearchCloudArea(kit.Ctx, kit.Header, input)
		if err != nil {
			blog.Errorf("search cloud area failed, err: %v, rid: %s", err, kit.Rid)
			return nil, nil, err
		}

		cloudArea = append(cloudArea, rsp.Info...)
		if len(rsp.Info) < common.BKMaxPageSize {
			break
		}

		start += common.BKMaxPageSize
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

// GetInst get instance
func (d *Dao) GetInst(kit *rest.Kit, objID string, cond mapstr.MapStr) ([]mapstr.MapStr, error) {
	result, err := d.Engine.CoreAPI.ApiServer().GetInstDetail(kit.Ctx, kit.Header, objID, cond)
	if nil != err {
		blog.Errorf("get inst data detail error: %v , search condition: %#v, rid: %s", err, cond, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !result.Result {
		blog.Errorf("get inst data detail error: %v ,condition: %#v, rid: %s", result.ErrMsg, cond, kit.Rid)
		return nil, kit.CCError.Error(result.Code)
	}

	if result.Data.Count == 0 {
		blog.Errorf("get inst data detail, but got 0 instances, condition: %#v, rid: %s", cond, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrAPINoObjectInstancesIsFound)
	}

	return result.Data.Info, nil
}

// GetHost get host instance
func (d *Dao) GetHost(kit *rest.Kit, cond mapstr.MapStr) ([]mapstr.MapStr, error) {
	result, err := d.Engine.CoreAPI.ApiServer().GetHostData(kit.Ctx, kit.Header, cond)
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

	if len(objIDs) == 0 {
		return hosts, nil
	}

	hosts, err = d.handleCustomTopo(kit, hosts, setInfo, objIDs)
	if err != nil {
		blog.Errorf("handle custom topo data failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
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

type hostSetInfo struct {
	setIDs     []int64
	hostSetMap map[int64][]int64
}

func (d *Dao) getHostSetInfo(kit *rest.Kit, hosts []mapstr.MapStr) (*hostSetInfo, error) {
	res, err := d.Engine.CoreAPI.CoreService().System().SearchPlatformSetting(kit.Ctx, kit.Header)
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
			return nil, errs.New("from host data get set map, not exist")
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

type topoInstData struct {
	parentIDs         []int64
	instIdParentIDMap map[int64]int64
	instIdNameMap     map[int64]string
}

func (d *Dao) handleCustomTopo(kit *rest.Kit, hosts []mapstr.MapStr, hostSetInfo *hostSetInfo, objIDs []string) (
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
			names := make([]string, 0)

			for _, topoID := range topoIDs {
				parentID := instIDParentIDMap[topoID]
				parentName := parentData.instIdNameMap[parentID]
				names = append(names, parentName)

				hostIDMap[hostID] = append(hostIDMap[hostID], instIDParentIDMap[topoID])
			}

			nameStr := ""
			for _, name := range names {
				if nameStr == "" {
					nameStr = name
					continue
				}

				nameStr += "," + name
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
	}

	return hosts, nil
}

func (d *Dao) getTopoInstData(kit *rest.Kit, instIDs []int64, objID string) (*topoInstData, error) {
	idField := common.GetInstIDField(objID)
	nameField := common.GetInstNameField(objID)

	query := &metadata.QueryCondition{
		Condition: mapstr.MapStr{idField: mapstr.MapStr{common.BKDBIN: instIDs}},
		Fields:    []string{idField, nameField, common.BKInstParentStr},
	}

	insts, err := d.Engine.CoreAPI.ApiServer().ReadInstance(kit.Ctx, kit.Header, objID, query)
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

// HandleEnumQuoteInst handle enum quote instance
func (d *Dao) HandleEnumQuoteInst(kit *rest.Kit, infos []mapstr.MapStr, colProps []ColProp) (
	[]mapstr.MapStr, error) {

	for _, rowMap := range infos {
		if err := d.getEnumQuoteInstNames(kit, colProps, rowMap); err != nil {
			blog.Errorf("get enum quote inst name list failed, err: %v, rid: %s", err, kit.Rid)
			return nil, err
		}
	}

	return infos, nil
}

// getEnumQuoteInstNames search inst detail and return a bk_inst_name
func (d *Dao) getEnumQuoteInstNames(kit *rest.Kit, colProps []ColProp, rowMap mapstr.MapStr) error {
	for _, property := range colProps {
		if property.PropertyType != common.FieldTypeEnumQuote {
			continue
		}

		enumQuoteIDInterface, exist := rowMap[property.ID]
		if !exist || enumQuoteIDInterface == nil {
			continue
		}

		enumQuoteIDList, ok := enumQuoteIDInterface.([]interface{})
		if !ok {
			blog.Errorf("rowMap[%s] type to array failed, rowMap: %v, rowMap type: %T, rid: %s", property,
				rowMap[property.ID], rowMap[property.ID], kit.Rid)
			return fmt.Errorf("convert variable rowMap[%v] type to int array failed", property)
		}

		enumQuoteIDMap := make(map[int64]interface{}, 0)
		for _, enumQuoteID := range enumQuoteIDList {
			id, err := util.GetInt64ByInterface(enumQuoteID)
			if err != nil {
				blog.Errorf("convert enumQuoteID[%d] to int64 failed, type: %T, err: %v, rid: %s", enumQuoteID,
					enumQuoteID, kit.Rid)
				return err
			}

			if id == 0 {
				return fmt.Errorf("enum quote instID is %d, it is illegal", id)
			}

			enumQuoteIDMap[id] = struct{}{}
		}

		if len(enumQuoteIDMap) == 0 {
			continue
		}

		quoteObjID, err := getEnumQuoteObjID(kit, property.Option)
		if err != nil {
			blog.Errorf("get enum quote option obj id failed, err: %s, rid: %s", err, kit.Rid)
			return fmt.Errorf("get enum quote option obj id failed, option: %v", property.Option)
		}

		enumQuoteIDs := make([]int64, 0)
		for enumQuoteID := range enumQuoteIDMap {
			enumQuoteIDs = append(enumQuoteIDs, enumQuoteID)
		}
		input := &metadata.QueryCondition{
			Fields:         []string{common.GetInstNameField(quoteObjID)},
			Condition:      mapstr.MapStr{common.GetInstIDField(quoteObjID): mapstr.MapStr{common.BKDBIN: enumQuoteIDs}},
			DisableCounter: true,
		}
		resp, err := d.Engine.CoreAPI.ApiServer().ReadInstance(kit.Ctx, kit.Header, quoteObjID, input)
		if err != nil {
			blog.Errorf("get quote inst name list failed, input: %+v, err: %v, rid: %s", input, err, kit.Rid)
			return err
		}

		enumQuoteNames := make([]string, 0)
		for _, info := range resp.Data.Info {
			var ok bool
			var enumQuoteName string
			if name, exist := info.Get(common.GetInstNameField(quoteObjID)); exist {
				enumQuoteName, ok = name.(string)
				if !ok {
					enumQuoteName = ""
				}
			}
			enumQuoteNames = append(enumQuoteNames, enumQuoteName)
		}

		rowMap[property.ID] = strings.Join(enumQuoteNames, "\n")
	}

	return nil
}

// getEnumQuoteObjID get enum quote field option bk_obj_id and bk_inst_id value
func getEnumQuoteObjID(kit *rest.Kit, option interface{}) (string, error) {
	if option == nil {
		return "", fmt.Errorf("enum quote option is nil")
	}

	arrOption, ok := option.([]interface{})
	if !ok {
		blog.Errorf("option %v not enum quote option, rid: %s", option, kit.Rid)
		return "", fmt.Errorf("enum quote option is unvalid")
	}

	for _, o := range arrOption {
		mapOption, ok := o.(map[string]interface{})
		if !ok || mapOption == nil {
			blog.Errorf("option %v not enum quote option, enum quote option item must bk_obj_id, rid: %s", option,
				kit.Rid)
			return "", fmt.Errorf("convert option map[string]interface{} failed")
		}

		objIDVal, objIDOk := mapOption[common.BKObjIDField]
		if !objIDOk || objIDVal == "" {
			blog.Errorf("enum quote option bk_obj_id can't be empty, rid: %s", option, kit.Rid)
			return "", fmt.Errorf("enum quote option bk_obj_id can't be empty")
		}

		objID, ok := objIDVal.(string)
		if !ok {
			blog.Errorf("objIDVal %v not string, rid: %s", objIDVal, kit.Rid)
			return "", fmt.Errorf("enum quote option bk_obj_id is not string")
		}

		return objID, nil
	}

	return "", nil
}

// GetUsernameMapWithPropertyList 依照"bk_obj_id"和"bk_property_type":"objuser"查询"cc_ObjAttDes"集合,得到"bk_property_id"
// 的值; 然后以它的值为key,取得Info中的value,然后以value作为param访问ESB,得到其中文名。
func (d *Dao) GetUsernameMapWithPropertyList(kit *rest.Kit, objID string, infoList []mapstr.MapStr) (map[string]string,
	[]string, error) {

	cond := metadata.QueryCondition{
		Fields: []string{metadata.AttributeFieldPropertyID},
		Condition: map[string]interface{}{
			metadata.AttributeFieldObjectID:     objID,
			metadata.AttributeFieldPropertyType: common.FieldTypeUser,
		},
	}
	attrRsp, err := d.Engine.CoreAPI.CoreService().Model().ReadModelAttr(kit.Ctx, kit.Header, objID, &cond)
	if err != nil {
		blog.Errorf("failed to request the object controller, err: %v, rid: %s", err, kit.Rid)
		return nil, nil, err
	}

	usernameList := []string{}
	propertyList := []string{}
	for _, info := range infoList {
		for _, item := range attrRsp.Info {
			propertyList = append(propertyList, item.PropertyID)
			if info[item.PropertyID] != nil {
				username, ok := info[item.PropertyID].(string)
				if !ok {
					err = fmt.Errorf("failed to cast %s instance from interface{} to string, rid: %s", objID, kit.Rid)
					blog.Errorf("failed to cast %s instance from interface{} to string", objID)
					return nil, nil, err
				}
				usernameList = append(usernameList, strings.Split(username, ",")...)
			}
		}
	}
	propertyList = util.RemoveDuplicatesAndEmpty(propertyList)
	userList := util.RemoveDuplicatesAndEmpty(usernameList)
	// get username from esb
	usernameMap, err := d.getUsernameFromEsb(kit, userList)
	if err != nil {
		blog.ErrorJSON("get username map from ESB failed, err: %s, rid: %s", err.Error(), kit.Rid)
		return nil, nil, err
	}

	return usernameMap, propertyList, nil
}

func (d *Dao) getUsernameFromEsb(kit *rest.Kit, userList []string) (map[string]string, error) {

	defErr := d.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(kit.Header))
	usernameMap := map[string]string{}

	if len(userList) == 0 {
		return usernameMap, nil
	}

	loginVersion, _ := cc.String("webServer.login.version")
	user := plugins.CurrentPlugin(loginVersion)

	// 处理请求的用户数据，将用户拼接成不超过500字节的字符串进行用户数据的获取
	userListStr := getUserListStr(userList)

	var wg sync.WaitGroup
	var lock sync.RWMutex
	var firstErr errors.CCErrorCoder
	pipeline := make(chan bool, 10)
	userListEsb := make([]*metadata.LoginSystemUserInfo, 0)

	for _, subStr := range userListStr {
		pipeline <- true
		wg.Add(1)
		go func(subStr string) {
			defer func() {
				wg.Done()
				<-pipeline
			}()

			lock.Lock()
			params := make(map[string]string)
			params["fields"] = "username,display_name"
			params["exact_lookups"] = subStr
			lock.Unlock()

			// todo 下面为原函数逻辑，但是这里没有 *gin.Context参数可以拿到，需要区分场景调用查询用户的方法
			userListEsbSub, errNew := user.GetUserList(nil, params)
			if errNew != nil {
				firstErr = errNew.ToCCError(defErr)
				blog.Errorf("get users(%s) list from ESB failed, err: %v, rid: %s", subStr, firstErr, kit.Rid)
				return
			}

			lock.Lock()
			userListEsb = append(userListEsb, userListEsbSub...)
			lock.Unlock()
		}(subStr)
	}
	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}

	for _, userInfo := range userListEsb {
		username := fmt.Sprintf("%s(%s)", userInfo.EnName, userInfo.CnName)
		usernameMap[userInfo.EnName] = username
	}
	return usernameMap, nil
}

const getUserMaxLength = 500

// getUserListStr get user list str
func getUserListStr(userList []string) []string {
	userListStr := make([]string, 0)

	userBuffer := bytes.Buffer{}
	for _, user := range userList {
		if userBuffer.Len()+len(user) > getUserMaxLength {
			userBuffer.WriteString(user)
			userStr := userBuffer.String()
			userListStr = append(userListStr, userStr)
			userBuffer.Reset()
			continue
		}

		userBuffer.WriteString(user)
		userBuffer.WriteByte(',')
	}

	if userBuffer.Len() == 0 {
		return userList
	}

	userStr := userBuffer.String()
	userListStr = append(userListStr, userStr[:len(userStr)-1])

	return userListStr
}

// GetInstWithOrgName get instance with organization name
func (d *Dao) GetInstWithOrgName(kit *rest.Kit, ccLang language.DefaultCCLanguageIf, objID string,
	insts []mapstr.MapStr, colProps []ColProp) ([]mapstr.MapStr, error) {

	orgIDList := make([]int64, 0)
	for _, inst := range insts {
		for _, property := range colProps {
			if property.PropertyType != common.FieldTypeOrganization || inst[property.ID] == nil {
				continue
			}

			orgIDs, ok := inst[property.ID].([]interface{})
			if !ok {
				return nil, fmt.Errorf("org id list type not []interface{}, real type is %T", inst[property.ID])
			}

			if len(orgIDs) == 0 {
				continue
			}

			for _, orgID := range orgIDs {
				id, err := util.GetInt64ByInterface(orgID)
				if err != nil {
					return nil, err
				}
				orgIDList = append(orgIDList, id)
			}
		}
	}

	orgIDs := util.IntArrayUnique(orgIDList)
	organizations, err := getAllOrganization(kit, orgIDs)
	if err != nil {
		blog.Errorf("get department failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}
	orgMap := make(map[int64]string)
	for _, item := range organizations.Results {
		orgMap[item.ID] = item.FullName
	}

	for idx, inst := range insts {
		for _, property := range colProps {
			if property.PropertyType != common.FieldTypeOrganization || inst[property.ID] == nil {
				continue
			}

			orgIDs, ok := inst[property.ID].([]interface{})
			if !ok {
				return nil, fmt.Errorf("org id list type not []interface{}, real type is %T", inst[property.ID])
			}

			if len(orgIDs) == 0 {
				continue
			}

			orgName := make([]string, 0)
			for _, orgID := range orgIDs {
				id, err := util.GetInt64ByInterface(orgID)
				if err != nil {
					return nil, err
				}

				name, exist := orgMap[id]
				if !exist {
					orgName = append(orgName, fmt.Sprintf("[%d]%s", id, ccLang.Language("nonexistent_org")))
					continue
				}

				orgName = append(orgName, fmt.Sprintf("[%d]%s", id, name))
			}

			insts[idx][property.ID] = strings.Join(orgName, ",")
		}
	}

	return insts, nil
}

// getAllOrganization get organization info from paas
func getAllOrganization(kit *rest.Kit, orgIDs []int64) (*metadata.DepartmentData, errors.CCErrorCoder) {

	loginVersion, _ := cc.String("webServer.login.version")
	if loginVersion == common.BKOpenSourceLoginPluginVersion || loginVersion == common.BKSkipLoginPluginVersion {
		return &metadata.DepartmentData{}, nil
	}

	orgIDList := getOrgListStr(orgIDs)
	departments := &metadata.DepartmentData{}
	var wg sync.WaitGroup
	var lock sync.RWMutex
	var firstErr error
	pipeline := make(chan bool, 10)

	for _, subStr := range orgIDList {
		pipeline <- true
		wg.Add(1)
		go func(subStr string) {
			defer func() {
				wg.Done()
				<-pipeline
			}()

			params := make(map[string]string)
			params["exact_lookups"] = subStr
			result, esbErr := esb.EsbClient().User().GetAllDepartment(kit.Ctx, kit.Header, params)
			if esbErr != nil {
				firstErr = esbErr
				blog.Errorf("get department by esb client failed, params: %+v, err: %v, rid: %s", params, esbErr,
					kit.Rid)
				return
			}
			if !result.Result {
				blog.Errorf("get department by esb client failed, params: %+v, rid: %s", params, kit.Rid)
				firstErr = fmt.Errorf("get department by esb failed, params: %v", params)
				return
			}

			lock.Lock()
			departments.Count += result.Data.Count
			departments.Results = append(departments.Results, result.Data.Results...)
			lock.Unlock()
		}(subStr)
	}
	wg.Wait()

	if firstErr != nil {
		return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	return departments, nil
}

const organizationMaxLength = 500

// getOrgListStr get org list str
func getOrgListStr(orgIDList []int64) []string {
	orgListStr := make([]string, 0)

	orgBuffer := bytes.Buffer{}
	for _, orgID := range orgIDList {
		if orgBuffer.Len()+len(strconv.FormatInt(orgID, 10)) > organizationMaxLength {
			orgBuffer.WriteString(strconv.FormatInt(orgID, 10))
			orgStr := orgBuffer.String()
			orgListStr = append(orgListStr, orgStr)
			orgBuffer.Reset()
			continue
		}

		orgBuffer.WriteString(strconv.FormatInt(orgID, 10))
		orgBuffer.WriteByte(',')
	}

	if orgBuffer.Len() == 0 {
		return []string{}
	}

	orgStr := orgBuffer.String()
	orgListStr = append(orgListStr, orgStr[:len(orgStr)-1])

	return orgListStr
}

// GetInstWithTable 第一个返回值是返回带有表格数据的实例信息，第二个返回值返回的每个实例数据所占用的excel行数
func (d *Dao) GetInstWithTable(kit *rest.Kit, objID string, insts []mapstr.MapStr, colProps []ColProp) (
	[]mapstr.MapStr, []int, error) {

	instHeights := make([]int, len(insts))
	for i := range instHeights {
		instHeights[i] = 1
	}

	// 1. 找出表格字段
	tableProperty := make([]ColProp, 0)
	for _, property := range colProps {
		if property.PropertyType == common.FieldTypeInnerTable {
			tableProperty = append(tableProperty, property)
		}
	}
	if len(tableProperty) == 0 {
		return insts, instHeights, nil
	}

	ids := make([]int64, 0)
	dataMap := make(map[int64]map[string][]mapstr.MapStr)
	for _, info := range insts {
		id, err := info.Int64(common.GetInstIDField(objID))
		if err != nil {
			return nil, nil, fmt.Errorf("data is invalid, err: %v", err)
		}
		ids = append(ids, id)
		dataMap[id] = make(map[string][]mapstr.MapStr)
	}

	// 2. 查询数据对应的表格字段的值
	queryOpt := metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{Filter: filtertools.GenAtomFilter(
			common.BKInstIDField, filter.In, ids)},
		Page: metadata.BasePage{Limit: common.BKMaxPageSize},
	}
	for _, property := range tableProperty {
		opt := &metadata.ListQuotedInstOption{ObjID: objID, PropertyID: property.ID, CommonQueryOption: queryOpt}
		instances, err := d.Engine.CoreAPI.ApiServer().ModelQuote().ListQuotedInstance(kit.Ctx, kit.Header, opt)
		if err != nil {
			return nil, nil, err
		}
		for _, inst := range instances.Info {
			instID, err := inst.Int64(common.BKInstIDField)
			if err != nil {
				return nil, nil, err
			}
			dataMap[instID][property.ID] = append(dataMap[instID][property.ID], inst)
		}
	}

	// 3. 整理返回带表格数据的结果, 以及每条数据需要占用excel多少行
	for idx := range insts {
		id, err := insts[idx].Int64(common.GetInstIDField(objID))
		if err != nil {
			return nil, nil, fmt.Errorf("data is invalid, err: %v", err)
		}

		for propertyID, data := range dataMap[id] {
			insts[idx][propertyID] = data

			if len(data) > instHeights[idx] {
				instHeights[idx] = len(data)
			}
		}
	}

	return insts, instHeights, nil
}

// TopoBriefMsg topo brief message
type TopoBriefMsg struct {
	ObjID string
	Name  string
}

// GetCustomTopoBriefMsg get custom topo brief message
func (d *Dao) GetCustomTopoBriefMsg(kit *rest.Kit) ([]TopoBriefMsg, error) {
	objIDs, err := d.getCustomTopoObjIDs(kit)
	if err != nil {
		blog.Errorf("get custom topo objID failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	if len(objIDs) == 0 {
		return nil, nil
	}

	input := &metadata.QueryCondition{
		Fields:    []string{common.BKObjIDField, common.BKObjNameField},
		Condition: mapstr.MapStr{common.BKObjIDField: mapstr.MapStr{common.BKDBIN: objIDs}},
	}

	objResult, err := d.Engine.CoreAPI.ApiServer().ReadModel(kit.Ctx, kit.Header, input)
	if err != nil {
		blog.Errorf("search mainline obj failed, objIDs: %#v, err: %v, rid: %s", objIDs, err, kit.Rid)
		return nil, err
	}

	result := make([]TopoBriefMsg, len(objResult.Info))
	for idx, val := range objResult.Info {
		result[idx] = TopoBriefMsg{ObjID: val.ObjectID, Name: val.ObjectName}
	}

	return result, nil
}

func (d *Dao) getCustomTopoObjIDs(kit *rest.Kit) ([]string, error) {
	query := &metadata.QueryCondition{
		Condition: mapstr.MapStr{common.AssociationKindIDField: common.AssociationKindMainline},
	}
	mainlineAsstRsp, err := d.Engine.CoreAPI.ApiServer().ReadModuleAssociation(kit.Ctx, kit.Header, query)
	if err != nil {
		blog.Errorf("search mainline association failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	mainlineObjChildMap := make(map[string]string, 0)
	for _, asst := range mainlineAsstRsp.Data.Info {
		if asst.ObjectID == common.BKInnerObjIDHost {
			continue
		}
		mainlineObjChildMap[asst.AsstObjID] = asst.ObjectID
	}

	// get all mainline custom object id
	objIDs := make([]string, 0)
	for objectID := common.BKInnerObjIDApp; len(objectID) != 0; objectID = mainlineObjChildMap[objectID] {
		if objectID == common.BKInnerObjIDApp || objectID == common.BKInnerObjIDSet ||
			objectID == common.BKInnerObjIDModule {
			continue
		}

		objIDs = append(objIDs, objectID)
	}

	return objIDs, nil
}
