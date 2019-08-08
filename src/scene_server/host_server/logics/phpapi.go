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
	"strings"

	"github.com/rs/xid"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// PHPAPI ee,ce version api
type PHPAPI struct {
	logic  *Logics
	header http.Header
	rid    string
	ccErr  errors.DefaultCCErrorIf
}

// NewPHPAPI return php api struct
func (lgc *Logics) NewPHPAPI() *PHPAPI {
	return &PHPAPI{
		logic:  lgc,
		header: lgc.header,
		rid:    lgc.rid,
		ccErr:  lgc.ccErr,
	}
}

func (lgc *Logics) UpdateHost(ctx context.Context, input map[string]interface{}, appID int64) (interface{}, int, errors.CCError) {

	updateData, ok := input["data"]
	if !ok {
		blog.Errorf("params data is required, input:%+v, rid:%s", input, lgc.rid)
		return nil, http.StatusBadRequest, lgc.ccErr.Errorf(common.CCErrCommParamsNeedSet, "data")
	}

	mapData, ok := updateData.(map[string]interface{})
	if !ok {
		blog.Errorf("UpdateHost params data must be object, input :%+v, rid:%s", input, lgc.rid)
		return nil, http.StatusBadRequest, lgc.ccErr.Errorf(common.CCErrCommParamsInvalid, "data")
	}

	dstPlat, ok := mapData[common.BKSubAreaField]
	if !ok {
		blog.Errorf("params data.bk_cloud_id is require, input:%+v, rid:%s", input, lgc.rid)
		return nil, http.StatusBadRequest, lgc.ccErr.Errorf(common.CCErrCommParamsNeedSet, common.BKSubAreaField)

	}

	innerIP, ok := input["condition"].(map[string]interface{})[common.BKHostInnerIPField]
	if !ok {
		blog.Errorf("params data.bk_ihost_innerip is require, input:%+v, rid:%s", input, lgc.rid)
		return nil, http.StatusBadRequest, lgc.ccErr.Errorf(common.CCErrCommParamsNeedSet, common.BKHostInnerIPField)
	}

	// dst host exist return success
	dstHostCondition := map[string]interface{}{
		common.BKHostInnerIPField: innerIP,
		common.BKCloudIDField:     dstPlat,
	}
	phpApi := lgc.NewPHPAPI()
	_, hostIDArr, err := phpApi.GetHostMapByCond(ctx, dstHostCondition)
	blog.V(5).Infof("hostIDArr:%+v,rid:%s", hostIDArr, lgc.rid)
	if nil != err {
		return nil, http.StatusInternalServerError, err
	}

	if len(hostIDArr) != 0 {
		return nil, 0, nil
	}

	hostCondition := map[string]interface{}{
		common.BKHostInnerIPField: input["condition"].(map[string]interface{})[common.BKHostInnerIPField],
		common.BKCloudIDField:     input["condition"].(map[string]interface{})[common.BKCloudIDField],
	}
	data := input["data"].(map[string]interface{})
	data[common.BKHostInnerIPField] = input["condition"].(map[string]interface{})[common.BKHostInnerIPField]

	res, err := phpApi.UpdateHostMain(ctx, hostCondition, data, appID)
	if nil != err {
		return nil, http.StatusInternalServerError, err
	}

	return res, 0, nil

}

