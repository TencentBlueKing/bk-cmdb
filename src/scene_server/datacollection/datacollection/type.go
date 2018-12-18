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
	"strings"
	"sync"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"

	"github.com/rs/xid"
	"gopkg.in/redis.v5"
)

type Analyzer interface {
	Analyze(mesg string) error
}

type Collector interface {
	Collect(ConllectorConfig)
}

type Porter interface {
	Name() string
	Run(ConllectorConfig) error
}

type ConllectorConfig mapstr.MapStr

type chanCollector struct {
	p Porter
}

type ChanHandler interface {
	Analyzer
}

func BuildChanPorter(name string, analyzer Analyzer, channels []string, redisCli, snapCli *redis.Client) *chanPorter {
	return &chanPorter{
		analyzer: analyzer,
		name:     name,
		pid:      xid.New().String(),
		isMaster: util.NewBool(false),
		redisCli: redisCli,
		snapCli:  snapCli,
		channels: channels,
		mesgC:    make(chan string, 1000),
	}
}

type chanPorter struct {
	analyzer Analyzer
	name     string
	pid      string
	isMaster *util.AtomicBool
	redisCli *redis.Client
	snapCli  *redis.Client
	channels []string
	mesgC    chan string

	lastMesgTs      time.Time
	healthCheckOnce sync.Once
	analyzeLoopOnce sync.Once
	runed           *util.AtomicBool
	popLock         sync.Mutex
	poping          bool
}

