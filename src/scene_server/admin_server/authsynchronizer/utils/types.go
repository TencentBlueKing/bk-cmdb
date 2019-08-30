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

package utils

import (
	"fmt"
	"net/http"

	"configcenter/src/auth/extensions"
	"configcenter/src/common"
)

var (
	SynchronizeDefaultUser = "synchronize_default_user"
)

// NewListBusinessAPIHeader new a api header for list all business
func NewListBusinessAPIHeader() *http.Header {
	header := http.Header{}
	header.Add(common.BKHTTPSupplierID, fmt.Sprintf("%d", common.BKDefaultSupplierID))
	header.Add(common.BKHTTPHeaderUser, SynchronizeDefaultUser)
	header.Add(common.BKHTTPOwnerID, common.BKSuperOwnerID)
	header.Add(common.BKHTTPOwner, common.BKSuperOwnerID)
	return &header
}

func NewAPIHeaderByBusiness(businessSimplify *extensions.BusinessSimplify) *http.Header {
	header := http.Header{}
	header.Add(common.BKHTTPSupplierID, fmt.Sprintf("%d", businessSimplify.BKSupplierIDField))
	header.Add(common.BKHTTPHeaderUser, SynchronizeDefaultUser)
	header.Add(common.BKHTTPOwnerID, fmt.Sprintf("%s", businessSimplify.BKOwnerIDField))
	header.Add(common.BKHTTPOwner, common.BKSuperOwnerID)
	return &header
}
