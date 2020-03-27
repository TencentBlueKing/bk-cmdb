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

package x19_05_16_01

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

func updateProcessBindIPProperty(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	type BindOption struct {
		Name      string `json:"name" bson:"name"`
		Type      string `json:"type" bson:"type"`
		IsDefault bool   `json:"is_default" bson:"is_default"`
		ID        string `json:"id" bson:"id"`
	}
	type Attribute struct {
		ID                int64      `field:"id" json:"id" bson:"id"`
		OwnerID           string     `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`
		ObjectID          string     `field:"bk_obj_id" json:"bk_obj_id" bson:"bk_obj_id"`
		PropertyID        string     `field:"bk_property_id" json:"bk_property_id" bson:"bk_property_id"`
		PropertyName      string     `field:"bk_property_name" json:"bk_property_name" bson:"bk_property_name"`
		PropertyGroup     string     `field:"bk_property_group" json:"bk_property_group" bson:"bk_property_group"`
		PropertyGroupName string     `field:"bk_property_group_name,ignoretomap" json:"bk_property_group_name" bson:"-"`
		PropertyIndex     int64      `field:"bk_property_index" json:"bk_property_index" bson:"bk_property_index"`
		Unit              string     `field:"unit" json:"unit" bson:"unit"`
		Placeholder       string     `field:"placeholder" json:"placeholder" bson:"placeholder"`
		IsEditable        bool       `field:"editable" json:"editable" bson:"editable"`
		IsPre             bool       `field:"ispre" json:"ispre" bson:"ispre"`
		IsRequired        bool       `field:"isrequired" json:"isrequired" bson:"isrequired"`
		IsReadOnly        bool       `field:"isreadonly" json:"isreadonly" bson:"isreadonly"`
		IsOnly            bool       `field:"isonly" json:"isonly" bson:"isonly"`
		IsSystem          bool       `field:"bk_issystem" json:"bk_issystem" bson:"bk_issystem"`
		IsAPI             bool       `field:"bk_isapi" json:"bk_isapi" bson:"bk_isapi"`
		PropertyType      string     `field:"bk_property_type" json:"bk_property_type" bson:"bk_property_type"`
		Option            string     `field:"option" json:"option" bson:"option"`
		Description       string     `field:"description" json:"description" bson:"description"`
		Creator           string     `field:"creator" json:"creator" bson:"creator"`
		CreateTime        *time.Time `json:"create_time" bson:"create_time"`
		LastTime          *time.Time `json:"last_time" bson:"last_time"`
	}

	var bindOptions = []BindOption{{
		Name:      "127.0.0.1",
		Type:      "text",
		IsDefault: true,
		ID:        "1",
	}, {
		Name:      "0.0.0.0",
		Type:      "text",
		IsDefault: false,
		ID:        "2",
	}, {
		Name:      "第一内网IP",
		Type:      "text",
		IsDefault: false,
		ID:        "3",
	}, {
		Name:      "第一外网IP",
		Type:      "text",
		IsDefault: false,
		ID:        "4",
	}}
	blog.InfoJSON("bindOptions removed from database, options: %s", bindOptions)
	now := time.Now()
	var bindIPProperty = Attribute{
		PropertyIndex: 0,
		IsEditable:    true,
		PropertyType:  "singlechar",
		Creator:       conf.User,
		ID:            56,
		PropertyName:  "绑定IP",
		IsSystem:      false,
		IsAPI:         false,
		PropertyID:    "bind_ip",
		ObjectID:      "process",
		PropertyGroup: "proc_port",
		Unit:          "",
		Placeholder:   "",
		IsPre:         true,
		IsRequired:    false,
		IsReadOnly:    false,
		OwnerID:       "0",
		LastTime:      &now,
		CreateTime:    &now,
		Option:        `^([0-9]{1,3}\.){3}[0-9]{1,3}$`,
	}

	uniqueFields := []string{common.BKObjIDField, common.BKPropertyIDField, common.BKOwnerIDField}
	_, _, err := upgrader.Upsert(ctx, db, common.BKTableNameObjAttDes, bindIPProperty, "id", uniqueFields, []string{})
	if nil != err {
		blog.Errorf("[upgrade v19.05.16.01] updateProcessBindIPProperty bind_ip failed, err: %+v", err)
		return err
	}

	return nil
}

