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
package y3_9_202106301910

import (
	"context"

	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

// changeHostAttrEditable change host attributes that has the property group of auto to editable
func changeHostAttrEditable(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	filter := map[string]interface{}{
		"bk_property_group": "auto",
		"bk_obj_id":         "host",
		"bk_biz_id":         0,
		"ispre":             true,
	}

	doc := map[string]interface{}{
		"editable": true,
	}

	return db.Table("cc_ObjAttDes").Update(ctx, filter, doc)
}
