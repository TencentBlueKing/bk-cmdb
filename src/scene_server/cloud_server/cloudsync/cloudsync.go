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
	"configcenter/src/common/mapstr"
	"configcenter/src/storage/dal/mongo/local"
	"context"
	"encoding/json"
	"fmt"
	//"strings"
	"sync"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/types"
	"configcenter/src/common/zkclient"
	"configcenter/src/scene_server/cloud_server/logics"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/reflector"
	stypes "configcenter/src/storage/stream/types"

	"github.com/samuel/go-zookeeper/zk"
	"stathat.com/c/consistent"
)

type taskProcessor struct {
	ctx       context.Context
	zkClient  *zkclient.ZkClient
	db        dal.DB
	logics    *logics.Logics
	addrport  string
	reflector reflector.Interface
	hashring  *consistent.Consistent
	tasklist  map[int64]bool
	taskChan  chan int64
	mu        sync.RWMutex
}

const (
	// 任务处理者数量
	processorNum int = 10
	// 循环检查任务列表的间隔
	checkInterval int = 5
)

var zkPath string

// The document created or modified by curd operations in mongo
// 可参考文档https://docs.mongodb.com/manual/reference/change-events/#change-stream-output
type FullDoc struct {
	TaskID int64 `json:"bk_task_id" bson:"bk_task_id"`
}

type SyncConf struct {
	ZKClient  *zkclient.ZkClient
	DB        dal.DB
	Logics    *logics.Logics
	AddrPort  string
	MongoConf local.MongoConf
}

// 处理云资源同步任务
func CloudSync(ctx context.Context, conf *SyncConf) error {
	t := &taskProcessor{
		ctx:      ctx,
		zkClient: conf.ZKClient,
		db:       conf.DB,
		logics:   conf.Logics,
		addrport: conf.AddrPort,
		hashring: consistent.New(),
		tasklist: make(map[int64]bool),
		taskChan: make(chan int64, 10),
	}

	var err error
	t.reflector, err = reflector.NewReflector(conf.MongoConf)
	if err != nil {
		blog.Errorf("NewReflector failed, mongoConf: %#v, err: %s", conf.MongoConf, err.Error())
		return err
	}

	// 监听任务进程节点
	if err := t.WatchTaskNode(ctx); err != nil {
		return err
	}
	// 监听任务表事件
	if err := t.WatchTaskTable(ctx); err != nil {
		return err
	}

	// 不断给任务channel提供任务数据
	t.TaskChanLoop()
	// 同步云资源
	t.SyncCloudResource()
	return nil
}

// 监听zk节点变化，有变化时重新分配当前进程的任务列表
func (t *taskProcessor) WatchTaskNode(ctx context.Context) error {
	zkPath = fmt.Sprintf("%s/%s", types.CC_SERV_BASEPATH, common.GetIdentification())
	go func() {
		var ch <-chan zk.Event
		var err error
		cnt := 0
		for {
			_, ch, err = t.zkClient.WatchChildren(zkPath)
			if err != nil {
				blog.Errorf("endpoints watch failed, will watch after 10s, path: %s, err: %s", zkPath, err.Error())
				switch err {
				case zk.ErrClosing, zk.ErrConnectionClosed:
					if conErr := t.zkClient.Connect(); conErr != nil {
						blog.Errorf("fail to watch register node(%s), reason: connect closed. retry connect err:%s", zkPath, conErr.Error())
						time.Sleep(5 * time.Second)
					}
				}
				time.Sleep(5 * time.Second)
				continue
			}

			// 启动时需执行任务重新分配
			if cnt == 0 {
				if err := t.dispatchTasks(ctx, zkPath); err != nil {
					continue
				}
			}

			select {
			case event := <-ch:
				cnt++
				if event.Type != zk.EventNodeChildrenChanged {
					continue
				}
				if err := t.dispatchTasks(ctx, zkPath); err != nil {
					continue
				}
			case <-ctx.Done():
				blog.Warnf("cloudserver endpoints watch stopped because of context done.")
				return
			}
		}
	}()
	return nil
}

