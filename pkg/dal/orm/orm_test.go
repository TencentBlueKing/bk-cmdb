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

package orm

import (
	"math/rand/v2"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/TencentBlueKing/bk-cmdb/pkg/tests"
)

func TestOrmRateLimit(t *testing.T) {
	orm, err := New(t.Context(), nil)
	assert.ErrorContains(t, err, "db is nil", "New() error")
	assert.Nil(t, orm)

	db, mock, err := tests.GetMockPG(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	orm, err = New(t.Context(), db)
	if !assert.NoError(t, err) {
		return
	}

	var atomicErr atomic.Pointer[error]
	wg := sync.WaitGroup{}

	concurrency := DefaultIngressLimit + 100
	for range concurrency {
		mock.ExpectQuery(`SELECT 1 FROM "a"`).
			WillDelayFor(150*time.Millisecond + (time.Duration(rand.Int64N(100)) * time.Millisecond)).
			WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))
	}
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ret := 0
			err := orm.DB().Table("a").Select("1").Find(&ret).Error
			if err != nil {
				assert.ErrorIs(t, err, ErrTooManyRequests, "should got ErrTooManyRequests at %d", i)
				atomicErr.CompareAndSwap(nil, &err)
				return
			}
			assert.Equalf(t, 1, ret, "should got 1 at %d", i)
		}(i)
	}
	wg.Wait()
	e := atomicErr.Load()
	if !assert.NotNil(t, e, "no error occurred") {
		return
	}
	assert.ErrorIs(t, *e, ErrTooManyRequests, "should got ErrTooManyRequests")
	// because of rate limiter, expectations will not be met
	assert.NotNil(t, mock.ExpectationsWereMet(), "expectations not met")

}
