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
	"encoding/json"
	"io"
	"net/http"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal"

	"github.com/rs/xid"
	redis "gopkg.in/redis.v5"
)

type Netcollect struct {
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
	db    dal.RDB

	routeCnt   int64
	handleCnt  int64
	handlelock sync.Mutex
	handlets   time.Time

	*backbone.Engine
	ctx     context.Context
	pheader http.Header

	wg *sync.WaitGroup
}

func NewNetcollect(ctx context.Context, chanName []string, maxSize int, redisCli, snapCli *redis.Client, db dal.RDB, backbone *backbone.Engine) *Netcollect {
	if 0 == maxSize {
		maxSize = 100
	}

	pheader := http.Header{}
	pheader.Add(common.BKHTTPOwnerID, common.BKDefaultOwnerID)
	pheader.Add(common.BKHTTPHeaderUser, common.CCSystemCollectorUserName)

	NetDeviceInstance := &Netcollect{
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
		routeCnt:      int64(0),
		handleCnt:     int64(0),
		handlelock:    sync.Mutex{},
		handlets:      time.Now(),

		ctx:     ctx,
		Engine:  backbone,
		pheader: pheader,
	}
	return NetDeviceInstance
}

func (h *Netcollect) Start() {

	go func() {
		h.Run()
		for {
			time.Sleep(time.Second * 10)
			NewNetcollect(h.ctx, h.hostChanName, h.maxSize, h.redisCli, h.snapCli, h.db, h.Engine).Run()
		}
	}()
}

func (h *Netcollect) Run() {
	defer func() {
		syserr := recover()
		if syserr != nil {
			blog.Errorf("[NetDevice] emergency error happened %s, we will try again 10s later, stack: \n%s", syserr, debug.Stack())
		}
		close(h.doneCh)
		h.isMaster = false
		return
	}()
	blog.Infof("[NetDevice] datacollection start with maxconcurrent: %d", h.maxconcurrent)
	ticker := time.NewTicker(getMasterProcIntervalTime)
	var err error
	var msg string
	var msgs []string
	var addCount int
	var waitCnt int

	if h.saveRunning() {
		go h.subChan(h.snapCli, h.hostChanName)
	} else {
		blog.Infof("[NetDevice] run: there is other master process exists, recheck after %v ", getMasterProcIntervalTime)
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
					blog.Warnf("[NetDevice] reset handlers")
					close(h.resetHandle)
					waitCnt = 0
					h.resetHandle = make(chan struct{})
				}
				if atomic.LoadInt64(&routeCnt) < int64(h.maxconcurrent) {
					atomic.AddInt64(&routeCnt, 1)
					h.handleMsg(msgs, h.resetHandle)
					break
				}
				waitCnt++
				time.Sleep(time.Millisecond * 100)
			}
		case err = <-h.interrupt:
			blog.Warn("[NetDevice] interrupted", err.Error())
			h.concede()
		}

	}
}

func (h *Netcollect) handleMsg(msgs []string, resetHandle chan struct{}) error {
	defer func() {
		syserr := recover()
		if syserr != nil {
			blog.Errorf("[NetDevice] handleMsg emergency error happened %s [recovered] stack: \n%s", syserr, debug.Stack())
		}
	}()
	defer atomic.AddInt64(&routeCnt, -1)
	blog.Infof("[NetDevice] handle %d num mesg, routines %d", len(msgs), atomic.LoadInt64(&routeCnt))
	for index, raw := range msgs {
		if raw == "" {
			continue
		}
		handlelock.Lock()
		handleCnt++
		if handleCnt%10000 == 0 {
			blog.Infof("[NetDevice] handle rate: %d/sec", int(float64(handleCnt)/time.Now().Sub(handlets).Seconds()))
			handleCnt = 0
			handlets = time.Now()
		}
		handlelock.Unlock()
		select {
		case <-resetHandle:
			blog.Warnf("[NetDevice] reset handler, handled %d, set maxSize to %d ", index, h.maxSize)
			return nil
		case <-h.doneCh:
			blog.Warnf("[NetDevice] close handler, handled %d")
			return nil
		default:
			blog.V(4).Infof("[NetDevice] received message: %s", raw)
			msg := NetcollectMessage{}
			err := json.Unmarshal([]byte(raw), &msg)
			if err != nil {
				blog.Errorf("[NetDevice] unmarshal message error: %v, raw: %s", err, raw)
				continue
			}

			for _, report := range msg.Data {
				if err = h.handleReport(&report); err != nil {
					blog.Errorf("handleData failed: %v,raw: %s", err, raw)
				}
			}
		}
	}

	return nil
}

