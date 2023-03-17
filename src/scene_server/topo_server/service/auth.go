/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package service

import (
	"fmt"

	"configcenter/src/ac/iam"
	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

func (s *Service) hasFindModelInstAuth(kit *rest.Kit, objIDs []string) (*metadata.BaseResp, bool, error) {
	if len(objIDs) == 0 {
		return nil, true, nil
	}

	cond := &metadata.QueryCondition{
		Fields: []string{common.BKFieldID},
		Page:   metadata.BasePage{Limit: common.BKNoLimit},
		Condition: map[string]interface{}{
			common.BKObjIDField: map[string]interface{}{
				common.BKDBIN: objIDs,
			},
		},
	}
	modelResp, err := s.Engine.CoreAPI.CoreService().Model().ReadModel(kit.Ctx, kit.Header, cond)
	if err != nil {
		return nil, false, err
	}

	if len(modelResp.Info) == 0 {
		return nil, true, nil
	}

	mainlineModels, err := s.getMainlineModel(kit)
	if err != nil {
		return nil, false, err
	}
	skipModels, err := s.getSkipFindAttrAuthModel(kit)
	if err != nil {
		return nil, false, err
	}

	authResources := make([]meta.ResourceAttribute, 0)
	for _, v := range modelResp.Info {
		if _, ok := skipModels[v.ObjectID]; !ok {
			authResources = append(authResources,
				meta.ResourceAttribute{Basic: meta.Basic{InstanceID: v.ID, Type: meta.Model, Action: meta.Find}})
		}

		instanceType, err := s.getInstanceTypeByObject(mainlineModels, v.ObjectID, v.ID)
		if err != nil {
			return nil, false, err
		}

		authResources = append(authResources,
			meta.ResourceAttribute{Basic: meta.Basic{Type: instanceType, Action: meta.Find}})
	}

	authResp, authorized := s.AuthManager.Authorize(kit, authResources...)
	return authResp, authorized, nil
}

// getSkipFindAttrAuthModel 主线模型和内置模型（不包括：交换机、路由器、防火墙、负载均衡）模型属性查看不鉴权
func (s *Service) getSkipFindAttrAuthModel(kit *rest.Kit) (map[string]struct{}, error) {
	models, err := s.getMainlineModel(kit)
	if err != nil {
		return nil, err
	}
	models[common.BKInnerObjIDProc] = struct{}{}
	models[common.BKInnerObjIDPlat] = struct{}{}
	models[common.BKInnerObjIDBizSet] = struct{}{}
	models[common.BKInnerObjIDProject] = struct{}{}
	return models, nil
}

func (s *Service) getInstanceTypeByObject(mainlineModel map[string]struct{}, objID string, id int64) (
	meta.ResourceType, error) {

	switch objID {
	case common.BKInnerObjIDPlat:
		return meta.CloudAreaInstance, nil
	case common.BKInnerObjIDHost:
		return meta.HostInstance, nil
	case common.BKInnerObjIDModule:
		return meta.ModelModule, nil
	case common.BKInnerObjIDSet:
		return meta.ModelSet, nil
	case common.BKInnerObjIDApp:
		return meta.Business, nil
	case common.BKInnerObjIDProc:
		return meta.Process, nil
	case common.BKInnerObjIDBizSet:
		return meta.BizSet, nil
	case common.BKInnerObjIDProject:
		return meta.Project, nil
	}

	if _, ok := mainlineModel[objID]; ok {
		return meta.MainlineInstance, nil
	}

	return iam.GenCMDBDynamicResType(id), nil
}

func (s *Service) getMainlineModel(kit *rest.Kit) (map[string]struct{}, error) {
	cond := mapstr.MapStr{common.AssociationKindIDField: common.AssociationKindMainline}
	asst, err := s.Engine.CoreAPI.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: cond})
	if err != nil {
		return nil, err
	}

	if len(asst.Info) <= 0 {
		return nil, fmt.Errorf("model association [%+v] not found", cond)
	}

	mainlineModel := make(map[string]struct{})
	for _, mainline := range asst.Info {
		mainlineModel[mainline.AsstObjID] = struct{}{}
		mainlineModel[mainline.ObjectID] = struct{}{}
	}

	return mainlineModel, nil
}

func (s *Service) hasFindModelAuth(kit *rest.Kit, objIDs []string) (*metadata.BaseResp, bool, error) {
	if len(objIDs) == 0 {
		return nil, true, nil
	}

	models, err := s.getSkipFindAttrAuthModel(kit)
	if err != nil {
		return nil, false, err
	}
	finalObjIDs := make([]string, 0)
	for _, obj := range objIDs {
		if _, ok := models[obj]; !ok {
			finalObjIDs = append(finalObjIDs, obj)
		}
	}

	if len(finalObjIDs) == 0 {
		return nil, true, nil
	}

	cond := &metadata.QueryCondition{
		Fields: []string{common.BKFieldID},
		Page:   metadata.BasePage{Limit: common.BKNoLimit},
		Condition: map[string]interface{}{
			common.BKObjIDField: map[string]interface{}{
				common.BKDBIN: finalObjIDs,
			},
		},
	}
	modelResp, err := s.Engine.CoreAPI.CoreService().Model().ReadModel(kit.Ctx, kit.Header, cond)
	if err != nil {
		return nil, false, err
	}
	if len(modelResp.Info) == 0 {
		return nil, true, nil
	}

	authResources := make([]meta.ResourceAttribute, len(modelResp.Info))
	for k, v := range modelResp.Info {
		authResources[k] = meta.ResourceAttribute{Basic: meta.Basic{InstanceID: v.ID, Type: meta.Model,
			Action: meta.Find}}
	}

	authResp, authorized := s.AuthManager.Authorize(kit, authResources...)
	return authResp, authorized, nil
}
