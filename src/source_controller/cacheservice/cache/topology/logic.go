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

package topology

import (
	"context"
	"errors"

	"configcenter/src/common"
	"configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
)

// genBusinessTopology generate a fully business topology data.
func (t *Topology) genBusinessTopology(ctx context.Context, biz int64) (*BizBriefTopology, error) {
	rid := ctx.Value(common.ContextRequestIDField)

	filter := mapstr.MapStr{
		common.BKAppIDField: biz,
	}
	detail := new(BizBase)
	err := t.db.Table(common.BKTableNameBaseApp).Find(filter).Fields(bizBaseFields...).One(ctx, detail)
	if err != nil {
		blog.Errorf("get biz: %d detail failed, err: %v, rid: %v", biz, err, rid)
		return nil, err
	}

	idle, common, err := t.getBusinessTopology(ctx, biz)
	if err != nil {
		blog.Errorf("get biz topology nodes from db failed, err: %v, rid: %v", err, rid)
		return nil, err
	}

	return &BizBriefTopology{
		Biz:   detail,
		Idle:  idle,
		Nodes: common,
	}, nil
}

// getBusinessTopology construct a business's fully Topology data, separated with inner set and
// other common nodes
func (t *Topology) getBusinessTopology(ctx context.Context, biz int64) ([]*Node, []*Node, error) {
	rid := ctx.Value(common.ContextRequestIDField)

	// read from secondary in mongodb cluster.
	ctx = util.SetDBReadPreference(ctx, common.SecondaryPreferredMode)

	reverseRank, err := t.getMainlineReverseRank(ctx)
	if err != nil {
		blog.Errorf("get biz: %d topology, but get mainline rank failed, err: %v, rid: %v", biz, err, rid)
		return nil, nil, err
	}

	var previousNodes map[int64][]*Node
	var idleSetNodes map[int64][]*Node
	// cycle from set to business
	for _, level := range reverseRank[2:] {
		if level == "set" {
			idleNodes, normalNodes, err := t.genSetNodes(ctx, biz)
			if err != nil {
				return nil, nil, err
			}

			idleSetNodes = idleNodes

			// update previous nodes
			previousNodes = normalNodes
			continue
		}

		if level == "biz" {
			break
		}

		customNodes, err := t.genCustomNodes(ctx, biz, level, previousNodes)
		if err != nil {
			return nil, nil, err
		}

		// update previous nodes with custom nodes
		previousNodes = customNodes
	}

	if previousNodes == nil || idleSetNodes == nil {
		return nil, nil, errors.New("invalid business topology")
	}

	inner, exists := idleSetNodes[biz]
	if !exists {
		blog.ErrorJSON("can not find biz inner set topology, origin: %s, rid: %s", idleSetNodes, rid)
		return nil, nil, errors.New("invalid biz inner set topology")
	}

	return inner, previousNodes[biz], nil
}

// genCustomNodes generate custom object's instance node with it's parent map
func (t *Topology) genCustomNodes(ctx context.Context, biz int64, object string, previousNodes map[int64][]*Node) (
	map[int64][]*Node, error) {

	if previousNodes == nil {
		previousNodes = make(map[int64][]*Node)
	}

	customList, err := t.listCustomInstance(ctx, biz, object)
	if err != nil {
		return nil, err
	}

	reminder := make(map[int64]struct{})
	customParentMap := make(map[int64][]*Node)
	for _, custom := range customList {
		if _, exists := customParentMap[custom.ParentID]; !exists {
			customParentMap[custom.ParentID] = make([]*Node, 0)
		}

		if _, exists := reminder[custom.ID]; exists {
			continue
		}
		reminder[custom.ID] = struct{}{}

		customParentMap[custom.ParentID] = append(customParentMap[custom.ParentID], &Node{
			Object: object,
			ID:     custom.ID,
			Name:   custom.Name,
			// fill it later
			SubNodes: previousNodes[custom.ID],
		})
	}

	return customParentMap, nil

}

// listCustomInstance list a biz's custom object's all instances
func (t *Topology) listCustomInstance(ctx context.Context, biz int64, object string) ([]*customBase, error) {
	filter := mapstr.MapStr{common.BKAppIDField: biz, "bk_obj_id": object}
	all := make([]*customBase, 0)
	start := uint64(0)
	for {
		oneStep := make([]*customBase, 0)
		err := t.db.Table(common.BKTableNameBaseInst).Find(filter).Fields(customBaseFields...).Start(start).
			Limit(step).Sort(common.BKInstIDField).All(ctx, &oneStep)
		if err != nil {
			blog.Errorf("get biz: %d custom object: %s instance list failed, err: %v", biz, object, err)
			return nil, err
		}

		all = append(all, oneStep...)

		if len(oneStep) < step {
			// we got all the data
			break
		}

		// update start position
		start += step
	}

	return all, nil
}

// genSetNodes generate set's node with it's parent map
func (t *Topology) genSetNodes(ctx context.Context, biz int64) (idle map[int64][]*Node, normal map[int64][]*Node,
	err error) {

	moduleNodes, err := t.genModulesNodes(ctx, biz)
	if err != nil {
		return nil, nil, err
	}

	setList, err := t.listSets(ctx, biz)
	if err != nil {
		return nil, nil, err
	}

	reminder := make(map[int64]struct{})
	idleSetNodes := make(map[int64][]*Node)
	normalSetParentMap := make(map[int64][]*Node)
	var current map[int64][]*Node
	for _, set := range setList {

		if set.Default > 0 {
			// not the common set
			current = idleSetNodes
		} else {
			// common sets
			current = normalSetParentMap
		}

		_, exists := current[set.ParentID]
		if !exists {
			current[set.ParentID] = make([]*Node, 0)
		}

		if _, exists = reminder[set.ID]; exists {
			continue
		}
		reminder[set.ID] = struct{}{}

		current[set.ParentID] = append(current[set.ParentID], &Node{
			Object:   "set",
			ID:       set.ID,
			Name:     set.Name,
			Default:  &set.Default,
			SubNodes: moduleNodes[set.ID],
		})
	}

	return idleSetNodes, normalSetParentMap, nil

}

