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
	"math/rand"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"configcenter/src/common/watch"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/driver/redis"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getLoopInternalMinutes() int {
	// give a random sleep interval to avoid clean different resource keys
	// at the same time.
	rand.Seed(time.Now().UnixNano())
	// time range [30, 60）minutes
	return int(30 + 6000*rand.Float32()/200)
}

// gc count for each time
const gcCount = 1000

func (f *Flow) cleanExpiredEvents() {
	blog.Infof("will clean expired events job for key: %s after sleep", f.key.Namespace())
	go func() {
		// sleep at first, so that user's can still consume events when we restart system or resume
		// from a fault.
		interval := time.Duration(getLoopInternalMinutes()) * time.Minute
		time.Sleep(interval)

		mainKey := f.key.MainHashKey()
		headKey := f.key.HeadKey()
		tailKey := f.key.TailKey()
		continueLoop := true
		rid := util.GenerateRID()
	loop:
		for {

			if !f.isMaster.IsMaster() {
				blog.V(4).Infof("clean expired events for key: %s, but not master, skip.", f.key.Namespace())
				time.Sleep(60 * time.Second)
				continue
			}

			if !continueLoop {
				// update loop interval
				interval = time.Duration(getLoopInternalMinutes())
				// update rid
				rid = util.GenerateRID()
				blog.Infof("clean expired events for key: %s, loop interval minutes: %d, rid: %s",
					f.key.Namespace(), interval, rid)

				time.Sleep(interval * time.Minute)
				continueLoop = true
			}

			blog.Infof("start clean expired events for key: %s, rid: %s", f.key.Namespace(), rid)

			// get operate lock to avoid concurrent revise the chain
			success, err := redis.Client().SetNX(context.Background(), f.key.LockKey(), 1, 15*time.Second).Result()
			if err != nil {
				blog.Errorf("clean expired events for key: %s, but get lock failed, err: %v, rid: %s", f.key.Namespace(), err, rid)
				time.Sleep(10 * time.Second)
				continue
			}

			if !success {
				blog.Errorf("clean expired events for key: %s, can not get lock, rid: %s", f.key.Namespace(), rid)
				time.Sleep(10 * time.Second)
				continue
			}

			// already get the lock. prepare to release the lock.
			releaseLock := func() {
				if err := redis.Client().Del(context.Background(), f.key.LockKey()).Err(); err != nil {
					blog.Errorf("clean expired events for key: %s, but delete lock failed, err: %v, rid: %s", f.key.Namespace(), err, rid)
				}
			}

			nodes, err := f.getNodesFromCursor(gcCount, headKey, f.key)
			if err != nil {
				blog.Errorf("clean expired events for key: %s, but get cursor node from head failed, err: %v, rid: %s",
					f.key.Namespace(), err, rid)
				continueLoop = false
				continue
			}

			if len(nodes) == 0 {
				// TODO: repair the node chain when this happen.
				blog.Errorf("clean expired events for key: %s, but got 0 nodes, at least we have tail node, rid: %s",
					f.key.Namespace(), rid)
				continueLoop = false
				continue
			}

			if len(nodes) == 1 {
				// no events is occurred
				if nodes[0].NextCursor == headKey {
					blog.Infof("clean expired events for key: %s success, * no events found *, rid: %s", f.key.Namespace(), rid)
					continueLoop = false
					continue
				}
				// have only one node, but not target to the head key.
				// something is wrong when this happens.
				// TODO: repair the event chain.
				blog.Errorf("clean expired events for key: %s, but something is wrong. rid: %s", f.key.Namespace(), rid)
				continueLoop = false
				continue
			}
			expiredNodes := make([]*watch.ChainNode, 0)
			for idx, node := range nodes {
				if idx == 0 {
					if node.Cursor == tailKey {
						// the first node is tail node. so there is no events.
						continueLoop = false
						blog.Infof("clean expired events for key: %s success, rid: %s", f.key.Namespace(), rid)
						goto loop
					}
					// if the first node is not expired, then no node is expired after it.
					if time.Now().Unix()-int64(node.ClusterTime.Sec) <= f.key.TTLSeconds() {
						blog.V(4).Infof("clean expired events for key: %s success, * no expired keys *, head cursor: %s. rid: %s",
							f.key.Namespace(), node.Cursor, rid)
						continueLoop = false
						goto loop
					}
					expiredNodes = append(expiredNodes, node)
					continue
				}

				// at least one event is expired.
				if node.Cursor != tailKey && time.Now().Unix()-int64(node.ClusterTime.Sec) > f.key.TTLSeconds() {
					expiredNodes = append(expiredNodes, node)
				} else {
					// if a node which is not expired occurred, break the loop immediately. nodes after it
					// is definitely not expired.
					continueLoop = false
					break
				}
			}

			lastNode := expiredNodes[len(expiredNodes)-1]
			// set head cursor to last node's next node.
			headNode := &watch.ChainNode{
				Cursor:     headKey,
				NextCursor: lastNode.NextCursor,
			}

			hBytes, err := json.Marshal(headNode)
			if err != nil {
				blog.Errorf("clean expired events for key: %s, but marshal head node failed, err: %v, rid: %s", f.key.Namespace(), err, rid)
				continue
			}

			expireCursor := make([]string, 0)
			pipe := redis.Client().Pipeline()
			// redirect the head key.
			pipe.HSet(mainKey, headKey, hBytes)
			for _, expire := range expiredNodes {
				// do not delete the last node before tail, in case no event occurred after the last event and watch 
				// with that cursor failed
				if expire.NextCursor == tailKey {
					continue
				}

				// delete expired chain node
				pipe.HDel(mainKey, expire.Cursor)
				// delete expired cursor targeted detail key
				pipe.Del(f.key.DetailKey(expire.Cursor))
				expireCursor = append(expireCursor, expire.Cursor)
			}

			// if last node is tail key, which means we have scan to the end.
			// then we need to update the tail key next_cursor to head key.
			if lastNode.NextCursor == tailKey {
				// update continue loop flag, we have to the end the loop,
				// and sleep wait for another loop.
				continueLoop = false

				tailNode := &watch.ChainNode{
					Cursor: tailKey,
					// the last token must be record,
					// we use this to resume watch.
					Token:      lastNode.Token,
					NextCursor: headKey,
				}
				tBytes, err := json.Marshal(tailNode)
				if err != nil {
					releaseLock()

					blog.Errorf("clean expired events for key: %s, but marshal tail node failed, err: %v, rid: %s",
						f.key.Namespace(), err, rid)
					continueLoop = false
					continue
				}
				pipe.HSet(mainKey, tailKey, string(tBytes))
			}

			// do the clean operation.
			_, err = pipe.Exec()
			if err != nil {
				releaseLock()

				blog.Errorf("clean expired events for key: %s, but do redis pipeline failed, err: %v, rid: %s",
					f.key.Namespace(), err, rid)
				time.Sleep(5 * time.Second)
				continue
			}
			releaseLock()

			blog.Infof("clean expired events for key: %s success, expire cursor: %v. rid: %s", f.key.Namespace(), expireCursor, rid)
			// sleep a while during the loop
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

// cleanDelArchiveData is to clean the table cc_DelArchive data which is a week ago.
// we do this everyday at a fixed time.
// we find the expired data with _id.
func (f *Flow) cleanDelArchiveData() {
	blog.Infof("start clean cc_DelArchive data job success.")
	go func() {
		for {
			if time.Now().Hour() != 1 {
				time.Sleep(5 * time.Minute)
				continue
			}
			rid := util.GenerateRID()
			blog.Infof("start do clean cc_DelArchive data job, rid: %s", rid)
			f.doClean(rid)
			blog.Infof("start do clean cc_DelArchive data job done, rid: %s", rid)
		}
	}()
}

func (f *Flow) doClean(rid string) {
	timeout := time.After(time.Hour)
	for {
		select {
		case <-timeout:
			blog.Errorf("do clean cc_DelArchive data job timeout, rid: %s", rid)
			return
		default:
		}
		time.Sleep(5 * time.Minute)

		if !f.isMaster.IsMaster() {
			blog.Infof("try to clean cc_DelArchive data job, but not master, skip, rid: %s", rid)
			continue
		}

		blog.Infof("do clean cc_DelArchive data job, rid: %s", rid)

		// it's time to do the clean job.
		// generate a ObjectID with a time.
		week := time.Now().Unix() - 7*24*60*60
		weekAgo := time.Unix(week, 0)
		oid := primitive.NewObjectIDFromTimestamp(weekAgo)

		// count the data older than this oid
		filter := mapstr.MapStr{
			"_id": mapstr.MapStr{
				common.BKDBLT: oid,
			},
		}

		count, err := mongodb.Client().Table(common.BKTableNameDelArchive).Find(filter).Count(context.Background())
		if err != nil {
			blog.Errorf("clean cc_DelArchive data, but count expired data in %s failed. rid: %s", common.BKTableNameDelArchive, rid)
			continue
		}

		blog.Infof("do clean cc_DelArchive data job, found %d expired docs, rid: %s", count, rid)

		pageSize := 500
		success := true
		for start := 0; start < int(count); start += pageSize {
			docs := make([]archived, pageSize)
			err = mongodb.Client().Table(common.BKTableNameDelArchive).Find(filter).Fields("oid").All(context.Background(), &docs)
			if err != nil {
				blog.Errorf("clean cc_DelArchive data, but find expired data failed, err: %v, rid: %s", err, rid)
				time.Sleep(10 * time.Second)
				success = false
				continue
			}

			oids := make([]string, len(docs))
			for idx, doc := range docs {
				oids[idx] = doc.Oid
			}

			delFilter := mapstr.MapStr{
				"oid": mapstr.MapStr{
					common.BKDBIN: oids,
				},
			}

			err = mongodb.Client().Table(common.BKTableNameDelArchive).Delete(context.Background(), delFilter)
			if err != nil {
				blog.Errorf("clean cc_DelArchive data, but delete data failed, err: %v, rid: %s", err, rid)
				time.Sleep(10 * time.Second)
				success = false
				continue
			}
			// sleep a while
			time.Sleep(10 * time.Second)
		}

		if success {
			blog.Infof("clean cc_DelArchive data success, delete %d docs, rid: %s", count, rid)
		} else {
			blog.Infof("clean cc_DelArchive data job done, but part of it is failed, rid: %s", rid)
		}

		// finished the for loop.
		return
	}
}

type archived struct {
	Oid string `bson:"oid"`
}