func updateProcessNameProperty(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	type Attribute struct {
		ID                int64       `field:"id" json:"id" bson:"id"`
		OwnerID           string      `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`
		ObjectID          string      `field:"bk_obj_id" json:"bk_obj_id" bson:"bk_obj_id"`
		PropertyID        string      `field:"bk_property_id" json:"bk_property_id" bson:"bk_property_id"`
		PropertyName      string      `field:"bk_property_name" json:"bk_property_name" bson:"bk_property_name"`
		PropertyGroup     string      `field:"bk_property_group" json:"bk_property_group" bson:"bk_property_group"`
		PropertyGroupName string      `field:"bk_property_group_name,ignoretomap" json:"bk_property_group_name" bson:"-"`
		PropertyIndex     int64       `field:"bk_property_index" json:"bk_property_index" bson:"bk_property_index"`
		Unit              string      `field:"unit" json:"unit" bson:"unit"`
		Placeholder       string      `field:"placeholder" json:"placeholder" bson:"placeholder"`
		IsEditable        bool        `field:"editable" json:"editable" bson:"editable"`
		IsPre             bool        `field:"ispre" json:"ispre" bson:"ispre"`
		IsRequired        bool        `field:"isrequired" json:"isrequired" bson:"isrequired"`
		IsReadOnly        bool        `field:"isreadonly" json:"isreadonly" bson:"isreadonly"`
		IsOnly            bool        `field:"isonly" json:"isonly" bson:"isonly"`
		IsSystem          bool        `field:"bk_issystem" json:"bk_issystem" bson:"bk_issystem"`
		IsAPI             bool        `field:"bk_isapi" json:"bk_isapi" bson:"bk_isapi"`
		PropertyType      string      `field:"bk_property_type" json:"bk_property_type" bson:"bk_property_type"`
		Option            interface{} `field:"option" json:"option" bson:"option"`
		Description       string      `field:"description" json:"description" bson:"description"`
		Creator           string      `field:"creator" json:"creator" bson:"creator"`
		CreateTime        *time.Time  `json:"create_time" bson:"create_time"`
		LastTime          *time.Time  `json:"last_time" bson:"last_time"`
	}

	now := time.Now()
	var ProcNameProperty = Attribute{
		Placeholder:   "对外显示的服务名</br> 比如程序的二进制名称为java的服务zookeeper，则填zookeeper",
		IsSystem:      false,
		Creator:       conf.User,
		LastTime:      &now,
		ID:            54,
		IsPre:         true,
		ObjectID:      "process",
		PropertyGroup: "default",
		IsReadOnly:    false,
		IsAPI:         false,
		PropertyType:  "singlechar",
		OwnerID:       "0",
		PropertyName:  "进程别名",
		PropertyIndex: -3,
		Unit:          "",
		IsEditable:    true,
		IsRequired:    true,
		Option:        "",
		PropertyID:    "bk_process_name",
		CreateTime:    &now,
	}

	uniqueFields := []string{common.BKObjIDField, common.BKPropertyIDField, common.BKOwnerIDField}
	_, _, err := upgrader.Upsert(ctx, db, common.BKTableNameObjAttDes, ProcNameProperty, "id", uniqueFields, []string{})
	if nil != err {
		blog.Errorf("[upgrade v19.05.16.01] updateProcessNameProperty bind_ip failed, err: %+v", err)
		return err
	}

	return nil
}

