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
	"context"
	"io"
	"net/http"
	"runtime"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/xid"
	"gopkg.in/redis.v5"

	bkc "configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	httpcli "configcenter/src/common/http/httpclient"
)

type Discover struct {
	sync.Mutex

	redisCli *redis.Client
	subCli   *redis.Client

	id          string
	chanName    string
	ts          time.Time
	msgChan     chan string
	interrupt   chan error
	doneCh      chan struct{}
	resetHandle chan struct{}
	isMaster    bool
	isSubing    bool

	maxConcurrent          int
	maxSize                int
	getMasterInterval      time.Duration
	masterProcLockLiveTime time.Duration

	requests *httpcli.HttpClient
	*backbone.Engine
	ctx     context.Context
	pheader http.Header
}

var msgHandlerCnt = int64(0)

func NewDiscover(ctx context.Context, chanName string, maxSize int, redisCli, subCli *redis.Client, backbone *backbone.Engine) *Discover {

	if 0 == maxSize {
		maxSize = 100
	}
	pheader := http.Header{}
	pheader.Add(bkc.BKHTTPOwnerID, bkc.BKDefaultOwnerID)
	pheader.Add(bkc.BKHTTPHeaderUser, bkc.CCSystemCollectorUserName)

	discover := &Discover{
		chanName:               chanName,
		msgChan:                make(chan string, maxSize*4),
		interrupt:              make(chan error),
		resetHandle:            make(chan struct{}),
		doneCh:                 make(chan struct{}),
		maxSize:                maxSize,
		redisCli:               redisCli,
		subCli:                 subCli,
		ts:                     time.Now(),
		id:                     xid.New().String()[5:],
		maxConcurrent:          runtime.NumCPU(),
		getMasterInterval:      time.Second * 11,
		masterProcLockLiveTime: getMasterProcIntervalTime + time.Second*10,
		ctx:     ctx,
		pheader: pheader,
	}
	discover.Engine = backbone
	return discover
}

func (d *Discover) Start() {

	go func() {
		d.Run()

		for {
			time.Sleep(10 * time.Second)
			NewDiscover(d.ctx, d.chanName, d.maxSize, d.redisCli, d.subCli, d.Engine).Run()
		}
	}()
}

func (d *Discover) Run() {
	defer func() {
		if err := recover(); err != nil {
			blog.Errorf("fatal error happened: %s, we will try again 10s later, stack: \n%s", err, debug.Stack())
		}

		close(d.doneCh)
		d.isMaster = false

		return
	}()

	blog.Infof("discover start with maxConcurrent: %d", d.maxConcurrent)

	ticker := time.NewTicker(d.getMasterInterval)

	var err error
	var msg string
	var msgs []string
	var addCount, delayHandleCnt int

	if d.lockMaster() {
		blog.Infof("lock master success, start subscribe channel: %s", d.chanName)
		go d.subChan()
	} else {
		blog.Infof("master process exists, recheck after %v ", d.getMasterInterval)
	}

	for {
		select {
		case <-ticker.C:
			if d.lockMaster() {
				if !d.isSubing {
					blog.Infof("try to subscribe channel: %s", d.chanName)
					go d.subChan()
				}
			}
		case msg = <-d.msgChan:

			d.Lock()

			msgs = make([]string, 0, d.maxSize*2)
			msgs = append(msgs, msg)

			addCount = 0
			d.ts = time.Now()

		RLoop:

			for {
				select {
				case <-time.After(time.Second):
					break RLoop
				case msg = <-d.msgChan:
					blog.Infof("continue read 1s from channel: %d", addCount)
					addCount++
					msgs = append(msgs, msg)
					if addCount > d.maxSize {
						break RLoop
					}
				}
			}
			d.Unlock()

			delayHandleCnt = 0
			for {

				if delayHandleCnt > d.maxConcurrent*2 {
					blog.Warnf("msg process delay %d times, reset handlers", delayHandleCnt)
					close(d.resetHandle)
					d.resetHandle = make(chan struct{})

					delayHandleCnt = 0
				}

				if atomic.LoadInt64(&msgHandlerCnt) < int64(d.maxConcurrent) {

					atomic.AddInt64(&msgHandlerCnt, 1)
					blog.Infof("start message handler: %d", msgHandlerCnt)

					go d.handleMsg(msgs, d.resetHandle)

					break
				}

				delayHandleCnt++
				blog.Warnf("msg process delay again(%d times)", delayHandleCnt)

				time.Sleep(time.Millisecond * 100)
			}
		case err = <-d.interrupt:
			blog.Warnf("release master, msg process interrupted by: %s", err.Error())
			d.releaseMaster()
		}

	}
}