// FindHostIDsByAppID is just a copy from which used for verify hostID before update host
func (lgc *Logics) FindHostIDsByAppID(ctx context.Context, input *meta.UpdateHostParams, appID int64) ([]int64, int, error) {

	blog.V(5).Infof("updateHostByAppID http body data: %+v, rid:%s", input, lgc.rid)

	phpApi := lgc.NewPHPAPI()

	moduleInfo, err := phpApi.GetDefaultModules(ctx, int(appID))

	if nil != err {
		blog.Errorf("getDefaultModules input: %v, error:%v, module:%v, rid: %s", input, err, moduleInfo, lgc.rid)
		return nil, http.StatusInternalServerError, lgc.ccErr.Error(common.CCErrTopoModuleSelectFailed)
	}

	defaultModuleID, err := moduleInfo.Int64(common.BKModuleIDField)
	if nil != err {
		blog.Errorf("getDefaultModules input: %v, error:%v, module:%v, rid: %s", input, err, moduleInfo, lgc.rid)
		return nil, http.StatusInternalServerError, lgc.ccErr.Error(common.CCErrTopoModuleSelectFailed)
	}
	hostIDArr := make([]int64, 0)
	for _, pro := range input.ProxyList {
		proMap := pro.(map[string]interface{})
		var hostID int64
		innerIP := proMap[common.BKHostInnerIPField]
		hostData, err := phpApi.GetHostByIPAndSource(ctx, innerIP.(string), input.CloudID)
		if nil != err {
			return nil, http.StatusInternalServerError, err
		}

		blog.V(5).Infof("hostData:%v, rid:%s", hostData, lgc.rid)
		if len(hostData) == 0 {
			platID, ok := proMap[common.BKCloudIDField]
			if ok {
				platConds := mapstr.MapStr{
					common.BKCloudIDField: platID,
				}

				bl, err := lgc.IsPlatExist(ctx, platConds)
				if nil != err {
					return nil, http.StatusInternalServerError, err
				}
				if !bl {
					blog.Errorf("is exist plat  not found platid :%v, input:%+v,rid:%s", platID, input, lgc.rid)
					return nil, http.StatusInternalServerError, lgc.ccErr.Error(common.CCErrTopoCloudNotFound)
				}
			}
			proMap["import_from"] = common.HostAddMethodAgent
			blog.V(5).Infof("procMap:%v, input:%+v,rid:%s", proMap, input, lgc.rid)
			hostIDNew, err := phpApi.AddHost(ctx, proMap)
			if nil != err {
				return nil, http.StatusInternalServerError, err
			}

			hostID = hostIDNew

			blog.V(5).Infof("addHost success, hostID: %d, input:%v,rid:%s", hostID, input, lgc.rid)
			err = phpApi.AddModuleHostConfig(ctx, hostID, int64(appID), []int64{defaultModuleID})
			if nil != err {
				return nil, http.StatusInternalServerError, err
			}

		} else {
			hostID, err = hostData[0].Int64(common.BKHostIDField)
			if nil != err {
				blog.Errorf("UpdateHostByAppID failed, getHostByIPAndSource result not found, hostInfo: %+v, input:%v, innerIP:%v, platID:%v error:%s, rid:%s",
					hostData[0], input, innerIP, input.CloudID, err.Error(), lgc.rid)
				return nil, http.StatusInternalServerError, lgc.ccErr.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDHost, common.BKHostIDField, "int", err.Error())
			}

		}
		hostIDArr = append(hostIDArr, hostID)
	}

	return hostIDArr, 0, nil
}

