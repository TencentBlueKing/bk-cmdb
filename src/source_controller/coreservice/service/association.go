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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

// CreateOneAssociationKind TODO
func (s *coreService) CreateOneAssociationKind(ctx *rest.Contexts) {
	inputData := metadata.CreateAssociationKind{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	result, err := s.core.AssociationOperation().CreateAssociationKind(ctx.Kit, inputData)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

// CreateManyAssociationKind TODO
func (s *coreService) CreateManyAssociationKind(ctx *rest.Contexts) {
	inputData := metadata.CreateManyAssociationKind{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	result, err := s.core.AssociationOperation().CreateManyAssociationKind(ctx.Kit, inputData)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

// SetOneAssociationKind TODO
func (s *coreService) SetOneAssociationKind(ctx *rest.Contexts) {
	inputData := metadata.SetAssociationKind{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	result, err := s.core.AssociationOperation().SetAssociationKind(ctx.Kit, inputData)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

// SetManyAssociationKind TODO
func (s *coreService) SetManyAssociationKind(ctx *rest.Contexts) {
	inputData := metadata.SetManyAssociationKind{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	result, err := s.core.AssociationOperation().SetManyAssociationKind(ctx.Kit, inputData)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

// UpdateAssociationKind TODO
func (s *coreService) UpdateAssociationKind(ctx *rest.Contexts) {
	inputData := metadata.UpdateOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	result, err := s.core.AssociationOperation().UpdateAssociationKind(ctx.Kit, inputData)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

// DeleteAssociationKind TODO
func (s *coreService) DeleteAssociationKind(ctx *rest.Contexts) {
	inputData := metadata.DeleteOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	result, err := s.core.AssociationOperation().DeleteAssociationKind(ctx.Kit, inputData)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

// CascadeDeleteAssociationKind TODO
func (s *coreService) CascadeDeleteAssociationKind(ctx *rest.Contexts) {
	inputData := metadata.DeleteOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	result, err := s.core.AssociationOperation().CascadeDeleteAssociationKind(ctx.Kit, inputData)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

// SearchAssociationKind TODO
func (s *coreService) SearchAssociationKind(ctx *rest.Contexts) {
	inputData := metadata.QueryCondition{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	result, err := s.core.AssociationOperation().SearchAssociationKind(ctx.Kit, inputData)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	// translate
	for idx := range result.Info {
		if result.Info[idx].IsPre != nil && *result.Info[idx].IsPre {
			s.TranslateAssociationType(s.Language(ctx.Kit.Header), &result.Info[idx])
		}
	}

	ctx.RespEntity(result)
}

// CreateModelAssociation TODO
func (s *coreService) CreateModelAssociation(ctx *rest.Contexts) {
	inputData := metadata.CreateModelAssociation{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	result, err := s.core.AssociationOperation().CreateModelAssociation(ctx.Kit, inputData)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

// CreateMainlineModelAssociation TODO
func (s *coreService) CreateMainlineModelAssociation(ctx *rest.Contexts) {
	inputData := metadata.CreateModelAssociation{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	result, err := s.core.AssociationOperation().CreateMainlineModelAssociation(ctx.Kit, inputData)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

// SetModelAssociation TODO
func (s *coreService) SetModelAssociation(ctx *rest.Contexts) {
	inputData := metadata.SetModelAssociation{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	result, err := s.core.AssociationOperation().SetModelAssociation(ctx.Kit, inputData)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

// UpdateModelAssociation TODO
func (s *coreService) UpdateModelAssociation(ctx *rest.Contexts) {
	inputData := metadata.UpdateOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	result, err := s.core.AssociationOperation().UpdateModelAssociation(ctx.Kit, inputData)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

// SearchModelAssociation TODO
func (s *coreService) SearchModelAssociation(ctx *rest.Contexts) {
	inputData := metadata.QueryCondition{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	result, err := s.core.AssociationOperation().SearchModelAssociation(ctx.Kit, inputData)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

// DeleteModelAssociation TODO
func (s *coreService) DeleteModelAssociation(ctx *rest.Contexts) {
	inputData := metadata.DeleteOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	result, err := s.core.AssociationOperation().DeleteModelAssociation(ctx.Kit, inputData)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

// CascadeDeleteModelAssociation TODO
func (s *coreService) CascadeDeleteModelAssociation(ctx *rest.Contexts) {
	inputData := metadata.DeleteOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	result, err := s.core.AssociationOperation().DeleteModelAssociation(ctx.Kit, inputData)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

// CreateOneInstanceAssociation TODO
func (s *coreService) CreateOneInstanceAssociation(ctx *rest.Contexts) {
	inputData := metadata.CreateOneInstanceAssociation{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	result, err := s.core.AssociationOperation().CreateOneInstanceAssociation(ctx.Kit, inputData)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

// CreateManyInstanceAssociation TODO
func (s *coreService) CreateManyInstanceAssociation(ctx *rest.Contexts) {
	inputData := metadata.CreateManyInstanceAssociation{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	result, err := s.core.AssociationOperation().CreateManyInstanceAssociation(ctx.Kit, inputData)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

// SearchInstanceAssociation TODO
func (s *coreService) SearchInstanceAssociation(ctx *rest.Contexts) {
	inputData := metadata.InstAsstQueryCondition{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	result, err := s.core.AssociationOperation().SearchInstanceAssociation(ctx.Kit, inputData.ObjID, inputData.Cond)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

// CountInstanceAssociations counts target model instance associations num.
func (s *coreService) CountInstanceAssociations(ctx *rest.Contexts) {
	objID := ctx.Request.PathParameter("bk_obj_id")

	// decode input parameter.
	input := &metadata.Condition{}
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	// count model instance associations num.
	result, err := s.core.AssociationOperation().CountInstanceAssociations(ctx.Kit, objID, input)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

// DeleteInstanceAssociation TODO
func (s *coreService) DeleteInstanceAssociation(ctx *rest.Contexts) {
	inputData := metadata.InstAsstDeleteOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	result, err := s.core.AssociationOperation().DeleteInstanceAssociation(ctx.Kit, inputData.ObjID, inputData.Opt)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}
