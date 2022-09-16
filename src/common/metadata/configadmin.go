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

// PlatformSettingResult the result of platform setting.
type PlatformSettingResult struct {
	BaseResp `json:",inline"`
	Data     PlatformSettingConfig `json:"data"`
}

// ConfigAdminParmams used to admin the cmdb config
type ConfigAdminParmams struct {
	Config ConfigAdmin `json:"config"`
}

// ConfigAdmin used to admin the cmdb config
type ConfigAdmin struct {
	Backend         BackendCfg         `json:"backend"`
	Site            SiteCfg            `json:"site"`
	ValidationRules ValidationRulesCfg `json:"validationRules"`
}

// Site site name and separator config.
type Site struct {
	SiteName  TitleItem `json:"name"`
	Separator string    `json:"separator"`
}

// Validate validate the fields of SiteCfg.
func (s Site) Validate() error {
	if strings.TrimSpace(s.SiteName.I18N.CN) == "" {
		return fmt.Errorf("site cn name can't be empty")
	}
	if strings.TrimSpace(s.SiteName.I18N.EN) == "" {
		return fmt.Errorf("site en name can't be empty")
	}

	if strings.TrimSpace(s.Separator) == "" {
		return fmt.Errorf("separator value can't be empty")
	}
	return nil
}

// ContactInfoItem contact information, markdown format.
type ContactInfoItem struct {
	BaseCfgItem `json:",inline"`
}

// CopyrightItem copyright information, markdown format.
type CopyrightItem struct {
	BaseCfgItem `json:",inline"`
}

// Footer footer information.
type Footer struct {
	ContactInfo ContactInfoItem `json:"contact"`
	Copyright   CopyrightItem   `json:"copyright"`
}

// Validate validate the fields of SiteCfg
func (s Footer) Validate() error {
	if strings.TrimSpace(s.ContactInfo.I18N.CN) == "" {
		return fmt.Errorf("contact info cn value can't be empty")
	}
	if strings.TrimSpace(s.ContactInfo.I18N.EN) == "" {
		return fmt.Errorf("contact info en value can't be empty")
	}
	if strings.TrimSpace(s.Copyright.I18N.CN) == "" {
		return fmt.Errorf("copyright cn value can't be empty")
	}
	if strings.TrimSpace(s.Copyright.I18N.EN) == "" {
		return fmt.Errorf("copyright en value can't be empty")
	}
	return nil
}

// UserModuleList custom section.
type UserModuleList struct {
	Key   string `json:"module_key"`
	Value string `json:"module_name"`
}

// GlobalModule Conifg, idleName, FaultName and RecycleName cannot be deleted.
type GlobalModule struct {
	IdleName    string           `json:"idle"`
	FaultName   string           `json:"fault"`
	RecycleName string           `json:"recycle"`
	UserModules []UserModuleList `json:"user_modules"`
}

// Validate validate the fields of IdleModule.
func (s GlobalModule) Validate() error {
	if strings.TrimSpace(s.RecycleName) == "" {
		return fmt.Errorf("site  value can't be empty")
	}
	if strings.TrimSpace(s.IdleName) == "" {
		return fmt.Errorf("separator  value can't be empty")
	}
	if strings.TrimSpace(s.FaultName) == "" {
		return fmt.Errorf("separator  value can't be empty")
	}
	return nil
}

// AdminBackendCfg TODO
type AdminBackendCfg struct {
	MaxBizTopoLevel int64  `json:"max_biz_topo_level"`
	SnapshotBizName string `json:"snapshot_biz_name"`
}

// Validate validate the fields of BackendCfg.
func (b AdminBackendCfg) Validate() error {

	if strings.TrimSpace(b.SnapshotBizName) == "" {
		return fmt.Errorf("snapshot biz name can't be empty")
	}

	if b.MaxBizTopoLevel < minBizTopoLevel || b.MaxBizTopoLevel > maxBizTopoLevel {
		return fmt.Errorf("max biz topo level value must in range [%d-%d]", minBizTopoLevel, maxBizTopoLevel)
	}
	return nil
}

