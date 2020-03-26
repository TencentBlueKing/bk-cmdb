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
	"fmt"
	"sync"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/zkclient"
	ccom "configcenter/src/scene_server/cloud_server/common"
	"configcenter/src/scene_server/cloud_server/logics"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/reflector"
	stypes "configcenter/src/storage/stream/types"

	"stathat.com/c/consistent"
)

type taskProcessor struct {
	zkClient  *zkclient.ZkClient
	db        dal.DB
	logics    *logics.Logics
	addrport  string
	reflector reflector.Interface
	hashring  *consistent.Consistent
	tasklist  map[int64]*metadata.CloudSyncTask
	taskChan  chan *metadata.CloudSyncTask
	mu        sync.RWMutex
}

const (
	// 同步器数量
	syncorNum int = 10
	// 循环检查任务列表的间隔
	checkInterval int = 5
)

var (
	// mongo server对于满足change stream查询的最大等待时间
	maxAwaitTime = time.Second * 10
	header       = ccom.GetHeader()
)
var kit *rest.Kit

type SyncConf struct {
	ZKClient  *zkclient.ZkClient
	DB        dal.DB
	Logics    *logics.Logics
	AddrPort  string
	MongoConf local.MongoConf
}

// 云同步接口
type CloudSyncInterface interface {
	Sync(task chan *metadata.CloudSyncTask) error
}

// 处理云资源同步任务
func CloudSync(ctx context.Context, conf *SyncConf) error {
	t := &taskProcessor{
		zkClient: conf.ZKClient,
		db:       conf.DB,
		logics:   conf.Logics,
		addrport: conf.AddrPort,
		hashring: consistent.New(),
		tasklist: make(map[int64]*metadata.CloudSyncTask),
		taskChan: make(chan *metadata.CloudSyncTask, 20),
	}

	kit = ccom.GetKit(header)

	var err error
	t.reflector, err = reflector.NewReflector(conf.MongoConf)
	if err != nil {
		blog.Errorf("NewReflector failed, mongoConf: %#v, err: %s", conf.MongoConf, err.Error())
		return err
	}

	// 监听任务表事件
	if err := t.WatchTaskTable(ctx); err != nil {
		return err
	}

	// 监听服务进程节点
	if err := t.WatchTaskNode(); err != nil {
		return err
	}

	// 不断给任务channel提供任务数据
	t.TaskChanLoop()
	// 同步云资源
	t.SyncCloudResource()
	return nil
}

// 监听zk的cloudserver节点变化，有变化时重新分配当前进程的任务列表
func (t *taskProcessor) WatchTaskNode() error {
	go func() {
		for servers := range t.logics.Discovery().CloudServer().GetServersChan() {
			t.setHashring(servers)
			t.dispatchTasks()
		}
	}()
	return nil
}

// 监控云资源同步任务表事件，发现有变更则判断是否将该任务加入到当前进程的任务列表里
func (t *taskProcessor) WatchTaskTable(ctx context.Context) error {
	opts := &stypes.WatchOptions{
		Options: stypes.Options{
			MaxAwaitTime: &maxAwaitTime,
			EventStruct:  new(metadata.CloudSyncTask),
			Collection:   common.BKTableNameCloudSyncTask,
		},
	}
	cap := &reflector.Capable{
		reflector.OnChangeEvent{
			OnAdd:    t.changeOnAdd,
			OnUpdate: t.changeOnUpdate,
			OnDelete: t.changeOnDelete,
		},
	}

	return t.reflector.Watcher(ctx, opts, cap)
}

// 表记录新增处理逻辑
func (t *taskProcessor) changeOnAdd(event *stypes.Event) {
	blog.V(4).Infof("OnAdd event, taskid:%d", event.Document.(*metadata.CloudSyncTask).TaskID)
	t.addTask(event.Document.(*metadata.CloudSyncTask))
}

// 表记录更新处理逻辑
func (t *taskProcessor) changeOnUpdate(event *stypes.Event) {
	blog.V(4).Infof("OnUpdate event, taskid:%d", event.Document.(*metadata.CloudSyncTask).TaskID)
	t.addTask(event.Document.(*metadata.CloudSyncTask))
}