func updateAutoTimeGapProperty(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	type Attribute struct {
		ID                int64       `field:"id" json:"id" bson:"id"`
		OwnerID           string      `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`
		ObjectID          string      `field:"bk_obj_id" json:"bk_obj_id" bson:"bk_obj_id"`
		PropertyID        string      `field:"bk_property_id" json:"bk_property_id" bson:"bk_property_id"`
		PropertyName      string      `field:"bk_property_name" json:"bk_property_name" bson:"bk_property_name"`
		PropertyGroup     string      `field:"bk_property_group" json:"bk_property_group" bson:"bk_property_group"`
		PropertyGroupName string      `field:"bk_property_group_name,ignoretomap" json:"bk_property_group_name" bson:"-"`
		PropertyIndex     int64       `field:"bk_property_index" json:"bk_property_index" bson:"bk_property_index"`
		Unit              string      `field:"unit" json:"unit" bson:"unit"`
		Placeholder       string      `field:"placeholder" json:"placeholder" bson:"placeholder"`
		IsEditable        bool        `field:"editable" json:"editable" bson:"editable"`
		IsPre             bool        `field:"ispre" json:"ispre" bson:"ispre"`
		IsRequired        bool        `field:"isrequired" json:"isrequired" bson:"isrequired"`
		IsReadOnly        bool        `field:"isreadonly" json:"isreadonly" bson:"isreadonly"`
		IsOnly            bool        `field:"isonly" json:"isonly" bson:"isonly"`
		IsSystem          bool        `field:"bk_issystem" json:"bk_issystem" bson:"bk_issystem"`
		IsAPI             bool        `field:"bk_isapi" json:"bk_isapi" bson:"bk_isapi"`
		PropertyType      string      `field:"bk_property_type" json:"bk_property_type" bson:"bk_property_type"`
		Option            interface{} `field:"option" json:"option" bson:"option"`
		Description       string      `field:"description" json:"description" bson:"description"`
		Creator           string      `field:"creator" json:"creator" bson:"creator"`
		CreateTime        *time.Time  `json:"create_time" bson:"create_time"`
		LastTime          *time.Time  `json:"last_time" bson:"last_time"`
	}

	now := time.Now()
	var property = Attribute{
		IsPre:         true,
		IsRequired:    false,
		PropertyID:    "auto_time_gap",
		PropertyName:  "拉起间隔",
		Unit:          "",
		IsAPI:         false,
		IsReadOnly:    false,
		IsSystem:      false,
		PropertyType:  "int",
		ID:            73,
		OwnerID:       "0",
		ObjectID:      "process",
		PropertyIndex: 0,
		IsEditable:    true,
		PropertyGroup: "none",
		Placeholder:   "",
		Option: map[string]int{
			"min": 1,
			"max": 10000,
		},
		Creator:    conf.User,
		CreateTime: &now,
		LastTime:   &now,
	}

	uniqueFields := []string{common.BKObjIDField, common.BKPropertyIDField, common.BKOwnerIDField}
	_, _, err := upgrader.Upsert(ctx, db, common.BKTableNameObjAttDes, property, "id", uniqueFields, []string{})
	if nil != err {
		blog.Errorf("[upgrade v19.05.16.01] updateAutoTimeGapProperty bind_ip failed, err: %+v", err)
		return err
	}

	return nil
}

func updateProcNumProperty(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	type Attribute struct {
		ID                int64       `field:"id" json:"id" bson:"id"`
		OwnerID           string      `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`
		ObjectID          string      `field:"bk_obj_id" json:"bk_obj_id" bson:"bk_obj_id"`
		PropertyID        string      `field:"bk_property_id" json:"bk_property_id" bson:"bk_property_id"`
		PropertyName      string      `field:"bk_property_name" json:"bk_property_name" bson:"bk_property_name"`
		PropertyGroup     string      `field:"bk_property_group" json:"bk_property_group" bson:"bk_property_group"`
		PropertyGroupName string      `field:"bk_property_group_name,ignoretomap" json:"bk_property_group_name" bson:"-"`
		PropertyIndex     int64       `field:"bk_property_index" json:"bk_property_index" bson:"bk_property_index"`
		Unit              string      `field:"unit" json:"unit" bson:"unit"`
		Placeholder       string      `field:"placeholder" json:"placeholder" bson:"placeholder"`
		IsEditable        bool        `field:"editable" json:"editable" bson:"editable"`
		IsPre             bool        `field:"ispre" json:"ispre" bson:"ispre"`
		IsRequired        bool        `field:"isrequired" json:"isrequired" bson:"isrequired"`
		IsReadOnly        bool        `field:"isreadonly" json:"isreadonly" bson:"isreadonly"`
		IsOnly            bool        `field:"isonly" json:"isonly" bson:"isonly"`
		IsSystem          bool        `field:"bk_issystem" json:"bk_issystem" bson:"bk_issystem"`
		IsAPI             bool        `field:"bk_isapi" json:"bk_isapi" bson:"bk_isapi"`
		PropertyType      string      `field:"bk_property_type" json:"bk_property_type" bson:"bk_property_type"`
		Option            interface{} `field:"option" json:"option" bson:"option"`
		Description       string      `field:"description" json:"description" bson:"description"`
		Creator           string      `field:"creator" json:"creator" bson:"creator"`
		CreateTime        *time.Time  `json:"create_time" bson:"create_time"`
		LastTime          *time.Time  `json:"last_time" bson:"last_time"`
	}

	now := time.Now()
	var property = Attribute{
		PropertyName:  "启动数量",
		PropertyGroup: "none",
		IsRequired:    false,
		Option: map[string]int{
			"min": 1,
			"max": 10000,
		},
		ID:            63,
		PropertyID:    "proc_num",
		Unit:          "",
		Placeholder:   "",
		IsReadOnly:    false,
		IsSystem:      false,
		IsPre:         true,
		CreateTime:    &now,
		LastTime:      &now,
		Creator:       conf.User,
		OwnerID:       "0",
		ObjectID:      "process",
		PropertyIndex: 0,
		IsEditable:    true,
		IsAPI:         false,
		PropertyType:  "int",
	}

	uniqueFields := []string{common.BKObjIDField, common.BKPropertyIDField, common.BKOwnerIDField}
	_, _, err := upgrader.Upsert(ctx, db, common.BKTableNameObjAttDes, property, "id", uniqueFields, []string{})
	if nil != err {
		blog.Errorf("[upgrade v19.05.16.01] updateProcNumProperty bind_ip failed, err: %+v", err)
		return err
	}

	return nil
}

