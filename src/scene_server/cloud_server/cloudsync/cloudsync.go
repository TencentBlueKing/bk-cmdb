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

package cloudsync

import (
	"context"

	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/common/zkclient"
	"configcenter/src/scene_server/cloud_server/logics"
	"configcenter/src/storage/dal/mongo/local"
)

const (
	// 同步器数量
	syncorNum int = 10
)

type SyncConf struct {
	ZKClient  *zkclient.ZkClient
	Logics    *logics.Logics
	AddrPort  string
	MongoConf local.MongoConf
}

// 云同步接口
type CloudSyncInterface interface {
	Sync(task chan *metadata.CloudSyncTask) error
}

// 进行云资源同步
func CloudSync(conf *SyncConf) error {
	errors.SetGlobalCCError(conf.Logics.CCErr)
	ctx := context.Background()

	schedulerConf := &SchedulerConf{
		ZKClient:  conf.ZKClient,
		Logics:    conf.Logics,
		AddrPort:  conf.AddrPort,
		MongoConf: conf.MongoConf,
	}
	scheduler, err := NewTaskScheduler(schedulerConf)
	if err != nil {
		blog.Errorf("NewTaskScheduler failed, schedulerConf:%#v, err:%+v", schedulerConf, err)
	}
	err = scheduler.Schedule(ctx)
	if err != nil {
		blog.Errorf("Schedule failed, err:%+v", err)
	}

	//冷启动获取完任务记录列表才继续往下执行
	<-scheduler.ListerDone()
	processor := NewTaskProcessor(scheduler)

	taskChan := make(chan *metadata.CloudSyncTask, 20)
	// 不断给taskChan提供任务数据
	processor.TaskChanLoop(taskChan)
	// 同步云资源
	SyncCloudResource(taskChan, conf)
	return nil
}

// 同步云资源
func SyncCloudResource(taskChan chan *metadata.CloudSyncTask, conf *SyncConf) {
	// 云主机channel
	hostChan := make(chan *metadata.CloudSyncTask, 10)

	// 云主机同步器处理同步任务
	for i := 1; i <= syncorNum; i++ {
		syncor := NewHostSyncor(conf.Logics)
		go func(syncor *HostSyncor) {
			for {
				task := <-hostChan
				syncor.Sync(task)
			}
		}(syncor)
	}

	// 根据任务类型，将任务放入不同的任务channel
	go func() {
		for {
			if task, ok := <-taskChan; ok {
				blog.V(4).Infof("processing taskid:%d, resource type:%s", task.TaskID, task.ResourceType)
				switch task.ResourceType {
				case "host":
					hostChan <- task
				default:
					blog.Errorf("unknown resource type:%s, ignore it!", task.ResourceType)
				}
			}
		}
	}()
}
