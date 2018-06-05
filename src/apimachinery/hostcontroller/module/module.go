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

package module

import (
    "context"
    
    "configcenter/src/apimachinery/rest"
    "configcenter/src/apimachinery/util"
    "configcenter/src/source_controller/hostcontroller/hostdata/actions/instdata"
    "configcenter/src/common/core/cc/api"
)

type ModuleInterface interface {
    GetHostModulesIDs(ctx context.Context, h util.Headers, dat *instdata.ModuleHostConfigParams) (resp *api.BKAPIRsp, err error)
    AddModuleHostConfig(ctx context.Context, h util.Headers, dat *instdata.ModuleHostConfigParams) (resp *api.BKAPIRsp, err error)
    DelModuleHostConfig(ctx context.Context, h util.Headers, dat *instdata.ModuleHostConfigParams) (resp *api.BKAPIRsp, err error)
    DelDefaultModuleHostConfig(ctx context.Context, h util.Headers, dat *instdata.ModuleHostConfigParams) (resp *api.BKAPIRsp, err error)
    MoveHost2ResourcePool(ctx context.Context, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error)
    AssignHostToApp(ctx context.Context, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error)
    GetModulesHostConfig(ctx context.Context, h util.Headers, dat map[string][]int) (resp *api.BKAPIRsp, err error)
}

func NewModuleInterface(client rest.ClientInterface) ModuleInterface {
    return &mod{client:client}
}

type mod struct {
    client rest.ClientInterface
}