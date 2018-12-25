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

package eventclient

import (
	"net/http"
	"testing"

	"configcenter/src/common"
	"configcenter/src/common/metadata"
)

func TestNewEventWithHeader(t *testing.T) {
	type args struct {
		header http.Header
	}

	header := http.Header{}
	header.Set(common.BKHTTPOwnerID, "0")
	header.Set(common.BKHTTPCCTransactionID, "123")
	header.Set(common.BKHTTPCCRequestID, "456")
	tests := []struct {
		name string
		args args
		want *metadata.EventInst
	}{
		{
			args: args{
				header: header,
			},
			want: &metadata.EventInst{
				OwnerID:   "0",
				TxnID:     "123",
				RequestID: "456",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEventWithHeader(tt.args.header); got.OwnerID != tt.want.OwnerID ||
				got.TxnID != tt.want.TxnID ||
				got.RequestID != tt.want.RequestID {
				t.Errorf("NewEventWithHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}
