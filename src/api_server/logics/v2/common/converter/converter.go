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

package converter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"configcenter/src/api_server/logics/v2/common/defs"
	"configcenter/src/api_server/logics/v2/common/utils"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	ccError "configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/coccyx/timeparser"
	"github.com/emicklei/go-restful"
)

// RespCommonResV2 turn the result without data into version V2
func RespCommonResV2(result bool, code int, message string, resp *restful.Response) {

	resV2 := make(mapstr.MapStr)

	if result {
		resV2["code"] = 0
		resV2["data"] = "success"
	} else {
		resV2["code"] = code
		resV2["msg"] = message
		resV2["extmsg"] = nil
	}
	s, _ := json.Marshal(resV2)
	io.WriteString(resp, string(s))
}

// RespSuccessV2 turn the result of successful data into V2 version
func RespSuccessV2(data interface{}, resp *restful.Response) {
	res_v2 := make(map[string]interface{})
	res_v2["code"] = 0
	res_v2["data"] = data
	s, err := json.Marshal(res_v2)
	if err != nil {
		blog.Errorf("ResToV2ForRoleApp error:%v, reply:%v", err, res_v2)
		RespFailV2(common.Json_Marshal_ERR, common.Json_Marshal_ERR_STR, resp)
	}

	io.WriteString(resp, string(s))
}

// RespFailV2Error convert the result of the failed data to V2
func RespFailV2Error(err ccError.CCError, resp *restful.Response) {
	res_v2 := make(map[string]interface{})

	if ccErr, ok := err.(ccError.CCErrorCoder); ok {
		res_v2["code"] = ccErr.GetCode()
	}
	res_v2["result"] = false
	res_v2["msg"] = err.Error()
	res_v2["extmsg"] = nil
	s, _ := json.Marshal(res_v2)
	io.WriteString(resp, string(s))
}

// RespSuccessV2 convert the result of the failed data to V2
func RespFailV2(code int, msg string, resp *restful.Response) {
	res_v2 := make(map[string]interface{})
	res_v2["code"] = code
	res_v2["msg"] = msg
	res_v2["extmsg"] = nil
	s, _ := json.Marshal(res_v2)
	io.WriteString(resp, string(s))
}

// DecorateUserName add suffixes to usernames to filter in roles
func DecorateUserName(originUserName string) string {
	return originUserName + ""
}

// ResToV2ForAppList  convert cc v3 json data to cc v2 for application list
func ResToV2ForAppList(resDataV3 metadata.InstResult) (interface{}, error) {

	resDataV2 := make([]map[string]interface{}, 0)
	for _, item := range resDataV3.Info {
		mapV2, err := convertOneApp(item)
		if nil != err {
			blog.Errorf("get app list error:%s, reply:%v", err.Error(), resDataV3)
			return nil, err
		}

		resDataV2 = append(resDataV2, mapV2)
	}

	return resDataV2, nil
}

//ResToV2ForAppList: convert cc v3 json data to cc v2 for application list
func ResToV2ForRoleApp(resDataV3 metadata.InstResult, uin string, roleArr []string) (interface{}, error) {

	resDataV2 := make(map[string][]interface{})

	for _, role := range roleArr {
		resDataV2[role] = make([]interface{}, 0)
	}

	resDataInfoV3 := resDataV3.Info
	for _, itemMap := range resDataInfoV3 {

		mapV2, err := convertOneApp(itemMap)
		if nil != err {
			blog.Errorf("ResToV2ForRoleApp error:%v, reply:%s", err, resDataV3)
			return nil, err
		}
		for _, roleStr := range roleArr {

			roleStrV3, ok := defs.RoleMap[roleStr]

			if !ok {
				continue
			}

			apps, ok := resDataV2[roleStr]
			if !ok {
				apps = make([]interface{}, 0)
				resDataV2[roleStr] = apps
			}
			roleUsers, ok := itemMap[roleStrV3]
			if !ok {
				continue
			}
			strUser, _ := roleUsers.(string)
			roleUsersList := strings.Split(strUser, ",")
			if util.InStrArr(roleUsersList, uin) {
				resDataV2[roleStr] = append(apps, mapV2)

			}

		}

	}

	return resDataV2, nil
}

//ResToV2ForModuleList: convert cc v3 json data to cc v2 for module
func ResToV2ForModuleList(result bool, message string, data interface{}) (interface{}, error) {

	resDataV2 := make([]string, 0)
	resDataV3, err := getResDataV3(result, message, data)
	if nil != err {
		return nil, err
	}

	resDataInfoV3 := (resDataV3.(map[string]interface{}))["info"].([]interface{})

	for _, item := range resDataInfoV3 {
		item_map := item.(map[string]interface{})
		resDataV2 = append(resDataV2, item_map["ModuleName"].(string))
	}

	return resDataV2, nil
}

//ResToV2ForModuleList: convert cc v3 json data to cc v2 for module map list
func ResToV2ForModuleMapList(data metadata.InstResult) ([]mapstr.MapStr, error) {
	resDataV2 := make([]mapstr.MapStr, 0)
	for _, itemMap := range data.Info {
		convMap, err := convertFieldsIntToStr(itemMap, []string{common.BKSetIDField, common.BKModuleIDField, common.BKAppIDField})
		if nil != err {
			return nil, err
		}
		if itemMap[common.BKModuleNameField].(string) == common.DefaultFaultModuleName {
			itemMap[common.BKDefaultField] = "1"
		}
		if itemMap[common.BKModuleNameField].(string) == common.DefaultResModuleName {
			itemMap[common.BKDefaultField] = "1"
		}
		moduleType, ok := itemMap[common.BKModuleTypeField]
		if false == ok || nil == moduleType {
			moduleType = "1"
		}
		moduleType = fmt.Sprintf("%v", moduleType)

		resDataV2 = append(resDataV2, mapstr.MapStr{
			"ModuleID":      convMap[common.BKModuleIDField],
			"ApplicationID": convMap[common.BKAppIDField],
			"ModuleName":    itemMap[common.BKModuleNameField],
			//"BakOperator": "",
			"CreateTime": convertToV2Time(itemMap[common.CreateTimeField]),
			"Default":    itemMap[common.BKDefaultField],
			//"Description": "",
			//"Operator": "",
			"ModuleType": moduleType,
			"SetID":      convMap[common.BKSetIDField],
		})
	}

	return resDataV2, nil
}

