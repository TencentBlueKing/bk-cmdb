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

package y3_10_202109181134

import (
	"context"
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
)

type attribute struct {
	BizID             int64       `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id"`
	ID                int64       `field:"id" json:"id" bson:"id"`
	OwnerID           string      `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`
	ObjectID          string      `field:"bk_obj_id" json:"bk_obj_id" bson:"bk_obj_id"`
	PropertyID        string      `field:"bk_property_id" json:"bk_property_id" bson:"bk_property_id"`
	PropertyName      string      `field:"bk_property_name" json:"bk_property_name" bson:"bk_property_name"`
	PropertyGroup     string      `field:"bk_property_group" json:"bk_property_group" bson:"bk_property_group"`
	PropertyGroupName string      `field:"bk_property_group_name,ignoretomap" json:"bk_property_group_name" bson:"-"`
	PropertyIndex     int64       `field:"bk_property_index" json:"bk_property_index" bson:"bk_property_index"`
	Unit              string      `field:"unit" json:"unit" bson:"unit"`
	Placeholder       string      `field:"placeholder" json:"placeholder" bson:"placeholder"`
	IsEditable        bool        `field:"editable" json:"editable" bson:"editable"`
	IsPre             bool        `field:"ispre" json:"ispre" bson:"ispre"`
	IsRequired        bool        `field:"isrequired" json:"isrequired" bson:"isrequired"`
	IsReadOnly        bool        `field:"isreadonly" json:"isreadonly" bson:"isreadonly"`
	IsOnly            bool        `field:"isonly" json:"isonly" bson:"isonly"`
	IsSystem          bool        `field:"bk_issystem" json:"bk_issystem" bson:"bk_issystem"`
	IsAPI             bool        `field:"bk_isapi" json:"bk_isapi" bson:"bk_isapi"`
	PropertyType      string      `field:"bk_property_type" json:"bk_property_type" bson:"bk_property_type"`
	Option            interface{} `field:"option" json:"option" bson:"option"`
	Description       string      `field:"description" json:"description" bson:"description"`
	Creator           string      `field:"creator" json:"creator" bson:"creator"`
	CreateTime        *time.Time  `json:"create_time" bson:"create_time"`
	LastTime          *time.Time  `json:"last_time" bson:"last_time"`
}

// addSvcInstIDAttrInProc copy the service_instance_id from process relation table to process table
func addSvcInstIDAttrInProc(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	svcInstIDFilter := mapstr.MapStr{
		common.BKObjIDField:      common.BKInnerObjIDProc,
		common.BKPropertyIDField: common.BKServiceInstanceIDField,
		common.BkSupplierAccount: conf.OwnerID,
	}

	cnt, err := db.Table(common.BKTableNameObjAttDes).Find(svcInstIDFilter).Count(ctx)
	if err != nil {
		blog.Errorf("count service instance id attribute failed, err: %v", err)
		return err
	}

	if cnt > 0 {
		return nil
	}

	attrID, err := db.NextSequence(ctx, common.BKTableNameObjAttDes)
	if err != nil {
		blog.Errorf("get new attribute id for service instance id failed, err: %v", err)
		return err
	}

	attrIndexFilter := map[string]interface{}{
		common.BKObjIDField: common.BKInnerObjIDProc,
	}
	sort := common.BKPropertyIndexField + ":-1"
	lastAttr := new(attribute)

	if err := db.Table(common.BKTableNameObjAttDes).Find(attrIndexFilter).Sort(sort).One(ctx, lastAttr); err != nil {
		blog.Errorf("get process attribute max property index id failed, err: %v", err)
		return err
	}

	now := time.Now()
	attr := &attribute{
		ID:            int64(attrID),
		OwnerID:       conf.OwnerID,
		ObjectID:      common.BKInnerObjIDProc,
		PropertyID:    common.BKServiceInstanceIDField,
		PropertyName:  "服务实例ID",
		PropertyGroup: common.BKDefaultField,
		PropertyIndex: lastAttr.PropertyIndex + 1,
		IsEditable:    false,
		IsPre:         true,
		IsRequired:    true,
		IsAPI:         true,
		PropertyType:  common.FieldTypeInt,
		Creator:       conf.User,
		CreateTime:    &now,
		LastTime:      &now,
	}

	return db.Table(common.BKTableNameObjAttDes).Insert(ctx, attr)
}

