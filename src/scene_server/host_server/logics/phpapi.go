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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"strings"
)

// PHPAPI ee,ce version api
type PHPAPI struct {
	logic  *Logics
	header http.Header
}

// NewPHPAPI return php api struct
func (lgc *Logics) NewPHPAPI(header http.Header) *PHPAPI {
	return &PHPAPI{
		logic:  lgc,
		header: header,
	}
}

func (lgc *Logics) UpdateHost(input map[string]interface{}, appID int64, header http.Header) (interface{}, int, error) {
	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))

	updateData, ok := input["data"]
	if !ok {
		blog.Errorf("params data is required:%v", input)
		return nil, http.StatusBadRequest, defErr.Errorf(common.CCErrCommParamsNeedSet, "data")
	}

	mapData, ok := updateData.(map[string]interface{})
	if !ok {
		blog.Errorf("UpdateHost params data must be object:%v", (input))
		return nil, http.StatusBadRequest, defErr.Errorf(common.CCErrCommParamsInvalid, "data")
	}

	dstPlat, ok := mapData[common.BKSubAreaField]
	if !ok {
		blog.Errorf("params data.bk_cloud_id is require:%v", input)
		return nil, http.StatusBadRequest, defErr.Errorf(common.CCErrCommParamsNeedSet, common.BKSubAreaField)

	}

	innerIP, ok := input["condition"].(map[string]interface{})[common.BKHostInnerIPField]
	if !ok {
		blog.Errorf("params data.bk_ihost_innerip is require:%v", input)
		return nil, http.StatusBadRequest, defErr.Errorf(common.CCErrCommParamsNeedSet, common.BKHostInnerIPField)
	}

	// dst host exist return souccess, hongsong tiyi
	dstHostCondition := map[string]interface{}{
		common.BKHostInnerIPField: innerIP,
		common.BKCloudIDField:     dstPlat,
	}
	phpapi := lgc.NewPHPAPI(header)
	_, hostIDArr, err := phpapi.GetHostMapByCond(dstHostCondition)
	blog.V(3).Infof("hostIDArr:%v", hostIDArr)
	if nil != err {
		blog.Errorf("updateHostMain error:%v", err)
		return nil, http.StatusBadGateway, defErr.Error(common.CCErrHostGetFail)
	}

	if len(hostIDArr) != 0 {
		return nil, 0, nil
	}

	blog.V(3).Infof(" input %s")
	hostCondition := map[string]interface{}{
		common.BKHostInnerIPField: input["condition"].(map[string]interface{})[common.BKHostInnerIPField],
		common.BKCloudIDField:     input["condition"].(map[string]interface{})[common.BKCloudIDField],
	}
	data := input["data"].(map[string]interface{})
	data[common.BKHostInnerIPField] = input["condition"].(map[string]interface{})[common.BKHostInnerIPField]

	res, err := phpapi.UpdateHostMain(hostCondition, data, appID)
	if nil != err {
		blog.Errorf("updateHostMain error:%v", err)
		return nil, http.StatusBadGateway, defErr.Error(common.CCErrHostModifyFail)
	}

	return res, 0, nil

}

