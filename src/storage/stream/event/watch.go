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

package event

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"configcenter/src/common/blog"
	"configcenter/src/storage/stream/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Watch mongodb change stream events
func (e *Event) Watch(ctx context.Context, opts *types.WatchOptions) (*types.Watcher, error) {
	if err := opts.CheckSetDefault(); err != nil {
		return nil, err
	}

	eventChan := make(chan *types.Event, types.DefaultEventChanSize)
	go func() {
		pipeline, streamOptions, collOptsInfo := generateOptions(&opts.Options)

		blog.InfoJSON("start watch db %s with pipeline: %s, options: %s, stream options: %s", e.DBName, pipeline, opts,
			streamOptions)

		stream, streamOptions, watchOpts, err := e.watch(ctx, pipeline, streamOptions, &opts.Options)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				// if error is context cancelled, then loop watch will exit at the same time
				return
			}
			blog.Fatalf("mongodb watch failed with conf: %+v, err: %v", *opts, err)
		}

		loopOpts := &loopWatchOpts{
			Options:       watchOpts,
			streamOptions: streamOptions,
			stream:        stream,
			pipeline:      pipeline,
			eventChan:     eventChan,
			collOptsInfo:  collOptsInfo,
		}
		go e.loopWatch(ctx, loopOpts)
	}()

	watcher := &types.Watcher{
		EventChan: eventChan,
	}
	return watcher, nil
}

func (e *Event) watch(ctx context.Context, pipeline mongo.Pipeline, streamOptions *options.ChangeStreamOptions,
	opts *types.Options) (*mongo.ChangeStream, *options.ChangeStreamOptions, *types.Options, error) {

	stream, err := e.client.
		Database(e.database).
		Watch(ctx, pipeline, streamOptions)

	if err != nil && isFatalError(err) {
		// TODO: send alarm immediately.
		blog.Errorf("mongodb watch db: %s got a fatal error, skip resume token and retry, err: %v", e.DBName, err)
		// reset the resume token, because we can not use the former resume token to watch success for now.
		streamOptions.StartAfter = nil
		opts.StartAfterToken = nil
		// cause we have already got a fatal error, we can not try to watch from where we lost.
		// so re-watch from 1 minutes ago to avoid lost events.
		// Note: apparently, we may got duplicate events with this re-watch
		startAtTime := uint32(time.Now().Unix()) - 60
		streamOptions.StartAtOperationTime = &primitive.Timestamp{
			T: startAtTime,
			I: 0,
		}
		opts.StartAtTime = &types.TimeStamp{Sec: startAtTime}

		if opts.WatchFatalErrorCallback != nil {
			err := opts.WatchFatalErrorCallback(types.TimeStamp{Sec: startAtTime})
			if err != nil {
				blog.Errorf("do watch fatal error callback for db %s failed, err: %v", e.DBName, err)
			}
		}

		blog.InfoJSON("start watch db %s with pipeline: %s, options: %s, stream options: %s", e.DBName, pipeline,
			opts, streamOptions)

		stream, err = e.client.
			Database(e.database).
			Watch(ctx, pipeline, streamOptions)
	}

	return stream, streamOptions, opts, err
}

type loopWatchOpts struct {
	*types.Options
	streamOptions *options.ChangeStreamOptions
	stream        *mongo.ChangeStream
	pipeline      mongo.Pipeline
	eventChan     chan *types.Event
	currentToken  types.EventToken
	collOptsInfo  *parsedCollOptsInfo
}

func (e *Event) loopWatch(ctx context.Context, opts *loopWatchOpts) {
	retry := false
	opts.currentToken = types.EventToken{Data: ""}

	e.setCleaner(ctx, opts.eventChan)

	// init collection to task ids map
	collTasksMap := make(map[string][]string)

	for {
		// no events, try cancel watch here.
		select {
		case <-ctx.Done():
			blog.Warnf("received stopped loop watch signal, stop watch db: %s, name: %s, err: %v", e.database, e.DBName,
				ctx.Err())

			if opts.stream != nil {
				opts.stream.Close(context.Background())
			}

			return
		default:
		}

		if retry {
			opts, retry = e.retryWatch(ctx, opts)
			if retry {
				continue
			}
		}

		for opts.stream.Next(ctx) {
			// still have events, try cancel steam here.
			select {
			case <-ctx.Done():
				blog.Warnf("received stopped loop watch signal, stop loop next, watch db: %s, db name: %s, err: %v",
					e.database, e.DBName, ctx.Err())
				opts.stream.Close(context.Background())
				return
			default:
			}

			opts, retry = e.handleStreamEvent(ctx, opts, collTasksMap)
			if retry {
				break
			}
		}

		if err := opts.stream.Err(); err != nil {
			blog.ErrorJSON("mongodb watch encountered a error, conf: %s, err: %s", *opts, err)
			opts.stream.Close(ctx)
			retry = true
			continue
		}
	}
}

