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
	"configcenter/src/storage/dal/types"
	"configcenter/src/tools/cmdb_ctl/app/config"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(NewMigrateCheckCommand())
}

// NewMigrateCheckCommand new tool command for checking before migration
func NewMigrateCheckCommand() *cobra.Command {
	checkAll := new(bool)
	cmd := &cobra.Command{
		Use:   "migrate-check",
		Short: "check if data satisfies migration constraints",
		RunE: func(cmd *cobra.Command, args []string) error {
			if *checkAll {
				if err := runUniqueCheck(); err != nil {
					return err
				}
				if err := runProcCheck(false); err != nil {
					return err
				}
				return nil
			} else {
				return cmd.Help()
			}
		},
	}
	cmd.Flags().BoolVar(checkAll, "check-all", false, "check if all data satisfies migration constraints")

	cmd.AddCommand(&cobra.Command{
		Use:   "unique",
		Short: "check if object instances satisfies unique constraints",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUniqueCheck()
		},
	})

	clearProc := new(bool)
	uniqueCmd := &cobra.Command{
		Use:   "process",
		Short: "check if process has relation",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runProcCheck(*clearProc)
		},
	}
	uniqueCmd.Flags().BoolVar(clearProc, "clear-proc", false, "clear process with no relation")
	cmd.AddCommand(uniqueCmd)

	return cmd
}

type migrateCheckService struct {
	service *config.Service
}

func newMigrateCheckService(mongoURI string, mongoRsName string) (*migrateCheckService, error) {
	service, err := config.NewMongoService(mongoURI, mongoRsName)
	if err != nil {
		return nil, err
	}
	return &migrateCheckService{
		service: service,
	}, nil
}

func runUniqueCheck() error {
	srv, err := newMigrateCheckService(config.Conf.MongoURI, config.Conf.MongoRsName)
	if err != nil {
		return err
	}
	return srv.checkUnique()
}

func runProcCheck(clearProc bool) error {
	srv, err := newMigrateCheckService(config.Conf.MongoURI, config.Conf.MongoRsName)
	if err != nil {
		return err
	}
	if clearProc {
		return srv.clearProc()
	} else {
		return srv.checkProc()
	}
}

func (s *migrateCheckService) checkUnique() error {
	fmt.Println("=================================")
	printInfo("start checking unique constraints\n")

	ctx := util.SetDBReadPreference(context.Background(), common.SecondaryPreferredMode)

	objIDs, err := s.getAllObjectIDs(ctx)
	if err != nil {
		return err
	}

	for _, objID := range objIDs {
		fmt.Println("=================================")
		if !model.SatisfyMongoCollLimit(objID) {
			printError("object id %s is invalid\n", objID)
		}

		attrMap, err := s.getObjAttrMap(ctx, objID)
		if err != nil {
			return err
		}

		printInfo("start checking unique constraints for object %s\n", objID)
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
					printError("object(%s) unique(%d) key(%d) duplicate with %d\n", objID, unique.ID, key.ID, id)
					continue
				}
				uniqueKeyMap[key.ID] = unique.ID
			}

			if !isValid {
				printError("object(%s) unique(%d) has invalid key, **skip checking instance**\n", objID, unique.ID)
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

	printInfo("checking unique constraints done\n")
	return nil
}

