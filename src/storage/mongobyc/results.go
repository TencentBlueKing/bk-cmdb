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

package mongobyc

import (
	"configcenter/src/common/mapstr"
)

// DeleteResult is a result of an DeleteOne operat
type DeleteResult struct {
	// The number of documents that were deleted.
	DeletedCount uint64 `json:"deletedCount"`
}

// UpdateResult is a result of an update operation.
type UpdateResult struct {
	MatchedCount  uint64 `json:"matchedCount"`
	UpsertedCount uint64 `json:"upsertedCount"`
	ModifiedCount uint64 `json:"modifiedCount"`
}

// ReplaceOneResult the  replace one function result
type ReplaceOneResult struct {
	UpdateResult `json:",inline"`
}

// IndexResult get collection index result
type IndexResult struct {
	Namespace string        `json:"ns"`
	Name      string        `json:"name"`
	Key       mapstr.MapStr `json:"key"`
}

// GetIndexResult get the indexex result
type GetIndexResult struct {
	Indexes []IndexResult
}
