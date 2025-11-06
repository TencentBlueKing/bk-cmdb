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

// Package idgenerator provides a id generator.
package idgenerator

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"unsafe"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/plugin/dbresolver"

	"github.com/TencentBlueKing/bk-cmdb/pkg/dal/orm"
	"github.com/TencentBlueKing/bk-cmdb/pkg/dal/table"
	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
)

// StrIDLength is the length of the string id.
const StrIDLength = 8

// StrIDBase is the base of the string format id.
const StrIDBase = 36

// Interface supplies all the method to generate a resource's unique identity id.
type Interface interface {
	// Batch return a list of resource's unique id as required.
	Batch(ctx context.Context, resource table.Name, count uint64) ([]string, error)
	// One return one unique id for this resource.
	One(ctx context.Context, resource table.Name) (string, error)
	// InitTable initialize the table for the resource
	InitTable(ctx context.Context, resource table.Name, initValue uint64) (err error)
}

var _ Interface = new(idGenerator)

// New create an id generator instance.
func New(db *gorm.DB) Interface {
	return &idGenerator{db: db}
}

type idGenerator struct {
	db *gorm.DB
}

// InitTable initialize row for the resource
func (ig *idGenerator) InitTable(ctx context.Context, resource table.Name, initValue uint64) (err error) {
	return ig.initTable(ctx, ig.db, resource, initValue)
}

// initTable initialize row for the resource
func (ig *idGenerator) initTable(ctx context.Context, txn *gorm.DB, resource table.Name, initValue uint64) (err error) {
	if err = resource.Validate(); err != nil {
		return err
	}

	value := &table.IDGenerator{
		Resource: resource,
		MaxID:    initValue,
	}
	err = txn.WithContext(ctx).Clauses(dbresolver.Write).Create(value).Error

	return err
}

// One generate one unique resource id.
func (ig *idGenerator) One(ctx context.Context, resource table.Name) (string, error) {
	list, err := ig.Batch(ctx, resource, 1)
	if err != nil {
		return "", err
	}

	if num := len(list); num != 1 {
		return "", fmt.Errorf("gen resource unique id, but %d returned", num)
	}

	return list[0], nil
}

// Batch is to generate distribute unique resource id list. returned with a number of unique ids as required.
func (ig *idGenerator) Batch(ctx context.Context, resource table.Name, count uint64) (ids []string, err error) {
	f := ig.BatchUpdateReturning
	if orm.IsPostgres(ig.db) {
		// only pg support update returning
		f = ig.BatchQueryUpdate
	}
	return retryOnDuplicate(f, ctx, resource, count)
}

func retryOnDuplicate(f func(ctx context.Context, resource table.Name, count uint64) ([]string, error),
	ctx context.Context, resource table.Name, count uint64) (ids []string, err error) {

	const mostRetry = 3
	const retryInterval = 10 * time.Millisecond

	for i := range mostRetry {
		ids, err = f(ctx, resource, count)
		if err == nil {
			return ids, err
		}
		if orm.IsDuplicatedError(err) {
			log.Warn(ctx, "got duplicate key err", "retry", i, log.E(err))
			time.Sleep(retryInterval)
			continue
		}
		// unknown error
		return ids, err
	}
	return nil, fmt.Errorf("too many retry: %d, last err: %w", mostRetry, err)
}

// BatchQueryUpdate is to generate distribute unique resource id list
func (ig *idGenerator) BatchQueryUpdate(ctx context.Context, resource table.Name, count uint64) ([]string, error) {
	if err := resource.Validate(); err != nil {
		return nil, err
	}

	oldMaxID := uint64(0)
	err := ig.db.WithContext(ctx).
		Clauses(dbresolver.Write).
		Table(table.IDGeneratorTable.String()).
		Transaction(func(txn *gorm.DB) error {

			// get current max id with for update, if not exist, create with default value 0
			result := txn.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
				Select("max_id").
				Where("resource= ? ", resource).
				Find(&oldMaxID)
			if result.Error != nil {
				return fmt.Errorf("get current max id for resource %s fail, err: %w", resource, result.Error)
			}
			if result.RowsAffected == 0 {
				return ig.initTable(ctx, txn, resource, count)
			}

			// generate new max id and update it
			newMaxID := oldMaxID + (count)
			// update with old max id to prevent id conflict
			updateResult := txn.Where("resource= ? and max_id = ?", resource, oldMaxID).
				Update("max_id", newMaxID)
			if updateResult.Error != nil {
				return fmt.Errorf("gen %s ids, but update max_id failed, err: %w", resource, updateResult.Error)
			}
			if updateResult.RowsAffected != 1 {
				return fmt.Errorf("gen %s ids, but update max_id failed, rows affected %d is not 1",
					resource, updateResult.RowsAffected)
			}

			// commit transaction
			return nil
		})
	if err != nil {
		return nil, fmt.Errorf("gen %s unique id, but transaction failed, err: %w", resource, err)
	}

	// generate the id list that can be used.
	ids := make([]string, count)
	for idx := range count {
		ids[idx] = formatStrID(oldMaxID + (idx + 1))
	}

	return ids, nil
}

// BatchUpdateReturning is to generate distribute unique resource id list using `UPDATE ... RETURNING` statement.
func (ig *idGenerator) BatchUpdateReturning(ctx context.Context, resource table.Name, count uint64) ([]string, error) {
	if err := resource.Validate(); err != nil {
		return nil, err
	}

	newMaxID := uint64(0)
	ret := ig.db.WithContext(ctx).Clauses(dbresolver.Write).
		Raw(fmt.Sprintf(`UPDATE %s SET max_id = max_id + %d WHERE resource = '%s' RETURNING max_id`,
			table.IDGeneratorTable.String(), count, resource)).
		Scan(&newMaxID)
	err := ret.Error
	if err != nil {
		return nil, fmt.Errorf("gen %s unique id, but update with returning failed, err: %w", resource, err)
	}
	if ret.RowsAffected == 0 || newMaxID == 0 {
		// only trigger when resource row is not exist
		newMaxID = count
		err = ig.initTable(ctx, ig.db, resource, newMaxID)
		if err != nil {
			return nil, fmt.Errorf("gen %s unique id, but create failed, err: %w", resource, err)
		}
	}

	oldMaxID := newMaxID - count
	// generate the id list that can be used.
	ids := make([]string, count)
	for idx := range count {
		id := oldMaxID + idx + 1
		ids[idx] = formatStrID(id)
	}

	return ids, nil
}

func formatStrID(id uint64) string {
	strID := strconv.FormatUint(id, StrIDBase)
	strID = paddingStr(strID, StrIDLength, '0')
	return strID
}

func paddingStr(str string, length int, fill byte) string {
	if len(str) >= length {
		return str
	}
	result := make([]byte, length)
	for i := 0; i < length-len(str); i++ {
		result[i] = fill
	}
	copy(result[length-len(str):length], str)
	return unsafe.String(unsafe.SliceData(result), length)
}
