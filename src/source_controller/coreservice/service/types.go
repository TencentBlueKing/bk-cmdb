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
	"configcenter/src/common/http/httpserver"
	"configcenter/src/common/mapstr"
	"configcenter/src/source_controller/coreservice/core"
)

// LogicFunc the core core function definition
type LogicFunc func(params core.ContextParams, parthParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error)

// ParamsGetter get param by key
type ParamsGetter func(name string) string

// ParseOriginDataFunc parse the origin data
type ParseOriginDataFunc func(data []byte) (mapstr.MapStr, error)

// Action the http action
type action struct {
	Method                     string
	Path                       string
	HandlerFunc                LogicFunc
	HandlerParseOriginDataFunc ParseOriginDataFunc
}

// API the API interface
type API interface {
	Actions() []*httpserver.Action
}
