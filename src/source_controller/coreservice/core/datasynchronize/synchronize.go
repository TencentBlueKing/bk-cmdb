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

package datasynchronize

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
)

type SynchronizeManager struct {
	dependent OperationDependencies
}

// New create a new model manager instance
func New(dependent OperationDependencies) core.DataSynchronizeOperation {
	return &SynchronizeManager{
		dependent: dependent,
	}
}

func (s *SynchronizeManager) SynchronizeInstanceAdapter(kit *rest.Kit, syncData *metadata.SynchronizeParameter) ([]metadata.ExceptionResult, error) {
	syncDataAdpater := NewSynchronizeInstanceAdapter(syncData)
	err := syncDataAdpater.PreSynchronizeFilter(kit)
	if err != nil {
		blog.Errorf("SynchronizeInstanceAdapter error, err:%s,rid:%s", err.Error(), kit.Rid)
		return nil, err
	}
	syncDataAdpater.SaveSynchronize(kit)
	return syncDataAdpater.GetErrorStringArr(kit)

}

func (s *SynchronizeManager) SynchronizeModelAdapter(kit *rest.Kit, syncData *metadata.SynchronizeParameter) ([]metadata.ExceptionResult, error) {
	syncDataAdpater := NewSynchronizeModelAdapter(syncData)
	err := syncDataAdpater.PreSynchronizeFilter(kit)
	if err != nil {
		return nil, err
	}
	syncDataAdpater.SaveSynchronize(kit)
	return syncDataAdpater.GetErrorStringArr(kit)

}

func (s *SynchronizeManager) SynchronizeAssociationAdapter(kit *rest.Kit, syncData *metadata.SynchronizeParameter) ([]metadata.ExceptionResult, error) {
	syncDataAdpater := NewSynchronizeAssociationAdapter(syncData)
	err := syncDataAdpater.PreSynchronizeFilter(kit)
	if err != nil {
		return nil, err
	}
	syncDataAdpater.SaveSynchronize(kit)
	return syncDataAdpater.GetErrorStringArr(kit)

}

func (s *SynchronizeManager) Find(kit *rest.Kit, input *metadata.SynchronizeFindInfoParameter) ([]mapstr.MapStr, uint64, error) {
	adapter := NewSynchronizeFindAdapter(input)
	return adapter.Find(kit)
}

func (s *SynchronizeManager) ClearData(kit *rest.Kit, input *metadata.SynchronizeClearDataParameter) error {

	adapter := NewClearData(input)
	if input.Sign == "" {
		blog.Errorf("clearData parameter synchronize_flag illegal, input:%#v,rid:%s", input, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsNeedSet, "synchronize_flag")
	}

	if !input.Legality(common.SynchronizeSignPrefix) {
		blog.Errorf("clearData parameter illegal, input:%#v,rid:%s", input, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsInvalid, input.Sign)
	}
	adapter.clearData(kit)
	return nil
}

// SetIdentifierFlag set cmdb synchronize identifier flag
func (s *SynchronizeManager) SetIdentifierFlag(kit *rest.Kit, input *metadata.SetIdenifierFlag) ([]metadata.ExceptionResult, error) {

	adapter := NewSetIdentifierFlag(input)
	if input.Flag == "" {
		blog.Errorf("SetIdentifierFlag parameter flag illegal, input:%#v,r id:%s", input, kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrCommParamsNeedSet, "flag")
	}
	if len(input.IdentifierID) == 0 {
		blog.Errorf("SetIdentifierFlag parameter identifier_id illegal, identifier_id empty. input:%#v,r id:%s", input, kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrCommParamsNeedSet, "identifier_id")
	}
	ccErr := adapter.Run(kit)
	if ccErr != nil {
		blog.Errorf("SetIdentifierFlag handle logic error. err:%s, input:%#v, rid:%s", ccErr.Error(), input, kit.Rid)
		return nil, ccErr
	}
	return nil, nil
}
