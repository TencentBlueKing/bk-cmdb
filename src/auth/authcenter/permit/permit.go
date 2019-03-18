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
func IsPermit(rsc *meta.ResourceAttribute) bool {

	if rsc.Type == meta.ModelClassification && rsc.Action == meta.FindMany {
		return true
	}

	if rsc.Type == meta.AssociationType && rsc.Action == meta.FindMany {
		return true
	}

	if rsc.Type == meta.Model && rsc.Action == meta.FindMany {
		return true
	}

	if rsc.Type == meta.ModelAssociation && rsc.Action == meta.FindMany {
		return true
	}

	// all the model instance association related operation is all authorized for now.
	if rsc.Type == meta.ModelInstanceAssociation {
		return true
	}

	// all the network data collector related operation is all authorized for now.
	if rsc.Type == meta.NetDataCollector {
		return true
	}

	return false
}
