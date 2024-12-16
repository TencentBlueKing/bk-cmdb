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
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/scene_server/admin_server/upgrader/tools"
	"configcenter/src/storage/dal"
)

func addSystemData(kit *rest.Kit, db dal.Dal) error {
	blog.Infof("start add init data for table: %s", common.BKTableNameSystem)

	data := map[string]interface{}{common.HostCrossBizField: common.HostCrossBizValue}
	needField := &tools.InsertOptions{
		UniqueFields: []string{common.HostCrossBizField},
		IgnoreKeys:   make([]string, 0),
	}
	_, err := tools.InsertData(kit, db.Shard(kit.SysShardOpts()), common.BKTableNameSystem, []interface{}{data},
		needField)
	if err != nil {
		blog.Errorf("insert data for table %s failed, err: %v", common.BKTableNameSystem, err)
		return err
	}

	if err := initConfigAdmin(kit, db); err != nil {
		blog.Errorf("add config admin failed, error: %v", err)
		return err
	}

	blog.Infof("end add init data for table: %s", common.BKTableNameSystem)
	return nil
}

func initConfigAdmin(kit *rest.Kit, db dal.Dal) error {
	configData := map[string]interface{}{
		common.BKFieldDBID:           common.ConfigAdminID,
		common.CreateTimeField:       time.Now(),
		common.LastTimeField:         time.Now(),
		common.ConfigAdminValueField: initConfig,
	}

	needField := &tools.InsertOptions{
		UniqueFields: []string{"_id"},
		IgnoreKeys:   make([]string, 0),
	}
	_, err := tools.InsertData(kit, db.Shard(kit.SysShardOpts()), common.BKTableNameSystem, []interface{}{configData},
		needField)
	if err != nil {
		blog.Errorf("insert data for table %s failed, err: %v", common.BKTableNameSystem, err)
		return err
	}

	return nil

}

var initConfig = "{\"backend\":{\"max_biz_topo_level\":7,\"snapshot_biz_id\":2},\"validation_rules\":{\"number\":{\"value\":\"XihcLXxcKyk/XGQrJA==\",\"description\":\"字段类型“数字”的验证规则\",\"i18n\":{\"cn\":\"请输入整数数字\",\"en\":\"Please enter integer number\"}},\"float\":{\"value\":\"XlsrLV0/KFswLTldKlsuXT9bMC05XSt8WzAtOV0rWy5dP1swLTldKikoW2VFXVsrLV0/WzAtOV0rKT8k\",\"description\":\"字段类型“浮点”的验证规则\",\"i18n\":{\"cn\":\"请输入浮点型数字\",\"en\":\"Please enter float bizSetData\"}},\"singlechar\":{\"value\":\"XFMq\",\"description\":\"字段类型“短字符”的验证规则\",\"i18n\":{\"cn\":\"请输入256长度以内的字符串\",\"en\":\"Please enter the string within 256 length\"}},\"longchar\":{\"value\":\"XFMq\",\"description\":\"字段类型“长字符”的验证规则\",\"i18n\":{\"cn\":\"请输入2000长度以内的字符串\",\"en\":\"Please enter the string within 2000 length\"}},\"associationId\":{\"value\":\"XlthLXpBLVpdW1x3XSok\",\"description\":\"关联类型唯一标识验证规则\",\"i18n\":{\"cn\":\"由英文字符开头，和下划线、数字或英文组合的字符\",\"en\":\"Start with lowercase or uppercase letter, followed by lowercase / uppercase / underscore / numbers characters\"}},\"classifyId\":{\"value\":\"XlthLXpBLVpdW1x3XSok\",\"description\":\"模型分组唯一标识验证规则\",\"i18n\":{\"cn\":\"由英文字符开头，和下划线、数字或英文组合的字符\",\"en\":\"Start with lowercase or uppercase letter, followed by lowercase / uppercase / underscore / numbers characters\"}},\"modelId\":{\"value\":\"XlthLXpBLVpdW1x3XSok\",\"description\":\"模型唯一标识验证规则\",\"i18n\":{\"cn\":\"由英文字符开头，和下划线、数字或英文组合的字符\",\"en\":\"Start with lowercase or uppercase letter, followed by lowercase / uppercase / underscore / numbers characters\"}},\"enumId\":{\"value\":\"XlthLXpBLVowLTlfLV0qJA==\",\"description\":\"字段类型“枚举”ID的验证规则\",\"i18n\":{\"cn\":\"由大小写英文字母，数字，_ 或 - 组成的字符\",\"en\":\"Composed of uppercase / lowercase / numbers / - or _ characters\"}},\"enumName\":{\"value\":\"XihbYS16QS1aMC05X118W1x1NGUwMC1cdTlmYTVdfFsoKSst44CK44CLLO+8jO+8mzvigJzigJ3igJjigJnjgIJcLlwiXCcgXC86XSkqJA==\",\"description\":\"字段类型“枚举”值的验证规则\",\"i18n\":{\"cn\":\"请输入枚举值\",\"en\":\"Please enter the enum value\"}},\"fieldId\":{\"value\":\"XlthLXpBLVpdW1x3XSok\",\"description\":\"模型字段唯一标识的验证规则\",\"i18n\":{\"cn\":\"由英文字符开头，和下划线、数字或英文组合的字符\",\"en\":\"Start with lowercase or uppercase letter, followed by lowercase / uppercase / underscore / numbers characters\"}},\"namedCharacter\":{\"value\":\"XlthLXpBLVowLTlcdTRlMDAtXHU5ZmE1X1wtOlwoXCldKyQ=\",\"description\":\"服务分类名称的验证规则\",\"i18n\":{\"cn\":\"请输入中英文或特殊字符 :_- 组成的名称\",\"en\":\"Special symbols only support(:_-)\"}},\"instanceTagKey\":{\"value\":\"XlthLXpBLVpdKFthLXowLTlBLVpcLV8uXSpbYS16MC05QS1aXSk/JA==\",\"description\":\"服务实例标签键的验证规则\",\"i18n\":{\"cn\":\"请输入以英文开头的英文+数字组合\",\"en\":\"Please enter letter / number starts with letter\"}},\"instanceTagValue\":{\"value\":\"XlthLXowLTlBLVpdKFthLXowLTlBLVpcLV8uXSpbYS16MC05QS1aXSk/JA==\",\"description\":\"服务实例标签值的验证规则\",\"i18n\":{\"cn\":\"请输入英文 / 数字\",\"en\":\"Please enter letter / number\"}},\"businessTopoInstNames\":{\"value\":\"XlteXCNcLyxcPlw8XHxdKyQ=\",\"description\":\"集群/模块/实例名称的验证规则\",\"i18n\":{\"cn\":\"请输入除 #/,\\u003e\\u003c| 以外的字符\",\"en\":\"Please enter characters other than #/,\\u003e\\u003c|\"}}},\"set\":\"空闲机池\",\"idle_pool\":{\"idle\":\"空闲机\",\"fault\":\"故障机\",\"recycle\":\"待回收\",\"user_modules\":null},\"id_generator\":{\"enabled\":false,\"step\":1}}"
