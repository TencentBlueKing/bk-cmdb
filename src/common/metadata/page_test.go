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

package metadata

import (
	"reflect"
	"testing"
)

func TestParsePage(t *testing.T) {
	type args struct {
		origin interface{}
	}
	tests := []struct {
		name string
		args args
		want BasePage
	}{
		{"", args{map[string]interface{}{"sort": "f", "limit": 9, "start": 3}}, BasePage{Sort: "f", Limit: 9, Start: 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParsePage(tt.args.origin); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParsePage() = %v, want %v", got, tt.want)
			}
		})
	}
}
