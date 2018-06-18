package logics

import (
	"fmt"
	"time"
	"sync"
	"sync/atomic"
	"gopkg.in/redis.v5"
	"github.com/rs/xid"
	"io"
	"runtime"
	"encoding/json"
	"github.com/tidwall/gjson"

	"configcenter/src/scene_server/datacollection/common"
	"configcenter/src/common/blog"
	bkc "configcenter/src/common"
	"configcenter/src/common/core/cc/api"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/framework/core/log"
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
	httpClient.SetHeader(bkc.BKHTTPHeaderUser, "collector")

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
		getMasterInterval:      time.Second * 10,
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

	//go d.fetchDB()
	go d.Run()
}

// Run discover main functionality
func (d *Discover) Run() {

	blog.Infof("datacollection start with maxConcurrent: %d", d.maxConcurrent)

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

			blog.Infof("[%v]: read messages before: %d", d.ts, len(msgs))

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

				blog.Infof("read messages after: %d", len(msgs))

				// 延迟处理的次数超过一定程度？
				if delayHandleCnt > d.maxConcurrent*2 {
					blog.Warnf("msg process delay %d times, reset handlers", delayHandleCnt)
					close(d.resetHandle)
					d.resetHandle = make(chan struct{})

					// todo 延迟处理计数清零？
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
		blog.Infof("start receive message: %v\n", received)
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
			blog.Warnf("receive message failed or empty!\n")
			continue
		}

		// todo 生产者生产消息速度大于消费者，自动清理超出的历史消息？？
		chanLen := len(d.msgChan)
		if d.maxSize*2 <= chanLen {
			//  if msgChan fulled, clear old msgs
			blog.Infof("msgChan full, maxsize: %d, current: %d", d.maxSize, chanLen)
			d.clearOldMsg()
		}

		d.msgChan <- msg.Payload
		cnt++
		blog.Infof("send %dth message to msgChan", cnt)

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

		// todo 清理时，若发生新的消息写入，则重新获取消息数量？
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

	// todo 确认最终清理完毕（清理时间等于最后一次的消息写入时间）
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
			blog.Errorf("master: lockMaster err %v", err)
		} else if val == d.id {
			blog.Infof("master check : i am still master")
			d.redisCli.Set(common.MasterDisLockKey, d.id, d.masterProcLockLiveTime)
			ok = true
		} else {
			blog.Infof("exit master,val = %v, id = %v", val, d.id)
			d.isMaster = false
			ok = false
		}
	} else {
		ok, err = d.redisCli.SetNX(common.MasterDisLockKey, d.id, d.masterProcLockLiveTime).Result()
		if err != nil {
			blog.Errorf("slave: lockMaster err %v", err)
		} else if ok {
			blog.Infof("slave: check ok, i am master from now")
			d.isMaster = true
		} else {
			blog.Infof("slave: check failed, there is other master process exists, recheck after %v ", d.getMasterInterval)
		}
	}

	close(setNXChan)

	return ok
}

func (d *Discover) handleMsg(msgs []string, resetHandle chan struct{}) error {

	defer atomic.AddInt64(&msgHandlerCnt, -1)

	blog.Infof("handle %d num message, routines %d", len(msgs), atomic.LoadInt64(&msgHandlerCnt))

	for index, msg := range msgs {

		if msg == "" {
			continue
		}

		select {
		case <-resetHandle:
			blog.Warnf("reset handler, handled %d, set maxSize to %d ", index, d.maxSize)
			return nil
		default:

			// todo 1- create model and attr
			err := d.CreateModel(msg)
			if err != nil {
				blog.Errorf("create model err: %s\n"+
					"##msg[%s]msg##\n", err, msg)
				continue
			}

			// todo 2- create inst
			err = d.UpdateOrCreateInst(msg)
			if err != nil {
				blog.Errorf("create inst err: %s\n"+
					"##msg[%s]msg##\n", err, msg)
				continue
			}
			blog.Infof("==============\n[%d/%d] datacollect message: %s\n==================\n", index, len(msgs), msg)
		}

	}

	return nil
}

