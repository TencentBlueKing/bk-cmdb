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

// Package blog TODO
package blog

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"reflect"
	"sync"
	"time"

	"configcenter/src/common/blog/glog"
)

// init TODO
// This is temporary until we agree on log dirs and put those into each cmd.
func init() {
	flag.Set("logtostderr", "true")
}

// GlogWriter serves as a bridge between the standard log package and the glog package.
type GlogWriter struct{}

// Write implements the io.Writer interface.
func (writer GlogWriter) Write(data []byte) (n int, err error) {
	glog.InfoDepth(1, string(data))
	return len(data), nil
}

// Output for mgo logger
func (writer GlogWriter) Output(calldepth int, s string) error {
	glog.InfoDepth(calldepth, s)
	return nil
}

// Print TODO
func (writer GlogWriter) Print(v ...interface{}) {
	glog.InfoDepth(1, v...)
}

// Printf TODO
func (writer GlogWriter) Printf(format string, v ...interface{}) {
	glog.InfoDepthf(1, format, v...)
}

// Println TODO
func (writer GlogWriter) Println(v ...interface{}) {
	glog.InfoDepth(1, v...)
}

var once sync.Once

// InitLogs initializes logs the way we want for blog.
func InitLogs() {
	once.Do(func() {
		log.SetOutput(GlogWriter{})
		log.SetFlags(0)
		// The default glog flush interval is 30 seconds, which is frighteningly long.
		go func() {
			d := time.Duration(5 * time.Second)
			tick := time.Tick(d)
			for {
				select {
				case <-tick:
					glog.Flush()
				}
			}
		}()
	})
}

// CloseLogs TODO
func CloseLogs() {
	glog.Flush()
}

var (
	// Info TODO
	Info = glog.Infof
	// Infof TODO
	Infof = glog.Infof
	// InfofDepthf TODO
	InfofDepthf = glog.InfoDepthf

	// Warn TODO
	Warn = glog.Warningf
	// Warnf TODO
	Warnf = glog.Warningf

	// Error TODO
	Error = glog.Errorf
	// Errorf TODO
	Errorf = glog.Errorf
	// ErrorfDepthf TODO
	ErrorfDepthf = glog.ErrorfDepthf

	// Fatal TODO
	Fatal = glog.Fatal
	// Fatalf TODO
	Fatalf = glog.Fatalf

	// V TODO
	V = glog.V
)

// Debug TODO
func Debug(args ...interface{}) {
	if format, ok := (args[0]).(string); ok {
		glog.InfoDepthf(1, format, args[1:]...)
	} else {
		glog.InfoDepth(1, args)
	}
}

// InfoJSON TODO
func InfoJSON(format string, args ...interface{}) {
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

		if arg == nil {
			params = append(params, []byte("null"))
			continue
		}

		kind := reflect.TypeOf(arg).Kind()
		if kind == reflect.Ptr {
			kind = reflect.TypeOf(arg).Elem().Kind()
		}
		if kind == reflect.Struct || kind == reflect.Interface ||
			kind == reflect.Array || kind == reflect.Map || kind == reflect.Slice {
			out, err := json.Marshal(arg)
			if err != nil {
				params = append(params, arg)
			} else {
				params = append(params, out)
			}
			continue
		}

		params = append(params, arg)
	}
	glog.InfoDepthf(1, format, params...)
}

// ErrorJSON TODO
func ErrorJSON(format string, args ...interface{}) {
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
	glog.ErrorDepth(1, fmt.Sprintf(format, params...))
}

// WarnJSON TODO
func WarnJSON(format string, args ...interface{}) {

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

		if arg == nil {
			params = append(params, []byte("null"))
			continue
		}

		kind := reflect.TypeOf(arg).Kind()
		if kind == reflect.Ptr {
			kind = reflect.TypeOf(arg).Elem().Kind()
		}
		if kind == reflect.Struct || kind == reflect.Interface ||
			kind == reflect.Array || kind == reflect.Map || kind == reflect.Slice {
			out, err := json.Marshal(arg)
			if err != nil {
				params = append(params, arg)
			} else {
				params = append(params, out)
			}
			continue
		}

		params = append(params, arg)
	}
	glog.WarningDepth(1, fmt.Sprintf(format, params...))
}

type errorFunc interface {
	Error() string
}
type stringFunc interface {
	String() string
}

// SetV TODO
func SetV(level int32) {
	glog.SetV(glog.Level(level))
}

// GetV TODO
func GetV() int32 {
	return int32(glog.GetV())
}