//ResToV2ForSetList: convert cc v3 json data to cc v2 for set
func ResToV2ForSetList(result bool, message string, data metadata.InstResult) (interface{}, error) {
	resDataV2 := make([]map[string]interface{}, 0)

	_, err := getResDataV3(result, message, data)
	if nil != err {
		return nil, err
	}
	for _, item := range data.Info {
		convMap, err := convertFieldsIntToStr(item, []string{common.BKSetIDField})
		if nil != err {
			return nil, err
		}
		setName, ok := item[common.BKSetNameField]
		if false == ok {
			return nil, errors.New("get set info error")
		}
		if setName == common.DefaultResSetName {
			setName = "空闲机池"
		}
		resDataV2 = append(resDataV2, map[string]interface{}{
			"SetID":   convMap[common.BKSetIDField],
			"SetName": setName, //itemMap[common.BKSetNameField],
		})
	}

	return resDataV2, nil
}

//ResToV2ForPlatList: convert cc v3 json data to cc v2 for plat
func ResToV2ForPlatList(data metadata.InstResult) (interface{}, error) {
	resDataV2 := make([]map[string]interface{}, 0)
	for _, itemMap := range data.Info {
		convMap, err := convertFieldsIntToStr(itemMap, []string{common.BKCloudIDField})
		if nil != err {
			return nil, err
		}

		resDataV2 = append(resDataV2, map[string]interface{}{
			"platId":      convMap[common.BKCloudIDField],
			"platName":    itemMap[common.BKCloudNameField],
			"platCompany": itemMap[common.BKOwnerIDField],
		})
	}

	return resDataV2, nil
}

//ResToV2ForHostList: convert cc v3 json data to cc v2 for host
func ResToV2ForHostList(result bool, message string, data interface{}) (interface{}, error) {

	resDataInfoV3, err := getResDataV3(result, message, data)
	if nil != err {
		blog.Errorf("ResToV2ForHostList reply:%v, error:%s", data, err.Error())
		return nil, err
	}

	return convertToV2HostListMain(resDataInfoV3)

}

func convertToV2HostListMain(resDataInfoV3 interface{}) (interface{}, error) {
	resDataV2 := make([]interface{}, 0)
	if nil == resDataInfoV3 {
		return resDataV2, nil
	}

	var dataArr []mapstr.MapStr
	switch realData := resDataInfoV3.(type) {
	case []mapstr.MapStr:
		dataArr = realData
	case []interface{}:
		for _, item := range realData {
			itemMap, err := mapstr.NewFromInterface(item)
			if nil != err {
				blog.Errorf("ResToV2ForHostList not map[string]interface resDataInfoV3 %v, error:%s", resDataInfoV3, err.Error())
				return nil, errors.New("http reply data error")
			}
			dataArr = append(dataArr, itemMap)
		}
	default:
		blog.Errorf("ResToV2ForHostList not []map[string]interface resDataInfoV3 %v", resDataInfoV3)
		return nil, errors.New("http reply data error")
	}

	for _, itemMap := range dataArr {

		hostID, err := util.GetInt64ByInterface(itemMap[common.BKHostIDField])
		if nil != err {
			blog.Warnf("convertToV2HostListMain hostID not found, appID:%s, hostInfo:%+v", itemMap[common.BKAppIDField], itemMap)
			continue
		}

		innerIP, ok := itemMap[common.BKHostInnerIPField].(string)
		if !ok {
			blog.Warnf("convertToV2HostListMain innerIP not found, appID:%s, hostInfo:%+v", itemMap[common.BKAppIDField], itemMap)
			continue
		}
		hostHard := convHostHardInfo(hostID, innerIP, itemMap)
		itemMap.Set("ExtInfo", hostHard)
		resDataV2 = append(resDataV2, GeneralV2Data(itemMap))
	}
	return resDataV2, nil
}

//ResToV2ForHostGroup: convert cc v3 json data to cc v2 for host group
func ResToV2ForHostGroup(result bool, message string, data interface{}) (interface{}, error) {
	resDataV2 := make(map[string]interface{}, 0)
	resDataInfoV3, err := getResDataV3(result, message, data)
	if nil != err {
		blog.Errorf("ResToV2ForHostList reply:%v, error:%s", data, err.Error())
		return nil, err
	}

	for k, v := range resDataInfoV3.(map[string]interface{}) {
		resDataV2[k] = GeneralV2Data(v)
	}

	return resDataV2, nil
}