func (t *Topology) listSets(ctx context.Context, biz int64) ([]*setBase, error) {
	filter := mapstr.MapStr{common.BKAppIDField: biz}
	all := make([]*setBase, 0)
	start := uint64(0)
	for {
		oneStep := make([]*setBase, 0)
		err := t.db.Table(common.BKTableNameBaseSet).Find(filter).Fields(setBaseFields...).Start(start).
			Limit(step).Sort(common.BKSetIDField).All(ctx, &oneStep)
		if err != nil {
			blog.Errorf("get biz: %d set list failed, err: %v", biz, err)
			return nil, err
		}

		all = append(all, oneStep...)

		if len(oneStep) < step {
			// we got all the data
			break
		}

		// update start position
		start += step
	}

	return all, nil
}

// genModulesNodes generate module's node with it's parent set map
func (t *Topology) genModulesNodes(ctx context.Context, biz int64) (map[int64][]*Node, error) {
	moduleList, err := t.listModules(ctx, biz)
	if err != nil {
		return nil, err
	}

	reminder := make(map[int64]struct{})
	moduleParentMap := make(map[int64][]*Node)
	for idx := range moduleList {
		module := moduleList[idx]
		_, exists := moduleParentMap[module.SetID]
		if !exists {
			moduleParentMap[module.SetID] = make([]*Node, 0)
		}

		if _, exists = reminder[module.ID]; exists {
			continue
		}
		reminder[module.ID] = struct{}{}

		moduleParentMap[module.SetID] = append(moduleParentMap[module.SetID], &Node{
			Object:   "module",
			ID:       module.ID,
			Name:     module.Name,
			Default:  &module.Default,
			SubNodes: nil,
		})
	}

	return moduleParentMap, nil

}

// listModules list a business's all modules
func (t *Topology) listModules(ctx context.Context, biz int64) ([]*moduleBase, error) {
	filter := mapstr.MapStr{common.BKAppIDField: biz}
	all := make([]*moduleBase, 0)
	start := uint64(0)
	for {
		oneStep := make([]*moduleBase, 0)
		err := t.db.Table(common.BKTableNameBaseModule).Find(filter).Fields(moduleBaseFields...).Start(start).
			Limit(step).Sort(common.BKModuleIDField).All(ctx, &oneStep)
		if err != nil {
			blog.Errorf("get biz: %d module list failed, err: %v", biz, err)
			return nil, err
		}

		all = append(all, oneStep...)

		if len(oneStep) < step {
			// we got all the data
			break
		}

		// update start position
		start += step
	}

	return all, nil
}

// getMainlineReverseRank rank mainline object from module to biz
func (t *Topology) getMainlineReverseRank(ctx context.Context) ([]string, error) {
	relations := make([]mainlineAssociation, 0)
	filter := mapstr.MapStr{
		common.AssociationKindIDField: common.AssociationKindMainline,
	}
	err := t.db.Table(common.BKTableNameObjAsst).Find(filter).Fields(mainlineAsstFields...).All(ctx, &relations)
	if err != nil {
		blog.Errorf("get mainline topology association failed, err: %v", err)
		return nil, err
	}

	// rank mainline object
	rank := make([]string, 0)
	next := "biz"
	rank = append(rank, next)
	for _, relation := range relations {
		if relation.AssociateTo == next {
			rank = append(rank, relation.ObjectID)
			next = relation.ObjectID
			continue
		} else {
			for _, rel := range relations {
				if rel.AssociateTo == next {
					rank = append(rank, rel.ObjectID)
					next = rel.ObjectID
					break
				}
			}
		}
	}

	return util.ReverseArrayString(rank), nil
}

// listAllBusiness list all business brief info
func (t *Topology) listAllBusiness(ctx context.Context) ([]*BizBase, error) {

	filter := mapstr.MapStr{}
	all := make([]*BizBase, 0)
	start := uint64(0)
	for {
		oneStep := make([]*BizBase, 0)
		err := t.db.Table(common.BKTableNameBaseApp).Find(filter).Fields(bizBaseFields...).Start(start).
			Limit(step).Sort(common.BKAppIDField).All(ctx, &oneStep)
		if err != nil {
			return nil, err
		}

		all = append(all, oneStep...)

		if len(oneStep) < step {
			// we got all the data
			break
		}

		// update start position
		start += step
	}

	return all, nil
}

func getBreifTopoCacheRefreshMinutes() int {
	duration, err := configcenter.Int("cacheService.briefTopologySyncIntervalMinutes")
	if err != nil {
		blog.Errorf("get brief biz topology cache refresh interval minutes failed, err: %v, use default value 15.", err)
		return defaultRefreshIntervalMinutes
	}

	if duration < 2 {
		blog.Warnf("got invalid brief biz topology cache refresh interval minutes %d, < 2min, use default value 15.")
		return defaultRefreshIntervalMinutes
	}

	return duration
}
