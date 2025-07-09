/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

package common

import "testing"

func TestSetIdentification(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
	}{
		{"", args{"id"}},
		{"", args{"unkown"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetIdentification(tt.args.id)
		})
	}
}

func TestGetIdentification(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"", "unkonw"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetIdentification(); got != tt.want {
				t.Errorf("GetIdentification() = %v, want %v", got, tt.want)
			}
		})
	}
}
