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
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

// TransferHostToInnerModule transfer host to inner module
// 转移到空闲机/故障机模块
func (hm *hostManager) TransferToInnerModule(kit *rest.Kit, input *metadata.TransferHostToInnerModule) error {
	return hm.hostTransfer.TransferToInnerModule(kit, input)
}

// TransferToNormalModule transfer host to normal module(modules except idle and fault module)
// 业务内主机转移
// 将主机转移到 input 表示的目标模块中
// IsIncrement 控制增量更新还是覆盖更新
func (hm *hostManager) TransferToNormalModule(kit *rest.Kit, input *metadata.HostsModuleRelation) error {
	return hm.hostTransfer.TransferToNormalModule(kit, input)
}

// TransferToAnotherBusiness transfer host to another business module
func (hm *hostManager) TransferToAnotherBusiness(kit *rest.Kit, input *metadata.TransferHostsCrossBusinessRequest) error {
	return hm.hostTransfer.TransferToAnotherBusiness(kit, input)
}

// DeleteHost delete host from cmdb
func (hm *hostManager) DeleteFromSystem(kit *rest.Kit, input *metadata.DeleteHostRequest) error {
	return hm.hostTransfer.DeleteFromSystem(kit, input)
}

// RemoveFromModule remove from one of original modules
func (hm *hostManager) RemoveFromModule(kit *rest.Kit, input *metadata.RemoveHostsFromModuleOption) error {
	return hm.hostTransfer.RemoveFromModule(kit, input)
}

func (hm *hostManager) GetHostModuleRelation(kit *rest.Kit, input *metadata.HostModuleRelationRequest) (*metadata.HostConfigData, error) {
	return hm.hostTransfer.GetHostModuleRelation(kit, input)
}

// GetDistinctHostIDsByTopoRelation get all  host ids by topology relation condition
func (hm *hostManager) GetDistinctHostIDsByTopoRelation(kit *rest.Kit, input *metadata.DistinctHostIDByTopoRelationRequest) ([]int64, error) {
	return hm.hostTransfer.GetDistinctHostIDsByTopoRelation(kit, input)
}

func (hm *hostManager) TransferResourceDirectory(kit *rest.Kit, input *metadata.TransferHostResourceDirectory) errors.CCErrorCoder {
	return hm.hostTransfer.TransferResourceDirectory(kit, input)
}
