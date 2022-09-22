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

package mainline

import "testing"

func TestBusinessKey(t *testing.T) {
	// test resume token key
	if bizKey.resumeTokenKey() != "cc:v3:biz:biz:resume_token" {
		t.Fatalf("invalid biz resume token key")
	}

	// test resume at time key
	if bizKey.resumeAtTimeKey() != "cc:v3:biz:biz:resume_at_time" {
		t.Fatalf("invalid biz resume at time key")
	}

	// test detail key
	if bizKey.detailKey(1) != "cc:v3:biz:biz_detail:1" {
		t.Fatalf("invalid biz resume at time key")
	}
}

func TestCustomKey(t *testing.T) {
	country := newCustomKey("country")
	// test resume token key
	if country.resumeTokenKey() != "cc:v3:biz:country:resume_token" {
		t.Fatalf("invalid custom country object instance resume token key")
	}

	// test resume at time key
	if country.resumeAtTimeKey() != "cc:v3:biz:country:resume_at_time" {
		t.Fatalf("invalid custom country object instance resume at time key")
	}

	// test detail key
	if country.detailKey(1) != "cc:v3:biz:country_detail:1" {
		t.Fatalf("invalid custom country object instance resume at time key")
	}

}
