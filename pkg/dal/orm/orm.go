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

// Package orm ...
package orm

import (
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// IsDuplicatedError is db duplicated error, will check:
// 1. Gorm ErrDuplicatedKey
// 2. PostgreSQL SQLSTATE 23505
// 3. MySQL Error Number 1062
func IsDuplicatedError(err error) bool {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}
	var pgErr *pgconn.PgError
	const PGStateUniqueViolation = "23505"
	if errors.As(err, &pgErr) && pgErr.Code == PGStateUniqueViolation {
		return true
	}
	var mysqlErr *mysql.MySQLError
	const MySQLNumberDuplicateEntry = 1062
	if errors.As(err, &mysqlErr) && mysqlErr.Number == MySQLNumberDuplicateEntry {
		return true
	}
	return false
}
