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
	"encoding/base64"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"configcenter/src/common"
)

// ConfigAdminResult the result with config admin
type ConfigAdminResult struct {
	BaseResp `json:",inline"`
	Data     ConfigAdmin `json:"data"`
}

// ConfigAdminParams used to admin the cmdb config
type ConfigAdminParmams struct {
	Config ConfigAdmin `json:"config"`
}

// ConfigAdmin used to admin the cmdb config
type ConfigAdmin struct {
	Backend         BackendCfg         `json:"backend"`
	Site            SiteCfg            `json:"site"`
	ValidationRules ValidationRulesCfg `json:"validationRules"`
}

// Validate validate the fields of ConfigAdmin
func (c *ConfigAdmin) Validate() error {
	vr := reflect.ValueOf(*c)
	vrt := reflect.TypeOf(*c)
	for i := 0; i < vr.NumField(); i++ {
		field := vr.Field(i)
		funcName := []string{"Validate"}
		for _, fn := range funcName {
			vf := field.MethodByName(fn)
			errVal := vf.Call(make([]reflect.Value, 0))
			if errVal[0].Interface() != nil {
				return fmt.Errorf("%s %s failed, error:%s", vrt.Field(i).Name, fn, errVal[0].Interface().(error).Error())
			}
		}
	}
	return nil
}

// EncodeWithBase64 encode the value of ValidationRules to base64
func (c *ConfigAdmin) EncodeWithBase64() error {
	vr := reflect.ValueOf(&c.ValidationRules).Elem()
	vrt := reflect.TypeOf(c.ValidationRules)
	for i := 0; i < vr.NumField(); i++ {
		field := vr.Field(i)
		bc := field.FieldByName("BaseCfgItem")
		value := bc.Interface().(BaseCfgItem).Value
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("%s can't be empty", vrt.Field(i).Name)
		}
		base64Val := base64.StdEncoding.EncodeToString([]byte(value))
		bc.FieldByName("Value").SetString(base64Val)
	}
	return nil
}

// DecodeWithBase64 decode the base64 value of ValidationRules
func (c *ConfigAdmin) DecodeWithBase64() error {
	vr := reflect.ValueOf(&c.ValidationRules).Elem()
	vrt := reflect.TypeOf(c.ValidationRules)
	for i := 0; i < vr.NumField(); i++ {
		field := vr.Field(i)
		bc := field.FieldByName("BaseCfgItem")
		value := bc.Interface().(BaseCfgItem).Value
		bytes, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			return fmt.Errorf("%s is invalid,base64 decode err: %v", vrt.Field(i).Name, err)
		}
		bc.FieldByName("Value").SetString(string(bytes))
	}
	return nil
}

// BackendCfg used to admin backend Config
type BackendCfg struct {
	SnapshotBizName string `json:"snapshotBizName"`
	MaxBizTopoLevel int64  `json:"maxBizTopoLevel"`
}

// Validate validate the fields of BackendCfg
func (b BackendCfg) Validate() error {
	if strings.TrimSpace(b.SnapshotBizName) == "" {
		return fmt.Errorf("snapshotBizName value can't be empty")
	}
	if b.MaxBizTopoLevel < 3 || b.MaxBizTopoLevel > 10 {
		return fmt.Errorf("maxBizTopoLevel value must in range [3-10]")
	}
	return nil
}

// BaseCfgItem the common base config item
type BaseCfgItem struct {
	Value       string `json:"value"`
	Description string `json:"description"`
	I18N        I18N   `json:"i18n"`
}

// ValidateValueFormat validate the value format
func (b *BaseCfgItem) ValidateValueFormat() error {
	if b.IsEmpty() {
		return fmt.Errorf("value cant't be empty")
	}
	if b.IsExceedMaxLength() {
		return fmt.Errorf("value length can't exceed %s", common.AttributeOptionMaxLength)
	}
	return nil
}

// IsEmpty judge whether the value is empty
func (b *BaseCfgItem) IsEmpty() bool {
	if strings.TrimSpace(b.Value) == "" {
		return true
	}
	return false
}

// IsExceedMaxLength judge whether the the length of value exceed the max value
func (b *BaseCfgItem) IsExceedMaxLength() bool {
	if len(strings.TrimSpace(b.Value)) > common.AttributeOptionMaxLength {
		return true
	}
	return false
}

// ValidateRegex validate regex
func (b *BaseCfgItem) ValidateRegex() error {
	bytes, err := base64.StdEncoding.DecodeString(b.Value)
	if err != nil {
		return fmt.Errorf("%#v is invalid,base64 decode err: %v", *b, err)
	}
	reg := string(bytes)
	// convert Chinese characters to satisfy go syntax
	expr := strings.Replace(reg, "\\u4e00-\\u9fa5", "\u4e00-\u9fa5", -1)
	if _, err := regexp.Compile(expr); err != nil {
		return fmt.Errorf("%s is not a valid regular expression，%s", reg, err.Error())
	}
	return nil
}

