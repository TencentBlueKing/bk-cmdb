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

package authorize

import (
	"configcenter/src/auth"
	"configcenter/src/auth/parser"

	restful "github.com/emicklei/go-restful"
)

func newResources(req *restful.Request, businessID int64, instanceIDs *[]int64, resourceName string, action auth.Action) (*[]auth.Resource, error) {
	apiVersion, err := parser.ParseAPIVersion(req)
	if err != nil {
		return nil, err
	}
	resources := make([]auth.Resource, len(*instanceIDs))
	for _, instanceID := range *instanceIDs {
		resource := auth.Resource{
			Name:       resourceName,
			InstanceID: instanceID,
			Action:     action,
			APIVersion: apiVersion,
			BusinessID: businessID,
		}
		resources = append(resources, resource)
	}
	return &resources, nil
}

func newHostTransferResources(req *restful.Request, businessID int64, hostIDs *[]int64) (*[]auth.Resource, error) {
	resources, err := newResources(req, businessID, hostIDs, "host", auth.TransferHost)
	return resources, err
}

func newHostViewResources(req *restful.Request, businessID int64, hostIDs *[]int64) (*[]auth.Resource, error) {
	resources, err := newResources(req, businessID, hostIDs, "host", auth.Find)
	return resources, err
}
