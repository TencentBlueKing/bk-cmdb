/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

package util

import (
	"reflect"
	"strings"
	"testing"
)

var (
	addrT  = "127.0.0.1,localhost"
	userT  = "demoUser"
	pwdT   = "demoPassword"
	rPortT = "redisPort"
	rAuthT = "redisAuth"
)

func TestSetGseConfig(t *testing.T) {

	SetGseConfig(addrT, userT, pwdT, rPortT, rAuthT)

	if !reflect.DeepEqual(strings.Split(addrT, ","), addr) {
		t.Error("test gse addr failed.")
	}

	if userT != usr {
		t.Error("test gse user failed.")
	}

	if pwd != pwdT {
		t.Error("test gse password failed.")
	}

	if rPort != rPortT {
		t.Error("test redis port failed.")
	}

	if rAuth != rAuthT {
		t.Error("test redis password failed.")
	}
}

func TestGetSetConfig(t *testing.T) {
	SetGseConfig(addrT, userT, pwdT, rPortT, rAuthT)
	addr, usr, pwd, rPort, rAuth := GetSetConfig()

	if !reflect.DeepEqual(strings.Split(addrT, ","), addr) {
		t.Error("test gse addr failed.")
	}

	if userT != usr {
		t.Error("test gse user failed.")
	}

	if pwd != pwdT {
		t.Error("test gse password failed.")
	}

	if rPort != rPortT {
		t.Error("test redis port failed.")
	}

	if rAuth != rAuthT {
		t.Error("test redis password failed.")
	}
}
