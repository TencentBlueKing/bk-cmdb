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
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"configcenter/src/common/valid/attribute/manager/register"
)

func init() {
	// Register the document attribute type
	register.Register(document{})
}

// document represents a document  attribute type.
type document struct {
}

// Name returns the name of the document attribute.
func (d document) Name() string {
	return "document"
}

// DisplayName returns the display name for user.
func (d document) DisplayName() string {
	return "文件"
}

// RealType returns the db type of the document attribute.
func (d document) RealType() string {
	return common.FieldTypeLongChar
}

// Info returns the tips for user.
func (d document) Info() string {
	return "文件,图片类型字段"
}

// Validate validate option and value
func (d document) Validate(ctx context.Context, objID string, propertyType string, required bool, option interface{}, value interface{}) error {

	val, err := d.ParseValue(value)
	if err != nil {
		return fmt.Errorf("document attribute %s.%s value must be a string, got %T", objID, propertyType, value)
	}
	if len(val.Value) == 0 {
		if required {
			blog.Errorf("document attribute %s.%s value is required, but got empty, rid: %s", objID, propertyType,
				util.ExtractRequestIDFromContext(ctx))
			return fmt.Errorf("document attribute %s.%s value is required, but got empty", objID, propertyType)
		}
		return nil
	}

	if len(val.Name) > common.FieldTypeLongLenChar || len(val.Value) > common.FieldTypeLongLenChar {
		return fmt.Errorf("document attribute %s.%s value length exceeds the maximum limit of %d characters", objID,
			propertyType, common.FieldTypeLongLenChar)
	}

	rid := util.ExtractRequestIDFromContext(ctx)
	// option compatible with the scene where the option is not set in the model attribute.
	dOption, err := d.ParseOption(option)
	if err != nil {
		blog.Errorf("parse document option failed, option: %v, error: %v, rid: %s", option, err, rid)
		return fmt.Errorf("document option is not a valid DocumentOption type: %v, error: %v", option, err)
	}
	match, err := regexp.MatchString(dOption.Regex, val.Value)
	if err != nil || !match {
		blog.Errorf("default value %s not matches string option %s, err: %v, rid: %s", val, dOption.Regex, err, rid)
		return fmt.Errorf("string default value not match regex")
	}

	return nil
}

// DocumentOption defines validation rules for a document or file.
// It specifies constraints such as allowed file suffixes, maximum file size,
// filename matching rules, and the logical file type.
type DocumentOption struct {
	// AllowSuffixes specifies the allowed file suffixes (without dot), e.g. ["jpg", "png", "pdf"].
	// If empty, no suffix restriction is applied.
	AllowSuffixes []string `json:"allow_suffixes,omitempty"`

	// AllowSize specifies the maximum allowed file size in bytes.
	// If zero, no size limit is enforced.
	AllowSize int64 `json:"allow_size,omitempty"`

	// Regex defines an optional regular expression used to validate
	// the file name or other string fields depending on the context.
	Regex string `json:"regex,omitempty"`

	// Type indicates the logical category of the file, such as image,
	// document, video, or audio. The value should exist in DocumentOptionTypeRela.
	Type string `json:"type,omitempty"`
}

// DocumentOptionTypeRela defines the supported file type categories.
// It is implemented as a set for efficient membership checking.
var DocumentOptionTypeRela = map[string]struct{}{
	"image":    {},
	"document": {},
	"video":    {},
	"audio":    {},
}

// DocumentValueDocument defines the document value format
type DocumentValueDocument struct {
	Value string `json:"value"`
	Name  string `json:"name"`
}

// ParseOption Parse document Option
func (d document) ParseOption(option interface{}) (DocumentOption, error) {
	if option == nil {
		return DocumentOption{}, nil
	}

	switch val := option.(type) {
	case *DocumentOption:
		return *val, nil
	case DocumentOption:
		return val, nil

	}
	optBytes, err := json.Marshal(option)
	if err != nil {
		return DocumentOption{}, fmt.Errorf("document option is not a valid DocumentOption type: %v", option)
	}
	res := DocumentOption{}
	if err := json.Unmarshal(optBytes, &res); err != nil {
		return DocumentOption{}, fmt.Errorf("document option is not a valid DocumentOption type: %v, error: %v",
			option, err)
	}

	return res, nil
}

