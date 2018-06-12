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

package application

import (
	"configcenter/src/api_server/ccapi/logics/v2/common/utils"
	"net/http"
	"strings"
	"testing"
)

var siteURL string = "http://127.0.0.1:50050/api/"

// TestAppAction_GetAppList 测试获取所有业务
func TestAppAction_GetAppList(t *testing.T) {
	rsp, err := http.Post(siteURL+"App/getapplist", "application/x-www-form-urlencoded", nil)
	if err != nil {
		t.Errorf("TestAppAction_GetAppList fail, %v", err)
		return
	}
	defer rsp.Body.Close()

	resMap, err := utils.GetResMap(rsp)
	if err != nil {
		t.Errorf("TestAppAction_GetAppList fail, %v", err)
		return
	}

	resCode, ok := resMap["code"].(float64)
	if !ok {
		t.Errorf("TestAppAction_GetAppList fail, err: resCode is not float64")
		return
	}

	if resCode != 0 {
		t.Errorf("TestAppAction_GetAppList fail, msg: %s", resMap["msg"])
	}
}

// TestAppAction_GetAppByID 测试获取单个APP
func TestAppAction_GetAppByID(t *testing.T) {
	rsp, err := http.Post(siteURL+"App/getAppByID", "application/x-www-form-urlencoded", strings.NewReader("ApplicationID=1"))
	if err != nil {
		t.Errorf("TestAppAction_GetAppByID fail, %v", err)
		return
	}
	defer rsp.Body.Close()

	resMap, err := utils.GetResMap(rsp)
	if err != nil {
		t.Errorf("TestAppAction_GetAppByID fail, %v", err)
		return
	}

	resCode, ok := resMap["code"].(float64)
	if !ok {
		t.Errorf("TestAppAction_GetAppByID fail, err: resCode is not float64")
		return
	}

	if resCode != 0 {
		t.Errorf("TestAppAction_GetAppByID fail, msg: %s", resMap["msg"])
		return
	}

	//单个AppID，返回结果数据必须等于1
	if len(resMap["data"].([]interface{})) != 1 {
		t.Errorf("TestAppAction_GetAppByID fail, msg: data length must eq 1")
		return
	}
}

// TestAppAction_GetAppByID2 测试获取不存在的APP
func TestAppAction_GetAppByID2(t *testing.T) {
	rsp, err := http.Post(siteURL+"App/getAppByID", "application/x-www-form-urlencoded", strings.NewReader("ApplicationID=9999"))
	if err != nil {
		t.Errorf("TestAppAction_GetAppByID2 fail, %v", err)
		return
	}
	defer rsp.Body.Close()

	resMap, err := utils.GetResMap(rsp)
	if err != nil {
		t.Errorf("TestAppAction_GetAppByID2 fail, %v", err)
		return
	}

	resCode, ok := resMap["code"].(float64)
	if !ok {
		t.Errorf("TestAppAction_GetAppByID2 fail, err: resCode is not float64")
		return
	}

	if resCode != 0 {
		t.Errorf("TestAppAction_GetAppByID2 fail, msg: %s", resMap["msg"])
		return
	}

	//不存在的AppID，返回结果数据必须等于0
	if len(resMap["data"].([]interface{})) != 0 {
		t.Errorf("TestAppAction_GetAppByID2 fail, msg: data length must eq 0")
		return
	}
}

// TestAppAction_GetAppByID3 测试获取多个APP,ApplicationID=1,3
func TestAppAction_GetAppByID3(t *testing.T) {
	rsp, err := http.Post(siteURL+"App/getAppByID", "application/x-www-form-urlencoded", strings.NewReader("ApplicationID=1,3"))
	if err != nil {
		t.Errorf("TestAppAction_GetAppByID3 fail, %v", err)
		return
	}
	defer rsp.Body.Close()

	resMap, err := utils.GetResMap(rsp)
	if err != nil {
		t.Errorf("TestAppAction_GetAppByID3 fail, %v", err)
		return
	}

	resCode, ok := resMap["code"].(float64)
	if !ok {
		t.Errorf("TestAppAction_GetAppByID3 fail, err: resCode is not float64")
		return
	}

	if resCode != 0 {
		t.Errorf("TestAppAction_GetAppByID3 fail, msg: %s", resMap["msg"])
		return
	}

	//两个AppID，返回结果数据必须等于2
	dataLength := len(resMap["data"].([]interface{}))
	if dataLength != 2 {
		t.Errorf("TestAppAction_GetAppByID3 fail, msg: data length must eq 2, act: %d", dataLength)
		return
	}
}

