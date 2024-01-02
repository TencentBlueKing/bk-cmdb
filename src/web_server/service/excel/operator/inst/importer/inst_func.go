package importer

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"
	"configcenter/src/web_server/service/excel/core"
)

// ExcelMsg excel file message
type ExcelMsg struct {
	propertyMap map[int]PropWithTable // key为字段在导入的excel所处的位置，value为字段属性
	mergeRowRes map[int]int           // 对同一列的行进行合并的开始和结束范围, 按行从1开始计数
}

// PropWithTable property with sub table property
type PropWithTable struct {
	core.ColProp
	subProperties map[int]PropWithTable // 只有表格类型才有这个属性，表示表格的每一列属性; key为在导入的excel所处位置，value为字段属性
}

var handleInstFieldFuncMap = make(map[string]handleInstFieldFunc)

var handleSpecialFieldFuncMap = make(map[string]handleInstFieldFunc)

func init() {
	handleInstFieldFuncMap[common.FieldTypeInt] = getHandleIntFieldFunc()
	handleInstFieldFuncMap[common.FieldTypeFloat] = getHandleFloatFieldFunc()
	handleInstFieldFuncMap[common.FieldTypeBool] = getHandleBoolFieldFunc()
	handleInstFieldFuncMap[common.FieldTypeEnum] = getHandleEnumFieldFunc()
	handleInstFieldFuncMap[common.FieldTypeEnumMulti] = getHandleEnumMultiFieldFunc()
	handleInstFieldFuncMap[common.FieldTypeEnumQuote] = getHandleEnumQuoteFieldFunc()
	handleInstFieldFuncMap[common.FieldTypeOrganization] = getHandleOrgFieldFunc()
	handleInstFieldFuncMap[common.FieldTypeUser] = getHandleUserFieldFunc()
	handleInstFieldFuncMap[common.FieldTypeInnerTable] = getHandleTableFieldFunc()

	handleSpecialFieldFuncMap[common.BKCloudIDField] = getCloudAreaFieldFunc()
}

func getHandleInstFieldFunc(prop *PropWithTable) handleInstFieldFunc {
	handleFunc, isSpecial := handleSpecialFieldFuncMap[prop.ID]
	if isSpecial {
		return handleFunc
	}

	handleFunc, ok := handleInstFieldFuncMap[prop.PropertyType]
	if !ok {
		handleFunc = getDefaultHandleFieldFunc()
	}

	return handleFunc
}

type handleInstFieldFunc func(i *Importer, property *PropWithTable, rows [][]string) (interface{}, error)

func getHandleIntFieldFunc() handleInstFieldFunc {
	return func(i *Importer, property *PropWithTable, rows [][]string) (interface{}, error) {
		if len(rows) == 0 || len(rows[0]) < property.ExcelColIndex {
			blog.Errorf("instance is invalid, data: %v, rid: %s", rows, i.GetKit().Rid)
			return nil, fmt.Errorf("instance is invalid")
		}

		val := rows[0][property.ExcelColIndex]
		intVal, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			blog.Errorf("failed to convert string type to int type, val: %v, err: %v, rid: %s", val, err,
				i.GetKit().Rid)
			return nil, err
		}

		return intVal, nil
	}
}

func getHandleFloatFieldFunc() handleInstFieldFunc {
	return func(i *Importer, property *PropWithTable, rows [][]string) (interface{}, error) {
		if len(rows) == 0 || len(rows[0]) < property.ExcelColIndex {
			blog.Errorf("instance is invalid, data: %v, rid: %s", rows, i.GetKit().Rid)
			return nil, fmt.Errorf("instance is invalid")
		}

		val := rows[0][property.ExcelColIndex]
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			blog.Errorf("failed to convert string type to float type, val: %v, err: %v, rid: %s", val, err,
				i.GetKit().Rid)
			return nil, err
		}

		return floatVal, nil
	}
}

