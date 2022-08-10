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

package types

import (
	"configcenter/src/storage/dal/table"
)

// WorkLoadSpecFieldsDescriptor workLoad spec's fields descriptors.
var WorkLoadSpecFieldsDescriptor = table.FieldsDescriptors{
	{Field: KubeNameField, Type: table.String, IsRequired: true, IsEditable: false},
	{Field: NamespaceField, Type: table.String, IsRequired: true, IsEditable: false},
	{Field: LabelsField, Type: table.MapString, IsRequired: false, IsEditable: true},
	{Field: SelectorField, Type: table.String, IsRequired: false, IsEditable: true},
	{Field: ReplicasField, Type: table.Numeric, IsRequired: true, IsEditable: true},
	{Field: StrategyTypeField, Type: table.String, IsRequired: false, IsEditable: true},
	{Field: MinReadySecondsField, Type: table.Numeric, IsRequired: false, IsEditable: true},
	{Field: RollingUpdateStrategyField, Type: table.Object, IsRequired: false, IsEditable: true},
}
