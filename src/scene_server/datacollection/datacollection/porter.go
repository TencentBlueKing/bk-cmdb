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
	"net"
	"runtime"
	"runtime/debug"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/xid"
	"gopkg.in/redis.v5"
)

type chanCollector struct {
	p Porter
}

// 14kB * 10000 = 140M
const cacheSize = 10000

func BuildChanPorter(name string, analyzer Analyzer, redisCli, snapCli *redis.Client, channels []string, mockmesg string, registry prometheus.Registerer, engine *backbone.Engine) *chanPorter {
	porter := &chanPorter{
		analyzer:        analyzer,
		name:            name,
		pid:             xid.New().String(),
		isMaster:        util.NewBool(false),
		redisCli:        redisCli,
		snapCli:         snapCli,
		channels:        channels,
		analyzeC:        make(chan string, cacheSize),
		slaveC:          make(chan string, cacheSize),
		analyzeCounterC: make(chan int, runtime.NumCPU()),
		runed:           util.NewBool(false),
		Engine:          engine,
	}

	ns := "cmdb_collector_" + name + "_"

	registry.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: ns + "is_master",
			Help: "describe whether this process is master.",
		},
		func() float64 { return float64(*porter.isMaster) },
	))

	registry.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: ns + "analyze_queue",
			Help: "current number of analyze queue.",
		},
		func() float64 { return float64(len(porter.analyzeC)) },
	))

	registry.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: ns + "slave_queue",
			Help: "current number of slave queue.",
		},
		func() float64 { return float64(len(porter.slaveC)) },
	))

	porter.analyseDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name: ns + "analyze_duration",
			Help: "analyze duration of each message.",
		},
	)
	registry.MustRegister(porter.analyseDuration)

	porter.receiveTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: ns + "receive_total",
			Help: "number of received message.",
		},
	)
	registry.MustRegister(porter.receiveTotal)

	porter.pushTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: ns + "push_total",
			Help: "number of pushed message.",
		},
	)
	registry.MustRegister(porter.pushTotal)

	porter.analyzeTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: ns + "analyze_total",
			Help: "number of analyzed message.",
		},
		[]string{LabelStatus},
	)
	registry.MustRegister(porter.analyzeTotal)

	return porter
}

const (
	LabelStatus = "status"
)

type chanPorter struct {
	analyseDuration prometheus.Histogram
	receiveTotal    prometheus.Counter
	pushTotal       prometheus.Counter
	analyzeTotal    *prometheus.CounterVec

	// 分析器
	analyzer Analyzer

	// porter 名称，用于打印日志
	name string

	// porter 的ID，用于抢master锁
	pid string

	// 标识当前协程是否是master协程，master协程负责从redis channel 取数据并把处理不过来的数据推送到slavequeue
	isMaster *util.AtomicBool

	// cc自己的redis，用于抢master锁，缓存slavequeue
	redisCli *redis.Client

	// 数据来源的redis，master 从这个redis读channel
	snapCli *redis.Client

	// redis channel 名称
	channels []string

	// 待处理队列，analyzer只消费这个队列的消息
	analyzeC chan string

	// 待推送到redis的消息队列
	slaveC chan string

	// 标识本porter是否已经运行过
	runed *util.AtomicBool

	// 最后一次收到消息的时间，用于健康检查
	lastMesgTs time.Time

	// 用于统计处理效率的
	analyzeCounterC chan int

	// 用于判断是否是master的
	*backbone.Engine
}

func (p *chanPorter) Name() string {
	return p.name
}

func (p *chanPorter) Mock(mesg string) error {
	select {
	case p.analyzeC <- mesg:
	default:
		return fmt.Errorf("message queue fulled")
	}
	return nil
}

func (p *chanPorter) Run() error {
	p.runed.Set()
	if p.runed.IsSet() {
		// 防止被上层manager重复执行, healthCheckLoop, analyzeLoop只需要运行一个即可
		go p.healthCheckLoop()
		go p.analyzeCount()
		for i := 0; i < runtime.NumCPU(); i++ {
			go p.analyzeLoop()
		}
		go p.popLoop()
		go p.pushLoop()
	}
	var err error
	for {
		if err = p.collect(); err != nil {
			blog.Errorf("[data-collection][%s] collect message failed: %v, retry 3s later", p.name, err)
		}
		// 睡3秒， 防止空跑导致CPU占用高涨
		time.Sleep(time.Second * 3)
	}
}