func (lgc *Logics) UpdateHostByAppID(input *meta.UpdateHostParams, appID int64, header http.Header) (interface{}, int, error) {

	blog.V(3).Infof("updateHostByAppID http body data: %v", input)
	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))

	phpapi := lgc.NewPHPAPI(header)

	moduleInfo, err := phpapi.GetDefaultModules(int(appID))

	if nil != err {
		blog.Errorf("getDefaultModules input: %v, error:%v, module:%v", input, err, moduleInfo)
		return nil, http.StatusBadGateway, defErr.Error(common.CCErrTopoModuleSelectFailed)
	}

	defaultModuleID, err := moduleInfo.Int64(common.BKModuleIDField)
	if nil != err {
		blog.Errorf("getDefaultModules input: %v, error:%v, module:%v", input, err, moduleInfo)
		return nil, http.StatusBadGateway, defErr.Error(common.CCErrTopoModuleSelectFailed)
	}
	for _, pro := range input.ProxyList {
		proMap := pro.(map[string]interface{})
		var hostID int64
		innerIP := proMap[common.BKHostInnerIPField]
		outerIP, ok := proMap[common.BKHostOuterIPField]
		if !ok {
			outerIP = ""
		}
		hostData, err := phpapi.GetHostByIPAndSource(innerIP.(string), input.CloudID)
		if nil != err {
			blog.Errorf("UpdateHostByAppID getHostByIPAndSource, input:%v, innerip:%v, platID:%v error:%v", input, innerIP, input.CloudID, err)
			return nil, http.StatusBadGateway, defErr.Error(common.CCErrHostGetFail)
		}

		blog.Errorf("hostData:%v", hostData)
		if len(hostData) == 0 {
			platID, ok := proMap[common.BKCloudIDField]
			if ok {
				platConds := common.KvMap{
					common.BKCloudIDField: platID,
				}

				bl, err := lgc.IsPlatExist(header, platConds)
				if nil != err {
					blog.Errorf("is exist plat  error:%s", err.Error())
					return nil, http.StatusBadGateway, defErr.Errorf(common.CCErrTopoGetCloudErrStrFaild, err.Error())
				}
				if !bl {
					blog.Errorf("is exist plat  not foud platid :%v", platID)
					return nil, http.StatusBadGateway, defErr.Error(common.CCErrTopoCloudNotFound)
				}
			}
			proMap["import_from"] = common.HostAddMethodAgent
			blog.V(3).Infof("procMap:%v", proMap)
			hostIDNew, err := phpapi.AddHost(proMap)
			if nil != err {
				blog.Errorf("addHost error:%v", err)
				return nil, http.StatusBadGateway, defErr.Error(common.CCErrHostCreateFail)
			}

			hostID = hostIDNew

			blog.V(3).Infof("addHost success, hostID: %d, input:%v", hostID, input)

			err = phpapi.AddModuleHostConfig(hostID, int64(appID), []int64{defaultModuleID})

			if nil != err {
				blog.Errorf("addModuleHostConfig error:%v, input:%v", err, input)
				return nil, http.StatusBadGateway, defErr.Error(common.CCErrHostTransferModule)
			}

		} else {
			hostID, err = hostData[0].Int64(common.BKHostIDField)
			if nil != err {
				blog.Errorf("UpdateHostByAppID getHostByIPAndSource not found hostid, hostinfo:%v, input:%v, innerip:%v, platID:%v error:%v", hostData[0], input, innerIP, input.CloudID, err)
				return nil, http.StatusBadGateway, defErr.Error(common.CCErrHostGetFail)
			}

		}

		if outerIP != "" {
			hostCondition := map[string]interface{}{
				common.BKHostIDField: hostID,
			}
			data := map[string]interface{}{
				// TODO 没有gse_proxy字段，暂时不修改;2018/03/09
				//common.BKGseProxyField: 1,
			}

			_, err := phpapi.UpdateHostMain(hostCondition, data, appID)
			if nil != err {
				blog.Errorf("updateHostMain error:%v", err)
				return nil, http.StatusBadGateway, defErr.Error(common.CCErrHostModifyFail)
			}
		}

	}

	return nil, 0, nil
}

