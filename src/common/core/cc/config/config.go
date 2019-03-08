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
 
package config

import (
	"fmt"
	"strconv"
	"strings"
)

// CCAPIConfig define configuration of ccapi server
type CCAPIConfig struct {
	AddrPort    string
	RegDiscover string
	ExConfig    string
	Qps int64
	Burst int64
}

// NewCCAPIConfig create ccapi config object
func NewCCAPIConfig() *CCAPIConfig {
	return &CCAPIConfig{
		AddrPort:    "127.0.0.1:8081",
		RegDiscover: "",
		Qps: 1000,
		Burst: 2000,
	}
}

// GetAddress get the address
func (conf *CCAPIConfig) GetAddress() (string, error) {
	arr := strings.Split(conf.AddrPort, ":")
	if len(arr) < 2 {
		return "", fmt.Errorf("the value of flag[AddrPort: %s] is wrong", conf.AddrPort)
	}

	return arr[0], nil
}

// GetPort get the port
func (conf *CCAPIConfig) GetPort() (uint, error) {
	arr := strings.Split(conf.AddrPort, ":")
	if len(arr) < 2 {
		return 0, fmt.Errorf("the value of flag[AddrPort: %s] is wrong", conf.AddrPort)
	}

	port, err := strconv.ParseUint(arr[1], 10, 0)
	if err != nil {
		return 0, err
	}

	return uint(port), nil
}
