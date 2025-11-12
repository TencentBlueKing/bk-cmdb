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

package kit

import (
	"testing"

	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
	_ "github.com/TencentBlueKing/bk-cmdb/pkg/trace" // for init trace
	"github.com/stretchr/testify/assert"
)

func TestNewKit(t *testing.T) {
	kt := NewKit(t.Context(), Metadata{})
	assert.Equal(t, len(kt.Rid()), 32)

	log.Info(kt, "test newkit")

	testDoBiz(kt)
}

func TestKitStartSpan(t *testing.T) {
	kt := NewKit(t.Context(), Metadata{})

	kt, span := kt.StartSpan("")
	defer span.End()

	log.Info(kt, "first level span")

	testDoBiz(kt)
}

func testDoBiz(kt *Kit) {
	kt, span := kt.StartSpan("")
	defer span.End()

	log.Info(kt, "doBiz")
}

func testBiz(_ *Kit, n int) int {
	return n + 1
}

func BenchmarkBiz(b *testing.B) {
	kt := NewKit(b.Context(), Metadata{})

	for b.Loop() {
		testBiz(kt, 1)
	}
}

func BenchmarkBizSpan(b *testing.B) {
	kt := NewKit(b.Context(), Metadata{})
	for b.Loop() {
		kt, span := kt.StartSpan("")
		defer span.End()

		testBiz(kt, 1)
	}
}

func BenchmarkBizSpanWithName(b *testing.B) {
	kt := NewKit(b.Context(), Metadata{})
	for b.Loop() {
		kt, span := kt.StartSpan("test")
		defer span.End()

		testBiz(kt, 1)
	}
}

func BenchmarkBizSpanWithoutEnd(b *testing.B) {
	kt := NewKit(b.Context(), Metadata{})
	for b.Loop() {
		kt, _ := kt.StartSpan("test")

		testBiz(kt, 1)
	}
}
