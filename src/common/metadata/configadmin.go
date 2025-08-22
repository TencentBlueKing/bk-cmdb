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

package metadata

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	idgen "configcenter/pkg/id-gen"
	"configcenter/src/common"
)

// ConfigAdminResult the result with config admin
type ConfigAdminResult struct {
	BaseResp `json:",inline"`
	Data     ConfigAdmin `json:"data"`
}

// GlobalConfOptions the options of global config
type GlobalConfOptions struct {
	Fields []string `json:"fields" bson:"fields"`
}

// GlobalSettingResult the result of global setting.
type GlobalSettingResult struct {
	BaseResp `json:",inline"`
	Data     GlobalSettingConfig `json:"data"`
}

// ConfigAdminParmams used to admin the cmdb config
type ConfigAdminParmams struct {
	Config ConfigAdmin `json:"config"`
}

// ConfigAdmin used to admin the cmdb config
type ConfigAdmin struct {
	Backend         BackendCfg         `json:"backend" bson:"backend"`
	ValidationRules ValidationRulesCfg `json:"validationRules" bson:"validationRules"`
}

// UserModuleList custom section.
type UserModuleList struct {
	Key   string `json:"module_key" bson:"module_key"`
	Value string `json:"module_name" bson:"module_name"`
}

// GlobalModule Conifg, idleName, FaultName and RecycleName cannot be deleted.
type GlobalModule struct {
	IdleName    string           `json:"idle" bson:"idle"`
	FaultName   string           `json:"fault" bson:"fault"`
	RecycleName string           `json:"recycle" bson:"recycle"`
	UserModules []UserModuleList `json:"user_modules" bson:"user_modules"`
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
	MaxBizTopoLevel int64 `json:"max_biz_topo_level" bson:"max_biz_topo_level"`
}

// Validate validate the fields of BackendCfg.
func (b AdminBackendCfg) Validate() error {
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

// IDGeneratorConf is id generator config
type IDGeneratorConf struct {
	Enabled bool                       `json:"enabled" bson:"enabled"`
	Step    int                        `json:"step" bson:"step"`
	InitID  map[idgen.IDGenType]uint64 `json:"init_id,omitempty" bson:"init_id,omitempty"`
	// CurrentID is the current id of each resource, this is only used for ui display
	CurrentID map[idgen.IDGenType]uint64 `json:"current_id,omitempty" bson:"current_id,omitempty"`
}

// platform setting config
const (
	// IDGeneratorConfig is id generator config for platform config
	IDGeneratorConfig = "id_generator"
)

// PlatformConfig platform config
type PlatformConfig struct {
	IDGenerator IDGeneratorConf `bson:"id_generator" json:"id_generator"`
}

// Validate id generator config
func (c IDGeneratorConf) Validate() error {
	if c.Step <= 0 {
		return fmt.Errorf("step is invalid")
	}

	if len(c.InitID) == 0 {
		return nil
	}

	for res, id := range c.InitID {
		if id <= 0 {
			return fmt.Errorf("%s init id %d is invalid", res, id)
		}
	}

	return nil
}

// GlobalSettingConfig  used to admin the global config. GlobalSettingConfig 每个成员对象都必须有"Validate"校验
// 函数，如果没有会panic.
type GlobalSettingConfig struct {
	TenantID            ObjectString       `json:"-" bson:"tenant_id"`
	Backend             AdminBackendCfg    `json:"backend" bson:"backend"`
	ValidationRules     ValidationRulesCfg `json:"validation_rules" bson:"validation_rules"`
	BuiltInSetName      ObjectString       `json:"set" bson:"set"`
	BuiltInModuleConfig GlobalModule       `json:"idle_pool" bson:"idle_pool"`
	CreateTime          Time               `json:"create_time" bson:"create_time"`
	LastTime            Time               `json:"last_time" bson:"last_time"`
}

// global configs
const (
	// BuiltInSetNameConfig is built-in set name config for platform config
	BuiltInSetNameConfig = "set"
	// BuiltInModuleConfig is built-in module config for platform config
	BuiltInModuleConfig = "idle_pool"
	// BackendConfig is backend config for platform config
	BackendConfig = "backend"
)

// Validate validate the fields of GlobalSettingReqOption is illegal .
func (c *GlobalSettingConfig) Validate() error {
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
				return fmt.Errorf("%s %s failed, error:%s", vrt.Field(i).Name, fn,
					errVal[0].Interface().(error).Error())
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
func (c *GlobalSettingConfig) EncodeWithBase64() error {
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
	MaxBizTopoLevel int64 `json:"max_biz_topo_level" bson:"max_biz_topo_level"`
}

// Validate validate the fields of BackendCfg
func (b BackendCfg) Validate() error {
	if b.MaxBizTopoLevel < 3 || b.MaxBizTopoLevel > 10 {
		return fmt.Errorf("maxBizTopoLevel value must in range [3-10]")
	}
	return nil
}

// BaseCfgItem the common base config item
type BaseCfgItem struct {
	Value       string `json:"value" bson:"value"`
	Description string `json:"description" bson:"description"`
	I18N        I18N   `json:"i18n" bson:"i18n"`
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
	CN string `json:"cn" bson:"cn"`
	EN string `json:"en" bson:"en"`
}

// ValidationRulesCfg used to admin valiedation rules Config
type ValidationRulesCfg struct {
	Number                NumberItem                `json:"number" bson:"number"`
	Float                 FloatItem                 `json:"float" bson:"float"`
	Singlechar            SinglecharItem            `json:"singlechar" bson:"singlechar"`
	Longchar              LongcharItem              `json:"longchar" bson:"longchar"`
	AssociationId         AssociationIdItem         `json:"associationId" bson:"associationId"`
	ClassifyId            ClassifyIdItem            `json:"classifyId" bson:"classifyId"`
	ModelId               ModelIdItem               `json:"modelId" bson:"modelId"`
	EnumId                EnumIdItem                `json:"enumId" bson:"enumId"`
	EnumName              EnumNameItem              `json:"enumName" bson:"enumName"`
	FieldId               FieldIdItem               `json:"fieldId" bson:"fieldId"`
	NamedCharacter        NamedCharacterItem        `json:"namedCharacter" bson:"namedCharacter"`
	InstanceTagKey        InstanceTagKeyItem        `json:"instanceTagKey" bson:"instanceTagKey"`
	InstanceTagValue      InstanceTagValueItem      `json:"instanceTagValue" bson:"instanceTagValue"`
	BusinessTopoInstNames BusinessTopoInstNamesItem `json:"businessTopoInstNames" bson:"businessTopoInstNames"`
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
				return fmt.Errorf("%s %s failed, error:%s", vrt.Field(i).Name, fn,
					errVal[0].Interface().(error).Error())
			}
		}
	}
	return nil
}

