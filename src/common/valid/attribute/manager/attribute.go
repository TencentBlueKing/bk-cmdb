// Package manager  provides a way to manage attributes in the system.
package manager

import (
	"configcenter/src/common/valid/attribute/manager/register"
	// import init to register all attributes
	_ "configcenter/src/common/valid/attribute/init"
)

// Get returns the Attribute by name
var Get = register.Get
