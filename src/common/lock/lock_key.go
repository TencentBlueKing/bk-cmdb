package lock

import (
	"fmt"
)

// not duplicate allow
const (
	// CreateModelFormat create model user format
	CreateModelFormat = "coreservice:create:model:%s"

	// CreateModuleAttrFormat create model  attribute format
	CreateModuleAttrFormat = "coreservice:create:model:%s:attr:%s"
)

// StrFormat  build  lock key format
type StrFormat string

// GetLockKey build lock key
func GetLockKey(format StrFormat, params ...interface{}) StrFormat {
	key := fmt.Sprintf(string(format), params...)
	return StrFormat(key)
}
