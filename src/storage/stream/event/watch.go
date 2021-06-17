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
	"reflect"
	"strings"
	"time"

	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	"configcenter/src/storage/stream/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (e *Event) Watch(ctx context.Context, opts *types.WatchOptions) (*types.Watcher, error) {
	if err := opts.CheckSetDefault(); err != nil {
		return nil, err
	}

	eventChan := make(chan *types.Event, types.DefaultEventChanSize)
	go func() {
		pipeline, streamOptions := generateOptions(&opts.Options)

		blog.InfoJSON("start watch with pipeline: %s, options: %s, stream options: %s", pipeline, opts, streamOptions)

		var stream *mongo.ChangeStream
		var err error

		if opts.Collection != "" {
			stream, err = e.client.
				Database(e.database).
				Collection(opts.Collection).
				Watch(ctx, pipeline, streamOptions)
		} else {
			stream, err = e.client.
				Database(e.database).
				Watch(ctx, pipeline, streamOptions)
		}

		if err != nil && isFatalError(err) {
			// TODO: send alarm immediately.
			blog.Errorf("mongodb watch collection: %s got a fatal error, skip resume token and retry, err: %v",
				opts.Collection, err)
			// reset the resume token, because we can not use the former resume token to watch success for now.
			streamOptions.StartAfter = nil
			// cause we have already got a fatal error, we can not try to watch from where we lost.
			// so re-watch from 1 minutes ago to avoid lost events.
			// Note: apparently, we may got duplicate events with this re-watch
			startAtTime := uint32(time.Now().Unix()) - 60
			streamOptions.StartAtOperationTime = &primitive.Timestamp{
				T: startAtTime,
				I: 0,
			}

			if opts.WatchFatalErrorCallback != nil {
				err := opts.WatchFatalErrorCallback(types.TimeStamp{Sec: startAtTime})
				if err != nil {
					blog.Errorf("do watch fatal error callback for coll %s failed, err: %v", opts.Collection, err)
				}
			}

			blog.InfoJSON("start watch with pipeline: %s, options: %s, stream options: %s", pipeline, opts, streamOptions)
			if opts.Collection != "" {
				stream, err = e.client.
					Database(e.database).
					Collection(opts.Collection).
					Watch(ctx, pipeline, streamOptions)
			} else {
				stream, err = e.client.
					Database(e.database).
					Watch(ctx, pipeline, streamOptions)
			}
		}

		if err != nil {
			if err == context.Canceled {
				// if error is context cancelled, then loop watch will exit at the same time
				return
			}
			blog.Fatalf("mongodb watch failed with conf: %+v, err: %v", *opts, err)
		}

		go e.loopWatch(ctx, &opts.Options, streamOptions, stream, pipeline, eventChan)

	}()

	watcher := &types.Watcher{
		EventChan: eventChan,
	}
	return watcher, nil
}

