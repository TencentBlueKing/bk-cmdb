/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package general

import (
	"fmt"
	"time"
)

var (
	// HostKey is the host detail cache key
	HostKey = newGeneralKey(Host, 6*time.Hour, [2]int{0, 30 * 60})
	// ModuleHostRelKey is the host relation cache key
	ModuleHostRelKey = newGeneralKey(ModuleHostRel, 6*time.Hour, [2]int{0, 30 * 60})
	// BizKey is the biz detail cache key
	BizKey = newGeneralKey(Biz, 6*time.Hour, [2]int{0, 30 * 60})
	// SetKey is the set detail cache key
	SetKey = newGeneralKey(Set, 6*time.Hour, [2]int{0, 30 * 60})
	// ModuleKey is the module detail cache key
	ModuleKey = newGeneralKey(Module, 6*time.Hour, [2]int{0, 30 * 60})
	// ProcessKey is the process detail cache key
	ProcessKey = newGeneralKey(Process, 6*time.Hour, [2]int{0, 30 * 60})
	// ProcessRelationKey is the process instance relation cache key
	ProcessRelationKey = newGeneralKey(ProcessRelation, 6*time.Hour, [2]int{0, 30 * 60})
	// BizSetKey is the biz set detail cache key
	BizSetKey = newGeneralKey(BizSet, 6*time.Hour, [2]int{0, 30 * 60})
	// PlatKey is the cloud area detail cache key
	PlatKey = newGeneralKey(Plat, 6*time.Hour, [2]int{0, 30 * 60})
	// ProjectKey is the  detail cache key
	ProjectKey = newGeneralKey(Project, 6*time.Hour, [2]int{0, 30 * 60})
	// ObjInstKey is the object instance detail cache key
	ObjInstKey = NewKey(ObjectInstance, 6*time.Hour, [2]int{0, 30 * 60}, genIDKeyByID, genDetailKeyWithoutSubRes)
	// MainlineInstKey is the mainline instance detail cache key
	MainlineInstKey = NewKey(MainlineInstance, 6*time.Hour, [2]int{0, 30 * 60}, genIDKeyByID, genDetailKeyWithoutSubRes)
	// InstAsstKey is the instance association detail cache key
	InstAsstKey = NewKey(InstAsst, 6*time.Hour, [2]int{0, 30 * 60}, genIDKeyByID, genDetailKeyWithoutSubRes)
)

// newGeneralKey new general Key
func newGeneralKey(resource ResType, expireSeconds time.Duration, expireRangeSeconds [2]int) *Key {
	genIDKey := genIDKeyByID
	genDetailKey := genDetailKeyWithoutSubRes

	_, needOid := ResTypeNeedOidMap[resource]
	if needOid {
		genIDKey = genIDKeyByOid
	}

	_, hasSubRes := ResTypeHasSubResMap[resource]
	if hasSubRes {
		genDetailKey = genDetailKeyWithSubRes
	}

	return NewKey(resource, expireSeconds, expireRangeSeconds, genIDKey, genDetailKey)
}

var cacheKeyMap = map[ResType]*Key{
	Host:             HostKey,
	ModuleHostRel:    ModuleHostRelKey,
	Biz:              BizKey,
	Set:              SetKey,
	Module:           ModuleKey,
	Process:          ProcessKey,
	ProcessRelation:  ProcessRelationKey,
	BizSet:           BizSetKey,
	Plat:             PlatKey,
	Project:          ProjectKey,
	ObjectInstance:   ObjInstKey,
	MainlineInstance: MainlineInstKey,
	InstAsst:         InstAsstKey,
}

// GetCacheKeyByResType get general resource detail cache key by resource type
func GetCacheKeyByResType(res ResType) (*Key, error) {
	key, exists := cacheKeyMap[res]
	if !exists {
		return nil, fmt.Errorf("resource type %s is invalid", res)
	}

	return key, nil
}