// setCleaner set up a monitor to close the cursor when the context is canceled.
// this is useful to release stream resource when this watch is canceled outside with context is canceled.
func (e *Event) setCleaner(ctx context.Context, eventChan chan *types.Event) {
	go func() {
		select {
		case <-ctx.Done():
			blog.Warnf("received stopped loop watch db: %s signal, close cursor now, err: %v", e.DBName, ctx.Err())

			// even though we may already close the stream, but there may still have events in the stream's
			// batch cursor, so we need to consume a event, so that we can release the stream resource
			select {
			// try consume a event, so that stream.Next(ctx) can be called to release the stream resources.
			case <-eventChan:
				blog.Warnf("received stopped loop watch db: %s signal, consumed a event", e.DBName)

			default:
				// no events, and stream resource will be recycled in the next round.
			}

			return
		}
	}()
}

func (e *Event) retryWatch(ctx context.Context, opts *loopWatchOpts) (*loopWatchOpts, bool) {
	streamOptions := opts.streamOptions

	time.Sleep(5 * time.Second)
	if len(opts.currentToken.Data) != 0 {
		// if error occurs, then retry watch and start from the last token.
		// so that we can continue the event from where it just broken.
		streamOptions.StartAtOperationTime = nil
		streamOptions.SetStartAfter(opts.currentToken)
	}

	// if start at operation time and start after token is both set, use resume token instead of start time
	if streamOptions.StartAtOperationTime != nil && streamOptions.StartAfter != nil {
		blog.Infof("resume token and time is both set, discard the resume time, option: %+v", streamOptions)
		streamOptions.StartAtOperationTime = nil
	}

	blog.InfoJSON("retry watch db %s with pipeline: %s, opts: %s, stream opts: %s", e.DBName, opts.pipeline,
		opts.Options, streamOptions)

	var err error
	opts.stream, err = e.client.
		Database(e.database).
		Watch(ctx, opts.pipeline, streamOptions)
	if err != nil {
		if isFatalError(err) {
			// TODO: send alarm immediately.
			blog.Errorf("mongodb watch db: %s got a fatal error, skip resume token and retry, err: %v",
				e.DBName, err)
			// reset the resume token, because we can not use the former resume token to watch success for now.
			streamOptions.StartAfter = nil
			opts.StartAfterToken = nil
			// because we have already got a fatal error, we can not try to watch from where we lost.
			// so re-watch from 1 minutes ago to avoid lost events.
			// Note: apparently, we may got duplicate events with this re-watch
			startAtTime := uint32(time.Now().Unix()) - 60
			streamOptions.StartAtOperationTime = &primitive.Timestamp{
				T: startAtTime,
				I: 0,
			}
			opts.StartAtTime = &types.TimeStamp{Sec: startAtTime}
			opts.currentToken.Data = ""

			if opts.WatchFatalErrorCallback != nil {
				err := opts.WatchFatalErrorCallback(types.TimeStamp{Sec: startAtTime})
				if err != nil {
					blog.Errorf("do watch fatal error callback for db %s failed, err: %v", e.DBName, err)
				}
			}
		}

		blog.ErrorJSON("mongodb watch db %s failed with opts: %s, pipeline: %s, streamOpts: %s, err: %s",
			e.DBName, opts, opts.pipeline, streamOptions, err)
		return opts, true
	}

	// re-watch success, now we clean start at operation time options
	streamOptions.StartAtOperationTime = nil
	return opts, false
}

