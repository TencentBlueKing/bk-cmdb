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

package idgenerator

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm/logger"

	"github.com/TencentBlueKing/bk-cmdb/pkg/dal/table"
	"github.com/TencentBlueKing/bk-cmdb/pkg/tests"
)

func Benchmark_paddingStr(b *testing.B) {
	type args struct {
		strID string
	}
	tests := []struct {
		name string
		args args
		want string
	}{

		{
			name: "1",
			args: args{
				strID: "1",
			},
			want: "00000001",
		},
		{
			name: "12",
			args: args{
				strID: "12",
			},
			want: "00000012",
		},
		{
			name: "123",
			args: args{
				strID: "123",
			},
			want: "00000123",
		},
		{
			name: "123456",
			args: args{
				strID: "123456",
			},
			want: "00123456",
		},
		{
			name: "123456789",
			args: args{
				strID: "123456789",
			},
			want: "123456789",
		},
	}
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.Run("paddingStr", func(t *testing.B) {
				for i := 0; i < b.N; i++ {
					// much faster
					if got := paddingStr(tt.args.strID, StrIDLength, '0'); got != tt.want {
						t.Errorf("paddingStr() = %v, want %v", got, tt.want)
					}
				}

			})
			b.Run("fmt", func(t *testing.B) {
				for i := 0; i < b.N; i++ {
					if got := fmt.Sprintf("%08s", tt.args.strID); got != tt.want {
						t.Errorf("fmt.Sprintf() = %v, want %v", got, tt.want)
					}
				}
			})
		})
	}
}

func TestIdGenerator_Batch(t *testing.T) {
	db, err := tests.GetTestGORM(t)
	if err != nil {
		t.Errorf("failed to open gorm db: %v", err)
		return
	}
	err = db.Migrator().AutoMigrate(&table.IDGenerator{})
	if err != nil {
		t.Fatalf("failed to auto migrate: idgen table, err: %v", err)
	}

	ctx := context.Background()
	idGen := New(db)
	resName := table.Name(fmt.Sprintf("test_%d", time.Now().Unix()))
	t.Cleanup(func() {
		delResult := db.Model(&table.IDGenerator{}).Where("resource = ?", resName).Delete(&table.IDGenerator{})
		if delResult.Error != nil {
			t.Fatalf("failed to delete idgen record for %s, err: %v", resName, delResult.Error)
			return
		}
	})
	t.Run("test invalid table name", func(t *testing.T) {
		err = resName.Validate()
		assert.ErrorContains(t, err, "table name is invalid:", "table name should be invalid without register")
	})
	resName.Register(&table.IDGenerator{})
	t.Run("test valid table name after register", func(t *testing.T) {
		err = resName.Validate()
		assert.NoError(t, err, "table name should be valid after register")
	})

	t.Run("test zero id", func(t *testing.T) {
		ids, err := idGen.Batch(ctx, resName, 0)
		if err != nil {
			t.Fatalf("fial to generate zero id for %s, err: %v", resName, err)
			return
		}
		assert.Len(t, ids, 0)
	})

	t.Run("test one id", func(t *testing.T) {
		firstID, err := idGen.One(ctx, resName)
		if err != nil {
			t.Fatalf("fial to generate first id for %s, err: %v", resName, err)
			return
		}
		assert.Equal(t, "00000001", firstID)
	})

	t.Run("test batch id", func(t *testing.T) {
		idList, err := idGen.Batch(ctx, resName, 4)
		if err != nil {
			t.Fatalf("fial to generate id list for %s, err: %v", resName, err)
			return
		}
		expected := []string{"00000002", "00000003", "00000004", "00000005"}
		assert.Equal(t, expected, idList)
	})

}
func TestParallelBatch(t *testing.T) {
	db, err := tests.GetTestGORM(t)
	if err != nil {
		t.Errorf("failed to open gorm db: %v", err)
		return
	}
	err = db.Migrator().AutoMigrate(&table.IDGenerator{})
	if err != nil {
		t.Fatalf("failed to auto migrate: idgen table, err: %v", err)
		return
	}

	ctx := context.Background()

	resName := table.Name(fmt.Sprintf("test_%d", time.Now().Unix()))
	clean := func() {
		delResult := db.Model(&table.IDGenerator{}).
			Where("resource in ?", []table.Name{resName}).
			Delete(&table.IDGenerator{})
		if delResult.Error != nil {
			t.Fatalf("failed to delete idgen record for %s, err: %v", resName, delResult.Error)
			return
		}
	}
	t.Cleanup(clean)
	resName.Register(&table.IDGenerator{})

	idGen := &idGenerator{db: db}
	t.Run("test parallel batch", func(t *testing.T) {
		clean()
		if t.Failed() {
			return
		}
		testParallelGen(t, idGen.BatchQueryUpdate, ctx, resName, 0, 2000, 1000)
	})
	t.Run("test parallel batch atomic", func(t *testing.T) {
		clean()
		if t.Failed() {
			return
		}
		testParallelGen(t, idGen.BatchUpdateReturning, ctx, resName, 0, 2000, 1000)
	})
}

