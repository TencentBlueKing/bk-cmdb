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
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/common/watch"
	etypes "configcenter/src/scene_server/event_server/types"
	ewatcher "configcenter/src/scene_server/event_server/watcher"
	"configcenter/src/source_controller/cacheservice/event"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/storage/reflector"
	"configcenter/src/storage/stream/types"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// defaultListWatchPageSize is default page size of list watcher.
	defaultListWatchPageSize = 500

	// defaultWatchEventStepSize is default step size of watch event.
	defaultWatchEventStepSize = 100

	// defaultWatchEventLoopInterval is default watch event loop interval.
	defaultWatchEventLoopInterval = 500 * time.Millisecond

	// defaultMasterCheckInterval is default master state check interval.
	defaultMasterCheckInterval = 3 * time.Second
)

// Distributor is event subscription distributer.
type Distributor struct {
	ctx    context.Context
	engine *backbone.Engine

	// db is cc main database.
	db dal.RDB

	// watchDB is cc event watch database.
	watchDB dal.DB

	// cache is cc redis client.
	cache redis.Client

	// subWatcher is subscription watcher.
	subWatcher reflector.Interface

	// subscriptions is local subscriptions records, update by listwatcher, subscriptionid -> subscription.
	subscriptions map[int64]*metadata.Subscription

	// subscriptionsMu is subscriptions mutex.
	subscriptionsMu sync.RWMutex

	// subscribers is event subscribers map, key(event_type) -> subscriberIDs.
	// could find subids by event type in records.
	subscribers map[string][]int64

	// subscribersMu is subscribers mutex.
	subscribersMu sync.RWMutex

	// waitForHandleResourceCursors is channel to wait for handling resource cursors.
	waitForHandleResourceCursors chan struct{}

	// resourceCursors is cursors for resource, CursorType -> oldest cursor.
	resourceCursors map[watch.CursorType]*watch.Cursor

	// resourceCursorsMu is resourceCursors mutex.
	resourceCursorsMu sync.RWMutex

	// eventHandler is event handler that handles all event senders.
	eventHandler *EventHandler

	// metrics.
	// watchAndDistributeTotal is event watch and distribute total stat.
	watchAndDistributeTotal *prometheus.CounterVec

	// watchAndDistributeDuration is watch and distribute cost duration stat.
	watchAndDistributeDuration *prometheus.HistogramVec
}

// NewDistributer creates a new Distributor instance.
func NewDistributer(ctx context.Context, engine *backbone.Engine, db dal.RDB, cache redis.Client, subWatcher reflector.Interface) *Distributor {
	return &Distributor{
		ctx:                          ctx,
		engine:                       engine,
		db:                           db,
		cache:                        cache,
		subWatcher:                   subWatcher,
		subscriptions:                make(map[int64]*metadata.Subscription),
		subscribers:                  make(map[string][]int64),
		resourceCursors:              make(map[watch.CursorType]*watch.Cursor),
		waitForHandleResourceCursors: make(chan struct{}),
	}
}

// registerMetrics registers prometheus metrics.
func (d *Distributor) registerMetrics() {
	d.watchAndDistributeTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_watch_dist_total", etypes.MetricsNamespacePrefix),
			Help: "total number stats of watched and distributed event.",
		},
		[]string{"status"},
	)

	d.engine.Metric().Registry().MustRegister(d.watchAndDistributeTotal)

	d.engine.Metric().Registry().MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: fmt.Sprintf("%s_subscriptions", etypes.MetricsNamespacePrefix),
			Help: "current number of subscriptions.",
		},
		func() float64 {
			d.subscriptionsMu.RLock()
			defer d.subscriptionsMu.RUnlock()
			return float64(len(d.subscriptions))
		},
	))

	d.watchAndDistributeDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    fmt.Sprintf("%s_watch_dist_duration", etypes.MetricsNamespacePrefix),
			Help:    "watch and distribute duration of events.",
			Buckets: []float64{10, 30, 50, 70, 100, 200, 300, 400, 500, 1000, 2000, 5000},
		},
		[]string{"status"},
	)
	d.engine.Metric().Registry().MustRegister(d.watchAndDistributeDuration)
}

