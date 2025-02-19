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

package y3_13_202410311500

import (
	"context"
	"errors"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader/history"
	"configcenter/src/storage/dal"
)

// attribute object attribute
type attribute struct {
	ID       int64        `json:"id" bson:"id"`
	ObjectID string       `json:"bk_obj_id" bson:"bk_obj_id"`
	Option   []enumOption `json:"option" bson:"option"`
}

// enumOption enum option
type enumOption struct {
	ID        string `json:"id" bson:"id"`
	Name      string `json:"name" bson:"name"`
	Type      string `json:"type" bson:"type"`
	IsDefault bool   `json:"is_default" bson:"is_default"`
}

func getEnumMap(vendors []string) map[string]string {
	enumMap := make(map[string]string, len(vendors))
	for index, cloudVendor := range vendors {
		enumMap[strconv.Itoa(index+1)] = cloudVendor
	}
	return enumMap
}

func getAddEnumMap() map[string]enumOption {
	return map[string]enumOption{
		"17": {
			ID:        "17",
			Name:      "腾讯自研云",
			Type:      "text",
			IsDefault: false,
		},
		"18": {
			ID:        "18",
			Name:      "Zenlayer",
			Type:      "text",
			IsDefault: false,
		},
	}
}

var (
	oldCloudVendors = []string{"AWS", "腾讯云", "GCP", "Azure", "企业私有云", "SalesForce", "Oracle Cloud", "IBM Cloud",
		"阿里云", "中国电信", "UCloud", "美团云", "金山云", "百度云", "华为云", "首都云"}

	newCloudVendors = []string{"亚马逊云", "腾讯云", "谷歌云", "微软云", "企业私有云", "SalesForce", "Oracle Cloud",
		"IBM Cloud",
		"阿里云", "中国电信", "UCloud", "美团云", "金山云", "百度云", "华为云", "首都云", "腾讯自研云", "Zenlayer"}

	oldEnumMap = getEnumMap(oldCloudVendors)

	newOptionMap = getEnumMap(newCloudVendors)
)

func updateCloudVendor(ctx context.Context, db dal.RDB, conf *history.Config) error {
	cond := map[string]interface{}{
		common.BKObjIDField: map[string]interface{}{
			common.BKDBIN: []string{common.BKInnerObjIDHost, common.BKInnerObjIDPlat},
		},
		common.BKPropertyIDField:   common.BKCloudVendor,
		common.BKPropertyTypeField: common.FieldTypeEnum,
	}

	objAttrs := make([]attribute, 0)
	err := db.Table(common.BKTableNameObjAttDes).Find(cond).All(ctx, &objAttrs)
	if err != nil {
		blog.Errorf("get cloud vendor field failed, err: %v", err)
		return err
	}
	if len(objAttrs) != 2 {
		blog.Errorf("get cloud vendor field failed, count cloud vendor field not equal 2, objAttrs: %v", objAttrs)
		return errors.New("count cloud vendor field not equal 2")
	}

	for _, attr := range objAttrs {
		addEnumMap := getAddEnumMap()
		for idx, currentEnum := range attr.Option {
			id, err := strconv.ParseInt(currentEnum.ID, 10, 64)
			if err != nil {
				blog.Errorf("parse enum id to int failed, err: %v, enum: %v", err, currentEnum)
				return errors.New("parse enum id to int failed")
			}
			if id > 18 {
				continue
			}

			if addEnum, exists := addEnumMap[currentEnum.ID]; exists {
				if currentEnum.Name != addEnum.Name {
					blog.Errorf("enum key %s already exists, enum: %v", currentEnum.ID, currentEnum)
					return errors.New("enum key already exists")
				}
				delete(addEnumMap, currentEnum.ID)
				continue
			}

			if currentEnum.Name != oldEnumMap[currentEnum.ID] && currentEnum.Name != newOptionMap[currentEnum.ID] {
				blog.Errorf("compare enum failed, currentEnum: %v", currentEnum)
				return errors.New("compare enum failed")
			}
			attr.Option[idx] = enumOption{
				ID:        currentEnum.ID,
				Name:      newOptionMap[currentEnum.ID],
				Type:      "text",
				IsDefault: false,
			}
		}

		if len(addEnumMap) != 0 {
			for _, enum := range addEnumMap {
				attr.Option = append(attr.Option, enum)
			}
		}

		updateCond := map[string]interface{}{
			common.BKObjIDField:        attr.ObjectID,
			common.BKPropertyIDField:   common.BKCloudVendor,
			common.BKPropertyTypeField: common.FieldTypeEnum,
		}
		updateData := map[string]interface{}{common.BKOptionField: attr.Option}
		if err := db.Table(common.BKTableNameObjAttDes).Update(ctx, updateCond, updateData); err != nil {
			blog.Errorf("update cloud vendor attribute failed, err: %v, cond: %v, updateData: %v", err, cond,
				updateData)
			return err
		}
	}
	return nil
}
