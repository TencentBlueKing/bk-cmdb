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

package history

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	idgen "configcenter/pkg/id-gen"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal"
)

// configChangeItem 配置变更项
type configChangeItem struct {
	// 变更目录，和upgrader下的变更程序目录保持一致，便于查找升级程序
	dir string
	// 变更说明
	description string
	// 变更后的配置，json格式
	config string
}

var initConfig = `
{
    "backend": {
        "snapshotBizName": "蓝鲸",
        "maxBizTopoLevel": 7
    },
    "validationRules": {
        "number": {
            "value": "^(\\-|\\+)?\\d+$",
            "description": "字段类型“数字”的验证规则",
            "i18n": {
                "cn": "请输入正确的数字",
                "en": "Please enter the correct number"
            }
        },
        "float": {
            "value": "^[+-]?([0-9]*[.]?[0-9]+|[0-9]+[.]?[0-9]*)([eE][+-]?[0-9]+)?$",
            "description": "字段类型“浮点”的验证规则",
            "i18n": {
                "cn": "请输入正确的浮点数",
                "en": "Please enter the correct float data"
            }
        },
        "singlechar": {
            "value": "\\S*",
            "description": "字段类型“短字符”的验证规则",
            "i18n": {
                "cn": "请输入正确的短字符内容",
                "en": "Please enter the correct content"
            }
        },
        "longchar": {
            "value": "\\S*",
            "description": "字段类型“长字符”的验证规则",
            "i18n": {
                "cn": "请输入正确的长字符内容",
                "en": "Please enter the correct content"
            }
        },
        "associationId": {
            "value": "^[a-zA-Z][\\w]*$",
            "description": "关联类型唯一标识验证规则",
            "i18n": {
                "cn": "格式不正确，请填写英文开头，下划线，数字，英文的组合",
                "en": "The format is incorrect, can only contain underscores, numbers, letter and start with a letter"
            }
        },
        "classifyId": {
            "value": "^[a-zA-Z][\\w]*$",
            "description": "模型分组唯一标识验证规则",
            "i18n": {
                "cn": "请输入正确的内容",
                "en": "Please enter the correct content"
            }
        },
        "modelId": {
            "value": "^[a-zA-Z][\\w]*$",
            "description": "模型唯一标识验证规则",
            "i18n": {
                "cn": "格式不正确，请填写英文开头，下划线，数字，英文的组合",
                "en": "The format is incorrect, can only contain underscores, numbers, letter and start with a letter"
            }
        },
        "enumId": {
            "value": "^[a-zA-Z0-9_-]*$",
            "description": "字段类型“枚举”ID的验证规则",
            "i18n": {
                "cn": "请输入正确的内容",
                "en": "Please enter the correct content"
            }
        },
        "enumName": {
            "value": "^([a-zA-Z0-9_]|[\\u4e00-\\u9fa5]|[()+-《》,，；;“”‘’。\\.\\\"\\' \\/:])*$",
            "description": "字段类型“枚举”值的验证规则",
            "i18n": {
                "cn": "请输入正确的内容",
                "en": "Please enter the correct content"
            }
        },
        "fieldId": {
            "value": "^[a-zA-Z][\\w]*$",
            "description": "模型字段唯一标识的验证规则",
            "i18n": {
                "cn": "请输入正确的内容",
                "en": "Please enter the correct content"
            }
        },
        "namedCharacter": {
            "value": "^[a-zA-Z0-9\\u4e00-\\u9fa5_\\-:\\(\\)]+$",
            "description": "服务分类名称的验证规则",
            "i18n": {
                "cn": "格式不正确，特殊符号仅支持(:_-)",
                "en": "Special symbols only support(:_-)"
            }
        },
        "instanceTagKey": {
            "value": "^[a-zA-Z]([a-z0-9A-Z\\-_.]*[a-z0-9A-Z])?$",
            "description": "服务实例标签键的验证规则",
            "i18n": {
                "cn": "请输入英文 / 数字, 以英文开头",
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
                "cn": "格式不正确，不能包含特殊字符 #/,><|",
                "en": "The format is incorrect and cannot contain special characters #/,><|"
            }
        }
    }
}
`