// ObjectString TODO
type ObjectString string

// Validate validate the fields of ObjectString
func (s ObjectString) Validate() error {
	if strings.TrimSpace(string(s)) == "" {
		return fmt.Errorf("site  value can't be empty")
	}
	return nil
}

// PlatformSettingConfig  used to admin the platform config. 结构体PlatformSettingConfig 每个成员对象都必须有"Validate"校验
// 函数，如果没有会panic.
type PlatformSettingConfig struct {
	Backend             AdminBackendCfg    `json:"backend"`
	SiteConfig          Site               `json:"site"`
	FooterConfig        Footer             `json:"footer"`
	ValidationRules     ValidationRulesCfg `json:"validation_rules"`
	BuiltInSetName      ObjectString       `json:"set"`
	BuiltInModuleConfig GlobalModule       `json:"idle_pool"`
}

// Validate validate the fields of PlatformSettingReqOption is illegal .
func (c *PlatformSettingConfig) Validate() error {
	vr := reflect.ValueOf(*c)
	vrt := reflect.TypeOf(*c)
	for i := 0; i < vr.NumField(); i++ {
		field := vr.Field(i)
		funcName := []string{"Validate"}
		for _, fn := range funcName {
			vf := field.MethodByName(fn)
			errVal := vf.Call(make([]reflect.Value, 0))
			if errVal[0].Interface() != nil {
				return fmt.Errorf("%s %s failed, error: %s", vrt.Field(i).Name, fn,
					errVal[0].Interface().(error).Error())
			}
		}
	}

	return nil
}

const (
	maxBizTopoLevel = 10
	minBizTopoLevel = 3
)