func (s *migrateCheckService) getAllObjectIDs(ctx context.Context) ([]string, error) {
	printInfo("start searching for all object ids\n")

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

func (s *migrateCheckService) getObjectUniques(ctx context.Context, objID string) ([]ObjectUnique, error) {
	printInfo("start searching unique constraints for object %s\n", objID)

	filter := map[string]interface{}{
		common.BKObjIDField: objID,
	}

	uniques := make([]ObjectUnique, 0)
	if err := s.service.DbProxy.Table(common.BKTableNameObjUnique).Find(filter).All(ctx, &uniques); err != nil {
		return nil, fmt.Errorf("get unique rules for object(%s) failed, err: %v", objID, err)
	}

	return uniques, nil
}

func (s *migrateCheckService) getObjAttrMap(ctx context.Context, objID string) (map[int64]metadata.Attribute, error) {
	printInfo("start searching object attributes for object %s\n", objID)

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
	for _, attr := range attributes {
		if !model.SatisfyMongoFieldLimit(attr.PropertyID) {
			printError("object(%s) attribute(%d) property id %s is invalid\n", objID, attr.ID, attr.PropertyID)
		}

		if objID == common.BKInnerObjIDHost && attr.PropertyID == common.BKCloudIDField {
			// NOTICE: 2021年03月12日 特殊逻辑。 现在主机的字段中类型未foreignkey 特殊的类型
			attr.PropertyType = common.FieldTypeInt
		}
		if objID == common.BKInnerObjIDHost &&
			(attr.PropertyID == common.BKHostInnerIPField || attr.PropertyID == common.BKHostOuterIPField ||
				attr.PropertyID == common.BKOperatorField || attr.PropertyID == common.BKBakOperatorField ||
				attr.PropertyID == common.BKHostInnerIPv6Field || attr.PropertyID == common.BKHostOuterIPv6Field) {
			// NOTICE: 2021年03月12日 特殊逻辑。 现在主机的字段中类型未innerIP,OuterIP 特殊的类型
			attr.PropertyType = common.FieldTypeList
		}

		attrMap[attr.ID] = attr
	}

	return attrMap, nil
}

func (s *migrateCheckService) checkObjectUnique(ctx context.Context, objID, supplierAccount string,
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
				printError("object(%s) unique(%d) key(%d) id(%d) not exist\n", objID, unique.ID, idx, key.ID)
				continue
			}

			if !index.ValidateCCFieldType(property.PropertyType, keyLen) {
				isValid = false
				printError("attribute(%d) type %s is invalid\n", property.ID, property.PropertyType)
			}

			uniqueFields = append(uniqueFields, property)
		default:
			isValid = false
			printError("object(%s) unique(%d) key(%d) kind %s is invalid\n", objID, unique.ID, idx, key.Kind)
			continue
		}
	}

	if !isValid {
		printError("object(%s) unique(%d) has invalid key, **skip checking instance**\n", objID, unique.ID)
		return nil
	}

	if objID == common.BKInnerObjIDPlat && len(uniqueFields) == 1 && uniqueFields[0].PropertyID == "bk_vpc_id" {
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
			printError("object(%s) attribute(%d) type %s is invalid\n", objID, attribute.ID, attribute.PropertyType)
			continue
		}

		filter[attribute.PropertyID] = map[string]interface{}{common.BKDBType: dbType}
		attrGroup[attribute.PropertyID] = "$" + attribute.PropertyID
	}

	if !isValid {
		printError("object(%s) unique(%d) has invalid key, **skip checking instance**\n", objID, unique.ID)
		return nil
	}

	pipeline := []map[string]interface{}{
		{common.BKDBMatch: filter},
		{common.BKDBGroup: map[string]interface{}{
			"_id":   attrGroup,
			"total": map[string]interface{}{common.BKDBSum: 1},
		}},
		{common.BKDBMatch: map[string]interface{}{"total": map[string]interface{}{common.BKDBGT: 1}}},
	}

	var tableName string
	if common.IsInnerModel(objID) {
		tableName = common.GetInstTableName(objID, supplierAccount)
	} else {
		tableName = common.BKTableNameBaseInst
	}

	items := make([]duplicateItems, 0)
	aggregateOpt := types.NewAggregateOpts().SetAllowDiskUse(true)
	if err := s.service.DbProxy.Table(tableName).AggregateAll(ctx, pipeline, &items, aggregateOpt); err != nil {
		return err
	}

	if len(items) > 0 {
		jsItem, _ := json.Marshal(items)
		printError("object(%s) unique(%d) has duplicate items(%s)\n", objID, unique.ID, string(jsItem))
		return nil
	}
	return nil
}

type duplicateItems struct {
	Attributes map[string]interface{} `json:"attributes" bson:"_id"`
	Total      int64                  `json:"total" bson:"total"`
}

// ObjectUnique TODO
type ObjectUnique struct {
	ID        uint64               `json:"id" bson:"id"`
	ObjID     string               `json:"bk_obj_id" bson:"bk_obj_id"`
	MustCheck bool                 `json:"must_check" bson:"must_check"`
	Keys      []metadata.UniqueKey `json:"keys" bson:"keys"`
	Ispre     bool                 `json:"ispre" bson:"ispre"`
	OwnerID   string               `json:"bk_supplier_account" bson:"bk_supplier_account"`
	LastTime  *time.Time           `json:"last_time" bson:"last_time"`
}

func printInfo(format string, a ...interface{}) {
	fmt.Printf("INFO: %s", fmt.Sprintf(format, a...))
}

