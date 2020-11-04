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

package event

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"configcenter/src/common/json"
	"configcenter/src/common/watch"
	"configcenter/src/storage/driver/redis"
)

const (
	eventStep           = 200
	cursorNotExistError = "cursor not exist error"

	// getNodeWithCursorScript is to get node start from a cursor, the return result
	// do not contain this cursor's value.
	// KEYS[1]: the hashmap's main key
	// KEYS[2]: the start from cursor
	// KEYS[3]: scan count
	// KEYS[4]: the tail cursor which is used to stop the loop to avoid dead lock loop.
	// ARGV[1]: cursor not exist error.
	getNodeWithCursorScript = `
local node = redis.pcall('hget', KEYS[1], KEYS[2]);

if (node == false) then
	return ARGV[1]
end;

local nodeJson = cjson.decode(node);

local elements = {};

local next = nodeJson.next_cursor;
for i = 1,KEYS[3] do
	local ele = redis.pcall('hget', KEYS[1], next);
	if (ele == false) then
		break
	end;

	local js = cjson.decode(ele);
	next = js.next_cursor
	elements[i] = ele;

	if(js.cursor == KEYS[4]) then
		break
	end;
end;

return elements

`
)

// the returned chain node may contain the tail node if can scan to it.
func (f *Flow) getNodesFromCursor(count int, cursor string, key Key) ([]*watch.ChainNode, error) {
	keys := []string{key.MainHashKey(), cursor, strconv.Itoa(count), key.TailKey()}
	return f.runScriptsWithArrayChainNode(getNodeWithCursorScript, keys, cursorNotExistError)
}

// runScripts run lua scripts that returns an string if an error occurs.
// or return a result array ChainNode
func (f *Flow) runScriptsWithArrayChainNode(script string, keys []string, args ...interface{}) ([]*watch.ChainNode, error) {
	result, err := redis.Client().Eval(context.Background(), script, keys, args...).Result()
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, fmt.Errorf("unsupported redis eval result value: %v", result)
	}

	switch reflect.TypeOf(result).Kind() {
	case reflect.String:
		err := result.(string)
		return nil, errors.New(err)

	case reflect.Slice:
		arrays, ok := result.([]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid redis eval response: %v", result)
		}

		nodes := make([]*watch.ChainNode, len(arrays))
		for idx, ele := range arrays {
			element, ok := ele.(string)
			if !ok {
				return nil, fmt.Errorf("invalid chain node details: %v", ele)
			}
			node := new(watch.ChainNode)
			if err := json.Unmarshal([]byte(element), node); err != nil {
				return nil, fmt.Errorf("unmarshal chain node failed, err: %v", err)
			}
			nodes[idx] = node
		}
		return nodes, nil
	default:
		return nil, fmt.Errorf("unsupported redis eval result value: %v", result)
	}
}
