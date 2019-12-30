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

package y3_6_201912241627

import (
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func sortBkPropertyIndex(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	cond := make(map[string]interface{})
	objects := make([]metadata.Object, 0)
	if err := db.Table(common.BKTableNameObjDes).Find(cond).All(ctx, &objects); err != nil {
		return fmt.Errorf("sortBkPropertyIndex, but find object info failed, err: %v", err)
	}

	for _, obj := range objects {
		cond[common.BKObjIDField] = obj.ObjectID
		attrGroups := make([]metadata.Group, 0)
		if err := db.Table(common.BKTableNamePropertyGroup).Find(cond).All(ctx, &attrGroups); err != nil {
			blog.Errorf("sortBkPropertyIndex, but find object attribute group info failed continue, objID: %s, err: %v", obj.ObjectID, err)
			continue
		}

		attributes := make([]metadata.Attribute, 0)
		if err := db.Table(common.BKTableNameObjAttDes).Find(cond).Sort(common.BKPropertyIndexField).All(ctx, &attributes); err != nil {
			blog.Errorf("sortBkPropertyIndex, but find object attribute info failed continue, objID: %s, err: %v", obj.ObjectID, err)
			continue
		}

		if err := updateSortedPropertyIndex(ctx, db, conf, attrGroups, attributes); err != nil {
			blog.Errorf("sortBkPropertyIndex, updateSortedPropertyIndex failed, objID: %v, err: %v", obj.ObjectID, err)
			return err
		}
	}

	return nil
}

func updateSortedPropertyIndex(ctx context.Context, db dal.RDB, conf *upgrader.Config, attrGroups []metadata.Group, attributes []metadata.Attribute) error {
	for _, group := range attrGroups {
		var index int64
		for _, attr := range attributes {
			if attr.PropertyGroup == "none" || attr.PropertyGroup == "" {
				attr.PropertyGroup = common.BKDefaultField
			}
			if attr.PropertyGroup != group.GroupID {
				continue
			}

			attr.PropertyIndex = index
			filter := make(map[string]interface{})
			filter["id"] = attr.ID
			if err := db.Table(common.BKTableNameObjAttDes).Update(ctx, filter, attr); err != nil {
				blog.Errorf("updateSortedPropertyIndex, update attribute failed continue, objID: %s, propertyID: %v err: %v", attr.ObjectID, attr.PropertyID, err)
				return fmt.Errorf("updateSortedPropertyIndex, update attribute failed continue, objID: %s, propertyID: %v err: %v", attr.ObjectID, attr.PropertyID, err)
			}

			index++
		}
	}

	return nil
}
