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
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

// SearchClassificationWithObjects search the classification with objects
func (s *Service) SearchTTWithObjects(ctx *rest.Contexts) {
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

	//resp, err := s.Core.ClassificationOperation().FindClassificationWithObjects(ctx.Kit, cond)
	resp, err := s.Core.TTOperation().FindTTWithObjects(ctx.Kit, cond)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(resp)
}