func (p *chanPorter) analyzeLoop() {
	for {
		p.analyze()
	}
}

func (p *chanPorter) analyze() {
	defer func() {
		if sysErr := recover(); sysErr != nil {
			blog.Errorf("[data-collection][%s] analyzeLoop panic by: %v, stack:\n %s", p.name, sysErr, debug.Stack())
		}
	}()

	var mesg string
	var err error

	for mesg = range p.analyzeC {
		before := time.Now()
		if err = p.analyzer.Analyze(mesg); err != nil {
			blog.Errorf("[data-collection][%s] analyze message failed: %v, raw msg: %s", p.name, err, mesg)
			p.analyzeTotal.WithLabelValues("failed").Inc()
		} else {
			p.analyzeTotal.WithLabelValues("success").Inc()
		}
		p.analyseDuration.Observe(time.Since(before).Seconds())
		p.analyzeCounterC <- 1
	}
}

func (p *chanPorter) analyzeCount() {
	var cnt int
	var ts = time.Now()
	var i int
	for i = range p.analyzeCounterC {
		cnt += i
		if time.Since(ts) > time.Minute*10 {
			blog.Infof("[data-collection][%s] analyze rate: %d message in last %v, analyzeC length: %d", p.name, cnt, time.Now().Sub(ts), len(p.analyzeC))
			cnt = 0
			ts = time.Now()
		}
	}
}

// collect 获取待处理消息，当是master时从redis channel获取，当是slave时从 redis queue 获取
func (p *chanPorter) collect() error {
	if !p.Engine.ServiceManageInterface.IsMaster() {
		p.isMaster.UnSet()
		blog.Infof("[data-collection][%s] %v", p.name, "there is other master")
		return nil
	}
	blog.Infof("[data-collection][%s] i am master(id: %s) from now", p.name, p.pid)
	defer blog.Infof("[data-collection][%s] exist master(id: %s) from now", p.name, p.pid)
	p.isMaster.Set()
	defer p.isMaster.UnSet()

	// 开始订阅
	err := p.subscribeLoop()
	if err != nil {
		return fmt.Errorf("subscribe channel return an error: %v", err)
	}

	return nil
}

func (p *chanPorter) subscribeLoop() error {
	subChan, err := p.snapCli.Subscribe(p.channels...)
	if nil != err {
		return fmt.Errorf("subscribe channel failed, %v", err)
	}
	defer func() {
		_ = subChan.Unsubscribe(p.channels...)
		_ = subChan.Close()
	}()

	blog.Info("[data-collection][%s] subscribing channel %v from redis", p.name, p.channels)
	defer blog.Info("[data-collection][%s] unsubscribe channel %v from redis", p.name, p.channels)

	ts := time.Now()
	var cnt int64
	var timeoutErr net.Error
	var ok bool
	var received interface{}
	var name = p.name + "[receive]"
	for p.Engine.ServiceManageInterface.IsMaster() {
		p.isMaster.Set()
		received, err = subChan.ReceiveTimeout(time.Second * 10)
		if err == redis.Nil || err == io.EOF {
			continue
		}
		if timeoutErr, ok = err.(net.Error); ok && timeoutErr.Timeout() {
			continue
		}
		if nil != err {
			return fmt.Errorf("receive message from redis failed: %v", err)
		}
		msg, ok := received.(*redis.Message)
		if !ok {
			continue
		}

		if "" == msg.Payload {
			continue
		}

		p.receiveTotal.Inc()

		// 当mesgC满时表明已达到本进程的处理速度上限，此时我们推送该消息到slavequeue让其他进程协助处理
		select {
		case p.analyzeC <- msg.Payload:
		default:
			select {
			case p.slaveC <- msg.Payload:
			default:
				writeOrClearChan(p.analyzeC, name, msg.Payload)
			}
		}

		cnt++

		p.lastMesgTs = time.Now()
		if time.Since(ts) > time.Minute*10 {
			blog.Infof("[data-collection][%s] receive rate: %d message in last %v", p.name, cnt, time.Now().Sub(ts))
			cnt = 0
			ts = time.Now()
		}
	}
	p.isMaster.UnSet()
	return nil
}