//ResToV2ForCpyHost: convert cc v3 json data to cc v2 for getCompanyIDByIps
func ResToV2ForCpyHost(result bool, message string, data interface{}) (interface{}, error) {
	resDataV2 := make(map[string]interface{})

	resDataV3, err := getResDataV3(result, message, data)
	if nil != err {
		return nil, err
	}

	resDataArrV3 := resDataV3.([]interface{})

	for _, item := range resDataArrV3 {
		itemMap := item.(map[string]interface{})

		appID, err := util.GetIntByInterface(itemMap[common.BKAppIDField])

		if nil != err {
			return resDataV2, nil
		}
		bkCloudID, err := util.GetIntByInterface(itemMap[common.BKCloudIDField])
		if nil != err {
			return resDataV2, nil
		}
		ownerID, err := util.GetIntByInterface(itemMap[common.BKOwnerIDField])
		if nil != err {
			return resDataV2, nil
		}
		buildStr := fmt.Sprintf("%d%d%d", bkCloudID, ownerID, appID)
		itemMap = convertFieldsNilToString(itemMap, []string{common.BKCloudIDField, common.BKOwnerIDField, common.BKAppIDField})

		resDataV2[itemMap[common.BKHostInnerIPField].(string)] = map[string]interface{}{
			buildStr: map[string]interface{}{
				"PlatID":        itemMap[common.BKCloudIDField],
				"CompanyID":     itemMap[common.BKOwnerIDField],
				"ApplicationID": itemMap[common.BKAppIDField],
			},
		}
	}

	return resDataV2, nil
}

func ResToV2ForPropertyList(result bool, message string, data interface{}, idName, idDisplayName string) (interface{}, error) {
	resDataV2 := map[string]interface{}{}
	standardMap := make(map[string]interface{})
	customerMap := make(map[string]interface{})

	resDataV3, err := getResDataV3(result, message, data)
	if nil != err {
		return nil, err
	}

	for _, item := range resDataV3.([]interface{}) {
		itemMap := item.(map[string]interface{})
		fileName := ConverterV3Fields(itemMap[common.BKPropertyIDField].(string), "")
		if itemMap[common.BKIsPre] != nil && itemMap[common.BKIsPre].(bool) {
			standardMap[fileName] = itemMap[common.BKPropertyNameField]
		} else {
			customerMap[fileName] = itemMap[common.BKPropertyNameField]
		}
	}
	standardMap[idName] = idDisplayName

	resDataV2["standard"] = standardMap
	resDataV2["customer"] = customerMap

	return resDataV2, nil
}

// ResToV2ForAppTree: convert cc v3 json data to cc v2 for topo tree
func ResToV2ForAppTree(result bool, message string, data interface{}) (interface{}, error) {
	resDataV3, err := getResDataV3(result, message, data)
	if nil != err {
		return nil, err
	}

	resDataV2 := getOneLevelData(resDataV3.([]interface{}), nil)
	if len(resDataV2) > 0 {
		return resDataV2[0], nil
	} else {
		return nil, nil
	}
}

//ResToV2ForCustomerGroup
func ResToV2ForCustomerGroup(result bool, message string, data interface{}, appID string) ([]common.KvMap, error) {
	resDataV3, err := getResDataV3(result, message, data)
	if nil != err {
		return nil, err
	}

	resDataArrV3, _ := (resDataV3.(map[string]interface{}))["info"].([]interface{})
	var ret []common.KvMap
	for _, item := range resDataArrV3 {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			return nil, errors.New("data format errors")
		}
		itemMap = convertFieldsNilToString(itemMap, []string{"bk_info"})
		ret = append(ret, common.KvMap{
			"ID":            itemMap["id"],
			"ApplicationID": appID,
			"GroupName":     itemMap["name"], //itemMap["Name"], //TODO 待确认
			"GroupContent":  itemMap["info"],
			"Type":          "host",
			"CreateTime":    convertToV2Time(itemMap[common.CreateTimeField]),
			"LastTime":      convertToV2Time(itemMap[common.LastTimeField]),
		})

	}
	return ret, nil
}

//ResToV2ForCustomerGroupResult return list, total, error
func ResToV2ForCustomerGroupResult(result bool, message string, dataInfo interface{}) ([]common.KvMap, int, error) {
	resDataV3, err := getResDataV3(result, message, dataInfo)
	if nil != err {
		return nil, 0, err
	}
	if "" == resDataV3 {
		return nil, 0, nil
	}

	data, ok := resDataV3.(map[string]interface{})
	if !ok {
		blog.Errorf("ResToV2ForCustomerGroupResult data item not found, %v", data)
		return nil, 0, errors.New("data format errors")
	}
	iCount, ok := data["count"]
	if !ok {
		blog.Error("ResToV2ForCustomerGroupResult count item not found")
		return nil, 0, errors.New("data format errors")
	}
	total, _ := util.GetIntByInterface(iCount)

	resDataArrV3, _ := data["info"].([]interface{})
	var ret []common.KvMap
	for _, item := range resDataArrV3 {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			blog.Errorf("ResToV2ForCustomerGroupResult data hostinfo item errors, %v", item)
			return nil, 0, errors.New("data format errors")
		}
		host, ok := itemMap[common.BKInnerObjIDHost].(map[string]interface{})
		if !ok {
			blog.Errorf("ResToV2ForCustomerGroupResult data hostinfo  host item errors, %v", itemMap)
			return nil, 0, errors.New("data format errors")
		}
		modules, ok := itemMap[common.BKInnerObjIDModule].([]interface{})
		if !ok {
			blog.Errorf("ResToV2ForCustomerGroupResult data hostinfo  module item errors, %v", itemMap)
			return nil, 0, errors.New("data format errors")
		}
		sets, ok := itemMap[common.BKInnerObjIDSet].([]interface{})
		if !ok {
			blog.Errorf("ResToV2ForCustomerGroupResult data hostinfo set item errors, %v", itemMap)
			return nil, 0, errors.New("data format errors")
		}
		innerIP, _ := host[common.BKHostInnerIPField]
		if !ok {
			blog.Errorf("ResToV2ForCustomerGroupResult data hostinfo host innerip item errors, %v", itemMap)
			return nil, 0, errors.New("data format errors")
		}
		hostName, _ := host[common.BKHostNameField]
		moduleName := "" // module[common.BKModuleNameField]
		setName := ""    //set[common.BKSetNameField]
		if 0 < len(modules) {
			for _, module := range modules {
				moduleMap, ok := module.(map[string]interface{})
				if false == ok {
					blog.Errorf("ResToV2ForCustomerGroupResult data hostinfo  module item errors, %v", itemMap)
					return nil, 0, errors.New("data format errors")
				}
				moduleName, _ = moduleMap[common.BKModuleNameField].(string)
				break

			}
		}
		if 0 < len(sets) {
			for _, set := range sets {
				setMap, ok := set.(map[string]interface{})
				if false == ok {
					blog.Errorf("ResToV2ForCustomerGroupResult data hostinfo set item errors, %v", itemMap)
					return nil, 0, errors.New("data format errors")
				}
				setName, _ = setMap[common.BKSetNameField].(string)
				break

			}
		}
		subArea, _ := host[common.BKSubAreaField].([]interface{}) //host["SubArea"].([]interface{})
		var source int64 = -1
		if nil != subArea && len(subArea) > 0 {
			sourceItem := subArea[0].(map[string]interface{})
			source, _ = util.GetInt64ByInterface(sourceItem[common.BKInstIDField])
		}
		if nil == hostName {
			hostName = ""
		}

		ret = append(ret, common.KvMap{
			"SetName":    setName,
			"ModuleName": moduleName,
			"Source":     fmt.Sprintf("%d", source),
			"HostName":   hostName,
			"InnerIP":    innerIP,
		})

	}
	return ret, total, nil
}

