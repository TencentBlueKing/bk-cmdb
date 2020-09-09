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

package resourcedir

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/metadata"
)

type ResourceDirectoryInterface interface {
	CreateResourceDirectory(ctx context.Context, header http.Header, data map[string]interface{}) (resp *metadata.CreatedOneOptionResult, err error)
	UpdateResourceDirectory(ctx context.Context, header http.Header, moduleID int64, data map[string]interface{}) (resp *metadata.Response, err error)
	SearchResourceDirectory(ctx context.Context, header http.Header, data map[string]interface{}) (resp *metadata.SearchResp, err error)
	DeleteResourceDirectory(ctx context.Context, header http.Header, moduleID int64) (resp *metadata.Response, err error)
}

func NewResourceDirectoryInterface(client rest.ClientInterface) ResourceDirectoryInterface {
	return &ResourceDirectory{client: client}
}

type ResourceDirectory struct {
	client rest.ClientInterface
}
