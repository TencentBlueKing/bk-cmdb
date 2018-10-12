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

package metadata

import (
	"configcenter/src/common/types"
	"database/sql/driver"
	"reflect"
	"testing"
	"time"
)

func TestSubscription_TableName(t *testing.T) {
	type fields struct {
		SubscriptionID   int64
		SubscriptionName string
		SystemName       string
		CallbackURL      string
		ConfirmMode      string
		ConfirmPattern   string
		TimeOut          int64
		SubscriptionForm string
		Operator         string
		OwnerID          string
		LastTime         *types.Time
		Statistics       *Statistics
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"", fields{}, "cc_Subscription"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Subscription{
				SubscriptionID:   tt.fields.SubscriptionID,
				SubscriptionName: tt.fields.SubscriptionName,
				SystemName:       tt.fields.SystemName,
				CallbackURL:      tt.fields.CallbackURL,
				ConfirmMode:      tt.fields.ConfirmMode,
				ConfirmPattern:   tt.fields.ConfirmPattern,
				TimeOut:          tt.fields.TimeOut,
				SubscriptionForm: tt.fields.SubscriptionForm,
				Operator:         tt.fields.Operator,
				OwnerID:          tt.fields.OwnerID,
				LastTime:         tt.fields.LastTime,
				Statistics:       tt.fields.Statistics,
			}
			if got := s.TableName(); got != tt.want {
				t.Errorf("Subscription.TableName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubscription_GetCacheKey(t *testing.T) {
	type fields struct {
		SubscriptionID   int64
		SubscriptionName string
		SystemName       string
		CallbackURL      string
		ConfirmMode      string
		ConfirmPattern   string
		TimeOut          int64
		SubscriptionForm string
		Operator         string
		OwnerID          string
		LastTime         *types.Time
		Statistics       *Statistics
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"", fields{SubscriptionForm: AssociationFieldAssociationForward}, `{"subscription_id":0,"subscription_name":"","system_name":"","callback_url":"","confirm_mode":"","confirm_pattern":"","time_out":0,"subscription_form":"bk_asst_forward","operator":"","bk_supplier_account":"","last_time":null,"statistics":null}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Subscription{
				SubscriptionID:   tt.fields.SubscriptionID,
				SubscriptionName: tt.fields.SubscriptionName,
				SystemName:       tt.fields.SystemName,
				CallbackURL:      tt.fields.CallbackURL,
				ConfirmMode:      tt.fields.ConfirmMode,
				ConfirmPattern:   tt.fields.ConfirmPattern,
				TimeOut:          tt.fields.TimeOut,
				SubscriptionForm: tt.fields.SubscriptionForm,
				Operator:         tt.fields.Operator,
				OwnerID:          tt.fields.OwnerID,
				LastTime:         tt.fields.LastTime,
				Statistics:       tt.fields.Statistics,
			}
			if got := s.GetCacheKey(); got != tt.want {
				t.Errorf("Subscription.GetCacheKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubscription_GetTimeout(t *testing.T) {
	type fields struct {
		SubscriptionID   int64
		SubscriptionName string
		SystemName       string
		CallbackURL      string
		ConfirmMode      string
		ConfirmPattern   string
		TimeOut          int64
		SubscriptionForm string
		Operator         string
		OwnerID          string
		LastTime         *types.Time
		Statistics       *Statistics
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Duration
	}{
		{"", fields{TimeOut: 1}, time.Second},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Subscription{
				SubscriptionID:   tt.fields.SubscriptionID,
				SubscriptionName: tt.fields.SubscriptionName,
				SystemName:       tt.fields.SystemName,
				CallbackURL:      tt.fields.CallbackURL,
				ConfirmMode:      tt.fields.ConfirmMode,
				ConfirmPattern:   tt.fields.ConfirmPattern,
				TimeOut:          tt.fields.TimeOut,
				SubscriptionForm: tt.fields.SubscriptionForm,
				Operator:         tt.fields.Operator,
				OwnerID:          tt.fields.OwnerID,
				LastTime:         tt.fields.LastTime,
				Statistics:       tt.fields.Statistics,
			}
			if got := s.GetTimeout(); got != tt.want {
				t.Errorf("Subscription.GetTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventInst_MarshalBinary(t *testing.T) {
	type fields struct {
		ID          int64
		EventType   string
		Action      string
		ActionTime  types.Time
		ObjType     string
		Data        []EventData
		OwnerID     string
		RequestID   string
		RequestTime types.Time
	}
	tests := []struct {
		name     string
		fields   fields
		wantData []byte
		wantErr  bool
	}{
		{"", fields{EventType: "custom"}, []byte(`{"event_type":"custom","action":"","action_time":"0001-01-01T00:00:00Z","obj_type":"","data":null,"bk_supplier_account":"","request_id":"","request_time":"0001-01-01T00:00:00Z"}`), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EventInst{
				ID:          tt.fields.ID,
				EventType:   tt.fields.EventType,
				Action:      tt.fields.Action,
				ActionTime:  tt.fields.ActionTime,
				ObjType:     tt.fields.ObjType,
				Data:        tt.fields.Data,
				OwnerID:     tt.fields.OwnerID,
				RequestID:   tt.fields.RequestID,
				RequestTime: tt.fields.RequestTime,
			}
			gotData, err := e.MarshalBinary()
			if (err != nil) != tt.wantErr {
				t.Errorf("EventInst.MarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("EventInst.MarshalBinary() = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}

func TestEventInst_GetType(t *testing.T) {
	type fields struct {
		ID          int64
		EventType   string
		Action      string
		ActionTime  types.Time
		ObjType     string
		Data        []EventData
		OwnerID     string
		RequestID   string
		RequestTime types.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"", fields{ObjType: "obj", Action: "act"}, `objact`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EventInst{
				ID:          tt.fields.ID,
				EventType:   tt.fields.EventType,
				Action:      tt.fields.Action,
				ActionTime:  tt.fields.ActionTime,
				ObjType:     tt.fields.ObjType,
				Data:        tt.fields.Data,
				OwnerID:     tt.fields.OwnerID,
				RequestID:   tt.fields.RequestID,
				RequestTime: tt.fields.RequestTime,
			}
			if got := e.GetType(); got != tt.want {
				t.Errorf("EventInst.GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfirmMode_Scan(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name    string
		n       *ConfirmMode
		args    args
		wantErr bool
	}{
		{"", new(ConfirmMode), args{[]byte(``)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.Scan(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("ConfirmMode.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfirmMode_Value(t *testing.T) {
	tests := []struct {
		name    string
		n       ConfirmMode
		want    driver.Value
		wantErr bool
	}{
		{"", ConfirmMode(`string`), `string`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("ConfirmMode.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConfirmMode.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
