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

// Package host define host types
package host

// HostBaseInfo host base info
type HostBaseInfo struct {
	BKAppIDField         int64    `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id"`
	BKModuleIDField      int64    `field:"bk_module_id" json:"bk_module_id" bson:"bk_module_id"`
	BKSetIDField         int64    `field:"bk_set_id" json:"bk_set_id" bson:"bk_set_id"`
	BKHostIDField        int64    `field:"bk_host_id" json:"bk_host_id" bson:"bk_host_id"`
	BKHostNameField      string   `field:"bk_host_name" json:"bk_host_name" bson:"bk_host_name"`
	BKHostInnerIPField   []string `field:"bk_host_innerip" json:"bk_host_innerip" bson:"bk_host_innerip"`
	BKHostInnerIPv6Field []string `field:"bk_host_innerip_v6" json:"bk_host_innerip_v6" bson:"bk_host_innerip_v6"`
	BKCloudID            int64    `field:"bk_cloud_id" json:"bk_cloud_id" bson:"bk_cloud_id"`
}
