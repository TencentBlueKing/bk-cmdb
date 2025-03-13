/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

// Package task defines event watch task logics
package task

import (
	"context"
	"fmt"
	"time"

	"configcenter/src/apimachinery/discovery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/stream/event"
	"configcenter/src/storage/stream/loop"
	"configcenter/src/storage/stream/types"
)

// Task is the event watch task that contains all resource watch tasks
type Task struct {
	// eventMap is the db uuid to event instance map
	eventMap map[string]*event.Event
	// loopWatch is the db uuid to loop watch instance map
	loopWatch map[string]*loop.LoopsWatch
	// dbClients is the db uuid to db client map
	dbClients map[string]local.DB
	// watchClients is the db uuid to watch client map
	watchClients map[string]*local.Mongo
	// watchTasks is the task name to watch task map
	watchTasks map[string]*watchTask

	// these options are used to generate loop watch options
	majorityCommitted *bool
	maxAwaitTime      *time.Duration

	// stopNotifier is used when user need to stop loop events and release related resources.
	// It's a optional option. when it's not set(as is nil), then the loop will not exit forever.
	// Otherwise, user can use it to stop loop events.
	// When a user want to stop the loop, the only thing that a user need to do is to just
	// **close** this stop notifier channel.
	// Attention:
	// Close this notifier channel is the only way to stop loop correctly.
	// Do not send data to this channel.
	stopNotifier <-chan struct{}
}

// New create a new watch task instance
func New(db, watchDB dal.Dal, isMaster discovery.ServiceManageInterface, opts *types.NewTaskOptions) (*Task, error) {
	if err := opts.Validate(); err != nil {
		blog.Errorf("validate new task options(%+v) failed, err: %v", opts, err)
		return nil, err
	}

	t := &Task{
		eventMap:     make(map[string]*event.Event),
		loopWatch:    make(map[string]*loop.LoopsWatch),
		dbClients:    make(map[string]local.DB),
		watchClients: make(map[string]*local.Mongo),
		watchTasks:   make(map[string]*watchTask),
		stopNotifier: opts.StopNotifier,
	}

	watchDBRelation, err := genWatchDBRelationMap(watchDB)
	if err != nil {
		return nil, err
	}

	// generate watch db uuid to watch db client map
	watchDBClientMap := make(map[string]*local.Mongo)
	err = watchDB.ExecForAllDB(func(db local.DB) error {
		dbClient, ok := db.(*local.Mongo)
		if !ok {
			return fmt.Errorf("watch db is not an instance of local mongo")
		}
		watchDBClientMap[dbClient.GetMongoClient().UUID()] = dbClient
		return nil
	})
	if err != nil {
		blog.Errorf("get all watch db client failed, err: %v", err)
		return nil, err
	}

	// generate db uuid to db client & watch db client & loop watch instance map
	err = db.ExecForAllDB(func(db local.DB) error {
		dbClient, ok := db.(*local.Mongo)
		if !ok {
			return fmt.Errorf("db to be watched is not an instance of local mongo")
		}
		mongoClient := dbClient.GetMongoClient()
		uuid := mongoClient.UUID()

		watchDBUUID, exists := watchDBRelation[uuid]
		if !exists {
			blog.Warnf("db %s has no watch db", uuid)
			return nil
		}

		watchClient, exists := watchDBClientMap[watchDBUUID]
		if !exists {
			return fmt.Errorf("db %s related watch db %s is invalid", uuid, watchDBUUID)
		}
		t.watchClients[uuid] = watchClient
		t.dbClients[uuid] = dbClient

		eventInst, err := event.NewEvent(mongoClient.Client(), mongoClient.DBName(), uuid)
		if err != nil {
			return fmt.Errorf("new event for db %s failed, err: %v", uuid, err)
		}
		t.eventMap[uuid] = eventInst

		loopWatch, err := loop.NewLoopWatch(eventInst, isMaster)
		if err != nil {
			return fmt.Errorf("new loop watch for db %s failed, err: %v", uuid, err)
		}
		t.loopWatch[uuid] = loopWatch
		return nil
	})
	if err != nil {
		blog.Errorf("generate db uuid related map failed, err: %v", err)
		return nil, err
	}

	return t, nil
}

