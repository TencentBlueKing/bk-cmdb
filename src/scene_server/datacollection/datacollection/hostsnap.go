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

package datacollection

import (
	"fmt"
	"io"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/xid"
	"github.com/tidwall/gjson"
	redis "gopkg.in/redis.v5"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/storage"
)

var (
	getMasterProcIntervalTime = time.Second * 10

	masterProcLockLiveTime       = getMasterProcIntervalTime + time.Second*10
	clearMsgChanTime       int64 = 120
	masterProcLockContent  string

	fetchDBInterval = time.Minute * 10
	maxconcurrent   = runtime.NumCPU()
)

type HostSnap struct {
	id           string
	hostChanName []string
	msgChan      chan string
	maxSize      int
	ts           time.Time
	lastMesgTs   time.Time
	isMaster     bool
	sync.Mutex
	interrupt     chan error
	maxconcurrent int

	doneCh chan struct{}

	resetHandle chan struct{}

	redisCli *redis.Client
	snapCli  *redis.Client

	subscribing bool

	cache *Cache
	db    storage.DI

	wg *sync.WaitGroup
}

type Cache struct {
	cache map[bool]*HostCache
	flag  bool
}

func NewHostSnap(chanName []string, maxSize int, redisCli, snapCli *redis.Client, db storage.DI) *HostSnap {
	if 0 == maxSize {
		maxSize = 100
	}
	hostSnapInstance := &HostSnap{
		hostChanName:  chanName,
		msgChan:       make(chan string, maxSize*4),
		interrupt:     make(chan error),
		resetHandle:   make(chan struct{}),
		maxSize:       maxSize,
		redisCli:      redisCli,
		snapCli:       snapCli,
		db:            db,
		ts:            time.Now(),
		id:            xid.New().String()[5:],
		maxconcurrent: maxconcurrent,
		wg:            &sync.WaitGroup{},
		doneCh:        make(chan struct{}),
		cache: &Cache{
			cache: map[bool]*HostCache{},
			flag:  false,
		},
	}
	return hostSnapInstance
}

func (h *HostSnap) Start() {

	go func() {
		h.Run()
		for {
			time.Sleep(time.Second * 10)
			NewHostSnap(h.hostChanName, h.maxSize, h.redisCli, h.snapCli, h.db).Run()
		}
	}()
}

func (h *HostSnap) Run() {
	defer func() {
		syserr := recover()
		if syserr != nil {
			blog.Errorf("emergency error happened %s, we will try again 10s later, stack: \n%s", syserr, debug.Stack())
		}
		close(h.doneCh)
		h.isMaster = false
		return
	}()
	blog.Infof("datacollection start with maxconcurrent: %d", h.maxconcurrent)
	ticker := time.NewTicker(getMasterProcIntervalTime)
	var err error
	var msg string
	var msgs []string
	var addCount int
	var waitCnt int

	go h.fetchDB()

	if h.saveRunning() {
		go h.subChan(h.snapCli, h.hostChanName)
	} else {
		blog.Infof("run: there is other master process exists, recheck after %v ", getMasterProcIntervalTime)
	}
	for {
		select {
		case <-ticker.C:
			if h.saveRunning() {
				if !h.subscribing {
					go h.subChan(h.snapCli, h.hostChanName)
				}
			}
		case msg = <-h.msgChan:

			h.Lock()
			h.ts = time.Now()
			msgs = make([]string, 0, h.maxSize*2)
			timeoutCh := time.After(time.Second)
			msgs = append(msgs, msg)
			addCount = 0
		f:
			for {
				select {
				case <-timeoutCh:
					break f
				case msg = <-h.msgChan:
					addCount++
					msgs = append(msgs, msg)
				}
				if addCount > h.maxSize {
					break f
				}
			}
			h.Unlock()

			waitCnt = 0
			for {
				if waitCnt > h.maxconcurrent*2 {
					blog.Warnf("reset handlers")
					close(h.resetHandle)
					waitCnt = 0
					h.resetHandle = make(chan struct{})
				}
				if atomic.LoadInt64(&routeCnt) < int64(h.maxconcurrent) {
					atomic.AddInt64(&routeCnt, 1)
					go h.handleMsg(msgs, h.resetHandle)
					break
				}
				waitCnt++
				time.Sleep(time.Millisecond * 100)
			}
		case err = <-h.interrupt:
			blog.Warn("interrupted", err.Error())
			h.concede()
		}

	}
}

