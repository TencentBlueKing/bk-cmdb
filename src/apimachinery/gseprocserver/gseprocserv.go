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
package gseprocserver

import (
    "context"
    "net/http"
    "fmt"
    
    "configcenter/src/apimachinery/rest"
    "configcenter/src/common/metadata"
    "configcenter/src/apimachinery/util"
)

type GseProcClientInterface interface {
    OperateProcess(ctx context.Context, h http.Header, namespace string, data interface{}) (resp *metadata.GseProcRespone, err error)
    QueryProcOperateResult(ctx context.Context, h http.Header, namespace, taskid string) (resp *metadata.GseProcRespone, err error)
    QueryProcStatus(ctx context.Context, h http.Header, namespace string, data interface{}) (resp *metadata.GseProcRespone, err error)
    RegisterProcInfo(ctx context.Context, h http.Header, namespace string, data interface{}) (resp *metadata.GseProcRespone, err error)
    UnRegisterProcInfo(ctx context.Context, h http.Header, namespace string, data interface{}) (resp *metadata.GseProcRespone, err error)
}

func NewGseProcClientInterface(c *util.Capability, version string) GseProcClientInterface {
    base := fmt.Sprintf("/procapi/%s", version)
    return &gseproc{
        client: rest.NewRESTClient(c, base),
    }
}

type gseproc struct {
    client rest.ClientInterface
}
