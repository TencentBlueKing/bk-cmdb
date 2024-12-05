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
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/common/valid"
	"configcenter/src/source_controller/coreservice/multilingual"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/driver/mongodb"
)

// InstanceMainline TODO
type InstanceMainline struct {
	bkBizID   int64
	modelTree *metadata.TopoModelNode

	lang language.DefaultCCLanguageIf

	modelIDs        []string
	objectParentMap map[string]string

	businessInstance  mapstr.MapStr
	setInstances      []mapstr.MapStr
	moduleInstances   []mapstr.MapStr
	mainlineInstances []mapstr.MapStr

	instanceMap      map[string]*metadata.TopoInstance
	allTopoInstances []*metadata.TopoInstance

	root *metadata.TopoInstanceNode
}

// NewInstanceMainline TODO
func NewInstanceMainline(lang language.DefaultCCLanguageIf, proxy dal.DB, bkBizID int64) (*InstanceMainline, error) {
	im := &InstanceMainline{
		lang:              lang,
		bkBizID:           bkBizID,
		objectParentMap:   map[string]string{},
		setInstances:      make([]mapstr.MapStr, 0),
		moduleInstances:   make([]mapstr.MapStr, 0),
		mainlineInstances: make([]mapstr.MapStr, 0),
		allTopoInstances:  make([]*metadata.TopoInstance, 0),
		instanceMap:       map[string]*metadata.TopoInstance{},
	}
	return im, nil
}

// SetModelTree TODO
func (im *InstanceMainline) SetModelTree(modelTree *metadata.TopoModelNode) {
	// step1
	im.modelTree = modelTree
}

// LoadModelParentMap TODO
func (im *InstanceMainline) LoadModelParentMap(kit *rest.Kit) {
	// step2
	im.modelIDs = im.modelTree.LeftestObjectIDList()
	for idx, objectID := range im.modelIDs {
		if idx == 0 {
			continue
		}
		im.objectParentMap[objectID] = im.modelIDs[idx-1]
	}
	blog.V(5).Infof("LoadModelParentMap mainline models: %#v, objectParentMap: %#v, rid: %s", im.modelIDs,
		im.objectParentMap, kit.Rid)
}

// LoadSetInstances TODO
func (im *InstanceMainline) LoadSetInstances(kit *rest.Kit) error {
	// set instance list of target business
	filter := map[string]interface{}{
		common.BKAppIDField: im.bkBizID,
	}
	err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseSet).Find(filter).All(kit.Ctx, &im.setInstances)
	if err != nil {
		blog.Errorf("get set instances by business: %d failed, %+v, cond: %#v, rid: %s", im.bkBizID, err, filter,
			kit.Rid)
		return fmt.Errorf("get set instances by business:%d failed, %+v", im.bkBizID, err)
	}
	multilingual.TranslateInstanceName(im.lang, common.BKInnerObjIDSet, im.setInstances)
	blog.V(5).Infof("get set instances by business: %d result: %+v, cond: %#v, rid: %s", im.bkBizID, im.setInstances,
		filter, kit.Rid)
	return nil
}

// LoadModuleInstances TODO
func (im *InstanceMainline) LoadModuleInstances(kit *rest.Kit) error {
	rid := util.ExtractRequestIDFromContext(kit.Ctx)
	// module instance list of target business
	filter := map[string]interface{}{
		common.BKAppIDField: im.bkBizID,
	}
	err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseModule).Find(filter).All(kit.Ctx,
		&im.moduleInstances)
	if err != nil {
		blog.Errorf("get module instances by business: %d failed, err: %v, cond: %#v, rid: %s", im.bkBizID, err,
			filter, rid)
		return fmt.Errorf("get module instances by business:%d failed, %+v", im.bkBizID, err)
	}
	multilingual.TranslateInstanceName(im.lang, common.BKInnerObjIDModule, im.moduleInstances)
	blog.V(5).Infof("get module instances by business: %d, result: %v, cond:%v, rid: %s", im.bkBizID,
		im.moduleInstances, filter, rid)
	return nil
}