// LoadSubscriptions loads all subscriptions in cc.
func (d *Distributor) LoadSubscriptions() error {
	// list and watch subscriptions.
	opts := types.Options{
		EventStruct: &metadata.Subscription{},
		Collection:  common.BKTableNameSubscription,
	}

	// set event handler callback funcs.
	listWatchCap := &reflector.Capable{
		OnChange: reflector.OnChangeEvent{
			OnLister:     d.onUpsertSubscriptions,
			OnListerDone: d.onListSubscriptionsDone,
			OnAdd:        d.onUpsertSubscriptions,
			OnUpdate:     d.onUpsertSubscriptions,
			OnDelete:     d.onDeleteSubscriptions,
		},
	}

	// set list watch options.
	listOpts := &types.ListWatchOptions{
		Options:  opts,
		PageSize: &defaultListWatchPageSize,
	}

	// run to list and keep watching subscriptions.
	return d.subWatcher.ListWatcher(context.Background(), listOpts, listWatchCap)
}

// onUpsertSubscriptions handles event that target subscription inserted or updated.
// It add or update subscription metadata and subscriber in local chains.
func (d *Distributor) onUpsertSubscriptions(e *types.Event) {
	d.subscriptionsMu.Lock()
	defer d.subscriptionsMu.Unlock()

	subscription := e.Document.(*metadata.Subscription)
	blog.Info("on upsert subscription, %+v", subscription)

	oldSubscription, isExist := d.subscriptions[subscription.SubscriptionID]
	if !isExist {
		// insert.
		d.subscriptions[subscription.SubscriptionID] = subscription
		eventTypes := strings.Split(subscription.SubscriptionForm, ",")

		// adds.
		for _, eventType := range eventTypes {
			d.addSubscriber(eventType, subscription.SubscriptionID)
		}
	} else {
		// update.
		if subscription.GetCacheKey() != oldSubscription.GetCacheKey() {
			d.subscriptions[subscription.SubscriptionID] = subscription
		}

		// update event types.
		oldEventTypes := strings.Split(oldSubscription.SubscriptionForm, ",")
		eventTypes := strings.Split(subscription.SubscriptionForm, ",")
		subs, plugs := util.CalSliceDiff(oldEventTypes, eventTypes)

		// removes.
		for _, eventType := range subs {
			d.remSubscriber(eventType, subscription.SubscriptionID)
		}

		// adds.
		for _, eventType := range plugs {
			d.addSubscriber(eventType, subscription.SubscriptionID)
		}
	}
}

// onDeleteSubscriptions handles event that target subscription deleted.
// It delete local subscription metadata and removes subscriber in local chains.
func (d *Distributor) onDeleteSubscriptions(e *types.Event) {
	d.subscriptionsMu.Lock()
	defer d.subscriptionsMu.Unlock()

	subscription := e.Document.(*metadata.Subscription)
	blog.Info("on delete subscription, %+v", subscription)

	if _, isExist := d.subscriptions[subscription.SubscriptionID]; isExist {
		delete(d.subscriptions, subscription.SubscriptionID)
	}

	// removes.
	eventTypes := strings.Split(subscription.SubscriptionForm, ",")

	for _, eventType := range eventTypes {
		d.remSubscriber(eventType, subscription.SubscriptionID)
	}
}

// onListSubscriptionsDone handles event that after LIST-DONE at distributer first setup.
func (d *Distributor) onListSubscriptionsDone() {
	d.subscriptionsMu.RLock()
	defer d.subscriptionsMu.RUnlock()

	blog.Info("distributer listwatch subscriptions at first setup done, count[%d]", len(d.subscriptions))
	d.waitForHandleResourceCursors <- struct{}{}
}

