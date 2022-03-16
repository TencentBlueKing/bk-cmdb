/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/model"
)

// CreateObjectGroup create a new object group

func (s *Service) CreateObjectGroup(ctx *rest.Contexts) {
	dataWithModelBizID := MapStrWithModelBizID{}
	if err := ctx.DecodeInto(&dataWithModelBizID); err != nil {
		ctx.RespAutoError(err)
		return
	}

	var rsp model.GroupInterface
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		rsp, err = s.Core.GroupOperation().CreateObjectGroup(ctx.Kit, dataWithModelBizID.Data, dataWithModelBizID.ModelBizID)
		if nil != err {
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	retData := rsp.ToMapStr()
	ctx.RespEntity(retData)
}

// UpdateObjectGroup update the object group information
func (s *Service) UpdateObjectGroup(ctx *rest.Contexts) {
	cond := &metadata.UpdateGroupCondition{}
	err := ctx.DecodeInto(cond)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err := s.Core.GroupOperation().UpdateObjectGroup(ctx.Kit, cond)
		if nil != err {
			return err
		}

		// query attribute groups with given condition, so that update them to iam after updated
		searchCondition := condition.CreateCondition()
		if cond.Condition.ID != 0 {
			searchCondition.Field(common.BKFieldID).Eq(cond.Condition.ID)
		}
		result, err := s.Core.GroupOperation().FindObjectGroup(ctx.Kit, searchCondition, cond.ModelBizID)
		if err != nil {
			blog.Errorf("search attribute group by condition failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
			return err
		}
		attributeGroups := make([]metadata.Group, 0)
		for _, item := range result {
			attributeGroups = append(attributeGroups, item.Group())
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}

// DeleteObjectGroup delete the object group
func (s *Service) DeleteObjectGroup(ctx *rest.Contexts) {
	gid, err := strconv.ParseInt(ctx.Request.PathParameter("id"), 10, 64)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err := s.Core.GroupOperation().DeleteObjectGroup(ctx.Kit, gid)
		if nil != err {
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}

// UpdateObjectAttributeGroupProperty update the object attribute belongs to group information
func (s *Service) UpdateObjectAttributeGroupProperty(ctx *rest.Contexts) {
	requestBody := struct {
		Data       []metadata.PropertyGroupObjectAtt `json:"data" field:"json"`
		ModelBizID int64                             `json:"bk_biz_id"`
	}{}
	if err := ctx.DecodeInto(&requestBody); err != nil {
		ctx.RespAutoError(err)
		return
	}

	objectAtt := requestBody.Data
	if objectAtt == nil {
		ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrCommParamsIsInvalid, "param not set"))
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err := s.Core.GroupOperation().UpdateObjectAttributeGroup(ctx.Kit, objectAtt, requestBody.ModelBizID)
		if nil != err {
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}

// DeleteObjectAttributeGroup delete the object attribute belongs to group information

func (s *Service) DeleteObjectAttributeGroup(ctx *rest.Contexts) {
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err := s.Core.GroupOperation().DeleteObjectAttributeGroup(ctx.Kit, ctx.Request.PathParameter("bk_object_id"), ctx.Request.PathParameter("property_id"), ctx.Request.PathParameter("group_id"))
		if nil != err {
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}

// SearchGroupByObject search the groups by the object
func (s *Service) SearchGroupByObject(ctx *rest.Contexts) {
	cond := condition.CreateCondition()

	modelType := new(ModelType)
	if err := ctx.DecodeInto(modelType); err != nil {
		ctx.RespAutoError(err)
		return
	}
	resp, err := s.Core.GroupOperation().FindGroupByObject(ctx.Kit, ctx.Request.PathParameter("bk_obj_id"), cond, modelType.BizID)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(resp)

}