// LoadMainlineInstances TODO
func (im *InstanceMainline) LoadMainlineInstances(kit *rest.Kit) error {

	// load other mainline instance(except business,set,module) list of target business
	for _, objectID := range im.modelIDs {
		if valid.IsInnerObject(objectID) {
			continue
		}

		filter := map[string]interface{}{
			common.BKObjIDField: objectID,
			common.BKAppIDField: im.bkBizID,
		}

		mainlineInstances := []mapstr.MapStr{}

		err := mongodb.Shard(kit.ShardOpts()).
			Table(common.GetObjectInstTableName(objectID, kit.TenantID)).
			Find(filter).
			All(kit.Ctx, &mainlineInstances)

		if err != nil {
			blog.Errorf("get other mainline instances by business: %d failed, err: %v, cond: %#v, rid: %s",
				im.bkBizID, err, filter, kit.Rid)
			return fmt.Errorf("get other mainline instances by business:%d failed, %+v", im.bkBizID, err)
		}
		im.mainlineInstances = append(im.mainlineInstances, mainlineInstances...)
	}

	return nil
}

// ConstructBizTopoInstance TODO
func (im *InstanceMainline) ConstructBizTopoInstance(kit *rest.Kit, withDetail bool) error {
	// enqueue business instance to allTopoInstances, instanceMap
	bizTopoInstance := &metadata.TopoInstance{
		ObjectID:         common.BKInnerObjIDApp,
		InstanceID:       im.bkBizID,
		ParentInstanceID: 0,
		Detail:           map[string]interface{}{},
	}

	// get business detail here
	bizFilter := map[string]interface{}{
		common.BKAppIDField: im.bkBizID,
	}
	err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseApp).Find(bizFilter).One(kit.Ctx,
		&im.businessInstance)
	if err != nil {
		blog.Errorf("get business instances by business: %d failed, err: %v, cond: %#v, rid: %s", im.bkBizID, err,
			kit.Rid)
		return fmt.Errorf("get business instances by business:%d failed, err: %+v", im.bkBizID, err)
	}
	multilingual.TranslateInstanceName(im.lang, common.BKInnerObjIDApp, []mapstr.MapStr{im.businessInstance})
	blog.V(5).Infof("SearchMainlineInstanceTopo businessInstances: %+v, rid: %s", im.businessInstance, kit.Rid)
	bizTopoInstance.InstanceName = util.GetStrByInterface(im.businessInstance[common.BKAppNameField])
	if withDetail {
		bizTopoInstance.Detail = im.businessInstance
	}

	im.allTopoInstances = append(im.allTopoInstances, bizTopoInstance)
	im.instanceMap[bizTopoInstance.Key()] = bizTopoInstance
	return nil
}