func (lgc *Logics) UpdateHostByAppID(ctx context.Context, input *meta.UpdateHostParams, appID int64) (interface{}, int, error) {

	blog.V(5).Infof("updateHostByAppID http body data: %+v, rid:%s", input, lgc.rid)

	phpApi := lgc.NewPHPAPI()

	moduleInfo, err := phpApi.GetDefaultModules(ctx, int(appID))

	if nil != err {
		blog.Errorf("getDefaultModules input: %v, error:%v, module:%v, rid: %s", input, err, moduleInfo, lgc.rid)
		return nil, http.StatusInternalServerError, lgc.ccErr.Error(common.CCErrTopoModuleSelectFailed)
	}

	defaultModuleID, err := moduleInfo.Int64(common.BKModuleIDField)
	if nil != err {
		blog.Errorf("getDefaultModules input: %v, error:%v, module:%v, rid: %s", input, err, moduleInfo, lgc.rid)
		return nil, http.StatusInternalServerError, lgc.ccErr.Error(common.CCErrTopoModuleSelectFailed)
	}
	for _, pro := range input.ProxyList {
		proMap := pro.(map[string]interface{})
		var hostID int64
		innerIP := proMap[common.BKHostInnerIPField]
		outerIP, ok := proMap[common.BKHostOuterIPField]
		if !ok {
			outerIP = ""
		}
		hostData, err := phpApi.GetHostByIPAndSource(ctx, innerIP.(string), input.CloudID)
		if nil != err {
			return nil, http.StatusInternalServerError, err
		}

		blog.V(5).Infof("hostData:%v, rid:%s", hostData, lgc.rid)
		if len(hostData) == 0 {
			platID, ok := proMap[common.BKCloudIDField]
			if ok {
				platConds := mapstr.MapStr{
					common.BKCloudIDField: platID,
				}

				bl, err := lgc.IsPlatExist(ctx, platConds)
				if nil != err {
					return nil, http.StatusInternalServerError, err
				}
				if !bl {
					blog.Errorf("is exist plat  not found platID :%v, input:%+v,rid:%s", platID, input, lgc.rid)
					return nil, http.StatusInternalServerError, lgc.ccErr.Error(common.CCErrTopoCloudNotFound)
				}
			}
			proMap["import_from"] = common.HostAddMethodAgent
			blog.V(5).Infof("procMap:%v, input:%+v,rid:%s", proMap, input, lgc.rid)
			hostIDNew, err := phpApi.AddHost(ctx, proMap)
			if nil != err {
				return nil, http.StatusInternalServerError, err
			}

			hostID = hostIDNew

			blog.V(5).Infof("addHost success, hostID: %d, input:%v,rid:%s", hostID, input, lgc.rid)
			err = phpApi.AddModuleHostConfig(ctx, hostID, int64(appID), []int64{defaultModuleID})
			if nil != err {
				return nil, http.StatusInternalServerError, err
			}

		} else {
			hostID, err = hostData[0].Int64(common.BKHostIDField)
			if nil != err {
				blog.Errorf("UpdateHostByAppID getHostByIPAndSource not found hostID, hostInfo:%v, input:%v, innerIP:%v, platID:%v error:%s, rid:%s", hostData[0], input, innerIP, input.CloudID, err.Error(), lgc.rid)
				return nil, http.StatusInternalServerError, lgc.ccErr.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDHost, common.BKHostIDField, "int", err.Error())
			}

		}

		if outerIP != "" {
			hostCondition := map[string]interface{}{
				common.BKHostIDField: hostID,
			}
			data := map[string]interface{}{
				// TODO 没有gse_proxy字段，暂时不修改;2018/03/09
				// common.BKGseProxyField: 1,
			}

			_, err := phpApi.UpdateHostMain(ctx, hostCondition, data, appID)
			if nil != err {
				return nil, http.StatusInternalServerError, err
			}
		}

	}

	return nil, 0, nil
}