// AddLoopOneTask add a loop watch task that handles one event at one time
func (t *Task) AddLoopOneTask(opts *types.LoopOneTaskOptions) error {
	if err := opts.Validate(); err != nil {
		blog.Errorf("validate loop batch task options(%s) failed, err: %v", opts.Name, err)
		return err
	}

	batchOpts := &types.LoopBatchTaskOptions{
		WatchTaskOptions: opts.WatchTaskOptions,
		BatchSize:        1,
		EventHandler: &types.TaskBatchHandler{
			DoBatch: func(dbInfo *types.DBInfo, es []*types.Event) bool {
				for _, e := range es {
					var retry bool
					switch e.OperationType {
					case types.Insert:
						retry = opts.EventHandler.DoAdd(dbInfo, e)
					case types.Update, types.Replace:
						retry = opts.EventHandler.DoUpdate(dbInfo, e)
					case types.Delete:
						retry = opts.EventHandler.DoDelete(dbInfo, e)
					default:
						blog.Warnf("received unsupported operation type for %s job, doc: %s", opts.Name, e.DocBytes)
						continue
					}
					if retry {
						return retry
					}
				}
				return false
			},
		},
	}

	return t.addWatchTask(batchOpts, false)
}

// AddLoopBatchTask add a loop watch task that handles batch events
func (t *Task) AddLoopBatchTask(opts *types.LoopBatchTaskOptions) error {
	if err := opts.Validate(); err != nil {
		blog.Errorf("validate loop batch task options(%s) failed, err: %v", opts.Name, err)
		return err
	}
	return t.addWatchTask(opts, false)
}

// AddListWatchTask add a list watch task
func (t *Task) AddListWatchTask(opts *types.LoopBatchTaskOptions) error {
	if err := opts.Validate(); err != nil {
		blog.Errorf("validate list watch task options(%s) failed, err: %v", opts.Name, err)
		return err
	}
	return t.addWatchTask(opts, true)
}

func (t *Task) addWatchTask(opts *types.LoopBatchTaskOptions, needList bool) error {
	_, exists := t.watchTasks[opts.Name]
	if exists {
		return fmt.Errorf("loop watch task %s already exists", opts.Name)
	}

	if opts.MajorityCommitted != nil && *opts.MajorityCommitted {
		t.majorityCommitted = opts.MajorityCommitted
	}
	if opts.MaxAwaitTime != nil && (t.maxAwaitTime == nil || *opts.MaxAwaitTime > *t.maxAwaitTime) {
		t.maxAwaitTime = opts.MaxAwaitTime
	}

	t.watchTasks[opts.Name] = &watchTask{
		name:         opts.Name,
		collOptions:  opts.CollOpts,
		eventHandler: opts.EventHandler,
		tokenHandler: opts.TokenHandler,
		needList:     needList,
		retryOptions: opts.RetryOptions,
		batchSize:    opts.BatchSize,
	}

	return nil
}

