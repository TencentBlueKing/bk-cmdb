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

package permit

import "configcenter/src/auth/meta"

// this function is used to check where this request attribute is permitted as default,
// so that it is not need to check permission status in auth center.
func ShouldSkipAuthorize(rsc *meta.ResourceAttribute) bool {

	switch {

	case rsc.Type == meta.AuditLog:
		return true
	case rsc.Type == meta.ResourceSync:
		return true
	case rsc.Type == meta.DynamicGrouping && rsc.Action == meta.FindMany || rsc.Action == meta.Find:
		return true

	case rsc.Type == meta.ModelClassification && rsc.Action == meta.FindMany:
		return true

	case rsc.Type == meta.AssociationType && rsc.Action == meta.FindMany:
		return true

	case rsc.Type == meta.Model && rsc.Action == meta.FindMany:
		return true
	case rsc.Type == meta.ModelAttribute && rsc.Action == meta.FindMany:
		return true
	case rsc.Type == meta.ModelUnique && rsc.Action == meta.FindMany:
		return true
	case rsc.Type == meta.ModelAttributeGroup && rsc.Action == meta.FindMany:
		return true

	case rsc.Type == meta.UserCustom:
		return true

	case rsc.Type == meta.ModelAssociation && rsc.Action == meta.FindMany:
		return true

	// all the model instance association related operation is all authorized for now.
	case rsc.Type == meta.ModelInstanceAssociation:
		return true

	case rsc.Type == meta.ModelInstance && (rsc.Action == meta.Find || rsc.Action == meta.FindMany):
		return true

	// all the network data collector related operation is all authorized for now.
	case rsc.Type == meta.NetDataCollector:
		return true
	// host search operation, skip.
	case rsc.Type == meta.HostInstance && rsc.Action == meta.FindMany:
		return true

	// topology instance resource types.
	case rsc.Type == meta.ModelModule || rsc.Type == meta.ModelSet || rsc.Type == meta.ModelInstanceTopology:
		return true

	default:
		return false
	}

	return false
}
