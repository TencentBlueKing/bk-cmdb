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

package distribution

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	ccjson "configcenter/src/common/json"
	"configcenter/src/common/metadata"
	"configcenter/src/common/watch"
	etypes "configcenter/src/scene_server/event_server/types"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/storage/stream/types"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
)

const (
	// defaultTransTimeout is default timeout for trans event data from base queue to identifier
	// duplicate queue.
	defaultTransTimeout = 60 * time.Second

	// defaultHandleTimeout is default timeout for event handle.
	defaultHandleTimeout = time.Second

	// defaultFusingEventExpireSec is default expire second num for event.
	defaultFusingEventExpireSec = 3 * 60 * 60

	// defaultCleanCheckInterval is default interval for cleaning check.
	defaultCleanCheckInterval = time.Minute

	// defaultPushersCheckInterval is default interval for pushers check.
	defaultPushersCheckInterval = 10 * time.Second

	// defaultCleanTrimThreshold is default threshold for cleaning trim.
	defaultCleanTrimThreshold = 2 * 10000

	// defaultCleanDelThreshold is default threshold for cleaning del.
	defaultCleanDelThreshold = 5 * 10000

	// defaultCleanUnit is default clean unit count for cleaning.
	defaultCleanUnit = 10000
)

// EventHandler manages all event pushers, and update pushers in dynamic mode,
// when there are events need to be sent, the pusher would check master state and send
// message to subscriber in callback or not.
type EventHandler struct {
	ctx    context.Context
	engine *backbone.Engine

	// cache is cc redis client.
	cache redis.Client

	// pushers is local event pushers, update in dynamic mode, subid -> EventPushers.
	pushers map[int64]*EventPusher

	// pushersMu make pushers ops safe.
	pushersMu sync.RWMutex

	// distributer handles all events distribution.
	distributer *Distributor

	// metrics.
	// eventHandleTotal is event handle total stat.
	eventHandleTotal *prometheus.CounterVec

	// eventHandleDuration is event handle cost duration stat.
	eventHandleDuration *prometheus.HistogramVec

	// pusherHandleTotal is event pusher handle total stat.
	pusherHandleTotal *prometheus.CounterVec

	// pusherHandleDuration is event pusher cost duration stat.
	pusherHandleDuration *prometheus.HistogramVec
}

// NewEventHandler creates new EventHandler object.
func NewEventHandler(ctx context.Context, engine *backbone.Engine, cache redis.Client) *EventHandler {
	return &EventHandler{
		ctx:     ctx,
		engine:  engine,
		cache:   cache,
		pushers: make(map[int64]*EventPusher),
	}
}

// registerMetrics registers prometheus metrics.
func (h *EventHandler) registerMetrics() {
	h.eventHandleTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_event_handle_total", etypes.MetricsNamespacePrefix),
			Help: "total number stats of handle event.",
		},
		[]string{"status"},
	)
	h.engine.Metric().Registry().MustRegister(h.eventHandleTotal)

	h.engine.Metric().Registry().MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: fmt.Sprintf("%s_pushers", etypes.MetricsNamespacePrefix),
			Help: "current number of pushers.",
		},
		func() float64 {
			h.pushersMu.RLock()
			defer h.pushersMu.RUnlock()
			return float64(len(h.pushers))
		},
	))

	h.eventHandleDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    fmt.Sprintf("%s_event_handle_duration", etypes.MetricsNamespacePrefix),
			Help:    "handle duration of events.",
			Buckets: []float64{10, 30, 50, 70, 100, 200, 300, 400, 500, 1000, 2000, 5000},
		},
		[]string{"status"},
	)
	h.engine.Metric().Registry().MustRegister(h.eventHandleDuration)

	h.pusherHandleTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_pusher_handle_total", etypes.MetricsNamespacePrefix),
			Help: "total number stats of event pusher.",
		},
		[]string{"status"},
	)
	h.engine.Metric().Registry().MustRegister(h.pusherHandleTotal)

	h.pusherHandleDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    fmt.Sprintf("%s_pusher_handle_duration", etypes.MetricsNamespacePrefix),
			Help:    "pusher duration of events.",
			Buckets: []float64{10, 30, 50, 70, 100, 200, 300, 400, 500, 1000, 2000, 5000},
		},
		[]string{"status"},
	)
	h.engine.Metric().Registry().MustRegister(h.pusherHandleDuration)
}

