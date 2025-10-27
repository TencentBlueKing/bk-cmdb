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
	"fmt"
	"os"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
)

// GetTestGORM open a gorm db for test from env TEST_GORM_DSN, should only be used in test code
func GetTestGORM(t testing.TB) (*gorm.DB, error) {
	var gormDSN = os.Getenv("TEST_GORM_PG_DSN")
	if gormDSN == "" {
		log.Warn(t.Context(), "TEST_GORM_PG_DSN is not set, skip", "test", t.Name())
		t.Skip("TEST_GORM_PG_DSN is not set")
		return nil, fmt.Errorf("TEST_GORM_PG_DSN is not set")
	}
	gormDB, err := gorm.Open(
		postgres.Open(gormDSN),
		&gorm.Config{
			NowFunc:     func() time.Time { return time.Now().UTC().Round(time.Microsecond) },
			PrepareStmt: true,
		})
	if err != nil {
		return nil, fmt.Errorf("failed to open gorm db: %w", err)
	}
	return gormDB.Debug(), nil
}
