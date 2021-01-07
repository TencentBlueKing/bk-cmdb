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

package RegisterDiscover

// RegDiscvServer define the register and discover function interface
type RegDiscvServer interface {
	// Ping to ping server
	Ping() error
	// start the register and discover service
	Start() error
	// stop the register and discover service
	Stop() error
	// RegisterAndWatch register server info into registe-discover service platform, and watch the info, if not exist, then register again
	RegisterAndWatch(key string, data []byte) error
	// GetServNodes get server nodes
	GetServNodes(key string) ([]string, error)
	// discover server from the registe-discover service platform
	Discover(key string) (<-chan *DiscoverEvent, error)
}
