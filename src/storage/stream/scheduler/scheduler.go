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

// Package scheduler defines event watch task scheduler logics
package scheduler

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
	"configcenter/src/storage/stream/task"
	"configcenter/src/storage/stream/types"
)

// Scheduler is the event watch task scheduler
type Scheduler struct {
	isMaster discovery.ServiceManageInterface

	// watchTasks is the task name to watch task map
	watchTasks map[string]*task.Task

	// eventMap is the db uuid to event instance map
	eventMap map[string]*event.Event
	// dbClients is the db uuid to db client map
	dbClients map[string]local.DB
	// watchClients is the db uuid to watch client map
	watchClients map[string]*local.Mongo

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
	stopNotifier chan struct{}
}

// New creates a new watch task scheduler
func New(db, watchDB dal.Dal, isMaster discovery.ServiceManageInterface) (*Scheduler, error) {
	t := &Scheduler{
		isMaster:     isMaster,
		eventMap:     make(map[string]*event.Event),
		dbClients:    make(map[string]local.DB),
		watchClients: make(map[string]*local.Mongo),
		watchTasks:   make(map[string]*task.Task),
		stopNotifier: make(chan struct{}),
	}

	watchDBRelation, defaultWatchDBUUID, err := genWatchDBRelationInfo(watchDB)
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
			blog.Warnf("db %s has no watch db, use default watch db for new db", uuid)
			watchDBUUID = defaultWatchDBUUID
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
		return nil
	})
	if err != nil {
		blog.Errorf("generate db uuid related map failed, err: %v", err)
		return nil, err
	}

	return t, nil
}

// AddTasks add watch tasks to scheduler
func (s *Scheduler) AddTasks(tasks ...*task.Task) error {
	for _, t := range tasks {
		_, exists := s.watchTasks[t.Name]
		if exists {
			blog.Errorf("add watch task %s to scheduler failed, task already exists", t.Name)
			return fmt.Errorf("loop watch task %s already exists", t.Name)
		}

		if t.MajorityCommitted != nil && *t.MajorityCommitted {
			s.majorityCommitted = t.MajorityCommitted
		}
		if t.MaxAwaitTime != nil && (s.maxAwaitTime == nil || *t.MaxAwaitTime > *s.maxAwaitTime) {
			s.maxAwaitTime = t.MaxAwaitTime
		}

		s.watchTasks[t.Name] = t
	}
	return nil
}

// Start execute all watch tasks
func (s *Scheduler) Start() error {
	if len(s.watchTasks) == 0 {
		blog.Errorf("no watch task to start")
		return fmt.Errorf("no watch task to start")
	}

	// generate task name to collection options map and db uuid to task name to db watch tasks map by watch task info
	taskCollOptsMap := make(map[string]types.WatchCollOptions)
	taskListCollOptsMap := make(map[string]types.CollectionOptions)
	dbWatchTasks := make(map[string]map[string]*task.DBWatchTask)
	var batchSize int
	for taskName, watchTask := range s.watchTasks {
		taskCollOptsMap[taskName] = *watchTask.CollOptions
		if watchTask.NeedList {
			taskListCollOptsMap[taskName] = watchTask.CollOptions.CollectionOptions
		}
		if watchTask.BatchSize > batchSize {
			batchSize = watchTask.BatchSize
		}
		for uuid, dbClient := range s.dbClients {
			dbTask, err := task.NewDBWatchTask(watchTask, &types.DBInfo{
				UUID:    uuid,
				WatchDB: s.watchClients[uuid],
				DB:      dbClient,
			})
			if err != nil {
				return err
			}
			if _, exists := dbWatchTasks[uuid]; !exists {
				dbWatchTasks[uuid] = make(map[string]*task.DBWatchTask)
			}
			dbWatchTasks[uuid][taskName] = dbTask
		}
	}

	// list data for all list watch tasks
	if len(taskListCollOptsMap) > 0 {
		err := s.startList(taskListCollOptsMap, batchSize, dbWatchTasks)
		if err != nil {
			return err
		}
	}

	// loop watch all db events for all tasks
	err := s.startLoopWatch(taskCollOptsMap, dbWatchTasks, batchSize)
	if err != nil {
		return err
	}

	// run watch tasks for all dbs
	for _, dbTaskMap := range dbWatchTasks {
		for _, dbTask := range dbTaskMap {
			dbTask.Start(s.stopNotifier)
		}
	}

	return nil
}

func (s *Scheduler) startList(taskListCollOptsMap map[string]types.CollectionOptions, batchSize int,
	dbWatchTasks map[string]map[string]*task.DBWatchTask) error {

	for uuid, eventInst := range s.eventMap {
		ctx := util.SetDBReadPreference(context.Background(), common.SecondaryPreferredMode)
		opt := &types.ListOptions{
			CollOpts:  taskListCollOptsMap,
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
				task.ListChan <- e
			}
		}(uuid)
	}
	return nil
}

func (s *Scheduler) startLoopWatch(taskCollOptsMap map[string]types.WatchCollOptions,
	dbWatchTasks map[string]map[string]*task.DBWatchTask, batchSize int) error {

	for uuid, dbTaskMap := range dbWatchTasks {
		watcher, err := s.newDBWatcher(uuid, dbTaskMap)
		if err != nil {
			blog.Errorf("new db watcher for db %s failed, err: %v", uuid, err)
			return err
		}

		opts := &types.WatchOptions{
			Options: types.Options{
				MajorityCommitted:       s.majorityCommitted,
				MaxAwaitTime:            s.maxAwaitTime,
				TaskCollOptsMap:         taskCollOptsMap,
				WatchFatalErrorCallback: watcher.resetWatchToken,
			},
		}
		err = watcher.loopWatch(opts, batchSize)
		if err != nil {
			blog.Errorf("start loop watch for db %s failed, err: %v", uuid, err)
			return err
		}
	}
	return nil
}

// Stop the task scheduler
func (s *Scheduler) Stop() {
	close(s.stopNotifier)
}
