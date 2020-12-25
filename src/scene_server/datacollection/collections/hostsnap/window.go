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

/*
 主机快照属性，如cpu,bk_cpu_mhz,bk_disk,bk_mem等数据的处理时间窗口,用于限制在指定周期的前多少分钟可以让请求通过,超过限定时间将不会处理请求。
 有三个参数，atTime,checkIntervalHours,windowMinute

 当不配置windowMinute，窗口不生效。当配置了windowMinute,至少配置atTime或者checkIntervalHours中的一个,否则不生效。
 当atTime和checkIntervalHours都配置时，取atTime这个配置的语义功能

 如果窗口生效,启动的时候,会先跑完windowMinutes。然后再生效

 atTime,设置一天中,几点开启时间窗口,如配置成14:40,表示14:40开启窗口,如果配置格式不正确,默认值为1:00

 checkIntervalHours,规定每隔几个小时窗口开启,单位为小时,如配置成 3,表示每隔3个小时,开启时间窗口,如果配置格式不正确,默认值为 1。
 注：窗口可以通过的时间为整时，即如果配置成1，那么每隔一个小时的整点开启窗口

 windowMinutes,代表开启时间窗口后,多长时间内请求可以通过,单位为分钟。如配置成 60,表示开启窗口时间60分钟内请求可以通过。
 注意：该时间不能大于窗口每次开启的间隔时间，取值范围不能小于等于0，如果配置不正确，默认值为15
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
	// oneHourToMinutes the value of one hour converted into minutes
	oneHourToMinutes = 60
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

	w.doTimedTasks(oneDayToHours)
}

// doTimedTasks do tasks according to a certain time interval
func (w *Window) doTimedTasks(intervalTimeIntVal int) {
	w.setStartTime(time.Now())
	timer := time.NewTicker(time.Hour*time.Duration(intervalTimeIntVal))
	for {
		select {
		case <-timer.C:
			w.setStartTime(time.Now())
		}
	}
}

func (w *Window) doIntervalHoursTasks(checkIntervalHours int) {
	now := time.Now()
	w.setStartTime(now)
	time.Sleep(time.Minute * time.Duration(oneHourToMinutes - now.Minute()))
	w.doTimedTasks(checkIntervalHours)
}

// canPassWindow used to judge whether the request can be passed
func (w *Window) canPassWindow() bool {
	// exist time window requires time judgment
	if w.exist {
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
	if !cc.IsExist("datacollection.hostsnap.timeWindow.windowMinutes") {
		return w
	}
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
		go w.doIntervalHoursTasks(checkIntervalHours)

		w.exist = true
		return w
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