func getHandleBoolFieldFunc() handleInstFieldFunc {
	return func(i *Importer, property *PropWithTable, rows [][]string) (interface{}, error) {
		if len(rows) == 0 || len(rows[0]) < property.ExcelColIndex {
			blog.Errorf("instance is invalid, data: %v, rid: %s", rows, i.GetKit().Rid)
			return nil, fmt.Errorf("instance is invalid")
		}

		val := rows[0][property.ExcelColIndex]
		bolVal, err := strconv.ParseBool(val)
		if err != nil {
			blog.Errorf("failed to convert string type to bool type, val: %v, err: %v, rid: %s", val, err,
				i.GetKit().Rid)
			return nil, err
		}

		return bolVal, nil
	}
}

func getHandleEnumFieldFunc() handleInstFieldFunc {
	return func(i *Importer, property *PropWithTable, rows [][]string) (interface{}, error) {
		if len(rows) == 0 || len(rows[0]) < property.ExcelColIndex {
			blog.Errorf("instance is invalid, data: %v, rid: %s", rows, i.GetKit().Rid)
			return nil, fmt.Errorf("instance is invalid")
		}

		val := rows[0][property.ExcelColIndex]

		if option, optionOk := property.Option.([]interface{}); optionOk {
			return getEnumIDByName(val, option), nil
		}

		return val, nil
	}
}

func getHandleEnumMultiFieldFunc() handleInstFieldFunc {
	return func(i *Importer, property *PropWithTable, rows [][]string) (interface{}, error) {
		if len(rows) == 0 || len(rows[0]) < property.ExcelColIndex {
			blog.Errorf("instance is invalid, data: %v, rid: %s", rows, i.GetKit().Rid)
			return nil, fmt.Errorf("instance is invalid")
		}

		val := rows[0][property.ExcelColIndex]

		option, optionOk := property.Option.([]interface{})
		if !optionOk {
			return val, nil
		}

		nameList := strings.Split(val, "\n")
		idList := make([]string, 0)
		for _, name := range nameList {
			id := getEnumIDByName(name, option)
			idList = append(idList, id)
		}

		return idList, nil
	}
}

func getHandleEnumQuoteFieldFunc() handleInstFieldFunc {
	return func(i *Importer, property *PropWithTable, rows [][]string) (interface{}, error) {
		if len(rows) == 0 || len(rows[0]) < property.ExcelColIndex {
			blog.Errorf("instance is invalid, data: %v, rid: %s", rows, i.GetKit().Rid)
			return nil, fmt.Errorf("instance is invalid")
		}

		names := strings.Split(rows[0][property.ExcelColIndex], "\n")
		ids, err := i.GetClient().TransEnumQuoteNameToID(i.GetKit(), names, &property.ColProp)
		if err != nil {
			blog.Errorf("transfer enum quote name to id failed, names: %v, err: %v, rid: %s", names, err,
				i.GetKit().Rid)
			return nil, err
		}

		return ids, nil
	}
}

// getEnumIDByName get enum id from option name
func getEnumIDByName(name string, items []interface{}) string {
	id := name

	for _, valRow := range items {
		mapVal, ok := valRow.(map[string]interface{})
		if !ok {
			continue
		}

		enumName, ok := mapVal[common.BKFieldName].(string)
		if !ok {
			continue
		}

		if enumName != name {
			continue
		}

		enumID, ok := mapVal[common.BKFieldID].(string)
		if ok {
			id = enumID
		}
	}

	return id
}

const (
	organizationBracketsPattern = `\[(\d+)\]([^\s]+)`
	userBracketsPattern         = `\([a-zA-Z0-9\@\p{Han} .,_-]*\)`
)

var (
	orgBracketsRegexp  = regexp.MustCompile(organizationBracketsPattern)
	userBracketsRegexp = regexp.MustCompile(userBracketsPattern)
)

