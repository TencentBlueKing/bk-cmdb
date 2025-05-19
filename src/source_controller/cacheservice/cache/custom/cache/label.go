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

package cache

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"configcenter/pkg/conv"
	"configcenter/pkg/tenant"
	"configcenter/src/apimachinery/discovery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/lock"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	kubetypes "configcenter/src/kube/types"
	"configcenter/src/source_controller/cacheservice/cache/custom/types"
	"configcenter/src/source_controller/cacheservice/cache/tools"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/driver/redis"
)

// PodLabelCache is pod label cache
type PodLabelCache struct {
	isMaster   discovery.ServiceManageInterface
	keyCache   *CountCache
	valueCache *CountCache
}

// NewPodLabelCache new pod label cache
func NewPodLabelCache(isMaster discovery.ServiceManageInterface) *PodLabelCache {
	return &PodLabelCache{
		isMaster:   isMaster,
		keyCache:   NewCountCache(Key{resType: types.PodLabelKeyType, ttl: 3 * 24 * time.Hour}),
		valueCache: NewCountCache(Key{resType: types.PodLabelValueType, ttl: 3 * 24 * time.Hour}),
	}
}

// genKeyRedisKey generate redis key for pod label key cache
func (c *PodLabelCache) genKeyRedisKey(bizID int64) string {
	return strconv.FormatInt(bizID, 10)
}

// genValueRedisKey generate redis key for pod label value cache
func (c *PodLabelCache) genValueRedisKey(bizID int64, key string) string {
	return fmt.Sprintf("%d:%s", bizID, key)
}

// GetKeys get biz pod label keys
func (c *PodLabelCache) GetKeys(kit *rest.Kit, bizID int64) ([]string, error) {
	redisKey := c.genKeyRedisKey(bizID)

	existRes, err := redis.Client().Exists(kit.Ctx, c.keyCache.key.Key(kit.TenantID, redisKey)).Result()
	if err != nil {
		blog.Errorf("check if biz %d pod label cache exists failed, err: %v, rid: %s", bizID, err, kit.Rid)
		return nil, err
	}

	// get pod label keys from cache if cache exists
	if existRes == 1 {
		keys, err := c.keyCache.GetDataList(kit, redisKey)
		if err != nil {
			blog.Errorf("get pod label keys from cache %s failed, err: %v, rid: %s", redisKey, err, kit.Rid)
			return nil, err
		}
		return keys, nil
	}

	// get pod label keys from db and refresh the cache
	refreshOpt := &RefreshPodLabelOpt{
		BizID:      bizID,
		ReturnType: LabelKeyReturnType,
	}
	keys, err := c.RefreshPodLabel(kit, refreshOpt)
	if err != nil {
		blog.Errorf("refresh biz: %d pod label cache failed, err: %v, rid: %s", bizID, err, kit.Rid)
		if len(keys) > 0 {
			// do not return error if keys are acquired but cache update failed
			return keys, nil
		}
		return nil, err
	}

	return keys, nil
}

// GetValues get biz pod label values for specified key
func (c *PodLabelCache) GetValues(kit *rest.Kit, bizID int64, key string) ([]string, error) {
	existRes, err := redis.Client().Exists(kit.Ctx, c.keyCache.key.Key(kit.TenantID, c.genKeyRedisKey(bizID))).Result()
	if err != nil {
		blog.Errorf("check if biz %d pod label cache exists failed, err: %v, rid: %s", bizID, err, kit.Rid)
		return nil, err
	}

	// get pod label values from cache if cache exists
	if existRes == 1 {
		values, err := c.valueCache.GetDataList(kit, c.genValueRedisKey(bizID, key))
		if err != nil {
			blog.Errorf("get biz %d pod label key %s values from cache failed, err: %v, rid: %s", bizID, key, err,
				kit.Rid)
			return nil, err
		}
		return values, nil
	}

	// get pod label values from db and refresh the cache
	refreshOpt := &RefreshPodLabelOpt{
		BizID:      bizID,
		ReturnType: LabelValueReturnType,
		LabelKey:   key,
	}
	values, err := c.RefreshPodLabel(kit, refreshOpt)
	if err != nil {
		blog.Errorf("refresh biz: %d pod label cache failed, err: %v, rid: %s", bizID, err, kit.Rid)
		if len(values) > 0 {
			// do not return error if values are acquired but cache update failed
			return values, nil
		}
		return nil, err
	}

	return values, nil
}

// UpdateKeyCount update pod label key count cache by map[bizID]map[labelKey]count
func (c *PodLabelCache) UpdateKeyCount(kit *rest.Kit, keyCntMap map[int64]map[string]int64) error {
	cntMap := make(map[string]map[string]int64)

	for bizID, keyCnt := range keyCntMap {
		for key, cnt := range keyCnt {
			if cnt == 0 {
				delete(keyCnt, key)
			}
		}
		if len(keyCnt) > 0 {
			cntMap[c.genKeyRedisKey(bizID)] = keyCnt
		}
	}

	if len(cntMap) == 0 {
		return nil
	}

	if err := c.keyCache.UpdateCount(kit, cntMap); err != nil {
		blog.Errorf("update pod label count failed, err: %v, count info: %+v, rid: %s", err, cntMap, kit.Rid)
		return err
	}
	return nil
}

