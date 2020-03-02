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

package taskprocess

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/hashring"
	"configcenter/src/common/types"
	"configcenter/src/common/zkclient"
	"configcenter/src/storage/dal"

	"github.com/samuel/go-zookeeper/zk"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type taskProcessor struct {
	client   *zkclient.ZkClient
	db       dal.RDB
	addrport string
	hashring *hashring.HashRing
	tasklist map[int64]bool
	taskChan chan int64
}

const (
	replicas   int = 10
	processNum int = 10
)

func ProcessTask(ctx context.Context, client *zkclient.ZkClient, db dal.RDB, addrport string) error {
	t := &taskProcessor{
		client:   client,
		db:       db,
		addrport: addrport,
		hashring: hashring.NewHashRing(replicas, nil),
		tasklist: make(map[int64]bool),
		taskChan: make(chan int64, 10),
	}
	if err := t.DispatchTask(ctx); err != nil {
		return err
	}
	if err := t.WatchTaskTable(ctx); err != nil {
		return err
	}
	t.TaskChanLoop()
	t.SyncCloudHost()
	return nil
}

// 监听zk节点变化，有变化时重新分配当前进程的任务列表
func (t *taskProcessor) DispatchTask(ctx context.Context) error {
	zkPath := fmt.Sprintf("%s/%s", types.CC_SERV_BASEPATH, common.GetIdentification())

	go func() {
		var children []string
		var ch <-chan zk.Event
		var err error
		cnt := 0
		for {
			children, ch, err = t.client.WatchChildren(zkPath)
			if err != nil {
				blog.Errorf("endpoints watch failed, will watch after 10s, path: %s, err: %s", zkPath, err.Error())
				switch err {
				case zk.ErrClosing, zk.ErrConnectionClosed:
					if conErr := t.client.Connect(); conErr != nil {
						blog.Errorf("fail to watch register node(%s), reason: connect closed. retry connect err:%s", zkPath, conErr.Error())
						time.Sleep(5 * time.Second)
					}
				}
				time.Sleep(5 * time.Second)
				continue
			}
			children, err = t.client.GetChildren(zkPath)
			if err != nil {
				blog.Errorf("fail to GetChildren(%s), err:%s", zkPath, err.Error())
				continue
			}
			blog.V(3).Infof("children:%#v", children)

			// 启动时需执行任务重新分配
			if cnt == 0 {
				t.ReDispatchTask(ctx, zkPath, children)
			}

			select {
			case event := <-ch:
				cnt++
				if event.Type != zk.EventNodeChildrenChanged {
					continue
				}
				blog.V(3).Infof("event:%#v", event)
				t.ReDispatchTask(ctx, zkPath, children)
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
	col := t.db.Table(common.BKTableNameCloudSyncTask)

	matchStage := bson.D{{"$match", bson.D{{"operationType", bson.D{{"$in",
		bson.A{"insert", "update", "replace"}}}}}}}
	opts := options.ChangeStream().SetMaxAwaitTime(2 * time.Second).SetFullDocument(options.UpdateLookup)
	changeStream, err := col.Watch(context.TODO(), mongo.Pipeline{matchStage}, opts)
	if err != nil {
		blog.Errorf("WatchTaskTable Watch err:%s", err.Error())
	}
	go func() {
		for changeStream.Next(context.TODO()) {
			blog.V(3).Infof("changeStream.Current:%s", changeStream.Current)
			doc := changeStream.Current.Lookup("fullDocument")
			fulldoc := map[string]interface{}{}
			err = doc.Unmarshal(&fulldoc)
			opType := changeStream.Current.Lookup("operationType").String()
			blog.V(3).Infof("fulldoc:%#v, err:%v, optype:%s", fulldoc, err, opType)
			var taskid int64
			var ok bool
			if _, ok = fulldoc["bk_task_id"]; ok {
				taskid = int64(fulldoc["bk_task_id"].(float64))
			}
			// 判断是否要添加到当前任务队列
			if t.hashring.Get(fmt.Sprintf("%d", taskid)) == t.addrport {
				if _, ok := t.tasklist[taskid]; !ok {
					t.tasklist[taskid] = true
				}
			}
		}
	}()

	return nil
}

// 获取资源同步任务表的所有任务
func (t *taskProcessor) GetAllTasks(ctx context.Context) ([]int64, error) {
	condition := map[string]interface{}{"del": map[string]interface{}{"$ne": 0}}
	taskResult := make([]map[string]interface{}, 0)
	err := t.db.Table(common.BKTableNameCloudSyncTask).Find(condition).All(ctx, &taskResult)
	if err != nil {
		return nil, err
	}
	blog.V(3).Infof("taskResult:%#v", taskResult)
	taskids := []int64{}
	for _, v := range taskResult {
		taskids = append(taskids, int64(v["bk_task_id"].(float64)))
	}
	return taskids, nil
}

// 重新分配任务
func (t *taskProcessor) ReDispatchTask(ctx context.Context, zkPath string, children []string) error {
	addrArr := []string{}
	for _, child := range children {
		childpath := zkPath + "/" + child
		data, err := t.client.Get(childpath)
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
	t.hashring.Clear()
	t.hashring.Add(addrArr...)
	taskids, err := t.GetAllTasks(ctx)
	if err != nil {
		blog.Errorf("GetAllTasks err:%s", err.Error())
		return err
	}
	for _, taskid := range taskids {
		if t.hashring.Get(fmt.Sprintf("%d", taskid)) == t.addrport {
			if _, ok := t.tasklist[taskid]; !ok {
				t.tasklist[taskid] = true
			}
		}
	}
	blog.Info("DispatchTask, tasklist:%#v", t.tasklist)
	return nil
}

// 同步云主机
func (t *taskProcessor) SyncCloudHost() {
	for i := 0; i < processNum; i++ {
		go func() {
			for {
				if taskid, ok := <-t.taskChan; ok {
					blog.Infof("******processing taskid:%d", taskid)
				}
			}
		}()
	}
}

// 不断给任务channel提供任务数据
func (t *taskProcessor) TaskChanLoop() {
	go func() {
		for {
			for taskid, _ := range t.tasklist {
				t.taskChan <- taskid
			}
			time.Sleep(time.Second)
		}
	}()
}