func (lgc *Logics) GetIPAndProxyByCompany(ipArr []string, cloudID, appID int64, header http.Header) (interface{}, error) {
	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	// 获取不合法的IP列表
	param := &meta.QueryInput{
		Condition: map[string]interface{}{
			common.BKHostInnerIPField: map[string]interface{}{common.BKDBIN: ipArr},
			common.BKCloudIDField:     cloudID,
		},
		Fields: fmt.Sprintf("%s,%s", common.BKHostIDField, common.BKHostInnerIPField),
	}
	phpapi := lgc.NewPHPAPI(header)
	hosts, err := phpapi.GetHostByCond(param)
	if nil != err {
		blog.Errorf("getHostByIPArrAndSource failed, error:%s, ip:%v, cloudID:%d, appID:%d", err.Error(), ipArr, cloudID, appID)
		return nil, defErr.Error(common.CCErrHostGetFail)
	}

	hostIDArr := make([]int64, 0)
	hostMap := make(map[string]mapstr.MapStr)

	for _, host := range hosts {
		hostID, err := host.Int64(common.BKHostIDField)
		if nil != err {
			blog.Errorf("getHostByIPArrAndSource failed, error:%v, ip:%s, cloudID:%d, appID:%d", err.Error(), ipArr, cloudID, appID)
			return nil, defErr.Error(common.CCErrHostGetFail)
		}
		hostIDArr = append(hostIDArr, hostID)
		hostMap[fmt.Sprintf("%v", hostID)] = host
	}

	blog.V(3).Infof("hostIDArr:%v", hostIDArr)
	muduleHostConfigs, err := lgc.GetConfigByCond(header, map[string][]int64{
		common.BKHostIDField: hostIDArr,
	})
	if nil != err {
		blog.Errorf("getHostByIPArrAndSource failed, error:%s, ip:%v, cloudID:%d, appID:%d", err.Error(), ipArr, cloudID, appID)
		return nil, defErr.Errorf(common.CCErrHostModuleConfigFaild, err.Error())
	}

	blog.V(3).Infof("vaildIPArr:%v", muduleHostConfigs)

	validIpArr := make([]interface{}, 0)
	appMap, err := lgc.GetAppMapByCond(header, "", nil)
	if nil != err {
		blog.Errorf("getHostByIPArrAndSource failed, error:%v, ip:%s, cloudID:%d, appID:%d", err.Error(), ipArr, cloudID, appID)
		return nil, defErr.Errorf(common.CCErrHostGetAPPFail, err.Error())
	}

	invalidIpMap := make(map[string]map[string]interface{})

	for _, config := range muduleHostConfigs {
		appIDTemp := fmt.Sprintf("%v", config[common.BKAppIDField])
		appIDIntTemp := config[common.BKAppIDField]
		hostID := config[common.BKHostIDField]
		ip, err := hostMap[fmt.Sprintf("%v", hostID)].String(common.BKHostInnerIPField)
		if nil != err {
			blog.Warnf("getHostByIPArrAndSource get host error, error:%s, appinfo:%v, ip:%v, cloudID:%d, appID:%d", err.Error(), appMap[appIDIntTemp], ipArr, cloudID, appID)
		}

		appName, err := appMap[appIDIntTemp].String(common.BKAppNameField)
		if nil != err {
			blog.Warnf("getHostByIPArrAndSource get appName error, error:%s, appinfo:%v, ip:%v, cloudID:%d, appID:%d", err.Error(), appMap[appIDIntTemp], ipArr, cloudID, appID)
		}

		if appIDIntTemp != appID {

			_, ok := invalidIpMap[appIDTemp]
			if !ok {
				invalidIpMap[appIDTemp] = make(map[string]interface{})
				invalidIpMap[appIDTemp][common.BKAppNameField] = appName
				invalidIpMap[appIDTemp]["ips"] = make([]string, 0)
			}

			invalidIpMap[appIDTemp]["ips"] = append(invalidIpMap[appIDTemp]["ips"].([]string), ip)

		} else {
			validIpArr = append(validIpArr, ip)
		}
	}

	// 获取所有的proxy ip列表
	paramProxy := &meta.QueryInput{
		Condition: map[string]interface{}{
			common.BKGseProxyField: 1,
			common.BKCloudIDField:  cloudID,
		},
		Fields: fmt.Sprintf("%s,%s", common.BKHostIDField, common.BKHostInnerIPField),
	}
	hostProxys, err := phpapi.GetHostByCond(paramProxy)
	if nil != err {
		blog.Errorf("getHostByIPArrAndSource failed, error:%v, ip:%s, cloudID:%d, appID:%d", err.Error(), ipArr, cloudID, appID)
		return nil, defErr.Error(common.CCErrHostGetFail)
	}
	proxyIpArr := make([]interface{}, 0)

	for _, host := range hostProxys {
		h := make(map[string]interface{})
		h[common.BKHostInnerIPField], _ = host.String(common.BKHostInnerIPField)
		h[common.BKHostOuterIPField] = ""
		proxyIpArr = append(proxyIpArr, h)
	}
	blog.V(3).Infof("proxyIpArr:%v", proxyIpArr)

	resData := make(map[string]interface{})
	resData[common.BKIPListField] = validIpArr
	resData[common.BKProxyListField] = proxyIpArr
	resData[common.BKInvalidIPSField] = invalidIpMap
	return resData, nil
}

