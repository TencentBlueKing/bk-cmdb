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

// TestModuleAction_GetModulesByApp 测试获取AppID为1下的所有模块，数量为2
func TestModuleAction_GetModulesByApp(t *testing.T) {
	rsp, err := http.Post(siteURL+"Module/getmodules", "application/x-www-form-urlencoded", strings.NewReader("ApplicationID=1"))
	if err != nil {
		t.Errorf("TestModuleAction_GetModulesByApp fail, %v", err)
		return
	}
	defer rsp.Body.Close()

	resMap, err := utils.GetResMap(rsp)
	if err != nil {
		t.Errorf("TestModuleAction_GetModulesByApp fail, %v", err)
		return
	}

	resCode, ok := resMap["code"].(float64)
	if !ok {
		t.Errorf("TestModuleAction_GetModulesByApp fail, err: resCode is not float64")
		return
	}

	if resCode != 0 {
		t.Errorf("TestModuleAction_GetModulesByApp fail, msg: %s", resMap["msg"])
		return
	}

	if len(resMap["data"].([]interface{})) != 2 {
		t.Errorf("TestModuleAction_GetModulesByApp fail, msg: data length must eq 2")
		return
	}
}

// TestModuleAction_GetModulesByApp2 测试获取AppID为999下的所有模块，数量为0
func TestModuleAction_GetModulesByApp2(t *testing.T) {
	rsp, err := http.Post(siteURL+"Module/getmodules", "application/x-www-form-urlencoded", strings.NewReader("ApplicationID=999"))
	if err != nil {
		t.Errorf("TestModuleAction_GetModulesByApp2 fail, %v", err)
		return
	}
	defer rsp.Body.Close()

	resMap, err := utils.GetResMap(rsp)
	if err != nil {
		t.Errorf("TestModuleAction_GetModulesByApp2 fail, %v", err)
		return
	}

	resCode, ok := resMap["code"].(float64)
	if !ok {
		t.Errorf("TestModuleAction_GetModulesByApp2 fail, err: resCode is not float64")
		return
	}

	if resCode != 0 {
		t.Errorf("TestModuleAction_GetModulesByApp2 fail, msg: %s", resMap["msg"])
		return
	}

	if len(resMap["data"].([]interface{})) != 0 {
		t.Errorf("TestModuleAction_GetModulesByApp2 fail, msg: data length must eq 0")
		return
	}
}

// TestModuleAction_GetModulesByApp3 测试获取AppID为空的情况，code大于0
func TestModuleAction_GetModulesByApp3(t *testing.T) {
	rsp, err := http.Post(siteURL+"Module/getmodules", "application/x-www-form-urlencoded", strings.NewReader("ApplicationID="))
	if err != nil {
		t.Errorf("TestModuleAction_GetModulesByApp3 fail, %v", err)
		return
	}
	defer rsp.Body.Close()

	resMap, err := utils.GetResMap(rsp)
	if err != nil {
		t.Errorf("TestModuleAction_GetModulesByApp3 fail, %v", err)
		return
	}

	resCode, ok := resMap["code"].(float64)
	if !ok {
		t.Errorf("TestModuleAction_GetModulesByApp3 fail, err: resCode is not float64")
		return
	}

	if resCode == 0 {
		t.Errorf("TestModuleAction_GetModulesByApp3 fail, msg: resCode must gt 0")
		return
	}

}

// TestModuleAction_GetModulesByProperty 测试根据SetID获取模块，SetID, AppID为1，结果应大于0
func TestModuleAction_GetModulesByProperty(t *testing.T) {
	rsp, err := http.Post(siteURL+"Set/getmodulesbyproperty", "application/x-www-form-urlencoded", strings.NewReader("ApplicationID=1&SetID=1"))
	if err != nil {
		t.Errorf("TestModuleAction_GetModulesByProperty fail, %v", err)
		return
	}
	defer rsp.Body.Close()

	resMap, err := utils.GetResMap(rsp)
	if err != nil {
		t.Errorf("TestModuleAction_GetModulesByProperty fail, %v", err)
		return
	}

	resCode, ok := resMap["code"].(float64)
	if !ok {
		t.Errorf("TestModuleAction_GetModulesByProperty fail, err: resCode is not float64")
		return
	}

	if resCode != 0 {
		t.Errorf("TestModuleAction_GetModulesByProperty fail, msg: resCode must eq 0")
		return
	}

	if len(resMap["data"].([]interface{})) == 0 {
		t.Errorf("TestModuleAction_GetModulesByProperty fail, msg: resCode must gt 0")
		return
	}
}

// TestModuleAction_GetModulesByProperty 测试根据SetID获取模块，SetID为999，结果应等于0
func TestModuleAction_GetModulesByProperty2(t *testing.T) {
	rsp, err := http.Post(siteURL+"Set/getmodulesbyproperty", "application/x-www-form-urlencoded", strings.NewReader("ApplicationID=1&SetID=999"))
	if err != nil {
		t.Errorf("TestModuleAction_GetModulesByProperty2 fail, %v", err)
		return
	}
	defer rsp.Body.Close()

	resMap, err := utils.GetResMap(rsp)
	if err != nil {
		t.Errorf("TestModuleAction_GetModulesByProperty2 fail, %v", err)
		return
	}

	resCode, ok := resMap["code"].(float64)
	if !ok {
		t.Errorf("TestModuleAction_GetModulesByProperty2 fail, err: resCode is not float64")
		return
	}

	if resCode != 0 {
		t.Errorf("TestModuleAction_GetModulesByProperty2 fail, msg: resCode must eq 0")
		return
	}

	if len(resMap["data"].([]interface{})) != 0 {
		t.Errorf("TestModuleAction_GetModulesByProperty2 fail, msg: resCode must eq 0")
		return
	}
}

// TestModuleAction_UpdateModule 测试更新模块名称，code应等于0
//func TestModuleAction_UpdateModule(t *testing.T) {
//	rsp, err := http.Post(siteURL + "module/editmodule", "application/x-www-form-urlencoded", strings.NewReader("ApplicationID=1&ModuleID=1&ModuleName=空闲机"))
//	if err != nil {
//		t.Errorf("TestModuleAction_UpdateModule fail, %v", err)
//		return
//	}
//	defer rsp.Body.Close()
//
//	resMap, err := utils.GetResMap(rsp)
//	if err != nil {
//		t.Errorf("TestModuleAction_UpdateModule fail, %v", err)
//		return
//	}
//
//	resCode, ok := resMap["code"].(float64)
//	if !ok {
//		t.Errorf("TestModuleAction_UpdateModule fail, err: resCode is not float64")
//		return
//	}
//
//	if resCode != 0 {
//		t.Errorf("TestModuleAction_UpdateModule fail, msg: resCode must eq 0")
//		return
//	}
//
//}
