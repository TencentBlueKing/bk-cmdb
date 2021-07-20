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

package configcenter

import (
	"configcenter/src/common"
	crd "configcenter/src/common/confregdiscover"
)

type ConfigCenter struct {
	Type       string // type of configuration center
	ConfigCenterDetail crd.ConfRegDiscvIf
}

var (
	configCenterGroup []*ConfigCenter
	configCenterType = common.BKDefaultConfigCenter // the default configuration center is zookeeper.
)

// SetConfigCenterType use this function to change the type of configuration center.
func SetConfigCenterType(serverType string){
	configCenterType = serverType
}

// AddConfigCenter add the configuration center you want to replace.
func AddConfigCenter(configCenter *ConfigCenter) {
	configCenterGroup = append(configCenterGroup, configCenter)
}

// CurrentConfigCenter use this method to return to the configuration center you want to use.
func CurrentConfigCenter() crd.ConfRegDiscvIf {
	var defaultConfigCenter *ConfigCenter
	for _, center := range configCenterGroup {
		if center.Type == configCenterType {
			return center.ConfigCenterDetail
		}
		if common.BKDefaultConfigCenter == center.Type {
			defaultConfigCenter = center
		}
	}
	if nil != defaultConfigCenter {
		return defaultConfigCenter.ConfigCenterDetail
	}
	return nil
}
