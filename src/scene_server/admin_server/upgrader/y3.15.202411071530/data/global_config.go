/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
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

package data

import (
	"fmt"
	"reflect"

	tenanttmp "configcenter/pkg/types/tenant-template"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/scene_server/admin_server/upgrader/tools"
	"configcenter/src/storage/dal/mongo/local"
)

func addGlobalConfigData(kit *rest.Kit, db local.DB) error {

	dbConfig := make([]PlatformSettingConfig, 0)
	err := db.Table(common.BKTableNameGlobalConfig).Find(mapstr.MapStr{}).All(kit.Ctx, &dbConfig)
	if err != nil {
		blog.Errorf("find data for table %s failed, err: %v", common.BKTableNameGlobalConfig, err)
		return err
	}

	insertConfig := InitGlobalConfig
	switch len(dbConfig) {
	case 0:
		insertConfig.Time = *tools.NewTime()
	case 1:
		if !cmpGlobalConfig(&dbConfig[0], &insertConfig) {
			blog.Errorf("exist data is not equal to insert, exist: %v, insert: %v", dbConfig, insertConfig)
			return err
		}
		return nil
	default:
		blog.Errorf("invalid exist data for global config, count: %d", len(dbConfig))
		return fmt.Errorf("invalid exist data for global config, count: %d", len(dbConfig))
	}

	err = db.Table(common.BKTableNameGlobalConfig).Insert(kit.Ctx, insertConfig)
	if err != nil {
		blog.Errorf("upsert data for table %s failed, data: %v, err: %v", common.BKTableNameGlobalConfig, insertConfig,
			err)
		return err
	}

	err = tools.InsertTemplateData(kit, db, []interface{}{InitGlobalConfig}, &tools.InsertOptions{}, &tools.IDOptions{
		RemoveKeys: []string{"tenant_id"}}, tenanttmp.TemplateTypeGlobalConfig)
	if err != nil {
		blog.Errorf("insert template data failed, err: %v", err)
		return err
	}

	// add audit data
	auditField := &tools.AuditStruct{
		AuditTypeData: &tools.AuditResType{
			AuditType:    "platform_setting",
			ResourceType: "platform_setting",
		},
		AuditDataField: &tools.AuditDataField{},
	}
	globalConfigMap, err := tools.ConvStructToMap(InitGlobalConfig)
	if err != nil {
		blog.Errorf("convert struct to map failed, err: %v", err)
		return err
	}
	if err = tools.AddCreateAuditLog(kit, db, []map[string]interface{}{globalConfigMap}, auditField); err != nil {
		blog.Errorf("add audit log failed, err: %v", err)
		return err
	}

	return nil
}

