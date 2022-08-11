/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package dal

import (
	"configcenter/src/kube/types"
)

// ClusterFields merge the fields of the cluster and the details corresponding to the fields together.
var ClusterFields = mergeFields(ClusterFieldsDescriptor)

// ClusterFieldsDescriptor cluster's fields descriptors.
var ClusterFieldsDescriptor = mergeFieldDescriptors(
	FieldsDescriptors{
		{Field: types.BKIDField, Type: Numeric, IsRequired: true, IsEditable: false},
		{Field: types.BKBizIDField, Type: Numeric, IsRequired: true, IsEditable: false},
		{Field: types.BKSupplierAccountField, Type: String, IsRequired: true, IsEditable: false},
		{Field: types.CreatorField, Type: String, IsRequired: true, IsEditable: false},
		{Field: types.ModifierField, Type: String, IsRequired: true, IsEditable: true},
		{Field: types.CreateTimeField, Type: Numeric, IsRequired: true, IsEditable: false},
		{Field: types.LastTimeField, Type: Numeric, IsRequired: true, IsEditable: true},
	},
	mergeFieldDescriptors(ClusterSpecFieldsDescriptor),
)

// ClusterSpecFieldsDescriptor cluster spec's fields descriptors.
var ClusterSpecFieldsDescriptor = FieldsDescriptors{
	{Field: types.KubeNameField, Type: String, IsRequired: true, IsEditable: false},
	{Field: types.SchedulingEngineField, Type: String, IsRequired: false, IsEditable: false},
	{Field: types.UidField, Type: String, IsRequired: true, IsEditable: false},
	{Field: types.XidField, Type: String, IsRequired: false, IsEditable: false},
	{Field: types.VersionField, Type: String, IsRequired: false, IsEditable: true},
	{Field: types.NetworkTypeField, Type: Enum, IsRequired: false, IsEditable: true},
	{Field: types.RegionField, Type: String, IsRequired: false, IsEditable: true},
	{Field: types.VpcField, Type: String, IsRequired: false, IsEditable: false},
	{Field: types.NetworkField, Type: String, IsRequired: false, IsEditable: false},
	{Field: types.TypeField, Type: String, IsRequired: false, IsEditable: true},
}
