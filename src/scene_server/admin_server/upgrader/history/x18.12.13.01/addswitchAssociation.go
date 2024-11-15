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
	switchAsst := Association{
		OwnerID:         conf.TenantID,
		AsstKindID:      "connect",
		ObjectID:        "bk_switch",
		AsstObjID:       "host",
		AssociationName: "bk_switch_connect_host",
		Mapping:         metadata.OneToManyMapping,
		OnDelete:        metadata.NoAction,
		IsPre:           &falseVar,
	}

	_, _, err := upgrader.Upsert(ctx, db, common.BKTableNameObjAsst, switchAsst, "id",
		[]string{"bk_obj_id", "bk_asst_obj_id"}, []string{"id"})
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

// Association defines the association between two objects.
type Association struct {
	ID      int64  `field:"id" json:"id" bson:"id"`
	OwnerID string `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`

	// the unique id belongs to  this association, should be generated with rules as follows:
	// "$ObjectID"_"$AsstID"_"$AsstObjID"
	AssociationName string `field:"bk_obj_asst_id" json:"bk_obj_asst_id" bson:"bk_obj_asst_id"`
	// the alias name of this association, which is a substitute name in the association kind $AsstKindID
	AssociationAliasName string `field:"bk_obj_asst_name" json:"bk_obj_asst_name" bson:"bk_obj_asst_name"`

	// describe which object this association is defined for.
	ObjectID string `field:"bk_obj_id" json:"bk_obj_id" bson:"bk_obj_id"`
	// describe where the Object associate with.
	AsstObjID string `field:"bk_asst_obj_id" json:"bk_asst_obj_id" bson:"bk_asst_obj_id"`
	// the association kind used by this association.
	AsstKindID string `field:"bk_asst_id" json:"bk_asst_id" bson:"bk_asst_id"`

	// defined which kind of association can be used between the source object and destination object.
	Mapping metadata.AssociationMapping `field:"mapping" json:"mapping" bson:"mapping"`
	// describe the action when this association is deleted.
	OnDelete metadata.AssociationOnDeleteAction `field:"on_delete" json:"on_delete" bson:"on_delete"`
	// describe whether this association is a pre-defined association or not,
	// if true, it means this association is used by cmdb itself.
	IsPre *bool `field:"ispre" json:"ispre" bson:"ispre"`
}
