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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

// CreateClassification create a new object classification
func (s *Service) CreateClassification(ctx *rest.Contexts) {
	data := make(map[string]interface{})
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	var cls *metadata.Classification
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		cls, err = s.Logics.ClassificationOperation().CreateClassification(ctx.Kit, data)
		if nil != err {
			return err
		}

		// register object classification resource creator action to iam
		if auth.EnableAuthorize() {
			iamInstance := metadata.IamInstanceWithCreator{
				Type:    string(iam.SysModelGroup),
				ID:      strconv.FormatInt(cls.ID, 10),
				Name:    cls.ClassificationName,
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
	ctx.RespEntity(cls)
}

// SearchClassificationWithObjects search the classification with objects
func (s *Service) SearchClassificationWithObjects(ctx *rest.Contexts) {
	data := mapstr.MapStr{}
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	resp, err := s.Logics.ClassificationOperation().FindClassificationWithObjects(ctx.Kit, data)
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

	resp, err := s.Logics.ClassificationOperation().FindClassification(ctx.Kit, data)
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
	data.Remove(metadata.BKMetadata)

	id, err := strconv.ParseInt(ctx.Request.PathParameter("id"), 10, 64)
	if nil != err {
		blog.Errorf("parse id(%s) failed, err: %v, rid: %s", ctx.Request.PathParameter("id"), err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err := s.Logics.ClassificationOperation().UpdateClassification(ctx.Kit, data, id)
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
	idParam := ctx.Request.PathParameter("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		blog.Errorf("parse id(%s) failed, err: %v, rid: %s", ctx.Request.PathParameter("id"), err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err = s.Logics.ClassificationOperation().DeleteClassification(ctx.Kit, id)
		if nil != err {
			blog.Errorf("delete classification with id(%s) failed, err: %v, rid: %s", idParam, err, ctx.Kit.Rid)
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
