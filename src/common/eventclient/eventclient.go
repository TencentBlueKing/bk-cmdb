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

package eventclient

import (
	"encoding/json"
	"net/http"

	"gopkg.in/redis.v5"

	"configcenter/src/common"
	"configcenter/src/common/metadata"
	commontypes "configcenter/src/common/types"
	"configcenter/src/common/util"
	"configcenter/src/framework/core/errors"
)

type EventContext struct {
	RequestID   string
	RequestTime commontypes.Time
	ownerID     string
	CacheCli    *redis.Client
}

func NewEventContextByReq(pheader http.Header, cacheCli *redis.Client) *EventContext {
	// TODO get reqid and time from req
	ownerID := util.GetOwnerID(pheader)
	return &EventContext{
		ownerID:     ownerID,
		RequestID:   "xxx-xxxx-xxx-xxx",
		RequestTime: commontypes.Now(),
		CacheCli:    cacheCli,
	}
}

func (c *EventContext) InsertEvent(eventType, objType, action string, curData interface{}, preData interface{}) (err error) {
	if c.CacheCli == nil {
		return errors.New("invalid event context with nil cache client")
	}
	eventID, err := c.CacheCli.Incr(common.EventCacheEventIDKey).Result()
	if err != nil {
		return err
	}
	ei := &metadata.EventInst{
		ID:         int64(eventID),
		EventType:  eventType,
		Action:     action,
		ActionTime: commontypes.Now(),
		ObjType:    objType,
		Data: []metadata.EventData{
			{
				CurData: curData,
				PreData: preData,
			},
		},
		OwnerID:     c.ownerID,
		RequestID:   c.RequestID,
		RequestTime: c.RequestTime,
	}

	value, err := json.Marshal(ei)
	if err != nil {
		return err
	}

	err = c.CacheCli.LPush(common.EventCacheEventQueueKey, value).Err()
	if err != nil {
		return
	}
	return nil
}
