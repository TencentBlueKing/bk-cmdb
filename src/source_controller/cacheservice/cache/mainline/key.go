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

package mainline

import (
	"fmt"
	"time"

	"configcenter/src/common"
)

// key file is used to manage all the business's mainline topology
// cache keys with keyGenerator.

const (
	bizNamespace      = common.BKCacheKeyV3Prefix + "biz"
	detailTTLDuration = 180 * time.Minute
)

// business's cache key generator instance.
var bizKey = keyGenerator{
	namespace:            bizNamespace,
	name:                 bizKeyName,
	detailExpireDuration: detailTTLDuration,
}

// module's cache key generator instance.
var moduleKey = keyGenerator{
	namespace:            bizNamespace,
	name:                 moduleKeyName,
	detailExpireDuration: detailTTLDuration,
}

// set's cache key generator instance.
var setKey = keyGenerator{
	namespace:            bizNamespace,
	name:                 setKeyName,
	detailExpireDuration: detailTTLDuration,
}

type keyName string

const (
	bizKeyName    keyName = common.BKInnerObjIDApp
	setKeyName    keyName = common.BKInnerObjIDSet
	moduleKeyName keyName = common.BKInnerObjIDModule
)

// newCustomKey initialize a custom object's key generator with object.
func newCustomKey(objID string) *keyGenerator {
	return &keyGenerator{
		namespace:            bizNamespace,
		name:                 keyName(objID),
		detailExpireDuration: detailTTLDuration,
	}
}

// keyGenerator is a mainline instance's cache key generator.
type keyGenerator struct {
	namespace            string
	name                 keyName
	detailExpireDuration time.Duration
}

// resumeTokenKey is used to store the event resume token data
func (k keyGenerator) resumeTokenKey() string {
	return fmt.Sprintf("%s:%s:resume_token", k.namespace, k.name)
}

// resumeAtTimeKey is used to store the time where to resume the event stream.
func (k keyGenerator) resumeAtTimeKey() string {
	return fmt.Sprintf("%s:%s:resume_at_time", k.namespace, k.name)
}

// detailKey is to generate the key to store the instance's detail information.
func (k keyGenerator) detailKey(instID int64) string {
	return fmt.Sprintf("%s:%s_detail:%d", k.namespace, k.name, instID)
}
