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

package process

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
    "configcenter/src/common/metadata"
)

type ProcessClientInterface interface {
	GetProcessDetailByID(ctx context.Context, ownerID string, appID string, procID string, h http.Header) (resp *metadata.Response, err error)
	GetProcessBindModule(ctx context.Context, ownerID string, businessID string, procID string, h http.Header) (resp *metadata.Response, err error)
	BindModuleProcess(ctx context.Context, ownerID string, businessID string, procID string, moduleName string, h http.Header) (resp *metadata.Response, err error)
	DeleteModuleProcessBind(ctx context.Context, ownerID string, businessID string, procID string, moduleName string, h http.Header) (resp *metadata.Response, err error)
	CreateProcess(ctx context.Context, ownerID string, businessID string, h http.Header, dat interface{}) (resp *metadata.Response, err error)
	DeleteProcess(ctx context.Context, ownerID string, businessID string, procID string, h http.Header) (resp *metadata.Response, err error)
	SearchProcess(ctx context.Context, ownerID string, businessID string, h http.Header) (resp *metadata.Response, err error)
	UpdateProcess(ctx context.Context, ownerID string, businessID string, procID string, h http.Header, dat interface{}) (resp *metadata.Response, err error)
    BatchUpdateProcess(ctx context.Context, ownerID, businessID string, h http.Header, dat interface{}) (resp *metadata.Response, err error)
    OperateProcessInstance(ctx context.Context, namespace string, h http.Header, dat interface{}) (resp *metadata.Response, err error)
	QueryProcessOperateResult(ctx context.Context, namespace string, h http.Header, dat interface{}) (resp *metadata.Response, err error)
	CreateConfigTemp(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error)
	UpdateConfigTemp(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error)
	DeleteConfigTemp(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error)
	QueryConfigTemp(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error)
}

func NewProcessClientInterface(client rest.ClientInterface) ProcessClientInterface {
	return &process{client: client}
}

type process struct {
	client rest.ClientInterface
}