// Start execute all watch tasks
func (t *Task) Start() error {
	if len(t.watchTasks) == 0 {
		return nil
	}

	// generate task name to collection options map and db uuid to task name to db watch tasks map by watch task info
	collOptions := make(map[string]types.WatchCollOptions)
	listCollOptions := make(map[string]types.CollectionOptions)
	dbWatchTasks := make(map[string]map[string]*dbWatchTask)
	var batchSize int
	for taskName, task := range t.watchTasks {
		collOptions[taskName] = *task.collOptions
		if task.needList {
			listCollOptions[taskName] = task.collOptions.CollectionOptions
		}
		if task.batchSize > batchSize {
			batchSize = task.batchSize
		}
		for uuid, dbClient := range t.dbClients {
			dbTask, err := newDBWatchTask(task, &types.DBInfo{
				UUID:    uuid,
				WatchDB: t.watchClients[uuid],
				CcDB:    dbClient,
			})
			if err != nil {
				return err
			}
			if _, exists := dbWatchTasks[uuid]; !exists {
				dbWatchTasks[uuid] = make(map[string]*dbWatchTask)
			}
			dbWatchTasks[uuid][taskName] = dbTask
		}
	}

	// list data for all list watch tasks
	if len(listCollOptions) > 0 {
		err := t.startList(listCollOptions, batchSize, dbWatchTasks)
		if err != nil {
			return err
		}
	}

	// loop watch all db events for all tasks
	err := t.startLoopWatch(collOptions, dbWatchTasks, batchSize)
	if err != nil {
		return err
	}

	// run watch tasks for all dbs
	for _, dbTaskMap := range dbWatchTasks {
		for _, dbTask := range dbTaskMap {
			dbTask.start(t.stopNotifier)
		}
	}

	return nil
}

func (t *Task) startList(listCollOptions map[string]types.CollectionOptions, batchSize int,
	dbWatchTasks map[string]map[string]*dbWatchTask) error {

	for uuid, eventInst := range t.eventMap {
		ctx := util.SetDBReadPreference(context.Background(), common.SecondaryPreferredMode)
		opt := &types.ListOptions{
			CollOpts:  listCollOptions,
			PageSize:  &batchSize,
			WithRetry: true,
		}
		listCh, err := eventInst.List(ctx, opt)
		if err != nil {
			blog.Errorf("list db %s failed, err: %v, options: %+v", uuid, err, *opt)
			return err
		}

		go func(uuid string) {
			for e := range listCh {
				task, exists := dbWatchTasks[uuid][e.TaskID]
				if !exists {
					blog.Warnf("loop watch task %s not exists, event: %+v", e.TaskID, *e)
					continue
				}
				task.listChan <- e
			}
		}(uuid)
	}
	return nil
}

func (t *Task) startLoopWatch(collOptions map[string]types.WatchCollOptions,
	dbWatchTasks map[string]map[string]*dbWatchTask, batchSize int) error {

	for uuid, loopWatch := range t.loopWatch {
		uuid := uuid
		tokenHandler, err := newDBTokenHandler(uuid, t.watchClients[uuid], dbWatchTasks[uuid])
		if err != nil {
			return err
		}
		opts := &types.LoopBatchOptions{
			LoopOptions: types.LoopOptions{
				Name: uuid,
				WatchOpt: &types.WatchOptions{
					Options: types.Options{
						MajorityCommitted: t.majorityCommitted,
						MaxAwaitTime:      t.maxAwaitTime,
						CollOpts:          collOptions,
					},
				},
				TokenHandler: tokenHandler,
				RetryOptions: &types.RetryOptions{
					MaxRetryCount: types.DefaultRetryCount,
					RetryDuration: types.DefaultRetryDuration,
				},
				StopNotifier: t.stopNotifier,
			},
			EventHandler: &types.BatchHandler{DoBatch: func(es []*types.Event) (retry bool) {
				taskLastTokenMap := make(map[string]string)
				for _, e := range es {
					task, exists := dbWatchTasks[uuid][e.TaskID]
					if !exists {
						blog.Warnf("loop watch task %s not exists, event: %+v", e.TaskID, *e)
						continue
					}
					task.eventChan <- e
					taskLastTokenMap[e.TaskID] = e.Token.Data
				}
				tokenHandler.setTaskLastTokenInfo(taskLastTokenMap)
				return false
			}},
			BatchSize: batchSize,
		}

		err = loopWatch.WithBatch(opts)
		if err != nil {
			blog.Errorf("start loop watch for db failed, err: %v", err)
			return err
		}
	}
	return nil
}
