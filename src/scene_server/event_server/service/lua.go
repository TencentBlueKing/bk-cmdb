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

package service

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"configcenter/src/common/json"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/coreservice/event"
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

local next = nodeJson.next_cursor
for i = 1,KEYS[3] do
	local ele = redis.pcall('hget', KEYS[1], next);
	if (ele == false) then
		break
	end;

	local js = cjson.decode(ele);
	next = js.next_cursor
	if(next == KEYS[4]) then
		break
	end;

	elements[i] = ele;

end;

return elements

`
)

// getNodesFromCursor get node start from a cursor, the return result
// do not contain this cursor's value.
func (s *Service) getNodesFromCursor(count int, cursor string, key event.Key) ([]*watch.ChainNode, error) {
	keys := []string{key.MainHashKey(), cursor, strconv.Itoa(count), key.TailKey()}
	return s.runScriptsWithArrayChainNode(getNodeWithCursorScript, keys, cursorNotExistError)
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

func (s *Service) getHeadTailNodeTargetNode(key event.Key) (*watch.ChainNode, *watch.ChainNode, error) {
	keys := []string{key.MainHashKey(), key.HeadKey(), key.TailKey()}
	headTail, err := s.runScriptsWithArrayChainNode(getHeadTailTargetNode, keys, headOrTailNodeNotExistError, headOrTailTargetedNodeNotExistError)
	if err != nil {
		return nil, nil, err
	}
	if len(headTail) != 2 {
		return nil, nil, errors.New("invalid head or tail node response from redis")
	}

	return headTail[0], headTail[1], nil
}

// runScripts run lua scripts that returns an string if an error occurs.
// or return a result array ChainNode
func (s *Service) runScriptsWithArrayChainNode(script string, keys []string, args ...interface{}) ([]*watch.ChainNode, error) {
	result, err := s.cache.Eval(script, keys, args...).Result()
	if err != nil {
		return nil, err
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

const (
	tailNodeNotExistError       = "tail node not exist error"
	tailNodeTargetNotExistError = "tail node target detail not exist error"

	// KEYS[1]: the resource's main hashmap key
	// KEYS[2]: tail cursor
	// KEYS[3]: event detail key prefix.
	// ARGV[1]: tail node not exist error
	// ARGV[2]: tail node target detail not exist error
	getTailTargetScript = `
local node = redis.pcall('hget', KEYS[1], KEYS[2]);

if (node == false) then
	return ARGV[1]
end;

local nodeJson = cjson.decode(node);

local lastNode = redis.pcall('hget', KEYS[1], nodeJson.next_cursor);
if (lastNode == false) then
	return ARGV[1]
end;

local nodeDetail = redis.pcall('get', KEYS[3]..nodeJson.next_cursor);
if (nodeDetail == false) then
	return ARGV[2]
end;

local rtn = {};
rtn[1] = lastNode
rtn[2] = nodeDetail

return rtn
`
)

func (s *Service) getLatestEventDetail(key event.Key) (node *watch.ChainNode, detail string, err error) {
	keys := []string{key.MainHashKey(), key.TailKey(), key.DetailKey("")}

	result, err := s.runScriptsWithArrayString(getTailTargetScript, keys, tailNodeNotExistError, tailNodeTargetNotExistError)
	if err != nil {
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
func (s *Service) runScriptsWithArrayString(script string, keys []string, args ...interface{}) ([]string, error) {
	result, err := s.cache.Eval(script, keys, args...).Result()
	if err != nil {
		return nil, err
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
