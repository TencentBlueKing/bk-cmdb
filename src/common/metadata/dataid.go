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

type GseConfigAddStreamToParams struct {
	Metadata  GseConfigAddStreamToMetadata `json:"metadata"`
	Operation GseConfigOperation           `json:"operation"`
	StreamTo  GseConfigStreamTo            `json:"stream_to"`
}

type GseConfigAddStreamToMetadata struct {
	PlatName GseConfigPlatName       `json:"plat_name"`
	Label    *GseConfigStreamToLabel `json:"label,omitempty"`
}

type GseConfigPlatName string

const (
	GseConfigPlatTglog     GseConfigPlatName = "tglog"
	GseConfigPlatTdm       GseConfigPlatName = "tdm"
	GseConfigPlatBkmonitor GseConfigPlatName = "bkmonitor"
	GseConfigPlatTgdp      GseConfigPlatName = "tgdp"
)

type GseConfigStreamToLabel struct {
	BizID   int64  `json:"bk_biz_id,omitempty"`
	BizName string `json:"bk_biz_name,omitempty"`
}

type GseConfigOperation struct {
	OperatorName string `json:"operator_name"`
}

type GseConfigStreamTo struct {
	Name        string                   `json:"name"`
	ReportMode  GseConfigReportMode      `json:"report_mode"`
	DataLogPath string                   `json:"data_log_path,omitempty"`
	Kafka       *GseConfigStreamToKafka  `json:"kafka,omitempty"`
	Redis       *GseConfigStreamToRedis  `json:"redis,omitempty"`
	Pulsar      *GseConfigStreamToPulsar `json:"pulsar,omitempty"`
}

type GseConfigReportMode string

const (
	GseConfigReportModeKafka  GseConfigReportMode = "kafka"
	GseConfigReportModeRedis  GseConfigReportMode = "redis"
	GseConfigReportModeFile   GseConfigReportMode = "file"
	GseConfigReportModePulsar GseConfigReportMode = "pulsar"
)

type GseConfigStreamToKafka struct {
	StorageAddresses []GseConfigStorageAddress `json:"storage_address"`
	SaslUsername     string                    `json:"sasl_usename,omitempty"`
	SaslPassword     string                    `json:"sasl_passwd,omitempty"`
	SaslMechanisms   string                    `json:"sasl_mechanisms,omitempty"`
	SecurityProtocol string                    `json:"security_protocol,omitempty"`
}

type GseConfigStorageAddress struct {
	IP   string `json:"ip"`
	Port int64  `json:"port"`
}

type GseConfigStreamToRedis struct {
	StorageAddresses []GseConfigStorageAddress `json:"storage_address"`
	Password         string                    `json:"passwd,omitempty"`
	MasterName       string                    `json:"master_name,omitempty"`
}

type GseConfigStreamToPulsar struct {
	StorageAddresses []GseConfigStorageAddress `json:"storage_address"`
	Token            string                    `json:"token,omitempty"`
}

type GseConfigAddStreamToResp struct {
	EsbBaseResponse `json:",inline"`
	Data            *GseConfigAddStreamToResult `json:"data"`
}

type GseConfigAddStreamToResult struct {
	StreamToID int64  `json:"stream_to_id"`
	Name       string `json:"name"`
}

type GseConfigUpdateStreamToParams struct {
	Condition     GseConfigStreamToCondition           `json:"condition"`
	Operation     GseConfigOperation                   `json:"operation"`
	Specification GseConfigUpdateStreamToSpecification `json:"specification"`
}

type GseConfigStreamToCondition struct {
	StreamToID int64             `json:"stream_to_id"`
	PlatName   GseConfigPlatName `json:"plat_name"`
}

type GseConfigUpdateStreamToSpecification struct {
	StreamTo GseConfigStreamTo `json:"stream_to"`
}

type GseConfigDeleteStreamToParams struct {
	Condition GseConfigStreamToCondition `json:"condition"`
	Operation GseConfigOperation         `json:"operation"`
}

type GseConfigQueryStreamToParams struct {
	Condition GseConfigQueryStreamToCondition `json:"condition"`
	Operation GseConfigOperation              `json:"operation"`
}

type GseConfigQueryStreamToCondition struct {
	GseConfigStreamToCondition `json:",inline"`
	Label                      *GseConfigStreamToLabel `json:"label,omitempty"`
}

type GseConfigQueryStreamToResp struct {
	EsbBaseResponse `json:",inline"`
	Data            []GseConfigAddStreamToParams `json:"data"`
}