func (p *chanPorter) renewalMasterLoop() {
	var err error
	for range time.Tick(time.Second * 3) {
		if err = renewalMaster(p.redisCli, p.name, p.pid); err != nil {
			blog.Warnf("[data-collection][%s] renewal master failed: %v", p.name, err)
			p.isMaster.UnSet()
			return
		}
	}
}

func (p *chanPorter) healthCheckLoop() {
	for {
		p.healthCheck()
	}
}

// healthCheck 报告自己的状态
func (p *chanPorter) healthCheck() {
	ticker := time.NewTicker(time.Minute)
	defer func() {
		ticker.Stop()
		if sysErr := recover(); sysErr != nil {
			blog.Errorf("[data-collection][%s] panic by: %v, stack:\n %s", p.name, sysErr, debug.Stack())
		}
	}()

	var err error
	var now time.Time
	for now = range ticker.C {
		var channelStatus int
		if err = p.snapCli.Ping().Err(); err != nil {
			channelStatus = common.CCErrHostGetSnapshotChannelClose
			blog.Errorf("[data-collection][%s][healthCheck] snap redis server connection error: %s", p.name, err.Error())
		} else if err = p.redisCli.Ping().Err(); err != nil {
			channelStatus = common.CCErrHostGetSnapshotChannelClose
			blog.Errorf("[data-collection][%s][healthCheck] cc redis server connection error: %s", p.name, err.Error())
		} else if p.isMaster.IsSet() && now.Sub(p.lastMesgTs) > time.Minute {
			blog.Warnf("[data-collection][%s][healthCheck] snapshot channel was empty in last 1 min", p.name)
			channelStatus = common.CCErrHostGetSnapshotChannelEmpty
		} else {
			channelStatus = common.CCSuccess
		}
		if err = p.redisCli.Set(channelStatusKey(p.name), channelStatus, time.Minute*2).Err(); err != nil {
			blog.Errorf("[data-collection][%s][healthCheck] set channel status failed: %v", p.name, err)
		}
	}
}

// popLoop 从slave处理队列获取消息，从而协助master处理
// 因为有可能单机部署，所以即使是master也要处理slavequeue
func (p *chanPorter) popLoop() {
	for {
		p.pop()
	}
}

func (p *chanPorter) pop() {
	blog.Info("[data-collection][%s] start popLoop from redis", p.name)
	defer blog.Info("[data-collection][%s] stop popLoop from redis", p.name)

	defer func() {
		if sysErr := recover(); sysErr != nil {
			blog.Errorf("[data-collection][%s] panic by: %v, stack:\n %s", p.name, sysErr, debug.Stack())
		}
	}()

	// 推消息到slave处理队列
	var msg []string
	var err error
	var timeoutErr net.Error
	var ok bool
	var llen int64
	var key = slaveQueueKey(p.name)
	var name = p.name + "[pop]"
	for {
		msg, err = p.redisCli.BRPop(time.Second*30, key).Result()
		if err == redis.Nil {
			continue
		}
		if timeoutErr, ok = err.(net.Error); ok && timeoutErr.Timeout() {
			continue
		}
		if err != nil {
			blog.Errorf("[data-collection][%s] pop message from redis failed: %v, retry 3s later", p.name, err)
			// 睡3秒，防止空跑导致CPU占用高涨
			time.Sleep(time.Second * 3)
		}
		if len(msg) > 1 && msg[1] != "nil" && msg[1] != "" {
			writeOrClearChan(p.analyzeC, name, msg[1])
		}
		if p.isMaster.IsSet() {
			llen, err = p.redisCli.LLen(key).Result()
			if err != nil {
				blog.Errorf("[data-collection][%s] llen failed: %v", p.name, err)
				continue
			}
			if llen > cacheSize {
				// 清理超过处理能力的未处理消息
				blog.Errorf("[data-collection][%s] slave queue %v fulled, clear it", p.name, key)
				if err = p.redisCli.Del(key).Err(); err != nil {
					blog.Errorf("[data-collection][%s] llen failed: %v", p.name, err)
					continue
				}
			}
			// 是master时，sleep可以让slave pop更多的消息
			time.Sleep(time.Second)
		}
	}
}

func (p *chanPorter) pushLoop() {
	for {
		p.push()
	}
}

