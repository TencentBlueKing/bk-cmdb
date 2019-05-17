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
	"reflect"
	"testing"

	"configcenter/src/common"
)

func TestSetQueryOwner(t *testing.T) {
	type args struct {
		condition interface{}
		ownerID   string
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			"",
			args{nil, "ownerid"},
			map[string]interface{}{
				common.BKOwnerIDField: map[string]interface{}{common.BKDBIN: []string{common.BKDefaultOwnerID, "ownerid"}},
			},
		},
		{
			"",
			args{nil, common.BKSuperOwnerID},
			map[string]interface{}{},
		},
		{
			"",
			args{struct{ Name string }{Name: "haha"}, common.BKSuperOwnerID},
			map[string]interface{}{
				"name": "haha",
			},
		},
		{
			"",
			args{struct{ Name string }{Name: "haha"}, "ownerid"},
			map[string]interface{}{
				"name":                "haha",
				common.BKOwnerIDField: map[string]interface{}{common.BKDBIN: []string{common.BKDefaultOwnerID, "ownerid"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetQueryOwner(tt.args.condition, tt.args.ownerID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetQueryOwner() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetModOwner(t *testing.T) {
	type args struct {
		condition interface{}
		ownerID   string
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{"", args{nil, "ownerid"}, map[string]interface{}{
			common.BKOwnerIDField: "ownerid",
		}},
		{
			"",
			args{nil, common.BKSuperOwnerID},
			map[string]interface{}{},
		},
		{
			"",
			args{struct{ Name string }{Name: "haha"}, common.BKSuperOwnerID},
			map[string]interface{}{
				"name": "haha",
			},
		},
		{
			"",
			args{struct{ Name string }{Name: "haha"}, "ownerid"},
			map[string]interface{}{
				"name":                "haha",
				common.BKOwnerIDField: "ownerid",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetModOwner(tt.args.condition, tt.args.ownerID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetModOwner() = %v, want %v", got, tt.want)
			}
		})
	}
}
