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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (e *Event) Watch(ctx context.Context, opts *types.WatchOptions) (*types.Watcher, error) {
	if err := opts.CheckSetDefault(); err != nil {
		return nil, err
	}
	pipeline, streamOptions := generateOptions(&opts.Options)

	var stream *mongo.ChangeStream
	var err error
	stream, err = e.client.
		Database(e.database).
		Collection(opts.Collection).
		Watch(ctx, pipeline, streamOptions)

	if err != nil {
		blog.Errorf("mongodb watch failed with conf: %+v, err: %v", *opts, err)
		return nil, fmt.Errorf("watch collection: %s failed, err: %v", opts.Collection, err)
	}
	eventChan := make(chan *types.Event, types.DefaultEventChanSize)
	go e.loopWatch(ctx, &opts.Options, streamOptions, stream, pipeline, eventChan)

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

	for {
		select {
		case <-ctx.Done():
			blog.Warnf("received stopped loop watch signal, stop watch db: %s, collection: %s, err: %v", e.database,
				opts.Collection, ctx.Err())
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
			var err error
			stream, err = e.client.
				Database(e.database).
				Collection(opts.Collection).
				Watch(ctx, pipeline, streamOptions)
			if err != nil {
				blog.Warnf("mongodb watch collection: %s failed with conf: %v, err: %v", opts.Collection, *opts, err)
				retry = true
				continue
			}
		}

		for stream.Next(ctx) {
			newStruct := newEventStruct(typ)
			if err := stream.Decode(newStruct.Addr().Interface()); err != nil {
				blog.Errorf("watch collection %s, but decode to event struct: %s failed, err: %v", opts.Collection, reflect.TypeOf(opts.EventStruct).Name(), err)
				continue
			}

			base := newStruct.Field(0).Interface().(types.EventStream)
			currentToken.Data = base.Token.Data
			byt, _ := json.Marshal(newStruct.Field(1).Addr().Interface())

			eventChan <- &types.Event{
				Oid:           base.DocumentKey.ID.Hex(),
				OperationType: base.OperationType,
				Document:      newStruct.Field(1).Addr().Interface(),
				DocBytes:      byt,
				ClusterTime: types.TimeStamp{
					Sec:  base.ClusterTime.T,
					Nano: base.ClusterTime.I,
				},
				Token: base.Token,
			}
		}

		if err := stream.Err(); err != nil {
			blog.Errorf("mongodb watch encountered a error, conf: %v, err: %v", *opts, err)
			stream.Close(ctx)
			retry = true
			continue
		}
	}
}
