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

	"configcenter/src/ac/iam"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/model"
)

// CreateClassification create a new object classification
func (s *Service) CreateClassification(ctx *rest.Contexts) {
	data := make(map[string]interface{})
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	var cls model.Classification
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		cls, err = s.Core.ClassificationOperation().CreateClassification(ctx.Kit, data)
		if nil != err {
			return err
		}

		// register object classification resource creator action to iam
		if auth.EnableAuthorize() {
			iamInstance := metadata.IamInstanceWithCreator{
				Type:    string(iam.SysModelGroup),
				ID:      strconv.FormatInt(cls.Classify().ID, 10),
				Name:    cls.Classify().ClassificationName,
				Creator: ctx.Kit.User,
			}
			_, err = s.AuthManager.Authorizer.RegisterResourceCreatorAction(ctx.Kit.Ctx, ctx.Kit.Header, iamInstance)
			if err != nil {
				blog.Errorf("register created object classification to iam failed, err: %s, rid: %s", err, ctx.Kit.Rid)
				return err
			}
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(cls.ToMapStr())
}

// SearchClassificationWithObjects search the classification with objects
func (s *Service) SearchClassificationWithObjects(ctx *rest.Contexts) {
	data := mapstr.MapStr{}
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	cond := condition.CreateCondition()
	if data.Exists(metadata.PageName) {
		page, err := data.MapStr(metadata.PageName)
		if nil != err {
			blog.Errorf("failed to get the page , error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}

		if err = cond.SetPage(page); nil != err {
			blog.Errorf("failed to parse the page, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}

		data.Remove(metadata.PageName)
	}

	if err := cond.Parse(data); nil != err {
		blog.Errorf("failed to parse the condition, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	resp, err := s.Core.ClassificationOperation().FindClassificationWithObjects(ctx.Kit, cond)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(resp)
}

// SearchClassification search the classifications
func (s *Service) SearchClassification(ctx *rest.Contexts) {
	data := mapstr.MapStr{}
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	cond := condition.CreateCondition()
	if data.Exists(metadata.PageName) {

		page, err := data.MapStr(metadata.PageName)
		if nil != err {
			blog.Errorf("failed to get the page , error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}

		if err = cond.SetPage(page); nil != err {
			blog.Errorf("failed to parse the page, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}

		data.Remove(metadata.PageName)
	}
	if err := cond.Parse(data); err != nil {
		blog.Errorf("parse condition from data failed, err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	resp, err := s.Core.ClassificationOperation().FindClassification(ctx.Kit, cond)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(resp)
}

// UpdateClassification update the object classification
func (s *Service) UpdateClassification(ctx *rest.Contexts) {
	data := make(mapstr.MapStr)
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	cond := condition.CreateCondition()
	paramPath := mapstr.MapStr{}
	paramPath.Set("id", ctx.Request.PathParameter("id"))
	id, err := paramPath.Int64("id")
	if nil != err {
		blog.Errorf("[api-cls] failed to parse the path params id(%s), error info is %s , rid: %s", ctx.Request.PathParameter("id"), err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	data.Remove(metadata.BKMetadata)

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err := s.Core.ClassificationOperation().UpdateClassification(ctx.Kit, data, id, cond)
		if err != nil {
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

// DeleteClassification delete the object classification
func (s *Service) DeleteClassification(ctx *rest.Contexts) {
	cond := condition.CreateCondition()
	id, err := strconv.ParseInt(ctx.Request.PathParameter("id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-cls] failed to parse the path params id(%s), error info is %s , rid: %s", ctx.Request.PathParameter("id"), err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err = s.Core.ClassificationOperation().DeleteClassification(ctx.Kit, id, cond)
		if nil != err {
			blog.Errorf("[api-cls] failed to parse the path params id(%s), error info is %s , rid: %s", ctx.Request.PathParameter("id"), err.Error(), ctx.Kit.Rid)
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
