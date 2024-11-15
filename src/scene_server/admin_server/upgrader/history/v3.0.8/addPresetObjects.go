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
	"context"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	mCommon "configcenter/src/scene_server/admin_server/common"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func addPresetObjects(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	err = addClassifications(ctx, db, conf)
	if err != nil {
		return err
	}
	err = addPropertyGroupData(ctx, db, conf)
	if err != nil {
		return err
	}
	err = addObjDesData(ctx, db, conf)
	if err != nil {
		return err
	}

	err = addObjAttDescData(ctx, db, conf)
	if err != nil {
		return err
	}

	err = addAsstData(ctx, db, conf)
	if err != nil {
		return err
	}

	return nil
}

func addAsstData(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	tablename := common.BKTableNameObjAsst
	blog.Errorf("add data for  %s table ", tablename)
	rows := getAddAsstData(conf.TenantID)
	for _, row := range rows {
		// topo mainline could be changed,so need to ignore bk_asst_obj_id
		_, _, err := upgrader.Upsert(ctx, db, tablename, row, "id",
			[]string{common.BKObjIDField, common.BKObjAttIDField, "bk_supplier_account"},
			[]string{"id", "bk_asst_obj_id"})
		if nil != err {
			blog.Errorf("add data for  %s table error  %s", tablename, err)
			return err
		}
	}

	blog.Errorf("add data for  %s table  ", tablename)
	return nil
}

func addObjAttDescData(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	tablename := common.BKTableNameObjAttDes
	blog.Infof("add data for  %s table ", tablename)
	rows := getObjAttDescData(conf.TenantID)
	for _, row := range rows {
		_, _, err := upgrader.Upsert(ctx, db, tablename, row, "id",
			[]string{common.BKObjIDField, common.BKPropertyIDField, "bk_supplier_account"}, []string{})
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
			},
		},
		common.BKPropertyIDField: "bk_name",
	}

	db.Table(tablename).Delete(ctx, selector)

	return nil
}

func addObjDesData(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	tablename := common.BKTableNameObjDes
	blog.Errorf("add data for  %s table ", tablename)
	rows := getObjectDesData(conf.TenantID)
	for _, row := range rows {
		if _, _, err := upgrader.Upsert(ctx, db, tablename, row, "id",
			[]string{common.BKObjIDField, common.BKClassificationIDField, "bk_supplier_account"},
			[]string{"id"}); err != nil {
			blog.Errorf("add data for  %s table error  %s", tablename, err)
			return err
		}
	}

	return nil
}

func addClassifications(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	tablename := common.BKTableNameObjClassification
	blog.Infof("add %s rows", tablename)
	for _, row := range classificationRows {
		if _, _, err = upgrader.Upsert(ctx, db, tablename, row, "id", []string{common.BKClassificationIDField},
			[]string{"id"}); err != nil {
			blog.Errorf("add data for  %s table error  %s", tablename, err)
			return err
		}
	}
	return
}