// 监控云资源同步任务表事件，发现有变更则判断是否将该任务加入到当前进程的任务列表里
func (t *taskProcessor) WatchTaskTable(ctx context.Context) error {
	opts := &stypes.WatchOptions{
		Options: stypes.Options{
			EventStruct: new(metadata.CloudSyncTask),
			Collection:  common.BKTableNameCloudSyncTask,
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
	t.addTask(event.Document.(*metadata.CloudSyncTask).TaskID)
}

// 表记录更新处理逻辑
func (t *taskProcessor) changeOnUpdate(event *stypes.Event) {
	blog.V(4).Infof("OnUpdate event, taskid:%d", event.Document.(*metadata.CloudSyncTask).TaskID)
	t.addTask(event.Document.(*metadata.CloudSyncTask).TaskID)
}

// 表记录删除处理逻辑
func (t *taskProcessor) changeOnDelete(event *stypes.Event) {
	blog.V(4).Info("OnDelete event")
	t.dispatchTasks(t.ctx, zkPath)
}

// 不断给任务channel提供任务数据
func (t *taskProcessor) TaskChanLoop() {
	go func() {
		for {
			taskids := t.getTaskList()
			for _, taskid := range taskids {
				t.taskChan <- taskid
			}
			time.Sleep(time.Second * time.Duration(checkInterval))
		}
	}()
}

// 同步云资源
func (t *taskProcessor) SyncCloudResource() {
	for i := 0; i < processorNum; i++ {
		go func() {
			for {
				if taskid, ok := <-t.taskChan; ok {
					task, err := t.getTaskDetail(taskid)
					if err != nil {
						blog.V(3).Infof("getTaskDetail err:%v", err)
						continue
					}
					blog.V(3).Infof("processing taskid:%d, resource type:%s", taskid, task.ResourceType)
					switch task.ResourceType {
					case "host":
						t.SyncCloudHost(task)
					default:
						blog.V(3).Infof("unknown resource type:%s, ignore it!", task.ResourceType)
					}
				}
			}
		}()
	}
}

// 获取资源同步任务表的所有任务
func (t *taskProcessor) getTasksFromTable(ctx context.Context) ([]int64, error) {
	result := make([]*metadata.CloudSyncTask, 0)
	err := t.db.Table(common.BKTableNameCloudSyncTask).Find(nil).All(ctx, &result)
	if err != nil {
		blog.Errorf("getTasksFromTable failed, err: %v", err)
		return nil, err
	}
	taskids := []int64{}
	for _, v := range result {
		taskids = append(taskids, v.TaskID)
	}
	blog.V(3).Infof("getTasksFromTable len(taskids):%d", len(taskids))
	return taskids, nil
}

// 监听zk路径上的任务节点，有变动时重新分配任务
func (t *taskProcessor) dispatchTasks(ctx context.Context, zkPath string) error {
	// 获取监控路径下的所有子节点
	children, err := t.zkClient.GetChildren(zkPath)
	if err != nil {
		blog.Errorf("fail to GetChildren(%s), err:%s", zkPath, err.Error())
		return err
	}
	blog.V(3).Infof("dispatchTasks zkPath:%s, children:%#v", zkPath, children)
	addrArr := []string{}
	for _, child := range children {
		childpath := zkPath + "/" + child
		data, err := t.zkClient.Get(childpath)
		if err != nil {
			blog.Errorf("fail to get node(%s), err:%s", childpath, err.Error())
			continue
		}
		info := types.ServerInfo{}
		err = json.Unmarshal([]byte(data), &info)
		if err != nil {
			blog.Errorf("fail to unmarshal data(%v), err:%s", data, err.Error())
			return err
		}
		addrArr = append(addrArr, info.Instance())

	}
	// 清空哈希环
	t.hashring.Set([]string{})
	// 添加所有子节点
	for _, addr := range addrArr {
		t.hashring.Add(addr)

	}
	// 清空任务列表后，将表中所有任务里属于自己的放入任务队列
	t.clearTaskList()
	taskids, err := t.getTasksFromTable(ctx)
	if err != nil {
		blog.Errorf("getTasksFromTable err:%s", err.Error())
		return err
	}
	for _, taskid := range taskids {
		t.addTask(taskid)
	}
	blog.V(3).Infof("finished dispatchTasks, tasklist:%#v", t.tasklist)
	return nil
}

// 添加属于自己的任务到当前任务队列
func (t *taskProcessor) addTask(taskid int64) error {
	if node, err := t.hashring.Get(fmt.Sprintf("%d", taskid)); err != nil {
		blog.Errorf("hashring Get err:%s", err.Error())
		return err
	} else {
		if node == t.addrport {
			t.mu.Lock()
			defer t.mu.Unlock()
			if _, ok := t.tasklist[taskid]; !ok {
				t.tasklist[taskid] = true
			}
		}
	}
	return nil
}

// 获取任务列表的所有任务
func (t *taskProcessor) getTaskList() []int64 {
	taskids := []int64{}
	t.mu.RLock()
	defer t.mu.RUnlock()
	for taskid, _ := range t.tasklist {
		taskids = append(taskids, taskid)
	}
	return taskids
}

// 清空任务列表
func (t *taskProcessor) clearTaskList() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.tasklist = map[int64]bool{}
}

// 根据任务id获取任务详情
func (t *taskProcessor) getTaskDetail(taskid int64) (*metadata.CloudSyncTask, error) {
	cond := mapstr.MapStr{common.BKCloudSyncTaskID: taskid}
	result := make([]*metadata.CloudSyncTask, 0)
	err := t.db.Table(common.BKTableNameCloudSyncTask).Find(cond).All(context.Background(), &result)
	if err != nil {
		blog.Errorf("getTaskDetail err:%v", err.Error())
		return nil, err
	}
	if len(result) > 0 {
		return result[0], nil
	}
	return nil, nil
}

// 根据账号id获取账号详情
func (t *taskProcessor) getAccountDetail(accountID int64) (*metadata.CloudAccount, error) {
	cond := mapstr.MapStr{common.BKCloudAccountID: accountID}
	result := make([]*metadata.CloudAccount, 0)
	err := t.db.Table(common.BKTableNameCloudAccount).Find(cond).All(context.Background(), &result)
	if err != nil {
		blog.Errorf("getAccountDetail err:%v", err.Error())
		return nil, err
	}
	if len(result) > 0 {
		return result[0], nil
	}
	return nil, nil
}

// 更新任务同步状态
func (t *taskProcessor) updateTaskState(taskid int64, status int) error {
	cond := mapstr.MapStr{common.BKCloudSyncTaskID: taskid}
	option := mapstr.MapStr{common.BKCloudSyncStatus: status}
	if status == metadata.CloudSyncSuccess || status == metadata.CloudSyncFail {

		option.Set(common.BKCloudLastSyncTime, time.Now().Format("2006-01-02 15:04:05"))
	}
	if err := t.db.Table(common.BKTableNameCloudSyncTask).Update(context.Background(), cond, option); err != nil {
		if err != nil {
			blog.Errorf("updateTaskState err:%v", err.Error())
			return err
		}
	}
	return nil
}
