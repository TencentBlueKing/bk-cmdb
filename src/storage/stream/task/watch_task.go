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

package task

import (
	"context"
	"time"

	"configcenter/src/common/blog"
	"configcenter/src/storage/stream/types"
)

// watchTask is the resource watch task
type watchTask struct {
	// name is the watch task name that uniquely identifies the watch task
	name string
	// collOptions is the watch collection options
	collOptions *types.WatchCollOptions
	// eventHandler is the batch event handler
	eventHandler *types.TaskBatchHandler
	// tokenHandler is the token handler
	tokenHandler types.TaskTokenHandler
	// needList defines whether to list all data before watch
	needList bool

	retryOptions *types.RetryOptions
	batchSize    int
}

type dbWatchTask struct {
	*watchTask
	dbInfo    *types.DBInfo
	eventChan chan *types.Event
	listChan  chan *types.Event
	lastToken *types.TokenInfo
	tokenChan chan struct{}
}

// maxUnhandledEventLimit if the number of unhandled events exceeds this value, block the event watch process
const maxUnhandledEventLimit = 2000

func newDBWatchTask(task *watchTask, dbInfo *types.DBInfo) (*dbWatchTask, error) {
	lastToken, err := task.tokenHandler.GetStartWatchToken(context.Background(), dbInfo.UUID, dbInfo.WatchDB)
	if err != nil {
		blog.Errorf("get task %s db %s last watch token failed, err: %v", task.name, dbInfo.UUID, err)
		return nil, err
	}

	return &dbWatchTask{
		watchTask: task,
		dbInfo:    dbInfo,
		eventChan: make(chan *types.Event, maxUnhandledEventLimit+task.batchSize),
		listChan:  make(chan *types.Event, task.batchSize),
		lastToken: lastToken,
	}, nil
}

// start execute watch task
func (t *dbWatchTask) start(stopNotifier <-chan struct{}) {
	go func() {
		// list all data before watch if this task is a list watch task
		if t.needList {
			t.lastToken = &types.TokenInfo{
				StartAtTime: &types.TimeStamp{
					Sec: uint32(time.Now().Unix()),
				},
			}

			events := make([]*types.Event, 0)
			for event := range t.listChan {
				events = append(events, event)
				if len(events) == t.batchSize {
					t.eventHandler.DoBatch(t.dbInfo, events)
				}
				if event.OperationType == types.ListDone {
					break
				}
			}
			if len(events) > 0 {
				t.eventHandler.DoBatch(t.dbInfo, events)
			}
		}

		ticker := time.NewTicker(50 * time.Millisecond)
		for {
			// get events to be handled
			events := make([]*types.Event, 0)
			for {
				select {
				case one := <-t.eventChan:
					// skip previous event with smaller token
					if !compareToken(one, t.lastToken) {
						blog.V(4).Infof("%s-%s job, skip previous event(%s)", t.name, t.dbInfo.UUID, one.String())
						continue
					}
					events = append(events, one)
					if len(events) < t.batchSize {
						continue
					}
				case <-ticker.C:
					if len(events) == 0 {
						continue
					}
				case <-stopNotifier:
					ticker.Stop()
					return
				}
				break
			}

			// handle events
			t.handleEvents(events)
		}
	}()
}

func (t *dbWatchTask) handleEvents(events []*types.Event) {
	ctx := context.Background()
	first, last := events[0], events[len(events)-1]
	rid := first.ID()
	blog.Infof("%s-%s job, received %d events, first op-time: %s, fist token: %s, rid: %s", t.name, t.dbInfo.UUID,
		len(events), first.ClusterTime.String(), first.Token.Data, rid)

	needRetry := false
	retryCnt := 0
	for {
		// get start watch token after retry to avoid conflict with another watch task
		if needRetry {
			time.Sleep(t.retryOptions.RetryDuration)
			lastToken, err := t.tokenHandler.GetStartWatchToken(ctx, t.dbInfo.UUID, t.dbInfo.WatchDB)
			if err != nil {
				blog.Errorf("get task %s db %s token failed, err: %v, rid: %s", t.name, t.dbInfo.UUID, err, rid)
				time.Sleep(t.retryOptions.RetryDuration)
				continue
			}
			t.lastToken = lastToken

			// if current token is greater than last token, return
			if !compareToken(last, lastToken) {
				return
			}

			// remove events with smaller token that are already handled
			index := 0
			for i, event := range events {
				if compareToken(event, lastToken) {
					break
				}
				index = i + 1
			}
			events = events[index:]
		}

		// handle events, if all events are handled, just update last watch token
		if len(events) > 0 {
			needRetry = t.eventHandler.DoBatch(t.dbInfo, events)
			if needRetry {
				if retryCnt < t.retryOptions.MaxRetryCount {
					retryCnt++
					continue
				}
			}
		}

		// update last watch token, retry if failed
		lastToken := &types.TokenInfo{
			Token:       last.Token.Data,
			StartAtTime: &last.ClusterTime,
		}
		if err := t.tokenHandler.SetLastWatchToken(ctx, t.dbInfo.UUID, t.dbInfo.WatchDB, lastToken); err != nil {
			blog.Errorf("set task %s db %s last watch token(%+v) failed, err: %v, rid: %s", t.name, t.dbInfo.UUID,
				*lastToken, err, rid)
			needRetry = true
			continue
		}
		t.lastToken = lastToken
		select {
		case t.tokenChan <- struct{}{}:
		default:
		}
		return
	}
}
