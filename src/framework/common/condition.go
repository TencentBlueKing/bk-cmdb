package common

import "configcenter/src/framework/core/types"

// CreateCondition create a condition object
func CreateCondition(tableName string) Condition {
	return &condition{tableName: tableName}
}

// Condition condition interface
type Condition interface {
	CreateField(filedName string) Field
}

// Condition the condition definition
type condition struct {
	tableName string
	fields    []Field
}

// CreateField create a field
func (cli *condition) CreateField(filedName string) Field {
	field := &field{
		fieldName: filedName,
		condition: cli,
	}
	cli.fields = append(cli.fields, field)
	return field
}

// ToMapStr to MapStr object
func (cli *condition) ToMapStr() types.MapStr {
	return nil
}