func ResToV2ForHostDataList(result bool, message string, data interface{}) (common.KvMap, error) {
	resDataV3, err := getResDataV3(result, message, data)
	if nil != err {
		return nil, err
	}
	convFields := []string{common.BKAppNameField, common.BKModuleNameField, common.BKBakOperatorField, common.BKSetNameField, common.BKOperatorField, common.BKSetIDField, common.BKAppIDField, common.BKModuleIDField}
	var ret common.KvMap

	if "" != resDataV3 {
		resDataArrV3, ok := resDataV3.([]interface{})
		if !ok {
			blog.Errorf("ResToV2ForHostDataList not array data :%#v", data)
			return nil, fmt.Errorf("data is not array %#v", resDataV3)
		}
		var operators []string
		var bakOperators []string
		var moduleIDs []string
		var moduleNames []string
		var setIDs []string
		var setNames []string

		for _, item := range resDataArrV3 {
			itemMap, ok := item.(map[string]interface{})
			if !ok {
				blog.Warnf("ResToV2ForHostDataList item %+v not map[string]interface{}, raw data", item, data)
				continue
			}
			itemMap = convertFieldsNilToString(itemMap, convFields)
			moduleName, ok := itemMap[common.BKModuleNameField].(string)
			if ok && "" != moduleName {
				moduleNames = append(moduleNames, moduleName)
				moduleIDs = append(moduleIDs, fmt.Sprintf("%v", itemMap[common.BKModuleIDField]))
			}
			setName, ok := itemMap[common.BKSetNameField].(string)
			if ok && "" != setName {
				setNames = append(setNames, setName)
				setIDs = append(setIDs, fmt.Sprintf("%v", itemMap[common.BKSetIDField]))
			}
			operator, ok := itemMap[common.BKOperatorField].(string)
			if ok && "" != operator {
				operators = append(operators, operator)
			}
			bakOperator, ok := itemMap[common.BKBakOperatorField].(string)
			if ok && "" != bakOperator {
				bakOperators = append(bakOperators, bakOperator)
			}
			ret = common.KvMap{
				"ApplicationName": itemMap[common.BKAppNameField],
				"ApplicationID":   itemMap[common.BKAppIDField],
			}
		}
		ret["ModuleName"] = strings.Join(moduleNames, ",")
		ret["ModuleID"] = strings.Join(moduleIDs, ",")
		ret["SetName"] = strings.Join(moduleNames, ",")
		ret["SetID"] = strings.Join(setIDs, ",")
		ret["Operator"] = strings.Join(operators, ",")
		ret["BakOperator"] = strings.Join(bakOperators, ",")
	}
	if 1 <= len(ret) {

		return ret, nil
	}
	return nil, nil

}

// ResToV2ForEnterIP get enterip result  for v2
func ResToV2ForEnterIP(result bool, message string, data interface{}) error {
	_, err := getResDataV3(result, message, data)
	return err
}

// ResV2ToForProcList get process info for v2
func ResV2ToForProcList(resDataV3 interface{}, defLang language.DefaultCCLanguageIf) interface{} {
	resDataArrV3 := resDataV3.([]interface{})
	ret := make([]interface{}, 0)
	for _, item := range resDataArrV3 {
		itemMap := item.(map[string]interface{})
		itemMap = convertFieldsNilToString(itemMap, []string{common.BKAppIDField, common.BKAppNameField, common.BKHostInnerIPField, common.BKHostOuterIPField})

		ret = append(ret, common.KvMap{
			"ApplicationID":   itemMap[common.BKAppIDField],
			"ApplicationName": itemMap[common.BKAppNameField],
			"InnerIP":         itemMap[common.BKHostInnerIPField],
			"OuterIP":         itemMap[common.BKHostOuterIPField],
			"process":         getOneProcData(itemMap["process"], defLang),
		})

	}

	return ret
}

