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
	"configcenter/src/common/blog"
	"fmt"
	"testing"
)

var (
	kvMap = map[string]string{"key": "value"}
)

func getKit() Kit {
	return Kit{
		Rid: "cc000rid1212",
	}
}
func TestFatalLog(t *testing.T) {
	kit := getKit()

	kit.LogFatalJSON("map: %s", kvMap)
	kit.LogFatalf("msg: %s", "fatalf")
	kit.LogFatal("msg: fatal")
	fmt.Println("set v=3")

}

func TestErrorLog(t *testing.T) {
	kit := getKit()

	kit.LogErrorJSON("map: %s", kvMap)
	kit.LogErrorf("msg: %s", "error")
	kit.LogError("msg: error")

}

func TestWarnLog(t *testing.T) {
	kit := getKit()

	kit.LogWarnJSON("map: %s", kvMap)
	kit.LogWarnf("msg: %s", "warn")
	kit.LogWarn("msg: warn")

}

func TestInfoLog(t *testing.T) {
	kit := getKit()

	kit.LogInfoJSON("map: %s", kvMap)
	kit.LogInfof("msg: %s", "Info")
	kit.LogInfo("msg: Info")
	blog.SetV(3)
	fmt.Println("set log level= 3")
	kit.LogInfoJSON("map: %s", kvMap)
	kit.LogInfof("msg: %s", "Info")
	kit.LogInfo("msg: Info")
}

func TestDebugLog(t *testing.T) {
	kit := getKit()

	kit.LogDebugJSON("map: %s", kvMap)
	kit.LogDebugf("msg: %s", "Debug")
	kit.LogDebug("msg: Debug")
	blog.SetV(3)
	fmt.Println("set log level= 3")
	kit.LogDebugJSON("map: %s", kvMap)
	kit.LogDebugf("msg: %s", "Debug")
	kit.LogDebug("msg: Debug")
	blog.SetV(5)
	fmt.Println("set log level= 5")
	kit.LogDebugJSON("map: %s", kvMap)
	kit.LogDebugf("msg: %s", "Debug")
	kit.LogDebug("msg: Debug")
}

func TestTraceLog(t *testing.T) {
	kit := getKit()

	kit.LogTraceJSON("map: %s", kvMap)
	kit.LogTracef("msg: %s", "Trace")
	kit.LogTrace("msg: Trace")
	blog.SetV(3)
	fmt.Println("set log level= 3")
	kit.LogTraceJSON("map: %s", kvMap)
	kit.LogTracef("msg: %s", "Trace")
	kit.LogTrace("msg: Trace")
	blog.SetV(5)
	fmt.Println("set log level= 5")
	kit.LogTraceJSON("map: %s", kvMap)
	kit.LogTracef("msg: %s", "Trace")
	kit.LogTrace("msg: Trace")
	blog.SetV(10)
	fmt.Println("set log level= 10")
	kit.LogTraceJSON("map: %s", kvMap)
	kit.LogTracef("msg: %s", "Trace")
	kit.LogTrace("msg: Trace")
}
