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

package business

import (
	"configcenter/src/framework/clientset/types"
	"configcenter/src/framework/common/rest"
	types2 "configcenter/src/framework/core/types"
)

type Interface interface {
	CreateBusiness(info *types.CreateBusinessCtx) (types2.MapStr, error)
	UpdateBusiness(info *types.UpdateBusinessCtx) error
	DeleteBusiness(info *types.DeleteBusinessCtx) error
	ListBusiness(info *types.ListBusinessCtx) ([]types2.MapStr, error)
}

func NewBusinessClient(client rest.ClientInterface) Interface {
	return &biz{
		client: client,
	}
}
