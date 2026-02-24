/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package util

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCurrentTimeStr(t *testing.T) {
	now := time.Now()
	val := GetCurrentTimeStr()
	valTime, err := time.ParseInLocation("2006-01-02 15:04:05", val, time.Local)
	require.NoError(t, err)
	require.InDelta(t, now.Unix(), valTime.Unix(), 1)
}

func TestConvParamsTime(t *testing.T) {
	// strJSON := `{"page":{"start":0,"limit":10,"sort":"bk_host_id"},"pattern":"","bk_biz_id":2,"ip":{"flag":"bk_host_innerip|bk_host_outerip","exact":0,"data":[]},"condition":[{"bk_obj_id":"host","fields":[],"condition":[{"create_time":["2018-03-04","2018-03-17"]}]},{"bk_obj_id":"biz","fields":[],"condition":[{"field":"default","operator":"$ne","value":1}]},{"bk_obj_id":"module","fields":[],"condition":[]},{"bk_obj_id":"set","fields":[],"condition":[]}]}`
	strJSON := `{"bk_host_id":{"$in":[99,100,101,102,103,104]},"create_time":{"$in":["2018-03-16 02:45:28","2018-03-16"]}}`
	var a interface{}
	err := json.Unmarshal([]byte(strJSON), &a)
	if nil != err {
		t.Error(err.Error())
	}
	fmt.Println("====================")
	a = ConvParamsTime(a)
	fmt.Println(a)

}

func TestFormatPeriod(t *testing.T) {
	period := "000290S"
	periodFormated, err := FormatPeriod(period)
	if nil != err {
		t.Error(err.Error())
		return
	}
	if periodFormated != "290S" {
		t.Errorf("error formated period %s", periodFormated)
	}
	fmt.Println(periodFormated)
}

func TestConvertTimeToUserTZ(t *testing.T) {
	tests := []struct {
		name     string
		val      string
		timeZone string
		want     string
		wantErr  bool
	}{
		{
			name:     "UTC ISO8601 Z suffix to Asia/Shanghai (+8)",
			val:      "2019-04-28T09:43:13Z",
			timeZone: "Asia/Shanghai",
			want:     "2019-04-28 17:43:13+0800",
		},
		{
			name:     "UTC ISO8601 fractional seconds to Asia/Shanghai (+8)",
			val:      "2019-04-28T09:43:13.985Z",
			timeZone: "Asia/Shanghai",
			want:     "2019-04-28 17:43:13+0800",
		},
		{
			name:     "UTC ISO8601 to America/New_York (EDT -4 in April)",
			val:      "2019-04-28T09:43:13Z",
			timeZone: "America/New_York",
			want:     "2019-04-28 05:43:13-0400",
		},
		{
			name:     "UTC ISO8601 to Europe/London (BST +1 in April)",
			val:      "2019-04-28T09:43:13Z",
			timeZone: "Europe/London",
			want:     "2019-04-28 10:43:13+0100",
		},
		{
			name:     "Europe/London in winter (GMT +0)",
			val:      "2026-02-10T03:44:50.993Z",
			timeZone: "Europe/London",
			want:     "2026-02-10 03:44:50+0000",
		},
		{
			name:     "ISO8601 with offset to Asia/Shanghai",
			val:      "2026-02-10T05:11:10+08:00",
			timeZone: "Asia/Shanghai",
			want:     "2026-02-10 05:11:10+0800",
		},
		{
			name:     "ISO8601 with offset to America/Los_Angeles (PST -8 in Feb)",
			val:      "2026-02-10T05:11:10+08:00",
			timeZone: "America/Los_Angeles",
			// UTC = 2026-02-09T21:11:10Z, LA PST = UTC-8 => 2026-02-09 13:11:10
			want: "2026-02-09 13:11:10-0800",
		},
		{
			name:     "timezone-naive format as UTC to Asia/Tokyo (+9)",
			val:      "2019-04-28 09:43:13",
			timeZone: "Asia/Tokyo",
			want:     "2019-04-28 18:43:13+0900",
		},
		{
			name:     "UTC midnight cross-day to Asia/Shanghai (+8)",
			val:      "2019-04-27T23:00:00Z",
			timeZone: "Asia/Shanghai",
			want:     "2019-04-28 07:00:00+0800",
		},
		{
			name:     "empty value returns empty",
			val:      "",
			timeZone: "Asia/Shanghai",
			want:     "",
		},
		{
			name:     "empty timezone returns original",
			val:      "2019-04-28T09:43:13Z",
			timeZone: "",
			want:     "2019-04-28T09:43:13Z",
		},
		{
			name:     "invalid timezone returns error",
			val:      "2019-04-28T09:43:13Z",
			timeZone: "Invalid/TimeZone",
			wantErr:  true,
		},
		{
			name:     "unparseable value returns error",
			val:      "not-a-time",
			timeZone: "Asia/Shanghai",
			wantErr:  true,
		},
		{
			name:     "high precision nano to Asia/Shanghai",
			val:      "2019-04-28T09:43:13.123456789Z",
			timeZone: "Asia/Shanghai",
			want:     "2019-04-28 17:43:13+0800",
		},
		{
			name:     "UTC to UTC remains same",
			val:      "2019-04-28T09:43:13Z",
			timeZone: "UTC",
			want:     "2019-04-28 09:43:13+0000",
		},
		{
			name:     "negative offset: Pacific/Honolulu (UTC-10)",
			val:      "2019-04-28T09:43:13Z",
			timeZone: "Pacific/Honolulu",
			want:     "2019-04-27 23:43:13-1000",
		},
		{
			name:     "half-hour offset: Asia/Kolkata (UTC+5:30)",
			val:      "2019-04-28T09:43:13Z",
			timeZone: "Asia/Kolkata",
			want:     "2019-04-28 15:13:13+0530",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertTimeToUserTZ(tt.val, tt.timeZone)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
