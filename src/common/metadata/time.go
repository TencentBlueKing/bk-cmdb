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
	"strconv"
	"strings"
	"time"

	"database/sql/driver"

	"github.com/coccyx/timeparser"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type Time struct {
	time.Time `bson:",inline" json:",inline"`
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
	return []byte(t.Format(`"2006-01-02 15:04:05"`)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be a quoted string in RFC 3339 format.
func (t *Time) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}

	parsed, err := time.ParseInLocation(`"2006-01-02 15:04:05"`, string(data), time.UTC)
	if err == nil {
		*t = Time{parsed}
		return nil
	}

	parsed, err = time.Parse(time.RFC3339, strings.Trim(string(data), "\""))
	if err == nil {
		*t = Time{parsed}
		return nil
	}

	parsed, err = timeparser.TimeParser(strings.Trim(string(data), "\""))
	if err == nil {
		*t = Time{parsed}
		return nil
	}

	timestamp, err := strconv.ParseInt(fmt.Sprintf("%s", data), 10, 64)
	if err == nil {
		*t = Time{time.Unix(timestamp, 0)}
	}
	return fmt.Errorf("parse unknow time format: %s, %v", data, err)
}

// MarshalBSONValue implements bson.MarshalBSON interface
func (t Time) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsonx.Time(t.Time).MarshalBSONValue()
}

// UnmarshalBSONValue implements bson.UnmarshalBSONValue interface
func (t *Time) UnmarshalBSONValue(typo bsontype.Type, raw []byte) error {
	switch typo {
	case bsontype.Timestamp:
		// 0x11 timestamp
		return bson.Unmarshal(raw, &t.Time)
	case bsontype.Double:
		rv := bson.RawValue{Type: bsontype.Double, Value: raw}
		if dt, ok := rv.DoubleOK(); ok {
			t.Time = time.Unix(int64(dt/1000), int64(uint64(dt)%1000*1000000))
			return nil
		}
		return nil
	case bsontype.DateTime:
		rv := bson.RawValue{Type: bsontype.DateTime, Value: raw}
		t.Time = rv.Time()
		return nil
	case bsontype.String:
		rawStr := bson.RawValue{Type: bsontype.String, Value: raw}
		strTime := strings.TrimSpace(strings.Trim(rawStr.String(), "\""))

		vTime, err := time.Parse(time.RFC3339Nano, strTime)
		if err == nil {
			t.Time = vTime
			return nil
		}

		vTime, err = time.Parse(time.RFC3339, strTime)
		if err == nil {
			t.Time = vTime
			return nil
		}

		vTime, err = timeparser.TimeParser(strTime)
		if err == nil {
			t.Time = vTime
			return nil
		}

		return fmt.Errorf("cannot decode %v into a metadata.Time, err:%s", bsontype.String, err.Error())

	}

	// for compatibility purpose
	tt := tmptime{}
	err := bson.Unmarshal(raw, &tt)
	t.Time = tt.Time

	return err
}

type tmptime struct {
	time.Time
}

// Now retruns now
func Now() Time {
	return Time{time.Now().UTC()}
}

var local = Now()

func StringToTimeDurationHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		isLocalTimeType := t == reflect.TypeOf(local)
		isStdTimeType := t == reflect.TypeOf(local.Time)
		if !isLocalTimeType && !isStdTimeType {
			return data, nil
		}
		parsed, err := ParseTime(data)
		if err != nil {
			return nil, err
		}
		if isLocalTimeType {
			return Time{parsed}, nil
		}
		return parsed, nil
	}
}

func ParseTime(data interface{}) (time.Time, error) {
	// Convert it by parsing
	var parsed time.Time
	var err error
	parsed, err = time.ParseInLocation(`"2006-01-02 15:04:05"`, data.(string), time.UTC)
	if err == nil {
		return parsed, nil
	}

	parsed, err = time.Parse(time.RFC3339, strings.Trim(data.(string), "\""))
	if err == nil {
		return parsed, nil
	}

	parsed, err = timeparser.TimeParser(strings.Trim(data.(string), "\""))
	if err == nil {
		return parsed, nil
	}

	timestamp, err := strconv.ParseInt(fmt.Sprintf("%s", data), 10, 64)
	if err == nil {
		return time.Unix(timestamp, 0), nil
	}
	return time.Now(), err
}

// ParseTimeInUnixTS return Time object base on unix timestamp.
func ParseTimeInUnixTS(sec, nsec int64) Time {
	return Time{time.Unix(sec, nsec)}
}
