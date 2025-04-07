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

package modelquote

import (
	"net/http"

	"configcenter/src/apimachinery"
	"configcenter/src/common/http/rest"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/source_controller/coreservice/service/capability"
)

type service struct {
	core      core.Core
	clientSet apimachinery.ClientSetInterface
}

// InitModelQuote init model quote related service
func InitModelQuote(c *capability.Capability, clientSet apimachinery.ClientSetInterface) {
	s := &service{
		core:      c.Core,
		clientSet: clientSet,
	}

	// model quote relation
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/list/model/quote/relation",
		Handler: s.ListModelQuoteRelation})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/createmany/model/quote/relation",
		Handler: s.CreateModelQuoteRelation})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/deletemany/model/quote/relation",
		Handler: s.DeleteModelQuoteRelation})

	// quoted instance
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/createmany/quoted/model/{bk_obj_id}/instance",
		Handler: s.BatchCreateQuotedInstance})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/quoted/model/{bk_obj_id}/instance",
		Handler: s.ListQuotedInstance})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/updatemany/quoted/model/{bk_obj_id}/instance",
		Handler: s.BatchUpdateQuotedInstance})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/deletemany/quoted/model/{bk_obj_id}/instance",
		Handler: s.BatchDeleteQuotedInstance})
}
