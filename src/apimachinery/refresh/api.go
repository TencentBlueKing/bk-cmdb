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

package refresh

import (
	"context"
	"fmt"

	"configcenter/pkg/tenant/types"
	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/http/header/util"
	commontypes "configcenter/src/common/types"
)

// RefreshTenant refresh tenant info
func (r *refresh) RefreshTenant(moduleName string) ([]types.Tenant, error) {

	switch moduleName {

	case commontypes.CC_MODULE_APISERVER:
		r.capability.Discover = r.disc.ApiServer()
	case commontypes.CC_MODULE_TASK:
		r.capability.Discover = r.disc.TaskServer()
	case commontypes.CC_MODULE_CACHESERVICE:
		r.capability.Discover = r.disc.CacheService()
	default:
		return nil, fmt.Errorf("unsupported refresh module: %s", moduleName)
	}

	resp := new(types.AllTenantsResult)
	client := rest.NewRESTClient(r.capability, "/")
	err := client.Post().
		WithContext(context.Background()).
		SubResourcef("/refresh/tenants").
		Body(nil).
		WithHeaders(util.GenDefaultHeader()).
		Do().
		Into(resp)

	if err != nil {
		return nil, err
	}

	if ccErr := resp.CCError(); ccErr != nil {
		return nil, ccErr
	}

	return resp.Data, nil
}
