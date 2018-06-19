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

package util

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"gopkg.in/mgo.v2/bson"
)

// SetQueryOwner returns condition that in default ownerid and request ownerid
func SetQueryOwner(condition interface{}, ownerID string) map[string]interface{} {

	switch cond := condition.(type) {
	case map[string]interface{}:
		// if val, ok := cond[common.BKOwnerIDField]; ok {
		// 	switch val.(type) {
		// 	case string:
		// 		cond[common.BKOwnerIDField] = map[string]interface{}{common.BKDBIN: []string{common.BKDefaultOwnerID, ownerID}}
		// 	case map[string]interface{}:
		// 		for op, v := range val.(map[string]interface{}) {
		// 			if op == common.BKDBIN {
		// 				switch owners := v.(type) {
		// 				case []string:
		// 					val.(map[string]interface{})[op] = append(owners, common.BKDefaultOwnerID)
		// 				}
		// 			}
		// 		}
		// 	}
		// } else {
		cond[common.BKOwnerIDField] = map[string]interface{}{common.BKDBIN: []string{common.BKDefaultOwnerID, ownerID}}
		// }
		return cond
	case nil:
		return map[string]interface{}{
			common.BKOwnerIDField: map[string]interface{}{common.BKDBIN: []string{common.BKDefaultOwnerID, ownerID}},
		}
	default:
		return map[string]interface{}{
			common.BKOwnerIDField: map[string]interface{}{common.BKDBIN: []string{common.BKDefaultOwnerID, ownerID}},
		}
	}
}

// SetModOwner set condition equal owner id, the condition must be a map or struct
func SetModOwner(condition interface{}, ownerID string) map[string]interface{} {
	switch condition.(type) {
	case map[string]interface{}:
		if ownerID == common.BKSuperOwnerID {
			return condition.(map[string]interface{})
		}
		condition.(map[string]interface{})[common.BKOwnerIDField] = ownerID
		return condition.(map[string]interface{})
	case common.KvMap:
		if ownerID == common.BKSuperOwnerID {
			return condition.(map[string]interface{})
		}
		condition.(common.KvMap)[common.BKOwnerIDField] = ownerID
		return condition.(map[string]interface{})
	default:
		out, err := bson.Marshal(condition)
		if err != nil {
			blog.Fatalf("SetModOwner faile condition %#v, error %s", condition, err.Error())
		}
		val := map[string]interface{}{}
		err = bson.Unmarshal(out, val)
		if err != nil {
			blog.Fatalf("SetModOwner faile condition %#v, error %s", condition, err.Error())
		}
		if ownerID == common.BKSuperOwnerID {
			return val
		}
		val[common.BKOwnerIDField] = ownerID
		return val
	}
}
