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

package host

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	"configcenter/src/common/mapstr"
	params "configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/reflector"
	"configcenter/src/storage/stream/types"
	"github.com/tidwall/gjson"
	"gopkg.in/redis.v5"
)

type hostCache struct {
	key   hostKeyGenerator
	rds   *redis.Client
	event reflector.Interface
	db    dal.DB
}

func (h *hostCache) Run() error {

	opts := types.Options{
		EventStruct: new(map[string]interface{}),
		Collection:  common.BKTableNameBaseHost,
	}

	_, err := h.rds.Get(h.key.ListDoneKey()).Result()
	if err != nil {
		if err != redis.Nil {
			blog.Errorf("get host list done redis key failed, err: %v", err)
			return fmt.Errorf("get host list done redis key failed, err: %v", err)
		}
		listCap := &reflector.Capable{
			OnChange: reflector.OnChangeEvent{
				OnLister:     h.onUpsert,
				OnAdd:        h.onUpsert,
				OnUpdate:     h.onUpsert,
				OnListerDone: h.onListDone,
				OnDelete:     h.onDelete,
			},
		}
		// do with list watcher.
		page := 500
		listOpts := &types.ListWatchOptions{
			Options:  opts,
			PageSize: &page,
		}
		blog.Info("do host cache with list watcher.")
		return h.event.ListWatcher(context.Background(), listOpts, listCap)
	}

	watchCap := &reflector.Capable{
		OnChange: reflector.OnChangeEvent{
			OnAdd:    h.onUpsert,
			OnUpdate: h.onUpsert,
			OnDelete: h.onDelete,
		},
	}
	// do with watcher only.
	watchOpts := &types.WatchOptions{
		Options: opts,
	}
	blog.Info("do host cache with only watcher")
	return h.event.Watcher(context.Background(), watchOpts, watchCap)
}

func (h *hostCache) onUpsert(e *types.Event) {
	blog.V(4).Infof("received host upsert event, oid: %s, doc: %s", e.Oid, e.DocBytes)

	hostID := gjson.GetBytes(e.DocBytes, "bk_host_id").Int()
	if hostID <= 0 {
		blog.Errorf("received host upsert event, but got invalid host id, doc: %s", e.DocBytes)
		return
	}

	h.upsertOid(e.Oid, hostID)
	// get host details from db again to avoid dirty data.
	ips, cloudID, detail, err := getHostDetailsFromMongoWithHostID(h.db, hostID)
	if err != nil {
		blog.Errorf("received host %d upsert event, but get detail from mongodb failed, err: %v", hostID, err)
		return
	}
	refreshHostDetailCache(h.rds, hostID, ips, cloudID, detail)
}

func (h *hostCache) onDelete(e *types.Event) {
	blog.Infof("received host delete event, oid: %s", e.Oid)
	// get host id with oid
	hostIDStr, err := h.rds.HGet(h.key.HostOidKey(), e.Oid).Result()
	if err != nil {
		blog.Errorf("received host delete event, oid: %s, but get host id failed, err: %v", e.Oid, err)
		return
	}

	hostID, err := strconv.ParseInt(hostIDStr, 64, 10)
	if err != nil {
		blog.Errorf("received host delete event, oid: %s, but parse host id failed, err: %v", e.Oid, err)
		return
	}
	host, err := h.rds.Get(h.key.HostDetailKey(hostID)).Result()
	if err != nil {
		blog.Errorf("received host delete event, oid: %s, but get host cache detail failed, err: %v", e.Oid, err)
		return
	}
	elements := gjson.GetMany(host, common.BKCloudIDField, common.BKHostInnerIPField)
	cloudID := elements[0].Int()
	ips := elements[1].String()

	pipe := h.rds.Pipeline()
	// delete cloud id and ip pair
	for _, ip := range strings.Split(ips, ",") {
		pipe.Del(h.key.IPCloudIDKey(ip, cloudID))
	}

	// delete oid relation
	pipe.HDel(h.key.HostOidKey(), e.Oid)

	// delete host details
	pipe.Del(h.key.HostDetailKey(hostID))
	_, err = pipe.Exec()
	if err != nil {
		blog.Errorf("received host delete event, oid: %s, but delete oid and detail failed, err: %v", e.Oid, err)
		return
	}
	blog.Infof("received host delete event, oid: %s, delete oid and host id: %d detail success", e.Oid, hostID)
}

func (h *hostCache) onListDone() {
	if err := h.rds.Set(h.key.ListDoneKey(), "done", 0).Err(); err != nil {
		blog.Errorf("list host data to cache and list done, but set list done key failed, err: %v", err)
		return
	}
	blog.Info("list host data to cache and list done")
}