// OrganizeSetInstance TODO
func (im *InstanceMainline) OrganizeSetInstance(kit *rest.Kit, withDetail bool) error {
	for _, instance := range im.setInstances {
		instanceID, err := util.GetInt64ByInterface(instance[common.BKSetIDField])
		if err != nil {
			blog.Errorf("parse instanceID:%+v to int64 failed, err: %v, rid: %s", instance[common.BKSetIDField], err,
				kit.Rid)
			return fmt.Errorf("parse instanceID:%+v to int64 failed, %+v", instance[common.BKSetIDField], err)
		}
		parentInstanceID, err := util.GetInt64ByInterface(instance[common.BKInstParentStr])
		if err != nil {
			blog.Errorf("parse instanceID:%+v to int64 failed, err: %v, rid: %s", instance[common.BKInstParentStr],
				err, kit.Rid)
			return fmt.Errorf("parse instanceID:%+v to int64 failed, %+v", instance[common.BKInstParentStr], err)
		}

		defaultFieldValue, err := util.GetInt64ByInterface(instance[common.BKDefaultField])
		if err != nil {
			blog.Errorf("parse set instance default field failed, default: %+v, err: %v, rid: %s",
				instance[common.BKDefaultField], err, kit.Rid)
			return fmt.Errorf("parse set instance default field failed, default: %+v, err: %+v",
				instance[common.BKDefaultField], err)
		}

		instanceName := util.GetStrByInterface(instance[common.BKSetNameField])
		topoInstance := &metadata.TopoInstance{
			Default:          defaultFieldValue,
			ObjectID:         common.BKInnerObjIDSet,
			InstanceID:       instanceID,
			InstanceName:     instanceName,
			ParentInstanceID: parentInstanceID,
			Detail:           map[string]interface{}{},
		}
		if withDetail {
			topoInstance.Detail = instance
		}
		im.allTopoInstances = append(im.allTopoInstances, topoInstance)
		im.instanceMap[topoInstance.Key()] = topoInstance
	}
	return nil
}

// OrganizeModuleInstance TODO
func (im *InstanceMainline) OrganizeModuleInstance(kit *rest.Kit, withDetail bool) error {
	for _, instance := range im.moduleInstances {
		instanceID, err := util.GetInt64ByInterface(instance[common.BKModuleIDField])
		if err != nil {
			blog.Errorf("parse instanceID: %+v to int64 failed, %+v, rid: %s", instance[common.BKModuleIDField], err,
				kit.Rid)
			return fmt.Errorf("parse instanceID:%+v to int64 failed, %+v", instance[common.BKModuleIDField], err)
		}
		parentInstanceID, err := util.GetInt64ByInterface(instance[common.BKInstParentStr])
		if err != nil {
			blog.Errorf("parse instanceID: %+v to int64 failed, %+v, rid: %s", instance[common.BKInstParentStr], err,
				kit.Rid)
			return fmt.Errorf("parse instanceID: %+v to int64 failed, %+v", instance[common.BKInstParentStr], err)
		}

		defaultFieldValue, err := util.GetInt64ByInterface(instance[common.BKDefaultField])
		if err != nil {
			blog.Errorf("parse module instance default field failed, default: %+v, err: %+v, rid: %s",
				instance[common.BKDefaultField], err, kit.Rid)
			return fmt.Errorf("parse module instance default field failed, default: %+v, err: %+v",
				instance[common.BKDefaultField], err)
		}

		instanceName := util.GetStrByInterface(instance[common.BKModuleNameField])

		topoInstance := &metadata.TopoInstance{
			Default:          defaultFieldValue,
			ObjectID:         common.BKInnerObjIDModule,
			InstanceID:       instanceID,
			InstanceName:     instanceName,
			ParentInstanceID: parentInstanceID,
			Detail:           map[string]interface{}{},
		}
		if withDetail {
			topoInstance.Detail = instance
		}
		im.allTopoInstances = append(im.allTopoInstances, topoInstance)
		im.instanceMap[topoInstance.Key()] = topoInstance
	}
	return nil
}

// OrganizeMainlineInstance TODO
func (im *InstanceMainline) OrganizeMainlineInstance(kit *rest.Kit, withDetail bool) error {
	for _, instance := range im.mainlineInstances {
		instanceID, err := util.GetInt64ByInterface(instance[common.BKInstIDField])
		if err != nil {
			blog.Errorf("parse instanceID: %+v to int64 failed, err: %v, rid: %s", instance[common.BKInstIDField], err,
				kit.Rid)
			return fmt.Errorf("parse instanceID:%+v to int64 failed, %+v", instance[common.BKInstIDField], err)
		}
		parentInstanceID, err := util.GetInt64ByInterface(instance[common.BKInstParentStr])
		if err != nil {
			blog.Errorf("parse instanceID: %+v to int64 failed, err: %v, rid: %s", instance[common.BKInstParentStr],
				err, kit.Rid)
			return fmt.Errorf("parse instanceID:%+v to int64 failed, %+v", instance[common.BKInstParentStr], err)
		}
		instanceName := util.GetStrByInterface(instance[common.BKInstNameField])
		topoInstance := &metadata.TopoInstance{
			ObjectID:         instance[common.BKObjIDField].(string),
			InstanceID:       instanceID,
			InstanceName:     instanceName,
			ParentInstanceID: parentInstanceID,
			Detail:           map[string]interface{}{},
		}
		if withDetail {
			topoInstance.Detail = instance
		}
		im.allTopoInstances = append(im.allTopoInstances, topoInstance)
		im.instanceMap[topoInstance.Key()] = topoInstance
	}
	return nil
}

