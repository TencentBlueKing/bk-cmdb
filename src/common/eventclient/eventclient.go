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
	"reflect"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/framework/core/errors"

	"gopkg.in/redis.v5"
)

type EventContext struct {
	RequestID   string
	RequestTime metadata.Time
	TxnID       string
	ownerID     string
	CacheCli    *redis.Client
}

func NewEventContextByReq(pheader http.Header, cacheCli *redis.Client) *EventContext {
	return &EventContext{
		ownerID:     util.GetOwnerID(pheader),
		TxnID:       util.GetHTTPCCTransaction(pheader),
		RequestID:   util.GetHTTPCCRequestID(pheader),
		RequestTime: metadata.Now(),
		CacheCli:    cacheCli,
	}
}

func (c *EventContext) InsertEvent(eventType, objType, action string, curData interface{}, preData interface{}) (err error) {
	if c.CacheCli == nil {
		return errors.New("invalid event context with nil cache client")
	}
	equal, err := instEqual(curData, preData)
	if err != nil {
		return err
	}
	if equal {
		return nil
	}

	eventID, err := c.CacheCli.Incr(common.EventCacheEventIDKey).Result()
	if err != nil {
		return err
	}
	ei := &metadata.EventInst{
		ID:         int64(eventID),
		TxnID:      c.TxnID,
		EventType:  eventType,
		Action:     action,
		ActionTime: metadata.Now(),
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

	if c.TxnID != "" {
		z := redis.Z{
			Score:  float64(time.Now().UTC().Unix()),
			Member: c.TxnID,
		}
		if err = c.CacheCli.ZAddNX(common.EventCacheEventTxnSet, z).Err(); err != nil {
			return err
		}
		return c.CacheCli.LPush(common.EventCacheEventTxnQueuePrefix+c.TxnID, value).Err()
	}
	return c.CacheCli.LPush(common.EventCacheEventQueueKey, value).Err()
}

// instEqual Determine whether the data is consistent before and after the change
func instEqual(predata, curdata interface{}) (bool, error) {
	switch {
	case predata == nil && curdata != nil:
		return false, nil
	case curdata == nil && predata != nil:
		return false, nil
	}

	preData, err := toMap(predata)
	if err != nil {
		return false, err
	}
	curData, err := toMap(curdata)
	if err != nil {
		return false, err
	}
	delete(preData, common.LastTimeField)
	delete(curData, common.LastTimeField)

	return reflect.DeepEqual(preData, curData), nil
}

func toMap(data interface{}) (map[string]interface{}, error) {
	out, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	m := map[string]interface{}{}
	err = json.Unmarshal(out, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}
