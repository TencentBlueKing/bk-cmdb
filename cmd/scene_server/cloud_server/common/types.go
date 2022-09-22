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

package common

import (
	"errors"
)

const (
	// MaxLimit 最大数量限制
	MaxLimit int64 = 999999

	// MaxLoopCnt 最大循环次数
	MaxLoopCnt int = 10
)

var (
	// ErrorLoopCnt 循环过多错误
	ErrorLoopCnt = errors.New("too much loop")
)

// BaseOpt 云厂商接口请求条件的公共部分
type BaseOpt struct {
	// 请求过滤条件
	Filters []*Filter

	// 返回数量限制
	Limit int64
}

// Filter 过滤条件
type Filter struct {
	// 需要过滤的字段。
	Name *string

	// 字段的过滤值。
	Values []*string
}

// VpcOpt VPC请求条件
type VpcOpt struct {
	BaseOpt
}

// InstanceOpt 实例请求条件
type InstanceOpt struct {
	BaseOpt
}

// GetDefaultVpcOpt 获取默认的Vpc请求条件
func GetDefaultVpcOpt() *VpcOpt {
	return &VpcOpt{
		BaseOpt{
			Limit: MaxLimit,
		},
	}
}

// GetDefaultInstanceOpt 获取默认的Instance请求条件
func GetDefaultInstanceOpt() *InstanceOpt {
	return &InstanceOpt{
		BaseOpt{
			Limit: MaxLimit,
		},
	}
}

// Int64Ptr 获取int64的指针
func Int64Ptr(v int64) *int64 {
	return &v
}

// StringPtr 获取string的指针
func StringPtr(v string) *string {
	return &v
}

// StringPtrs 获取[]string的指针
func StringPtrs(vals []string) []*string {
	ptrs := make([]*string, len(vals))
	for i := 0; i < len(vals); i++ {
		ptrs[i] = &vals[i]
	}
	return ptrs
}