// CheckAndFillingMissingModels TODO
func (im *InstanceMainline) CheckAndFillingMissingModels(kit *rest.Kit, withDetail bool) error {

	for _, topoInstance := range im.allTopoInstances {
		blog.V(5).Infof("topo instance: %#v, rid: %s", topoInstance, kit.Rid)
		if topoInstance.ParentInstanceID == 0 {
			continue
		}
		var parentKey string
		if topoInstance.ObjectID == common.BKInnerObjIDSet && topoInstance.Default == 1 {
			// `空闲机池` 是一种特殊的set，它用来包含空闲机和故障机两个模块，它的父节点直接是业务（不论是否有自定义层级）
			// 这类特殊情况的结点是业务，不需要重复获取，ConstructInstanceTopoTree 会做进一步处理
			parentKey = fmt.Sprintf("%s:%d", common.BKInnerObjIDApp, topoInstance.ParentInstanceID)
		} else {
			parentObjectID := im.objectParentMap[topoInstance.ObjectID]
			parentKey = fmt.Sprintf("%s:%d", parentObjectID, topoInstance.ParentInstanceID)
		}
		// check whether parent instance exist, if not, try to get it at best.
		_, exist := im.instanceMap[parentKey]
		if exist {
			continue
		}
		blog.Warnf("get parent of %+v with key=%s failed, not Found, rid: %s", topoInstance, parentKey, kit.Rid)
		// There is a bug in legacy code that business before mainline model be created in cc_ObjectBase table has
		// no bk_biz_id field and therefore find parentInstance failed. in this case current algorithm
		// degenerate in to o(n) query cost.

		topoInstance, needSkip, err := im.getMissingModelInstance(kit, kit.TenantID, *topoInstance, withDetail)
		if err != nil {
			blog.Errorf("check and filling missing models failed, topoInstance: %v, err: %v, rid: %s", topoInstance,
				err, kit.Rid)
			return err
		}
		if needSkip {
			continue
		}

		im.allTopoInstances = append(im.allTopoInstances, &topoInstance)
		im.instanceMap[topoInstance.Key()] = &topoInstance
	}
	return nil
}

