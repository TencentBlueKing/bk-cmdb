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

package time

import (
	"strings"
	"time"

	"configcenter/src/common/json"
	"configcenter/src/common/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type Time struct {
	time.Time `bson:",inline" json:",inline"`
}

func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Time) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if len(data) == 0 || string(data) == "null" {
		return nil
	}

	dataStr := strings.Trim(string(data), "\"")

	timeType, isTime := util.IsTime(dataStr)
	if isTime {
		t.Time = util.Str2Time(dataStr, timeType)
		return nil
	}

	return json.Unmarshal(data, &t.Time)
}

// MarshalBSONValue implements bson.MarshalBSON interface
func (t Time) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsonx.Time(t.Time).MarshalBSONValue()
}

// UnmarshalBSONValue implements bson.UnmarshalBSONValue interface
func (t *Time) UnmarshalBSONValue(typo bsontype.Type, raw []byte) error {
	switch typo {
	case bsontype.DateTime:
		rv := bson.RawValue{Type: bsontype.DateTime, Value: raw}
		t.Time = rv.Time()
		return nil
	case bsontype.String:
		rawStr := bson.RawValue{Type: bsontype.String, Value: raw}
		return t.UnmarshalJSON([]byte(rawStr.String()))
	}

	return bson.Unmarshal(raw, &t.Time)
}
