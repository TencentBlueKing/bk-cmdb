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

	"configcenter/src/common/blog"
	types2 "configcenter/src/common/types"
	"configcenter/src/common/watch"
	"configcenter/src/thirdparty/monitor"
	"configcenter/src/thirdparty/monitor/meta"
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

// ReduceChainNode remove conflict chain node, returns reduced chain nodes
func ReduceChainNode(chainNodeMap map[string][]*watch.ChainNode, tenantID, flowKey string, conflictErr error,
	metrics *EventMetrics, rid string) map[string][]*watch.ChainNode {

	chainNodes := chainNodeMap[tenantID]

	rid = rid + ":" + chainNodes[0].Oid
	monitor.Collect(&meta.Alarm{
		RequestID: rid,
		Type:      meta.EventFatalError,
		Detail: fmt.Sprintf("run event flow, but got conflict %s tenant %s cursor with chain nodes",
			flowKey, tenantID),
		Module:    types2.CC_MODULE_CACHESERVICE,
		Dimension: map[string]string{"retry_conflict_nodes": "yes"},
	})

	if len(chainNodes) <= 1 {
		delete(chainNodeMap, tenantID)
		return chainNodeMap
	}

	for index, reducedChainNode := range chainNodes {
		if isConflictChainNode(reducedChainNode, conflictErr) {
			metrics.CollectConflict()
			chainNodes = append(chainNodes[:index], chainNodes[index+1:]...)

			// need do with retry with reduce
			blog.ErrorJSON("run flow, insert %s tenant %s event with reduce node %s, remain nodes: %s, rid: %s",
				flowKey, tenantID, reducedChainNode, chainNodes, rid)
			chainNodeMap[tenantID] = chainNodes
			return chainNodeMap
		}
	}

	// when no cursor conflict node is found, discard the first node and try to insert the others
	blog.ErrorJSON("run flow, insert %s tenant %s event with reduce node %s, remain nodes: %s, rid: %s",
		flowKey, tenantID, chainNodes[0], chainNodes[1:], rid)
	chainNodeMap[tenantID] = chainNodes[1:]
	return chainNodeMap
}

func isConflictChainNode(chainNode *watch.ChainNode, err error) bool {
	return strings.Contains(err.Error(), chainNode.Cursor) && strings.Contains(err.Error(), "index_cursor")
}
