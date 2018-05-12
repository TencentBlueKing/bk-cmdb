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

package system

import (
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/migrate_service/models"
	"configcenter/src/scene_server/admin_server/migrateregister"
	dbStorage "configcenter/src/storage"
)

type MigrateSystem struct {
	TableName string
}

// createTable create table
func (m *MigrateSystem) createTable(ownerID string, metaData dbStorage.DI, instData dbStorage.DI) error {
	blog.Infof("start create %s table", m.TableName)

	isExist, err := instData.HasTable(m.TableName)
	if nil != err {
		blog.Errorf("create %s table error %v", m.TableName, err)
		return err
	}
	if !isExist {
		// add instant data table
		err = instData.CreateTable(m.TableName)
		if nil != err {
			blog.Errorf("create %s table error %v", m.TableName, err)
			return err
		}
	}

	blog.Infof("end create %s table", m.TableName)

	return nil
}

func (m *MigrateSystem) addData(ownerID string, metaData dbStorage.DI, instData dbStorage.DI) error {
	err := models.InitSystemData(m.TableName, instData)
	return err
}

func (m *MigrateSystem) ModifyData(ownerID string, instData dbStorage.DI) error {
	err := models.ModifySystemData(m.TableName, ownerID, instData)
	return err
}

func init() {
	m := &MigrateSystem{TableName: "cc_System"}
	migrateregister.RegisterMigrateAction(m.createTable, migrateregister.MigrateTypeCreateTable)
	migrateregister.RegisterMigrateAction(m.addData, migrateregister.MigrateTypeAddData)
}
