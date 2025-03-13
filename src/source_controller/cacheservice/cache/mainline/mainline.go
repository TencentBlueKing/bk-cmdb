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

package mainline

import (
	"strings"
	"time"

	"configcenter/pkg/tenant"
	"configcenter/src/apimachinery/discovery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/driver/redis"
)

// mainlineCache is an instance to refresh mainline topology cache.
type mainlineCache struct {
	isMaster discovery.ServiceManageInterface
}

// Run start to watch and refresh the mainline topology cache.
func (m *mainlineCache) Run() error {
	kit := rest.NewKit()
	if err := m.refreshMainlineTopoCache(kit); err != nil {
		return err
	}

	go func() {
		// wait for a moment and then start loop.
		time.Sleep(5 * time.Minute)
		for {
			if !m.isMaster.IsMaster() {
				blog.V(4).Infof("loop refresh mainline topology cache, but not master, skip.")
				time.Sleep(time.Minute)
				continue
			}

			kit := rest.NewKit()
			blog.V(4).Infof("start refresh mainline topology cache, rid: %s", kit.Rid)
			if err := m.refreshMainlineTopoCache(kit); err != nil {
				time.Sleep(time.Minute)
				continue
			}

			blog.V(4).Infof("finished refresh mainline topology cache, rid: %s", kit.Rid)
			time.Sleep(5 * time.Minute)
		}
	}()

	return nil
}

// refreshMainlineTopoCache refresh mainline topology cache for all tenants.
func (m *mainlineCache) refreshMainlineTopoCache(kit *rest.Kit) error {
	err := tenant.ExecForAllTenants(func(tenantID string) error {
		kit = kit.WithTenant(tenantID)
		_, err := refreshAndGetTopologyRank(kit)
		if err != nil {
			blog.Errorf("refresh tenant: %s mainline topology cache failed, err: %v, rid: %s", tenantID, err, kit.Rid)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// getMainlineTopology get mainline topology's association details.
func getMainlineTopology(kit *rest.Kit) ([]mainlineAssociation, error) {
	relations := make([]mainlineAssociation, 0)
	filter := mapstr.MapStr{
		common.AssociationKindIDField: common.AssociationKindMainline,
	}
	err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameObjAsst).Find(filter).All(kit.Ctx, &relations)
	if err != nil {
		blog.Errorf("get mainline topology association failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}
	return relations, nil
}

// rankMainlineTopology ranks the biz topology to an array, starting from biz to host
func rankMainlineTopology(relations []mainlineAssociation) []string {
	asstMap := make(map[string]string)
	for _, relation := range relations {
		asstMap[relation.AssociateTo] = relation.ObjectID
	}

	rank := make([]string, 0)

	for next := "biz"; next != ""; next = asstMap[next] {
		rank = append(rank, next)
	}

	return rank
}

// refreshTopologyRank is to refresh the business's rank information.
func refreshTopologyRank(kit *rest.Kit, rank []string) {
	// then set the rank to cache
	value := strings.Join(rank, ",")
	err := redis.Client().Set(kit.Ctx, genTopologyKey(kit), value, detailTTLDuration).Err()
	if err != nil {
		blog.Errorf("refresh mainline topology rank, but update to cache failed, err: %v", err)
		// do not return, it will be refreshed next round.
	}
}

// refreshAndGetTopologyRank refresh the business's topology rank to cache, from biz model to host model.
func refreshAndGetTopologyRank(kit *rest.Kit) ([]string, error) {
	// read information from mongodb
	relations, err := getMainlineTopology(kit)
	if err != nil {
		blog.Errorf("refresh mainline topology rank, but get it from mongodb failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}
	// rank start from biz to host
	rank := rankMainlineTopology(relations)
	refreshTopologyRank(kit, rank)

	return rank, nil
}
