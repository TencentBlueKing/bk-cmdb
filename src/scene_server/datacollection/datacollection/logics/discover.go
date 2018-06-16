package logics

import (
	"time"
	"sync"
	"sync/atomic"
	"gopkg.in/redis.v5"
	"github.com/rs/xid"

	"configcenter/src/scene_server/datacollection/common"
	"configcenter/src/common/blog"

	"io"
	"runtime"
	"fmt"
)

type Discover struct {
	sync.Mutex

	redisCli *redis.Client
	subCli   *redis.Client

	id          string
	chanName    string
	ts          time.Time // life cycle timestamp
	msgChan     chan string
	interrupt   chan error
	resetHandle chan struct{}
	isMaster    bool
	isSubing    bool
	cache       *DiscoverCache

	maxConcurrent             int
	maxSize                   int
	getMasterProcIntervalTime time.Duration
	masterProcLockLiveTime    time.Duration

	wg *sync.WaitGroup
}

type DiscoverInterface interface {
	HandleMsg(string) error
	Parse()
	Update()
}

type DiscoverCache struct {
	cache map[bool]map[string]map[string]interface{}
	flag bool
}

func NewDiscover(chanName string, maxSize int, redisCli, subCli *redis.Client, wg *sync.WaitGroup) *Discover {
	if 0 == maxSize {
		maxSize = 100
	}
	return &Discover{
		chanName:                  chanName,
		msgChan:                   make(chan string, maxSize*4),
		interrupt:                 make(chan error),
		resetHandle:               make(chan struct{}),
		maxSize:                   maxSize,
		redisCli:                  redisCli,
		subCli:                    subCli,
		ts:                        time.Now(),
		id:                        xid.New().String()[5:],
		maxConcurrent:             runtime.NumCPU(),
		getMasterProcIntervalTime: time.Second * 10,
		masterProcLockLiveTime:    getMasterProcIntervalTime + time.Second*10,
		wg:                        wg,
		cache: &DiscoverCache{
			cache: map[bool]map[string]map[string]interface{}{},
			flag: false,
		},
	}
}

// Start start main handle routines
func (d *Discover) Start() {
	defer d.wg.Done()

	//go d.fetchDB()
	go d.Run()
}

// Run hostsnap main functionality
func (d *Discover) Run() {

	blog.Infof("datacollection start with maxConcurrent: %d", d.maxConcurrent)

	ticker := time.NewTicker(d.getMasterProcIntervalTime)

	var err error
	var msg string
	var msgs []string
	var addCount, waitCnt int

	if d.saveRunning() {
		blog.Infof("saveRunning subChan")
		go d.subChan()
	} else {
		blog.Infof("run: there is other master process exists, recheck after %v ", d.getMasterProcIntervalTime)
	}

	for {
		select {
		case <-ticker.C:
			if d.saveRunning() {
				if !d.isSubing {
					blog.Infof("SELECT subChan")
					go d.subChan()
				}
			}
		case msg = <-d.msgChan:
			// read all from msgChan and lock to prevent clear operation
			d.Lock()
			d.ts = time.Now()
			msgs = make([]string, 0, d.maxSize*2)
			timeoutCh := time.After(time.Second)
			msgs = append(msgs, msg)
			addCount = 0
		f:
			for {
				select {
				case <-timeoutCh:
					break f
				case msg = <-d.msgChan:
					addCount++
					msgs = append(msgs, msg)
				}
				if addCount > d.maxSize {
					break f
				}
			}
			d.Unlock()

			// handle them
			waitCnt = 0
			for {
				if waitCnt > d.maxConcurrent*2 {
					blog.Warnf("reset handlers")
					close(d.resetHandle)
					d.resetHandle = make(chan struct{})
				}
				if atomic.LoadInt64(&routeCnt) < int64(d.maxConcurrent) {
					atomic.AddInt64(&routeCnt, 1)
					go d.handleMsg(msgs, d.resetHandle)
					break
				}
				waitCnt++
				time.Sleep(time.Millisecond * 100)
			}
		case err = <-d.interrupt:
			blog.Warn("interrupted", err.Error())
			d.concede()
		}

	}
}

