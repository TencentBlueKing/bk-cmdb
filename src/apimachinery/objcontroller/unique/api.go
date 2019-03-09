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

package unique

import (
	"context"
	"net/http"

	"configcenter/src/common/metadata"
)

type UniqueInterface interface {
	Search(ctx context.Context, h http.Header, objectID string) (resp *metadata.SearchUniqueResult, err error)
	Create(ctx context.Context, h http.Header, objectID string, request *metadata.CreateUniqueRequest) (resp *metadata.CreateUniqueResult, err error)
	Update(ctx context.Context, h http.Header, objectID string, id uint64, request *metadata.UpdateUniqueRequest) (resp *metadata.UpdateUniqueResult, err error)
	Delete(ctx context.Context, h http.Header, objectID string, id uint64) (resp *metadata.DeleteUniqueResult, err error)
}