func (lgc *Logics) UpdateCustomProperty(hostID, appID int64, proeprtyJson map[string]interface{}, header http.Header) (interface{}, error) {

	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))

	phpapi := lgc.NewPHPAPI(header)
	propertys, err := phpapi.GetCustomerPropertyByOwner(common.BKInnerObjIDHost)
	if nil != err {
		blog.Errorf("UpdateCustomProperty error:%s,, hostID:%d, appID:%d, property:%s", hostID, appID, proeprtyJson)
		return nil, defErr.Error(common.CCErrCommSearchPropertyFailed)
	}
	params := make(common.KvMap)
	for _, attrMap := range propertys {
		PropertyId := attrMap.PropertyID

		blog.V(3).Infof("input[PropertyId]:%v", proeprtyJson[PropertyId])
		if _, ok := proeprtyJson[PropertyId]; ok {
			params[PropertyId] = proeprtyJson[PropertyId]
		}
	}
	blog.V(3).Infof("params:%v", params)
	hostCondition := map[string]interface{}{
		common.BKHostIDField: hostID,
	}
	res, err := phpapi.UpdateHostMain(hostCondition, params, appID)
	if nil != err {
		blog.Errorf("UpdateCustomProperty error:%s,, hostID:%d, appID:%d, property:%s", hostID, appID, proeprtyJson)
		return nil, defErr.Error(common.CCErrHostModifyFail)
	}

	return res, nil
}

