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
	"encoding/json"
	"fmt"

	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/driver/mongodb"
)

// SearchMainlineInstanceTopo get topo tree of mainline model
func (m *topoManager) SearchMainlineInstanceTopo(kit *rest.Kit, bkBizID int64,
	withDetail bool) (*metadata.TopoInstanceNode, error) {

	bizTopoNode, err := m.SearchMainlineModelTopo(kit, false)
	if err != nil {
		blog.Errorf("get mainline model topo info failed, err: %v, rid: %s", err, kit.Rid)
		return nil, fmt.Errorf("get mainline model topo info failed, %+v", err)
	}
	blog.V(9).Infof("model mainline: %+v, rid: %s", bizTopoNode, kit.Rid)

	im, err := NewInstanceMainline(m.lang.CreateDefaultCCLanguageIf(httpheader.GetLanguage(kit.Header)), mongodb.Shard(
		kit.ShardOpts()), bkBizID)
	if err != nil {
		blog.Errorf("new instance mainline failed, bizID: %d, err: %v, rid: %s", bkBizID, err, kit.Rid)
		return nil, fmt.Errorf("new mainline instance by business:%d failed, %+v", bkBizID, err)
	}

	im.SetModelTree(bizTopoNode)
	im.LoadModelParentMap(kit)

	if err := im.OrganizeTopo(kit, bkBizID, withDetail); err != nil {
		blog.Errorf("organize topo instances failed, err: %v, bizID: %d, rid: %s", err, bkBizID, kit.Rid)
		return nil, fmt.Errorf("get set instances by business:%d failed, %+v", bkBizID, err)
	}

	instanceMap := im.GetInstanceMap()
	instanceMapStr, err := json.Marshal(instanceMap)
	if err != nil {
		blog.Errorf("json encode instanceMap:%+v failed, err: %v, rid: %s", instanceMap, err, kit.Rid)
		return nil, fmt.Errorf("json encode instanceMap:%+v failed, %+v", instanceMap, err)
	}
	blog.V(5).Infof("instanceMap before check is: %s, rid: %s", instanceMapStr, kit.Rid)

	if err := im.CheckAndFillingMissingModels(kit, withDetail); err != nil {
		blog.Errorf("check and filling missing models failed, business:%d, err: %v, rid: %s", bkBizID, err, kit.Rid)
		return nil, fmt.Errorf("check and filling missing models failed, business:%d %+v", bkBizID, err)
	}

	instanceMapStr, err = json.Marshal(im.GetInstanceMap())
	if err != nil {
		blog.Errorf("json encode instanceMap failed, err: %v, rid: %s", err, kit.Rid)
		return nil, fmt.Errorf("json encode instanceMap failed, %+v", err)
	}
	blog.V(5).Infof("instanceMap after check: %s, rid: %s", instanceMapStr, kit.Rid)

	if err := im.ConstructInstanceTopoTree(kit, withDetail); err != nil {
		blog.Errorf("get other mainline instances by business:%d failed, err: %v, rid: %s", bkBizID, err, kit.Rid)
		return nil, fmt.Errorf("get other mainline instances by business:%d failed, %+v", bkBizID, err)
	}

	root := im.GetRoot()
	blog.V(9).Infof("topo instance tree root is: %+v, rid: %s", root, kit.Rid)
	treeData, err := json.Marshal(root)
	if err != nil {
		blog.Errorf("get other mainline instances by business:%d failed, err: %v, rid: %s", bkBizID, err, kit.Rid)
		return root, nil
	}
	blog.V(9).Infof("topo instance tree root data is: %s, rid: %s", treeData, kit.Rid)
	return root, nil
}
