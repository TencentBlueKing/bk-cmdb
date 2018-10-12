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

package v3v0v8

import (
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	mCommon "configcenter/src/scene_server/admin_server/common"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage"
)

func addPresetObjects(db storage.DI, conf *upgrader.Config) (err error) {
	err = addClassifications(db, conf)
	if err != nil {
		return err
	}
	err = addPropertyGroupData(db, conf)
	if err != nil {
		return err
	}
	err = addObjDesData(db, conf)
	if err != nil {
		return err
	}

	err = addObjAttDescData(db, conf)
	if err != nil {
		return err
	}

	err = addAsstData(db, conf)
	if err != nil {
		return err
	}

	return nil
}

func addAsstData(db storage.DI, conf *upgrader.Config) error {
	tablename := common.BKTableNameObjAsst
	blog.Errorf("add data for  %s table ", tablename)
	rows := getAddAsstData(conf.OwnerID)
	for _, row := range rows {
		// topo mainline could be changed,so need to ignore bk_asst_obj_id
		_, _, err := upgrader.Upsert(db, tablename, row, "id", []string{common.BKObjIDField, common.BKObjAttIDField, common.BKOwnerIDField}, []string{"id", "bk_asst_obj_id"})
		if nil != err {
			blog.Errorf("add data for  %s table error  %s", tablename, err)
			return err
		}
	}

	blog.Errorf("add data for  %s table  ", tablename)
	return nil
}

