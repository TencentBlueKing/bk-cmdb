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

// 云厂商接口请求选项
type RequestOpt struct {
	// 请求过滤条件
	Filters []*Filter

	// 偏移量，默认为0，分页时使用，对于某些云厂商适用，如腾讯云
	Offset *int64

	// 返回数量限制
	Limit *int64

	// 用来获取下一页数据时的请求参数，对于某些云厂商适用，如AWS
	NextToken *string
}

type Filter struct {
	// 需要过滤的字段。
	Name *string

	// 字段的过滤值。
	Values []*string
}

const (
	MaxLimit   int64 = 99999
	MaxLoopCnt int   = 10
)

var (
	ErrorLoopCnt = errors.New("too much loop")
)

func Int64Ptr(v int64) *int64 {
	return &v
}

func StringPtr(v string) *string {
	return &v
}

func StringPtrs(vals []string) []*string {
	ptrs := make([]*string, len(vals))
	for i := 0; i < len(vals); i++ {
		ptrs[i] = &vals[i]
	}
	return ptrs
}
