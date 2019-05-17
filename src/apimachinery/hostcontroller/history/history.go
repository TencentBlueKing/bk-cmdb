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

package history

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/metadata"
)

type HistoryInterface interface {
	AddHistory(ctx context.Context, user string, h http.Header, dat *metadata.HistoryContent) (resp *metadata.AddHistoryResult, err error)
	GetHistorys(ctx context.Context, user string, start string, limit string, h http.Header) (resp *metadata.GetHistoryResult, err error)
}

func NewHistoryInterface(client rest.ClientInterface) HistoryInterface {
	return &history{client: client}
}

type history struct {
	client rest.ClientInterface
}