// SetDistributer setups distributer to event handler.
func (h *EventHandler) SetDistributer(distributer *Distributor) {
	h.distributer = distributer
}

// Handle handles events distributed by distributer, add events to real handle queue to
// handle host identifier infos. Handler would find all relate subscribers and send event message
// to target subscribers by callback.
func (h *EventHandler) Handle(nodes []*watch.ChainNode, eventDetailStrs []string, opts *watch.WatchEventOptions) error {
	blog.V(4).Infof("handle new events now, count[%d]", len(eventDetailStrs))
	defer blog.V(4).Infof("handle new events done, count[%d]", len(eventDetailStrs))

	for idx, eventDetailStr := range eventDetailStrs {
		updateFields := []string{}
		for updateField := range gjson.Get(eventDetailStr, "update_fields").Map() {
			updateFields = append(updateFields, updateField)
		}

		deletedFields := []string{}
		for _, deletedField := range gjson.Get(eventDetailStr, "deleted_fields").Array() {
			deletedFields = append(deletedFields, deletedField.String())
		}

		jsonDetailStr := types.GetEventDetail(&eventDetailStr)
		cut := ccjson.CutJsonDataWithFields(jsonDetailStr, opts.Fields)

		event := &watch.WatchEventDetail{
			Cursor:    nodes[idx].Cursor,
			Resource:  opts.Resource,
			EventType: nodes[idx].EventType,
			Detail:    watch.JsonString(*cut),
		}

		eventInst := &metadata.EventInst{
			ID:            h.cache.Incr(h.ctx, etypes.EventCacheEventIDKey).Val(),
			Cursor:        event.Cursor,
			UpdateFields:  updateFields,
			DeletedFields: deletedFields,
		}

		cursor := &watch.Cursor{}
		if err := cursor.Decode(eventInst.Cursor); err == nil {
			timeNow := metadata.ParseTimeInUnixTS(int64(cursor.ClusterTime.Sec), int64(cursor.ClusterTime.Nano))
			eventInst.ActionTime = timeNow
		}

		switch event.Resource {
		case watch.Host:
			eventInst.EventType = metadata.EventTypeInstData
			eventInst.ObjType = common.BKInnerObjIDHost

		case watch.ModuleHostRelation:
			eventInst.EventType = metadata.EventTypeRelation
			eventInst.ObjType = metadata.EventObjTypeModuleTransfer

		case watch.Biz:
			eventInst.EventType = metadata.EventTypeInstData
			eventInst.ObjType = common.BKInnerObjIDApp

		case watch.Set:
			eventInst.EventType = metadata.EventTypeInstData
			eventInst.ObjType = common.BKInnerObjIDSet

		case watch.Module:
			eventInst.EventType = metadata.EventTypeInstData
			eventInst.ObjType = common.BKInnerObjIDModule

		case watch.Process:
			eventInst.EventType = metadata.EventTypeInstData
			eventInst.ObjType = common.BKInnerObjIDProc

		case watch.ProcessInstanceRelation:
			eventInst.EventType = metadata.EventTypeRelation
			eventInst.ObjType = metadata.EventObjTypeProcModule

		case watch.ObjectBase:
			eventInst.EventType = metadata.EventTypeInstData
			jsonStr := event.Detail.(watch.JsonString)
			detailBytes, _ := jsonStr.MarshalJSON()
			eventInst.ObjType = gjson.Get(string(detailBytes), "bk_obj_id").String()

		default:
			continue
		}

		switch event.EventType {
		case watch.Create:
			eventInst.Action = metadata.EventActionCreate
			eventInst.Data = []metadata.EventData{{CurData: event.Detail}}

		case watch.Update:
			eventInst.Action = metadata.EventActionUpdate
			eventInst.Data = []metadata.EventData{{CurData: event.Detail}}

		case watch.Delete:
			eventInst.Action = metadata.EventActionDelete
			eventInst.Data = []metadata.EventData{{PreData: event.Detail}}

		default:
			continue
		}

		eventData, err := json.Marshal(eventInst)
		if err != nil {
			blog.Errorf("marshal event data failed, %+v, %+v", eventInst, err)
			continue
		}
		blog.Info("handle new event %+v", eventInst)

		// push event data to main event queue.
		cost := time.Now()
		err = h.cache.LPush(h.ctx, etypes.EventCacheEventQueueKey, eventData).Err()
		h.eventHandleDuration.WithLabelValues("PushEventInst").Observe(time.Since(cost).Seconds())

		if err != nil {
			h.eventHandleTotal.WithLabelValues("PushEventInstFailed").Inc()
			blog.Errorf("push new event data to event queue failed, %+v, %+v", eventInst, err)
			continue
		}
	}

	return nil
}

