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
	"sync"
	"time"

	"configcenter/pkg/filter"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/storage/stream/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// List is a wrapper to list all the data with a collection and filter.
// when the list is done, it will send a event with a operation type is types.ListDone.
// when an error occurred, the returned ch will be closed.
func (e *Event) List(ctx context.Context, opts *types.ListOptions) (ch chan *types.Event, err error) {
	if err := opts.CheckSetDefault(); err != nil {
		return nil, err
	}

	// list collections
	collections, err := e.client.Database(e.database).ListCollectionNames(ctx, bson.M{})
	if err != nil {
		blog.Errorf("list db: %s collections failed, err :%v", e.database, err)
		return nil, err
	}

	collOpts := make(map[string]types.WatchCollOptions)
	for id, collOpt := range opts.CollOpts {
		collOpts[id] = types.WatchCollOptions{CollectionOptions: collOpt}
	}

	listOpts := &listOptions{
		collections:  collections,
		collOptsInfo: parseCollOpts(collOpts),
		pageSize:     opts.PageSize,
	}

	eventChan := make(chan *types.Event, types.DefaultEventChanSize)

	go func() {
		e.lister(ctx, opts.WithRetry, listOpts, eventChan)
	}()

	return eventChan, nil

}

type listOptions struct {
	collections  []string
	collOptsInfo *parsedCollOptsInfo
	pageSize     *int
}

// lister try to list data with filter. withRetry controls whether you need to retry list when an error encountered.
func (e *Event) lister(ctx context.Context, withRetry bool, opts *listOptions, ch chan *types.Event) {
	reset := func() {
		// sleep a while and retry later
		time.Sleep(3 * time.Second)
	}

	var wg sync.WaitGroup
	needReturn := false
	pipeline := make(chan struct{}, 10)
	for _, collection := range opts.collections {
		taskIDs, findOpts, cond, needSkip, err := e.parseCollListOpts(collection, opts)
		if err != nil {
			if withRetry {
				continue
			}
			return
		}

		if needSkip {
			continue
		}

		// list data from this collection
		pipeline <- struct{}{}
		wg.Add(1)
		listOpt := &listOneCollOptions{collection: collection, taskIDs: taskIDs, filter: cond,
			findOpts: findOpts, ch: ch, withRetry: withRetry, reset: reset}
		go func(listOpt *listOneCollOptions) {
			defer func() {
				wg.Done()
				<-pipeline
			}()

			collNeedReturn := e.listOneColl(ctx, listOpt)
			if collNeedReturn {
				needReturn = true
			}
		}(listOpt)

		if needReturn {
			return
		}
	}

	wg.Wait()

	// tell the user that the list operation has already done.
	// we only send for once.
	ch <- &types.Event{
		OperationType: types.ListDone,
	}
}

// get collection related task ids, find options and filters
func (e *Event) parseCollListOpts(collection string, opts *listOptions) ([]string, *options.FindOptions, mapstr.MapStr,
	bool, error) {
	taskIDs, fields, filters := make([]string, 0), make([]string, 0), make([]filter.RuleFactory, 0)
	needAllFilter, needAllFields := false, false
	for collRegex, regex := range opts.collOptsInfo.collRegexMap {
		if !regex.MatchString(collection) {
			continue
		}

		taskIDs = append(taskIDs, opts.collOptsInfo.collRegexTasksMap[collRegex]...)
		if opts.collOptsInfo.collCondMap[collRegex] == nil {
			needAllFilter = true
		} else if !needAllFilter {
			filters = append(filters, opts.collOptsInfo.collCondMap[collRegex])
		}
		if len(opts.collOptsInfo.collFieldsMap[collRegex]) == 0 {
			needAllFields = true
		} else if !needAllFields {
			fields = append(fields, opts.collOptsInfo.collFieldsMap[collRegex]...)
		}
	}

	if len(taskIDs) == 0 {
		return nil, nil, nil, true, nil
	}

	findOpts := new(options.FindOptions)
	findOpts.SetLimit(int64(*opts.pageSize))
	if !needAllFields && len(fields) != 0 {
		projection := make(map[string]int)
		for _, field := range fields {
			if len(field) <= 0 {
				continue
			}
			projection[field] = 1
		}
		projection["_id"] = 1
		findOpts.Projection = projection
	}

	cond := make(mapstr.MapStr)
	if !needAllFilter && len(filters) != 0 {
		expr := filter.Expression{RuleFactory: &filter.CombinedRule{Condition: filter.Or, Rules: filters}}
		var err error
		cond, err = expr.ToMgo()
		if err != nil {
			return nil, nil, nil, false, err
		}
	}

	return taskIDs, findOpts, cond, false, nil
}

