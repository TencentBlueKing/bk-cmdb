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

func TestGetInstNameField(t *testing.T) {
	type args struct {
		objID string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"", args{BKInnerObjIDApp}, BKAppNameField},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetInstNameField(tt.args.objID); got != tt.want {
				t.Errorf("GetInstNameField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetInstIDField(t *testing.T) {
	type args struct {
		objType string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"", args{"not found"}, BKInstIDField},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetInstIDField(tt.args.objType); got != tt.want {
				t.Errorf("GetInstIDField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetObjByType(t *testing.T) {
	type args struct {
		objType string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"", args{BKInnerObjIDHost}, BKInnerObjIDHost},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetObjByType(tt.args.objType); got != tt.want {
				t.Errorf("GetObjByType() = %v, want %v", got, tt.want)
			}
		})
	}
}