// configChangeHistory 配置变更历史
// 每次配置有变更时，需要在configChangeHistory里增加一项configChangeItem，表示当前要变更的配置
var configChangeHistory = []configChangeItem{
	{
		dir:         "y3.8.202006092135",
		description: "第一次初始化后的配置",
		config:      initConfig,
	},
}

// AddConfigAdminChangeItem add config admin change item to the change history for later upgrading use
func AddConfigAdminChangeItem(dir string, description string, config string) {
	configChangeHistory = append(configChangeHistory, configChangeItem{
		dir:         dir,
		description: description,
		config:      config,
	})
}

// UpgradeConfigAdmin 升级配置管理
// 每次升级变更配置，需要在configChangeHistory最后追加一项要变更的配置
// 将configChangeHistory里的最后一项作为当前配置项curCfg，倒数第二项作为上一次配置项preCfg
// 需要将preCfg和db存在的配置dbCfg进行对比，对于不一致的（说明有用户调过配置管理接口做过更改）,curCfg里对应的配置不做覆盖，仍为db里的数据
func UpgradeConfigAdmin(ctx context.Context, db dal.RDB, dir string) error {
	preCfg, curCfg, dbCfgStr, err := getConfigs(ctx, db, dir)
	if err != nil {
		blog.Errorf("upgradeConfigAdmin failed, getConfigs err: %v", err)
		return err
	}
	dbCfg := new(ConfigAdmin)

	if dbCfgStr != "" {
		// dbNewCfg 用来保存新的 PlatformSettingConfig 结构数据
		dbNewCfg := new(PlatformSettingConfig)
		if err := json.Unmarshal([]byte(dbCfgStr), dbNewCfg); err != nil {
			blog.Errorf("get dbConfig failed, unmarshal err: %v, config: %v", err, dbCfgStr)
			return err
		}

		// 此时说明数据库中是最新的配置结构不需要后面的动作升级
		if dbNewCfg.Backend.MaxBizTopoLevel != 0 {
			return nil
		}
		// dbCfg 用来保存老的ConfigAdmin结构数据
		if err := json.Unmarshal([]byte(dbCfgStr), dbCfg); err != nil {
			blog.Errorf("get dbConfig failed, unmarshal err: %v, config: %v", err, dbCfgStr)
			return err
		}
		// 说明db中的config有问题，直接报错
		if dbCfg.Backend.MaxBizTopoLevel == 0 {
			return fmt.Errorf("config is error")
		}
	}

	if err := curCfg.EncodeWithBase64(); err != nil {
		blog.Errorf("UpgradeConfigAdmin failed, EncodeWithBase64 err: %v, curCfg:%#v", err, *curCfg)
		return err
	}
	if err := curCfg.Validate(); err != nil {
		blog.Errorf("UpgradeConfigAdmin failed, Validate err: %v, curCfg:%#v", err, *curCfg)
		return err
	}
	// 如果是首次进行初始化，直接保存当前配置
	if preCfg == nil {
		err = updateConfig(ctx, db, curCfg)
		if err != nil {
			blog.Errorf("UpgradeConfigAdmin failed, update config err: %v, config:%#v", err, *curCfg)
			return err
		}
		return nil
	}

	// 对比上一次配置和db配置的差异，获取最终配置
	if err := preCfg.EncodeWithBase64(); err != nil {
		blog.Errorf("UpgradeConfigAdmin failed, EncodeWithBase64 err: %v, preCfg:%#v", err, *preCfg)
		return err
	}

	config := getFinalConfig(preCfg, curCfg, dbCfg)

	err = updateConfig(ctx, db, config)
	if err != nil {
		blog.Errorf("UpgradeConfigAdmin failed, update config err: %v, config:%#v", err, *config)
		return err
	}

	return nil
}

