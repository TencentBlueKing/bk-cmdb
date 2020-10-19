/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package redis

import (
	"time"
)

// baseResult is the base result for redis commands
type baseResult interface {
	Err() error
}

// Result is the common result for redis commands
type Result interface {
	baseResult
	Val() interface{}
	Result() (interface{}, error)
}

// StringResult is the string result for redis commands
type StringResult interface {
	baseResult
	Val() string
	Result() (string, error)
	Scan(val interface{}) error
}

// FloatResult is the float result for redis commands
type FloatResult interface {
	baseResult
	Val() float64
	Result() (float64, error)
}

// IntResult is the int result for redis commands
type IntResult interface {
	baseResult
	Val() int64
	Result() (int64, error)
}

// SliceResult is the slice result for redis commands
type SliceResult interface {
	baseResult
	Val() []interface{}
	Result() ([]interface{}, error)
}

// StatusResult is the status result for redis commands
type StatusResult interface {
	baseResult
	Val() string
	Result() (string, error)
}

// BoolResult the bool result for redis commands
type BoolResult interface {
	baseResult
	Val() bool
	Result() (bool, error)
}

// IntSliceResult is the int slice result for redis commands
type IntSliceResult interface {
	baseResult
	Val() []int64
	Result() ([]int64, error)
}

// StringSliceResult is the string slice result for redis commands
type StringSliceResult interface {
	baseResult
	Val() []string
	Result() ([]string, error)
}

// BoolSliceResult is the bool slice result for redis commands
type BoolSliceResult interface {
	baseResult
	Val() []bool
	Result() ([]bool, error)
}

// StringStringMapResult is the string string map result for redis commands
type StringStringMapResult interface {
	baseResult
	Val() map[string]string
	Result() (map[string]string, error)
}

// StringIntMapResult is the string int map result for redis commands
type StringIntMapResult interface {
	baseResult
	Val() map[string]int64
	Result() (map[string]int64, error)
}

// StringStructMapResult is the string struct map result for redis commands
type StringStructMapResult interface {
	baseResult
	Val() map[string]struct{}
	Result() (map[string]struct{}, error)
}

// DurationResult is the duration result for redis commands
type DurationResult interface {
	baseResult
	Val() time.Duration
	Result() (time.Duration, error)
}

// ScanResult is the duration result for redis commands
type ScanResult interface {
	baseResult
	Val() (keys []string, cursor uint64)
	Result() (keys []string, cursor uint64, err error)
}
