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
	"configcenter/src/storage"

	"configcenter/src/source_controller/api/metadata"
	dbStorage "configcenter/src/storage"
	"time"
)

func AddObjDesData(tableName, ownerID string, metaCli dbStorage.DI) error {
	blog.Errorf("add data for  %s table ", tableName)
	rows := getObjectDesData(ownerID)
	for _, row := range rows {
		selector :=
			map[string]interface{}{
				common.BKClassificationIDField: row.ObjCls,
				common.BKObjIDField:            row.ObjectID,
				common.BKOwnerIDField:          row.OwnerID,
			}
		isExist, err := metaCli.GetCntByCondition(tableName, selector)
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

	blog.Errorf("add data for  %s table  ", tableName)
	return nil
}

func AlterObjDesTable(tableName string, metaCli dbStorage.DI) error {
	addCols := []*storage.Column{
		storage.GetMongoColumn("bk_obj_icon", ""),
		storage.GetMongoColumn("position", ""),
		//&storage.Column{Name: "ObjIcon", Ext: " varchar(50) NOT NULL COMMENT 'icon图标'"},
		//&storage.Column{Name: "Position", Ext: " varchar(1024) DEFAULT NULL COMMENT '用于存储前端显示的模型的位置信息'"},
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

func getObjectDesData(ownerID string) []*metadata.ObjectDes {

	dataRows := []*metadata.ObjectDes{
		&metadata.ObjectDes{ObjCls: "bk_host_manage", ObjectID: common.BKInnerObjIDHost, ObjectName: "主机", ObjIcon: "icon-cc-host", Position: `{"bk_host_manage":{"x":-600,"y":-650}}`},
		&metadata.ObjectDes{ObjCls: "bk_biz_topo", ObjectID: common.BKInnerObjIDModule, ObjectName: "模块", ObjIcon: "icon-cc-module", Position: ``},
		&metadata.ObjectDes{ObjCls: "bk_biz_topo", ObjectID: common.BKInnerObjIDSet, ObjectName: "集群", ObjIcon: "icon-cc-set", Position: ``},
		&metadata.ObjectDes{ObjCls: "bk_organization", ObjectID: common.BKInnerObjIDApp, ObjectName: "业务", ObjIcon: "icon-cc-business", Position: `{"bk_organization":{"x":-100,"y":-100}}`},
		&metadata.ObjectDes{ObjCls: "bk_host_manage", ObjectID: common.BKInnerObjIDProc, ObjectName: "进程", ObjIcon: "icon-cc-process", Position: `{"bk_host_manage":{"x":-450,"y":-650}}`},
		&metadata.ObjectDes{ObjCls: "bk_host_manage", ObjectID: common.BKInnerObjIDPlat, ObjectName: "子网区域", ObjIcon: "icon-cc-subnet", Position: `{"bk_host_manage":{"x":-600,"y":-500}}`},

		&metadata.ObjectDes{ObjCls: "bk_network", ObjectID: common.BKInnerObjIDSwitch, ObjectName: "交换机", ObjIcon: "icon-cc-switch2", Position: `{"bk_network":{"x":-200,"y":-50}}`},
		&metadata.ObjectDes{ObjCls: "bk_network", ObjectID: common.BKInnerObjIDRouter, ObjectName: "路由器", ObjIcon: "icon-cc-router", Position: `{"bk_network":{"x":-350,"y":-50}}`},
		&metadata.ObjectDes{ObjCls: "bk_network", ObjectID: common.BKInnerObjIDBlance, ObjectName: "防火墙", ObjIcon: "icon-cc-balance", Position: `{"bk_network":{"x":-500,"y":-50}}`},
		&metadata.ObjectDes{ObjCls: "bk_network", ObjectID: common.BKInnerObjIDFirewall, ObjectName: "负载均衡", ObjIcon: "icon-cc-firewall", Position: `{"bk_network":{"x":-650,"y":-50}}`},

		&metadata.ObjectDes{ObjCls: "bk_middleware", ObjectID: common.BKInnerObjIDWeblogic, ObjectName: "weblogic", ObjIcon: "icon-cc-weblogic", Position: `{"bk_middleware":{"x":-200,"y":-50}}`},
		&metadata.ObjectDes{ObjCls: "bk_middleware", ObjectID: common.BKInnerObjIDApache, ObjectName: "apache", ObjIcon: "icon-cc-tomcat", Position: `{"bk_middleware":{"x":-500,"y":-50}}`},
		&metadata.ObjectDes{ObjCls: "bk_middleware", ObjectID: common.BKInnerObjIDTomcat, ObjectName: "tomcat", ObjIcon: "icon-cc-apache", Position: `{"bk_middleware":{"x":-350,"y":-50}}`},
	}
	t := new(time.Time)
	*t = time.Now()
	for _, r := range dataRows {
		r.CreateTime = t
		r.LastTime = t
		r.IsPaused = false
		r.IsPre = true
		r.Creator = common.CCSystemOperatorUserName
		r.OwnerID = ownerID
		r.Description = ""
		r.Modifier = ""
	}

	return dataRows

}
