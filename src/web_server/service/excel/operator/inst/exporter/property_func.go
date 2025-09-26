package exporter

import (
	"fmt"

	"configcenter/pkg/excel"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/metadata"
	"configcenter/src/web_server/service/excel/core"
)

const emptyCell = ""

var handleColPropFuncMap = make(map[string]handleColPropFunc)

var handleSpecialPropFuncMap = make(map[string]handleColPropFunc)

func init() {
	handleColPropFuncMap[common.FieldTypeInt] = getHandleNumericTypeFunc()
	handleColPropFuncMap[common.FieldTypeFloat] = getHandleNumericTypeFunc()
	handleColPropFuncMap[common.FieldTypeEnum] = getHandleEnumTypeFunc()
	handleColPropFuncMap[common.FieldTypeEnumMulti] = getHandleEnumTypeFunc()
	handleColPropFuncMap[common.FieldTypeBool] = getHandleBoolTypeFunc()
	handleColPropFuncMap[common.FieldTypeInnerTable] = getHandleTableTypeFunc()

	handleSpecialPropFuncMap[common.BKCloudIDField] = getHandleCloudAreaPropFunc()
}

func getHandleColPropFunc(property *core.ColProp) handleColPropFunc {
	handleFunc, isSpecial := handleSpecialPropFuncMap[property.ID]
	if isSpecial {
		return handleFunc
	}

	handleFunc, ok := handleColPropFuncMap[property.PropertyType]
	if !ok {
		handleFunc = getDefaultHandleTypeFunc()
	}

	return handleFunc
}

type handleColPropFunc func(t *TmplOp, property *core.ColProp) ([][]excel.Cell, error)

func getHandleNumericTypeFunc() handleColPropFunc {
	return func(t *TmplOp, property *core.ColProp) ([][]excel.Cell, error) {

		sqref, err := core.GetSingleColSqref(property.ExcelColIndex)
		if err != nil {
			return nil, err
		}
		err = t.GetExcel().AddValidation(t.GetObjID(), &excel.ValidationParam{Type: excel.Decimal, Sqref: sqref})
		if err != nil {
			return nil, err
		}

		handleFunc := getDefaultHandleTypeFunc()
		return handleFunc(t, property)
	}
}

func createSheetWithData(t *TmplOp, sheet string, rowIdx int, data [][]excel.Cell) error {
	if err := t.GetExcel().CreateSheet(sheet); err != nil {
		return err
	}

	if err := t.GetExcel().StreamingWrite(sheet, rowIdx, data); err != nil {
		return err
	}

	if err := t.GetExcel().Flush([]string{sheet}); err != nil {
		return err
	}

	if err := t.GetExcel().Save(); err != nil {
		return err
	}

	return nil
}

func getHandleEnumTypeFunc() handleColPropFunc {
	return func(t *TmplOp, property *core.ColProp) ([][]excel.Cell, error) {
		optionArr, ok := property.Option.([]interface{})
		if ok {

			data := make([][]excel.Cell, len(optionArr))
			for idx, name := range getEnumNames(optionArr) {
				data[idx] = append(data[idx], excel.Cell{Value: name})
			}

			if err := createSheetWithData(t, property.RefSheet, core.NameRowIdx, data); err != nil {
				return nil, err
			}

			if property.PropertyType == common.FieldTypeEnum {
				sqref, err := core.GetSingleColSqref(property.ExcelColIndex)
				if err != nil {
					return nil, err
				}
				if err := t.GetExcel().AddValidation(t.GetObjID(),
					&excel.ValidationParam{Type: excel.Ref, Sqref: sqref, Option: property.RefSheet}); err != nil {
					return nil, err
				}
			}
		}

		handleFunc := getDefaultHandleTypeFunc()
		return handleFunc(t, property)
	}
}

// getEnumNames get enum name from option
func getEnumNames(items []interface{}) []string {
	var names []string
	for _, valRow := range items {
		mapVal, ok := valRow.(map[string]interface{})
		if !ok {
			continue
		}
		name, ok := mapVal[common.BKFieldName].(string)
		if ok {
			names = append(names, name)
		}
	}

	return names
}

func getHandleBoolTypeFunc() handleColPropFunc {
	return func(t *TmplOp, property *core.ColProp) ([][]excel.Cell, error) {
		sqref, err := core.GetSingleColSqref(property.ExcelColIndex)
		if err != nil {
			return nil, err
		}
		if err = t.GetExcel().AddValidation(t.GetObjID(),
			&excel.ValidationParam{Type: excel.Bool, Sqref: sqref, Option: property.Name}); err != nil {
			return nil, err
		}

		handleFunc := getDefaultHandleTypeFunc()
		return handleFunc(t, property)
	}
}