// TestAppAction_GetAppByUIN 测试正常获取有权限业务
// 确保帐号test有业务权限
func TestAppAction_GetAppByUIN(t *testing.T) {
	rsp, err := http.Post(siteURL+"App/getappbyuin", "application/x-www-form-urlencoded", strings.NewReader("userName=test"))
	if err != nil {
		t.Errorf("TestAppAction_GetAppByUIN fail, %v", err)
		return
	}
	defer rsp.Body.Close()

	resMap, err := utils.GetResMap(rsp)
	if err != nil {
		t.Errorf("TestAppAction_GetAppByUIN fail, %v", err)
		return
	}

	resCode, ok := resMap["code"].(float64)
	if !ok {
		t.Errorf("TestAppAction_GetAppByUIN fail, err: resCode is not float64")
		return
	}

	if resCode != 0 {
		t.Errorf("TestAppAction_GetAppByUIN fail, msg: %s", resMap["msg"])
		return
	}

	if len(resMap["data"].([]interface{})) == 0 {
		t.Errorf("TestAppAction_GetAppByUIN fail, msg: data length must gt 0")
		return
	}
}

// TestAppAction_GetAppByUIN2 测试没有任何业务权限的帐号
// 确保帐号test2没有业务权限
func TestAppAction_GetAppByUIN2(t *testing.T) {
	rsp, err := http.Post(siteURL+"App/getappbyuin", "application/x-www-form-urlencoded", strings.NewReader("userName=test2"))
	if err != nil {
		t.Errorf("TestAppAction_GetAppByUIN2 fail, %v", err)
		return
	}
	defer rsp.Body.Close()

	resMap, err := utils.GetResMap(rsp)
	if err != nil {
		t.Errorf("TestAppAction_GetAppByUIN2 fail, %v", err)
		return
	}

	resCode, ok := resMap["code"].(float64)
	if !ok {
		t.Errorf("TestAppAction_GetAppByUIN2 fail, err: resCode is not float64")
		return
	}

	if resCode != 0 {
		t.Errorf("TestAppAction_GetAppByUIN2 fail, msg: %s", resMap["msg"])
		return
	}

	if len(resMap["data"].([]interface{})) > 0 {
		t.Errorf("TestAppAction_GetAppByUIN2 fail, msg: data length must eq 0")
		return
	}
}

// TestAppAction_GetAppByUIN3 测试帐号为空的情况
func TestAppAction_GetAppByUIN3(t *testing.T) {
	rsp, err := http.Post(siteURL+"App/getappbyuin", "application/x-www-form-urlencoded", strings.NewReader("userName="))
	if err != nil {
		t.Errorf("TestAppAction_GetAppByUIN3 fail, %v", err)
		return
	}
	defer rsp.Body.Close()

	resMap, err := utils.GetResMap(rsp)
	if err != nil {
		t.Errorf("TestAppAction_GetAppByUIN3 fail, %v", err)
		return
	}

	resCode, ok := resMap["code"].(float64)
	if !ok {
		t.Errorf("TestAppAction_GetAppByUIN3 fail, err: resCode is not float64")
		return
	}

	if resCode == 0 {
		t.Errorf("TestAppAction_GetAppByUIN3 fail, msg: resCode must gt 0")
		return
	}
}

// TestAppAction_GetAppSetModuleTreeByAppID 测试正常获取topo树，AppID为1
func TestAppAction_GetAppSetModuleTreeByAppID(t *testing.T) {
	rsp, err := http.Post(siteURL+"TopSetModule/getappsetmoduletreebyappid", "application/x-www-form-urlencoded", strings.NewReader("ApplicationID=1"))
	if err != nil {
		t.Errorf("TestAppAction_GetAppSetModuleTreeByAppID fail, %v", err)
		return
	}
	defer rsp.Body.Close()

	resMap, err := utils.GetResMap(rsp)
	if err != nil {
		t.Errorf("TestAppAction_GetAppSetModuleTreeByAppID fail, %v", err)
		return
	}

	resCode, ok := resMap["code"].(float64)
	if !ok {
		t.Errorf("TestAppAction_GetAppSetModuleTreeByAppID fail, err: resCode is not float64")
		return
	}

	if resCode != 0 {
		t.Errorf("TestAppAction_GetAppSetModuleTreeByAppID fail, msg: resCode must eq 0")
		return
	}
}

// TestAppAction_GetAppSetModuleTreeByAppID2 测试不存在的业务，AppID为9999
func TestAppAction_GetAppSetModuleTreeByAppID2(t *testing.T) {
	rsp, err := http.Post(siteURL+"TopSetModule/getappsetmoduletreebyappid", "application/x-www-form-urlencoded", strings.NewReader("ApplicationID=9999"))
	if err != nil {
		t.Errorf("TestAppAction_GetAppSetModuleTreeByAppID2 fail, %v", err)
		return
	}
	defer rsp.Body.Close()

	resMap, err := utils.GetResMap(rsp)
	if err != nil {
		t.Errorf("TestAppAction_GetAppSetModuleTreeByAppID2 fail, %v", err)
		return
	}

	resCode, ok := resMap["code"].(float64)
	if !ok {
		t.Errorf("TestAppAction_GetAppSetModuleTreeByAppID fail, err: resCode is not float64")
		return
	}

	if resCode != 0 {
		t.Errorf("TestAppAction_GetAppSetModuleTreeByAppID2 fail, msg: resCode must eq 0")
		return
	}

	if resMap["data"] != nil {
		t.Errorf("TestAppAction_GetAppSetModuleTreeByAppID2 fail, msg: resData must eq nil")
		return
	}

}
