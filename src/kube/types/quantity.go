/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package types

import (
	"fmt"
	"math/big"

	"configcenter/src/common/json"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// Scale is used for getting and setting the base-10 scaled value.
type Scale int32

// Format describes how a quantity is formatted.
type Format string

// Quantity represents a Kubernetes resource quantity. Its JSON and BSON
// representations are strings such as "500m" and "128Mi".
type Quantity struct {
	i int64Amount
	d infDecAmount
	s string
	Format
}

type int64Amount struct {
	value int64
	scale Scale
}

type infDecAmount struct {
	*Dec
}

// Dec represents a signed arbitrary-precision decimal.
type Dec struct {
	unscaled big.Int
	scale    Scale
}

// String returns the string representation of the quantity.
func (q Quantity) String() string {
	return q.s
}

// MarshalJSON implements the json.Marshaler interface.
func (q Quantity) MarshalJSON() ([]byte, error) {
	return json.Marshal(q.s)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (q *Quantity) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*q = Quantity{}
		return nil
	}

	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	q.s = value
	return nil
}

// MarshalBSONValue implements the bson.ValueMarshaler interface.
func (q Quantity) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(q.String())
}

// UnmarshalBSONValue implements the bson.ValueUnmarshaler interface.
func (q *Quantity) UnmarshalBSONValue(typ bsontype.Type, raw []byte) error {
	if typ == bson.TypeNull {
		*q = Quantity{}
		return nil
	}

	value, ok := bson.RawValue{Type: typ, Value: raw}.StringValueOK()
	if !ok {
		return fmt.Errorf("cannot decode BSON type %s into Quantity", typ)
	}

	q.s = value
	return nil
}