// UpgradePlatConfigAdmin 对平台配置进行升级:
// 1、获取三份配置文件分别是上一份初始化配置 preCfg。数据库中保存的配置文件dbCfg。当前要升级的配置文件curCfg。
// 2、原则是如果对于原有字段，如果preCfg与dbCfg字段相同那么直接升级到curCfg即可。 如果preCfg中的字段和dbCfg中的字段不同，说明用户改过
// 本字段,本次升级保持用户改后的配置不变.
func UpgradePlatConfigAdmin(ctx context.Context, db dal.RDB, dir string) error {

	curCfg, preCfg, dbCfgStr, err := getAllConfigs(ctx, db, dir)
	if err != nil {
		blog.Errorf("upgrade platform config failed, err: %v", err)
		return err
	}

	// dbNewCfg 用来保存新的 GlobalSettingConfig 结构数据，对于此场景dbCfgStr 必然不为空如果为空直接报错就好
	dbNewCfg := new(PlatformSettingConfig)

	if err := json.Unmarshal([]byte(dbCfgStr), dbNewCfg); err != nil {
		blog.Errorf("get dbConfig failed, unmarshal err: %v, config: %v", err, dbCfgStr)
		return err
	}

	// 此时说明数据库中是最新的配置结构不需要后面的动作升级
	if dbNewCfg.Backend.MaxBizTopoLevel != 0 {
		return nil
	}

	// dbCfg 用来保存老的ConfigAdmin结构数据
	dbCfg := new(ConfigAdmin)
	if err := json.Unmarshal([]byte(dbCfgStr), dbCfg); err != nil {
		blog.Errorf("get dbConfig failed, unmarshal err: %v, config: %v", err, dbCfgStr)
		return err
	}

	// 此场景说明数据库中的配置文件有问题，直接报错
	if dbCfg.Backend.MaxBizTopoLevel == 0 {
		return fmt.Errorf("config is error")
	}

	// 从这里开始场景是老的配置文件升级到新的配置文件场景
	if err := curCfg.EncodeWithBase64(); err != nil {
		blog.Errorf("upgrade platform config encode base64 failed, curCfg: %#v, err: %v", *curCfg, err)
		return err
	}
	if err := curCfg.Validate(); err != nil {
		blog.Errorf("upgrade platform config validate failed, curCfg: %#v, err: %v", *curCfg, err)
		return err
	}

	// 如果是首次进行初始化，直接保存当前配置
	if preCfg == nil {
		err = updatePlatformConfig(ctx, db, curCfg)
		if err != nil {
			blog.Errorf("upgrade platform config failed, config %v, err: %v", *curCfg, err)
			return err
		}
		return nil
	}

	// 对比上一次配置和db配置的差异，获取最终配置
	if err := preCfg.EncodeWithBase64(); err != nil {
		blog.Errorf("get preCfg config encode base64 failed, preCfg: %#v, err: %v", *preCfg, err)
		return err
	}
	// 获取最终的配置文件
	config := getFinalPlatformConfig(preCfg, dbCfg, curCfg)

	err = updatePlatformConfig(ctx, db, config)
	if err != nil {
		blog.Errorf("update platform config failed, config: %#v, err: %v", *config, err)
		return err
	}

	return nil
}

// getConfigs 获取preCfg, curCfg
func getConfigs(ctx context.Context, db dal.RDB, dir string) (preCfg, curCfg *ConfigAdmin, dbCfg string,
	err error) {
	var pre string
	for index, config := range configChangeHistory {
		if config.dir == dir {
			cur := config.config
			curCfg = new(ConfigAdmin)
			if err := json.Unmarshal([]byte(cur), curCfg); err != nil {
				blog.Errorf("get all config failed, Unmarshal err: %v, config: %v", err, cur)
				return nil, nil, "", err
			}

			// 第一次初始化,preCfg为nil
			if index == 0 {
				return nil, curCfg, "", nil
			}

			pre = configChangeHistory[index-1].config
			break
		}
	}
	cond := map[string]interface{}{
		"_id": common.ConfigAdminID,
	}
	ret := make(map[string]interface{})
	err = db.Table(common.BKTableNameSystem).Find(cond).Fields(common.ConfigAdminValueField).One(ctx, &ret)
	if nil != err {
		blog.Errorf("get db config failed, err: %v", err)
		return nil, nil, "", err
	}

	if ret[common.ConfigAdminValueField] == nil {
		blog.Errorf(" db config type is error")
		return nil, nil, "", err
	}

	if _, ok := ret[common.ConfigAdminValueField].(string); !ok {
		blog.Errorf("get db config type is error")
		return nil, nil, "", err
	}
	preCfg = new(ConfigAdmin)
	if err = json.Unmarshal([]byte(pre), preCfg); err != nil {
		blog.Errorf("get all config failed, Unmarshal err: %v, config: %v", err, pre)
		return nil, nil, "", err
	}

	return preCfg, curCfg, ret[common.ConfigAdminValueField].(string), nil
}

