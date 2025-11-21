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

package database

import (
	"database/sql"
	"fmt"
	"net"
	"strconv"

	mysqldriver "github.com/go-sql-driver/mysql"

	"github.com/TencentBlueKing/bk-cmdb/pkg/config-center/config"
)

// MySQLOption mysql连接选项
type MySQLOption func(cfg *mysqldriver.Config)

// OpenMySQL open a mysql connection
func OpenMySQL(cfg *config.DBConfig, options ...MySQLOption) (*sql.DB, error) {
	mysqlCfg, err := buildMySQLConfig(cfg, options)
	if err != nil {
		return nil, err
	}

	connector, err := mysqldriver.NewConnector(mysqlCfg)
	if err != nil {
		return nil, fmt.Errorf("fail to create mysql connector: %w", err)
	}
	conn := sql.OpenDB(connector)
	return conn, nil
}

func buildMySQLConfig(cfg *config.DBConfig, options []MySQLOption) (*mysqldriver.Config, error) {
	mysqlCfg := mysqldriver.NewConfig()
	mysqlCfg.Addr = cfg.Host
	if cfg.Port > 0 {
		mysqlCfg.Addr = net.JoinHostPort(cfg.Host, strconv.Itoa(int(cfg.Port)))
	}
	mysqlCfg.User = cfg.Username
	mysqlCfg.Passwd = cfg.Password
	mysqlCfg.DBName = cfg.Database
	_ = mysqlCfg.Apply(mysqldriver.Charset("utf8mb4", ""))

	tlsCfg, tlsEnabled, err := cfg.TLS.ToClientConf()
	if err != nil {
		return nil, fmt.Errorf("fail to get tls config: %w", err)
	}
	if tlsEnabled {
		mysqlCfg.TLS = tlsCfg
	}
	for _, opt := range options {
		opt(mysqlCfg)
	}
	return mysqlCfg, err
}
