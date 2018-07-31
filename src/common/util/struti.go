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
	"time"
)

const (
	chinaMobilePattern = `^1[34578][0-9]{9}$`
	charPattern        = `^[a-zA-Z]*$`
	numCharPattern     = `^[a-zA-Z0-9]*$`
	mailPattern        = `^[a-z0-9A-Z]+([\-_\.][a-z0-9A-Z]+)*@([a-z0-9A-Z]+(-[a-z0-9A-Z]+)*\.)+[a-zA-Z]{2,4}$`
	datePattern        = `^[0-9]{4}[\-]{1}[0-9]{2}[\-]{1}[0-9]{2}$`
	dateTimePattern    = `^[0-9]{4}[\-]{1}[0-9]{2}[\-]{1}[0-9]{2}[\s]{1}[0-9]{2}[\:]{1}[0-9]{2}[\:]{1}[0-9]{2}$`
	//timeZonePattern    = `^[a-zA-Z]+/[a-z\-\_+\-A-Z]+$`
	timeZonePattern = `^[a-zA-Z0-9\-−_\/\+]+$`
)

var (
	chinaMobileRegexp = regexp.MustCompile(chinaMobilePattern)
	charRegexp        = regexp.MustCompile(charPattern)
	numCharRegexp     = regexp.MustCompile(numCharPattern)
	mailRegexp        = regexp.MustCompile(mailPattern)
	dateRegexp        = regexp.MustCompile(datePattern)
	dateTimeRegexp    = regexp.MustCompile(dateTimePattern)
	timeZoneRegexp    = regexp.MustCompile(timeZonePattern)
)

//字符串输入长度
func CheckLen(sInput string, min, max int) bool {
	if len(sInput) >= min && len(sInput) <= max {
		return true
	}
	return false
}

//是否大、小写字母组合
func IsChar(sInput string) bool {
	return charRegexp.MatchString(sInput)
}

//是否字母、数字组合
func IsNumChar(sInput string) bool {
	return numCharRegexp.MatchString(sInput)
}

//是否日期
func IsDate(sInput string) bool {
	return dateRegexp.MatchString(sInput)
}

//是否时间
func IsTime(sInput string) bool {
	return dateTimeRegexp.MatchString(sInput)
}

//是否时区
func IsTimeZone(sInput string) bool {
	return timeZoneRegexp.MatchString(sInput)
}

//str2time
func Str2Time(timeStr string) time.Time {
	fTime, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, time.Local)
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
