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

// GseConfigAddStreamToParams TODO
type GseConfigAddStreamToParams struct {
	Metadata  GseConfigAddStreamToMetadata `json:"metadata"`
	Operation GseConfigOperation           `json:"operation"`
	StreamTo  GseConfigStreamTo            `json:"stream_to"`
}

// GseConfigAddStreamToMetadata TODO
type GseConfigAddStreamToMetadata struct {
	PlatName GseConfigPlatName       `json:"plat_name"`
	Label    *GseConfigStreamToLabel `json:"label,omitempty"`
}

// GseConfigPlatName TODO
type GseConfigPlatName string

const (
	// GseConfigPlatTglog TODO
	GseConfigPlatTglog GseConfigPlatName = "tglog"
	// GseConfigPlatTdm TODO
	GseConfigPlatTdm GseConfigPlatName = "tdm"
	// GseConfigPlatBkmonitor TODO
	GseConfigPlatBkmonitor GseConfigPlatName = "bkmonitor"
	// GseConfigPlatTgdp TODO
	GseConfigPlatTgdp GseConfigPlatName = "tgdp"
)

// GseConfigStreamToLabel TODO
type GseConfigStreamToLabel struct {
	BizID   int64  `json:"bk_biz_id,omitempty"`
	BizName string `json:"bk_biz_name,omitempty"`
}

// GseConfigOperation TODO
type GseConfigOperation struct {
	OperatorName string `json:"operator_name"`
}

// GseConfigStreamTo TODO
type GseConfigStreamTo struct {
	Name        string                   `json:"name"`
	ReportMode  GseConfigReportMode      `json:"report_mode"`
	DataLogPath string                   `json:"data_log_path,omitempty"`
	Kafka       *GseConfigStreamToKafka  `json:"kafka,omitempty"`
	Redis       *GseConfigStreamToRedis  `json:"redis,omitempty"`
	Pulsar      *GseConfigStreamToPulsar `json:"pulsar,omitempty"`
}

// GseConfigReportMode TODO
type GseConfigReportMode string

const (
	// GseConfigReportModeKafka TODO
	GseConfigReportModeKafka GseConfigReportMode = "kafka"
	// GseConfigReportModeRedis TODO
	GseConfigReportModeRedis GseConfigReportMode = "redis"
	// GseConfigReportModeFile TODO
	GseConfigReportModeFile GseConfigReportMode = "file"
	// GseConfigReportModePulsar TODO
	GseConfigReportModePulsar GseConfigReportMode = "pulsar"
)

// GseConfigStreamToKafka TODO
type GseConfigStreamToKafka struct {
	StorageAddresses []GseConfigStorageAddress `json:"storage_address"`
	SaslUsername     string                    `json:"sasl_usename,omitempty"`
	SaslPassword     string                    `json:"sasl_passwd,omitempty"`
	SaslMechanisms   string                    `json:"sasl_mechanisms,omitempty"`
	SecurityProtocol string                    `json:"security_protocol,omitempty"`
}

// GseConfigStorageAddress TODO
type GseConfigStorageAddress struct {
	IP   string `json:"ip"`
	Port int64  `json:"port"`
}

// GseConfigStreamToRedis TODO
type GseConfigStreamToRedis struct {
	StorageAddresses []GseConfigStorageAddress `json:"storage_address"`
	Password         string                    `json:"passwd,omitempty"`
	MasterName       string                    `json:"master_name,omitempty"`
	Mode             string                    `json:"mode,omitempty"`
}

// GseConfigStreamToPulsar TODO
type GseConfigStreamToPulsar struct {
	StorageAddresses []GseConfigStorageAddress `json:"storage_address"`
	Token            string                    `json:"token,omitempty"`
}

// GseConfigAddStreamToResp TODO
type GseConfigAddStreamToResp struct {
	EsbBaseResponse `json:",inline"`
	Data            *GseConfigAddStreamToResult `json:"data"`
}

// GseConfigAddStreamToResult TODO
type GseConfigAddStreamToResult struct {
	StreamToID int64  `json:"stream_to_id"`
	Name       string `json:"name"`
}

// GseConfigUpdateStreamToParams TODO
type GseConfigUpdateStreamToParams struct {
	Condition     GseConfigStreamToCondition           `json:"condition"`
	Operation     GseConfigOperation                   `json:"operation"`
	Specification GseConfigUpdateStreamToSpecification `json:"specification"`
}

// GseConfigStreamToCondition TODO
type GseConfigStreamToCondition struct {
	StreamToID int64             `json:"stream_to_id"`
	PlatName   GseConfigPlatName `json:"plat_name"`
}

// GseConfigUpdateStreamToSpecification TODO
type GseConfigUpdateStreamToSpecification struct {
	StreamTo GseConfigStreamTo `json:"stream_to"`
}

// GseConfigDeleteStreamToParams TODO
type GseConfigDeleteStreamToParams struct {
	Condition GseConfigStreamToCondition `json:"condition"`
	Operation GseConfigOperation         `json:"operation"`
}

// GseConfigQueryStreamToParams TODO
type GseConfigQueryStreamToParams struct {
	Condition GseConfigQueryStreamToCondition `json:"condition"`
	Operation GseConfigOperation              `json:"operation"`
}

// GseConfigQueryStreamToCondition TODO
type GseConfigQueryStreamToCondition struct {
	GseConfigStreamToCondition `json:",inline"`
	Label                      *GseConfigStreamToLabel `json:"label,omitempty"`
}

