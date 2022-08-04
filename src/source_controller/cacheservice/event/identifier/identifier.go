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

// Package identifier TODO
package identifier

import (
	"context"
	"fmt"

	"configcenter/src/apimachinery/discovery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/source_controller/cacheservice/event"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/stream"
)

// NewIdentity TODO
func NewIdentity(
	watch stream.LoopInterface,
	isMaster discovery.ServiceManageInterface,
	watchDB dal.DB,
	ccDB dal.DB) error {

	watchMongoDB, ok := watchDB.(*local.Mongo)
	if !ok {
		blog.Errorf("watch event, but watch db is not an instance of local mongo to start transaction")
		return fmt.Errorf("watch db is not an instance of local mongo")
	}

	base := identityOptions{
		watch:    watch,
		isMaster: isMaster,
		watchDB:  watchMongoDB,
		ccDB:     ccDB,
	}

	host := base
	host.key = event.HostKey
	host.watchFields = needCareHostFields
	if err := newIdentity(context.Background(), host); err != nil {
		blog.Errorf("new host identify host event failed, err: %v", err)
		return err
	}
	blog.Info("host identity events, watch host success.")

	relation := base
	relation.key = event.ModuleHostRelationKey
	relation.watchFields = []string{common.BKHostIDField}
	if err := newIdentity(context.Background(), relation); err != nil {
		blog.Errorf("new host identify host relation event failed, err: %v", err)
		return err
	}
	blog.Info("host identity events, watch host relation success.")

	process := base
	process.key = event.ProcessKey
	process.watchFields = []string{common.BKProcessIDField}
	if err := newIdentity(context.Background(), process); err != nil {
		blog.Errorf("new host identify process event failed, err: %v", err)
		return err
	}
	blog.Info("host identity events, watch process success.")

	return nil
}
