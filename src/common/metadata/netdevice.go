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
	"time"
)

type NetcollectDevice struct {
	DeviceID    int64      `json:"device_id,omitempty" bson:"device_id,omitempty"`
	DeviceName  string     `json:"device_name" bson:"device_name,omitempty"`
	DeviceModel string     `json:"device_model" bson:"device_model,omitempty"`
	ObjectID    string     `json:"bk_obj_id" bson:"bk_obj_id,omitempty"`
	ObjectName  string     `json:"bk_obj_name" bson:"-"`
	BkVendor    string     `json:"bk_vendor" bson:"bk_vendor,omitempty"`
	OwnerID     string     `json:"-" bson:"bk_supplier_account"`
	CreateTime  *time.Time `field:"create_time" json:"create_time,omitempty" bson:"create_time,omitempty"`
	LastTime    *time.Time `field:"last_time" json:"last_time,omitempty" bson:"last_time,omitempty"`
}

type NetcollectProperty struct {
	NetcollectPropertyID int64      `json:"netcollect_property_id,omitempty" bson:"netcollect_property_id,omitempty"`
	PropertyID           string     `json:"bk_property_id" bson:"bk_property_id,omitempty"`
	PropertyName         string     `json:"bk_property_name" bson:"-"`
	DeviceID             int64      `json:"device_id" bson:"device_id,omitempty"`
	DeviceName           string     `json:"device_name" bson:"-"`
	ObjectID             string     `json:"bk_obj_id,omitempty" bson:"-"`
	ObjectName           string     `json:"bk_obj_name,omitempty" bson:"-"`
	OID                  string     `json:"oid" bson:"oid"`
	Period               string     `json:"period" bson:"period"`
	Action               string     `json:"action" bson:"action"`
	OwnerID              string     `json:"-" bson:"bk_supplier_account"`
	CreateTime           *time.Time `field:"create_time,omitempty" json:"create_time,omitempty" bson:"create_time,omitempty"`
	LastTime             *time.Time `field:"last_time,omitempty" json:"last_ time,omitempty" bson:"last_time,omitempty"`
	Unit                 string     `json:"unit,omitempty" bson:"-"`
	DeviceModel          string     `json:"device_model" bson:"device_model"`
}
