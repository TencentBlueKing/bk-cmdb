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

package metadata

import (
	types "configcenter/src/common/mapstr"
)

const (
	ClassificationFieldID        = "id"
	ClassFieldClassificationID   = "bk_classification_id"
	ClassFieldClassificationName = "bk_classification_name"
	ClassFieldClassificationType = "bk_classification_type"
	ClassFieldClassificationIcon = "bk_classification_icon"
)

// Classification the classification metadata definition
type Classification struct {
	ID                 int64  `field:"id" json:"id" bson:"id"`
	ClassificationID   string `field:"bk_classification_id"  json:"bk_classification_id" bson:"bk_classification_id"`
	ClassificationName string `field:"bk_classification_name" json:"bk_classification_name" bson:"bk_classification_name"`
	ClassificationType string `field:"bk_classification_type" json:"bk_classification_type" bson:"bk_classification_type"`
	ClassificationIcon string `field:"bk_classification_icon" json:"bk_classification_icon" bson:"bk_classification_icon"`
	OwnerID            string `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"  `
}

// Parse load the data from mapstr classification into classification instance
func (cli *Classification) Parse(data types.MapStr) (*Classification, error) {

	err := SetValueToStructByTags(cli, data)
	if nil != err {
		return nil, err
	}

	return cli, err
}

// ToMapStr to mapstr
func (cli *Classification) ToMapStr() types.MapStr {
	return SetValueToMapStrByTags(cli)
}
