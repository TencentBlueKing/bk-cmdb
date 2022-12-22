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
	"errors"
	"fmt"
	"reflect"
	"strings"

	"configcenter/src/common/criteria/enumor"
	"configcenter/src/storage/dal/table"
)

const (
	// JsonTag tag type is json
	JsonTag = "json"
)

// CommonSpecFieldsDescriptor public field properties
var CommonSpecFieldsDescriptor = table.FieldsDescriptors{
	{Field: BKIDField, Type: enumor.Numeric, IsRequired: true, IsEditable: false},
	{Field: BKSupplierAccountField, Type: enumor.String, IsRequired: true, IsEditable: false},
	{Field: CreatorField, Type: enumor.String, IsRequired: true, IsEditable: false},
	{Field: ModifierField, Type: enumor.String, IsRequired: true, IsEditable: true},
	{Field: CreateTimeField, Type: enumor.Timestamp, IsRequired: true, IsEditable: false},
	{Field: LastTimeField, Type: enumor.Timestamp, IsRequired: true, IsEditable: true},
}

// BizIDDescriptor bizID descriptor is taken out separately and not placed in CommonSpecFieldsDescriptor because
// bk_biz_id does not exist in the container table and needs to be processed separately.
var BizIDDescriptor = table.FieldsDescriptors{
	{Field: BKBizIDField, Type: enumor.Numeric, IsRequired: true, IsEditable: false},
}

// HostIDDescriptor the description used when other resources refer to the host.
var HostIDDescriptor = table.FieldsDescriptors{
	{Field: "bk_host_id", Type: enumor.Numeric, IsRequired: false, IsEditable: false},
}

// isRequiredField check if field is required Field is not filled
func isRequiredField(tag string, value reflect.Value, i int) error {
	// do a pocket inspection to prevent the front from being blocked
	if isCommonField(tag) {
		return nil
	}
	fieldValue := value.Field(i)
	if fieldValue.Kind() != reflect.Ptr && fieldValue.Kind() != reflect.UnsafePointer {
		return nil
	}
	if fieldValue.IsNil() {
		return fmt.Errorf("required fields cannot be empty, %s", tag)
	}
	return nil
}

func isNotEditableField(tag string, value reflect.Value, i int) bool {

	if isCommonField(tag) {
		return true
	}

	field := value.Field(i)
	if field.Kind() != reflect.Ptr && field.Kind() != reflect.UnsafePointer {
		return true
	}

	// check each variable for a null pointer.
	// if it is a null pointer, it means that
	// this field will not be updated, skip it directly.
	if field.IsNil() {
		return true
	}
	return false
}

// getFieldTag a variable with a non-null pointer gets the corresponding tag.
// for example, it needs to be compatible when the tag is "name,omitempty"
func getFieldTag(typeOfOption reflect.Type, tag string, i int) (string, bool) {

	tagTmp := typeOfOption.Field(i).Tag.Get(tag)
	tags := strings.Split(tagTmp, ",")

	if tags[0] == "" {
		return "", true
	}

	return tags[0], false
}

// isCommonField judges whether the field is a special field.
// If it belongs to common fields or biz_id field, it needs
// to skip unified verification whether it is creating a scene or updating a scene.
func isCommonField(field string) bool {
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

	for _, hostDescriptor := range HostIDDescriptor {
		if field == hostDescriptor.Field {
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
	case KubeFolder, KubePod:
		return "", nil
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

// IsKubeTopoResource determine whether it is a container object type.
func IsKubeTopoResource(object string) bool {
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

// KubeAttrsRsp kube resource attribute response
type KubeAttrsRsp struct {
	Field    string `json:"field"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
}

// SpecSimpleInfo information about container fields in cmdb.
type SpecSimpleInfo struct {
	// ClusterID cluster id in cc
	ClusterID int64 `json:"bk_cluster_id" bson:"bk_cluster_id"`
	// NamespaceID namespace id in cc
	NamespaceID int64 `json:"bk_namespace_id" bson:"bk_namespace_id"`
	// Ref the workload at the upper level of the pod
	Ref Reference `json:"ref" bson:"ref"`
	// NodeID node id in cc
	NodeID int64 `json:"bk_node_id" bson:"bk_node_id"`
}

// validate validate the spec info
func (option *SpecSimpleInfo) validate() error {

	if option.ClusterID == 0 {
		return errors.New("cluster id must be set")
	}

	if option.NamespaceID == 0 {
		return errors.New("namespace id must be set")
	}

	if option.NodeID == 0 {
		return errors.New("node id must be set")
	}

	if option.Ref.Kind == "" || option.Ref.ID == 0 {
		return errors.New("workload must be set")
	}

	if err := WorkloadType(option.Ref.Kind).Validate(); err != nil {
		return errors.New("workload is illegal type")
	}

	return nil
}
