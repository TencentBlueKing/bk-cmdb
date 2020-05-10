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

import "configcenter/src/common"

const watchCacheNamespace = common.BKCacheKeyV3Prefix + "watch:"

var HostKey = Key{
	namespace:  watchCacheNamespace + "host",
	ttlSeconds: 3 * 60 * 60,
}

var ModuleHostRelationKey = Key{
	namespace:  watchCacheNamespace + "host_relation",
	ttlSeconds: 3 * 60 * 60,
}

type Key struct {
	namespace string
	// the valid event's life time.
	// if the event is exist longer than this, it will be deleted.
	// if use's watch start from value is older than time.Now().Unix() - startFrom value,
	// that means use's is watching event that has already deleted, it's not allowed.
	ttlSeconds int64
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
