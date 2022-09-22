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

// GseConfigAddStreamToParams add gse stream_to config params
type GseConfigAddStreamToParams struct {
	Metadata  GseConfigAddStreamToMetadata `json:"metadata"`
	Operation GseConfigOperation           `json:"operation"`
	StreamTo  GseConfigStreamTo            `json:"stream_to"`
}

// GseConfigAddStreamToMetadata add gse stream_to config metadata
type GseConfigAddStreamToMetadata struct {
	PlatName GseConfigPlatName       `json:"plat_name"`
	Label    *GseConfigStreamToLabel `json:"label,omitempty"`
}

// GseConfigPlatName gse plat name
type GseConfigPlatName string

const (
	// GseConfigPlatTglog gse tglog plat config
	GseConfigPlatTglog GseConfigPlatName = "tglog"
	// GseConfigPlatTdm gse tdm plat config
	GseConfigPlatTdm GseConfigPlatName = "tdm"
	// GseConfigPlatBkmonitor gse bkmonitor plat config
	GseConfigPlatBkmonitor GseConfigPlatName = "bkmonitor"
	// GseConfigPlatTgdp gse tgdp plat config
	GseConfigPlatTgdp GseConfigPlatName = "tgdp"
)

// GseConfigStreamToLabel gse stream_to label config
type GseConfigStreamToLabel struct {
	BizID   int64  `json:"bk_biz_id,omitempty"`
	BizName string `json:"bk_biz_name,omitempty"`
}

// GseConfigOperation gse operation config
type GseConfigOperation struct {
	OperatorName string `json:"operator_name"`
}

// GseConfigStreamTo gse stream_to config
type GseConfigStreamTo struct {
	Name        string                   `json:"name"`
	ReportMode  GseConfigReportMode      `json:"report_mode"`
	DataLogPath string                   `json:"data_log_path,omitempty"`
	Kafka       *GseConfigStreamToKafka  `json:"kafka,omitempty"`
	Redis       *GseConfigStreamToRedis  `json:"redis,omitempty"`
	Pulsar      *GseConfigStreamToPulsar `json:"pulsar,omitempty"`
}

// GseConfigReportMode gse report mode config
type GseConfigReportMode string

const (
	// GseConfigReportModeKafka gse kafka report mode
	GseConfigReportModeKafka GseConfigReportMode = "kafka"
	// GseConfigReportModeRedis gse redis report mode
	GseConfigReportModeRedis GseConfigReportMode = "redis"
	// GseConfigReportModeFile gse file report mode
	GseConfigReportModeFile GseConfigReportMode = "file"
	// GseConfigReportModePulsar gse pulsar report mode
	GseConfigReportModePulsar GseConfigReportMode = "pulsar"
)

// GseConfigStreamToKafka gse kafka stream_to config
type GseConfigStreamToKafka struct {
	StorageAddresses []GseConfigStorageAddress `json:"storage_address"`
	SaslUsername     string                    `json:"sasl_usename,omitempty"`
	SaslPassword     string                    `json:"sasl_passwd,omitempty"`
	SaslMechanisms   string                    `json:"sasl_mechanisms,omitempty"`
	SecurityProtocol string                    `json:"security_protocol,omitempty"`
}

// GseConfigStorageAddress gse storage address config
type GseConfigStorageAddress struct {
	IP   string `json:"ip"`
	Port int64  `json:"port"`
}

// GseConfigStreamToRedis gse redis stream_to config
type GseConfigStreamToRedis struct {
	StorageAddresses []GseConfigStorageAddress `json:"storage_address"`
	Password         string                    `json:"passwd,omitempty"`
	MasterName       string                    `json:"master_name,omitempty"`
	SentinelPasswd   string                    `json:"sentinel_passwd,omitempty"`
	Mode             string                    `json:"mode,omitempty"`
}