func addPropertyGroupData(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	tablename := common.BKTableNamePropertyGroup
	blog.Infof("add data for %s table", tablename)
	rows := getPropertyGroupData(conf.TenantID)
	for _, row := range rows {
		_, _, err := upgrader.Upsert(ctx, db, tablename, row, "id",
			[]string{common.BKObjIDField, "bk_group_id"}, []string{"id"})
		if err != nil && !db.IsDuplicatedError(err) {
			blog.Errorf("add data for %s table error  %s", tablename, err)
			return err
		}
	}
	return nil
}
func getObjectDesData(ownerID string) []*Object {

	dataRows := []*Object{
		&Object{ObjCls: "bk_host_manage", ObjectID: common.BKInnerObjIDHost, ObjectName: "主机", IsPre: true,
			ObjIcon: "icon-cc-host", Position: `{"bk_host_manage":{"x":-600,"y":-650}}`},
		&Object{ObjCls: "bk_biz_topo", ObjectID: common.BKInnerObjIDModule, ObjectName: "模块", IsPre: true,
			ObjIcon: "icon-cc-module", Position: ``},
		&Object{ObjCls: "bk_biz_topo", ObjectID: common.BKInnerObjIDSet, ObjectName: "集群", IsPre: true,
			ObjIcon: "icon-cc-set", Position: ``},
		&Object{ObjCls: "bk_organization", ObjectID: common.BKInnerObjIDApp, ObjectName: "业务", IsPre: true,
			ObjIcon: "icon-cc-business", Position: `{"bk_organization":{"x":-100,"y":-100}}`},
		&Object{ObjCls: "bk_host_manage", ObjectID: common.BKInnerObjIDProc, ObjectName: "进程", IsPre: true,
			ObjIcon: "icon-cc-process", Position: `{"bk_host_manage":{"x":-450,"y":-650}}`},
		&Object{ObjCls: "bk_host_manage", ObjectID: common.BKInnerObjIDPlat, ObjectName: "云区域", IsPre: true,
			ObjIcon: "icon-cc-subnet", Position: `{"bk_host_manage":{"x":-600,"y":-500}}`},
		&Object{ObjCls: "bk_network", ObjectID: common.BKInnerObjIDSwitch, ObjectName: "交换机",
			ObjIcon: "icon-cc-switch2", Position: `{"bk_network":{"x":-200,"y":-50}}`},
		&Object{ObjCls: "bk_network", ObjectID: common.BKInnerObjIDRouter, ObjectName: "路由器",
			ObjIcon: "icon-cc-router", Position: `{"bk_network":{"x":-350,"y":-50}}`},
		&Object{ObjCls: "bk_network", ObjectID: common.BKInnerObjIDBlance, ObjectName: "负载均衡",
			ObjIcon: "icon-cc-balance", Position: `{"bk_network":{"x":-500,"y":-50}}`},
		&Object{ObjCls: "bk_network", ObjectID: common.BKInnerObjIDFirewall, ObjectName: "防火墙",
			ObjIcon: "icon-cc-firewall", Position: `{"bk_network":{"x":-650,"y":-50}}`},
	}
	t := metadata.Now()
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

// Association for purpose of this structure not change by other, copy here
type Association struct {
	ID               int64  `field:"id" json:"id" bson:"id"`
	ObjectID         string `field:"bk_obj_id" json:"bk_obj_id" bson:"bk_obj_id"`
	OwnerID          string `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`
	AsstForward      string `field:"bk_asst_forward" json:"bk_asst_forward" bson:"bk_asst_forward"`
	AsstObjID        string `field:"bk_asst_obj_id" json:"bk_asst_obj_id" bson:"bk_asst_obj_id"`
	AsstName         string `field:"bk_asst_name" json:"bk_asst_name" bson:"bk_asst_name"`
	ObjectAttID      string `field:"bk_object_att_id" json:"bk_object_att_id" bson:"bk_object_att_id"`
	ClassificationID string `field:"bk_classification_id" bson:"-"`
	ObjectIcon       string `field:"bk_obj_icon" bson:"-"`
	ObjectName       string `field:"bk_obj_name" bson:"-"`
}

func getAddAsstData(ownerID string) []Association {
	dataRows := []Association{
		{OwnerID: ownerID, ObjectID: common.BKInnerObjIDSet, ObjectAttID: "bk_childid",
			AsstObjID: common.BKInnerObjIDApp},
		{OwnerID: ownerID, ObjectID: common.BKInnerObjIDModule, ObjectAttID: "bk_childid",
			AsstObjID: common.BKInnerObjIDSet},
		{OwnerID: ownerID, ObjectID: common.BKInnerObjIDHost, ObjectAttID: "bk_childid",
			AsstObjID: common.BKInnerObjIDModule},
		{OwnerID: ownerID, ObjectID: common.BKInnerObjIDHost, ObjectAttID: common.BKCloudIDField,
			AsstObjID: common.BKInnerObjIDPlat},
	}
	return dataRows
}

func getObjAttDescData(ownerID string) []*Attribute {

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

func getPropertyGroupData(ownerID string) []*group {
	objectIDs := make(map[string]map[string]string)

	dataRows := []*group{
		// app
		&group{ObjectID: common.BKInnerObjIDApp, GroupID: mCommon.BaseInfo, GroupName: mCommon.BaseInfoName,
			GroupIndex: 1, OwnerID: ownerID, IsDefault: true},
		&group{ObjectID: common.BKInnerObjIDApp, GroupID: mCommon.AppRole, GroupName: mCommon.AppRoleName,
			GroupIndex: 2, OwnerID: ownerID, IsDefault: true},

		// set
		&group{ObjectID: common.BKInnerObjIDSet, GroupID: mCommon.BaseInfo, GroupName: mCommon.BaseInfoName,
			GroupIndex: 1, OwnerID: ownerID, IsDefault: true},

		// module
		&group{ObjectID: common.BKInnerObjIDModule, GroupID: mCommon.BaseInfo, GroupName: mCommon.BaseInfoName,
			GroupIndex: 1, OwnerID: ownerID, IsDefault: true},

		// host
		&group{ObjectID: common.BKInnerObjIDHost, GroupID: mCommon.BaseInfo, GroupName: mCommon.BaseInfoName,
			GroupIndex: 1, OwnerID: ownerID, IsDefault: true},
		&group{ObjectID: common.BKInnerObjIDHost, GroupID: mCommon.HostAutoFields,
			GroupName: mCommon.HostAutoFieldsName, GroupIndex: 3, OwnerID: ownerID, IsDefault: true},

		// proc
		&group{ObjectID: common.BKInnerObjIDProc, GroupID: mCommon.BaseInfo, GroupName: mCommon.BaseInfoName,
			GroupIndex: 1, OwnerID: ownerID, IsDefault: true},
		&group{ObjectID: common.BKInnerObjIDProc, GroupID: mCommon.ProcPort, GroupName: mCommon.ProcPortName,
			GroupIndex: 2, OwnerID: ownerID, IsDefault: true},
		&group{ObjectID: common.BKInnerObjIDProc, GroupID: mCommon.ProcGsekitBaseInfo,
			GroupName: mCommon.ProcGsekitBaseInfoName, GroupIndex: 3, OwnerID: ownerID, IsDefault: true},
		&group{ObjectID: common.BKInnerObjIDProc, GroupID: mCommon.ProcGsekitManageInfo,
			GroupName: mCommon.ProcGsekitManageInfoName, GroupIndex: 4, OwnerID: ownerID, IsDefault: true},

		// plat
		&group{ObjectID: common.BKInnerObjIDPlat, GroupID: mCommon.BaseInfo, GroupName: mCommon.BaseInfoName,
			GroupIndex: 1, OwnerID: ownerID, IsDefault: true},

		// bk_switch
		&group{ObjectID: common.BKInnerObjIDSwitch, GroupID: mCommon.BaseInfo, GroupName: mCommon.BaseInfoName,
			GroupIndex: 1, OwnerID: ownerID, IsDefault: true},
		&group{ObjectID: common.BKInnerObjIDRouter, GroupID: mCommon.BaseInfo, GroupName: mCommon.BaseInfoName,
			GroupIndex: 1, OwnerID: ownerID, IsDefault: true},
		&group{ObjectID: common.BKInnerObjIDBlance, GroupID: mCommon.BaseInfo, GroupName: mCommon.BaseInfoName,
			GroupIndex: 1, OwnerID: ownerID, IsDefault: true},
		&group{ObjectID: common.BKInnerObjIDFirewall, GroupID: mCommon.BaseInfo,
			GroupName: mCommon.BaseInfoName, GroupIndex: 1, OwnerID: ownerID, IsDefault: true},
	}
	for objID, kv := range objectIDs {
		index := int64(1)
		for id, name := range kv {
			row := &group{ObjectID: objID, GroupID: id, GroupName: name, GroupIndex: index, OwnerID: ownerID,
				IsDefault: true}
			dataRows = append(dataRows, row)
			index++
		}

	}

	return dataRows

}

var classificationRows = []*metadata.Classification{
	&metadata.Classification{ClassificationID: "bk_host_manage", ClassificationName: "主机管理",
		ClassificationType: "inner", ClassificationIcon: "icon-cc-host"},
	&metadata.Classification{ClassificationID: "bk_biz_topo", ClassificationName: "业务拓扑",
		ClassificationType: "inner", ClassificationIcon: "icon-cc-business"},
	&metadata.Classification{ClassificationID: "bk_organization", ClassificationName: "组织架构",
		ClassificationType: "inner", ClassificationIcon: "icon-cc-organization"},
	&metadata.Classification{ClassificationID: "bk_network", ClassificationName: "网络", ClassificationType: "inner",
		ClassificationIcon: "icon-cc-network-equipment"},
}

// Group group metadata definition
type group struct {
	BizID      int64  `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id"`
	ID         int64  `field:"id" json:"id" bson:"id"`
	GroupID    string `field:"bk_group_id" json:"bk_group_id" bson:"bk_group_id"`
	GroupName  string `field:"bk_group_name" json:"bk_group_name" bson:"bk_group_name"`
	GroupIndex int64  `field:"bk_group_index" json:"bk_group_index" bson:"bk_group_index"`
	ObjectID   string `field:"bk_obj_id" json:"bk_obj_id" bson:"bk_obj_id"`
	OwnerID    string `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`
	IsDefault  bool   `field:"bk_isdefault" json:"bk_isdefault" bson:"bk_isdefault"`
	IsPre      bool   `field:"ispre" json:"ispre" bson:"ispre"`
	IsCollapse bool   `field:"is_collapse" json:"is_collapse" bson:"is_collapse"`
}

// Object object metadata definition
type Object struct {
	ID         int64  `field:"id" json:"id" bson:"id" mapstructure:"id"`
	ObjCls     string `field:"bk_classification_id" json:"bk_classification_id" bson:"bk_classification_id" mapstructure:"bk_classification_id"`
	ObjIcon    string `field:"bk_obj_icon" json:"bk_obj_icon" bson:"bk_obj_icon" mapstructure:"bk_obj_icon"`
	ObjectID   string `field:"bk_obj_id" json:"bk_obj_id" bson:"bk_obj_id" mapstructure:"bk_obj_id"`
	ObjectName string `field:"bk_obj_name" json:"bk_obj_name" bson:"bk_obj_name" mapstructure:"bk_obj_name"`

	// IsHidden front-end don't display the object if IsHidden is true
	IsHidden bool `field:"bk_ishidden" json:"bk_ishidden" bson:"bk_ishidden" mapstructure:"bk_ishidden"`

	IsPre         bool           `field:"ispre" json:"ispre" bson:"ispre" mapstructure:"ispre"`
	IsPaused      bool           `field:"bk_ispaused" json:"bk_ispaused" bson:"bk_ispaused" mapstructure:"bk_ispaused"`
	Position      string         `field:"position" json:"position" bson:"position" mapstructure:"position"`
	OwnerID       string         `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account" mapstructure:"bk_supplier_account"`
	Description   string         `field:"description" json:"description" bson:"description" mapstructure:"description"`
	Creator       string         `field:"creator" json:"creator" bson:"creator" mapstructure:"creator"`
	Modifier      string         `field:"modifier" json:"modifier" bson:"modifier" mapstructure:"modifier"`
	CreateTime    *metadata.Time `field:"create_time" json:"create_time" bson:"create_time" mapstructure:"create_time"`
	LastTime      *metadata.Time `field:"last_time" json:"last_time" bson:"last_time" mapstructure:"last_time"`
	ObjSortNumber int64          `field:"obj_sort_number" json:"obj_sort_number" bson:"obj_sort_number" mapstructure:"obj_sort_number"`
}