// getMissingModelInstance get missing model instance
func (im *InstanceMainline) getMissingModelInstance(kit *rest.Kit, tenantID string,
	topoInstance metadata.TopoInstance, withDetail bool) (metadata.TopoInstance, bool, error) {

	filter := map[string]interface{}{common.BKInstIDField: topoInstance.ParentInstanceID}

	missedInstances := make([]mapstr.MapStr, 0)
	err := mongodb.Shard(kit.ShardOpts()).Table(common.GetObjectInstTableName(topoInstance.ObjectID, tenantID)).
		Find(filter).All(kit.Ctx, &missedInstances)
	if err != nil {
		blog.Errorf("get common instances failed, err: %v, rid: %s", topoInstance.ParentInstanceID, err, kit.Rid)
		return topoInstance, false, err
	}
	blog.V(5).Infof("get missed instances by id: %d results: %+v, rid: %s", topoInstance.ParentInstanceID,
		missedInstances, kit.Rid)

	if len(missedInstances) == 0 {
		if topoInstance.ObjectID == common.BKInnerObjIDSet && im.bkBizID == topoInstance.ParentInstanceID {
			// `空闲机池` 是一种特殊的set，它用来包含空闲机和故障机两个模块，它的父节点直接是业务（不论是否有自定义层级）
			// 这类特殊情况的结点是业务，不需要重复获取，ConstructInstanceTopoTree 会做进一步处理
			return topoInstance, true, nil
		} else {
			// parent id not found, ignore node
			blog.Warnf("found unexpected count of missedInstances: %#v, cond: %#v, rid: %s", missedInstances,
				filter, kit.Rid)
			return topoInstance, true, nil
		}
	}
	if len(missedInstances) > 1 {
		blog.Errorf("found too many(%d) missedInstances: %#v by id: %d, cond: %#v, rid: %s", len(missedInstances),
			missedInstances, topoInstance.ParentInstanceID, filter, kit.Rid)
		return topoInstance, false, fmt.Errorf("found too many(%d) missedInstances: %+v by id: %d",
			len(missedInstances), missedInstances, topoInstance.ParentInstanceID)
	}
	instance := missedInstances[0]
	instanceID, err := util.GetInt64ByInterface(instance[common.BKInstIDField])
	if err != nil {
		blog.Errorf("parse instanceID: %+v to int64 failed, err: %v, instanceID:%v, rid: %s", err,
			instance[common.BKInstIDField], kit.Rid)
		return topoInstance, false, err
	}

	var parentInstanceID int64
	parentValue, existed := instance[common.BKInstParentStr]
	if existed {
		parentInstanceID, err = util.GetInt64ByInterface(parentValue)
		if err != nil {
			blog.Errorf("parse instanceID to int64 failed, err: %v, instanceID: %v, rid: %s", err,
				instance[common.BKInstParentStr], kit.Rid)
			return topoInstance, false, err
		}
	} else {
		// `空闲机池` 是一种特殊的set，它用来包含空闲机和故障机两个模块，它的父节点直接是业务（不论是否有自定义层级）
		// 这类特殊情况的结点是业务，不需要重复获取，ConstructInstanceTopoTree 会做进一步处理
		if topoInstance.ObjectID == common.BKInnerObjIDSet && im.bkBizID == topoInstance.ParentInstanceID {
			return topoInstance, true, nil
		}
		blog.Errorf("construct biz topo tree, instance doesn't have field %s, instance: %+v, rid: %s",
			common.BKInstParentStr, instance, kit.Rid)
		return topoInstance, false, fmt.Errorf("construct biz topo tree, instance doesn't have field %s, "+
			"instance: %+v", common.BKInstParentStr, instance)
	}
	blog.V(7).Infof("model: %s, instance: %d, parent: %d, rid: %s", topoInstance.ObjectID, topoInstance.InstanceID,
		parentInstanceID, kit.Rid)

	topoInstance = metadata.TopoInstance{
		ObjectID:         util.GetStrByInterface(instance[common.BKObjIDField]),
		InstanceName:     util.GetStrByInterface(instance[common.BKInstNameField]),
		InstanceID:       instanceID,
		ParentInstanceID: parentInstanceID,
		Detail:           map[string]interface{}{},
	}
	if withDetail {
		topoInstance.Detail = instance
	}
	return topoInstance, false, nil
}

