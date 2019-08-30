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
// 转移到空闲机/故障机模块
func (hm *hostManager) TransferToInnerModule(ctx core.ContextParams, input *metadata.TransferHostToInnerModule) ([]metadata.ExceptionResult, error) {
	return hm.hostTransfer.TransferToInnerModule(ctx, input)
}

// TransferToNormalModule transfer host to normal module(modules except idle and fault module)
// 业务内主机转移
// 将主机转移到 input 表示的目标模块中
// IsIncrement 控制增量更新还是覆盖更新
func (hm *hostManager) TransferToNormalModule(ctx core.ContextParams, input *metadata.HostsModuleRelation) ([]metadata.ExceptionResult, error) {
	return hm.hostTransfer.TransferToNormalModule(ctx, input)
}

// TransferToAnotherBusiness transfer host to another business module
func (hm *hostManager) TransferToAnotherBusiness(ctx core.ContextParams, input *metadata.TransferHostsCrossBusinessRequest) ([]metadata.ExceptionResult, error) {
	return hm.hostTransfer.TransferToAnotherBusiness(ctx, input)
}

// DeleteHost delete host from cmdb
func (hm *hostManager) DeleteFromSystem(ctx core.ContextParams, input *metadata.DeleteHostRequest) ([]metadata.ExceptionResult, error) {
	return hm.hostTransfer.DeleteFromSystem(ctx, input)
}

// RemoveFromModule remove from one of original modules
func (hm *hostManager) RemoveFromModule(ctx core.ContextParams, input *metadata.RemoveHostsFromModuleOption) ([]metadata.ExceptionResult, error) {
	return hm.hostTransfer.RemoveFromModule(ctx, input)
}

func (hm *hostManager) GetHostModuleRelation(ctx core.ContextParams, input *metadata.HostModuleRelationRequest) (*metadata.HostConfigData, error) {
	return hm.hostTransfer.GetHostModuleRelation(ctx, input)
}
