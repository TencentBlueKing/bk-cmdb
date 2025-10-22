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

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/plugin/dbresolver"

	"github.com/TencentBlueKing/bk-cmdb/pkg/dal/table"
)

// StrIDLength is the length of the string id.
const StrIDLength = 8

// Interface supplies all the method to generate a resource's unique identity id.
type Interface interface {
	// Batch return a list of resource's unique id as required.
	Batch(ctx context.Context, resource table.Name, count uint) ([]string, error)
	// One return one unique id for this resource.
	One(ctx context.Context, resource table.Name) (string, error)
}

var _ Interface = new(idGenerator)

// New create an id generator instance.
func New(db *gorm.DB) Interface {
	return &idGenerator{db: db}
}

type idGenerator struct {
	db *gorm.DB
}

// One generate one unique resource id.
func (ig idGenerator) One(ctx context.Context, resource table.Name) (string, error) {
	list, err := ig.Batch(ctx, resource, 1)
	if err != nil {
		return "", err
	}

	if num := len(list); num != 1 {
		return "", fmt.Errorf("gen resource unique id, but %d returned", num)
	}

	return list[0], nil
}

// Batch is to generate distribute unique resource id list.
// returned with a number of unique ids as required.
func (ig idGenerator) Batch(ctx context.Context, resource table.Name, count uint) ([]string, error) {
	if err := resource.Validate(); err != nil {
		return nil, err
	}

	idLine := table.IDGenerator{
		Resource: resource,
		MaxID:    "0",
	}
	maxID := uint64(0)
	err := ig.db.WithContext(ctx).
		Clauses(dbresolver.Write).
		Table(table.IDGeneratorTable.String()).
		Transaction(func(txn *gorm.DB) error {

			// get current max id with for update, if not exist, create with default value 0
			result := txn.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
				Where(map[string]any{"resource": resource}).
				FirstOrCreate(&idLine)
			if result.Error != nil {
				return fmt.Errorf("get current max id for resource %s fail, err: %w", resource, result.Error)
			}

			var err error
			// generate new max id and update it
			maxID, err = strconv.ParseUint(idLine.MaxID, 36, 64)
			if err != nil {
				return fmt.Errorf("gen %s unique id, but parse max id failed, err: %w", resource, err)
			}

			newMaxID := formatStrID(maxID + uint64(count))

			// update with old max id to prevent id conflict
			updateResult := txn.Where(map[string]any{"resource": resource, "max_id": idLine.MaxID}).
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
		ids[idx] = formatStrID(maxID + uint64(idx+1))
	}

	return ids, nil
}

func formatStrID(id uint64) string {
	strID := strconv.FormatUint(id, 36)
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
	return string(result)
}