// GeneralV2Data  general convertor v2 funcation
func GeneralV2Data(data interface{}) interface{} {

	switch realData := data.(type) {
	case []mapstr.MapStr:
		mapItem := make([]map[string]interface{}, 0)
		for _, item := range realData {
			mapItem = append(mapItem, convMapInterface(item))
		}
		return mapItem
	case []interface{}:
		mapItem := make([]interface{}, 0)
		for _, item := range realData {
			if nil == item {
				continue
			}
			mapItem = append(mapItem, GeneralV2Data(item))
		}
		return mapItem
	case mapstr.MapStr:
		return convMapInterface(realData)
	case map[string]interface{}:
		return convMapInterface(realData)
	}

	if nil == data {
		return ""
	}

	return convToV2ValStr(data) //fmt.Sprintf("%v", data)

}

func GetHostHardInfo(appID int64, hostInfoArr []mapstr.MapStr) []mapstr.MapStr {
	hostHardInfoArr := make([]mapstr.MapStr, 0)
	alreadyIPMap := make(map[int64]bool, 0)
	for _, host := range hostInfoArr {
		hostID, err := host.Int64(common.BKHostIDField)
		if nil != err {
			blog.Warnf("GetHostHardInfo hostID not found, appID:%s, hostInfo:%+v", appID, host)
			continue
		}
		innerIP, err := host.String(common.BKHostInnerIPField)
		if nil != err {
			blog.Warnf("GetHostHardInfo innerIP not found, appID:%s, hostInfo:%+v", appID, host)
			continue
		}
		_, isHandle := alreadyIPMap[hostID]
		if isHandle {
			continue
		}
		alreadyIPMap[hostID] = true
		hostHardInfoArr = append(hostHardInfoArr, convHostHardInfo(hostID, innerIP, host))
	}
	return hostHardInfoArr
}
func convHostHardInfo(hostID int64, innerIP string, host mapstr.MapStr) (hostHardInfo mapstr.MapStr) {
	hostHardInfo = mapstr.New()
	osVersionEnumID, err := host.String(common.BKOSTypeField)
	osVersion := ""
	if nil != err {
		hostHardInfo.Set("PlatformOS", "")
	} else {
		osVersion = getOSTypeByEnumID(osVersionEnumID)
		hostHardInfo.Set("PlatformOS", osVersion)
	}
	network := mapstr.New()
	innerMac, ok := host.Get("bk_mac")
	if ok {
		network.Set(innerIP, innerMac)
	}
	outerIP, err := host.String(common.BKHostOuterIPField)
	if nil == err {
		if 0 < len(outerIP) {
			outerMac, err := host.String("bk_outer_mac")
			if nil == err {
				network.Set(outerIP, outerMac)
			}
		}
	}
	system := mapstr.New()
	system.Set("OS", osVersion)
	dockerClientVersion, clientOk := host.Get(common.HostFieldDockerClientVersion)
	dockerServerVersion, serverOk := host.Get(common.HostFieldDockerServerVersion)
	if clientOk || serverOk {
		system.Set("clientDockerVersion", dockerClientVersion)
		system.Set("serverDockerVersion", dockerServerVersion)
	}
	system.Set("kernelVersion", "")
	mem, err := host.Int64("bk_mem")
	if nil != err {
		mem = 0
	} else {
		mem = mem * 1024 * 1024
	}
	disk, err := host.Int64("bk_disk")
	if nil != err {
		disk = 0
	} else {
		disk = disk * 1024 * 1024 * 1024
	}
	cpu, err := host.Int64("bk_cpu")
	if nil != err {
		cpu = 0
	}
	hostHardInfo.Set("System", system)
	hostHardInfo.Set("network", network)
	hostHardInfo.Set("InnerIP", innerIP)
	hostHardInfo.Set("OuterIP", outerIP)
	hostHardInfo.Set("HostID", fmt.Sprintf("%d", hostID))
	hostHardInfo.Set("Memory", mapstr.MapStr{"Total": mem})
	hostHardInfo.Set("Disk", mapstr.MapStr{"Total": disk})
	hostHardInfo.Set("Cpu", mapstr.MapStr{"CpuNum": cpu})
	return hostHardInfo
}

func convMapInterface(data map[string]interface{}) map[string]interface{} {
	mapItem := make(map[string]interface{})
	for v3key, val := range data {
		key := ConverterV3Fields(v3key, "")
		if key == "CreateTime" || key == "LastTime" || key == common.CreateTimeField || key == common.LastTimeField {
			ts, ok := val.(time.Time)
			if ok {
				mapItem[key] = ts.Format("2006-01-02 15:04:05")

			} else {
				mapItem[key] = ""
			}
		} else if common.BKProtocol == key || "Protocol" == key {
			//v2 api erturn use protocol name
			protocol, ok := val.(string)
			if false == ok {
				protocol = ""
			} else {
				switch protocol {
				case "1":
					protocol = "TCP"
				case "2":
					protocol = "UDP"
				default:
					protocol = ""
				}
			}
			mapItem[key] = protocol
		} else if key == "osType" {

			switch realVal := val.(type) {
			case string:
				mapItem[key] = getOSTypeByEnumID(realVal)

			case nil:
				mapItem[key] = ""
			default:
				mapItem[key] = realVal
			}

			mapItem["OSType"] = mapItem[key]
		} else if v3key == common.BKCloudIDField {
			switch rawVal := val.(type) {
			case []mapstr.MapStr:
				if len(rawVal) == 0 {
					mapItem[key] = ""
				}
				strVal, err := rawVal[0].String(common.BKInstIDField)
				if err != nil {
					mapItem[key] = ""
				}
				mapItem[key] = strVal
			case []interface{}:
				if len(rawVal) == 0 {
					mapItem[key] = ""
				}
				cloudInfo, err := mapstr.NewFromInterface(rawVal[0])
				if err != nil {
					mapItem[key] = ""
				}
				strVal, err := cloudInfo.String(common.BKInstIDField)
				if err != nil {
					mapItem[key] = ""
				}
				mapItem[key] = strVal
			default:
				intVal, err := util.GetInt64ByInterface(rawVal)
				if err != nil {
					mapItem[key] = ""
				}
				mapItem[key] = strconv.FormatInt(intVal, 10)
			}

		} else {
			mapItem[key] = GeneralV2Data(val)
		}

	}
	return mapItem
}

