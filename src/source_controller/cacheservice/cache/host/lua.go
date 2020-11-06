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
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/driver/redis"
)

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
	result, err := redis.Client().Eval(context.Background(), getHostWithIpScript, []string{keys}, hostCloudIdRelationNotExitError,
		hostDetailNotExitError).Result()

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
	hostID, detail, err := getHostDetailsFromMongoWithIP(innerIP, cloudID)
	if err != nil {
		return nil, fmt.Errorf("get host detail with ip failed, err: %v", err)
	}

	c.tryRefreshHostDetail(hostID, innerIP, cloudID, detail)
	detailStr := string(detail)
	return &detailStr, nil
}

var keyNotExistError = errors.New(notExistError)

const (
	notExistError = "key not exist error"

	// getPagedHostIDListScript is used to check if a zset key is exist or not, if not then return with a error.
	// if yes, then return the total keys of this zset and keys sorted with key's scores.
	// KEYS[1]:  zset key's name
	// KEYS[2]: page started at position.
	// KEYS[3]: page stopped at position.
	// KEYS[4]: host detail key prefix
	// ARGV[1]: zset key not exist error
	getPagedHostIDListScript = `
local exist = redis.pcall('exists', KEYS[1]);

if (exist == 0) then
	return ARGV[1]
end;

local total = redis.pcall('zcard', KEYS[1]);
local keys = redis.pcall('zrange', KEYS[1], KEYS[2], KEYS[3]);

if table.getn(keys) == 0 then
	local elements = {};
	elements[1] = total;
	elements[2] = keys;
	elements[3] = {};
	return elements
end;

local list ={}
for _,key in ipairs(keys) do
        table.insert(list,KEYS[4]..key)
end

local details = redis.pcall('MGET', unpack(list));

local elements = {};
elements[1] = total;
elements[2] = keys;
elements[3] = details;

return elements
`
)

// getPagedHostIDListScript get paged host id list from redis, it returns all the paged host
// id list if zset key is exist.
// Otherwise, return with a not exist error
// Note: the returned host detail string may be empty when the host detail is not exist.
func (c *Client) getPagedHostDetailList(page metadata.BasePage) (int64, []int64, []string, error) {
	if page.Limit <= 0 || page.Start < 0 {
		return 0, nil, nil, errors.New("invalid page parameter")
	}

	// calculate the end score
	page.Limit = page.Start + page.Limit - 1

	keys := []string{hostKey.HostIDListKey(), strconv.Itoa(page.Start), strconv.Itoa(page.Limit),
		hostKey.HostDetailKeyPrefix()}

	result, err := redis.Client().Eval(context.Background(), getPagedHostIDListScript, keys, notExistError).Result()
	if err != nil {
		return 0, nil, nil, err
	}

	switch reflect.TypeOf(result).Kind() {
	case reflect.String:
		err := result.(string)
		if err == notExistError {
			return 0, nil, nil, keyNotExistError
		} else {
			return 0, nil, nil, errors.New(err)
		}

	case reflect.Slice:
		element, ok := result.([]interface{})
		if !ok {
			return 0, nil, nil, fmt.Errorf("invalid redis eval response: %v", result)
		}

		if len(element) != 3 {
			return 0, nil, nil, fmt.Errorf("invalid redis eval response: %v, not two element", result)
		}

		// first element is the total count of host ids
		total, ok := element[0].(int64)
		if !ok {
			return 0, nil, nil, fmt.Errorf("invalid total key: %v in redis lua response", element[0])
		}

		idStr, ok := element[1].([]interface{})
		if !ok {
			return 0, nil, nil, fmt.Errorf("invalid id list: %v in redis lua response", element[1])
		}

		var err error
		idList := make([]int64, len(idStr))
		for idx, str := range idStr {
			id, ok := str.(string)
			if !ok {
				return 0, nil, nil, fmt.Errorf("invalid host detail: %v", idStr)
			}
			idList[idx], err = strconv.ParseInt(id, 10, 64)
			if err != nil {
				return 0, nil, nil, fmt.Errorf("invalid host id: %v in redis lua response", id)
			}
		}

		if element[2] == nil {
			// no detail is found
			return total, idList, make([]string, len(idList)), nil
		}

		detailStr, ok := element[2].([]interface{})
		if !ok {
			return 0, nil, nil, fmt.Errorf("invalid detail key: %v in redis lua response", element[2])
		}

		detailList := make([]string, len(detailStr))
		for idx, idStr := range detailStr {

			if idStr == nil {
				detailList[idx] = ""
				continue
			}

			detail, ok := idStr.(string)
			if !ok {
				return 0, nil, nil, fmt.Errorf("invalid host detail: %v", idStr)
			}
			detailList[idx] = detail
		}

		return total, idList, detailList, nil

	default:
		return 0, nil, nil, fmt.Errorf("unsupported redis eval result value with get paged host id list, response: %v",
			result)
	}

}