// InitAdminConfig factory configuration.
var InitAdminConfig = `
{
    "backend":{
        "max_biz_topo_level":7,
        "snapshot_biz_name":"蓝鲸"
    },
    "site":{
         "name":{
            "value":"配置平台 | 蓝鲸",
            "description":"网站标题",
            "i18n":{
                "cn":"配置平台 | 蓝鲸",
                "en":"CMDB | BlueKing"
            }
        },
        "separator":"|"
    },
    "footer":{
         "contact":{
            "value":"http://127.0.0.1",
            "description":"联系BK助手",
            "i18n":{
                "cn":"http://127.0.0.1",
                "en":"http://127.0.0.1"
            }
        },
         "copyright":{
            "value":"Copyright © 2012-{{current_year}} Tencent BlueKing. All Rights Reserved.",
            "description":"版权信息",
            "i18n":{
                "cn":"Copyright © 2012-{{current_year}} Tencent BlueKing. All Rights Reserved.",
                "en":"Copyright © 2012-{{current_year}} Tencent BlueKing. All Rights Reserved."
            }
        }
    },
  "validation_rules": {
        "number": {
            "value": "^(\\-|\\+)?\\d+$",
            "description": "字段类型“数字”的验证规则",
            "i18n": {
                "cn": "请输入整数数字",
                "en": "Please enter integer number"
            }
        },
        "float": {
            "value": "^[+-]?([0-9]*[.]?[0-9]+|[0-9]+[.]?[0-9]*)([eE][+-]?[0-9]+)?$",
            "description": "字段类型“浮点”的验证规则",
            "i18n": {
                "cn": "请输入浮点型数字",
                "en": "Please enter float data"
            }
        },
        "singlechar": {
            "value": "\\S*",
            "description": "字段类型“短字符”的验证规则",
            "i18n": {
                "cn": "请输入256长度以内的字符串",
                "en": "Please enter the string within 256 length"
            }
        },
        "longchar": {
            "value": "\\S*",
            "description": "字段类型“长字符”的验证规则",
            "i18n": {
                "cn": "请输入2000长度以内的字符串",
                "en": "Please enter the string within 2000 length"
            }
        },
        "associationId": {
            "value": "^[a-zA-Z][\\w]*$",
            "description": "关联类型唯一标识验证规则",
            "i18n": {
                "cn": "由英文字符开头，和下划线、数字或英文组合的字符",
"en": "Start with lowercase or uppercase letter, followed by lowercase / uppercase / underscore / numbers characters"
            }
        },
        "classifyId": {
            "value": "^[a-zA-Z][\\w]*$",
            "description": "模型分组唯一标识验证规则",
            "i18n": {
                "cn": "由英文字符开头，和下划线、数字或英文组合的字符",
"en": "Start with lowercase or uppercase letter, followed by lowercase / uppercase / underscore / numbers characters"
            }
        },
        "modelId": {
            "value": "^[a-zA-Z][\\w]*$",
            "description": "模型唯一标识验证规则",
            "i18n": {
                "cn": "由英文字符开头，和下划线、数字或英文组合的字符",
"en": "Start with lowercase or uppercase letter, followed by lowercase / uppercase / underscore / numbers characters"
            }
        },
        "enumId": {
            "value": "^[a-zA-Z0-9_-]*$",
            "description": "字段类型“枚举”ID的验证规则",
            "i18n": {
                "cn": "由大小写英文字母，数字，_ 或 - 组成的字符",
                "en": "Composed of uppercase / lowercase / numbers / - or _ characters"
            }
        },
        "enumName": {
            "value": "^([a-zA-Z0-9_]|[\\u4e00-\\u9fa5]|[()+-《》,，；;“”‘’。\\.\\\"\\' \\/:])*$",
            "description": "字段类型“枚举”值的验证规则",
            "i18n": {
                "cn": "请输入枚举值",
                "en": "Please enter the enum value"
            }
        },
        "fieldId": {
            "value": "^[a-zA-Z][\\w]*$",
            "description": "模型字段唯一标识的验证规则",
            "i18n": {
                "cn": "由英文字符开头，和下划线、数字或英文组合的字符",
"en": "Start with lowercase or uppercase letter, followed by lowercase / uppercase / underscore / numbers characters"
            }
        },
        "namedCharacter": {
            "value": "^[a-zA-Z0-9\\u4e00-\\u9fa5_\\-:\\(\\)]+$",
            "description": "服务分类名称的验证规则",
            "i18n": {
                "cn": "请输入中英文或特殊字符 :_- 组成的名称",
                "en": "Special symbols only support(:_-)"
            }
        },
        "instanceTagKey": {
            "value": "^[a-zA-Z]([a-z0-9A-Z\\-_.]*[a-z0-9A-Z])?$",
            "description": "服务实例标签键的验证规则",
            "i18n": {
                "cn": "请输入以英文开头的英文+数字组合",
                "en": "Please enter letter / number starts with letter"
            }
        },
        "instanceTagValue": {
            "value": "^[a-z0-9A-Z]([a-z0-9A-Z\\-_.]*[a-z0-9A-Z])?$",
            "description": "服务实例标签值的验证规则",
            "i18n": {
                "cn": "请输入英文 / 数字",
                "en": "Please enter letter / number"
            }
        },
        "businessTopoInstNames": {
            "value": "^[^\\#\\/,\\>\\<\\|]+$",
            "description": "集群/模块/实例名称的验证规则",
            "i18n": {
                "cn": "请输入除 #/,><| 以外的字符",
                "en": "Please enter characters other than #/,><|"
            }
        }
    },
    "set":"空闲机池",
    "idle_pool":{
        "idle":"空闲机",
        "fault":"故障机",
        "recycle":"待回收"
    }
}
`

// BuiltInModuleDeleteOption  used to admin the idle module config
type BuiltInModuleDeleteOption struct {
	ModuleKey  string `json:"module_key"`
	ModuleName string `json:"module_name"`
}

