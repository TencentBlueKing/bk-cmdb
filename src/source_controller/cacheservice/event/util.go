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

package event

import (
	"fmt"
	"strings"

	"configcenter/pkg/conv"
	"configcenter/src/common"
	"configcenter/src/common/metadata"
	"configcenter/src/common/watch"
	kubetypes "configcenter/src/kube/types"
)

var resourceKeyMap = map[watch.CursorType]Key{
	watch.Host:                    HostKey,
	watch.ModuleHostRelation:      ModuleHostRelationKey,
	watch.Biz:                     BizKey,
	watch.Set:                     SetKey,
	watch.Module:                  ModuleKey,
	watch.ObjectBase:              ObjectBaseKey,
	watch.Process:                 ProcessKey,
	watch.ProcessInstanceRelation: ProcessInstanceRelationKey,
	watch.HostIdentifier:          HostIdentityKey,
	watch.MainlineInstance:        MainlineInstanceKey,
	watch.InstAsst:                InstAsstKey,
	watch.BizSet:                  BizSetKey,
	watch.BizSetRelation:          BizSetRelationKey,
	watch.Plat:                    PlatKey,
	watch.KubeCluster:             KubeClusterKey,
	watch.KubeNode:                KubeNodeKey,
	watch.KubeNamespace:           KubeNamespaceKey,
	watch.KubeWorkload:            KubeWorkloadKey,
	watch.KubePod:                 KubePodKey,
	watch.Project:                 ProjectKey,
}

// GetResourceKeyWithCursorType get resource key
func GetResourceKeyWithCursorType(res watch.CursorType) (Key, error) {
	key, exists := resourceKeyMap[res]
	if !exists {
		return key, fmt.Errorf("unsupported cursor type %s", res)
	}

	return key, nil
}

// IsConflictError check if a error is event cursor conflict/duplicate error
func IsConflictError(err error) bool {
	if strings.Contains(err.Error(), "duplicate key error") {
		return true
	}

	if strings.Contains(err.Error(), "index_cursor dup key") {
		return true
	}

	return false
}

// HostArchive TODO
type HostArchive struct {
	Oid    string              `bson:"oid"`
	Detail metadata.HostMapStr `bson:"detail"`
}

// ObjInstTablePrefixRegex TODO
const ObjInstTablePrefixRegex = "^" + common.BKObjectInstShardingTablePrefix

// InstAsstTablePrefixRegex TODO
const InstAsstTablePrefixRegex = "^" + common.BKObjectInstAsstShardingTablePrefix

// ConvertLabel 由于目前使用版本的mongodb不支持key中包含的.的查询，存入db的时候是将.以编码的方式存入，这里需要进行解码
func ConvertLabel(podDetail map[string]interface{}) map[string]interface{} {
	labels, ok := podDetail[kubetypes.LabelsField]
	if !ok {
		return podDetail
	}

	labelMap, ok := labels.(map[string]string)
	if !ok {
		return podDetail
	}

	newLabels := make(map[string]string)
	for key, val := range labelMap {
		newLabels[conv.DecodeDot(key)] = val
	}
	podDetail[kubetypes.LabelsField] = newLabels

	return podDetail
}
