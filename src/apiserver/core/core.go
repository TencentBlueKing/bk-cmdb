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

package core

import (
	"github.com/emicklei/go-restful"

	"configcenter/src/common/backbone"
)

// AuthOperation auth methods
type AuthOperation interface {
}

// CompatibleV2Operation v2 api
type CompatibleV2Operation interface {
	WebService() *restful.WebService
	SetConfig(engine *backbone.Engine)
}

// Core core methods
type Core interface {
	AuthOperation() AuthOperation
}

type core struct {
	auth AuthOperation
	v2   CompatibleV2Operation
}

// New create a new core instance
func New(auth AuthOperation) Core {
	return &core{
		auth: auth,
	}
}

func (c *core) AuthOperation() AuthOperation {
	return c.auth
}
