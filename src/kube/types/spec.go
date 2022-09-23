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
	"errors"
)

func init() {
	initClusterFieldsType()
	initNodeFieldsType()
	initPodFieldsType()
	initContainerFieldsType()
}

// ClusterSpec describes the common attributes of cluster, it is used by the structure below it.
type ClusterSpec struct {
	// BizID business id in cc
	BizID int64 `json:"bk_biz_id,omitempty" bson:"bk_biz_id"`

	// ClusterID cluster id in cc
	ClusterID int64 `json:"bk_cluster_id,omitempty" bson:"bk_cluster_id"`

	// ClusterUID cluster id in third party platform
	ClusterUID string `json:"cluster_uid,omitempty" bson:"cluster_uid"`
}

// NamespaceSpec describes the common attributes of namespace, it is used by the structure below it.
type NamespaceSpec struct {
	ClusterSpec `json:",inline" bson:",inline"`

	// NamespaceID namespace id in cc
	NamespaceID int64 `json:"bk_namespace_id,omitempty" bson:"bk_namespace_id"`

	// Namespace namespace name in third party platform
	Namespace string `json:"namespace,omitempty" bson:"namespace"`
}

// Reference store pod-related workload related information
type Reference struct {
	// Kind workload kind
	Kind WorkloadType `json:"kind" bson:"kind"`

	// Name workload name
	Name string `json:"name" bson:"name"`

	// ID workload id in cc
	ID int64 `json:"id" bson:"id"`
}

// WorkloadSpec describes the common attributes of workload,
// it is used by the structure below it.
type WorkloadSpec struct {
	NamespaceSpec `json:",inline" bson:",inline"`
	Ref           Reference `json:"ref" bson:"ref"`
}

// CommonSpecFieldsDescriptor public field properties
var CommonSpecFieldsDescriptor = table.FieldsDescriptors{
	{Field: BKIDField, IsRequired: true, IsEditable: false},
	{Field: BKSupplierAccountField, IsRequired: true, IsEditable: false},
	{Field: CreatorField, IsRequired: true, IsEditable: false},
	{Field: ModifierField, IsRequired: true, IsEditable: true},
	{Field: CreateTimeField, IsRequired: true, IsEditable: false},
	{Field: LastTimeField, IsRequired: true, IsEditable: true},
}

// BizIDDescriptor bizID descriptor is taken out separately and not placed in CommonSpecFieldsDescriptor because
// bk_biz_id does not exist in the container table and needs to be processed separately.
var BizIDDescriptor = table.FieldsDescriptors{
	{Field: BKBizIDField, IsRequired: true, IsEditable: false},
}

// IsCommonField judges whether the field is a special field.
// If it belongs to common fields or biz_id field, it needs
// to skip unified verification whether it is creating a scene or updating a scene.
func IsCommonField(field string) bool {
	for _, commonDescriptor := range CommonSpecFieldsDescriptor {
		if commonDescriptor.Field == field {
			return true
		}
	}

	for _, bizDescriptor := range BizIDDescriptor {
		if field == bizDescriptor.Field {
			return true
		}
	}
	return false
}

// GetKubeSubTopoObject get the next-level topology resource object of the specified resource
func GetKubeSubTopoObject(object string, id int64, bizID int64) (string, map[string]interface{}) {

	switch object {
	case KubeBusiness:
		return KubeCluster, map[string]interface{}{
			BKBizIDField: bizID,
		}
	case KubeCluster:
		return KubeNamespace, map[string]interface{}{
			BKClusterIDFiled: id,
		}
	case KubeNamespace:
		return KubeWorkload, map[string]interface{}{
			BKNamespaceIDField: id,
		}
	case KubeFolder:
		return "", map[string]interface{}{}
	default:
		return KubePod, map[string]interface{}{}
	}
}

// GetWorkLoadTables get the table name of the full workload.
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

// IsContainerTopoResource determine whether it is a container object type.
func IsContainerTopoResource(object string) bool {
	switch object {
	case KubeBusiness, KubeCluster, KubeNode, KubeNamespace, KubeWorkload, KubePod, KubeContainer, KubeFolder:
		return true
	default:
		return false
	}
}

