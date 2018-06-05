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

package host

import (
    "context"

    "configcenter/src/apimachinery/rest"
    "configcenter/src/apimachinery/util"
    "configcenter/src/common/core/cc/api"
    "configcenter/src/source_controller/common/commondata"
)

type HostInterface interface {
    GetHostByID(ctx context.Context, hostID string, h util.Headers) (resp *api.BKAPIRsp, err error)
    GetHosts(ctx context.Context, h util.Headers, opt *commondata.ObjQueryInput) (resp *api.BKAPIRsp, err error)
    AddHost(ctx context.Context, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error)
    GetHostSnap(ctx context.Context, hostID string, h util.Headers) (resp *api.BKAPIRsp, err error)
}

func NewHostInterface(client rest.ClientInterface) HostInterface {
    return &hostctrl{client:client}
}

type hostctrl struct {
    client rest.ClientInterface
}
