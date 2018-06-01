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
 
package migrateregister

import (
	dbStorage "configcenter/src/storage"
)

type MigrateType int

const (
	MigrateTypeCreateTable MigrateType = iota
	MigrateTypeAlterTable
	MigrateTypeDropTable

	MigrateTypeAddData
	MigrateTypeUpdateData
	MigrateTypeDelData
)

type MigrateLogic struct {
	TableName string
}

var createTableAction []func(ownerID string, mysql dbStorage.DI, mgo dbStorage.DI) error
var alterTableAction []func(ownerID string, mysql dbStorage.DI, mgo dbStorage.DI) error
var dropTableAction []func(ownerID string, mysql dbStorage.DI, mgo dbStorage.DI) error

var addDataAction []func(ownerID string, mysql dbStorage.DI, mgo dbStorage.DI) error
var updateDataAction []func(ownerID string, mysql dbStorage.DI, mgo dbStorage.DI) error
var delDataAction []func(ownerID string, mysql dbStorage.DI, mgo dbStorage.DI) error

func RegisterMigrateAction(f func(ownerID string, mysql dbStorage.DI, mgo dbStorage.DI) error, mType MigrateType) {
	switch mType {
	case MigrateTypeCreateTable:
		createTableAction = append(createTableAction, f)
	case MigrateTypeAlterTable:
		alterTableAction = append(alterTableAction, f)
	case MigrateTypeDropTable:
		dropTableAction = append(dropTableAction, f)
	case MigrateTypeAddData:
		addDataAction = append(addDataAction, f)
	case MigrateTypeUpdateData:
		updateDataAction = append(updateDataAction, f)
	case MigrateTypeDelData:
		delDataAction = append(delDataAction, f)
	}
}

func GetMigrateAction(mType MigrateType) []func(ownerID string, mysql dbStorage.DI, mgo dbStorage.DI) error {
	switch mType {
	case MigrateTypeCreateTable:
		return createTableAction
	case MigrateTypeAlterTable:
		return alterTableAction
	case MigrateTypeDropTable:
		return dropTableAction
	case MigrateTypeAddData:
		return addDataAction
	case MigrateTypeUpdateData:
		return updateDataAction
	case MigrateTypeDelData:
		return delDataAction
	}

	return nil
}