type Model struct {
	BkClassificationID string `json:"bk_classification_id"`
	BkObjID            string `json:"bk_obj_id"`
	BkObjName          string `json:"bk_obj_name"`
	Keys               string `json:"keys"`
}

type Field struct {
	BkPropertyName string `json:"bk_property_name"`
	BkPropertyType string `json:"bk_property_type"`
}

func (d *Discover) parseModel(msg string) (model *Model, err error) {

	//create model
	model = &Model{}
	modelStr := gjson.Get(msg, "data.meta.model").String()

	if err = json.Unmarshal([]byte(modelStr), &model); err != nil {
		blog.Errorf("unmarshal error: %s", err)
		return
	}

	return
}

func (d *Discover) CreateModel(msg string) error {

	//create model
	model, err := d.parseModel(msg)
	if err != nil {
		return fmt.Errorf("parse model error: %s", err.Error())
	}

	body := map[string]interface{}{
		bkc.BKClassificationIDField: model.BkClassificationID,
		bkc.BKObjIDField:            model.BkObjID,
		bkc.BKObjNameField:          model.BkObjName,
		bkc.BKOwnerIDField:          bkc.BKDefaultOwnerID,
		bkc.BKObjIconField:          "icon-cc-middleware",
		bkc.CreatorField:            bkc.CCSystemOperatorUserName,
	}

	bodyJs, _ := json.Marshal(body)
	blog.Infof("create model body: %s", bodyJs)
	res, err := d.requests.POST(fmt.Sprintf("%s/topo/v1/object", d.cc.TopoAPI()), nil, []byte(bodyJs))
	if nil != err {
		return fmt.Errorf("create model error: %s", err.Error())
	}
	blog.Infof("create model result: %s\n", res)

	// TODO verify res
	//js, err := simplejson.NewJson([]byte(res))
	//jsMap, _ := js.Map()

	//create model attr
	fields := make(map[string]Field, 0)
	fieldsStr := gjson.Get(msg, "data.meta.fields").String()
	blog.Infof("create model attr fieldsStr: %s\n", fieldsStr)
	if err := json.Unmarshal([]byte(fieldsStr), &fields); err != nil {
		blog.Errorf("create model attr unmarshal error: %s", err)
		return err
	}

	// batch create model attrs
	for k, v := range fields {
		blog.Infof("attr: %s -> %v\n", k, v)

		// skip default field
		if model.BkObjID == bkc.BKInstNameField {
			log.Infof("skip default field: %s", model.BkObjID)
			continue
		}

		attrBody := map[string]interface{}{
			bkc.BKObjIDField:         model.BkObjID,
			bkc.BKPropertyGroupField: bkc.BKDefaultField,
			bkc.BKPropertyIDField:    k,
			bkc.BKPropertyNameField:  v.BkPropertyName,
			bkc.BKPropertyTypeField:  v.BkPropertyType,
			bkc.BKOwnerIDField:       bkc.BKDefaultOwnerID,
		}

		attrBodyJs, _ := json.Marshal(attrBody)
		blog.Infof("create model attr body: %s", attrBody)

		res, err = d.requests.POST(fmt.Sprintf("%s/topo/v1/objectattr", d.cc.TopoAPI()), nil, []byte(attrBodyJs))
		if nil != err {
			return fmt.Errorf("create model attr error: %s", err.Error())
		}

		blog.Infof("create model attr result: %s\n", res)
	}

	return nil
}

func (d *Discover) UpdateOrCreateInst(msg string) error {

	//create inst
	bodyJs := gjson.Get(msg, "data.data").String()
	objID := gjson.Get(msg, "data.meta.model.bk_obj_id").String()
	blog.Infof("create model body: %s", bodyJs)

	res, err := d.requests.POST(fmt.Sprintf("%s/topo/v1/inst/%s/%s", d.cc.TopoAPI(), bkc.BKDefaultOwnerID, objID), nil, []byte(bodyJs))
	if nil != err {
		return fmt.Errorf("create inst error: %s", err.Error())
	}
	blog.Infof("create inst result: %s\n", res)

	return nil
}
