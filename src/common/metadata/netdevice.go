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

	"configcenter/src/common/types"
)

type NetcollectDevice struct {
	DeviceID    int64      `json:"device_id,omitempty" bson:"device_id"`
	DeviceName  string     `json:"device_name" bson:"device_name"`
	DeviceModel string     `json:"device_model" bson:"device_model"`
	ObjectID    string     `json:"bk_obj_id" bson:"bk_obj_id"`
	ObjectName  string     `json:"bk_obj_name" bson:"-"`
	BkVendor    string     `json:"bk_vendor" bson:"bk_vendor"`
	OwnerID     string     `json:"-" bson:"bk_supplier_account"`
	CreateTime  *time.Time `json:"create_time,omitempty" bson:"create_time"`
	LastTime    *time.Time `json:"last_time,omitempty" bson:"last_time"`
}

type NetcollectProperty struct {
	NetcollectPropertyID int64      `json:"netcollect_property_id,omitempty" bson:"netcollect_property_id"`
	PropertyID           string     `json:"bk_property_id" bson:"bk_property_id"`
	PropertyName         string     `json:"bk_property_name" bson:"-"`
	DeviceID             int64      `json:"device_id" bson:"device_id"`
	DeviceName           string     `json:"device_name" bson:"-"`
	ObjectID             string     `json:"bk_obj_id,omitempty" bson:"-"`
	ObjectName           string     `json:"bk_obj_name,omitempty" bson:"-"`
	OID                  string     `json:"oid" bson:"oid"`
	Period               string     `json:"period" bson:"period"`
	Action               string     `json:"action" bson:"action"`
	OwnerID              string     `json:"-" bson:"bk_supplier_account"`
	CreateTime           *time.Time `json:"create_time,omitempty" bson:"create_time"`
	LastTime             *time.Time `json:"last_time,omitempty" bson:"last_time"`
	Unit                 string     `json:"unit,omitempty" bson:"-"`
	DeviceModel          string     `json:"device_model" bson:"device_model"`
}

type ParamNetcollectorSearch struct {
	Query string   `json:"query"`
	Page  BasePage `json:"page"`
}
type RspNetcollectorSearch struct {
	Count int64          `json:"count"`
	Info  []Netcollector `json:"info"`
}

type Netcollector struct {
	CloudID    int64              `json:"bk_cloud_id"`
	CloudName  string             `json:"bk_cloud_name"`
	InnerIP    string             `json:"bk_host_innerip"`
	Status     NetcollectorStatus `json:"status"`
	DeployTime time.Time          `json:"deploy_time"`
	Version    string             `json:"version"`
	Config     NetcollectConfig   `json:"config"`
}

type NetcollectorConfig struct {
	CloudID int64            `json:"bk_cloud_id" bson:"bk_cloud_id" `
	InnerIP string           `json:"bk_host_innerip" bson:"bk_host_innerip"`
	Config  NetcollectConfig `json:"config"`
}

type NetcollectorStatus struct {
	CollectorStatus string `json:"collector_status"`
	ConfigStatus    string `json:"config_status"`
	ReportStatus    string `json:"report_status"`
}

// NetcollectorStatus define
const (
	CollectorStatusNormal   = "normal"
	CollectorStatusAbnormal = "abnormal"

	CollectorConfigStatusNormal   = "normal"
	CollectorConfigStatusAbnormal = "abnormal"

	CollectorReportStatusNormal   = "normal"
	CollectorReportStatusAbnormal = "abnormal"
)

type NetcollectConfig struct {
	ScanRange []string `json:"scan_range"`
	Period    string   `json:"period"`
	Community string   `json:"community"`
}

type ParamSearchNetcollectReport struct {
	Action    string   `json:"action"`
	ObjectID  string   `json:"bk_object_id"`
	CloudName string   `json:"bk_cloud_name"`
	CloudID   int64    `json:"bk_cloud_id"`
	InnerIP   string   `json:"bk_host_innerip"`
	Page      BasePage `json:"page"`
}

type NetcollectReportSummary struct {
	CloudID    int64          `json:"bk_cloud_id"`
	CloudName  string         `json:"bk_cloud_name"`
	LastTime   *types.Time    `json:"last_time"`
	Statistics map[string]int `json:"statistics"`
}

type RspNetcollectReport struct {
	Count int64              `json:"count"`
	Info  []NetcollectReport `json:"info"`
}

type NetcollectReport struct {
	Action        string `json:"action" bson:"action"`
	CloudID       int64  `json:"bk_cloud_id" bson:"bk_cloud_id"`
	CloudName     string `json:"bk_cloud_name" bson:"-"`
	ObjectID      string `json:"bk_obj_id" bson:"bk_obj_id"`
	ObjectName    string `json:"bk_obj_name" bson:"-"`
	InnerIP       string `json:"bk_host_innerip" bson:"bk_host_innerip"`
	Configuration string `json:"configuration" bson:"configuration"`

	InstID       int64                         `json:"bk_inst_id" bson:"bk_inst_id"`
	InstKey      string                        `json:"bk_inst_key" bson:"bk_inst_key"`
	LastTime     *types.Time                   `json:"last_time" bson:"last_time"`
	Attributes   []NetcollectReportAttribute   `json:"attributes" bson:"attributes"`
	Associations []NetcollectReportAssociation `json:"associations" bson:"associations"`
}

type NetcollectReportAttribute struct {
	PropertyID   string      `json:"bk_property_id" bson:"bk_property_id"`
	PropertyName string      `json:"bk_property_name" bson:"-"`
	IsRequired   bool        `json:"isrequired" bson:"-"`
	CurValue     interface{} `json:"value" bson:"value"`
	PreValue     interface{} `json:"pre_value" bson:"-"`
}

type NetcollectReportAssociation struct {
	CloudID    int64  `json:"bk_cloud_id" bson:"bk_cloud_id"`
	InstID     int64  `json:"bk_inst_id" bson:"-"`
	InstKey    string `json:"bk_inst_key" bson:"bk_inst_key"`
	ObjectID   string `json:"bk_obj_id" bson:"bk_obj_id"`
	ObjectName string `json:"bk_obj_name" bson:"bk_obj_name"`

	AsstInstID     int64                      `json:"bk_asst_inst_id" bson:"bk_asst_inst_id`
	AsstInstName   int64                      `json:"bk_asst_inst_name" bson:"bk_asst_inst_name`
	AsstObjectID   string                     `json:"bk_asst_obj_id" bson:"bk_asst_obj_id`
	AsstObjectName string                     `json:"bk_asst_obj_name" bson:"bk_asst_obj_name`
	AsstCond       []NetcollectReportAsstCond `json:"asst_cond" bson:"asst_cond`
}

type NetcollectReportAsstCond struct {
	PropertyID   string      `json:"bk_property_id" bson:"bk_property_id`
	PropertyName string      `json:"bk_property_name" bson:"bk_property_name`
	Value        interface{} `json:"value" bson:"value`
}

const (
	ReporctActionCreate = "create"
	ReporctActionUpdate = "update"
	ReporctActionDelete = "delete"
)
