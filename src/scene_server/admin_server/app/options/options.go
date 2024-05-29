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

// Package options TODO
package options

import (
	"configcenter/src/ac/iam"
	"configcenter/src/common/auth"
	"configcenter/src/common/core/cc/config"
	"configcenter/src/storage/dal/kafka"
	"configcenter/src/storage/dal/mongo"
	"configcenter/src/storage/dal/redis"

	"github.com/spf13/pflag"
)

// ServerOption define option of server in flags
type ServerOption struct {
	ServConf *config.CCAPIConfig
}

// NewServerOption create a ServerOption object
func NewServerOption() *ServerOption {
	s := ServerOption{
		ServConf: config.NewCCAPIConfig(),
	}

	return &s
}

// AddFlags add flags
func (s *ServerOption) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.ServConf.AddrPort, "addrport", "127.0.0.1:60005", "The ip address and port for the serve on")
	fs.StringVar(&s.ServConf.ExConfig, "config", "conf/api.conf", "The config path. e.g conf/api.conf")
	fs.StringVar(&s.ServConf.RegisterIP, "register-ip", "", "the ip address registered on zookeeper, it can be domain")
	fs.Var(auth.EnableAuthFlag, "enable-auth", "The auth center enable status, true for enabled, false for disabled")
}

// Config TODO
type Config struct {
	MongoDB        mongo.Config
	WatchDB        mongo.Config
	Errors         ErrorConfig
	Language       LanguageConfig
	Configures     ConfConfig
	Register       RegisterConfig
	Redis          redis.Config
	SnapRedis      redis.Config
	SnapKafka      kafka.Config
	IAM            iam.AuthConfig
	SnapDataID     int64
	SnapReportMode string
	ShardingTable  ShardingTableConfig
	// SyncIAMPeriodMinutes the period for sync IAM resources
	SyncIAMPeriodMinutes int
	// 通过何种方式调用gse接口注册dataid
	DataIdMigrateWay MigrateWay
}

// MigrateWay 通过何种方式调用gse接口注册dataid
type MigrateWay string

const (
	// MigrateWayESB 通过esb调用gse接口
	MigrateWayESB MigrateWay = "esb"
	// MigrateWayApiGW 通过api gateway调用gse接口
	MigrateWayApiGW MigrateWay = "apigw"
)

// LanguageConfig TODO
type LanguageConfig struct {
	Res string
}

// ErrorConfig TODO
type ErrorConfig struct {
	Res string
}

// ConfConfig TODO
type ConfConfig struct {
	Dir string
}

// RegisterConfig TODO
type RegisterConfig struct {
	Address string
}

// ShardingTableConfig TODO
type ShardingTableConfig struct {
	// 表中同步索引间隔时间，单位分钟， 最小30分钟， 默认60分钟， 最大720分钟
	IndexesInterval int64
	// 模型shardingTable 对比和处理， 单位秒， 最小60秒，默认 120秒， 最大1800s
	TableInterval int64
}
