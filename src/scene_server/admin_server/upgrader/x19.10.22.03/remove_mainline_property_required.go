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

package x19_10_22_03

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func removeMainlinePropertyRequired(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	// find mainline model
	cond := map[string]interface{}{common.AssociationKindIDField: common.AssociationKindMainline}
	assts := make([]metadata.Association, 0)
	if err := db.Table(common.BKTableNameObjAsst).Find(cond).All(ctx, &assts); err != nil {
		blog.ErrorJSON("upgrade find mainline model association error. cond:%s, err:%s", cond, err.Error())
		return err
	}
	objIDs := make([]string, 0)
	for _, asst := range assts {
		objIDs = append(objIDs, asst.ObjectID)
		objIDs = append(objIDs, asst.AsstObjID)
	}
	objIDs = util.StrArrayUnique(objIDs)
	// update user-defined property remove required
	filter := map[string]interface{}{
		common.BKObjIDField: map[string]interface{}{
			"$in": objIDs,
		},
		common.BKIsPre: false,
	}
	doc := map[string]interface{}{
		common.BKIsRequiredField: false,
	}
	if err := db.Table(common.BKTableNameObjAttDes).Update(ctx, filter, doc); err != nil {
		blog.ErrorJSON("upgrade remove mainline model user-defined property required error. filter:%s, doc:%s, err:%s", filter, doc, err.Error())
		return err
	}
	return nil
}
