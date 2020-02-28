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
	"time"

	"configcenter/src/common/mapstr"
)

type CloudAccount struct {
	AccountName string      `json:"bk_account_name" bson:"bk_account_name"`
	CloudVendor AccountType `json:"bk_cloud_vendor" bson:"bk_cloud_vendor"`
	AccountID   int64       `json:"bk_account_id" bson:"bk_account_id"`
	SecreteID   string      `json:"bk_secret_id" bson:"bk_secret_id"`
	SecreteKey  string      `json:"bk_secret_key" bson:"bk_secret_key"`
	Description string      `json:"bk_description" bson:"bk_description"`
	OwnerID     string      `json:"bk_supplier_account" bson:"bk_supplier_account"`
	Creator     string      `json:"bk_creator" bson:"bk_creator"`
	LastEditor  string      `json:"bk_last_editor" bson:"bk_last_editor"`
	CreateTime  time.Time   `json:"create_time" bson:"create_time"`
	LastTime    time.Time   `json:"last_time" bson:"last_time"`
}

type AccountType string

const (
	AWS          AccountType = "aws"
	TencentCloud AccountType = "tencent_cloud"
)

var SupportCloudVendors = []string{"aws", "tencent_cloud"}

type SearchCloudAccountOption struct {
	Condition mapstr.MapStr `json:"condition" bson:"condition" field:"condition"`
	Page      BasePage      `json:"page" bson:"page" field:"page"`
	Fields    []string      `json:"fields,omitempty" bson:"fields,omitempty"`
	Exact     bool          `json:"exact" bson:"exact"`
}

type MultipleCloudAccount struct {
	Count int64          `json:"count"`
	Info  []CloudAccount `json:"info"`
}

type CloudAccountVerify struct {
	SecretID    string      `json:"bk_secret_id"`
	SecretKey   string      `json:"bk_secret_key"`
	CloudVendor AccountType `json:"bk_cloud_vendor"`
}

type VpcInfo struct {
	VpcName string `json:"bk_vpc_name"`
	VpcID   string `json:"bk_vpc_id"`
	Region  string `json:"bk_region"`
}
