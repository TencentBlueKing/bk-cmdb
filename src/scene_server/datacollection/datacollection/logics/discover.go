package logics

import (
	"time"
	"sync"
	"sync/atomic"
	"gopkg.in/redis.v5"
	"github.com/rs/xid"
	"io"
	"runtime"
	"configcenter/src/scene_server/datacollection/common"
	"configcenter/src/common/blog"
	bkc "configcenter/src/common"
	"configcenter/src/common/core/cc/api"
	httpcli "configcenter/src/common/http/httpclient"
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

	maxConcurrent          int
	maxSize                int
	getMasterInterval      time.Duration
	masterProcLockLiveTime time.Duration

	requests *httpcli.HttpClient
	cc       *api.APIResource
	wg       *sync.WaitGroup
}

type DiscoverCache struct {
	cache map[bool]map[string]map[string]interface{}
	flag bool
}

var msgHandlerCnt = int64(0)

func NewDiscover(chanName string, maxSize int, redisCli, subCli *redis.Client, wg *sync.WaitGroup, cc *api.APIResource) *Discover {

	if 0 == maxSize {
		maxSize = 100
	}

	httpClient := httpcli.NewHttpClient()
	httpClient.SetHeader("Content-Type", "application/json")
	httpClient.SetHeader("Accept", "application/json")
	httpClient.SetHeader(bkc.BKHTTPOwnerID, bkc.BKDefaultOwnerID)
	httpClient.SetHeader(bkc.BKHTTPHeaderUser, bkc.CCSystemCollectorUserName)

	return &Discover{
		chanName:               chanName,
		msgChan:                make(chan string, maxSize*4),
		interrupt:              make(chan error),
		resetHandle:            make(chan struct{}),
		maxSize:                maxSize,
		redisCli:               redisCli,
		subCli:                 subCli,
		ts:                     time.Now(),
		id:                     xid.New().String()[5:],
		maxConcurrent:          runtime.NumCPU(),
		getMasterInterval:      time.Second * 11,
		masterProcLockLiveTime: getMasterProcIntervalTime + time.Second*10,
		wg:                     wg,
		cc:                     cc,
		requests:               httpClient,
		cache: &DiscoverCache{
			cache: map[bool]map[string]map[string]interface{}{},
			flag: false,
		},
	}
}

// Start start main handle routines
func (d *Discover) Start() {
	defer d.wg.Done()

	go d.Run()
}