// ParseValue Parse document Value
func (d document) ParseValue(value interface{}) (DocumentValueDocument, error) {
	if value == nil {
		return DocumentValueDocument{}, nil
	}

	switch val := value.(type) {
	case *DocumentValueDocument:
		return *val, nil
	case DocumentValueDocument:
		return val, nil
	}

	valBytes, err := json.Marshal(value)
	if err != nil {
		return DocumentValueDocument{}, fmt.Errorf("document value is not a valid value type: %v", value)
	}
	res := DocumentValueDocument{}
	if err := json.Unmarshal(valBytes, &res); err != nil {
		return DocumentValueDocument{}, fmt.Errorf("document value is not a valid value type: %v, error: %v", value,
			err)
	}

	return res, nil
}

// FillLostValue Fill document LostValue
func (d document) FillLostValue(ctx context.Context, valData mapstr.MapStr, name string,
	defaultValue, option interface{}) error {

	rid := util.ExtractRequestIDFromContext(ctx)

	valData[name] = nil
	if defaultValue == nil {
		return nil
	}

	defaultVal, err := d.ParseValue(defaultValue)
	if err != nil {
		return fmt.Errorf("single char default value not string, value: %v, rid: %s", defaultValue, rid)
	}

	if len(defaultVal.Value) == 0 && len(defaultVal.Name) == 0 {
		return nil
	}

	// option compatible with the scene where the option is not set in the model attribute.
	dOption, err := d.ParseOption(option)
	if err != nil {
		return err
	}

	match, err := regexp.MatchString(dOption.Regex, defaultVal.Value)
	if err != nil || !match {
		return fmt.Errorf("the current string does not conform to regular verification rules")
	}
	valData[name] = defaultVal
	return nil
}

// ValidateOption Validate document Option
func (d document) ValidateOption(ctx context.Context, option interface{}, defaultVal interface{}) error {

	rid := util.ExtractRequestIDFromContext(ctx)
	// option compatible with the scene where the option is not set in the model attribute.
	dOption, err := d.ParseOption(option)
	if err != nil {
		return err
	}

	// allow suffixes is required, if not set, return error
	if len(dOption.AllowSuffixes) == 0 {
		blog.Errorf("document option allow_suffixes is required, but not set, rid: %s", rid)
		return fmt.Errorf("document option allow_suffixes is required, but not set")
	}
	if len(dOption.Type) == 0 {
		blog.Errorf("document option type is required, but not set, rid: %s", rid)
		return fmt.Errorf("document option type is required, but not set")
	}
	if _, ok := DocumentOptionTypeRela[dOption.Type]; !ok {
		blog.Errorf("document option type %s is not supported, rid: %s", dOption.Type, rid)
		return fmt.Errorf("document option type %s is not supported", dOption.Type)
	}

	if len(dOption.Regex) == 0 {
		return nil
	}

	value := DocumentValueDocument{}
	if defaultVal != nil {
		value, err = d.ParseValue(defaultVal)
		if err != nil {
			blog.Errorf("string type default value %+v type %T is invalid, err: %s, rid: %s", defaultVal, defaultVal, err, rid)
			return fmt.Errorf("field default value, not string type")
		}
	}

	if _, err := regexp.Compile(dOption.Regex); err != nil {
		blog.Errorf("regular expression %s is invalid, err: %, rid: %s", dOption.Regex, err, rid)
		return fmt.Errorf("regular is wrong")
	}

	if defaultVal == nil {
		return nil
	}

	match, err := regexp.MatchString(dOption.Regex, value.Value)
	if err != nil || !match {
		blog.Errorf("default value %s not matches string option %s, err: %v, rid: %s", value.Value, dOption.Regex, err, rid)
		return fmt.Errorf("string default value not match regex")
	}

	return nil
}

var _ register.AttributeTypeI = &document{}