var routeCnt = int64(0)
var handleCnt = int64(0)
var handlelock = sync.Mutex{}
var handlets = time.Now()

func (h *HostSnap) handleMsg(msgs []string, resetHandle chan struct{}) error {
	defer atomic.AddInt64(&routeCnt, -1)
	blog.Infof("handle %d num mesg, routines %d", len(msgs), atomic.LoadInt64(&routeCnt))
	for index, msg := range msgs {
		if msg == "" {
			continue
		}
		handlelock.Lock()
		handleCnt++
		if handleCnt%10000 == 0 {
			blog.Infof("handle rate: %d/sec", int(float64(handleCnt)/time.Now().Sub(handlets).Seconds()))
			handleCnt = 0
			handlets = time.Now()
		}
		handlelock.Unlock()
		select {
		case <-resetHandle:
			blog.Warnf("reset handler, handled %d, set maxSize to %d ", index, h.maxSize)
			return nil
		case <-h.doneCh:
			blog.Warnf("close handler, handled %d")
			return nil
		default:
			var data = msg
			if !gjson.Get(msg, "cloudid").Exists() {
				data = gjson.Get(msg, "data").String()
			}
			val := gjson.Parse(data)
			host := h.getHostByVal(&val)
			if host == nil {
				blog.Warnf("host not found, continue, %s", val.String())
				continue
			}
			hostid := fmt.Sprint(host.get(common.BKHostIDField))
			if hostid == "" {
				blog.Warnf("host id not found, continue, %s", val.String())
				continue
			}

			if err := h.redisCli.Set(common.RedisSnapKeyPrefix+hostid, data, time.Minute*10).Err(); err != nil {
				blog.Errorf("save snapshot %s to redis faile: %s", common.RedisSnapKeyPrefix+hostid, err.Error())
			}

			condition := map[string]interface{}{common.BKHostIDField: host.get(common.BKHostIDField)}
			innerip, ok := host.get(common.BKHostInnerIPField).(string)
			if !ok {
				blog.Infof("innerip is empty, continue, %s", val.String())
				continue
			}
			outip, ok := host.get(common.BKHostOuterIPField).(string)
			if !ok {
				blog.Warnf("outip is not string, %s", val.String())
			}
			setter := parseSetter(&val, innerip, outip)
			if needToUpdate(setter, host) {
				blog.Infof("update by %v, to %v", condition, setter)
				if err := h.db.UpdateByCondition(common.BKTableNameBaseHost, setter, condition); err != nil {
					blog.Error("update host error:", err.Error())
					continue
				}
				copyVal(setter, host)
			}
		}
	}

	return nil
}

func copyVal(a map[string]interface{}, b *HostInst) {
	for k, v := range a {
		b.set(k, v)
	}
}
func needToUpdate(a map[string]interface{}, b *HostInst) bool {
	for k, v := range a {
		if b.get(k) != v {
			return true
		}
	}
	return false
}

