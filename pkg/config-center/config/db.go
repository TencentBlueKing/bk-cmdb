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

// DBType ...
type DBType string

const (
	// PostgreSQLType ...
	PostgreSQLType DBType = "postgres"
	// MySQLType ...
	MySQLType DBType = "mysql"
)

// DBConfig 数据库配置
type DBConfig struct {
	DBType DBType `yaml:"db_type" json:"db_type"`

	/* connection params */

	Host     string     `yaml:"host" json:"host"`
	Port     uint16     `yaml:"port" json:"port"`
	Database string     `yaml:"database" json:"database"`
	Username string     `yaml:"username" json:"username"`
	Password string     `yaml:"password" json:"password"`
	TLS      *TLSConfig `yaml:"tls" json:"tls,omitempty"`

	/* sql.DB connection pool settings */

	MaxIdleConns   uint `yaml:"max_idle_conns" json:"max_idle_conns"`
	MaxOpenConns   uint `yaml:"max_open_conns" json:"max_open_conns"`
	ConnMaxLifeSec uint `yaml:"conn_max_life_sec" json:"conn_max_life_sec,format:units"`
	ConnMaxIdleSec uint `yaml:"conn_max_idle_sec" json:"conn_max_idle_sec,format:units"`

	/* gorm settings */

	// Debug enable gorm debug mode, will print sql to log
	Debug bool `yaml:"debug" json:"debug"`

	SlowLogThresholdMS uint `yaml:"slow_log_threshold_ms" json:"slow_log_threshold_ms,format:units"`
}