func (d *Discover) subChan() {
	defer func() {
		if err := recover(); err != nil {
			blog.Errorf("subChan fatal error happened %s, we will try again 10s later, stack: \n%s", err, debug.Stack())
		}
		d.isSubing = false
	}()

	d.isSubing = true

	subChan, err := d.subCli.Subscribe(d.chanName)
	if nil != err {
		d.interrupt <- err
		blog.Errorf("subscribe [%s] failed: %s", d.chanName, err.Error())
	}

	defer func() {
		subChan.Unsubscribe(d.chanName)
		d.isSubing = false
		blog.Infof("close subscribe channel: %s", d.chanName)
	}()

	var ts = time.Now()
	var cnt int64
	blog.Infof("start subscribe channel %s", d.chanName)

	for {

		if !d.isMaster {

			blog.Infof("i am not master, stop subscribe")
			return
		}

		received, err := subChan.Receive()

		if nil != err {

			if err == redis.Nil || err == io.EOF {
				continue
			}

			blog.Warnf("receive message err: %s", err.Error())
			d.interrupt <- err
			continue
		}

		msg, ok := received.(*redis.Message)
		if !ok || "" == msg.Payload {
			blog.Warnf("receive message failed(%v) or empty!", ok)
			continue
		}

		chanLen := len(d.msgChan)
		if d.maxSize*2 <= chanLen {

			blog.Infof("msgChan full, maxsize: %d, current: %d", d.maxSize, chanLen)
			d.clearOldMsg()
		}

		d.msgChan <- msg.Payload
		cnt++

		blog.Infof("send %d message to discover channel", cnt)

		if cnt%10000 == 0 {
			blog.Infof("receive rate: %d/sec", int(float64(cnt)/time.Now().Sub(ts).Seconds()))
			cnt = 0
			ts = time.Now()
		}
	}
}

func (d *Discover) clearOldMsg() {

	ts := d.ts
	msgCnt := len(d.msgChan) - d.maxSize

	blog.Warnf("start msgChan clear: %d", msgCnt)

	var cnt int
	for cnt < msgCnt {

		d.Lock()
		cnt++

		if ts != d.ts {
			msgCnt = len(d.msgChan) - d.maxSize
		} else {
			select {
			case <-time.After(time.Second * 10):
			case <-d.msgChan:
			}
		}

		d.Unlock()
	}

	if ts == d.ts {
		close(d.resetHandle)
	}

	blog.Warnf("msgChan cleared: %d", cnt)
}

func (d *Discover) releaseMaster() {

	val := d.redisCli.Get(MasterDisLockKey).Val()
	if val != d.id {
		d.redisCli.Del(MasterDisLockKey)
	}

	d.isMaster, d.isSubing = false, false
}

func (d *Discover) lockMaster() (ok bool) {
	var err error

	if d.isMaster {
		var val string
		val, err = d.redisCli.Get(MasterDisLockKey).Result()
		if err != nil {
			d.isMaster = false
			blog.Errorf("discover-master: lock master err %v", err)
		} else if val == d.id {
			blog.Infof("discover-master check : i am still master")
			d.redisCli.Set(MasterDisLockKey, d.id, d.masterProcLockLiveTime)
			ok = true
			d.isMaster = true
		} else {
			blog.Infof("discover-master: exit, val = %v, id = %v", val, d.id)
			d.isMaster = false
			ok = false
		}
	} else {
		ok, err = d.redisCli.SetNX(MasterDisLockKey, d.id, d.masterProcLockLiveTime).Result()
		if err != nil {
			d.isMaster = false
			blog.Errorf("discover-slave: lock master err %v", err)
		} else if ok {
			blog.Infof("discover-slave: check ok, i am master from now")
			d.isMaster = true
		} else {
			d.isMaster = false
			blog.Infof("discover-slave: check failed, there is other master process exists, recheck after %v ", d.getMasterInterval)
		}
	}

	return ok
}

func (d *Discover) handleMsg(msgs []string, resetHandle chan struct{}) error {

	defer atomic.AddInt64(&msgHandlerCnt, -1)

	blog.Infof("discover-master: handle %d num message, routines %d", len(msgs), atomic.LoadInt64(&msgHandlerCnt))

	for index, msg := range msgs {

		if msg == "" {
			continue
		}

		select {
		case <-resetHandle:
			blog.Warnf("reset handler, handled %d, set maxSize to %d ", index, d.maxSize)
			return nil
		case <-d.doneCh:
			blog.Warnf("close handler, handled %d")
			return nil
		default:

			err := d.TryCreateModel(msg)
			if err != nil {
				blog.Errorf("create model err: %s"+
					"##msg[%s]msg##", err, msg)
				continue
			}

			err = d.UpdateOrAppendAttrs(msg)
			if err != nil {
				blog.Errorf("create attr err: %s"+
					"##msg[%s]msg##", err, msg)
				continue
			}

			err = d.UpdateOrCreateInst(msg)
			if err != nil {
				blog.Errorf("create inst err: %s"+
					"##msg[%s]msg##", err, msg)
				continue
			}

			blog.Infof("==============[%d/%d] discover message finished", index, len(msgs))
		}

	}
	return nil
}
