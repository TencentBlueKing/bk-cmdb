/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package authcenter

var ResourceTypeMap = make(map[ResourceTypeID]ResourceType)

func init() {
	for _, bizResourceType := range expectBizResourceType {
		ResourceTypeMap[bizResourceType.ResourceTypeID] = bizResourceType
	}
	for _, sysResourceType := range expectSystemResourceType {
		ResourceTypeMap[sysResourceType.ResourceTypeID] = sysResourceType
	}
}

// IsRelatedToResourceID check whether authorization on this resourceType need resourceID
func IsRelatedToResourceID(resourceTypeID ResourceTypeID) bool {
	resourceType, exist := ResourceTypeMap[resourceTypeID]
	if exist == false {
		return false
	}
	for _, action := range resourceType.Actions {
		if action.IsRelatedResource == true {
			return true
		}
	}
	return false
}
