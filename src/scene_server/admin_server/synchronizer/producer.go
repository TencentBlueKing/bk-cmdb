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

package synchronizer

import (
	"context"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/synchronizer/meta"
	"configcenter/src/scene_server/admin_server/synchronizer/utils"
)

// Producer producer WorkRequest and enqueue it
type Producer struct {
	*backbone.Engine
	ID          int
	WorkerQueue chan meta.WorkRequest
	QuitChan    chan bool
}

// NewProducer make a producer
func NewProducer(engine *backbone.Engine, workerQueue chan meta.WorkRequest) *Producer {
	// Create, and return the producer.
	producer := Producer{
		Engine:      engine,
		ID:          0,
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool),
	}

	return &producer
}

// Start do main loop
func (p *Producer) Start() {
	start := time.Now()
	finished := false
	go func(producer *Producer) {
		for {
			if start.Add(time.Minute * 5).Before(time.Now()) {
				start = start.Add(time.Minute * 5)
				finished = false
			}

			if finished == false {
				// get jobs
				jobs := producer.generateJobs()

				for _, job := range *jobs {
					// pass
					p.WorkerQueue <- job
				}
				finished = true
			}
			time.Sleep(time.Millisecond * 100)
		}
	}(p)
}

func (p *Producer) generateJobs() *[]meta.WorkRequest {
	// split all jobs
	jobs := make([]meta.WorkRequest, 0)

	businessSyncJob := meta.WorkRequest{
		ResourceType: meta.BusinessResource,
		Data:         map[string]interface{}{},
	}
	jobs = append(jobs, businessSyncJob)

	// list all business
	header := utils.NewListBusinessAPIHeader()
	condition := metadata.QueryCondition{}
	result, err := p.CoreAPI.CoreService().Instance().ReadInstance(context.TODO(), *header, common.BKInnerObjIDApp, &condition)
	if err != nil {
		blog.Errorf("list business failed, err: %v", err)
		return &jobs
	}
	// businessIDArr := make([]int64, 0)
	businessList := make([]meta.BusinessSimplify, 0)
	for _, business := range result.Data.Info {
		businessSimplify := meta.BusinessSimplify{
			BKAppIDField:      int64(business[common.BKAppIDField].(float64)),
			BKSupplierIDField: int64(business[common.BKSupplierIDField].(float64)),
			BKOwnerIDField:    business[common.BKOwnerIDField].(string),
		}
		// businessID := business[common.BKAppIDField].(int64)
		// businessIDArr = append(businessIDArr, businessID)
		businessList = append(businessList, businessSimplify)
	}
	blog.Info("list business businessList: %+v", businessList)

	// job of synchronize business host to iam
	for _, businessSimplify := range businessList {
		jobs = append(jobs, meta.WorkRequest{
			ResourceType: meta.HostResource,
			Data:         businessSimplify,
		})
	}

	return &jobs
}
