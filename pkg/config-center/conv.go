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
	"encoding/json/v2"
	"fmt"

	"github.com/spf13/cast"
)

// convertBasic converts config value to specified basic type value.
func convertBasic[T cast.Basic](data any) (T, error) {
	return cast.ToE[T](data)
}

// convert converts config value to pointer of specified type.
func convert[T any](data any) (*T, error) {
	marshal, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshal failed: %w", err)
	}

	result := new(T)
	if err = json.Unmarshal(marshal, result); err != nil {
		return nil, fmt.Errorf("unmarshal failed: %w", err)
	}

	return result, nil
}

// GetBasic get config value as specified basic type.
func GetBasic[T cast.Basic](cc ConfigCenter, conf ConfigType, key string) (T, bool, error) {
	raw, exists := cc.Get(conf, key)
	if !exists {
		return *new(T), false, nil
	}

	data, err := convertBasic[T](raw)
	if err != nil {
		return *new(T), false, fmt.Errorf("convert %s config %s failed: %v", conf, key, err)
	}
	return data, true, nil
}

// GetPtr get config value as pointer of specified type.
func GetPtr[T any](cc ConfigCenter, conf ConfigType, key string) (*T, bool, error) {
	raw, exists := cc.Get(conf, key)
	if !exists {
		return new(T), false, nil
	}

	data, err := convert[T](raw)
	if err != nil {
		return nil, false, fmt.Errorf("convert %s config %s failed: %v", conf, key, err)
	}
	return data, true, nil
}
