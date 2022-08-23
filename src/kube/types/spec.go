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

import "errors"

// ClusterSpec describes the common attributes of cluster, it is used by the structure below it.
type ClusterSpec struct {
	// BizID business id in cc
	BizID *int64 `json:"bk_biz_id" bson:"bk_biz_id"`

	// ClusterID cluster id in cc
	ClusterID *int64 `json:"bk_cluster_id" bson:"bk_cluster_id"`

	// ClusterUID cluster id in third party platform
	ClusterUID *string `json:"cluster_uid" bson:"cluster_uid"`
}

// NamespaceSpec describes the common attributes of namespace, it is used by the structure below it.
type NamespaceSpec struct {
	ClusterSpec `json:",inline" bson:",inline"`

	// NamespaceID namespace id in cc
	NamespaceID *int64 `json:"bk_namespace_id" bson:"bk_namespace_id"`

	// Namespace namespace name in third party platform
	Namespace *string `json:"namespace" bson:"namespace"`
}

// Reference store pod-related workload related information
type Reference struct {
	// Kind workload kind
	Kind *string `json:"kind" bson:"kind"`

	// Name workload name
	Name *string `json:"name" bson:"name"`

	// ID workload id in cc
	ID *int64 `json:"id" bson:"id"`
}

// WorkloadSpec describes the common attributes of workload, it is used by the structure below it.
type WorkloadSpec struct {
	NamespaceSpec `json:",inline" bson:",inline`
	Ref           *Reference `json:"ref" bson:"ref"`
}

// PodSpec describes the common attributes of pod, it is used by the structure below it.
type PodSpec struct {
	WorkloadSpec `json:",inline" bson:",inline`

	// NodeID node id in cc
	NodeID *int64 `json:"bk_node_id" bson:"bk_node_id"`

	// Node node name in third party platform
	Node *string `json:"node" bson:"node"`

	// HostID host id in cc
	HostID *int64 `json:"bk_host_id" bson:"bk_host_id"`

	// PodID pod id in cc
	PodID *int64 `json:"bk_pod_id" bson:"bk_pod_id"`

	// Pod pod name in third party platform
	Pod *string `json:"pod" bson:"pod"`
}

// GetKubeSubTopoObject 获取指定资源的下一级拓扑资源对象，需要首先判断是否是
func GetKubeSubTopoObject(object string) string {

	switch object {
	case KubeBusiness:
		return KubeCluster
	case KubeCluster:
		return KubeNamespace
	case KubeNamespace:
		return KubeWorkload
	default:
		return KubePod
	}
}

// GetWorkLoadTables 获取workload子项
func GetWorkLoadTables() []string {

	return []string{
		BKTableNameBaseDeployment,
		BKTableNameGameDeployment,
		BKTableNameBaseJob,
		BKTableNameBaseCronJob,
		BKTableNameGameStatefulSet,
		BKTableNameBaseStatefulSet,
		BKTableNameBaseDaemonSet,
		BKTableNameBasePodWorkload,
	}
}

// IsContainerTopoResource 判断是否是容器拓扑对象
func IsContainerTopoResource(object string) bool {
	switch object {
	case KubeBusiness, KubeCluster, KubeNode, KubeNamespace, KubeWorkload:
		return true
	default:
		return false
	}
}

// GetCollectionWithObject 根据容器对象获取对应的collection
func GetCollectionWithObject(object string) ([]string, error) {
	switch object {
	case KubeCluster:
		return []string{BKTableNameBaseCluster}, nil
	case KubeNamespace:
		return []string{BKTableNameBaseNamespace}, nil
	case KubeNode:
		return []string{BKTableNameBaseNode}, nil
	case KubePod:
		return []string{BKTableNameBasePod}, nil
	case KubeWorkload:
		return GetWorkLoadTables(), nil
	default:
		return []string{}, errors.New("no corresponding table found")
	}
}

// IsKubeResourceKind 判断是否是容器资源对象
func IsKubeResourceKind(object string) bool {
	switch object {
	case KubeBusiness, KubeCluster, KubeNode, KubeNamespace, string(KubeDeployment),
		string(KubeStatefulSet), string(KubeDaemonSet), string(KubeGameStatefulSet), string(KubeGameDeployment),
		string(KubeCronJob), string(KubeJob), string(KubePodWorkload):
		return true
	default:
		return false
	}
}
