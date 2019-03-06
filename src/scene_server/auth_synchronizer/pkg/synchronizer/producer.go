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
	"configcenter/src/scene_server/auth_synchronizer/pkg/synchronizer/meta"
	"time"
)

// Producer producer WorkRequest and enqueue it
type Producer struct {
	ID          int
	WorkerQueue chan meta.WorkRequest
	QuitChan    chan bool
}

// NewProducer make a producer
func NewProducer(workerQueue chan meta.WorkRequest) *Producer {
	// Create, and return the producer.
	producer := Producer{
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
	go func() {
		for {
			if start.Add(time.Minute * 5).Before(time.Now()) {
				start = start.Add(time.Minute * 5)
				finished = false
			}

			if finished == false {
				// split all jobs
				jobs := make([]meta.WorkRequest, 5)

				for _, job := range jobs {
					// pass
					p.WorkerQueue <- job
				}
				finished = true
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()
}
