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

package x19_09_03_02

import (
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

type Idgen struct {
	ID         string `bson:"_id"`
	SequenceID uint64 `bson:"SequenceID"`
}

func SetModuleServiceTemplateFieldUneditable(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	filter := map[string]interface{}{
		common.BKObjIDField:      common.BKInnerObjIDModule,
		common.BKPropertyIDField: common.BKServiceTemplateIDField,
	}
	doc := map[string]interface{}{
		"editable": false,
	}
	if err := db.Table(common.BKTableNameObjAttDes).Update(ctx, filter, doc); err != nil {
		if db.IsNotFoundError(err) == false {
			return fmt.Errorf("upgrade x19_09_03_02, set module service_template_id field uneditable failed, err: %v", err)
		}
	}
	return nil
}
