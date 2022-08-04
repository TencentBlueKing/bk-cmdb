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

package service

import (
	"configcenter/src/common"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/multilingual"
	"configcenter/src/storage/driver/mongodb/instancemapping"
)

// CreateOneModelInstance TODO
func (s *coreService) CreateOneModelInstance(ctx *rest.Contexts) {
	inputData := metadata.CreateModelInstance{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.InstanceOperation().CreateModelInstance(ctx.Kit, ctx.Request.PathParameter("bk_obj_id"), inputData))
}

// CreateManyModelInstances TODO
func (s *coreService) CreateManyModelInstances(ctx *rest.Contexts) {
	inputData := metadata.CreateManyModelInstance{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.InstanceOperation().CreateManyModelInstance(ctx.Kit, ctx.Request.PathParameter(common.BKObjIDField), inputData))
}

// UpdateModelInstances TODO
func (s *coreService) UpdateModelInstances(ctx *rest.Contexts) {
	inputData := metadata.UpdateOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntityWithError(s.core.InstanceOperation().UpdateModelInstance(ctx.Kit, ctx.Request.PathParameter("bk_obj_id"), inputData))
}

// SearchModelInstances TODO
func (s *coreService) SearchModelInstances(ctx *rest.Contexts) {
	inputData := metadata.QueryCondition{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	objectID := ctx.Request.PathParameter("bk_obj_id")

	// 判断是否有要根据default字段，需要国际化的内容
	if _, ok := multilingual.BuildInInstanceNamePkg[objectID]; ok {
		// 大于两个字段
		if len(inputData.Fields) > 1 {
			inputData.Fields = append(inputData.Fields, common.BKDefaultField)
		} else if len(inputData.Fields) == 1 && inputData.Fields[0] != "" {
			// 只有一个字段，如果字段为空白字符，则不处理
			inputData.Fields = append(inputData.Fields, common.BKDefaultField)
		}
	}

	dataResult, err := s.core.InstanceOperation().SearchModelInstance(ctx.Kit, objectID, inputData)
	if nil != err {
		ctx.RespEntityWithError(dataResult, err)
		return
	}

	multilingual.TranslateInstanceName(s.Language(ctx.Kit.Header), objectID, dataResult.Info)

	ctx.RespEntity(dataResult)
}

// CountModelInstances counts target model instances num.
func (s *coreService) CountModelInstances(ctx *rest.Contexts) {
	objID := ctx.Request.PathParameter("bk_obj_id")

	// decode input parameter.
	input := &metadata.Condition{}
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	// count target model instances num.
	result, err := s.core.InstanceOperation().CountModelInstances(ctx.Kit, objID, input)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(result)
}

// DeleteModelInstances TODO
func (s *coreService) DeleteModelInstances(ctx *rest.Contexts) {
	inputData := metadata.DeleteOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.InstanceOperation().DeleteModelInstance(ctx.Kit, ctx.Request.PathParameter("bk_obj_id"), inputData))
}

// CascadeDeleteModelInstances TODO
func (s *coreService) CascadeDeleteModelInstances(ctx *rest.Contexts) {
	inputData := metadata.DeleteOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.InstanceOperation().CascadeDeleteModelInstance(ctx.Kit, ctx.Request.PathParameter("bk_obj_id"), inputData))
}

// GetInstanceObjectMapping TODO
func (s *coreService) GetInstanceObjectMapping(ctx *rest.Contexts) {
	inputData := metadata.GetInstanceObjectMappingsOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(instancemapping.GetInstanceObjectMapping(inputData.IDs))
}
