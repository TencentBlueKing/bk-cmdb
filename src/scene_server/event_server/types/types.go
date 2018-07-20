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

package types

import (
	"configcenter/src/common/types"
	"encoding/json"
	"time"
)

// Subscription define
type Subscription struct {
	SubscriptionID   int64       `bson:"subscription_id" json:"subscription_id"`
	SubscriptionName string      `bson:"subscription_name" json:"subscription_name"`
	SystemName       string      `bson:"system_name" json:"system_name"`
	CallbackURL      string      `bson:"callback_url" json:"callback_url"`
	ConfirmMode      string      `bson:"confirm_mode" json:"confirm_mode"`
	ConfirmPattern   string      `bson:"confirm_pattern" json:"confirm_pattern"`
	TimeOut          int64       `bson:"time_out" json:"time_out"`                   // second
	SubscriptionForm string      `bson:"subscription_form" json:"subscription_form"` // json format
	Operator         string      `bson:"operator" json:"operator"`
	OwnerID          string      `bson:"bk_supplier_account" json:"bk_supplier_account"`
	LastTime         *types.Time `bson:"last_time" json:"last_time"`
	Statistics       *Statistics `bson:"-" json:"statistics"`
}

// Report define sending statistic
type Statistics struct {
	Total   int64 `json:"total"`
	Failure int64 `json:"failure"`
}

func (Subscription) TableName() string {
	return TableNameSubscription
}

func (s Subscription) GetCacheKey() string {
	ns := &Subscription{
		SubscriptionID:   s.SubscriptionID,
		CallbackURL:      s.CallbackURL,
		ConfirmMode:      s.ConfirmMode,
		ConfirmPattern:   s.ConfirmPattern,
		SubscriptionForm: s.SubscriptionForm,
		TimeOut:          s.TimeOut,
	}
	b, _ := json.Marshal(ns)
	return string(b)
}

func (s Subscription) GetTimeout() time.Duration {
	return time.Second * time.Duration(s.TimeOut)
}

type EventInst struct {
	ID          int64       `json:"event_id,omitempty"`
	EventType   string      `json:"event_type"`
	Action      string      `json:"action"`
	ActionTime  types.Time  `json:"action_time"`
	ObjType     string      `json:"obj_type"`
	Data        []EventData `json:"data"`
	RequestID   string      `json:"request_id"`
	OwnerID     string      `json:"supplier_account"`
	RequestTime types.Time  `json:"request_time"`
}

func (e *EventInst) MarshalBinary() (data []byte, err error) {
	return json.Marshal(e)
}

type EventData struct {
	CurData interface{} `json:"cur_data"`
	PreData interface{} `json:"pre_data"`
}

func (e *EventInst) GetType() string {
	if e.EventType == EventTypeRelation {
		return e.ObjType
	}
	return e.ObjType + e.Action
}

type EventInstCtx struct {
	EventInst
	Raw string
}

type DistInst struct {
	EventInst
	DstbID         int64 `json:"distribution_id"`
	SubscriptionID int64 `json:"subscription_id"`
}

type DistInstCtx struct {
	DistInst
	Raw string
}
