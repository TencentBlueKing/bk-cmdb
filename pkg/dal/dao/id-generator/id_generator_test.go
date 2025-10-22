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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

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
		idList, err := idGen.Batch(ctx, resName, 5)
		if err != nil {
			t.Fatalf("fial to generate id list for %s, err: %v", resName, err)
			return
		}
		expected := []string{"00000002", "00000003", "00000004", "00000005", "00000006"}
		assert.Equal(t, expected, idList)
	})

}
