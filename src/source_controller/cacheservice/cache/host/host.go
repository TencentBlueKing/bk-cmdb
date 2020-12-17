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
	"errors"
	"fmt"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/driver/redis"
	"configcenter/src/storage/reflector"
	"configcenter/src/storage/stream/types"

	rawRedis "github.com/go-redis/redis/v7"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type hostCache struct {
	key   hostKeyGenerator
	event reflector.Interface
}

func (h *hostCache) Run() error {

	opts := types.Options{
		EventStruct: new(metadata.HostMapStr),
		Collection:  common.BKTableNameBaseHost,
	}

	_, err := redis.Client().Get(context.Background(), h.key.ListDoneKey()).Result()
	if err != nil {
		if !redis.IsNilErr(err) {
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

	// get host details from db again to avoid dirty data.
	ips, cloudID, detail, err := getHostDetailsFromMongoWithHostID(hostID)
	if err != nil {
		blog.Errorf("received host %d upsert event, but get detail from mongodb failed, err: %v", hostID, err)
		return
	}
	refreshHostDetailCache(hostID, ips, cloudID, detail)
}

func (h *hostCache) onDelete(e *types.Event) {
	blog.Infof("received host delete event, oid: %s", e.Oid)

	filter := mapstr.MapStr{
		"oid":  e.Oid,
		"coll": common.BKTableNameBaseHost,
	}
	doc := bsonx.Doc{}
	err := mongodb.Client().Table(common.BKTableNameDelArchive).Find(filter).One(context.Background(), &doc)
	if err != nil {
		blog.Errorf("received delete host event, but get archive deleted doc from mongodb failed, oid: %s, err: %v", e.Oid, err)
		return
	}

	byt, err := bson.MarshalExtJSON(doc.Lookup("detail"), false, false)
	if err != nil {
		blog.Errorf("received delete host event, but marshal doc to bytes failed, oid: %s, err: %v", e.Oid, err)
		return
	}

	elements := gjson.GetManyBytes(byt, common.BKCloudIDField, common.BKHostInnerIPField, common.BKHostIDField)
	cloudID := elements[0].Int()
	ips := elements[1].Array()
	hostID := elements[2].Int()

	pipe := redis.Client().Pipeline()
	// delete cloud id and ip pair
	for _, ip := range ips {
		pipe.Del(h.key.IPCloudIDKey(ip.String(), cloudID))
	}

	// delete host details
	pipe.Del(h.key.HostDetailKey(hostID))
	// remove host id from host id list.
	pipe.ZRem(h.key.HostIDListKey(), hostID)

	_, err = pipe.Exec()
	if err != nil {
		blog.Errorf("received host delete event, oid: %s, but delete oid and detail failed, err: %v", e.Oid, err)
		return
	}
	blog.Infof("received host delete event, oid: %s, delete oid and host id: %d, ip: %s detail success", e.Oid, hostID, ips)
}

func (h *hostCache) onListDone() {
	if err := redis.Client().Set(context.Background(), h.key.ListDoneKey(), "done", 0).Err(); err != nil {
		blog.Errorf("list host data to cache and list done, but set list done key failed, err: %v", err)
		return
	}
	blog.Info("list host data to cache and list done")
}

// refreshHostDetailCache refresh the host's detail cache
func refreshHostDetailCache(hostID int64, ips string, cloudID int64, hostDetail []byte) {
	// get refresh lock to avoid concurrent
	success, err := redis.Client().SetNX(context.Background(), hostKey.HostDetailLockKey(hostID), 1, 10*time.Second).Result()
	if err != nil {
		blog.Errorf("upsert host: %d %s detail cache, but got redis lock failed, err: %v", hostID, ips, err)
		return
	}

	if !success {
		blog.V(4).Infof("upsert host: %d %s detail cache, but do not get redis lock. skip", hostID, ips)
		return
	}

	defer func() {
		if err := redis.Client().Del(context.Background(), hostKey.HostDetailLockKey(hostID)).Err(); err != nil {
			blog.Errorf("upsert host: %d %s detail cache, but delete redis lock failed, err: %v", hostID, ips, err)
		}
	}()

	// we have get the lock, and now we can refresh the cache.
	pipeline := redis.Client().Pipeline()

	// upsert host ip and cloud id relation
	// a host can have multiple host inner ips
	ttl := hostKey.WithRandomExpireSeconds()
	for _, ip := range strings.Split(ips, ",") {
		pipeline.Set(hostKey.IPCloudIDKey(ip, cloudID), hostID, ttl)
	}

	// update host details
	pipeline.Set(hostKey.HostDetailKey(hostID), hostDetail, ttl)

	// add host id to id list.
	pipeline.ZAddNX(hostKey.HostIDListKey(), &rawRedis.Z{
		// set host id as it's score number
		Score:  float64(hostID),
		Member: hostID,
	})

	_, err = pipeline.Exec()
	if err != nil {
		blog.Errorf("upsert host: %d, ip: %s cache, but upsert to redis failed, err: %v", hostID, ips, err)
		return
	}
	blog.V(4).Infof("refresh host cache success, host id: %d, ips: %s, ttl: %ds", hostID, ips, ttl/time.Second)
}

func getHostDetailsFromMongoWithHostID(hostID int64) (ips string, cloudID int64, detail []byte, err error) {
	filter := mapstr.MapStr{
		common.BKHostIDField: hostID,
	}
	host := make(metadata.HostMapStr)
	err = mongodb.Client().Table(common.BKTableNameBaseHost).Find(filter).One(context.Background(), &host)
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

type hostBase struct {
	id      int64
	ip      string
	cloudID int64
	detail  string
}

func listHostDetailsFromMongoWithHostID(hostID []int64) (list []*hostBase, err error) {
	filter := mapstr.MapStr{
		common.BKHostIDField: mapstr.MapStr{
			common.BKDBIN: hostID,
		},
	}
	host := make([]metadata.HostMapStr, 0)
	err = mongodb.Client().Table(common.BKTableNameBaseHost).Find(filter).Sort(common.BKHostIDField).All(context.Background(), &host)
	if err != nil {
		blog.Errorf("get host data from mongodb for cache failed, err: %v", err)
		return nil, err
	}

	for _, h := range host {
		ips, ok := h[common.BKHostInnerIPField].(string)
		if !ok {
			blog.Errorf("get host: %v data from mongodb for cache, but got invalid ip, host: %v", hostID, h)
			return nil, errors.New("invalid host innerip")
		}

		js, _ := json.Marshal(h)
		ele := gjson.GetManyBytes(js, common.BKCloudIDField, common.BKHostIDField)
		if !ele[0].Exists() {
			blog.Errorf("get host from mongodb for cache, but cloud id not exist, host: %v", h)
			return nil, errors.New("host cloud id not exist")
		}
		if !ele[1].Exists() {
			blog.Errorf("get host from mongodb for cache, but host id not exist, host: %v", h)
			return nil, errors.New("host id not exist")
		}

		id := ele[1].Int()
		if id == 0 {
			blog.Errorf("get host from mongodb for cache, but host id is 0, host: %v", h)
			return nil, errors.New("host id is 0")
		}

		list = append(list, &hostBase{
			id:      id,
			ip:      ips,
			cloudID: ele[0].Int(),
			detail:  string(js),
		})
	}
	return list, nil
}

func getHostDetailsFromMongoWithIP(innerIP string, cloudID int64) (hostID int64, detail []byte, err error) {
	innerIPArr := strings.Split(innerIP, ",")
	filter := mapstr.MapStr{
		common.BKHostInnerIPField: map[string]interface{}{
			common.BKDBAll: innerIPArr,
		},
		common.BKCloudIDField: cloudID,
	}
	host := make(metadata.HostMapStr)
	err = mongodb.Client().Table(common.BKTableNameBaseHost).Find(filter).One(context.Background(), &host)
	if err != nil {
		blog.Errorf("get host data from mongodb with ip: %s, cloud: %d for cache failed, err: %v", innerIP, cloudID, err)
		return 0, nil, err
	}

	id, err := util.GetInt64ByInterface(host[common.BKHostIDField])
	if err != nil {
		return 0, nil, fmt.Errorf("get host data from mongodb ip: %s, cloud: %d for cache, but got invalid host id: %v, err: %v",
			innerIP, cloudID, host[common.BKHostIDField], err)
	}

	js, _ := json.Marshal(host)
	return id, js, nil
}
