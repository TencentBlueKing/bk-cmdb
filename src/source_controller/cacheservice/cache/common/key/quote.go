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

import (
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/stream/types"
)

func init() {
	addKeyGenerator(ModelQuoteRelType, modelQuoteRelKey)
}

var modelQuoteRelKey = KeyGenerator{
	namespace: fmt.Sprintf("%s%s:", common.BKCacheKeyV3Prefix, ModelQuoteRelType),
	watchOpt: types.Options{
		EventStruct: new(mapstr.MapStr),
		Collection:  common.BKTableNameModelQuoteRelation,
	},
	expireSeconds:      30 * 60 * time.Second,
	expireRangeSeconds: [2]int{-600, 600},
	idGen: func(data interface{}) (string, float64, error) {
		return commonOidGenerator(data)
	},
	keyGenMap: map[KeyKind]redisKeyGenerator{
		DestModelKind: func(data interface{}) (string, error) {
			return commonKeyGenerator(data, common.BKDestModelField)
		},
	},
	dataGetterMap: map[KeyKind]dataGetter{
		IDKind: func(db dal.DB, keys ...string) ([]interface{}, error) {
			return commonOidDataGetter(db, common.BKTableNameModelQuoteRelation, keys...)
		},
		DestModelKind: func(db dal.DB, keys ...string) ([]interface{}, error) {
			return commonKeyWithOidGetter(db, common.BKTableNameModelQuoteRelation, common.BKDestModelField, keys...)
		},
	},
}
