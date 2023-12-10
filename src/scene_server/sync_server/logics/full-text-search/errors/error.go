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

// Package errors defines the full-text search error handler
package errors

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/olivere/elastic/v7"
)

// BasicErrHandler is the basic err handler that retries after an increasing interval
// @param baseTime: the basic retry sleep time in milliseconds
// @param randTime: the random retry sleep time maximum value in milliseconds
func BasicErrHandler(baseTime, randTime int, operator func() (bool, error)) {
	retry := 0
	for {
		needRetry, err := operator()
		if err == nil {
			return
		}

		if !needRetry {
			return
		}

		retry++

		rand.Seed(time.Now().UnixNano())
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(randTime)+baseTime) * time.Duration(retry))
	}
}

// FatalErrHandler is the err handler for fatal error that must be retried
// strategy: always retry, because when there is a fatal error, every operation will be failed
func FatalErrHandler(baseTime, randTime int, operator func() error) {
	BasicErrHandler(baseTime, randTime, func() (bool, error) {
		err := operator()
		return true, err
	})
}

// EsStatusHandler is the es status handler, returns if the request should be retied & if it's a fatal error
func EsStatusHandler(status int) (bool, bool) {
	// the request is successful
	if status >= 200 && status <= 299 {
		return false, false
	}

	// skip the invalid requests
	if elastic.IsForbidden(status) || elastic.IsUnauthorized(status) ||
		elastic.IsStatusCode(status, http.StatusBadRequest) {
		return false, false
	}

	// ignores version conflict error
	if elastic.IsConflict(status) {
		return true, false
	}

	// this status mostly means index not exists, so we sleep for a long time to wait until index is recovered
	if elastic.IsNotFound(status) {
		time.Sleep(5 * time.Minute)
		return true, true
	}

	// sleep for a long time to lower the request num
	if elastic.IsTimeout(status) || elastic.IsStatusCode(status, http.StatusTooManyRequests) {
		time.Sleep(2 * time.Minute)
		return true, false
	}

	return true, false
}

// EsErrRetryCount is the retry count for es error
const EsErrRetryCount = 5

// EsRespErrHandler is the response err handler for es operation
func EsRespErrHandler(operator func() (bool, error)) error {
	retry := 1
	var err error

	BasicErrHandler(200, 100, func() (bool, error) {
		var fatal bool
		fatal, err = operator()
		if err == nil {
			return false, nil
		}

		if elastic.IsConnErr(err) {
			return true, err
		}

		if fatal {
			return true, err
		}

		if retry == EsErrRetryCount {
			return false, err
		}

		time.Sleep(time.Duration(retry) * 100 * time.Millisecond)
		retry++
		return true, err
	})

	return err
}
