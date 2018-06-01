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
 
package metadata

/** TODO:delete
import (
	"time"
)

// HostDes  host metata
type HostDes struct {
	ID              int       `gorm:"column:ID;auto_increment;primary_key" json:"id"`
	ObjClass        int       `gorm:"column:ObjCls" json:"objcls"`
	IsPre           int       `gorm:"column:IsPre;index:i_IsPre" json:"ispre"`
	OwnerID         string    `gorm:"column:OwnerID;type:varchar(32);index:i_OwnerID" json:"ownerid"`
	HostPropertyTag string    `gorm:"column:HostPropertyTag;type:varchar(64)" json:"hostpropertytag"`
	ObjID           string    `gorm:"column:ObjID;type:varchar(32);index:i_ObjID" json:"objid"`
	ObjName         string    `gorm:"column:ObjName;type:varchar(64)" json:"objname"`
	Description     string    `gorm:"column:Description;type:varchar(256)" json:"description"`
	Creator         string    `gorm:"column:Creator;type:varchar(64)" json:"creator"`
	CreateTime      time.Time `gorm:"column:CreateTime" json:"createtime"`
	LastTime        time.Time `gorm:"column:LastTime" json:"lasttime"`
	Page            *BasePage `gorm:"-" json:"Page,omitempty"`
}

func (HostDes) TableName() string {

	return "cc_HostDes"
}

// HostRealtimedata
type HostRealTimeData struct {
	ID            int       `gorm:"column:ID;auto_increment;primary_key" json:"id"`
	HostID        int       `gorm:"column:HostID;index:i_HostID" json:"hostid"`
	GroupName     string    `gorm:"column:GroupName;index:i_GroupName;type:varchar(32)" json:"groupname"`
	PropertyKey   string    `gorm:"column:PropertyKey;index:i_PropertyKey;type:varchar(32)" json:"propertykey"`
	PropertyName  string    `gorm:"column:PropertyName;type:varchar(64)" json:"propertyname"`
	PropertyValue string    `gorm:"column:PropertyValue;type:varchar(64);index:i_PropertyValue" json:"propertyvalue"`
	Page          *BasePage `gorm:"-" json:"Page,omitempty"`
}

func (HostRealTimeData) TableName() string {
	return "cc_HostRealTimeData"
}
*/

// ModuleHostConfig module host map config
type ModuleHostConfig struct {
	ID            int       `bson:"id"        json:"id"`
	HostID        int       `bson:"bk_host_id"   json:"bk_host_id"`
	ModuleID      int       `bson:"bk_module_id" json:"bk_module_id"`
	SetID         int       `bson:"bk_set_id"    json:"bk_set_id"`
	ApplicationID int       `bson:"bk_biz_id"    json:"bk_biz_id"`
	Page          *BasePage `gorm:"-"            json:"page,omitempty"`
}

// TableName return the table name
func (ModuleHostConfig) TableName() string {
	return "cc_ModuleHostConfig"
}
