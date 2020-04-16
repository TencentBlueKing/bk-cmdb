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
	"regexp"
	"strings"

	"configcenter/src/source_controller/coreservice/cache/business"
)

type TopologyTree struct {
	bizCache *business.Client
}

func (t *TopologyTree) Search(opt *SearchOption) ([]Topology, error) {
	if err := opt.Validate(); err != nil {
		return nil, err
	}
	bizList := make([]int64, 0)
	// get business id
	if opt.BusinessID > 0 {
		// ignore business name, because you have set a business id.
		opt.BusinessName = ""
		bizList = append(bizList, opt.BusinessID)
	} else if opt.BusinessID == -1 {
		// all the business
		base, err := t.bizCache.GetBizBaseList()
		if err != nil {
			return nil, err
		}
		for _, biz := range base {
			bizList = append(bizList, biz.BusinessID)
		}
		
	} else {
		// find business id with business name.
		base, err := t.bizCache.GetBizBaseList()
		if err != nil {
			return nil, err
		}
		
		for _, biz := range base {
			matched, err := t.matchName(opt.BusinessName, biz.BusinessName)
			if err != nil {
				return nil, err
			}
			if matched {
				bizList = append(bizList, biz.BusinessID)
			}
		}
		// check if the hit biz is too much.
		if len(bizList) >= overHead {
			return nil, OverHeadError
		}
	}
	
	// we have already find all the business need to search.
	// now we search them each.
	
	
	

	return nil, nil
}

// match instance name with regexp and case insensitive.
func (t *TopologyTree) matchName(src, toMatch string) (bool, error) {
	reg, err := regexp.Compile("(?i)"+strings.Replace(src, " ", "[ \\._-]", -1))
	if err != nil {
		return false, err
	}
	return reg.MatchString(toMatch), nil
}

// Obviously, we search it from top to bottom.
func (t *TopologyTree) SearchWithBusiness(opt *SearchOption) {
	
}

type searchOpt struct {
	opt *SearchOption
	topoLevel []string
}

func (s *searchOpt) SearchCustomLevel() {
	
}


