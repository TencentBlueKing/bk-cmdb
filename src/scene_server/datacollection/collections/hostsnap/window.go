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
	defaultWindowMinutes = 15
	defaultCheckIntervalHours = 1
	defaultAtTime = "1:00"
	defaultAtTimeHour = 1
	defaultAtTimeMin = 0
)

type Window struct {
	// exist time window exists as true, does not exist as false
	exist bool
	// startTime the start time of each timed task
	startTime time.Time
	// windowMinutes duration of each timed task
	windowMinutes  int
	changeTimeLock sync.RWMutex
}

func (w *Window) setStartTime(startTime time.Time) {
	w.changeTimeLock.Lock()
	w.startTime = startTime
	w.changeTimeLock.Unlock()
}

func (w *Window) getStartTime() time.Time {
	w.changeTimeLock.RLock()
	defer w.changeTimeLock.RUnlock()
	return w.startTime
}

func (w *Window) setWindowMinutes(windowMinutes, intervalTime int) {
	w.windowMinutes = windowMinutes
	if windowMinutes <= 0 {
		blog.Errorf("windowMinutes val %d can not be less than or equal to 0, set the default value %d", windowMinutes, defaultWindowMinutes)
		w.windowMinutes = defaultWindowMinutes
	}

	if windowMinutes > intervalTime {
		blog.Errorf("the value of the task cycle %d min is less than the window time value %d min, set the default value %d .", intervalTime, windowMinutes, defaultWindowMinutes)
		w.windowMinutes = defaultWindowMinutes
	}
}

// doFixedTimeTasks do tasks based on a fixed hour each day
func (w *Window) doFixedTimeTasks(atTime string) {
	hour, min := parseTime(atTime)
	now := time.Now()
	targetTime := time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, now.Location())

	if hour > now.Hour() || (now.Hour() == hour && min > now.Minute()) {
		time.Sleep(targetTime.Sub(now))
	} else {
		time.Sleep(targetTime.Add(time.Hour*24).Sub(now))
	}

	w.setStartTime(time.Now())
	w.doTimedTasks(oneDayToHours)
}

// doTimedTasks do tasks according to a certain time interval
func (w *Window) doTimedTasks(intervalTimeIntVal int) {
	timer := time.NewTicker(time.Hour*time.Duration(intervalTimeIntVal))
	for {
		select {
		case <-timer.C:
			w.setStartTime(time.Now())
		}
	}
}

// canPassWindow used to judge whether the request can be passed
func (w *Window) canPassWindow() bool {
	// exist time window requires time judgment
	if w.exist == true {
		now := time.Now()
		stopTime := w.getStartTime().Add(time.Duration(w.windowMinutes) * time.Minute)
		// not within the time window, unable to pass through the window
		if now.Before(w.getStartTime()) || now.After(stopTime) {
			return false
		}
	}

	return true
}

// NewWindow create a time window
func newWindow() *Window {
	w := &Window{}

	// if the parameters are configured, the time window is valid
	if cc.IsExist("datacollection.hostsnap.timeWindow.windowMinutes") {
		windowMinutes, _ := cc.Int("datacollection.hostsnap.timeWindow.windowMinutes")

		if cc.IsExist("datacollection.hostsnap.timeWindow.atTime") {
			atTime, _ := cc.String("datacollection.hostsnap.timeWindow.atTime")
			w.setWindowMinutes(windowMinutes, oneDayToMinutes)

			w.setStartTime(time.Now())
			go w.doFixedTimeTasks(atTime)

			w.exist = true
			return w
		}

		if cc.IsExist("datacollection.hostsnap.timeWindow.checkIntervalHours") {
			checkIntervalHours, _ := cc.Int("datacollection.hostsnap.timeWindow.checkIntervalHours")
			if checkIntervalHours <= 0 {
				blog.Errorf("checkIntervalHours val %d can not be less than or equal to 0, set the default value %d", checkIntervalHours, defaultCheckIntervalHours)
				checkIntervalHours = defaultCheckIntervalHours
			}
			w.setWindowMinutes(windowMinutes, checkIntervalHours*60)

			w.setStartTime(time.Now())
			go w.doTimedTasks(checkIntervalHours)

			w.exist = true
			return w
		}
	}

	return w
}

func parseTime(atTime string) (hour, min int) {
	timeVal := strings.Split(atTime,":")
	if len(timeVal) != 2 {
		blog.Errorf("the format of atTime value %s is wrong, set the default value %s", atTime, defaultAtTime)
		return defaultAtTimeHour, defaultAtTimeMin
	}

	var err error
	hour, err = strconv.Atoi(timeVal[0])
	if err != nil || hour < 0 || hour > 24 {
		blog.Errorf("the format of atTime value %s is wrong, set the default value %s", atTime, defaultAtTime)
		return defaultAtTimeHour, defaultAtTimeMin
	}

	min, err = strconv.Atoi(timeVal[1])
	if err != nil || min < 0 || min > 60 {
		blog.Errorf("the format of atTime value %s is wrong, set the default value %s", atTime, defaultAtTime)
		return defaultAtTimeHour, defaultAtTimeMin
	}
	return hour, min
}
