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

package proccontroller

import (
	"fmt"
    "context"
    "net/http"
	
	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/util"
    "configcenter/src/common/metadata"
)

type ProcCtrlClientInterface interface {
    CreateProc2Module(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error)
    GetProc2Module(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.ProcModuleResult, err error)
    DeleteProc2Module(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error)
    CreateConfTemp(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error)
    UpdateConfTemp(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error)
    DeleteConfTemp(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error)
    QueryConfTemp(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error)
    CreateProcInstanceModel(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error)
    GetProcInstanceModel(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.ProcInstModelResult, err error)
    DeleteProcInstanceModel(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error)
}

func NewProcCtrlClientInterface(c *util.Capability, version string) ProcCtrlClientInterface {
	base := fmt.Sprintf("/process/%s", version)
	return &procctrl{
		client: rest.NewRESTClient(c, base),
	}
}

type procctrl struct {
	client rest.ClientInterface
}
