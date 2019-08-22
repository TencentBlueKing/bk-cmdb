package querybuilder

import (
	"fmt"
	"reflect"
	"time"

	"configcenter/src/common/util"
)

var (
	TypeNumeric = "numeric"
	TypeBoolean = "boolean"
	TypeString  = "string"
	TypeUnknown = "unknown"
)

func getType(value interface{}) string {
	switch value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float64, float32:
		return TypeNumeric
	case bool:
		return TypeBoolean
	case string:
		return TypeString
	default:
		return TypeUnknown
	}
}

func validateBasicType(value interface{}) error {
	if t := getType(value); t == TypeUnknown {
		return fmt.Errorf("unknow value type with value: %+v", value)
	}
	return nil
}

func validateNumericType(value interface{}) error {
	if t := getType(value); t != TypeNumeric {
		return fmt.Errorf("unknow value type: %s, value: %+v", t, value)
	}
	return nil
}

func validateBoolType(value interface{}) error {
	if t := getType(value); t != TypeBoolean {
		return fmt.Errorf("unknow value type: %s, value: %+v", t, value)
	}
	return nil
}

func validateStringType(value interface{}) error {
	if t := getType(value); t != TypeString {
		return fmt.Errorf("unknow value type of: %s, value: %+v", t, value)
	}
	return nil
}
func validateNotEmptyStringType(value interface{}) error {
	if err := validateStringType(value); err != nil {
		return err
	}
	if len(value.(string)) == 0 {
		return fmt.Errorf("value shouldn't be empty")
	}
	return nil
}

func validateDatetimeStringType(value interface{}) error {
	if err := validateStringType(value); err != nil {
		return err
	}
	if _, err := time.Parse(time.RFC3339, value.(string)); err != nil {
		return err
	}
	return nil
}

func validateSliceOfBasicType(value interface{}, requireSameType bool) error {
	t := reflect.TypeOf(value)
	if t.Kind() != reflect.Array && t.Kind() != reflect.Slice {
		return fmt.Errorf("unexpected value type: %s, expect array", t.Kind().String())
	}
	v := reflect.ValueOf(value)
	for i := 0; i < v.Len(); i++ {
		item := v.Index(i).Interface()
		if err := validateBasicType(item); err != nil {
			return err
		}
	}
	if requireSameType == true {
		vTypes := make([]string, 0)
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i).Interface()
			vTypes = append(vTypes, getType(item))
		}
		vTypes = util.StrArrayUnique(vTypes)
		if len(vTypes) > 1 {
			return fmt.Errorf("slice element type not unique, types: %+v", vTypes)
		}
	}
	return nil
}
