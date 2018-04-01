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

/** TODO: DELETE
import (
	"time"
)


// APIUser
type APIUser struct {
	ID            int      `gorm:"primary_key;column:ID;auto_increment"`
	UserName      string   `gorm:"type:varchar(32);column:UserName;index:i_UserName"`
	UserPwd       string   `gorm:"type:varchar(32);column:UserPwd;index:i_UserPwd"`
	UserType      int      `gorm:"column:UserType"`
	AllowedIP     string   `gorm:"column:AllowedIP;type:varchar(128)"`
	AllowedMethod string   `gorm:"column:AllowedMethod;type:varchar(512)"`
	AllowedAppID  string   `gorm:"column:AllowedAppID;type:varchar(32)"`
	Description   string   `gorm:"column:Description;type:varchar(256)"`
	Page          BasePage `gorm:"-" json:"Page"`
}

func (cli APIUser) TableName() string {
	return "cc_APIUser"
}

// BaseParaData
type BaseParaData struct {
	ID            int       `gorm:"column:ID;auto_increment;primary_key"`
	DataType      string    `gorm:"type:varchar(50);column:DataType;index:i_DataType"`
	ParameterCode string    `gorm:"type:varchar(50);column:ParameterCode"`
	ParameterName string    `gorm:"type:varchar(50);column:ParameterName"`
	ParentCode    string    `gorm:"type:varchar(50);column:ParentCode"`
	CreateTime    time.Time `gorm:"column:CreateTime"`
	LastTime      time.Time `gorm:"column:LastTime"`
	Page          BasePage `gorm:"-" json:"Page"`
}

func (cli BaseParaData) TableName() string {
	return "cc_BaseParaData"
}

// UserCustom
type UserCustom struct {
	ID            int       `gorm:primary_key;column:ID`
	UserName      string    `gorm: "type:varchar(32);column:UserName"`
	ApplicationID int       `gorm:"column:ApplicationID"`
	CustomType    string    `gorm:"type:varchar(32);column:CustomType;"`
	CustomContent string    `gorm:"column:CustomContent"`
	CreateTime    time.Time `gorm:"column:CreateTime"`
	LastTime      time.Time `gorm:"column:LastTime"`
	Page          BasePage `gorm:"-" json:"Page"`
}

func (cli UserCustom) TableName() string {
	return "cc_UserCustom"
}

*/
