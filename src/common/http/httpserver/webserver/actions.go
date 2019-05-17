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
 
package webserver

import (
	"github.com/emicklei/go-restful"
	"github.com/gin-gonic/gin"
)

//Action restful action struct
type Action struct {
	Verb    string               // Verb identifying the action ("GET", "POST", "WATCH", PROXY", etc).
	Path    string               // The path of the action
	Params  []*restful.Parameter // List of parameters associated with the action.
	Handler gin.HandlerFunc
}

func NewAction(verb, path string, params []*restful.Parameter, handler gin.HandlerFunc) *Action {
	return &Action{
		Verb:    verb,
		Path:    path,
		Params:  params,
		Handler: handler,
	}
}