// InitGlobalConfig init global config
var InitGlobalConfig = PlatformSettingConfig{
	Backend: AdminBackendCfg{
		MaxBizTopoLevel: 7,
	},
	ValidationRules: ValidationRulesCfg{
		Number: NumberItem{
			BaseCfgItem: BaseCfgItem{
				Value:       "XihcLXxcKyk/XGQrJA==",
				Description: "字段类型'数字'的验证规则",
				I18N: I18N{
					CN: "请输入整数数字",
					EN: "Please enter integer number",
				},
			},
		},
		Float: FloatItem{
			BaseCfgItem: BaseCfgItem{
				Value:       "XlsrLV0/KFswLTldKlsuXT9bMC05XSt8WzAtOV0rWy5dP1swLTldKikoW2VFXVsrLV0/WzAtOV0rKT8k",
				Description: "字段类型'浮点'的验证规则",
				I18N: I18N{
					CN: "请输入浮点型数字",
					EN: "Please enter float number",
				},
			},
		},
		Singlechar: SinglecharItem{
			BaseCfgItem: BaseCfgItem{
				Value:       "XFMq",
				Description: "字段类型'短字符'的验证规则",
				I18N: I18N{
					CN: "请输入256长度以内的字符串",
					EN: "Please enter the string within 256 length",
				},
			},
		},
		Longchar: LongcharItem{
			BaseCfgItem: BaseCfgItem{
				Value:       "XFMq",
				Description: "字段类型'长字符'的验证规则",
				I18N: I18N{
					CN: "请输入2000长度以内的字符串",
					EN: "Please enter the string within 2000 length",
				},
			},
		},
		AssociationId: AssociationIdItem{
			BaseCfgItem: BaseCfgItem{
				Value:       "XlthLXpBLVpdW1x3XSok",
				Description: "关联类型唯一标识验证规则",
				I18N: I18N{
					CN: "由英文字符开头，和下划线、数字或英文组合的字符",
					EN: "Start with lowercase or uppercase letter, followed by lowercase / uppercase / underscore / numbers characters",
				},
			},
		},
		ClassifyId: ClassifyIdItem{
			BaseCfgItem: BaseCfgItem{
				Value:       "XlthLXpBLVpdW1x3XSok",
				Description: "模型分组唯一标识验证规则",
				I18N: I18N{
					CN: "由英文字符开头，和下划线、数字或英文组合的字符",
					EN: "Start with lowercase or uppercase letter, followed by lowercase / uppercase / underscore / numbers characters",
				},
			},
		},
		ModelId: ModelIdItem{
			BaseCfgItem: BaseCfgItem{
				Value:       "XlthLXpBLVpdW1x3XSok",
				Description: "模型分组唯一标识验证规则",
				I18N: I18N{
					CN: "由英文字符开头，和下划线、数字或英文组合的字符",
					EN: "Start with lowercase or uppercase letter, followed by lowercase / uppercase / underscore / numbers characters",
				},
			},
		},
		EnumId: EnumIdItem{
			BaseCfgItem: BaseCfgItem{
				Value:       "XlthLXpBLVowLTlfLV0qJA==",
				Description: "字段类型'枚举'ID的验证规则",
				I18N: I18N{
					CN: "由大小写英文字母，数字，_ 或 - 组成的字符",
					EN: "Composed of uppercase / lowercase / numbers / - or _ characters",
				},
			},
		},
		EnumName: EnumNameItem{
			BaseCfgItem: BaseCfgItem{
				Value:       "XihbYS16QS1aMC05X118W1x1NGUwMC1cdTlmYTVdfFsoKSst44CK44CLLO+8jO+8mzvigJzigJ3igJjigJnjgIJcLlwiXCcgXC86XSkqJA==",
				Description: "字段类型'枚举'值的验证规则",
				I18N: I18N{
					CN: "请输入枚举值",
					EN: "Please enter the enum value",
				},
			},
		},
		FieldId: FieldIdItem{
			BaseCfgItem: BaseCfgItem{
				Value:       "XlthLXpBLVpdW1x3XSok",
				Description: "模型字段唯一标识的验证规则",
				I18N: I18N{
					CN: "由英文字符开头，和下划线、数字或英文组合的字符",
					EN: "Start with lowercase or uppercase letter, followed by lowercase / uppercase / underscore / numbers characters",
				},
			},
		},
		NamedCharacter: NamedCharacterItem{
			BaseCfgItem: BaseCfgItem{
				Value:       "XlthLXpBLVowLTlcdTRlMDAtXHU5ZmE1X1wtOlwoXCldKyQ=",
				Description: "服务分类名称的验证规则",
				I18N: I18N{
					CN: "请输入中英文或特殊字符 :_- 组成的名称",
					EN: "Special symbols only support(:_-)",
				},
			},
		},
		InstanceTagKey: InstanceTagKeyItem{
			BaseCfgItem: BaseCfgItem{
				Value:       "XlthLXpBLVpdKFthLXowLTlBLVpcLV8uXSpbYS16MC05QS1aXSk/JA==",
				Description: "服务实例标签键的验证规则",
				I18N: I18N{
					CN: "请输入以英文开头的英文+数字组合",
					EN: "Please enter letter / number starts with letter",
				},
			},
		},
		InstanceTagValue: InstanceTagValueItem{
			BaseCfgItem: BaseCfgItem{
				Value:       "XlthLXowLTlBLVpdKFthLXowLTlBLVpcLV8uXSpbYS16MC05QS1aXSk/JA==",
				Description: "服务实例标签值的验证规则",
				I18N: I18N{
					CN: "请输入英文 / 数字",
					EN: "Please enter letter / number",
				},
			},
		},
		BusinessTopoInstNames: BusinessTopoInstNamesItem{
			BaseCfgItem: BaseCfgItem{
				Value:       "XlteXCNcLyxcPlw8XHxdKyQ=",
				Description: "集群/模块/实例名称的验证规则",
				I18N: I18N{
					CN: "请输入除 #/,><| 以外的字符",
					EN: "Please enter characters other than #/,><|",
				},
			},
		},
	},
	BuiltInSetName: "空闲机池",
	BuiltInModuleConfig: GlobalModule{
		IdleName:    "空闲机",
		FaultName:   "故障机",
		RecycleName: "待回收",
		UserModules: nil,
	},
}

