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
	"fmt"
	"time"

	"configcenter/pkg/tenant"
	"configcenter/src/apimachinery"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
)

// InitTenant init tenant, refresh tenants info while server is starting
func InitTenant(apiMachineryCli apimachinery.ClientSetInterface) error {
	coreExist := false
	for retry := 0; retry < 10; retry++ {
		if _, err := apiMachineryCli.Healthz().HealthCheck(types.CC_MODULE_CORESERVICE); err != nil {
			blog.Errorf("connect core server failed: %v", err)
			time.Sleep(time.Second * 2)
			continue
		}
		coreExist = true
		break
	}
	if !coreExist {
		blog.Errorf("core server not exist")
		return fmt.Errorf("core server not exist")
	}
	err := tenant.Init(&tenant.Options{ApiMachineryCli: apiMachineryCli})
	if err != nil {
		return err
	}
	return nil
}
