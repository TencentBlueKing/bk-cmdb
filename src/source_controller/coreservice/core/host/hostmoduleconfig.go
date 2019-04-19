/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package host

import (
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
)

// TransferHostToInnerModule transfer host to inner module
func (hm *hostManager) TransferHostToInnerModule(ctx core.ContextParams, input *metadata.TransferHostToInnerModule) ([]metadata.ExceptionResult, error) {
	return hm.moduleHost.TransferHostToInnerModule(ctx, input)
}

// TransferHostModule transfer host to  module
func (hm *hostManager) TransferHostModule(ctx core.ContextParams, input *metadata.HostsModuleRelation) ([]metadata.ExceptionResult, error) {
	return hm.moduleHost.TransferHostModule(ctx, input)
}

// TransferHostCrossBusiness transfer host to other business module
func (hm *hostManager) TransferHostCrossBusiness(ctx core.ContextParams, input *metadata.TransferHostsCrossBusinessRequest) ([]metadata.ExceptionResult, error) {
	return hm.moduleHost.TransferHostCrossBusiness(ctx, input)
}

func (hm *hostManager) GetHostModuleRelation(ctx core.ContextParams, input *metadata.HostModuleRelationRequest) ([]metadata.ModuleHost, error) {
	return hm.moduleHost.GetHostModuleRelation(ctx, input)
}

// DeleteHost delete host module relation and host info
func (hm *hostManager) DeleteHost(ctx core.ContextParams, input *metadata.DeleteHostRequest) ([]metadata.ExceptionResult, error) {
	return hm.moduleHost.DeleteHost(ctx, input)
}