// UpdateValueCount update pod label value count cache by map[bizID]map[labelKey]map[labelValue]count
func (c *PodLabelCache) UpdateValueCount(kit *rest.Kit, valueCntMap map[int64]map[string]map[string]int64) error {
	cntMap := make(map[string]map[string]int64)

	for bizID, keyValueCnt := range valueCntMap {
		for key, valueCnt := range keyValueCnt {
			for value, cnt := range valueCnt {
				if cnt == 0 {
					delete(valueCnt, value)
				}
			}

			if len(valueCnt) > 0 {
				cntMap[c.genValueRedisKey(bizID, key)] = valueCnt
			}
		}
	}

	if len(cntMap) == 0 {
		return nil
	}

	if err := c.valueCache.UpdateCount(kit, cntMap); err != nil {
		blog.Errorf("update pod label count failed, err: %v, count info: %+v, rid: %s", err, cntMap, kit.Rid)
		return err
	}
	return nil
}

// RefreshPodLabelOpt is refresh pod label cache options
type RefreshPodLabelOpt struct {
	BizID int64
	// ReturnType is the return data type, returns no data if not set. options: keys, values
	ReturnType string
	// LabelKey is the label key to return its values when ReturnType is 'values'
	LabelKey string
}

const (
	EmptyReturnType      = ""
	LabelKeyReturnType   = "keys"
	LabelValueReturnType = "values"
)

// Validate RefreshPodLabelOpt
func (opt *RefreshPodLabelOpt) Validate() error {
	if opt == nil {
		return errors.New("refresh pod label cache options not set")
	}

	if opt.BizID <= 0 {
		return fmt.Errorf("biz id %d is invalid", opt.BizID)
	}

	switch opt.ReturnType {
	case EmptyReturnType:
	case LabelKeyReturnType:
	case LabelValueReturnType:
		if opt.LabelKey == "" {
			return errors.New("label key is not set")
		}
	default:
		return fmt.Errorf("return type %s is invalid", opt.ReturnType)
	}

	return nil
}

// RefreshPodLabel refresh pod label key and value cache
func (c *PodLabelCache) RefreshPodLabel(kit *rest.Kit, opt *RefreshPodLabelOpt) ([]string, error) {
	if err := opt.Validate(); err != nil {
		return nil, err
	}

	// lock refresh pod label key and value cache operation, returns error if it is already locked
	lockKey := fmt.Sprintf("%s:pod_label_refresh:lock:%d", Namespace, opt.BizID)

	locker := lock.NewLocker(redis.Client())
	locked, err := locker.Lock(lock.StrFormat(lockKey), 5*time.Minute)
	defer locker.Unlock()
	if err != nil {
		blog.Errorf("get %s lock failed, err: %v, rid: %s", lockKey, err, kit.Rid)
		return nil, err
	}

	if !locked {
		blog.Errorf("%s task is already lock, rid: %s", lockKey, kit.Rid)
		return nil, errors.New("there's a same refreshing task running, please retry later")
	}

	kit.Ctx = util.SetDBReadPreference(kit.Ctx, common.SecondaryPreferredMode)

	keyCntMap, keyValueCntMap, err := c.getBizPodLabelCountInfo(kit, opt.BizID)
	if err != nil {
		return nil, err
	}

	// generate return data
	results := make([]string, 0)
	switch opt.ReturnType {
	case LabelKeyReturnType:
		for key := range keyCntMap {
			results = append(results, key)
		}
	case LabelValueReturnType:
		valueCntMap, exists := keyValueCntMap[opt.LabelKey]
		if !exists {
			break
		}
		for value := range valueCntMap {
			results = append(results, value)
		}
	}

	// refresh label key and value count cache
	delLabelKeys, err := c.keyCache.RefreshCount(kit, c.genKeyRedisKey(opt.BizID), keyCntMap)
	if err != nil {
		blog.Errorf("refresh pod label key count failed, err: %v, count info: %+v, rid: %s", err, keyCntMap, kit.Rid)
		return results, err
	}

	for key, valueCntMap := range keyValueCntMap {
		if _, err = c.valueCache.RefreshCount(kit, c.genValueRedisKey(opt.BizID, key), valueCntMap); err != nil {
			blog.Errorf("refresh pod label key count failed, err: %v, count info: %+v, rid: %s", err, keyCntMap,
				kit.Rid)
			return results, err
		}
	}

	for _, key := range delLabelKeys {
		if err = c.valueCache.Delete(kit, c.genValueRedisKey(opt.BizID, key)); err != nil {
			return nil, err
		}
	}

	return results, nil
}

