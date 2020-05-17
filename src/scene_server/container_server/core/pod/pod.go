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
	"gopkg.in/redis.v5"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/language"
	"configcenter/src/common/metadata"
)

// PodManager pod manager
type PodManager struct {
	clientSet apimachinery.ClientSetInterface
	language  language.CCLanguageIf
	cache     *redis.Client
}

// New create pod manager
func New(
	clientSet apimachinery.ClientSetInterface,
	language language.CCLanguageIf,
	cache *redis.Client,
) *PodManager {
	return &PodManager{
		clientSet: clientSet,
		language:  language,
		cache:     cache,
	}
}

// CreatePod implements core PodOperation
func (p *PodManager) CreatePod(kit *rest.Kit, inputParam metadata.CreatePod) (*metadata.CreatedOneOptionResult, error) {
	blog.V(3).Infof("Rid [%s] CreatePod params %#v", kit.Rid, inputParam)
	// get bk_module_id
	moduleID, err := inputParam.Pod.Int64(common.BKModuleIDField)
	if err != nil {
		blog.Errorf("get module id failed of pod %#v, err %s", inputParam.Pod, err.Error())
		err = kit.CCError.CCError(common.CCErrContainerGetPodModuleFail)
		return &metadata.CreatedOneOptionResult{
			BaseResp: metadata.BaseResp{
				Result: false,
				ErrMsg: err.Error(),
				Code:   common.CCErrContainerGetPodModuleFail,
			},
		}, err
	}
	isExisted, err := p.checkModuleIDs(kit, inputParam.BizID, []int64{moduleID})
	if err != nil || !isExisted {
		blog.Errorf("check module failed, err %s", err.Error())
		err = kit.CCError.CCError(common.CCErrContainerQueryPodModuleFail)
		return &metadata.CreatedOneOptionResult{
			BaseResp: metadata.BaseResp{
				Result: false,
				ErrMsg: err.Error(),
				Code:   common.CCErrContainerGetPodModuleFail,
			},
		}, err
	}
	// set biz id
	inputParam.Pod[common.BKAppIDField] = inputParam.BizID

	return p.clientSet.CoreService().Instance().CreateInstance(
		kit.Ctx, kit.Header, common.BKInnerObjIDPod,
		&metadata.CreateModelInstance{
			Data: inputParam.Pod,
		})
}

// CreateManyPod implements core PodOperation
func (p *PodManager) CreateManyPod(kit *rest.Kit, inputParam metadata.CreateManyPod) (*metadata.CreatedManyOptionResult, error) {
	blog.V(3).Infof("Rid [%s] CreateManyPod params %#v", kit.Rid, inputParam)
	var moduleIDArr []int64
	// check and collect module
	for _, pod := range inputParam.PodList {
		// get bk_module_id
		moduleID, err := pod.Int64(common.BKModuleIDField)
		if err != nil {
			blog.Errorf("get module id failed of pod %#v, err %s", pod, err.Error())
			return nil, kit.CCError.CCError(common.CCErrContainerGetPodModuleFail)
		}
		// set biz id
		pod[common.BKAppIDField] = inputParam.BizID
		moduleIDArr = append(moduleIDArr, moduleID)
	}

	isValid, err := p.checkModuleIDs(kit, inputParam.BizID, moduleIDArr)
	if err != nil || !isValid {
		blog.Errorf("check module failed, err %s", err.Error())
		return nil, kit.CCError.CCError(common.CCErrContainerQueryPodModuleFail)
	}

	return p.clientSet.CoreService().Instance().CreateManyInstance(
		kit.Ctx, kit.Header, common.BKInnerObjIDPod,
		&metadata.CreateManyModelInstance{
			Datas: inputParam.PodList,
		})
}

// UpdatePod implements core PodOperation
func (p *PodManager) UpdatePod(kit *rest.Kit, inputParam metadata.UpdatePod) (*metadata.UpdatedOptionResult, error) {
	blog.V(3).Infof("Rid [%s] UpdatePod params %#v", kit.Rid, inputParam)
	// get pod attr
	attrs, err := p.getPodAttrDes(kit)
	if err != nil {
		return nil, kit.CCError.CCError(common.CCErrContainerInternalError)
	}
	// get pod unique
	uniques, err := p.getPodUnique(kit)
	if err != nil {
		return nil, kit.CCError.CCError(common.CCErrContainerInternalError)
	}
	// validate update condition
	// update condition should be unique attr
	isValid := validateUpdateCondition(inputParam.Condition, uniques, attrs)
	if !isValid {
		return nil, kit.CCError.CCError(common.CCErrContainerUpdatePodConditionValidateFail)
	}

	// set biz id
	inputParam.Condition[common.BKAppIDField] = inputParam.BizID

	return p.clientSet.CoreService().Instance().UpdateInstance(kit.Ctx, kit.Header, common.BKInnerObjIDPod, &inputParam.UpdateOption)
}

// DeletePod implements core PodOperation
func (p *PodManager) DeletePod(kit *rest.Kit, inputParam metadata.DeletePod) (*metadata.DeletedOptionResult, error) {
	blog.V(3).Infof("Rid [%s] DeletePod params %#v", kit.Rid, inputParam)
	// get pod attr
	attrs, err := p.getPodAttrDes(kit)
	if err != nil {
		return nil, kit.CCError.CCError(common.CCErrContainerInternalError)
	}
	// get pod unique
	uniques, err := p.getPodUnique(kit)
	if err != nil {
		return nil, kit.CCError.CCError(common.CCErrContainerInternalError)
	}
	// validate update condition
	// update condition should be unique attr
	isValid := validateUpdateCondition(inputParam.Condition, uniques, attrs)
	if !isValid {
		return nil, kit.CCError.CCError(common.CCErrContainerUpdatePodConditionValidateFail)
	}

	// set biz id
	inputParam.Condition[common.BKAppIDField] = inputParam.BizID

	return p.clientSet.CoreService().Instance().DeleteInstance(kit.Ctx, kit.Header, common.BKInnerObjIDPod, &inputParam.DeleteOption)
}

// // DeleteManyPod implements core PodOperation
// func (p *PodManager) DeleteManyPod(kit *rest.Kit, inputParam metadata.DeleteManyPod) (*metadata.DeleteManyPodResult, error) {

// 	return nil, nil
// }

// ListPods implements core PodOperation
func (p *PodManager) ListPods(kit *rest.Kit, inputParam metadata.ListPods) (*metadata.ListPodsResult, error) {
	lister := NewLister(p.clientSet)
	return lister.ListPod(kit, inputParam)
}
