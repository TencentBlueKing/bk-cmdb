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
	"math/rand"
	"time"

	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	"configcenter/src/common/util"
	"configcenter/src/common/watch"
)

func getLoopInternalMinutes() int {
	// give a random sleep interval to avoid clean different resource keys
	// at the same time.
	rand.Seed(time.Now().UnixNano())
	// time range [60, 120] minutes
	return 60 + (rand.Intn(60-1) + 1)
}

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
			if !continueLoop {
				time.Sleep(interval)
				// update rid
				rid = util.GenerateRID()
				// update loop interval
				minute := getLoopInternalMinutes()
				interval = time.Duration(minute) * time.Minute
				continueLoop = true
				blog.Infof("clean expired events for key: %s, loop interval minutes: %d, rid: %s", f.key.Namespace(), minute, rid)
			}

			blog.Infof("start clean expired events for key: %s, rid: ", f.key.Namespace(), rid)

			nodes, err := f.getNodesFromCursor(100, headKey, f.key)
			if err != nil {
				blog.Errorf("clean expired events for key: %s, but get cursor node from head failed, err: %v, rid: %s",
					f.key.Namespace(), err, rid)
				continue
			}

			if len(nodes) == 0 {
				// TODO: repair the node chain when this happen.
				blog.Errorf("clean expired events for key: %s, but got 0 nodes, at least we have tail node, rid: %s",
					f.key.Namespace(), rid)
				continue
			}

			if len(nodes) == 1 {
				// no events is occurred
				if nodes[0].Cursor == headKey {
					blog.Infof("clean expired events for key: %s success, * no events found *, rid: %s",
						f.key.Namespace(), rid)
					continue
				}
				// have only one node, but not target to the head key.
				// something is wrong when this happens.
				// TODO: repair the event chain.
				blog.Errorf("clean expired events for key: %s, but something is wrong. rid: %s", f.key.Namespace(), rid)
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
						blog.Infof("clean expired events for key: %s success, * no expired keys *, head cursor: %s. rid: %s",
							f.key.Namespace(), node.Cursor, rid)
						goto loop
					}
					expiredNodes = append(expiredNodes, node)
					continue
				}

				// at least one event is expired.
				if time.Now().Unix()-int64(node.ClusterTime.Sec) > f.key.TTLSeconds() {
					expiredNodes = append(expiredNodes, node)
				} else {
					// if a node which is not expired occurred, break the loop immediately. nodes after it
					// is definitely not expired.
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

			pipe := f.rds.Pipeline()
			// redirect the head key.
			pipe.HSet(mainKey, headKey, hBytes)
			for _, expire := range expiredNodes {
				// delete expired chain node
				pipe.HDel(mainKey, expire.Cursor)
				// delete expired cursor targeted detail key
				pipe.Del(f.key.DetailKey(expire.Cursor))
			}

			// if last node is tail key, which means we have scan to the end.
			// then we need to update the tail key next_cursor to head key.
			if lastNode.Cursor == tailKey {
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
					blog.Errorf("clean expired events for key: %s, but marshal tail node failed, err: %v, rid: %s",
						f.key.Namespace(), err, rid)
					continue
				}
				pipe.HSet(mainKey, tailKey, string(tBytes))
			}

			// do the clean operation.
			_, err = pipe.Exec()
			if err != nil {
				blog.Errorf("clean expired events for key: %s, but do redis pipeline failed, err: %v, rid: %s",
					f.key.Namespace(), err, rid)
				continue
			}

			if !continueLoop {
				blog.Infof("clean expired events for key: %s success. rid: %s", f.key.Namespace(), rid)
			}
			// sleep a while during the loop
			time.Sleep(30 * time.Second)
		}
	}()
}
