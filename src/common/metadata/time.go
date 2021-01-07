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
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"
)

type Time struct {
	time.Time
}

// Scan implement sql driver's Scan interface
func (t *Time) Scan(value interface{}) error {
	t.Time = value.(time.Time)
	return nil
}

// Value implement sql driver's Value interface
func (t Time) Value() (driver.Value, error) {
	return t.Time, nil
}

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(t.UTC().Format(`"2006-01-02T15:04:05Z"`)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be a quoted string in RFC 3339 format.
func (t *Time) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}
	parsed, err := time.Parse(`"`+time.RFC3339+`"`, string(data))
	if err == nil {
		*t = Time{parsed}
		return nil
	}
	parsed, err = time.ParseInLocation(`"2006-01-02 15:04:05"`, string(data), time.UTC)
	if err == nil {
		*t = Time{parsed}
		return nil
	}

	timestamp, err := strconv.ParseInt(fmt.Sprintf("%s", data), 10, 64)
	if err == nil {
		*t = Time{time.Unix(timestamp, 0)}
	}
	return err
}

func Now() Time {
	return Time{time.Now()}
}
