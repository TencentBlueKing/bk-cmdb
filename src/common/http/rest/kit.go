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

package rest

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/common/blog"
	"configcenter/src/common/blog/glog"
	"configcenter/src/common/errors"
	"configcenter/src/common/json"
	"configcenter/src/common/util"
)

/*

FILE: kit 座位request 的接入函数

功能：


log： 功能描述

fatal: 致命性错误， 需要报警或者影响任务执行
error: 业务逻辑错误， 用表示处理任务中逻辑错误或者数据错误
warn: 警告信息，用来提示处理任务中非预期的情况的信息，但是该情况可以忽略的


info:  提示信息，默认情况下该信息不输出， 用来标志处理任务中重要节点任务开始，结束或者提示性信息，启动flag中v>=3 输出内容

debug: 调试信息， 默认情况下该信息输出，主要用于线上出现问题的时候获取更多的日志， 启动flag中v>=5 输出内容
trace: 调用洗洗，默认情况下该信息输出，主要用于线上出现问题的时候获取更多的日志， 启动flag中v>10输出内容

*/

const (
	fatalFlag  = "[fatal: cc panic] "
	infoLevel  = 3
	debugLevel = 5
	traceLevel = 10
)

type Kit struct {
	Rid             string
	Header          http.Header
	Ctx             context.Context
	CCError         errors.DefaultCCErrorIf
	User            string
	SupplierAccount string
}

// NewKit 产生一个新的kit， 一般用于在创建新的协程的时候，这个时候会对header 做处理，删除不必要的http header。
func (kit *Kit) NewKit() *Kit {
	newHeader := util.CCHeader(kit.Header)
	newKit := *kit
	newKit.Header = newHeader
	return &newKit
}

// NewHeader 产生一个新的header， 一般用于在创建新的协程的时候，这个时候会对header 做处理，删除不必要的http header。
func (kit *Kit) NewHeader() http.Header {
	return util.CCHeader(kit.Header)
}

func (kit *Kit) LogFatalJSON(format string, args ...interface{}) {
	kit.logJSON(true, format, args)
}

func (kit *Kit) LogFatalf(format string, args ...interface{}) {
	format = fatalFlag + kit.logprefix() + format + kit.logSuffix()
	glog.ErrorfDepthf(1, format, args...)
}

func (kit *Kit) LogFatal(msg string) {
	msg = fatalFlag + kit.logprefix() + msg + kit.logSuffix()
	glog.ErrorDepth(1, msg)
}

func (kit *Kit) LogErrorJSON(format string, args ...interface{}) {
	kit.logJSON(true, format, args)
}

func (kit *Kit) LogErrorf(format string, args ...interface{}) {
	format = kit.logprefix() + format + kit.logSuffix()
	glog.ErrorfDepthf(1, format, args...)
}

func (kit *Kit) LogError(msg string) {
	msg = kit.logprefix() + msg + kit.logSuffix()
	glog.ErrorDepth(1, msg)
}

func (kit *Kit) LogWarnJSON(format string, args ...interface{}) {
	kit.logJSON(true, format, args)
}

func (kit *Kit) LogWarnf(format string, args ...interface{}) {
	format = kit.logprefix() + format + kit.logSuffix()
	glog.InfoDepthf(1, format, args...)
}

func (kit *Kit) LogWarn(msg string) {
	msg = kit.logprefix() + msg + kit.logSuffix()
	glog.InfoDepth(1, msg)
}

func (kit *Kit) LogInfoJSON(format string, args ...interface{}) {
	if blog.GetV() >= infoLevel {
		kit.logJSON(true, format, args)
	}
}

func (kit *Kit) LogInfof(format string, args ...interface{}) {
	if blog.GetV() >= infoLevel {
		format = kit.logprefix() + format + kit.logSuffix()
		glog.InfoDepthf(1, format, args...)
	}
}

func (kit *Kit) LogInfo(msg string) {
	if blog.GetV() >= infoLevel {
		msg = kit.logprefix() + msg + kit.logSuffix()
		glog.InfoDepth(1, msg)
	}
}

func (kit *Kit) LogDebugJSON(format string, args ...interface{}) {
	if blog.GetV() >= debugLevel {
		kit.logJSON(true, format, args)
	}
}

func (kit *Kit) LogDebugf(format string, args ...interface{}) {
	if blog.GetV() >= debugLevel {
		format = kit.logprefix() + format + kit.logSuffix()
		glog.InfoDepthf(1, format, args...)
	}
}

func (kit *Kit) LogDebug(msg string) {
	if blog.GetV() >= debugLevel {
		msg = kit.logprefix() + msg + kit.logSuffix()
		glog.InfoDepth(1, msg)
	}
}

func (kit *Kit) LogTraceJSON(format string, args ...interface{}) {
	if blog.GetV() >= traceLevel {
		kit.logJSON(true, format, args)
	}
}

func (kit *Kit) LogTracef(format string, args ...interface{}) {
	if blog.GetV() >= traceLevel {
		format = kit.logprefix() + format + kit.logSuffix()
		glog.InfoDepthf(1, format, args...)
	}
}

func (kit *Kit) LogTrace(msg string) {
	if blog.GetV() >= traceLevel {
		msg = kit.logprefix() + msg + kit.logSuffix()
		glog.InfoDepth(1, msg)
	}
}

func (kit *Kit) logJSON(isInfo bool, format string, args ...interface{}) {
	params := []interface{}{}
	for _, arg := range args {
		if f, ok := arg.(errorFunc); ok {
			params = append(params, f.Error())
			continue
		}
		if f, ok := arg.(stringFunc); ok {
			params = append(params, f.String())
			continue
		}
		out, err := json.Marshal(arg)
		if err != nil {
			params = append(params, err.Error())
		}
		params = append(params, out)
	}
	format = kit.logprefix() + format + kit.logSuffix()
	if isInfo {
		glog.InfoDepth(2, fmt.Sprintf(format, params...))
	} else {
		glog.ErrorDepth(2, fmt.Sprintf(format, params...))
	}
}

func (kit *Kit) logSuffix() string {
	return " rid: " + kit.Rid
}

func (kit *Kit) logprefix() string {
	return ""
}

type errorFunc interface {
	Error() string
}
type stringFunc interface {
	String() string
}