// subChan subscribe message from redis channel
func (d *Discover) subChan() {

	d.isSubing = true

	subChan, err := d.subCli.Subscribe(d.chanName)
	if nil != err {
		d.interrupt <- err
		blog.Errorf("subscribe [%s] failed: %s", d.chanName, err.Error())
	}

	defer func() {
		subChan.Unsubscribe(d.chanName)
		d.isSubing = false
		blog.Infof("subChan Close: %s", d.chanName)
	}()

	var ts = time.Now()
	var cnt int64
	blog.Infof("isSubing channel %s", d.chanName)

	for {

		if false == d.isMaster {
			// not master again, close subscribe to prevent unnecessary subscribe
			blog.Info("This is not master process, subChan Close")
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
			continue
		}

		// todo 这是什么逻辑？
		chanLen := len(d.msgChan)
		if d.maxSize*2 <= chanLen {
			//  if msgChan fulled, clear old msgs
			blog.Infof("msgChan full, maxsize %d, len %d", d.maxSize, chanLen)
			d.clearMsgChan()
		}

		d.msgChan <- msg.Payload

		cnt++
		if cnt%10000 == 0 {
			blog.Infof("receive rate: %d/sec", int(float64(cnt)/time.Now().Sub(ts).Seconds()))
			cnt = 0
			ts = time.Now()
		}
	}
}

//clearMsgChan clear msgchan when msgchan is twice length of maxsize
func (d *Discover) clearMsgChan() {

	ts := d.ts
	msgCnt := len(d.msgChan) - d.maxSize

	blog.Warnf("start clear %d", msgCnt)

	var cnt int
	for msgCnt > cnt {
		d.Lock()

		cnt++

		// todo 不理解
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

	blog.Warnf("cleared %d", cnt)
}

// concede concede when buffer fulled
func (d *Discover) concede() {
	blog.Info("concede")
	d.isMaster = false
	d.isSubing = false
	val := d.redisCli.Get(common.MasterDisLockKey).Val()
	if val != d.id {
		d.redisCli.Del(common.MasterDisLockKey)
	}
}

// saveRunning lock master process
func (d *Discover) saveRunning() (ok bool) {
	var err error
	setNXChan := make(chan struct{})

	go func() {
		select {
		case <-time.After(d.masterProcLockLiveTime):
			blog.Fatalf("saveRunning check: set nx time out!! the network may be broken, redis stats: %v ", d.redisCli.PoolStats())
		case <-setNXChan:
		}
	}()

	if d.isMaster {
		var val string
		val, err = d.redisCli.Get(common.MasterDisLockKey).Result()
		if err != nil {
			blog.Errorf("master: saveRunning err %v", err)
		} else if val == d.id {
			blog.Infof("master check : i am still master")
			d.redisCli.Set(common.MasterDisLockKey, d.id, masterProcLockLiveTime)
			ok = true
		} else {
			blog.Infof("exit master,val = %v, id = %v", val, d.id)
			d.isMaster = false
			ok = false
		}
	} else {
		ok, err = d.redisCli.SetNX(common.MasterDisLockKey, d.id, d.masterProcLockLiveTime).Result()
		if err != nil {
			blog.Errorf("slave: saveRunning err %v", err)
		} else if ok {
			blog.Infof("slave: check ok, i am master from now")
			d.isMaster = true
		} else {
			blog.Infof("slave: check failed, there is other master process exists, recheck after %v ", d.getMasterProcIntervalTime)
		}
	}

	close(setNXChan)

	return ok
}

func (d *Discover) handleMsg(msgs []string, resetHandle chan struct{}) error {
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
			blog.Warnf("reset handler, handled %d, set maxSize to %d ", index, d.maxSize)
			return nil
		default:
			fmt.Printf("[%d/%d] datacollect message: %s", index, len(msgs), msg)
		}
	}

	return nil
}

func (d *Discover) HandleMsg(string) error {
	//body := `{"bk_inst_name":"asdfasdf1","bk_obj_id":"man","bk_supplier_account":"0","create_time":"2018-06-16 22:40:31","default":0,"name":"asdfasdf1"}`
	//
	////gHostAttrURL := "http://" + cli.CC.ObjCtrl + "/object/v1/meta/objectatts"
	//gHostAttrURL := cli.CC.ObjCtrl() + "/object/v1/meta/objectatts"
	//searchBody := make(map[string]interface{})
	//searchBody[bkcommon.BKObjIDField] = bkcommon.BKInnerObjIDHost
	//searchBody[bkcommon.BKOwnerIDField] = ownerID
	//searchJson, _ := json.Marshal(searchBody)
	//gHostAttrRe, err := httpcli.ReqHttp(req, gHostAttrURL, bkcommon.HTTPSelectPost, []byte(searchJson))
	//if nil != err {
	//	blog.Error("GetHostDetailByID  attr error :%v", err)
	//	return http.StatusInternalServerError, nil, defErr.Error(bkcommon.CCErrHostDetailFail)
	//}
	return nil
}