func parseSetter(val *gjson.Result, innerIP, outerIP string) map[string]interface{} {
	var cpumodule = val.Get("data.cpu.cpuinfo.0.modelName").String()
	var cupnum int64
	for _, core := range val.Get("data.cpu.cpuinfo.#.cores").Array() {
		cupnum += core.Int()
	}
	var CPUMhz = val.Get("data.cpu.cpuinfo.0.mhz").Int()
	var disk int64
	for _, disktotal := range val.Get("data.disk.usage.#.total").Array() {
		disk += disktotal.Int()
	}
	var mem = val.Get("data.mem.meminfo.total").Int()
	var hostname = val.Get("data.system.info.hostname").String()
	var ostype = val.Get("data.system.info.os").String()
	var osname string
	platform := val.Get("data.system.info.platform").String()
	version := val.Get("data.system.info.platformVersion").String()
	switch strings.ToLower(ostype) {
	case "linux":
		version = strings.Replace(version, ".x86_64", "", 1)
		version = strings.Replace(version, ".i386", "", 1)
		osname = fmt.Sprintf("%s %s", ostype, platform)
		ostype = common.HostOSTypeEnumLinux
	case "windows":
		version = strings.Replace(version, "Microsoft ", "", 1)
		platform = strings.Replace(platform, "Microsoft ", "", 1)
		osname = fmt.Sprintf("%s", platform)
		ostype = common.HostOSTypeEnumWindows
	default:
		osname = fmt.Sprintf("%s", platform)
	}
	var OuterMAC, InnerMAC string
	for _, inter := range val.Get("data.net.interface").Array() {
		for _, addr := range inter.Get("addrs.#.addr").Array() {
			ip := strings.Split(addr.String(), "/")[0]
			if ip == innerIP {
				InnerMAC = inter.Get("hardwareaddr").String()
			} else if ip == outerIP {
				OuterMAC = inter.Get("hardwareaddr").String()
			}
		}
	}

	osbit := val.Get("data.system.info.systemtype").String()

	setter := map[string]interface{}{
		"bk_cpu":        cupnum,
		"bk_cpu_module": cpumodule,
		"bk_cpu_mhz":    CPUMhz,
		"bk_disk":       disk / 1024 / 1024 / 1024,
		"bk_mem":        mem / 1024 / 1024,
		"bk_os_type":    ostype,
		"bk_os_name":    osname,
		"bk_os_version": version,
		"bk_host_name":  hostname,
		"bk_outer_mac":  OuterMAC,
		"bk_mac":        InnerMAC,
		"bk_os_bit":     osbit,
	}

	if cupnum <= 0 {
		blog.Infof("bk_cpu not found in message for %s", innerIP)
	}
	if cpumodule == "" {
		blog.Infof("bk_cpu_module not found in message for %s", innerIP)
	}
	if CPUMhz <= 0 {
		blog.Infof("bk_cpu_mhz not found in message for %s", innerIP)
	}
	if disk <= 0 {
		blog.Infof("bk_disk not found in message for %s", innerIP)
	}
	if mem <= 0 {
		blog.Infof("bk_mem not found in message for %s", innerIP)
	}
	if ostype == "" {
		blog.Infof("bk_os_type not found in message for %s", innerIP)
	}
	if osname == "" {
		blog.Infof("bk_os_name not found in message for %s", innerIP)
	}
	if version == "" {
		blog.Infof("bk_os_version not found in message for %s", innerIP)
	}
	if hostname == "" {
		blog.Infof("bk_host_name not found in message for %s", innerIP)
	}
	if outerIP != "" && OuterMAC == "" {
		blog.Infof("bk_outer_mac not found in message for %s", innerIP)
	}
	if InnerMAC == "" {
		blog.Infof("bk_mac not found in message for %s", innerIP)
	}

	return setter
}
func getIPS(val *gjson.Result) (ips []string) {
	interfaces := val.Get("data.net.interface.#.addrs.#.addr").Array()
	for _, addrs := range interfaces {
		for _, addr := range addrs.Array() {
			ip := strings.Split(addr.String(), "/")[0]
			if strings.HasPrefix(ip, "127.0.0.") {
				continue
			}
			ips = append(ips, ip)
		}
	}
	if !strings.HasPrefix(val.Get("ip").String(), "127.0.0.") {
		ips = append(ips, val.Get("ip").String())
	}
	return ips
}

