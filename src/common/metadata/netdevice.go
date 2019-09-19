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
	DeviceID    uint64     `json:"device_id,omitempty" bson:"device_id,omitempty"`
	DeviceName  string     `json:"device_name,omitempty" bson:"device_name,omitempty"`
	DeviceModel string     `json:"device_model,omitempty" bson:"device_model,omitempty"`
	ObjectID    string     `json:"bk_obj_id" bson:"bk_obj_id,omitempty"`
	ObjectName  string     `json:"bk_obj_name,omitempty" bson:"-"`
	BkVendor    string     `json:"bk_vendor,omitempty" bson:"bk_vendor,omitempty"`
	OwnerID     string     `json:"-" bson:"bk_supplier_account,omitempty"`
	CreateTime  *time.Time `field:"create_time,omitempty" json:"create_time,omitempty" bson:"create_time,omitempty"`
	LastTime    *time.Time `field:"last_time" json:"last_time,omitempty" bson:"last_time,omitempty"`
}

type NetcollectProperty struct {
	NetcollectPropertyID uint64     `json:"netcollect_property_id,omitempty" bson:"netcollect_property_id,omitempty"`
	PropertyID           string     `json:"bk_property_id" bson:"bk_property_id,omitempty"`
	PropertyName         string     `json:"bk_property_name" bson:"-"`
	DeviceID             uint64     `json:"device_id" bson:"device_id,omitempty"`
	DeviceName           string     `json:"device_name" bson:"-"`
	ObjectID             string     `json:"bk_obj_id,omitempty" bson:"-"`
	ObjectName           string     `json:"bk_obj_name,omitempty" bson:"-"`
	OID                  string     `json:"oid,omitempty" bson:"oid,omitempty"`
	Period               string     `json:"period,omitempty" bson:"period,omitempty"`
	Action               string     `json:"action,omitempty" bson:"action,omitempty"`
	OwnerID              string     `json:"-" bson:"bk_supplier_account,omitempty"`
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
	CloudID       int64              `json:"bk_cloud_id" bson:"bk_cloud_id"`
	CloudName     string             `json:"bk_cloud_name" bson:"-"`
	BizID         int64              `json:"bk_biz_id" bson:"bk_biz_id" bson:"bk_biz_id"`
	InnerIP       string             `json:"bk_host_innerip" bson:"bk_host_innerip"`
	Status        NetcollectorStatus `json:"status" bson:"status"`
	TaskID        int64              `json:"task_id" bson:"task_id"`
	DeployTime    *time.Time         `json:"deploy_time" bson:"deploy_time"`
	Version       string             `json:"version" bson:"version"`
	LatestVersion string             `json:"latest_ersion" bson:"latest_ersion"`
	ReportTotal   int64              `json:"report_total" bson:"report_total"`
	Config        NetcollectConfig   `json:"config" bson:"config"`
}

type ParamNetcollectDiscover struct {
	Collectors []Netcollector `json:"collectors"`
}

type NetcollectorStatus struct {
	CollectorStatus string `json:"collector_status" bson:"collector_status"`
	ConfigStatus    string `json:"config_status" bson:"config_status"`
	ReportStatus    string `json:"report_status" bson:"report_status"`
}

// NetcollectorStatus define
const (
	CollectorStatusNormal   = "normal"
	CollectorStatusAbnormal = "abnormal"

	CollectorConfigStatusNormal   = "normal"
	CollectorConfigStatusPending  = "pending"
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
	Action    string   `json:"action"`
	ObjectID  string   `json:"bk_object_id"`
	Query     string   `json:"query"`
	CloudName string   `json:"bk_cloud_name"`
	CloudID   int64    `json:"bk_cloud_id"`
	InnerIP   string   `json:"bk_host_innerip"`
	LastTime  []Time   `json:"last_time"`
	Page      BasePage `json:"page"`
}

type NetcollectReportSummary struct {
	CloudID    int64          `json:"bk_cloud_id"`
	CloudName  string         `json:"bk_cloud_name"`
	LastTime   Time           `json:"last_time"`
	Statistics map[string]int `json:"operation"`
}

type RspNetcollectReport struct {
	Count uint64             `json:"count"`
	Info  []NetcollectReport `json:"info"`
}
type RspNetcollectHistory struct {
	Count uint64              `json:"count"`
	Info  []NetcollectHistory `json:"info"`
}

type ParamNetcollectComfirm struct {
	Reports []NetcollectReport `json:"reports"`
}

type NetcollectReport struct {
	Action        string `json:"action" bson:"action"`
	CloudID       int64  `json:"bk_cloud_id" bson:"bk_cloud_id"`
	CloudName     string `json:"bk_cloud_name" bson:"-"`
	ObjectID      string `json:"bk_obj_id" bson:"bk_obj_id"`
	ObjectName    string `json:"bk_obj_name" bson:"-"`
	InnerIP       string `json:"bk_host_innerip" bson:"bk_host_innerip"`
	Configuration string `json:"configuration" bson:"configuration"`
	OwnerID       string `json:"bk_supplier_account" bson:"bk_supplier_account"`

	InstID       int64                         `json:"bk_inst_id" bson:"bk_inst_id"`
	InstKey      string                        `json:"bk_inst_key" bson:"bk_inst_key"`
	LastTime     Time                          `json:"last_time" bson:"last_time"`
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
	Action       string `json:"action" bson:"-"`
	AsstInstName string `json:"bk_asst_inst_name" bson:"bk_asst_inst_name"`
	// AsstPropertyID string `json:"bk_asst_property_id" bson:"bk_asst_property_id"`
	AsstObjectID   string `json:"bk_asst_obj_id" bson:"bk_asst_obj_id"`
	AsstObjectName string `json:"bk_asst_obj_name" bson:"bk_asst_obj_name"`

	ObjectAsstID  string `json:"bk_obj_asst_id" bson:"bk_obj_asst_id"`
	Configuration string `json:"configuration" bson:"configuration"`
}

type NetcollectReportAsstCond struct {
	PropertyID   string      `json:"bk_property_id" bson:"bk_property_id"`
	PropertyName string      `json:"bk_property_name" bson:"bk_property_name"`
	Value        interface{} `json:"value" bson:"value"`
}

type RspNetcollectConfirm struct {
	ChangeAttributeSuccess    int      `json:"change_attribute_success"`
	ChangeAttributeFailure    int      `json:"change_attribute_failure"`
	ChangeAssociationsSuccess int      `json:"change_associations_success"`
	ChangeAssociationsFailure int      `json:"change_associations_failure"`
	Errors                    []string `json:"errors"`
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
