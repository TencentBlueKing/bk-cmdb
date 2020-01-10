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
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

func (s *coreService) SynchronizeInstance(ctx *rest.Contexts) {
	inputData := &metadata.SynchronizeParameter{}
	if err := ctx.DecodeInto(inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	inputData.OperateDataType = metadata.SynchronizeOperateDataTypeInstance
	exceptionArr, err := s.core.DataSynchronizeOperation().SynchronizeInstanceAdapter(ctx.Kit, inputData)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(metadata.SynchronizeDataResult{Exceptions: exceptionArr})
}

func (s *coreService) SynchronizeModel(ctx *rest.Contexts) {
	inputData := &metadata.SynchronizeParameter{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	inputData.OperateDataType = metadata.SynchronizeOperateDataTypeModel
	exceptionArr, err := s.core.DataSynchronizeOperation().SynchronizeModelAdapter(ctx.Kit, inputData)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(metadata.SynchronizeDataResult{Exceptions: exceptionArr})
}

func (s *coreService) SynchronizeAssociation(ctx *rest.Contexts) {
	inputData := &metadata.SynchronizeParameter{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	inputData.OperateDataType = metadata.SynchronizeOperateDataTypeAssociation
	exceptionArr, err := s.core.DataSynchronizeOperation().SynchronizeAssociationAdapter(ctx.Kit, inputData)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(metadata.SynchronizeDataResult{Exceptions: exceptionArr})
}

func (s *coreService) SynchronizeFind(ctx *rest.Contexts) {
	inputData := &metadata.SynchronizeFindInfoParameter{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	info, cnt, err := s.core.DataSynchronizeOperation().Find(ctx.Kit, inputData)
	if err != nil {
		blog.Errorf("SynchronizeFind Find error, err:%s,input:%v,rid:%s", err.Error(), inputData, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(map[string]interface{}{"info": info, "count": cnt})
}

func (s *coreService) SynchronizeClearData(ctx *rest.Contexts) {
	inputData := &metadata.SynchronizeClearDataParameter{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	err := s.core.DataSynchronizeOperation().ClearData(ctx.Kit, inputData)
	if err != nil {
		blog.Errorf("SynchronizeClearData ClearData error, err:%s,input:%v,rid:%s", err.Error(), inputData, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

func (s *coreService) SetIdentifierFlag(ctx *rest.Contexts) {
	inputData := &metadata.SetIdenifierFlag{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	exceptionArr, err := s.core.DataSynchronizeOperation().SetIdentifierFlag(ctx.Kit, inputData)
	if err != nil {
		blog.Errorf("SetIdentifierFlag SetIdentifierFlag error, err:%s,input:%v,rid:%s", err.Error(), inputData, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(metadata.SynchronizeDataResult{Exceptions: exceptionArr})
}