func (h *Netcollect) handleReport(report *metadata.NetcollectReport) (err error) {
	// TODO compare 若有变化才插入
	if err = h.upsertReport(report); err != nil {
		blog.Errorf("[NetDevice] upsert association error: %v", err)
		return err
	}

	return nil
}

func buildReport(metric *NetcollectMetric) *metadata.NetcollectReport {
	report := metadata.NetcollectReport{}
	report.CloudID = metric.CloudID
	switch metric.ObjectID {
	case common.BKInnerObjIDSwitch:
	}
	return &report
}

func (h *Netcollect) upsertReport(report *metadata.NetcollectReport) error {
	existCond := condition.CreateCondition()
	existCond.Field(common.BKCloudIDField).Eq(report.CloudID)
	existCond.Field(common.BKObjIDField).Eq(report.ObjectID)
	existCond.Field(common.BKInstKeyField).Eq(report.InstKey)

	count, err := h.db.Table(common.BKTableNameNetcollectReport).Find(existCond.ToMapStr()).Count(h.ctx)
	if err != nil {
		return err
	}
	if count <= 0 {
		err = h.db.Table(common.BKTableNameNetcollectReport).Insert(h.ctx, report)
		return err
	}

	return h.db.Table(common.BKTableNameNetcollectReport).Update(h.ctx, existCond, report)
}

func (h *Netcollect) findInst(objectID string, query *metadata.QueryInput) ([]mapstr.MapStr, error) {
	resp, err := h.CoreAPI.ObjectController().Instance().SearchObjects(h.ctx, objectID, h.pheader, query)
	if err != nil {
		blog.Infof("[NetDevice] findInst error: %v", err)
		return nil, err
	}
	if !resp.Result {
		blog.Infof("[NetDevice] findInst error: %v", resp.ErrMsg)
		return nil, err
	}
	return resp.Data.Info, nil
}

func (h *Netcollect) concede() {
	blog.Info("[NetDevice] concede")
	h.isMaster = false
	h.subscribing = false
	val := h.redisCli.Get(MasterNetLockKey).Val()
	if val != h.id {
		h.redisCli.Del(MasterNetLockKey)
	}
}

func (h *Netcollect) saveRunning() (ok bool) {
	var err error
	if h.isMaster {
		var val string
		val, err = h.redisCli.Get(MasterNetLockKey).Result()
		if err != nil {
			blog.Errorf("[NetDevice] master: saveRunning err %v", err)
			h.isMaster = false
		} else if val == h.id {
			blog.Infof("[NetDevice] master check : i am still master")
			h.redisCli.Set(MasterNetLockKey, h.id, masterProcLockLiveTime)
			ok = true
			h.isMaster = true
		} else {
			blog.Infof("[NetDevice] exit master,val = %v, id = %v", val, h.id)
			h.isMaster = false
			ok = false
		}
	} else {
		ok, err = h.redisCli.SetNX(MasterNetLockKey, h.id, masterProcLockLiveTime).Result()
		if err != nil {
			blog.Errorf("[NetDevice] slave: saveRunning err %v", err)
			h.isMaster = false
		} else if ok {
			blog.Infof("[NetDevice] slave check: ok")
			blog.Infof("[NetDevice] i am master from now")
			h.isMaster = true
		} else {
			blog.Infof("[NetDevice] slave check: there is other master process exists, recheck after %v ", getMasterProcIntervalTime)
			h.isMaster = false
		}
	}
	return ok
}

