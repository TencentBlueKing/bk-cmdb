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

// Package database ...
package database

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/TencentBlueKing/bk-cmdb/pkg/config-center/config"
)

const defaultContextTimeout = time.Second * 20

// NewGORMClient 创建新的GORM客户端
func NewGORMClient(cfg *config.DBConfig) (db *gorm.DB, err error) {
	var dialector gorm.Dialector
	switch cfg.DBType {
	case config.PostgreSQLType:
		pgConn, err := OpenStdlibPGX(cfg)
		if err != nil {
			return nil, fmt.Errorf("fail to create pgx std client for new gorm client: %w", err)
		}
		pgCfg := postgres.Config{Conn: pgConn}
		dialector = postgres.New(pgCfg)
	case config.MySQLType:
		conn, err := OpenMySQL(cfg)
		if err != nil {
			return nil, fmt.Errorf("fail to create mysql client for new gorm client: %w", err)
		}
		mysqlCfg := mysql.Config{Conn: conn}
		dialector = mysql.New(mysqlCfg)
	default:
		return nil, fmt.Errorf("unsupported db type: %s", cfg.DBType)
	}
	gormCfg := &gorm.Config{DefaultContextTimeout: defaultContextTimeout}
	db, err = gorm.Open(dialector, gormCfg)
	if err != nil {
		return nil, fmt.Errorf("open gorm failed, err: %w", err)
	}

	// 设置连接池参数
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("fail to get sql.DB during init gorm database, err: %w", err)
	}
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.ConnMaxIdleSec) * time.Second)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifeSec) * time.Second)
	sqlDB.SetMaxIdleConns(int(cfg.MaxIdleConns))
	sqlDB.SetMaxOpenConns(int(cfg.MaxOpenConns))

	return db, nil
}
