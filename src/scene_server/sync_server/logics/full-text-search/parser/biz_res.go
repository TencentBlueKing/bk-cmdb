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

package parser

import (
	"errors"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/sync_server/logics/full-text-search/cache"
)

// bizResParser is the data parser for biz resource
type bizResParser struct {
	*objInstParser
}

func newBizResParser(index string, cs *ClientSet) *bizResParser {
	return &bizResParser{newObjInstParser(index, cs)}
}

// ParseData parse mongo data to es data
func (p *bizResParser) ParseData(info []mapstr.MapStr, coll string, rid string) (bool, mapstr.MapStr, error) {
	if len(info) == 0 {
		return false, nil, errors.New("data is empty")
	}
	data := info[0]

	// do not sync resource pool resource to es
	bizID, err := util.GetIntByInterface(data[common.BKAppIDField])
	if err != nil {
		blog.Errorf("parse %s biz id failed, err: %v, data: %+v, rid: %s", p.index, err, data, rid)
		return false, nil, err
	}

	if _, exists := cache.ResPoolBizIDMap.Load(bizID); exists {
		blog.Errorf("%s biz id %d is resource pool, skip, data: %+v, rid: %s", p.index, bizID, data, rid)
		return true, nil, nil
	}

	return p.objInstParser.ParseData(info, coll, rid)
}