func printError(format string, a ...interface{}) {
	fmt.Printf("ERROR: %s", fmt.Sprintf(format, a...))
}

func (s *migrateCheckService) checkProc() error {
	fmt.Println("=================================")
	printInfo("start checking process with no relation\n")

	ctx := util.SetDBReadPreference(context.Background(), common.SecondaryPreferredMode)

	procIDs, err := s.getProcWithNoRelation(ctx)
	if err != nil {
		return err
	}

	if len(procIDs) > 0 {
		procLen := len(procIDs)
		allProcesses := make([]metadata.Process, 0)
		for start := 0; start < procLen; start += common.BKMaxPageSize {
			limit := start + common.BKMaxPageSize
			if limit > procLen {
				limit = procLen
			}
			procFilter := map[string]interface{}{
				common.BKProcessIDField: map[string]interface{}{common.BKDBIN: procIDs[start:limit]},
			}

			processes := make([]metadata.Process, 0)
			err := s.service.DbProxy.Table(common.BKTableNameBaseProcess).Find(procFilter).All(ctx, &processes)
			if err != nil {
				return fmt.Errorf("get processes with no relation failed, procIDs: %+v, err: %v", procIDs, err)
			}

			allProcesses = append(allProcesses, processes...)
		}

		allProcessesJson, err := json.Marshal(allProcesses)
		if err != nil {
			printError("processes has no relations, need to delete, data: %#v\n", allProcesses)
		} else {
			printError("processes has no relations, need to delete, data: %s\n", string(allProcessesJson))
		}
	}

	printInfo("checking process with no relation done\n")
	return nil
}

func (s *migrateCheckService) clearProc() error {
	fmt.Println("=================================")
	printInfo("start clearing processes with no instances\n")

	ctx := util.SetDBReadPreference(context.Background(), common.SecondaryPreferredMode)

	procIDs, err := s.getProcWithNoRelation(ctx)
	if err != nil {
		return err
	}

	procLen := len(procIDs)
	if procLen > 0 {
		for start := 0; start < procLen; start += common.BKMaxPageSize {
			limit := start + common.BKMaxPageSize
			if limit > procLen {
				limit = procLen
			}
			procFilter := map[string]interface{}{
				common.BKProcessIDField: map[string]interface{}{common.BKDBIN: procIDs[start:limit]},
			}

			err := s.service.DbProxy.Table(common.BKTableNameBaseProcess).Delete(ctx, procFilter)
			if err != nil {
				return fmt.Errorf("delete processes with no relation failed, procIDs: %+v, err: %v", procIDs, err)
			}
		}

		printInfo("clear processes successful, ids: %+v\n", procIDs)
	}

	printInfo("clearing process with no relation done\n")
	return nil
}

func (s *migrateCheckService) getProcWithNoRelation(ctx context.Context) ([]int64, error) {
	noRelationProcIDs := make([]int64, 0)

	for start := uint64(0); ; start += common.BKMaxPageSize {
		processes := make([]struct {
			ProcessID int64 `bson:"bk_process_id"`
		}, 0)
		err := s.service.DbProxy.Table(common.BKTableNameBaseProcess).Find(nil).Fields(common.BKProcessIDField).
			Start(start).Limit(common.BKMaxPageSize).All(ctx, &processes)
		if err != nil {
			return nil, fmt.Errorf("get process ids failed, err: %v", err)
		}

		if len(processes) == 0 {
			break
		}

		procIDs := make([]int64, len(processes))
		procIDMap := make(map[int64]struct{})
		for idx, process := range processes {
			procIDs[idx] = process.ProcessID
			procIDMap[process.ProcessID] = struct{}{}
		}

		// get process id to service instance id relations
		relationFilter := map[string]interface{}{
			common.BKProcessIDField: map[string]interface{}{common.BKDBIN: procIDs},
		}
		procRelations := make([]metadata.ProcessInstanceRelation, 0)
		if err := s.service.DbProxy.Table(common.BKTableNameProcessInstanceRelation).Find(relationFilter).Fields(
			common.BKProcessIDField).All(ctx, &procRelations); err != nil {
			return nil, fmt.Errorf("get process relations failed, procIDs: %+v, err: %v", procIDs, err)
		}

		for _, relation := range procRelations {
			delete(procIDMap, relation.ProcessID)
		}

		for procID := range procIDMap {
			noRelationProcIDs = append(noRelationProcIDs, procID)
		}
	}

	return noRelationProcIDs, nil
}