// getAllConfigs get preCfg, curCfg.
func getAllConfigs(ctx context.Context, db dal.RDB, dir string) (curCfg *OldPlatformSettingConfig,
	preCfg *ConfigAdmin, dbCfg string, err error) {
	var pre string

	for index, config := range configChangeHistory {
		if config.dir == dir {
			cur := config.config
			curCfg = new(OldPlatformSettingConfig)
			if err := json.Unmarshal([]byte(cur), curCfg); err != nil {
				blog.Errorf("get all config failed, unmarshal err: %v, config: %v", err, cur)
				return nil, nil, "", err
			}

			// 第一次初始化,preCfg为nil
			if index == 0 {
				return curCfg, nil, "", nil
			}

			pre = configChangeHistory[index-1].config
			break
		}
	}

	cond := map[string]interface{}{
		"_id": common.ConfigAdminID,
	}
	ret := make(map[string]interface{})
	err = db.Table(common.BKTableNameSystem).Find(cond).Fields(common.ConfigAdminValueField).One(ctx, &ret)
	if nil != err {
		blog.Errorf("get all config failed, db find err: %v", err)
		return nil, nil, "", err
	}

	if ret[common.ConfigAdminValueField] == nil {
		blog.Errorf("get config failed, db config type is error")
		return nil, nil, "", nil
	}

	if _, ok := ret[common.ConfigAdminValueField].(string); !ok {
		blog.Errorf("get config failed, db config type is error")
		return nil, nil, "", nil
	}

	preCfg = new(ConfigAdmin)
	if err = json.Unmarshal([]byte(pre), preCfg); err != nil {
		blog.Errorf("get all config failed, unmarshal err: %v, config: %v", err, pre)
		return nil, nil, "", err
	}

	return curCfg, preCfg, ret[common.ConfigAdminValueField].(string), nil
}

// getFinalConfig 获取最终需要保存的配置
// 将preCfg和db存在的配置dbCfg进行对比，对于不一致的（说明有用户调过配置管理接口做过更改）,curCfg里对应的配置不做覆盖，仍为db里的数据
func getFinalConfig(preCfg, curCfg, dbCfg *ConfigAdmin) *ConfigAdmin {
	if preCfg.Backend.SnapshotBizName != dbCfg.Backend.SnapshotBizName {
		curCfg.Backend.SnapshotBizName = dbCfg.Backend.SnapshotBizName
	}
	if preCfg.Backend.MaxBizTopoLevel != dbCfg.Backend.MaxBizTopoLevel {
		curCfg.Backend.MaxBizTopoLevel = dbCfg.Backend.MaxBizTopoLevel
	}

	preRuleType := reflect.TypeOf(preCfg.ValidationRules)
	preRuleVal := reflect.ValueOf(&preCfg.ValidationRules).Elem()
	dbRuleVal := reflect.ValueOf(&dbCfg.ValidationRules).Elem()
	curRuleVal := reflect.ValueOf(&curCfg.ValidationRules).Elem()

	for i := 0; i < preRuleType.NumField(); i++ {
		fieldName := preRuleType.Field(i).Name
		preRuleValStr := preRuleVal.FieldByName(fieldName).FieldByName("BaseCfgItem").FieldByName("Value").String()
		dbRuleValStr := dbRuleVal.FieldByName(fieldName).FieldByName("BaseCfgItem").FieldByName("Value").String()
		if preRuleValStr != dbRuleValStr {
			curRuleVal.FieldByName(fieldName).FieldByName("BaseCfgItem").FieldByName("Value").SetString(dbRuleValStr)
		}
	}

	return curCfg
}

