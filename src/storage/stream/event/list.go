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
	"fmt"
	"reflect"
	"time"

	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	"configcenter/src/storage/stream/types"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// List is a wrapper to list all the data with a collection and filter.
// when the list is done, it will send a event with a operation type is types.ListDone.
// when an error occurred, the returned ch will be closed.
func (e *Event) List(ctx context.Context, opts *types.ListOptions) (ch chan *types.Event, err error) {
	if err := opts.CheckSetDefault(); err != nil {
		return nil, err
	}
	// prepare for list all the data.
	totalCnt, err := e.client.Database(e.database).
		Collection(opts.Collection).
		CountDocuments(ctx, opts.Filter)
	if err != nil {
		return nil, fmt.Errorf("count db %s, collection: %s with filter: %+v failed, err: %v",
			e.database, opts.Collection, opts.Filter, err)
	}

	eventChan := make(chan *types.Event, types.DefaultEventChanSize)

	go func() {
		e.lister(ctx, false, totalCnt, opts, eventChan)
	}()

	return eventChan, nil

}

// lister is try to list data with filter. withRetry is to control whether you need to retry list when an error encountered.
func (e *Event) lister(ctx context.Context, withRetry bool, cnt int64, opts *types.ListOptions, ch chan *types.Event) {

	pageSize := *opts.PageSize
	reset := func() {
		// sleep a while and retry later
		time.Sleep(3 * time.Second)
	}
	for start := 0; start < int(cnt); start += pageSize {

		findOpts := new(options.FindOptions)
		findOpts.SetSkip(int64(start))
		findOpts.SetLimit(int64(pageSize))
		projection := make(map[string]int)
		if len(opts.Fields) != 0 {
			for _, field := range opts.Fields {
				if len(field) <= 0 {
					continue
				}
				projection[field] = 1
			}
			findOpts.Projection = projection
		}

	retry:
		cursor, err := e.client.Database(e.database).
			Collection(opts.Collection).
			Find(ctx, opts.Filter, findOpts)
		if err != nil {
			blog.Errorf("list watch operation, but list db: %s, collection: %s failed, will *retry later*, err: %v",
				e.database, opts.Collection, err)
			reset()
			continue
		}

		for cursor.Next(ctx) {
			select {
			case <-ctx.Done():
				blog.Errorf("received stopped lister signal, stop list db: %s, collection: %s, err: %v", e.database,
					opts.Collection, ctx.Err())
				return
			default:

			}

			// create a new event struct for use
			result := reflect.New(reflect.TypeOf(opts.EventStruct)).Elem()
			err := cursor.Decode(result.Addr().Interface())
			if err != nil {
				blog.Errorf("list watch operation, but list db: %s, collection: %s with cursor failed, will *retry later*, err: %v",
					e.database, opts.Collection, err)

				cursor.Close(ctx)
				if !withRetry {
					blog.Warnf("list watch operation, but list db: %s, collection: %s with cursor failed, will exit list immediately.",
						e.database, opts.Collection)
					close(ch)
					return
				}

				reset()
				goto retry
			}

			byt, _ := json.Marshal(result.Addr().Interface())
			oid := gjson.GetBytes(byt, "_id").String()

			// send the event now
			ch <- &types.Event{
				Oid:           oid,
				Document:      result.Interface(),
				OperationType: types.Lister,
				DocBytes:      byt,
			}
		}

		if err := cursor.Err(); err != nil {
			blog.Errorf("list watch operation, but list db: %s, collection: %s with cursor failed, will *retry later*, err: %v",
				e.database, opts.Collection, err)
			cursor.Close(ctx)
			if !withRetry {
				blog.Warnf("list watch operation, but list db: %s, collection: %s with cursor failed, will exit list immediately.",
					e.database, opts.Collection)
				close(ch)
				return
			}
			reset()
			goto retry
		}
		cursor.Close(ctx)
	}

	// tell the user that the list operation has already done.
	// we only send for once.
	ch <- &types.Event{
		Oid:           "",
		Document:      reflect.New(reflect.TypeOf(opts.EventStruct)).Elem().Interface(),
		OperationType: types.ListDone,
	}

}
