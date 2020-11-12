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
	"fmt"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/watch"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/driver/redis"
	"configcenter/src/storage/stream/types"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func newFlow(ctx context.Context, opts FlowOptions) error {
	flow := Flow{
		FlowOptions: opts,
		retrySignal: make(chan struct{}),
		metrics:     initialMetrics(opts.Collection),
	}
	flow.cleanExpiredEvents()
	flow.cleanDelArchiveData()

	return flow.RunFlow(ctx)
}

type Flow struct {
	FlowOptions
	retrySignal chan struct{}
	metrics     *eventMetrics
}

func (f *Flow) RunFlow(ctx context.Context) error {
	blog.Infof("start run flow for key: %s.", f.key.Namespace())

	startToken, err := f.getStartToken()
	if err != nil {
		blog.Errorf("run flow, but get start token failed, err: %v", err)
		return err
	}

	blog.Infof("start run flow for key: %s with token: %s.", f.key.Namespace(), startToken)
	event := make(map[string]interface{})
	watchOpts := &types.WatchOptions{
		Options: types.Options{
			EventStruct:     &event,
			Collection:      f.Collection,
			StartAfterToken: nil,
		},
	}
	if f.Collection == common.BKTableNameBaseHost {
		watchOpts.EventStruct = &metadata.HostMapStr{}
	}

	if len(startToken) != 0 {
		blog.Warnf("run flow, we got a start token: %s, and we start watch from here.", startToken)
		watchOpts.StartAfterToken = &types.EventToken{
			Data: startToken,
		}
	}

	return f.tryLoopFlow(ctx, watchOpts)
}

func (f *Flow) tryLoopFlow(ctx context.Context, opts *types.WatchOptions) error {
	var cancel func()
	var cancelCtx context.Context
	cancelCtx, cancel = context.WithCancel(ctx)
	if err := f.loopFlow(cancelCtx, opts); err != nil {
		cancel()
		return err
	}

	go func() {
		for {
			select {
			// wait for another retry
			case <-f.retrySignal:
				// collect retry error metrics
				f.metrics.collectRetryError()

				// wait for a second and then do the retry work.
				time.Sleep(1 * time.Second)

				// initialize a new retry signal immediately for next usage.
				f.retrySignal = make(chan struct{})

				// cancel the former watch
				cancel()

				// use the last token to resume so that we can start again from where we stopped.
				lastToken, err := f.getStartToken()
				if err != nil {
					blog.Errorf("run flow, but get last event token failed, err: %v", err)
					// notify retry signal
					close(f.retrySignal)
					continue
				}
				blog.Warnf("the former loop flow: %s failed, start retry again from token: %s.", f.Collection, lastToken)

				// set start after token if needed.
				if lastToken != "" {
					// we have already received the new event and handle it success,
					// so we need to use this token. otherwise, we should still use the opts.StartAfterToken
					opts.StartAfterToken = &types.EventToken{Data: lastToken}
				}

				cancelCtx, cancel = context.WithCancel(ctx)
				if err := f.loopFlow(cancelCtx, opts); err != nil {
					cancel()
					// notify retry signal
					close(f.retrySignal)

					blog.Errorf("loop flow %s failed, err: %v", f.Collection, err)
					continue
				}

				blog.Warnf("retry loop flow: %s from token: %s success.", f.Collection, lastToken)
			}
		}

	}()

	return nil
}

