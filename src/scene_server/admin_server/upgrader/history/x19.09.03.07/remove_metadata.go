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

package x19_09_03_07

import (
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

var tables = []string{
	common.BKTableNameServiceInstance,
	common.BKTableNameServiceTemplate,
	common.BKTableNameServiceCategory,
	common.BKTableNameProcessTemplate,
}

const PageSize = 200

// ,
// common.BKTableNameProcessInstanceRelation,

func RemoveMetadata(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	for _, tableName := range tables {
		if err := removeMetadata(ctx, db, conf, tableName); err != nil {
			blog.Errorf("removeMetadata on table: %s failed, err: %+v", tableName, err)
			return fmt.Errorf("remove metadata from table: %s failed, err: %+v", tableName, err)
		}
	}
	return nil
}

func removeMetadata(ctx context.Context, db dal.RDB, conf *upgrader.Config, tableName string) error {
	type MetaDataItem struct {
		ID       int64             `bson:"id"`
		Metadata metadata.Metadata `bson:"metadata"`
		BizID    int64             `bson:"bk_biz_id"`
	}
	items := make([]MetaDataItem, 0)
	start := uint64(0)
	for {
		if err := db.Table(tableName).Find(map[string]interface{}{}).Start(start).Limit(PageSize).All(ctx, &items); err != nil {
			if db.IsNotFoundError(err) {
				break
			}
			blog.Errorf("read records failed, table: %s, err: %+v", tableName, err)
			return fmt.Errorf("read records failed, err: %+v", err)
		}
		if len(items) == 0 {
			break
		}
		start += PageSize
		for _, item := range items {
			if item.BizID != 0 {
				continue
			}
			bizID, err := item.Metadata.ParseBizID()
			if err != nil {
				blog.Errorf("parse bk_biz_id field failed, table: %s, item, err: %+v", tableName, item, err)
				return fmt.Errorf("parse bk_biz_id field failed, err: %+v", err)
			}
			itemFilter := map[string]interface{}{
				"id": item.ID,
			}
			doc := map[string]interface{}{
				common.BKAppIDField: bizID,
			}
			if err := db.Table(tableName).Update(ctx, itemFilter, doc); err != nil {
				blog.Errorf("set bk_biz_id field failed, table: %s, itemFilter: %+v, doc: %+v, err: %+v", tableName, itemFilter, doc, err)
				return fmt.Errorf("set bk_biz_id field failed, err: %+v", err)
			}
		}
	}

	// replace with bk_biz_id
	if err := db.Table(tableName).DropColumn(ctx, common.MetadataField); err != nil {
		blog.Errorf("drop metadata field failed, table: %s, err: %+v", tableName, err)
		return fmt.Errorf("drop metadata field failed, err: %v", err)
	}
	return nil
}

func RemoveMetadataProcess(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	tableName := common.BKTableNameBaseProcess
	type MetaDataItem struct {
		ProcessID int64             `bson:"bk_process_id"`
		Metadata  metadata.Metadata `bson:"metadata"`
		BizID     int64             `bson:"bk_biz_id"`
	}
	items := make([]MetaDataItem, 0)
	start := uint64(0)
	for {
		if err := db.Table(tableName).Find(map[string]interface{}{}).Start(start).Limit(PageSize).All(ctx, &items); err != nil {
			if db.IsNotFoundError(err) {
				break
			}
			blog.Errorf("read records failed, table: %s, err: %+v", tableName, err)
			return fmt.Errorf("read records failed, err: %+v", err)
		}
		if len(items) == 0 {
			break
		}
		start += PageSize
		for _, item := range items {
			if item.BizID != 0 {
				continue
			}
			bizID, err := item.Metadata.ParseBizID()
			if err != nil {
				blog.Errorf("parse bk_biz_id field failed, table: %s, item, err: %+v", tableName, item, err)
				return fmt.Errorf("parse bk_biz_id field failed, err: %+v", err)
			}
			itemFilter := map[string]interface{}{
				"bk_process_id": item.ProcessID,
			}
			doc := map[string]interface{}{
				common.BKAppIDField: bizID,
			}
			if err := db.Table(tableName).Update(ctx, itemFilter, doc); err != nil {
				blog.Errorf("set bk_biz_id field failed, table: %s, itemFilter: %+v, doc: %+v, err: %+v", tableName, itemFilter, doc, err)
				return fmt.Errorf("set bk_biz_id field failed, err: %+v", err)
			}
		}
	}

	// replace with bk_biz_id
	if err := db.Table(tableName).DropColumn(ctx, common.MetadataField); err != nil {
		blog.Errorf("drop metadata field failed, table: %s, err: %+v", tableName, err)
		return fmt.Errorf("drop metadata field failed, err: %v", err)
	}
	return nil
}

func RemoveMetadataFromRelation(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	tableName := common.BKTableNameProcessInstanceRelation
	type MetaDataItem struct {
		BkProcessID       int64             `bson:"bk_process_id"`
		ServiceInstanceID int64             `bson:"service_instance_id"`
		ProcessTemplateID int64             `bson:"process_template_id"`
		BkHostID          int64             `bson:"bk_host_id"`
		Metadata          metadata.Metadata `bson:"metadata"`
		BizID             int64             `bson:"bk_biz_id"`
	}
	items := make([]MetaDataItem, 0)
	start := uint64(0)
	for {
		if err := db.Table(tableName).Find(map[string]interface{}{}).Start(start).Limit(PageSize).All(ctx, &items); err != nil {
			if db.IsNotFoundError(err) {
				break
			}
			blog.Errorf("read records failed, table: %s, err: %+v", tableName, err)
			return fmt.Errorf("read records failed, err: %+v", err)
		}
		if len(items) == 0 {
			break
		}
		start += PageSize
		for _, item := range items {
			if item.BizID != 0 {
				continue
			}
			bizID, err := item.Metadata.ParseBizID()
			if err != nil {
				blog.Errorf("parse bk_biz_id field failed, table: %s, item, err: %+v", tableName, item, err)
				return fmt.Errorf("parse bk_biz_id field failed, err: %+v", err)
			}
			itemFilter := map[string]interface{}{
				"bk_process_id":       item.BkProcessID,
				"service_instance_id": item.ServiceInstanceID,
				"process_template_id": item.ProcessTemplateID,
				"bk_host_id":          item.BkHostID,
			}
			doc := map[string]interface{}{
				common.BKAppIDField: bizID,
			}
			if err := db.Table(tableName).Update(ctx, itemFilter, doc); err != nil {
				blog.Errorf("set bk_biz_id field failed, table: %s, itemFilter: %+v, doc: %+v, err: %+v", tableName, itemFilter, doc, err)
				return fmt.Errorf("set bk_biz_id field failed, err: %+v", err)
			}
		}
	}

	// replace with bk_biz_id
	if err := db.Table(tableName).DropColumn(ctx, common.MetadataField); err != nil {
		blog.Errorf("drop metadata field failed, table: %s, err: %+v", tableName, err)
		return fmt.Errorf("drop metadata field failed, err: %v", err)
	}
	return nil
}
