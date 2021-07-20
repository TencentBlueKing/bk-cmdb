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

package topo_tree

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"configcenter/src/source_controller/cacheservice/cache/business"

	"github.com/tidwall/gjson"
)

func NewTopologyTree(client *business.Client) *TopologyTree {
	return &TopologyTree{bizCache: client}
}

type TopologyTree struct {
	bizCache *business.Client
}

func (t *TopologyTree) SearchTopologyTree(opt *SearchOption) ([]*Topology, error) {
	if err := opt.Validate(); err != nil {
		return nil, err
	}
	bizList := make([]business.BizBaseInfo, 0)
	// get business id
	if opt.BusinessID > 0 {
		// ignore business name, because you have set a business id.
		if len(opt.BusinessName) == 0 {
			biz, err := t.bizCache.GetBusiness(opt.BusinessID)
			if err != nil {
				return nil, err
			}
			opt.BusinessName = gjson.Get(biz, "bk_biz_name").String()
		}

		bizList = append(bizList, business.BizBaseInfo{
			BusinessID:   opt.BusinessID,
			BusinessName: opt.BusinessName,
		})
	} else if opt.BusinessID == -1 {
		// all the business
		base, err := t.bizCache.GetBizBaseList()
		if err != nil {
			return nil, err
		}
		for _, biz := range base {
			bizList = append(bizList, biz)
		}

	} else {
		// find business id with business name.
		base, err := t.bizCache.GetBizBaseList()
		if err != nil {
			return nil, err
		}

		for _, biz := range base {
			matched, err := matchName(biz.BusinessName, opt.BusinessName)
			if err != nil {
				return nil, err
			}
			if matched {
				bizList = append(bizList, biz)
			}
		}
		// check if the hit biz is too much.
		if len(bizList) >= overHead {
			return nil, OverHeadError
		}
	}

	// we have already find all the business need to search.
	// now we search them each.
	allTopology := make([]*Topology, 0)
	for _, biz := range bizList {
		topo, err := t.SearchWithBusiness(&SearchOption{
			BusinessID:   biz.BusinessID,
			BusinessName: biz.BusinessName,
			SetName:      opt.SetName,
			ModuleName:   opt.ModuleName,
			Level:        opt.Level,
		})
		if err != nil {
			return nil, err
		}
		if topo != nil {
			allTopology = append(allTopology, topo)
		}
	}

	return allTopology, nil
}

// match instance name with regexp and case insensitive.
func matchName(src, toMatch string) (bool, error) {
	reg, err := regexp.Compile("(?i)" + strings.TrimSpace(toMatch))
	if err != nil {
		return false, err
	}
	return reg.MatchString(src), nil
}

// Obviously, we search it from top to bottom.
func (t *TopologyTree) SearchWithBusiness(opt *SearchOption) (*Topology, error) {
	do := doSearch{
		bizCache: t.bizCache,
		opt:      opt,
	}
	return do.createTopology()
}

type doSearch struct {
	bizCache       *business.Client
	opt            *SearchOption
	customInstance []instance
	setInstance    []instance
	moduleInstance []instance
	hitError       error
}

func (do *doSearch) searchCustomLevel() {

	if len(do.opt.Level.Object) == 0 && len(do.opt.Level.InstName) == 0 {
		return
	}

	if len(do.opt.Level.Object) == 0 {
		do.hitError = errors.New("empty level.object")
		return
	}

	if len(do.opt.Level.InstName) == 0 {
		do.hitError = errors.New("level.inst_name is empty")
		return
	}

	// search custom instance.
	list, err := do.bizCache.GetCustomLevelBaseList(do.opt.Level.Object, do.opt.BusinessID)
	if err != nil {
		do.hitError = fmt.Errorf("get custom base list failed, err: %v", err)
		return
	}

	custom := make([]instance, 0)
	for _, inst := range list {
		hit, err := matchName(inst.InstanceName, do.opt.Level.InstName)
		if err != nil {
			do.hitError = fmt.Errorf("match custom object instance name failed, err: %v", err)
			return
		}

		if !hit {
			continue
		}

		custom = append(custom, instance{
			name:     inst.InstanceName,
			id:       inst.InstanceID,
			parentID: inst.ParentID,
		})
	}

	if len(custom) > overHead {
		do.hitError = OverHeadError
		return
	}

	do.customInstance = custom
}

func (do *doSearch) searchSet() {
	if do.hitError != nil {
		return
	}

	if len(do.opt.SetName) == 0 {
		return
	}

	list, err := do.bizCache.GetSetBaseList(do.opt.BusinessID)
	if err != nil {
		do.hitError = fmt.Errorf("get set base list failed, err: %v", err)
		return
	}

	sets := make([]instance, 0)
	for _, set := range list {
		hit, err := matchName(set.SetName, do.opt.SetName)
		if err != nil {
			do.hitError = fmt.Errorf("match set name failed, err: %v", err)
			return
		}

		if !hit {
			continue
		}

		sets = append(sets, instance{
			name:     set.SetName,
			id:       set.SetID,
			parentID: set.ParentID,
		})
	}

	if len(sets) > overHead {
		do.hitError = OverHeadError
		return
	}

	do.setInstance = sets

}