func cmpGlobalConfig(exist, insert *PlatformSettingConfig) bool {

	if exist.Backend.MaxBizTopoLevel != insert.Backend.MaxBizTopoLevel {
		return false
	}

	if exist.BuiltInModuleConfig.FaultName != insert.BuiltInModuleConfig.FaultName {
		return false
	}

	if exist.BuiltInModuleConfig.RecycleName != insert.BuiltInModuleConfig.RecycleName {
		exist.BuiltInModuleConfig.RecycleName = insert.BuiltInModuleConfig.RecycleName
	}

	if exist.BuiltInModuleConfig.IdleName != insert.BuiltInModuleConfig.IdleName {
		exist.BuiltInModuleConfig.IdleName = insert.BuiltInModuleConfig.IdleName
	}

	if exist.BuiltInSetName != insert.BuiltInSetName {
		exist.BuiltInSetName = insert.BuiltInSetName
	}

	existRuleType := reflect.TypeOf(exist.ValidationRules)
	existRuleVal := reflect.ValueOf(&exist.ValidationRules).Elem()
	insertRuleVal := reflect.ValueOf(&insert.ValidationRules).Elem()

	for i := 0; i < existRuleType.NumField(); i++ {
		fieldName := existRuleType.Field(i).Name
		existRuleValStr := existRuleVal.FieldByName(fieldName).FieldByName("BaseCfgItem").FieldByName("Value").String()
		insertRuleValStr := insertRuleVal.FieldByName(fieldName).FieldByName("BaseCfgItem").FieldByName("Value").String()
		if existRuleValStr != insertRuleValStr {
			return false
		}
	}

	return true
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

type PlatformSettingConfig struct {
	TenantID            string             `json:"-" bson:"tenant_id"`
	Backend             AdminBackendCfg    `json:"backend" bson:"backend"`
	ValidationRules     ValidationRulesCfg `json:"validation_rules" bson:"validation_rules"`
	BuiltInSetName      ObjectString       `json:"set" bson:"set"`
	BuiltInModuleConfig GlobalModule       `json:"idle_pool" bson:"idle_pool"`
	Time                tools.Time         `bson:",inline"`
}

// AdminBackendCfg used to admin backend config
type AdminBackendCfg struct {
	MaxBizTopoLevel int64 `json:"max_biz_topo_level" bson:"max_biz_topo_level"`
}

// ObjectString used to admin set config
type ObjectString string

// GlobalModule used to admin global module config
type GlobalModule struct {
	IdleName    string           `json:"idle" bson:"idle"`
	FaultName   string           `json:"fault" bson:"fault"`
	RecycleName string           `json:"recycle" bson:"recycle"`
	UserModules []UserModuleList `json:"user_modules" bson:"user_modules"`
}

// UserModuleList custom section.
type UserModuleList struct {
	Key   string `json:"module_key" bson:"module_key"`
	Value string `json:"module_name" bson:"module_name"`
}

// BaseCfgItem the common base config item
type BaseCfgItem struct {
	Value       string `json:"value" bson:"value"`
	Description string `json:"description" bson:"description"`
	I18N        I18N   `json:"i18n" bson:"i18n"`
}

// I18N use English or Chinese
type I18N struct {
	CN string `json:"cn" bson:"cn"`
	EN string `json:"en" bson:"en"`
}

// NumberItem number item
type NumberItem struct {
	BaseCfgItem `json:",inline" bson:",inline"`
}

// FloatItem float item
type FloatItem struct {
	BaseCfgItem `json:",inline" bson:",inline"`
}

// SinglecharItem single char item
type SinglecharItem struct {
	BaseCfgItem `json:",inline" bson:",inline"`
}

// LongcharItem long char item
type LongcharItem struct {
	BaseCfgItem `json:",inline" bson:",inline"`
}

// AssociationIdItem association id item
type AssociationIdItem struct {
	BaseCfgItem `json:",inline" bson:",inline"`
}

// ClassifyIdItem classify id item
type ClassifyIdItem struct {
	BaseCfgItem `json:",inline" bson:",inline"`
}

// ModelIdItem model id item
type ModelIdItem struct {
	BaseCfgItem `json:",inline" bson:",inline"`
}

// EnumIdItem enum id item
type EnumIdItem struct {
	BaseCfgItem `json:",inline" bson:",inline"`
}

// EnumNameItem enum name item
type EnumNameItem struct {
	BaseCfgItem `json:",inline" bson:",inline"`
}

// FieldIdItem field id item
type FieldIdItem struct {
	BaseCfgItem `json:",inline" bson:",inline"`
}

// NamedCharacterItem named character item
type NamedCharacterItem struct {
	BaseCfgItem `json:",inline" bson:",inline"`
}

// InstanceTagKeyItem instance tag key item
type InstanceTagKeyItem struct {
	BaseCfgItem `json:",inline" bson:",inline"`
}

// InstanceTagValueItem instance tag value item
type InstanceTagValueItem struct {
	BaseCfgItem `json:",inline" bson:",inline"`
}

// BusinessTopoInstNamesItem business topo inst names item
type BusinessTopoInstNamesItem struct {
	BaseCfgItem `json:",inline" bson:",inline"`
}