// ConstructInstanceTopoTree TODO
func (im *InstanceMainline) ConstructInstanceTopoTree(kit *rest.Kit, withDetail bool) error {

	topoInstanceNodeMap := map[string]*metadata.TopoInstanceNode{}
	for index := 0; index < len(im.allTopoInstances); index++ {
		topoInstance := im.allTopoInstances[index]
		blog.V(5).Infof("topoInstance: %+v, rid: %s", topoInstance, kit.Rid)
		if topoInstance.ParentInstanceID == 0 {
			continue
		}

		parentObjectID := im.objectParentMap[topoInstance.ObjectID]
		parentKey := fmt.Sprintf("%s:%d", parentObjectID, topoInstance.ParentInstanceID)
		if _, exist := topoInstanceNodeMap[parentKey]; !exist {
			parentInstance, needSkip, err := im.getParentInstance(kit, kit.TenantID, parentObjectID, parentKey,
				*topoInstance)
			if err != nil {
				blog.Errorf("get parent instance failed, err: %v, topoInstance: %v, rid: %s", err, topoInstance,
					kit.Rid)
				return err
			}
			if needSkip {
				continue
			}
			topoInstanceNode := &metadata.TopoInstanceNode{
				ObjectID:     parentInstance.ObjectID,
				InstanceID:   parentInstance.InstanceID,
				InstanceName: parentInstance.InstanceName,
				Detail:       parentInstance.Detail,
				Children:     []*metadata.TopoInstanceNode{},
			}
			topoInstanceNodeMap[parentKey] = topoInstanceNode
		}

		parentInstanceNode := topoInstanceNodeMap[parentKey]

		// extract tree root node pointer
		if parentInstanceNode.ObjectID == common.BKInnerObjIDApp {
			im.root = parentInstanceNode
		}

		if _, exist := topoInstanceNodeMap[topoInstance.Key()]; !exist {
			childTopoInstanceNode := &metadata.TopoInstanceNode{
				ObjectID:     topoInstance.ObjectID,
				InstanceID:   topoInstance.InstanceID,
				InstanceName: topoInstance.InstanceName,
				Detail:       topoInstance.Detail,
				Children:     []*metadata.TopoInstanceNode{},
			}
			topoInstanceNodeMap[topoInstance.Key()] = childTopoInstanceNode
		}
		childTopoInstanceNode := topoInstanceNodeMap[topoInstance.Key()]
		parentInstanceNode.Children = append(parentInstanceNode.Children, childTopoInstanceNode)
	}
	return nil
}

// getParentInstance get parent instance
func (im *InstanceMainline) getParentInstance(kit *rest.Kit, tenantID string, parentObjectID string,
	parentKey string, topoInstance metadata.TopoInstance) (*metadata.TopoInstance, bool, error) {

	parentInstance, exist := im.instanceMap[parentKey]
	if !exist {
		// 空闲机池 是一种特殊的set，它用来包含空闲机和故障机两个模块，它的父节点直接是业务（不论是否有自定义层级）
		if topoInstance.ObjectID == common.BKInnerObjIDSet && im.bkBizID == topoInstance.ParentInstanceID {
			parentObjectID = common.BKInnerObjIDApp
			parentKey = fmt.Sprintf("%s:%d", parentObjectID, im.bkBizID)
			parentInstance, exist = im.instanceMap[parentKey]
		}
		if !exist {
			cond := map[string]interface{}{
				common.BKObjIDField:  parentObjectID,
				common.BKInstIDField: topoInstance.ParentInstanceID,
			}

			inst := mapstr.MapStr{}
			err := mongodb.Shard(kit.ShardOpts()).
				Table(common.GetObjectInstTableName(parentObjectID, tenantID)).
				Find(cond).
				One(context.Background(), &inst)

			if err != nil {
				if isNotFound := mongodb.IsNotFoundError(err); !isNotFound {
					blog.Errorf("get mainline instances failed, filter: %+v, err: %+v, rid: %s", cond, err, kit.Rid)
					return parentInstance, false, err
				} else {
					im.mainlineInstances = append(im.mainlineInstances, inst)
					blog.Errorf("unexpected err, parent instance not found, instance: %+v, rid: %s", topoInstance,
						kit.Rid)
					return parentInstance, true, nil
				}
			}

			parentValue, existed := inst[common.BKInstParentStr]
			if !existed {
				blog.Errorf("get mainline instances failed, field %s not in db data, data: %+v, rid: %s",
					common.BKInstParentStr, inst, kit.Rid)
				return parentInstance, false, fmt.Errorf("get mainline instances failed, field %s not in db data,"+
					" data: %+v", common.BKInstParentStr, inst)
			}
			parentParentID, err := util.GetInt64ByInterface(parentValue)
			if err != nil {
				blog.Errorf("get mainline instances failed, field %s parse into int failed, data: %+v, err: %+v,"+
					" rid:  %s", common.BKInstParentStr, inst, err, kit.Rid)
				return parentInstance, false, fmt.Errorf("get mainline instances failed, field %s parse into int"+
					" failed, data: %+v, err: %+v", common.BKInstParentStr, inst, err)
			}
			parentInstance = &metadata.TopoInstance{
				ObjectID:         parentObjectID,
				InstanceID:       topoInstance.ParentInstanceID,
				InstanceName:     util.GetStrByInterface(inst),
				ParentInstanceID: parentParentID,
				Detail:           inst,
			}
			im.instanceMap[parentKey] = parentInstance
			im.allTopoInstances = append(im.allTopoInstances, parentInstance)
		}
	}

	return parentInstance, false, nil
}