func (p *chanPorter) Run() error {
	p.runed.Set()
	if p.runed.IsSet() {
		// 防止被上层manager重复执行, healthCheckLoop, analyzeLoop只需要运行一个即可
		go p.healthCheckLoop()
		for i := 0; i < runtime.NumCPU(); i++ {
			go p.analyzeLoop()
		}
	}
	for {
		if err := p.collect(); err != nil {
			blog.Errorf("[datacollect][%s] collect message failed: %v, retry 3s later", p.name, err)
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
		if syserr := recover(); syserr != nil {
			blog.Errorf("[datacollect][%s] analyzeLoop panic by: %v, stack:\n %s", p.name, syserr, debug.Stack())
		}
	}()

	var mesg string
	var err error
	for mesg = range p.mesgC {
		if err = p.analyzer.Analyze(mesg); err != nil {
			blog.Errorf("[datacollect][%s] analyze message failed: %v, raw mesg: %s", p.name, err, mesg)
		}
	}
}

// collect 获取待处理消息，当是master时从redis channel获取，当是slave时从 redis queue 获取
func (p *chanPorter) collect() error {
	// 抢master锁
	err := loginMaster(p.redisCli, p.name, p.pid)
	if err != nil {
		// 抢失败，成为slave，开始从slave处理队列获取消息
		go p.popLoop()
		if strings.HasPrefix(err.Error(), "there is other master") {
			blog.Infof("[datacollect][%s] %v", p.name, err)
			return nil
		} else {
			blog.Errorf("[datacollect][%s] %v", p.name, err)
			return err
		}
	}

	// 抢成功，成为master，开始读redis channel，并推送处理不过来的消息到slave处理队列
	p.isMaster.Set()
	defer p.isMaster.UnSet()

	var wg = &sync.WaitGroup{}
	wg.Add(1)
	go pushLoop(p.redisCli, p.name, p.mesgC, p.isMaster, wg)

	err = p.subscribeLoop()
	if err != nil {
		return fmt.Errorf("subscribe channel return an error: %v", err)
	}

	// 读 redis channel 异常， 退出 master 状态
	p.isMaster.UnSet()
	err = logoutMaster(p.redisCli, p.name, p.pid)
	wg.Wait()

	return err
}

func (p *chanPorter) subscribeLoop() error {
	subChan, err := p.snapCli.Subscribe(p.channels...)
	if nil != err {
		return fmt.Errorf("subscribe channel failed, %v", err)
	}
	defer subChan.Unsubscribe(p.channels...)

	blog.Info("[datacollect][%s] subcribing channel %v from redis", p.name, p.channels)
	defer blog.Info("[datacollect][%s] unsubcribe channel %v from redis", p.name, p.channels)

	p.lastMesgTs = time.Now()
	var cnt int64
	for {
		received, err := subChan.Receive()
		if err == redis.Nil || err == io.EOF {
			continue
		}
		msg, ok := received.(*redis.Message)
		if !ok {
			continue
		}
		if nil != err {
			return fmt.Errorf("receive message from redis failed: %v", err)
		}

		if "" == msg.Payload {
			continue
		}

		writeOrClearChan(p.mesgC, p.name, msg.Payload)

		cnt++
		p.lastMesgTs = time.Now()
		if time.Since(p.lastMesgTs) > time.Minute {
			blog.Infof("[datacollect][%s] receive rate: %d/sec", p.name, int(float64(cnt)/time.Now().Sub(p.lastMesgTs).Seconds()))
			cnt = 0
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
		if syserr := recover(); syserr != nil {
			blog.Errorf("[datacollect][%s] panic by: %v, stack:\n %s", p.name, syserr, debug.Stack())
		}
	}()

	var err error
	for now := range ticker.C {
		channelstatus := 0
		if err := p.snapCli.Ping().Err(); err != nil {
			channelstatus = common.CCErrHostGetSnapshotChannelClose
			blog.Errorf("[datacollect][%s] snap redis server connection error: %s", p.name, err.Error())
		} else if err := p.redisCli.Ping().Err(); err != nil {
			channelstatus = common.CCErrHostGetSnapshotChannelClose
			blog.Errorf("[datacollect][%s] cc redis server connection error: %s", p.name, err.Error())
		} else if p.isMaster.IsSet() && now.Sub(p.lastMesgTs) > time.Minute {
			blog.Errorf("[datacollect][%s] snapchannel was empty in last 1 min", p.name)
			channelstatus = common.CCErrHostGetSnapshotChannelEmpty
		} else {
			channelstatus = common.CCSuccess
		}
		if err = p.redisCli.Set(channelStatusKey(p.name), channelstatus, time.Minute*2).Err(); err != nil {
			blog.Errorf("[datacollect][%s] set channelstatus failed: %v", err)
		}
	}
}

// popLoop 从slave处理队列获取消息，从而协助master处理
func (p *chanPorter) popLoop() {
	blog.Info("[datacollect][%s] start popLoop from redis", p.name)
	defer blog.Info("[datacollect][%s] stop popLoop from redis", p.name)

	// 加锁是为了防止执行到下面的 ```if p.isMaster.IsSet() ``` 即将退出而新协程判断 p.poping == true 也退出掉
	p.popLock.Lock()
	if p.poping {
		p.popLock.Unlock()
		return
	}
	p.poping = true
	p.popLock.Unlock()
	defer func() {
		// 防止panic后poping标志未重置
		p.popLock.Lock()
		p.poping = false
		p.popLock.Unlock()
	}()

	// 推消息到slave处理队列
	var mesg []string
	var err error
	key := slavequeueKey(p.name)
	for {
		p.popLock.Lock()
		if p.isMaster.IsSet() {
			// master 不需要从slave处理队列里取消息，所以退出
			p.poping = false
			return
		}
		p.popLock.Lock()
		mesg, err = p.redisCli.BRPop(time.Second*30, key).Result()
		if err == redis.Nil {
			continue
		}
		if err != nil {
			blog.Errorf("[datacollect][%s] pop message from redis failed: %v, retry 3s later", p.name, err)
			// 睡3秒，防止空跑导致CPU占用高涨
			time.Sleep(time.Second * 3)
		}
		if len(mesg) > 1 && mesg[1] != "nil" {
			writeOrClearChan(p.mesgC, p.name, mesg[1])
		}
	}
}

// pushLoop 把master处理不过来的消息推到slave处理队列，让slave协助处理
func pushLoop(redisCli *redis.Client, name string, mesgC chan string, isMaster *util.AtomicBool, wg *sync.WaitGroup) {
	defer wg.Done()

	blog.Info("[datacollect][%s] start pushLoop to redis", name)
	defer blog.Info("[datacollect][%s] stop pushLoop to redis", name)

	var mesg string
	var err error
	key := slavequeueKey(name)
	for {
		if !isMaster.IsSet() {
			// 当不再是master时就不需要再推消息到slave处理队列了
			return
		}
		select {
		case mesg = <-mesgC:
			if err = redisCli.LPush(key, mesg).Err(); err != nil {
				blog.Errorf("[datacollect][%s] push message to redis failed: %v", name, err)
			}
		default:
			blog.V(5).Infof("[datacollect][%s] pushLoop idle in last 3s", name)
			// 睡3秒，防止空跑导致CPU占用高涨
			time.Sleep(time.Second * 3)
		}
	}
}

// writeOrClearChan 利用非阻塞读channel达到清里channel的目的
func writeOrClearChan(mesgC chan string, name, mesg string) {
	select {
	case mesgC <- mesg:
	default:
		// channel fulled, so we drop 200 oldest events from queue
		blog.Infof("[datacollect][%s] msgChan full, len %d. drop 200 oldest from queue", name, len(mesgC))
		var ok bool
		for i := 0; i < 200; i-- {
			_, ok = <-mesgC
			if !ok {
				break
			}
		}
		mesgC <- mesg
	}
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
	return nil
}

// masterLockKey master锁的key
func masterLockKey(name string) string {
	return common.BKCacheKeyV3Prefix + name + ":masterlock"
}

// slavequeueKey 交给slave处理的消息待处理队列的key
func slavequeueKey(name string) string {
	return common.BKCacheKeyV3Prefix + name + ":queue"
}

// channelStatusKey 通道状态的key
func channelStatusKey(name string) string {
	return common.BKCacheKeyV3Prefix + name + ":channelstatus"
}