// push 把master处理不过来的消息推到slave处理队列，让slave协助处理
func (p *chanPorter) push() {
	blog.Info("[data-collection][%s] start pushLoop to redis", p.name)
	defer blog.Info("[data-collection][%s] stop pushLoop to redis", p.name)

	defer func() {
		if sysErr := recover(); sysErr != nil {
			blog.Errorf("[data-collection][%s] panic by: %v, stack:\n %s", p.name, sysErr, debug.Stack())
		}
	}()

	var msg string
	var err error
	key := slaveQueueKey(p.name)
	for msg = range p.slaveC {
		p.pushTotal.Inc()
		if err = p.redisCli.LPush(key, msg).Err(); err != nil {
			blog.Errorf("[data-collection][%s] push message to redis failed: %v", p.name, err)
		}
	}
}

// writeOrClearChan 利用非阻塞读channel达到清里channel的目的
func writeOrClearChan(msgC chan string, name, msg string) {
	select {
	case msgC <- msg:
	default:
		// channel fulled, so we drop 200 oldest events from queue
		blog.Infof("[data-collection][%s] msgChan full, len %d. clear 200 oldest from queue", name, len(msgC))
		defer blog.Infof("[data-collection][%s] msgChan full, len %d. cleared 200 oldest from queue", name, len(msgC))
		var ok bool
		for i := 0; i < 200; i++ {
			_, ok = <-msgC
			if !ok {
				break
			}
		}
		select {
		case msgC <- msg:
		default:
		}
	}
}

func renewalMaster(redisCli *redis.Client, name string, procID string) error {
	lockKey := masterLockKey(name)
	masterPID, err := redisCli.Get(lockKey).Result()
	if err != nil {
		return fmt.Errorf("key [%s] value nil: %v", lockKey, err)
	}
	if masterPID != procID {
		return fmt.Errorf("there is other master(id: %s) running", masterPID)
	}
	err = redisCli.Set(lockKey, procID, masterProcLockLiveTime).Err()
	if err != nil {
		return fmt.Errorf("renewal failed: %v, key(%s)", err, lockKey)
	}
	return nil
}

// loginMaster 抢master锁，当已经是master时给锁续期
func loginMaster(redisCli *redis.Client, name string, procID string) error {
	lockKey := masterLockKey(name)
	var err error
	var ok bool
	var masterPID string
	for {
		masterPID, err = redisCli.Get(lockKey).Result()
		if err == redis.Nil {
			break
		}
		if err != nil {
			return fmt.Errorf("get master failed: %v, key(%s)", err, lockKey)
		}
		if masterPID == "" {
			break
		}
		if masterPID != procID {
			return fmt.Errorf("there is other master(id: %s) running", masterPID)
		}
		err = redisCli.Set(lockKey, procID, masterProcLockLiveTime).Err()
		if err != nil {
			return fmt.Errorf("renewal failed: %v, key(%s)", err, lockKey)
		}
		return nil
	}

	ok, err = redisCli.SetNX(lockKey, procID, masterProcLockLiveTime).Result()
	if err != nil {
		return fmt.Errorf("race to master failed: %v, key(%s)", err, lockKey)
	}
	if !ok {
		masterPID, err = redisCli.Get(lockKey).Result()
		if err != nil {
			return fmt.Errorf("get master failed: %v, key(%s)", err, lockKey)
		}
		return fmt.Errorf("there is other master(id: %s) running", masterPID)
	}
	return nil
}

// logoutMaster 主动退出master
func logoutMaster(redisCli *redis.Client, name string, procID string) error {
	lockKey := masterLockKey(name)
	masterPID, err := redisCli.Get(lockKey).Result()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return fmt.Errorf("get master failed: %v, key(%s)", err, lockKey)
	}
	if masterPID == "" {
		return nil
	}
	if masterPID != procID {
		return fmt.Errorf("there is other master(id: %s) running", masterPID)
	}
	err = redisCli.Del(lockKey).Err()
	if err != nil {
		return fmt.Errorf("release master failed: %v, key(%s)", err, lockKey)
	}
	blog.Infof("[data-collection][%s] logout master success", name)
	return nil
}

// masterLockKey master锁的key
func masterLockKey(name string) string {
	return common.BKCacheKeyV3Prefix + name + ":masterlock"
}

// slaveQueueKey 交给slave处理的消息待处理队列的key
func slaveQueueKey(name string) string {
	return common.BKCacheKeyV3Prefix + name + ":queue"
}

// channelStatusKey 通道状态的key
func channelStatusKey(name string) string {
	return common.BKCacheKeyV3Prefix + name + ":channelstatus"
}
