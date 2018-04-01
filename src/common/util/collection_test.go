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
)

func TestCalSliceDiff(t *testing.T) {
	type args struct {
		oldslice []string
		newslice []string
	}
	tests := []struct {
		name      string
		args      args
		wantSubs  []string
		wantPlugs []string
	}{
		{
			args: args{
				oldslice: []string{"a", "b", "c"},
				newslice: []string{"b", "c", "d"},
			},
			wantSubs:  []string{"a"},
			wantPlugs: []string{"d"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSubs, gotPlugs := CalSliceDiff(tt.args.oldslice, tt.args.newslice)
			if !reflect.DeepEqual(gotSubs, tt.wantSubs) {
				t.Errorf("CalSliceDiff() gotSubs = %v, want %v", gotSubs, tt.wantSubs)
			}
			if !reflect.DeepEqual(gotPlugs, tt.wantPlugs) {
				t.Errorf("CalSliceDiff() gotPlugs = %v, want %v", gotPlugs, tt.wantPlugs)
			}
		})
	}
}

func TestContains(t *testing.T) {
	type args struct {
		set    []string
		substr string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			args: args{
				set:    []string{"a", "b", "c"},
				substr: "a",
			},
			want: true,
		},
		{
			args: args{
				set:    []string{"a", "b", "c"},
				substr: "d",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.args.set, tt.args.substr); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainsInt64(t *testing.T) {
	type args struct {
		set []int64
		sub int64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			args: args{
				set: []int64{1, 2, 3},
				sub: 1,
			},
			want: true,
		},
		{
			args: args{
				set: []int64{1, 2, 3},
				sub: 4,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsInt64(tt.args.set, tt.args.sub); got != tt.want {
				t.Errorf("ContainsInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalSliceInt64Diff(t *testing.T) {
	type args struct {
		oldslice []int64
		newslice []int64
	}
	tests := []struct {
		name      string
		args      args
		wantSubs  []int64
		wantInter []int64
		wantPlugs []int64
	}{
		{
			args: args{
				oldslice: []int64{1, 2, 3},
				newslice: []int64{2, 3, 4},
			},
			wantSubs:  []int64{1},
			wantInter: []int64{2, 3},
			wantPlugs: []int64{4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSubs, gotInter, gotPlugs := CalSliceInt64Diff(tt.args.oldslice, tt.args.newslice)
			if !reflect.DeepEqual(gotSubs, tt.wantSubs) {
				t.Errorf("CalSliceInt64Diff() gotSubs = %v, want %v", gotSubs, tt.wantSubs)
			}
			if !reflect.DeepEqual(gotInter, tt.wantInter) {
				t.Errorf("CalSliceInt64Diff() gotInter = %v, want %v", gotInter, tt.wantInter)
			}
			if !reflect.DeepEqual(gotPlugs, tt.wantPlugs) {
				t.Errorf("CalSliceInt64Diff() gotPlugs = %v, want %v", gotPlugs, tt.wantPlugs)
			}
		})
	}
}