func (f *Flow) loopFlow(ctx context.Context, opts *types.WatchOptions) error {

	watcher, err := f.watch.Watch(ctx, opts)
	if err != nil {
		blog.Errorf("run flow, but watch failed, err: %v", err)
		return err
	}

	go func() {
		for e := range watcher.EventChan {

			select {
			case <-ctx.Done():
				blog.Warnf("received stop watch flow collection: %s, err: %v", opts.Collection, ctx.Err())
				return
			default:

			}

			if !f.isMaster.IsMaster() {
				blog.V(4).Infof("run flow, received collection %s event, type: %s, oid: %s, but not master, skip.", f.Collection,
					e.OperationType, e.Oid)
				continue
			}

			blog.V(4).Infof("run flow, received collection %s event, type: %s, oid :%s", f.Collection, e.OperationType, e.Oid)
			if blog.V(5) {
				blog.Infof("event doc details: %s, oid: %s", e.Oid)
			}

			// collect event's basic metrics
			f.metrics.collectBasic(e)

			switch e.OperationType {
			case types.Insert, types.Update, types.Replace:
				if retry, err := f.doUpsert(e); err != nil {
					blog.Errorf("run flow, but do upsert failed, operation type: %s, doc: %s, err: %v, oid: %s", e.OperationType,
						e.DocBytes, err, e.Oid)

					if !retry {
						// an error and do not need to retry do update.
						continue
					}
					// an error occurred, we need to retry do upsert later.
					// tell the schedule to re-watch again.
					close(f.retrySignal)
					// exist this goroutine.
					return
				}

			case types.Delete:
				if retry, err := f.doDelete(e); err != nil {
					blog.Errorf("run flow, but do delete failed, doc: %s, err: %v, oid: %s", e.DocBytes, err, e.Oid)

					if !retry {
						// an error and do not need to retry do update.
						continue
					}

					// tell the schedule to re-watch again.
					close(f.retrySignal)
					// exist this goroutine.
					return
				}

			case types.Invalidate:
				blog.Errorf("loop flow, received invalid event operation type, doc: %s", e.DocBytes)
				continue
			default:
				blog.Errorf("loop flow, received unsupported event operation type, doc: %s", e.DocBytes)
				continue
			}
		}
	}()

	return nil
}

// if an error occurred, then the caller should check the retry is true or not, if true, then this event
// need to be retry later.
func (f *Flow) doUpsert(e *types.Event) (retry bool, err error) {
	start := time.Now()
	retry, err = f.do(e)
	if err != nil {
		return retry, err
	}

	f.metrics.collectCycleDuration(time.Since(start))

	return
}

func (f *Flow) do(e *types.Event) (retry bool, err error) {
	blog.Infof("run flow, received %s %s event, oid: %s", f.Collection, e.OperationType, e.Oid)
	blog.V(5).Infof("event doc detail: %s, oid: %s", e.DocBytes, e.Oid)

	// validate the event is valid or not.
	// the invalid event will be dropped.
	if err := f.key.Validate(e.DocBytes); err != nil {
		blog.Errorf("run flow, received %s event, but got invalid event, doc: %s err: %v, oid: %s", f.Collection, e.DocBytes, err, e.Oid)
		return false, err
	}

	keys, err := redis.Client().HMGet(context.Background(), f.key.MainHashKey(), f.key.HeadKey(), f.key.TailKey()).Result()
	if err != nil {
		f.metrics.collectRedisError()
		return true, err
	}

	var head, tail string
	var ok bool

	if keys[0] != nil {
		head, ok = keys[0].(string)
		if !ok {
			return false, fmt.Errorf("got invalid head value: %v", keys[0])
		}
	}

	if keys[1] != nil {
		tail, ok = keys[1].(string)
		if !ok {
			return false, fmt.Errorf("got invalid tail value: %v", keys[1])
		}
	}

	switch {
	case head == "" && tail == "":
		// this is ok, we initialize the head and tail.
		return f.initializeHeadTailNode(e)

	case head == "" && tail != "":
		// the event hashmap chain has broken, this should not happen.
		// we need to repair the hashmap now.
		// TODO: repair the hashmap chain.

	case head != "" && tail == "":
		// the event hashmap chain has broken, this should not happen.
		// we need to repair the hashmap now.
		// TODO: repair the hashmap chain.
	}

	// head and tail is all not empty, it's good.
	// now it's time to insert the event node.

	// get the previous node from the tail node's next cursor
	prevCursor := gjson.Get(tail, "next_cursor").String()
	if prevCursor == "" {
		blog.Errorf("get previous cursor from tail node, but next_cursor is empty, tail: %s, oid: %s", tail, e.Oid)
		return false, errors.New("next_cursor is empty in tail node")
	}

	// get previous node with previous cursor
	prev, err := redis.Client().HGet(context.Background(), f.key.MainHashKey(), prevCursor).Result()
	if err != nil {
		if redis.IsNilErr(err) {
			blog.Errorf("get previous cursor: %s node from redis failed, err: %v, oid: %s", prevCursor, err, e.Oid)
			return false, err
		}

		f.metrics.collectRedisError()
		blog.Errorf("get previous cursor: %s node from redis failed, err: %v, oid: %s", prevCursor, err, e.Oid)
		return true, err
	}

	prevNode := new(watch.ChainNode)
	if err := json.Unmarshal([]byte(prev), prevNode); err != nil {
		blog.Errorf("run flow, unmarshal previous node failed, node: %s, err: %v, oid: %s", prev, err, e.Oid)
		return false, err
	}

	// insert a node and update the tail node.
	return f.insertNewNode(prevCursor, prevNode, e)
}

