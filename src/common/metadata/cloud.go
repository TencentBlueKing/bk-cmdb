// Package metadata TODO
/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017,-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */
package metadata

import (
	"configcenter/src/common/mapstr"
	"time"
)

// CloudArea 管控区域
type CloudArea struct {
	CloudID     int64     `json:"bk_cloud_id" bson:"bk_cloud_id"`
	CloudName   string    `json:"bk_cloud_name" bson:"bk_cloud_name"`
	Status      string    `json:"bk_status" bson:"bk_status"`
	CloudVendor string    `json:"bk_cloud_vendor" bson:"bk_cloud_vendor"`
	OwnerID     string    `json:"bk_supplier_account" bson:"bk_supplier_account"`
	VpcID       string    `json:"bk_vpc_id" bson:"bk_vpc_id"`
	VpcName     string    `json:"bk_vpc_name" bson:"bk_vpc_name"`
	Region      string    `json:"bk_region" bson:"bk_region"`
	AccountID   int64     `json:"bk_account_id" bson:"bk_account_id"`
	Creator     string    `json:"bk_creator" bson:"bk_creator"`
	LastEditor  string    `json:"bk_last_editor" bson:"bk_last_editor"`
	CreateTime  time.Time `json:"create_time" bson:"create_time"`
	LastTime    time.Time `json:"last_time" bson:"last_time"`
	Default     int64     `json:"default" bson:"default"`
}

// SearchCloudOption TODO
type SearchCloudOption struct {
	Condition mapstr.MapStr `json:"condition" bson:"condition" field:"condition"`
	Page      BasePage      `json:"page" bson:"page" field:"page"`
	Fields    []string      `json:"fields,omitempty" bson:"fields,omitempty"`
	// 对于condition里的属性值是否模糊匹配，默认为false，即不采用模糊匹配，而使用精确匹配
	IsFuzzy bool `json:"is_fuzzy" bson:"is_fuzzy"`
}

// SecretKeyResult TODO
type SecretKeyResult struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Result  bool          `json:"result"`
	Data    SecretKeyInfo `json:"data"`
}

// SecretKeyInfo TODO
type SecretKeyInfo struct {
	Content SecretContent `json:"content"`
}

// SecretContent TODO
type SecretContent struct {
	SecretKey string `json:"secret_key"`
}

// CloudMapping cloud mapping
type CloudMapping struct {
	CreateTime Time   `json:"create_time" bson:"create_time"`
	LastTime   Time   `json:"last_time" bson:"lsat_time"`
	CloudName  string `json:"bk_cloud_name" bson:"bk_cloud_name"`
	OwnerID    string `json:"bk_supplier_account" bson:"bk_supplier_account"`
	CloudID    int64  `json:"bk_cloud_id" bson:"bk_cloud_id"`
}