func addObjAttDescData(db storage.DI, conf *upgrader.Config) error {
	tablename := common.BKTableNameObjAttDes
	blog.Infof("add data for  %s table ", tablename)
	rows := getObjAttDescData(conf.OwnerID)
	for _, row := range rows {
		_, _, err := upgrader.Upsert(db, tablename, row, "id", []string{common.BKObjIDField, common.BKPropertyIDField, common.BKOwnerIDField}, []string{})
		if nil != err {
			blog.Errorf("add data for  %s table error  %s", tablename, err)
			return err
		}
	}
	selector := map[string]interface{}{
		common.BKObjIDField: map[string]interface{}{
			common.BKDBIN: []string{"bk_switch",
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

	db.DelByCondition(tablename, selector)

	return nil
}

func addObjDesData(db storage.DI, conf *upgrader.Config) error {
	tablename := common.BKTableNameObjDes
	blog.Errorf("add data for  %s table ", tablename)
	rows := getObjectDesData(conf.OwnerID)
	for _, row := range rows {
		if _, _, err := upgrader.Upsert(db, tablename, row, "id", []string{common.BKObjIDField, common.BKClassificationIDField, common.BKOwnerIDField}, []string{"id"}); err != nil {
			blog.Errorf("add data for  %s table error  %s", tablename, err)
			return err
		}
	}

	return nil
}

func addClassifications(db storage.DI, conf *upgrader.Config) (err error) {
	tablename := common.BKTableNameObjClassifiction
	blog.Infof("add %s rows", tablename)
	for _, row := range classificationRows {
		if _, _, err = upgrader.Upsert(db, tablename, row, "id", []string{common.BKClassificationIDField}, []string{"id"}); err != nil {
			blog.Errorf("add data for  %s table error  %s", tablename, err)
			return err
		}
	}
	return
}

func addPropertyGroupData(db storage.DI, conf *upgrader.Config) error {
	tablename := common.BKTableNamePropertyGroup
	blog.Errorf("add data for  %s table ", tablename)
	rows := getPropertyGroupData(conf.OwnerID)
	for _, row := range rows {
		if _, _, err := upgrader.Upsert(db, tablename, row, "id", []string{common.BKObjIDField, "bk_group_id"}, []string{"id"}); err != nil {
			blog.Errorf("add data for  %s table error  %s", tablename, err)
			return err
		}
	}
	return nil
}
func getObjectDesData(ownerID string) []*metadata.ObjectDes {

	dataRows := []*metadata.ObjectDes{
		&metadata.ObjectDes{ObjCls: "bk_host_manage", ObjectID: common.BKInnerObjIDHost, ObjectName: "主机", IsPre: true, ObjIcon: "icon-cc-host", Position: `{"bk_host_manage":{"x":-600,"y":-650}}`},
		&metadata.ObjectDes{ObjCls: "bk_biz_topo", ObjectID: common.BKInnerObjIDModule, ObjectName: "模块", IsPre: true, ObjIcon: "icon-cc-module", Position: ``},
		&metadata.ObjectDes{ObjCls: "bk_biz_topo", ObjectID: common.BKInnerObjIDSet, ObjectName: "集群", IsPre: true, ObjIcon: "icon-cc-set", Position: ``},
		&metadata.ObjectDes{ObjCls: "bk_organization", ObjectID: common.BKInnerObjIDApp, ObjectName: "业务", IsPre: true, ObjIcon: "icon-cc-business", Position: `{"bk_organization":{"x":-100,"y":-100}}`},
		&metadata.ObjectDes{ObjCls: "bk_host_manage", ObjectID: common.BKInnerObjIDProc, ObjectName: "进程", IsPre: true, ObjIcon: "icon-cc-process", Position: `{"bk_host_manage":{"x":-450,"y":-650}}`},
		&metadata.ObjectDes{ObjCls: "bk_host_manage", ObjectID: common.BKInnerObjIDPlat, ObjectName: "云区域", IsPre: true, ObjIcon: "icon-cc-subnet", Position: `{"bk_host_manage":{"x":-600,"y":-500}}`},
		&metadata.ObjectDes{ObjCls: "bk_network", ObjectID: common.BKInnerObjIDSwitch, ObjectName: "交换机", ObjIcon: "icon-cc-switch2", Position: `{"bk_network":{"x":-200,"y":-50}}`},
		&metadata.ObjectDes{ObjCls: "bk_network", ObjectID: common.BKInnerObjIDRouter, ObjectName: "路由器", ObjIcon: "icon-cc-router", Position: `{"bk_network":{"x":-350,"y":-50}}`},
		&metadata.ObjectDes{ObjCls: "bk_network", ObjectID: common.BKInnerObjIDBlance, ObjectName: "负载均衡", ObjIcon: "icon-cc-balance", Position: `{"bk_network":{"x":-500,"y":-50}}`},
		&metadata.ObjectDes{ObjCls: "bk_network", ObjectID: common.BKInnerObjIDFirewall, ObjectName: "防火墙", ObjIcon: "icon-cc-firewall", Position: `{"bk_network":{"x":-650,"y":-50}}`},

		&metadata.ObjectDes{ObjCls: "bk_middleware", ObjectID: common.BKInnerObjIDWeblogic, ObjectName: "weblogic", ObjIcon: "icon-cc-weblogic", Position: `{"bk_middleware":{"x":-200,"y":-50}}`},
		&metadata.ObjectDes{ObjCls: "bk_middleware", ObjectID: common.BKInnerObjIDApache, ObjectName: "apache", ObjIcon: "icon-cc-apache", Position: `{"bk_middleware":{"x":-500,"y":-50}}`},
		&metadata.ObjectDes{ObjCls: "bk_middleware", ObjectID: common.BKInnerObjIDTomcat, ObjectName: "tomcat", ObjIcon: "icon-cc-tomcat", Position: `{"bk_middleware":{"x":-350,"y":-50}}`},
	}
	t := time.Now()
	for _, r := range dataRows {
		r.CreateTime = &t
		r.LastTime = &t
		r.IsPaused = false
		r.Creator = common.CCSystemOperatorUserName
		r.OwnerID = ownerID
		r.Description = ""
		r.Modifier = ""
	}

	return dataRows
}

func getAddAsstData(ownerID string) []*metadata.Association {
	dataRows := []*metadata.Association{
		&metadata.Association{OwnerID: ownerID, ObjectID: common.BKInnerObjIDSet, ObjectAttID: common.BKChildStr, AsstObjID: common.BKInnerObjIDApp},
		&metadata.Association{OwnerID: ownerID, ObjectID: common.BKInnerObjIDModule, ObjectAttID: common.BKChildStr, AsstObjID: common.BKInnerObjIDSet},
		&metadata.Association{OwnerID: ownerID, ObjectID: common.BKInnerObjIDHost, ObjectAttID: common.BKChildStr, AsstObjID: common.BKInnerObjIDModule},
		&metadata.Association{OwnerID: ownerID, ObjectID: common.BKInnerObjIDHost, ObjectAttID: common.BKCloudIDField, AsstObjID: common.BKInnerObjIDPlat},
	}
	return dataRows
}

func getObjAttDescData(ownerID string) []*metadata.Attribute {

	predataRows := AppRow()
	predataRows = append(predataRows, SetRow()...)
	predataRows = append(predataRows, ModuleRow()...)
	predataRows = append(predataRows, HostRow()...)
	predataRows = append(predataRows, ProcRow()...)
	predataRows = append(predataRows, PlatRow()...)

	dataRows := SwitchRow()
	dataRows = append(dataRows, RouterRow()...)
	dataRows = append(dataRows, LoadBalanceRow()...)
	dataRows = append(dataRows, FirewallRow()...)
	dataRows = append(dataRows, WeblogicRow()...)
	dataRows = append(dataRows, ApacheRow()...)
	dataRows = append(dataRows, TomcatRow()...)

	t := new(time.Time)
	*t = time.Now()
	for _, r := range predataRows {
		r.OwnerID = ownerID
		r.IsPre = true
		if false != r.IsEditable {
			r.IsEditable = true
		}
		r.IsReadOnly = false
		r.CreateTime = t
		r.Creator = common.CCSystemOperatorUserName
		r.LastTime = r.CreateTime
		r.Description = ""
	}
	for _, r := range dataRows {
		r.OwnerID = ownerID
		if false != r.IsEditable {
			r.IsEditable = true
		}
		r.IsReadOnly = false
		r.CreateTime = t
		r.Creator = common.CCSystemOperatorUserName
		r.LastTime = r.CreateTime
		r.Description = ""
	}

	return append(predataRows, dataRows...)
}

func getPropertyGroupData(ownerID string) []*metadata.Group {
	objectIDs := make(map[string]map[string]string)

	dataRows := []*metadata.Group{
		//app
		&metadata.Group{ObjectID: common.BKInnerObjIDApp, GroupID: mCommon.BaseInfo, GroupName: mCommon.BaseInfoName, GroupIndex: 1, OwnerID: ownerID, IsDefault: true},
		&metadata.Group{ObjectID: common.BKInnerObjIDApp, GroupID: mCommon.AppRole, GroupName: mCommon.AppRoleName, GroupIndex: 2, OwnerID: ownerID, IsDefault: true},

		//set
		&metadata.Group{ObjectID: common.BKInnerObjIDSet, GroupID: mCommon.BaseInfo, GroupName: mCommon.BaseInfoName, GroupIndex: 1, OwnerID: ownerID, IsDefault: true},

		//module
		&metadata.Group{ObjectID: common.BKInnerObjIDModule, GroupID: mCommon.BaseInfo, GroupName: mCommon.BaseInfoName, GroupIndex: 1, OwnerID: ownerID, IsDefault: true},

		//host
		&metadata.Group{ObjectID: common.BKInnerObjIDHost, GroupID: mCommon.BaseInfo, GroupName: mCommon.BaseInfoName, GroupIndex: 1, OwnerID: ownerID, IsDefault: true},
		&metadata.Group{ObjectID: common.BKInnerObjIDHost, GroupID: mCommon.HostAutoFields, GroupName: mCommon.HostAutoFieldsName, GroupIndex: 3, OwnerID: ownerID, IsDefault: true},

		//proc
		&metadata.Group{ObjectID: common.BKInnerObjIDProc, GroupID: mCommon.BaseInfo, GroupName: mCommon.BaseInfoName, GroupIndex: 1, OwnerID: ownerID, IsDefault: true},
		&metadata.Group{ObjectID: common.BKInnerObjIDProc, GroupID: mCommon.ProcPort, GroupName: mCommon.ProcPortName, GroupIndex: 2, OwnerID: ownerID, IsDefault: true},
		&metadata.Group{ObjectID: common.BKInnerObjIDProc, GroupID: mCommon.ProcGsekitBaseInfo, GroupName: mCommon.ProcGsekitBaseInfoName, GroupIndex: 3, OwnerID: ownerID, IsDefault: true},
		&metadata.Group{ObjectID: common.BKInnerObjIDProc, GroupID: mCommon.ProcGsekitManageInfo, GroupName: mCommon.ProcGsekitManageInfoName, GroupIndex: 4, OwnerID: ownerID, IsDefault: true},

		//plat
		&metadata.Group{ObjectID: common.BKInnerObjIDPlat, GroupID: mCommon.BaseInfo, GroupName: mCommon.BaseInfoName, GroupIndex: 1, OwnerID: ownerID, IsDefault: true},

		//bk_switch
		&metadata.Group{ObjectID: common.BKInnerObjIDSwitch, GroupID: mCommon.BaseInfo, GroupName: mCommon.BaseInfoName, GroupIndex: 1, OwnerID: ownerID, IsDefault: true},
		//plat
		&metadata.Group{ObjectID: common.BKInnerObjIDRouter, GroupID: mCommon.BaseInfo, GroupName: mCommon.BaseInfoName, GroupIndex: 1, OwnerID: ownerID, IsDefault: true},
		//plat
		&metadata.Group{ObjectID: common.BKInnerObjIDBlance, GroupID: mCommon.BaseInfo, GroupName: mCommon.BaseInfoName, GroupIndex: 1, OwnerID: ownerID, IsDefault: true},
		//plat
		&metadata.Group{ObjectID: common.BKInnerObjIDFirewall, GroupID: mCommon.BaseInfo, GroupName: mCommon.BaseInfoName, GroupIndex: 1, OwnerID: ownerID, IsDefault: true},
		//plat
		&metadata.Group{ObjectID: common.BKInnerObjIDWeblogic, GroupID: mCommon.BaseInfo, GroupName: mCommon.BaseInfoName, GroupIndex: 1, OwnerID: ownerID, IsDefault: true},
		//plat
		&metadata.Group{ObjectID: common.BKInnerObjIDTomcat, GroupID: mCommon.BaseInfo, GroupName: mCommon.BaseInfoName, GroupIndex: 1, OwnerID: ownerID, IsDefault: true},
		//plat
		&metadata.Group{ObjectID: common.BKInnerObjIDApache, GroupID: mCommon.BaseInfo, GroupName: mCommon.BaseInfoName, GroupIndex: 1, OwnerID: ownerID, IsDefault: true},
	}
	for objID, kv := range objectIDs {
		index := int64(1)
		for id, name := range kv {
			row := &metadata.Group{ObjectID: objID, GroupID: id, GroupName: name, GroupIndex: index, OwnerID: ownerID, IsDefault: true}
			dataRows = append(dataRows, row)
			index++
		}

	}

	return dataRows

}

var classificationRows = []*metadata.Classification{
	&metadata.Classification{ClassificationID: "bk_host_manage", ClassificationName: "主机管理", ClassificationType: "inner", ClassificationIcon: "icon-cc-host"},
	&metadata.Classification{ClassificationID: "bk_biz_topo", ClassificationName: "业务拓扑", ClassificationType: "inner", ClassificationIcon: "icon-cc-business"},
	&metadata.Classification{ClassificationID: "bk_organization", ClassificationName: "组织架构", ClassificationType: "inner", ClassificationIcon: "icon-cc-organization"},
	&metadata.Classification{ClassificationID: "bk_network", ClassificationName: "网络", ClassificationType: "inner", ClassificationIcon: "icon-cc-network-equipment"},
	&metadata.Classification{ClassificationID: "bk_middleware", ClassificationName: "中间件", ClassificationType: "inner", ClassificationIcon: "icon-cc-middleware"},
}