type listOneCollOptions struct {
	collection    string
	taskIDs       []string
	filter        mapstr.MapStr
	findOpts      *options.FindOptions
	taskTypeMap   map[string]reflect.Type
	taskFilterMap map[string]*filter.Expression
	ch            chan *types.Event
	withRetry     bool
	reset         func()
}

type mongoID struct {
	Oid primitive.ObjectID `bson:"_id"`
}

// listOneColl try to list data with filter from one collection, returns if list operation needs to exit
func (e *Event) listOneColl(ctx context.Context, opts *listOneCollOptions) bool {
	collInfo, err := parseCollInfo(opts.collection)
	if err != nil {
		blog.Errorf("parse collection info for list operation failed, opt: %+v, err: %v", *opts, err)
		return false
	}

	for {
	retry:
		cursor, err := e.client.Database(e.database).
			Collection(opts.collection).
			Find(ctx, opts.filter, opts.findOpts)
		if err != nil {
			blog.Errorf("list db: %s, coll: %s failed, will *retry later*, err: %v", e.database, opts.collection, err)
			opts.reset()
			continue
		}

		hasData := false
		for cursor.Next(ctx) {
			hasData = true
			select {
			case <-ctx.Done():
				blog.Errorf("received stopped lister signal, stop list db: %s, collection: %s, err: %v", e.database,
					opts.collection, ctx.Err())
				return true
			default:
			}

			rawDoc := bson.Raw{}
			if err := cursor.Decode(&rawDoc); err != nil {
				blog.Errorf("list db: %s, coll: %s with cursor failed, err: %v", e.database, opts.collection, err)
				cursor.Close(ctx)
				if !opts.withRetry {
					blog.Warnf("list db: %s, coll: %s failed, will exit list immediately", e.database, opts.collection)
					close(opts.ch)
					return true
				}

				opts.reset()
				goto retry
			}

			oidInfo := new(mongoID)
			if err := bson.Unmarshal(rawDoc, &oidInfo); err != nil {
				blog.Errorf("decode mongodb oid failed, err: %v, data: %s", err, rawDoc)
				continue
			}
			opts.filter["_id"] = mapstr.MapStr{common.BKDBGT: oidInfo.Oid}

			for _, taskID := range opts.taskIDs {
				parsed, isValid := parseDataForTask(rawDoc, taskID, opts.taskFilterMap, opts.taskTypeMap)
				if !isValid {
					continue
				}

				parsed.Oid = oidInfo.Oid.Hex()
				parsed.OperationType = types.Lister
				parsed.CollectionInfo = collInfo
				opts.ch <- parsed
			}
		}

		if err := cursor.Err(); err != nil {
			blog.Errorf("list db: %s, coll: %s with cursor failed, err: %v", e.database, opts.collection, err)
			cursor.Close(ctx)
			if !opts.withRetry {
				blog.Warnf("list db: %s, coll: %s failed, will exit list immediately", e.database, opts.collection)
				close(opts.ch)
				return true
			}
			opts.reset()
			goto retry
		}
		cursor.Close(ctx)

		if !hasData {
			return false
		}
	}
}