// getResourceCursor parses all subscribers and find oldest watcher cursor for distributer which
// would watchs from the oldest cursor and distribute events to all subscribers by event handler.
func (d *Distributor) getResourceCursor(cursorType watch.CursorType) (*watch.Cursor, error) {
	d.subscribersMu.RLock()
	defer d.subscribersMu.RUnlock()

	for eventType, subids := range d.subscribers {
		if metadata.ParseCursorTypeFromEventType(eventType) != cursorType {
			blog.Info("get resource cursor type from event type[%+v], not %+v, ignore it", eventType, cursorType)
			continue
		}

		// range subscribers on target resource event.
		for _, subid := range subids {
			// get all cursors of subscribers, and try to find oldest cursor for target resource.
			cursorKey := etypes.EventCacheSubscriberCursorKey(eventType, subid)

			// get target subscriber cursor from cache.
			val, err := d.cache.Get(d.ctx, cursorKey).Result()
			if err != nil || len(val) == 0 {
				blog.Warnf("get resource[%+v] cursor failed, query target subscriber cursor[%s], %+v, %+v", cursorType, cursorKey, val, err)
				continue
			}

			// decode target cursor from cache.
			subCursor := &watch.Cursor{}
			if err := subCursor.Decode(val); err != nil {
				blog.Warnf("get resource[%+v] cursor failed, invalid cursor[%s], %+v, %+v", cursorType, cursorKey, val, err)
				continue
			}

			// update local resource oldest cursors.
			d.resourceCursorsMu.Lock()
			if oldestCursor, isExist := d.resourceCursors[cursorType]; !isExist {
				d.resourceCursors[cursorType] = subCursor
			} else {
				// compare to get oldest cursor.
				if subCursor.ClusterTime.Sec < oldestCursor.ClusterTime.Sec {
					d.resourceCursors[cursorType] = subCursor
				}
			}
			d.resourceCursorsMu.Unlock()
		}
	}

	// return final oldest cursor for this cursortype.
	return d.resourceCursors[cursorType], nil
}

func (d *Distributor) watchAndDistribute(cursorType watch.CursorType) error {
	// get inner resource key.
	resourcekey, err := event.GetResourceKeyWithCursorType(cursorType)
	if err != nil {
		return err
	}

	go func() {
		for {
			if !d.engine.ServiceManageInterface.IsMaster() {
				blog.Warnf("not master eventserver node, skip watch and distribute!")
				time.Sleep(defaultMasterCheckInterval)
				continue
			}

			opts := &watch.WatchEventOptions{Resource: cursorType, Cursor: watch.NoEventCursor}

			// try get oldest watch cursor every time
			// this may cause push the duplicated event data when restart eventserver
			cursor, err := d.getResourceCursor(cursorType)
			if err != nil {
				blog.Warnf("watch and distribute for resource[%+v] failed, can't get subscriber cursor, %+v", cursorType, err)
			}

			if cursor != nil {
				cursorStr, err := cursor.Encode()
				if err != nil {
					blog.Warnf("watch and distribute for resource[%+v] failed, can't encode subscriber cursor, %+v", cursorType, err)
				} else {
					opts.Cursor = cursorStr
				}
			}

			// watch resource with cursor.
			if err := d.watchAndDistributeWithCursor(cursorType, resourcekey, opts); err != nil {
				d.watchAndDistributeTotal.WithLabelValues("WatchDistributeFailed").Inc()
				blog.Errorf("watch and distribute for resource[%+v] failed, retry now, %+v", cursorType, err)
				time.Sleep(defaultWatchEventLoopInterval)
			}
		}
	}()

	return nil
}

