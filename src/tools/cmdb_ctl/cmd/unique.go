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

package cmd

import (
	"context"
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/index"
	"configcenter/src/common/json"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core/model"
	"configcenter/src/tools/cmdb_ctl/app/config"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(NewUniqueCheckCommand())
}

func NewUniqueCheckCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "unique",
		Short: "check if object instances satisfies unique constraints",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUniqueCheck()
		},
	}
}

type uniqueCheckService struct {
	service *config.Service
}

func newUniqueCheckService(mongoURI string, mongoRsName string) (*uniqueCheckService, error) {
	service, err := config.NewMongoService(mongoURI, mongoRsName)
	if err != nil {
		return nil, err
	}
	return &uniqueCheckService{
		service: service,
	}, nil
}

func runUniqueCheck() error {
	srv, err := newUniqueCheckService(config.Conf.MongoURI, config.Conf.MongoRsName)
	if err != nil {
		return err
	}
	return srv.checkUnique()
}

func (s *uniqueCheckService) checkUnique() error {
	fmt.Println("=================================\nstart checking unique constraints")

	ctx := util.SetDBReadPreference(context.Background(), common.SecondaryPreferredMode)

	objIDs, err := s.getAllObjectIDs(ctx)
	if err != nil {
		return err
	}

	for _, objID := range objIDs {
		if !model.SatisfyMongoCollLimit(objID) {
			fmt.Printf("object id %s is invalid\n", objID)
		}

		attrMap, err := s.getObjAttrMap(ctx, objID)
		if err != nil {
			return err
		}

		fmt.Printf("=================================\nstart checking unique constraints for object %s\n", objID)
		uniques, err := s.getObjectUniques(ctx, objID)
		if err != nil {
			return err
		}

		if len(uniques) == 0 {
			continue
		}

		uniqueKeyMap := make(map[uint64]uint64)
		for _, unique := range uniques {
			isValid := true
			for _, key := range unique.Keys {
				id, exists := uniqueKeyMap[key.ID]
				if exists {
					isValid = false
					fmt.Printf("object(%s) unique(%d) key(%d) duplicate with %d\n", objID, unique.ID, key.ID, id)
					continue
				}
				uniqueKeyMap[key.ID] = unique.ID
			}

			if !isValid {
				fmt.Printf("object(%s) unique(%d) has invalid key, **skip checking instance**\n", objID, unique.ID)
				continue
			}

			if !unique.MustCheck {
				fmt.Printf("WARNING: not must check object(%s) unique(%d) will not be supported\n", objID, unique.ID)
			}

			if err := s.checkObjectUnique(ctx, objID, unique.OwnerID, unique, attrMap); err != nil {
				return err
			}
		}
	}

	fmt.Println("checking unique constraints done")
	return nil
}

func (s *uniqueCheckService) getAllObjectIDs(ctx context.Context) ([]string, error) {
	fmt.Println("start searching for all object ids")

	rawObjIDs, err := s.service.DbProxy.Table(common.BKTableNameObjDes).Distinct(ctx, common.BKObjIDField, nil)
	if err != nil {
		return nil, fmt.Errorf("get all object ids from db failed, err: %v", err)
	}

	objIDs, err := util.SliceInterfaceToString(rawObjIDs)
	if err != nil {
		return nil, fmt.Errorf("parse all object ids failed, err: %v", err)
	}

	return objIDs, nil
}

func (s *uniqueCheckService) getObjectUniques(ctx context.Context, objID string) ([]ObjectUnique, error) {
	fmt.Printf("start searching unique constraints for object %s\n", objID)

	filter := map[string]interface{}{
		common.BKObjIDField: objID,
	}

	uniques := make([]ObjectUnique, 0)
	if err := s.service.DbProxy.Table(common.BKTableNameObjUnique).Find(filter).All(ctx, &uniques); err != nil {
		return nil, fmt.Errorf("get unique rules for object(%s) failed, err: %v", objID, err)
	}

	return uniques, nil
}