// 表记录删除处理逻辑
func (t *taskProcessor) changeOnDelete(event *stypes.Event) {
	blog.V(4).Info("OnDelete event")
	// 由于不知道删除的是哪一个任务，故进行任务的重新分配
	t.dispatchTasks()
}

// 不断给任务channel提供任务数据
func (t *taskProcessor) TaskChanLoop() {
	go func() {
		for {
			tasks := t.getTaskList()
			for i, _ := range tasks {
				t.taskChan <- tasks[i]
			}
			time.Sleep(time.Second * time.Duration(checkInterval))
		}
	}()
}

// 同步云资源
func (t *taskProcessor) SyncCloudResource() {
	hostChan := make(chan *metadata.CloudSyncTask, 10)
	// 根据任务类型，将任务放入不同的任务channel
	go func() {
		for {
			if task, ok := <-t.taskChan; ok {
				blog.V(3).Infof("processing taskid:%d, resource type:%s", task.TaskID, task.ResourceType)
				switch task.ResourceType {
				case "host":
					hostChan <- task
				default:
					blog.V(3).Infof("unknown resource type:%s, ignore it!", task.ResourceType)
				}
			}
		}
	}()

	// 云主机同步器处理同步任务
	for i := 1; i <= syncorNum; i++ {
		syncor := NewHostSyncor(i, t.logics, t.db)
		go func(syncor *HostSyncor) {
			for {
				task := <-hostChan
				syncor.Sync(task)
			}
		}(syncor)
	}
}

// 获取资源同步任务表的所有任务
func (t *taskProcessor) getTasksFromTable() ([]*metadata.CloudSyncTask, error) {
	option := &metadata.SearchCloudOption{Page: metadata.BasePage{
		Limit: common.BKNoLimit,
	}}
	result, err := t.logics.CoreAPI.CoreService().Cloud().SearchSyncTask(context.Background(), header, option)
	if err != nil {
		blog.Errorf("getTasksFromTable failed, err: %v", err)
		return nil, err
	}
	res := make([]*metadata.CloudSyncTask, 0)
	for i, _ := range result.Info {
		res = append(res, &result.Info[i])
	}
	blog.V(3).Infof("getTasksFromTable len(tasks):%d", len(res))
	return res, nil
}

// 根据服务节点设置哈希环
func (t *taskProcessor) setHashring(serversAddrs []string) {
	// 清空哈希环
	t.hashring.Set([]string{})
	// 添加所有子节点
	for _, addr := range serversAddrs {
		t.hashring.Add(addr)
	}
}

// 分配任务，清空任务列表后，将表中所有任务里属于自己的放入任务队列
func (t *taskProcessor) dispatchTasks() error {
	t.clearTaskList()
	tasks, err := t.getTasksFromTable()
	if err != nil {
		blog.Errorf("getTasksFromTable err:%s", err.Error())
		return err
	}
	for i, _ := range tasks {
		t.addTask(tasks[i])
	}
	blog.V(3).Infof("dispatchTasks is done, taskids:%#v", t.getTaskids())
	return nil
}

// 添加属于自己的任务到当前任务队列
func (t *taskProcessor) addTask(task *metadata.CloudSyncTask) error {
	if node, err := t.hashring.Get(fmt.Sprintf("%d", task.TaskID)); err != nil {
		blog.Errorf("hashring Get err:%s", err.Error())
		return err
	} else {
		if node == t.addrport {
			t.mu.Lock()
			defer t.mu.Unlock()
			t.tasklist[task.TaskID] = task
		}
	}
	return nil
}

// 获取任务列表的所有任务
func (t *taskProcessor) getTaskList() []*metadata.CloudSyncTask {
	tasks := []*metadata.CloudSyncTask{}
	t.mu.RLock()
	defer t.mu.RUnlock()
	for taskid, _ := range t.tasklist {
		tasks = append(tasks, t.tasklist[taskid])
	}
	return tasks
}

// 清空任务列表
func (t *taskProcessor) clearTaskList() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.tasklist = map[int64]*metadata.CloudSyncTask{}
}

// 获取任务列表的所有任务id
func (t *taskProcessor) getTaskids() []int64 {
	taskids := []int64{}
	t.mu.RLock()
	defer t.mu.RUnlock()
	for taskid, _ := range t.tasklist {
		taskids = append(taskids, taskid)
	}
	return taskids
}
