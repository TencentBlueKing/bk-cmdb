package v3v0v8

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/paraparse"
	"configcenter/src/scene_server/admin_server/migrate_service/upgrader"
	"configcenter/src/storage"
	"time"
)

func addDefaultBiz(db storage.DI, conf *upgrader.Config) error {
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
	createBiz(defaultBiz)
	return nil
}

func createBiz(data map[string]interface{}) {
	tablename := "cc_ApplicationBase"
}

func fillEmptyFields(data map[string]interface{}) {
	type intOptionType struct {
		Min int
		Max int
	}
	type EnumOptionType struct {
		Name string
		Type string
	}
	for _, field := range AppRow() {
		fieldName := field.PropertyID
		fieldType := field.PropertyType
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
			continue
		}

	}

}
