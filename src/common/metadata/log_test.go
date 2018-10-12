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
	"testing"
	"time"
)

func TestOperationLog_TableName(t *testing.T) {
	type fields struct {
		OwnerID       string
		ApplicationID int64
		ExtKey        string
		OpDesc        string
		OpType        int
		OpTarget      string
		Content       interface{}
		User          string
		OpFrom        string
		ExtInfo       string
		CreateTime    time.Time
		InstID        int64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"", fields{}, "cc_OperationLog"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := OperationLog{
				OwnerID:       tt.fields.OwnerID,
				ApplicationID: tt.fields.ApplicationID,
				ExtKey:        tt.fields.ExtKey,
				OpDesc:        tt.fields.OpDesc,
				OpType:        tt.fields.OpType,
				OpTarget:      tt.fields.OpTarget,
				Content:       tt.fields.Content,
				User:          tt.fields.User,
				OpFrom:        tt.fields.OpFrom,
				ExtInfo:       tt.fields.ExtInfo,
				CreateTime:    tt.fields.CreateTime,
				InstID:        tt.fields.InstID,
			}
			if got := o.TableName(); got != tt.want {
				t.Errorf("OperationLog.TableName() = %v, want %v", got, tt.want)
			}
		})
	}
}