func updatePriorityProperty(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	type Attribute struct {
		ID                int64       `field:"id" json:"id" bson:"id"`
		OwnerID           string      `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`
		ObjectID          string      `field:"bk_obj_id" json:"bk_obj_id" bson:"bk_obj_id"`
		PropertyID        string      `field:"bk_property_id" json:"bk_property_id" bson:"bk_property_id"`
		PropertyName      string      `field:"bk_property_name" json:"bk_property_name" bson:"bk_property_name"`
		PropertyGroup     string      `field:"bk_property_group" json:"bk_property_group" bson:"bk_property_group"`
		PropertyGroupName string      `field:"bk_property_group_name,ignoretomap" json:"bk_property_group_name" bson:"-"`
		PropertyIndex     int64       `field:"bk_property_index" json:"bk_property_index" bson:"bk_property_index"`
		Unit              string      `field:"unit" json:"unit" bson:"unit"`
		Placeholder       string      `field:"placeholder" json:"placeholder" bson:"placeholder"`
		IsEditable        bool        `field:"editable" json:"editable" bson:"editable"`
		IsPre             bool        `field:"ispre" json:"ispre" bson:"ispre"`
		IsRequired        bool        `field:"isrequired" json:"isrequired" bson:"isrequired"`
		IsReadOnly        bool        `field:"isreadonly" json:"isreadonly" bson:"isreadonly"`
		IsOnly            bool        `field:"isonly" json:"isonly" bson:"isonly"`
		IsSystem          bool        `field:"bk_issystem" json:"bk_issystem" bson:"bk_issystem"`
		IsAPI             bool        `field:"bk_isapi" json:"bk_isapi" bson:"bk_isapi"`
		PropertyType      string      `field:"bk_property_type" json:"bk_property_type" bson:"bk_property_type"`
		Option            interface{} `field:"option" json:"option" bson:"option"`
		Description       string      `field:"description" json:"description" bson:"description"`
		Creator           string      `field:"creator" json:"creator" bson:"creator"`
		CreateTime        *time.Time  `json:"create_time" bson:"create_time"`
		LastTime          *time.Time  `json:"last_time" bson:"last_time"`
	}

	now := time.Now()
	var property = Attribute{
		PropertyGroup: "none",
		Unit:          "",
		IsSystem:      false,
		Creator:       conf.User,
		ID:            64,
		OwnerID:       "0",
		PropertyIndex: 0,
		CreateTime:    &now,
		PropertyID:    "priority",
		Placeholder:   "",
		IsEditable:    true,
		IsReadOnly:    false,
		IsAPI:         false,
		PropertyType:  "int",
		Option: map[string]int{
			"min": 1,
			"max": 10000,
		},
		LastTime:     &now,
		ObjectID:     "process",
		PropertyName: "启动优先级",
		IsPre:        true,
		IsRequired:   false,
	}

	uniqueFields := []string{common.BKObjIDField, common.BKPropertyIDField, common.BKOwnerIDField}
	_, _, err := upgrader.Upsert(ctx, db, common.BKTableNameObjAttDes, property, "id", uniqueFields, []string{})
	if nil != err {
		blog.Errorf("[upgrade v19.05.16.01] updatePriorityProperty bind_ip failed, err: %+v", err)
		return err
	}

	return nil
}

