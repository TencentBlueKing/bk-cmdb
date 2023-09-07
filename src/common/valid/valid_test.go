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

package valid

import (
	"testing"

	"configcenter/src/common"
)

func TestIsInnerObject(t *testing.T) {
	type args struct {
		objID string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"", args{"id"}, false},
		{"", args{common.BKInnerObjIDApp}, true},
		{"", args{common.BKInnerObjIDHost}, true},
		{"", args{common.BKInnerObjIDModule}, true},
		{"", args{common.BKInnerObjIDPlat}, true},
		{"", args{common.BKInnerObjIDProc}, true},
		{"", args{common.BKInnerObjIDSet}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInnerObject(tt.args.objID); got != tt.want {
				t.Errorf("IsInnerObject() = %v, want %v", got, tt.want)
			}
		})
	}
}