type I18N struct {
	CN string `json:"cn"`
	EN string `json:"en"`
}

// SiteCfg used to admin Site Config
type SiteCfg struct {
	Title  TitleItem  `json:"title"`
	Footer FooterItem `json:"footer"`
}

// Validate validate the fields of SiteCfg
func (s SiteCfg) Validate() error {
	if err := s.Title.ValidateValueFormat(); err != nil {
		return fmt.Errorf("title format err:%s", err.Error())
	}
	if err := s.Footer.Validate(); err != nil {
		return fmt.Errorf("footer validate err:%s", err.Error())
	}
	return nil
}

type TitleItem struct {
	BaseCfgItem `json:",inline"`
}

type FooterItem struct {
	Links []LinksItem `json:"links"`
}

type LinksItem struct {
	BaseCfgItem `json:",inline"`
	Enabled     bool `json:"enabled"`
}

// Validate validate the fields of FooterItem
func (s FooterItem) Validate() error {
	if len(s.Links) == 0 {
		return fmt.Errorf("links can't be empty")
	}
	for _, link := range s.Links {
		if err := link.ValidateValueFormat(); err != nil {
			return fmt.Errorf("link %#v ValidateValueFormat err, %s", link, err.Error())
		}
	}
	return nil
}

// ValidationRulesCfg used to admin valiedation rules Config
type ValidationRulesCfg struct {
	Number                NumberItem                `json:"number"`
	Float                 FloatItem                 `json:"float"`
	Singlechar            SinglecharItem            `json:"singlechar"`
	Longchar              LongcharItem              `json:"longchar"`
	AssociationId         AssociationIdItem         `json:"associationId"`
	ClassifyId            ClassifyIdItem            `json:"classifyId"`
	ModelId               ModelIdItem               `json:"modelId"`
	EnumId                EnumIdItem                `json:"enumId"`
	EnumName              EnumNameItem              `json:"enumName"`
	FieldId               FieldIdItem               `json:"fieldId"`
	NamedCharacter        NamedCharacterItem        `json:"namedCharacter"`
	InstanceTagKey        InstanceTagKeyItem        `json:"instanceTagKey"`
	InstanceTagValue      InstanceTagValueItem      `json:"instanceTagValue"`
	BusinessTopoInstNames BusinessTopoInstNamesItem `json:"businessTopoInstNames"`
}

// Validate validate the fields of ValidationRulesCfg
func (v ValidationRulesCfg) Validate() error {
	vr := reflect.ValueOf(v)
	vrt := reflect.TypeOf(v)
	for i := 0; i < vr.NumField(); i++ {
		field := vr.Field(i)
		bc := field.FieldByName("BaseCfgItem").Interface().(BaseCfgItem)
		bcr := reflect.ValueOf(&bc)
		funcName := []string{"ValidateValueFormat", "ValidateRegex"}
		for _, fn := range funcName {
			vf := bcr.MethodByName(fn)
			errVal := vf.Call(make([]reflect.Value, 0))
			if errVal[0].Interface() != nil {
				return fmt.Errorf("%s %s failed, error:%s", vrt.Field(i).Name, fn, errVal[0].Interface().(error).Error())
			}
		}
	}
	return nil
}

type NumberItem struct {
	BaseCfgItem `json:",inline"`
}

type FloatItem struct {
	BaseCfgItem `json:",inline"`
}

type SinglecharItem struct {
	BaseCfgItem `json:",inline"`
}

type LongcharItem struct {
	BaseCfgItem `json:",inline"`
}

type AssociationIdItem struct {
	BaseCfgItem `json:",inline"`
}

type ClassifyIdItem struct {
	BaseCfgItem `json:",inline"`
}

type ModelIdItem struct {
	BaseCfgItem `json:",inline"`
}

type EnumIdItem struct {
	BaseCfgItem `json:",inline"`
}

type EnumNameItem struct {
	BaseCfgItem `json:",inline"`
}

type FieldIdItem struct {
	BaseCfgItem `json:",inline"`
}

type NamedCharacterItem struct {
	BaseCfgItem `json:",inline"`
}

type InstanceTagKeyItem struct {
	BaseCfgItem `json:",inline"`
}

type InstanceTagValueItem struct {
	BaseCfgItem `json:",inline"`
}

type BusinessTopoInstNamesItem struct {
	BaseCfgItem `json:",inline"`
}