func updateTimeoutProperty(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	type Attribute struct {
		ID                int64       `field:"id" json:"id" bson:"id"`
		OwnerID           string      `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`
		ObjectID          string      `field:"bk_obj_id" json:"bk_obj_id" bson:"bk_obj_id"`
		PropertyID        string      `field:"bk_property_id" json:"bk_property_id" bson:"bk_property_id"`
		PropertyName      string      `field:"bk_property_name" json:"bk_property_name" bson:"bk_property_name"`
		PropertyGroup     string      `field:"bk_property_group" json:"bk_property_group" bson:"bk_property_group"`
		PropertyGroupName string      `field:"bk_property_group_name,ignoretomap" json:"bk_property_group_name" bson:"-"`
		PropertyIndex     int64       `field:"bk_property_index" json:"bk_property_index" bson:"bk_property_index"`
		Unit              string      `field:"unit" json:"unit" bson:"unit"`
		Placeholder       string      `field:"placeholder" json:"placeholder" bson:"placeholder"`
		IsEditable        bool        `field:"editable" json:"editable" bson:"editable"`
		IsPre             bool        `field:"ispre" json:"ispre" bson:"ispre"`
		IsRequired        bool        `field:"isrequired" json:"isrequired" bson:"isrequired"`
		IsReadOnly        bool        `field:"isreadonly" json:"isreadonly" bson:"isreadonly"`
		IsOnly            bool        `field:"isonly" json:"isonly" bson:"isonly"`
		IsSystem          bool        `field:"bk_issystem" json:"bk_issystem" bson:"bk_issystem"`
		IsAPI             bool        `field:"bk_isapi" json:"bk_isapi" bson:"bk_isapi"`
		PropertyType      string      `field:"bk_property_type" json:"bk_property_type" bson:"bk_property_type"`
		Option            interface{} `field:"option" json:"option" bson:"option"`
		Description       string      `field:"description" json:"description" bson:"description"`
		Creator           string      `field:"creator" json:"creator" bson:"creator"`
		CreateTime        *time.Time  `json:"create_time" bson:"create_time"`
		LastTime          *time.Time  `json:"last_time" bson:"last_time"`
	}

	now := time.Now()
	var property = Attribute{
		Creator:       conf.User,
		PropertyIndex: 0,
		IsReadOnly:    false,
		Unit:          "",
		IsPre:         true,
		IsRequired:    false,
		PropertyType:  "int",
		Option: map[string]int{
			"min": 1,
			"max": 10000,
		},
		CreateTime:    &now,
		OwnerID:       "0",
		PropertyName:  "操作超时时长",
		PropertyGroup: "none",
		IsEditable:    true,
		ID:            65,
		PropertyID:    "timeout",
		IsSystem:      false,
		IsAPI:         false,
		LastTime:      &now,
		ObjectID:      "process",
		Placeholder:   "",
	}

	uniqueFields := []string{common.BKObjIDField, common.BKPropertyIDField, common.BKOwnerIDField}
	_, _, err := upgrader.Upsert(ctx, db, common.BKTableNameObjAttDes, property, "id", uniqueFields, []string{})
	if nil != err {
		blog.Errorf("[upgrade v19.05.16.01] updateTimeoutProperty bind_ip failed, err: %+v", err)
		return err
	}

	return nil
}

func updateProcessNamePropertyIndex(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	processNameFilter := map[string]interface{}{
		common.BKPropertyIDField: "bk_process_name",
	}
	processNameIndex := map[string]interface{}{
		common.BKPropertyIndexField: -2,
	}
	if err := db.Table(common.BKTableNameObjAttDes).Update(ctx, processNameFilter, processNameIndex); err != nil {
		blog.Errorf("[upgrade v19.05.16.01] updatePropertyIndex bk_process_name index failed, err: %+v", err)
		return err
	}
	return nil
}

func updateFuncNamePropertyIndex(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	funcNameFilter := map[string]interface{}{
		common.BKPropertyIDField: "bk_func_name",
	}
	funcNameIndex := map[string]interface{}{
		common.BKPropertyIndexField: -3,
	}
	if err := db.Table(common.BKTableNameObjAttDes).Update(ctx, funcNameFilter, funcNameIndex); err != nil {
		blog.Errorf("[upgrade v19.05.16.01] updateFuncNamePropertyIndex bk_func_name index failed, err: %+v", err)
		return err
	}
	return nil
}