// getBizPodLabelCountInfo generate map[label_key]count & map[label_key]map[label_value]count by biz pods
func (c *PodLabelCache) getBizPodLabelCountInfo(kit *rest.Kit, bizID int64) (map[string]int64,
	map[string]map[string]int64, error) {

	cond, err := tools.GenKubeSharedNsCond(kit, bizID, kubetypes.BKNamespaceIDField)
	if err != nil {
		blog.Errorf("generate shared namespace condition failed, err: %v, biz: %d, rid: %v", err, bizID, kit.Rid)
		return nil, nil, err
	}

	keyCntMap := make(map[string]int64)
	keyValueCntMap := make(map[string]map[string]int64)

	for {
		pods := make([]kubetypes.Pod, 0)

		err = mongodb.Shard(kit.ShardOpts()).Table(kubetypes.BKTableNameBasePod).Find(cond).Fields(kubetypes.BKIDField,
			kubetypes.LabelsField).Sort(kubetypes.BKIDField).Limit(types.DBPage).All(kit.Ctx, &pods)
		if err != nil {
			blog.Errorf("get pods to refresh label cache failed, cond: %+v, err: %v, rid: %s", cond, err, kit.Rid)
			return nil, nil, err
		}

		for _, pod := range pods {
			if pod.Labels == nil || len(*pod.Labels) == 0 {
				continue
			}

			for key, value := range *pod.Labels {
				key = conv.DecodeDot(key)
				keyCntMap[key]++
				_, exists := keyValueCntMap[key]
				if !exists {
					keyValueCntMap[key] = make(map[string]int64)
				}
				keyValueCntMap[key][value]++
			}
		}

		if len(pods) < types.DBPage {
			break
		}

		cond[kubetypes.BKIDField] = mapstr.MapStr{common.BKDBGT: pods[len(pods)-1].ID}
	}
	return keyCntMap, keyValueCntMap, nil
}

// loopRefreshCache loop refresh pod label key and value cache every day at 3am
func (c *PodLabelCache) loopRefreshCache() {
	lastRefreshTime := -1

	for {
		time.Sleep(time.Hour)

		now := time.Now().Local()
		if lastRefreshTime == now.Day() {
			continue
		}

		if now.Hour() != 3 {
			continue
		}

		for {
			if !c.isMaster.IsMaster() {
				blog.V(4).Infof("loop refresh pod label cache, but not master, skip.")
				time.Sleep(time.Minute)
				continue
			}

			rid := util.GenerateRID()
			blog.Infof("start loop refresh pod label cache task, rid: %s", rid)
			c.RefreshCache(rid)
			lastRefreshTime = now.Day()
			break
		}
	}
}

// RefreshCache loop refresh pod label cache for all bizs
func (c *PodLabelCache) RefreshCache(rid string) {
	kit := rest.NewKit().WithCtx(util.SetDBReadPreference(context.Background(), common.SecondaryPreferredMode)).
		WithRid(rid)

	_ = tenant.ExecForAllTenants(func(tenantID string) error {
		kit = kit.WithTenant(tenantID)

		bizIDs, err := c.getAllBizID(kit)
		if err != nil {
			blog.Errorf("list all biz id for refresh pod label cache task failed, err: %v, rid: %s", err, rid)
			return nil
		}

		for _, bizID := range bizIDs {
			time.Sleep(100 * time.Millisecond)

			kit = kit.WithRid(fmt.Sprintf("%s:%d", rid, bizID))

			blog.Infof("start refresh biz %d pod label cache task, rid: %s", bizID, rid)

			refreshOpt := &RefreshPodLabelOpt{BizID: bizID}
			_, err = c.RefreshPodLabel(kit, refreshOpt)
			if err != nil {
				blog.Errorf("refresh biz %d pod label cache task failed, err: %v, rid: %s", bizID, err, rid)
				continue
			}

			blog.Infof("refresh biz %d pod label cache task successfully, rid: %s", bizID, rid)
		}
		return nil
	})
}

// getAllBizID get all biz id
func (c *PodLabelCache) getAllBizID(kit *rest.Kit) ([]int64, error) {
	cond := mapstr.MapStr{
		common.BKDefaultField:    mapstr.MapStr{common.BKDBNE: common.DefaultAppFlag},
		common.BKDataStatusField: mapstr.MapStr{common.BKDBNE: common.DataStatusDisabled},
	}

	bizIDs := make([]int64, 0)

	for {
		bizs := make([]metadata.BizInst, 0)
		err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseApp).Find(cond).Fields(common.BKAppIDField).
			Limit(types.DBPage).Sort(common.BKAppIDField).All(kit.Ctx, &bizs)
		if err != nil {
			return nil, err
		}

		for _, biz := range bizs {
			bizIDs = append(bizIDs, biz.BizID)
		}

		if len(bizs) < types.DBPage {
			break
		}

		cond[common.BKAppIDField] = mapstr.MapStr{common.BKDBGT: bizs[len(bizs)-1].BizID}
	}

	return bizIDs, nil
}
