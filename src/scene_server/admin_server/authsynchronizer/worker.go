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
	"fmt"
	"time"

	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/authsynchronizer/meta"
)

// NewWorker creates, and returns a new Worker object. Its only argument
// is a channel that the worker can add itself to whenever it is done its
// work.
func NewWorker(id int, workerQueue chan meta.WorkRequest, handler meta.SyncHandler) *Worker {
	// Create, and return the worker.
	worker := Worker{
		ID:          id,
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool),
		SyncHandler: handler,
	}

	return &worker
}

// Worker represent a worker
type Worker struct {
	ID          int
	WorkerQueue chan meta.WorkRequest
	QuitChan    chan bool
	SyncHandler meta.SyncHandler
}

// Start "starts" the worker by starting a goroutine, that is
// an infinite "for-select" loop.
func (w *Worker) Start() {
	go func() {
		for {
			select {
			case work := <-w.WorkerQueue:
				// Receive a work request.
				if err := w.doWork(&work); err != nil {
					blog.Errorf("do work failed, err: %v", err)
				}
				if len(w.WorkerQueue) == 0 {
					blog.Infof("finished all auth synchronize jobs")
				}
				// time interval between two job
				time.Sleep(time.Millisecond * meta.JobIntervalMillisecond)

			case <-w.QuitChan:
				// We have been asked to stop.
				blog.Infof("worker%d stopping\n", w.ID)
				return
			}
		}
	}()
}

// Stop tells the worker to stop listening for work requests.
//
// Note that the worker will only stop *after* it has finished its work.
func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

func (w *Worker) doWork(work *meta.WorkRequest) error {
	var err error
	switch work.ResourceType {
	case meta.BusinessResource:
		err = w.SyncHandler.HandleBusinessSync(work)
	case meta.HostBizResource:
		err = w.SyncHandler.HandleHostSync(work)
	case meta.HostResourcePool:
		err = w.SyncHandler.HandleHostResourcePoolSync(work)
	case meta.SetResource:
		err = w.SyncHandler.HandleSetSync(work)
	case meta.ModuleResource:
		err = w.SyncHandler.HandleModuleSync(work)
	case meta.ModelResource:
		err = w.SyncHandler.HandleModelSync(work)
	case meta.InstanceResource:
		err = w.SyncHandler.HandleInstanceSync(work)
	case meta.ProcessResource:
		// err = w.SyncHandler.HandleAuditSync(work)
		// case meta.AuditCategory:
		err = w.SyncHandler.HandleProcessSync(work)
	case meta.DynamicGroupResource:
		err = w.SyncHandler.HandleDynamicGroupSync(work)
	case meta.ClassificationResource:
		err = w.SyncHandler.HandleClassificationSync(work)
	// case meta.UserGroupSyncResource:
	// 	err = w.SyncHandler.HandleUserGroupSync(work)
	case meta.ServiceTemplateResource:
		err = w.SyncHandler.HandleServiceTemplateSync(work)
	case meta.PlatResource:
		err = w.SyncHandler.HandlePlatSync(work)
	case meta.SetTemplateResource:
		err = w.SyncHandler.HandleSetTemplateSync(work)
	default:
		return fmt.Errorf("unsupported work resource type: %s", work.ResourceType)
	}
	return err
}
