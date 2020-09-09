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

package y3_9_202002131522

import (
	"context"
	"fmt"

	"configcenter/src/storage/dal"
)

func removeOldTables(ctx context.Context, db dal.RDB, tableNames []string) error {
	for _, tableName := range tableNames {
		hasTable, err := db.HasTable(ctx, tableName)
		if err != nil {
			return fmt.Errorf("removeOldTables failed, tableName: %s, err: %+v", tableName, err)
		}
		if hasTable == false {
			continue
		}
		if err := db.DropTable(ctx, tableName); err != nil {
			return fmt.Errorf("removeOldTables failed, tableName: %s, err: %+v", tableName, err)
		}
	}

	return nil
}
