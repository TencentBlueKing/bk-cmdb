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

import "configcenter/src/common/mapstr"

type SearchInstResult struct {
	BaseResp `json:",inline"`
	Data     InstResult `json:"data"`
}

type AppBasicInfoResult struct {
	BaseResp
	Data BizBasicInfo `json:"data"`
}

type CreateModelResult struct {
	BaseResp `json:",inline"`
	Data     Object `json:"data"`
}
type SearchModelResult struct {
	BaseResp `json:",inline"`
	Data     []Object `json:"data"`
}

type SearchInnterAppTopoResult struct {
	BaseResp `json:",inline"`
	Data     InnterAppTopo
}

type MainlineObjectTopoResult struct {
	BaseResp `json:",inline"`
	Data     []MainlineObjectTopo `json:"data"`
}

type CommonInstTopo struct {
	InstNameAsst
	Count    int            `json:"count"`
	Children []InstNameAsst `json:"children"`
}

type CommonInstTopoV2 struct {
	Prev []*CommonInstTopo `json:"prev"`
	Next []*CommonInstTopo `json:"next"`
	Curr interface{}       `json:"curr"`
}
type SearchAssociationTopoResult struct {
	BaseResp `json:",inline"`
	Data     []CommonInstTopoV2 `json:"data"`
}

type SearchTopoResult struct {
	BaseResp `json:",inline"`
	Data     []*CommonInstTopo `json:"data"`
}

type QueryBusinessRequest struct {
	Fields    []string      `json:"fields"`
	Page      BasePage      `json:"page"`
	Condition mapstr.MapStr `json:"condition"`
}
