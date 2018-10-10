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
	DeviceID    int64      `json:"device_id,omitempty" bson:"device_id,omitempty"`
	DeviceName  string     `json:"device_name,omitempty" bson:"device_name,omitempty"`
	DeviceModel string     `json:"device_model,omitempty" bson:"device_model,omitempty"`
	ObjectID    string     `json:"bk_obj_id" bson:"bk_obj_id,omitempty"`
	ObjectName  string     `json:"bk_obj_name,omitempty" bson:"-"`
	BkVendor    string     `json:"bk_vendor,omitempty" bson:"bk_vendor,omitempty"`
	OwnerID     string     `json:"-" bson:"bk_supplier_account"`
	CreateTime  *time.Time `field:"create_time,omitempty" json:"create_time,omitempty" bson:"create_time,omitempty"`
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
	OID                  string     `json:"oid,omitempty" bson:"oid"`
	Period               string     `json:"period,omitempty" bson:"period"`
	Action               string     `json:"action,omitempty" bson:"action"`
	OwnerID              string     `json:"-" bson:"bk_supplier_account"`
	CreateTime           *time.Time `field:"create_time,omitempty" json:"create_time,omitempty" bson:"create_time,omitempty"`
	LastTime             *time.Time `field:"last_time,omitempty" json:"last_time,omitempty" bson:"last_time,omitempty"`
	Unit                 string     `json:"unit,omitempty" bson:"-"`
	DeviceModel          string     `json:"device_model,omitempty" bson:"-"`
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
	CloudID       int64              `json:"bk_cloud_id"`
	CloudName     string             `json:"bk_cloud_name"`
	InnerIP       string             `json:"bk_host_innerip"`
	Status        NetcollectorStatus `json:"status"`
	DeployTime    time.Time          `json:"deploy_time"`
	Version       string             `json:"version"`
	LatestVersion string             `json:"latest_ersion"`
	ReportTotal   int64              `json:"report_total"`
	Config        NetcollectConfig   `json:"config"`
}

type NetcollectorConfig struct {
	CloudID int64            `json:"bk_cloud_id" bson:"bk_cloud_id" `
	InnerIP string           `json:"bk_host_innerip" bson:"bk_host_innerip"`
	Config  NetcollectConfig `json:"config" bson:"config"`
}

type ParamNetcollectDiscover struct {
	Collectors []NetcollectorConfig `json:"collectors"`
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
	ScanRange []string `json:"scan_range" bson:"scan_range"`
	Period    string   `json:"period" bson:"period"`
	Community string   `json:"community" bson:"community"`
}

type ParamSearchNetcollectReport struct {
	Action    string       `json:"action"`
	ObjectID  string       `json:"bk_object_id"`
	Query     string       `json:"query"`
	CloudName string       `json:"bk_cloud_name"`
	CloudID   int64        `json:"bk_cloud_id"`
	InnerIP   string       `json:"bk_host_innerip"`
	LastTime  []types.Time `json:"last_time"`
	Page      BasePage     `json:"page"`
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
type RspNetcollectHistory struct {
	Count int64               `json:"count"`
	Info  []NetcollectHistory `json:"info"`
}

type ParamNetcollectComfirm struct {
	Reports []NetcollectReport `json:reports`
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

type NetcollectHistory struct {
	NetcollectReport `json:",inline" bson:",inline"`
	Success          bool `json:"success" bson:"success"`
}

type NetcollectReportAttribute struct {
	PropertyID   string      `json:"bk_property_id" bson:"bk_property_id"`
	PropertyName string      `json:"bk_property_name" bson:"-"`
	IsRequired   bool        `json:"isrequired" bson:"-"`
	CurValue     interface{} `json:"value" bson:"value"`
	PreValue     interface{} `json:"pre_value" bson:"-"`

	Method  string `json:"method,omitempty" bson:"-"`
	Success bool   `json:"success,omitempty" bson:"-"`
	Error   string `json:"error,omitempty" bson:"-"`
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

type RspNetcollectConfirm struct {
	ChangeAttributeSuccess int      `json:"change_attribute_success"`
	ChangeAttributeFailure int      `json:"change_attribute_failure"`
	ChangeRelationSuccess  int      `json:"change_relation_success"`
	ChangeRelationFailure  int      `json:"change_relation_failure"`
	Errors                 []string `json:"errors"`
}

const (
	ReporctActionCreate = "create"
	ReporctActionUpdate = "update"
	ReporctActionDelete = "delete"
)
const (
	ReporctMethodAccept = "accept"
	ReporctMethodIgnore = "ignore"
)