// Run discover main functionality
func (d *Discover) Run() {

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

	// 尝试成为master/订阅消息并处理
	for {
		select {
		case <-ticker.C:
			if d.lockMaster() {
				if !d.isSubing {
					blog.Infof("try to subscribe channel: %s\n", d.chanName)
					go d.subChan()
				}
			}
		case msg = <-d.msgChan:
			// read all from msgChan and lock to prevent clear operation
			d.Lock()

			msgs = make([]string, 0, d.maxSize*2)
			msgs = append(msgs, msg)

			addCount = 0
			d.ts = time.Now()

		f:
		// 持续读取1s通道内的消息，最多读取d.maxSize个
			for {
				select {
				case <-time.After(time.Second):
					break f
				case msg = <-d.msgChan:
					blog.Infof("continue read 1s from channel: %d", addCount)
					addCount++
					msgs = append(msgs, msg)
					if addCount > d.maxSize {
						break f
					}
				}
			}
			d.Unlock()

			// 消息处理逻辑？
			delayHandleCnt = 0
			for {

				// 延迟处理的次数超过一定程度？
				if delayHandleCnt > d.maxConcurrent*2 {
					blog.Warnf("msg process delay %d times, reset handlers", delayHandleCnt)
					close(d.resetHandle)
					d.resetHandle = make(chan struct{})

					// 延迟处理计数清零
					//delayHandleCnt = 0
				}

				if atomic.LoadInt64(&msgHandlerCnt) < int64(d.maxConcurrent) {

					atomic.AddInt64(&msgHandlerCnt, 1)
					blog.Infof("start message handler: %d", msgHandlerCnt)

					go d.handleMsg(msgs, d.resetHandle)

					break
				}

				// 消息处理进程数超限，延迟处理
				delayHandleCnt++
				blog.Warnf("msg process delay again(%d times)\n", delayHandleCnt)

				time.Sleep(time.Millisecond * 100)
			}
		case err = <-d.interrupt:
			blog.Warnf("release master, msg process interrupted by: %s\n", err.Error())
			d.releaseMaster()
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
		blog.Infof("close subscribe channel: %s", d.chanName)
	}()

	var ts = time.Now()
	var cnt int64
	blog.Infof("start subscribe channel %s", d.chanName)

	for {

		if !d.isMaster {
			// not master again, close subscribe to prevent unnecessary subscribe
			blog.Infof("i am not master, stop subscribe\n")
			return
		}

		received, err := subChan.Receive()
		//blog.Debug("start receive message: %v\n", received)
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
			blog.Warnf("receive message failed(%v) or empty!\n", ok)
			continue
		}

		// 生产者生产消息速度大于消费者，自动清理超出的历史消息
		chanLen := len(d.msgChan)
		if d.maxSize*2 <= chanLen {
			//  if msgChan fulled, clear old msgs
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

//clearOldMsg clear old message when msgChan is twice length of maxsize
func (d *Discover) clearOldMsg() {

	ts := d.ts
	msgCnt := len(d.msgChan) - d.maxSize

	blog.Warnf("start msgChan clear: %d\n", msgCnt)

	var cnt int
	for cnt < msgCnt {

		d.Lock()
		cnt++

		// 清理时，若发生新的消息写入，则重新获取消息数量？
		if ts != d.ts {
			blog.Infof("clearOldMsg")
			msgCnt = len(d.msgChan) - d.maxSize
		} else {
			select {
			case <-time.After(time.Second * 10):
			case <-d.msgChan:
			}
		}

		d.Unlock()
	}

	// 确认最终清理完毕（清理时间等于最后一次的消息写入时间）
	if ts == d.ts {
		close(d.resetHandle)
	}

	blog.Warnf("msgChan cleared: %d\n", cnt)
}

// releaseMaster releaseMaster when buffer fulled
func (d *Discover) releaseMaster() {

	val := d.redisCli.Get(common.MasterDisLockKey).Val()
	if val != d.id {
		d.redisCli.Del(common.MasterDisLockKey)
	}

	d.isMaster, d.isSubing = false, false
}

// lockMaster lock master process
func (d *Discover) lockMaster() (ok bool) {
	var err error
	setNXChan := make(chan struct{})

	go func() {
		select {
		case <-time.After(d.masterProcLockLiveTime):
			blog.Fatalf("lockMaster check: set nx time out!! the network may be broken, redis stats: %v ", d.redisCli.PoolStats())
		case <-setNXChan:
		}
	}()

	if d.isMaster {
		var val string
		val, err = d.redisCli.Get(common.MasterDisLockKey).Result()
		if err != nil {
			blog.Errorf("discover-master: lock master err %v", err)
		} else if val == d.id {
			blog.Infof("discover-master check : i am still master")
			d.redisCli.Set(common.MasterDisLockKey, d.id, d.masterProcLockLiveTime)
			ok = true
		} else {
			blog.Infof("discover-master: exit, val = %v, id = %v", val, d.id)
			d.isMaster = false
			ok = false
		}
	} else {
		ok, err = d.redisCli.SetNX(common.MasterDisLockKey, d.id, d.masterProcLockLiveTime).Result()
		if err != nil {
			blog.Errorf("discover-slave: lock master err %v", err)
		} else if ok {
			blog.Infof("discover-slave: check ok, i am master from now")
			d.isMaster = true
		} else {
			blog.Infof("discover-slave: check failed, there is other master process exists, recheck after %v ", d.getMasterInterval)
		}
	}

	close(setNXChan)

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
		default:

			// 1- try create model
			err := d.TryCreateModel(msg)
			if err != nil {
				blog.Errorf("create model err: %s\n"+
					"##msg[%s]msg##\n", err, msg)
				continue
			}

			// 2- try create model attr
			err = d.UpdateOrAppendAttrs(msg)
			if err != nil {
				blog.Errorf("create attr err: %s\n"+
					"##msg[%s]msg##\n", err, msg)
				continue
			}

			// 3- create inst
			err = d.UpdateOrCreateInst(msg)
			if err != nil {
				blog.Errorf("create inst err: %s\n"+
					"##msg[%s]msg##\n", err, msg)
				continue
			}

			blog.Infof("==============\n[%d/%d] discover message finished\n", index, len(msgs))
		}

	}

	return nil
}
