package common

import (
	"configcenter/src/framework/core/types"
)

// CreateCondition create a condition object
func CreateCondition() Condition {
	return &condition{}
}

// Condition condition interface
type Condition interface {
	Field(fieldName string) Field
	ToMapStr() types.MapStr
}

// Condition the condition definition
type condition struct {
	fields []Field
}

// CreateField create a field
func (cli *condition) Field(fieldName string) Field {
	field := &field{
		fieldName: fieldName,
		condition: cli,
	}
	cli.fields = append(cli.fields, field)
	return field
}

// ToMapStr to MapStr object
func (cli *condition) ToMapStr() types.MapStr {
	tmpResult := types.MapStr{}
	for _, item := range cli.fields {
		tmpResult.Merge(item.ToMapStr())
	}
	return tmpResult
}
