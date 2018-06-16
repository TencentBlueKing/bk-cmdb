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

package logics

import (
	bkcommon "configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/datacollection/common"
	"configcenter/src/source_controller/common/instdata"
	"fmt"
	"github.com/rs/xid"
	"github.com/tidwall/gjson"
	"io"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"time"

	redis "gopkg.in/redis.v5"
)

// const
var (
	// check master lifetime interval
	getMasterProcIntervalTime = time.Second * 10

	// locker expire duration in redis
	masterProcLockLiveTime       = getMasterProcIntervalTime + time.Second*10
	clearMsgChanTime       int64 = 120
	masterProcLockContent  string

	fetchDBInterval = time.Minute * 10
	maxconcurrent   = runtime.NumCPU()
)

// HostSnap define HostSnap
type HostSnap struct {
	id       string
	chanName string
	msgChan  chan string
	maxSize  int
	ts       time.Time // life cycle timestamp
	isMaster bool
	sync.Mutex
	interrupt     chan error
	maxconcurrent int

	resetHandle chan struct{}

	redisCli *redis.Client
	snapCli  *redis.Client

	subscribing bool

	cache *hostcache

	wg *sync.WaitGroup
}

type hostcache struct {
	cache map[bool]map[string]map[string]interface{}
	flag  bool
}

// NewHostSnap  returns new hostsnap object
//
// chanName: redis channel name，maxSize: max buffer cache, redisCli: CC redis cli, snapCli: snap redis cli
func NewHostSnap(chanName string, maxSize int, redisCli, snapCli *redis.Client, wg *sync.WaitGroup) *HostSnap {
	if 0 == maxSize {
		maxSize = 100
	}
	hostSnapInstance := &HostSnap{
		chanName:      chanName,
		msgChan:       make(chan string, maxSize*4),
		interrupt:     make(chan error),
		resetHandle:   make(chan struct{}),
		maxSize:       maxSize,
		redisCli:      redisCli,
		snapCli:       snapCli,
		ts:            time.Now(),
		id:            xid.New().String()[5:],
		maxconcurrent: maxconcurrent,
		wg:            wg,
		cache: &hostcache{
			cache: map[bool]map[string]map[string]interface{}{},
			flag:  false,
		},
	}
	return hostSnapInstance
}

// Start start main handle routines
func (h *HostSnap) Start() {

	defer h.wg.Done()

	go h.fetchDB()
	go h.Run()
}

