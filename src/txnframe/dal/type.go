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

package dal

import "github.com/mongodb/mongo-go-driver/mongo"

// mongodb collection names
const (
	ApplicationBaseCollection    = "cc_ApplicationBase"
	HistoryCollection            = "cc_History"
	HostBaseCollection           = "cc_HostBase"
	HostFavouriteCollection      = "cc_HostFavourite"
	InstAssociationCollection    = "cc_InstAsst"
	ModuleBaseCollection         = "cc_ModuleBase"
	ModuleHostConfigCollection   = "cc_ModuleHostConfig"
	ObjectAssociationCollection  = "cc_ObjAsst"
	ObjectAttributeCollection    = "cc_ObjAttDes"
	ObjectClassifyCollection     = "cc_ObjClassification"
	ObjectDescriptionCollection  = "cc_ObjDes"
	ObjectBaseCollection         = "cc_ObjectBase"
	OperationLogCollection       = "cc_OperationLog"
	PlatBaseCollection           = "cc_PlatBase"
	PrivilegeCollection          = "cc_Privilege"
	Proc2ModuleCollection        = "cc_Proc2Module"
	ProcessCollection            = "cc_Process"
	PropertyGroupCollection      = "cc_PropertyGroup"
	BaseSetCollection            = "cc_SetBase"
	SubscriptionCollection       = "cc_Subscription"
	SystemCollection             = "cc_System"
	TopologyCollection           = "cc_TopoGraphics"
	UserAPICollection            = "cc_UserAPI"
	UserCustomCollection         = "cc_UserCustom"
	UserGroupCollection          = "cc_UserGroup"
	UserGroupPrivilegeCollection = "cc_UserGroupPrivilege"
	IDGeneratorCollection        = "cc_idgenerator"
)

type Cursor mongo.Cursor
type DeleteResult mongo.DeleteResult
type DocumentResult mongo.DocumentResult
type InsertManyResult mongo.InsertManyResult
type InsertOneResult mongo.InsertOneResult
type UpdateResult mongo.UpdateResult
type IndexView mongo.IndexView
