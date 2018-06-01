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

package host

import (
	"configcenter/src/api_server/ccapi/logics/v2/common/utils"
	"net/http"
	"strings"
	"testing"
)

var siteURL string = "http://127.0.0.1:50050/api/"

// TestHostAction_GetHostListByIP 测试获取存在的IP
func TestHostAction_GetHostListByIP(t *testing.T) {
	rsp, err := http.Post(siteURL+"Host/gethostlistbyip", "application/x-www-form-urlencoded", strings.NewReader("ApplicationID=1&IP=127.0.0.1"))
	if err != nil {
		t.Errorf("TestHostAction_GetHostListByIP fail, %v", err)
		return
	}
	defer rsp.Body.Close()

	resMap, err := utils.GetResMap(rsp)
	if err != nil {
		t.Errorf("TestHostAction_GetHostListByIP fail, %v", err)
		return
	}

	resCode, ok := resMap["code"].(float64)
	if !ok {
		t.Errorf("TestHostAction_GetHostListByIP fail, err: resCode is not float64")
		return
	}

	if resCode != 0 {
		t.Errorf("TestHostAction_GetHostListByIP fail, msg: %s", resMap["msg"])
		return
	}

	if len(resMap["data"].([]interface{})) == 0 {
		t.Errorf("TestHostAction_GetHostListByIP fail, msg: data length must gt 0")
	}
}

// TestHostAction_GetHostListByIP2 测试获取不存在的IP
func TestHostAction_GetHostListByIP2(t *testing.T) {
	rsp, err := http.Post(siteURL+"Host/gethostlistbyip", "application/x-www-form-urlencoded", strings.NewReader("ApplicationID=1&IP=127.0.0.1"))
	if err != nil {
		t.Errorf("TestHostAction_GetHostListByIP2 fail, %v", err)
		return
	}
	defer rsp.Body.Close()

	resMap, err := utils.GetResMap(rsp)
	if err != nil {
		t.Errorf("TestHostAction_GetHostListByIP2 fail, %v", err)
		return
	}

	resCode, ok := resMap["code"].(float64)
	if !ok {
		t.Errorf("TestHostAction_GetHostListByIP2 fail, err: resCode is not float64")
		return
	}

	if resCode != 0 {
		t.Errorf("TestHostAction_GetHostListByIP2 fail, msg: %s", resMap["msg"])
	}

	if len(resMap["data"].([]interface{})) != 0 {
		t.Errorf("TestHostAction_GetHostListByIP2 fail, msg: data length must eq 0")
	}
}

// TestHostAction_GetModuleHostList 测试获取存在主机的模块，业务ID：1，模块ID：1
func TestHostAction_GetModuleHostList(t *testing.T) {
	rsp, err := http.Post(siteURL+"Host/getmodulehostlist", "application/x-www-form-urlencoded", strings.NewReader("ApplicationID=1&ModuleID=1"))
	if err != nil {
		t.Errorf("TestHostAction_GetModuleHostList fail, %v", err)
		return
	}
	defer rsp.Body.Close()

	resMap, err := utils.GetResMap(rsp)
	if err != nil {
		t.Errorf("TestHostAction_GetModuleHostList fail, %v", err)
		return
	}

	resCode, ok := resMap["code"].(float64)
	if !ok {
		t.Errorf("TestHostAction_GetModuleHostList fail, err: resCode is not float64")
		return
	}

	if resCode != 0 {
		t.Errorf("TestHostAction_GetModuleHostList fail, msg: %s", resMap["msg"])
	}

	if len(resMap["data"].([]interface{})) == 0 {
		t.Errorf("TestHostAction_GetModuleHostList fail, msg: data length must gt 0")
	}
}

// TestHostAction_GetModuleHostList2 测试获取不存在主机的模块，业务ID：1，模块ID：3
func TestHostAction_GetModuleHostList2(t *testing.T) {
	rsp, err := http.Post(siteURL+"Host/getmodulehostlist", "application/x-www-form-urlencoded", strings.NewReader("ApplicationID=1&ModuleID=3"))
	if err != nil {
		t.Errorf("TestHostAction_GetModuleHostList2 fail, %v", err)
		return
	}
	defer rsp.Body.Close()

	resMap, err := utils.GetResMap(rsp)
	if err != nil {
		t.Errorf("TestHostAction_GetModuleHostList2 fail, %v", err)
		return
	}

	resCode, ok := resMap["code"].(float64)
	if !ok {
		t.Errorf("TestHostAction_GetModuleHostList2 fail, err: resCode is not float64")
		return
	}

	if resCode != 0 {
		t.Errorf("TestHostAction_GetModuleHostList2 fail, msg: %s", resMap["msg"])
	}

	if len(resMap["data"].([]interface{})) != 0 {
		t.Errorf("TestHostAction_GetModuleHostList2 fail, msg: data length must eq 0")
	}
}

