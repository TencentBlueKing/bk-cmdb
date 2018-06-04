package v3v0v8

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/migrate_service/upgrader"
	"configcenter/src/scene_server/validator"
	"configcenter/src/storage"
)

func addDefaultBiz(db storage.DI, conf *upgrader.Config) error {
	// add default biz
	defaultBiz := map[string]interface{}{}
	defaultBiz[common.BKAppNameField] = common.DefaultAppName
	defaultBiz[common.BKMaintainersField] = "admin"
	defaultBiz[common.BKProductPMField] = "admin"
	defaultBiz[common.BKTimeZoneField] = "Asia/Shanghai"
	defaultBiz[common.BKLanguageField] = "1" //中文
	defaultBiz[common.BKLifeCycleField] = common.DefaultAppLifeCycleNormal
	defaultBiz[common.BKOwnerIDField] = conf.OwnerID
	defaultBiz[common.BKSupplierIDField] = common.BKDefaultSupplierID
	defaultBiz[common.BKDefaultField] = common.DefaultAppFlag
	bizID, err := db.GetIncID("cc_ApplicationBase")
	if err != nil {
		return err
	}
	defaultBiz[common.BKAppIDField] = bizID
	filled := fillEmptyFields(defaultBiz)
	if err := storage.Upsert(db, "cc_ApplicationBase", defaultBiz, []string{common.BKOwnerIDField, common.BKDefaultField}, append(filled, common.BKAppIDField)); err != nil {
		blog.Error("add defaultBiz error ", err.Error())
		return err
	}

	// add default set
	defaultSet := make(map[string]interface{})
	defaultSet[common.BKAppIDField] = bizID
	defaultSet[common.BKInstParentStr] = bizID
	defaultSet[common.BKSetNameField] = common.DefaultResSetName
	defaultSet[common.BKDefaultField] = common.DefaultResSetFlag
	defaultSet[common.BKOwnerIDField] = conf.OwnerID
	filled = fillEmptyFields(defaultSet)
	setID, err := db.GetIncID("cc_SetBase")
	if err != nil {
		return err
	}
	defaultSet[common.BKSetIDField] = setID
	if err := storage.Upsert(db, "cc_SetBase", defaultSet, []string{common.BKOwnerIDField, common.BKDefaultField}, append(filled, common.BKSetIDField)); err != nil {
		blog.Error("add defaultSet error ", err.Error())
		return err
	}

	// add default module
	defaultResModule := make(map[string]interface{})
	defaultResModule[common.BKSetIDField] = setID
	defaultResModule[common.BKInstParentStr] = setID
	defaultResModule[common.BKAppIDField] = bizID
	defaultResModule[common.BKModuleNameField] = common.DefaultResModuleName
	defaultResModule[common.BKDefaultField] = common.DefaultResModuleFlag
	defaultResModule[common.BKOwnerIDField] = conf.OwnerID
	filled = fillEmptyFields(defaultResModule)
	defaultResModuleID, err := db.GetIncID("cc_ModuleBase")
	if err != nil {
		return err
	}
	defaultResModule[common.BKModuleIDField] = defaultResModuleID
	if err := storage.Upsert(db, "cc_ModuleBase", defaultResModule, []string{common.BKOwnerIDField, common.BKDefaultField}, append(filled, common.BKModuleIDField)); err != nil {
		blog.Error("add defaultResModule error ", err.Error())
		return err
	}
	defaultFaultModule := make(map[string]interface{})
	defaultFaultModule[common.BKSetIDField] = setID
	defaultFaultModule[common.BKInstParentStr] = setID
	defaultFaultModule[common.BKAppIDField] = bizID
	defaultFaultModule[common.BKModuleNameField] = common.DefaultFaultModuleName
	defaultFaultModule[common.BKDefaultField] = common.DefaultFaultModuleFlag
	defaultFaultModule[common.BKOwnerIDField] = conf.OwnerID
	filled = fillEmptyFields(defaultFaultModule)
	defaultFaultModuleID, err := db.GetIncID("cc_ModuleBase")
	if err != nil {
		return err
	}
	defaultFaultModule[common.BKModuleIDField] = defaultFaultModuleID
	if err := storage.Upsert(db, "cc_ModuleBase", defaultFaultModule, []string{common.BKOwnerIDField, common.BKDefaultField}, append(filled, common.BKModuleIDField)); err != nil {
		blog.Error("add defaultFaultModule error ", err.Error())
		return err
	}

	return nil
}

func fillEmptyFields(data map[string]interface{}) []string {
	filled := []string{}
	for _, field := range AppRow() {
		fieldName := field.PropertyID
		fieldType := field.PropertyType
		if _, ok := data[fieldName]; ok {
			continue
		}
		option := field.Option
		switch fieldType {
		case common.FieldTypeSingleChar:
			data[fieldName] = ""
		case common.FieldTypeLongChar:
			data[fieldName] = ""
		case common.FieldTypeInt:
			data[fieldName] = nil
		case common.FieldTypeEnum:
			enumOptions := validator.ParseEnumOption(option)
			v := ""
			if len(enumOptions) > 0 {
				var defaultOption *validator.EnumVal
				for _, k := range enumOptions {
					if k.IsDefault {
						defaultOption = &k
						break
					}
				}
				if nil != defaultOption {
					v = defaultOption.ID
				}
			}
			data[fieldName] = v
		case common.FieldTypeDate:
			data[fieldName] = ""
		case common.FieldTypeTime:
			data[fieldName] = ""
		case common.FieldTypeUser:
			data[fieldName] = ""
		case common.FieldTypeMultiAsst:
			data[fieldName] = nil
		case common.FieldTypeTimeZone:
			data[fieldName] = nil
		case common.FieldTypeBool:
			data[fieldName] = false
		default:
			data[fieldName] = nil
		}
		filled = append(filled, fieldName)
	}
	return filled
}
