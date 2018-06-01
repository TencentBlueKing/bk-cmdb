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
 
package api

import (
	"time"
)

// ObjectDes 模型
type ObjectDes struct {
	ID               int        `json:"id"`
	ClassificationID string     `json:"bk_classification_id"`
	ObjID            string     `json:"bk_obj_id"`
	ObjName          string     `json:"bk_obj_name"`
	IsPre            bool       `json:"ispre"`
	IsPaused         bool       `json:"bk_ispaused"`
	OwnerID          string     `json:"bk_supplier_account"`
	Description      string     `json:"description"`
	Creator          string     `json:"creator"`
	Modifier         string     `json:"modifier"`
	CreateTime       *time.Time `json:"create_time"`
	LastTime         *time.Time `json:"last_time"`
}

// ObjectAttDes 模型属性
type ObjectAttDes struct {
	ID            int        `json:"id"`
	OwnerID       string     `json:"bk_supplier_account"`
	ObjID         string     `json:"bk_obj_id"`
	PropertyID    string     `json:"bk_property_id"`
	PropertyName  string     `json:"bk_property_name"`
	AssociationID string     `json:"bk_association_id"`
	PropertyGroup string     `json:"bk_property_group"`
	Editable      bool       `json:"editable"`
	IsRequired    bool       `json:"isrequired"`
	IsReadOnly    bool       `json:"isreadonly"`
	IsOnly        bool       `json:"isonly"`
	PropertyType  string     `json:"bk_property_type"`
	Option        string     `json:"option"`
	Description   string     `json:"description"`
	Creator       string     `json:"creator"`
	CreateTime    *time.Time `json:"create_time"`
	LastTime      *time.Time `json:"last_time"`
}

// ObjectClsDes 模型分类
type ObjectClsDes struct {
	ID      int    `json:"id"`
	ClsID   string `json:"bk_classification_id"`
	ClsName string `json:"bk_classification_name"`
	ClsType string `json:"bk_classification_type"`
	ClsIcon string `json:"bk_classification_icon"`
}

// ObjectAsstDes 模型关联关系
type ObjectAsstDes struct {
	ID          int    `json:"id"`
	ObjectID    string `json:"bk_obj_id"`
	ObjectAttID string `json:"bk_object_att_id"`
	OwnerID     string `json:"bk_supplier_account"`
	AsstObjID   string `json:"bk_asst_obj_id"`
}