func (e *Event) loopWatch(ctx context.Context,
	opts *types.Options,
	streamOptions *options.ChangeStreamOptions,
	stream *mongo.ChangeStream,
	pipeline mongo.Pipeline,
	eventChan chan *types.Event) {

	retry := false
	currentToken := types.EventToken{Data: ""}
	typ := reflect.Indirect(reflect.ValueOf(opts.EventStruct)).Type()

	e.setCleaner(ctx, eventChan, opts.Collection)

	for {
		// no events, try cancel watch here.
		select {
		case <-ctx.Done():
			blog.Warnf("received stopped loop watch signal, stop watch db: %s, collection: %s, err: %v", e.database,
				opts.Collection, ctx.Err())

			if stream != nil {
				stream.Close(context.Background())
			}

			return
		default:

		}

		if retry {
			time.Sleep(5 * time.Second)
			if len(currentToken.Data) != 0 {
				// if error occurs, then retry watch and start from the last token.
				// so that we can continue the event from where it just broken.
				streamOptions.SetStartAfter(currentToken)
			}

			blog.InfoJSON("retry watch with pipeline: %s, options: %s, stream options: %s", pipeline, opts, streamOptions)

			var err error
			if opts.Collection != "" {
				stream, err = e.client.
					Database(e.database).
					Collection(opts.Collection).
					Watch(ctx, pipeline, streamOptions)
			} else {
				stream, err = e.client.
					Database(e.database).
					Watch(ctx, pipeline, streamOptions)
			}
			if err != nil {
				if isFatalError(err) {
					// TODO: send alarm immediately.
					blog.Errorf("mongodb watch collection: %s got a fatal error, skip resume token and retry, err: %v",
						opts.Collection, err)
					// reset the resume token, because we can not use the former resume token to watch success for now.
					streamOptions.StartAfter = nil
					// cause we have already got a fatal error, we can not try to watch from where we lost.
					// so re-watch from 1 minutes ago to avoid lost events.
					// Note: apparently, we may got duplicate events with this re-watch
					startAtTime := uint32(time.Now().Unix()) - 60
					streamOptions.StartAtOperationTime = &primitive.Timestamp{
						T: startAtTime,
						I: 0,
					}
					currentToken.Data = ""

					if opts.WatchFatalErrorCallback != nil {
						err := opts.WatchFatalErrorCallback(types.TimeStamp{Sec: startAtTime})
						if err != nil {
							blog.Errorf("do watch fatal error callback for coll %s failed, err: %v", opts.Collection, err)
						}
					}
				}

				blog.ErrorJSON("mongodb watch %s failed with opts: %s, pipeline: %s, streamOpts: %s, err: %s",
					opts.Collection, opts, pipeline, streamOptions, err)

				retry = true
				continue
			}

			// re-watch success, now we clean start at operation time options
			streamOptions.StartAtOperationTime = nil
		}

		for stream.Next(ctx) {
			// still have events, try cancel steam here.
			select {
			case <-ctx.Done():
				blog.Warnf("received stopped loop watch signal, stop loop next, watch db: %s, collection: %s, err: %v",
					e.database, opts.Collection, ctx.Err())
				stream.Close(context.Background())
				return
			default:

			}

			newStruct := newEventStruct(typ)
			if err := stream.Decode(newStruct.Addr().Interface()); err != nil {
				blog.Errorf("watch collection %s, but decode to event struct: %v failed, err: %v",
					opts.Collection, reflect.TypeOf(opts.EventStruct), err)
				continue
			}

			base := newStruct.Field(0).Interface().(types.EventStream)

			// if we received a invalid event, which is caused by collection drop, rename or drop database operation,
			// we have to try re-watch again. otherwise, this may cause this process CPU high because of continue
			// for loop cursor.
			// https://docs.mongodb.com/manual/reference/change-events/#invalidate-event
			if base.OperationType == types.Invalidate {
				blog.ErrorJSON("mongodb watch received a invalid event, will retry watch again, options: %s", *opts)

				// clean the last resume token to force the next try watch from the beginning. otherwise we will
				// receive the invalid event again.
				streamOptions.StartAfter = nil
				// cause we have already got a fatal error, we can not try to watch from where we lost.
				// so re-watch from 1 minutes ago to avoid lost events.
				// Note: apparently, we may got duplicate events with this re-watch
				streamOptions.StartAtOperationTime = &primitive.Timestamp{
					T: uint32(time.Now().Unix()) - 60,
					I: 0,
				}
				currentToken.Data = ""

				stream.Close(ctx)
				retry = true
				break
			}

			currentToken.Data = base.Token.Data
			byt, _ := json.Marshal(newStruct.Field(1).Addr().Interface())

			eventChan <- &types.Event{
				Oid:           base.DocumentKey.ID.Hex(),
				OperationType: base.OperationType,
				Document:      newStruct.Field(1).Addr().Interface(),
				DocBytes:      byt,
				Collection:    base.Namespace.Collection,
				ClusterTime: types.TimeStamp{
					Sec:  base.ClusterTime.T,
					Nano: base.ClusterTime.I,
				},
				Token: base.Token,
				ChangeDesc: &types.ChangeDescription{
					UpdatedFields: base.UpdateDesc.UpdatedFields,
					RemovedFields: base.UpdateDesc.RemovedFields,
				},
			}
		}

		if err := stream.Err(); err != nil {
			blog.ErrorJSON("mongodb watch encountered a error, conf: %s, err: %s", *opts, err)
			stream.Close(ctx)
			retry = true
			continue
		}
	}
}

// setCleaner set up a monitor to close the cursor when the context is canceled.
// this is useful to release stream resource when this watch is canceled outside with context is canceled.
func (e *Event) setCleaner(ctx context.Context, eventChan chan *types.Event, coll string) {
	go func() {
		select {
		case <-ctx.Done():
			blog.Warnf("received stopped loop watch collection: %s signal, close cursor now, err: %v",
				coll, ctx.Err())

			// even though we may already close the stream, but there may still have events in the stream's
			// batch cursor, so we need to consume a event, so that we can release the stream resource
			select {
			// try consume a event, so that stream.Next(ctx) can be called to release the stream resources.
			case <-eventChan:
				blog.Warnf("received stopped loop watch collection: %s signal, consumed a event", coll)

			default:
				// no events, and stream resource will be recycled in the next round.
			}

			return
		}
	}()
}

// if watch encountered a fatal error, we should watch without resume token, which means from now.
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

	return false
}