// getOneLevelData  get one level data
func getOneLevelData(data []interface{}, appID interface{}) []map[string]interface{} {
	dataArrTemp := make([]map[string]interface{}, 0)
	for _, item := range data {
		itemMap, ok := item.(map[string]interface{})
		if false == ok {
			blog.Errorf("Assert error item is not map[string]interface{},item %v", item)
			continue
		}
		dataTemp := make(map[string]interface{})
		dataTemp["ObjID"] = itemMap[common.BKObjIDField]
		InstId, _ := util.GetIntByInterface(itemMap[common.BKInstIDField])
		appIdInt, _ := util.GetIntByInterface(appID)
		strInstId := strconv.Itoa(InstId)
		appIdStr := strconv.Itoa(appIdInt)

		switch itemMap[common.BKObjIDField] {
		case common.BKInnerObjIDApp:
			//dataTemp = itemMap
			dataTemp["Level"] = 3
			dataTemp["ApplicationID"] = strInstId
			dataTemp["ApplicationName"] = itemMap[common.BKInstNameField]
			appID = itemMap[common.BKInstIDField]
		case common.BKInnerObjIDSet:
			//dataTemp = itemMap
			dataTemp["SetID"] = strInstId
			dataTemp["SetName"] = itemMap[common.BKInstNameField]
		case common.BKInnerObjIDModule:
			//dataTemp = itemMap
			dataTemp["ApplicationID"] = appIdStr
			dataTemp["ModuleID"] = itemMap[common.BKInstIDField]
			dataTemp["ModuleName"] = itemMap[common.BKInstNameField]

		default:
			if nil != itemMap["child"] {
				children := getOneLevelData(itemMap["child"].([]interface{}), appID)
				for _, child := range children {
					dataArrTemp = append(dataArrTemp, child)
				}
			}
			continue
		}

		if nil != itemMap["child"] {
			children := getOneLevelData(itemMap["child"].([]interface{}), appID)
			if len(children) > 0 {
				dataTemp["Children"] = children
			}
		}
		dataArrTemp = append(dataArrTemp, dataTemp)
	}
	return dataArrTemp
}

// getOneProcData get one process data
func getOneProcData(data interface{}, defLang language.DefaultCCLanguageIf) interface{} {
	var ret interface{}

	itemMap := data.(map[string]interface{})

	createTime, _ := itemMap[common.CreateTimeField]

	switch createTime.(type) {
	case time.Time:
		createTime = createTime.(time.Time).Format("2006-01-02 15:04:05")
	case string:
		ts, _ := timeparser.TimeParser(createTime.(string))
		createTime = ts.Format("2006-01-02 15:04:05")
	default:
		createTime = ""
	}

	updateTime, _ := itemMap[common.LastTimeField]
	switch createTime.(type) {
	case time.Time:
		updateTime = updateTime.(time.Time).Format("2006-01-02 15:04:05")
	case string:
		ts, _ := timeparser.TimeParser(updateTime.(string))
		updateTime = ts.Format("2006-01-02 15:04:05")
	default:
		updateTime = ""
	}
	protocol, ok := itemMap[common.BKProtocol].(string)
	if false == ok {
		protocol = ""
	} else {
		switch protocol {
		case "1":
			protocol = "TCP"
		case "2":
			protocol = "UDP"
		default:
			protocol = ""
		}
	}

	intAtuotimeGap, err := util.GetIntByInterface(itemMap["auto_time_gap"])
	atuotimeGap := ""
	if err == nil {
		atuotimeGap = fmt.Sprintf("%d", intAtuotimeGap)
	}

	convFields := []string{common.BKWorkPath, common.BKFuncIDField, common.BKFuncName,
		common.BKBindIP, common.BKUser, "start_cmd", "stop_cmd", common.BKProcessNameField, common.BKPort,
		common.BKProtocol, "pid_file", "restart_cmd", "face_stop_cmd", "auto_start", "timeout", "priority", "proc_num"}
	itemMap = convertFieldsNilToString(itemMap, convFields)

	ret = map[string]interface{}{
		"WorkPath":    itemMap[common.BKWorkPath],
		"AutoTimeGap": atuotimeGap,
		"LastTime":    updateTime,
		"StartCmd":    itemMap["start_cmd"],
		"FuncID":      itemMap[common.BKFuncIDField],
		"BindIP":      itemMap[common.BKBindIP],
		"FuncName":    itemMap[common.BKFuncName],
		"Flag":        "",
		"User":        itemMap[common.BKUser],
		"StopCmd":     itemMap["stop_cmd"],
		"ProNum":      itemMap["proc_num"],
		"ReloadCmd":   itemMap["reload_cmd"],
		"ProcessName": itemMap[common.BKProcessNameField],
		"OpTimeout":   itemMap["timeout"],       //"0",
		"KillCmd":     itemMap["face_stop_cmd"], //"",
		"Protocol":    protocol,
		"Seq":         itemMap["priority"], //0",
		"ProcGrp":     "",
		"Port":        itemMap[common.BKPort],
		"ReStartCmd":  itemMap["restart_cmd"], //"",
		"AutoStart":   itemMap["auto_start"],
		"CreateTime":  createTime,
		"PidFile":     itemMap["pid_file"],
	}

	return ret
}