func (lgc *Logics) CloneHostProperty(input *meta.CloneHostPropertyParams, appID, cloudID int64, header http.Header) (interface{}, error) {
	defError := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))

	condition := common.KvMap{
		common.BKHostInnerIPField: input.OrgIP,
		common.BKCloudIDField:     cloudID,
	}

	phpapi := lgc.NewPHPAPI(header)
	// 处理源IP
	hostMap, hostIdArr, err := phpapi.GetHostMapByCond(condition)

	blog.V(3).Infof("CloneHostPropertyhostMapData:%v", hostMap)
	if err != nil {
		blog.Errorf("CloneHostPropertygetHostMapByCond error : %v, input:%v", err, input)

		return nil, defError.Error(common.CCErrHostDetailFail)
	}

	if len(hostIdArr) == 0 {
		blog.Errorf("CloneHostProperty clone host getHostMapByCond not found  input:%v", input)
		return nil, defError.Error(common.CCErrHostDetailFail)
	}
	hostMapData, ok := hostMap[hostIdArr[0]]
	if false == ok {
		blog.Errorf("CloneHostProperty getHostMapByCond not source ip , raw data format hostMap:%v, input:%v", hostMap, input)
		return nil, defError.Error(common.CCErrHostDetailFail)
	}

	srcHostID, err := util.GetInt64ByInterface(hostMapData[common.BKHostIDField])
	if nil != err {
		blog.Errorf("CloneHostProperty clone source host host id  not found hostmap:%v input:%v", hostMapData, input)
		return nil, defError.Error(common.CCErrHostDetailFail)
	}
	configCond := map[string][]int64{
		common.BKHostIDField: []int64{srcHostID},
		common.BKAppIDField:  []int64{appID},
	}
	// 判断源IP是否存在
	configDataArr, err := lgc.GetConfigByCond(header, configCond)
	blog.V(3).Infof("configData:%v", configDataArr)
	if nil != err {
		blog.Errorf("CloneHostProperty clone host property error : %v, input:%v", err, input)
		return nil, defError.Errorf(common.CCErrHostModuleConfigFaild, err.Error())
	}
	if len(configDataArr) == 0 {
		blog.Errorf("CloneHostProperty clone host property error not found src host  input:%v", input)
		return nil, defError.Error(common.CCErrCommNotFound)
	}

	// 处理目标IP
	dstIpArr := strings.Split(input.DstIP, ",")
	// 获得已存在的主机
	dstCondition := map[string]interface{}{
		common.BKHostInnerIPField: map[string]interface{}{
			common.BKDBIN: dstIpArr,
		},
		common.BKCloudIDField: cloudID,
	}

	dstHostMap, dstHostIdArr, err := phpapi.GetHostMapByCond(dstCondition)
	blog.V(3).Infof("dstHostMap:%v, input:%v", dstHostMap, input)

	dstConfigCond := map[string][]int64{
		common.BKAppIDField:  []int64{appID},
		common.BKHostIDField: dstHostIdArr,
	}
	dstHostIdArrV, err := lgc.GetHostIDByCond(header, dstConfigCond)
	existIPMap := make(map[string]int64, 0)
	for _, id := range dstHostIdArrV {
		if dstHostMapData, ok := dstHostMap[id]; ok {
			ip, ok := dstHostMapData[common.BKHostInnerIPField].(string)
			if false == ok {
				blog.Errorf("CloneHostProperty not found innerip , raw data format hostMap:%v, input:%v", dstHostMapData, input)
				return nil, defError.Error(common.CCErrHostDetailFail)

			}

			hostID, err := util.GetInt64ByInterface(dstHostMapData[common.BKHostIDField])
			if nil != err {
				blog.Errorf("CloneHostProperty not found host id  , raw data format hostMap:%v, input:%v", dstHostMapData, input)
				return nil, defError.Error(common.CCErrHostDetailFail)
			}
			existIPMap[ip] = hostID
		} else {
			blog.Errorf("CloneHostProperty not host id , host id:%v, input:%v", id, input)
			return nil, defError.Error(common.CCErrHostDetailFail)
		}
	}

	//更新的时候，不修改为nil的数据
	updateHostData := make(map[string]interface{})
	for key, val := range hostMapData {
		if nil != val {
			updateHostData[key] = val
		}
	}
	// remote duplication ip
	dstIPMap := make(map[string]bool, len(dstIpArr))
	for _, ip := range dstIpArr {
		dstIPMap[ip] = true
	}

	blog.V(3).Infof("configData[0]:%v, input:%v", configDataArr[0], input)
	moduleIDs := make([]int64, 0)
	for _, configData := range configDataArr {

		moduleID, err := util.GetInt64ByInterface(configData[common.BKModuleIDField])
		if nil != err {
			blog.Errorf("CloneHostProperty not host module relation error, not found module id: raw config:%v, input:%v", configData, input)
			return nil, defError.Error(common.CCErrGetOriginHostModuelRelationship)
		}
		moduleIDs = append(moduleIDs, moduleID)
	}

	// 克隆主机, 已存在的修改，不存在的新增；dstIpArr: 全部要克隆的主机，existIpArr：已存在的要克隆的主机
	blog.V(3).Infof("existIpArr:%v, input:%v", existIPMap, input)
	for dstIpV, _ := range dstIPMap {
		if dstIpV == input.OrgIP {
			blog.V(3).Infof("clone host updateHostMain err:dstIp and orgIp cannot be the same,srcIP:%s, dstIP:%s, input:%v", input.OrgIP, dstIpV, input)
			continue
		}
		blog.V(3).Infof("hostMapData:%v", hostMapData)
		hostID, oK := existIPMap[dstIpV]
		if true == oK {
			blog.V(3).Infof("clone update")
			hostCondition := map[string]interface{}{
				common.BKHostInnerIPField: dstIpV,
				common.BKHostIDField:      hostID,
			}

			updateHostData[common.BKHostInnerIPField] = dstIpV
			delete(updateHostData, common.BKHostIDField)
			res, err := phpapi.UpdateHostMain(hostCondition, updateHostData, appID)
			if nil != err {
				blog.Errorf("CloneHostProperty  update dst host error, error:%s, currentIP:%s, input:%v", err.Error(), dstIpV, input)
				return nil, defError.Error(common.CCErrHostModifyFail)
			}
			blog.V(3).Infof("CloneHostPropertyclone host updateHostMain res:%v", res)
			params := new(meta.ModuleHostConfigParams)
			params.HostID = hostID
			params.ApplicationID = appID

			resDelRelation, err := lgc.CoreAPI.HostController().Module().DelModuleHostConfig(context.Background(), header, params)
			if nil != err || (nil == err && false == resDelRelation.Result) {
				if nil == err {
					err = fmt.Errorf(resDelRelation.ErrMsg)
				}
				blog.Errorf("CloneHostProperty remove hosthostconfig error, params:%v, error:%s, input:%v", params, err.Error(), input)
				return nil, defError.Error(common.CCErrHostTransferModule)
			}
		} else {
			hostMapData[common.BKHostInnerIPField] = dstIpV
			addHostMapData := hostMapData
			delete(addHostMapData, common.BKHostIDField)
			cloneHostId, err := phpapi.AddHost(addHostMapData)
			if nil != err {
				blog.Errorf("CloneHostProperty remove hosthostconfig error, addHostMapData:%v, error:%s, input:%v", addHostMapData, err.Error(), input)
				return nil, defError.Error(common.CCErrHostCreateFail)
			}
			blog.V(3).Infof("CloneHostProperty dstIP:%s, cloneHostId:%v, input:%v", dstIpV, cloneHostId, input)
			hostID = cloneHostId

		}
		err := phpapi.AddModuleHostConfig(hostID, appID, moduleIDs)
		if nil != err {
			blog.Errorf("CloneHostProperty remove hosthostconfig error, hostID:%d, moduleID:%v, appID:%d, error:%s, input:%v", hostID, moduleIDs, appID, err.Error(), input)
			return nil, defError.Error(common.CCErrHostModuleRelationAddFailed)
		}
	}

	return nil, nil
}
