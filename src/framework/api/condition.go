package api

import "configcenter/src/framework/common"

// CreateCondition create a condition object
func CreateCondition(tableName string) *common.Condition {
	return common.CreateCondition(tableName)
}
