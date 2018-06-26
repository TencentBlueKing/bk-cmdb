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
	"errors"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	types "configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (phpapi *PHPAPI) GetDefaultModules(appID int) (types.MapStr, error) {

	param := &meta.QueryInput{
		Condition: map[string]interface{}{
			common.BKAppIDField:   appID,
			common.BKDefaultField: 1,
		},
		Fields: fmt.Sprintf("%s,%s", common.BKSetIDField, common.BKModuleIDField),
	}

	resMap, err := phpapi.getObjByCondition(param, common.BKInnerObjIDModule)

	if nil != err {
		return nil, err
	}

	blog.V(3).Infof("getDefaultModules complete, res: %v", resMap)

	if false == resMap.Result {
		return nil, errors.New(resMap.ErrMsg)
	}

	if resMap.Data.Count == 0 {
		return nil, errors.New(fmt.Sprintf("can not found default module, appid: %d", appID))
	}

	return resMap.Data.Info[0], nil

}

func (phpapi *PHPAPI) GetHostByIPAndSource(innerIP string, platID int64) ([]types.MapStr, error) {

	param := &meta.QueryInput{
		Condition: map[string]interface{}{
			common.BKHostInnerIPField: innerIP,
			common.BKCloudIDField:     platID,
		},
		Fields: common.BKHostIDField,
	}

	resMap, err := phpapi.getObjByCondition(param, common.BKInnerObjIDHost)

	if nil != err {
		return nil, err
	}

	if !resMap.Result {
		return nil, errors.New(resMap.ErrMsg)
	}

	blog.V(3).Infof("getHostByIPAndSource res: %v", resMap)

	return resMap.Data.Info, nil
}

func (phpapi *PHPAPI) GetHostByCond(param *meta.QueryInput) ([]types.MapStr, error) {
	blog.V(3).Infof("GetHostByCond param:%v", param)
	resMap, err := phpapi.getObjByCondition(param, common.BKInnerObjIDHost)
	if nil != err {
		return nil, err
	}

	if false == resMap.Result {
		return nil, errors.New(resMap.ErrMsg)
	}

	blog.V(3).Infof("getHostByIPArrAndSource res: %v", resMap)
	return resMap.Data.Info, nil
}

//search host helpers
func (phpapi *PHPAPI) GetHostMapByCond(condition map[string]interface{}) (map[int64]map[string]interface{}, []int64, error) {
	hostMap := make(map[int64]map[string]interface{})
	hostIDArr := make([]int64, 0)

	// build host controller url
	searchParams := &meta.QueryInput{
		Fields:    "",
		Condition: condition,
	}
	res, err := phpapi.logic.CoreAPI.HostController().Host().GetHosts(context.Background(), phpapi.header, searchParams)

	if nil != err {
		blog.Errorf("getHostMapByCond error params:%s, error:%s", condition, err.Error())
		return hostMap, hostIDArr, err
	}

	blog.V(3).Infof("appInfo:%v", res)

	if false == res.Result {
		return nil, nil, errors.New(res.ErrMsg)
	}

	for _, host := range res.Data.Info {

		host_id, err := util.GetInt64ByInterface(host[common.BKHostIDField])
		if nil != err {
			return nil, nil, err
		}

		hostMap[host_id] = host
		hostIDArr = append(hostIDArr, host_id)
	}
	return hostMap, hostIDArr, nil
}

// GetHostDataByConfig  get host info
func (phpapi *PHPAPI) GetHostDataByConfig(configData []map[string]int64) ([]interface{}, error) {
	hostIDArr := make([]int64, 0)
	for _, config := range configData {
		hostIDArr = append(hostIDArr, config[common.BKHostIDField])
	}

	hostMapCondition := map[string]interface{}{
		common.BKHostIDField: map[string]interface{}{
			common.BKDBIN: hostIDArr,
		},
	}

	hostMap, _, err := phpapi.GetHostMapByCond(hostMapCondition)
	if nil != err {
		return nil, err
	}

	hostData, err := phpapi.SetHostData(configData, hostMap)
	if nil != err {
		return hostData, err
	}

	return hostData, nil
}

func (phpapi *PHPAPI) GetCustomerPropertyByOwner(objType string) ([]meta.Attribute, error) {

	blog.V(3).Infof("getCustomerPropertyByOwner start")
	searchBody := make(map[string]interface{})
	searchBody[common.BKObjIDField] = common.BKInnerObjIDHost
	searchBody[common.BKOwnerIDField] = util.GetOwnerID(phpapi.header)
	res, err := phpapi.logic.CoreAPI.ObjectController().Meta().SelectObjectAttWithParams(context.Background(), phpapi.header, searchBody)
	if nil != err {
		blog.Errorf("GetHostDetailById  attr error :%v", err)
		return nil, err
	}

	if false == res.Result {
		blog.Errorf("GetHostDetailById  attr error :%v", err)
		return nil, fmt.Errorf(res.ErrMsg)
	}
	customAttrArr := make([]meta.Attribute, 0)
	for _, attr := range res.Data { //hostAttrArr {
		if false == attr.IsPre {
			customAttrArr = append(customAttrArr, attr)
		}
	}
	return customAttrArr, nil
}

// In_existIpArr exsit ip in array
func (phpapi *PHPAPI) In_existIpArr(arr []string, ip string) bool {
	for _, v := range arr {
		if ip == v {
			return true
		}
	}
	return false
}

func (phpapi *PHPAPI) getObjByCondition(dat *meta.QueryInput, objType string) (*meta.QueryInstResult, error) {

	res, err := phpapi.logic.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), objType, phpapi.header, dat)
	if nil != err {
		return nil, err
	}

	return res, nil
}