// GetInstanceMap TODO
func (im *InstanceMainline) GetInstanceMap() map[string]*metadata.TopoInstance {
	return im.instanceMap
}

// GetRoot TODO
func (im *InstanceMainline) GetRoot() *metadata.TopoInstanceNode {
	return im.root
}

// OrganizeTopo organize topo instance
func (im *InstanceMainline) OrganizeTopo(kit *rest.Kit, bkBizID int64, withDetail bool) error {
	if err := im.LoadSetInstances(kit); err != nil {
		blog.Errorf("get set instances by business: %d failed, err: %v, rid: %s", bkBizID, err, kit.Rid)
		return fmt.Errorf("get set instances by business:%d failed, %+v", bkBizID, err)
	}

	if err := im.LoadModuleInstances(kit); err != nil {
		blog.Errorf("get module instances by business: %d failed, err: %v, rid: %s", bkBizID, err, kit.Rid)
		return fmt.Errorf("get module instances by business:%d failed, %+v", bkBizID, err)
	}

	if err := im.LoadMainlineInstances(kit); err != nil {
		blog.Errorf("get other mainline instances by business: %d failed, err: %v, rid: %s", bkBizID, err, kit.Rid)
		return fmt.Errorf("get other mainline instances by business:%d failed, %+v", bkBizID, err)
	}

	if err := im.ConstructBizTopoInstance(kit, withDetail); err != nil {
		blog.Errorf("construct business: %d detail as topo instance failed, err: %v, rid: %s", bkBizID, err, kit.Rid)
		return fmt.Errorf("construct business:%d detail as topo instance failed, %+v", bkBizID, err)
	}

	if err := im.OrganizeSetInstance(kit, withDetail); err != nil {
		blog.Errorf("organize set instance failed, businessID: %d, err: %v, rid: %s", bkBizID, err, kit.Rid)
		return fmt.Errorf("organize set instance failed, businessID:%d, %+v", bkBizID, err)
	}

	if err := im.OrganizeModuleInstance(kit, withDetail); err != nil {
		blog.Errorf("organize module instance failed, businessID: %d, err: %v, rid: %s", bkBizID, err, kit.Rid)
		return fmt.Errorf("organize module instance failed, businessID:%d, %+v", bkBizID, err)
	}

	if err := im.OrganizeMainlineInstance(kit, withDetail); err != nil {
		blog.Errorf("organize other mainline instance failed, businessID: %d, err: %v, rid: %s", bkBizID, err, kit.Rid)

		return fmt.Errorf("organize other mainline instance failed, businessID:%d, %+v", bkBizID, err)
	}
	return nil
}
