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
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestRespError_Error(t *testing.T) {
	type fields struct {
		Msg     error
		ErrCode int
		Data    interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"", fields{fmt.Errorf("error"), 0, "data"}, `{"result":false,"bk_error_code":0,"bk_error_msg":"error","data":"data"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RespError{
				Msg:     tt.fields.Msg,
				ErrCode: tt.fields.ErrCode,
				Data:    tt.fields.Data,
			}
			if got := r.Error(); got != tt.want {
				t.Errorf("RespError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSuccessResp(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name string
		args args
		want *Response
	}{
		{"", args{"data"}, &Response{
			BaseResp: BaseResp{true, 0, "success"},
			Data:     "data",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSuccessResp(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSuccessResp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryInput_ConvTime(t *testing.T) {
	type fields struct {
		Condition interface{}
		Fields    string
		Start     int
		Limit     int
		Sort      string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"", fields{
			Condition: "",
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &QueryInput{
				Condition: tt.fields.Condition,
				Fields:    tt.fields.Fields,
				Start:     tt.fields.Start,
				Limit:     tt.fields.Limit,
				Sort:      tt.fields.Sort,
			}
			if err := o.ConvTime(); (err != nil) != tt.wantErr {
				t.Errorf("QueryInput.ConvTime() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQueryInput_convTimeItem(t *testing.T) {
	type fields struct {
		Condition interface{}
		Fields    string
		Start     int
		Limit     int
		Sort      string
	}
	type args struct {
		item interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{"", fields{}, args{nil}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &QueryInput{
				Condition: tt.fields.Condition,
				Fields:    tt.fields.Fields,
				Start:     tt.fields.Start,
				Limit:     tt.fields.Limit,
				Sort:      tt.fields.Sort,
			}
			got, err := o.convTimeItem(tt.args.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryInput.convTimeItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryInput.convTimeItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryInput_convInterfaceToTime(t *testing.T) {
	type fields struct {
		Condition interface{}
		Fields    string
		Start     int
		Limit     int
		Sort      string
	}
	type args struct {
		val interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{"", fields{}, args{"2010-01-01"}, time.Date(2010, 1, 1, 0, 0, 0, 0, time.Local).UTC(), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &QueryInput{
				Condition: tt.fields.Condition,
				Fields:    tt.fields.Fields,
				Start:     tt.fields.Start,
				Limit:     tt.fields.Limit,
				Sort:      tt.fields.Sort,
			}
			got, err := o.convInterfaceToTime(tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryInput.convInterfaceToTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryInput.convInterfaceToTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
