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
)

// CreateObjectGroup create a new object group

func (s *Service) CreateObjectGroup(ctx *rest.Contexts) {
	dataWithMetadata := MapStrWithMetadata{}
	if err := ctx.DecodeInto(&dataWithMetadata); err != nil {
		ctx.RespAutoError(err)
		return
	}

	err := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, s.EnableTxn, ctx.Kit.Header, func() error {
		rsp, err := s.Core.GroupOperation().CreateObjectGroup(ctx.Kit, dataWithMetadata.Data, dataWithMetadata.Metadata)
		if nil != err {
			ctx.RespAutoError(err)
			return err
		}

		// auth: register attribute group
		if err := s.AuthManager.RegisterModelAttributeGroup(ctx.Kit.Ctx, ctx.Kit.Header, rsp.Group()); err != nil {
			blog.Errorf("create object group success, but register attribute group to iam failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommRegistResourceToIAMFailed))
			return err
		}
		ctx.RespEntity(rsp.ToMapStr())
		return nil
	})

	if err != nil {
		blog.Errorf("CreateObjectGroup failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		return
	}
}

// UpdateObjectGroup update the object group information
func (s *Service) UpdateObjectGroup(ctx *rest.Contexts) {
	cond := &metadata.UpdateGroupCondition{}
	err := ctx.DecodeInto(cond)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}

	err = s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, s.EnableTxn, ctx.Kit.Header, func() error {
		err := s.Core.GroupOperation().UpdateObjectGroup(ctx.Kit, cond)
		if nil != err {
			ctx.RespAutoError(err)
			return err
		}

		// query attribute groups with given condition, so that update them to iam after updated
		searchCondition := condition.CreateCondition()
		if cond.Condition.ID != 0 {
			searchCondition.Field(common.BKFieldID).Eq(cond.Condition.ID)
		}
		result, err := s.Core.GroupOperation().FindObjectGroup(ctx.Kit, searchCondition, cond.Metadata)
		if err != nil {
			blog.Errorf("search attribute group by condition failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return err
		}
		attributeGroups := make([]metadata.Group, 0)
		for _, item := range result {
			attributeGroups = append(attributeGroups, item.Group())
		}

		// auth: register attribute group
		if err := s.AuthManager.UpdateRegisteredModelAttributeGroup(ctx.Kit.Ctx, ctx.Kit.Header, attributeGroups...); err != nil {
			blog.Errorf("update object group success, but update attribute group to iam failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommRegistResourceToIAMFailed))
			return err
		}
		ctx.RespEntity(nil)
		return nil
	})

	if err != nil {
		blog.Errorf("UpdateObjectGroup failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		return
	}
}

// DeleteObjectGroup delete the object group
func (s *Service) DeleteObjectGroup(ctx *rest.Contexts) {
	gid, err := strconv.ParseInt(ctx.Request.PathParameter("id"), 10, 64)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}

	err = s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, s.EnableTxn, ctx.Kit.Header, func() error {
		err := s.Core.GroupOperation().DeleteObjectGroup(ctx.Kit, gid)
		if nil != err {
			ctx.RespAutoError(err)
			return err
		}
		// auth: deregister attribute group
		if err := s.AuthManager.DeregisterModelAttributeGroupByID(ctx.Kit.Ctx, ctx.Kit.Header, gid); err != nil {
			blog.Errorf("delete object group failed, deregister attribute group to iam failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommUnRegistResourceToIAMFailed))
			return err
		}

		ctx.RespEntity(nil)
		return nil
	})

	if err != nil {
		blog.Errorf("DeleteObjectGroup failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		return
	}
}

// UpdateObjectAttributeGroupProperty update the object attribute belongs to group information
func (s *Service) UpdateObjectAttributeGroupProperty(ctx *rest.Contexts) {
	requestBody := struct {
		Data               []metadata.PropertyGroupObjectAtt `json:"data" field:"json"`
		*metadata.Metadata `json:"metadata"`
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

	err := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, s.EnableTxn, ctx.Kit.Header, func() error {
		err := s.Core.GroupOperation().UpdateObjectAttributeGroup(ctx.Kit, objectAtt, requestBody.Metadata)
		if nil != err {
			ctx.RespAutoError(err)
			return err
		}

		ctx.RespEntity(nil)
		return nil
	})

	if err != nil {
		blog.Errorf("UpdateObjectAttributeGroupProperty failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		return
	}
}

// DeleteObjectAttributeGroup delete the object attribute belongs to group information

func (s *Service) DeleteObjectAttributeGroup(ctx *rest.Contexts) {
	err := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, s.EnableTxn, ctx.Kit.Header, func() error {
		err := s.Core.GroupOperation().DeleteObjectAttributeGroup(ctx.Kit, ctx.Request.PathParameter("bk_object_id"), ctx.Request.PathParameter("property_id"), ctx.Request.PathParameter("group_id"))
		if nil != err {
			ctx.RespAutoError(err)
			return err
		}
		ctx.RespEntity(nil)
		return nil
	})

	if err != nil {
		blog.Errorf("DeleteObjectAttributeGroup failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		return
	}
}

// SearchGroupByObject search the groups by the object
func (s *Service) SearchGroupByObject(ctx *rest.Contexts) {
	cond := condition.CreateCondition()

	md := new(MetaShell)
	if err := ctx.DecodeInto(md); err != nil {
		ctx.RespAutoError(err)
		return
	}
	resp, err := s.Core.GroupOperation().FindGroupByObject(ctx.Kit, ctx.Request.PathParameter("bk_obj_id"), cond, md.Metadata)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(resp)

}
