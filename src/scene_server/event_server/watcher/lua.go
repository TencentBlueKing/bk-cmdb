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

package watcher

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"configcenter/src/common/json"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/cacheservice/event"
)

const (
	eventStep = 200
	// which means the the start cursor is not exist error, may be a head cursor.
	startCursorNotExistError = "start cursor not exist error"

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
local next = nodeJson.next_cursor

local elements = {};
if(next == KEYS[4]) then
	return elements
end;

for i = 1,KEYS[3] do
	local ele = redis.pcall('hget', KEYS[1], next);
	if (ele == false) then
		break
	end;

	elements[i] = ele;

	local js = cjson.decode(ele);

	next = js.next_cursor
	if(next == KEYS[4]) then
		break
	end;

end;

return elements
`
)

// GetNodesFromCursor get node start from a cursor, the return result
// do not contain this cursor's value.
func (w *Watcher) GetNodesFromCursor(count int, startCursor string, key event.Key) ([]*watch.ChainNode, error) {
	keys := []string{key.MainHashKey(), startCursor, strconv.Itoa(count), key.TailKey()}
	nodes, err := w.runScriptsWithArrayChainNode(getNodeWithCursorScript, keys, startCursorNotExistError)
	if err != nil {

		if strings.Contains(err.Error(), startCursorNotExistError) {
			if startCursor == key.HeadKey() {
				return nil, HeadNodeNotExistError
			}
			return nil, StartCursorNotExistError
		}

		return nil, err
	}

	return nodes, nil
}

const (
	headOrTailNodeNotExistError         = "head or tail node not exist error"
	headOrTailTargetedNodeNotExistError = "head or tail targeted not not exist error"

	// getHeadTailTargetNode get both head and tail node and then get their targeted node value.
	// the head and tail targeted node may be itself, if there is no event.
	// if head or tail node itself not exist, then an headOrTailNodeNotExistError error will return.
	// if head or tail's targeted is not exist, then an headOrTailTargetedNodeNotExistError error will return.
	// if all success, then an array string will return. the first one is head targeted node value. the second one
	// is tail targeted node.
	// KEYS[1]: the resource's main hashmap key.
	// KEYS[2]: head key.
	// KEYS[3]: tail key.
	// ARGV[1]: headOrTailNodeNotExistError
	// ARGV[2]: headOrTailTargetedNodeNotExistError
	getHeadTailTargetNode = `
local node = redis.pcall('hget', KEYS[1], KEYS[2]);

if (node == false) then
	return ARGV[1]
end;

local headJson = cjson.decode(node);
local headTo = redis.pcall('hget', KEYS[1], headJson.next_cursor);

if (headTo == false) then
	return ARGV[2]
end;

local tailNode = redis.pcall('hget', KEYS[1], KEYS[3]);

if (tailNode == false) then
	return ARGV[1]
end;

local tailJson = cjson.decode(tailNode);
local tailTo = redis.pcall('hget', KEYS[1], tailJson.next_cursor);

if (tailTo == false) then
	return ARGV[2]
end;

local rtn = {};
rtn[1] = headTo
rtn[2] = tailTo

return rtn

`
)

func (w *Watcher) GetHeadTailNodeTargetNode(key event.Key) (*watch.ChainNode, *watch.ChainNode, error) {
	keys := []string{key.MainHashKey(), key.HeadKey(), key.TailKey()}
	headTail, err := w.runScriptsWithArrayChainNode(getHeadTailTargetNode, keys, headOrTailNodeNotExistError, headOrTailTargetedNodeNotExistError)
	if err != nil {

		if strings.Contains(err.Error(), headOrTailNodeNotExistError) {
			return nil, nil, TailNodeNotExistError
		}

		return nil, nil, err
	}
	if len(headTail) != 2 {
		return nil, nil, errors.New("invalid head or tail node response from redis")
	}

	return headTail[0], headTail[1], nil
}

// runScripts run lua scripts that returns an string if an error occurs.
// or return a result array ChainNode
func (w *Watcher) runScriptsWithArrayChainNode(script string, keys []string, args ...interface{}) ([]*watch.ChainNode, error) {
	result, err := w.cache.Eval(w.ctx, script, keys, args...).Result()
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

var (
	NoEventsError               = errors.New(noEventWarning)
	HeadNodeNotExistError       = errors.New(headNodeNotExistError)
	TailNodeNotExistError       = errors.New(tailNodeNotExistError)
	TailNodeTargetNotExistError = errors.New(tailNodeTargetNotExistError)
	StartCursorNotExistError    = errors.New(startCursorNotExistError)
)

const (
	headNodeNotExistError       = "head node not exist error"
	tailNodeNotExistError       = "tail node not exist error"
	tailNodeTargetNotExistError = "tail node target detail not exist error"
	noEventWarning              = "no events"

	// KEYS[1]: the resource's main hashmap key
	// KEYS[2]: head cursor
	// KEYS[3]: tail cursor
	// KEYS[4]: event detail key prefix.
	// ARGV[1]: tail node not exist error
	// ARGV[2]: tail node target detail not exist error
	// ARGV[3]: no event warning
	getTailTargetScript = `
local node = redis.pcall('hget', KEYS[1], KEYS[3]);

if (node == false) then
	return ARGV[1]
end;

local nodeJson = cjson.decode(node);

if (nodeJson.next_cursor == KEYS[2]) then
	return ARGV[3]
end;

local lastNode = redis.pcall('hget', KEYS[1], nodeJson.next_cursor);
if (lastNode == false) then
	return ARGV[1]
end;

local nodeDetail = redis.pcall('get', KEYS[4]..nodeJson.next_cursor);
if (nodeDetail == false) then
	return ARGV[2]
end;

local rtn = {};
rtn[1] = lastNode
rtn[2] = nodeDetail

return rtn
`
)

func (w *Watcher) GetLatestEventDetail(key event.Key) (node *watch.ChainNode, detail string, err error) {
	keys := []string{key.MainHashKey(), key.HeadKey(), key.TailKey(), key.DetailKey("")}

	result, err := w.runScriptsWithArrayString(getTailTargetScript, keys, tailNodeNotExistError,
		tailNodeTargetNotExistError, noEventWarning)
	if err != nil {

		if strings.Contains(err.Error(), tailNodeNotExistError) {
			return nil, "", TailNodeNotExistError
		}

		if strings.Contains(err.Error(), noEventWarning) {
			return nil, "", NoEventsError
		}

		if strings.Contains(err.Error(), tailNodeTargetNotExistError) {
			return nil, "", TailNodeTargetNotExistError
		}

		return nil, "", err
	}

	if len(result) != 2 {
		return nil, "", errors.New("invalid tail target script response from redis")
	}

	node = new(watch.ChainNode)
	if err := json.Unmarshal([]byte(result[0]), node); err != nil {
		return nil, "", fmt.Errorf("unmarshal chain node failed, err: %v", err)
	}

	return node, result[1], nil

}

// runScripts run lua scripts that returns an string if an error occurs.
// or return a result array string
func (w *Watcher) runScriptsWithArrayString(script string, keys []string, args ...interface{}) ([]string, error) {
	result, err := w.cache.Eval(w.ctx, script, keys, args...).Result()
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

		details := make([]string, len(arrays))
		for idx, ele := range arrays {
			element, ok := ele.(string)
			if !ok {
				return nil, fmt.Errorf("invalid element type: %v", ele)
			}
			details[idx] = element

		}
		return details, nil
	default:
		return nil, fmt.Errorf("unsupported redis eval result value: %v", result)
	}
}