func (h *hostCache) upsertOid(oid string, hostID int64) {
	if err := h.rds.HSet(h.key.HostOidKey(), oid, hostID).Err(); err != nil {
		blog.Errorf("upsert host: %d cache oid key: %s, but hset failed, err: %v", hostID, h.key.HostOidKey(), err)
		return
	}
}

func (h *hostCache) deleteOid(oid string) {
	if err := h.rds.HDel(h.key.HostOidKey(), oid).Err(); err != nil {
		blog.Errorf("delete host oid: %s cache oid key: %s, but hdel failed, err: %v", oid, h.key.HostOidKey(), err)
		return
	}
}

// refreshHostDetailCache refresh the host's detail cache
func refreshHostDetailCache(rds *redis.Client, hostID int64, ips string, cloudID int64, hostDetail []byte) {
	// get refresh lock to avoid concurrent
	success, err := rds.SetNX(hostKey.HostDetailLockKey(hostID), 1, 10*time.Second).Result()
	if err != nil {
		blog.Errorf("upsert host: %d %s detail cache, but got redis lock failed, err: %v", hostID, ips, err)
		return
	}

	if !success {
		blog.V(4).Infof("upsert host: %d %s detail cache, but do not get redis lock. skip", hostID, ips)
		return
	}

	defer func() {
		if err := rds.Del(hostKey.HostDetailLockKey(hostID)).Err(); err != nil {
			blog.Errorf("upsert host: %d %s detail cache, but delete redis lock failed, err: %v", hostID, ips, err)
		}
	}()

	// we have get the lock, and now we can refresh the cache.
	pipeline := rds.Pipeline()
	// upsert host ip and cloud id relation
	// a host can have multiple host inner ips
	ttl := hostKey.WithRandomExpireSeconds()
	for _, ip := range strings.Split(ips, ",") {
		pipeline.Set(hostKey.IPCloudIDKey(ip, cloudID), hostID, ttl)
	}

	// update host details
	pipeline.Set(hostKey.HostDetailKey(hostID), hostDetail, ttl)
	_, err = pipeline.Exec()
	if err != nil {
		blog.Errorf("upsert host: %d, ip: %s cache, but upsert to redis failed, err: %v", hostID, ips, err)
		return
	}
	blog.V(4).Infof("refresh host cache success, host id: %d, ips: %s, ttl: %ds", hostID, ips, ttl/time.Second)
}

func getHostDetailsFromMongoWithHostID(db dal.DB, hostID int64) (ips string, cloudID int64, detail []byte, err error) {
	filter := mapstr.MapStr{
		common.BKHostIDField: hostID,
	}
	host := make(map[string]interface{})
	err = db.Table(common.BKTableNameBaseHost).Find(filter).One(context.Background(), &host)
	if err != nil {
		blog.Errorf("get host data from mongodb for cache failed, err: %v", err)
		return "", 0, nil, err
	}

	var ok bool
	ips, ok = host[common.BKHostInnerIPField].(string)
	if !ok {
		blog.Errorf("get host: %d data from mongodb for cache, but got invalid ip, host: %v", hostID, host)
		return "", 0, nil, fmt.Errorf("invalid host: %d innerip", hostID)
	}

	js, _ := json.Marshal(host)

	ele := gjson.GetBytes(js, common.BKCloudIDField)
	if !ele.Exists() {
		blog.Errorf("get host: %d data from mongodb for cache, but cloud id not exist, host: %v", hostID, host)
		return "", 0, nil, fmt.Errorf("host %d cloud id not exist", hostID)
	}
	return ips, ele.Int(), js, nil
}

const exactIPRegexp = `(^IP_PLACEHOLDER$)|(^IP_PLACEHOLDER[,]{1})|([,]{1}IP_PLACEHOLDER[,]{1})|([,]{1}IP_PLACEHOLDER$)`

func getHostDetailsFromMongoWithIP(db dal.DB, innerIP string, cloudID int64) (hostID int64, detail []byte, err error) {
	filter := mapstr.MapStr{
		common.BKHostInnerIPField: mapstr.MapStr{
			common.BKDBLIKE: strings.Replace(exactIPRegexp, "IP_PLACEHOLDER", params.SpecialCharChange(innerIP), -1),
		},
		common.BKCloudIDField: cloudID,
	}
	host := make(map[string]interface{})
	err = db.Table(common.BKTableNameBaseHost).Find(filter).One(context.Background(), &host)
	if err != nil {
		blog.Errorf("get host data from mongodb with ip: %s, cloud: %d for cache failed, err: %v", innerIP, cloudID, err)
		return 0, nil, err
	}

	id, err := util.GetInt64ByInterface(host[common.BKHostIDField])
	if err != nil {
		return 0, nil, fmt.Errorf("get host data from mongodb ip: %s, cloud: %d for cache, but got invalid host id: %v, err: %v",
			innerIP, hostID, host, err)
	}

	js, _ := json.Marshal(host)
	return id, js, nil
}