type GseConfigAddRouteParams struct {
	GseConfigChannel `json:",inline"`
	Operation        GseConfigOperation `json:"operation"`
}

type GseConfigChannel struct {
	Metadata      GseConfigAddRouteMetadata    `json:"metadata"`
	Route         []GseConfigRoute             `json:"route,omitempty"`
	StreamFilters []GseConfigRouteStreamFilter `json:"stream_filters,omitempty"`
}

type GseConfigAddRouteMetadata struct {
	PlatName  GseConfigPlatName    `json:"plat_name"`
	Label     *GseConfigRouteLabel `json:"label,omitempty"`
	ChannelID int64                `json:"channel_id,omitempty"`
}

type GseConfigRouteLabel struct {
	Odm       string `json:"odm,omitempty"`
	BkBizID   int64  `json:"bk_biz_id,omitempty"`
	BkBizName string `json:"bk_biz_name,omitempty"`
}

type GseConfigRoute struct {
	Name          string                 `json:"name"`
	StreamTo      GseConfigRouteStreamTo `json:"stream_to"`
	FilterNameAnd []string               `json:"filter_name_and,omitempty"`
	FilterNameOr  []string               `json:"filter_name_or,omitempty"`
}

type GseConfigRouteStreamTo struct {
	StreamToID int64                 `json:"stream_to_id"`
	Kafka      *GseConfigRouteKafka  `json:"kafka,omitempty"`
	Redis      *GseConfigRouteRedis  `json:"redis,omitempty"`
	Pulsar     *GseConfigRoutePulsar `json:"pulsar,omitempty"`
}

type GseConfigRouteKafka struct {
	TopicName string `json:"topic_name"`
	DataSet   string `json:"data_set,omitempty"`
	BizID     int64  `json:"biz_id,omitempty"`
	Partition int64  `json:"partition,omitempty"`
}

type GseConfigRouteRedis struct {
	ChannelName string `json:"channel_name"`
	DataSet     string `json:"data_set,omitempty"`
	BizID       int64  `json:"biz_id,omitempty"`
}

type GseConfigRoutePulsar struct {
	Name      string `json:"name"`
	Tenant    string `json:"tenant,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

type GseConfigRouteStreamFilter struct {
	Name           string `json:"name"`
	FieldIndex     int64  `json:"field_index"`
	FieldDataType  string `json:"field_data_type"`
	FieldDataValue string `json:"field_data_value"`
	FieldSeparator string `json:"field_separator,omitempty"`
	FieldIn        string `json:"field_in,omitempty"`
}

type GseConfigAddRouteResp struct {
	EsbBaseResponse `json:",inline"`
	Data            *GseConfigAddRouteResult `json:"data"`
}

type GseConfigAddRouteResult struct {
	ChannelID  int64 `json:"channel_id"`
	ProofingID int64 `json:"proofing_id"`
}

type GseConfigUpdateRouteParams struct {
	Condition     GseConfigRouteCondition           `json:"condition"`
	Operation     GseConfigOperation                `json:"operation"`
	Specification GseConfigUpdateRouteSpecification `json:"specification"`
}

type GseConfigRouteCondition struct {
	ChannelID int64                `json:"channel_id"`
	PlatName  GseConfigPlatName    `json:"plat_name"`
	Label     *GseConfigRouteLabel `json:"label,omitempty"`
}

type GseConfigUpdateRouteSpecification struct {
	Route         []GseConfigRoute             `json:"route,omitempty"`
	StreamFilters []GseConfigRouteStreamFilter `json:"stream_filters,omitempty"`
}

type GseConfigDeleteRouteParams struct {
	Condition GseConfigRouteCondition `json:"condition"`
	Operation GseConfigOperation      `json:"operation"`
}

type GseConfigDeleteRouteOperation struct {
	GseConfigOperation `json:",inline"`
	Method             GseConfigDeleteRouteMethod `json:"method"`
}

type GseConfigDeleteRouteMethod string

const (
	All           GseConfigDeleteRouteMethod = "all"
	Specification GseConfigDeleteRouteMethod = "specification"
)

type GseConfigQueryRouteParams struct {
	Condition GseConfigRouteCondition `json:"condition"`
	Operation GseConfigOperation      `json:"operation"`
}

type GseConfigQueryRouteResp struct {
	EsbBaseResponse `json:",inline"`
	Data            []GseConfigChannel `json:"data"`
}
