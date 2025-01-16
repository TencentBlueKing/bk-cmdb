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
	"time"

	"configcenter/src/common/blog"
	"configcenter/src/storage/stream/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ListWatch list all data and watch change stream events
func (e *Event) ListWatch(ctx context.Context, opts *types.ListWatchOptions) (*types.Watcher, error) {
	if err := opts.CheckSetDefault(); err != nil {
		return nil, err
	}

	// list collections
	collections, err := e.client.Database(e.database).ListCollectionNames(ctx, bson.M{})
	if err != nil {
		blog.Errorf("list db: %s collections failed, err :%v", e.database, err)
		return nil, err
	}

	eventChan := make(chan *types.Event, types.DefaultEventChanSize)

	go func() {
		pipeline, streamOptions, collOptsInfo := generateOptions(&opts.Options)

		// TODO: should use the mongodb cluster timestamp, if the time is not synchronise with
		// mongodb cluster time, then we may have to lost some events.
		// A better way is to get the mongodb cluster timestamp and set it
		// as the "startTime".
		startAt := time.Now()
		streamOptions.StartAtOperationTime = &primitive.Timestamp{
			// normally, a unix time seconds is a int64 value,
			// but mongodb has a 32 bit T represent a unix seconds time.
			// calculate this: time.Duration(math.MaxInt32 - int32(time.Now().Unix()))/(time.Hour * 24 *365/time.Second))
			// the value is 17 years, it's okay for now.
			// reference: https://docs.mongodb.com/manual/reference/bson-types/#timestamps
			T: uint32(startAt.Unix()),
			I: 1,
		}

		// we watch the stream at first, so that we can know if we can watch success.
		// and, we do not read the event stream immediately, we wait until all the data
		// has been listed from database.
		stream, streamOptions, watchOpts, err := e.watch(ctx, pipeline, streamOptions, &opts.Options)
		if err != nil {
			blog.Fatalf("mongodb watch failed with conf: %+v, err: %v", *opts, err)
		}

		listOpts := &listOptions{
			collections:  collections,
			collOptsInfo: collOptsInfo,
			pageSize:     opts.PageSize,
		}

		go func() {
			// list all the data from the collection and send it as an event now.
			e.lister(ctx, true, listOpts, eventChan)

			select {
			case <-ctx.Done():
				blog.Errorf("received stopped watch signal, stop list db: %s, name: %s, err: %v", e.database, e.DBName,
					ctx.Err())
				return
			default:

			}

			// all the data has already listed and send the event.
			// now, it's time to watch the event stream.
			loopOpts := &loopWatchOpts{
				Options:       watchOpts,
				streamOptions: streamOptions,
				stream:        stream,
				pipeline:      pipeline,
				eventChan:     eventChan,
				collOptsInfo:  collOptsInfo,
			}
			e.loopWatch(ctx, loopOpts)
		}()
	}()

	watcher := &types.Watcher{
		EventChan: eventChan,
	}
	return watcher, nil

}
