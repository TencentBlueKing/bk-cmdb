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
 
package obj

import (
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/migrate_service/models"
	"configcenter/src/scene_server/admin_server/migrateregister"

	dbStorage "configcenter/src/storage"
)

type migrateObjDes struct {
	tableName string
}

// createTable create table
func (m *migrateObjDes) createTable(ownerID string, metaData dbStorage.DI, instData dbStorage.DI) error {
	blog.Infof("start create %s table", m.tableName)

	isExist, err := instData.HasTable(m.tableName)
	if nil != err {
		blog.Errorf("create %s table error %v", m.tableName, err)
		return err
	}
	if !isExist {
		// add instant data table
		err = instData.CreateTable(m.tableName)
		if nil != err {
			blog.Errorf("create %s table error %v", m.tableName, err)
			return err
		}
	}

	blog.Infof("end create %s table", m.tableName)

	return nil
}

func (m *migrateObjDes) addData(ownerID string, metaData dbStorage.DI, instData dbStorage.DI) error {
	err := models.AddObjDesData(m.tableName, ownerID, metaData)
	if nil != err {
		return err
	}
	return nil
}

func (m *migrateObjDes) alterTable(ownerID string, metaData dbStorage.DI, instData dbStorage.DI) error {
	err := models.AlterObjDesTable(m.tableName, metaData)
	if nil != err {
		return err
	}
	return nil
}

func init() {
	m := &migrateObjDes{tableName: "cc_ObjDes"}
	migrateregister.RegisterMigrateAction(m.createTable, migrateregister.MigrateTypeCreateTable)
	migrateregister.RegisterMigrateAction(m.alterTable, migrateregister.MigrateTypeAlterTable)
	migrateregister.RegisterMigrateAction(m.addData, migrateregister.MigrateTypeAddData)
}
