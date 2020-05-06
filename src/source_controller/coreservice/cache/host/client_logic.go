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

package host

import (
	"bytes"
	"fmt"

	"configcenter/src/common/blog"
	"github.com/tidwall/gjson"
)

// cutJsonDataWithFields cut jsonData and only return the "fields" be targeted.
// jsonData can not be nil, and must be a json string
func cutJsonDataWithFields(jsonData *string, fields []string) *string {
	if jsonData == nil {
		empty := ""
		return &empty
	}
	if len(fields) == 0 || *jsonData == "" {
		return jsonData
	}
	elements := gjson.GetMany(*jsonData, fields...)
	last := len(fields) - 1
	jsonBuffer := bytes.Buffer{}
	jsonBuffer.Write([]byte{'{'})
	for idx, field := range fields {
		jsonBuffer.Write([]byte{'"'})
		jsonBuffer.Write([]byte(field))
		jsonBuffer.Write([]byte{'"'})
		jsonBuffer.Write([]byte{':'})
		if elements[idx].Raw == "" {
			jsonBuffer.Write([]byte("null"))
		} else {
			jsonBuffer.Write([]byte(elements[idx].Raw))
		}
		if idx != last {
			jsonBuffer.Write([]byte{','})
		}
	}
	jsonBuffer.Write([]byte{'}'})
	cutOff := jsonBuffer.String()
	return &cutOff
}

func (c *Client) tryRefreshHostDetail(hostID int64, ips string, cloudID int64, detail []byte) {
	hostKey := hostKey.HostDetailKey(hostID)
	if !c.lock.CanRefresh(hostKey) {
		return
	}
	// set refreshing status
	c.lock.SetRefreshing(hostKey)

	// now, we check whether we can refresh the host detail cache
	go func() {
		defer func() {
			c.lock.SetUnRefreshing(hostKey)
		}()

		refreshHostDetailCache(c.rds, hostID, ips, cloudID, detail)
	}()
}

// NOTE: this script is fragile, the key depends on the way that host key
// been generated.
// So, when you change the key pattern, then you need to change this script.
const getHostWithIpScript = `
local host_id = redis.pcall('get', KEYS[1]); 

if (host_id == false) then 
	return ARGV[1] 
end;

local key = string.format("cc:v3:host:detail:%d", host_id)

local detail = redis.pcall('get', key)

if (detail == false) then
	return ARGV[2]
end;

return detail
`

const hostCloudIdRelationNotExitError = "host cloud id relation not exist"
const hostDetailNotExitError = "host detail not exist"

func (c *Client) getHostDetailWithIP(innerIP string, cloudID int64) (*string, error) {
	keys := hostKey.IPCloudIDKey(innerIP, cloudID)
	result, err := c.rds.Eval(getHostWithIpScript, []string{keys}, hostCloudIdRelationNotExitError, hostDetailNotExitError).Result()
	if err != nil {
		return nil, fmt.Errorf("run getHostWithIpScript in redis failed, err: %v", err)
	}

	resp, ok := result.(string)
	if !ok {
		return nil, fmt.Errorf("run getHostWithIpScript in redis, but get invalid result data: %v", result)
	}

	switch resp {
	case hostCloudIdRelationNotExitError:
		// host inner ip and cloud id relation not exist
		blog.V(5).Infof("run getHostWithIpScript in redis, but not find key: %s", keys)
	case hostDetailNotExitError:
		blog.V(5).Infof("run getHostWithIpScript in redis, but not find host detail key pattern: %s", hostKey.HostDetailKey(-1))
		// host detail not exist
	default:
		// we have find the data, return directly.
		return &resp, nil
	}

	// now, we need to refresh the cache.
	hostID, detail, err := getHostDetailsFromMongoWithIP(c.db, innerIP, cloudID)
	if err != nil {
		return nil, fmt.Errorf("get host detail with ip failed, err: %v", err)
	}

	c.tryRefreshHostDetail(hostID, innerIP, cloudID, detail)
	detailStr := string(detail)
	return &detailStr, nil
}