// GseConfigQueryStreamToResp TODO
type GseConfigQueryStreamToResp struct {
	EsbBaseResponse `json:",inline"`
	Data            []GseConfigAddStreamToParams `json:"data"`
}

// GseConfigAddRouteParams TODO
type GseConfigAddRouteParams struct {
	GseConfigChannel `json:",inline"`
	Operation        GseConfigOperation `json:"operation"`
}

// GseConfigChannel TODO
type GseConfigChannel struct {
	Metadata      GseConfigAddRouteMetadata    `json:"metadata"`
	Route         []GseConfigRoute             `json:"route,omitempty"`
	StreamFilters []GseConfigRouteStreamFilter `json:"stream_filters,omitempty"`
}

// GseConfigAddRouteMetadata TODO
type GseConfigAddRouteMetadata struct {
	PlatName  GseConfigPlatName    `json:"plat_name"`
	Label     *GseConfigRouteLabel `json:"label,omitempty"`
	ChannelID int64                `json:"channel_id,omitempty"`
}

// GseConfigRouteLabel TODO
type GseConfigRouteLabel struct {
	Odm       string `json:"odm,omitempty"`
	BkBizID   int64  `json:"bk_biz_id,omitempty"`
	BkBizName string `json:"bk_biz_name,omitempty"`
}

// GseConfigRoute TODO
type GseConfigRoute struct {
	Name          string                 `json:"name"`
	StreamTo      GseConfigRouteStreamTo `json:"stream_to"`
	FilterNameAnd []string               `json:"filter_name_and,omitempty"`
	FilterNameOr  []string               `json:"filter_name_or,omitempty"`
}

// GseConfigRouteStreamTo TODO
type GseConfigRouteStreamTo struct {
	StreamToID int64                 `json:"stream_to_id"`
	Kafka      *GseConfigRouteKafka  `json:"kafka,omitempty"`
	Redis      *GseConfigRouteRedis  `json:"redis,omitempty"`
	Pulsar     *GseConfigRoutePulsar `json:"pulsar,omitempty"`
}

// GseConfigRouteKafka TODO
type GseConfigRouteKafka struct {
	TopicName string `json:"topic_name"`
	DataSet   string `json:"data_set,omitempty"`
	BizID     int64  `json:"biz_id,omitempty"`
	Partition int64  `json:"partition,omitempty"`
}

// GseConfigRouteRedis TODO
type GseConfigRouteRedis struct {
	ChannelName string `json:"channel_name"`
	DataSet     string `json:"data_set,omitempty"`
	BizID       int64  `json:"biz_id,omitempty"`
}

// GseConfigRoutePulsar TODO
type GseConfigRoutePulsar struct {
	Name      string `json:"name"`
	Tenant    string `json:"tenant,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

// GseConfigRouteStreamFilter TODO
type GseConfigRouteStreamFilter struct {
	Name           string `json:"name"`
	FieldIndex     int64  `json:"field_index"`
	FieldDataType  string `json:"field_data_type"`
	FieldDataValue string `json:"field_data_value"`
	FieldSeparator string `json:"field_separator,omitempty"`
	FieldIn        string `json:"field_in,omitempty"`
}

// GseConfigAddRouteResp TODO
type GseConfigAddRouteResp struct {
	EsbBaseResponse `json:",inline"`
	Data            *GseConfigAddRouteResult `json:"data"`
}

// GseConfigAddRouteResult TODO
type GseConfigAddRouteResult struct {
	ChannelID  int64 `json:"channel_id"`
	ProofingID int64 `json:"proofing_id"`
}

// GseConfigUpdateRouteParams TODO
type GseConfigUpdateRouteParams struct {
	Condition     GseConfigRouteCondition           `json:"condition"`
	Operation     GseConfigOperation                `json:"operation"`
	Specification GseConfigUpdateRouteSpecification `json:"specification"`
}

// GseConfigRouteCondition TODO
type GseConfigRouteCondition struct {
	ChannelID int64                `json:"channel_id"`
	PlatName  GseConfigPlatName    `json:"plat_name"`
	Label     *GseConfigRouteLabel `json:"label,omitempty"`
}

// GseConfigUpdateRouteSpecification TODO
type GseConfigUpdateRouteSpecification struct {
	Route         []GseConfigRoute             `json:"route,omitempty"`
	StreamFilters []GseConfigRouteStreamFilter `json:"stream_filters,omitempty"`
}

// GseConfigDeleteRouteParams TODO
type GseConfigDeleteRouteParams struct {
	Condition GseConfigRouteCondition `json:"condition"`
	Operation GseConfigOperation      `json:"operation"`
}

// GseConfigDeleteRouteOperation TODO
type GseConfigDeleteRouteOperation struct {
	GseConfigOperation `json:",inline"`
	Method             GseConfigDeleteRouteMethod `json:"method"`
}

// GseConfigDeleteRouteMethod TODO
type GseConfigDeleteRouteMethod string

const (
	// All TODO
	All GseConfigDeleteRouteMethod = "all"
	// Specification TODO
	Specification GseConfigDeleteRouteMethod = "specification"
)

// GseConfigQueryRouteParams TODO
type GseConfigQueryRouteParams struct {
	Condition GseConfigRouteCondition `json:"condition"`
	Operation GseConfigOperation      `json:"operation"`
}

// GseConfigQueryRouteResp TODO
type GseConfigQueryRouteResp struct {
	EsbBaseResponse `json:",inline"`
	Data            []GseConfigChannel `json:"data"`
}