func (lgc *Logics) GetIPAndProxyByCompany(ctx context.Context, ipArr []string, cloudID, appID int64) (interface{}, error) {
	// 获取不合法的IP列表
	param := &meta.QueryCondition{
		Condition: mapstr.MapStr{
			common.BKHostInnerIPField: mapstr.MapStr{common.BKDBIN: ipArr},
			common.BKCloudIDField:     cloudID,
		},
		Fields: []string{common.BKHostIDField, common.BKHostInnerIPField},
	}
	phpApi := lgc.NewPHPAPI()
	hosts, err := phpApi.GetHostByCond(ctx, param)
	if nil != err {
		return nil, err
	}

	hostIDArr := make([]int64, 0)
	hostMap := make(map[string]mapstr.MapStr)

	for _, host := range hosts {
		hostID, err := host.Int64(common.BKHostIDField)
		if nil != err {
			blog.Errorf("GetIPAndProxyByCompany hostID not integer, error:%v, ip:%s, cloudID:%d, appID:%d, hostInfo:%+v,rid:%s", err.Error(), ipArr, cloudID, appID, host, lgc.rid)
			return nil, lgc.ccErr.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDHost, common.BKHostIDField, "int", err.Error())
		}
		hostIDArr = append(hostIDArr, hostID)
		hostMap[fmt.Sprintf("%v", hostID)] = host
	}

	blog.V(5).Infof("hostIDArr:%v,rid:%s", hostIDArr, lgc.rid)
	moduleHostConfigs, err := lgc.GetConfigByCond(ctx, meta.HostModuleRelationRequest{
		HostIDArr: hostIDArr,
	})
	if nil != err {
		return nil, err
	}

	blog.V(5).Infof("validIPArr:%v,rid:%s", moduleHostConfigs, lgc.rid)

	validIpArr := make([]interface{}, 0)
	appMap, err := lgc.GetAppMapByCond(ctx, nil, nil)
	if nil != err {
		return nil, err
	}

	invalidIpMap := make(map[string]map[string]interface{})

	for _, config := range moduleHostConfigs {
		appIDTemp := fmt.Sprintf("%v", config.AppID)
		appIDIntTemp := config.AppID
		hostID := config.HostID
		ip, err := hostMap[fmt.Sprintf("%v", hostID)].String(common.BKHostInnerIPField)
		if nil != err {
			blog.Warnf("getHostByIPArrAndSource get host error, error:%s, appInfo:%v, ip:%v, cloudID:%d, appID:%d,rid:%s", err.Error(), appMap[appIDIntTemp], ipArr, cloudID, appID, lgc.rid)
		}

		appName, err := appMap[appIDIntTemp].String(common.BKAppNameField)
		if nil != err {
			blog.Warnf("getHostByIPArrAndSource get appName error, error:%s, appInfo:%v, ip:%v, cloudID:%d, appID:%d,rid:%s", err.Error(), appMap[appIDIntTemp], ipArr, cloudID, appID, lgc.rid)
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
	paramProxy := &meta.QueryCondition{
		Condition: mapstr.MapStr{
			common.BKGseProxyField: 1,
			common.BKCloudIDField:  cloudID,
		},
		Fields: []string{common.BKHostIDField, common.BKHostInnerIPField},
	}
	hostProxys, err := phpApi.GetHostByCond(ctx, paramProxy)
	if nil != err {
		return nil, err
	}
	proxyIpArr := make([]interface{}, 0)

	for _, host := range hostProxys {
		h := make(map[string]interface{})
		h[common.BKHostInnerIPField], _ = host.String(common.BKHostInnerIPField)
		h[common.BKHostOuterIPField] = ""
		proxyIpArr = append(proxyIpArr, h)
	}
	blog.V(5).Infof("proxyIpArr:%v,rid:%s", proxyIpArr, lgc.rid)

	resData := make(map[string]interface{})
	resData[common.BKIPListField] = validIpArr
	resData[common.BKProxyListField] = proxyIpArr
	resData[common.BKInvalidIPSField] = invalidIpMap
	return resData, nil
}

func (lgc *Logics) UpdateCustomProperty(ctx context.Context, hostID, appID int64, propertyJson map[string]interface{}) (interface{}, error) {

	phpApi := lgc.NewPHPAPI()
	propertys, err := phpApi.GetCustomerPropertyByOwner(ctx, common.BKInnerObjIDHost)
	if nil != err {
		return nil, err
	}
	params := make(common.KvMap)
	for _, attrMap := range propertys {
		PropertyId := attrMap.PropertyID

		blog.V(5).Infof("input[PropertyId]:%v, rid:%s", propertyJson[PropertyId], lgc.rid)
		if _, ok := propertyJson[PropertyId]; ok {
			params[PropertyId] = propertyJson[PropertyId]
		}
	}
	blog.V(5).Infof("params:%v,rid:%s", params, lgc.rid)
	hostCondition := map[string]interface{}{
		common.BKHostIDField: hostID,
	}
	res, err := phpApi.UpdateHostMain(ctx, hostCondition, params, appID)
	if nil != err {
		return nil, err
	}

	return res, nil
}

func (lgc *Logics) CloneHostProperty(ctx context.Context, input *meta.CloneHostPropertyParams, appID, cloudID int64) (interface{}, error) {

	condition := common.KvMap{
		common.BKHostInnerIPField: input.OrgIP,
		common.BKCloudIDField:     cloudID,
	}

	appConf := map[string]interface{}{common.BKAppIDField: input.AppID}
	appInfo, err := lgc.GetAppDetails(ctx, "", appConf)
	if err != nil {
		return nil, err
	}
	if len(appInfo) == 0 {
		return nil, lgc.ccErr.Errorf(common.CCErrCoreServiceBusinessNotExist, input.AppID)
	}

	phpApi := lgc.NewPHPAPI()
	// 处理源IP
	hostMap, hostIdArr, err := phpApi.GetHostMapByCond(ctx, condition)

	blog.V(5).Infof("CloneHostProperty hostMapData:%v,rid:%s", hostMap, lgc.rid)
	if err != nil {
		return nil, err
	}

	if len(hostIdArr) == 0 {
		blog.Errorf("CloneHostProperty clone host getHostMapByCond not found  input:%+v,rid:%s", input, lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrHostNotFound)
	}
	hostMapData, ok := hostMap[hostIdArr[0]]
	if false == ok {
		blog.Errorf("CloneHostProperty getHostMapByCond not source ip , raw data format hostMap:%+v, input:%+v,rid:%s", hostMap, input, lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrHostDetailFail)
	}

	srcHostID, err := util.GetInt64ByInterface(hostMapData[common.BKHostIDField])
	if nil != err {
		blog.Errorf("CloneHostProperty clone source host host id  not found hostmap:%+v input:%+v,rid:%s", hostMapData, input, lgc.rid)
		return nil, lgc.ccErr.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDHost, common.BKHostIDField, "int", err.Error())
	}
	configCond := meta.HostModuleRelationRequest{
		HostIDArr:     []int64{srcHostID},
		ApplicationID: appID,
	}
	// 判断源IP是否存在
	configDataArr, err := lgc.GetConfigByCond(ctx, configCond)
	blog.V(5).Infof("configData:%+v,rid:%s", configDataArr, lgc.rid)
	if nil != err {
		return nil, err
	}
	if len(configDataArr) == 0 {
		blog.Errorf("CloneHostProperty clone host property error not found src host  input:%+v, param:%+v,rid:%s", input, configCond, lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrCommNotFound)
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

	dstHostMap, dstHostIDArr, err := phpApi.GetHostMapByCond(ctx, dstCondition)
	blog.V(5).Infof("dstHostMap:%+v, input:%+v,rid:%s", dstHostMap, input, lgc.rid)

	var dstHostIDArrV []int64
	if len(dstHostIDArr) > 0 {
		dstConfigCond := meta.HostModuleRelationRequest{
			ApplicationID: appID,
			HostIDArr:     dstHostIDArr,
		}
		dstHostIDArrV, err = lgc.GetHostIDByCond(ctx, dstConfigCond)
		if err != nil {
			return nil, err
		}
	}

	existIPMap := make(map[string]int64, 0)
	for _, id := range dstHostIDArrV {
		if dstHostMapData, ok := dstHostMap[id]; ok {
			ip, ok := dstHostMapData[common.BKHostInnerIPField].(string)
			if false == ok {
				blog.Errorf("CloneHostProperty not found innerIP , raw data format hostMap:%+v, input:%+v, rid:%s", dstHostMapData, input, lgc.rid)
				return nil, lgc.ccErr.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDHost, common.BKHostInnerIPField, "string", "convert fail")

			}

			hostID, err := util.GetInt64ByInterface(dstHostMapData[common.BKHostIDField])
			if nil != err {
				blog.Errorf("CloneHostProperty not found host id  , raw data format hostMap:%+v, input:%+v, rid:%s", dstHostMapData, input, lgc.rid)
				return nil, lgc.ccErr.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDHost, common.BKHostIDField, "int", "convert fail")
			}
			existIPMap[ip] = hostID
		} else {
			blog.Errorf("CloneHostProperty not host id , host id:%+v, hostMap:%+v, input:%+v,rid:%s", id, dstHostMapData, input, lgc.rid)
			return nil, lgc.ccErr.Error(common.CCErrHostDetailFail)
		}
	}

	hostMapData, err = lgc.removeHostBadField(ctx, hostMapData)
	if nil != err {
		blog.Errorf("CloneHostProperty clone host property error : %v, input:%#v,rid:%s", err, input, lgc.rid)
		return nil, lgc.ccErr.Errorf(common.CCErrHostDetailFail, err.Error())
	}
	// 更新的时候，不修改为nil的数据
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

	blog.V(5).Infof("configData[0]:%+v, input:%+v, rid: %s", configDataArr[0], input, lgc.rid)
	moduleIDs := make([]int64, 0)
	for _, configData := range configDataArr {

		moduleIDs = append(moduleIDs, configData.ModuleID)
	}

	// 克隆主机, 已存在的修改，不存在的新增；dstIpArr: 全部要克隆的主机，existIpArr：已存在的要克隆的主机
	blog.V(5).Infof("existIpArr:%+v, input:%+v,rid:%s", existIPMap, input, lgc.rid)
	for dstIPV := range dstIPMap {
		if dstIPV == input.OrgIP {
			blog.V(5).Infof("clone host updateHostMain err:dstIp and orgIp cannot be the same,srcIP:%s, dstIP:%s, input:%+v,rid:%s", input.OrgIP, dstIPV, input, lgc.rid)
			continue
		}
		blog.V(5).Infof("hostMapData:%+v,rid:%s", hostMapData, lgc.rid)
		hostID, oK := existIPMap[dstIPV]
		if true == oK {
			blog.V(5).Infof("clone update, rid: %s", lgc.rid)
			hostCondition := map[string]interface{}{
				common.BKHostInnerIPField: dstIPV,
				common.BKHostIDField:      hostID,
			}

			updateHostData[common.BKHostInnerIPField] = dstIPV
			delete(updateHostData, common.BKHostIDField)
			delete(updateHostData, common.BKAssetIDField)
			delete(updateHostData, common.CreateTimeField)
			res, err := phpApi.UpdateHostMain(ctx, hostCondition, updateHostData, appID)
			if nil != err {
				return nil, err
			}
			blog.V(5).Infof("CloneHostProperty clone host updateHostMain res:%v, rid:%s", res, lgc.rid)
		} else {
			hostMapData[common.BKHostInnerIPField] = dstIPV
			addHostMapData := hostMapData
			delete(addHostMapData, common.BKHostIDField)
			delete(addHostMapData, common.CreateTimeField)
			addHostMapData[common.BKAssetIDField] = xid.New().String()
			cloneHostID, err := phpApi.AddHost(ctx, addHostMapData)
			if nil != err {
				return nil, err
			}
			blog.V(5).Infof("CloneHostProperty dstIP:%s, cloneHostId:%+v, input:%+v,rid:%s", dstIPV, cloneHostID, input, lgc.rid)
			hostID = cloneHostID

		}
		err := phpApi.AddModuleHostConfig(ctx, hostID, appID, moduleIDs)
		if nil != err {
			return nil, err
		}
	}

	return nil, nil
}

// removeHostBadField remove host bad field, host module delete field
func (lgc *Logics) removeHostBadField(ctx context.Context, hostInfo map[string]interface{}) (mapstr.MapStr, error) {
	defError := lgc.ccErr

	newHostInfo := mapstr.New()
	hostAttributeArr, err := lgc.GetHostAttributes(ctx, lgc.ownerID, nil)
	if err != nil {
		blog.Errorf("CloneHostProperty GetHostAttributes, err:%s, rid:%s", err.Error(), lgc.rid)
		return nil, defError.Error(common.CCErrHostDetailFail)
	}
	hostAttributeMap := make(map[string]string, 0)
	for _, attr := range hostAttributeArr {
		hostAttributeMap[attr.PropertyID] = attr.PropertyID
	}
	// delete bad field
	for key, val := range hostInfo {
		_, ok := hostAttributeMap[key]
		if ok {
			newHostInfo.Set(key, val)
		}
	}
	return newHostInfo, nil
}