// popEvent keeps poping event from main event queue and add event to duplicated queue. identifier would
// handle host events and add back to main event queue with another level EVENT.
func (h *EventHandler) popEvent() (*metadata.EventInst, error) {
	var eventStr string

	// pop from main event queue and re-add to duplicated queue for identifier.
	cost := time.Now()
	err := h.cache.BRPopLPush(h.ctx, etypes.EventCacheEventQueueKey,
		etypes.EventCacheEventQueueDuplicateKey, defaultTransTimeout).Scan(&eventStr)
	h.eventHandleDuration.WithLabelValues("PopNewEvent").Observe(time.Since(cost).Seconds())

	if err != nil {
		h.eventHandleTotal.WithLabelValues("PopNewEventFailed").Inc()
		return nil, err
	}

	// ignore empty event.
	if len(eventStr) == 0 || eventStr == etypes.NilStr {
		return nil, nil
	}
	eventData := []byte(eventStr)

	// marshal to EventInst.
	eventInst := &metadata.EventInst{}
	if err := json.Unmarshal(eventData, eventInst); err != nil {
		return nil, err
	}

	// check expire event.
	if time.Now().Unix()-eventInst.ActionTime.Unix() > defaultFusingEventExpireSec {
		// old event, expire it.
		h.eventHandleTotal.WithLabelValues("ExpireEventNum").Inc()
		return nil, nil
	}
	blog.V(4).Infof("pop new event, %s", eventStr)

	return eventInst, nil
}

// nextDistID returns new event dist id for target subscriber.
func (h *EventHandler) nextDistID(subid int64) (int64, error) {
	return h.cache.Incr(h.ctx, etypes.EventCacheDistIDPrefix+fmt.Sprint(subid)).Result()
}

// sendToPushers sends new event instance to pusher of target subscriber.
func (h *EventHandler) sendToPusher(subid int64, dist *metadata.DistInst) error {
	h.pushersMu.Lock()

	if _, isExist := h.pushers[subid]; !isExist {
		// create new pusher for the subscriber.
		newPusher := NewEventPusher(h.ctx, h.engine, subid, h.cache, h.distributer,
			h.pusherHandleTotal, h.pusherHandleDuration)

		// run new pusher.
		newPusher.Run()
		h.pushers[subid] = newPusher
	}

	// pusher of subscriber.
	pusher := h.pushers[subid]
	h.pushersMu.Unlock()

	dstbID, err := h.nextDistID(subid)
	if err != nil {
		return err
	}
	dist.DstbID = dstbID
	dist.SubscriptionID = subid

	// add to subscriber pusher.
	return pusher.Handle(dist)
}

// handleEvent handles target event.
func (h *EventHandler) handleEvent(event *metadata.EventInst) error {
	blog.V(4).Infof("handle event inst, %+v", event)
	defer blog.V(4).Infof("handle event inst done, %+v", event)

	// find all subscribers on this event type.
	subscribers := h.distributer.FindSubscribers(event.GetType())

	// distribute to subscribers.
	for _, subscriber := range subscribers {
		// push to subscriber pusher.
		cost := time.Now()
		err := h.sendToPusher(subscriber, &metadata.DistInst{EventInst: *event})
		h.eventHandleDuration.WithLabelValues("PushEventToPusher").Observe(time.Since(cost).Seconds())

		if err != nil {
			blog.Errorf("push to pusher for subid[%d] failed, %+v", subscriber, err)
			h.eventHandleTotal.WithLabelValues("PushEventToPusherFailed").Inc()
			continue
		}
		h.eventHandleTotal.WithLabelValues("Success").Inc()
	}

	return nil
}

