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

package util

import (
	"regexp"
	"strings"
	"time"
)

const (
	// chinaMobilePattern = `^1[34578][0-9]{9}$`
	charPattern    = `^[a-zA-Z]*$`
	numCharPattern = `^[a-zA-Z0-9]*$`
	// mailPattern     = `^[a-z0-9A-Z]+([\-_\.][a-z0-9A-Z]+)*@([a-z0-9A-Z]+(-[a-z0-9A-Z]+)*\.)+[a-zA-Z]{2,4}$`
	datePattern             = `^[0-9]{4}[\-]{1}[0-9]{2}[\-]{1}[0-9]{2}$`
	dateTimePattern         = `^[0-9]{4}[\-]{1}[0-9]{2}[\-]{1}[0-9]{2}[\s]{1}[0-9]{2}[\:]{1}[0-9]{2}[\:]{1}[0-9]{2}$`
	timeWithLocationPattern = `^[0-9]{4}[\-]{1}[0-9]{2}[\-]{1}[0-9]{2}[T]{1}[0-9]{2}[\:]{1}[0-9]{2}[\:]{1}[0-9]{2}([\.]{1}[0-9]+)?[\+]{1}[0-9]{2}[\:]{1}[0-9]{2}$`
	// timeZonePattern    = `^[a-zA-Z]+/[a-z\-\_+\-A-Z]+$`
	timeZonePattern = `^[a-zA-Z0-9\-−_\/\+]+$`
	// userPattern the user names regex expression
	userPattern = `^(\d|[a-zA-Z])([a-zA-Z0-9\@.,_-])*$`
)

var (
	// chinaMobileRegexp = regexp.MustCompile(chinaMobilePattern)
	charRegexp    = regexp.MustCompile(charPattern)
	numCharRegexp = regexp.MustCompile(numCharPattern)
	// mailRegexp        = regexp.MustCompile(mailPattern)
	dateRegexp             = regexp.MustCompile(datePattern)
	dateTimeRegexp         = regexp.MustCompile(dateTimePattern)
	timeWithLocationRegexp = regexp.MustCompile(timeWithLocationPattern)
	timeZoneRegexp         = regexp.MustCompile(timeZonePattern)
	userRegexp             = regexp.MustCompile(userPattern)
)

// CheckLen 字符串输入长度
func CheckLen(sInput string, min, max int) bool {
	if len(sInput) >= min && len(sInput) <= max {
		return true
	}
	return false
}

// IsChar 是否大、小写字母组合
func IsChar(sInput string) bool {
	return charRegexp.MatchString(sInput)
}

// IsNumChar 是否字母、数字组合
func IsNumChar(sInput string) bool {
	return numCharRegexp.MatchString(sInput)
}

// IsDate 是否日期
func IsDate(sInput string) bool {
	return dateRegexp.MatchString(sInput)
}

// DateTimeFieldType TODO
type DateTimeFieldType string

const (
	// timeWithoutLocationType the common date time type which is used by front end and api
	timeWithoutLocationType DateTimeFieldType = "time_without_location"
	// timeWithLocationType the date time type compatible for values from db which is marshaled with time zone
	timeWithLocationType DateTimeFieldType = "time_with_location"
	invalidDateTimeType  DateTimeFieldType = "invalid"
)

// IsTime 是否时间
func IsTime(sInput string) (DateTimeFieldType, bool) {
	if dateTimeRegexp.MatchString(sInput) {
		return timeWithoutLocationType, true
	}
	if timeWithLocationRegexp.MatchString(sInput) {
		return timeWithLocationType, true
	}
	return invalidDateTimeType, false
}

// IsTimeZone 是否时区
func IsTimeZone(sInput string) bool {
	return timeZoneRegexp.MatchString(sInput)
}

// IsUser 是否用户
func IsUser(sInput string) bool {
	return userRegexp.MatchString(sInput)
}

// Str2Time TODO
// str2time
func Str2Time(timeStr string, timeType DateTimeFieldType) time.Time {
	var layout string
	switch timeType {
	case timeWithoutLocationType:
		layout = "2006-01-02 15:04:05"
	case timeWithLocationType:
		layout = "2006-01-02T15:04:05+08:00"
	default:
		return time.Time{}
	}

	fTime, err := time.ParseInLocation(layout, timeStr, time.Local)
	if nil != err {
		return fTime
	}
	return fTime.UTC()
}

// FirstNotEmptyString return the first string in slice strs that is not empty
func FirstNotEmptyString(strs ...string) string {
	for _, str := range strs {
		if str != "" {
			return str
		}
	}
	return ""
}

// ContainsAnyString TODO
func ContainsAnyString(s string, subs ...string) bool {
	for index := range subs {
		if strings.Contains(s, subs[index]) {
			return true
		}
	}
	return false
}

// Normalize to trim space of the str and get it's upper format
// for example, Normalize(" hello world") ==> "HELLO WORLD"
func Normalize(str string) string {
	return strings.ToUpper(strings.TrimSpace(str))
}