//convertFieldsNilToString  convertor nil to empty string in map field
func convertFieldsNilToString(itemMap map[string]interface{}, fields []string) map[string]interface{} {

	for _, field := range fields {

		val, ok := itemMap[field]
		if !ok || nil == val {
			itemMap[field] = ""
		} else {
			itemMap[field] = convToV2ValStr(val) //fmt.Sprintf("%v", val)
		}
	}

	return itemMap
}

// getResDataV3 get res data v3
func getResDataV3(result bool, message string, data interface{}) (interface{}, error) {
	if result {
		return data, nil
	} else {
		return nil, errors.New(message)
	}
}

// convertOneApp convert one len app
func convertOneApp(itemMap map[string]interface{}) (map[string]interface{}, error) {

	convMap, err := convertFieldsIntToStr(itemMap, []string{common.BKAppIDField, common.BKDefaultField})
	if nil != err {
		return nil, err
	}
	maintainer := ""
	productPm := ""
	operator := ""
	developer := ""
	tester := ""
	if nil != itemMap[common.BKMaintainersField] {
		maintainer, _ = itemMap[common.BKMaintainersField].(string)
	}
	if nil != itemMap[common.BKProductPMField] {
		productPm, _ = itemMap[common.BKProductPMField].(string)
	}
	if nil != itemMap[common.BKOperatorField] {
		operator, _ = itemMap[common.BKOperatorField].(string)
	}
	if nil != itemMap[common.BKDeveloperField] {
		developer, _ = itemMap[common.BKDeveloperField].(string)
	}
	if nil != itemMap[common.BKTesterField] {
		tester, _ = itemMap[common.BKTesterField].(string)
	}
	maintainer = strings.Replace(maintainer, ",", ";", -1)
	productPm = strings.Replace(productPm, ",", ";", -1)
	operator = strings.Replace(operator, ",", ";", -1)
	developer = strings.Replace(developer, ",", ";", -1)
	tester = strings.Replace(tester, ",", ";", -1)
	lifecycle := ""
	if nil != itemMap["life_cycle"] {
		lifecycle, _ = itemMap["life_cycle"].(string)
	}
	language := "zh-cn"
	if nil != itemMap["language"] {
		language, _ = itemMap["language"].(string)
		language = utils.ConvLanguageToV3(language)

	}

	timeZone := "Asia/Shanghai"
	if nil != itemMap[common.BKTimeZoneField] {
		timeZone, _ = itemMap[common.BKTimeZoneField].(string)
	}
	itemMapV2 := map[string]interface{}{
		"ApplicationName": itemMap[common.BKAppNameField],
		//"Description": "",
		//"BusinessDeptName": "",
		//"Creator": "",
		"Default":       convMap[common.BKDefaultField],
		"ApplicationID": convMap[common.BKAppIDField],
		"Level":         "3",
		//"Display":"",
		//"Source": "",
		//"GroupName": "",
		"Operator":    operator,
		"Developer":   developer,
		"Maintainers": maintainer,
		"CompanyID":   "0",
		"Owner":       "",
		"ProductPm":   productPm,
		"LifeCycle":   lifecycle,
		"Language":    language,
		"TimeZone":    timeZone,
		"Tester":      tester,
		"LastTime":    convertToV2Time(itemMap[common.LastTimeField]),
		"DeptName":    "",
		"CreateTime":  convertToV2Time(itemMap[common.CreateTimeField]),
	}
	return itemMapV2, nil
}

//convertToV2Time time string convertor 2018-01-23 01:02:03 format
func convertToV2Time(val interface{}) string {
	strTm, _ := val.(string)
	if "" == strTm {
		return ""
	}
	createTime, err := timeparser.TimeParser(strTm)
	if nil != err {
		return ""
	}
	m := createTime.Month()
	d := createTime.Day()
	h := createTime.Hour()
	minute := createTime.Minute()
	s := createTime.Second()

	strM := fmt.Sprintf("%d", m)
	strD := fmt.Sprintf("%d", d)
	strH := fmt.Sprintf("%d", h)
	strMinute := fmt.Sprintf("%d", minute)
	strS := fmt.Sprintf("%d", s)
	if 10 > m {
		strM = "0" + strM
	}
	if 10 > d {
		strD = "0" + strD
	}
	if 10 > h {
		strH = "0" + strH
	}
	if 10 > minute {
		strMinute = "0" + strMinute
	}
	if 10 > s {
		strS = "0" + strS
	}

	return fmt.Sprintf("%d-%s-%s %s:%s:%s", createTime.Year(), strM, strD, strH, strMinute, strS)
}

//  convertToString interface{} to string
func convertToString(itemMap map[string]interface{}) map[string]interface{} {
	tempMap := make(map[string]interface{})
	for key, val := range itemMap {
		filedInt, err := util.GetInt64ByInterface(val)
		if nil != err {
			blog.Errorf("convert field %s to number fail!value:%v", key, val)
		}
		tempMap[key] = strconv.FormatInt(filedInt, 10)
	}

	return tempMap
}

// convertFieldsIntToStr convert fields int to str
func convertFieldsIntToStr(itemMap map[string]interface{}, fields []string) (map[string]interface{}, error) {

	tempMap := make(map[string]interface{})
	for _, field := range fields {
		item, ok := itemMap[field]
		if !ok {
			continue
		}
		if nil == item {
			tempMap[field] = ""
			continue
		}

		switch item.(type) {
		case string:
		case nil:
		default:
			filedInt, err := util.GetInt64ByInterface(item)
			if nil != err {
				blog.Warnf("convert field %s to number fail!", field)
				return nil, err
			}
			tempMap[field] = strconv.FormatInt(filedInt, 10)
		}

	}

	return tempMap, nil
}

