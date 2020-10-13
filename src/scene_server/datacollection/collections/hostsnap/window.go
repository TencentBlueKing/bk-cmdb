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

package hostsnap

import (
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
)

const (
	// oneDayToMinutes the value of one day converted into minutes. 24 * 60
	oneDayToMinutes = 1440

	// oneDayToHours the value of one day converted into hours.
	oneDayToHours = 24
)

type Window struct {
	// exist time window exists as true, does not exist as false
	exist bool
	// startTime the start time of each timed task
	startTime time.Time
	// durationInMinutes duration of each timed task
	durationInMinutes  int
	changeTimeLock sync.RWMutex
}

// doFixedTimeTasks do tasks based on a fixed hour each day
func (w *Window) doFixedTimeTasks(hour int) {
	now := time.Now()
	targetTime := time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, now.Location())
	if hour > now.Hour() {
		beforeOneDay, _ := time.ParseDuration("-24h")
		w.changeTimeLock.Lock()
		w.startTime = targetTime.Add(beforeOneDay)
		w.changeTimeLock.Unlock()

		time.Sleep(targetTime.Sub(now))

		w.changeTimeLock.Lock()
		w.startTime = targetTime
		w.changeTimeLock.Unlock()
	} else {
		w.changeTimeLock.Lock()
		w.startTime = targetTime
		w.changeTimeLock.Unlock()

		afterOneDay, _ := time.ParseDuration("24h")
		time.Sleep(targetTime.Add(afterOneDay).Sub(now))

		w.changeTimeLock.Lock()
		w.startTime = targetTime.Add(afterOneDay)
		w.changeTimeLock.Unlock()
	}

	w.doTimedTasks(oneDayToHours)
}

// doTimedTasks do tasks according to a certain time interval
func (w *Window) doTimedTasks(intervalTimeIntVal int) {
	timer := time.NewTicker(time.Hour*time.Duration(intervalTimeIntVal))
	for {
		select {
		case <-timer.C:
			w.changeTimeLock.Lock()
			w.startTime = time.Now()
			w.changeTimeLock.Unlock()
		}
	}
}

// NewWindow create a time window
func newWindow() *Window {
	w := &Window{}
	// if the parameters are configured, the time window is valid
	if cc.IsExist("datacollection.hostsnap.timeWindow.hourRule") && cc.IsExist("datacollection.hostsnap.timeWindow.durationInMinutes") {

		hourRule, _ := cc.String("datacollection.hostsnap.timeWindow.hourRule")
		durationInMinutes, _ := cc.Int("datacollection.hostsnap.timeWindow.durationInMinutes")

		w.durationInMinutes = durationInMinutes
		w.exist = true

		splitVal := strings.Split(hourRule, "/")

		// means to be executed regularly at a fixed time every day
		if len(splitVal) == 1 {
			hour, err := strconv.Atoi(hourRule)
			if err != nil {
				blog.Errorf("parse hourRule config %s error: %s", hourRule, err.Error())
				os.Exit(1)
			}

			if hour > 24 || hour < 0 {
				blog.Errorf("the current time value is %d, the value range is 0-24, need to reset the value.", hour)
				os.Exit(1)
			}

			checkDurationInMinutes(durationInMinutes, oneDayToMinutes)

			go w.doFixedTimeTasks(hour)
			return w
		}

		// means that the execution is an interval of cyclic tasks
		if len(splitVal) == 2 {
			intervalTimeIntVal, err := strconv.Atoi(splitVal[1])
			if err != nil {
				blog.Errorf("parse hourRule config %s error: %s", hourRule, err.Error())
				os.Exit(1)
			}

			if intervalTimeIntVal <= 0 {
				blog.Errorf("the current time value is %d, value needs to be greater than 0", intervalTimeIntVal)
				os.Exit(1)
			}

			checkDurationInMinutes(durationInMinutes, intervalTimeIntVal*60)

			w.startTime = time.Now()
			go w.doTimedTasks(intervalTimeIntVal)

			return w

		}

		// means that the rules are wrong and need to exit
		blog.Errorf("the hourRule rule is wrong and needs to be checked.")
		os.Exit(1)
	}

	return w
}

// judge whether the time of each cycle of the timed task is less than the time that the window can pass
func checkDurationInMinutes(durationInMinutes, intervalTime int) {
	if durationInMinutes > intervalTime {
		blog.Errorf("the value of the task cycle %d min is less than the window time value %d min.", intervalTime, durationInMinutes)
		os.Exit(1)
	}
}

// canPassWindow used to judge whether the request can be passed
func (w *Window) canPassWindow() bool {
	w.changeTimeLock.RLock()
	defer w.changeTimeLock.RUnlock()
	// exist time window requires time judgment
	if w.exist == true {
		now := time.Now()
		unit, _ := time.ParseDuration("1m")
		stopTime := w.startTime.Add(time.Duration(w.durationInMinutes)*unit)
		// not within the time window, unable to pass through the window
		if now.Before(w.startTime) || now.After(stopTime) {
			return false
		}
	}

	return true
}
