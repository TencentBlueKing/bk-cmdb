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

// PropertyGroupCondition used to reflect the property group json
type PropertyGroupCondition struct {
	Condition map[string]interface{} `json:"condition"`
	Data      map[string]interface{} `json:"data"`
}

// PropertyGroupObjectAtt uset to update or delete the property group object attribute
type PropertyGroupObjectAtt struct {
	Condition struct {
		OwnerID    string `json:"bk_supplier_account"`
		ObjectID   string `json:"bk_obj_id"`
		PropertyID string `json:"bk_property_id"`
	} `json:"condition"`
	Data struct {
		PropertyGroupID string `json:"bk_property_group"`
		PropertyIndex   int    `json:"bk_property_index"`
	} `json:"data"`
}