// watchAndDistributeWithCursor watches with oldest cursor(NoEventCursor when find oldest cursor faild or no-exist) and
// distributes events to event handler which would send event messages to subscribers.
func (d *Distributor) watchAndDistributeWithCursor(cursorType watch.CursorType, key event.Key, opts *watch.WatchEventOptions) error {
	blog.Info("start watching and distribute for resource[%+v] with opts[%+v]", cursorType, opts)
	defer blog.Info("stop watching and distribute for resource[%+v] with opts[%+v]", cursorType, opts)

	// build a resource watcher.
	header := util.BuildHeader(common.GetIdentification(), common.BKDefaultOwnerID)
	watcher := ewatcher.NewWatcher(d.ctx, header, d.cache, d.engine.CoreAPI.CacheService().Cache())

	// start from this cursor.
	startCursor := opts.Cursor

	// can't find target old cursor, try to watch from last event node.
	if startCursor == watch.NoEventCursor {
		node, err := watcher.GetLatestEvent(cursorType)
		if err != nil {
			if err != ewatcher.NoEventsError && err != ewatcher.TailNodeTargetNotExistError {
				return err
			}
			// event chain list is empty, which means no event, just watch and wait from the first node.
			blog.Info("watching for resource[%+v] from lastest event node, but no events now, try to watch from head node", cursorType)
		} else {
			// watch from latest event node.
			startCursor = node.Cursor
			blog.Info("watching for resource[%+v] from cursor[%+v]", cursorType, startCursor)
		}
	}

	// keep watching.
	for {
		// check master state.
		if !d.engine.ServiceManageInterface.IsMaster() {
			blog.Warnf("master state changed, stop watching!")
			return errors.New("master state changed")
		}
		watcher.ResetRequestID()

		// new watching round.
		cost := time.Now()
		blog.V(4).Infof("watching for resource[%+v] from cursor[%+v] now! rid:%s", cursorType, startCursor, watcher.GetRid())
		targetNodes, err := watcher.GetNodesFromCursor(defaultWatchEventStepSize, startCursor, opts.Resource)

		d.watchAndDistributeDuration.WithLabelValues("WatchWithCursor").Observe(time.Since(cost).Seconds())
		if err != nil {
			if err == ewatcher.StartCursorNotExistError {
				// target cursor not exist, re-watch from head event node.
				startCursor = watch.NoEventCursor
				time.Sleep(defaultWatchEventLoopInterval)
				continue
			}
			return err
		}

		if len(targetNodes) == 0 {
			// no events from origin.
			time.Sleep(defaultWatchEventLoopInterval)
			continue
		}
		lastNode := targetNodes[len(targetNodes)-1]
		blog.V(4).Infof("watching for resource[%+v], cursor[%+v], new nodes num[%d], rid:%s", cursorType, startCursor, len(targetNodes), watcher.GetRid())

		// hit target event types.
		hitNodes := watcher.GetHitNodeWithEventType(targetNodes, opts.EventTypes)
		if len(hitNodes) == 0 {
			// range some nodes but no any hits, so reset the cursor to current cursor,
			// and try to get hit nodes in next round.
			startCursor = lastNode.Cursor
			time.Sleep(defaultWatchEventLoopInterval)
			continue
		}
		blog.V(4).Infof("watching for resource[%+v], cursor[%+v], new hit nodes num[%d], rid:%s", cursorType, startCursor, len(hitNodes), watcher.GetRid())

		// get final hit event datas.
		cost = time.Now()
		eventDetailStrs, err := watcher.GetEventDetailsWithCursorNodes(opts.Resource, hitNodes)
		d.watchAndDistributeDuration.WithLabelValues("GetEventDetailsWithCursor").Observe(time.Since(cost).Seconds())

		if err != nil {
			d.watchAndDistributeTotal.WithLabelValues("GetEventDetailsWithCursorFailed").Inc()
			blog.Errorf("watching for resource[%+v], get event details with cursor[%+v] failed, %+v, rid:%s", cursorType, startCursor, err, watcher.GetRid())

			// can't get final target event datas, and try again later.
			time.Sleep(defaultWatchEventLoopInterval)
			continue
		}

		// distribute to subscriber senders.
		cost = time.Now()
		err = d.eventHandler.Handle(hitNodes, eventDetailStrs, opts)
		d.watchAndDistributeDuration.WithLabelValues("HandleEvents").Observe(time.Since(cost).Seconds())

		if err != nil {
			d.watchAndDistributeTotal.WithLabelValues("HandleEventsFailed").Inc()
			blog.Errorf("distribute resource[%+v] %d events to handler failed, %+v, rid:%s", cursorType, len(eventDetailStrs), err, watcher.GetRid())

			// get hit nodes success, but can't distribute to event handler, do not reset cursor,
			// try to re-distribute to handler in next round.
			time.Sleep(defaultWatchEventLoopInterval)
			continue
		}
		d.watchAndDistributeTotal.WithLabelValues("Success").Inc()

		// distribute success, reset cursor and watch in next round.
		blog.Infof("watching and distribute for resource[%+v], handled %d nodes successfully, rid:%s", cursorType, len(eventDetailStrs), watcher.GetRid())
		startCursor = lastNode.Cursor
	}

	return nil
}