// TestHostAction_GetSetHostList 测试获取存在主机的Set，业务ID：1，SetID：1
func TestHostAction_GetSetHostList(t *testing.T) {
	rsp, err := http.Post(siteURL+"Host/getsethostlist", "application/x-www-form-urlencoded", strings.NewReader("ApplicationID=1&SetID=1"))
	if err != nil {
		t.Errorf("TestHostAction_GetSetHostList fail, %v", err)
		return
	}
	defer rsp.Body.Close()

	resMap, err := utils.GetResMap(rsp)
	if err != nil {
		t.Errorf("TestHostAction_GetSetHostList fail, %v", err)
		return
	}

	resCode, ok := resMap["code"].(float64)
	if !ok {
		t.Errorf("TestHostAction_GetSetHostList fail, err: resCode is not float64")
		return
	}

	if resCode != 0 {
		t.Errorf("TestHostAction_GetSetHostList fail, msg: %s", resMap["msg"])
	}

	if len(resMap["data"].([]interface{})) == 0 {
		t.Errorf("TestHostAction_GetSetHostList fail, msg: data length must gt 0")
	}
}

// TestHostAction_GetSetHostList2 测试获取不存在主机的Set，业务ID：1，SetID：3
func TestHostAction_GetSetHostList2(t *testing.T) {
	rsp, err := http.Post(siteURL+"Host/getsethostlist", "application/x-www-form-urlencoded", strings.NewReader("ApplicationID=1&SetID=3"))
	if err != nil {
		t.Errorf("TestHostAction_GetSetHostList2 fail, %v", err)
		return
	}
	defer rsp.Body.Close()

	resMap, err := utils.GetResMap(rsp)
	if err != nil {
		t.Errorf("TestHostAction_GetSetHostList2 fail, %v", err)
		return
	}

	resCode, ok := resMap["code"].(float64)
	if !ok {
		t.Errorf("TestHostAction_GetSetHostList2 fail, err: resCode is not float64")
		return
	}

	if resCode != 0 {
		t.Errorf("TestHostAction_GetSetHostList2 fail, msg: %s", resMap["msg"])
	}

	if len(resMap["data"].([]interface{})) != 0 {
		t.Errorf("TestHostAction_GetSetHostList2 fail, msg: data length must eq 0")
	}
}

// TestHostAction_GetAppHostList 测试获取存在主机的业务，业务ID：1
func TestHostAction_GetAppHostList(t *testing.T) {
	rsp, err := http.Post(siteURL+"Host/getapphostlist", "application/x-www-form-urlencoded", strings.NewReader("ApplicationID=1"))
	if err != nil {
		t.Errorf("TestHostAction_GetAppHostList fail, %v", err)
		return
	}
	defer rsp.Body.Close()

	resMap, err := utils.GetResMap(rsp)
	if err != nil {
		t.Errorf("TestHostAction_GetAppHostList fail, %v", err)
		return
	}

	resCode, ok := resMap["code"].(float64)
	if !ok {
		t.Errorf("TestHostAction_GetAppHostList fail, err: resCode is not float64")
		return
	}

	if resCode != 0 {
		t.Errorf("TestHostAction_GetAppHostList fail, msg: %s", resMap["msg"])
	}

	if len(resMap["data"].([]interface{})) == 0 {
		t.Errorf("TestHostAction_GetAppHostList fail, msg: data length must gt 0")
	}
}

// TestHostAction_GetHostsByProperty 测试通过SetID属性获取主机，业务ID：1，SetID：1
func TestHostAction_GetHostsByProperty(t *testing.T) {
	rsp, err := http.Post(siteURL+"Host/gethostsbyproperty", "application/x-www-form-urlencoded", strings.NewReader("ApplicationID=1&SetID=1"))
	if err != nil {
		t.Errorf("TestHostAction_GetHostsByProperty fail, %v", err)
		return
	}
	defer rsp.Body.Close()

	resMap, err := utils.GetResMap(rsp)
	if err != nil {
		t.Errorf("TestHostAction_GetHostsByProperty fail, %v", err)
		return
	}

	resCode, ok := resMap["code"].(float64)
	if !ok {
		t.Errorf("TestHostAction_GetHostsByProperty fail, err: resCode is not float64")
		return
	}

	if resCode != 0 {
		t.Errorf("TestHostAction_GetHostsByProperty fail, msg: %s", resMap["msg"])
	}

	if len(resMap["data"].([]interface{})) == 0 {
		t.Errorf("TestHostAction_GetHostsByProperty fail, msg: data length must gt 0")
	}
}
