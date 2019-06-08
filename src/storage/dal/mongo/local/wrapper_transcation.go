/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package local

import (
	"context"
	"net/http"

	"configcenter/src/storage/dal"
)

// AutoRun Interface for automatic processing of encapsulated transactions
// f func return error, abort commit, other commit transcation. transcation commit can be error.
// f func parameter http.header, the handler must be accepted and processed. Subsequent passthrough to call subfunctions and APIs
func (c *Mongo) AutoRun(ctx context.Context, opt dal.TxnWrapperOption, f func(header http.Header) error) error {
	panic("transcation wrapper not implemented")
}