// getFinalPlatformConfig 获取最终需要保存的配置
// 1、将preCfg和db存在的配置dbCfg进行对比，对于不一致的（说明有用户调过配置管理接口做过更改）,curCfg里对应的配置不做覆盖，仍为db里的数据
// 2、如果preCfg和dbCfg如果一样的话，那么如果本次curCfg不一样，则需要升级覆盖.
func getFinalPlatformConfig(preCfg, dbCfg *ConfigAdmin,
	curCfg *OldPlatformSettingConfig) *OldPlatformSettingConfig {

	if preCfg.Backend.SnapshotBizName != dbCfg.Backend.SnapshotBizName {
		curCfg.Backend.SnapshotBizName = dbCfg.Backend.SnapshotBizName
	}

	if preCfg.Backend.MaxBizTopoLevel != dbCfg.Backend.MaxBizTopoLevel {
		curCfg.Backend.MaxBizTopoLevel = dbCfg.Backend.MaxBizTopoLevel
	}

	preRuleType := reflect.TypeOf(preCfg.ValidationRules)
	preRule := reflect.ValueOf(&preCfg.ValidationRules).Elem()
	dbRule := reflect.ValueOf(&dbCfg.ValidationRules).Elem()
	curRule := reflect.ValueOf(&curCfg.ValidationRules).Elem()

	for i := 0; i < preRuleType.NumField(); i++ {
		fieldName := preRuleType.Field(i).Name
		preVal := preRule.FieldByName(fieldName).FieldByName("BaseCfgItem").FieldByName("Value").String()
		dbVal := dbRule.FieldByName(fieldName).FieldByName("BaseCfgItem").FieldByName("Value").String()
		if preVal != dbVal {
			curRule.FieldByName(fieldName).FieldByName("BaseCfgItem").FieldByName("Value").SetString(dbVal)
		}
	}

	return curCfg
}

// updateConfig 将配置更新到db里
func updateConfig(ctx context.Context, db dal.RDB, config *ConfigAdmin) error {
	bytes, err := json.Marshal(config)
	if err != nil {
		blog.Errorf("update config failed, Marshal err: %v, config:%+v", err, config)
		return err
	}

	cond := map[string]interface{}{
		"_id": common.ConfigAdminID,
	}
	data := map[string]interface{}{
		common.ConfigAdminValueField: string(bytes),
		common.LastTimeField:         time.Now(),
	}

	err = db.Table(common.BKTableNameSystem).Update(ctx, cond, data)
	if err != nil {
		blog.Errorf("update config failed, update err: %+v", err)
		return err
	}

	return nil
}

// updatePlatformConfig update configuration to database.
func updatePlatformConfig(ctx context.Context, db dal.RDB, config *OldPlatformSettingConfig) error {
	bytes, err := json.Marshal(config)
	if err != nil {
		return err
	}

	cond := map[string]interface{}{
		"_id": common.ConfigAdminID,
	}
	data := map[string]interface{}{
		common.ConfigAdminValueField: string(bytes),
		common.LastTimeField:         time.Now(),
	}

	err = db.Table(common.BKTableNameSystem).Update(ctx, cond, data)
	if err != nil {
		return err
	}

	return nil
}

