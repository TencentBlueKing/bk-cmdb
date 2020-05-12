/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under,
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package pod

import (
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/language"
	"configcenter/src/common/metadata"
)

//var _ core.PodOperation = (*PodManager)(nil)

// PodManager pod manager
type PodManager struct {
	clientSet apimachinery.ClientSetInterface
	language  language.CCLanguageIf
}

// New create pod manager
func New(
	clientSet apimachinery.ClientSetInterface,
	language language.CCLanguageIf,
) *PodManager {
	return &PodManager{
		clientSet: clientSet,
		language:  language,
	}
}

// CreatePod implements core PodOperation
func (p *PodManager) CreatePod(kit *rest.Kit, inputParam metadata.CreatePod) (*metadata.CreatePodResult, error) {
	// check business id

	// check module id

	// check pod id

	return nil, nil
}

// CreateManyPod implements core PodOperation
func (p *PodManager) CreateManyPod(kit *rest.Kit, inputParam metadata.CreateManyPod) (*metadata.CreateManyPodResult, error) {

	var createdResults []metadata.CreatedDataResult
	var repeatedResults []metadata.RepeatedDataResult
	var exceptionResults []metadata.ExceptionResult

	for index, pod := range inputParam.PodList {
		// get bk_module_id
		moduleID, err := pod.String(common.BKModuleIDField)
		if err != nil {
			blog.Errorf("get module id failed, err %s", pod, err.Error())
			exceptionResults = append(exceptionResults, metadata.ExceptionResult{
				Message:     kit.CCError.CCError(common.CCErrContainerGetPodModuleFail).Error(),
				Code:        common.CCErrContainerGetPodModuleFail,
				OriginIndex: int64(index),
			})
			continue
		}
		isExisted, err := p.checkModuleExist(kit, moduleID)
		if err != nil || !isExisted {
			blog.Errorf("check module failed, err %s", err.Error())
			exceptionResults = append(exceptionResults, metadata.ExceptionResult{
				Message:     kit.CCError.CCError(common.CCErrContainerGetPodModuleFail).Error(),
				Code:        common.CCErrContainerGetPodModuleFail,
				OriginIndex: int64(index),
			})
			continue
		}
		// create pod instance
		createInsResult, err := p.clientSet.CoreService().Instance().CreateInstance(
			kit.Ctx, kit.Header, common.BKInnerObjIDPod,
			&metadata.CreateModelInstance{
				Data: pod,
			})
		if err != nil {
			blog.Errorf("CreateManyPod call CreateInstance failed, err %s", err.Error())
			exceptionResults = append(exceptionResults, metadata.ExceptionResult{
				Message:     kit.CCError.CCError(common.CCErrContainerCreatePodInstanceFail).Error(),
				Code:        common.CCErrContainerCreatePodInstanceFail,
				OriginIndex: int64(index),
			})
			continue
		}
		if !createInsResult.Result {
			blog.Errorf("CreateManyPod CreateInstance return failed, err %s", err.Error())
			if createInsResult.Code == common.CCErrCoreServiceInstanceAlreadyExist {
				repeatedResults = append(repeatedResults, metadata.RepeatedDataResult{
					OriginIndex: int64(index),
				})
			} else {
				exceptionResults = append(exceptionResults, metadata.ExceptionResult{
					Message:     kit.CCError.CCError(common.CCErrContainerCreatePodInstanceFail).Error(),
					Code:        common.CCErrContainerCreatePodInstanceFail,
					OriginIndex: int64(index),
				})
			}
			continue
		}
		createdResults = append(createdResults, metadata.CreatedDataResult{
			ID:          createInsResult.Data.Created.ID,
			OriginIndex: int64(index),
		})
	}

	result := true
	code := common.CCSuccess
	if len(exceptionResults) != 0 || len(repeatedResults) != 0 {
		result = false
		code = common.CCErrContainerCreateManyPodPartialFail
	}
	if len(createdResults) == 0 {
		result = false
		code = common.CCErrContainerCreateManyPodAllFail
	}

	return &metadata.CreateManyPodResult{
		CreatedManyOptionResult: metadata.CreatedManyOptionResult{
			BaseResp: metadata.BaseResp{
				Result: result,
				Code:   code,
				ErrMsg: kit.CCError.CCError(code).Error(),
			},
			Data: metadata.CreateManyDataResult{
				CreateManyInfoResult: metadata.CreateManyInfoResult{},
			},
		},
	}, nil
}

// UpdatePod implements core PodOperation
func (p *PodManager) UpdatePod(kit *rest.Kit, podID string, inputParam metadata.UpdatePod) (*metadata.UpdatePodResult, error) {

	return nil, nil
}

// UpdateManyPod implements core PodOperation
func (p *PodManager) UpdateManyPod(kit *rest.Kit, inputParam metadata.UpdateManyPod) (*metadata.UpdateManyPodResult, error) {

	return nil, nil
}

// DeletePod implements core PodOperation
func (p *PodManager) DeletePod(kit *rest.Kit, podID string) (*metadata.DeletePodResult, error) {

	return nil, nil
}

// DeleteManyPod implements core PodOperation
func (p *PodManager) DeleteManyPod(kit *rest.Kit, inputParam metadata.DeleteManyPod) (*metadata.DeleteManyPodResult, error) {

	return nil, nil
}

// ListPod implements core PodOperation
func (p *PodManager) ListPod(kit *rest.Kit, inputParam metadata.ListPod) (*metadata.ListPodResult, error) {

	return nil, nil
}
