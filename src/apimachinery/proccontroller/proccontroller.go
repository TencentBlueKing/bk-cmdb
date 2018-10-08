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
	"context"
	"fmt"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common/metadata"
)

type ProcCtrlClientInterface interface {
	CreateProc2Module(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error)
	GetProc2Module(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.ProcModuleResult, err error)
	DeleteProc2Module(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error)
	CreateProc2Template(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error)
	SearchProc2Template(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.MapArrayResponse, err error)
	DeleteProc2Template(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error)
	CreateProcInstanceModel(ctx context.Context, h http.Header, dat []*metadata.ProcInstanceModel) (resp *metadata.Response, err error)
	GetProcInstanceModel(ctx context.Context, h http.Header, dat *metadata.QueryInput) (resp *metadata.ProcInstModelResult, err error)
	DeleteProcInstanceModel(ctx context.Context, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error)
	RegisterProcInstanceDetail(ctx context.Context, h http.Header, dat *metadata.GseProcRequest) (resp *metadata.Response, err error)
	ModifyProcInstanceDetail(ctx context.Context, h http.Header, dat *metadata.ModifyProcInstanceDetail) (resp *metadata.Response, err error)
	GetProcInstanceDetail(ctx context.Context, h http.Header, dat *metadata.QueryInput) (resp *metadata.ProcInstanceDetailResult, err error)
	DeleteProcInstanceDetail(ctx context.Context, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error)
	AddOperateTaskInfo(ctx context.Context, h http.Header, dat []*metadata.ProcessOperateTask) (resp *metadata.Response, err error)
	UpdateOperateTaskInfo(ctx context.Context, h http.Header, dat *metadata.UpdateParams) (resp *metadata.Response, err error)
	SearchOperateTaskInfo(ctx context.Context, h http.Header, dat *metadata.QueryInput) (resp *metadata.ProcessOperateTaskResult, err error)
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
