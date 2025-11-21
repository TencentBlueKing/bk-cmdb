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

	mysqldriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/TencentBlueKing/bk-cmdb/pkg/config-center/config"
)

// NewGORMClient 创建新的GORM客户端
func NewGORMClient(cfg *config.DBConfig) (db *gorm.DB, err error) {
	var dialector gorm.Dialector
	switch cfg.DBType {
	case config.PostgreSQLType:
		dsn := buildPostgreSQLDSN(cfg)
		dialector = postgres.Open(dsn)
	case config.MySQLType:
		tlsCfg, tlsEnabled, err := cfg.TLS.ToClientConf()
		if err != nil {
			return nil, fmt.Errorf("fail to get tls config: %w", err)
		}
		if tlsEnabled {
			if err = mysqldriver.RegisterTLSConfig("cmdb", tlsCfg); err != nil {
				return nil, fmt.Errorf("fail to register mysql tls config: %w", err)
			}
		}
		dsn := buildMySQLDSN(cfg, tlsEnabled)
		mysqlCfg := mysql.Config{DSN: dsn}
		dialector = mysql.New(mysqlCfg)
	default:
		return nil, fmt.Errorf("unsupported db type: %s", cfg.DBType)
	}
	gormCfg := &gorm.Config{
		SkipDefaultTransaction:                   false,
		DefaultTransactionTimeout:                0,
		DefaultContextTimeout:                    cfg.DefaultConnectionTimeout,
		NamingStrategy:                           nil,
		NowFunc:                                  nil,
		PrepareStmt:                              false,
		PrepareStmtMaxSize:                       0,
		PrepareStmtTTL:                           0,
		DisableAutomaticPing:                     false,
		DisableForeignKeyConstraintWhenMigrating: false,
		IgnoreRelationshipsWhenMigrating:         false,
		DisableNestedTransaction:                 false,
		AllowGlobalUpdate:                        false,
		QueryFields:                              false,
		CreateBatchSize:                          0,
		TranslateError:                           false,
	}
	db, err = gorm.Open(dialector, gormCfg)
	if err != nil {
		return nil, fmt.Errorf("open gorm failed, err: %w", err)
	}

	// 设置连接池参数
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("fail to get sql.DB during init gorm database, err: %w", err)
	}
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

	return db, nil
}

func buildPostgreSQLDSN(cfg *config.DBConfig) string {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s ",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database)
	if cfg.TLS == nil {
		dsn += "sslmode=disable"
	} else {
		if cfg.TLS.InsecureSkipVerify {
			dsn += "sslmode=allow "
		} else {
			dsn += "sslmode=verify-ca "
		}
		if cfg.TLS.CAFile != "" {
			dsn += fmt.Sprintf("sslrootcert=%s ", cfg.TLS.CAFile)
		}
		// cert and key should be set together
		if cfg.TLS.CertFile != "" {
			dsn += fmt.Sprintf("sslkey=%s sslcert=%s ", cfg.TLS.KeyFile, cfg.TLS.CertFile)
		}
	}
	return dsn
}

// buildMySQLDSN 构建MySQL的DSN
func buildMySQLDSN(cfg *config.DBConfig, tlsEnabled bool) string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	if tlsEnabled {
		dsn += "&tls=cmdb"
	}

	return dsn
}