func (h *Netcollect) subChan(snapcli *redis.Client, chanName []string) {
	defer func() {
		syserr := recover()
		if syserr != nil {
			blog.Errorf("[NetDevice] subChan emergency error happened %s, we will try again 10s later, stack: \n%s", syserr, debug.Stack())
		}
		h.subscribing = false
	}()
	h.subscribing = true
	var chanlen int
	subChan, err := snapcli.Subscribe(chanName...)
	if nil != err {
		h.interrupt <- err
		blog.Error("[NetDevice] subscribe channel faile ", err.Error())
	}
	closeChan := make(chan struct{})
	go h.healthCheck(closeChan)
	defer func() {
		h.subscribing = false
		close(closeChan)
		blog.Infof("[NetDevice] subChan Close")
		subChan.Unsubscribe(chanName...)
	}()

	var ts = time.Now()
	var cnt int64
	blog.Infof("[NetDevice] subcribing channel %v", chanName)
	for {
		if false == h.isMaster {

			blog.Info("[NetDevice] This is not master process, subChan Close")
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
			blog.Debug("[NetDevice] receive messave  err", err.Error())
			h.interrupt <- err
			continue
		}

		if "" == msg.Payload {
			continue
		}

		chanlen = len(h.msgChan)
		if h.maxSize*2 <= chanlen {

			blog.Infof("[NetDevice] msgChan full, maxsize %d, len %d", h.maxSize, chanlen)
			h.clearMsgChan()
		}
		if chanlen != 0 && chanlen%10 == 0 {
			blog.Infof("[NetDevice] buff len %d", chanlen)
		}
		h.lastMesgTs = time.Now()
		h.msgChan <- msg.Payload
		cnt++
		if cnt%10000 == 0 {
			blog.Infof("[NetDevice] receive rate: %d/sec", int(float64(cnt)/time.Now().Sub(ts).Seconds()))
			cnt = 0
			ts = time.Now()
		}
	}
}

func (h *Netcollect) clearMsgChan() {
	ts := h.ts
	msgCnt := len(h.msgChan) - h.maxSize
	blog.Warnf("[NetDevice] start clear %d", msgCnt)
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
	blog.Warnf("[NetDevice] cleared %d", cnt)
}

func (h *Netcollect) healthCheck(closeChan chan struct{}) {
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
				blog.Errorf("[NetDevice] snap redis server connection error: %s", err.Error())
			} else if time.Now().Sub(h.lastMesgTs) > time.Minute {
				blog.Errorf("[NetDevice] %v was empty in last 1 min ", h.hostChanName)
				channelstatus = common.CCErrHostGetSnapshotChannelEmpty
			} else {
				channelstatus = common.CCSuccess
			}
			h.redisCli.Set(RedisNetcollectKeyChannelStatus, channelstatus, time.Minute*2)
		}
	}
}

type NetcollectMessage struct {
	Timestamp time.Time                   `json:"timestamp"`
	Dataid    int                         `json:"dataid"`
	Type      string                      `json:"type"`
	Counter   int                         `json:"counter"`
	Build     CollectorBuild              `json:"build"`
	Data      []metadata.NetcollectReport `json:"data"`
}

type CollectorBuild struct {
	Version     string `json:"version"`
	buildCommit string `json:"build_commit"`
	buildTime   string `json:"build_time"`
	goVersion   string `json:"go_version"`
}

type NetcollectMetric struct {
	CloudID      int64                                  `json:"bk_cloud_id"`
	ObjectID     string                                 `json:"bk_obj_id"`
	Attributes   []metadata.NetcollectReportAttribute   `json:"attributes"`
	Associations []metadata.NetcollectReportAssociation `json:"associations"`
}

const netcollectMockMsg = `{
    "dataid": 1014,
    "type": "netdevicebeat",
    "counter": 1,
    "Build": {
        "version": "1.0.0",
        "build_commit": "3fb6cb0b5a55cffae028d3df7bee71f90155a2f5",
        "buildtime": "2018-10-03 17:09:00",
        "go_version": "1.11.2"
    },
    "data": [
        {
            "bk_obj_id": "bk_switch",
            "bk_inst_key": "huawei 5789#56-79-9a-ii",
            "bk_host_innerip": "192.168.1.1",
			"bk_cloud_id": 0,
			"last_time": "2018-10-03 17:09:00",
            "attributes": [
                {
                    "bk_property_id": "bk_inst_name",
                    "value": "huawei 5789#56-79-9a-ii"
                }
            ],
            "associations": [
				{
					"bk_asst_inst_name": "192.168.1.1",
                    "bk_asst_obj_id": "bk_host",
                    "bk_asst_obj_name": "主机",
                    "bk_asst_property_id": "bk_host_id"
				}
			]
        }
    ]
}
`