func (do *doSearch) searchModule() {
	if do.hitError != nil {
		return
	}

	if len(do.opt.ModuleName) == 0 {
		return
	}

	list, err := do.bizCache.GetModuleBaseList(do.opt.BusinessID)
	if err != nil {
		do.hitError = fmt.Errorf("get module base list failed, err: %v", err)
		return
	}

	modules := make([]instance, 0)
	for _, module := range list {
		hit, err := matchName(module.ModuleName, do.opt.ModuleName)
		if err != nil {
			do.hitError = fmt.Errorf("match module name failed, err: %v", err)
			return
		}

		if !hit {
			continue
		}

		modules = append(modules, instance{
			name:     module.ModuleName,
			id:       module.ModuleID,
			parentID: module.SetID,
		})
	}

	if len(modules) > overHead {
		do.hitError = OverHeadError
		return
	}

	do.moduleInstance = modules
}

func (do *doSearch) createTopology() (*Topology, error) {
	if do.opt.BusinessID == 0 {
		return nil, errors.New("create topology, but biz id is 0")
	}
	do.searchCustomLevel()
	do.searchSet()
	do.searchModule()

	if do.hitError != nil {
		return nil, do.hitError
	}

	// prepare the tree from bottom to biz
	if len(do.customInstance) != 0 {
		return do.withCustomLevel()
	}

	if len(do.setInstance) != 0 {
		return do.withSet()
	}

	if len(do.moduleInstance) != 0 {
		return do.withModule()
	}

	return nil, nil
}

func (do *doSearch) withCustomLevel() (*Topology, error) {
	if len(do.customInstance) == 0 {
		return nil, nil
	}

	topo, err := do.bizCache.GetTopology()
	if err != nil {
		return nil, fmt.Errorf("construct custom level, but get topo failed, err: %v", err)
	}
	// find next level
	hit := false
	next := 0
	for idx, level := range topo {
		if level == do.opt.Level.Object {
			hit = true
			next = idx
		}
	}
	if !hit {
		return nil, fmt.Errorf("construct custom level, but can not find this custom level: %s", do.opt.Level.Object)
	}

	if next == 0 || (next >= (len(topo) - 1)) {
		return nil, fmt.Errorf("construct custom level, but got invalid topo %v", topo)
	}
	next += 1
	nextLevel := ""
	switch topo[next] {
	case "set":
		nextLevel = "set"
	case "biz", "module", "host":
		return nil, fmt.Errorf("construct custom level, but got invalid topo %v", topo)
	default:
		nextLevel = topo[next]
	}

	// construct next tree
	trees := make([]Tree, 0)
	for _, cur := range do.customInstance {
		tree := Tree{
			Object:   do.opt.Level.Object,
			InstName: cur.name,
			InstID:   cur.id,
		}

		subTree := make([]Tree, 0)
		if nextLevel == "set" {
			sets, err := do.bizCache.GetSetBaseList(do.opt.BusinessID)
			if err != nil {
				return nil, fmt.Errorf("get set list failed, err: %v", err)
			}
			for _, set := range sets {
				if set.ParentID == cur.id {
					subTree = append(subTree, Tree{
						Object:   "set",
						InstName: set.SetName,
						InstID:   set.SetID,
						Children: nil,
					})
				}
			}
		} else {
			customs, err := do.bizCache.GetCustomLevelBaseList(nextLevel, do.opt.BusinessID)
			if err != nil {
				return nil, fmt.Errorf("get object base list failed, err: %v", err)
			}

			for _, cut := range customs {
				if cut.ParentID == cur.id {
					subTree = append(subTree, Tree{
						Object:   cut.ObjectID,
						InstName: cut.InstanceName,
						InstID:   cut.InstanceID,
						Children: nil,
					})
				}
			}

		}

		tree.Children = subTree

		trees = append(trees, tree)
	}

	return &Topology{
		BusinessID:   do.opt.BusinessID,
		BusinessName: do.opt.BusinessName,
		Trees:        trees,
	}, nil
}

func (do *doSearch) withSet() (*Topology, error) {
	if len(do.setInstance) == 0 {
		return nil, nil
	}
	trees := make([]Tree, 0)
	for _, set := range do.setInstance {
		tree := Tree{
			Object:   "set",
			InstName: set.name,
			InstID:   set.id,
		}

		modules, err := do.bizCache.GetModuleBaseList(do.opt.BusinessID)
		if err != nil {
			return nil, fmt.Errorf("get module base list failed, err: %v", err)
		}
		subTrees := make([]Tree, 0)
		for _, mod := range modules {
			if mod.SetID != set.id {
				continue
			}
			subTrees = append(subTrees, Tree{
				Object:   "module",
				InstName: mod.ModuleName,
				InstID:   mod.ModuleID,
				Children: nil,
			})
		}
		tree.Children = subTrees
		trees = append(trees, tree)
	}

	return &Topology{
		BusinessID:   do.opt.BusinessID,
		BusinessName: do.opt.BusinessName,
		Trees:        trees,
	}, nil

}

func (do *doSearch) withModule() (*Topology, error) {
	if len(do.moduleInstance) == 0 {
		return nil, nil
	}

	trees := make([]Tree, 0)
	for _, mod := range do.moduleInstance {
		trees = append(trees, Tree{
			Object:   "module",
			InstName: mod.name,
			InstID:   mod.id,
			Children: nil,
		})
	}
	return &Topology{
		BusinessID:   do.opt.BusinessID,
		BusinessName: do.opt.BusinessName,
		Trees:        trees,
	}, nil
}