func (s *uniqueCheckService) getObjAttrMap(ctx context.Context, objID string) (map[int64]metadata.Attribute, error) {
	fmt.Printf("start searching object attributes for object %s\n", objID)

	filter := map[string]interface{}{
		common.BKObjIDField: objID,
	}

	attributes := make([]metadata.Attribute, 0)
	err := s.service.DbProxy.Table(common.BKTableNameObjAttDes).Find(filter).Fields(common.BKFieldID,
		common.BKPropertyIDField, common.BKPropertyTypeField).All(ctx, &attributes)
	if err != nil {
		return nil, fmt.Errorf("get attributes for object(%s) failed, err: %v", objID, err)
	}

	attrMap := make(map[int64]metadata.Attribute)
	for _, attribute := range attributes {
		if !model.SatisfyMongoFieldLimit(attribute.PropertyID) {
			fmt.Printf("object(%s) attribute(%d) property id %s is invalid\n", objID, attribute.ID, attribute.PropertyID)
		}

		if objID == common.BKInnerObjIDHost && attribute.PropertyID == common.BKCloudIDField {
			// NOTICE: 2021年03月12日 特殊逻辑。 现在主机的字段中类型未foreignkey 特殊的类型
			attribute.PropertyType = common.FieldTypeInt
		}
		if objID == common.BKInnerObjIDHost &&
			(attribute.PropertyID == common.BKHostInnerIPField || attribute.PropertyID == common.BKHostOuterIPField ||
				attribute.PropertyID == common.BKOperatorField || attribute.PropertyID == common.BKBakOperatorField) {
			// NOTICE: 2021年03月12日 特殊逻辑。 现在主机的字段中类型未innerIP,OuterIP 特殊的类型
			attribute.PropertyType = common.FieldTypeList
		}

		attrMap[attribute.ID] = attribute
	}

	return attrMap, nil
}

func (s *uniqueCheckService) checkObjectUnique(ctx context.Context, objID, supplierAccount string,
	unique ObjectUnique, attrMap map[int64]metadata.Attribute) error {

	// check if all unique keys are valid
	uniqueFields := make([]metadata.Attribute, 0)
	isValid := true

	keyLen := len(unique.Keys)
	for idx, key := range unique.Keys {
		switch key.Kind {
		case metadata.UniqueKeyKindProperty:
			property, ok := attrMap[int64(key.ID)]
			if !ok {
				isValid = false
				fmt.Printf("object(%s) unique(%d) key(%d) id(%d) not exist\n", objID, unique.ID, idx, key.ID)
				continue
			}

			if !index.ValidateCCFieldType(property.PropertyType, keyLen) {
				isValid = false
				fmt.Printf("attribute(%d) type %s is invalid\n", property.ID, property.PropertyType)
			}

			uniqueFields = append(uniqueFields, property)
		default:
			isValid = false
			fmt.Printf("object(%s) unique(%d) key(%d) kind %s is invalid\n", objID, unique.ID, idx, key.Kind)
			continue
		}
	}

	if !isValid {
		fmt.Printf("object(%s) unique(%d) has invalid key, **skip checking instance**\n", objID, unique.ID)
		return nil
	}

	//  check the uniqueness of the object instances
	filter := make(map[string]interface{})
	if !common.IsInnerModel(objID) {
		filter[common.BKObjIDField] = objID
	}

	attrGroup := make(map[string]interface{})
	for _, attribute := range uniqueFields {
		dbType := index.CCFieldTypeToDBType(attribute.PropertyType)
		if dbType == "" {
			isValid = false
			fmt.Printf("object(%s) attribute(%d) type %s is invalid\n", objID, attribute.ID, attribute.PropertyType)
			continue
		}

		filter[attribute.PropertyID] = map[string]interface{}{common.BKDBType: dbType}
		attrGroup[attribute.PropertyID] = "$" + attribute.PropertyID
	}

	if !isValid {
		fmt.Printf("object(%s) unique(%d) has invalid key, **skip checking instance**\n", objID, unique.ID)
		return nil
	}

	pipeline := []map[string]interface{}{
		{common.BKDBMatch: filter},
		{common.BKDBGroup: map[string]interface{}{
			"_id":   attrGroup,
			"total": map[string]interface{}{common.BKDBSum: 1},
		}},
		{common.BKDBMatch: map[string]interface{}{
			"total": map[string]interface{}{common.BKDBGT: 1},
		}},
	}

	var tableName string
	if common.IsInnerModel(objID) {
		tableName = common.GetInstTableName(objID, supplierAccount)
	} else {
		tableName = common.BKTableNameBaseInst
	}

	items := make([]duplicateItems, 0)
	if err := s.service.DbProxy.Table(tableName).AggregateAll(ctx, pipeline, &items); err != nil {
		return err
	}

	if len(items) > 0 {
		jsItem, _ := json.Marshal(items)
		fmt.Printf("object(%s) unique(%d) has duplicate items(%s)\n", objID, unique.ID, string(jsItem))
		return nil
	}
	return nil
}

type duplicateItems struct {
	Attributes map[string]interface{} `json:"attributes" bson:"_id"`
	Total      int64                  `json:"total" bson:"total"`
}

type ObjectUnique struct {
	ID        uint64               `json:"id" bson:"id"`
	ObjID     string               `json:"bk_obj_id" bson:"bk_obj_id"`
	MustCheck bool                 `json:"must_check" bson:"must_check"`
	Keys      []metadata.UniqueKey `json:"keys" bson:"keys"`
	Ispre     bool                 `json:"ispre" bson:"ispre"`
	OwnerID   string               `json:"bk_supplier_account" bson:"bk_supplier_account"`
	LastTime  *time.Time           `json:"last_time" bson:"last_time"`
}
