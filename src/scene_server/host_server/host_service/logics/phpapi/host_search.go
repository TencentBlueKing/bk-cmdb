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

package phpapi

import (
	"encoding/json"
	"errors"
	"fmt"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/util"
)

func GetDefaultModules(req *restful.Request, appID int, objURL string) (interface{}, error) {

	param := map[string]interface{}{
		"condition": map[string]interface{}{
			common.BKAppIDField:   appID,
			common.BKDefaultField: 1,
		},
		"fields": fmt.Sprintf("%s,%s", common.BKSetIDField, common.BKModuleIDField),
	}

	resMap, err := getObjByCondition(req, param, common.BKInnerObjIDModule, objURL)

	if nil != err {
		return nil, err
	}

	blog.Debug("getDefaultModules complete, res: %v", resMap)

	if !resMap["result"].(bool) {
		return nil, errors.New(resMap["message"].(string))
	}

	resDataMap := resMap["data"].(map[string]interface{})

	if resDataMap["count"] == 0 {
		return nil, errors.New(fmt.Sprintf("can not found default module, appid: %d", appID))
	}

	return (resDataMap["info"].([]interface{}))[0], nil

}

func GetHostByIPAndSource(req *restful.Request, innerIP string, platID int, objURL string) (interface{}, error) {

	param := map[string]interface{}{
		"condition": map[string]interface{}{
			common.BKHostInnerIPField: innerIP,
			common.BKCloudIDField:     platID,
		},
		"fields": common.BKHostIDField,
	}

	resMap, err := getObjByCondition(req, param, common.BKInnerObjIDHost, objURL)

	if nil != err {
		return nil, err
	}

	if !resMap["result"].(bool) {
		return nil, errors.New(resMap["message"].(string))
	}

	resDataMap := resMap["data"].(map[string]interface{})

	blog.Debug("getHostByIPAndSource res: %v", resDataMap)

	return resDataMap["info"], nil
}

func GetHostByCond(req *restful.Request, param map[string]interface{}, objURL string) (interface{}, error) {

	blog.Debug("param:%v", param)
	resMap, err := getObjByCondition(req, param, common.BKInnerObjIDHost, objURL)

	if nil != err {
		return nil, err
	}

	if !resMap["result"].(bool) {
		return nil, errors.New(resMap["message"].(string))
	}

	resDataMap := resMap["data"].(map[string]interface{})

	blog.Debug("getHostByIPArrAndSource res: %v", resDataMap)

	return resDataMap["info"], nil
}

//search host helpers

func GetHostMapByCond(req *restful.Request, condition map[string]interface{}) (map[int]interface{}, []int, error) {
	hostMap := make(map[int]interface{})
	hostIDArr := make([]int, 0)

	// build host controller url
	url := host.CC.HostCtrl() + "/host/v1/hosts/search"
	searchParams := map[string]interface{}{
		"fields":    "",
		"condition": condition,
	}
	inputJson, err := json.Marshal(searchParams)
	if nil != err {
		return nil, nil, err
	}
	hostInfo, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(inputJson))
	blog.Debug("appInfo:%v", hostInfo)
	if nil != err {
		blog.Errorf("getHostMapByCond error:%s, params:%s, error:%s", url, string(inputJson), err.Error())
		return hostMap, hostIDArr, err
	}

	js, err := simplejson.NewJson([]byte(hostInfo))
	if nil != err {
		return nil, nil, err
	}

	resDataInfo, err := js.Get("data").Get("info").Array() //res["data"].(map[string]interface{})
	if nil != err {
		return nil, nil, err
	}

	for _, item := range resDataInfo {
		host := item.(map[string]interface{})
		host_id, err := util.GetIntByInterface(host[common.BKHostIDField])
		if nil != err {
			return nil, nil, err
		}

		hostMap[host_id] = host
		hostIDArr = append(hostIDArr, host_id)
	}
	return hostMap, hostIDArr, nil
}

// GetHostDataByConfig  get host info
func GetHostDataByConfig(req *restful.Request, configData []map[string]int) ([]interface{}, error) {

	hostIDArr := make([]int, 0)

	for _, config := range configData {
		hostIDArr = append(hostIDArr, config[common.BKHostIDField])
	}

	hostMapCondition := map[string]interface{}{
		common.BKHostIDField: map[string]interface{}{
			common.BKDBIN: hostIDArr,
		},
	}

	hostMap, _, err := GetHostMapByCond(req, hostMapCondition)
	if nil != err {
		return nil, err
	}

	hostData, err := SetHostData(req, configData, hostMap)
	if nil != err {
		return hostData, err
	}

	return hostData, nil
}

func GetCustomerPropertyByOwner(req *restful.Request, OwnerId interface{}, ObjCtrl string) ([]map[string]interface{}, error) {
	blog.Debug("getCustomerPropertyByOwner start")
	gHostAttrUrl := ObjCtrl + "/object/v1/meta/objectatts"
	searchBody := make(map[string]interface{})
	searchBody[common.BKObjIDField] = common.BKInnerObjIDHost
	searchBody[common.BKOwnerIDField] = OwnerId
	searchJson, _ := json.Marshal(searchBody)
	gHostAttrRe, err := httpcli.ReqHttp(req, gHostAttrUrl, common.HTTPSelectPost, []byte(searchJson))
	if nil != err {
		blog.Error("GetHostDetailById  attr error :%v", err)
		return nil, err
	}
	js, err := simplejson.NewJson([]byte(gHostAttrRe))
	gHostAttr, _ := js.Map()

	gAttrResult := gHostAttr["result"].(bool)
	if false == gAttrResult {
		blog.Error("GetHostDetailById  attr error :%v", err)
		return nil, err
	}
	hostAttrArr := gHostAttr["data"].([]interface{})
	customAttrArr := make([]map[string]interface{}, 0)
	for _, attr := range hostAttrArr {
		if !attr.(map[string]interface{})[common.BKIsPre].(bool) {
			customAttrArr = append(customAttrArr, attr.(map[string]interface{}))
		}
	}
	return customAttrArr, nil
}

// In_existIpArr exsit ip in array
func In_existIpArr(arr []string, ip string) bool {
	for _, v := range arr {
		if ip == v {
			return true
		}
	}
	return false
}

func getObjByCondition(req *restful.Request, param map[string]interface{}, objType, objURL string) (map[string]interface{}, error) {
	resMap := make(map[string]interface{})

	url := objURL + "/object/v1/insts/" + objType + "/search"
	inputJson, _ := json.Marshal(param)
	res, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(inputJson))
	if nil != err {
		return nil, err
	}

	err = json.Unmarshal([]byte(res), &resMap)
	if nil != err {
		return nil, err
	}

	return resMap, nil
}