type PlatformSettingConfig struct {
	Backend             AdminBackendCfg    `json:"backend"`
	ValidationRules     ValidationRulesCfg `json:"validation_rules"`
	BuiltInSetName      ObjectString       `json:"set"`
	BuiltInModuleConfig GlobalModule       `json:"idle_pool"`
	IDGenerator         IDGeneratorConf    `json:"id_generator"`
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

// AdminBackendCfg TODO
type AdminBackendCfg struct {
	MaxBizTopoLevel int64 `json:"max_biz_topo_level"`
}

// ValidationRulesCfg used to admin valiedation rules Config
type ValidationRulesCfg struct {
	Number                metadata.NumberItem                `json:"number"`
	Float                 metadata.FloatItem                 `json:"float"`
	Singlechar            metadata.SinglecharItem            `json:"singlechar"`
	Longchar              metadata.LongcharItem              `json:"longchar"`
	AssociationId         metadata.AssociationIdItem         `json:"associationId"`
	ClassifyId            metadata.ClassifyIdItem            `json:"classifyId"`
	ModelId               metadata.ModelIdItem               `json:"modelId"`
	EnumId                metadata.EnumIdItem                `json:"enumId"`
	EnumName              metadata.EnumNameItem              `json:"enumName"`
	FieldId               metadata.FieldIdItem               `json:"fieldId"`
	NamedCharacter        metadata.NamedCharacterItem        `json:"namedCharacter"`
	InstanceTagKey        metadata.InstanceTagKeyItem        `json:"instanceTagKey"`
	InstanceTagValue      metadata.InstanceTagValueItem      `json:"instanceTagValue"`
	BusinessTopoInstNames metadata.BusinessTopoInstNamesItem `json:"businessTopoInstNames"`
}

// Validate the fields of ValidationRulesCfg
func (v ValidationRulesCfg) Validate() error {
	vr := reflect.ValueOf(v)
	vrt := reflect.TypeOf(v)
	for i := 0; i < vr.NumField(); i++ {
		field := vr.Field(i)
		bc := field.FieldByName("BaseCfgItem").Interface().(metadata.BaseCfgItem)
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

// ObjectString used to admin object string Config
type ObjectString string

// Validate the fields of ObjectString
func (s ObjectString) Validate() error {
	if strings.TrimSpace(string(s)) == "" {
		return fmt.Errorf("site  value can't be empty")
	}
	return nil
}

// GlobalModule Conifg, idleName, FaultName and RecycleName cannot be deleted.
type GlobalModule struct {
	IdleName    string                    `json:"idle"`
	FaultName   string                    `json:"fault"`
	RecycleName string                    `json:"recycle"`
	UserModules []metadata.UserModuleList `json:"user_modules"`
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

// IDGeneratorConf is id generator config
type IDGeneratorConf struct {
	Enabled bool                       `json:"enabled"`
	Step    int                        `json:"step"`
	InitID  map[idgen.IDGenType]uint64 `json:"init_id,omitempty"`
	// CurrentID is the current id of each resource, this is only used for ui display
	CurrentID map[idgen.IDGenType]uint64 `json:"current_id,omitempty"`
}

// ConfigAdmin used to admin the cmdb config
type ConfigAdmin struct {
	Backend         BackendCfg         `json:"backend"`
	ValidationRules ValidationRulesCfg `json:"validationRules"`
}

// EncodeWithBase64 encode the value of ValidationRules to base64
func (c *ConfigAdmin) EncodeWithBase64() error {
	vr := reflect.ValueOf(&c.ValidationRules).Elem()
	vrt := reflect.TypeOf(c.ValidationRules)
	for i := 0; i < vr.NumField(); i++ {
		field := vr.Field(i)
		bc := field.FieldByName("BaseCfgItem")
		value := bc.Interface().(metadata.BaseCfgItem).Value
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("%s can't be empty", vrt.Field(i).Name)
		}
		base64Val := base64.StdEncoding.EncodeToString([]byte(value))
		bc.FieldByName("Value").SetString(base64Val)
	}
	return nil
}

// Validate validate the fields of BackendCfg.
func (b AdminBackendCfg) Validate() error {
	if b.MaxBizTopoLevel < minBizTopoLevel || b.MaxBizTopoLevel > maxBizTopoLevel {
		return fmt.Errorf("max biz topo level value must in range [%d-%d]", minBizTopoLevel, maxBizTopoLevel)
	}
	return nil
}

const (
	maxBizTopoLevel = 10
	minBizTopoLevel = 3
)

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

// OldPlatformSettingConfig old platform setting config
type OldPlatformSettingConfig struct {
	Backend             OldAdminBackendCfg `json:"backend"`
	ValidationRules     ValidationRulesCfg `json:"validation_rules"`
	BuiltInSetName      ObjectString       `json:"set"`
	BuiltInModuleConfig GlobalModule       `json:"idle_pool"`
}

// OldAdminBackendCfg old admin backend config
type OldAdminBackendCfg struct {
	MaxBizTopoLevel int64  `json:"max_biz_topo_level"`
	SnapshotBizName string `json:"snapshot_biz_name"`
	SnapshotBizID   int64  `json:"snapshot_biz_id"`
}

// Validate validate the fields of BackendCfg.
func (o OldAdminBackendCfg) Validate() error {
	if strings.TrimSpace(o.SnapshotBizName) == "" {
		return fmt.Errorf("snapshot biz name can't be empty")
	}

	if o.MaxBizTopoLevel < minBizTopoLevel || o.MaxBizTopoLevel > maxBizTopoLevel {
		return fmt.Errorf("max biz topo level value must in range [%d-%d]", minBizTopoLevel, maxBizTopoLevel)
	}
	return nil
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
		value := bc.Interface().(metadata.BaseCfgItem).Value
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("%s can't be empty", vrt.Field(i).Name)
		}
		base64Val := base64.StdEncoding.EncodeToString([]byte(value))
		bc.FieldByName("Value").SetString(base64Val)
	}
	return nil
}
