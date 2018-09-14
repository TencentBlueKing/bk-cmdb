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

package blog

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang/glog"

	"configcenter/src/common"
)

type CCLog struct {
	header http.Header
}

// CCLogHeader get cclog with header
func CCLogHeader(header http.Header) *CCLog {
	return &CCLog{
		header: header,
	}
}

// Errorf  write log with cc info
func (l *CCLog) Errorf(format string, args ...interface{}) {
	l.errorf(format, args...)
}

// Error  write log with cc info
func (l *CCLog) Error(format string, args ...interface{}) {
	l.errorf(format, args...)
}

func (l *CCLog) errorf(format string, args ...interface{}) {
	logStr := fmt.Sprintf(format, args...)
	suffixStr := l.logSuffixStr()
	if "" != suffixStr {
		logStr = fmt.Sprintf("%s %s", logStr, suffixStr)
	}
	glog.ErrorDepth(2, logStr)
}

// Warnf  write log with cc info
func (l *CCLog) Warnf(format string, args ...interface{}) {
	l.warnf(format, args...)
}

// Warn  write log with cc info
func (l *CCLog) Warn(format string, args ...interface{}) {
	l.warnf(format, args...)
}

// warnf  write log with cc info
func (l *CCLog) warnf(format string, args ...interface{}) {
	logStr := fmt.Sprintf(format, args...)
	suffixStr := l.logSuffixStr()
	if "" != suffixStr {
		logStr = fmt.Sprintf("%s %s", logStr, suffixStr)
	}
	glog.WarningDepth(2, logStr)
}

// Debug  write log with cc info
func (l *CCLog) Debug(format string, args ...interface{}) {
	logStr := fmt.Sprintf(format, args...)
	suffixStr := l.logSuffixStr()
	if "" != suffixStr {
		logStr = fmt.Sprintf("%s %s", logStr, suffixStr)
	}
	glog.InfoDepthf(1, logStr)
}

func (l *CCLog) Info(format string, args ...interface{}) {
	l.infof(format, args...)
}

func (l *CCLog) Infof(format string, args ...interface{}) {
	l.infof(format, args...)
}

func (l *CCLog) infof(format string, args ...interface{}) {
	logStr := fmt.Sprintf(format, args...)
	suffixStr := l.logSuffixStr()
	if "" != suffixStr {
		logStr = fmt.Sprintf("%s %s", logStr, suffixStr)
	}
	glog.InfoDepthf(1, logStr)
}

// CCLog write log with cc info
func (l *CCLog) InfoJSON(format string, args ...interface{}) {
	params := []interface{}{}
	for _, arg := range args {
		out, err := json.Marshal(arg)
		if err != nil {
			params = append(params, err.Error())
		}
		params = append(params, out)
	}
	logStr := fmt.Sprintf(format, params...)
	suffixStr := l.logSuffixStr()
	if "" != suffixStr {
		logStr = fmt.Sprintf("%s %s", logStr, suffixStr)
	}
	glog.InfoDepthf(1, logStr)
}

func (l *CCLog) logSuffixStr() string {
	if nil != l.header {
		return "logID:" + l.getHTTPCCRequestID()
	}
	return ""
}

func (l *CCLog) getHTTPCCRequestID() string {
	return l.header.Get(common.BKHTTPCCRequestID)
}
