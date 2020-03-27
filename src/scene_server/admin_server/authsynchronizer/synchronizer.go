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

package authsynchronizer

import (
	"context"
	"fmt"

	"configcenter/src/apimachinery"
	"configcenter/src/auth"
	"configcenter/src/auth/authcenter"
	"configcenter/src/auth/extensions"
	enableauth "configcenter/src/common/auth"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/authsynchronizer/handler"
	"configcenter/src/scene_server/admin_server/authsynchronizer/meta"

	"github.com/prometheus/client_golang/prometheus"
)

// AuthSynchronizer stores all related resource
type AuthSynchronizer struct {
	AuthConfig          authcenter.AuthConfig
	clientSet           apimachinery.ClientSetInterface
	ctx                 context.Context
	WorkerCount         int
	Workers             *[]Worker
	WorkerQueue         chan meta.WorkRequest
	Producer            *Producer
	Engine              *backbone.Engine
	SyncIntervalMinutes int

	reg prometheus.Registerer
}

// NewSynchronizer new a synchronizer object
func NewSynchronizer(ctx context.Context, authConfig *authcenter.AuthConfig,
	clientSet apimachinery.ClientSetInterface, reg prometheus.Registerer,
	engine *backbone.Engine) *AuthSynchronizer {
	return &AuthSynchronizer{
		ctx:                 ctx,
		AuthConfig:          *authConfig,
		clientSet:           clientSet,
		reg:                 reg,
		Engine:              engine,
		WorkerCount:         authConfig.SyncWorkerCount,
		SyncIntervalMinutes: authConfig.SyncIntervalMinutes,
	}
}

// Run do start synchronize
func (d *AuthSynchronizer) Run() error {
	if !enableauth.IsAuthed() {
		blog.Info("authConfig is disabled, exit now")
		return nil
	}

	blog.Infof("auth synchronize start..., worker count: %d", d.WorkerCount)

	// init queue
	d.WorkerQueue = make(chan meta.WorkRequest, 1000)

	// make fake handler
	blog.Infof("new auth client with config: %+v", d.AuthConfig)
	authorize, err := auth.NewAuthorize(nil, d.AuthConfig, d.reg)
	if err != nil {
		blog.Errorf("new auth client failed, err: %+v", err)
		return fmt.Errorf("new auth client failed, err: %+v", err)
	}
	authManager := extensions.NewAuthManager(d.clientSet, authorize)
	workerHandler := handler.NewIAMHandler(d.clientSet, authManager)

	// init worker
	workers := make([]Worker, d.WorkerCount)
	for w := 1; w <= d.WorkerCount; w++ {
		worker := NewWorker(w, d.WorkerQueue, workerHandler)
		workers = append(workers, *worker)
		worker.Start()
	}
	d.Workers = &workers

	// init producer
	d.Producer = NewProducer(d.clientSet, authManager, d.WorkerQueue, d.Engine, d.SyncIntervalMinutes)
	d.Producer.Start()
	blog.Infof("auth synchronize started")
	return nil
}
