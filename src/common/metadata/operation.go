/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017,-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package metadata

import (
	"time"
)

type ChartConfig struct {
	ConfigID   uint64 `json:"config_id" bson:"config_id"`
	Metadata   `field:"metadata" json:"metadata" bson:"metadata"`
	ReportType string    `json:"report_type" bson:"report_type"`
	Name       string    `json:"name" bson:"name"`
	CreateTime time.Time `json:"create_time" bson:"create_time"`
	OwnerID    string    `json:"bk_supplier_account" bson:"bk_supplier_account"`
	ObjID      string    `json:"bk_obj_id" bson:"bk_obj_id"`
	Width      string    `json:"width" bson:"width"`
	ChartType  string    `json:"chart_type" bson:"chart_type"`
	Field      string    `json:"field" bson:"field"`
}

type ChartPosition struct {
	BizID    int64               `json:"bk_biz_id" bson:"bk_biz_id"`
	Position map[string][]uint64 `json:"position" bson:"position"`
	OwnerID  string              `json:"bk_supplier_account" bson:"bk_supplier_account"`
}

type ModelInstChange map[string]*InstChangeCount

type InstChangeCount struct {
	Create int64 `json:"create" bson:"create"`
	Update int64 `json:"update" bson:"update"`
	Delete int64 `json:"delete" bson:"delete"`
}

type AggregateIntResponse struct {
	BaseResp `json:",inline"`
	Data     []IntIDCount `json:"data"`
}

type IntIDCount struct {
	Id    int64 `json:"id" bson:"_id"`
	Count int64 `json:"count" bson:"count"`
}

type AggregateStringResponse struct {
	BaseResp `json:",inline"`
	Data     []StringIDCount `json:"data"`
}

type StringIDCount struct {
	Id    string `json:"id" bson:"_id"`
	Count int64  `json:"count" bson:"count"`
}

type UpdateInstCount struct {
	Id    UpdateID `json:"id" bson:"_id"`
	Count int64    `json:"count" bson:"count"`
}

type UpdateID struct {
	ObjID  string `json:"bk_obj_id" bson:"bk_obj_id"`
	InstID int64  `json:"bk_inst_id" bson:"bk_inst_id"`
}

type HostChangeChartData struct {
	ReportType string                    `json:"report_type" bson:"report_type"`
	Data       map[string][]BizHostChart `json:"data" bson:"data"`
	OwnerID    string                    `json:"bk_supplier_account" bson:"bk_supplier_account"`
}

type BizHostChart struct {
	Id    time.Time `json:"id" bson:"id"`
	Count int64     `json:"count" bson:"count"`
}

type IDStringCountInn64 struct {
	Id    string `json:"id" bson:"id"`
	Count int64  `json:"count" bson:"count"`
}

type ChartData struct {
	ReportType string      `json:"report_type" bson:"report_type"`
	Data       interface{} `json:"data" data:"data"`
	OwnerID    string      `json:"bk_supplier_account" bson:"bk_supplier_account"`
}

type SearchChartResponse struct {
	BaseResp `json:",inline"`
	Data     SearchChartConfig `json:"data"`
}

type SearchChartCommon struct {
	BaseResp `json:",inline"`
	Data     CommonSearchChart `json:"data"`
}

type CommonSearchChart struct {
	Count uint64      `json:"count"`
	Info  ChartConfig `json:"info"`
}

type SearchChartConfig struct {
	Count uint64                   `json:"count"`
	Info  map[string][]ChartConfig `json:"info"`
}

type CloudMapping struct {
	CreateTime time.Time `json:"create_time" bson:"create_time"`
	LastTime   time.Time `json:"last_time" bson:"lsat_time"`
	CloudName  string    `json:"bk_cloud_name" bson:"bk_cloud_name"`
	OwnerID    string    `json:"bk_supplier_account" bson:"bk_supplier_account"`
	CloudID    int64     `json:"bk_cloud_id" bson:"bk_cloud_id"`
}

type AttributesOptions []AttributesOption

type AttributesOption struct {
	Id        string `json:"id" bson:"id"`
	Name      string `json:"name" bson:"name"`
	Type      string `json:"type" bson:"type"`
	IsDefault string `json:"is_default" bson:"is_default"`
}
