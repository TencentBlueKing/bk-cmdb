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

package app

import (
	"sync"

	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/storage/dal/mongo"
	"configcenter/src/storage/tmserver/app/options"
	"configcenter/src/storage/tmserver/service"
)

// Server tmserver definition
type Server struct {
	configLock  sync.Mutex
	engin       *backbone.Engine
	config      *options.Config
	coreService service.Service
}

func (s *Server) onConfigUpdate(previous, current cc.ProcessConfig) {

	s.configLock.Lock()
	defer s.configLock.Unlock()

	if len(current.ConfigMap) > 0 {
		if nil == s.config {
			s.config = &options.Config{}
		}

		s.config.MongoDB = mongo.ParseConfigFromKV("mongodb", current.ConfigMap)

		s.config.Transaction.Enable = current.ConfigMap["transaction.enable"]
		s.config.Transaction.TransactionLifetimeSecond = current.ConfigMap["transaction.transactionLifetimeSecond"]
	}
}