func (e *Event) handleStreamEvent(ctx context.Context, opts *loopWatchOpts, collTasksMap map[string][]string) (
	*loopWatchOpts, bool) {

	event := new(types.RawEvent)
	if err := opts.stream.Decode(event); err != nil {
		blog.Errorf("watch db %s, but decode to raw event struct failed, err: %v", e.DBName, err)
		return opts, true
	}

	// if we received a invalid event, which is caused by collection drop, rename or drop database operation,
	// we have to try re-watch again. otherwise, this may cause this process CPU high because of continue
	// for loop cursor.
	// https://docs.mongodb.com/manual/reference/change-events/#invalidate-event
	if event.EventStream.OperationType == types.Invalidate {
		blog.ErrorJSON("mongodb watch received a invalid event, will retry watch again, options: %s", *opts)

		// clean the last resume token to force the next try watch from the beginning. otherwise we will
		// receive the invalid event again.
		opts.streamOptions.StartAfter = nil
		opts.StartAfterToken = nil
		// cause we have already got a fatal error, we can not try to watch from where we lost.
		// so re-watch from 1 minutes ago to avoid lost events.
		// Note: apparently, we may got duplicate events with this re-watch
		startAtTime := uint32(time.Now().Unix()) - 60
		opts.streamOptions.StartAtOperationTime = &primitive.Timestamp{
			T: startAtTime,
			I: 0,
		}
		opts.StartAtTime = &types.TimeStamp{Sec: startAtTime}
		opts.currentToken.Data = ""

		opts.stream.Close(ctx)
		return opts, true
	}

	opts.currentToken.Data = event.EventStream.Token.Data

	e.parseEvent(event, opts.eventChan, opts.collOptsInfo, collTasksMap)

	return opts, false
}

func (e *Event) parseEvent(event *types.RawEvent, eventChan chan *types.Event, collOptsInfo *parsedCollOptsInfo,
	collTasksMap map[string][]string) {

	base := event.EventStream

	collInfo, err := parseCollInfo(base.Namespace.Collection)
	if err != nil {
		blog.Errorf("parse event(%+v) collection info failed, err: %v", base, err)
		return
	}

	// get the event task ids matching the collection name
	taskIDs, exists := collTasksMap[base.Namespace.Collection]
	if !exists {
		for collRegex, regex := range collOptsInfo.collRegexMap {
			if regex.MatchString(base.Namespace.Collection) {
				taskIDs = append(taskIDs, collOptsInfo.collRegexTasksMap[collRegex]...)
			}
		}
		collTasksMap[base.Namespace.Collection] = taskIDs
	}

	if len(taskIDs) == 0 {
		blog.Errorf("watch db %s, but get invalid event not matching any task, base: %+v", e.DBName, base)
		return
	}

	// decode the event data to the event data struct, use pre data for delete event
	rawDoc := event.FullDoc
	if base.OperationType == types.Delete || event.FullDoc == nil {
		rawDoc = event.PreFullDoc
	}

	if rawDoc == nil {
		blog.Errorf("watch db %s, but get invalid event with no detail, base: %+v", e.DBName, base)
		return
	}

	var wg sync.WaitGroup
	for _, taskID := range taskIDs {
		wg.Add(1)
		go func(taskID string) {
			defer wg.Done()

			parsed, isValid := parseDataForTask(rawDoc, taskID, collOptsInfo.taskFilterMap, collOptsInfo.taskTypeMap)
			if !isValid {
				return
			}

			parsed.Oid = base.DocumentKey.ID.Hex()
			parsed.OperationType = base.OperationType
			parsed.CollectionInfo = collInfo
			parsed.ClusterTime = types.TimeStamp{
				Sec:  base.ClusterTime.T,
				Nano: base.ClusterTime.I,
			}
			parsed.Token = base.Token
			parsed.ChangeDesc = &types.ChangeDescription{
				UpdatedFields: base.UpdateDesc.UpdatedFields,
				RemovedFields: base.UpdateDesc.RemovedFields,
			}

			eventChan <- parsed
		}(taskID)
	}
	wg.Wait()
}

// isFatalError if watch encountered a fatal error, we should watch without resume token, which means from now.
// errors like:
// https://jira.mongodb.org/browse/SERVER-44610
// https://jira.mongodb.org/browse/SERVER-44733
func isFatalError(err error) bool {
	if strings.Contains(err.Error(), "ChangeStreamFatalError") {
		return true
	}

	if strings.Contains(err.Error(), "the resume point may no longer be in the oplog") {
		return true
	}

	if strings.Contains(err.Error(), "the resume token was not found") {
		return true
	}

	if strings.Contains(err.Error(), "CappedPositionLost") {
		return true
	}

	return false
}