func (h *EventHandler) handlePushers() {
	ticker := time.NewTicker(defaultPushersCheckInterval)
	defer ticker.Stop()

	for {
		<-ticker.C
		blog.Info("handle pushers, check base on subscriptions now")

		subids := h.distributer.ListSubscriptionIDs()

		h.pushersMu.Lock()
		for _, subid := range subids {
			if _, isExist := h.pushers[subid]; !isExist {
				// create new pusher for the subscriber.
				newPusher := NewEventPusher(h.ctx, h.engine, subid, h.cache, h.distributer,
					h.pusherHandleTotal, h.pusherHandleDuration)

				// run new pusher.
				newPusher.Run()
				h.pushers[subid] = newPusher
			}
		}
		pushersCount := len(h.pushers)
		h.pushersMu.Unlock()

		blog.Info("handle pushers, current subscription count[%d], pusher count[%d]", len(subids), pushersCount)
	}
}

// cleaning keeps cleaning expire or redundancy events in event cache queue.
func (h *EventHandler) cleaning() {
	ticker := time.NewTicker(defaultCleanCheckInterval)
	defer ticker.Stop()

	for {
		if !h.engine.ServiceManageInterface.IsMaster() {
			blog.Warnf("not master eventserver node, skip cleaning events!")
			time.Sleep(defaultMasterCheckInterval)
			continue
		}

		<-ticker.C
		blog.Info("cleaning expire and redundancy events now")

		// clean.
		eventCount, err := h.cache.LLen(h.ctx, etypes.EventCacheEventQueueKey).Result()
		if err != nil {
			blog.Errorf("fetch expire and redundancy events failed, %+v", err)
			continue
		}
		blog.Info("cleaning events, current events count[%d], trim threshold[%d], delete threshold[%d], clean unit count[%d]",
			eventCount, defaultCleanTrimThreshold, defaultCleanDelThreshold, defaultCleanUnit)

		if eventCount < defaultCleanTrimThreshold {
			continue
		}

		if eventCount < defaultCleanDelThreshold {
			if err := h.cache.LTrim(h.ctx, etypes.EventCacheEventQueueKey, 0, eventCount-defaultCleanUnit).Err(); err != nil {
				blog.Errorf("trim expire and redundancy events failed, %+v", err)
				continue
			}
		}

		// too many events.
		if err := h.cache.Del(h.ctx, etypes.EventCacheEventQueueKey).Err(); err != nil {
			blog.Errorf("delete expire and redundancy queue failed, %+v", err)
			continue
		}
		blog.Info("cleaning expire and redundancy events done")
	}
}

// start starts the poping and handling logic really.
func (h *EventHandler) start() {
	// keep cleaning.
	go h.cleaning()

	// keep check and handle pushers.
	go h.handlePushers()

	for {
		if !h.engine.ServiceManageInterface.IsMaster() {
			blog.Warnf("not master eventserver node, skip distribute events!")
			time.Sleep(defaultMasterCheckInterval)
			continue
		}

		// pop.
		event, err := h.popEvent()
		if err != nil {
			blog.Errorf("pop new event failed, %+v", err)
			time.Sleep(defaultHandleTimeout)
			continue
		}

		// ignore empty event.
		if event == nil {
			continue
		}

		// handle.
		if err := h.handleEvent(event); err != nil {
			blog.Errorf("handle new event failed, %+v, %+v", event, err)
			time.Sleep(defaultHandleTimeout)
			continue
		}
	}
}

// Start starts event handler and keep processing event from distributer.
func (h *EventHandler) Start() error {
	if h.distributer == nil {
		return errors.New("distributer not inited")
	}
	blog.Info("event handler starting now!")

	// register metrics.
	h.registerMetrics()

	// keep poping events and handle distribution.
	go h.start()

	return nil
}
