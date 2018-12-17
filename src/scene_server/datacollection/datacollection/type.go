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
	"strings"
	"sync"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/datacollection/app/options"
	ccredis "configcenter/src/storage/dal/redis"

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
	Config() (ChanHandlerConfig, error)
}

type ChanHandlerConfig struct {
	Channels []string
	CCRedis  ccredis.Config
	SRCRedis ccredis.Config
}

func BuildChanPorter(name string, h ChanHandler, redisCli, snapCli *redis.Client) *chanPorter {
	return &chanPorter{
		h:        h,
		name:     name,
		pid:      xid.New().String(),
		isMaster: util.NewBool(false),
	}
}

type chanPorter struct {
	h        ChanHandler
	name     string
	pid      string
	isMaster *util.AtomicBool
	cfg      ChanHandlerConfig
	redisCli *redis.Client
	snapCli  *redis.Client
}

func (p *chanPorter) Run(conf options.Config) error {
	cfg, err := p.h.Config()
	if err != nil {
		return fmt.Errorf("get config error %v", err)
	}
	p.cfg = cfg

	blog.Infof("[datacollect][%s] connecting to redis: %+v", p.name, cfg.CCRedis)
	redisCli, err := ccredis.NewFromConfig(cfg.CCRedis)
	if err != nil {
		return fmt.Errorf("connect to redis failed: %v, cfg: %+v", err, cfg.CCRedis)
	}
	p.redisCli = redisCli

	mesgC := make(chan string, 1000)
	for i := 0; i < runtime.NumCPU(); i++ {
		go p.analyzeLoop(mesgC)
	}
	for {
		if err := p.Collect(mesgC); err != nil {
			blog.Errorf("[datacollect][%s] collect message failed: %v, retry 3s later", p.name, err)
		}
		time.Sleep(time.Second * 3)
	}
}

func (p *chanPorter) analyzeLoop(mesgC chan string) {
	var mesg string
	var err error
	for mesg = range mesgC {
		if err = p.h.Analyze(mesg); err != nil {
			blog.Errorf("[datacollect][%s]analyze message failed: %v, raw mesg: %s", p.name, err, mesg)
		}
	}
}

func (p *chanPorter) Collect(mesgC chan string) error {

	err = loginMaster(p.redisCli, p.name, p.pid)
	if err != nil {
		// this process is slave, so pop message from redis list
		go popLoop(p, redisCli, p.name, mesgC, p.isMaster)
		if strings.HasPrefix(err.Error(), "there is other master") {
			blog.Infof("[datacollect][%s] %v", p.name, err)
			return nil
		} else {
			blog.Errorf("[datacollect][%s] %v", p.name, err)
			return err
		}
	}
	p.isMaster.Set()
	defer p.isMaster.UnSet()

	// this process is slave, so receive message from redis channel
	blog.Infof("[datacollect][%s] connecting to redis: %+v", p.name, cfg.SRCRedis)
	snapCli, err := ccredis.NewFromConfig(cfg.SRCRedis)
	if err != nil {
		return fmt.Errorf("connect to redis failed: %v, cfg: %+v", err, cfg.SRCRedis)
	}

	var wg = &sync.WaitGroup{}
	go pushLoop(redisCli, p.name, mesgC, p.isMaster, wg)

	err = p.subscribeLoop(snapCli, cfg.Channels, mesgC)
	if err != nil {
		return fmt.Errorf("subscribe channel return an error: %v", err)
	}
	p.isMaster.UnSet()
	err = logoutMaster(redisCli, p.name, p.pid)
	wg.Wait()

	return err
}

func (p *chanPorter) subscribeLoop(snapCli *redis.Client, channels []string, mesgC chan string) error {

	subChan, err := snapCli.Subscribe(channels...)
	if nil != err {
		return fmt.Errorf("subscribe channel failed, %v", err)
	}
	defer subChan.Unsubscribe(channels...)

	blog.Info("[datacollect][%s] subcribing channel %v from redis", p.name, channels)
	defer blog.Info("[datacollect][%s] unsubcribe channel %v from redis", p.name, channels)

	var ts = time.Now()
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

		writeOrClearChan(mesgC, p.name, msg.Payload)

		cnt++
		if time.Since(ts) > time.Minute {
			blog.Infof("[datacollect][%s] receive rate: %d/sec", p.name, int(float64(cnt)/time.Now().Sub(ts).Seconds()))
			cnt = 0
			ts = time.Now()
		}
	}
}

func pushLoop(redisCli *redis.Client, name string, mesgC chan string, isMaster *util.AtomicBool, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	blog.Info("[datacollect][%s] start pushLoop to redis", name)
	defer blog.Info("[datacollect][%s] stop pushLoop to redis", name)

	var mesg string
	var err error
	queueKey := queueKey(name)
	for {
		if !isMaster.IsSet() {
			return
		}
		select {
		case mesg = <-mesgC:
			if err = redisCli.LPush(queueKey, mesg).Err(); err != nil {
				blog.Errorf("[datacollect][%s] push message to redis failed: %v", name, err)
			}
		default:
			time.Sleep(time.Second)
			blog.V(5).Infof("[datacollect][%s] pushLoop idle in last 1s", name)
		}
	}
}

func popLoop(redisCli *redis.Client, name string, mesgC chan string, isMaster *util.AtomicBool) {
	blog.Info("[datacollect][%s] start popLoop from redis", name)
	defer blog.Info("[datacollect][%s] stop popLoop from redis", name)

	var mesg []string
	var err error
	queueKey := queueKey(name)

	for {
		if isMaster.IsSet() {
			return
		}
		mesg, err = redisCli.BRPop(time.Second*30, queueKey).Result()
		if err == redis.Nil {
			continue
		}
		if err != nil {
			blog.Errorf("[datacollect][%s] pop message from redis failed: %v, retry 3s later", name, err)
			time.Sleep(time.Second * 3)
		}
		if len(mesg) > 1 && mesg[1] != "nil" {
			writeOrClearChan(mesgC, name, mesg[1])
		}
	}
}

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

func masterLockKey(name string) string {
	return common.BKCacheKeyV3Prefix + name + ":masterlock"
}

func queueKey(name string) string {
	return common.BKCacheKeyV3Prefix + name + ":queue"
}
