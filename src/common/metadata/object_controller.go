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

// QueryObjectResult query object result
type QueryObjectResult struct {
	BaseResp `json:",inline"`
	Data     []Object `json:"data"`
}

// CreateObjectResult create object result
type CreateObjectResult struct {
	BaseResp `json:",inline"`
	Data     ID
}

// CreateObjectAttributeResult create object attribute result
type CreateObjectAttributeResult struct {
	BaseResp `json:",inline"`
	Data     ID
}

// AttributeWrapper  wrapper, expansion field
type AttributeWrapper struct {
	Attribute         `json:",inline"`
	AssoType          int    `json:"bk_asst_type"`
	AsstForward       string `json:"bk_asst_forward"`
	AssociationID     string `json:"bk_asst_obj_id"`
	PropertyGroupName string `json:"bk_property_group_name"`
}

// QueryObjectAttributeWrapperResult query object attribute with association info result
type QueryObjectAttributeWrapperResult struct {
	BaseResp `json:",inline"`
	Data     []AttributeWrapper `json:"data"`
}

// QueryObjectAttributeResult query object attribute result
type QueryObjectAttributeResult struct {
	BaseResp `json:",inline"`
	Data     []Attribute `json:"data"`
}

// CreateObjectGroupResult create the object group result
type CreateObjectGroupResult struct {
	BaseResp `json:",inline"`
	Data     ID
}

// QueryObjectGroupResult query the object group result
type QueryObjectGroupResult struct {
	BaseResp `json:",inline"`
	Data     []Group `json:"data"`
}

// CreateObjectClassificationResult create the object classification result
type CreateObjectClassificationResult struct {
	BaseResp `json:",inline"`
	Data     ID
}

// QueryObjectClassificationResult query the object classification result
type QueryObjectClassificationResult struct {
	BaseResp `json:",inline"`
	Data     []Classification `json:"data"`
}
