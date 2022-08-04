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

package hostidentifier

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/util"
	getstatus "configcenter/src/thirdparty/gse/get_agent_state_forsyncdata"

	"github.com/tidwall/gjson"
)

const (
	// maxSecondForSleep is the maximum number of seconds of sleep
	maxSecondForSleep = 60
)

func strMd5(str string) (retMd5 string) {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// sleepForFail sleep due to failure
func sleepForFail(failCount int) {
	if failCount > maxSecondForSleep {
		failCount = maxSecondForSleep
	}
	time.Sleep(time.Duration(failCount) * time.Second)
}

func newHeaderWithRid() (http.Header, string) {
	header := http.Header{}
	header.Add(common.BKHTTPOwnerID, common.BKDefaultOwnerID)
	header.Add(common.BKHTTPHeaderUser, common.CCSystemOperatorUserName)
	rid := util.GenerateRID()
	header.Add(common.BKHTTPCCRequestID, rid)
	return header, rid
}

// buildForStatus 构造查询agent状态的主机信息
func buildForStatus(cloudID, innerIP string) []*getstatus.CacheIPInfo {
	var hostInfos []*getstatus.CacheIPInfo
	// 对于多ip的情况需要特殊处理，agent可能仅有一个ip处于on状态，需要将ip数组里的ip分别查询
	ips := strings.Split(innerIP, ",")
	for _, ip := range ips {
		hostInfos = append(hostInfos, &getstatus.CacheIPInfo{
			GseCompositeID: cloudID,
			IP:             ip,
		})
	}
	return hostInfos
}

// getStatusOnAgentIP 只需要拿到主机的其中一个处于on状态的ip即可
func getStatusOnAgentIP(cloudID, innerIP string, agentStatus map[string]string) (bool, string) {
	ips := strings.Split(innerIP, ",")
	for _, ip := range ips {
		key := HostKey(cloudID, ip)
		if gjson.Get(agentStatus[key], "bk_agent_alive").Int() == agentOnStatus {
			return true, ip
		}
	}
	return false, ""
}

// buildTaskResultMap 根据与gse约定，需要根据content的内容拿到对应的ip和cloudID，但是现在接口还未提供相关内容，这里作兼容，如果拿不到，就从key中截取相关的信息
func buildTaskResultMap(originMap map[string]string) map[string]string {
	taskResultMap := make(map[string]string)
	for key, val := range originMap {
		if gjson.Get(val, "content.dest").Exists() && gjson.Get(val, "content.dest_cloudid").Exists() {
			key = HostKey(gjson.Get(val, "content.dest_cloudid").String(),
				gjson.Get(val, "content.dest").String())
			taskResultMap[key] = val
			continue
		}

		split := strings.Split(key, ":")
		if len(split) < 2 {
			continue
		}
		key = HostKey(split[len(split)-2], split[len(split)-1])
		taskResultMap[key] = val
	}
	return taskResultMap
}

// HostKey return the host key to represent a unique host
func HostKey(cloudID, hostIP string) string {
	return fmt.Sprintf("%s:%s", cloudID, hostIP)
}