type idGenFunc func(ctx context.Context, resource table.Name, count uint) ([]string, error)

func testParallelGen(t *testing.T, idGenFunc idGenFunc, ctx context.Context, resName table.Name, currentMax,
	total, parallel int) {

	exceptedMaxID := currentMax + total
	batchSize := uint(total / parallel)

	wg := sync.WaitGroup{}
	allChan := make(chan []string, parallel)
	allResult := make([]string, 0)
	go func() {
		for idList := range allChan {
			allResult = append(allResult, idList...)
		}
	}()
	resultMap := sync.Map{}

	for i := range parallel {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			t.Logf("start to generate id list for idx:%d", idx)
			idList, err := retryOnDuplicate(idGenFunc, ctx, resName, batchSize)
			if err != nil {
				t.Fatalf("fial to generate id list for %s, idx:%d, err: %v", resName, idx, err)
				return
			}
			t.Logf("idx:%d, len:%d: %s", idx, len(idList), idList)
			resultMap.Store(idx, idList)
			allChan <- idList
		}(i)
	}

	wg.Wait()
	close(allChan)
	for i := range parallel {
		_, ok := resultMap.Load(i)
		if !ok {
			t.Fatalf("failed to load id list for idx:%d", i)
			return
		}
	}
	if !assert.Len(t, allResult, total, "total id should be equal to total") {
		return
	}
	allResult = lo.Uniq(allResult)
	if !assert.Len(t, allResult, total, "total id should be equal to total after unique") {
		return
	}
	slices.Sort(allResult)
	minID, _ := strconv.ParseUint(allResult[0], 36, 64)
	maxID, _ := strconv.ParseUint(allResult[len(allResult)-1], 36, 64)
	expectedFirst := currentMax + 1
	assert.Equal(t, uint64(expectedFirst), minID, "min id should be current max id + 1")
	assert.Equal(t, uint64(exceptedMaxID), maxID, "max id should be currentMax + total")
}

type IDGeneratorINT struct {
	Resource table.Name `json:"resource,omitempty" gorm:"resource;primaryKey;size:64"`
	MaxID    int64      `json:"max_id,omitempty" gorm:"max_id;default:0"`
	Prefix   string     `json:"prefix,omitempty" gorm:"prefix;size:64;default:''"`
}

func (i *IDGeneratorINT) TableName() string {
	return "id_generator_int"
}

func BenchmarkIdGenerator_Batch(b *testing.B) {
	db, err := tests.GetTestGORM(b)
	if err != nil {
		b.Errorf("failed to open gorm db: %v", err)
		return
	}
	db.Logger = logger.Default.LogMode(logger.Error)
	err = db.Migrator().AutoMigrate(&table.IDGenerator{}, &IDGeneratorINT{})
	if err != nil {
		b.Fatalf("failed to auto migrate: idgen table, err: %v", err)
	}

	resName := table.Name(fmt.Sprintf("test_%d", time.Now().Unix()))
	resName.Register(&table.IDGenerator{})

	resNameInt := table.Name(fmt.Sprintf("test_%d", time.Now().Unix()))
	resNameInt.Register(&IDGeneratorINT{})

	b.Cleanup(func() {
		delResult := db.Model(&table.IDGenerator{}).
			Where("resource in ?", []table.Name{resName}).
			Delete(&table.IDGenerator{})
		if delResult.Error != nil {
			b.Fatalf("failed to delete idgen record for %s, err: %v", resName, delResult.Error)
			return
		}
	})
	resName.Register(&table.IDGenerator{})
	idGen := &idGenerator{db: db}
	ctx := context.Background()

	b.Run("UpdateReturning-serial", func(b *testing.B) {
		for b.Loop() {
			batch, err := retryOnDuplicate(idGen.BatchUpdateReturning, ctx, resName, 5)
			if err != nil {
				b.Errorf("failed to generate id batch, err: %v", err)
				return
			}
			_ = batch
		}
	})
	b.Run("UpdateReturning-parallel", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				batch, err := retryOnDuplicate(idGen.BatchUpdateReturning, ctx, resName, 5)
				if err != nil {
					b.Errorf("failed to generate id batch, err: %v", err)
					return
				}
				_ = batch
			}
		})
	})

	b.Run("QueryUpdate-serial", func(b *testing.B) {
		for b.Loop() {
			batch, err := retryOnDuplicate(idGen.BatchQueryUpdate, ctx, resName, 5)
			if err != nil {
				b.Errorf("failed to generate id batch, err: %v", err)
				return
			}
			_ = batch
		}
	})
	b.Run("QueryUpdate-parallel", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				batch, err := retryOnDuplicate(idGen.BatchQueryUpdate, ctx, resName, 5)
				if err != nil {
					b.Errorf("failed to generate id batch, err: %v", err)
					return
				}
				_ = batch
			}
		})
	})

}