// GetCollectionWithObject get the corresponding collection name based on the container object resource
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
	case KubeContainer:
		return []string{BKTableNameBaseContainer}, nil
	case string(KubeDeployment):
		return []string{BKTableNameBaseDeployment}, nil
	case string(KubeDaemonSet):
		return []string{BKTableNameBaseDaemonSet}, nil
	case string(KubeStatefulSet):
		return []string{BKTableNameBaseStatefulSet}, nil
	case string(KubeGameStatefulSet):
		return []string{BKTableNameGameStatefulSet}, nil
	case string(KubeGameDeployment):
		return []string{BKTableNameGameDeployment}, nil
	case string(KubeCronJob):
		return []string{BKTableNameBaseCronJob}, nil
	case string(KubeJob):
		return []string{BKTableNameBaseJob}, nil
	case string(KubePodWorkload):
		return []string{BKTableNameBasePodWorkload}, nil
	case KubeWorkload:
		return GetWorkLoadTables(), nil
	default:
		return []string{}, errors.New("no corresponding table found")
	}
}

// IsKubeResourceKind determine whether it is a container resource object.
func IsKubeResourceKind(object string) bool {
	switch object {
	case KubeBusiness, KubeCluster, KubeNode, KubeFolder, KubeNamespace, string(KubeDeployment),
		string(KubeStatefulSet), string(KubeDaemonSet), string(KubeGameStatefulSet), string(KubeGameDeployment),
		string(KubeCronJob), string(KubeJob), string(KubePodWorkload):
		return true
	default:
		return false
	}
}

// GetKindByWorkLoadTableNameMap get the corresponding workload type according to the database table name.
func GetKindByWorkLoadTableNameMap(table string) (map[string]string, error) {
	switch table {
	case BKTableNameBaseDeployment:
		return map[string]string{
			table: string(KubeDeployment),
		}, nil
	case BKTableNameBaseStatefulSet:
		return map[string]string{
			table: string(KubeStatefulSet),
		}, nil
	case BKTableNameBaseDaemonSet:
		return map[string]string{
			table: string(KubeDaemonSet),
		}, nil
	case BKTableNameGameStatefulSet:
		return map[string]string{
			table: string(KubeGameStatefulSet),
		}, nil
	case BKTableNameGameDeployment:
		return map[string]string{
			table: string(KubeGameDeployment),
		}, nil
	case BKTableNameBaseCronJob:
		return map[string]string{
			table: string(KubeCronJob),
		}, nil
	case BKTableNameBaseJob:
		return map[string]string{
			table: string(KubeJob),
		}, nil
	case BKTableNameBasePodWorkload:
		return map[string]string{
			table: string(KubePodWorkload),
		}, nil
	default:
		return nil, errors.New("this table name does not exist")
	}

}

// IsWorkLoadKind whether the resource type is workload
func IsWorkLoadKind(kind string) bool {
	switch kind {
	case string(KubeDeployment), string(KubeStatefulSet), string(KubeDaemonSet), string(KubeJob),
		string(KubeCronJob), string(KubeGameStatefulSet), string(KubeGameDeployment), string(KubePodWorkload):
		return true
	default:
		return false
	}
}

// KubeAttrsRsp 容器资源属性回应
type KubeAttrsRsp struct {
	Field    string `json:"field"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
}

// SpecInfo information about container fields in cmdb.
type SpecInfo struct {
	// ClusterID cluster id in cc
	ClusterID *int64 `json:"bk_cluster_id" bson:"bk_cluster_id"`
	// NamespaceID namespace id in cc
	NamespaceID *int64 `json:"bk_namespace_id" bson:"bk_namespace_id"`
	//// WorkloadKind workload kind
	//WorkloadKind *string `json:"workload_kind" bson:"workload_kind"`
	//// WorkloadID workload id in cc
	//WorkloadID *int64 `json:"workload_id" bson:"workload_id"`
	Ref Ref `json:"ref" bson:"ref"`
	// NodeID node id in cc
	NodeID *int64 `json:"bk_node_id" bson:"bk_node_id"`
}

// validate validate the spec info
func (option *SpecInfo) validate() error {

	if option.ClusterID == nil || *option.ClusterID == 0 {
		return errors.New("cluster id must be set")
	}

	if option.NamespaceID == nil || *option.NamespaceID == 0 {
		return errors.New("namespace id must be set")
	}

	if option.NodeID == nil || *option.NodeID == 0 {
		return errors.New("node id must be set")
	}

	if option.Ref.Kind == "" {
		return errors.New("workload kind must be set")
	}

	if err := WorkloadType(option.Ref.Kind).Validate(); err != nil {
		return errors.New("workload is illegal type")
	}

	if option.Ref.ID == 0 {
		return errors.New("workload id must be set")
	}

	return nil
}