// NumberItem TODO
type NumberItem struct {
	BaseCfgItem `json:",inline" bson:",inline"`
}

// FloatItem TODO
type FloatItem struct {
	BaseCfgItem `json:",inline" bson:",inline"`
}

// SinglecharItem TODO
type SinglecharItem struct {
	BaseCfgItem `json:",inline" bson:",inline"`
}

// LongcharItem TODO
type LongcharItem struct {
	BaseCfgItem `json:",inline" bson:",inline"`
}

// AssociationIdItem TODO
type AssociationIdItem struct {
	BaseCfgItem `json:",inline" bson:",inline"`
}

// ClassifyIdItem TODO
type ClassifyIdItem struct {
	BaseCfgItem `json:",inline" bson:",inline"`
}

// ModelIdItem TODO
type ModelIdItem struct {
	BaseCfgItem `json:",inline" bson:",inline"`
}

// EnumIdItem TODO
type EnumIdItem struct {
	BaseCfgItem `json:",inline" bson:",inline"`
}

// EnumNameItem TODO
type EnumNameItem struct {
	BaseCfgItem `json:",inline" bson:",inline"`
}

// FieldIdItem TODO
type FieldIdItem struct {
	BaseCfgItem `json:",inline" bson:",inline"`
}

// NamedCharacterItem TODO
type NamedCharacterItem struct {
	BaseCfgItem `json:",inline" bson:",inline"`
}

// InstanceTagKeyItem TODO
type InstanceTagKeyItem struct {
	BaseCfgItem `json:",inline" bson:",inline"`
}

// InstanceTagValueItem TODO
type InstanceTagValueItem struct {
	BaseCfgItem `json:",inline" bson:",inline"`
}

// BusinessTopoInstNamesItem TODO
type BusinessTopoInstNamesItem struct {
	BaseCfgItem `json:",inline" bson:",inline"`
}

// OldAdminBackendCfg old admin backend config
type OldAdminBackendCfg struct {
	MaxBizTopoLevel int64  `json:"max_biz_topo_level"`
	SnapshotBizName string `json:"snapshot_biz_name"`
	SnapshotBizID   int64  `json:"snapshot_biz_id"`
}

// Validate validate the fields of BackendCfg.
func (o OldAdminBackendCfg) Validate() error {
	if o.MaxBizTopoLevel < minBizTopoLevel || o.MaxBizTopoLevel > maxBizTopoLevel {
		return fmt.Errorf("max biz topo level value must in range [%d-%d]", minBizTopoLevel, maxBizTopoLevel)
	}
	return nil
}

// OldPlatformSettingConfig old platform setting config
type OldPlatformSettingConfig struct {
	Backend             OldAdminBackendCfg `json:"backend"`
	ValidationRules     ValidationRulesCfg `json:"validation_rules"`
	BuiltInSetName      ObjectString       `json:"set"`
	BuiltInModuleConfig GlobalModule       `json:"idle_pool"`
	IDGenerator         IDGeneratorConf    `json:"id_generator"`
}

// Validate validate the fields of OldPlatformSettingConfig is illegal .
func (c *OldPlatformSettingConfig) Validate() error {
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

// EncodeWithBase64 encode the value of ValidationRules to base64.
func (c *OldPlatformSettingConfig) EncodeWithBase64() error {
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