// migrateSvcInstIDToProc copy the service_instance_id from process relation table to process table
func migrateSvcInstIDToProc(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	for {
		// get processes that has no service instance id in one page
		procFilter := mapstr.MapStr{
			common.BKServiceInstanceIDField: mapstr.MapStr{
				common.BKDBExists: false,
			},
		}
		processes := make([]metadata.Process, 0)
		err := db.Table(common.BKTableNameBaseProcess).Find(procFilter).Fields(common.BKProcessIDField).
			Limit(common.BKMaxPageSize).All(ctx, &processes)
		if err != nil {
			blog.Errorf("get process ids that do not have service instance id field failed, err: %v", err)
			return err
		}

		if len(processes) == 0 {
			return nil
		}

		procIDs := make([]int64, len(processes))
		for index, process := range processes {
			procIDs[index] = process.ProcessID
		}

		// get process id to service instance id relations
		relationFilter := map[string]interface{}{
			common.BKProcessIDField: map[string]interface{}{common.BKDBIN: procIDs},
		}
		procRelations := make([]metadata.ProcessInstanceRelation, 0)
		if err := db.Table(common.BKTableNameProcessInstanceRelation).Find(relationFilter).Fields(
			common.BKProcessIDField, common.BKServiceInstanceIDField).All(ctx, &procRelations); err != nil {
			blog.Errorf("get process relations failed, err: %v", err)
			return err
		}

		if len(procRelations) != len(processes) {
			blog.Errorf("process count differs with relation count, process ids: %+v", procIDs)
			return fmt.Errorf("process count differs with relation count")
		}

		// set service instance id to corresponding process
		for _, relation := range procRelations {
			updateFilter := mapstr.MapStr{
				common.BKProcessIDField: relation.ProcessID,
			}
			svcInstIDInfo := mapstr.MapStr{
				common.BKServiceInstanceIDField: relation.ServiceInstanceID,
			}
			if err := db.Table(common.BKTableNameBaseProcess).Update(ctx, updateFilter, svcInstIDInfo); err != nil {
				blog.Errorf("update process(%d) failed, err: %v", relation.ProcessID, err)
				return err
			}
		}
	}
}

// addProcUniqueIndex add process unique indexes
func addProcUniqueIndex(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	processUniqueIndexes := []types.Index{
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "svcInstID_bkProcName",
			Keys: bson.D{
				{common.BKServiceInstanceIDField, 1},
				{common.BKProcessNameField, 1},
			},
			Unique:     true,
			Background: true,
			PartialFilterExpression: map[string]interface{}{
				common.BKServiceInstanceIDField: map[string]string{common.BKDBType: "number"},
				common.BKProcessNameField:       map[string]string{common.BKDBType: "string"},
			},
		},
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "svcInstID_bkFuncName_bkStartParamRegex",
			Keys: bson.D{
				{common.BKServiceInstanceIDField, 1},
				{common.BKFuncName, 1},
				{common.BKStartParamRegex, 1},
			},
			Unique:     true,
			Background: true,
			PartialFilterExpression: map[string]interface{}{
				common.BKServiceInstanceIDField: map[string]string{common.BKDBType: "number"},
				common.BKFuncName:               map[string]string{common.BKDBType: "string"},
				common.BKStartParamRegex:        map[string]string{common.BKDBType: "string"},
			},
		},
	}

	dbIndexes, err := db.Table(common.BKTableNameBaseProcess).Indexes(ctx)
	if err != nil {
		blog.Errorf("get process index failed. err: %v", err)
		return err
	}
	existIndexMap := make(map[string]struct{})
	for _, dbIndex := range dbIndexes {
		existIndexMap[dbIndex.Name] = struct{}{}
	}

	for _, index := range processUniqueIndexes {
		if _, exists := existIndexMap[index.Name]; exists {
			continue
		}
		err = db.Table(common.BKTableNameBaseProcess).CreateIndex(ctx, index)
		if err != nil && !db.IsDuplicatedError(err) {
			blog.Errorf("add process unique index(%#v) failed, err: %+v", index, err)
			return err
		}
	}

	return nil
}
