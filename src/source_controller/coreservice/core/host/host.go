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

package host

import (
	redis "gopkg.in/redis.v5"

	"configcenter/src/common/eventclient"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/source_controller/coreservice/core/host/modulehost"
	"configcenter/src/storage/dal"
)

var _ core.HostOperation = (*hostManager)(nil)

type hostManager struct {
	DbProxy    dal.RDB
	Cache      *redis.Client
	EventC     eventclient.Client
	moduleHost *modulehost.ModuleHost
}

// New create a new model manager instance
func New(dbProxy dal.RDB, cache *redis.Client) core.HostOperation {

	coreMgr := &hostManager{
		DbProxy: dbProxy,
		Cache:   cache,
		EventC:  eventclient.NewClientViaRedis(cache, dbProxy),
	}
	coreMgr.moduleHost = modulehost.New(dbProxy, cache, coreMgr.EventC)
	return coreMgr
}
