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
	"testing"

	"github.com/TencentBlueKing/bk-cmdb/pkg/config-center/config"
	"github.com/TencentBlueKing/bk-cmdb/pkg/tests"
)

func TestNewGORMClient(t *testing.T) {
	t.Run("pgsql", func(t *testing.T) {
		cfg, err := tests.GetTestPGSSLConfig(t)
		if err != nil {
			t.Errorf("GetTestPGSSLConfig() error = %v", err)
			return
		}
		testNewGorm(t, cfg)
	})
	t.Run("mysql", func(t *testing.T) {
		cfg, err := tests.GetTestMySQLSSLConfig(t)
		if err != nil {
			t.Errorf("GetTestMySQLSSLConfig() error = %v", err)
			return
		}
		testNewGorm(t, cfg)
	})
}
func testNewGorm(t *testing.T, config *config.DBConfig) {
	gotDb, err := NewGORMClient(config)
	if err != nil {
		t.Errorf("NewGORMClient() error = %v, wantErr %v", err, false)
		return
	}
	db, err := gotDb.DB()
	if err != nil {
		t.Errorf("get DB() error = %v, config: %s", err, config.DBType)
		return
	}
	if err := db.Ping(); err != nil {
		t.Errorf("ping DB() error = %v", err)
		return
	}

}
