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

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/TencentBlueKing/bk-cmdb/pkg/config-center/config"
)

// PGXOption pgx连接选项
type PGXOption func(connConfig *pgx.ConnConfig)

// OpenStdlibPGX open a stdlib version pgx connection
func OpenStdlibPGX(cfg *config.DBConfig, options ...PGXOption) (*sql.DB, error) {
	pgCfg, err := buildPGXConfig(cfg, options...)
	if err != nil {
		return nil, fmt.Errorf("fail to build pgx config: %w", err)
	}
	db := stdlib.OpenDB(*pgCfg)
	return db, nil
}

func buildPGXConfig(cfg *config.DBConfig, options ...PGXOption) (*pgx.ConnConfig, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s", cfg.Host, cfg.Username, cfg.Password, cfg.Database)
	// pgx requires that connConfig must be created by pgx.ParseConfig.
	pgCfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	// port could be zero for unix domain socket directory
	pgCfg.Port = cfg.Port
	if cfg.TLS != nil {
		tlsCfg, enabled, err := cfg.TLS.ToClientConf()
		if err != nil {
			return nil, err
		}
		if enabled {
			pgCfg.TLSConfig = tlsCfg
		}
	}
	for _, opt := range options {
		opt(pgCfg)
	}
	return pgCfg, nil
}
