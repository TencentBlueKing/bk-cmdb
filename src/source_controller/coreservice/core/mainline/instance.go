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
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
	"encoding/json"
	"fmt"
)

// SearchMainlineBusinessTopo get topo tree of model on mainline
func (m *topoManager) SearchMainlineInstanceTopo(bkBizID int64, withDetail bool) (*metadata.TopoInstanceNode, error) {
	bizTopoNode, err := m.SearchMainlineModelTopo()
	if err != nil {
		blog.Errorf("get mainline model topo info failed, %+v", err)
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
	err = m.dbProxy.Table(common.BKTableNameBaseSet).Find(mongoCondition.ToMapStr()).All(ctx, &setInstances)
	if err != nil {
		blog.Errorf("get set instances failed, %+v", err)
		return nil, err
	}
	blog.V(5).Infof("SearchMainlineInstanceTopo setInstances: %+v", setInstances)

	// module instance list of target business
	mongoCondition = mongo.NewCondition()
	mongoCondition.Element(&mongo.Eq{Key: common.BKAppIDField, Val: bkBizID})

	moduleInstances := make([]mapstr.MapStr, 0)
	err = m.dbProxy.Table(common.BKTableNameBaseModule).Find(mongoCondition.ToMapStr()).All(ctx, &moduleInstances)
	if err != nil {
		blog.Errorf("get module instances failed, %+v", err)
		return nil, err
	}
	blog.V(5).Infof("SearchMainlineInstanceTopo moduleInstances: %+v", moduleInstances)

	// other mainline instance list of target business
	mongoCondition = mongo.NewCondition()
	mongoCondition.Element(&mongo.Eq{Key: common.BKAppIDField, Val: bkBizID})

	commonInstances := make([]mapstr.MapStr, 0)
	err = m.dbProxy.Table(common.BKTableNameBaseInst).Find(mongoCondition.ToMapStr()).All(ctx, &commonInstances)
	if err != nil {
		blog.Errorf("get common instances failed, %+v", err)
		return nil, err
	}
	blog.V(5).Infof("SearchMainlineInstanceTopo commonInstances: %+v", commonInstances)

	instanceMap := map[string]*metadata.TopoInstance{}
	allTopoInstances := []*metadata.TopoInstance{}
	bizTopoInstance := &metadata.TopoInstance{
		ObjectID:         common.BKInnerObjIDApp,
		InstanceID:       bkBizID,
		ParentInstanceID: 0,
		Detail:           map[string]interface{}{},
	}
	if withDetail == true {
		// get business detail here
		mongoCondition = mongo.NewCondition()
		mongoCondition.Element(&mongo.Eq{Key: common.BKAppIDField, Val: bkBizID})

		businessInstances := make([]mapstr.MapStr, 0)
		err = m.dbProxy.Table(common.BKTableNameBaseApp).Find(mongoCondition.ToMapStr()).All(ctx, &businessInstances)
		if err != nil {
			blog.Errorf("get business instances failed, %+v", err)
			return nil, err
		}
		blog.V(5).Infof("SearchMainlineInstanceTopo businessInstances: %+v", businessInstances)
		if len(businessInstances) == 0 {
			blog.Error("get business instances failed, not found")
			return nil, fmt.Errorf("business with bk_biz_id=%d not found", bkBizID)
		}
		bizTopoInstance.Detail = businessInstances[0]
	}
	allTopoInstances = append(allTopoInstances, bizTopoInstance)
	instanceMap[bizTopoInstance.Key()] = bizTopoInstance

	for _, instance := range setInstances {
		instanceID, err := util.GetInt64ByInterface(instance[common.BKSetIDField])
		if err != nil {
			blog.Errorf("parse isntanceID:%+v failed, %+v", instance[common.BKSetIDField], err)
			return nil, fmt.Errorf("parse isntanceID:%+v failed, %+v", instance[common.BKSetIDField], err)
		}
		parentInstanceID, err := util.GetInt64ByInterface(instance[common.BKInstParentStr])
		if err != nil {
			blog.Errorf("parse isntanceID:%+v failed, %+v", instance[common.BKInstParentStr], err)
			return nil, fmt.Errorf("parse isntanceID:%+v failed, %+v", instance[common.BKInstParentStr], err)
		}
		topoInstance := &metadata.TopoInstance{
			ObjectID:         common.BKInnerObjIDSet,
			InstanceID:       instanceID,
			ParentInstanceID: parentInstanceID,
			Detail:           map[string]interface{}{},
		}
		if withDetail == true {
			topoInstance.Detail = instance
		}
		allTopoInstances = append(allTopoInstances, topoInstance)
		instanceMap[topoInstance.Key()] = topoInstance
	}
	for _, instance := range moduleInstances {
		instanceID, err := util.GetInt64ByInterface(instance[common.BKModuleIDField])
		if err != nil {
			blog.Errorf("parse isntanceID:%+v failed, %+v", instance[common.BKModuleIDField], err)
			return nil, fmt.Errorf("parse isntanceID:%+v failed, %+v", instance[common.BKModuleIDField], err)
		}
		parentInstanceID, err := util.GetInt64ByInterface(instance[common.BKInstParentStr])
		if err != nil {
			blog.Errorf("parse isntanceID:%+v failed, %+v", instance[common.BKInstParentStr], err)
			return nil, fmt.Errorf("parse isntanceID:%+v failed, %+v", instance[common.BKInstParentStr], err)
		}
		topoInstance := &metadata.TopoInstance{
			ObjectID:         common.BKInnerObjIDModule,
			InstanceID:       instanceID,
			ParentInstanceID: parentInstanceID,
			Detail:           map[string]interface{}{},
		}
		if withDetail == true {
			topoInstance.Detail = instance
		}
		allTopoInstances = append(allTopoInstances, topoInstance)
		instanceMap[topoInstance.Key()] = topoInstance
	}
	for _, instance := range commonInstances {
		instanceID, err := util.GetInt64ByInterface(instance[common.BKInstIDField])
		if err != nil {
			blog.Errorf("parse isntanceID:%+v failed, %+v", instance[common.BKInstIDField], err)
			return nil, fmt.Errorf("parse isntanceID:%+v failed, %+v", instance[common.BKInstIDField], err)
		}
		parentInstanceID, err := util.GetInt64ByInterface(instance[common.BKInstParentStr])
		if err != nil {
			blog.Errorf("parse isntanceID:%+v failed, %+v", instance[common.BKInstParentStr], err)
			return nil, fmt.Errorf("parse isntanceID:%+v failed, %+v", instance[common.BKInstParentStr], err)
		}
		topoInstance := &metadata.TopoInstance{
			ObjectID:         instance[common.BKObjIDField].(string),
			InstanceID:       instanceID,
			ParentInstanceID: parentInstanceID,
			Detail:           map[string]interface{}{},
		}
		if withDetail == true {
			topoInstance.Detail = instance
		}
		allTopoInstances = append(allTopoInstances, topoInstance)
		instanceMap[topoInstance.Key()] = topoInstance
	}

	instanceMapStr, err := json.Marshal(instanceMap)
	if err != nil {
		blog.Errorf("json encode instanceMap failed, %+v", err)
		return nil, fmt.Errorf("json encode instanceMap failed, %+v", err)
	}
	blog.V(3).Infof("instanceMap: %s", instanceMapStr)

	// prepare loop that make sure all node's parent are exist in allTopoInstances
	for _, topoInstance := range allTopoInstances {
		blog.V(5).Infof("topoInstance: %+v", topoInstance)
		if topoInstance.ParentInstanceID == 0 {
			continue
		}
		parentObjectID := objectParentMap[topoInstance.ObjectID]
		parentKey := fmt.Sprintf("%s:%d", parentObjectID, topoInstance.ParentInstanceID)
		// check whether parent instance exist, if not, try to get it at best.
		_, exist := instanceMap[parentKey]
		if exist == true {
			continue
		}
		blog.Warnf("get parent of %+v with key=%s failed, not Found", topoInstance, parentKey)
		// There is a bug in legacy code that business before mainline model be created in cc_ObjectBase table has no bk_biz_id field
		// and therefore find parentInstance failed.
		// In this case current algorithm degenerate in to o(n) query cost.

		mongoCondition = mongo.NewCondition()
		mongoCondition.Element(&mongo.Eq{Key: common.BKInstIDField, Val: topoInstance.ParentInstanceID})

		missedInstances := make([]mapstr.MapStr, 0)
		err = m.dbProxy.Table(common.BKTableNameBaseInst).Find(mongoCondition.ToMapStr()).All(ctx, &missedInstances)
		if err != nil {
			blog.Errorf("get common instances failed, %+v", err)
			return nil, err
		}
		blog.V(5).Infof("SearchMainlineInstanceTopo missedInstances: %+v", missedInstances)
		if len(missedInstances) != 1 {
			blog.Errorf("SearchMainlineInstanceTopo foud unexpected count of missedInstances: %+v", missedInstances)
			return nil, fmt.Errorf("SearchMainlineInstanceTopo foud %d missedInstances with instanceID=%d", len(missedInstances), topoInstance.ParentInstanceID)
		}
		instance := missedInstances[0]
		instanceID, err := util.GetInt64ByInterface(instance[common.BKInstIDField])
		if err != nil {
			blog.Errorf("parse isntanceID:%+v failed, %+v", instance[common.BKInstIDField], err)
			return nil, fmt.Errorf("parse isntanceID:%+v failed, %+v", instance[common.BKInstIDField], err)
		}
		parentInstanceID, err := util.GetInt64ByInterface(instance[common.BKInstParentStr])
		if err != nil {
			blog.Errorf("parse isntanceID:%+v failed, %+v", instance[common.BKInstParentStr], err)
			return nil, fmt.Errorf("parse isntanceID:%+v failed, %+v", instance[common.BKInstParentStr], err)
		}
		topoInstance := &metadata.TopoInstance{
			ObjectID:         util.GetStrByInterface(instance[common.BKObjIDField]),
			InstanceID:       instanceID,
			ParentInstanceID: parentInstanceID,
			Detail:           map[string]interface{}{},
		}
		if withDetail == true {
			topoInstance.Detail = instance
		}
		allTopoInstances = append(allTopoInstances, topoInstance)
		instanceMap[topoInstance.Key()] = topoInstance
	}

	instanceMapStr, err = json.Marshal(instanceMap)
	if err != nil {
		blog.Errorf("json encode instanceMap failed, %+v", err)
		return nil, fmt.Errorf("json encode instanceMap failed, %+v", err)
	}
	blog.V(3).Infof("instanceMap: %s", instanceMapStr)

	var bizTopoInstanceNode *metadata.TopoInstanceNode
	topoInstanceNodeMap := map[string]*metadata.TopoInstanceNode{}
	for _, topoInstance := range allTopoInstances {
		blog.V(5).Infof("topoInstance: %+v", topoInstance)
		if topoInstance.ParentInstanceID == 0 {
			continue
		}
		parentObjectID := objectParentMap[topoInstance.ObjectID]
		parentKey := fmt.Sprintf("%s:%d", parentObjectID, topoInstance.ParentInstanceID)
		if _, exist := topoInstanceNodeMap[parentKey]; exist == false {
			parentInstance := instanceMap[parentKey]
			topoInstanceNode := &metadata.TopoInstanceNode{
				ObjectID:   parentInstance.ObjectID,
				InstanceID: parentInstance.InstanceID,
				Detail:     parentInstance.Detail,
				Children:   []*metadata.TopoInstanceNode{},
			}
			topoInstanceNodeMap[parentKey] = topoInstanceNode
		}

		parentInstanceNode := topoInstanceNodeMap[parentKey]

		// extract tree root node pointer
		if parentInstanceNode.ObjectID == common.BKInnerObjIDApp {
			bizTopoInstanceNode = parentInstanceNode
		}

		if _, exist := topoInstanceNodeMap[topoInstance.Key()]; exist == false {
			childTopoInstanceNode := &metadata.TopoInstanceNode{
				ObjectID:   topoInstance.ObjectID,
				InstanceID: topoInstance.InstanceID,
				Detail:     topoInstance.Detail,
				Children:   []*metadata.TopoInstanceNode{},
			}
			topoInstanceNodeMap[topoInstance.Key()] = childTopoInstanceNode
		}
		childTopoInstanceNode, _ := topoInstanceNodeMap[topoInstance.Key()]
		parentInstanceNode.Children = append(parentInstanceNode.Children, childTopoInstanceNode)
	}
	return bizTopoInstanceNode, nil
}
