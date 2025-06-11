/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
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

package logics

import (
	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
)

// IsDefaultAreaStaticHost check if the host is default area host with static addressing type
func IsDefaultAreaStaticHost(item mapstr.MapStr) (bool, error) {
	cloudID, isExist := item[common.BKCloudIDField]
	if !isExist {
		return false, nil
	}

	cloudIDInt, err := util.GetIntByInterface(cloudID)
	if err != nil {
		return false, err
	}
	if cloudIDInt != common.BKDefaultDirSubArea {
		return false, nil
	}

	// if the addressing mode is not specified, it defaults to static addressing
	addressType, _ := item[common.BKAddressingField]
	if addressType == common.BKAddressingDynamic {
		return false, nil
	}

	return true, nil
}
