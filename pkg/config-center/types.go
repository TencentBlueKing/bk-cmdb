/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - CMDB) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
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

package cc

import (
	"fmt"
	"strings"
)

const (
	// pgsqlConfType is the pgsql config type.
	pgsqlConfType = "pgsql"
	// redisConfType is the redis config type.
	redisConfType = "redis"
	// commonConfType is the common config type.
	commonConfType = "common"
	// extraConfType is the extra config type.
	extraConfType = "extra"
)

// allConfTypes stores all config types.
var allConfTypes = []string{redisConfType, pgsqlConfType, commonConfType, extraConfType}

const (
	// configPath is the config register path prefix.
	configPath = "/cc/services/config"
)

// getConfigRegisterPath generate config register path by config key.
func getConfigRegisterPath(key string) string {
	return fmt.Sprintf("%s/%s", configPath, key)
}

// getConfigTypeByRegisterPath get config type by config register path.
func getConfigTypeByRegisterPath(path string) string {
	return strings.TrimPrefix(path, configPath+"/")
}