// addSubscriber adds new subscriber with target event type.
func (d *Distributor) addSubscriber(eventType string, subid int64) {
	d.subscribersMu.Lock()
	defer d.subscribersMu.Unlock()

	blog.Info("add event subscriber, type[%s] subid[%d]", eventType, subid)

	if _, isExist := d.subscribers[eventType]; !isExist {
		d.subscribers[eventType] = []int64{}
	}
	subscribers := d.subscribers[eventType]

	for _, id := range subscribers {
		if subid == id {
			// already exist.
			return
		}
	}
	d.subscribers[eventType] = append(d.subscribers[eventType], subid)
}

// remSubscriber removes subscriber with target event type.
func (d *Distributor) remSubscriber(eventType string, subid int64) {
	d.subscribersMu.Lock()
	defer d.subscribersMu.Unlock()

	blog.Info("remove event subscriber, type[%s] subid[%d]", eventType, subid)

	if _, isExist := d.subscribers[eventType]; !isExist {
		return
	}
	subscribers := d.subscribers[eventType]

	updated := []int64{}
	for _, id := range subscribers {
		if subid != id {
			updated = append(updated, id)
		}
	}
	d.subscribers[eventType] = updated
}

// FindSubscribers returns all subscribers on event type.
func (d *Distributor) FindSubscribers(eventType string) []int64 {
	d.subscribersMu.RLock()
	defer d.subscribersMu.RUnlock()

	subscribers := d.subscribers[eventType]
	if subscribers == nil {
		return []int64{}
	}
	return subscribers
}

// FindSubscription return target subscription base on subid.
func (d *Distributor) FindSubscription(subid int64) *metadata.Subscription {
	d.subscriptionsMu.RLock()
	defer d.subscriptionsMu.RUnlock()

	return d.subscriptions[subid]
}

// ListSubscriptionIDs return all subscription ids.
func (d *Distributor) ListSubscriptionIDs() []int64 {
	d.subscriptionsMu.RLock()
	defer d.subscriptionsMu.RUnlock()

	subids := []int64{}
	for subid := range d.subscriptions {
		subids = append(subids, subid)
	}
	return subids
}

// Start starts the Distributor, it would load all subscriptions in listwatch mode, and handle runtime
// subscription update messages, push event to subscribers when tatget event happend.
func (d *Distributor) Start(eventHandler *EventHandler) error {
	d.eventHandler = eventHandler

	// register metrics.
	d.registerMetrics()

	// list and keep watching subscriptions.
	if err := d.LoadSubscriptions(); err != nil {
		return fmt.Errorf("load subscriptions at first setups failed, %+v", err)
	}

	// wait for LIST-DONE to handle resource cursors.
	<-d.waitForHandleResourceCursors

	// run event hander.
	if err := d.eventHandler.Start(); err != nil {
		return fmt.Errorf("start event handler failed, %+v", err)
	}

	// range all resource cursors and watch to distribute.
	for _, cursorType := range watch.ListEventCallbackCursorTypes() {
		if err := d.watchAndDistribute(cursorType); err != nil {
			return fmt.Errorf("watch and distribute resource events failed, %+v", err)
		}
	}

	return nil
}