func getHandleOrgFieldFunc() handleInstFieldFunc {
	return func(i *Importer, property *PropWithTable, rows [][]string) (interface{}, error) {
		if len(rows) == 0 || len(rows[0]) < property.ExcelColIndex {
			blog.Errorf("instance is invalid, data: %v, rid: %s", rows, i.GetKit().Rid)
			return nil, fmt.Errorf("instance is invalid")
		}

		val := rows[0][property.ExcelColIndex]
		if len(val) <= 0 {
			blog.Errorf("instance organization is invalid, inst: %v, id: %s, rid: %s", rows, property.ID,
				i.GetKit().Rid)
			return nil, fmt.Errorf("instance origiazation is invalid, id: %s", property.ID)
		}

		orgItems := strings.Split(val, ",")
		org := make([]int64, len(orgItems))

		for idx, item := range orgItems {
			orgID := orgBracketsRegexp.FindStringSubmatch(item)
			if len(orgID) != 3 {
				blog.Errorf("instance organization is invalid, val: %v, id: %s, rid: %s", val, property.ID,
					i.GetKit().Rid)
				return nil, fmt.Errorf("instance organization is invalid, val: %v, id: %s", val, property.ID)
			}

			var err error
			if org[idx], err = strconv.ParseInt(orgID[1], 10, 64); err != nil {
				blog.Errorf("instance organization is invalid, val: %v, id: %s, rid: %s", val, property.ID,
					i.GetKit().Rid)
				return nil, fmt.Errorf("instance organization is invalid, val: %v, id: %s", val, property.ID)
			}
		}

		return org, nil
	}
}

func getHandleUserFieldFunc() handleInstFieldFunc {
	return func(i *Importer, property *PropWithTable, rows [][]string) (interface{}, error) {
		if len(rows) == 0 || len(rows[0]) < property.ExcelColIndex {
			blog.Errorf("instance is invalid, data: %v, rid: %s", rows, i.GetKit().Rid)
			return nil, fmt.Errorf("instance is invalid")
		}

		userNames := util.GetStrByInterface(rows[0][property.ExcelColIndex])
		userNames = userBracketsRegexp.ReplaceAllString(userNames, "")
		userNames = strings.Trim(strings.Trim(userNames, " "), ",")

		return userNames, nil
	}
}

func getHandleTableFieldFunc() handleInstFieldFunc {
	return func(i *Importer, property *PropWithTable, rows [][]string) (interface{}, error) {
		result := make([]map[string]interface{}, len(rows))
		for idx := range result {
			result[idx] = make(map[string]interface{})
		}

		for idx, row := range rows {
			for subIdx, cell := range row {
				if cell == "" {
					continue
				}

				subProp, ok := property.subProperties[subIdx]
				if !ok {
					continue
				}

				handleFunc, ok := handleInstFieldFuncMap[subProp.PropertyType]
				if !ok {
					handleFunc = getDefaultHandleFieldFunc()
				}

				val, err := handleFunc(i, &subProp, [][]string{row})
				if err != nil {
					return nil, err
				}
				result[idx][subProp.ID] = val
			}
		}

		return result, nil
	}
}

func getDefaultHandleFieldFunc() handleInstFieldFunc {
	return func(i *Importer, property *PropWithTable, rows [][]string) (interface{}, error) {
		if len(rows) == 0 || len(rows[0]) < property.ExcelColIndex {
			blog.Errorf("instance is invalid, data: %v, rid: %s", rows, i.GetKit().Rid)
			return nil, fmt.Errorf("instance is invalid")
		}

		return rows[0][property.ExcelColIndex], nil
	}
}

func getCloudAreaFieldFunc() handleInstFieldFunc {
	return func(i *Importer, property *PropWithTable, rows [][]string) (interface{}, error) {
		if len(rows) == 0 || len(rows[0]) < property.ExcelColIndex {
			blog.Errorf("cloud area data is invalid, data: %v, rid: %s", rows, i.GetKit().Rid)
			return nil, fmt.Errorf("cloud area data is invalid")
		}

		// 查找 "[" 的位置,如果找到 "["，则截取字符串到 "[" 的位置
		value := rows[0][property.ExcelColIndex]
		openBracketIndex := strings.Index(value, "[")
		if openBracketIndex != -1 {
			value = value[:openBracketIndex]
		}

		return value, nil
	}
}
