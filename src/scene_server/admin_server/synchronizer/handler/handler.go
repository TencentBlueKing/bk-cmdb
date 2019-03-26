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

package handler

import (
	"fmt"
	"strings"

	"configcenter/src/apimachinery"
	"configcenter/src/auth/authcenter"
	"configcenter/src/auth/extensions"
	authmeta "configcenter/src/auth/meta"
)

// IAMHandler sync resource to iam
type IAMHandler struct {
	clientSet apimachinery.ClientSetInterface
	authManager *extensions.AuthManager
}


// NewIAMHandler new a IAMHandler
func NewIAMHandler(clientSet apimachinery.ClientSetInterface, authManager *extensions.AuthManager) *IAMHandler {
	iamHandler := &IAMHandler{
		authManager: authManager,
		clientSet: clientSet,
	}
	return iamHandler
}

func generateCMDBResourceKey(resource *authcenter.ResourceEntity) string {
	resourcesIDs := make([]string, 0)
	for _, resourceID := range resource.ResourceID {
		resourcesIDs = append(resourcesIDs, fmt.Sprintf("%s:%s", resourceID.ResourceType, resourceID.ResourceID))
	}
	key := strings.Join(resourcesIDs, "-")
	return key
}

func generateIAMResourceKey(iamResource authmeta.BackendResource) string {
	resourcesIDs := make([]string, 0)
	for _, iamLayer := range iamResource {
		resourcesIDs = append(resourcesIDs, fmt.Sprintf("%s:%s", iamLayer.ResourceType, iamLayer.ResourceID))
	}
	key := strings.Join(resourcesIDs, "-")
	return key
}