func (h *HostSnap) getHostByVal(val *gjson.Result) *HostInst {
	cloudid := val.Get("cloudid").String()
	ownerID := val.Get("bizid").String()

	ips := getIPS(val)
	if len(ips) > 0 {
		blog.Infof("handle clouid: %s ips: %v", cloudid, ips)
		for _, ip := range ips {
			if host := h.getCache().get(cloudid + "::" + ip); host != nil {
				return host
			}
		}

		blog.Infof("ips not in cache clouid: %s,ip: %v", cloudid, ips)
		clouidInt, err := strconv.Atoi(cloudid)
		if nil != err {
			blog.Infof("cloudid \"%s\" not integer", cloudid)
			return nil
		}
		condition := map[string]interface{}{
			common.BKCloudIDField: clouidInt,
			common.BKHostInnerIPField: map[string]interface{}{
				common.BKDBIN: ips,
			},
			common.BKOwnerIDField: ownerID,
		}
		result := []map[string]interface{}{}
		err = h.db.GetMutilByCondition(common.BKTableNameBaseHost, nil, condition, &result, "", 0, 0)
		if err != nil {
			blog.Errorf("fetch db error %v", err)
		}
		for index := range result {
			cloudid := fmt.Sprint(result[index][common.BKCloudIDField])
			innerip := fmt.Sprint(result[index][common.BKHostInnerIPField])
			inst := &HostInst{data: result[index]}
			h.setCache(cloudid+"::"+innerip, inst)
			return inst
		}
		blog.Infof("ips not in cache and db, clouid: %v, ip: %v", cloudid, ips)
	} else {
		blog.Errorf("message has no ip, message:%s", val.String())
	}
	return nil
}

func (h *HostSnap) concede() {
	blog.Info("concede")
	h.isMaster = false
	h.subscribing = false
	val := h.redisCli.Get(MasterProcLockKey).Val()
	if val != h.id {
		h.redisCli.Del(MasterProcLockKey)
	}
}

func (h *HostSnap) saveRunning() (ok bool) {
	var err error
	if h.isMaster {
		var val string
		val, err = h.redisCli.Get(MasterProcLockKey).Result()
		if err != nil {
			blog.Errorf("master: saveRunning err %v", err)
			h.isMaster = false
		} else if val == h.id {
			blog.Infof("master check : i am still master")
			h.redisCli.Set(MasterProcLockKey, h.id, masterProcLockLiveTime)
			ok = true
			h.isMaster = true
		} else {
			blog.Infof("exit master,val = %v, id = %v", val, h.id)
			h.isMaster = false
			ok = false
		}
	} else {
		ok, err = h.redisCli.SetNX(MasterProcLockKey, h.id, masterProcLockLiveTime).Result()
		if err != nil {
			blog.Errorf("slave: saveRunning err %v", err)
			h.isMaster = false
		} else if ok {
			blog.Infof("slave check: ok")
			blog.Infof("i am master from now")
			h.isMaster = true
		} else {
			blog.Infof("slave check: there is other master process exists, recheck after %v ", getMasterProcIntervalTime)
			h.isMaster = false
		}
	}
	return ok
}

func (h *HostSnap) subChan(snapcli *redis.Client, chanName []string) {
	defer func() {
		syserr := recover()
		if syserr != nil {
			blog.Errorf("subChan emergency error happened %s, we will try again 10s later, stack: \n%s", syserr, debug.Stack())
		}
		h.subscribing = false
	}()
	h.subscribing = true
	var chanlen int
	subChan, err := snapcli.Subscribe(chanName...)
	if nil != err {
		h.interrupt <- err
		blog.Error("subscribe channel faile ", err.Error())
	}
	closeChan := make(chan struct{})
	go h.healthCheck(closeChan)
	defer func() {
		h.subscribing = false
		close(closeChan)
		blog.Infof("subChan Close")
		subChan.Unsubscribe(chanName...)
	}()

	var ts = time.Now()
	var cnt int64
	blog.Infof("subcribing channel %v", chanName)
	for {
		if false == h.isMaster {

			blog.Info("This is not master process, subChan Close")
			return
		}
		received, err := subChan.Receive()
		if err == redis.Nil || err == io.EOF {
			continue
		}
		msg, ok := received.(*redis.Message)
		if !ok {
			continue
		}
		if nil != err {
			blog.Debug("receive messave  err", err.Error())
			h.interrupt <- err
			continue
		}

		if "" == msg.Payload {
			continue
		}

		chanlen = len(h.msgChan)
		if h.maxSize*2 <= chanlen {

			blog.Infof("msgChan full, maxsize %d, len %d", h.maxSize, chanlen)
			h.clearMsgChan()
		}
		if chanlen != 0 && chanlen%10 == 0 {
			blog.Infof("buff len %d", chanlen)
		}
		h.lastMesgTs = time.Now()
		h.msgChan <- msg.Payload
		cnt++
		if cnt%10000 == 0 {
			blog.Infof("receive rate: %d/sec", int(float64(cnt)/time.Now().Sub(ts).Seconds()))
			cnt = 0
			ts = time.Now()
		}
	}
}

