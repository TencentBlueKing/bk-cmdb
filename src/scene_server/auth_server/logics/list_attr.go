/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package logics

import (
	"configcenter/src/ac/iam"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/auth_server/types"

	"gopkg.in/redis.v5"
)

// list enumeration attributes of instance type resource
func (lgc *Logics) ListAttr(kit *rest.Kit, resourceType iam.ResourceTypeID) ([]types.AttrResource, error) {
	attrs := make([]types.AttrResource, 0)
	objID := GetInstanceResourceObjID(resourceType)
	if objID == "" {
		return attrs, nil
	}

	// get all attribute cache keys by object ID TODO use database if error occurred using redis
	keys, err := lgc.cache.Keys(common.BKCacheKeyV3Prefix + "attribute:object" + objID + "*").Result()
	if err != nil {
		blog.ErrorJSON("get attribute cache keys for object %s failed, error: %s, rid: %s", objID, err.Error(), kit.Rid)
		return nil, err
	}

	// get all attributes' id and name field from cache TODO use a separate function for all redis reading
	pipeline := lgc.cache.Pipeline()
	for _, key := range keys {
		pipeline.HGetAll(key)
	}
	results, err := pipeline.Exec()
	if err != nil {
		blog.ErrorJSON("get cached attributes for object %s using keys: %v failed, error: %s, rid: %s", objID, keys, err.Error(), kit.Rid)
		return nil, err
	}
	for _, result := range results {
		cmd := result.(*redis.StringStringMapCmd)
		attribute, err := cmd.Result()
		if err != nil {
			blog.ErrorJSON("get cached attribute result failed, error: %s, keys: %v, rid: %s", objID, keys, err.Error(), kit.Rid)
			return nil, err
		}
		// only returns enumeration attributes
		if attribute[common.BKPropertyTypeField] == common.FieldTypeEnum {
			attrs = append(attrs, types.AttrResource{
				ID:          util.GetStrByInterface(attribute[common.BKPropertyIDField]),
				DisplayName: util.GetStrByInterface(attribute[common.BKPropertyNameField]),
			})
		}
	}

	return attrs, nil
}
