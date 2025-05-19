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

package scheduler

import (
	"context"
	"sync"
	"time"

	"configcenter/src/apimachinery/discovery"
	"configcenter/src/common/blog"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/stream/event"
	"configcenter/src/storage/stream/task"
	"configcenter/src/storage/stream/types"
)

type dbWatcher struct {
	uuid        string
	watchDB     local.DB
	streamWatch *event.Event
	isMaster    discovery.ServiceManageInterface

	// taskMap is the task id to
	taskMap map[string]*task.DBWatchTask

	// lastToken is the last db watch token that has been set
	lastToken *types.TokenInfo
	// lastTokenInfo is the task id to last event token map, which is used to judge if all events have been handled
	lastTokenInfo map[string]string
	mu            sync.RWMutex

	// retryChan is used to notify the db watcher to retry watch
	retryChan chan struct{}
	// tokenChan is used to notify the db watcher to set the last token
	tokenChan chan struct{}
	// stopNotifier is used to notify the db watcher to stop
	stopNotifier <-chan struct{}
}

// newDBWatcher new db event watcher for task scheduler
func (s *Scheduler) newDBWatcher(uuid string, taskMap map[string]*task.DBWatchTask) (*dbWatcher, error) {
	watcher := &dbWatcher{
		uuid:          uuid,
		watchDB:       s.watchClients[uuid],
		streamWatch:   s.eventMap[uuid],
		isMaster:      s.isMaster,
		taskMap:       taskMap,
		lastTokenInfo: make(map[string]string),
		retryChan:     make(chan struct{}, 1),
		tokenChan:     make(chan struct{}, 1),
		stopNotifier:  s.stopNotifier,
	}

	lastToken, err := watcher.getStartWatchToken(context.Background())
	if err != nil {
		return nil, err
	}
	watcher.lastToken = lastToken

	for taskID, task := range taskMap {
		if task.LastToken != nil {
			watcher.lastTokenInfo[taskID] = task.LastToken.Token
		}
		task.TokenChan = watcher.tokenChan
	}

	go func() {
		for range watcher.tokenChan {
			watcher.setLastWatchToken()
		}
	}()
	return watcher, nil
}

func (w *dbWatcher) loopWatch(watchOpt *types.WatchOptions, batchSize int) error {
	ctx, cancel := context.WithCancel(context.Background())

	watcher, watchOpt, err := w.watch(ctx, watchOpt)
	if err != nil {
		cancel()
		return err
	}

	go func() {
		ticker := time.NewTicker(50 * time.Millisecond)
		observer := &watchObserver{
			isMaster:       w.isMaster,
			previousStatus: w.isMaster.IsMaster(),
		}

		for {
			select {
			case <-w.stopNotifier:
				cancel()
				close(w.tokenChan)
				blog.Warnf("received stop %s loop watch job notify, stopping now.", w.uuid)
				return
			case <-w.retryChan:
				// wait for a while and then do the retry work.
				time.Sleep(types.DefaultRetryDuration)

				// initialize a new retry signal, exit loop immediately for next usage.
				w.retryChan = make(chan struct{}, 1)

				// cancel the former watch
				cancel()

				ctx, cancel = context.WithCancel(context.Background())
				watcher, watchOpt, err = w.watch(ctx, watchOpt)
				if err != nil {
					// notify retry signal, exit loop
					blog.Errorf("watch db %s with opt(%+v) failed, err: %v, retry again", w.uuid, *watchOpt)
					w.notifyRetry()
					continue
				}

				blog.Errorf("%s job, the former watch loop: %s failed, start retry again from token: %+v.", w.uuid,
					w.streamWatch.DBName, watchOpt.StartAfterToken)

				ticker = time.NewTicker(50 * time.Millisecond)
				observer = &watchObserver{
					isMaster:       w.isMaster,
					previousStatus: w.isMaster.IsMaster(),
				}
			default:
				reWatch, loop := observer.canLoop()
				if reWatch {
					// stop the tick to release resource.
					ticker.Stop()
					blog.Warnf("%s job, master status has changed, try to re-watch again, db：%s", w.uuid,
						w.streamWatch.DBName)
					// trigger re-watch action now.
					w.notifyRetry()
					// exit the for loop
					continue
				}

				if !loop {
					blog.V(5).Infof("%s job, loop %s event, but not master, skip.", w.uuid, w.streamWatch.DBName)
					time.Sleep(5 * time.Second)
					continue
				}

				w.handleEvents(watcher, ticker, batchSize)
			}
		}
	}()

	return nil
}

func (w *dbWatcher) watch(ctx context.Context, watchOpt *types.WatchOptions) (*types.Watcher, *types.WatchOptions,
	error) {

	startToken, err := w.getStartWatchToken(context.Background())
	if err != nil {
		blog.Errorf("%s job, loop watch db %s, but get start watch token failed, err: %v", w.uuid, w.streamWatch.DBName,
			err)
		return nil, watchOpt, err
	}
	w.lastToken = startToken

	// update the start token of the watch options
	if len(startToken.Token) != 0 {
		watchOpt.StartAfterToken = &types.EventToken{Data: startToken.Token}
	}
	if startToken.StartAtTime != nil {
		watchOpt.StartAtTime = startToken.StartAtTime
	}

	watcher, err := w.streamWatch.Watch(ctx, watchOpt)
	if err != nil {
		blog.Errorf("%s job, run loop, but watch failed, err: %v", w.uuid, err)
		return nil, watchOpt, err
	}

	return watcher, watchOpt, nil
}

// handleEvents handle events
func (w *dbWatcher) handleEvents(watcher *types.Watcher, ticker *time.Ticker, batchSize int) {
	batchEventsLen := 0
	var first *types.Event

	// get event with loop.
	for {
		select {
		case one := <-watcher.EventChan:
			task, exists := w.taskMap[one.TaskID]
			if !exists {
				blog.Warnf("loop watch task %s not exists, event: %+v", one.TaskID, *one)
				continue
			}
			task.EventChan <- one
			w.setTaskLastTokenInfo(one.TaskID, one.Token.Data)

			if blog.V(4) {
				blog.Infof("%s job, received %s event, detail: %s, op-time: %s, rid: %s", w.uuid,
					w.streamWatch.DBName, one.String(), one.ClusterTime.String(), one.ID())
			}

			// calculate event count, try to get more event for a batch
			batchEventsLen++
			if first == nil {
				first = one
			}
			if batchEventsLen < batchSize {
				// continue to get more events
				continue
			}
		case <-ticker.C:
			// handle with batch event.
			if batchEventsLen == 0 {
				// ticks, but no events received, loop next round to get events.
				return
			}
		case <-w.stopNotifier:
			ticker.Stop()
			blog.Warnf("received stop %s loop watch job notify, stopping now.", w.uuid)
			return
		}

		// break the for loop to handle event for now.
		break
	}

	blog.Infof("%s job, received %s batch %d events, first op-time: %s, rid: %s.", w.uuid, w.streamWatch.DBName,
		batchEventsLen, first.ClusterTime.String(), first.ID())
}

func (w *dbWatcher) notifyRetry() {
	select {
	case w.retryChan <- struct{}{}:
	default:
	}
}
