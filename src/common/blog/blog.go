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
	"flag"
	"log"
	"sync"
	"time"

	"github.com/golang/glog"
)

// This is temporary until we agree on log dirs and put those into each cmd.
func init() {
	flag.Set("logtostderr", "true")
}

// GlogWriter serves as a bridge between the standard log package and the glog package.
type GlogWriter struct{}

// Write implements the io.Writer interface.
func (writer GlogWriter) Write(data []byte) (n int, err error) {
	glog.Info(string(data))
	return len(data), nil
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

func CloseLogs() {
	glog.Flush()
}

var (
	Info  = glog.Infof
	Infof = glog.Infof

	Warn  = glog.Warningf
	Warnf = glog.Warningf

	Error  = glog.Errorf
	Errorf = glog.Errorf

	Fatal  = glog.Fatal
	Fatalf = glog.Fatalf

	V = glog.V
)

func Debug(args ...interface{}) {
	if format, ok := (args[0]).(string); ok {
		glog.InfoDepthf(1, format, args[1:]...)
	} else {
		glog.InfoDepth(1, args)
	}
}

func InfoJSON(format string, args ...interface{}) {
	params := []interface{}{}
	for _, arg := range args {
		out, err := json.Marshal(arg)
		if err != nil {
			params = append(params, err.Error())
		}
		params = append(params, out)
	}
	glog.InfoDepthf(1, format, params...)
}
