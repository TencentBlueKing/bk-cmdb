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

// Cluster container cluster table structure
type Cluster struct {
	// ID cluster auto-increment ID in cc
	ID *int64 `json:"id" bson:"id"`
	// BizID the business ID to which the cluster belongs
	BizID *int64 `json:"bk_biz_id" bson:"bk_biz_id"`
	// ClusterFields cluster base fields
	ClusterBaseFields `json:",inline" bson:",inline"`
	// SupplierAccount the supplier account' that this resource belongs to.
	SupplierAccount string `json:"bk_supplier_account" bson:"bk_supplier_account"`
}

// ClusterBaseFields 创建集群请求字段
type ClusterBaseFields struct {
	Name *string `json:"name" bson:"name"`
	// SchedulingEngine scheduling engines, such as k8s, tke, etc.
	SchedulingEngine *string `json:"scheduling_engine" bson:"scheduling_engine"`
	// Uid ID of the cluster itself
	Uid *string `json:"uid" bson:"uid"`
	// Xid The underlying cluster ID it depends on
	Xid *string `json:"xid" bson:"xid"`
	// Version cluster version
	Version *string `json:"version" bson:"version"`
	// NetworkType network type, such as overlay or underlay
	NetworkType *string `json:"network_type" bson:"network_type"`
	// Region the region where the cluster is located
	Region *string `json:"region" bson:"region"`
	// Vpc vpc network
	Vpc *string `json:"vpc" bson:"vpc"`
	// NetWork global routing network address (container overlay network) For example: ["1.1.1.0/21"]
	NetWork *[]string `json:"network" bson:"network"`
	// Type cluster network type, e.g. public clusters, private clusters, etc.
	Type *string `json:"type" bson:"type"`
}
