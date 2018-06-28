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
	"configcenter/src/scene_server/admin_server/migrate_service/data"
	"configcenter/src/source_controller/api/metadata"
	dbStorage "configcenter/src/storage"
	"time"
)

func AddObjAttDescData(tableName, ownerID string, metaCli dbStorage.DI) error {
	blog.Infof("add data for  %s table ", tableName)
	rows := getObjAttDescData(ownerID)
	for _, row := range rows {
		selector := map[string]interface{}{
			common.BKObjIDField:      row.ObjectID,
			common.BKPropertyIDField: row.PropertyID,
			common.BKOwnerIDField:    row.OwnerID,
		}
		exist := []metadata.ObjectAttDes{}
		err := metaCli.GetMutilByCondition(tableName, nil, selector, &exist, "", 0, 1)
		if nil != err {
			blog.Errorf("add data for  %s table error  %s", tableName, err)
			return err
		}
		update := false
		if len(exist) > 0 {
			row.ID = exist[0].ID
			update = true
		}
		if row.ID <= 0 {
			id, err := metaCli.GetIncID(tableName)
			if nil != err {
				blog.Errorf("add data for  %s table error  %s", tableName, err)
				return err
			}
			row.ID = int(id)
		}
		if update {
			err = metaCli.UpdateByCondition(tableName, row, selector)
		} else {
			_, err = metaCli.Insert(tableName, row)
		}
		if nil != err {
			blog.Errorf("add data for  %s table error  %s", tableName, err)
			return err
		}
	}
	selector := map[string]interface{}{
		common.BKObjIDField: map[string]interface{}{
			"$in": []string{"bk_switch",
				"bk_router",
				"bk_load_balance",
				"bk_firewall",
				"bk_weblogic",
				"bk_tomcat",
				"bk_apache",
			},
		},
		common.BKPropertyIDField: "bk_name",
	}

	metaCli.DelByCondition(tableName, selector)

	return nil
}

func getObjAttDescData(ownerID string) []*metadata.ObjectAttDes {

	predataRows := data.AppRow()
	predataRows = append(predataRows, data.SetRow()...)
	predataRows = append(predataRows, data.ModuleRow()...)
	predataRows = append(predataRows, data.HostRow()...)
	predataRows = append(predataRows, data.ProcRow()...)
	predataRows = append(predataRows, data.PlatRow()...)

	dataRows := data.SwitchRow()
	dataRows = append(dataRows, data.RouterRow()...)
	dataRows = append(dataRows, data.LoadBalanceRow()...)
	dataRows = append(dataRows, data.FirewallRow()...)
	dataRows = append(dataRows, data.WeblogicRow()...)
	dataRows = append(dataRows, data.ApacheRow()...)
	dataRows = append(dataRows, data.TomcatRow()...)

	t := new(time.Time)
	*t = time.Now()
	for _, r := range predataRows {
		r.OwnerID = ownerID
		r.IsPre = true
		if false != r.Editable {
			r.Editable = true
		}
		r.IsReadOnly = false
		r.CreateTime = t
		r.Creator = common.CCSystemOperatorUserName
		r.LastTime = r.CreateTime
		r.Description = ""
	}
	for _, r := range dataRows {
		r.OwnerID = ownerID
		if false != r.Editable {
			r.Editable = true
		}
		r.IsReadOnly = false
		r.CreateTime = t
		r.Creator = common.CCSystemOperatorUserName
		r.LastTime = r.CreateTime
		r.Description = ""
	}

	return append(predataRows, dataRows...)
}

func AlterObjAttrDesTable(tableName string, metaCli dbStorage.DI) error {
	addCols := []*dbStorage.Column{
		dbStorage.GetMongoColumn("unit", ""),        //{Name: "Unit", Ext: " varchar(32) NOT NULL COMMENT '单位'"},
		dbStorage.GetMongoColumn("placeholder", ""), //.Column{Name: "Placeholder", Ext: " varchar(512) NOT NULL COMMENT '提示信息'"},
	}
	for _, c := range addCols {
		bl, err := metaCli.HasFields(tableName, c.Name)

		if nil != err {
			blog.Errorf("check  column %s is exist for  %s table error  %s", c.Name, tableName, err)
			return err
		}
		if !bl {
			err = metaCli.AddColumn(tableName, c)
			if nil != err {
				blog.Errorf("add column for  %s table error  %s", tableName, err)
				return err
			}
		}

	}
	return nil
}