// GseConfigStreamToPulsar gse pulsar stream_to config
type GseConfigStreamToPulsar struct {
	StorageAddresses []GseConfigStorageAddress `json:"storage_address"`
	Token            string                    `json:"token,omitempty"`
}

// GseConfigAddStreamToResp add gse stream_to config response
type GseConfigAddStreamToResp struct {
	EsbBaseResponse `json:",inline"`
	Data            *GseConfigAddStreamToResult `json:"data"`
}

// GseConfigAddStreamToResult add gse stream_to config result
type GseConfigAddStreamToResult struct {
	StreamToID int64  `json:"stream_to_id"`
	Name       string `json:"name"`
}

// GseConfigUpdateStreamToParams update gse stream_to config params
type GseConfigUpdateStreamToParams struct {
	Condition     GseConfigStreamToCondition           `json:"condition"`
	Operation     GseConfigOperation                   `json:"operation"`
	Specification GseConfigUpdateStreamToSpecification `json:"specification"`
}

// GseConfigStreamToCondition gse stream_to config condition
type GseConfigStreamToCondition struct {
	StreamToID int64             `json:"stream_to_id"`
	PlatName   GseConfigPlatName `json:"plat_name"`
}

// GseConfigUpdateStreamToSpecification gse stream_to update specification config
type GseConfigUpdateStreamToSpecification struct {
	StreamTo GseConfigStreamTo `json:"stream_to"`
}

// GseConfigDeleteStreamToParams delete gse stream_to config params
type GseConfigDeleteStreamToParams struct {
	Condition GseConfigStreamToCondition `json:"condition"`
	Operation GseConfigOperation         `json:"operation"`
}

// GseConfigQueryStreamToParams query gse stream_to config params
type GseConfigQueryStreamToParams struct {
	Condition GseConfigQueryStreamToCondition `json:"condition"`
	Operation GseConfigOperation              `json:"operation"`
}

// GseConfigQueryStreamToCondition query gse stream_to config condition
type GseConfigQueryStreamToCondition struct {
	GseConfigStreamToCondition `json:",inline"`
	Label                      *GseConfigStreamToLabel `json:"label,omitempty"`
}

// GseConfigQueryStreamToResp query gse stream_to config response
type GseConfigQueryStreamToResp struct {
	EsbBaseResponse `json:",inline"`
	Data            []GseConfigAddStreamToParams `json:"data"`
}

// GseConfigAddRouteParams add gse route config params
type GseConfigAddRouteParams struct {
	GseConfigChannel `json:",inline"`
	Operation        GseConfigOperation `json:"operation"`
}

// GseConfigChannel gse channel config
type GseConfigChannel struct {
	Metadata      GseConfigAddRouteMetadata    `json:"metadata"`
	Route         []GseConfigRoute             `json:"route,omitempty"`
	StreamFilters []GseConfigRouteStreamFilter `json:"stream_filters,omitempty"`
}

// GseConfigAddRouteMetadata add gse route config metadata
type GseConfigAddRouteMetadata struct {
	PlatName  GseConfigPlatName    `json:"plat_name"`
	Label     *GseConfigRouteLabel `json:"label,omitempty"`
	ChannelID int64                `json:"channel_id,omitempty"`
}

// GseConfigRouteLabel gse route label config
type GseConfigRouteLabel struct {
	Odm       string `json:"odm,omitempty"`
	BkBizID   int64  `json:"bk_biz_id,omitempty"`
	BkBizName string `json:"bk_biz_name,omitempty"`
}

// GseConfigRoute gse route config
type GseConfigRoute struct {
	Name          string                 `json:"name"`
	StreamTo      GseConfigRouteStreamTo `json:"stream_to"`
	FilterNameAnd []string               `json:"filter_name_and,omitempty"`
	FilterNameOr  []string               `json:"filter_name_or,omitempty"`
}

