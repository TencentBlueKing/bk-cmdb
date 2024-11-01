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

package y3_14_202410311500

import (
	"context"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func updateCloudVendor(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	cloudVendors := []string{"亚马逊云", "腾讯云", "谷歌云", "微软云", "企业私有云", "SalesForce", "Oracle Cloud", "IBM Cloud",
		"阿里云", "中国电信", "UCloud", "美团云", "金山云", "百度云", "华为云", "首都云", "腾讯自研云", "Zenlayer"}

	option := make([]map[string]interface{}, len(cloudVendors))
	for index, cloudVendor := range cloudVendors {
		option[index] = map[string]interface{}{
			common.BKFieldID:   strconv.Itoa(index + 1),
			common.BKFieldName: cloudVendor,
			"type":             "text",
			"is_default":       false,
		}
	}

	cond := map[string]interface{}{
		common.BKObjIDField: map[string]interface{}{
			common.BKDBIN: []string{common.BKInnerObjIDHost, common.BKInnerObjIDPlat},
		},
		common.BKPropertyIDField: common.BKCloudVendor,
	}

	updateData := map[string]interface{}{common.BKOptionField: option}

	if err := db.Table(common.BKTableNameObjAttDes).Update(ctx, cond, updateData); err != nil {
		blog.Errorf("update cloud vendor attribute failed, err: %v, cond: %v, updateData: %v", err, cond, updateData)
		return err
	}
	return nil
}