func getHandleTableTypeFunc() handleColPropFunc {
	return func(t *TmplOp, property *core.ColProp) ([][]excel.Cell, error) {
		nameStyle, err := t.styleCreator.getStyle(firstRow)
		if err != nil {
			return nil, err
		}
		headerStyle, err := t.styleCreator.getStyle(generalHeader)
		if err != nil {
			return nil, err
		}

		ccLang := t.GetLang().CreateDefaultCCLanguageIf(httpheader.GetLanguage(t.GetKit().Header))
		propertyType := core.GetTypeAliasName(ccLang, property.PropertyType)

		result := make([][]excel.Cell, core.InstHeaderLen)
		result[core.NameRowIdx] = append(result[core.NameRowIdx], excel.Cell{Value: property.Name, StyleID: nameStyle})
		result[core.TypeRowIdx] = append(result[core.TypeRowIdx],
			excel.Cell{Value: propertyType, StyleID: headerStyle})
		result[core.IDRowIdx] = append(result[core.IDRowIdx], excel.Cell{Value: property.ID, StyleID: headerStyle})

		option, err := metadata.ParseTableAttrOption(property.Option)
		if err != nil {
			return nil, err
		}

		// 设置属性字段相关行，空白单元格的样式
		for i := 1; i < len(option.Header); i++ {
			result[core.NameRowIdx] = append(result[core.NameRowIdx], excel.Cell{StyleID: nameStyle})
			result[core.TypeRowIdx] = append(result[core.TypeRowIdx], excel.Cell{StyleID: headerStyle})
			result[core.IDRowIdx] = append(result[core.IDRowIdx], excel.Cell{StyleID: headerStyle})
		}

		for _, attr := range option.Header {
			colProp := &core.ColProp{ID: attr.PropertyID, Name: attr.PropertyName, PropertyType: attr.PropertyType,
				IsRequire: attr.IsRequired, Option: attr.Option, Group: attr.PropertyGroup, RefSheet: attr.PropertyName}

			if colProp.PropertyType == common.FieldTypeEnumMulti {
				colProp.RefSheet = property.Name + "##" + colProp.Name
			}

			colPropFunc, ok := handleColPropFuncMap[colProp.PropertyType]
			if !ok {
				colPropFunc = getDefaultHandleTypeFunc()
			}

			properyResult, err := colPropFunc(t, colProp)
			if err != nil {
				return nil, err
			}

			if len(properyResult) < core.HeaderTableLen {
				return nil, fmt.Errorf("table type property %s is invalid, option attr: %v", property.ID, attr)
			}

			result[core.TableNameRowIdx] = append(result[core.TableNameRowIdx], properyResult[core.NameRowIdx]...)
			result[core.TableTypeRowIdx] = append(result[core.TableTypeRowIdx], properyResult[core.TypeRowIdx]...)
			result[core.TableIDRowIdx] = append(result[core.TableIDRowIdx], properyResult[core.IDRowIdx]...)
		}

		tableHeaderStyle, err := t.styleCreator.getStyle(tableHeader)
		if err != nil {
			return nil, err
		}

		for _, rowIdx := range tableRowIndexes {
			for idx := range result[rowIdx] {
				result[rowIdx][idx].StyleID = tableHeaderStyle
			}
		}

		return result, nil
	}
}

func getDefaultHandleTypeFunc() handleColPropFunc {
	return func(t *TmplOp, property *core.ColProp) ([][]excel.Cell, error) {
		nameStyleType := firstRow
		headerStyleType := generalHeader
		if property.NotEditable {
			nameStyleType = noEditHeader
			headerStyleType = noEditHeader
		}

		nameStyle, err := t.styleCreator.getStyle(nameStyleType)
		if err != nil {
			return nil, err
		}
		headerStyle, err := t.styleCreator.getStyle(headerStyleType)
		if err != nil {
			return nil, err
		}

		ccLang := t.GetLang().CreateDefaultCCLanguageIf(httpheader.GetLanguage(t.GetKit().Header))
		propertyType := core.GetTypeAliasName(ccLang, property.PropertyType)

		result := make([][]excel.Cell, core.InstHeaderLen)
		result[core.NameRowIdx] = append(result[core.NameRowIdx], excel.Cell{Value: property.Name, StyleID: nameStyle})
		result[core.TypeRowIdx] = append(result[core.TypeRowIdx], excel.Cell{Value: propertyType, StyleID: headerStyle})
		result[core.IDRowIdx] = append(result[core.IDRowIdx], excel.Cell{Value: property.ID, StyleID: headerStyle})

		result[core.TableNameRowIdx] = append(result[core.TableNameRowIdx],
			excel.Cell{Value: emptyCell, StyleID: headerStyle})
		result[core.TableTypeRowIdx] = append(result[core.TableTypeRowIdx],
			excel.Cell{Value: emptyCell, StyleID: headerStyle})
		result[core.TableIDRowIdx] = append(result[core.TableIDRowIdx],
			excel.Cell{Value: emptyCell, StyleID: headerStyle})

		return result, nil
	}
}

func getHandleCloudAreaPropFunc() handleColPropFunc {
	return func(t *TmplOp, property *core.ColProp) ([][]excel.Cell, error) {
		cloudAreaArr, cloudAreaMap, err := t.GetClient().GetCloudArea(t.GetKit())
		if err != nil {
			blog.Errorf("get cloud area failed, err: %v, rid: %s", err, t.GetKit().Rid)
			return nil, err
		}

		if err := t.GetExcel().CreateSheet(property.RefSheet); err != nil {
			return nil, err
		}

		data := make([][]excel.Cell, len(cloudAreaArr))
		for idx, cloudArea := range cloudAreaArr {
			data[idx] = append(data[idx], excel.Cell{Value: spliceCloudArea(cloudArea, cloudAreaMap[cloudArea])})
		}
		if err := t.GetExcel().StreamingWrite(property.RefSheet, core.NameRowIdx, data); err != nil {
			return nil, err
		}

		if err := t.GetExcel().Flush([]string{property.RefSheet}); err != nil {
			return nil, err
		}

		if err := t.GetExcel().Save(); err != nil {
			return nil, err
		}

		sqref, err := core.GetSingleColSqref(property.ExcelColIndex)
		if err != nil {
			return nil, err
		}
		if err := t.GetExcel().AddValidation(t.GetObjID(),
			&excel.ValidationParam{Type: excel.Ref, Sqref: sqref, Option: property.RefSheet}); err != nil {
			return nil, err
		}

		handleFunc := getDefaultHandleTypeFunc()
		return handleFunc(t, property)
	}
}

func spliceCloudArea(cloudAreaName interface{}, cloudAreaID interface{}) string {
	return fmt.Sprintf("%v[%v]", cloudAreaName, cloudAreaID)
}
