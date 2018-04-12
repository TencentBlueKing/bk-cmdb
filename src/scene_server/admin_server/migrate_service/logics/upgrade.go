package logics

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/migrate_service/data"
	"configcenter/src/scene_server/validator"
	"configcenter/src/source_controller/api/metadata"
	"configcenter/src/source_controller/common/commondata"
	dbStorage "configcenter/src/storage"
	"sort"
	"strconv"
	"strings"
)

// Upgrade upgrade
func Upgrade(instData dbStorage.DI) error {
	err := upgradeAppfield(instData)
	if err != nil {
		return err
	}
	err = upgradeGlobalization(instData)
	if err != nil {
		return err
	}
	return nil
}

func upgradeGlobalization(db dbStorage.DI) error {
	presetRows := data.AppRow()
	presetRows = append(presetRows, data.HostRow()...)
	presetRows = append(presetRows, data.ModuleRow()...)
	presetRows = append(presetRows, data.PlatRow()...)
	presetRows = append(presetRows, data.ProcRow()...)
	presetRows = append(presetRows, data.SetRow()...)
	presetRowsMap := map[string]*metadata.ObjectAttDes{}
	for _, row := range presetRows {
		presetRowsMap[row.ObjectID+"::"+row.PropertyID] = row
	}

	attdes := []metadata.ObjectAttDes{}
	err := db.GetMutilByCondition(common.BKTableNameObjAttDes, nil, map[string]interface{}{common.BKPropertyTypeField: common.FiledTypeEnum}, &attdes, "", 0, 0)
	if err != nil {
		blog.Errorf("upgradeGlobalization get attdes error: %v", err)
		return err
	}
	for _, curRow := range attdes {
		curOptions := validator.ParseEnumOption(curRow.Option)
		if len(curOptions) <= 0 {
			continue
		}

		expectRow := presetRowsMap[curRow.ObjectID+"::"+curRow.PropertyID]
		expectOptions, ok := expectRow.Option.([]validator.EnumVal)
		if !ok {
			expectOptions = []validator.EnumVal{}
		}

		sort.SliceStable(curOptions, func(i, j int) bool {
			return strings.Compare(curOptions[i].Name, curOptions[j].Name) < 0
		})
		// get max id
		newID := len(expectOptions)
		for _, option := range curOptions {
			id, _ := strconv.Atoi(option.ID) // ignore err cause we just want the max id
			if id > newID {
				newID = id
			}
		}

		newOptions := []validator.EnumVal{}
		// get custom options
		for _, curOption := range curOptions {
			exists := false
			for _, expectOption := range expectOptions {
				if expectOption.Name == curOption.Name {
					newOptions = append(newOptions, expectOption)
					exists = true
					break
				}
			}
			if !exists {
				if curOption.ID != "" {
					newOptions = append(newOptions, curOption)
					continue
				}
				newID++
				curOption.ID = strconv.Itoa(newID)
				newOptions = append(newOptions, curOption)
			}
		}

		if len(newOptions) <= 0 {
			continue
		}

		// update inst
		tablename := commondata.GetInstTableName(expectRow.ObjectID)
		for _, option := range newOptions {
			updateinstdata := map[string]interface{}{
				expectRow.PropertyID: option.ID,
			}
			updateinstcondition := map[string]interface{}{
				expectRow.PropertyID: option.Name,
			}
			if tablename == common.BKTableNameBaseInst {
				updateinstcondition[common.BKObjIDField] = expectRow.ObjectID
			}

			db.UpdateByCondition(tablename, updateinstdata, updateinstcondition)
			if err != nil {
				blog.Errorf("upgradeGlobalization update preset inst error: %v", err)
				return err
			}
		}

		// update property's option fields
		selector := map[string]interface{}{
			common.BKObjIDField:      curRow.ObjectID,
			common.BKPropertyIDField: curRow.PropertyID,
		}
		updatedata := map[string]interface{}{
			common.BKOptionField: newOptions,
		}
		err = db.UpdateByCondition(common.BKTableNameObjAttDes, updatedata, selector)
		if err != nil {
			blog.Errorf("upgradeGlobalization update preset option field error: %v", err)
			return err
		}
	}
	return nil
}

func upgradeAppfield(instData dbStorage.DI) error {
	condition := map[string]interface{}{
		common.BKObjIDField: common.BKInnerObjIDApp,
		common.BKPropertyIDField: map[string]interface{}{
			"$in": []string{
				"time_zone",
				"language",
			},
		},
	}
	data := map[string]interface{}{
		"isrequired": true,
	}
	instData.UpdateByCondition(common.BKTableNameObjAttDes, data, condition)
	return nil
}
