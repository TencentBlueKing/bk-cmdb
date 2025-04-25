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

package y3_9_202112071431

import (
	"context"

	"configcenter/src/scene_server/admin_server/upgrader/history"
	"configcenter/src/storage/dal"
)

// config 上一个版本的初始化配置文件
var config = `
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
            "value": "^[^\\\\\\|\\/:\\*,<>\"\\?#\\s]+$",
            "description": "集群/模块/拓扑节点/集群模板/服务模板名称的验证规则",
            "i18n": {
                "cn": "格式不正确，不能包含特殊字符\\ | / : * , < > \" ? #及空格",
                "en": "The format is incorrect and cannot contain special characters \\ | / : * , < > \" ? # and space"
            }
        }
    }
}
`

func updatePlatformConfigAdmin(ctx context.Context, db dal.RDB, conf *history.Config) error {

	// 将上一份配置文件写入到history中
	history.AddConfigAdminChangeItem("y3.9.202105261459", "调整拓扑实例和模板名称的校验规则的配置", config)

	// 将本次平台管理配置文件写入history中
	history.AddConfigAdminChangeItem("y3.9.202112071431", "调整平台管理配置", InitAdminConfig)

	return history.UpgradePlatConfigAdmin(ctx, db, "y3.9.202112071431")
}

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
    },
    "id_generator": {
        "enabled": false,
        "step": 1
    }
}
`