func deleteProcessUnique(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	filter := map[string]interface{}{
		common.BKObjIDField: common.BKInnerObjIDProc,
	}
	if err := db.Table(common.BKTableNameObjUnique).Delete(ctx, filter); err != nil {
		blog.Errorf("[upgrade v19.05.16.01] deleteProcessUnique failed, err: %+v", err)
		return err
	}
	return nil
}

func updateFuncIDProperty(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	type Attribute struct {
		ID                int64       `field:"id" json:"id" bson:"id"`
		OwnerID           string      `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`
		ObjectID          string      `field:"bk_obj_id" json:"bk_obj_id" bson:"bk_obj_id"`
		PropertyID        string      `field:"bk_property_id" json:"bk_property_id" bson:"bk_property_id"`
		PropertyName      string      `field:"bk_property_name" json:"bk_property_name" bson:"bk_property_name"`
		PropertyGroup     string      `field:"bk_property_group" json:"bk_property_group" bson:"bk_property_group"`
		PropertyGroupName string      `field:"bk_property_group_name,ignoretomap" json:"bk_property_group_name" bson:"-"`
		PropertyIndex     int64       `field:"bk_property_index" json:"bk_property_index" bson:"bk_property_index"`
		Unit              string      `field:"unit" json:"unit" bson:"unit"`
		Placeholder       string      `field:"placeholder" json:"placeholder" bson:"placeholder"`
		IsEditable        bool        `field:"editable" json:"editable" bson:"editable"`
		IsPre             bool        `field:"ispre" json:"ispre" bson:"ispre"`
		IsRequired        bool        `field:"isrequired" json:"isrequired" bson:"isrequired"`
		IsReadOnly        bool        `field:"isreadonly" json:"isreadonly" bson:"isreadonly"`
		IsOnly            bool        `field:"isonly" json:"isonly" bson:"isonly"`
		IsSystem          bool        `field:"bk_issystem" json:"bk_issystem" bson:"bk_issystem"`
		IsAPI             bool        `field:"bk_isapi" json:"bk_isapi" bson:"bk_isapi"`
		PropertyType      string      `field:"bk_property_type" json:"bk_property_type" bson:"bk_property_type"`
		Option            interface{} `field:"option" json:"option" bson:"option"`
		Description       string      `field:"description" json:"description" bson:"description"`
		Creator           string      `field:"creator" json:"creator" bson:"creator"`
		CreateTime        *time.Time  `json:"create_time" bson:"create_time"`
		LastTime          *time.Time  `json:"last_time" bson:"last_time"`
	}

	now := time.Now()
	var property = Attribute{
		Option:        "",
		PropertyIndex: 0,
		Unit:          "",
		IsRequired:    false,
		PropertyType:  "singlechar",
		ID:            59,
		IsEditable:    true,
		CreateTime:    &now,
		IsAPI:         false,
		Creator:       conf.User,
		OwnerID:       "0",
		ObjectID:      "process",
		PropertyName:  "功能ID",
		PropertyGroup: "none",
		IsPre:         true,
		PropertyID:    "bk_func_id",
		Placeholder:   "对进程的数字型标注，便于快速检索",
		IsReadOnly:    false,
		IsSystem:      false,
		LastTime:      &now,
	}

	uniqueFields := []string{common.BKObjIDField, common.BKPropertyIDField, common.BKOwnerIDField}
	_, _, err := upgrader.Upsert(ctx, db, common.BKTableNameObjAttDes, property, "id", uniqueFields, []string{})
	if nil != err {
		blog.Errorf("[upgrade v19.05.16.01] updateFuncIDProperty bk_func_id failed, err: %+v", err)
		return err
	}

	return nil
}

func UpdateProcPortPropertyGroupName(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	// update proc_port group name
	row := &metadata.Group{
		ObjectID:   common.BKInnerObjIDProc,
		GroupID:    mCommon.ProcPort,
		GroupName:  mCommon.ProcPortName,
		GroupIndex: 2,
		OwnerID:    conf.OwnerID,
		IsDefault:  true,
	}
	if _, _, err := upgrader.Upsert(ctx, db, common.BKTableNamePropertyGroup, row, "id", []string{common.BKObjIDField, "bk_group_id"}, []string{"id"}); err != nil {
		blog.Errorf("add data for  %s table error  %s", common.BKTableNamePropertyGroup, err)
		return err
	}
	return nil
}