// Validate whether the option parameter is legal.
func (option *BuiltInModuleDeleteOption) Validate() error {
	if option.ModuleKey == "" || option.ModuleName == "" {
		return fmt.Errorf("module key and name must be set")
	}

	if option.ModuleKey == common.SystemIdleModuleKey || option.ModuleKey == common.SystemFaultModuleKey ||
		option.ModuleKey == common.SystemRecycleModuleKey {
		return fmt.Errorf("systen module key can not be delete")
	}
	return nil
}

// ModuleOption  used to modify the idle module config.
type ModuleOption struct {
	Key  string `json:"module_key"`
	Name string `json:"module_name"`
}

// SetOption  used to modify the idle set config.
type SetOption struct {
	Key  string `json:"set_key"`
	Name string `json:"set_name"`
}

const (
	// ConfigUpdateTypeSet TODO
	ConfigUpdateTypeSet ConfigUpdateType = "set"
	// ConfigUpdateTypeModule TODO
	ConfigUpdateTypeModule ConfigUpdateType = "module"
)

// ConfigUpdateType TODO
type ConfigUpdateType string

// ConfigUpdateSettingOption  used to modify the idle update config
type ConfigUpdateSettingOption struct {
	Set    SetOption    `json:"set"`
	Module ModuleOption `json:"module"`

	// Type request type: set or module
	Type ConfigUpdateType `json:"type"`
}

// Validate whether the option parameter is legal.
func (option *ConfigUpdateSettingOption) Validate() error {

	switch option.Type {
	case ConfigUpdateTypeModule:
		if option.Module.Key == "" || option.Module.Name == "" {
			return fmt.Errorf("option module param error")
		}
	case ConfigUpdateTypeSet:
		if option.Set.Key == "" || option.Set.Name == "" {
			return fmt.Errorf("option set param error")
		}
	default:
		return fmt.Errorf("input param type error")
	}
	return nil
}

// RestoreSettingOption  used to restore platform config.
type RestoreSettingOption struct {
	RestoreItem string `json:"restore_item"`
}

// UserModulesSettingReqOption  used to admin the cmdb config.
type UserModulesSettingReqOption struct {
	UserModule string `json:"user_module"`
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

// EncodeWithBase64 encode the value of ValidationRules to base64.
func (c *PlatformSettingConfig) EncodeWithBase64() error {
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

// I18N TODO
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

// TitleItem TODO
type TitleItem struct {
	BaseCfgItem `json:",inline"`
}

// FooterItem TODO
type FooterItem struct {
	Links []LinksItem `json:"links"`
}

// LinksItem TODO
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

// NumberItem TODO
type NumberItem struct {
	BaseCfgItem `json:",inline"`
}

// FloatItem TODO
type FloatItem struct {
	BaseCfgItem `json:",inline"`
}

// SinglecharItem TODO
type SinglecharItem struct {
	BaseCfgItem `json:",inline"`
}

// LongcharItem TODO
type LongcharItem struct {
	BaseCfgItem `json:",inline"`
}

// AssociationIdItem TODO
type AssociationIdItem struct {
	BaseCfgItem `json:",inline"`
}

// ClassifyIdItem TODO
type ClassifyIdItem struct {
	BaseCfgItem `json:",inline"`
}

// ModelIdItem TODO
type ModelIdItem struct {
	BaseCfgItem `json:",inline"`
}

// EnumIdItem TODO
type EnumIdItem struct {
	BaseCfgItem `json:",inline"`
}

// EnumNameItem TODO
type EnumNameItem struct {
	BaseCfgItem `json:",inline"`
}

// FieldIdItem TODO
type FieldIdItem struct {
	BaseCfgItem `json:",inline"`
}

// NamedCharacterItem TODO
type NamedCharacterItem struct {
	BaseCfgItem `json:",inline"`
}

// InstanceTagKeyItem TODO
type InstanceTagKeyItem struct {
	BaseCfgItem `json:",inline"`
}

// InstanceTagValueItem TODO
type InstanceTagValueItem struct {
	BaseCfgItem `json:",inline"`
}

// BusinessTopoInstNamesItem TODO
type BusinessTopoInstNamesItem struct {
	BaseCfgItem `json:",inline"`
}
