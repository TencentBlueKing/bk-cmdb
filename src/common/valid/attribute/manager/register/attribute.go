// Package register  for attribute type registration and constraint
package register

import (
	"context"
	"fmt"

	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
)

var (
	attrTypeMap = map[string]AttributeTypeI{}
)

// AttributeTypeI interface defines the methods for attribute types
type AttributeTypeI interface {
	// Name 展示名字
	Name() string
	// DisplayName  类型的唯一名字
	DisplayName() string
	// RealType  common.FieldTypeXXX , 必须是基础之一
	RealType() string

	// Info 描述信息
	Info() string

	//  预留暂未使用， 在 Validate 前处理， 做数据转换,
	// Transform(kit *rest.Kit，value interface{})(interface{}, error)

	// Validate 实际校验方法
	Validate(ctx context.Context, objID string, propertyType string, required bool, option interface{}, value interface{}) error
	// FillLostValue 填充默认值
	FillLostValue(ctx context.Context, valData mapstr.MapStr, propertyId string, defaultValue, option interface{}) error
	// ValidateOption 校验 Option
	ValidateOption(ctx context.Context, option interface{}, defaultVal interface{}) error
}

// Register attribute type
// attrTypeMap is a map of attribute name to Attribute interface
// this function is used to register a new attribute type
// it will panic if the attribute name is empty or already exists
// it is called in init() function of each attribute type file
// so that all attribute types are registered when the program starts, not at runtime
func Register(attr AttributeTypeI) {
	if attr == nil {
		blog.Errorf("register attribute is nil")
		panic("register attribute is nil")
	}

	name := attr.Name()
	if name == "" {
		blog.Errorf("register attribute name is empty")
		panic("register attribute name is empty")
	}

	if _, exists := attrTypeMap[name]; exists {
		blog.Errorf("attribute %s already exists", name)
		panic(fmt.Sprintf("attribute %s already exists", name))
	}
	attrTypeMap[name] = attr
}

// Get returns the Attribute by name
func Get(name string) (AttributeTypeI, bool) {
	a, ok := attrTypeMap[name]
	return a, ok
}
