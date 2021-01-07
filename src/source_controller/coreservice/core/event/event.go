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

package event

import (
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/redis"
)

type eventOperation struct {
	dbProxy dal.DB
	cache   redis.Client
}

// New create a new event manager instance
func New(dbProxy dal.DB, cache redis.Client) core.EventOperation {
	eventOps := &eventOperation{
		dbProxy:  dbProxy,
		cache:    cache,
	}
	return eventOps
}