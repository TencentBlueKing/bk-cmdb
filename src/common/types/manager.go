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

package types

var (
	// needDiscoveryServiceName 服务依赖的第三方服务名字的配置
	needDiscoveryServiceName map[string]struct{} = make(map[string]struct{}, 0)
)

// DiscoveryAllService 发现所有定义的服务
func DiscoveryAllService() {
	for name := range AllModule {
		needDiscoveryServiceName[name] = struct{}{}
	}
}

// AddDiscoveryService 新加需要发现服务的名字
func AddDiscoveryService(name ...string) {
	for _, name := range name {
		needDiscoveryServiceName[name] = struct{}{}
	}
}

func GetDiscoveryService() map[string]struct{} {
	// compatible 如果没配置,发现所有的服务
	if len(needDiscoveryServiceName) == 0 {
		DiscoveryAllService()
	}
	return needDiscoveryServiceName
}