// ConverterV3Fields  converter v3 fields
func ConverterV3Fields(fields, objType string) string {

	fieldsMap := getFieldsMap(objType)
	oldFields, ok := fieldsMap[fields]
	if true == ok {
		return oldFields
	}
	return fields
}

// ConverterV2FieldsToV3 converter v2 field to v3
func ConverterV2FieldsToV3(fields, objType string) string {

	reMap := make(map[string]string)
	fieldsMap := getFieldsMap(objType)
	for k, v := range fieldsMap {
		reMap[v] = k
	}

	fieldsV3, ok := reMap[fields]
	if ok {
		return fieldsV3
	}
	return fields
}

//getV2KeyVal  convert v2 to v3.(key, val)
func getV2KeyVal(key string, val interface{}) (string, string) {
	fieldsMap := getFieldsMap("")

	v2Key, ok := fieldsMap[key]
	var v2Val string
	if !ok {
		v2Key = key
	}

	if nil != val {
		if key == common.CreateTimeField || key == common.LastTimeField {
			ts, ok := val.(time.Time)
			if ok {
				v2Val = ts.Format("2006-01-02 15:04:05")

			} else {
				v2Val = ""
			}
		} else {
			return v2Key, convToV2ValStr(val)
		}
	}

	return v2Key, v2Val
}

func convToV2ValStr(val interface{}) string {
	switch realVal := val.(type) {
	case int:
		return strconv.FormatInt(int64(realVal), 10)
	case int8:
		return strconv.FormatInt(int64(realVal), 10)
	case int16:
		return strconv.FormatInt(int64(realVal), 10)
	case int32:
		return strconv.FormatInt(int64(realVal), 10)
	case int64:
		return strconv.FormatInt(realVal, 10)
	case uint:
		return strconv.FormatInt(int64(realVal), 10)
	case uint8:
		return strconv.FormatInt(int64(realVal), 10)
	case uint16:
		return strconv.FormatInt(int64(realVal), 10)
	case uint32:
		return strconv.FormatInt(int64(realVal), 10)
	case uint64:
		return strconv.FormatInt(int64(realVal), 10)
	case float32:
		return strconv.FormatInt(int64(realVal), 10)
	case float64:
		return strconv.FormatInt(int64(realVal), 10)
	case json.Number:
		jsVal, err := realVal.Int64()
		if err != nil {
			return realVal.String()
		}
		return strconv.FormatInt(jsVal, 10)
	case string:
		return realVal
	}
	return fmt.Sprintf("%v", val)
}

func getFieldsMap(objType string) map[string]string {
	fieldsMap := map[string]string{

		common.BKAppIDField:   "ApplicationID",
		common.BKAppNameField: "ApplicationName",
		"life_cycle":          "LifeCycle",
		"language":            "Language",
		"time_zone":           "TimeZone",
		"bk_biz_developer":    "Developer",
		"bk_biz_tester":       "Tester",
		"bk_biz_maintainer":   "Maintainers",
		"bk_biz_productor":    "ProductPm",
		common.BKOwnerIDField: "Owner",
		"creator":             "Creator",

		common.BKSetIDField:   "SetID",
		common.BKSetNameField: "SetName",
		"bk_set_env":          "SetEnv",
		"bk_service_status":   "ServiceStatus",
		"description":         "Description",
		"bk_capacity":         "Capacity",

		common.BKModuleIDField:   "ModuleID",
		common.BKModuleNameField: "ModuleName",
		"bk_module_type":         "ModuleType",

		common.BKHostIDField:      "HostID",
		common.BKHostNameField:    "HostName",
		"bk_assetId":              "AssetID",
		"bk_sn":                   "SN",
		common.BKCloudIDField:     "Source",
		"bk_os_type":              "osType",
		"bk_os_name":              "OSName",
		"bk_cpu":                  "Cpu",
		"bk_mem":                  "Mem",
		common.BKHostInnerIPField: "InnerIP",
		common.BKHostOuterIPField: "OuterIP",

		"operator":             "BakOperator",
		"bk_bak_operator":      "Operator",
		common.BKDefaultField:  "Default",
		common.CreateTimeField: "LastTime",
		common.LastTimeField:   "CreateTime",

		common.BKProcIDField:   "ProcessID",
		common.BKProcNameField: "ProcessName",
		common.BKWorkPath:      "WorkPath",
		"start_cmd":            "StartCmd",
		common.BKFuncIDField:   "FuncID",
		common.BKBindIP:        "BindIP",
		common.BKFuncName:      "FuncName",
		common.BKUser:          "User",
		"stop_cmd":             "StopCmd",
		"proc_num":             "ProNum",
		"reload_cmd":           "ReloadCmd",
		//common.BKProcessNameField:    "ProcessName",
		"bk_timeout":       "OpTimeout",
		"kill_cmd":         "KillCmd",
		common.BKProcField: "Process",
		common.BKProtocol:  "Protocol",
		"priority":         "Seq",
		"seq":              "Seq",
		common.BKPort:      "Port",
		"restart_cmd":      "ReStartCmd",
		"auto_start":       "AutoStart",
		"pid_file":         "PidFile",
		"face_stop_cmd":    "KillCmd",
		"timeout":          "OpTimeout",
		"auto_time_gap":    "AutoTimeGap",
	}
	return fieldsMap
}

func getOSTypeByEnumID(enumID string) (OSType string) {
	switch enumID {
	case common.HostOSTypeEnumLinux:
		OSType = "linux"
	case common.HostOSTypeEnumWindows:
		OSType = "windows"
	case common.HostOSTypeEnumAIX:
		OSType = "aix"
	default:
		OSType = enumID
	}
	return
}
