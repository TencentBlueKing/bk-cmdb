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

	"configcenter/src/auth"
	"configcenter/src/auth/authcenter"
	authmeta "configcenter/src/auth/meta"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
)

// IAMHandler sync resource to iam
type IAMHandler struct {
	*backbone.Engine
	AuthConfig authcenter.AuthConfig
	Authorizer auth.Authorize
}

func (ih *IAMHandler) InitAuthClient() error {
	blog.Infof("new auth client with config: %+v", ih.AuthConfig)
	authorize, err := auth.NewAuthorize(nil, ih.AuthConfig)
	if err != nil {
		blog.Errorf("new auth client failed, err: %+v", err)
		return fmt.Errorf("new auth client failed, err: %+v", err)
	}
	ih.Authorizer = authorize
	return nil
}

// NewIAMHandler new a IAMHandler
func NewIAMHandler(engine *backbone.Engine, authConfig authcenter.AuthConfig) *IAMHandler {
	iamHandler := new(IAMHandler)
	iamHandler.Engine = engine
	iamHandler.AuthConfig = authConfig
	iamHandler.InitAuthClient()
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
