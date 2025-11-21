/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - CMDB) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package config

import "time"

// DBType ...
type DBType string

const (
	// PostgreSQLType ...
	PostgreSQLType = "postgres"
	// MySQLType ...
	MySQLType = "mysql"
)

// DBConfig 数据库配置
type DBConfig struct {
	DBType string `yaml:"db_type" json:"db_type"`

	// connection params
	Host     string     `yaml:"host" json:"host"`
	Port     int        `yaml:"port" json:"port"`
	Database string     `yaml:"database" json:"database"`
	Username string     `yaml:"username" json:"username"`
	Password string     `yaml:"password" json:"password"`
	TLS      *TLSConfig `yaml:"tls" json:"tls,omitempty"`

	// sql.DB connection pool settings

	MaxIdleConns    int           `yaml:"max_idle_conns" json:"max_idle_conns"`
	MaxOpenConns    int           `yaml:"max_open_conns" json:"max_open_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" json:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time" json:"conn_max_idle_time"`

	/* gorm settings */

	// if > 0, all query will add this default timeout to context if given context timeout not specified
	DefaultConnectionTimeout time.Duration `yaml:"default_connection_timeout" json:"default_connection_timeout"`
	Debug                    bool          `yaml:"debug" json:"debug"`

	SlowLogThreshold time.Duration       `yaml:"slow_log_threshold" json:"slow_log_threshold"`
	IngressLimit     *TokenBucketLimiter `yaml:"ingress_limit" json:"ingress_limit"`
}

// TokenBucketLimiter 令牌桶限流配置
type TokenBucketLimiter struct {
	Bucket  uint `yaml:"bucket" json:"bucket"`
	RateQPS uint `yaml:"rate_qps" json:"rate_qps"`
}
