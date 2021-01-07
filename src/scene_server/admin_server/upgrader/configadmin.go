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

package upgrader

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/version"
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
    "site": {
        "title": {
            "value": "SITE_TITLE_VAL",
            "description": "网站标题",
            "i18n": {
                "cn": "SITE_TITLE_VAL",
                "en": "CMDB | BlueKing"
            }
        },
        "footer": {
            "links": [
                {
                    "value": "http://wpa.b.qq.com/cgi/wpa.php?ln=1&key=XzgwMDgwMjAwMV80NDMwOTZfODAwODAyMDAxXzJf",
                    "description": "QQ咨询",
                    "i18n": {
                        "cn": "QQ咨询(800802001)",
                        "en": "QQ(800802001)"
                    },
                    "enabled": true
                },
                {
                    "value": "https://bk.tencent.com/s-mart/community/",
                    "description": "蓝鲸论坛",
                    "i18n": {
                        "cn": "蓝鲸论坛",
                        "en": "Blueking Forum"
                    },
                    "enabled": true
                },
                {
                    "value": "https://bk.tencent.com/index/",
                    "description": "蓝鲸官网",
                    "i18n": {
                        "cn": "蓝鲸官网",
                        "en": "BlueKing Official"
                    },
                    "enabled": true
                },
                {
                    "value": "http://your-bk-desktop.com",
                    "description": "蓝鲸桌面",
                    "i18n": {
                        "cn": "蓝鲸桌面",
                        "en": "Blueking Desktop"
                    },
                    "enabled": false
                }
            ]
        }
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
		description: "第一次初始化前的配置",
		config:      ``,
	},
	{
		dir:         "y3.8.202006092135",
		description: "第一次初始化后的配置",
		config:      initConfig,
	},
}

// setInitConfigSiteTitle 根据版本编译参数设置网站title
func setInitConfigSiteTitle() {
	switch version.CCDistro {
	case version.CCDistrCommunity:
		initConfig = strings.ReplaceAll(initConfig, "SITE_TITLE_VAL", "配置平台 | 蓝鲸智云社区版")
	case version.CCDistrEnterprise:
		initConfig = strings.ReplaceAll(initConfig, "SITE_TITLE_VAL", "配置平台 | 蓝鲸智云企业版")
	default:
		initConfig = strings.ReplaceAll(initConfig, "SITE_TITLE_VAL", "配置平台 | 蓝鲸智云社区版")
	}

	configChangeHistory[1].config = initConfig
}

// UpgradeConfigAdmin 升级配置管理
// 每次升级变更配置，需要在configChangeHistory最后追加一项要变更的配置
// 将configChangeHistory里的最后一项作为当前配置项curCfg，倒数第二项作为上一次配置项preCfg
// 需要将preCfg和db存在的配置dbCfg进行对比，对于不一致的（说明有用户调过配置管理接口做过更改）,curCfg里对应的配置不做覆盖，仍为db里的数据
func UpgradeConfigAdmin(ctx context.Context, db dal.RDB) error {
	setInitConfigSiteTitle()

	preCfg, curCfg, dbCfg, err := getConfigs(ctx, db)
	if err != nil {
		blog.Errorf("upgradeConfigAdmin failed, getConfigs err: %v", err)
		return err
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
		err = updateCofig(ctx, db, curCfg)
		if err != nil {
			blog.Errorf("UpgradeConfigAdmin failed, updateCofig err: %v, config:%#v", err, *curCfg)
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

	err = updateCofig(ctx, db, config)
	if err != nil {
		blog.Errorf("UpgradeConfigAdmin failed, updateCofig err: %v, config:%#v", err, *config)
		return err
	}

	return nil
}

// getConfigs 获取preCfg, curCfg, dbCfg
func getConfigs(ctx context.Context, db dal.RDB) (preCfg, curCfg, dbCfg *metadata.ConfigAdmin, err error) {
	length := len(configChangeHistory)
	pre := configChangeHistory[length-2].config
	cur := configChangeHistory[length-1].config

	curCfg = new(metadata.ConfigAdmin)
	if err := json.Unmarshal([]byte(cur), curCfg); err != nil {
		blog.Errorf("getAllCofnig failed, Unmarshal err: %v, config:%+v", err, cur)
		return nil, nil, nil, err
	}

	// pre为空，说明是第一次初始化,preCfg为nil
	if pre == "" {
		preCfg = nil
		return preCfg, curCfg, dbCfg, nil
	}
	preCfg = new(metadata.ConfigAdmin)
	if err = json.Unmarshal([]byte(pre), preCfg); err != nil {
		blog.Errorf("getAllCofnig failed, Unmarshal err: %v, config:%+v", err, pre)
		return nil, nil, nil, err
	}

	cond := map[string]interface{}{
		"_id": common.ConfigAdminID,
	}
	ret := struct {
		Config string `json:"config"`
	}{}
	err = db.Table(common.BKTableNameSystem).Find(cond).Fields(common.ConfigAdminValueField).One(ctx, &ret)
	if nil != err {
		blog.Errorf("getAllCofnig failed, db find err: %+v, rid: %s", err)
		return nil, nil, nil, err
	}

	dbCfg = new(metadata.ConfigAdmin)
	if err := json.Unmarshal([]byte(ret.Config), dbCfg); err != nil {
		blog.Errorf("getAllCofnig failed, Unmarshal err: %v, config:%+v", err, ret.Config)
		return nil, nil, nil, err
	}

	return preCfg, curCfg, dbCfg, nil
}

// getFinalConfig 获取最终需要保存的配置
// 将preCfg和db存在的配置dbCfg进行对比，对于不一致的（说明有用户调过配置管理接口做过更改）,curCfg里对应的配置不做覆盖，仍为db里的数据
func getFinalConfig(preCfg, curCfg, dbCfg *metadata.ConfigAdmin) *metadata.ConfigAdmin {
	if preCfg.Backend.SnapshotBizName != dbCfg.Backend.SnapshotBizName {
		curCfg.Backend.SnapshotBizName = dbCfg.Backend.SnapshotBizName
	}
	if preCfg.Backend.MaxBizTopoLevel != dbCfg.Backend.MaxBizTopoLevel {
		curCfg.Backend.MaxBizTopoLevel = dbCfg.Backend.MaxBizTopoLevel
	}
	if preCfg.Site.Title != dbCfg.Site.Title {
		curCfg.Site.Title = dbCfg.Site.Title
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

// updateCofig 将配置更新到db里
func updateCofig(ctx context.Context, db dal.RDB, config *metadata.ConfigAdmin) error {
	bytes, err := json.Marshal(config)
	if err != nil {
		blog.Errorf("updateCofig failed, Marshal err: %v, config:%+v", err, config)
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
		blog.Errorf("updateCofig failed, update err: %+v", err)
		return err
	}

	return nil
}
