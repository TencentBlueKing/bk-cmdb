/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
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

package local

import (
	"context"

	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/storage/dal/types"
)

type errDB struct {
	err error
}

// Table return error for method chaining
func (e *errDB) Table(_ string) types.Table {
	return &errColl{err: e.err}
}

// NextSequence return error for method chaining
func (e *errDB) NextSequence(_ context.Context, _ string) (uint64, error) { return 0, e.err }

// NextSequences return error for method chaining
func (e *errDB) NextSequences(_ context.Context, _ string, _ int) ([]uint64, error) {
	return nil, e.err
}

// Ping return error for method chaining
func (e *errDB) Ping() error { return e.err }

// HasTable return error for method chaining
func (e *errDB) HasTable(_ context.Context, _ string) (bool, error) { return false, e.err }

// ListTables return error for method chaining
func (e *errDB) ListTables(_ context.Context) ([]string, error) { return nil, e.err }

// DropTable return error for method chaining
func (e *errDB) DropTable(_ context.Context, _ string) error { return e.err }

// CreateTable return error for method chaining
func (e *errDB) CreateTable(_ context.Context, _ string) error { return e.err }

// RenameTable return error for method chaining
func (e *errDB) RenameTable(_ context.Context, _, _ string) error { return e.err }

// IsDuplicatedError return error for method chaining
func (e *errDB) IsDuplicatedError(_ error) bool { return false }

// IsNotFoundError return error for method chaining
func (e *errDB) IsNotFoundError(_ error) bool { return false }

// Close return error for method chaining
func (e *errDB) Close() error { return e.err }

// CommitTransaction return error for method chaining
func (e *errDB) CommitTransaction(_ context.Context, _ *metadata.TxnCapable) error { return e.err }

// AbortTransaction return error for method chaining
func (e *errDB) AbortTransaction(_ context.Context, _ *metadata.TxnCapable) (bool, error) {
	return false, e.err
}

// InitTxnManager return error for method chaining
func (e *errDB) InitTxnManager(_ redis.Client) error { return e.err }

type errColl struct {
	err error
}

// Find return error for method chaining
func (e *errColl) Find(_ types.Filter, _ ...*types.FindOpts) types.Find { return &errFind{err: e.err} }

// AggregateOne return error for method chaining
func (e *errColl) AggregateOne(_ context.Context, _ interface{}, _ interface{}) error { return e.err }

// AggregateAll return error for method chaining
func (e *errColl) AggregateAll(_ context.Context, _ interface{}, _ interface{}, _ ...*types.AggregateOpts) error {
	return e.err
}

// Insert return error for method chaining
func (e *errColl) Insert(_ context.Context, _ interface{}) error { return e.err }

// Update return error for method chaining
func (e *errColl) Update(_ context.Context, _ types.Filter, _ interface{}) error { return e.err }

// Upsert return error for method chaining
func (e *errColl) Upsert(_ context.Context, _ types.Filter, _ interface{}) error { return e.err }

// UpdateMultiModel return error for method chaining
func (e *errColl) UpdateMultiModel(_ context.Context, _ types.Filter, _ ...types.ModeUpdate) error {
	return e.err
}

// Delete return error for method chaining
func (e *errColl) Delete(_ context.Context, _ types.Filter) error { return e.err }

// CreateIndex return error for method chaining
func (e *errColl) CreateIndex(_ context.Context, _ types.Index) error { return e.err }

// BatchCreateIndexes return error for method chaining
func (e *errColl) BatchCreateIndexes(_ context.Context, _ []types.Index) error { return e.err }

// DropIndex return error for method chaining
func (e *errColl) DropIndex(_ context.Context, _ string) error { return e.err }

// Indexes return error for method chaining
func (e *errColl) Indexes(_ context.Context) ([]types.Index, error) { return nil, e.err }

// AddColumn return error for method chaining
func (e *errColl) AddColumn(_ context.Context, _ string, _ interface{}) error { return e.err }

// RenameColumn return error for method chaining
func (e *errColl) RenameColumn(_ context.Context, _ types.Filter, _, _ string) error { return e.err }

// DropColumn return error for method chaining
func (e *errColl) DropColumn(_ context.Context, _ string) error { return e.err }

// DropColumns return error for method chaining
func (e *errColl) DropColumns(_ context.Context, _ types.Filter, _ []string) error { return e.err }

// DropDocsColumn return error for method chaining
func (e *errColl) DropDocsColumn(_ context.Context, _ string, _ types.Filter) error { return e.err }

// Distinct return error for method chaining
func (e *errColl) Distinct(_ context.Context, _ string, _ types.Filter) ([]interface{}, error) {
	return nil, e.err
}

// DeleteMany return error for method chaining
func (e *errColl) DeleteMany(_ context.Context, _ types.Filter) (uint64, error) { return 0, e.err }

// UpdateMany return error for method chaining
func (e *errColl) UpdateMany(_ context.Context, _ types.Filter, _ interface{}) (uint64, error) {
	return 0, e.err
}

type errFind struct {
	err error
}

// Fields return error for method chaining
func (e *errFind) Fields(_ ...string) types.Find { return e }

// Sort return error for method chaining
func (e *errFind) Sort(_ string) types.Find { return e }

// Start return error for method chaining
func (e *errFind) Start(_ uint64) types.Find { return e }

// Limit return error for method chaining
func (e *errFind) Limit(_ uint64) types.Find { return e }

// All return error for method chaining
func (e *errFind) All(_ context.Context, _ interface{}) error { return e.err }

// One return error for method chaining
func (e *errFind) One(_ context.Context, _ interface{}) error { return e.err }

// Count return error for method chaining
func (e *errFind) Count(_ context.Context) (uint64, error) { return 0, e.err }

// List return error for method chaining
func (e *errFind) List(_ context.Context, _ interface{}) (int64, error) { return 0, e.err }

// Option return error for method chaining
func (e *errFind) Option(_ ...*types.FindOpts) {}
