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

// Package host TODO
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
)

type hostCache struct {
	key   hostKeyGenerator
	event reflector.Interface
}

// Run TODO
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
	if blog.V(4) {
		blog.Infof("received host upsert event, oid: %s, doc: %s", e.Oid, e.DocBytes)
	}

	hostID := gjson.GetBytes(e.DocBytes, "bk_host_id").Int()
	if hostID <= 0 {
		blog.Errorf("received host upsert event, but got invalid host id, doc: %s", e.DocBytes)
		return
	}

	// event refresh cache does not need to generate rid.
	rid := ""
	// get host details from db again to avoid dirty data.
	host, err := getHostDetailsFromMongoWithHostID(rid, hostID)
	if err != nil {
		blog.Errorf("received host %d upsert event, but get detail from mongodb failed, err: %v", hostID, err)
		return
	}
	refreshHostDetailCache(rid, host)
}

func (h *hostCache) onDelete(e *types.Event) {
	blog.Infof("received host delete event, oid: %s", e.Oid)

	filter := mapstr.MapStr{
		"oid":  e.Oid,
		"coll": common.BKTableNameBaseHost,
	}
	doc := make(map[string]interface{})
	err := mongodb.Client().Table(common.BKTableNameDelArchive).Find(filter).One(context.Background(), &doc)
	if err != nil {
		blog.Errorf("received delete host event, but get archive deleted doc from mongodb failed, oid: %s, err: %v", e.Oid, err)
		return
	}

	byt, err := json.Marshal(doc["detail"])
	if err != nil {
		blog.Errorf("received delete host event, but marshal doc to bytes failed, oid: %s, err: %v", e.Oid, err)
		return
	}

	elements := gjson.GetManyBytes(byt, common.BKCloudIDField, common.BKHostInnerIPField, common.BKHostIDField,
		common.BKAgentIDField)
	cloudID := elements[0].Int()
	ips := elements[1].Array()
	hostID := elements[2].Int()
	agentID := elements[3].String()
	pipe := redis.Client().Pipeline()
	// delete cloud id and ip pair
	for _, ip := range ips {
		pipe.Del(h.key.IPCloudIDKey(ip.String(), cloudID))
	}
	if agentID != "" {
		pipe.Del(h.key.AgentIDKey(agentID))
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
func refreshHostDetailCache(rid string, host *hostBase) {

	// If there is an agentID in the static ip scenario, the agentID will be used as the key first. If there is no
	// agentID, it must be a static IP scenario to be updated with ip+cloud; if it is a dynamic ip and If the agentID is
	// empty, there is a problem with the theoretical data, and no cache update will be done.
	if host.agentID == "" && host.addressType != common.BKAddressingStatic {
		return
	}

	hostDetailKey := hostKey.HostDetailLockKey(host.hostID)
	// get refresh lock to avoid concurrent
	success, err := redis.Client().SetNX(context.Background(), hostDetailKey, 1, 10*time.Second).Result()
	if err != nil {
		blog.Errorf("upsert hostID: %d, ip: %s, agentID: %s, detail cache, but got redis lock failed, err: %v, rid: %s",
			host.hostID, host.ip, host.agentID, err, rid)
		return
	}

	if !success {
		blog.V(4).Infof("upsert hostID: %d, ip: %s, agentID: %s, rid: %s detail cache, but do not get redis lock. skip",
			host.hostID, host.ip, host.agentID, rid)
		return
	}

	defer func() {
		if err := redis.Client().Del(context.Background(), hostDetailKey).Err(); err != nil {
			blog.Errorf("upsert host id: %d, ip: %s, agentID: %s, detail cache, but delete redis lock failed, "+
				"err: %v, rid: %s", host.hostID, host.ip, host.agentID, err, rid)
		}
	}()

	// we have get the lock, and now we can refresh the cache.
	pipeline := redis.Client().Pipeline()

	// upsert host ip and cloud id relation
	// a host can have multiple host inner ips
	ttl := hostKey.WithRandomExpireSeconds()
	if host.agentID != "" {
		pipeline.Set(hostKey.AgentIDKey(host.agentID), host.hostID, ttl)
	}

	if host.addressType == common.BKAddressingStatic && host.ip != "" {
		for _, ip := range strings.Split(host.ip, ",") {
			pipeline.Set(hostKey.IPCloudIDKey(ip, host.cloudID), host.hostID, ttl)
		}
	}

	// update host details
	pipeline.Set(hostKey.HostDetailKey(host.hostID), host.detail, ttl)

	// add host id to id list.
	pipeline.ZAddNX(hostKey.HostIDListKey(), &rawRedis.Z{
		// set host id as it's score number
		Score:  float64(host.hostID),
		Member: host.hostID,
	})

	_, err = pipeline.Exec()
	if err != nil {
		blog.Errorf("upsert hostID: %d, ip: %s, agentID: %s, cache, but upsert to redis failed, err: %v, rid: %s",
			host.hostID, host.ip, host.agentID, err, rid)
		return
	}
	if blog.V(4) {
		blog.Infof("refresh host cache success, host id: %d, ips: %s, agentID: %s, ttl: %ds, rid: %s",
			host.hostID, host.ip, host.agentID, ttl/time.Second, rid)
	}
}

func getHostDetailsFromMongoWithHostID(rid string, hostID int64) (*hostBase, error) {
	filter := mapstr.MapStr{
		common.BKHostIDField: hostID,
	}
	host := make(metadata.HostMapStr)
	err := mongodb.Client().Table(common.BKTableNameBaseHost).Find(filter).One(context.Background(), &host)
	if err != nil {
		blog.Errorf("get host data from mongodb for cache failed, err: %v,rid: %s", err, rid)
		return nil, err
	}

	ips := util.GetStrByInterface(host[common.BKHostInnerIPField])
	agentID := ""
	if host[common.BKAgentIDField] != nil {
		id, ok := host[common.BKAgentIDField].(string)
		if !ok {
			blog.Errorf("get host: %d data from mongodb for cache, but got invalid agentID, host: %v, rid: %s",
				hostID, host, rid)
			return nil, fmt.Errorf("invalid host: %d innerip", hostID)
		}
		agentID = id
	}

	addressType := common.BKAddressingStatic
	if host[common.BKAddressingField] != nil {
		field, ok := host[common.BKAddressingField].(string)
		if !ok {
			return nil, errors.New("bk_addressing type error")
		}
		addressType = field
	}
	js, _ := json.Marshal(host)

	ele := gjson.GetBytes(js, common.BKCloudIDField)
	if !ele.Exists() {
		blog.Errorf("get host: %d data from mongodb for cache, but cloud id not exist, host: %v, rid: %s",
			hostID, host, rid)
		return nil, fmt.Errorf("host %d cloud id not exist", hostID)
	}
	result := &hostBase{
		ip:          ips,
		hostID:      hostID,
		cloudID:     ele.Int(),
		agentID:     agentID,
		detail:      string(js),
		addressType: addressType,
	}
	return result, nil
}

// hostBase host base info
type hostBase struct {
	hostID      int64
	ip          string
	agentID     string
	cloudID     int64
	detail      string
	addressType string
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
		ips := util.GetStrByInterface(h[common.BKHostInnerIPField])
		js, _ := json.Marshal(h)
		ele := gjson.GetManyBytes(js, common.BKCloudIDField, common.BKHostIDField, common.BKAddressingField)
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
		addressType := ele[2].String()

		list = append(list, &hostBase{
			hostID:      id,
			ip:          ips,
			cloudID:     ele[0].Int(),
			detail:      string(js),
			addressType: addressType,
		})
	}
	return list, nil
}

// getHostDetailsFromMongoWithIP In the static IP scenario, this function is allowed to query, because the IP can be
// used as a unique index in this scenario.
func getHostDetailsFromMongoWithIP(innerIP string, cloudID int64) (int64, []byte, error) {
	innerIPArr := strings.Split(innerIP, ",")
	filter := mapstr.MapStr{
		common.BKHostInnerIPField: map[string]interface{}{
			common.BKDBAll: innerIPArr,
		},
		common.BKCloudIDField:    cloudID,
		common.BKAddressingField: common.BKAddressingStatic,
	}
	host := make(metadata.HostMapStr)
	err := mongodb.Client().Table(common.BKTableNameBaseHost).Find(filter).One(context.Background(), &host)
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

func getHostDetailsFromMongoWithAgentID(rid string, agentID string) (int64, string, []byte, error) {
	filter := mapstr.MapStr{
		common.BKAgentIDField: agentID,
	}
	host := make(metadata.HostMapStr)
	err := mongodb.Client().Table(common.BKTableNameBaseHost).Find(filter).One(context.Background(), &host)
	if err != nil {
		blog.Errorf("get host data from mongodb with agentID: %s failed, err: %v, rid: %v", agentID, err, rid)
		return 0, "", nil, err
	}
	id, err := util.GetInt64ByInterface(host[common.BKHostIDField])
	if err != nil {
		return 0, "", nil, fmt.Errorf("get host data from mongodb with agentID: %s failed, host: %+v, err: %v, rid: %s",
			agentID, host, err, rid)
	}

	addressType := common.BKAddressingStatic
	if host[common.BKAddressingField] != nil {
		if field, ok := host[common.BKAddressingField].(string); ok {
			addressType = field
		}
	}
	js, _ := json.Marshal(host)
	return id, addressType, js, nil
}
