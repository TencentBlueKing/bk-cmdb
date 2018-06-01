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
 
package util

import (
	"testing"
	"time"
)

func TestMapMatch(t *testing.T) {
	type args struct {
		src interface{}
		tar interface{}
	}
	now := time.Now()
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			args: args{
				src: map[string]interface{}{"a": "1", "time": Time(time.Unix(now.Unix(), 0))},
				tar: map[string]interface{}{"a": "1", "time": &now, "n": "m"},
			},
			want: true,
		},
		{
			args: args{
				src: map[string]interface{}{"a": "1", "time": time.Unix(now.Unix(), 0), "k": "j"},
				tar: map[string]interface{}{"a": "1", "time": time.Unix(now.Unix(), 0), "n": "m"},
			},
			want: false,
		},
		{
			args: args{
				src: map[string]interface{}{"a": "1", "time": time.Unix(now.Unix(), 0)},
				tar: map[string]interface{}{"a": "1", "time": nil, "n": "m"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MapMatch(tt.args.src, tt.args.tar); got != tt.want {
				t.Errorf("MapMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

type Time time.Time