func (f *Flow) initializeHeadTailNode(e *types.Event) (bool, error) {
	name := f.key.Name(e.DocBytes)

	currentCursor, err := watch.GetEventCursor(f.Collection, e)
	if err != nil {
		blog.Errorf("initialize head tail node, but get cursor failed, name: %s, err: %v, oid: %s", name, err, e.Oid)
		return false, err
	}

	newNode := &watch.ChainNode{
		ClusterTime: e.ClusterTime,
		Oid:         e.Oid,
		EventType:   watch.ConvertOperateType(e.OperationType),
		Token:       e.Token.Data,
		Cursor:      currentCursor,
		NextCursor:  f.key.TailKey(),
	}

	nByte, err := json.Marshal(newNode)
	if err != nil {
		blog.Errorf("run flow, marshal new node failed, node: %v, name: %s, err: %v oid: %s", *newNode, name, err, e.Oid)
		return false, err
	}

	headNode := &watch.ChainNode{
		Cursor:     f.key.HeadKey(),
		NextCursor: currentCursor,
	}
	hByte, err := json.Marshal(headNode)
	if err != nil {
		blog.Errorf("run flow, marshal head node failed, node: %v, name: %s, err: %v, oid: %s", *headNode, name, err, e.Oid)
		return false, err
	}

	// when tail node is all empty, then the head node is must next to the tail node
	tailNode := &watch.ChainNode{
		Token:      e.Token.Data,
		Cursor:     f.key.TailKey(),
		NextCursor: currentCursor,
	}

	tByte, err := json.Marshal(tailNode)
	if err != nil {
		blog.Errorf("run flow, marshal tail node failed, node: %v, name: %s, err: %v, oid: %s", *tailNode, name, err, e.Oid)
		return false, err
	}

	val := map[string]interface{}{
		currentCursor:   string(nByte),
		f.key.HeadKey(): string(hByte),
		f.key.TailKey(): string(tByte),
	}

	detail := types.EventDetail{
		Detail:        types.JsonString(e.DocBytes),
		UpdatedFields: e.ChangeDesc.UpdatedFields,
		RemovedFields: e.ChangeDesc.RemovedFields,
	}
	detailBytes, err := json.Marshal(detail)
	if err != nil {
		blog.Errorf("run flow, marshal detail failed, name: %s, detail: %+v, err: %v, oid: %s", name, detail, err, e.Oid)
		return false, err
	}

	getLock := f.getLockWithRetry(f.key.LockKey(), e.Oid)
	if !getLock {
		blog.Errorf("run flow, set head and tail key, name: %s, op: %s, do not get lock and return, oid: %s", name, e.OperationType, e.Oid)
		return true, errors.New("get lock failed")
	}

	// already get the lock. prepare to release the lock.
	releaseLock := func() {
		if err := redis.Client().Del(context.Background(), f.key.LockKey()).Err(); err != nil {
			blog.ErrorfDepthf(1, "run flow, set head and tail key, name: %s, op: %s, release lock failed, err: %v, oid: %s", name, e.OperationType, err, e.Oid)
		}
	}

	pipe := redis.Client().Pipeline()
	pipe.HMSet(f.key.MainHashKey(), val)
	pipe.Set(f.key.DetailKey(currentCursor), string(detailBytes), 0)
	if _, err := pipe.Exec(); err != nil {
		f.metrics.collectRedisError()
		releaseLock()
		blog.Errorf("run flow, set head and tail key failed, name: %s, err: %v, oid: %s", name, err, e.Oid)
		return true, err
	}
	releaseLock()
	blog.Infof("insert watch event for %s success, name: %s, cursor: %s, oid: %s", f.Collection, name, currentCursor, e.Oid)
	return false, nil
}

