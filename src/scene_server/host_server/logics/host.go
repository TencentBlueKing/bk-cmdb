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

package logics

import (
	"configcenter/src/common"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (lgc *Logics) GetConfigByCond(pheader http.Header, cond map[string][]int64) ([]map[string]int64, error) {
	configArr := make([]map[string]int64, 0)

	if 0 == len(cond) {
		return configArr, nil
	}

	result, err := lgc.CoreAPI.HostController().Module().GetModulesHostConfig(context.Background(), pheader, cond)
	if err != nil || (err == nil && !result.Result) {
		return nil, fmt.Errorf("get module host config failed, err: %v, %v", err, result.ErrMsg)
	}

	for _, infos := range result.Data {
		info := infos.(map[string]interface{})
		hostID, err := info[common.BKHostIDField].(json.Number).Int64()
		if err != nil {
			return nil, err
		}

		setID, err := info[common.BKSetIDField].(json.Number).Int64()
		if err != nil {
			return nil, err
		}

		moduleID, err := info[common.BKModuleIDField].(json.Number).Int64()
		if err != nil {
			return nil, err
		}

		appID, err := info[common.BKAppIDField].(json.Number).Int64()
		if err != nil {
			return nil, err
		}

		data := make(map[string]int)
		data[common.BKAppIDField] = int(appID)
		data[common.BKSetIDField] = int(setID)
		data[common.BKModuleIDField] = int(moduleID)
		data[common.BKHostIDField] = int(hostID)
		configArr = append(configArr, data)
	}
	return configArr, nil
}
