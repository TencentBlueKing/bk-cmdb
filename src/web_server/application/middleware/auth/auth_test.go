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

package auth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidSysConfigPrivi(t *testing.T) {
	sysPrivi := "event"
	config := "[\"audit\",\"model\"]"
	result := validSysConfigPrivi(config, sysPrivi)
	require.Equal(t, result, false)
}

func TestValidAppConfigPrivi(t *testing.T) {
	method := "post"
	appID := "2"
	userPriviAppStr := "{\"2\":[\"bk_biz_productor\",\"bk_biz_maintainer\"]}"
	rolePrivilege := "{\"2\":[\"resource\",\"hostupdate\"]}"
	pathStr := "/api/v3/host/search"
	result := validAppConfigPrivi(method, appID, userPriviAppStr, rolePrivilege, pathStr)
	require.Equal(t, result, true)
}

func TestValidModelConfigPrivi(t *testing.T) {
	pathArr := []string{"api", "v1"}
	method := "update"
	modelPrivi := "{\"v1\":[\"create\",\"delete\"]}"
	result := validModelConfigPrivi(modelPrivi, method, pathArr)
	require.Equal(t, result, false)
}
