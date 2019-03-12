/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mainline

import (
	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/source_controller/coreservice/core"
	"fmt"
)

// SearchMainlineBusinessTopo get topo tree of model on mainline
func (m *topoManager) SearchMainlineInstanceTopo(bkBizID int64, withDetail bool) (*metadata.TopoInstanceNode, error) {
	bizTopoNode, err := m.SearchMainlineModelTopo()
	if err != nil {
		return nil, err
	}

	objectIDs := bizTopoNode.LeftestObjectIDList()
	objectParentMap := map[string]string{}
	for idx, objectID := range objectIDs {
		if idx == 0 {
			continue
		}
		objectParentMap[objectID] = objectIDs[idx-1]
	}

	ctx := core.ContextParams{}

	// set instance list of target business
	mongoCondition := mongo.NewCondition()
	mongoCondition.Element(&mongo.Eq{Key: common.BKAppIDField, Val: bkBizID})

	setInstances := make([]mapstr.MapStr, 0)
	err = m.dbProxy.Table(common.BKTableNameBaseSet).Find(mongoCondition.ToMapStr()).All(ctx, setInstances)
	if err != nil {
		return nil, err
	}

	// module instance list of target business
	mongoCondition = mongo.NewCondition()
	mongoCondition.Element(&mongo.Eq{Key: common.BKAppIDField, Val: bkBizID})

	moduleInstances := make([]mapstr.MapStr, 0)
	err = m.dbProxy.Table(common.BKTableNameBaseModule).Find(mongoCondition.ToMapStr()).All(ctx, moduleInstances)
	if err != nil {
		return nil, err
	}

	// other mainline instance list of target business
	mongoCondition = mongo.NewCondition()
	mongoCondition.Element(&mongo.Eq{Key: common.BKAppIDField, Val: bkBizID})

	commonInstances := make([]mapstr.MapStr, 0)
	err = m.dbProxy.Table(common.BKTableNameBaseInst).Find(mongoCondition.ToMapStr()).All(ctx, commonInstances)
	if err != nil {
		return nil, err
	}

	instanceMap := map[string]*metadata.TopoInstance{}
	allTopoInstances := []*metadata.TopoInstance{}
	bizTopoInstance := &metadata.TopoInstance{
		ObjectID:         "biz",
		InstanceID:       bkBizID,
		ParentInstanceID: 0,
	}
	if withDetail == true {
		// get business detail here
		mongoCondition = mongo.NewCondition()
		mongoCondition.Element(&mongo.Eq{Key: common.BKAppIDField, Val: bkBizID})

		businessInstances := make([]mapstr.MapStr, 0)
		err = m.dbProxy.Table(common.BKTableNameBaseApp).Find(mongoCondition.ToMapStr()).All(ctx, businessInstances)
		if err != nil {
			return nil, err
		}
		if len(businessInstances) == 0 {
			return nil, fmt.Errorf("business with bk_biz_id=%d not found", bkBizID)
		}
		bizTopoInstance.Detail = businessInstances[0]
	}
	allTopoInstances = append(allTopoInstances, bizTopoInstance)
	instanceMap[bizTopoInstance.Key()] = bizTopoInstance

	for _, instance := range setInstances {
		topoInstance := &metadata.TopoInstance{
			ObjectID:         "set",
			InstanceID:       instance[common.BKSetIDField].(int64),
			ParentInstanceID: instance[common.BKInstParentStr].(int64),
		}
		if withDetail == true {
			topoInstance.Detail = instance
		}
		allTopoInstances = append(allTopoInstances, topoInstance)
		instanceMap[topoInstance.Key()] = topoInstance
	}
	for _, instance := range moduleInstances {
		topoInstance := &metadata.TopoInstance{
			ObjectID:         "module",
			InstanceID:       instance[common.BKModuleIDField].(int64),
			ParentInstanceID: instance[common.BKInstParentStr].(int64),
		}
		if withDetail == true {
			topoInstance.Detail = instance
		}
		allTopoInstances = append(allTopoInstances, topoInstance)
		instanceMap[topoInstance.Key()] = topoInstance
	}
	for _, instance := range commonInstances {
		topoInstance := &metadata.TopoInstance{
			ObjectID:         instance[common.BKObjIDField].(string),
			InstanceID:       instance[common.BKModuleIDField].(int64),
			ParentInstanceID: instance[common.BKInstParentStr].(int64),
		}
		if withDetail == true {
			topoInstance.Detail = instance
		}
		allTopoInstances = append(allTopoInstances, topoInstance)
		instanceMap[topoInstance.Key()] = topoInstance
	}

	var bizTopoInstanceNode *metadata.TopoInstanceNode
	topoInstanceNodeMap := map[string]*metadata.TopoInstanceNode{}
	for _, topoInstance := range allTopoInstances {
		parentObjectID := objectParentMap[topoInstance.ObjectID]
		parentKey := fmt.Sprintf("%s:%d", parentObjectID, topoInstance.ParentInstanceID)
		if _, exist := topoInstanceNodeMap[parentKey]; exist == false {
			parentInstance := instanceMap[parentKey]
			topoInstanceNode := &metadata.TopoInstanceNode{
				ObjectID:   parentInstance.ObjectID,
				InstanceID: parentInstance.InstanceID,
				Detail:     parentInstance.Detail,
			}
			topoInstanceNodeMap[parentKey] = topoInstanceNode
		}

		parentInstanceNode := topoInstanceNodeMap[parentKey]

		// extract tree root node pointer
		if parentInstanceNode.ObjectID == "biz" {
			bizTopoInstanceNode = parentInstanceNode
		}

		if _, exist := topoInstanceNodeMap[topoInstance.Key()]; exist == false {
			childTopoInstanceNode := &metadata.TopoInstanceNode{
				ObjectID:   topoInstance.ObjectID,
				InstanceID: topoInstance.InstanceID,
				Detail:     topoInstance.Detail,
			}
			topoInstanceNodeMap[topoInstance.Key()] = childTopoInstanceNode
		}
		childTopoInstanceNode, _ := topoInstanceNodeMap[topoInstance.Key()]
		parentInstanceNode.Children = append(parentInstanceNode.Children, childTopoInstanceNode)
	}
	return bizTopoInstanceNode, nil
}
