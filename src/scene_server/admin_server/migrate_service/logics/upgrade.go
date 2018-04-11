package logics

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/migrate_service/data"
	"configcenter/src/scene_server/validator"
	dbStorage "configcenter/src/storage"
	"strconv"
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
	dataRows := data.AppRow()
	dataRows = append(dataRows, data.SetRow()...)
	dataRows = append(dataRows, data.ModuleRow()...)
	dataRows = append(dataRows, data.HostRow()...)
	dataRows = append(dataRows, data.ProcRow()...)
	dataRows = append(dataRows, data.PlatRow()...)

	for _, expectRow := range dataRows {
		expectOptions, ok := expectRow.Option.([]validator.EnumVal)
		if !ok {
			continue
		}
		selector := map[string]interface{}{
			common.BKObjIDField:      expectRow.ObjectID,
			common.BKPropertyIDField: expectRow.PropertyID,
			common.BKOwnerIDField:    expectRow.OwnerID,
		}
		curRow := map[string]interface{}{}
		err := db.GetOneByCondition(common.BKTableNameObjAttDes, nil, selector, &curRow)
		if err != nil {
			blog.Errorf("upgradeGlobalization get row error: %v", err)
			return err
		}

		curOptions := validator.ParseEnumOption(curRow["option"])
		newOptions := []validator.EnumVal{}
		newID := len(expectOptions)
		// get custom options
		for _, curOption := range curOptions {
			if curOption.ID != "" {
				// if ID!="" then we believe this property has upgraded so we just ignore it
				continue
			}
			exists := false
			for _, expectOption := range expectOptions {
				if expectOption.Name == curOption.Name {
					newOptions = append(newOptions, expectOption)
					exists = true
					break
				}
			}
			if !exists {
				newID++
				curOption.ID = strconv.Itoa(newID)
				newOptions = append(newOptions, curOption)
			}
		}

		if len(newOptions) <= 0 {
			continue
		}

		// update property's option fields
		updatedata := map[string]interface{}{
			common.BKOptionField: newOptions,
		}
		err = db.UpdateByCondition(common.BKTableNameObjAttDes, updatedata, selector)
		if err != nil {
			blog.Errorf("upgradeGlobalization update option field error: %v", err)
			return err
		}

		isExist, err := db.GetCntByCondition(common.BKTableNameObjAttDes, selector)
		if nil != err {
			blog.Errorf("add data for  %s table error  %s", common.BKTableNameObjAttDes, err)
			return err
		}
		if isExist > 0 {
			continue
		}
		id, err := db.GetIncID(common.BKTableNameObjAttDes)
		if nil != err {
			blog.Errorf("add data for  %s table error  %s", common.BKTableNameObjAttDes, err)
			return err
		}
		expectRow.ID = int(id)
		_, err = db.Insert(common.BKTableNameObjAttDes, expectRow)
		if nil != err {
			blog.Errorf("add data for  %s table error  %s", common.BKTableNameObjAttDes, err)
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
