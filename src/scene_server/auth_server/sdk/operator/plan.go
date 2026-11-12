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

package operator

import (
	"bytes"
	"encoding/json"

	ccjson "configcenter/src/common/json"
)

// PlanKind is the authorization result kind of the hybrid plan apis.
type PlanKind string

const (
	// AlwaysAllowedKind means that the subject can access all the resources.
	AlwaysAllowedKind PlanKind = "ALWAYS_ALLOWED"
	// AlwaysDeniedKind means that the subject can access none of the resources.
	AlwaysDeniedKind PlanKind = "ALWAYS_DENIED"
	// ConditionalKind means that the resources should be filtered by the returned condition.
	ConditionalKind PlanKind = "CONDITIONAL"
)

// Plan is one action's authorization plan returned by the hybrid plan apis.
type Plan struct {
	Kind PlanKind `json:"kind"`
	// Condition is the filter expression when Kind is ConditionalKind, otherwise it is nil.
	Condition *AuthCondition `json:"condition"`
}

// AuthCondition is a node of the authorization plan expression.
type AuthCondition struct {
	Operator OperType `json:"op"`
	// Element is a pointer interface point to the implements struct,
	// which should be one of Content or FieldValue.
	Element
}

// UnmarshalJSON unmarshal the authorization plan condition from the standard expression protocol.
func (c *AuthCondition) UnmarshalJSON(i []byte) error {
	if string(i) == "{}" {
		return nil
	}

	broker := new(conditionBroker)
	err := ccjson.Unmarshal(i, broker)
	if err != nil {
		return err
	}

	c.Operator = broker.Operator

	if broker.Operator.IsLogical() {
		content := new(Content)
		if err := ccjson.Unmarshal(broker.Content, &content.Content); err != nil {
			return err
		}
		c.Element = content
		return nil
	}

	if broker.Operator == In {
		to := make([]interface{}, 0)
		if err := ccjson.Unmarshal(broker.Value, &to); err != nil {
			return err
		}

		c.Element = &FieldValue{
			Field: broker.Field,
			Value: to,
		}
		return nil
	}

	to := new(interface{})
	if err := ccjson.Unmarshal(broker.Value, &to); err != nil {
		return err
	}

	c.Element = &FieldValue{
		Field: broker.Field,
		Value: *to,
	}

	return nil
}

type conditionBroker struct {
	Operator OperType        `json:"op"`
	Content  json.RawMessage `json:"content"`
	Field    Field           `json:"field"`
	Value    json.RawMessage `json:"value"`
}

// MarshalJSON is used to marshal the condition to the standard
// iam policy protocol, which is not correspond to the struct
// we defined here.
// Note: when you marshal the condition, the condition must be a pointer,
// otherwise, the marshaled json struct is wrong.
func (c *AuthCondition) MarshalJSON() ([]byte, error) {
	js, err := ccjson.Marshal(c.Element)
	if err != nil {
		return nil, err
	}
	buf := bytes.Buffer{}
	buf.WriteString(`{"op":"`)
	buf.WriteString(string(c.Operator))
	buf.WriteString(`",`)
	buf.Write(js[1 : len(js)-1])
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// Element is the payload of an AuthCondition node.
type Element interface {
	EleName() string
}

// Content TODO
type Content struct {
	// Content is only exist when OperType is a logical operator.
	Content []*AuthCondition `json:"content"`
}

// EleName TODO
func (e *Content) EleName() string {
	return "content"
}

// FieldValue is a compare node of the authorization plan expression.
type FieldValue struct {
	// Field and Value is only exist when OperType is not a logical operator.
	Field Field       `json:"field"`
	Value interface{} `json:"value"`
}

// EleName TODO
func (f *FieldValue) EleName() string {
	return "field_value"
}
