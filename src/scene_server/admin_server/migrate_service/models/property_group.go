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

package models

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	mCommon "configcenter/src/scene_server/admin_server/common"
	"configcenter/src/source_controller/api/metadata"
	dbStorage "configcenter/src/storage"
)

func AddPropertyGroupData(tableName, ownerID string, metaCli dbStorage.DI) error {
	blog.Errorf("add data for  %s table ", tableName)
	rows := getPropertyGroupData(ownerID)
	for _, row := range rows {
		selectorRow :=
			map[string]interface{}{
				common.BKObjIDField: row.ObjectID,
				"bk_group_id":       row.GroupID,
			}
		isExist, err := metaCli.GetCntByCondition(tableName, selectorRow)
		if nil != err {
			blog.Errorf("add data for  %s table error  %s", tableName, err)
			return err
		}
		if isExist > 0 {
			continue
		}
		id, err := metaCli.GetIncID(tableName)
		if nil != err {
			blog.Errorf("add data for  %s table error  %s", tableName, err)
			return err
		}
		row.ID = int(id)
		_, err = metaCli.Insert(tableName, row)
		if nil != err {
			blog.Errorf("add data for  %s table error  %s", tableName, err)
			return err
		}
	}

	return nil
}

func getPropertyGroupData(ownerID string) []*metadata.PropertyGroup {
	//]string{"app", "set", "module", "host", "process"}
	objectIDs := make(map[string]map[string]string)
	//=============== 注意 =========
	//分组在map顺序，决定其显示顺序。注意
	//=============================
	/*objectIDs[common.BKInnerObjIDApp] = map[string]string{
		mCommon.Base_info: mCommon.BaseInfoName,
		mCommon.App_role:  mCommon.App_role_name,
	}
	objectIDs[common.BKInnerObjIDSet] = map[string]string{
		mCommon.Base_info: mCommon.BaseInfoName,
	}
	objectIDs[common.BKInnerObjIDModule] = map[string]string{
		mCommon.Base_info: mCommon.BaseInfoName,
	}
	objectIDs[common.BKInnerObjIDHost] = map[string]string{
		mCommon.Base_info:        mCommon.BaseInfoName,
		mCommon.Host_topology:    mCommon.Host_topology_name,
		mCommon.Host_auto_fields: mCommon.Host_auto_fields_name,
	}
	objectIDs[common.BKInnerObjIDProc] = map[string]string{
		mCommon.Base_info:               mCommon.BaseInfoName,
		mCommon.Proc_port:               mCommon.Proc_port_name,
		mCommon.Proc_gsekit_base_info:   mCommon.Proc_gsekit_base_info_name,
		mCommon.Proc_gsekit_manage_info: mCommon.Proc_gsekit_manage_info_name,
	}
	objectIDs[common.BKInnerObjIDPlat] = map[string]string{
		mCommon.Base_info: mCommon.BaseInfoName,
	}*/

	dataRows := []*metadata.PropertyGroup{
		//app
		&metadata.PropertyGroup{ObjectID: common.BKInnerObjIDApp, GroupID: mCommon.BaseInfo, GroupName: mCommon.BaseInfoName, GroupIndex: 1, OwnerID: ownerID, IsDefault: true},
		&metadata.PropertyGroup{ObjectID: common.BKInnerObjIDApp, GroupID: mCommon.AppRole, GroupName: mCommon.AppRoleName, GroupIndex: 2, OwnerID: ownerID, IsDefault: true},

		//set
		&metadata.PropertyGroup{ObjectID: common.BKInnerObjIDSet, GroupID: mCommon.BaseInfo, GroupName: mCommon.BaseInfoName, GroupIndex: 1, OwnerID: ownerID, IsDefault: true},

		//module
		&metadata.PropertyGroup{ObjectID: common.BKInnerObjIDModule, GroupID: mCommon.BaseInfo, GroupName: mCommon.BaseInfoName, GroupIndex: 1, OwnerID: ownerID, IsDefault: true},

		//host
		&metadata.PropertyGroup{ObjectID: common.BKInnerObjIDHost, GroupID: mCommon.BaseInfo, GroupName: mCommon.BaseInfoName, GroupIndex: 1, OwnerID: ownerID, IsDefault: true},
		&metadata.PropertyGroup{ObjectID: common.BKInnerObjIDHost, GroupID: mCommon.HostAutoFields, GroupName: mCommon.HostAutoFieldsName, GroupIndex: 3, OwnerID: ownerID, IsDefault: true},

		//proc
		&metadata.PropertyGroup{ObjectID: common.BKInnerObjIDProc, GroupID: mCommon.BaseInfo, GroupName: mCommon.BaseInfoName, GroupIndex: 1, OwnerID: ownerID, IsDefault: true},
		&metadata.PropertyGroup{ObjectID: common.BKInnerObjIDProc, GroupID: mCommon.ProcPort, GroupName: mCommon.ProcPortName, GroupIndex: 2, OwnerID: ownerID, IsDefault: true},
		&metadata.PropertyGroup{ObjectID: common.BKInnerObjIDProc, GroupID: mCommon.ProcGsekitBaseInfo, GroupName: mCommon.ProcGsekitBaseInfoName, GroupIndex: 3, OwnerID: ownerID, IsDefault: true},
		&metadata.PropertyGroup{ObjectID: common.BKInnerObjIDProc, GroupID: mCommon.ProcGsekitManageInfo, GroupName: mCommon.ProcGsekitManageInfoName, GroupIndex: 4, OwnerID: ownerID, IsDefault: true},

		//plat
		&metadata.PropertyGroup{ObjectID: common.BKInnerObjIDPlat, GroupID: mCommon.BaseInfo, GroupName: mCommon.BaseInfoName, GroupIndex: 1, OwnerID: ownerID, IsDefault: true},
	}
	for objID, kv := range objectIDs {
		//dataRows = append(dataRows, &metadata.PropertyGroup{ObjectID: objID, GroupID: "default", GroupName: "Default", GroupIndex: -1, OwnerID: ownerID, IsDefault: true})
		index := 1
		for id, name := range kv {
			row := &metadata.PropertyGroup{ObjectID: objID, GroupID: id, GroupName: name, GroupIndex: index, OwnerID: ownerID, IsDefault: true}
			dataRows = append(dataRows, row)
			index += 1
		}

	}

	return dataRows

}
