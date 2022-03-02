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

package hostidentifier

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"
	"configcenter/src/common/watch"
	"configcenter/src/storage/dal/redis"

	"github.com/tidwall/gjson"
)

// Event operate watch host identifier
type Event struct {
	engine   *backbone.Engine
	cursor   string
	errFreq  util.ErrFrequencyInterface
	redisCli redis.Client
}

// getEvent watch to get hostIdentifier event
func (e *Event) getEvent(header http.Header, rid string, preStatus bool) ([]*IdentifierEvent, string, bool) {
	// 如果之前节点的状态为从节点（或刚启动时为false），需要从redis中拿到最新的cursor或者直接从当前开始watch
	if !preStatus {
		newCursor, err := e.redisCli.Get(context.Background(), hostIdentifierCursor).Result()
		if err != nil {
			blog.Warnf("get host identity cursor from redis fail, start watch from now, err: %v， rid: %s", err, rid)
			// 从redis里拿不到cursor，从当前时间watch
			newCursor = ""
		}
		e.cursor = newCursor
	}

	options := &watch.WatchEventOptions{
		EventTypes: []watch.EventType{watch.Create, watch.Update},
		Resource:   watch.HostIdentifier,
		Cursor:     e.cursor,
	}

	events, err := e.engine.CoreAPI.CacheService().Cache().Event().WatchEvent(context.Background(), header, options)
	if err != nil {
		if err.GetCode() == common.CCErrEventChainNodeNotExist {
			// 设置从当前时间开始watch
			e.cursor = ""
			if err := e.redisCli.Del(context.Background(), hostIdentifierCursor).Err(); err != nil {
				blog.Errorf("delete redis key failed, key: %s, err: %v, rid: %s", hostIdentifierCursor, err, rid)
			}

			blog.Errorf("watch identifier event failed, reset watch it from current time, err: %v, rid: %s", err, rid)
			return nil, "", false
		}

		// 同一个错误如果出现太频繁，设置从当前时间开始watch
		if e.errFreq.IsErrAlwaysAppear(err) {
			e.cursor = ""
			if err := e.redisCli.Del(context.Background(), hostIdentifierCursor).Err(); err != nil {
				blog.Errorf("delete redis key failed, key: %s, err: %v, rid: %s", hostIdentifierCursor, err, rid)
			}
			e.errFreq.Release()
			blog.Errorf("watch frequent same error, reset watch it from current time, err: %v, rid: %s", err, rid)
			return nil, "", false
		}

		blog.Errorf("watch identifier event failed, err: %v, rid: %s", err, rid)
		return nil, "", false
	}

	e.errFreq.Release()

	if !gjson.Get(*events, "bk_watched").Bool() {
		e.cursor = gjson.Get(*events, "bk_events.0.bk_cursor").String()
		blog.Warnf("can not get host identifier watch event, rid: %s", rid)
		return nil, "", false
	}

	rawEvents := gjson.Get(*events, "bk_events").Array()
	if len(rawEvents) == 0 {
		blog.Errorf("can not get events from bk_events, watchEvent: %v, rid: %s", *events, rid)
		return nil, "", false
	}

	eventArr := make([]*IdentifierEvent, 0)

	for _, rawEvent := range rawEvents {
		eventDetail := rawEvent.Map()["bk_detail"]
		if !eventDetail.Exists() {
			continue
		}

		eventDetailMap := eventDetail.Map()
		if !eventDetailMap[common.BKCloudIDField].Exists() || !eventDetailMap[common.BKHostInnerIPField].Exists() ||
			!eventDetailMap[common.BKHostIDField].Exists() || !eventDetailMap[common.BKOSTypeField].Exists() {
			blog.Errorf("the eventDetail message is error, event: %v, rid: %s", eventDetailMap, rid)
			continue
		}

		event := &IdentifierEvent{}
		err := json.Unmarshal([]byte(eventDetail.String()), event)
		if err != nil {
			blog.Errorf("unmarshal hostIdentifier error, val: %v, err: %v, rid: %s", eventDetail.String(), err, rid)
			continue
		}
		event.RawEvent = eventDetail.String()
		eventArr = append(eventArr, event)
	}

	cursor := rawEvents[len(rawEvents)-1].Get("bk_cursor").String()

	return eventArr, cursor, len(eventArr) > 0
}

// setCursor set watch cursor
func (e *Event) setCursor(cursor string, rid string) {
	// 保存新的cursor到内存和redis中
	e.cursor = cursor

	redisFailCount := 0
	for redisFailCount < retryTimes {
		if err := e.redisCli.Set(context.Background(), hostIdentifierCursor, cursor, 3*time.Hour).Err(); err != nil {
			blog.Errorf("set redis error, key: %s, val: %s, err: %v, rid: %s", hostIdentifierCursor, cursor, err, rid)
			redisFailCount++
			sleepForFail(redisFailCount)
			continue
		}
		break
	}
}

// IdentifierEvent host identifier event
type IdentifierEvent struct {
	HostID   int64  `json:"bk_host_id" bson:"bk_host_id"`
	CloudID  int64  `json:"bk_cloud_id" bson:"bk_cloud_id"`
	InnerIP  string `json:"bk_host_innerip" bson:"bk_host_innerip"`
	OSType   string `json:"bk_os_type" bson:"bk_os_type"`
	RawEvent string `json:"raw_event" bson:"raw_event"`
}
