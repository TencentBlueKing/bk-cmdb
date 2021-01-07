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
	commontypes "configcenter/src/common/types"
	"gopkg.in/redis.v5"
	"net/http"
	"testing"
)

func TestNewEventContextByReq(t *testing.T) {
	type args struct {
		pheader  http.Header
		cacheCli *redis.Client
	}
	tests := []struct {
		name string
		args args
		want *EventContext
	}{
		{"", args{http.Header{}, redis.NewClient(&redis.Options{DB: 0})}, &EventContext{
			RequestID:   "xxx-xxxx-xxx-xxx",
			RequestTime: commontypes.Now(),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEventContextByReq(tt.args.pheader, tt.args.cacheCli); got.RequestID != "xxx-xxxx-xxx-xxx" {
				t.Errorf("NewEventContextByReq() = %v, want %v", got.RequestID, tt.want.RequestID)
			}
		})
	}
}

func TestEventContext_InsertEvent(t *testing.T) {
	type fields struct {
		RequestID   string
		RequestTime commontypes.Time
		ownerID     string
		CacheCli    *redis.Client
	}
	type args struct {
		eventType string
		objType   string
		action    string
		curData   interface{}
		preData   interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &EventContext{
				RequestID:   tt.fields.RequestID,
				RequestTime: tt.fields.RequestTime,
				ownerID:     tt.fields.ownerID,
				CacheCli:    tt.fields.CacheCli,
			}
			if err := c.InsertEvent(tt.args.eventType, tt.args.objType, tt.args.action, tt.args.curData, tt.args.preData); (err != nil) != tt.wantErr {
				t.Errorf("EventContext.InsertEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
