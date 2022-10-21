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

// ResponseModuleInstance TODO
//  只有模块的企业版本部分属性，使用前请注意
type ResponseModuleInstance struct {
	BaseResp `json:",inline"`
	Data     ModuleInstanceData `json:"data"`
}

// ModuleInstanceData TODO
//  只有模块的企业版本部分属性，使用前请注意
type ModuleInstanceData struct {
	Count int          `json:"count"`
	Info  []ModuleInst `json:"info"`
}

// ResponseSetInstance TODO
//  只有集群的企业版本部分属性，使用前请注意
type ResponseSetInstance struct {
	BaseResp `json:",inline"`
	Data     SetInstanceData `json:"data"`
}

// SetInstanceData TODO
//  只有集群的企业版本部分属性，使用前请注意
type SetInstanceData struct {
	Count int       `json:"count"`
	Info  []SetInst `json:"info"`
}

// BizSetInstanceResponse 业务集查询接口响应，只有业务集的企业版本部分属性，使用前请注意
type BizSetInstanceResponse struct {
	BaseResp `json:",inline"`
	Data     BizSetInstanceData `json:"data"`
}

// BizSetInstanceData 业务集查询接口响应数据，只有业务集的企业版本部分属性，使用前请注意
type BizSetInstanceData struct {
	Count int          `json:"count"`
	Info  []BizSetInst `json:"info"`
}
