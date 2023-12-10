/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package key

// KeyType is the key type for common resource cache
type KeyType string

var (
	// ModelType is the key type for model cache
	ModelType KeyType = "model"
	// AttributeType is the key type for model attribute cache
	AttributeType KeyType = "attribute"
	// ModelQuoteRelType is the key type for model quote relation cache
	ModelQuoteRelType KeyType = "model_quote_relation"
)

// KeyKind defines the cache key's different kind of caching aspects
type KeyKind string

var (
	// IDKind is the key kind for data id that stores the detail of the cache key type
	// other kind of keys only stores the id to get detail from IDKind
	IDKind KeyKind = "id"
	// ObjIDKind is the obj id key kind
	ObjIDKind KeyKind = "bk_obj_id"
	// DestModelKind is the destination model id key kind
	DestModelKind KeyKind = "dest_model"
)
