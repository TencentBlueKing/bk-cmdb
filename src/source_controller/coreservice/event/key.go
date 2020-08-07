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
	"fmt"

	"configcenter/src/common"

	"github.com/tidwall/gjson"
)

const watchCacheNamespace = common.BKCacheKeyV3Prefix + "watch:"

var hostFields = []string{common.BKHostIDField, common.BKHostInnerIPField, common.BKCloudIDField}

var HostKey = Key{
	namespace:  watchCacheNamespace + "host",
	ttlSeconds: 6 * 60 * 60,
	validator: func(doc []byte) error {
		fields := gjson.GetManyBytes(doc, hostFields...)
		for idx := range hostFields {
			if !fields[idx].Exists() {
				return fmt.Errorf("field %s not exist", hostFields[idx])
			}
		}
		return nil
	},
	instName: func(doc []byte) string {
		fields := gjson.GetManyBytes(doc, hostFields...)
		return fields[1].String() + ":" + fields[2].String()
	},
}

var ModuleHostRelationKey = Key{
	namespace:  watchCacheNamespace + "host_relation",
	ttlSeconds: 6 * 60 * 60,
	instName: func(doc []byte) string {
		fields := gjson.GetManyBytes(doc, "bk_module_id", "bk_host_id")
		return fmt.Sprintf("module id: %s, host id: %s", fields[0].String(), fields[1].String())
	},
}

var bizFields = []string{common.BKAppIDField, common.BKAppNameField}
var BizKey = Key{
	namespace:  watchCacheNamespace + common.BKInnerObjIDApp,
	ttlSeconds: 6 * 60 * 60,
	validator: func(doc []byte) error {
		fields := gjson.GetManyBytes(doc, bizFields...)
		for idx := range bizFields {
			if !fields[idx].Exists() {
				return fmt.Errorf("field %s not exist", bizFields[idx])
			}
		}
		return nil
	},
	instName: func(doc []byte) string {
		fields := gjson.GetManyBytes(doc, bizFields...)
		return fields[1].String()
	},
}

var setFields = []string{common.BKSetIDField, common.BKSetNameField}
var SetKey = Key{
	namespace:  watchCacheNamespace + common.BKInnerObjIDSet,
	ttlSeconds: 6 * 60 * 60,
	validator: func(doc []byte) error {
		fields := gjson.GetManyBytes(doc, setFields...)
		for idx := range setFields {
			if !fields[idx].Exists() {
				return fmt.Errorf("field %s not exist", setFields[idx])
			}
		}
		return nil
	},
	instName: func(doc []byte) string {
		fields := gjson.GetManyBytes(doc, setFields...)
		return fields[1].String()
	},
}

var moduleFields = []string{common.BKModuleIDField, common.BKModuleNameField}
var ModuleKey = Key{
	namespace:  watchCacheNamespace + common.BKInnerObjIDModule,
	ttlSeconds: 6 * 60 * 60,
	validator: func(doc []byte) error {
		fields := gjson.GetManyBytes(doc, moduleFields...)
		for idx := range moduleFields {
			if !fields[idx].Exists() {
				return fmt.Errorf("field %s not exist", moduleFields[idx])
			}
		}
		return nil
	},
	instName: func(doc []byte) string {
		fields := gjson.GetManyBytes(doc, moduleFields...)
		return fields[1].String()
	},
}

type Key struct {
	namespace string
	// the valid event's life time.
	// if the event is exist longer than this, it will be deleted.
	// if use's watch start from value is older than time.Now().Unix() - startFrom value,
	// that means use's is watching event that has already deleted, it's not allowed.
	ttlSeconds int64

	// validator validate whether the event data is valid or not.
	// if not, then this event should not be handle, should be dropped.
	validator func(doc []byte) error

	// instance name returns a name which can describe the event's instances
	instName func(doc []byte) string
}

// MainKey is the hashmap key
func (k Key) MainHashKey() string {
	return k.namespace + ":chain"
}

func (k Key) HeadKey() string {
	return "head"
}

func (k Key) TailKey() string {
	return "tail"
}

// Note: do not change the format, it will affect the way in event server to
// get the details with lua scripts.
func (k Key) DetailKey(cursor string) string {
	return k.namespace + ":detail:" + cursor
}

func (k Key) Namespace() string {
	return k.namespace
}

func (k Key) TTLSeconds() int64 {
	return k.ttlSeconds
}

func (k Key) Validate(doc []byte) error {
	if k.validator != nil {
		return k.validator(doc)
	}

	return nil
}

func (k Key) Name(doc []byte) string {
	if k.instName != nil {
		return k.instName(doc)
	}
	return ""
}

func (k Key) LockKey() string {
	return k.namespace + ":lock"
}