func (f *Flow) doDelete(e *types.Event) (retry bool, err error) {
	blog.Infof("received delete %s event, oid: %s", f.Collection, e.Oid)

	start := time.Now()
	defer func() {
		if err != nil {
			return
		}
		f.metrics.collectCycleDuration(time.Since(start))
	}()

	filter := mapstr.MapStr{
		"oid": e.Oid,
	}

	if f.Collection == common.BKTableNameBaseHost {
		doc := new(hostArchive)
		err = mongodb.Client().Table(common.BKTableNameDelArchive).Find(filter).One(context.Background(), doc)
		if err != nil {
			f.metrics.collectMongoError()
			blog.Errorf("received delete %s event, but get archive deleted doc from mongodb failed, oid: %s, err: %v",
				f.Collection, e.Oid, err)
			if strings.Contains(err.Error(), "document not found") {
				return false, err
			}
			return true, err
		}

		byt, err := json.Marshal(doc.Detail)
		if err != nil {
			blog.Errorf("received delete %s event, but marshal detail to bytes failed, oid: %s, err: %v", f.Collection, e.Oid, err)
			return false, err
		}
		e.DocBytes = byt

	} else {

		doc := bsonx.Doc{}
		err = mongodb.Client().Table(common.BKTableNameDelArchive).Find(filter).One(context.Background(), &doc)
		if err != nil {
			f.metrics.collectMongoError()

			blog.Errorf("received delete %s event, but get archive deleted doc from mongodb failed, oid: %s, err: %v",
				f.Collection, e.Oid, err)
			if strings.Contains(err.Error(), "document not found") {
				return false, err
			}
			return true, err
		}

		byt, err := bson.MarshalExtJSON(doc.Lookup("detail"), false, false)
		if err != nil {
			f.metrics.collectMongoError()

			blog.Errorf("received delete %s event, but marshal detail to bytes failed, oid: %s, err: %v", f.Collection, e.Oid, err)
			return false, err
		}
		e.DocBytes = byt
	}

	return f.do(e)
}

// getStartToken get the started token when the system is started.
// if this token is empty, then system need to watch from now on.
func (f *Flow) getStartToken() (token string, err error) {
	tail, err := redis.Client().HGet(context.Background(), f.key.MainHashKey(), f.key.TailKey()).Result()
	if err != nil {
		if !redis.IsNilErr(err) {
			f.metrics.collectRedisError()
			return "", err
		}
		// the tail key is not exist.
		return "", nil
	}

	tailNode := new(watch.ChainNode)
	if err := json.Unmarshal([]byte(tail), tailNode); err != nil {
		return "", err
	}

	return tailNode.Token, nil
}

