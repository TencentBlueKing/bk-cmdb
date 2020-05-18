/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under,
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package y3_8_202005121212

import (
	"time"

	"configcenter/src/common"
	"configcenter/src/common/metadata"
	mCommon "configcenter/src/scene_server/admin_server/common"
)

// Attribute for object
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

// PodRow return pod attr description data
func PodRow() []*Attribute {
	// TODO: labels, annotations, containers, volumes
	podRows := []*Attribute{
		// BaseInfo class
		&Attribute{
			ObjectID:      common.BKInnerObjIDPod,
			PropertyID:    common.BKAppIDField,
			PropertyName:  "业务ID",
			IsRequired:    true,
			IsOnly:        false,
			IsEditable:    false,
			PropertyGroup: mCommon.BaseInfo,
			PropertyType:  common.FieldTypeInt,
			Option:        metadata.IntOption{},
		},
		&Attribute{
			ObjectID:      common.BKInnerObjIDPod,
			PropertyID:    common.BKModuleIDField,
			PropertyName:  "模块ID",
			IsRequired:    true,
			IsOnly:        false,
			IsEditable:    true,
			PropertyGroup: mCommon.BaseInfo,
			PropertyType:  common.FieldTypeInt,
			Option:        metadata.IntOption{},
		},
		&Attribute{
			ObjectID:      common.BKInnerObjIDPod,
			PropertyID:    common.BKCloudIDField,
			PropertyName:  "云区域ID",
			IsRequired:    true,
			IsOnly:        false,
			IsEditable:    true,
			PropertyGroup: mCommon.BaseInfo,
			PropertyType:  common.FieldTypeInt,
			Option:        metadata.IntOption{},
		},
		&Attribute{
			ObjectID:      common.BKInnerObjIDPod,
			PropertyID:    common.BKHostInnerIPField,
			PropertyName:  "主机内网IP",
			IsRequired:    true,
			IsOnly:        false,
			IsEditable:    true,
			PropertyGroup: mCommon.BaseInfo,
			PropertyType:  common.FieldTypeSingleChar,
			Option:        common.PatternIP,
		},
		&Attribute{
			ObjectID:      common.BKInnerObjIDPod,
			PropertyID:    common.BKPodNameField,
			PropertyName:  "Pod名称",
			IsRequired:    true,
			IsOnly:        false,
			IsEditable:    false,
			PropertyGroup: mCommon.BaseInfo,
			PropertyType:  common.FieldTypeSingleChar,
			Option:        "",
		},
		&Attribute{
			ObjectID:      common.BKInnerObjIDPod,
			PropertyID:    common.BKPodNamespaceField,
			PropertyName:  "命名空间",
			IsRequired:    true,
			IsOnly:        false,
			IsEditable:    false,
			PropertyGroup: mCommon.BaseInfo,
			PropertyType:  common.FieldTypeSingleChar,
			Option:        "",
		},
		&Attribute{
			ObjectID:      common.BKInnerObjIDPod,
			PropertyID:    common.BKPodClusterField,
			PropertyName:  "集群ID",
			IsRequired:    true,
			IsOnly:        false,
			IsEditable:    false,
			PropertyGroup: mCommon.BaseInfo,
			PropertyType:  common.FieldTypeSingleChar,
			Option:        "",
		},
		&Attribute{
			ObjectID:      common.BKInnerObjIDPod,
			PropertyID:    common.BKPodWorkloadType,
			PropertyName:  "工作负载类型",
			IsRequired:    true,
			IsOnly:        false,
			IsEditable:    false,
			PropertyGroup: mCommon.BaseInfo,
			PropertyType:  common.FieldTypeSingleChar,
			Option:        "",
		},
		&Attribute{
			ObjectID:      common.BKInnerObjIDPod,
			PropertyID:    common.BKPodWorkloadName,
			PropertyName:  "工作负载名称",
			IsRequired:    true,
			IsOnly:        false,
			IsEditable:    false,
			PropertyGroup: mCommon.BaseInfo,
			PropertyType:  common.FieldTypeSingleChar,
			Option:        "",
		},
		&Attribute{
			ObjectID:      common.BKInnerObjIDPod,
			PropertyID:    "bk_pod_labels",
			PropertyName:  "Pod标签",
			IsRequired:    false,
			IsOnly:        false,
			IsEditable:    true,
			PropertyGroup: mCommon.BaseInfo,
			PropertyType:  common.FieldTypeLongChar,
			Option:        common.PatternPodLabels,
		},
		&Attribute{
			ObjectID:      common.BKInnerObjIDPod,
			PropertyID:    "bk_pod_annotations",
			PropertyName:  "Pod注解",
			IsRequired:    false,
			IsOnly:        false,
			IsEditable:    true,
			PropertyGroup: mCommon.BaseInfo,
			PropertyType:  common.FieldTypeLongChar,
			Option:        common.PatternPodAnnotations,
		},
		&Attribute{
			ObjectID:      common.BKInnerObjIDPod,
			PropertyID:    "bk_pod_ip",
			PropertyName:  "Pod IP",
			IsRequired:    false,
			IsOnly:        false,
			IsEditable:    true,
			PropertyGroup: mCommon.BaseInfo,
			PropertyType:  common.FieldTypeSingleChar,
			Option:        common.PatternIP,
		},
		&Attribute{
			ObjectID:      common.BKInnerObjIDPod,
			PropertyID:    common.BKPodUUIDField,
			PropertyName:  "UUID",
			IsRequired:    true,
			IsOnly:        false,
			IsEditable:    false,
			PropertyGroup: mCommon.BaseInfo,
			PropertyType:  common.FieldTypeSingleChar,
			Option:        "",
		},
		&Attribute{
			ObjectID:      common.BKInnerObjIDPod,
			PropertyID:    "bk_pod_clustertype",
			PropertyName:  "集群类型",
			IsRequired:    false,
			IsOnly:        false,
			IsEditable:    false,
			PropertyGroup: mCommon.BaseInfo,
			PropertyType:  common.FieldTypeSingleChar,
			Option:        "",
		},
		&Attribute{
			ObjectID:      common.BKInnerObjIDPod,
			PropertyID:    "bk_pod_networkmode",
			PropertyName:  "网络模式",
			IsRequired:    false,
			IsOnly:        false,
			IsEditable:    true,
			PropertyGroup: mCommon.BaseInfo,
			PropertyType:  common.FieldTypeSingleChar,
			Option:        "",
		},
		&Attribute{
			ObjectID:      common.BKInnerObjIDPod,
			PropertyID:    "bk_pod_networktype",
			PropertyName:  "网络类型",
			IsRequired:    false,
			IsOnly:        false,
			IsEditable:    true,
			PropertyGroup: mCommon.BaseInfo,
			PropertyType:  common.FieldTypeSingleChar,
			Option:        "",
		},
		&Attribute{
			ObjectID:      common.BKInnerObjIDPod,
			PropertyID:    "bk_pod_containers",
			PropertyName:  "容器数据",
			IsRequired:    false,
			IsOnly:        false,
			IsEditable:    true,
			PropertyGroup: mCommon.BaseInfo,
			PropertyType:  common.FieldTypeLongChar,
			Option:        "",
		},
		&Attribute{
			ObjectID:      common.BKInnerObjIDPod,
			PropertyID:    "bk_pod_volumes",
			PropertyName:  "卷数据",
			IsRequired:    false,
			IsOnly:        false,
			IsEditable:    true,
			PropertyGroup: mCommon.BaseInfo,
			PropertyType:  common.FieldTypeLongChar,
			Option:        "",
		},
		&Attribute{
			ObjectID:      common.BKInnerObjIDPod,
			PropertyID:    "bk_pod_status",
			PropertyName:  "运行状态",
			IsRequired:    false,
			IsOnly:        false,
			IsEditable:    true,
			PropertyGroup: mCommon.BaseInfo,
			PropertyType:  common.FieldTypeSingleChar,
			Option:        "",
		},
		&Attribute{
			ObjectID:      common.BKInnerObjIDPod,
			PropertyID:    "bk_pod_create_time",
			PropertyName:  "创建时间",
			IsRequired:    false,
			IsOnly:        false,
			IsEditable:    false,
			PropertyGroup: mCommon.BaseInfo,
			PropertyType:  common.FieldTypeTime,
			Option:        "",
		},
		&Attribute{
			ObjectID:      common.BKInnerObjIDPod,
			PropertyID:    "bk_pod_start_time",
			PropertyName:  "运行时间",
			IsRequired:    false,
			IsOnly:        false,
			IsEditable:    true,
			PropertyGroup: mCommon.BaseInfo,
			PropertyType:  common.FieldTypeTime,
			Option:        "",
		},
	}
	return podRows
}
