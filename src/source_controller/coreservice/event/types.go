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

	"configcenter/src/apimachinery/discovery"
	"configcenter/src/common/metadata"
	"configcenter/src/common/watch"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/stream"
	"gopkg.in/redis.v5"
)

// get resource key
func GetResourceKeyWithCursorType(res watch.CursorType) (Key, error) {
	var key Key
	switch res {
	case watch.Host:
		key = HostKey
	case watch.ModuleHostRelation:
		key = ModuleHostRelationKey
	case watch.Biz:
		key = BizKey
	case watch.Set:
		key = SetKey
	case watch.Module:
		key = ModuleKey
	default:
		return key, fmt.Errorf("unsupported cursor type %s", res)
	}
	return key, nil
}

type FlowOptions struct {
	Collection string
	key        Key
	rds        *redis.Client
	watch      stream.Interface
	db         dal.DB
	isMaster   discovery.ServiceManageInterface
}

type hostArchive struct {
	Oid    string              `bson:"oid"`
	Detail metadata.HostMapStr `bson:"detail"`
}