func (h *HostSnap) clearMsgChan() {
	ts := h.ts
	msgCnt := len(h.msgChan) - h.maxSize
	blog.Warnf("start clear %d", msgCnt)
	var cnt int
	for msgCnt > cnt {
		h.Lock()
		cnt++
		if ts != h.ts {
			msgCnt = len(h.msgChan) - h.maxSize
		} else {
			select {
			case <-time.After(time.Second * 10):
			case <-h.msgChan:
			}
		}
		h.Unlock()
	}
	if ts == h.ts {
		close(h.resetHandle)
	}
	blog.Warnf("cleared %d", cnt)
}

var cachelock = sync.RWMutex{}

func (h *HostSnap) getCache() *HostCache {
	cachelock.RLock()
	defer cachelock.RUnlock()
	return h.cache.cache[h.cache.flag]
}

func (h *HostSnap) setCache(key string, val *HostInst) {
	cachelock.Lock()
	h.cache.cache[h.cache.flag].set(key, val)
	cachelock.Unlock()
}

func (h *HostSnap) fetchDB() {
	cachelock.Lock()
	h.cache.cache[h.cache.flag] = h.fetch()
	cachelock.Unlock()
	go func() {
		ticker := time.NewTicker(fetchDBInterval)
		for {
			select {
			case <-ticker.C:
				cache := h.fetch()
				cachelock.Lock()
				h.cache.cache[!h.cache.flag] = cache
				h.cache.flag = !h.cache.flag
				cachelock.Unlock()
			case <-h.doneCh:
				blog.Warnf("close fetchDB")
				return
			}
		}
	}()
}

func (h *HostSnap) fetch() *HostCache {
	result := []map[string]interface{}{}
	err := h.db.GetMutilByCondition(common.BKTableNameBaseHost, nil, nil, &result, "", 0, 0)
	if err != nil {
		blog.Errorf("fetch db error %v", err)
	}
	hostcache := &HostCache{data: map[string]*HostInst{}}
	for index := range result {
		cloudid := fmt.Sprint(result[index][common.BKCloudIDField])
		innerip := fmt.Sprint(result[index][common.BKHostInnerIPField])
		hostcache.data[cloudid+"::"+innerip] = &HostInst{data: result[index]}
	}
	blog.Infof("success fetch %d collections to cache", len(result))
	return hostcache
}

type HostInst struct {
	sync.RWMutex
	data map[string]interface{}
}

func (h *HostInst) get(key string) interface{} {
	h.RLock()
	value := h.data[key]
	h.RUnlock()
	return value
}

func (h *HostInst) set(key string, value interface{}) {
	h.Lock()
	h.data[key] = value
	h.Unlock()
}

type HostCache struct {
	sync.RWMutex
	data map[string]*HostInst
}

func (h *HostCache) get(key string) *HostInst {
	h.RLock()
	value := h.data[key]
	h.RUnlock()
	return value
}

func (h *HostCache) set(key string, value *HostInst) {
	h.Lock()
	h.data[key] = value
	h.Unlock()
}

func (h *HostSnap) healthCheck(closeChan chan struct{}) {
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-h.doneCh:
			ticker.Stop()
			return
		case <-closeChan:
			ticker.Stop()
			return
		case <-ticker.C:
			channelstatus := 0
			if err := h.snapCli.Ping().Err(); err != nil {
				channelstatus = common.CCErrHostGetSnapshotChannelClose
				blog.Errorf("snap redis server connection error: %s", err.Error())
			} else if time.Now().Sub(h.lastMesgTs) > time.Minute {
				blog.Errorf("snapchannel was empty in last 1 min ")
				channelstatus = common.CCErrHostGetSnapshotChannelEmpty
			} else {
				channelstatus = common.CCSuccess
			}
			h.redisCli.Set(RedisSnapKeyChannelStatus, channelstatus, time.Minute*2)
		}
	}
}
