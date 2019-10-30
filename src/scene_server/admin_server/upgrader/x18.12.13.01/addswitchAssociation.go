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

package x18_12_05_01

import (
	"context"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func addswitchAssociation(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	falseVar := false
	switchAsst := metadata.Association{
		OwnerID:         conf.OwnerID,
		AsstKindID:      "connect",
		ObjectID:        "bk_switch",
		AsstObjID:       "host",
		AssociationName: "bk_switch_connect_host",
		Mapping:         metadata.OneToManyMapping,
		OnDelete:        metadata.NoAction,
		IsPre:           &falseVar,
	}

	_, _, err := upgrader.Upsert(ctx, db, common.BKTableNameObjAsst, switchAsst, "id", []string{"bk_obj_id", "bk_asst_obj_id"}, []string{"id"})
	if err != nil {
		return err
	}

	return nil
}

func changeNetDeviceTableName(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	err := db.DropTable(ctx, "cc_Netcollect_Device")
	if err != nil && !strings.Contains(err.Error(), "ns not found") {
		return err
	}
	err = db.DropTable(ctx, "cc_Netcollect_Property")
	if err != nil && !strings.Contains(err.Error(), "ns not found") {
		return err
	}

	tablenames := []string{"cc_NetcollectDevice", "cc_NetcollectProperty"}
	for _, tablename := range tablenames {
		exists, err := db.HasTable(ctx, tablename)
		if err != nil {
			return err
		}
		if !exists {
			if err = db.CreateTable(ctx, tablename); err != nil && !db.IsDuplicatedError(err) {
				return err
			}
		}
	}
	return nil
}
