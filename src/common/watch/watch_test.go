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

// Package watch TODO
package watch

import (
	"testing"
	"time"
)

func TestTimeDrift(t *testing.T) {
	type args struct {
		duration string
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		//{"empty", args{""}, 0 * time.Hour},//err
		{"0s", args{"0s"}, 0 * time.Hour},
		{"-0s", args{"-0s"}, 0 * time.Hour},
		{"1h1m59s", args{"1h1m59s"}, time.Hour + time.Minute + 59*time.Second},
		{"3h", args{"3h"}, 3 * time.Hour},
		{"-3h", args{"-3h"}, -3 * time.Hour},
		{"-4h", args{"-4h"}, -4 * time.Hour},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			duration, err := time.ParseDuration(tt.args.duration)
			if err != nil {
				t.Error(err)
			}
			if duration != tt.want {
				t.Errorf("got %v, want %v", duration, tt.want)
			}
		})
	}
}
