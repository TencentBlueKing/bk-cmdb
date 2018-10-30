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

import "testing"

func TestGetInstTableName(t *testing.T) {
	type args struct {
		objID string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"", args{BKInnerObjIDApp}, BKTableNameBaseApp},
		{"", args{BKInnerObjIDSet}, BKTableNameBaseSet},
		{"", args{BKInnerObjIDModule}, BKTableNameBaseModule},
		{"", args{BKInnerObjIDObject}, BKTableNameBaseInst},
		{"", args{BKInnerObjIDHost}, BKTableNameBaseHost},
		{"", args{BKInnerObjIDProc}, BKTableNameBaseProcess},
		{"", args{BKInnerObjIDPlat}, BKTableNameBasePlat},
		{"", args{BKTableNameInstAsst}, BKTableNameInstAsst},
		{"", args{""}, BKTableNameBaseInst},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetInstTableName(tt.args.objID); got != tt.want {
				t.Errorf("GetInstTableName() = %v, want %v", got, tt.want)
			}
		})
	}
}
