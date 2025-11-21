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

// Package tests common functions for test
package tests

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/TencentBlueKing/bk-cmdb/pkg/config-center/config"
)

// TestGORMDSNMySQL env key for test gorm mysql dsn
const TestGORMDSNMySQL = "TEST_GORM_DSN_MYSQL"

// TestGORMDSNPG env key for test gorm pg dsn
const TestGORMDSNPG = "TEST_GORM_DSN_PG"

// GetRealDB open a real gorm db for test from env TEST_GORM_DSN_PG or TEST_GORM_DSN_MYSQL,
// if both TEST_GORM_DSN_PG and TEST_GORM_DSN_MYSQL are not set, it will skip the test
func GetRealDB(t testing.TB) (*gorm.DB, error) {
	if os.Getenv(TestGORMDSNMySQL) != "" {
		return GetRealMySQL(t)
	}
	return GetRealPG(t)
}

// GetRealMySQL open a real gorm db for test from env TEST_GORM_DSN_MYSQL,
// if TEST_GORM_DSN_MYSQL is not set, it will skip the test
func GetRealMySQL(t testing.TB) (*gorm.DB, error) {
	var mysqlDsn = os.Getenv(TestGORMDSNMySQL)
	if mysqlDsn == "" {
		t.Skip("TEST_GORM_DSN_MYSQL is not set")
		return nil, fmt.Errorf("TEST_GORM_DSN_MYSQL is not set")
	}
	return getGORM(mysql.Open(mysqlDsn))
}

// GetRealPG open a real gorm db for test from env TEST_GORM_DSN_PG or TEST_GORM_DSN_MYSQL,
// if TEST_GORM_DSN_PG is not set, it will skip the test
func GetRealPG(t testing.TB) (*gorm.DB, error) {
	var pgDSN = os.Getenv(TestGORMDSNPG)
	if pgDSN == "" {
		t.Skip("TEST_GORM_DSN_PG is not set")
		return nil, fmt.Errorf("TEST_GORM_DSN_PG is not set")
	}
	return getGORM(postgres.Open(pgDSN))

}

// GetMockPG open a mock pg gorm db for test from sqlmock
func GetMockPG(t testing.TB) (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	gormDB, err := getGORM(getMockGormPG(db))
	return gormDB, mock, err
}

// GetMockMySQL open a mock mysql gorm db for test from sqlmock
func GetMockMySQL(t testing.TB) (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	gormDB, err := getGORM(getMockGormMySQL(db))
	return gormDB, mock, err
}

func getGORM(dialector gorm.Dialector) (*gorm.DB, error) {
	gormDB, err := gorm.Open(
		dialector,
		&gorm.Config{
			NowFunc:     func() time.Time { return time.Now().UTC().Round(time.Microsecond) },
			PrepareStmt: false,
		})
	if err != nil {
		return nil, fmt.Errorf("failed to open gorm db with dialector %v, err: %w", dialector.Name(), err)
	}
	return gormDB.Debug(), nil
}

func getMockGormMySQL(db *sql.DB) gorm.Dialector {
	return mysql.New(mysql.Config{
		DriverName: "mysql",
		Conn:       db,
	})
}

func getMockGormPG(db *sql.DB) gorm.Dialector {
	return postgres.New(postgres.Config{
		DriverName: "postgres",
		Conn:       db,
	})
}

const (
	// TestGormSSLPgConfig is the env name for ssl pg config
	TestGormSSLPgConfig = "TEST_GORM_SSL_PG_CONFIG"
	// TestGormSSLMySQLConfig is the env name for ssl mysql config
	TestGormSSLMySQLConfig = "TEST_GORM_SSL_MYSQL_CONFIG"
)

// GetTestPGSSLConfig get the test postgresql ssl config
func GetTestPGSSLConfig(t testing.TB) (*config.DBConfig, error) {
	return GetTestDBSSLConfig(t, TestGormSSLPgConfig)
}

// GetTestMySQLSSLConfig get the test mysql ssl config
func GetTestMySQLSSLConfig(t testing.TB) (*config.DBConfig, error) {
	return GetTestDBSSLConfig(t, TestGormSSLMySQLConfig)
}

// GetTestDBSSLConfig get the test db ssl config
func GetTestDBSSLConfig(t testing.TB, configKey string) (*config.DBConfig, error) {
	sslConfigJSON := os.Getenv(configKey)
	if sslConfigJSON == "" {
		t.Skipf("%s is not set", configKey)
		return nil, fmt.Errorf("%s is not set", configKey)
	}
	var sslConfig = new(config.DBConfig)
	err := json.Unmarshal([]byte(sslConfigJSON), sslConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s, err: %w", configKey, err)
	}
	return sslConfig, nil
}
