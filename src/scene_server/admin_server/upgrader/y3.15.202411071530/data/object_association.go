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

package data

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader/tools"
	"configcenter/src/storage/dal"
)

var associationMap = map[string]string{
	"set":    "biz",
	"module": "set",
	"host":   "module",
}

func addObjAssociationData(kit *rest.Kit, db dal.Dal) error {

	asstData := make([]interface{}, 0)
	for obj, asstObj := range associationMap {
		asst := association{
			AsstKindID:      "bk_mainline",
			ObjectID:        obj,
			AsstObjID:       asstObj,
			AssociationName: obj + "_bk_mainline_" + asstObj,
			Mapping:         metadata.OneToOneMapping,
			OnDelete:        metadata.NoAction,
			IsPre:           &trueVar,
		}
		asstData = append(asstData, asst)
	}

	needField := &tools.InsertOptions{
		UniqueFields: []string{"bk_obj_asst_id"},
		IgnoreKeys:   []string{"id"},
		IDField:      []string{common.BKFieldID},
		AuditDataField: &tools.AuditDataField{
			ResIDField:   "id",
			ResNameField: "bk_obj_asst_id",
		},
		AuditTypeField: &tools.AuditResType{
			AuditType:    metadata.AssociationKindType,
			ResourceType: metadata.MainlineInstanceRes,
		},
	}

	_, err := tools.InsertData(kit, db.Shard(kit.ShardOpts()), common.BKTableNameObjAsst, asstData, needField)
	if err != nil {
		blog.Errorf("insert data for table %s failed, err: %v", common.BKTableNameObjAsst, err)
		return err
	}

	return nil
}

// association defines the association between two objects.
type association struct {
	ID                   int64                              `bson:"id"`
	AssociationName      string                             `bson:"bk_obj_asst_id"`
	AssociationAliasName string                             `bson:"bk_obj_asst_name"`
	ObjectID             string                             `bson:"bk_obj_id"`
	AsstObjID            string                             `bson:"bk_asst_obj_id"`
	AsstKindID           string                             `bson:"bk_asst_id"`
	Mapping              metadata.AssociationMapping        `bson:"mapping"`
	OnDelete             metadata.AssociationOnDeleteAction `bson:"on_delete"`
	IsPre                *bool                              `bson:"ispre"`
}