// GseConfigRouteStreamTo gse stream_to route config
type GseConfigRouteStreamTo struct {
	StreamToID int64                 `json:"stream_to_id"`
	Kafka      *GseConfigRouteKafka  `json:"kafka,omitempty"`
	Redis      *GseConfigRouteRedis  `json:"redis,omitempty"`
	Pulsar     *GseConfigRoutePulsar `json:"pulsar,omitempty"`
}

// GseConfigRouteKafka gse kafka route config
type GseConfigRouteKafka struct {
	TopicName string `json:"topic_name"`
	DataSet   string `json:"data_set,omitempty"`
	BizID     int64  `json:"biz_id,omitempty"`
	Partition int64  `json:"partition,omitempty"`
}

// GseConfigRouteRedis gse redis route config
type GseConfigRouteRedis struct {
	ChannelName string `json:"channel_name"`
	DataSet     string `json:"data_set,omitempty"`
	BizID       int64  `json:"biz_id,omitempty"`
}

// GseConfigRoutePulsar gse pulsar route config
type GseConfigRoutePulsar struct {
	Name      string `json:"name"`
	Tenant    string `json:"tenant,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

// GseConfigRouteStreamFilter gse stream route filter config
type GseConfigRouteStreamFilter struct {
	Name           string `json:"name"`
	FieldIndex     int64  `json:"field_index"`
	FieldDataType  string `json:"field_data_type"`
	FieldDataValue string `json:"field_data_value"`
	FieldSeparator string `json:"field_separator,omitempty"`
	FieldIn        string `json:"field_in,omitempty"`
}

// GseConfigAddRouteResp add gse route config response
type GseConfigAddRouteResp struct {
	EsbBaseResponse `json:",inline"`
	Data            *GseConfigAddRouteResult `json:"data"`
}

// GseConfigAddRouteResult add gse route config result
type GseConfigAddRouteResult struct {
	ChannelID  int64 `json:"channel_id"`
	ProofingID int64 `json:"proofing_id"`
}

// GseConfigUpdateRouteParams update gse route config params
type GseConfigUpdateRouteParams struct {
	Condition     GseConfigRouteCondition           `json:"condition"`
	Operation     GseConfigOperation                `json:"operation"`
	Specification GseConfigUpdateRouteSpecification `json:"specification"`
}

// GseConfigRouteCondition gse route config condition
type GseConfigRouteCondition struct {
	ChannelID int64                `json:"channel_id"`
	PlatName  GseConfigPlatName    `json:"plat_name"`
	Label     *GseConfigRouteLabel `json:"label,omitempty"`
}

// GseConfigUpdateRouteSpecification update gse route specification config
type GseConfigUpdateRouteSpecification struct {
	Route         []GseConfigRoute             `json:"route,omitempty"`
	StreamFilters []GseConfigRouteStreamFilter `json:"stream_filters,omitempty"`
}

// GseConfigDeleteRouteParams delete gse route config params
type GseConfigDeleteRouteParams struct {
	Condition GseConfigRouteCondition `json:"condition"`
	Operation GseConfigOperation      `json:"operation"`
}

// GseConfigDeleteRouteOperation delete gse route operation config
type GseConfigDeleteRouteOperation struct {
	GseConfigOperation `json:",inline"`
	Method             GseConfigDeleteRouteMethod `json:"method"`
}

// GseConfigDeleteRouteMethod delete gse route config method
type GseConfigDeleteRouteMethod string

const (
	// All delete all gse route config method
	All GseConfigDeleteRouteMethod = "all"
	// Specification delete specification gse route config method
	Specification GseConfigDeleteRouteMethod = "specification"
)

// GseConfigQueryRouteParams gse query route params
type GseConfigQueryRouteParams struct {
	Condition GseConfigRouteCondition `json:"condition"`
	Operation GseConfigOperation      `json:"operation"`
}

// GseConfigQueryRouteResp gse query route response
type GseConfigQueryRouteResp struct {
	EsbBaseResponse `json:",inline"`
	Data            []GseConfigChannel `json:"data"`
}
