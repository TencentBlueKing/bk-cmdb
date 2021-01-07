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
	"configcenter/src/common/metadata"
	"configcenter/src/common/zkclient"
	"configcenter/src/scene_server/cloud_server/logics"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/reflector"
	stypes "configcenter/src/storage/stream/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"stathat.com/c/consistent"
)

var (
	// mongo server对于满足change stream查询的最大等待时间
	maxAwaitTime = time.Second * 10
)

// 任务调度器
type taskScheduler struct {
	zkClient   *zkclient.ZkClient
	logics     *logics.Logics
	addrport   string
	reflector  reflector.Interface
	hashring   *consistent.Consistent
	tasklist   map[string]*metadata.CloudSyncTask
	mu         sync.RWMutex
	listerDone chan bool
}

// 调度器配置
type SchedulerConf struct {
	ZKClient  *zkclient.ZkClient
	Logics    *logics.Logics
	AddrPort  string
	MongoConf local.MongoConf
}

// 调度器实例创建
func NewTaskScheduler(conf *SchedulerConf) (*taskScheduler, error) {
	reflector, err := reflector.NewReflector(conf.MongoConf)
	if err != nil {
		blog.Errorf("NewReflector failed, mongoConf: %#v, err: %s", conf.MongoConf, err.Error())
		return nil, err
	}
	return &taskScheduler{
		zkClient:   conf.ZKClient,
		logics:     conf.Logics,
		addrport:   conf.AddrPort,
		hashring:   consistent.New(),
		tasklist:   make(map[string]*metadata.CloudSyncTask),
		reflector:  reflector,
		listerDone: make(chan bool),
	}, nil
}

// 调度云同步任务
func (t *taskScheduler) Schedule(ctx context.Context) error {
	// 监听任务表事件
	if err := t.watchTaskTable(ctx); err != nil {
		return err
	}

	// 监听服务进程节点
	if err := t.watchServerNode(); err != nil {
		return err
	}

	return nil
}

// 监听zk的cloudserver节点变化，有变化时重置哈希环
func (t *taskScheduler) watchServerNode() error {
	go func() {
		for servers := range t.logics.Discovery().CloudServer().GetServersChan() {
			t.setHashring(servers)
		}
	}()
	return nil
}

// 任务表事件结构
type taskEvent struct {
	metadata.CloudSyncTask `json:",inline" bson:",inline"`
	Oid                    primitive.ObjectID `json:"_id" bson:"_id"`
}

// 监听云资源同步任务表事件，有变更时进行相应的处理
func (t *taskScheduler) watchTaskTable(ctx context.Context) error {
	opts := &stypes.ListWatchOptions{
		Options: stypes.Options{
			MaxAwaitTime: &maxAwaitTime,
			EventStruct:  new(taskEvent),
			Collection:   common.BKTableNameCloudSyncTask,
		},
	}
	cap := &reflector.Capable{
		reflector.OnChangeEvent{
			OnAdd:        t.changeOnAdd,
			OnUpdate:     t.changeOnUpdate,
			OnDelete:     t.changeOnDelete,
			OnLister:     t.changeOnLister,
			OnListerDone: t.changeOnListerDone,
		},
	}

	return t.reflector.ListWatcher(ctx, opts, cap)
}

// 表记录新增处理逻辑
func (t *taskScheduler) changeOnAdd(event *stypes.Event) {
	blog.V(4).Infof("OnAdd event, taskid:%d", event.Document.(*taskEvent).TaskID)
	t.addTask(event.Oid, &event.Document.(*taskEvent).CloudSyncTask)
}

// 表记录更新处理逻辑
func (t *taskScheduler) changeOnUpdate(event *stypes.Event) {
	blog.V(4).Infof("OnUpdate event, taskid:%d", event.Document.(*taskEvent).TaskID)
	t.addTask(event.Oid, &event.Document.(*taskEvent).CloudSyncTask)
}

// 表记录删除处理逻辑
func (t *taskScheduler) changeOnDelete(event *stypes.Event) {
	blog.V(4).Infof("OnDelete event, oid:%s", event.Oid)
	t.delTask(event.Oid)
}

// 冷启动时已有表记录的处理逻辑
func (t *taskScheduler) changeOnLister(event *stypes.Event) {
	blog.V(4).Infof("changeOnLister event, taskid:%d", event.Document.(*taskEvent).TaskID)
	t.addTask(event.Oid, &event.Document.(*taskEvent).CloudSyncTask)
}

// 冷启动时已有表记录获取完成时的处理逻辑
func (t *taskScheduler) changeOnListerDone() {
	blog.V(4).Info("changeOnListerDone event")
	close(t.listerDone)
}

// 根据服务节点设置哈希环
func (t *taskScheduler) setHashring(serversAddrs []string) {
	blog.V(4).Infof("setHashring, serversAddrs:%#v", serversAddrs)
	t.mu.Lock()
	defer t.mu.Unlock()
	// 清空哈希环
	t.hashring.Set([]string{})
	// 添加所有子节点
	for _, addr := range serversAddrs {
		t.hashring.Add(addr)
	}
}

// 添加任务
func (t *taskScheduler) addTask(oid string, task *metadata.CloudSyncTask) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.tasklist[oid] = task
	return nil
}

// 删除任务
func (t *taskScheduler) delTask(oid string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.tasklist, oid)
	return nil
}

// 获取属于当前进程的任务列表
func (t *taskScheduler) GetTaskList() ([]*metadata.CloudSyncTask, error) {
	tasks := []*metadata.CloudSyncTask{}
	t.mu.RLock()
	defer t.mu.RUnlock()
	for oid, _ := range t.tasklist {
		if node, err := t.hashring.Get(fmt.Sprintf("%d", t.tasklist[oid].TaskID)); err != nil {
			blog.Errorf("hashring Get err:%s", err.Error())
			return nil, err
		} else {
			if node == t.addrport {
				task := *t.tasklist[oid]
				tasks = append(tasks, &task)
			}
		}
	}
	blog.V(5).Infof("GetTaskList, len(tasks):%d", len(tasks))
	return tasks, nil
}

// 记录列表获取是否完成的channel
func (t *taskScheduler) ListerDone() chan bool {
	return t.listerDone
}
