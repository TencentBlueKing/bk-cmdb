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

package common

// GetInstNameField returns the inst name field
func GetInstNameField(objID string) string {
	switch objID {
	case BKInnerObjIDApp:
		return BKAppNameField
	case BKInnerObjIDSet:
		return BKSetNameField
	case BKInnerObjIDModule:
		return BKModuleNameField
	case BKINnerObjIDObject:
		return BKInstNameField
	case BKInnerObjIDHost:
		return BKHostNameField
	case BKInnerObjIDProc:
		return BKProcNameField
	case BKInnerObjIDPlat:
		return BKCloudNameField
	case BKTableNameInstAsst:
		return BKFieldID
	default:
		return BKInstNameField
	}
}

func GetInstIDField(objType string) string {
	switch objType {
	case BKInnerObjIDApp:
		return BKAppIDField
	case BKInnerObjIDSet:
		return BKSetIDField
	case BKInnerObjIDModule:
		return BKModuleIDField
	case BKINnerObjIDObject:
		return BKInstIDField
	case BKInnerObjIDHost:
		return BKHostIDField
	case BKInnerObjIDProc:
		return BKProcIDField
	case BKInnerObjIDPlat:
		return BKCloudIDField
	case BKTableNameInstAsst:
		return BKFieldID
	default:
		return BKInstIDField
	}
}

func GetObjByType(objType string) string {
	switch objType {
	case BKInnerObjIDApp, BKInnerObjIDSet,
		BKInnerObjIDModule, BKInnerObjIDProc,
		BKInnerObjIDHost, BKInnerObjIDPlat:
		return objType
	default:
		return BKINnerObjIDObject
	}
}