// insertNewNode insert the new event to hashmap's node chain list.
// update the previous node to link the new node.
// set the tail node line to the new node.
// if an error occurred, then the caller should check the retry is true or not, if true, then this event
// need to be retry later.
func (f *Flow) insertNewNode(prevCursor string, prevNode *watch.ChainNode, e *types.Event) (retry bool, err error) {
	name := f.key.Name(e.DocBytes)

	currentCursor, err := watch.GetEventCursor(f.Collection, e)
	if err != nil {
		blog.Errorf("get event cursor failed, name: %s, err: %v, oid: %s ", name, err, e.Oid)
		return false, err
	}
	prevNode.NextCursor = currentCursor

	// create a new node
	newNode := &watch.ChainNode{
		ClusterTime: e.ClusterTime,
		Oid:         e.Oid,
		EventType:   watch.ConvertOperateType(e.OperationType),
		Token:       e.Token.Data,
		Cursor:      currentCursor,
		NextCursor:  f.key.TailKey(),
	}

	nBytes, err := json.Marshal(newNode)
	if err != nil {
		blog.Errorf("run flow, marshal new node failed, name: %s, node: %+v err: %v, oid: %s", name, *newNode, err, e.Oid)
		return false, err
	}

	// new tail node
	tailNode := &watch.ChainNode{
		Token:      e.Token.Data,
		Cursor:     f.key.TailKey(),
		NextCursor: currentCursor,
	}
	tBytes, err := json.Marshal(tailNode)
	if err != nil {
		blog.Errorf("run flow, marshal tail node failed, name: %s, node: %+v, err: %v, oid: %s", name, *tailNode, err, e.Oid)
		return false, err
	}

	pBytes, err := json.Marshal(prevNode)
	if err != nil {
		blog.Errorf("run flow, marshal previous node failed, name: %s, node： %+v err: %v, oid: %s", name, *prevNode, err, e.Oid)
		return false, err
	}

	values := map[string]interface{}{
		prevCursor:      string(pBytes),
		currentCursor:   string(nBytes),
		f.key.TailKey(): string(tBytes),
	}

	detail := types.EventDetail{
		Detail:        types.JsonString(e.DocBytes),
		UpdatedFields: e.ChangeDesc.UpdatedFields,
		RemovedFields: e.ChangeDesc.RemovedFields,
	}
	detailBytes, err := json.Marshal(detail)
	if err != nil {
		blog.Errorf("run flow, marshal detail failed, name: %s, detail: %+v, err: %v, oid: %s", name, detail, err, e.Oid)
		return false, err
	}

	getLock := f.getLockWithRetry(f.key.LockKey(), e.Oid)
	if !getLock {
		blog.Errorf("run flow, insert node, name: %s, op: %s, do not get lock and return, oid: %s", name, e.OperationType, e.Oid)
		return true, errors.New("get lock failed")
	}

	// already get the lock. prepare to release the lock.
	releaseLock := func() {
		if err := redis.Client().Del(context.Background(), f.key.LockKey()).Err(); err != nil {
			f.metrics.collectLockError()
			blog.ErrorfDepthf(1, "run flow, insert node, name: %s, op: %s, release lock failed, err: %v, oid: %s", name, e.OperationType, err, e.Oid)
		}
	}

	pipe := redis.Client().Pipeline()
	pipe.HMSet(f.key.MainHashKey(), values)
	pipe.Set(f.key.DetailKey(currentCursor), string(detailBytes), 0)
	if _, err := pipe.Exec(); err != nil {
		f.metrics.collectRedisError()
		releaseLock()
		blog.Errorf("run flow, but insert node failed, name: %s, op: %s, err: %v, oid: %s", name, e.OperationType, err, e.Oid)
		return true, err
	}

	// release the lock immediately.
	releaseLock()
	blog.Infof("insert watch event for %s success, op: %s cursor: %s, name: %s, oid: %s",
		f.Collection, e.OperationType, currentCursor, name, e.Oid)
	return false, nil
}

func (f *Flow) getLockWithRetry(name string, oid string) bool {
	getLock := false
	for retry := 0; retry < 10; retry++ {
		// get operate lock to avoid concurrent revise the chain
		success, err := redis.Client().SetNX(context.Background(), f.key.LockKey(), "lock", 5*time.Second).Result()
		if err != nil {
			f.metrics.collectLockError()
			blog.Errorf("get lock failed, err: %v, oid: %s", name, err, oid)
			time.Sleep(500 * time.Millisecond)
			continue
		}

		if !success {
			blog.Warnf("do not get lock, oid: %s", name, oid)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		getLock = true
		break
	}
	return getLock
}