// Run hostsnap main functionality
func (h *HostSnap) Run() {
	blog.Infof("datacollection start with maxconcurrent: %d", h.maxconcurrent)
	ticker := time.NewTicker(getMasterProcIntervalTime)
	var err error
	var msg string
	var msgs []string
	var addCount int
	var waitCnt int

	if h.saveRunning() {
		go h.subChan()
	} else {
		blog.Infof("run: there is other master process exists, recheck after %v ", getMasterProcIntervalTime)
	}
	for {
		select {
		case <-ticker.C:
			if h.saveRunning() {
				if !h.subscribing {
					go h.subChan()
				}
			}
		case msg = <-h.msgChan:
			// read all from msgChan and lock to prevent clear operation
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

			// handle them
			waitCnt = 0
			for {
				if waitCnt > h.maxconcurrent*2 {
					blog.Warnf("reset handlers")
					close(h.resetHandle)
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
		default:
			data := gjson.Get(msg, "data").String()
			val := gjson.Parse(data)
			host := h.getHostByVal(&val)
			if host == nil {
				blog.Infof("host not found, continue, %s", val.String())
				continue
			}
			hostid := fmt.Sprint(host[bkcommon.BKHostIDField])
			if hostid == "" {
				blog.Infof("host id not found, continue, %s", val.String())
				continue
			}

			// set snap cache
			h.redisCli.Set(common.RedisSnapKeyPrefix+hostid, data, time.Minute*10)

			// update host fields value
			condition := map[string]interface{}{bkcommon.BKHostIDField: host[bkcommon.BKHostIDField]}
			innerip, _ := host[bkcommon.BKHostInnerIPField].(string)
			outip, _ := host[bkcommon.BKHostOuterIPField].(string)
			setter := parseSetter(&val, innerip, outip)
			if needToUpdate(setter, host) {
				blog.Infof("update by %v, to %v", condition, setter)
				if err := instdata.UpdateHostByCondition(setter, condition); err != nil {
					blog.Error("update host error:", err.Error())
					continue
				}
				copyVal(setter, host)
			}
		}
	}

	return nil
}

func copyVal(a, b map[string]interface{}) {
	for k, v := range a {
		b[k] = v
	}
}
func needToUpdate(a, b map[string]interface{}) bool {
	for k, v := range a {
		if b[k] != v {
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
		ostype = bkcommon.HostOSTypeEnumLinux //"Linux"
	case "windows":
		version = strings.Replace(version, "Microsoft ", "", 1)
		platform = strings.Replace(platform, "Microsoft ", "", 1)
		osname = fmt.Sprintf("%s", platform)
		ostype = bkcommon.HostOSTypeEnumWindows // "Windows"
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

	// TODO add log when fields empty
	return map[string]interface{}{
		"bk_cpu":        cupnum,
		"bk_cpu_module": cpumodule,
		"bk_cpu_mhz":    CPUMhz,                    //Mhz
		"bk_disk":       disk / 1024 / 1024 / 1024, //GB
		"bk_mem":        mem / 1024 / 1024,         //MB
		"bk_os_type":    ostype,
		"bk_os_name":    osname,
		"bk_os_version": version,
		"bk_host_name":  hostname,
		"bk_outer_mac":  OuterMAC,
		"bk_mac":        InnerMAC,
		"bk_os_bit":     osbit,
	}
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

func (h *HostSnap) getHostByVal(val *gjson.Result) map[string]interface{} {
	cloudid := val.Get("cloudid").String()
	/*if cloudid == "0" || cloudid == "" {
		cloudid := common.BKCloudIDField
	}*/
	ips := getIPS(val)
	if len(ips) > 0 {
		blog.Infof("handle clouid: %s ips: %v", cloudid, ips)
		for _, ip := range ips {
			if host := h.getCache()[cloudid+"::"+ip]; host != nil {
				return host
			}
		}

		blog.Infof("ips not in cache clouid: %s,ip: %v", cloudid, ips)
		clouidInt, _ := strconv.Atoi(cloudid)
		condition := map[string]interface{}{
			bkcommon.BKCloudIDField: clouidInt,
			bkcommon.BKHostInnerIPField: map[string]interface{}{
				bkcommon.BKDBIN: ips,
			},
		}
		result := []map[string]interface{}{}
		err := instdata.GetHostByCondition(nil, condition, &result, "", 0, 0)
		if err != nil {
			blog.Errorf("fetch db error %v", err)
		}
		for _, host := range result {
			cloudid := fmt.Sprint(host[bkcommon.BKCloudIDField])
			innerip := fmt.Sprint(host[bkcommon.BKHostInnerIPField])
			h.setCache(cloudid+"::"+innerip, host)
			return h.getCache()[cloudid+"::"+innerip]
		}
		blog.Infof("ips not in cache and db, clouid: %v, ip: %v", cloudid, ips)
	} else {
		blog.Errorf("message has no ip, message:%s", val.String())
	}
	return nil
}

// concede concede when buffer fulled
func (h *HostSnap) concede() {
	blog.Info("concede")
	h.isMaster = false
	h.subscribing = false
	val := h.redisCli.Get(common.MasterProcLockKey).Val()
	if val != h.id {
		h.redisCli.Del(common.MasterProcLockKey)
	}
}

// saveRunning lock master process
func (h *HostSnap) saveRunning() (ok bool) {
	var err error
	setNXChan := make(chan struct{})
	go func() {
		select {
		case <-time.After(masterProcLockLiveTime):
			blog.Fatalf("saveRunning check: set nx time out!! the network may be broken, redis stats: %v ", h.redisCli.PoolStats())
		case <-setNXChan:
		}
	}()
	if h.isMaster {
		var val string
		val, err = h.redisCli.Get(common.MasterProcLockKey).Result()
		if err != nil {
			blog.Errorf("master: saveRunning err %v", err)
		} else if val == h.id {
			blog.Infof("master check : i am still master")
			h.redisCli.Set(common.MasterProcLockKey, h.id, masterProcLockLiveTime)
			ok = true
		} else {
			blog.Infof("exit master,val = %v, id = %v", val, h.id)
			h.isMaster = false
			ok = false
		}
	} else {
		ok, err = h.redisCli.SetNX(common.MasterProcLockKey, h.id, masterProcLockLiveTime).Result()
		if err != nil {
			blog.Errorf("slave: saveRunning err %v", err)
		} else if ok {
			blog.Infof("slave check: ok")
			blog.Infof("i am master from now")
			h.isMaster = true
		} else {
			blog.Infof("slave check: there is other master process exists, recheck after %v ", getMasterProcIntervalTime)
		}
	}
	close(setNXChan)
	return ok
}

const mockmsg = "{\"localTime\": \"2017-09-19 16:57:00\", \"data\": \"{\\\"ip\\\":\\\"127.0.0.1\\\",\\\"bizid\\\":0,\\\"cloudid\\\":1,\\\"data\\\":{\\\"timezone\\\":8,\\\"datetime\\\":\\\"2017-09-19 16:57:07\\\",\\\"utctime\\\":\\\"2017-09-19 08:57:07\\\",\\\"country\\\":\\\"Asia\\\",\\\"city\\\":\\\"Shanghai\\\",\\\"cpu\\\":{\\\"cpuinfo\\\":[{\\\"cpu\\\":0,\\\"vendorID\\\":\\\"GenuineIntel\\\",\\\"family\\\":\\\"6\\\",\\\"model\\\":\\\"63\\\",\\\"stepping\\\":2,\\\"physicalID\\\":\\\"0\\\",\\\"coreID\\\":\\\"0\\\",\\\"cores\\\":1,\\\"modelName\\\":\\\"Intel(R) Xeon(R) CPU E5-26xx v3\\\",\\\"mhz\\\":2294.01,\\\"cacheSize\\\":4096,\\\"flags\\\":[\\\"fpu\\\",\\\"vme\\\",\\\"de\\\",\\\"pse\\\",\\\"tsc\\\",\\\"msr\\\",\\\"pae\\\",\\\"mce\\\",\\\"cx8\\\",\\\"apic\\\",\\\"sep\\\",\\\"mtrr\\\",\\\"pge\\\",\\\"mca\\\",\\\"cmov\\\",\\\"pat\\\",\\\"pse36\\\",\\\"clflush\\\",\\\"mmx\\\",\\\"fxsr\\\",\\\"sse\\\",\\\"sse2\\\",\\\"ss\\\",\\\"ht\\\",\\\"syscall\\\",\\\"nx\\\",\\\"lm\\\",\\\"constant_tsc\\\",\\\"up\\\",\\\"rep_good\\\",\\\"unfair_spinlock\\\",\\\"pni\\\",\\\"pclmulqdq\\\",\\\"ssse3\\\",\\\"fma\\\",\\\"cx16\\\",\\\"pcid\\\",\\\"sse4_1\\\",\\\"sse4_2\\\",\\\"x2apic\\\",\\\"movbe\\\",\\\"popcnt\\\",\\\"tsc_deadline_timer\\\",\\\"aes\\\",\\\"xsave\\\",\\\"avx\\\",\\\"f16c\\\",\\\"rdrand\\\",\\\"hypervisor\\\",\\\"lahf_lm\\\",\\\"abm\\\",\\\"xsaveopt\\\",\\\"bmi1\\\",\\\"avx2\\\",\\\"bmi2\\\"],\\\"microcode\\\":\\\"1\\\"}],\\\"per_usage\\\":[3.0232169701043103],\\\"total_usage\\\":3.0232169701043103,\\\"per_stat\\\":[{\\\"cpu\\\":\\\"cpu0\\\",\\\"user\\\":5206.09,\\\"system\\\":6107.04,\\\"idle\\\":337100.84,\\\"nice\\\":6.68,\\\"iowait\\\":528.24,\\\"irq\\\":0.02,\\\"softirq\\\":13.48,\\\"steal\\\":0,\\\"guest\\\":0,\\\"guestNice\\\":0,\\\"stolen\\\":0}],\\\"total_stat\\\":{\\\"cpu\\\":\\\"cpu-total\\\",\\\"user\\\":5206.09,\\\"system\\\":6107.04,\\\"idle\\\":337100.84,\\\"nice\\\":6.68,\\\"iowait\\\":528.24,\\\"irq\\\":0.02,\\\"softirq\\\":13.48,\\\"steal\\\":0,\\\"guest\\\":0,\\\"guestNice\\\":0,\\\"stolen\\\":0}},\\\"env\\\":{\\\"crontab\\\":[{\\\"user\\\":\\\"root\\\",\\\"content\\\":\\\"#secu-tcs-agent monitor, install at Fri Sep 15 16:12:02 CST 2017\\\\n* * * * * /usr/local/sa/agent/secu-tcs-agent-mon-safe.sh /usr/local/sa/agent \\\\u003e /dev/null 2\\\\u003e\\\\u00261\\\\n*/1 * * * * /usr/local/qcloud/stargate/admin/start.sh \\\\u003e /dev/null 2\\\\u003e\\\\u00261 \\\\u0026\\\\n*/20 * * * * /usr/sbin/ntpdate ntpupdate.tencentyun.com \\\\u003e/dev/null \\\\u0026\\\\n*/1 * * * * cd /usr/local/gse/gseagent; ./cron_agent.sh 1\\\\u003e/dev/null 2\\\\u003e\\\\u00261\\\\n\\\"}],\\\"host\\\":\\\"127.0.0.1  localhost  localhost.localdomain  VM_0_31_centos\\\\n::1         localhost localhost.localdomain localhost6 localhost6.localdomain6\\\\n\\\",\\\"route\\\":\\\"Kernel IP routing table\\\\nDestination     Gateway         Genmask         Flags Metric Ref    Use Iface\\\\n10.0.0.0        0.0.0.0         255.255.255.0   U     0      0        0 eth0\\\\n169.254.0.0     0.0.0.0         255.255.0.0     U     1002   0        0 eth0\\\\n0.0.0.0         127.0.0.1        0.0.0.0         UG    0      0        0 eth0\\\\n\\\"},\\\"disk\\\":{\\\"diskstat\\\":{\\\"vda1\\\":{\\\"major\\\":252,\\\"minor\\\":1,\\\"readCount\\\":24347,\\\"mergedReadCount\\\":570,\\\"writeCount\\\":696357,\\\"mergedWriteCount\\\":4684783,\\\"readBytes\\\":783955968,\\\"writeBytes\\\":22041231360,\\\"readSectors\\\":1531164,\\\"writeSectors\\\":43049280,\\\"readTime\\\":80626,\\\"writeTime\\\":12704736,\\\"iopsInProgress\\\":0,\\\"ioTime\\\":822057,\\\"weightedIoTime\\\":12785026,\\\"name\\\":\\\"vda1\\\",\\\"serialNumber\\\":\\\"\\\",\\\"speedIORead\\\":0,\\\"speedByteRead\\\":0,\\\"speedIOWrite\\\":2.9,\\\"speedByteWrite\\\":171144.53333333333,\\\"util\\\":0.0025666666666666667,\\\"avgrq_sz\\\":115.26436781609195,\\\"avgqu_sz\\\":0.06568333333333334,\\\"await\\\":22.649425287356323,\\\"svctm\\\":0.8850574712643678}},\\\"partition\\\":[{\\\"device\\\":\\\"/dev/vda1\\\",\\\"mountpoint\\\":\\\"/\\\",\\\"fstype\\\":\\\"ext3\\\",\\\"opts\\\":\\\"rw,noatime,acl,user_xattr\\\"}],\\\"usage\\\":[{\\\"path\\\":\\\"/\\\",\\\"fstype\\\":\\\"ext2/ext3\\\",\\\"total\\\":52843638784,\\\"free\\\":47807447040,\\\"used\\\":2351915008,\\\"usedPercent\\\":4.4507060113962345,\\\"inodesTotal\\\":3276800,\\\"inodesUsed\\\":29554,\\\"inodesFree\\\":3247246,\\\"inodesUsedPercent\\\":0.9019165039062501}]},\\\"load\\\":{\\\"load_avg\\\":{\\\"load1\\\":0,\\\"load5\\\":0,\\\"load15\\\":0}},\\\"mem\\\":{\\\"meminfo\\\":{\\\"total\\\":1044832256,\\\"available\\\":805912576,\\\"used\\\":238919680,\\\"usedPercent\\\":22.866797864249705,\\\"free\\\":92041216,\\\"active\\\":521183232,\\\"inactive\\\":352964608,\\\"wired\\\":0,\\\"buffers\\\":110895104,\\\"cached\\\":602976256,\\\"writeback\\\":0,\\\"dirty\\\":151552,\\\"writebacktmp\\\":0},\\\"vmstat\\\":{\\\"total\\\":0,\\\"used\\\":0,\\\"free\\\":0,\\\"usedPercent\\\":0,\\\"sin\\\":0,\\\"sout\\\":0}},\\\"net\\\":{\\\"interface\\\":[{\\\"mtu\\\":65536,\\\"name\\\":\\\"lo\\\",\\\"hardwareaddr\\\":\\\"28:31:52:1d:c6:0a\\\",\\\"flags\\\":[\\\"up\\\",\\\"loopback\\\"],\\\"addrs\\\":[{\\\"addr\\\":\\\"127.0.0.1/8\\\"}]},{\\\"mtu\\\":1500,\\\"name\\\":\\\"eth0\\\",\\\"hardwareaddr\\\":\\\"52:54:00:19:2e:e8\\\",\\\"flags\\\":[\\\"up\\\",\\\"broadcast\\\",\\\"multicast\\\"],\\\"addrs\\\":[{\\\"addr\\\":\\\"127.0.0.1/24\\\"}]}],\\\"dev\\\":[{\\\"name\\\":\\\"lo\\\",\\\"speedSent\\\":0,\\\"speedRecv\\\":0,\\\"speedPacketsSent\\\":0,\\\"speedPacketsRecv\\\":0,\\\"bytesSent\\\":604,\\\"bytesRecv\\\":604,\\\"packetsSent\\\":2,\\\"packetsRecv\\\":2,\\\"errin\\\":0,\\\"errout\\\":0,\\\"dropin\\\":0,\\\"dropout\\\":0,\\\"fifoin\\\":0,\\\"fifoout\\\":0},{\\\"name\\\":\\\"eth0\\\",\\\"speedSent\\\":574,\\\"speedRecv\\\":214,\\\"speedPacketsSent\\\":3,\\\"speedPacketsRecv\\\":2,\\\"bytesSent\\\":161709123,\\\"bytesRecv\\\":285910298,\\\"packetsSent\\\":1116625,\\\"packetsRecv\\\":1167796,\\\"errin\\\":0,\\\"errout\\\":0,\\\"dropin\\\":0,\\\"dropout\\\":0,\\\"fifoin\\\":0,\\\"fifoout\\\":0}],\\\"netstat\\\":{\\\"established\\\":2,\\\"syncSent\\\":1,\\\"synRecv\\\":0,\\\"finWait1\\\":0,\\\"finWait2\\\":0,\\\"timeWait\\\":0,\\\"close\\\":0,\\\"closeWait\\\":0,\\\"lastAck\\\":0,\\\"listen\\\":2,\\\"closing\\\":0},\\\"protocolstat\\\":[{\\\"protocol\\\":\\\"udp\\\",\\\"stats\\\":{\\\"inDatagrams\\\":176253,\\\"inErrors\\\":0,\\\"noPorts\\\":1,\\\"outDatagrams\\\":199569,\\\"rcvbufErrors\\\":0,\\\"sndbufErrors\\\":0}}]},\\\"system\\\":{\\\"info\\\":{\\\"hostname\\\":\\\"VM_0_31_centos\\\",\\\"uptime\\\":348315,\\\"bootTime\\\":1505463112,\\\"procs\\\":142,\\\"os\\\":\\\"linux\\\",\\\"platform\\\":\\\"centos\\\",\\\"platformFamily\\\":\\\"rhel\\\",\\\"platformVersion\\\":\\\"6.2\\\",\\\"kernelVersion\\\":\\\"2.6.32-504.30.3.el6.x86_64\\\",\\\"virtualizationSystem\\\":\\\"\\\",\\\"virtualizationRole\\\":\\\"\\\",\\\"hostid\\\":\\\"96D0F4CA-2157-40E6-BF22-6A7CD9B6EB8C\\\",\\\"systemtype\\\":\\\"64-bit\\\"}}}}\", \"timestamp\": 1505811427, \"dtEventTime\": \"2017-09-19 16:57:07\", \"dtEventTimeStamp\": 1505811427000}"

// subChan subscribe message from redis channel
func (h *HostSnap) subChan() {
	h.subscribing = true
	var chanlen int
	subChan, err := h.snapCli.Subscribe(h.chanName)
	if nil != err {
		h.interrupt <- err
		blog.Error("subscribe channel faile ", err.Error())
	}
	defer func() {
		subChan.Unsubscribe(h.chanName)
		h.subscribing = false
		blog.Infof("subChan Close")
	}()

	var ts = time.Now()
	var cnt int64
	blog.Infof("subcribing channel %s", h.chanName)
	for {
		if false == h.isMaster {
			// not master again, close subscribe to prevent unnecessary subscript
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
			//  if msgChan fulled, clear old msgs
			blog.Infof("msgChan full, maxsize %d, len %d", h.maxSize, chanlen)
			h.clearMsgChan()
		}
		if chanlen != 0 && chanlen%10 == 0 {
			blog.Infof("buff len %d", chanlen)
		}

		h.msgChan <- msg.Payload
		cnt++
		if cnt%10000 == 0 {
			blog.Infof("receive rate: %d/sec", int(float64(cnt)/time.Now().Sub(ts).Seconds()))
			cnt = 0
			ts = time.Now()
		}
	}
}

//clearMsgChan clear msgchan when msgchan is twice length of maxsize
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

var cachelock = sync.Mutex{}

func (h *HostSnap) getCache() map[string]map[string]interface{} {
	cachelock.Lock()
	defer cachelock.Unlock()
	return h.cache.cache[h.cache.flag]
}

func (h *HostSnap) setCache(key string, val map[string]interface{}) {
	cachelock.Lock()
	h.cache.cache[h.cache.flag][key] = val
	cachelock.Unlock()
}

func (h *HostSnap) fetchDB() {
	cachelock.Lock()
	h.cache.cache[h.cache.flag] = fetch()
	cachelock.Unlock()
	go func() {
		ticker := time.NewTicker(fetchDBInterval)
		for range ticker.C {
			cache := fetch()
			cachelock.Lock()
			h.cache.cache[!h.cache.flag] = cache
			h.cache.flag = !h.cache.flag
			cachelock.Unlock()
		}
	}()
}

func fetch() map[string]map[string]interface{} {
	result := []map[string]interface{}{}
	err := instdata.GetHostByCondition(nil, nil, &result, "", 0, 0)
	if err != nil {
		blog.Errorf("fetch db error %v", err)
	}
	hostcache := map[string]map[string]interface{}{}
	for _, host := range result {
		cloudid := fmt.Sprint(host[bkcommon.BKCloudIDField])
		innerip := fmt.Sprint(host[bkcommon.BKHostInnerIPField])
		hostcache[cloudid+"::"+innerip] = host
	}
	blog.Infof("success fetch %d collections to cache", len(result))
	return hostcache
}
